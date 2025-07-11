#!/usr/bin/env python3
"""
Error Recovery System
Date: 2024-12-28

Phase 3.3: 에러 복구 흐름 구현
- 에러 타입별 처리 전략
- 자동 복구 메커니즘
- 사용자 알림 시스템
- 에러 통계 및 분석
"""

import os
import json
import logging
import time
from typing import Dict, Any, Optional, List, Callable, Type
from datetime import datetime, timedelta
from collections import defaultdict, deque
from dataclasses import dataclass, asdict
from enum import Enum
import threading

logger = logging.getLogger(__name__)


class ErrorSeverity(Enum):
    """에러 심각도"""
    LOW = "LOW"  # 재시도로 해결 가능
    MEDIUM = "MEDIUM"  # 일시적 문제
    HIGH = "HIGH"  # 영구적 문제
    CRITICAL = "CRITICAL"  # 즉시 중단 필요


class RecoveryAction(Enum):
    """복구 액션"""
    RETRY = "RETRY"  # 재시도
    WAIT = "WAIT"  # 대기 후 재시도
    REFRESH_TOKEN = "REFRESH_TOKEN"  # 토큰 갱신
    NOTIFY_USER = "NOTIFY_USER"  # 사용자 알림
    FAIL_FAST = "FAIL_FAST"  # 즉시 실패
    CIRCUIT_BREAK = "CIRCUIT_BREAK"  # Circuit Breaker 활성화


@dataclass
class ErrorPattern:
    """에러 패턴 정의"""
    error_type: str
    error_code: Optional[str]
    severity: ErrorSeverity
    recovery_actions: List[RecoveryAction]
    max_retries: int
    wait_time: Optional[float]
    message_template: str
    
    def matches(self, error: Exception, response: Optional[Dict] = None) -> bool:
        """에러가 이 패턴과 일치하는지 확인"""
        # 에러 타입 체크
        if not isinstance(error, Exception):
            return self.error_type == type(error).__name__
        
        # 응답 코드 체크 (있는 경우)
        if self.error_code and response:
            return response.get('rt_cd') == self.error_code or \
                   response.get('msg_cd') == self.error_code
        
        # Exception 타입 체크
        return type(error).__name__ == self.error_type


@dataclass
class ErrorEvent:
    """에러 이벤트"""
    timestamp: datetime
    error_type: str
    error_code: Optional[str]
    error_message: str
    severity: ErrorSeverity
    recovery_attempted: bool
    recovery_success: bool
    context: Dict[str, Any]


