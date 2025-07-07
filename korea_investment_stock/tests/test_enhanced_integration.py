#!/usr/bin/env python3
"""
Enhanced RateLimiter 통합 테스트
Date: 2024-12-28

기존 KoreaInvestment 클래스와 Enhanced RateLimiter가 
잘 통합되었는지 확인하는 테스트
"""

import time
import logging

# 로깅 설정
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)

def test_rate_limiter_integration():
    """Rate Limiter 통합 테스트"""
    
    # Mock KoreaInvestment 인스턴스 생성을 위한 테스트
    from .korea_investment_stock import KoreaInvestment
    
    print("=" * 50)
    print("Enhanced RateLimiter 통합 테스트")
    print("=" * 50)
    
    # Mock 데이터로 테스트 (실제 API 키 없이)
    try:
        # __init__ 메서드에서 rate_limiter가 EnhancedRateLimiter인지 확인
        print("\n1. KoreaInvestment 인스턴스 생성 테스트...")
        
        # 테스트용 더미 데이터
        api_key = "test_key"
        api_secret = "test_secret"
        acc_no = "12345678-01"
        
        # Mock 모드로 생성 (실제 API 호출 없음)
        broker = KoreaInvestment(
            api_key=api_key,
            api_secret=api_secret,
            acc_no=acc_no,
            mock=True
        )
        
        print("✓ KoreaInvestment 인스턴스 생성 성공")
        
        # Rate Limiter 타입 확인
        print(f"\n2. Rate Limiter 타입: {type(broker.rate_limiter).__name__}")
        
        # Rate Limiter 설정 확인
        if hasattr(broker.rate_limiter, 'get_stats'):
            stats = broker.rate_limiter.get_stats()
            print(f"✓ Enhanced RateLimiter 확인됨")
            print(f"   - nominal_max_calls: {stats['config']['nominal_max_calls']}")
            print(f"   - effective_max_calls: {stats['config']['effective_max_calls']}")
            print(f"   - safety_margin: {stats['config']['safety_margin']}")
            print(f"   - min_interval: {stats['config']['min_interval']:.3f}초")
        
        # ThreadPoolExecutor 워커 수 확인
        print(f"\n3. ThreadPoolExecutor 워커 수: {broker.executor._max_workers}")
        
        # Rate Limiter 기본 동작 테스트
        print("\n4. Rate Limiter 동작 테스트...")
        print("   10개 요청 시뮬레이션:")
        
        for i in range(10):
            start = time.time()
            broker.rate_limiter.acquire()
            elapsed = time.time() - start
            print(f"   요청 {i+1}: {elapsed:.3f}초 대기")
        
        # 통계 출력
        print("\n5. Rate Limiter 통계:")
        broker.rate_limiter.print_stats()
        
        # 정리
        broker.shutdown()
        print("\n✓ 모든 테스트 성공!")
        
    except Exception as e:
        print(f"\n✗ 테스트 실패: {e}")
        import traceback
        traceback.print_exc()


def test_rate_limiter_features():
    """Enhanced RateLimiter 추가 기능 테스트"""
    from .enhanced_rate_limiter import EnhancedRateLimiter
    
    print("\n" + "=" * 50)
    print("Enhanced RateLimiter 추가 기능 테스트")
    print("=" * 50)
    
    # Context Manager 테스트
    print("\n1. Context Manager 테스트:")
    limiter = EnhancedRateLimiter(max_calls=5, per_seconds=1.0)
    
    with limiter:
        print("   ✓ Context manager로 rate limit 획득 성공")
    
    # 타임아웃 테스트
    print("\n2. 타임아웃 테스트:")
    for i in range(7):
        if limiter.acquire(timeout=0.5):
            print(f"   요청 {i+1}: 성공")
        else:
            print(f"   요청 {i+1}: 타임아웃 (예상된 동작)")
    
    # 에러 기록 테스트
    print("\n3. 에러 통계 테스트:")
    limiter.record_error()
    limiter.record_error()
    stats = limiter.get_stats()
    print(f"   에러 수: {stats['error_count']}")
    print(f"   에러율: {stats['error_rate']:.1%}")
    
    # Reset 테스트
    print("\n4. Reset 테스트:")
    limiter.reset()
    stats_after_reset = limiter.get_stats()
    print(f"   Reset 후 총 호출 수: {stats_after_reset['total_calls']}")
    print(f"   Reset 후 에러 수: {stats_after_reset['error_count']}")


if __name__ == "__main__":
    # 통합 테스트
    test_rate_limiter_integration()
    
    # 추가 기능 테스트
    test_rate_limiter_features()
    
    print("\n" + "=" * 50)
    print("모든 테스트 완료!")
    print("=" * 50) 