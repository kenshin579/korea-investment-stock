'''
한국투자증권 OpenAPI Python Wrapper

Simple, transparent, and flexible Python wrapper for Korea Investment Securities OpenAPI.
'''

# 메인 클래스
from .korea_investment_stock import KoreaInvestment

# 상수 정의
from .constants import (
    MARKET_TYPE_MAP,
    API_RETURN_CODE,
)

# 설정 관리
from .config import Config
from .config_resolver import ConfigResolver

# 캐시 기능 (서브패키지)
from .cache import CacheManager, CacheEntry, CachedKoreaInvestment

# 토큰 저장소 (서브패키지)
from .token_storage import TokenStorage, FileTokenStorage, RedisTokenStorage

# Rate Limiting (서브패키지)
from .rate_limit import RateLimiter, RateLimitedKoreaInvestment

# 파서 (서브패키지)
from .parsers import parse_kospi_master, parse_kosdaq_master

# IPO 헬퍼 (서브패키지)
from .ipo import (
    validate_date_format,
    validate_date_range,
    parse_ipo_date_range,
    format_ipo_date,
    calculate_ipo_d_day,
    get_ipo_status,
    format_number,
)

# Git tag에서 버전 자동 추출 (setuptools-scm)
try:
    from importlib.metadata import version
    __version__ = version("korea-investment-stock")
except Exception:
    # Fallback for development without git tags
    __version__ = "0.0.0.dev0"

__all__ = [
    # 메인 API
    "KoreaInvestment",

    # 상수 정의
    "MARKET_TYPE_MAP",
    "API_RETURN_CODE",

    # 설정 관리
    "Config",
    "ConfigResolver",

    # 캐시 기능
    "CacheManager",
    "CacheEntry",
    "CachedKoreaInvestment",

    # 토큰 저장소
    "TokenStorage",
    "FileTokenStorage",
    "RedisTokenStorage",

    # Rate Limiting
    "RateLimiter",
    "RateLimitedKoreaInvestment",

    # 파서
    "parse_kospi_master",
    "parse_kosdaq_master",

    # IPO 헬퍼
    "validate_date_format",
    "validate_date_range",
    "parse_ipo_date_range",
    "format_ipo_date",
    "calculate_ipo_d_day",
    "get_ipo_status",
    "format_number",
]
