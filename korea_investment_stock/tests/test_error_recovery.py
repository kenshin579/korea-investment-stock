#!/usr/bin/env python3
"""
Error Recovery System Test
Date: 2024-12-28

Phase 3.3 테스트: 에러 복구 시스템 테스트
"""

import unittest
import time
from unittest.mock import Mock, patch
from datetime import datetime, timedelta

from ..error_handling.error_recovery_system import (
    ErrorRecoverySystem,
    ErrorPattern,
    ErrorSeverity,
    RecoveryAction,
    handle_api_error,
    get_error_recovery_system
)
from ..rate_limiting.enhanced_retry_decorator import (
    retry_on_rate_limit,
    enhanced_retry,
    RateLimitError,
    TokenExpiredError,
    APIError
)


class TestErrorRecoverySystem(unittest.TestCase):
    """에러 복구 시스템 테스트"""
    
    def setUp(self):
        """테스트 환경 설정"""
        self.recovery_system = ErrorRecoverySystem()
        
    def test_error_pattern_matching(self):
        """에러 패턴 매칭 테스트"""
        print("\n=== 에러 패턴 매칭 테스트 ===")
        
        # Rate Limit 에러
        error = RateLimitError("Rate limit exceeded")
        response = {'rt_cd': 'EGW00201'}
        strategy = self.recovery_system.handle_error(error, {}, response)
        
        self.assertTrue(strategy['should_retry'])
        self.assertEqual(strategy['pattern'].severity, ErrorSeverity.MEDIUM)
        self.assertIn(RecoveryAction.RETRY, strategy['recovery_actions'])
        print(f"✅ Rate Limit 에러 패턴 매칭: {strategy['user_message']}")
        
        # 토큰 만료 에러
        error = TokenExpiredError("Token expired")
        response = {'rt_cd': '1'}
        strategy = self.recovery_system.handle_error(error, {}, response)
        
        self.assertTrue(strategy['should_retry'])
        self.assertEqual(strategy['pattern'].severity, ErrorSeverity.LOW)
        self.assertIn(RecoveryAction.REFRESH_TOKEN, strategy['recovery_actions'])
        print(f"✅ 토큰 만료 에러 패턴 매칭: {strategy['user_message']}")
        
        # 네트워크 에러
        error = ConnectionError("Network error")
        strategy = self.recovery_system.handle_error(error)
        
        self.assertTrue(strategy['should_retry'])
        self.assertEqual(strategy['pattern'].severity, ErrorSeverity.LOW)
        print(f"✅ 네트워크 에러 패턴 매칭: {strategy['user_message']}")
        
    def test_retry_on_rate_limit_decorator(self):
        """Rate Limit 재시도 데코레이터 테스트"""
        print("\n=== Rate Limit 재시도 데코레이터 테스트 ===")
        
        call_count = [0]  # 리스트로 만들어 클로저에서 수정 가능하게 함
        
        @retry_on_rate_limit(max_retries=3)
        def api_call():
            call_count[0] += 1
            if call_count[0] < 3:
                raise RateLimitError("Rate limit exceeded")
            return {"status": "success"}
        
        # 테스트 실행
        result = api_call()
        
        self.assertEqual(call_count[0], 3)
        self.assertEqual(result["status"], "success")
        print(f"✅ 3회 시도 후 성공 (총 호출 수: {call_count[0]})")
        
    def test_enhanced_retry_with_circuit_breaker(self):
        """Circuit Breaker를 포함한 향상된 재시도 테스트"""
        print("\n=== Circuit Breaker 테스트 ===")
        
        call_count = [0]  # 리스트로 만들어 클로저에서 수정 가능하게 함
        
        @enhanced_retry(max_retries=5, enable_circuit_breaker=True)
        def flaky_api_call():
            call_count[0] += 1
            if call_count[0] < 2:
                raise APIError("Server error", code="500")
            return {"status": "success"}
        
        # 테스트 실행
        result = flaky_api_call()
        
        self.assertEqual(result["status"], "success")
        print(f"✅ Circuit Breaker 활성화 상태에서 재시도 성공 (시도 횟수: {call_count[0]})")
        
    def test_error_statistics(self):
        """에러 통계 수집 테스트"""
        print("\n=== 에러 통계 수집 테스트 ===")
        
        # 다양한 에러 시뮬레이션
        errors = [
            (RateLimitError("Rate limit"), {'rt_cd': 'EGW00201'}),
            (RateLimitError("Rate limit"), {'rt_cd': 'EGW00201'}),
            (TokenExpiredError("Token"), {'rt_cd': '1'}),
            (ConnectionError("Network"), None),
        ]
        
        for error, response in errors:
            self.recovery_system.handle_error(error, response=response)
        
        # 통계 확인
        summary = self.recovery_system.get_error_summary(hours=1)
        
        self.assertEqual(summary['total_errors'], 4)
        self.assertEqual(summary['by_severity']['MEDIUM'], 2)  # Rate limit errors
        self.assertEqual(summary['by_severity']['LOW'], 2)  # Token + Network
        
        print(f"✅ 총 에러 수: {summary['total_errors']}")
        print(f"✅ 심각도별 분포: {summary['by_severity']}")
        print(f"✅ 에러 타입별 분포: {summary['by_type']}")
        
    def test_recovery_callbacks(self):
        """복구 콜백 테스트"""
        print("\n=== 복구 콜백 테스트 ===")
        
        self.callback_called = False
        
        def refresh_token_callback(pattern, event):
            self.callback_called = True
            print(f"✅ 토큰 갱신 콜백 호출됨: {pattern.message_template}")
        
        # 콜백 등록
        self.recovery_system.register_callback(
            RecoveryAction.REFRESH_TOKEN,
            refresh_token_callback
        )
        
        # 토큰 만료 에러 발생
        error = TokenExpiredError("Token expired")
        response = {'rt_cd': '1'}
        self.recovery_system.handle_error(error, response=response)
        
        self.assertTrue(self.callback_called)
        
    def test_non_retryable_errors(self):
        """재시도 불가능한 에러 테스트"""
        print("\n=== 재시도 불가능한 에러 테스트 ===")
        
        # 인증 에러 (재시도 불가)
        error = APIError("Authentication failed", code="AUTH_FAILED")
        self.recovery_system.error_patterns.append(
            ErrorPattern(
                error_type="APIError",
                error_code="AUTH_FAILED",
                severity=ErrorSeverity.HIGH,
                recovery_actions=[RecoveryAction.NOTIFY_USER, RecoveryAction.FAIL_FAST],
                max_retries=0,
                wait_time=0,
                message_template="인증 실패. API 키를 확인해주세요."
            )
        )
        
        strategy = self.recovery_system.handle_error(error, response={'rt_cd': 'AUTH_FAILED'})
        
        self.assertFalse(strategy['should_retry'])
        self.assertEqual(strategy['max_retries'], 0)
        print(f"✅ 재시도 불가능한 에러 감지: {strategy['user_message']}")


def run_error_recovery_tests():
    """테스트 실행 헬퍼"""
    print("\n" + "="*60)
    print("에러 복구 시스템 테스트 시작")
    print("="*60)
    
    unittest.main(TestErrorRecoverySystem(), argv=[''], exit=False)
    
    print("\n" + "="*60)
    print("테스트 완료")
    print("="*60)


if __name__ == "__main__":
    run_error_recovery_tests() 