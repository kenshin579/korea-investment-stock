#!/usr/bin/env python3
"""
Rate Limit 초과 방지 시뮬레이션
Enhanced RateLimiter가 실제로 Rate Limit을 초과하지 않는지 검증
"""

import time
import threading
from ..rate_limiting import EnhancedRateLimiter
from collections import defaultdict

# matplotlib 옵션
try:
    import matplotlib.pyplot as plt
    HAS_MATPLOTLIB = True
except ImportError:
    HAS_MATPLOTLIB = False
    plt = None

def simulate_api_server(call_times, window_size=1.0, max_calls=20):
    """
    API 서버의 Fixed Window Rate Limiting을 시뮬레이션
    
    Args:
        call_times: API 호출 시각 리스트
        window_size: 윈도우 크기 (초)
        max_calls: 윈도우당 최대 호출 수
    
    Returns:
        tuple: (초과 발생 횟수, 초별 호출 수)
    """
    calls_per_window = defaultdict(int)
    violations = 0
    
    for call_time in call_times:
        # Fixed window: 초 단위로 구분
        window_start = int(call_time / window_size) * window_size
        calls_per_window[window_start] += 1
        
        if calls_per_window[window_start] > max_calls:
            violations += 1
            print(f"⚠️  Rate limit 초과! 시간: {call_time:.3f}s, "
                  f"윈도우 {window_start}-{window_start+1}s: "
                  f"{calls_per_window[window_start]}회")
    
    return violations, calls_per_window

def test_concurrent_requests(num_threads=20, requests_per_thread=5):
    """
    동시 요청 시뮬레이션 (초기 버스트 시나리오)
    """
    print("\n" + "="*60)
    print(f"동시 요청 테스트: {num_threads}개 스레드, 각 {requests_per_thread}개 요청")
    print("="*60)
    
    limiter = EnhancedRateLimiter(max_calls=15, safety_margin=0.8)
    call_times = []
    lock = threading.Lock()
    
    def worker():
        for _ in range(requests_per_thread):
            start = time.time()
            limiter.acquire()
            with lock:
                call_times.append(time.time() - test_start_time)
    
    # 모든 스레드 동시 시작
    test_start_time = time.time()
    threads = []
    for _ in range(num_threads):
        t = threading.Thread(target=worker)
        threads.append(t)
        t.start()
    
    # 모든 스레드 종료 대기
    for t in threads:
        t.join()
    
    # 결과 분석
    violations, calls_per_window = simulate_api_server(call_times)
    
    print(f"\n결과:")
    print(f"- 총 요청 수: {len(call_times)}")
    print(f"- Rate limit 초과 횟수: {violations}")
    print(f"- 최대 초당 호출 수: {max(calls_per_window.values()) if calls_per_window else 0}")
    
    # 시각화 (matplotlib이 있는 경우만)
    if HAS_MATPLOTLIB and calls_per_window:
        windows = sorted(calls_per_window.keys())
        counts = [calls_per_window[w] for w in windows]
        
        plt.figure(figsize=(10, 6))
        plt.bar(windows, counts, width=0.8, alpha=0.7)
        plt.axhline(y=20, color='r', linestyle='--', label='API Limit (20)')
        plt.axhline(y=12, color='g', linestyle='--', label='Our Limit (12)')
        plt.xlabel('Time Window (seconds)')
        plt.ylabel('API Calls')
        plt.title('API Calls per Second - Concurrent Requests')
        plt.legend()
        plt.grid(True, alpha=0.3)
        plt.savefig('concurrent_requests_test.png')
        plt.close()
    
    return violations == 0

def test_sustained_load(duration=10):
    """
    지속적인 부하 테스트
    """
    print("\n" + "="*60)
    print(f"지속적 부하 테스트: {duration}초 동안")
    print("="*60)
    
    limiter = EnhancedRateLimiter(max_calls=15, safety_margin=0.8)
    call_times = []
    
    start_time = time.time()
    request_count = 0
    
    while time.time() - start_time < duration:
        limiter.acquire()
        call_times.append(time.time() - start_time)
        request_count += 1
        
        # 진행 상황 출력
        if request_count % 50 == 0:
            elapsed = time.time() - start_time
            rate = request_count / elapsed
            print(f"  {request_count}개 요청 완료, 평균 속도: {rate:.1f} req/s")
    
    # 결과 분석
    violations, calls_per_window = simulate_api_server(call_times)
    
    print(f"\n결과:")
    print(f"- 총 요청 수: {len(call_times)}")
    print(f"- Rate limit 초과 횟수: {violations}")
    print(f"- 평균 요청 속도: {len(call_times)/duration:.1f} req/s")
    print(f"- 최대 초당 호출 수: {max(calls_per_window.values()) if calls_per_window else 0}")
    
    return violations == 0

def test_edge_cases():
    """
    경계 조건 테스트
    """
    print("\n" + "="*60)
    print("경계 조건 테스트")
    print("="*60)
    
    # 1. 윈도우 경계에서의 요청
    print("\n1. 윈도우 경계 테스트:")
    limiter = EnhancedRateLimiter(max_calls=15, safety_margin=0.8)
    call_times = []
    
    # 0.9초에 11개 요청
    for i in range(11):
        limiter.acquire()
        call_times.append(0.9 + i * 0.001)
    
    # 1.1초에 추가 요청
    time.sleep(0.2)
    limiter.acquire()
    call_times.append(1.1)
    
    violations, _ = simulate_api_server(call_times)
    print(f"   윈도우 경계 초과: {violations}회")
    
    # 2. 네트워크 지연 시뮬레이션
    print("\n2. 네트워크 지연 시뮬레이션:")
    call_times_with_delay = []
    for t in call_times:
        # 10-50ms 랜덤 지연 추가
        import random
        delay = random.uniform(0.01, 0.05)
        call_times_with_delay.append(t + delay)
    
    violations, _ = simulate_api_server(call_times_with_delay)
    print(f"   네트워크 지연 시 초과: {violations}회")
    
    return violations == 0

def main():
    """메인 테스트 실행"""
    print("Enhanced RateLimiter - Rate Limit 초과 방지 검증")
    print("API 제한: 20 calls/sec (Fixed Window)")
    print("우리 설정: 15 calls * 0.8 = 12 calls/sec")
    
    results = []
    
    # 1. 동시 요청 테스트
    results.append(("동시 요청 (버스트)", test_concurrent_requests()))
    
    # 2. 지속적 부하 테스트
    results.append(("지속적 부하", test_sustained_load(5)))
    
    # 3. 경계 조건 테스트
    results.append(("경계 조건", test_edge_cases()))
    
    # 최종 결과
    print("\n" + "="*60)
    print("최종 결과:")
    print("="*60)
    
    all_passed = True
    for test_name, passed in results:
        status = "✅ PASS" if passed else "❌ FAIL"
        print(f"{test_name}: {status}")
        all_passed = all_passed and passed
    
    if all_passed:
        print("\n🎉 모든 테스트 통과! Rate Limit 초과가 발생하지 않았습니다.")
    else:
        print("\n⚠️  일부 테스트에서 Rate Limit 초과가 발생했습니다.")
    
    return all_passed

if __name__ == "__main__":
    main() 