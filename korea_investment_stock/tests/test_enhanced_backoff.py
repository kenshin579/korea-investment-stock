#!/usr/bin/env python3
"""
Enhanced Backoff Strategy 테스트
Date: 2024-12-28

Phase 3.2: 고급 Exponential Backoff 테스트
"""

import time
import os
from unittest.mock import Mock, patch

# 환경 변수 설정 (테스트용)
os.environ['BACKOFF_BASE_DELAY'] = '0.5'  # 테스트 속도를 위해 짧게
os.environ['BACKOFF_MAX_DELAY'] = '4.0'
os.environ['CIRCUIT_FAILURE_THRESHOLD'] = '3'  # 3번 실패 시 Circuit Open
os.environ['CIRCUIT_TIMEOUT'] = '5.0'  # 5초 후 Half Open

from ..rate_limiting.enhanced_backoff_strategy import (
    EnhancedBackoffStrategy, 
    BackoffConfig,
    CircuitState,
    get_backoff_strategy
)


def test_basic_exponential_backoff():
    """기본 Exponential Backoff 테스트"""
    print("\n" + "="*60)
    print("1. 기본 Exponential Backoff 테스트")
    print("="*60)
    
    config = BackoffConfig(
        base_delay=1.0,
        max_delay=10.0,
        exponential_base=2.0,
        jitter_factor=0
    )
    strategy = EnhancedBackoffStrategy(config)
    
    print("\n재시도 횟수별 백오프 시간:")
    for i in range(6):
        delay, reason = strategy.calculate_backoff(i)
        expected = min(2**i, 10.0)
        print(f"  시도 {i+1}: {delay:.1f}초 (예상: {expected}초)")
        assert delay == expected


def test_circuit_breaker():
    """Circuit Breaker 패턴 테스트"""
    print("\n" + "="*60)
    print("2. Circuit Breaker 패턴 테스트")
    print("="*60)
    
    config = BackoffConfig(
        failure_threshold=3,
        success_threshold=2,
        timeout=2.0  # 2초
    )
    strategy = EnhancedBackoffStrategy(config)
    
    # 1. 초기 상태는 CLOSED
    print("\n✓ 초기 상태 확인")
    assert strategy.state == CircuitState.CLOSED
    print(f"  - Circuit 상태: {strategy.state}")
    
    # 2. 3번 실패 시 OPEN
    print("\n✓ 3번 실패 후 OPEN 전환")
    for i in range(3):
        strategy.record_attempt(success=False)
        print(f"  - 실패 {i+1}회: Circuit {strategy.state}")
    assert strategy.state == CircuitState.OPEN
    
    # 3. OPEN 상태에서는 즉시 차단
    delay, reason = strategy.calculate_backoff(0)
    print(f"\n✓ OPEN 상태에서 백오프: {reason}")
    assert "Circuit OPEN" in reason
    
    # 4. 타임아웃 후 HALF_OPEN
    print("\n✓ 2초 대기 후 HALF_OPEN 전환")
    time.sleep(2.1)
    delay, reason = strategy.calculate_backoff(0)
    assert strategy.state == CircuitState.HALF_OPEN
    print(f"  - Circuit 상태: {strategy.state}")
    
    # 5. HALF_OPEN에서 성공 시 CLOSED
    print("\n✓ HALF_OPEN에서 2번 성공 후 CLOSED")
    strategy.record_attempt(success=True)
    print(f"  - 성공 1회: Circuit {strategy.state}")
    strategy.record_attempt(success=True)
    print(f"  - 성공 2회: Circuit {strategy.state}")
    assert strategy.state == CircuitState.CLOSED


def test_adaptive_backoff():
    """Adaptive Backoff 테스트"""
    print("\n" + "="*60)
    print("3. Adaptive Backoff 테스트")
    print("="*60)
    
    config = BackoffConfig(
        base_delay=1.0,
        enable_adaptive=True,
        min_success_rate=0.5,
        jitter_factor=0
    )
    strategy = EnhancedBackoffStrategy(config)
    
    # 성공률 100%일 때
    print("\n✓ 성공률 100%")
    for _ in range(10):
        strategy.record_attempt(success=True)
    
    delay1, _ = strategy.calculate_backoff(1)
    print(f"  - 백오프 시간: {delay1:.1f}초")
    
    # 성공률 20%로 낮춤
    print("\n✓ 성공률 20%")
    for _ in range(80):
        strategy.record_attempt(success=False)
    
    delay2, _ = strategy.calculate_backoff(1)
    print(f"  - 백오프 시간: {delay2:.1f}초 (증가됨)")
    assert delay2 > delay1  # 성공률이 낮으면 더 긴 대기


