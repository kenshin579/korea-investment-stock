"""
Enhanced RateLimiter 구현
Date: 2024-12-28
Issue: #27 - Rate Limiting 개선

특징:
- Token Bucket + Sliding Window 하이브리드 방식
- 보수적 설정값 적용 (max_calls=15, safety_margin=0.8)
- 최소 간격 보장 메커니즘
- 향상된 에러 처리 및 통계 수집
"""

import threading
import time
from collections import deque, defaultdict
from typing import Optional, Dict, Any
import datetime
import logging
import os
import json
from pathlib import Path

# 로거 설정
logger = logging.getLogger(__name__)


class EnhancedRateLimiter:
    """향상된 Rate Limiter 구현
    
    Token Bucket과 Sliding Window를 병합한 하이브리드 방식으로
    더 안정적이고 예측 가능한 Rate Limiting을 제공합니다.
    """
    
    def __init__(
        self, 
        max_calls: int = 15,  # 기본값 20 -> 15로 감소
        per_seconds: float = 1.0,
        safety_margin: float = 0.8,  # 기본값 0.9 -> 0.8로 감소
        enable_min_interval: bool = True,
        enable_stats: bool = True
    ):
        """
        Args:
            max_calls: 시간 윈도우 내 최대 호출 횟수
            per_seconds: 시간 윈도우 (초)
            safety_margin: 안전 마진 (0.8 = 실제 제한의 80%만 사용)
            enable_min_interval: 최소 간격 보장 기능 활성화
            enable_stats: 통계 수집 기능 활성화
        """
        # 설정값 적용
        self.nominal_max_calls = max_calls
        self.max_calls = int(max_calls * safety_margin)
        self.per_seconds = per_seconds
        self.safety_margin = safety_margin
        self.enable_min_interval = enable_min_interval
        self.enable_stats = enable_stats
        
        # 동시성 제어
        self.lock = threading.Lock()
        
        # Sliding Window 관련
        self.call_timestamps = deque()
        
        # Token Bucket 관련
        self.tokens = float(self.max_calls)
        self.last_update = time.time()
        self.refill_rate = self.max_calls / self.per_seconds
        
        # 통계 수집
        if self.enable_stats:
            self.calls_per_second = defaultdict(int)
            self.wait_times = []
            self.error_count = 0
            self.total_calls = 0
        
        # 최소 간격 계산
        if self.enable_min_interval:
            # 균등 분산을 위한 최소 간격 (추가 여유 20%)
            self.min_interval = self.per_seconds / (self.max_calls * 1.2)
        else:
            self.min_interval = 0
        
        # 환경 변수로 설정 오버라이드 가능
        self._load_env_config()
        
        logger.info(
            f"EnhancedRateLimiter 초기화: "
            f"max_calls={self.max_calls} (nominal={self.nominal_max_calls}), "
            f"safety_margin={self.safety_margin}, "
            f"min_interval={self.min_interval:.3f}s"
        )
    
    def _load_env_config(self):
        """환경 변수에서 설정값 로드"""
        if max_calls := os.getenv('RATE_LIMIT_MAX_CALLS'):
            self.nominal_max_calls = int(max_calls)
            self.max_calls = int(int(max_calls) * self.safety_margin)
            logger.info(f"환경 변수에서 max_calls 로드: {self.nominal_max_calls}")
        
        if safety_margin := os.getenv('RATE_LIMIT_SAFETY_MARGIN'):
            self.safety_margin = float(safety_margin)
            self.max_calls = int(self.nominal_max_calls * self.safety_margin)
            logger.info(f"환경 변수에서 safety_margin 로드: {self.safety_margin}")
    
    def _refill_tokens(self, now: float) -> None:
        """Token Bucket 리필 로직"""
        time_passed = now - self.last_update
        tokens_to_add = time_passed * self.refill_rate
        
        # 토큰 추가 (최대값 제한)
        self.tokens = min(self.max_calls, self.tokens + tokens_to_add)
        self.last_update = now
        
        logger.debug(
            f"토큰 리필: +{tokens_to_add:.2f} tokens, "
            f"현재 토큰: {self.tokens:.2f}/{self.max_calls}"
        )
    
    def _clean_old_timestamps(self, now: float) -> None:
        """만료된 타임스탬프 제거"""
        cutoff_time = now - self.per_seconds
        removed_count = 0
        
        while self.call_timestamps and self.call_timestamps[0] <= cutoff_time:
            self.call_timestamps.popleft()
            removed_count += 1
        
        if removed_count > 0:
            logger.debug(f"만료된 타임스탬프 {removed_count}개 제거")
    
    def acquire(self, timeout: Optional[float] = None) -> bool:
        """Rate limit 획득
        
        Args:
            timeout: 최대 대기 시간 (초). None이면 무한 대기
            
        Returns:
            bool: 성공 시 True, 타임아웃 시 False
        """
        start_time = time.time()
        
        with self.lock:
            while True:
                now = time.time()
                
                # 타임아웃 체크
                if timeout is not None and (now - start_time) > timeout:
                    logger.warning(f"Rate limiter 타임아웃: {timeout}초 초과")
                    return False
                
                # Token Bucket 리필
                self._refill_tokens(now)
                
                # Sliding Window 정리
                self._clean_old_timestamps(now)
                
                # 두 가지 조건 모두 확인
                has_token = self.tokens >= 1
                under_window_limit = len(self.call_timestamps) < self.max_calls
                
                if has_token and under_window_limit:
                    # 토큰 사용
                    self.tokens -= 1
                    
                    # 타임스탬프 기록
                    self.call_timestamps.append(now)
                    
                    # 통계 수집
                    if self.enable_stats:
                        wait_time = now - start_time
                        self.wait_times.append(wait_time)
                        self.calls_per_second[int(now)] += 1
                        self.total_calls += 1
                        
                        if wait_time > 0.01:  # 10ms 이상 대기한 경우
                            logger.info(f"Rate limit 대기: {wait_time:.3f}초")
                    
                    # 최소 간격 보장
                    if self.enable_min_interval and self.min_interval > 0:
                        time.sleep(self.min_interval)
                    
                    return True
                
                # 대기 시간 계산
                wait_time = self._calculate_wait_time(now, has_token)
                
                if wait_time > 0:
                    # 최대 100ms씩 대기 (응답성 향상)
                    actual_wait = min(wait_time, 0.1)
                    logger.debug(
                        f"Rate limit 대기 중: {actual_wait:.3f}초 "
                        f"(tokens={self.tokens:.1f}, window={len(self.call_timestamps)})"
                    )
                    time.sleep(actual_wait)
    
    def _calculate_wait_time(self, now: float, has_token: bool) -> float:
        """필요한 대기 시간 계산"""
        if not has_token:
            # 토큰이 부족한 경우: 다음 토큰까지 대기
            tokens_needed = 1 - self.tokens
            wait_for_token = tokens_needed / self.refill_rate
        else:
            wait_for_token = 0
        
        if self.call_timestamps and len(self.call_timestamps) >= self.max_calls:
            # 윈도우가 가득 찬 경우: 가장 오래된 타임스탬프 만료까지 대기
            wait_for_window = self.per_seconds - (now - self.call_timestamps[0])
        else:
            wait_for_window = 0
        
        # 두 조건 중 더 긴 시간 대기
        return max(wait_for_token, wait_for_window, 0)
    
    def record_error(self) -> None:
        """에러 발생 기록"""
        if self.enable_stats:
            with self.lock:
                self.error_count += 1
    
    def get_stats(self) -> Dict[str, Any]:
        """통계 정보 반환"""
        if not self.enable_stats:
            return {"error": "통계 수집이 비활성화되어 있습니다."}
        
        with self.lock:
            if not self.calls_per_second:
                return {
                    "total_calls": 0,
                    "error_count": 0,
                    "message": "아직 호출 데이터가 없습니다."
                }
            
            max_calls_per_sec = max(self.calls_per_second.values())
            avg_wait_time = sum(self.wait_times) / len(self.wait_times) if self.wait_times else 0
            
            return {
                "total_calls": self.total_calls,
                "error_count": self.error_count,
                "error_rate": self.error_count / self.total_calls if self.total_calls > 0 else 0,
                "max_calls_per_second": max_calls_per_sec,
                "avg_wait_time": avg_wait_time,
                "current_tokens": self.tokens,
                "current_window_size": len(self.call_timestamps),
                "config": {
                    "nominal_max_calls": self.nominal_max_calls,
                    "effective_max_calls": self.max_calls,
                    "safety_margin": self.safety_margin,
                    "min_interval": self.min_interval
                }
            }
    
    def print_stats(self) -> None:
        """통계 정보 출력"""
        stats = self.get_stats()
        
        if "error" in stats:
            print(stats["error"])
            return
        
        print("\n" + "=" * 50)
        print("Enhanced Rate Limiter 통계")
        print("=" * 50)
        print(f"총 호출 수: {stats['total_calls']}")
        print(f"에러 수: {stats['error_count']} ({stats['error_rate']:.1%})")
        print(f"최대 초당 호출: {stats['max_calls_per_second']} / {self.max_calls}")
        print(f"평균 대기 시간: {stats['avg_wait_time']:.3f}초")
        print(f"현재 토큰: {stats['current_tokens']:.1f} / {self.max_calls}")
        print(f"현재 윈도우 크기: {stats['current_window_size']}")
        print("\n설정값:")
        for key, value in stats['config'].items():
            print(f"  {key}: {value}")
        print("=" * 50 + "\n")
    
    def save_stats(self, filepath: Optional[str] = None, include_timestamp: bool = True) -> str:
        """통계를 JSON 파일로 저장
        
        Args:
            filepath: 저장할 파일 경로. None이면 기본 경로 사용
            include_timestamp: 파일명에 타임스탬프 포함 여부
            
        Returns:
            str: 저장된 파일 경로
        """
        stats = self.get_stats()
        
        if "error" in stats:
            logger.warning(f"통계 저장 실패: {stats['error']}")
            return ""
        
        # 타임스탬프 추가
        stats['timestamp'] = datetime.datetime.now().isoformat()
        stats['timestamp_epoch'] = time.time()
        
        # 파일 경로 결정
        if filepath is None:
            stats_dir = Path("logs/rate_limiter_stats")
            stats_dir.mkdir(parents=True, exist_ok=True)
            
            if include_timestamp:
                timestamp = datetime.datetime.now().strftime("%Y%m%d_%H%M%S")
                filename = f"rate_limiter_stats_{timestamp}.json"
            else:
                filename = "rate_limiter_stats_latest.json"
            
            filepath = stats_dir / filename
        else:
            filepath = Path(filepath)
            filepath.parent.mkdir(parents=True, exist_ok=True)
        
        # JSON 파일로 저장
        try:
            with open(filepath, 'w', encoding='utf-8') as f:
                json.dump(stats, f, indent=2, ensure_ascii=False)
            
            logger.info(f"통계 저장 완료: {filepath}")
            return str(filepath)
            
        except Exception as e:
            logger.error(f"통계 저장 실패: {e}")
            return ""
    
    def enable_auto_save(self, interval_seconds: int = 300, filepath: Optional[str] = None):
        """주기적으로 통계를 자동 저장
        
        Args:
            interval_seconds: 저장 주기 (초), 기본값 5분
            filepath: 저장할 파일 경로
        """
        def save_periodically():
            while getattr(self, '_auto_save_enabled', False):
                time.sleep(interval_seconds)
                if getattr(self, '_auto_save_enabled', False):
                    self.save_stats(filepath, include_timestamp=False)
        
        # 이미 실행 중이면 중지
        self.disable_auto_save()
        
        # 새로운 스레드 시작
        self._auto_save_enabled = True
        self._auto_save_thread = threading.Thread(
            target=save_periodically,
            daemon=True,
            name="RateLimiterAutoSave"
        )
        self._auto_save_thread.start()
        
        logger.info(f"통계 자동 저장 활성화: {interval_seconds}초마다 저장")
    
    def disable_auto_save(self):
        """자동 저장 비활성화"""
        if hasattr(self, '_auto_save_enabled'):
            self._auto_save_enabled = False
            if hasattr(self, '_auto_save_thread'):
                # 스레드가 종료될 때까지 대기 (최대 1초)
                self._auto_save_thread.join(timeout=1.0)
            logger.info("통계 자동 저장 비활성화")
    
    def reset(self) -> None:
        """Rate limiter 상태 초기화"""
        with self.lock:
            self.call_timestamps.clear()
            self.tokens = float(self.max_calls)
            self.last_update = time.time()
            
            if self.enable_stats:
                self.calls_per_second.clear()
                self.wait_times.clear()
                self.error_count = 0
                self.total_calls = 0
            
            logger.info("Rate limiter 초기화 완료")
    
    def __enter__(self):
        """Context manager 지원"""
        self.acquire()
        return self
    
    def __exit__(self, exc_type, exc_val, exc_tb):
        """Context manager 종료"""
        if exc_type is not None and self.enable_stats:
            self.record_error()
        return False


# 테스트 코드
if __name__ == "__main__":
    # 로깅 설정
    logging.basicConfig(
        level=logging.INFO,
        format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
    )
    
    # Enhanced RateLimiter 테스트
    limiter = EnhancedRateLimiter(max_calls=10, per_seconds=1.0)
    
    print("Enhanced RateLimiter 테스트 시작...")
    
    # 20개 요청 테스트
    for i in range(20):
        start = time.time()
        if limiter.acquire(timeout=5.0):
            elapsed = time.time() - start
            print(f"요청 {i+1}: 성공 (대기: {elapsed:.3f}초)")
        else:
            print(f"요청 {i+1}: 타임아웃")
    
    # 통계 출력
    limiter.print_stats() 