#!/usr/bin/env python3
"""
Enhanced Backoff Strategy
Date: 2024-12-28

Phase 3.2: 고급 Exponential Backoff 전략 구현
- Circuit Breaker 패턴
- Adaptive Backoff
- 통계 수집
- 환경 변수 설정
"""

import os
import time
import logging
import random
from typing import Optional, Dict, Any, Tuple
from collections import deque
from datetime import datetime, timedelta
from dataclasses import dataclass
import threading

logger = logging.getLogger(__name__)


@dataclass
class BackoffConfig:
    """백오프 설정"""
    base_delay: float = 1.0  # 기본 대기 시간
    max_delay: float = 32.0  # 최대 대기 시간
    exponential_base: float = 2.0  # 지수 베이스
    jitter_factor: float = 0.1  # Jitter 비율 (0~10%)
    
    # Circuit Breaker 설정
    failure_threshold: int = 10  # 실패 임계치
    success_threshold: int = 3  # 복구 임계치
    timeout: float = 60.0  # Circuit open 시간
    
    # Adaptive 설정
    enable_adaptive: bool = True
    min_success_rate: float = 0.2  # 최소 성공률
    
    @classmethod
    def from_env(cls) -> 'BackoffConfig':
        """환경 변수에서 설정 로드"""
        config = cls()
        
        if base_delay := os.getenv('BACKOFF_BASE_DELAY'):
            config.base_delay = float(base_delay)
        
        if max_delay := os.getenv('BACKOFF_MAX_DELAY'):
            config.max_delay = float(max_delay)
        
        if exp_base := os.getenv('BACKOFF_EXPONENTIAL_BASE'):
            config.exponential_base = float(exp_base)
        
        if jitter := os.getenv('BACKOFF_JITTER_FACTOR'):
            config.jitter_factor = float(jitter)
        
        if failure_threshold := os.getenv('CIRCUIT_FAILURE_THRESHOLD'):
            config.failure_threshold = int(failure_threshold)
        
        if success_threshold := os.getenv('CIRCUIT_SUCCESS_THRESHOLD'):
            config.success_threshold = int(success_threshold)
        
        if timeout := os.getenv('CIRCUIT_TIMEOUT'):
            config.timeout = float(timeout)
        
        return config


class CircuitState:
    """Circuit Breaker 상태"""
    CLOSED = "CLOSED"  # 정상 작동
    OPEN = "OPEN"  # 차단됨
    HALF_OPEN = "HALF_OPEN"  # 테스트 중


