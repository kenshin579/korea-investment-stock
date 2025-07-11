"""
Legacy RateLimiter 클래스 백업
Date: 2024-12-28
Branch: feat/#27-rate-limit

이 파일은 Rate Limiting 개선 작업 전의 원본 RateLimiter 클래스를 보관합니다.
원본 위치: korea_investment_stock/koreainvestmentstock.py (line 778~)

특징:
- Sliding Window + Token Bucket 하이브리드 방식
- safety_margin 0.9 적용 (max_calls의 90%만 사용)
- 통계 수집 기능 포함
- 추가 안전장치로 최소 간격 보장
"""

import threading
import time
from collections import deque, defaultdict
import datetime


class RateLimiter:
    def __init__(self, max_calls, per_seconds, safety_margin=0.9):
        """
        Args:
            max_calls: 시간 윈도우 내 최대 호출 횟수
            per_seconds: 시간 윈도우 (초)
            safety_margin: 안전 마진 (0.9 = 실제 제한의 90%만 사용)
        """
        self.max_calls = int(max_calls * safety_margin)  # 안전 마진 적용
        self.per_seconds = per_seconds
        self.lock = threading.Lock()
        self.call_timestamps = deque()
        
        # 호출 통계 추적
        self.calls_per_second = defaultdict(int)
        
        # 버킷 토큰 방식 추가
        self.tokens = self.max_calls
        self.last_update = time.time()
        self.refill_rate = self.max_calls / self.per_seconds

    def acquire(self):
        with self.lock:
            now = time.time()
            
            # 토큰 리필
            time_passed = now - self.last_update
            self.tokens = min(self.max_calls, self.tokens + time_passed * self.refill_rate)
            self.last_update = now
            
            # 슬라이딩 윈도우 방식도 병행
            while self.call_timestamps and self.call_timestamps[0] <= now - self.per_seconds:
                self.call_timestamps.popleft()
            
            # 두 가지 조건 모두 확인
            while self.tokens < 1 or len(self.call_timestamps) >= self.max_calls:
                # 다음 토큰이 사용 가능할 때까지 대기
                if self.tokens < 1:
                    wait_time = (1 - self.tokens) / self.refill_rate
                else:
                    wait_time = self.per_seconds - (now - self.call_timestamps[0])
                
                if wait_time > 0:
                    time.sleep(min(wait_time + 0.001, 0.1))  # 최소 1ms, 최대 100ms 대기
                    now = time.time()
                    
                    # 토큰 리필
                    time_passed = now - self.last_update
                    self.tokens = min(self.max_calls, self.tokens + time_passed * self.refill_rate)
                    self.last_update = now
                    
                    # 만료된 타임스탬프 제거
                    while self.call_timestamps and self.call_timestamps[0] <= now - self.per_seconds:
                        self.call_timestamps.popleft()
            
            # 토큰 사용
            self.tokens -= 1
            
            # 호출 기록
            current_second = int(now)
            self.calls_per_second[current_second] += 1
            self.call_timestamps.append(now)
            
            # 추가 안전장치: 마지막 호출 후 최소 간격 보장
            min_interval = self.per_seconds / (self.max_calls * 1.2)  # 20% 여유
            time.sleep(min_interval)

    def get_stats(self):
        """호출 통계 분석 결과 반환"""
        stats = {
            "calls_per_second": dict(self.calls_per_second),
            "max_calls_in_one_second": max(self.calls_per_second.values()) if self.calls_per_second else 0,
            "total_calls": sum(self.calls_per_second.values()),
            "seconds_tracked": len(self.calls_per_second)
        }
        return stats

    def print_stats(self):
        """호출 통계 출력"""
        if not self.calls_per_second:
            print("호출 데이터가 없습니다.")
            return

        print("\n===== 초당 API 호출 횟수 분석 =====")
        max_calls = max(self.calls_per_second.values())

        for second, count in sorted(self.calls_per_second.items()):
            timestamp = datetime.datetime.fromtimestamp(second).strftime('%H:%M:%S')
            print(f"시간: {timestamp}, 호출 수: {count}")

        print(f"\n최대 초당 호출 횟수: {max_calls}")
        print(f"설정된 max_calls: {self.max_calls}")
        print(f"제한 준수 여부: {'준수' if max_calls <= self.max_calls else '초과'}")
        print(f"총 호출 횟수: {sum(self.calls_per_second.values())}")
        print("================================\n") 