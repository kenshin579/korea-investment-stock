'''
한국투자증권 python wrapper
'''
import datetime
import json
import os
import pickle
import random
import threading
import time
import zipfile
import logging
import re
from collections import deque, defaultdict
from concurrent.futures import ThreadPoolExecutor, as_completed
from functools import wraps
from pathlib import Path
from typing import Literal, Optional, List
from zoneinfo import ZoneInfo  # Requires Python 3.9+
from datetime import datetime, timedelta
import atexit

import pandas as pd
import requests
from typing import Dict, Any

# 로거 설정
logger = logging.getLogger(__name__)

# Enhanced RateLimiter import
try:
    from .rate_limiting.enhanced_rate_limiter import EnhancedRateLimiter
    from .rate_limiting.enhanced_backoff_strategy import get_backoff_strategy
    from .rate_limiting.enhanced_retry_decorator import retry_on_rate_limit, retry_on_network_error
    from .error_handling.error_recovery_system import get_error_recovery_system
    from .monitoring.stats_manager import get_stats_manager
    from .caching import TTLCache, cacheable, CACHE_TTL_CONFIG
except ImportError:
    from rate_limiting.enhanced_rate_limiter import EnhancedRateLimiter
    from rate_limiting.enhanced_backoff_strategy import get_backoff_strategy
    from rate_limiting.enhanced_retry_decorator import retry_on_rate_limit, retry_on_network_error
    from error_handling.error_recovery_system import get_error_recovery_system
    from monitoring.stats_manager import get_stats_manager
    from caching import TTLCache, cacheable, CACHE_TTL_CONFIG

# Visualization 모듈
try:
    from .visualization import PlotlyVisualizer, DashboardManager
    VISUALIZATION_AVAILABLE = True
except ImportError:
    VISUALIZATION_AVAILABLE = False
    logger.warning("Visualization 모듈을 사용할 수 없습니다. plotly를 설치하세요.")

EXCHANGE_CODE = {
    "홍콩": "HKS",
    "뉴욕": "NYS",
    "나스닥": "NAS",
    "아멕스": "AMS",
    "도쿄": "TSE",
    "상해": "SHS",
    "심천": "SZS",
    "상해지수": "SHI",
    "심천지수": "SZI",
    "호치민": "HSX",
    "하노이": "HNX"
}

# 해외주식 주문
# 해외주식 잔고
EXCHANGE_CODE2 = {
    "미국전체": "NASD",
    "나스닥": "NAS",
    "뉴욕": "NYSE",
    "아멕스": "AMEX",
    "홍콩": "SEHK",
    "상해": "SHAA",
    "심천": "SZAA",
    "도쿄": "TKSE",
    "하노이": "HASE",
    "호치민": "VNSE"
}

EXCHANGE_CODE3 = {
    "나스닥": "NASD",
    "뉴욕": "NYSE",
    "아멕스": "AMEX",
    "홍콩": "SEHK",
    "상해": "SHAA",
    "심천": "SZAA",
    "도쿄": "TKSE",
    "하노이": "HASE",
    "호치민": "VNSE"
}

EXCHANGE_CODE4 = {
    "나스닥": "NAS",
    "뉴욕": "NYS",
    "아멕스": "AMS",
    "홍콩": "HKS",
    "상해": "SHS",
    "심천": "SZS",
    "도쿄": "TSE",
    "하노이": "HNX",
    "호치민": "HSX",
    "상해지수": "SHI",
    "심천지수": "SZI"
}

CURRENCY_CODE = {
    "나스닥": "USD",
    "뉴욕": "USD",
    "아멕스": "USD",
    "홍콩": "HKD",
    "상해": "CNY",
    "심천": "CNY",
    "도쿄": "JPY",
    "하노이": "VND",
    "호치민": "VND"
}

MARKET_TYPE_MAP = {
    "KR": ["300"],  # "301", "302"
    "KRX": ["300"],  # "301", "302"
    "NASDAQ": ["512"],
    "NYSE": ["513"],
    "AMEX": ["529"],
    "US": ["512", "513", "529"],
    "TYO": ["515"],
    "JP": ["515"],
    "HKEX": ["501"],
    "HK": ["501", "543", "558"],
    "HNX": ["507"],
    "HSX": ["508"],
    "VN": ["507", "508"],
    "SSE": ["551"],
    "SZSE": ["552"],
    "CN": ["551", "552"]
}

MARKET_TYPE = Literal[
    "KRX",
    "NASDAQ",
    "NYSE",
    "AMEX",
    "TYO",
    "HKEX",
    "HNX",
    "HSX",
    "SSE",
    "SZSE",
]

EXCHANGE_TYPE = Literal[
    "NAS",
    "NYS",
    "AMS"
]

MARKET_CODE_MAP: dict[str, MARKET_TYPE] = {
    "300": "KRX",
    "301": "KRX",
    "302": "KRX",
    "512": "NASDAQ",
    "513": "NYSE",
    "529": "AMEX",
    "515": "TYO",
    "501": "HKEX",
    "543": "HKEX",
    "558": "HKEX",
    "507": "HNX",
    "508": "HSX",
    "551": "SSE",
    "552": "SZSE",
}

EXCHANGE_CODE_MAP: dict[str, EXCHANGE_TYPE] = {
    "NASDAQ": "NAS",
    "NYSE": "NYS",
    "AMEX": "AMS"
}

API_RETURN_CODE = {
    "SUCCESS": "0",  # 조회되었습니다
    "EXPIRED_TOKEN": "1",  # 기간이 만료된 token 입니다
    "NO_DATA": "7",  # 조회할 자료가 없습니다
    "RATE_LIMIT_EXCEEDED": "EGW00201",  # Rate limit 초과
}


# Note: retry_on_rate_limit decorator는 enhanced_retry_decorator 모듈에서 import됨


