"""
시장 시간대 관리 유틸리티

각 시장의 거래 시간을 관리하고 현재 시장 상태를 판별합니다.

휴일 처리를 위해 holidays 패키지를 사용할 수 있지만,
선택적 의존성으로 처리합니다.
"""

from datetime import datetime, time, timezone
from typing import Literal, Dict, Optional
from zoneinfo import ZoneInfo
import logging

# holidays 모듈을 선택적으로 import
try:
    import holidays
    HOLIDAYS_AVAILABLE = True
except ImportError:
    holidays = None
    HOLIDAYS_AVAILABLE = False
    logging.warning("holidays 모듈이 설치되지 않았습니다. 휴일 검사가 제한됩니다.")

logger = logging.getLogger(__name__)

# 시장 상태 타입
MarketStatus = Literal['regular', 'after_hours', 'weekend']

# 시장별 시간대 정보
MARKET_TIMEZONES = {
    'KR': ZoneInfo('Asia/Seoul'),
    'US': ZoneInfo('America/New_York'),
    'JP': ZoneInfo('Asia/Tokyo'),
    'HK': ZoneInfo('Asia/Hong_Kong'),
    'CN': ZoneInfo('Asia/Shanghai'),
    'VN': ZoneInfo('Asia/Ho_Chi_Minh'),
}

# 시장별 거래 시간 (현지 시간 기준)
MARKET_HOURS = {
    'KR': {
        'pre_market': (time(8, 0), time(9, 0)),      # 장전 시간외
        'regular': (time(9, 0), time(15, 30)),       # 정규장
        'after_hours': (time(15, 30), time(18, 0)),  # 장후 시간외
    },
    'US': {
        'pre_market': (time(4, 0), time(9, 30)),     # 프리마켓
        'regular': (time(9, 30), time(16, 0)),       # 정규장
        'after_hours': (time(16, 0), time(20, 0)),   # 애프터마켓
    },
    'JP': {
        'pre_market': (time(8, 0), time(9, 0)),
        'regular': (time(9, 0), time(15, 0)),        # 점심시간 제외
        'after_hours': (time(15, 0), time(16, 30)),
    },
    'HK': {
        'pre_market': (time(9, 0), time(9, 30)),
        'regular': (time(9, 30), time(16, 0)),       # 점심시간 제외
        'after_hours': (time(16, 0), time(17, 0)),
    },
}


def get_market_status(market: str = 'KR') -> MarketStatus:
    """
    현재 시장 상태 판별
    
    Args:
        market: 시장 코드 (KR, US, JP, HK 등)
        
    Returns:
        'regular' (장중), 'after_hours' (장외), 'weekend' (주말/공휴일)
    """
    # 시장 시간대 가져오기
    tz = MARKET_TIMEZONES.get(market, ZoneInfo('Asia/Seoul'))
    now = datetime.now(tz)
    
    # 주말 확인
    if now.weekday() >= 5:  # 토요일(5) 또는 일요일(6)
        return 'weekend'
    
    # 공휴일 확인
    if market == 'KR':
        if HOLIDAYS_AVAILABLE:
            kr_holidays = holidays.KR()
            if now.date() in kr_holidays:
                return 'weekend'
    elif market == 'US':
        if HOLIDAYS_AVAILABLE:
            us_holidays = holidays.US()
            if now.date() in us_holidays:
                return 'weekend'
    
    # 시장 시간 확인
    current_time = now.time()
    hours = MARKET_HOURS.get(market, MARKET_HOURS['KR'])
    
    # 정규장 시간
    if hours['regular'][0] <= current_time <= hours['regular'][1]:
        # 일본/홍콩 점심시간 체크
        if market == 'JP' and time(11, 30) <= current_time <= time(12, 30):
            return 'after_hours'
        if market == 'HK' and time(12, 0) <= current_time <= time(13, 0):
            return 'after_hours'
        return 'regular'
    
    # 시간외 거래
    pre_start, pre_end = hours.get('pre_market', (time(0, 0), time(0, 0)))
    after_start, after_end = hours.get('after_hours', (time(0, 0), time(0, 0)))
    
    if (pre_start <= current_time <= pre_end) or (after_start <= current_time <= after_end):
        return 'after_hours'
    
    # 그 외 시간
    return 'weekend'


def get_dynamic_ttl(method_name: str, market: str = 'KR') -> int:
    """
    시장 상태에 따른 동적 TTL 계산
    
    Args:
        method_name: API 메서드명
        market: 시장 코드
        
    Returns:
        조정된 TTL (초)
    """
    from .ttl_cache import CACHE_TTL_CONFIG
    
    # 기본 TTL 가져오기
    base_ttl = CACHE_TTL_CONFIG.get(method_name, 300)
    
    # 시장 상태 확인
    market_status = get_market_status(market)
    
    # 상태별 TTL 조정
    if market_status == 'regular':
        return base_ttl
    elif market_status == 'after_hours':
        return base_ttl * 3  # 장외 시간은 3배
    elif market_status == 'weekend':
        return base_ttl * 10  # 주말/공휴일은 10배
    else:
        return base_ttl


def get_all_market_status() -> Dict[str, MarketStatus]:
    """모든 주요 시장의 현재 상태 조회"""
    markets = ['KR', 'US', 'JP', 'HK']
    return {market: get_market_status(market) for market in markets}


def is_market_open(market: str = 'KR') -> bool:
    """시장이 열려있는지 확인 (정규장 시간)"""
    return get_market_status(market) == 'regular' 