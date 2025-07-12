"""A Python port of Korea-Investment-Stock API"""

__version__ = "0.4.1"

# Core imports
from .korea_investment_stock import KoreaInvestment, MARKET_CODE_MAP, EXCHANGE_CODE_MAP, API_RETURN_CODE

# Rate limiting imports
from .rate_limiting.enhanced_rate_limiter import EnhancedRateLimiter
from .rate_limiting.enhanced_retry_decorator import retry_on_rate_limit, retry_on_network_error
from .rate_limiting.enhanced_backoff_strategy import EnhancedBackoffStrategy, get_backoff_strategy

# Error handling imports
from .error_handling.error_recovery_system import ErrorRecoverySystem, get_error_recovery_system

# Batch processing imports
from .batch_processing.dynamic_batch_controller import DynamicBatchController

# Monitoring imports
from .monitoring.stats_manager import StatsManager, get_stats_manager

# Make main class easily accessible
__all__ = [
    'KoreaInvestment',
    'MARKET_CODE_MAP',
    'EXCHANGE_CODE_MAP',
    'API_RETURN_CODE',
    'EnhancedRateLimiter',
    'retry_on_rate_limit',
    'retry_on_network_error',
    'get_backoff_strategy',
    'get_error_recovery_system',
    'DynamicBatchController',
    'get_stats_manager',
]
