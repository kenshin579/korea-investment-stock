# ThreadPoolExecutor 패턴 개선안

## 현재 패턴의 문제점

### 1. 리소스 관리 문제
```python
# 현재: 수동 shutdown 필요
broker = KoreaInvestment(...)
# 사용
broker.shutdown()  # 잊기 쉬움
```

### 2. 동시성 제어 분리
- ThreadPoolExecutor와 RateLimiter가 독립적으로 동작
- 배치 내 5개 스레드가 동시에 RateLimiter에 접근하여 경쟁

### 3. 에러 처리 미흡
```python
# 현재: 개별 future에서 발생한 예외 처리 없음
batch_results = [future.result() for future in futures]
```

### 4. 순차적 결과 수집
- `future.result()` 호출이 순차적으로 블로킹
- 첫 번째 요청이 느리면 전체 지연

## 개선안

### 개선안 1: 컨텍스트 매니저 패턴

```python
from contextlib import contextmanager
import atexit

class KoreaInvestment:
    def __init__(self, ...):
        # ... 기존 코드 ...
        self.executor = ThreadPoolExecutor(max_workers=4)
        # 프로그램 종료 시 자동 정리
        atexit.register(self.shutdown)
    
    def __enter__(self):
        return self
    
    def __exit__(self, exc_type, exc_val, exc_tb):
        self.shutdown()
        return False
    
    def shutdown(self):
        if hasattr(self, 'executor') and self.executor:
            self.executor.shutdown(wait=True)
            self.executor = None

# 사용법
with KoreaInvestment(...) as broker:
    results = broker.fetch_price_list(stocks)
    # 자동으로 리소스 정리됨
```

### 개선안 2: 세마포어 기반 동시성 제어

```python
import threading
from concurrent.futures import as_completed

class KoreaInvestment:
    def __init__(self, ...):
        # 동시 실행 제한을 위한 세마포어
        self.concurrent_limit = threading.Semaphore(3)  # 동시 3개만
        self.executor = ThreadPoolExecutor(max_workers=8)
    
    def __execute_concurrent_requests(self, method, stock_list):
        """개선된 병렬 요청 실행"""
        futures = {}
        results = []
        
        def wrapped_method(symbol, market):
            # 세마포어로 동시 실행 제한
            with self.concurrent_limit:
                return method(symbol, market)
        
        # 모든 작업 제출
        for symbol, market in stock_list:
            future = self.executor.submit(wrapped_method, symbol, market)
            futures[future] = (symbol, market)
        
        # 완료되는 대로 결과 수집
        for future in as_completed(futures):
            symbol, market = futures[future]
            try:
                result = future.result()
                results.append(result)
            except Exception as e:
                print(f"Error processing {symbol}: {e}")
                # 에러 정보도 결과에 포함
                results.append({
                    'error': str(e),
                    'symbol': symbol,
                    'market': market
                })
        
        return results
```

### 개선안 3: Queue 기반 Producer-Consumer 패턴

```python
import queue
from threading import Thread

class KoreaInvestment:
    def __init__(self, ...):
        self.request_queue = queue.Queue()
        self.result_queue = queue.Queue()
        self.workers = []
        self.start_workers(num_workers=3)
    
    def start_workers(self, num_workers):
        for i in range(num_workers):
            worker = Thread(target=self._worker, daemon=True)
            worker.start()
            self.workers.append(worker)
    
    def _worker(self):
        """워커 스레드 - Rate Limiting 포함"""
        while True:
            try:
                method, symbol, market, future = self.request_queue.get(timeout=1)
                if method is None:  # 종료 신호
                    break
                
                # Rate limiting
                self.rate_limiter.acquire()
                
                try:
                    result = method(symbol, market)
                    future.set_result(result)
                except Exception as e:
                    future.set_exception(e)
                
                self.request_queue.task_done()
            except queue.Empty:
                continue
    
    def __execute_concurrent_requests(self, method, stock_list):
        """Queue 기반 병렬 처리"""
        from concurrent.futures import Future
        
        futures = []
        for symbol, market in stock_list:
            future = Future()
            self.request_queue.put((method, symbol, market, future))
            futures.append(future)
        
        # 결과 수집
        results = []
        for future in futures:
            try:
                result = future.result(timeout=30)  # 타임아웃 설정
                results.append(result)
            except Exception as e:
                results.append({'error': str(e)})
        
        return results
```

### 개선안 4: asyncio 기반 비동기 처리

