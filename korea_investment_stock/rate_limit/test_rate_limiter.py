import time
import threading
import pytest
from .rate_limiter import RateLimiter


def test_rate_limiter_basic():
    """기본 속도 제한 테스트"""
    limiter = RateLimiter(calls_per_second=10)

    start = time.time()
    for _ in range(20):
        limiter.wait()
    elapsed = time.time() - start

    # 10회/초로 20번 호출하면 약 2초 소요
    # 허용 오차: ±10%
    assert 1.8 <= elapsed <= 2.2, f"Expected 1.8-2.2s, got {elapsed:.2f}s"


def test_rate_limiter_thread_safe():
    """멀티스레드 안전성 테스트"""
    limiter = RateLimiter(calls_per_second=10)
    results = []

    def worker():
        for _ in range(10):
            limiter.wait()
            results.append(time.time())

    # 3개 스레드가 각각 10번씩 호출 = 총 30회
    threads = [threading.Thread(target=worker) for _ in range(3)]
    for t in threads:
        t.start()
    for t in threads:
        t.join()

    # 호출 횟수 검증
    assert len(results) == 30, f"Expected 30 calls, got {len(results)}"

    # 시간순 정렬 후 간격 확인
    results.sort()
    intervals = [results[i+1] - results[i] for i in range(len(results)-1)]

    # 0.1초 간격 (10회/초), 허용 오차: -10%
    # 일부 호출은 동시에 처리될 수 있으므로 최소 간격만 검증
    min_interval = 1.0 / 10  # 0.1초
    tolerance = 0.09  # 허용 오차

    # 대부분의 간격이 최소 간격을 만족하는지 검증
    valid_intervals = [interval >= tolerance for interval in intervals]
    valid_percentage = sum(valid_intervals) / len(valid_intervals)

    assert valid_percentage >= 0.9, \
        f"Expected 90% of intervals >= {tolerance}s, got {valid_percentage*100:.1f}%"


def test_rate_limiter_stats():
    """통계 정보 테스트"""
    limiter = RateLimiter(calls_per_second=15)

    # 5번 호출
    for _ in range(5):
        limiter.wait()

    stats = limiter.get_stats()

    # 통계 검증
    assert stats['calls_per_second'] == 15, \
        f"Expected calls_per_second=15, got {stats['calls_per_second']}"
    assert stats['total_calls'] == 5, \
        f"Expected total_calls=5, got {stats['total_calls']}"
    assert stats['min_interval'] == pytest.approx(1.0 / 15), \
        f"Expected min_interval={1.0/15:.4f}, got {stats['min_interval']:.4f}"


def test_rate_limiter_adjust():
    """동적 속도 조정 테스트"""
    limiter = RateLimiter(calls_per_second=10)

    # 초기 설정 확인
    stats = limiter.get_stats()
    assert stats['calls_per_second'] == 10, \
        f"Expected calls_per_second=10, got {stats['calls_per_second']}"

    # 속도 조정
    limiter.adjust_rate_limit(calls_per_second=20)

    # 변경 확인
    stats = limiter.get_stats()
    assert stats['calls_per_second'] == 20, \
        f"Expected calls_per_second=20, got {stats['calls_per_second']}"
    assert stats['min_interval'] == pytest.approx(1.0 / 20), \
        f"Expected min_interval={1.0/20:.4f}, got {stats['min_interval']:.4f}"


def test_rate_limiter_invalid_input():
    """잘못된 입력 테스트"""
    # 0 이하의 값은 ValueError 발생
    with pytest.raises(ValueError, match="calls_per_second must be positive"):
        RateLimiter(calls_per_second=0)

    with pytest.raises(ValueError, match="calls_per_second must be positive"):
        RateLimiter(calls_per_second=-1)

    # adjust_rate_limit에서도 동일한 검증
    limiter = RateLimiter(calls_per_second=10)

    with pytest.raises(ValueError, match="calls_per_second must be positive"):
        limiter.adjust_rate_limit(calls_per_second=0)

    with pytest.raises(ValueError, match="calls_per_second must be positive"):
        limiter.adjust_rate_limit(calls_per_second=-1)


def test_rate_limiter_precision():
    """속도 제한 정밀도 테스트"""
    limiter = RateLimiter(calls_per_second=15)

    start = time.time()
    for _ in range(30):
        limiter.wait()
    elapsed = time.time() - start

    # 15회/초로 30번 호출 = 2초, 허용 오차: ±10%
    expected = 30 / 15  # 2초
    assert expected * 0.9 <= elapsed <= expected * 1.1, \
        f"Expected {expected*0.9:.2f}-{expected*1.1:.2f}s, got {elapsed:.2f}s"


def test_rate_limiter_zero_wait():
    """첫 호출 시 대기 없음 테스트"""
    limiter = RateLimiter(calls_per_second=10)

    # 첫 호출은 즉시 반환되어야 함
    start = time.time()
    limiter.wait()
    elapsed = time.time() - start

    # 첫 호출은 10ms 이내에 완료되어야 함
    assert elapsed < 0.01, f"First call should be immediate, took {elapsed:.4f}s"
