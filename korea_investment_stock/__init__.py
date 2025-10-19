"""A Python port of Korea-Investment-Stock API

Pure wrapper for Korea Investment Securities OpenAPI
"""

__version__ = "0.6.0"

# Core imports
from .korea_investment_stock import KoreaInvestment, MARKET_CODE_MAP, EXCHANGE_CODE_MAP, API_RETURN_CODE

# Public API
__all__ = [
    'KoreaInvestment',
    'MARKET_CODE_MAP',
    'EXCHANGE_CODE_MAP',
    'API_RETURN_CODE',
]
