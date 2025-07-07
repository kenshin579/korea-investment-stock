#!/usr/bin/env python3
"""
Phase 4.1: 배치 처리 파라미터화 테스트
Date: 2024-12-28

테스트 항목:
1. 기본 동작 (배치 없이)
2. 배치 크기 설정
3. 배치 간 대기 시간
4. 진행 상황 출력 간격
"""

import time
from unittest.mock import Mock, patch
from concurrent.futures import ThreadPoolExecutor
import threading
from ..rate_limiting import EnhancedRateLimiter


class MockKoreaInvestment:
    """테스트를 위한 Mock 클래스"""
    
    def __init__(self):
        self.rate_limiter = EnhancedRateLimiter(max_calls=10, per_seconds=1.0)
        self.executor = ThreadPoolExecutor(max_workers=3)
        self.concurrent_limit = threading.Semaphore(3)
        self.api_call_count = 0
        self.api_call_times = []
    
    def __fetch_mock_price(self, symbol: str, market: str = "KR") -> dict:
        """Mock API 호출"""
        self.rate_limiter.acquire()
        
        # API 호출 시간 기록
        call_time = time.time()
        self.api_call_times.append(call_time)
        self.api_call_count += 1
        
        # 실제 API 호출 시뮬레이션 (10-50ms)
        time.sleep(0.02)
        
        # 가끔 에러 발생 시뮬레이션
        if symbol == "ERROR":
            raise Exception("Mock API Error")
        
        return {
            'rt_cd': '0',
            'msg1': '정상처리',
            'output': {
                'symbol': symbol,
                'price': 50000 + (self.api_call_count * 100),
                'market': market
            }
        }
    
    # koreainvestmentstock.py의 __execute_concurrent_requests 메서드를 복사
    def __execute_concurrent_requests(self, method, stock_list, 
                                     batch_size=None,
                                     batch_delay=0.0,
                                     progress_interval=10):
        """배치 처리가 가능한 병렬 요청 실행"""
        from concurrent.futures import as_completed
        
        futures = {}
        results = []
        
        def wrapped_method(symbol, market):
            """세마포어로 동시 실행 제한"""
            with self.concurrent_limit:
                return method(symbol, market)
        
        # 배치 처리 설정
        if batch_size is None:
            batches = [stock_list]  # 전체를 하나의 배치로
        else:
            # stock_list를 batch_size 크기로 나누기
            batches = [stock_list[i:i + batch_size] for i in range(0, len(stock_list), batch_size)]
            print(f"📦 배치 처리 모드: {len(stock_list)}개 항목을 {len(batches)}개 배치로 처리 (배치 크기: {batch_size})")
        
        # 배치별로 처리
        for batch_idx, batch in enumerate(batches):
            if len(batches) > 1:
                print(f"\n🔄 배치 {batch_idx + 1}/{len(batches)} 처리 중... ({len(batch)}개 항목)")
            
            # 배치 내 모든 작업 제출
            batch_futures = {}
            for symbol, market in batch:
                future = self.executor.submit(wrapped_method, symbol, market)
                batch_futures[future] = (symbol, market)
                futures[future] = (symbol, market)
            
            # 배치 내 작업 완료 대기
            batch_completed = 0
            batch_total = len(batch)
            
            for future in as_completed(batch_futures, timeout=30):
                symbol, market = batch_futures[future]
                batch_completed += 1
                
                try:
                    result = future.result()
                    results.append(result)
                    
                    # 진행 상황 출력
                    if batch_completed % progress_interval == 0 or batch_completed == batch_total:
                        if len(batches) > 1:
                            print(f"  배치 진행률: {batch_completed}/{batch_total} ({batch_completed/batch_total*100:.1f}%)")
                        else:
                            total = len(stock_list)
                            completed = len(results)
                            print(f"처리 진행률: {completed}/{total} ({completed/total*100:.1f}%)")
                            
                except Exception as e:
                    print(f"❌ 에러 발생 - {symbol} ({market}): {e}")
                    results.append({
                        'rt_cd': '9',
                        'error': True,
                        'symbol': symbol,
                        'market': market,
                        'error_type': type(e).__name__
                    })
            
            # 배치 간 대기 (마지막 배치 제외)
            if batch_delay > 0 and batch_idx < len(batches) - 1:
                print(f"⏱️ 다음 배치까지 {batch_delay}초 대기...")
                time.sleep(batch_delay)
        
        # 성공/실패 요약
        success_count = sum(1 for r in results if not r.get('error', False))
        error_count = len(results) - success_count
        
        print(f"\n📊 처리 완료 - 성공: {success_count}, 실패: {error_count}")
        
        # 배치 처리 통계
        if len(batches) > 1:
            print(f"   배치 수: {len(batches)}, 배치 크기: {batch_size}")
        
        return results
    
    def fetch_price_list(self, stock_list):
        """가격 조회 (배치 처리 지원)"""
        return self.__execute_concurrent_requests(self.__fetch_mock_price, stock_list)
    
    def fetch_price_list_batched(self, stock_list, batch_size=10, batch_delay=1.0):
        """가격 조회 (배치 처리 파라미터 지정)"""
        return self.__execute_concurrent_requests(
            self.__fetch_mock_price, 
            stock_list,
            batch_size=batch_size,
            batch_delay=batch_delay,
            progress_interval=5
        )
    
    def shutdown(self):
        """리소스 정리"""
        self.executor.shutdown(wait=True)
        if hasattr(self.rate_limiter, 'print_stats'):
            self.rate_limiter.print_stats()