```python
import asyncio
import aiohttp
from asyncio import Semaphore

class AsyncKoreaInvestment:
    def __init__(self, ...):
        self.semaphore = Semaphore(3)  # 동시 요청 제한
        self.rate_limiter = AsyncRateLimiter(12, 1)
    
    async def fetch_price_list(self, stock_list):
        """비동기 일괄 가격 조회"""
        async with aiohttp.ClientSession() as session:
            tasks = [
                self._fetch_price_async(session, symbol, market)
                for symbol, market in stock_list
            ]
            return await asyncio.gather(*tasks, return_exceptions=True)
    
    async def _fetch_price_async(self, session, symbol, market):
        """비동기 개별 가격 조회"""
        async with self.semaphore:  # 동시 실행 제한
            await self.rate_limiter.acquire()  # Rate limiting
            
            url = f"{self.base_url}/..."
            headers = {...}
            
            try:
                async with session.get(url, headers=headers) as response:
                    return await response.json()
            except Exception as e:
                return {'error': str(e), 'symbol': symbol}
```

### 개선안 5: 통합 Rate-Limited Executor

```python
from concurrent.futures import ThreadPoolExecutor, Future
import heapq
import time

class RateLimitedExecutor:
    """Rate Limiting이 통합된 Executor"""
    
    def __init__(self, max_workers=4, rate_limit=12, per_seconds=1):
        self.executor = ThreadPoolExecutor(max_workers=max_workers)
        self.rate_limit = rate_limit
        self.per_seconds = per_seconds
        self.pending_tasks = []  # 우선순위 큐
        self.call_times = deque()
        self.lock = threading.Lock()
        self.scheduler = Thread(target=self._schedule_tasks, daemon=True)
        self.scheduler.start()
    
    def submit(self, fn, *args, priority=0, **kwargs):
        """작업 제출 (우선순위 지원)"""
        future = Future()
        task = (priority, time.time(), fn, args, kwargs, future)
        
        with self.lock:
            heapq.heappush(self.pending_tasks, task)
        
        return future
    
    def _schedule_tasks(self):
        """Rate limit을 고려하여 작업 스케줄링"""
        while True:
            with self.lock:
                # Rate limit 확인
                now = time.time()
                self.call_times = deque(
                    t for t in self.call_times 
                    if t > now - self.per_seconds
                )
                
                if len(self.call_times) >= self.rate_limit:
                    time.sleep(0.01)
                    continue
                
                if not self.pending_tasks:
                    time.sleep(0.01)
                    continue
                
                # 다음 작업 실행
                _, _, fn, args, kwargs, future = heapq.heappop(self.pending_tasks)
                self.call_times.append(now)
                
                # 실제 실행
                self.executor.submit(self._execute_task, fn, args, kwargs, future)
    
    def _execute_task(self, fn, args, kwargs, future):
        """작업 실행 및 결과 설정"""
        try:
            result = fn(*args, **kwargs)
            future.set_result(result)
        except Exception as e:
            future.set_exception(e)
```

## 권장 사항

### 1. 단기 개선 (현재 구조 유지)
- **개선안 1**: 컨텍스트 매니저 추가
- **개선안 2**: 세마포어로 동시성 제어
- 에러 처리 강화
- `as_completed` 사용으로 효율성 개선

### 2. 중기 개선
- **개선안 3**: Queue 기반 패턴으로 전환
- **개선안 5**: 통합 Rate-Limited Executor 도입
- 재시도 로직 통합

### 3. 장기 개선
- **개선안 4**: asyncio 기반으로 전면 재작성
- 더 효율적이고 확장 가능한 구조
- 현대적인 Python 비동기 패턴 활용

## 성능 비교

| 패턴 | 장점 | 단점 | 복잡도 |
|------|-----|------|--------|
| 현재 (ThreadPoolExecutor) | 간단함 | 동시성 제어 미흡 | 낮음 |
| 세마포어 추가 | 동시성 제어 개선 | 여전히 분리된 구조 | 중간 |
| Queue 기반 | 정교한 제어 | 구현 복잡도 증가 | 높음 |
| asyncio | 최고 효율성 | 전면 재작성 필요 | 매우 높음 |
| 통합 Executor | Rate Limiting 통합 | 새로운 추상화 | 높음 |

## 결론

현재 패턴은 작동하지만 개선의 여지가 많습니다. 단기적으로는 컨텍스트 매니저와 세마포어를 추가하여 안정성을 높이고, 중장기적으로는 더 정교한 동시성 제어 메커니즘으로 전환하는 것을 권장합니다. 