class KoreaInvestment:
    '''
    한국투자증권 REST API
    '''

    def __init__(self, api_key: str, api_secret: str, acc_no: str,
                 mock: bool = False, cache_enabled: bool = True, 
                 cache_config: Optional[dict] = None):
        """생성자
        Args:
            api_key (str): 발급받은 API key
            api_secret (str): 발급받은 API secret
            acc_no (str): 계좌번호 체계의 앞 8자리-뒤 2자리
            exchange (str): "서울", "나스닥", "뉴욕", "아멕스", "홍콩", "상해", "심천", # todo: exchange는 제거 예정
                            "도쿄", "하노이", "호치민"
            mock (bool): True (mock trading), False (real trading)
            cache_enabled (bool): True if cache is enabled, False otherwise
            cache_config (dict, optional): Configuration for the cache
        """
        self.mock = mock
        self.set_base_url(mock)
        self.api_key = api_key
        self.api_secret = api_secret

        # account number
        self.acc_no = acc_no
        self.acc_no_prefix = acc_no.split('-')[0]
        self.acc_no_postfix = acc_no.split('-')[1]
        
        # Enhanced RateLimiter 설정
        self.rate_limiter = EnhancedRateLimiter(
            max_calls=15,  # 기본값 20에서 15로 감소
            per_seconds=1.0,
            safety_margin=0.8,  # 실제로는 12회/초
            enable_min_interval=True,  # 최소 간격 보장
            enable_stats=True  # 통계 수집 활성화
        )
        
        # ThreadPoolExecutor 개선
        # 동시 실행 제한을 위한 세마포어 (최대 3개만 동시 실행)
        self.concurrent_limit = threading.Semaphore(3)
        # 워커 수 감소 (8 -> 3)
        self.executor = ThreadPoolExecutor(max_workers=3)
        # 프로그램 종료 시 자동 정리
        atexit.register(self.shutdown)

        # access token
        self.token_file = Path("~/.cache/mojito2/token.dat").expanduser()
        self.access_token = None
        if self.check_access_token():
            self.load_access_token()
        else:
            self.issue_access_token()

        # Cache configuration
        self._cache_enabled = cache_enabled
        if cache_enabled:
            # 캐시 설정 병합
            default_cache_config = {
                'default_ttl': 300,  # 5분
                'max_size': 10000,
                'ttl_config': CACHE_TTL_CONFIG
            }
            if cache_config:
                default_cache_config.update(cache_config)
            
            # TTLCache 인스턴스 생성
            self._cache = TTLCache(
                default_ttl=default_cache_config['default_ttl'],
                max_size=default_cache_config['max_size']
            )
            logger.info(f"TTL 캐시 활성화 (기본 TTL: {default_cache_config['default_ttl']}초, "
                       f"최대 크기: {default_cache_config['max_size']})")
        else:
            self._cache = None
            logger.info("TTL 캐시 비활성화")
        
        # Visualization 초기화
        self.visualizer = None
        self.dashboard_manager = None
        if VISUALIZATION_AVAILABLE:
            try:
                self.visualizer = PlotlyVisualizer()
                self.dashboard_manager = DashboardManager(self.visualizer)
                logger.info("Visualization 모듈 초기화 완료")
            except Exception as e:
                logger.warning(f"Visualization 모듈 초기화 실패: {e}")

    def __enter__(self):
        """컨텍스트 매니저 진입"""
        return self
    
    def __exit__(self, exc_type, exc_val, exc_tb):
        """컨텍스트 매니저 종료 - 리소스 정리"""
        self.shutdown()
        return False  # 예외를 전파

    def __execute_concurrent_requests(self, method, stock_list, 
                                     batch_size: Optional[int] = None,
                                     batch_delay: float = 0.0,
                                     progress_interval: int = 10,
                                     dynamic_batch_controller=None):
        """병렬 요청 실행 (개선된 버전 with 에러 처리 강화 및 배치 처리)
        
        Phase 3.4: ThreadPoolExecutor 에러 처리 통합
        Phase 4.1: 배치 크기 파라미터화
        Phase 4.2: 동적 배치 조정
        
        Args:
            method: 실행할 메서드
            stock_list: (symbol, market) 튜플 리스트
            batch_size: 배치 크기 (None이면 전체를 한 번에 처리)
            batch_delay: 배치 간 대기 시간 (초)
            progress_interval: 진행 상황 출력 간격
            dynamic_batch_controller: DynamicBatchController 인스턴스 (동적 조정용)
        """
        from .rate_limiting.enhanced_retry_decorator import RateLimitError, APIError
        from .rate_limiting.enhanced_backoff_strategy import get_backoff_strategy
        
        futures = {}
        results = []
        
        # Rate Limit 에러 발생 시 전체 작업 중단 플래그
        rate_limit_error_occurred = False
        rate_limit_error = None
        
        def wrapped_method(symbol, market):
            """세마포어로 동시 실행 제한"""
            with self.concurrent_limit:
                return method(symbol, market)
        
        # 배치 처리 설정
        if dynamic_batch_controller:
            # 동적 배치 조정 사용
            current_batch_size, current_batch_delay = dynamic_batch_controller.get_current_parameters()
            batch_size = current_batch_size
            batch_delay = current_batch_delay
            print(f"🎯 동적 배치 조정 모드: 초기 배치 크기={batch_size}, 대기 시간={batch_delay:.1f}s")
        
        if batch_size is None:
            batches = [stock_list]  # 전체를 하나의 배치로
        else:
            # stock_list를 batch_size 크기로 나누기
            batches = [stock_list[i:i + batch_size] for i in range(0, len(stock_list), batch_size)]
            print(f"📦 배치 처리 모드: {len(stock_list)}개 항목을 {len(batches)}개 배치로 처리 (배치 크기: {batch_size})")
        
        # 전체 작업을 재시도 가능하도록 감싸기
        max_retries = 3  # 전체 작업 재시도 횟수
        retry_count = 0
        
        while retry_count < max_retries:
            futures.clear()
            results.clear()
            rate_limit_error_occurred = False
            rate_limit_error = None
            
            # 배치별로 처리
            for batch_idx, batch in enumerate(batches):
                # 동적 배치 조정: 각 배치마다 새로운 파라미터 가져오기
                if dynamic_batch_controller and batch_idx > 0:
                    new_batch_size, new_batch_delay = dynamic_batch_controller.get_current_parameters()
                    if new_batch_size != batch_size or new_batch_delay != batch_delay:
                        batch_size = new_batch_size
                        batch_delay = new_batch_delay
                        print(f"📊 배치 파라미터 업데이트: 크기={batch_size}, 대기={batch_delay:.1f}s")
                        # 새로운 배치 크기로 재구성이 필요한 경우 (다음 루프에서 적용)
                
                if len(batches) > 1:
                    print(f"\n🔄 배치 {batch_idx + 1}/{len(batches)} 처리 중... ({len(batch)}개 항목)")
                
                # 배치 시작 시간 기록
                batch_start_time = time.time()
                
                # 배치 내 순차적 제출로 초기 버스트 방지
                batch_futures = {}
                submit_delay = 0.01  # 각 제출 간 10ms 대기
                
                # 배치 통계 초기화
                batch_stats = {
                    'batch_idx': batch_idx,
                    'batch_size': len(batch),
                    'submit_start': time.time(),
                    'symbols': []
                }
                
                for idx, (symbol, market) in enumerate(batch):
                    # 순차적 제출로 초기 버스트 방지
                    if idx > 0 and submit_delay > 0:
                        time.sleep(submit_delay)
                    
                    future = self.executor.submit(wrapped_method, symbol, market)
                    batch_futures[future] = (symbol, market)
                    futures[future] = (symbol, market)
                    batch_stats['symbols'].append(symbol)
                
                batch_stats['submit_end'] = time.time()
                batch_stats['submit_duration'] = batch_stats['submit_end'] - batch_stats['submit_start']
                
                # 배치 내 작업 완료 대기
                batch_completed = 0
                batch_total = len(batch)
                batch_success_count = 0
                batch_error_count = 0
                
                try:
                    for future in as_completed(batch_futures, timeout=30):  # 30초 타임아웃
                        symbol, market = batch_futures[future]
                        batch_completed += 1
                        
                        try:
                            result = future.result()
                            results.append(result)
                            batch_success_count += 1
                            
                            # 진행 상황 출력
                            if batch_completed % progress_interval == 0 or batch_completed == batch_total:
                                if len(batches) > 1:
                                    print(f"  배치 진행률: {batch_completed}/{batch_total} ({batch_completed/batch_total*100:.1f}%)")
                                else:
                                    total = len(stock_list)
                                    completed = len(results)
                                    print(f"처리 진행률: {completed}/{total} ({completed/total*100:.1f}%)")
                                
                        except Exception as e:
                            error_info = {
                                'rt_cd': '9',  # 에러 코드
                                'msg1': f'Error: {str(e)}',
                                'error': True,
                                'symbol': symbol,
                                'market': market,
                                'error_type': type(e).__name__,
                                'error_details': str(e)
                            }
                            
                            # Rate Limit 에러 감지
                            if (isinstance(e, RateLimitError) or 
                                (hasattr(e, 'response') and isinstance(e.response, dict) and 
                                 e.response.get('rt_cd') == 'EGW00201') or
                                'EGW00201' in str(e)):
                                
                                print(f"⚠️ Rate Limit 에러 감지 - {symbol} ({market})")
                                rate_limit_error_occurred = True
                                rate_limit_error = e
                                
                                # 남은 작업들 취소
                                for future in futures:
                                    if not future.done():
                                        future.cancel()
                                break
                            
                            # 일반 에러 처리
                            print(f"❌ 에러 발생 - {symbol} ({market}): {e}")
                            results.append(error_info)
                            batch_error_count += 1
                            
                            # Rate limit 에러인 경우 기록
                            if hasattr(self.rate_limiter, 'record_error'):
                                self.rate_limiter.record_error()
                    
                    # Rate Limit 에러가 발생한 경우 배치 처리 중단
                    if rate_limit_error_occurred:
                        break
                    
                    # 동적 배치 조정: 배치 결과 기록
                    if dynamic_batch_controller:
                        batch_elapsed_time = time.time() - batch_start_time
                        dynamic_batch_controller.record_batch_result(
                            batch_size=len(batch),
                            success_count=batch_success_count,
                            error_count=batch_error_count,
                            elapsed_time=batch_elapsed_time
                        )
                    
                    # 배치별 결과 통계 수집 및 로깅
                    batch_elapsed_time = time.time() - batch_start_time
                    batch_stats['process_end'] = time.time()
                    batch_stats['total_duration'] = batch_elapsed_time
                    batch_stats['success_count'] = batch_success_count
                    batch_stats['error_count'] = batch_error_count
                    batch_stats['throughput'] = (batch_success_count + batch_error_count) / batch_elapsed_time if batch_elapsed_time > 0 else 0
                    
                    # 배치 처리 결과 로깅
                    if len(batches) > 1:
                        print(f"\n📊 배치 {batch_idx + 1} 통계:")
                        print(f"   - 제출 시간: {batch_stats['submit_duration']:.2f}초 ({len(batch)}개)")
                        print(f"   - 처리 시간: {batch_elapsed_time:.2f}초")
                        print(f"   - 성공/실패: {batch_success_count}/{batch_error_count}")
                        print(f"   - 처리량: {batch_stats['throughput']:.1f} TPS")
                        
                        # 에러가 있으면 에러 타입별로 분석
                        if batch_error_count > 0:
                            error_types = {}
                            for r in results[-len(batch):]:  # 현재 배치의 결과만
                                if r.get('error'):
                                    error_type = r.get('error_type', 'Unknown')
                                    error_types[error_type] = error_types.get(error_type, 0) + 1
                            print(f"   - 에러 타입: {error_types}")
                    
                    # 배치 간 대기 (마지막 배치 제외)
                    if batch_delay > 0 and batch_idx < len(batches) - 1:
                        print(f"⏱️ 다음 배치까지 {batch_delay}초 대기...")
                        time.sleep(batch_delay)
                        
                except TimeoutError:
                    print(f"⏱️ 타임아웃 발생 - 30초 내 완료되지 않은 작업이 있습니다.")
                    # 타임아웃된 작업들 처리
                    for future, (symbol, market) in batch_futures.items():
                        if not future.done():
                            future.cancel()
                            results.append({
                                'rt_cd': '9',
                                'msg1': 'Timeout - operation took too long',
                                'error': True,
                                'symbol': symbol,
                                'market': market,
                                'error_type': 'TimeoutError'
                            })
                    # 타임아웃이 발생하면 전체 처리 중단
                    rate_limit_error_occurred = True
                    break
                
                # Rate Limit 에러가 발생한 경우 배치 루프 종료
                if rate_limit_error_occurred:
                    break
            
            # Rate Limit 에러가 발생한 경우 재시도
            if rate_limit_error_occurred:
                retry_count += 1
                if retry_count < max_retries:
                    # Backoff 전략 사용
                    backoff = get_backoff_strategy()
                    wait_time, reason = backoff.calculate_backoff(retry_count - 1)
                    
                    print(f"\n⏳ Rate Limit 초과로 전체 작업 재시도 중...")
                    print(f"   대기 시간: {wait_time:.2f}초 ({reason})")
                    print(f"   재시도: {retry_count}/{max_retries}")
                    
                    time.sleep(wait_time)
                    continue  # 전체 작업 재시도
                else:
                    # 최대 재시도 횟수 초과
                    print(f"\n❌ 최대 재시도 횟수 초과. 부분 결과 반환.")
                    # 실패한 작업들도 에러 정보로 추가
                    for future, (symbol, market) in futures.items():
                        if not future.done() or future.cancelled():
                            results.append({
                                'rt_cd': 'EGW00201',
                                'msg1': 'Rate limit exceeded - max retries reached',
                                'error': True,
                                'symbol': symbol,
                                'market': market,
                                'error_type': 'RateLimitError'
                            })
            
            # 성공적으로 완료
            break
        
        # 통계 출력
        if hasattr(self.rate_limiter, 'print_stats'):
            self.rate_limiter.print_stats()
        
        # 성공/실패 요약
        success_count = sum(1 for r in results if not r.get('error', False))
        error_count = len(results) - success_count
        
        print(f"\n📊 처리 완료 - 성공: {success_count}, 실패: {error_count}")
        if error_count > 0:
            error_types = {}
            for r in results:
                if r.get('error'):
                    error_type = r.get('error_type', 'Unknown')
                    error_types[error_type] = error_types.get(error_type, 0) + 1
            print(f"   에러 타입별 분포: {error_types}")
        
        # 배치 처리 통계
        if len(batches) > 1:
            print(f"   배치 수: {len(batches)}, 배치 크기: {batch_size}")
        
        # 동적 배치 조정 통계
        if dynamic_batch_controller:
            controller_stats = dynamic_batch_controller.get_stats()
            print(f"\n🎯 동적 배치 조정 통계:")
            print(f"   최종 배치 크기: {controller_stats['current_batch_size']}")
            print(f"   최종 대기 시간: {controller_stats['current_batch_delay']:.1f}s")
            print(f"   파라미터 조정 횟수: {controller_stats['adjustment_count']}")
            print(f"   목표 에러율: {controller_stats['target_error_rate']:.1%}")
            print(f"   실제 에러율: {controller_stats['overall_error_rate']:.1%}")
        
        return results
    
    def __handle_rate_limit_error(self, retry_count: int):
        """Rate limit 에러 처리 (Exponential Backoff)
        
        DEPRECATED: Enhanced Backoff Strategy로 대체됨
        이 메서드는 하위 호환성을 위해 유지되며, 향후 제거될 예정입니다.
        
        Args:
            retry_count: 재시도 횟수 (0부터 시작)
        """
        # Exponential backoff: 1, 2, 4, 8, 16, 32초
        wait_time = min(2 ** retry_count, 32)
        
        # Jitter 추가 (0~10% 랜덤 추가 대기)
        jitter = random.uniform(0, 0.1 * wait_time)
        total_wait = wait_time + jitter
        
        print(f"Rate limit 초과. {total_wait:.2f}초 대기 후 재시도... (시도 {retry_count + 1}/5)")
        time.sleep(total_wait)

    def shutdown(self):
        """리소스 정리 - ThreadPoolExecutor 종료"""
        if hasattr(self, 'executor') and self.executor:
            print("ThreadPoolExecutor 종료 중...")
            self.executor.shutdown(wait=True)
            self.executor = None
            print("ThreadPoolExecutor 종료 완료")
        
        # Rate limiter 통계 최종 출력 및 저장
        if hasattr(self, 'rate_limiter'):
            if hasattr(self.rate_limiter, 'get_stats'):
                stats = self.rate_limiter.get_stats()
                if stats.get('total_calls', 0) > 0:
                    print(f"\n최종 Rate Limiter 통계:")
                    print(f"- 총 호출 수: {stats['total_calls']}")
                    print(f"- 에러 수: {stats['error_count']}")
                    print(f"- 에러율: {stats['error_rate']:.1%}")
            
            # 통계를 파일로 저장
            if hasattr(self.rate_limiter, 'save_stats'):
                filepath = self.rate_limiter.save_stats(include_timestamp=True)
                if filepath:
                    print(f"- 통계 저장됨: {filepath}")
            
            # 자동 저장 비활성화
            if hasattr(self.rate_limiter, 'disable_auto_save'):
                self.rate_limiter.disable_auto_save()
        
        # Backoff 전략 통계 출력
        backoff_strategy = get_backoff_strategy()
        backoff_stats = backoff_strategy.get_stats()
        if backoff_stats['total_attempts'] > 0:
            print(f"\n최종 Backoff 전략 통계:")
            print(f"- Circuit 상태: {backoff_stats['state']}")
            print(f"- 총 시도: {backoff_stats['total_attempts']}")
            print(f"- 총 실패: {backoff_stats['total_failures']}")
            print(f"- 성공률: {backoff_stats['success_rate']:.1%}")
            print(f"- Circuit Open 횟수: {backoff_stats['circuit_opens']}")
            print(f"- 평균 백오프 시간: {backoff_stats['avg_backoff_time']:.2f}초")
        
        # 캐시 통계 출력 (Phase 8.7)
        if self._cache_enabled and self._cache:
            cache_stats = self.get_cache_stats()
            if cache_stats['total_entries'] > 0 or cache_stats['hit_count'] > 0:
                print(f"\n최종 캐시 통계:")
                print(f"- 활성화 여부: {'예' if cache_stats['enabled'] else '아니오'}")
                print(f"- 총 항목 수: {cache_stats['total_entries']}")
                print(f"- 캐시 적중: {cache_stats['hit_count']}")
                print(f"- 캐시 미스: {cache_stats['miss_count']}")
                print(f"- 적중률: {cache_stats['hit_rate']:.1%}")
                print(f"- 메모리 사용량: {cache_stats['memory_usage']:.1f}MB")
                print(f"- 만료된 항목: {cache_stats['expired_count']}")
                print(f"- 제거된 항목: {cache_stats['eviction_count']}")
        
        # 에러 복구 시스템 통계 출력
        recovery_system = get_error_recovery_system()
        error_summary = recovery_system.get_error_summary(hours=24)
        if error_summary['total_errors'] > 0:
            print(f"\n최종 에러 복구 통계 (최근 24시간):")
            print(f"- 총 에러 수: {error_summary['total_errors']}")
            print(f"- 심각도별 분포: {error_summary['by_severity']}")
            print(f"- 복구 성공률: {error_summary['recovery_rate']:.1%}")
            print(f"- 가장 빈번한 에러:")
            for error_info in error_summary['most_common'][:3]:
                print(f"  - {error_info['error']}: {error_info['count']}회")
        
        # 에러 통계 파일로 저장
        recovery_system.save_stats()
        
        # 통합 통계 저장 (Phase 5.1)
        print("\n통합 통계 저장 중...")
        stats_manager = get_stats_manager()
        
        # DynamicBatchController가 있다면 포함
        batch_controller = None
        if hasattr(self, '_dynamic_batch_controller'):
            batch_controller = self._dynamic_batch_controller
        
        # 모든 모듈의 통계 수집
        all_stats = stats_manager.collect_all_stats(
            rate_limiter=self.rate_limiter if hasattr(self, 'rate_limiter') else None,
            backoff_strategy=backoff_strategy,
            error_recovery=recovery_system,
            batch_controller=batch_controller,
            cache=self._cache if self._cache_enabled and self._cache else None
        )
        
        # JSON 형식으로 저장
        json_path = stats_manager.save_stats(all_stats, format='json', include_timestamp=True)
        print(f"- 통합 통계 저장됨 (JSON): {json_path}")
        
        # CSV 형식으로도 저장 (요약 정보)
        csv_path = stats_manager.save_stats(all_stats, format='csv', include_timestamp=True)
        print(f"- 통합 통계 저장됨 (CSV): {csv_path}")
        
        # 압축된 JSON Lines 형식으로 저장 (장기 보관용)
        jsonl_gz_path = stats_manager.save_stats(
            all_stats, 
            format='jsonl', 
            compress=True,
            filename='stats_history',
            include_timestamp=False
        )
        print(f"- 통계 이력 추가됨 (JSONL.GZ): {jsonl_gz_path}")
        
        # 시스템 상태 요약 출력
        summary = all_stats.get('summary', {})
        print(f"\n시스템 최종 상태: {summary.get('system_health', 'UNKNOWN')}")
        print(f"- 전체 API 호출: {summary.get('total_api_calls', 0):,}")
        print(f"- 전체 에러: {summary.get('total_errors', 0):,}")
        print(f"- 전체 에러율: {summary.get('overall_error_rate', 0):.2%}")
        
        # 캐시 정리 (Phase 8.7)
        if self._cache_enabled and self._cache:
            # 백그라운드 스레드 정지
            if hasattr(self._cache, 'stop_cleanup_thread'):
                self._cache.stop_cleanup_thread()
            logger.info("캐시 백그라운드 스레드 정리 완료")

    def set_base_url(self, mock: bool = True):
        """테스트(모의투자) 서버 사용 설정
        Args:
            mock(bool, optional): True: 테스트서버, False: 실서버 Defaults to True.
        """
        if mock:
            self.base_url = "https://openapivts.koreainvestment.com:29443"
        else:
            self.base_url = "https://openapi.koreainvestment.com:9443"

    @retry_on_rate_limit(max_retries=3)  # 토큰 발급은 3회만 재시도
    def issue_access_token(self):
        """OAuth인증/접근토큰발급
        """
        path = "oauth2/tokenP"
        url = f"{self.base_url}/{path}"
        headers = {"content-type": "application/json"}
        data = {
            "grant_type": "client_credentials",
            "appkey": self.api_key,
            "appsecret": self.api_secret
        }

        resp = requests.post(url, headers=headers, json=data)
        resp_data = resp.json()
        self.access_token = f'Bearer {resp_data["access_token"]}'

        # 'expires_in' has no reference time and causes trouble:
        # The server thinks I'm expired but my token.dat looks still valid!
        # Hence, we use 'access_token_token_expired' here.
        # This error is quite big. I've seen 4000 seconds.
        timezone = ZoneInfo('Asia/Seoul')
        dt = datetime.strptime(resp_data['access_token_token_expired'], '%Y-%m-%d %H:%M:%S').replace(
            tzinfo=timezone)
        resp_data['timestamp'] = int(dt.timestamp())
        resp_data['api_key'] = self.api_key
        resp_data['api_secret'] = self.api_secret

        # dump access token
        self.token_file.parent.mkdir(parents=True, exist_ok=True)
        with self.token_file.open("wb") as f:
            pickle.dump(resp_data, f)

    def check_access_token(self) -> bool:
        """check access token

        Returns:
            Bool: True: token is valid, False: token is not valid
        """

        if not self.token_file.exists():
            return False

        with self.token_file.open("rb") as f:
            data = pickle.load(f)

        expire_epoch = data['timestamp']
        now_epoch = int(datetime.now().timestamp())
        status = False

        if (data['api_key'] != self.api_key) or (data['api_secret'] != self.api_secret):
            return False

        good_until = data['timestamp']
        ts_now = int(datetime.now().timestamp())
        return ts_now < good_until

    def load_access_token(self):
        """load access token
        """
        with self.token_file.open("rb") as f:
            data = pickle.load(f)
        self.access_token = f'Bearer {data["access_token"]}'

    def issue_hashkey(self, data: dict):
        """해쉬키 발급
        Args:
            data (dict): POST 요청 데이터
        Returns:
            _type_: _description_
        """
        path = "uapi/hashkey"
        url = f"{self.base_url}/{path}"
        headers = {
            "content-type": "application/json",
            "appKey": self.api_key,
            "appSecret": self.api_secret,
            "User-Agent": "Mozilla/5.0"
        }
        resp = requests.post(url, headers=headers, data=json.dumps(data))
        haskkey = resp.json()["HASH"]
        return haskkey

    def fetch_search_stock_info_list(self, stock_market_list):
        return self.__execute_concurrent_requests_with_cache(self.__fetch_search_stock_info, stock_market_list)

    def fetch_price_list(self, stock_list):
        return self.__execute_concurrent_requests_with_cache(self.__fetch_price, stock_list)

    def fetch_price_list_with_batch(self, stock_list, batch_size=50, batch_delay=1.0, progress_interval=10):
        """가격 목록 조회 (배치 처리 지원)
        
        Args:
            stock_list: (symbol, market) 튜플 리스트
            batch_size: 배치 크기 (기본값: 50)
            batch_delay: 배치 간 대기 시간 (초, 기본값: 1.0)
            progress_interval: 진행 상황 출력 간격 (기본값: 10)
        
        Returns:
            list: 조회 결과 리스트
        """
        return self.__execute_concurrent_requests(
            self.__fetch_price, 
            stock_list,
            batch_size=batch_size,
            batch_delay=batch_delay,
            progress_interval=progress_interval
        )
    
    def fetch_price_list_with_dynamic_batch(self, stock_list, dynamic_batch_controller=None):
        """가격 목록 조회 (동적 배치 조정)
        
        Args:
            stock_list: (symbol, market) 튜플 리스트
            dynamic_batch_controller: DynamicBatchController 인스턴스
                                     (None이면 자동 생성)
        
        Returns:
            list: 조회 결과 리스트
        """
        if dynamic_batch_controller is None:
            from .batch_processing.dynamic_batch_controller import DynamicBatchController
            dynamic_batch_controller = DynamicBatchController(
                initial_batch_size=50,
                initial_batch_delay=1.0,
                target_error_rate=0.01
            )
        
        return self.__execute_concurrent_requests(
            self.__fetch_price,
            stock_list,
            dynamic_batch_controller=dynamic_batch_controller
        )

    def __fetch_price(self, symbol: str, market: str = "KR") -> dict:
        """국내주식시세/주식현재가 시세
           해외주식현재가/해외주식 현재체결가

        Args:
            symbol (str): 종목코드

        Returns:
            dict: _description_
        """

        if market == "KR" or market == "KRX":
            stock_info = self.__fetch_stock_info(symbol, market)
            symbol_type = self.__get_symbol_type(stock_info)
            if symbol_type == "ETF":
                resp_json = self.fetch_etf_domestic_price("J", symbol)
            else:
                resp_json = self.fetch_domestic_price("J", symbol)
        elif market == "US":
            resp_json = self.fetch_oversea_price(symbol)
        else:
            raise ValueError("Unsupported market type")

        return resp_json

    def __get_symbol_type(self, symbol_info):
        symbol_type = symbol_info['output']['prdt_clsf_name']
        if symbol_type == '주권' or symbol_type == '상장REITS' or symbol_type == '사회간접자본투융자회사':
            return 'Stock'
        elif symbol_type == 'ETF':
            return 'ETF'

        return "Unknown"

    @cacheable(
        ttl=300,  # 5분
        key_generator=lambda self, market_code, symbol: f"fetch_etf_domestic_price:{market_code}:{symbol}"
    )
    @retry_on_rate_limit()
    def fetch_etf_domestic_price(self, market_code: str, symbol: str) -> dict:
        """주식현재가시세
        Args:
            market_code (str): 시장 분류코드
            symbol (str): 종목코드
        Returns:
            dict: API 개발 가이드 참조
        """
        path = "uapi/domestic-stock/v1/quotations/inquire-price"
        url = f"{self.base_url}/{path}"
        headers = {
            "content-type": "application/json",
            "authorization": self.access_token,
            "appKey": self.api_key,
            "appSecret": self.api_secret,
            "tr_id": "FHPST02400000"
        }
        params = {
            "fid_cond_mrkt_div_code": market_code,
            "fid_input_iscd": symbol
        }
        resp = requests.get(url, headers=headers, params=params)
        return resp.json()

    @cacheable(
        ttl=300,  # 5분
        key_generator=lambda self, market_code, symbol: f"fetch_domestic_price:{market_code}:{symbol}"
    )
    @retry_on_rate_limit()
    def fetch_domestic_price(self, market_code: str, symbol: str) -> dict:
        """주식현재가시세
        Args:
            market_code (str): 시장 분류코드
            symbol (str): 종목코드
        Returns:
            dict: API 개발 가이드 참조
        """
        path = "uapi/domestic-stock/v1/quotations/inquire-price"
        url = f"{self.base_url}/{path}"
        headers = {
            "content-type": "application/json",
            "authorization": self.access_token,
            "appKey": self.api_key,
            "appSecret": self.api_secret,
            "tr_id": "FHKST01010100"
        }
        params = {
            "fid_cond_mrkt_div_code": market_code,
            "fid_input_iscd": symbol
        }
        resp = requests.get(url, headers=headers, params=params)
        return resp.json()

    @cacheable(
        ttl=259200,  # 3일
        key_generator=lambda self: "fetch_kospi_symbols"
    )
    def fetch_kospi_symbols(self):
        """코스피 종목 코드

        실제 필요한 종목: ST, RT, EF, IF

        ST	주권
        MF	증권투자회사
        RT	부동산투자회사
        SC	선박투자회사
        IF	사회간접자본투융자회사
        DR	주식예탁증서
        EW	ELW
        EF	ETF
        SW	신주인수권증권
        SR	신주인수권증서
        BC	수익증권
        FE	해외ETF
        FS	외국주권


        Returns:
            DataFrame:
        """
        base_dir = os.getcwd()
        file_name = "kospi_code.mst.zip"
        url = "https://new.real.download.dws.co.kr/common/master/" + file_name
        self.download_master_file(base_dir, file_name, url)
        df = self.parse_kospi_master(base_dir)
        return df

    @cacheable(
        ttl=259200,  # 3일
        key_generator=lambda self: "fetch_kosdaq_symbols"
    )
    def fetch_kosdaq_symbols(self):
        """코스닥 종목 코드

        Returns:
            DataFrame:
        """
        base_dir = os.getcwd()
        file_name = "kosdaq_code.mst.zip"
        url = "https://new.real.download.dws.co.kr/common/master/" + file_name
        self.download_master_file(base_dir, file_name, url)
        df = self.parse_kosdaq_master(base_dir)
        return df

    def fetch_symbols(self):
        """fetch symbols from the exchange

        Returns:
            pd.DataFrame: pandas dataframe
        """
        if self.exchange == "서울":  # todo: exchange는 제거 예정
            df = self.fetch_kospi_symbols()
            kospi_df = df[['단축코드', '한글명', '그룹코드']].copy()
            kospi_df['시장'] = '코스피'

            df = self.fetch_kosdaq_symbols()
            kosdaq_df = df[['단축코드', '한글명', '그룹코드']].copy()
            kosdaq_df['시장'] = '코스닥'

            df = pd.concat([kospi_df, kosdaq_df], axis=0)

        return df

    def download_master_file(self, base_dir: str, file_name: str, url: str):
        """download master file

        Args:
            base_dir (str): download directory
            file_name (str: filename
            url (str): url
        """
        os.chdir(base_dir)

        # delete legacy master file
        if os.path.exists(file_name):
            os.remove(file_name)

        # download master file
        resp = requests.get(url)
        with open(file_name, "wb") as f:
            f.write(resp.content)

        # unzip
        kospi_zip = zipfile.ZipFile(file_name)
        kospi_zip.extractall()
        kospi_zip.close()

    def parse_kospi_master(self, base_dir: str):
        """parse kospi master file

        Args:
            base_dir (str): directory where kospi code exists

        Returns:
            _type_: _description_
        """
        file_name = base_dir + "/kospi_code.mst"
        tmp_fil1 = base_dir + "/kospi_code_part1.tmp"
        tmp_fil2 = base_dir + "/kospi_code_part2.tmp"

        wf1 = open(tmp_fil1, mode="w", encoding="cp949")
        wf2 = open(tmp_fil2, mode="w")

        with open(file_name, mode="r", encoding="cp949") as f:
            for row in f:
                rf1 = row[0:len(row) - 228]
                rf1_1 = rf1[0:9].rstrip()
                rf1_2 = rf1[9:21].rstrip()
                rf1_3 = rf1[21:].strip()
                wf1.write(rf1_1 + ',' + rf1_2 + ',' + rf1_3 + '\n')
                rf2 = row[-228:]
                wf2.write(rf2)

        wf1.close()
        wf2.close()

        part1_columns = ['단축코드', '표준코드', '한글명']
        df1 = pd.read_csv(tmp_fil1, header=None, encoding='cp949', names=part1_columns)

        field_specs = [
            2, 1, 4, 4, 4,
            1, 1, 1, 1, 1,
            1, 1, 1, 1, 1,
            1, 1, 1, 1, 1,
            1, 1, 1, 1, 1,
            1, 1, 1, 1, 1,
            1, 9, 5, 5, 1,
            1, 1, 2, 1, 1,
            1, 2, 2, 2, 3,
            1, 3, 12, 12, 8,
            15, 21, 2, 7, 1,
            1, 1, 1, 1, 9,
            9, 9, 5, 9, 8,
            9, 3, 1, 1, 1
        ]

        part2_columns = [
            '그룹코드', '시가총액규모', '지수업종대분류', '지수업종중분류', '지수업종소분류',
            '제조업', '저유동성', '지배구조지수종목', 'KOSPI200섹터업종', 'KOSPI100',
            'KOSPI50', 'KRX', 'ETP', 'ELW발행', 'KRX100',
            'KRX자동차', 'KRX반도체', 'KRX바이오', 'KRX은행', 'SPAC',
            'KRX에너지화학', 'KRX철강', '단기과열', 'KRX미디어통신', 'KRX건설',
            'Non1', 'KRX증권', 'KRX선박', 'KRX섹터_보험', 'KRX섹터_운송',
            'SRI', '기준가', '매매수량단위', '시간외수량단위', '거래정지',
            '정리매매', '관리종목', '시장경고', '경고예고', '불성실공시',
            '우회상장', '락구분', '액면변경', '증자구분', '증거금비율',
            '신용가능', '신용기간', '전일거래량', '액면가', '상장일자',
            '상장주수', '자본금', '결산월', '공모가', '우선주',
            '공매도과열', '이상급등', 'KRX300', 'KOSPI', '매출액',
            '영업이익', '경상이익', '당기순이익', 'ROE', '기준년월',
            '시가총액', '그룹사코드', '회사신용한도초과', '담보대출가능', '대주가능'
        ]

        df2 = pd.read_fwf(tmp_fil2, widths=field_specs, names=part2_columns)
        df = pd.merge(df1, df2, how='outer', left_index=True, right_index=True)

        # clean temporary file and dataframe
        del (df1)
        del (df2)
        os.remove(tmp_fil1)
        os.remove(tmp_fil2)
        return df

    def parse_kosdaq_master(self, base_dir: str):
        """parse kosdaq master file

        Args:
            base_dir (str): directory where kosdaq code exists

        Returns:
            _type_: _description_
        """
        file_name = base_dir + "/kosdaq_code.mst"
        tmp_fil1 = base_dir + "/kosdaq_code_part1.tmp"
        tmp_fil2 = base_dir + "/kosdaq_code_part2.tmp"

        wf1 = open(tmp_fil1, mode="w", encoding="cp949")
        wf2 = open(tmp_fil2, mode="w")
        with open(file_name, mode="r", encoding="cp949") as f:
            for row in f:
                rf1 = row[0:len(row) - 222]
                rf1_1 = rf1[0:9].rstrip()
                rf1_2 = rf1[9:21].rstrip()
                rf1_3 = rf1[21:].strip()
                wf1.write(rf1_1 + ',' + rf1_2 + ',' + rf1_3 + '\n')

                rf2 = row[-222:]
                wf2.write(rf2)

        wf1.close()
        wf2.close()

        part1_columns = ['단축코드', '표준코드', '한글명']
        df1 = pd.read_csv(tmp_fil1, header=None, encoding="cp949", names=part1_columns)

        field_specs = [
            2, 1, 4, 4, 4,  # line 20
            1, 1, 1, 1, 1,  # line 27
            1, 1, 1, 1, 1,  # line 32
            1, 1, 1, 1, 1,  # line 38
            1, 1, 1, 1, 1,  # line 43
            1, 9, 5, 5, 1,  # line 48
            1, 1, 2, 1, 1,  # line 54
            1, 2, 2, 2, 3,  # line 64
            1, 3, 12, 12, 8,  # line 69
            15, 21, 2, 7, 1,  # line 75
            1, 1, 1, 9, 9,  # line 80
            9, 5, 9, 8, 9,  # line 85
            3, 1, 1, 1
        ]

        part2_columns = [
            '그룹코드', '시가총액규모', '지수업종대분류', '지수업종중분류', '지수업종소분류',  # line 20
            '벤처기업', '저유동성', 'KRX', 'ETP', 'KRX100',  # line 27
            'KRX자동차', 'KRX반도체', 'KRX바이오', 'KRX은행', 'SPAC',  # line 32
            'KRX에너지화학', 'KRX철강', '단기과열', 'KRX미디어통신', 'KRX건설',  # line 38
            '투자주의', 'KRX증권', 'KRX선박', 'KRX섹터_보험', 'KRX섹터_운송',  # line 43
            'KOSDAQ150', '기준가', '매매수량단위', '시간외수량단위', '거래정지',  # line 48
            '정리매매', '관리종목', '시장경고', '경고예고', '불성실공시',  # line 54
            '우회상장', '락구분', '액면변경', '증자구분', '증거금비율',  # line 64
            '신용가능', '신용기간', '전일거래량', '액면가', '상장일자',  # line 69
            '상장주수', '자본금', '결산월', '공모가', '우선주',  # line 75
            '공매도과열', '이상급등', 'KRX300', '매출액', '영업이익',  # line 80
            '경상이익', '당기순이익', 'ROE', '기준년월', '시가총액',  # line 85
            '그룹사코드', '회사신용한도초과', '담보대출가능', '대주가능'
        ]

        df2 = pd.read_fwf(tmp_fil2, widths=field_specs, names=part2_columns)
        df = pd.merge(df1, df2, how='outer', left_index=True, right_index=True)

        # clean temporary file and dataframe
        del (df1)
        del (df2)
        os.remove(tmp_fil1)
        os.remove(tmp_fil2)
        return df

    def fetch_price_detail_oversea_list(self, stock_market_list):
        return self.__execute_concurrent_requests_with_cache(self.__fetch_price_detail_oversea, stock_market_list)

    @cacheable(
        ttl=300,  # 5분 
        key_generator=lambda self, symbol, market: f"fetch_price_detail_oversea:{market}:{symbol}"
    )
    @retry_on_rate_limit()
    def __fetch_price_detail_oversea(self, symbol: str, market: str = "KR"):
        """해외주식 현재가상세

        Args:
            symbol (str): symbol
        """
        self.rate_limiter.acquire()

        path = "/uapi/overseas-price/v1/quotations/price-detail"
        url = f"{self.base_url}/{path}"

        headers = {
            "content-type": "application/json",
            "authorization": self.access_token,
            "appKey": self.api_key,
            "appSecret": self.api_secret,
            "tr_id": "HHDFS76200200"
        }

        if market == "KR" or market == "KRX":
            # API 호출해서 실제로 확인은 못해봄, overasea 이라서 안될 것으로 판단해서 조건문 추가함
            raise ValueError("Market cannot be either 'KR' or 'KRX'.")

        for market_code in MARKET_TYPE_MAP[market]:
            print("market_code", market_code)
            market_type = MARKET_CODE_MAP[market_code]
            exchange_code = EXCHANGE_CODE_MAP[market_type]
            params = {
                "AUTH": "",
                "EXCD": exchange_code,
                "SYMB": symbol
            }
            resp = requests.get(url, headers=headers, params=params)
            resp_json = resp.json()
            if resp_json['rt_cd'] != API_RETURN_CODE["SUCCESS"] or resp_json['output']['rsym'] == '':
                continue

            return resp_json

    def fetch_stock_info_list(self, stock_market_list):
        return self.__execute_concurrent_requests_with_cache(self.__fetch_stock_info, stock_market_list)

    @cacheable(
        ttl=18000,  # 5시간
        key_generator=lambda self, symbol, market: f"fetch_stock_info:{market}:{symbol}"
    )
    @retry_on_rate_limit()
    def __fetch_stock_info(self, symbol: str, market: str = "KR"):
        self.rate_limiter.acquire()

        path = "uapi/domestic-stock/v1/quotations/search-info"
        url = f"{self.base_url}/{path}"
        headers = {
            "content-type": "application/json",
            "authorization": self.access_token,
            "appKey": self.api_key,
            "appSecret": self.api_secret,
            "tr_id": "CTPF1604R"
        }

        for market_code in MARKET_TYPE_MAP[market]:
            try:
                params = {
                    "PDNO": symbol,
                    "PRDT_TYPE_CD": market_code
                }
                resp = requests.get(url, headers=headers, params=params)
                resp_json = resp.json()

                if resp_json['rt_cd'] == API_RETURN_CODE['NO_DATA']:
                    continue
                return resp_json

            except Exception as e:
                print(e)
                if resp_json['rt_cd'] != API_RETURN_CODE['SUCCESS']:
                    continue
                raise e

    def fetch_search_stock_info_list(self, stock_market_list):
        return self.__execute_concurrent_requests_with_cache(self.__fetch_search_stock_info, stock_market_list)

    @cacheable(
        ttl=18000,  # 5시간
        key_generator=lambda self, symbol, market: f"fetch_search_stock_info:{market}:{symbol}"
    )
    @retry_on_rate_limit()
    def __fetch_search_stock_info(self, symbol: str, market: str = "KR"):
        """
        국내 주식만 제공하는 API이다
        """

        self.rate_limiter.acquire()

        path = "uapi/domestic-stock/v1/quotations/search-stock-info"
        url = f"{self.base_url}/{path}"
        headers = {
            "content-type": "application/json",
            "authorization": self.access_token,
            "appKey": self.api_key,
            "appSecret": self.api_secret,
            "tr_id": "CTPF1002R"
        }

        if market != "KR" and market != "KRX":
            raise ValueError("Market must be either 'KR' or 'KRX'.")

        for market_ in MARKET_TYPE_MAP[market]:
            try:
                params = {
                    "PDNO": symbol,
                    "PRDT_TYPE_CD": market_
                }
                resp = requests.get(url, headers=headers, params=params)
                resp_json = resp.json()

                if resp_json['rt_cd'] == API_RETURN_CODE['NO_DATA']:
                    continue
                return resp_json

            except Exception as e:
                print(e)
                if resp_json['rt_cd'] != API_RETURN_CODE['SUCCESS']:
                    continue
                raise e

    def __execute_concurrent_requests_with_cache(self, method, stock_list,
                                                  batch_size: Optional[int] = None,
                                                  batch_delay: float = 0.0,
                                                  progress_interval: int = 10,
                                                  use_cache: bool = True):
        """캐시를 활용한 병렬 요청 실행
        
        Phase 8.4: 리스트 메서드 캐시 처리
        
        Args:
            method: 실행할 메서드
            stock_list: (symbol, market) 튜플 리스트
            batch_size: 배치 크기 (None이면 전체를 한 번에 처리)
            batch_delay: 배치 간 대기 시간 (초)
            progress_interval: 진행 상황 출력 간격
            use_cache: 캐시 사용 여부
        
        Returns:
            list: 조회 결과 리스트
        """
        if not use_cache or not hasattr(self, '_cache') or not self._cache_enabled:
            # 캐시가 비활성화되어 있으면 기존 메서드 사용
            return self.__execute_concurrent_requests(
                method, stock_list, batch_size, batch_delay, progress_interval
            )
        
        # 결과를 순서대로 저장하기 위한 딕셔너리
        cached_results = {}
        uncached_items = []
        result_order = []  # 원래 순서 유지용
        
        # 1. 캐시에서 먼저 조회
        cache_hits = 0
        cache_misses = 0
        
        for idx, (symbol, market) in enumerate(stock_list):
            result_order.append((idx, symbol, market))
            
            # 캐시 키 생성
            method_name = getattr(method, '__name__', str(method)).replace('_KoreaInvestment__', '')
            cache_key = f"{method_name}:{market}:{symbol}"
            
            # 캐시에서 조회
            cached_value = self._cache.get(cache_key)
            if cached_value is not None:
                cached_results[idx] = cached_value
                cache_hits += 1
            else:
                uncached_items.append((symbol, market))
                cache_misses += 1
        
        if cache_hits > 0:
            print(f"💾 캐시 적중: {cache_hits}개 ({cache_hits/len(stock_list)*100:.1f}%), "
                  f"캐시 미스: {cache_misses}개 ({cache_misses/len(stock_list)*100:.1f}%)")
        
        # 2. 캐시되지 않은 항목만 API 호출
        api_results = []
        if uncached_items:
            print(f"🔄 API 호출 필요: {len(uncached_items)}개 항목")
            api_results = self.__execute_concurrent_requests(
                method, uncached_items, batch_size, batch_delay, progress_interval
            )
            
            # 3. API 호출 결과를 캐시에 저장
            method_name = getattr(method, '__name__', str(method)).replace('_KoreaInvestment__', '')
            for (symbol, market), result in zip(uncached_items, api_results):
                # 성공한 경우만 캐싱
                if result.get('rt_cd') == '0':
                    cache_key = f"{method_name}:{market}:{symbol}"
                    # 메서드별 TTL은 cacheable 데코레이터에서 처리
                    # 여기서는 기본 TTL 사용
                    self._cache.set(cache_key, result)
        
        # 4. 전체 결과 조합 (원래 순서대로)
        final_results = []
        api_result_idx = 0
        
        for idx, symbol, market in result_order:
            if idx in cached_results:
                # 캐시된 결과 사용
                final_results.append(cached_results[idx])
            else:
                # API 호출 결과 사용
                if api_result_idx < len(api_results):
                    final_results.append(api_results[api_result_idx])
                    api_result_idx += 1
                else:
                    # 예외적인 경우 (API 호출 실패 등)
                    final_results.append({
                        'rt_cd': '9',
                        'error': True,
                        'symbol': symbol,
                        'market': market,
                        'msg1': 'Failed to fetch data'
                    })
        
        # 캐시 통계 업데이트
        if hasattr(self._cache, 'print_stats') and cache_hits > 0:
            print(f"📊 API 호출 절감: {cache_hits}회")
        
        return final_results
    
    # Phase 8.6: 캐시 관리 메서드
    def clear_cache(self, pattern: Optional[str] = None):
        """캐시 삭제
        
        Args:
            pattern: 삭제할 캐시 키 패턴 (None이면 전체 삭제)
                    예: "fetch_domestic_price:J:005930"
        """
        if not self._cache_enabled or not self._cache:
            return
        
        if pattern is None:
            # 전체 캐시 삭제
            self._cache.clear()
            logger.info("전체 캐시 삭제 완료")
        else:
            # 패턴에 맞는 캐시 삭제
            deleted_count = self._cache.delete_pattern(pattern)
            logger.info(f"{pattern} 패턴의 캐시 {deleted_count}개 삭제 완료")
    
    def get_cache_stats(self) -> dict:
        """캐시 통계 조회
        
        Returns:
            dict: 캐시 통계 정보
        """
        if not self._cache_enabled or not self._cache:
            return {
                'enabled': False,
                'hit_rate': 0.0,
                'total_entries': 0,
                'memory_usage': 0,
                'expired_count': 0
            }
        
        stats = self._cache.get_stats()
        return {
            'enabled': True,
            'hit_rate': stats.get('hit_rate', 0.0),
            'total_entries': stats.get('size', 0),
            'memory_usage': stats.get('memory_usage_mb', 0),
            'expired_count': stats.get('expired_count', 0),
            'hit_count': stats.get('hit_count', 0),
            'miss_count': stats.get('miss_count', 0),
            'eviction_count': stats.get('eviction_count', 0)
        }
    
    def set_cache_enabled(self, enabled: bool):
        """캐시 기능 on/off
        
        Args:
            enabled: True면 캐시 활성화, False면 비활성화
        """
        self._cache_enabled = enabled
        logger.info(f"캐시 {'활성화' if enabled else '비활성화'}")
    
    def preload_cache(self, symbols: List[str], market: str = "KR"):
        """자주 사용하는 종목 미리 캐싱
        
        Args:
            symbols: 종목 코드 리스트
            market: 시장 코드 (기본값: "KR")
        """
        if not self._cache_enabled or not self._cache:
            logger.warning("캐시가 비활성화되어 있어 preload를 수행할 수 없습니다")
            return
        
        print(f"🔄 {len(symbols)}개 종목 캐시 사전 로드 시작...")
        
        # 종목 정보 로드
        stock_info_list = [(symbol, market) for symbol in symbols]
        self.fetch_stock_info_list(stock_info_list)
        
        # 현재가 정보 로드
        price_list = [(symbol, market) for symbol in symbols]
        self.fetch_price_list(price_list)
        
        print(f"✅ {len(symbols)}개 종목 캐시 사전 로드 완료")
        
        # 캐시 통계 출력
        stats = self.get_cache_stats()
        print(f"📊 캐시 상태: {stats['total_entries']}개 항목, "
              f"메모리 사용량: {stats['memory_usage']:.1f}MB")
    
    # Visualization 메서드들
    def create_monitoring_dashboard(self, 
                                  stats_dir: str = "logs/integrated_stats",
                                  update_interval: int = 5000) -> Optional[Any]:
        """모니터링 대시보드 생성
        
        Args:
            stats_dir: 통계 파일 디렉토리
            update_interval: 업데이트 간격 (밀리초)
            
        Returns:
            대시보드 Figure 객체 또는 None
        """
        if not self.dashboard_manager:
            logger.error("Visualization 모듈이 초기화되지 않았습니다.")
            return None
        
        try:
            # 데이터 로드
            self.visualizer.stats_dir = Path(stats_dir)
            self.visualizer.load_history_data()
            self.visualizer.load_latest_stats()
            
            # 대시보드 생성
            dashboard = self.dashboard_manager.create_realtime_dashboard(update_interval)
            
            if dashboard:
                logger.info("모니터링 대시보드 생성 완료")
            
            return dashboard
            
        except Exception as e:
            logger.error(f"대시보드 생성 실패: {e}")
            return None
    
    def save_monitoring_dashboard(self, 
                                filename: str = "api_monitoring_dashboard.html") -> bool:
        """모니터링 대시보드를 파일로 저장
        
        Args:
            filename: 저장할 파일명
            
        Returns:
            성공 여부
        """
        if not self.dashboard_manager:
            logger.error("Visualization 모듈이 초기화되지 않았습니다.")
            return False
        
        try:
            path = self.dashboard_manager.save_dashboard(filename)
            return bool(path)
        except Exception as e:
            logger.error(f"대시보드 저장 실패: {e}")
            return False
    
    def create_stats_report(self, save_as: str = "monitoring_report") -> Dict[str, str]:
        """통계 리포트 생성
        
        Args:
            save_as: 저장할 파일명 (확장자 제외)
            
        Returns:
            생성된 파일 경로들
        """
        if not self.dashboard_manager:
            logger.error("Visualization 모듈이 초기화되지 않았습니다.")
            return {}
        
        try:
            paths = self.dashboard_manager.create_report(save_as)
            logger.info(f"통계 리포트 생성 완료: {len(paths)}개 파일")
            return paths
        except Exception as e:
            logger.error(f"리포트 생성 실패: {e}")
            return {}
    
    def get_system_health_chart(self) -> Optional[Any]:
        """시스템 헬스 차트 생성
        
        Returns:
            헬스 인디케이터 Figure 또는 None
        """
        if not self.visualizer:
            logger.error("Visualization 모듈이 초기화되지 않았습니다.")
            return None
        
        try:
            # 최신 통계 로드
            if not self.visualizer.latest_stats:
                self.visualizer.load_latest_stats()
            
            # 헬스 차트 생성
            chart = self.visualizer.create_system_health_indicator()
            return chart
        except Exception as e:
            logger.error(f"헬스 차트 생성 실패: {e}")
            return None
    
    def get_api_usage_chart(self, hours: int = 24) -> Optional[Any]:
        """API 사용량 차트 생성
        
        Args:
            hours: 표시할 시간 범위
            
        Returns:
            API 사용량 차트 Figure 또는 None
        """
        if not self.visualizer:
            logger.error("Visualization 모듈이 초기화되지 않았습니다.")
            return None
        
        try:
            # 히스토리 데이터 로드
            if not self.visualizer.history_data:
                self.visualizer.load_history_data()
            
            # 데이터프레임 생성
            df = self.visualizer.prepare_dataframe()
            
            # 시간 필터링
            if not df.empty and 'timestamp' in df.columns:
                from datetime import datetime, timedelta
                cutoff_time = datetime.now() - timedelta(hours=hours)
                df = df[df['timestamp'] >= cutoff_time]
            
            # API 호출 차트 생성
            chart = self.visualizer.create_api_calls_chart(df)
            return chart
        except Exception as e:
            logger.error(f"API 사용량 차트 생성 실패: {e}")
            return None
    
    def show_monitoring_dashboard(self):
        """모니터링 대시보드 표시 (브라우저에서 열기)"""
        if not self.dashboard_manager:
            logger.error("Visualization 모듈이 초기화되지 않았습니다.")
            return
        
        try:
            # 대시보드가 없으면 생성
            if not self.dashboard_manager.dashboard:
                self.create_monitoring_dashboard()
            
            # 대시보드 표시
            self.dashboard_manager.show_dashboard()
        except Exception as e:
            logger.error(f"대시보드 표시 실패: {e}")

    # IPO 관련 헬퍼 함수들
    def _validate_date_format(self, date_str: str) -> bool:
        """날짜 형식 검증 (YYYYMMDD)"""
        if len(date_str) != 8:
            return False
        try:
            datetime.strptime(date_str, "%Y%m%d")
            return True
        except ValueError:
            return False

    def _validate_date_range(self, from_date: str, to_date: str) -> bool:
        """날짜 범위 유효성 검증"""
        try:
            start = datetime.strptime(from_date, "%Y%m%d")
            end = datetime.strptime(to_date, "%Y%m%d")
            return start <= end
        except ValueError:
            return False

    @staticmethod
    def parse_ipo_date_range(date_range_str: str) -> tuple:
        """청약기간 문자열 파싱
        
        Args:
            date_range_str: "2024.01.15~2024.01.16" 형식의 문자열
            
        Returns:
            tuple: (시작일 datetime, 종료일 datetime) 또는 (None, None)
        """
        if not date_range_str:
            return (None, None)
        
        # "2024.01.15~2024.01.16" 형식 파싱
        pattern = r'(\d{4}\.\d{2}\.\d{2})~(\d{4}\.\d{2}\.\d{2})'
        match = re.match(pattern, date_range_str)
        
        if match:
            try:
                start_str = match.group(1).replace('.', '')
                end_str = match.group(2).replace('.', '')
                start_date = datetime.strptime(start_str, "%Y%m%d")
                end_date = datetime.strptime(end_str, "%Y%m%d")
                return (start_date, end_date)
            except ValueError:
                pass
        
        return (None, None)

    @staticmethod
    def format_ipo_date(date_str: str) -> str:
        """날짜 형식 변환 (YYYYMMDD -> YYYY-MM-DD)"""
        if len(date_str) == 8:
            return f"{date_str[:4]}-{date_str[4:6]}-{date_str[6:8]}"
        elif '.' in date_str:
            return date_str.replace('.', '-')
        return date_str

    @staticmethod
    def calculate_ipo_d_day(ipo_date_str: str) -> int:
        """청약일까지 남은 일수 계산"""
        if '~' in ipo_date_str:
            start_date, _ = KoreaInvestment.parse_ipo_date_range(ipo_date_str)
            if start_date:
                today = datetime.now()
                return (start_date - today).days
        return -999

    @staticmethod
    def get_ipo_status(subscr_dt: str) -> str:
        """청약 상태 판단
        
        Returns:
            str: "예정", "진행중", "마감", "알수없음"
        """
        start_date, end_date = KoreaInvestment.parse_ipo_date_range(subscr_dt)
        if not start_date or not end_date:
            return "알수없음"
        
        today = datetime.now()
        if today < start_date:
            return "예정"
        elif start_date <= today <= end_date:
            return "진행중"
        else:
            return "마감"

    @staticmethod
    def format_number(num_str: str) -> str:
        """숫자 문자열에 천단위 콤마 추가"""
        try:
            return f"{int(num_str):,}"
        except (ValueError, TypeError):
            return num_str

    # IPO Schedule API
    @cacheable(
        ttl=3600,  # 1시간
        key_generator=lambda self, from_date=None, to_date=None, symbol="": f"fetch_ipo_schedule:{from_date or 'DEFAULT'}:{to_date or 'DEFAULT'}:{symbol or 'ALL'}"
    )
    @retry_on_rate_limit()
    def fetch_ipo_schedule(self, from_date: str = None, to_date: str = None, symbol: str = "") -> dict:
        """공모주 청약 일정 조회
        
        예탁원정보(공모주청약일정) API를 통해 공모주 정보를 조회합니다.
        한국투자 HTS(eFriend Plus) > [0667] 공모주청약 화면과 동일한 기능입니다.
        
        Args:
            from_date: 조회 시작일 (YYYYMMDD, 기본값: 오늘)
            to_date: 조회 종료일 (YYYYMMDD, 기본값: 30일 후)
            symbol: 종목코드 (선택, 공백시 전체 조회)
            
        Returns:
            dict: 공모주 청약 일정 정보
                {
                    "rt_cd": "0",  # 성공여부
                    "msg_cd": "응답코드",
                    "msg1": "응답메시지",
                    "output1": [
                        {
                            "record_date": "기준일",
                            "sht_cd": "종목코드",
                            "isin_name": "종목명",
                            "fix_subscr_pri": "공모가",
                            "face_value": "액면가",
                            "subscr_dt": "청약기간",  # "2024.01.15~2024.01.16"
                            "pay_dt": "납입일",
                            "refund_dt": "환불일",
                            "list_dt": "상장/등록일",
                            "lead_mgr": "주간사",
                            "pub_bf_cap": "공모전자본금",
                            "pub_af_cap": "공모후자본금",
                            "assign_stk_qty": "당사배정물량"
                        }
                    ]
                }
                
        Raises:
            ValueError: 모의투자 사용시 또는 날짜 형식 오류시
            
        Note:
            - 모의투자는 지원하지 않습니다.
            - 예탁원에서 제공한 자료이므로 정보용으로만 사용하시기 바랍니다.
            - 실제 청약시에는 반드시 공식 공모주 청약 공고문을 확인하세요.
            
        Examples:
            >>> # 전체 공모주 조회 (오늘부터 30일)
            >>> ipos = broker.fetch_ipo_schedule()
            
            >>> # 특정 기간 조회
            >>> ipos = broker.fetch_ipo_schedule(
            ...     from_date="20240101",
            ...     to_date="20240131"
            ... )
            
            >>> # 특정 종목 조회
            >>> ipo = broker.fetch_ipo_schedule(symbol="123456")
        """
        # 모의투자 체크
        if self.mock:
            raise ValueError("공모주청약일정 조회는 모의투자를 지원하지 않습니다.")
            
        self.rate_limiter.acquire()
        
        # 날짜 기본값 설정
        if not from_date:
            from_date = datetime.now().strftime("%Y%m%d")
        if not to_date:
            to_date = (datetime.now() + timedelta(days=30)).strftime("%Y%m%d")
        
        # 날짜 유효성 검증
        if not self._validate_date_format(from_date) or not self._validate_date_format(to_date):
            raise ValueError("날짜 형식은 YYYYMMDD 이어야 합니다.")
        
        if not self._validate_date_range(from_date, to_date):
            raise ValueError("시작일은 종료일보다 이전이어야 합니다.")
        
        path = "uapi/domestic-stock/v1/ksdinfo/pub-offer"
        url = f"{self.base_url}/{path}"
        headers = {
            "content-type": "application/json",
            "authorization": self.access_token,
            "appKey": self.api_key,
            "appSecret": self.api_secret,
            "tr_id": "HHKDB669108C0",
            "custtype": "P"  # 개인
        }
        
        params = {
            "SHT_CD": symbol,
            "CTS": "",
            "F_DT": from_date,
            "T_DT": to_date
        }
        
        resp = requests.get(url, headers=headers, params=params)
        resp_json = resp.json()
        
        # 에러 처리
        if resp_json.get('rt_cd') != '0':
            logger.error(f"공모주 조회 실패: {resp_json.get('msg1', 'Unknown error')}")
            return resp_json
        
        return resp_json