class EnhancedBackoffStrategy:
    """고급 백오프 전략"""
    
    def __init__(self, config: Optional[BackoffConfig] = None):
        self.config = config or BackoffConfig.from_env()
        self.lock = threading.Lock()
        
        # Circuit Breaker 상태
        self.state = CircuitState.CLOSED
        self.failure_count = 0
        self.success_count = 0
        self.last_failure_time: Optional[datetime] = None
        self.circuit_open_time: Optional[datetime] = None
        
        # 통계
        self.attempt_history = deque(maxlen=100)  # 최근 100개 시도 기록
        self.backoff_times = []
        self.total_attempts = 0
        self.total_failures = 0
        self.total_circuit_opens = 0
        
        logger.info(f"EnhancedBackoffStrategy 초기화: {self.config}")
    
    def calculate_backoff(self, retry_count: int) -> Tuple[float, str]:
        """
        백오프 시간 계산
        
        Returns:
            tuple: (대기시간, 이유)
        """
        with self.lock:
            # Circuit Breaker 체크
            circuit_status = self._check_circuit()
            if circuit_status == CircuitState.OPEN:
                remaining = self._get_circuit_remaining_time()
                return (remaining, f"Circuit OPEN (남은시간: {remaining:.1f}초)")
            
            # 기본 Exponential Backoff
            base_delay = self.config.base_delay * (
                self.config.exponential_base ** retry_count
            )
            delay = min(base_delay, self.config.max_delay)
            
            # Adaptive Backoff
            if self.config.enable_adaptive:
                success_rate = self._get_success_rate()
                if success_rate < self.config.min_success_rate:
                    # 성공률이 낮으면 더 긴 대기
                    delay *= (1 + (self.config.min_success_rate - success_rate))
            
            # Jitter 추가
            jitter = random.uniform(0, delay * self.config.jitter_factor)
            total_delay = delay + jitter
            
            # 통계 기록
            self.backoff_times.append(total_delay)
            
            reason = f"Exponential Backoff (시도 {retry_count + 1})"
            if circuit_status == CircuitState.HALF_OPEN:
                reason += " [Circuit HALF_OPEN]"
            
            return (total_delay, reason)
    
    def record_attempt(self, success: bool) -> None:
        """시도 결과 기록"""
        with self.lock:
            self.total_attempts += 1
            self.attempt_history.append((datetime.now(), success))
            
            if success:
                self._handle_success()
            else:
                self._handle_failure()
    
    def _handle_success(self) -> None:
        """성공 처리"""
        self.success_count += 1
        self.failure_count = 0  # 실패 카운트 리셋
        
        # Circuit Breaker 상태 업데이트
        if self.state == CircuitState.HALF_OPEN:
            if self.success_count >= self.config.success_threshold:
                self.state = CircuitState.CLOSED
                logger.info("Circuit Breaker: HALF_OPEN → CLOSED")
    
    def _handle_failure(self) -> None:
        """실패 처리"""
        self.failure_count += 1
        self.total_failures += 1
        self.last_failure_time = datetime.now()
        
        # Circuit Breaker 상태 업데이트
        if self.state == CircuitState.CLOSED:
            if self.failure_count >= self.config.failure_threshold:
                self.state = CircuitState.OPEN
                self.circuit_open_time = datetime.now()
                self.total_circuit_opens += 1
                logger.warning(
                    f"Circuit Breaker: CLOSED → OPEN "
                    f"(실패 {self.failure_count}회)"
                )
        elif self.state == CircuitState.HALF_OPEN:
            # 테스트 중 실패하면 다시 OPEN
            self.state = CircuitState.OPEN
            self.circuit_open_time = datetime.now()
            logger.warning("Circuit Breaker: HALF_OPEN → OPEN")
    
    def _check_circuit(self) -> str:
        """Circuit Breaker 상태 확인 및 업데이트"""
        if self.state == CircuitState.OPEN and self.circuit_open_time:
            elapsed = (datetime.now() - self.circuit_open_time).total_seconds()
            if elapsed >= self.config.timeout:
                self.state = CircuitState.HALF_OPEN
                self.success_count = 0
                logger.info("Circuit Breaker: OPEN → HALF_OPEN")
        
        return self.state
    
    def _get_circuit_remaining_time(self) -> float:
        """Circuit이 열릴 때까지 남은 시간"""
        if self.state == CircuitState.OPEN and self.circuit_open_time:
            elapsed = (datetime.now() - self.circuit_open_time).total_seconds()
            remaining = self.config.timeout - elapsed
            return max(0, remaining)
        return 0
    
    def _get_success_rate(self) -> float:
        """최근 성공률 계산"""
        if not self.attempt_history:
            return 1.0
        
        recent_attempts = list(self.attempt_history)
        success_count = sum(1 for _, success in recent_attempts if success)
        return success_count / len(recent_attempts)
    
    def should_retry(self, error_type: str) -> bool:
        """재시도 여부 결정"""
        # Circuit이 OPEN이면 재시도 안함
        if self.state == CircuitState.OPEN:
            return False
        
        # 재시도 불가능한 에러 타입
        non_retryable_errors = {
            "AuthenticationError",
            "InvalidParameterError",
            "PermissionDeniedError",
            "AuthError"  # 테스트용 추가
        }
        
        return error_type not in non_retryable_errors
    
    def get_stats(self) -> Dict[str, Any]:
        """통계 정보 반환"""
        with self.lock:
            success_rate = self._get_success_rate()
            avg_backoff = (
                sum(self.backoff_times) / len(self.backoff_times)
                if self.backoff_times else 0
            )
            
            return {
                "state": self.state,
                "total_attempts": self.total_attempts,
                "total_failures": self.total_failures,
                "success_rate": success_rate,
                "circuit_opens": self.total_circuit_opens,
                "current_failure_count": self.failure_count,
                "avg_backoff_time": avg_backoff,
                "config": {
                    "base_delay": self.config.base_delay,
                    "max_delay": self.config.max_delay,
                    "failure_threshold": self.config.failure_threshold,
                    "success_threshold": self.config.success_threshold,
                    "circuit_timeout": self.config.timeout
                }
            }
    
    def reset(self) -> None:
        """상태 초기화"""
        with self.lock:
            self.state = CircuitState.CLOSED
            self.failure_count = 0
            self.success_count = 0
            self.last_failure_time = None
            self.circuit_open_time = None
            self.attempt_history.clear()
            self.backoff_times.clear()
            logger.info("BackoffStrategy 초기화 완료")


# 전역 백오프 전략 인스턴스
_global_backoff_strategy: Optional[EnhancedBackoffStrategy] = None


def get_backoff_strategy() -> EnhancedBackoffStrategy:
    """싱글톤 백오프 전략 인스턴스 반환"""
    global _global_backoff_strategy
    if _global_backoff_strategy is None:
        _global_backoff_strategy = EnhancedBackoffStrategy()
    return _global_backoff_strategy 