#!/usr/bin/env python3
"""
Enhanced Retry Decorator with Error Recovery
Date: 2024-12-28

Phase 3.3: 에러 복구 흐름 구현
- ErrorRecoverySystem과 통합
- 사용자 친화적 에러 메시지
- 자동 복구 메커니즘
"""

import functools
import logging
import time
from typing import Callable, Optional, Any, Dict, Tuple

from .enhanced_backoff_strategy import get_backoff_strategy
from ..error_handling.error_recovery_system import (
    get_error_recovery_system, 
    RecoveryAction,
    ErrorSeverity
)

logger = logging.getLogger(__name__)


class RateLimitError(Exception):
    """Rate Limit 에러"""
    pass


class TokenExpiredError(Exception):
    """토큰 만료 에러"""
    pass


class APIError(Exception):
    """API 에러 기본 클래스"""
    def __init__(self, message: str, code: Optional[str] = None, response: Optional[Dict] = None):
        super().__init__(message)
        self.code = code
        self.response = response


def enhanced_retry(
    max_retries: Optional[int] = None,
    initial_delay: Optional[float] = None,
    enable_circuit_breaker: bool = True,
    error_callback: Optional[Callable] = None
):
    """
    향상된 재시도 데코레이터
    
    Args:
        max_retries: 최대 재시도 횟수 (None이면 에러 패턴에서 결정)
        initial_delay: 초기 대기 시간 (None이면 backoff strategy가 결정)
        enable_circuit_breaker: Circuit Breaker 활성화 여부
        error_callback: 에러 발생 시 호출할 콜백
    """
    def decorator(func: Callable) -> Callable:
        @functools.wraps(func)
        def wrapper(*args, **kwargs) -> Any:
            backoff = get_backoff_strategy()
            recovery = get_error_recovery_system()
            
            retry_count = 0
            last_error = None
            
            while True:
                try:
                    # Circuit Breaker 체크
                    if enable_circuit_breaker:
                        circuit_state = backoff._check_circuit()
                        if circuit_state == "OPEN":
                            error_msg = "Circuit Breaker가 활성화되어 있습니다. 잠시 후 다시 시도해주세요."
                            logger.warning(error_msg)
                            raise APIError(error_msg, code="CIRCUIT_OPEN")
                    
                    # 함수 실행
                    result = func(*args, **kwargs)
                    
                    # 성공 시 backoff 리셋
                    if retry_count > 0:
                        backoff.record_attempt(success=True)
                        logger.info(f"{func.__name__} 재시도 성공 (시도 횟수: {retry_count + 1})")
                    
                    return result
                    
                except Exception as e:
                    # 에러 변환 및 응답 추출
                    error, response = _convert_error(e)
                    
                    # 에러 복구 전략 결정
                    context = {
                        'function': func.__name__,
                        'retry_count': retry_count,
                        'args': str(args)[:100],  # 너무 길면 잘라냄
                        'kwargs': str(kwargs)[:100]
                    }
                    strategy = recovery.handle_error(error, context, response)
                    
                    # 사용자 콜백 호출
                    if error_callback:
                        try:
                            error_callback(error, retry_count, strategy)
                        except Exception as cb_error:
                            logger.error(f"에러 콜백 실행 중 오류: {cb_error}")
                    
                    # 재시도 불가능한 에러인지 확인
                    if not strategy['should_retry']:
                        logger.error(
                            f"{func.__name__} 재시도 불가능한 에러: "
                            f"{strategy['user_message']}"
                        )
                        raise error
                    
                    # 최대 재시도 횟수 확인
                    effective_max_retries = max_retries or strategy['max_retries']
                    if retry_count >= effective_max_retries:
                        logger.error(
                            f"{func.__name__} 최대 재시도 횟수 초과 "
                            f"({effective_max_retries}회): {strategy['user_message']}"
                        )
                        raise error
                    
                    # 대기 시간 계산
                    if RecoveryAction.WAIT in strategy['recovery_actions']:
                        wait_time, reason = backoff.calculate_backoff(retry_count)
                        backoff.record_attempt(success=False)
                    else:
                        wait_time = strategy['wait_time'] or 0
                    
                    # 로그 및 대기
                    logger.warning(
                        f"{func.__name__} 재시도 {retry_count + 1}/{effective_max_retries}: "
                        f"{strategy['user_message']} "
                        f"(대기시간: {wait_time:.2f}초)"
                    )
                    
                    if wait_time > 0:
                        time.sleep(wait_time)
                    
                    retry_count += 1
                    last_error = error
        
        return wrapper
    return decorator


def _convert_error(e: Exception) -> Tuple[Exception, Optional[Dict]]:
    """
    에러를 표준 형식으로 변환
    
    Returns:
        (변환된 에러, 응답 딕셔너리)
    """
    # 이미 APIError인 경우
    if isinstance(e, APIError):
        return e, e.response
    
    # 응답에서 에러 정보 추출 시도
    response = None
    if hasattr(e, 'response'):
        response = e.response
    
    # 에러 코드 확인
    error_code = None
    if response and isinstance(response, dict):
        error_code = response.get('rt_cd') or response.get('msg_cd')
    
    # 에러 타입 결정
    if error_code == "EGW00201":
        return RateLimitError(str(e)), response
    elif error_code == "1":
        return TokenExpiredError(str(e)), response
    elif error_code:
        return APIError(str(e), code=error_code, response=response), response
    
    # 기본 에러
    return e, response


# 특화된 데코레이터들
def retry_on_rate_limit(max_retries: int = 5):
    """Rate Limit 에러에 특화된 재시도 데코레이터"""
    def error_callback(error, retry_count, strategy):
        if isinstance(error, RateLimitError):
            logger.info(f"Rate Limit 에러 감지. 자동 재시도 중... ({retry_count + 1}회)")
    
    return enhanced_retry(
        max_retries=max_retries,
        enable_circuit_breaker=True,
        error_callback=error_callback
    )


def retry_on_network_error(max_retries: int = 3):
    """네트워크 에러에 특화된 재시도 데코레이터"""
    def decorator(func: Callable) -> Callable:
        @functools.wraps(func)
        def wrapper(*args, **kwargs) -> Any:
            retry_count = 0
            
            while retry_count < max_retries:
                try:
                    return func(*args, **kwargs)
                except (ConnectionError, TimeoutError) as e:
                    retry_count += 1
                    if retry_count >= max_retries:
                        raise
                    
                    wait_time = min(2 ** retry_count, 10)  # 최대 10초
                    logger.warning(
                        f"네트워크 에러 발생. {wait_time}초 후 재시도... "
                        f"({retry_count}/{max_retries})"
                    )
                    time.sleep(wait_time)
            
            return func(*args, **kwargs)
        
        return wrapper
    return decorator


def auto_refresh_token(token_refresh_func: Callable):
    """토큰 자동 갱신 데코레이터"""
    def decorator(func: Callable) -> Callable:
        @functools.wraps(func)
        def wrapper(self, *args, **kwargs) -> Any:
            try:
                return func(self, *args, **kwargs)
            except TokenExpiredError:
                logger.info("토큰 만료 감지. 자동 갱신 중...")
                try:
                    token_refresh_func(self)
                    logger.info("토큰 갱신 완료. 요청 재시도...")
                    return func(self, *args, **kwargs)
                except Exception as e:
                    logger.error(f"토큰 갱신 실패: {e}")
                    raise
        
        return wrapper
    return decorator 