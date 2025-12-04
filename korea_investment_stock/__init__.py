'''
한국투자증권 OpenAPI Python Wrapper

Simple, transparent, and flexible Python wrapper for Korea Investment Securities OpenAPI.
'''

# 메인 클래스
from .korea_investment_stock import KoreaInvestment

# 상수 정의
from .constants import (
    EXCHANGE_CODE_QUOTE,
    EXCHANGE_CODE_ORDER,
    EXCHANGE_CODE_BALANCE,
    EXCHANGE_CODE_DETAIL,
    CURRENCY_CODE,
    MARKET_TYPE_MAP,
    API_RETURN_CODE,
    MARKET_CODE_MAP,
    EXCHANGE_CODE_MAP,
    # 하위 호환성
    EXCHANGE_CODE,
    EXCHANGE_CODE2,
    EXCHANGE_CODE3,
    EXCHANGE_CODE4,
)

# 설정 관리
from .config import Config

# 캐시 기능 (서브패키지)
from .cache import CacheManager, CacheEntry, CachedKoreaInvestment

# 토큰 저장소 (서브패키지)
from .token_storage import TokenStorage, FileTokenStorage, RedisTokenStorage

# Rate Limiting (서브패키지)
from .rate_limit import RateLimiter, RateLimitedKoreaInvestment

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
    "EXCHANGE_CODE_QUOTE",
    "EXCHANGE_CODE_ORDER",
    "EXCHANGE_CODE_BALANCE",
    "EXCHANGE_CODE_DETAIL",
    "CURRENCY_CODE",
    "MARKET_TYPE_MAP",
    "API_RETURN_CODE",
    "MARKET_CODE_MAP",
    "EXCHANGE_CODE_MAP",
    # 하위 호환성
    "EXCHANGE_CODE",
    "EXCHANGE_CODE2",
    "EXCHANGE_CODE3",
    "EXCHANGE_CODE4",

    # 설정 관리
    "Config",

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
]
