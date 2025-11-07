import time
import threading
from typing import Dict, Any


class RateLimiter:
    """
    API 호출 속도 제한을 위한 스레드 안전 Rate Limiter

    토큰 버킷 알고리즘 기반 단순 구현
    """

    def __init__(self, calls_per_second: float = 15.0):
        """
        Args:
            calls_per_second: 초당 최대 API 호출 수 (기본값: 15)

        Raises:
            ValueError: calls_per_second가 0 이하인 경우
        """
        if calls_per_second <= 0:
            raise ValueError("calls_per_second must be positive")

        self._calls_per_second = calls_per_second
        self._min_interval = 1.0 / calls_per_second
        self._last_call = 0.0
        self._total_calls = 0
        self._lock = threading.Lock()

    def wait(self) -> None:
        """
        다음 API 호출이 허용될 때까지 대기

        속도 제한을 초과하면 자동으로 sleep하여 속도 제한을 준수합니다.
        스레드 안전하게 동작합니다.
        """
        with self._lock:
            current_time = time.time()
            elapsed = current_time - self._last_call

            if elapsed < self._min_interval:
                sleep_time = self._min_interval - elapsed
                time.sleep(sleep_time)
                current_time = time.time()

            self._last_call = current_time
            self._total_calls += 1

    def get_stats(self) -> Dict[str, Any]:
        """
        속도 제한 통계 반환

        Returns:
            {
                'calls_per_second': float,  # 현재 설정된 초당 호출 수
                'min_interval': float,      # 최소 호출 간격 (초)
                'last_call': float,         # 마지막 호출 시각 (timestamp)
                'total_calls': int          # 총 호출 횟수
            }
        """
        with self._lock:
            return {
                'calls_per_second': self._calls_per_second,
                'min_interval': self._min_interval,
                'last_call': self._last_call,
                'total_calls': self._total_calls
            }

    def adjust_rate_limit(self, calls_per_second: float) -> None:
        """
        런타임 중 속도 제한 동적 조정

        Args:
            calls_per_second: 새로운 초당 호출 수

        Raises:
            ValueError: calls_per_second가 0 이하인 경우
        """
        if calls_per_second <= 0:
            raise ValueError("calls_per_second must be positive")

        with self._lock:
            self._calls_per_second = calls_per_second
            self._min_interval = 1.0 / calls_per_second