def test_default_behavior():
    """테스트 1: 기본 동작 (배치 없이)"""
    print("=== 1. 기본 동작 테스트 (배치 없이) ===")
    
    kis = MockKoreaInvestment()
    stock_list = [(f"00593{i}", "KR") for i in range(15)]
    
    start_time = time.time()
    results = kis.fetch_price_list(stock_list)
    elapsed = time.time() - start_time
    
    print(f"\n실행 시간: {elapsed:.2f}초")
    print(f"API 호출 수: {kis.api_call_count}")
    print(f"평균 TPS: {kis.api_call_count / elapsed:.2f}")
    
    kis.shutdown()
    print("\n✅ 기본 동작 테스트 완료\n")


def test_batch_size():
    """테스트 2: 배치 크기 설정"""
    print("=== 2. 배치 크기 설정 테스트 ===")
    
    kis = MockKoreaInvestment()
    stock_list = [(f"00593{i}", "KR") for i in range(25)]
    
    start_time = time.time()
    results = kis.fetch_price_list_batched(stock_list, batch_size=10, batch_delay=0.5)
    elapsed = time.time() - start_time
    
    print(f"\n실행 시간: {elapsed:.2f}초")
    print(f"API 호출 수: {kis.api_call_count}")
    
    # 배치별 시간 분석
    if len(kis.api_call_times) > 20:
        batch1_end = kis.api_call_times[9]
        batch2_start = kis.api_call_times[10]
        batch2_end = kis.api_call_times[19]
        batch3_start = kis.api_call_times[20]
        
        print(f"\n배치 간 대기 시간:")
        print(f"  배치 1→2: {batch2_start - batch1_end:.2f}초")
        print(f"  배치 2→3: {batch3_start - batch2_end:.2f}초")
    
    kis.shutdown()
    print("\n✅ 배치 크기 설정 테스트 완료\n")


def test_progress_interval():
    """테스트 3: 진행 상황 출력 간격"""
    print("=== 3. 진행 상황 출력 간격 테스트 ===")
    
    kis = MockKoreaInvestment()
    stock_list = [(f"00593{i}", "KR") for i in range(12)]
    
    # progress_interval=3으로 설정
    results = kis._MockKoreaInvestment__execute_concurrent_requests(
        kis._MockKoreaInvestment__fetch_mock_price,
        stock_list,
        progress_interval=3
    )
    
    print(f"\n처리 항목 수: {len(results)}")
    print("✅ 진행 상황 출력 간격 테스트 완료\n")
    
    kis.shutdown()


def test_error_handling():
    """테스트 4: 배치 처리 중 에러 처리"""
    print("=== 4. 배치 처리 중 에러 처리 ===")
    
    kis = MockKoreaInvestment()
    stock_list = [
        ("005930", "KR"),
        ("ERROR", "KR"),  # 에러 발생
        ("000660", "KR"),
        ("035720", "KR"),
        ("ERROR", "KR"),  # 에러 발생
    ]
    
    results = kis.fetch_price_list_batched(stock_list, batch_size=3, batch_delay=0.2)
    
    error_count = sum(1 for r in results if r.get('error', False))
    success_count = len(results) - error_count
    
    print(f"\n결과 검증:")
    print(f"  전체 결과: {len(results)}")
    print(f"  성공: {success_count}")
    print(f"  실패: {error_count}")
    
    if error_count == 2:
        print("✅ 에러 처리 테스트 성공")
    else:
        print("❌ 에러 처리 테스트 실패")
    
    kis.shutdown()
    print()


def test_performance_comparison():
    """테스트 5: 배치 처리 성능 비교"""
    print("=== 5. 배치 처리 성능 비교 ===")
    
    stock_list = [(f"00593{i:02d}", "KR") for i in range(30)]
    
    # 배치 없이 처리
    print("\n1) 배치 없이 처리:")
    kis1 = MockKoreaInvestment()
    start = time.time()
    results1 = kis1.fetch_price_list(stock_list)
    time1 = time.time() - start
    print(f"   실행 시간: {time1:.2f}초")
    kis1.shutdown()
    
    # 배치 크기 10, 대기 시간 0.5초
    print("\n2) 배치 크기 10, 대기 0.5초:")
    kis2 = MockKoreaInvestment()
    start = time.time()
    results2 = kis2.fetch_price_list_batched(stock_list, batch_size=10, batch_delay=0.5)
    time2 = time.time() - start
    print(f"   실행 시간: {time2:.2f}초")
    kis2.shutdown()
    
    # 배치 크기 5, 대기 시간 0.2초
    print("\n3) 배치 크기 5, 대기 0.2초:")
    kis3 = MockKoreaInvestment()
    start = time.time()
    results3 = kis3.fetch_price_list_batched(stock_list, batch_size=5, batch_delay=0.2)
    time3 = time.time() - start
    print(f"   실행 시간: {time3:.2f}초")
    kis3.shutdown()
    
    print("\n성능 비교 요약:")
    print(f"  배치 없음: {time1:.2f}초")
    print(f"  배치 10: {time2:.2f}초 ({(time2-time1)/time1*100:+.1f}%)")
    print(f"  배치 5: {time3:.2f}초 ({(time3-time1)/time1*100:+.1f}%)")
    print("\n✅ 성능 비교 완료\n")


if __name__ == "__main__":
    print("배치 처리 파라미터화 테스트\n")
    
    test_default_behavior()
    test_batch_size()
    test_progress_interval()
    test_error_handling()
    test_performance_comparison()
    
    print("모든 테스트 완료!") 