class ErrorRecoverySystem:
    """에러 복구 시스템"""
    
    def __init__(self):
        self.lock = threading.Lock()
        self.error_patterns = self._init_error_patterns()
        self.error_history = deque(maxlen=1000)  # 최근 1000개 에러
        self.error_stats = defaultdict(lambda: {'count': 0, 'last_seen': None})
        self.recovery_callbacks: Dict[RecoveryAction, List[Callable]] = defaultdict(list)
        
        # 설정
        self.enable_notifications = True
        self.enable_auto_recovery = True
        self.stats_file = os.getenv('ERROR_STATS_FILE', 'error_stats.json')
        
        logger.info("ErrorRecoverySystem 초기화 완료")
    
    def _init_error_patterns(self) -> List[ErrorPattern]:
        """에러 패턴 초기화"""
        return [
            # Rate Limit 에러
            ErrorPattern(
                error_type="RateLimitError",
                error_code="EGW00201",
                severity=ErrorSeverity.MEDIUM,
                recovery_actions=[RecoveryAction.WAIT, RecoveryAction.RETRY],
                max_retries=5,
                wait_time=None,  # Enhanced Backoff가 계산
                message_template="API 호출 한도 초과. 잠시 후 다시 시도해주세요."
            ),
            
            # 토큰 만료
            ErrorPattern(
                error_type="TokenExpiredError",
                error_code="1",
                severity=ErrorSeverity.LOW,
                recovery_actions=[RecoveryAction.REFRESH_TOKEN, RecoveryAction.RETRY],
                max_retries=1,
                wait_time=0,
                message_template="인증 토큰이 만료되었습니다. 자동으로 갱신 중..."
            ),
            
            # 네트워크 에러
            ErrorPattern(
                error_type="ConnectionError",
                error_code=None,
                severity=ErrorSeverity.LOW,
                recovery_actions=[RecoveryAction.WAIT, RecoveryAction.RETRY],
                max_retries=3,
                wait_time=1.0,
                message_template="네트워크 연결 오류. 재시도 중..."
            ),
            
            # 타임아웃
            ErrorPattern(
                error_type="Timeout",
                error_code=None,
                severity=ErrorSeverity.LOW,
                recovery_actions=[RecoveryAction.RETRY],
                max_retries=2,
                wait_time=0.5,
                message_template="요청 시간 초과. 다시 시도 중..."
            ),
            
            # 인증 에러
            ErrorPattern(
                error_type="AuthenticationError",
                error_code=None,
                severity=ErrorSeverity.HIGH,
                recovery_actions=[RecoveryAction.NOTIFY_USER, RecoveryAction.FAIL_FAST],
                max_retries=0,
                wait_time=0,
                message_template="인증 실패. API 키를 확인해주세요."
            ),
            
            # 잘못된 파라미터
            ErrorPattern(
                error_type="InvalidParameterError",
                error_code=None,
                severity=ErrorSeverity.HIGH,
                recovery_actions=[RecoveryAction.NOTIFY_USER, RecoveryAction.FAIL_FAST],
                max_retries=0,
                wait_time=0,
                message_template="잘못된 요청 파라미터입니다."
            ),
            
            # 서버 에러
            ErrorPattern(
                error_type="ServerError",
                error_code=None,
                severity=ErrorSeverity.MEDIUM,
                recovery_actions=[RecoveryAction.WAIT, RecoveryAction.RETRY, RecoveryAction.CIRCUIT_BREAK],
                max_retries=3,
                wait_time=5.0,
                message_template="서버 오류가 발생했습니다. 잠시 후 다시 시도해주세요."
            ),
            
            # API 일반 에러 (500 등)
            ErrorPattern(
                error_type="APIError",
                error_code="500",
                severity=ErrorSeverity.MEDIUM,
                recovery_actions=[RecoveryAction.WAIT, RecoveryAction.RETRY],
                max_retries=5,
                wait_time=2.0,
                message_template="서버 내부 오류. 잠시 후 재시도합니다."
            ),
            
            # 데이터 없음
            ErrorPattern(
                error_type="NoDataError",
                error_code="7",
                severity=ErrorSeverity.LOW,
                recovery_actions=[RecoveryAction.NOTIFY_USER],
                max_retries=0,
                wait_time=0,
                message_template="조회할 데이터가 없습니다."
            )
        ]
    
    def register_callback(self, action: RecoveryAction, callback: Callable) -> None:
        """복구 액션에 대한 콜백 등록"""
        with self.lock:
            self.recovery_callbacks[action].append(callback)
    
    def handle_error(
        self, 
        error: Exception, 
        context: Optional[Dict[str, Any]] = None,
        response: Optional[Dict] = None
    ) -> Dict[str, Any]:
        """
        에러 처리 및 복구 전략 결정
        
        Returns:
            dict: {
                'pattern': ErrorPattern,
                'should_retry': bool,
                'wait_time': float,
                'user_message': str,
                'recovery_actions': List[RecoveryAction]
            }
        """
        with self.lock:
            # 에러 패턴 매칭
            pattern = self._match_error_pattern(error, response)
            if not pattern:
                # 알 수 없는 에러
                pattern = self._get_default_pattern(error)
            
            # 에러 이벤트 기록
            event = ErrorEvent(
                timestamp=datetime.now(),
                error_type=pattern.error_type,
                error_code=pattern.error_code,
                error_message=str(error),
                severity=pattern.severity,
                recovery_attempted=False,
                recovery_success=False,
                context=context or {}
            )
            self.error_history.append(event)
            
            # 통계 업데이트
            self._update_stats(pattern)
            
            # 복구 전략 결정
            strategy = {
                'pattern': pattern,
                'should_retry': RecoveryAction.RETRY in pattern.recovery_actions,
                'wait_time': pattern.wait_time or 0,
                'user_message': pattern.message_template,
                'recovery_actions': pattern.recovery_actions,
                'max_retries': pattern.max_retries
            }
            
            # 복구 액션 실행
            if self.enable_auto_recovery:
                self._execute_recovery_actions(pattern, event)
            
            return strategy
    
    def _match_error_pattern(
        self, 
        error: Exception, 
        response: Optional[Dict] = None
    ) -> Optional[ErrorPattern]:
        """에러에 맞는 패턴 찾기"""
        for pattern in self.error_patterns:
            if pattern.matches(error, response):
                return pattern
        return None
    
    def _get_default_pattern(self, error: Exception) -> ErrorPattern:
        """기본 에러 패턴"""
        return ErrorPattern(
            error_type=type(error).__name__,
            error_code=None,
            severity=ErrorSeverity.MEDIUM,
            recovery_actions=[RecoveryAction.NOTIFY_USER],
            max_retries=0,
            wait_time=0,
            message_template=f"예기치 않은 오류: {str(error)}"
        )
    
    def _update_stats(self, pattern: ErrorPattern) -> None:
        """에러 통계 업데이트"""
        key = f"{pattern.error_type}:{pattern.error_code or 'None'}"
        self.error_stats[key]['count'] += 1
        self.error_stats[key]['last_seen'] = datetime.now().isoformat()
        self.error_stats[key]['severity'] = pattern.severity.value
    
    def _execute_recovery_actions(
        self, 
        pattern: ErrorPattern, 
        event: ErrorEvent
    ) -> None:
        """복구 액션 실행"""
        for action in pattern.recovery_actions:
            callbacks = self.recovery_callbacks.get(action, [])
            for callback in callbacks:
                try:
                    callback(pattern, event)
                except Exception as e:
                    logger.error(f"복구 콜백 실행 중 에러: {e}")
    
    def get_error_summary(self, hours: int = 24) -> Dict[str, Any]:
        """에러 요약 통계"""
        with self.lock:
            cutoff_time = datetime.now() - timedelta(hours=hours)
            recent_errors = [
                e for e in self.error_history 
                if e.timestamp >= cutoff_time
            ]
            
            # 심각도별 분류
            by_severity = defaultdict(int)
            for error in recent_errors:
                by_severity[error.severity.value] += 1
            
            # 타입별 분류
            by_type = defaultdict(int)
            for error in recent_errors:
                by_type[error.error_type] += 1
            
            # 복구 성공률
            recovery_attempts = [e for e in recent_errors if e.recovery_attempted]
            recovery_success = [e for e in recovery_attempts if e.recovery_success]
            recovery_rate = (
                len(recovery_success) / len(recovery_attempts) 
                if recovery_attempts else 0
            )
            
            return {
                'period_hours': hours,
                'total_errors': len(recent_errors),
                'by_severity': dict(by_severity),
                'by_type': dict(by_type),
                'recovery_rate': recovery_rate,
                'most_common': self._get_most_common_errors(5)
            }
    
    def _get_most_common_errors(self, limit: int = 5) -> List[Dict]:
        """가장 빈번한 에러들"""
        sorted_stats = sorted(
            self.error_stats.items(),
            key=lambda x: x[1]['count'],
            reverse=True
        )[:limit]
        
        return [
            {
                'error': key,
                'count': value['count'],
                'last_seen': value['last_seen'],
                'severity': value.get('severity', 'UNKNOWN')
            }
            for key, value in sorted_stats
        ]
    
    def save_stats(self) -> None:
        """통계를 파일로 저장"""
        try:
            with open(self.stats_file, 'w') as f:
                stats = {
                    'summary': self.get_error_summary(),
                    'detailed_stats': dict(self.error_stats),
                    'timestamp': datetime.now().isoformat()
                }
                json.dump(stats, f, indent=2, ensure_ascii=False)
            logger.info(f"에러 통계 저장됨: {self.stats_file}")
        except Exception as e:
            logger.error(f"통계 저장 실패: {e}")
    
    def clear_old_errors(self, days: int = 7) -> int:
        """오래된 에러 기록 삭제"""
        with self.lock:
            cutoff_time = datetime.now() - timedelta(days=days)
            initial_count = len(self.error_history)
            
            # deque는 자동으로 오래된 항목을 제거하므로
            # 여기서는 통계만 정리
            for key in list(self.error_stats.keys()):
                last_seen_str = self.error_stats[key].get('last_seen')
                if last_seen_str:
                    last_seen = datetime.fromisoformat(last_seen_str)
                    if last_seen < cutoff_time:
                        del self.error_stats[key]
            
            return initial_count - len(self.error_history)


# 전역 에러 복구 시스템
_global_error_recovery: Optional[ErrorRecoverySystem] = None


def get_error_recovery_system() -> ErrorRecoverySystem:
    """싱글톤 에러 복구 시스템 반환"""
    global _global_error_recovery
    if _global_error_recovery is None:
        _global_error_recovery = ErrorRecoverySystem()
    return _global_error_recovery


# 편의 함수들
def handle_api_error(
    error: Exception,
    context: Optional[Dict[str, Any]] = None,
    response: Optional[Dict] = None
) -> Dict[str, Any]:
    """API 에러 처리 편의 함수"""
    recovery_system = get_error_recovery_system()
    return recovery_system.handle_error(error, context, response) 