# RateLimiter 클래스는 enhanced_rate_limiter.py로 이동됨


if __name__ == "__main__":
    with open("../koreainvestment.key", encoding='utf-8') as key_file:
        lines = key_file.readlines()

    key = lines[0].strip()
    secret = lines[1].strip()
    acc_no = lines[2].strip()

    broker = KoreaInvestment(
        api_key=key,
        api_secret=secret,
        acc_no=acc_no,
        # exchange="나스닥" # todo: exchange는 제거 예정
    )

    balance = broker.fetch_present_balance()
    print(balance)

    # result = broker.fetch_oversea_day_night()
    # pprint.pprint(result)

    # minute1_ohlcv = broker.fetch_today_1m_ohlcv("005930")
    # pprint.pprint(minute1_ohlcv)

    # broker = KoreaInvestment(key, secret, exchange="나스닥")
    # import pprint
    # resp = broker.fetch_price("005930")
    # pprint.pprint(resp)
    #
    # b = broker.fetch_balance("63398082")
    # pprint.pprint(b)
    #
    # resp = broker.create_market_buy_order("63398082", "005930", 10)
    # pprint.pprint(resp)
    #
    # resp = broker.cancel_order("63398082", "91252", "0000117057", "00", 60000, 5, "Y")
    # print(resp)
    #
    # resp = broker.create_limit_buy_order("63398082", "TQQQ", 35, 1)
    # print(resp)



    # import pprint
    # broker = KoreaInvestment(key, secret, exchange="나스닥")
    # resp_ohlcv = broker.fetch_ohlcv("TSLA", '1d', to="")
    # print(len(resp_ohlcv['output2']))
    # pprint.pprint(resp_ohlcv['output2'][0])
    # pprint.pprint(resp_ohlcv['output2'][-1])
