#!/usr/bin/env python3
"""
Error Handling Unit Tests
Date: 2024-12-28

Phase 6.1: 단위 테스트 작성
- Exponential Backoff 테스트
- 재시도 로직 테스트
"""

import unittest
import time
from unittest.mock import Mock, patch

from ..rate_limiting.enhanced_backoff_strategy import (
    EnhancedBackoffStrategy,
    BackoffConfig,
    CircuitState
)
from ..rate_limiting.enhanced_retry_decorator import (
    retry_on_rate_limit,
    enhanced_retry,
    retry_on_network_error,
    RateLimitError,
    APIError
)


class TestExponentialBackoff(unittest.TestCase):
    """Exponential Backoff 단위 테스트"""
    
    def setUp(self):
        """테스트 환경 설정"""
        # 테스트용 설정
        self.config = BackoffConfig(
            base_delay=0.1,  # 빠른 테스트를 위해 짧게
            max_delay=1.0,
            exponential_base=2.0,
            jitter_factor=0.1,
            failure_threshold=3,
            success_threshold=2,
            timeout=1.0
        )
        self.backoff = EnhancedBackoffStrategy(self.config)
    
    def test_exponential_calculation(self):
        """지수 백오프 계산 테스트"""
        print("\n=== Exponential Backoff 계산 테스트 ===")
        
        # 첫 번째 재시도
        delay1, reason1 = self.backoff.calculate_backoff(0)
        self.assertGreaterEqual(delay1, 0.1)  # base_delay
        self.assertLessEqual(delay1, 0.11)  # base_delay + jitter
        print(f"✅ 첫 번째 재시도: {delay1:.3f}초")
        
        # 두 번째 재시도
        delay2, reason2 = self.backoff.calculate_backoff(1)
        self.assertGreaterEqual(delay2, 0.2)  # base_delay * 2
        self.assertLessEqual(delay2, 0.22)  # + jitter
        print(f"✅ 두 번째 재시도: {delay2:.3f}초")
        
        # 세 번째 재시도
        delay3, reason3 = self.backoff.calculate_backoff(2)
        self.assertGreaterEqual(delay3, 0.4)  # base_delay * 4
        self.assertLessEqual(delay3, 0.44)  # + jitter
        print(f"✅ 세 번째 재시도: {delay3:.3f}초")
        
        # 최대값 테스트
        delay_max, _ = self.backoff.calculate_backoff(10)
        self.assertLessEqual(delay_max, 1.1)  # max_delay + jitter
        print(f"✅ 최대 대기 시간 제한: {delay_max:.3f}초")
    
    def test_circuit_breaker_states(self):
        """Circuit Breaker 상태 전환 테스트"""
        print("\n=== Circuit Breaker 상태 전환 테스트 ===")
        
        # 초기 상태는 CLOSED
        self.assertEqual(self.backoff.state, CircuitState.CLOSED)
        print("✅ 초기 상태: CLOSED")
        
        # 3번 실패 -> OPEN
        for i in range(3):
            self.backoff.record_attempt(success=False)
        
        self.assertEqual(self.backoff.state, CircuitState.OPEN)
        print("✅ 3번 실패 후: OPEN")
        
        # OPEN 상태에서 백오프 계산
        delay, reason = self.backoff.calculate_backoff(0)
        self.assertIn("Circuit OPEN", reason)
        print(f"✅ OPEN 상태 메시지: {reason}")
        
        # 타임아웃 후 HALF_OPEN
        time.sleep(1.1)  # timeout=1.0초
        self.backoff._check_circuit()
        self.assertEqual(self.backoff.state, CircuitState.HALF_OPEN)
        print("✅ 타임아웃 후: HALF_OPEN")
        
        # 2번 성공 -> CLOSED
        self.backoff.record_attempt(success=True)
        self.backoff.record_attempt(success=True)
        self.assertEqual(self.backoff.state, CircuitState.CLOSED)
        print("✅ 2번 성공 후: CLOSED")
    
    def test_jitter(self):
        """Jitter 테스트"""
        print("\n=== Jitter 테스트 ===")
        
        # 같은 재시도 횟수에서 10번 호출
        delays = []
        for _ in range(10):
            delay, _ = self.backoff.calculate_backoff(1)
            delays.append(delay)
        
        # 모든 값이 다름 (Jitter 때문에)
        unique_delays = set(delays)
        self.assertGreater(len(unique_delays), 5)  # 최소 5개 이상 다른 값
        
        # 범위 확인
        min_delay = min(delays)
        max_delay = max(delays)
        self.assertGreaterEqual(min_delay, 0.2)  # base * 2
        self.assertLessEqual(max_delay, 0.22)  # base * 2 * 1.1
        
        print(f"✅ Jitter 범위: {min_delay:.3f} ~ {max_delay:.3f}")
        print(f"✅ 고유한 값 개수: {len(unique_delays)}/10")
    
    def test_adaptive_backoff(self):
        """Adaptive Backoff 테스트"""
        print("\n=== Adaptive Backoff 테스트 ===")
        
        # 여러 번 실패 기록
        for _ in range(10):
            self.backoff.record_attempt(success=False)
        
        # 낮은 성공률일 때 더 긴 대기
        delay_low_success, _ = self.backoff.calculate_backoff(1)
        
        # 성공률 개선
        self.backoff.reset()
        for _ in range(8):
            self.backoff.record_attempt(success=True)
        for _ in range(2):
            self.backoff.record_attempt(success=False)
        
        # 높은 성공률일 때 더 짧은 대기
        delay_high_success, _ = self.backoff.calculate_backoff(1)
        
        # Adaptive가 활성화되어 있으면 차이가 있어야 함
        if self.config.enable_adaptive:
            self.assertGreater(delay_low_success, delay_high_success)
            print(f"✅ 낮은 성공률 대기: {delay_low_success:.3f}초")
            print(f"✅ 높은 성공률 대기: {delay_high_success:.3f}초")
    
    def test_statistics(self):
        """통계 수집 테스트"""
        print("\n=== 통계 수집 테스트 ===")
        
        # 몇 번의 성공과 실패
        for _ in range(7):
            self.backoff.record_attempt(success=True)
        for _ in range(3):
            self.backoff.record_attempt(success=False)
        
        stats = self.backoff.get_stats()
        
        self.assertEqual(stats['total_attempts'], 10)
        self.assertEqual(stats['total_failures'], 3)
        self.assertEqual(stats['success_rate'], 0.7)
        
        print(f"✅ 총 시도: {stats['total_attempts']}")
        print(f"✅ 총 실패: {stats['total_failures']}")
        print(f"✅ 성공률: {stats['success_rate']:.1%}")