def test_non_retryable_errors():
    """재시도 불가능한 에러 테스트"""
    print("\n" + "="*60)
    print("4. 재시도 불가능한 에러 테스트")
    print("="*60)
    
    strategy = EnhancedBackoffStrategy()
    
    # 재시도 가능한 에러
    retryable = ["ConnectionError", "TimeoutError", "RateLimitError"]
    for error in retryable:
        assert strategy.should_retry(error) == True
        print(f"✓ {error}: 재시도 가능")
    
    # 재시도 불가능한 에러
    non_retryable = ["AuthenticationError", "InvalidParameterError"]
    for error in non_retryable:
        assert strategy.should_retry(error) == False
        print(f"✗ {error}: 재시도 불가능")


def test_statistics():
    """통계 수집 테스트"""
    print("\n" + "="*60)
    print("5. 통계 수집 테스트")
    print("="*60)
    
    strategy = EnhancedBackoffStrategy()
    
    # 여러 시도 시뮬레이션
    for i in range(10):
        success = i % 3 != 0  # 30% 실패
        strategy.record_attempt(success)
        if not success:
            strategy.calculate_backoff(i // 3)
    
    stats = strategy.get_stats()
    
    print("\n통계 정보:")
    print(f"  - 총 시도: {stats['total_attempts']}")
    print(f"  - 총 실패: {stats['total_failures']}")
    print(f"  - 성공률: {stats['success_rate']:.1%}")
    print(f"  - 평균 백오프: {stats['avg_backoff_time']:.2f}초")
    print(f"  - Circuit Opens: {stats['circuit_opens']}")
    
    assert stats['total_attempts'] == 10
    assert stats['total_failures'] == 4  # 30% 실패
    assert 0.6 <= stats['success_rate'] <= 0.7


def test_integration_with_decorator():
    """retry_on_rate_limit 데코레이터와의 통합 테스트"""
    print("\n" + "="*60)
    print("6. 데코레이터 통합 테스트")
    print("="*60)
    
    from .korea_investment_stock import API_RETURN_CODE
    from .enhanced_retry_decorator import retry_on_rate_limit
    
    # 백오프 전략 리셋
    strategy = get_backoff_strategy()
    strategy.reset()
    
    class MockAPI:
        def __init__(self):
            self.call_count = 0
            self.rate_limiter = Mock()
            self.rate_limiter.record_error = Mock()
        
        @retry_on_rate_limit(max_retries=3)
        def api_call(self):
            self.call_count += 1
            
            # 처음 2번은 rate limit 에러
            if self.call_count <= 2:
                return {
                    "rt_cd": API_RETURN_CODE["RATE_LIMIT_EXCEEDED"],
                    "msg1": "Rate limit exceeded"
                }
            else:
                return {"rt_cd": API_RETURN_CODE["SUCCESS"]}
    
    api = MockAPI()
    
    # 시간 측정
    start_time = time.time()
    result = api.api_call()
    elapsed = time.time() - start_time
    
    print(f"\n결과:")
    print(f"  - API 호출 횟수: {api.call_count}")
    print(f"  - 최종 결과: {'성공' if result['rt_cd'] == '0' else '실패'}")
    print(f"  - 소요 시간: {elapsed:.1f}초")
    
    # 백오프 통계 확인
    stats = strategy.get_stats()
    print(f"\n백오프 통계:")
    print(f"  - 총 시도: {stats['total_attempts']}")
    print(f"  - 성공률: {stats['success_rate']:.1%}")
    
    assert api.call_count == 3  # 2번 실패 + 1번 성공
    assert result['rt_cd'] == API_RETURN_CODE["SUCCESS"]


def main():
    """메인 테스트 실행"""
    print("Enhanced Backoff Strategy 테스트")
    print("Phase 3.2 구현 검증")
    
    test_basic_exponential_backoff()
    test_circuit_breaker()
    test_adaptive_backoff()
    test_non_retryable_errors()
    test_statistics()
    test_integration_with_decorator()
    
    print("\n" + "="*60)
    print("✅ 모든 테스트 통과!")
    print("Enhanced Backoff Strategy가 정상 작동합니다.")
    print("="*60)


if __name__ == "__main__":
    main() 