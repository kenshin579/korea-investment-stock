# Rate Limit 구현 문서

## 개요
한국투자증권 API 라이브러리는 API 호출 제한을 준수하기 위해 `RateLimiter` 클래스와 `ThreadPoolExecutor`를 사용한 정교한 Rate Limiting 시스템을 구현하고 있습니다.

## 핵심 구성 요소

### 1. RateLimiter 클래스
```python
class RateLimiter:
    def __init__(self, max_calls, per_seconds):
        self.max_calls = max_calls         # 시간 윈도우 내 최대 호출 횟수
        self.per_seconds = per_seconds     # 시간 윈도우 (초)
        self.lock = threading.Lock()       # 스레드 안전성을 위한 락
        self.call_timestamps = deque()     # 호출 타임스탬프 큐
        self.calls_per_second = defaultdict(int)  # 호출 통계 추적
```

### 2. ThreadPoolExecutor
- `max_workers=20`: 최대 20개의 동시 실행 스레드
- 병렬 API 호출을 위해 사용됨

## 작동 방식

### 1. 초기화 (KoreaInvestment.__init__)
```python
max_calls = 20
self.rate_limiter = RateLimiter(max_calls, 1)  # 초당 최대 20회 호출
self.executor = ThreadPoolExecutor(max_workers=max_calls)
```

### 2. Rate Limiting 메커니즘 (RateLimiter.acquire)
1. **타임스탬프 관리**:
   - 현재 시간 기준으로 만료된 타임스탬프 제거
   - `deque`에서 1초 이상 지난 타임스탬프 제거

2. **호출 제한 확인**:
   - 현재 큐에 있는 타임스탬프 개수가 `max_calls` 이상인지 확인
   - 제한에 도달한 경우, 가장 오래된 타임스탬프가 만료될 때까지 대기

3. **호출 기록**:
   - 현재 타임스탬프를 큐에 추가
   - 통계 정보 업데이트

### 3. 병렬 처리 (__execute_concurrent_requests)
```python
def __execute_concurrent_requests(self, method, stock_list):
    # 각 종목에 대해 병렬로 API 호출 실행
    futures = [self.executor.submit(method, symbol_id, market) 
               for symbol_id, market in stock_list]
    # 모든 결과 대기
    results = [future.result() for future in futures]
    # 통계 출력
    self.rate_limiter.print_stats()
    return results
```

### 4. 개별 API 호출
개별 API 메서드들은 실제 API 호출 전에 `rate_limiter.acquire()`를 호출:
- `__fetch_price_detail_oversea`
- `__fetch_stock_info`
- `__fetch_search_stock_info`

## 사용되는 메서드들

### 병렬 처리 메서드 (Public)
- `fetch_search_stock_info_list`: 여러 종목의 정보를 병렬로 검색
- `fetch_price_list`: 여러 종목의 가격을 병렬로 조회
- `fetch_price_detail_oversea_list`: 해외 주식 상세 가격 병렬 조회
- `fetch_stock_info_list`: 주식 기본 정보 병렬 조회

### Rate Limiting이 적용되는 Private 메서드
- `__fetch_price_detail_oversea`: 해외주식 현재가 상세 조회
- `__fetch_stock_info`: 주식 기본 정보 조회
- `__fetch_search_stock_info`: 국내 주식 정보 검색

## 통계 및 모니터링

### get_stats()
호출 통계 정보를 딕셔너리로 반환:
- `calls_per_second`: 초별 호출 횟수
- `max_calls_in_one_second`: 최대 초당 호출 횟수
- `total_calls`: 총 호출 횟수
- `seconds_tracked`: 추적된 시간(초)

### print_stats()
호출 통계를 콘솔에 출력:
```
===== 초당 API 호출 횟수 분석 =====
시간: 14:30:15, 호출 수: 15
시간: 14:30:16, 호출 수: 20
시간: 14:30:17, 호출 수: 18

최대 초당 호출 횟수: 20
설정된 max_calls: 20
제한 준수 여부: 준수
총 호출 횟수: 53
================================
```

## 장점

1. **스레드 안전성**: `threading.Lock`을 사용하여 동시성 문제 방지
2. **정확한 제한 관리**: 슬라이딩 윈도우 방식으로 정확한 초당 호출 횟수 제한
3. **병렬 처리 효율성**: ThreadPoolExecutor를 통한 효율적인 병렬 API 호출
4. **실시간 모니터링**: 호출 통계를 통한 Rate Limit 준수 확인
5. **자동 대기**: 제한 초과 시 자동으로 대기하여 API 제한 위반 방지

## 주의사항

1. **max_calls 설정**: 현재 20으로 설정되어 있으나, 실제 API 제한에 맞게 조정 필요
2. **병렬 처리 시**: 모든 병렬 요청이 동시에 rate limiter를 통과하려 할 수 있으므로 주의
3. **통계 메모리**: `calls_per_second` 딕셔너리가 계속 증가할 수 있으므로 장시간 실행 시 주의

## 예시 사용법

```python
# 여러 종목 정보를 병렬로 조회
stock_list = [("005930", "KR"), ("000660", "KR"), ("035720", "KR")]
results = broker.fetch_stock_info_list(stock_list)
# 자동으로 rate limiting이 적용되고 통계가 출력됨
```

## 개선된 구현 (v2)

### 발견된 문제점들

1. **동시성 문제**: 20개 스레드가 동시에 rate limiter를 통과하려 함
2. **슬라이딩 윈도우 vs 고정 윈도우**: 서버와 클라이언트의 시간 계산 방식 차이
3. **초기 버스트**: 처음 시작 시 과도한 요청 발생
4. **타이밍 경계 문제**: 초 단위 경계에서 발생하는 미묘한 차이

### 개선 사항

#### 1. 하이브리드 Rate Limiting (Token Bucket + Sliding Window)
```python
class RateLimiter:
    def __init__(self, max_calls, per_seconds, safety_margin=0.9):
        self.max_calls = int(max_calls * safety_margin)  # 안전 마진 적용
        self.tokens = self.max_calls  # 토큰 버킷
        self.refill_rate = self.max_calls / self.per_seconds
        # ... 슬라이딩 윈도우도 병행
```

#### 2. 보수적인 설정값
```python
# 기존: max_calls = 20
max_calls = 15  # 20에서 15로 감소
self.rate_limiter = RateLimiter(max_calls, 1, safety_margin=0.8)  # 실제로는 12회/초
self.executor = ThreadPoolExecutor(max_workers=min(max_calls // 2, 8))  # 최대 8개 워커
```

#### 3. 배치 처리
```python
def __execute_concurrent_requests(self, method, stock_list):
    batch_size = 5  # 한 번에 처리할 최대 요청 수
    # 배치 단위로 처리하고 배치 간 대기 시간 추가
```

#### 4. 최소 간격 보장
```python
# 각 호출 후 최소 간격 보장
min_interval = self.per_seconds / (self.max_calls * 1.2)  # 20% 여유
time.sleep(min_interval)
```

### 개선 효과

1. **안전 마진**: 실제 제한의 80%만 사용하여 여유 확보
2. **토큰 버킷**: 더 정확한 rate limiting
3. **배치 처리**: 동시 요청 수 제한으로 버스트 방지
4. **최소 간격**: 각 요청 사이 최소 시간 보장

### 추가 권장사항

1. **모니터링**: 실시간으로 API 응답 헤더의 rate limit 정보 확인
2. **백오프 전략**: 429 에러 시 지수 백오프 적용
3. **환경별 설정**: 개발/운영 환경별로 다른 제한값 사용 