class TestRetryLogic(unittest.TestCase):
    """재시도 로직 단위 테스트"""
    
    def test_retry_on_rate_limit_success(self):
        """Rate Limit 재시도 성공 테스트"""
        print("\n=== Rate Limit 재시도 성공 테스트 ===")
        
        call_count = [0]
        
        @retry_on_rate_limit(max_retries=3)
        def api_call():
            call_count[0] += 1
            if call_count[0] < 2:
                raise RateLimitError("Rate limit exceeded")
            return {"status": "success"}
        
        result = api_call()
        
        self.assertEqual(call_count[0], 2)
        self.assertEqual(result["status"], "success")
        print(f"✅ {call_count[0]}번 시도 후 성공")
    
    def test_retry_on_rate_limit_failure(self):
        """Rate Limit 재시도 실패 테스트"""
        print("\n=== Rate Limit 재시도 실패 테스트 ===")
        
        call_count = [0]
        
        @retry_on_rate_limit(max_retries=2)
        def api_call():
            call_count[0] += 1
            raise RateLimitError("Rate limit exceeded")
        
        with self.assertRaises(RateLimitError):
            api_call()
        
        self.assertEqual(call_count[0], 3)  # 초기 + 2회 재시도
        print(f"✅ {call_count[0]}번 시도 후 실패")
    
    def test_retry_on_network_error(self):
        """네트워크 에러 재시도 테스트"""
        print("\n=== 네트워크 에러 재시도 테스트 ===")
        
        call_count = [0]
        
        @retry_on_network_error(max_retries=2)
        def api_call():
            call_count[0] += 1
            if call_count[0] < 2:
                raise ConnectionError("Network error")
            return {"status": "success"}
        
        result = api_call()
        
        self.assertEqual(call_count[0], 2)
        self.assertEqual(result["status"], "success")
        print(f"✅ 네트워크 에러 {call_count[0]}번 시도 후 성공")
    
    def test_non_retryable_error(self):
        """재시도 불가능한 에러 테스트"""
        print("\n=== 재시도 불가능한 에러 테스트 ===")
        
        call_count = [0]
        
        # AuthenticationError 타입은 이미 시스템에 정의되어 있음
        # 새로운 에러 타입을 정의
        class AuthError(Exception):
            pass
        
        from ..rate_limiting.enhanced_backoff_strategy import get_backoff_strategy
        backoff = get_backoff_strategy()
        
        @enhanced_retry(max_retries=5)
        def api_call():
            call_count[0] += 1
            # AuthenticationError 타입은 재시도 불가
            raise AuthError("Authentication failed")
        
        with self.assertRaises(AuthError):
            api_call()
        
        self.assertEqual(call_count[0], 1)  # 한 번만 호출됨
        print(f"✅ 재시도 불가능한 에러로 즉시 실패 (호출 수: {call_count[0]})")
    
    def test_retry_with_circuit_breaker(self):
        """Circuit Breaker와 함께 재시도 테스트"""
        print("\n=== Circuit Breaker 재시도 테스트 ===")
        
        # Backoff 전략 리셋
        from ..rate_limiting.enhanced_backoff_strategy import get_backoff_strategy
        backoff = get_backoff_strategy()
        backoff.reset()
        
        call_count = [0]
        
        @enhanced_retry(max_retries=5, enable_circuit_breaker=True)
        def api_call():
            call_count[0] += 1
            if call_count[0] < 3:
                raise APIError("Server error", code="500")
            return {"status": "success"}
        
        result = api_call()
        
        self.assertEqual(result["status"], "success")
        print(f"✅ Circuit Breaker 활성 상태에서 {call_count[0]}번 시도 후 성공")
    
    def test_retry_with_custom_callback(self):
        """커스텀 콜백과 함께 재시도 테스트"""
        print("\n=== 커스텀 콜백 재시도 테스트 ===")
        
        callback_calls = []
        
        def error_callback(error, retry_count, strategy):
            callback_calls.append({
                'error': str(error),
                'retry_count': retry_count,
                'should_retry': strategy['should_retry']
            })
        
        call_count = [0]
        
        @enhanced_retry(max_retries=2, error_callback=error_callback)
        def api_call():
            call_count[0] += 1
            if call_count[0] < 3:
                raise APIError("Temporary error", code="TEMP_ERROR")
            return {"status": "success"}
        
        result = api_call()
        
        self.assertEqual(len(callback_calls), 2)  # 2번의 에러
        self.assertEqual(result["status"], "success")
        print(f"✅ 콜백 호출 횟수: {len(callback_calls)}")
        for i, call in enumerate(callback_calls):
            print(f"   - 재시도 {i+1}: {call['error']}")


def run_error_handling_tests():
    """테스트 실행 헬퍼"""
    print("\n" + "="*60)
    print("에러 핸들링 단위 테스트 시작")
    print("="*60)
    
    # 두 테스트 클래스 실행
    suite = unittest.TestSuite()
    
    # Exponential Backoff 테스트
    suite.addTests(unittest.TestLoader().loadTestsFromTestCase(TestExponentialBackoff))
    
    # 재시도 로직 테스트
    suite.addTests(unittest.TestLoader().loadTestsFromTestCase(TestRetryLogic))
    
    runner = unittest.TextTestRunner(verbosity=2)
    runner.run(suite)
    
    print("\n" + "="*60)
    print("테스트 완료")
    print("="*60)


if __name__ == "__main__":
    run_error_handling_tests() 