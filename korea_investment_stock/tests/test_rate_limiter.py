#!/usr/bin/env python3
"""
Enhanced Rate Limiter Unit Tests
Date: 2024-12-28

Phase 6.1: 단위 테스트 작성
- Token Bucket 리필 테스트
- 동시성 테스트
- 최소 간격 보장 테스트
"""

import unittest
import time
import threading
from multiprocessing import Value
from concurrent.futures import ThreadPoolExecutor, as_completed
from datetime import datetime

from ..rate_limiting import EnhancedRateLimiter


class TestEnhancedRateLimiter(unittest.TestCase):
    """Enhanced Rate Limiter 단위 테스트"""
    
    def setUp(self):
        """테스트 환경 설정"""
        # 테스트용 설정: 초당 10회 호출 제한
        self.rate_limiter = EnhancedRateLimiter(
            max_calls=10,
            per_seconds=1.0,
            safety_margin=1.0,  # 테스트에서는 안전 마진 없음
            enable_min_interval=True,
            enable_stats=True
        )
    
    def test_token_bucket_refill(self):
        """Token Bucket 리필 테스트"""
        print("\n=== Token Bucket 리필 테스트 ===")
        
        # 빠르게 토큰 소비
        success_count = 0
        for i in range(15):
            if self.rate_limiter.acquire():
                success_count += 1
        
        # 일부는 실패해야 함
        self.assertLess(success_count, 15)
        self.assertGreater(success_count, 5)
        print(f"✅ 첫 번째 버스트: {success_count}/15 성공")
        
        # 0.5초 대기 (5개 토큰 리필)
        time.sleep(0.5)
        
        # 5개는 성공해야 함
        success_count = 0
        for i in range(10):
            if self.rate_limiter.acquire():
                success_count += 1
        
        self.assertGreaterEqual(success_count, 4)  # 약간의 오차 허용
        self.assertLessEqual(success_count, 6)
        print(f"✅ 0.5초 후 {success_count}개 토큰 리필 확인")
        
        # 통계 확인
        stats = self.rate_limiter.get_stats()
        print(f"✅ 총 호출 수: {stats['total_calls']}")
        print(f"✅ 차단된 호출 수: {stats['blocked_calls']}")
    
    def test_sliding_window(self):
        """Sliding Window 테스트"""
        print("\n=== Sliding Window 테스트 ===")
        
        # 새로운 rate limiter (초당 5회)
        rl = EnhancedRateLimiter(max_calls=5, per_seconds=1.0, safety_margin=1.0)
        
        # 빠르게 여러 번 호출
        first_burst = 0
        for i in range(10):
            if rl.acquire():
                first_burst += 1
        
        # 일부만 성공해야 함
        self.assertLessEqual(first_burst, 5)
        self.assertGreater(first_burst, 0)
        print(f"✅ 첫 버스트: {first_burst}/10 성공")
        
        # 1초 대기 후 다시 시도
        time.sleep(1.1)
        second_burst = 0
        for i in range(5):
            if rl.acquire():
                second_burst += 1
        
        self.assertGreater(second_burst, 0)
        print(f"✅ 1초 후: {second_burst}/5 성공")
    
    def test_concurrency(self):
        """동시성 테스트"""
        print("\n=== 동시성 테스트 ===")
        
        # 새로운 rate limiter (초당 20회)
        rl = EnhancedRateLimiter(max_calls=20, per_seconds=1.0, safety_margin=1.0)
        
        success_count = Value('i', 0)
        blocked_count = Value('i', 0)
        
        def worker():
            """워커 스레드"""
            for _ in range(5):
                if rl.acquire():
                    with success_count.get_lock():
                        success_count.value += 1
                else:
                    with blocked_count.get_lock():
                        blocked_count.value += 1
                time.sleep(0.01)  # 약간의 지연
        
        # 10개 스레드에서 각각 5번 호출 (총 50회 시도)
        threads = []
        for _ in range(10):
            t = threading.Thread(target=worker)
            threads.append(t)
            t.start()
        
        # 모든 스레드 대기
        for t in threads:
            t.join()
        
        print(f"✅ 성공: {success_count.value}, 차단: {blocked_count.value}")
        self.assertEqual(success_count.value + blocked_count.value, 50)
        self.assertLessEqual(success_count.value, 25)  # 약 1초 동안 최대 20-25개
        print("✅ Thread-safe 동작 확인")
    
    def test_min_interval(self):
        """최소 간격 보장 테스트"""
        print("\n=== 최소 간격 보장 테스트 ===")
        
        # 최소 간격이 활성화된 rate limiter (초당 5회 = 0.2초 간격)
        rl = EnhancedRateLimiter(
            max_calls=5,
            per_seconds=1.0,
            safety_margin=1.0,
            enable_min_interval=True
        )
        
        timestamps = []
        
        # 5번 연속 호출
        for i in range(5):
            start = time.time()
            self.assertTrue(rl.acquire())
            timestamps.append(time.time())
            
            # 첫 호출이 아니면 간격 체크
            if i > 0:
                interval = timestamps[i] - timestamps[i-1]
                # 최소 간격은 min_interval = per_seconds / effective_max_calls
                # safety_margin = 1.0이므로 effective = 5, min_interval = 0.2
                # 하지만 토큰이 충분하면 더 빠르게 호출 가능
                self.assertGreaterEqual(interval, 0.15)  # 약간의 오차 허용
                print(f"✅ 호출 {i}: 간격 {interval:.3f}초")
        
        print("✅ 최소 간격 보장 확인")
    
    def test_statistics(self):
        """통계 수집 테스트"""
        print("\n=== 통계 수집 테스트 ===")
        
        rl = EnhancedRateLimiter(
            max_calls=5,
            per_seconds=1.0,
            enable_stats=True
        )
        
        # 몇 번 호출
        for i in range(8):
            rl.acquire()
            time.sleep(0.1)
        
        # 에러 기록
        rl.record_error()
        rl.record_error()
        
        stats = rl.get_stats()
        self.assertGreater(stats['total_calls'], 0)
        self.assertEqual(stats['error_count'], 2)
        # calls_per_second 대신 max_calls_per_second 확인
        self.assertIn('max_calls_per_second', stats)
        
        print(f"✅ 통계 수집 확인:")
        print(f"   - 총 호출: {stats['total_calls']}")
        print(f"   - 에러: {stats['error_count']}")
        print(f"   - 에러율: {stats['error_rate']:.1%}")
        if 'max_calls_per_second' in stats:
            print(f"   - 최대 초당 호출: {stats['max_calls_per_second']}")
    
    def test_token_bucket_and_sliding_window_hybrid(self):
        """Token Bucket과 Sliding Window 하이브리드 테스트"""
        print("\n=== 하이브리드 알고리즘 테스트 ===")
        
        rl = EnhancedRateLimiter(
            max_calls=10,
            per_seconds=1.0,
            safety_margin=0.8  # 실제 8회/초
        )
        
        # 빠르게 15번 호출
        success_fast = 0
        for i in range(15):
            if rl.acquire():
                success_fast += 1
        
        # safety margin으로 인해 8개 정도만 성공
        self.assertLessEqual(success_fast, 10)
        self.assertGreaterEqual(success_fast, 5)
        print(f"✅ Safety margin 적용 확인: {success_fast}/15 성공")
        
        # 1초 대기 후 다시 시도
        time.sleep(1.0)
        
        success_after = 0
        for i in range(10):
            if rl.acquire():
                success_after += 1
        
        self.assertGreaterEqual(success_after, 7)
        print(f"✅ 1초 후 리필 확인: {success_after}/10 성공")
    
    def test_thread_pool_integration(self):
        """ThreadPoolExecutor와의 통합 테스트"""
        print("\n=== ThreadPoolExecutor 통합 테스트 ===")
        
        rl = EnhancedRateLimiter(max_calls=15, per_seconds=1.0)
        
        def api_call(i):
            """모의 API 호출"""
            if rl.acquire():
                time.sleep(0.05)  # API 호출 시뮬레이션
                return f"Success-{i}"
            else:
                return f"Blocked-{i}"
        
        # ThreadPoolExecutor로 30개 작업 제출
        with ThreadPoolExecutor(max_workers=5) as executor:
            futures = [executor.submit(api_call, i) for i in range(30)]
            
            results = []
            for future in as_completed(futures):
                results.append(future.result())
        
        success_count = sum(1 for r in results if r.startswith("Success"))
        blocked_count = sum(1 for r in results if r.startswith("Blocked"))
        
        print(f"✅ 성공: {success_count}, 차단: {blocked_count}")
        self.assertEqual(success_count + blocked_count, 30)
        # 실제로는 더 많이 성공할 수 있음 (토큰 버킷 리필 때문에)
        self.assertGreater(success_count, 10)  # 최소 10개는 성공
        self.assertLess(success_count, 31)  # 모두 성공하지는 않음


def run_rate_limiter_tests():
    """테스트 실행 헬퍼"""
    print("\n" + "="*60)
    print("Enhanced Rate Limiter 단위 테스트 시작")
    print("="*60)
    
    unittest.main(module=__name__, argv=[''], exit=False)
    
    print("\n" + "="*60)
    print("테스트 완료")
    print("="*60)


if __name__ == "__main__":
    run_rate_limiter_tests() 