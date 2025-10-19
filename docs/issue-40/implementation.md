# Implementation Guide: Korea Investment Stock 단순화

> 이 문서는 [PRD](prd.md)의 구현 상세 가이드입니다.

## 📊 File Changes Matrix

### 완전 삭제 대상 (16개 파일)

**Modules (12 files)**
```
✗ korea_investment_stock/rate_limiting/
  ├── enhanced_rate_limiter.py (~400 lines)
  ├── enhanced_backoff_strategy.py (~300 lines)
  ├── enhanced_retry_decorator.py (~200 lines)
  └── __init__.py (~50 lines)

✗ korea_investment_stock/caching/
  ├── ttl_cache.py (~500 lines)
  ├── market_hours.py (~100 lines)
  └── __init__.py (~50 lines)

✗ korea_investment_stock/visualization/
  ├── plotly_visualizer.py (~400 lines)
  ├── dashboard.py (~350 lines)
  ├── charts.py (~250 lines)
  └── __init__.py (~50 lines)

✗ korea_investment_stock/batch_processing/
  ├── dynamic_batch_controller.py (~300 lines)
  └── __init__.py (~30 lines)

✗ korea_investment_stock/monitoring/
  ├── stats_manager.py (~600 lines)
  └── __init__.py (~30 lines)

✗ korea_investment_stock/error_handling/
  ├── error_recovery_system.py (~500 lines)
  └── __init__.py (~30 lines)
```

**Examples (4 files)**
```
✗ examples/rate_limiting_example.py
✗ examples/stats_management_example.py
✗ examples/stats_visualization_plotly.py
✗ examples/visualization_integrated_example.py
```

**총 삭제**: ~4,090 lines

---

### 수정 대상 파일 상세

#### 1. korea_investment_stock/korea_investment_stock.py (주요 수정)

**Import 문 제거 (10줄)**
```python
# 제거할 imports
from .rate_limiting.enhanced_rate_limiter import EnhancedRateLimiter
from .rate_limiting.enhanced_backoff_strategy import get_backoff_strategy
from .rate_limiting.enhanced_retry_decorator import retry_on_rate_limit, retry_on_network_error
from .error_handling.error_recovery_system import get_error_recovery_system
from .monitoring.stats_manager import get_stats_manager
from .caching import TTLCache, cacheable, CACHE_TTL_CONFIG
from .visualization import PlotlyVisualizer, DashboardManager

# VISUALIZATION_AVAILABLE 관련 try-except 블록 제거
```

**__init__() 메서드 간소화**

Before (~100 lines):
```python
def __init__(self, api_key: str, api_secret: str, acc_no: str,
             mock: bool = True, max_workers: int = 3, cache_enabled: bool = True):
    # 기본 설정
    self.api_key = api_key
    self.api_secret = api_secret
    self.acc_no = acc_no
    self.base_url = None
    self.set_base_url(mock)
    
    # Rate limiting 초기화
    self.rate_limiter = EnhancedRateLimiter(max_calls=15, time_window=1.0)
    self.backoff_strategy = get_backoff_strategy()
    self._rate_limit_semaphore = threading.Semaphore(max_workers)
    
    # Cache 초기화
    self.cache = TTLCache(max_size=10000, default_ttl=300)
    self.cache_enabled = cache_enabled
    
    # ThreadPoolExecutor 초기화
    self.executor = ThreadPoolExecutor(max_workers=max_workers)
    self._shutdown_event = threading.Event()
    atexit.register(self.shutdown)
    
    # Monitoring 초기화
    self.stats_manager = get_stats_manager()
    self.error_recovery = get_error_recovery_system()
    
    # Visualization 초기화
    if VISUALIZATION_AVAILABLE:
        self.visualizer = PlotlyVisualizer()
        self.dashboard_manager = DashboardManager()
    else:
        self.visualizer = None
        self.dashboard_manager = None
```

After (~20 lines):
```python
def __init__(self, api_key: str, api_secret: str, acc_no: str, mock: bool = True):
    """한국투자증권 API 클라이언트 초기화
    
    Args:
        api_key: API 키
        api_secret: API 시크릿
        acc_no: 계좌번호
        mock: Mock 서버 사용 여부 (기본값: True)
    """
    self.api_key = api_key
    self.api_secret = api_secret
    self.acc_no = acc_no
    self.base_url = None
    self.set_base_url(mock)
    self.access_token = None
```

**List 기반 메서드 제거 (6개 메서드, ~170 lines)**

Line 814-816:
```python
# 삭제
def fetch_search_stock_info_list(self, stock_market_list):
    return self.__execute_concurrent_requests_with_cache(...)
```

Line 817-819:
```python
# 삭제
def fetch_price_list(self, stock_list):
    return self.__execute_concurrent_requests_with_cache(...)
```

Line 820-838:
```python
# 삭제
def fetch_price_list_with_batch(self, stock_list, batch_size=50, batch_delay=1.0, progress_interval=10):
    return self.__execute_concurrent_requests(...)
```

Line 840-863:
```python
# 삭제
def fetch_price_list_with_dynamic_batch(self, stock_list, dynamic_batch_controller=None):
    # DynamicBatchController 사용
    ...
```

Line 1212-1218:
```python
# 삭제
def fetch_price_detail_oversea_list(self, stock_market_list):
    return self.__execute_concurrent_requests_with_cache(...)
```

Line 1262-1268:
```python
# 삭제
def fetch_stock_info_list(self, stock_market_list):
    return self.__execute_concurrent_requests_with_cache(...)
```

Line 1302-1308:
```python
# 삭제 (중복 정의)
def fetch_search_stock_info_list(self, stock_market_list):
    return self.__execute_concurrent_requests_with_cache(...)
```

**내부 실행 메서드 제거 (2개 메서드, ~230 lines)**

Line 290-582:
```python
# 삭제
def __execute_concurrent_requests(self, method, stock_list, 
                                   batch_size=50, batch_delay=1.0, 
                                   progress_interval=10):
    """ThreadPoolExecutor 기반 병렬 실행"""
    # 150 lines of concurrent execution logic
    ...
```

Line 1349-1450:
```python
# 삭제
def __execute_concurrent_requests_with_cache(self, method, stock_list,
                                              batch_size=50, batch_delay=1.0):
    """캐시 통합 병렬 실행"""
    # 80 lines of cache + concurrent logic
    ...
```

**Private → Public 메서드 전환 (8개 메서드)**

1. Line 865: `__fetch_price` → `fetch_price`
```python
# BEFORE
def __fetch_price(self, symbol: str, market: str = "KR") -> dict:
    """내부 메서드: 단일 주식 가격 조회"""
    if market == "KR":
        # 국내 주식 처리
        symbol_info = self.__get_symbol_type({"symbol": symbol})
        ...

# AFTER
def fetch_price(self, symbol: str, market: str = "KR") -> dict:
    """단일 주식 가격 조회
    
    Args:
        symbol: 종목 코드 (예: "005930" - 삼성전자, "AAPL" - Apple)
        market: 시장 구분 ("KR" 또는 "US", 기본값: "KR")
    
    Returns:
        dict: 가격 정보
            - KR: stck_prpr (현재가), prdy_vrss (전일대비), prdy_ctrt (등락률) 등
            - US: last (현재가), diff (전일대비), rate (등락률) 등
    
    Example:
        >>> broker = KoreaInvestment(api_key, secret, acc_no)
        >>> price = broker.fetch_price("005930", "KR")
        >>> print(f"현재가: {price['stck_prpr']}")
        
        >>> us_price = broker.fetch_price("AAPL", "US")
        >>> print(f"Last: {us_price['last']}")
    """
    if market == "KR":
        symbol_info = self.get_symbol_type({"symbol": symbol})  # Public 호출
        ...
```

2. Line 893: `__get_symbol_type` → `get_symbol_type`
```python
# BEFORE
def __get_symbol_type(self, symbol_info):
    """심볼 타입 판단 (주식/ETF)"""

# AFTER
def get_symbol_type(self, symbol_info):
    """심볼 타입 판단 (주식/ETF)
    
    Args:
        symbol_info: dict with 'symbol' key
    
    Returns:
        str: 'stock' 또는 'etf'
    """
```

3. Line 907: `__fetch_etf_domestic_price` → `fetch_etf_domestic_price`
```python
# BEFORE
@cacheable(ttl=300, key_generator=lambda self, market_code, symbol: f"etf_price:{market_code}:{symbol}")
@retry_on_rate_limit()
def __fetch_etf_domestic_price(self, market_code: str, symbol: str) -> dict:
    with self.rate_limiter.acquire():
        response = self._call(url, headers, params)
        return response.json()

# AFTER (데코레이터 제거, Rate limiter 제거)
def fetch_etf_domestic_price(self, market_code: str, symbol: str) -> dict:
    """국내 ETF 현재가 조회
    
    Args:
        market_code: 시장 코드 ("J" - 코스피, "Q" - 코스닥)
        symbol: ETF 종목 코드
    
    Returns:
        dict: ETF 가격 정보
    """
    response = self._call(url, headers, params)
    return response.json()
```

4. Line 940: `__fetch_domestic_price` → `fetch_domestic_price`
5. Line 1220: `__fetch_price_detail_oversea` → `fetch_price_detail_oversea`
6. Line 1270: `__fetch_stock_info` → `fetch_stock_info`
7. Line 1310: `__fetch_search_stock_info` → `fetch_search_stock_info`

8. Line 583: `__handle_rate_limit_error` → **삭제** (DEPRECATED)

**Cache 관련 메서드 제거 (5개 메서드, ~80 lines)**

Line 1452-1469:
```python
# 삭제
def clear_cache(self, pattern: Optional[str] = None):
    """캐시 초기화"""
    ...
```

Line 1471-1496:
```python
# 삭제
def get_cache_stats(self) -> dict:
    """캐시 통계 조회"""
    ...
```

Line 1498-1505:
```python
# 삭제
def set_cache_enabled(self, enabled: bool):
    """캐시 활성화/비활성화"""
    ...
```

Line 1507-1534:
```python
# 삭제
def preload_cache(self, symbols: List[str], market: str = "KR"):
    """캐시 사전 로딩"""
    ...
```

**Monitoring/Visualization 관련 메서드 제거 (7개 메서드, ~150 lines)**

Line 1536-1568:
```python
# 삭제
def create_monitoring_dashboard(self, ...):
    """모니터링 대시보드 생성"""
    ...
```

Line 1570-1589:
```python
# 삭제
def save_monitoring_dashboard(self, filename: str):
    """대시보드 HTML 저장"""
    ...
```

Line 1591-1610:
```python
# 삭제
def create_stats_report(self, save_as: str = "monitoring_report") -> Dict[str, str]:
    """통계 리포트 생성"""
    ...
```

Line 1612-1632:
```python
# 삭제
def get_system_health_chart(self) -> Optional[Any]:
    """시스템 건강 차트"""
    ...
```

Line 1634-1666:
```python
# 삭제
def get_api_usage_chart(self, hours: int = 24) -> Optional[Any]:
    """API 사용량 차트"""
    ...
```

Line 1668-1683:
```python
# 삭제
def show_monitoring_dashboard(self):
    """브라우저에서 대시보드 표시"""
    ...
```

**데코레이터 제거 (13개 위치)**

Line 731: `issue_access_token()`
```python
# BEFORE
@retry_on_rate_limit(max_retries=3)
def issue_access_token(self):

# AFTER
def issue_access_token(self):
```

Line 902-906: `__fetch_etf_domestic_price()` (이미 Public 전환에서 처리)
Line 935-939: `__fetch_domestic_price()` (이미 Public 전환에서 처리)
Line 968: `fetch_kospi_symbols()`
Line 1002: `fetch_kosdaq_symbols()`
Line 1215-1219: `__fetch_price_detail_oversea()` (이미 Public 전환에서 처리)
Line 1265-1269: `__fetch_stock_info()` (이미 Public 전환에서 처리)
Line 1305-1309: `__fetch_search_stock_info()` (이미 Public 전환에서 처리)
Line 1780-1784: `fetch_ipo_schedule()`

**shutdown() 메서드 간소화**

Line 602-620:
```python
# BEFORE
def shutdown(self):
    """리소스 정리 및 통계 저장"""
    if self._shutdown_event.is_set():
        return
    
    self._shutdown_event.set()
    
    # ThreadPoolExecutor 종료
    if hasattr(self, 'executor'):
        self.executor.shutdown(wait=True)
    
    # 통계 저장
    if hasattr(self, 'stats_manager'):
        self.stats_manager.save_all_stats()
    
    logger.info("KoreaInvestment 클라이언트 종료 완료")

# AFTER (간소화 또는 완전히 제거)
def shutdown(self):
    """리소스 정리 (단순화됨)"""
    logger.info("KoreaInvestment 클라이언트 종료")
```

**예상 변경**: 1,941 lines → ~800 lines (**~60% 감소**)

---

#### 2. korea_investment_stock/__init__.py

**Before (36 lines)**:
```python
"""A Python port of Korea-Investment-Stock API"""

__version__ = "0.5.0"

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
```

**After (10 lines)**:
```python
"""A Python port of Korea-Investment-Stock API"""

__version__ = "0.6.0"

from .korea_investment_stock import KoreaInvestment, MARKET_CODE_MAP, EXCHANGE_CODE_MAP, API_RETURN_CODE

__all__ = [
    'KoreaInvestment',
    'MARKET_CODE_MAP',
    'EXCHANGE_CODE_MAP',
    'API_RETURN_CODE',
]
```

---

#### 3. pyproject.toml

**변경 사항**:
```toml
# BEFORE
[project]
name = "korea-investment-stock"
version = "0.5.0"
dependencies = [
    "requests",
    "pandas",
    "websockets",
    "pycryptodome",
    "crypto>=1.4.1",
]

# AFTER
[project]
name = "korea-investment-stock"
version = "0.6.0"  # Major version bump (Breaking changes)
dependencies = [
    "requests",
    "pandas",
    "websockets",
    "pycryptodome",
    "crypto>=1.4.1",
]

# plotly는 선택적 의존성에서도 제거 (visualization 모듈 삭제로)
```

---

## 🧪 Testing Strategy 상세

### 테스트 삭제 대상 (12개 파일)

| 파일명 | 라인 수 | 제거 이유 |
|--------|---------|-----------|
| `test_rate_limiter.py` | ~300 lines | Rate limiting 모듈 삭제 |
| `test_enhanced_backoff.py` | ~200 lines | Backoff strategy 모듈 삭제 |
| `test_rate_limit_error_detection.py` | ~150 lines | Rate limit 에러 검출 삭제 |
| `test_rate_limit_simulation.py` | ~250 lines | Rate limit 시뮬레이션 삭제 |
| `test_ttl_cache.py` | ~400 lines | TTL Cache 모듈 삭제 |
| `test_cache_integration.py` | ~300 lines | Cache 통합 테스트 삭제 |
| `test_batch_processing.py` | ~200 lines | Batch processing 모듈 삭제 |
| `test_error_recovery.py` | ~250 lines | Error recovery 시스템 삭제 |
| `test_error_handling.py` | ~200 lines | Error handling 모듈 삭제 |
| `test_stats_save.py` | ~150 lines | Stats 저장 기능 삭제 |
| `test_enhanced_integration.py` | ~300 lines | Enhanced 기능 통합 테스트 삭제 |
| `test_threadpool_improvement.py` | ~200 lines | ThreadPool 개선 기능 삭제 |

**총 삭제**: ~2,900 lines

---

### 테스트 업데이트 상세

#### 1. test_korea_investment_stock.py

**삭제할 테스트 클래스/메서드**:
```python
class TestBatchProcessing:
    def test_fetch_price_list(self):  # ❌
    def test_fetch_price_list_with_batch(self):  # ❌
    def test_fetch_price_list_with_dynamic_batch(self):  # ❌
    def test_concurrent_requests(self):  # ❌
    def test_concurrent_requests_with_cache(self):  # ❌

class TestCaching:
    def test_cache_hit(self):  # ❌
    def test_cache_miss(self):  # ❌
    def test_clear_cache(self):  # ❌
    def test_get_cache_stats(self):  # ❌

class TestRateLimiting:
    def test_rate_limiter_acquire(self):  # ❌
    def test_backoff_on_error(self):  # ❌
```

**추가할 테스트 클래스/메서드**:
```python
class TestSingleFetch:
    """Public 전환된 단일 조회 메서드 테스트"""
    
    def test_fetch_price_kr(self):
        """국내 주식 단일 조회"""
        broker = KoreaInvestment(api_key, secret, acc_no)
        price = broker.fetch_price("005930", "KR")
        
        assert "stck_prpr" in price  # 현재가
        assert "prdy_vrss" in price  # 전일대비
        assert "prdy_ctrt" in price  # 등락률
    
    def test_fetch_price_us(self):
        """해외 주식 단일 조회"""
        broker = KoreaInvestment(api_key, secret, acc_no, mock=False)
        price = broker.fetch_price("AAPL", "US")
        
        assert "last" in price
        assert "diff" in price
    
    def test_fetch_domestic_price(self):
        """국내 주식 가격 직접 조회"""
        broker = KoreaInvestment(api_key, secret, acc_no)
        price = broker.fetch_domestic_price("J", "005930")
        
        assert "output" in price
    
    def test_fetch_etf_domestic_price(self):
        """국내 ETF 가격 조회"""
        broker = KoreaInvestment(api_key, secret, acc_no)
        price = broker.fetch_etf_domestic_price("J", "069500")  # KODEX 200
        
        assert "output" in price
    
    def test_fetch_stock_info(self):
        """주식 정보 조회"""
        broker = KoreaInvestment(api_key, secret, acc_no)
        info = broker.fetch_stock_info("005930", "KR")
        
        assert "output" in info
    
    def test_get_symbol_type(self):
        """심볼 타입 판단"""
        broker = KoreaInvestment(api_key, secret, acc_no)
        
        stock_type = broker.get_symbol_type({"symbol": "005930"})
        assert stock_type == "stock"
        
        etf_type = broker.get_symbol_type({"symbol": "069500"})
        assert etf_type == "etf"
```

---

#### 2. test_integration.py

**Before**:
```python
def test_fetch_price_list():
    """여러 주식 동시 조회 테스트"""
    with KoreaInvestment(api_key, secret, acc_no) as broker:
        stock_list = [("005930", "KR"), ("035420", "KR"), ("000660", "KR")]
        prices = broker.fetch_price_list(stock_list)
        
        assert len(prices) == 3
        # Rate limit 체크
        assert broker.rate_limiter.get_stats()['total_calls'] >= 3
```

**After**:
```python
def test_fetch_price():
    """단일 주식 조회 테스트"""
    broker = KoreaInvestment(api_key, secret, acc_no)
    
    # 단일 조회
    price = broker.fetch_price("005930", "KR")
    assert "stck_prpr" in price
    
    # 다중 조회는 사용자가 직접 제어
    symbols = ["005930", "035420", "000660"]
    prices = []
    
    for symbol in symbols:
        price = broker.fetch_price(symbol, "KR")
        prices.append(price)
    
    assert len(prices) == 3
```

---

#### 3. test_integration_us_stocks.py

**추가 테스트**:
```python
def test_fetch_price_detail_oversea():
    """해외 주식 상세 조회 (Public 전환)"""
    broker = KoreaInvestment(api_key, secret, acc_no, mock=False)
    
    price = broker.fetch_price_detail_oversea("AAPL", "US")
    
    assert "output" in price
    # PER, PBR 등 상세 정보 포함
```

---

#### 4. test_load.py

**Before** (ThreadPoolExecutor 사용):
```python
def test_concurrent_load():
    """병렬 부하 테스트"""
    with KoreaInvestment(api_key, secret, acc_no) as broker:
        stock_list = [(f"{i:06d}", "KR") for i in range(100)]
        
        start_time = time.time()
        prices = broker.fetch_price_list(stock_list)
        duration = time.time() - start_time
        
        assert len(prices) == 100
        assert duration < 20  # 병렬 처리로 빠름
```

**After** (단순 loop):
```python
def test_sequential_load():
    """순차 부하 테스트 (사용자 제어)"""
    broker = KoreaInvestment(api_key, secret, acc_no)
    symbols = ["005930", "035420", "000660", "005380", "068270"]
    
    prices = []
    errors = 0
    
    for symbol in symbols:
        try:
            price = broker.fetch_price(symbol, "KR")
            prices.append(price)
            time.sleep(0.1)  # 사용자가 직접 rate limit 제어
        except Exception as e:
            errors += 1
            print(f"Error fetching {symbol}: {e}")
    
    assert len(prices) >= 4  # 최소 4개 성공
    assert errors <= 1  # 에러 허용 범위
```

---

### 신규 테스트 파일: test_public_api.py

```python
"""
Public API 전환 후 기본 동작 검증

Private 메서드가 Public으로 전환되면서 사용자가 직접 호출 가능해진
메서드들의 기본 동작을 검증합니다.
"""

import time
import pytest
from korea_investment_stock import KoreaInvestment


class TestPublicAPI:
    """Public 전환 메서드 테스트"""
    
    @pytest.fixture
    def broker(self):
        """테스트용 broker 인스턴스"""
        return KoreaInvestment(
            api_key="test_key",
            api_secret="test_secret",
            acc_no="12345678-01",
            mock=True
        )
    
    def test_fetch_price_kr(self, broker):
        """국내 주식 단일 조회"""
        price = broker.fetch_price("005930", "KR")
        assert "stck_prpr" in price
    
    def test_fetch_price_us(self):
        """해외 주식 단일 조회 (실제 계정 필요)"""
        broker = KoreaInvestment(
            api_key="real_key",
            api_secret="real_secret",
            acc_no="real_acc",
            mock=False
        )
        price = broker.fetch_price("AAPL", "US")
        assert "last" in price
    
    def test_user_controlled_batch(self, broker):
        """사용자 제어 배치 조회"""
        symbols = ["005930", "035420", "000660"]
        prices = []
        
        for symbol in symbols:
            price = broker.fetch_price(symbol, "KR")
            prices.append(price)
            time.sleep(0.1)  # 사용자가 직접 rate limit 제어
        
        assert len(prices) == 3
    
    def test_user_controlled_retry(self, broker):
        """사용자 구현 재시도 로직"""
        def fetch_with_retry(symbol, retries=3):
            for i in range(retries):
                try:
                    return broker.fetch_price(symbol, "KR")
                except Exception as e:
                    if i == retries - 1:
                        raise
                    time.sleep(2 ** i)  # Exponential backoff
        
        price = fetch_with_retry("005930")
        assert "stck_prpr" in price


class TestUserImplementation:
    """사용자 직접 구현 패턴 테스트"""
    
    def test_user_caching(self):
        """사용자 구현 캐싱"""
        from datetime import datetime, timedelta
        
        broker = KoreaInvestment("key", "secret", "acc", mock=True)
        cache = {}
        cache_ttl = timedelta(minutes=5)
        
        def fetch_with_cache(symbol, market="KR"):
            cache_key = f"{symbol}:{market}"
            now = datetime.now()
            
            # Cache hit check
            if cache_key in cache:
                cached_time, cached_price = cache[cache_key]
                if now - cached_time < cache_ttl:
                    return cached_price
            
            # Cache miss
            price = broker.fetch_price(symbol, market)
            cache[cache_key] = (now, price)
            return price
        
        # First call - cache miss
        price1 = fetch_with_cache("005930", "KR")
        
        # Second call - cache hit
        price2 = fetch_with_cache("005930", "KR")
        
        assert price1 == price2  # Same result from cache
```

---

## 📚 Example 파일 수정

### 삭제 대상 (4개)
```bash
rm examples/rate_limiting_example.py
rm examples/stats_management_example.py
rm examples/stats_visualization_plotly.py
rm examples/visualization_integrated_example.py
```

### 업데이트 대상 (2개)

#### examples/ipo_schedule_example.py

**Before**:
```python
with KoreaInvestment(api_key, secret, acc_no) as broker:
    # @cacheable 데코레이터가 자동 캐싱
    ipo_data = broker.fetch_ipo_schedule()
```

**After**:
```python
broker = KoreaInvestment(api_key, secret, acc_no)

# 데코레이터 제거됨 - 필요시 사용자가 직접 캐싱 구현
ipo_data = broker.fetch_ipo_schedule()
```

#### examples/us_stock_price_example.py

**Before**:
```python
with KoreaInvestment(api_key, secret, acc_no, mock=False) as broker:
    stock_list = [("AAPL", "US"), ("TSLA", "US"), ("MSFT", "US")]
    prices = broker.fetch_price_list(stock_list)
```

**After**:
```python
broker = KoreaInvestment(api_key, secret, acc_no, mock=False)

# 단일 조회
price = broker.fetch_price("AAPL", "US")
print(f"AAPL: {price['last']}")

# 배치 조회 (사용자 제어)
symbols = ["AAPL", "TSLA", "MSFT"]
prices = []

for symbol in symbols:
    try:
        price = broker.fetch_price(symbol, "US")
        prices.append(price)
        time.sleep(0.1)  # Rate limiting
    except Exception as e:
        print(f"Error: {e}")
```

### 신규 Example: basic_usage_example.py

```python
"""
Korea Investment Stock 기본 사용 예시 (v0.6.0+)

v0.6.0부터 단순한 API Wrapper로 변경되었습니다.
Rate limiting, Caching, Batch processing은 사용자가 직접 구현합니다.
"""

import time
from korea_investment_stock import KoreaInvestment

# 1. 기본 초기화
broker = KoreaInvestment(
    api_key="YOUR_API_KEY",
    api_secret="YOUR_API_SECRET",
    acc_no="YOUR_ACCOUNT_NO",
    mock=True  # Mock 서버 (실제 거래는 False)
)

# 2. 단일 주식 조회 (Public 메서드)
print("=== 단일 조회 ===")
price = broker.fetch_price("005930", "KR")  # 삼성전자
print(f"현재가: {price['stck_prpr']}")
print(f"등락률: {price['prdy_ctrt']}%")

# 3. 배치 조회 (사용자 제어)
print("\n=== 배치 조회 (사용자 제어) ===")
symbols = ["005930", "035420", "000660"]
prices = []

for symbol in symbols:
    try:
        price = broker.fetch_price(symbol, "KR")
        prices.append(price)
        print(f"{symbol}: {price['stck_prpr']}")
        time.sleep(0.1)  # Rate limiting (초당 10개)
    except Exception as e:
        print(f"Error fetching {symbol}: {e}")

print(f"총 {len(prices)}개 조회 완료")

# 4. 재시도 로직 (사용자 구현)
print("\n=== 재시도 로직 ===")
def fetch_with_retry(symbol, market="KR", retries=3):
    """지수 백오프 재시도"""
    for i in range(retries):
        try:
            return broker.fetch_price(symbol, market)
        except Exception as e:
            if i == retries - 1:
                raise
            wait_time = 2 ** i
            print(f"Retry {i+1}/{retries} after {wait_time}s...")
            time.sleep(wait_time)

price = fetch_with_retry("005930")
print(f"재시도 성공: {price['stck_prpr']}")

# 5. 캐싱 (사용자 구현)
print("\n=== 캐싱 구현 ===")
from datetime import datetime, timedelta

cache = {}
cache_ttl = timedelta(minutes=5)

def fetch_with_cache(symbol, market="KR"):
    """TTL 기반 캐싱"""
    cache_key = f"{symbol}:{market}"
    now = datetime.now()
    
    # Cache hit check
    if cache_key in cache:
        cached_time, cached_price = cache[cache_key]
        if now - cached_time < cache_ttl:
            print(f"Cache HIT: {cache_key}")
            return cached_price
    
    # Cache miss
    print(f"Cache MISS: {cache_key}")
    price = broker.fetch_price(symbol, market)
    cache[cache_key] = (now, price)
    return price

# First call - cache miss
price1 = fetch_with_cache("005930")

# Second call - cache hit
price2 = fetch_with_cache("005930")

# 6. IPO 조회 (변경 없음)
print("\n=== IPO 일정 조회 ===")
ipo_data = broker.fetch_ipo_schedule(
    from_date="20250101",
    to_date="20250131"
)
print(f"IPO 건수: {len(ipo_data)}")
```

---

## 📝 Documentation 업데이트

### README.md 주요 변경사항

**Features 섹션**:
```markdown
# BEFORE
## Features

- ✅ Rate Limiting (Token Bucket + Sliding Window)
- ✅ Automatic Retry with Exponential Backoff
- ✅ TTL-based Caching (5min for prices, 5hrs for stock info)
- ✅ Batch Processing with Dynamic Adjustment
- ✅ Real-time Monitoring & Visualization
- ✅ Error Recovery System
- ✅ Domestic & US Stock Support
- ✅ IPO Schedule Lookup

# AFTER
## Features

- ✅ 순수 API Wrapper (한국투자증권 OpenAPI)
- ✅ 국내/해외 주식 조회
- ✅ IPO 일정 조회
- ✅ 간단하고 명확한 API
- ✅ 사용자 제어 가능 (Rate limiting, Caching, Batch processing)
```

**Usage 섹션**:
```markdown
# BEFORE
## Usage

```python
with KoreaInvestment(api_key, secret, acc_no) as broker:
    # Automatic batch processing, caching, rate limiting
    prices = broker.fetch_price_list(stock_list)
    
    # Built-in monitoring
    broker.save_monitoring_dashboard("dashboard.html")
```

# AFTER
## Usage

```python
broker = KoreaInvestment(api_key, secret, acc_no)

# Single query
price = broker.fetch_price("005930", "KR")

# Batch query (user-controlled)
for symbol in stock_list:
    price = broker.fetch_price(symbol, "KR")
    time.sleep(0.1)  # Rate limiting
```

**Migration Guide 링크 추가**:
```markdown
## Migration from 0.5.x to 0.6.0

**Breaking Changes**: v0.6.0은 주요 변경사항을 포함합니다.

자세한 마이그레이션 가이드는 [docs/issue-40/prd.md](docs/issue-40/prd.md)를 참조하세요.
```

---

### CLAUDE.md 주요 변경사항

**Architecture Overview 섹션 업데이트**:

```markdown
# BEFORE
## Architecture Overview

### Core Component Flow

```
User API Call
  ↓
@retry_on_rate_limit decorator (5 retries)
  ↓
@cacheable decorator (TTL-based)
  ↓
EnhancedRateLimiter.acquire() (Token Bucket + Sliding Window)
  ↓
HTTP Request to Korea Investment API
  ↓
Error Recovery System
  ↓
Circuit Breaker
```

# AFTER
## Architecture Overview

### Simplified Component Flow

```
User API Call
  ↓
KoreaInvestment Method
  ↓
HTTP Request to Korea Investment API
  ↓
JSON Response
```

사용자가 필요에 따라 직접 구현:
- Rate Limiting (time.sleep 등)
- Caching (dict, redis 등)
- Retry Logic (for loop + try-except)
- Monitoring (logging, metrics 등)
```

**Key Modules 섹션 업데이트**:
```markdown
# BEFORE
1. **rate_limiting/** - Rate limiting system
2. **caching/** - TTL cache system
3. **visualization/** - Monitoring dashboards
4. **batch_processing/** - Dynamic batch control
5. **monitoring/** - Statistics management
6. **error_handling/** - Error recovery

# AFTER
1. **korea_investment_stock.py** - Main API wrapper (800 lines)
2. **utils/** - Utility functions (if any)

모든 고급 기능은 제거되었습니다.
```

---

### CHANGELOG.md 추가

```markdown
## [0.6.0] - 2025-01-XX

### 🚨 Breaking Changes

**Major Simplification**: 한국투자증권 API의 순수 Wrapper로 단순화

**Removed Modules**:
- ❌ `rate_limiting/` - Rate limiting system
- ❌ `caching/` - TTL cache system
- ❌ `visualization/` - Monitoring & dashboards
- ❌ `batch_processing/` - Dynamic batch controller
- ❌ `monitoring/` - Statistics manager
- ❌ `error_handling/` - Error recovery system

**Removed Methods**:
- ❌ `fetch_price_list()` - Use `fetch_price()` in a loop
- ❌ `fetch_price_list_with_batch()`
- ❌ `fetch_price_list_with_dynamic_batch()`
- ❌ `fetch_stock_info_list()` - Use `fetch_stock_info()` in a loop
- ❌ `fetch_search_stock_info_list()`
- ❌ `fetch_price_detail_oversea_list()`
- ❌ `clear_cache()`, `get_cache_stats()`, `set_cache_enabled()`, `preload_cache()`
- ❌ `create_monitoring_dashboard()`, `save_monitoring_dashboard()`, etc.

**Changed Methods** (Private → Public):
- ✅ `__fetch_price()` → `fetch_price()` (now public)
- ✅ `__fetch_domestic_price()` → `fetch_domestic_price()`
- ✅ `__fetch_etf_domestic_price()` → `fetch_etf_domestic_price()`
- ✅ `__fetch_price_detail_oversea()` → `fetch_price_detail_oversea()`
- ✅ `__fetch_stock_info()` → `fetch_stock_info()`
- ✅ `__fetch_search_stock_info()` → `fetch_search_stock_info()`
- ✅ `__get_symbol_type()` → `get_symbol_type()`

**Migration Guide**: See [docs/issue-40/prd.md](docs/issue-40/prd.md)

### Changed
- `__init__()` simplified - No more ThreadPoolExecutor, RateLimiter, Cache initialization
- All decorators removed (@cacheable, @retry_on_rate_limit)
- Dependencies: plotly removed

### Fixed
- None (this is a simplification release)

---

## [0.5.0] - 2024-XX-XX

(Previous version with all features)
```

---

## 🔧 Implementation Tips

### 1. Git Workflow

```bash
# 1. Feature branch 생성
git checkout -b feat/issue-40-simplify

# 2. 단계별 커밋
git add korea_investment_stock/rate_limiting
git commit -m "[feat] #40 - Remove rate limiting module"

git add korea_investment_stock/caching
git commit -m "[feat] #40 - Remove caching module"

# ... (각 모듈 단위로 커밋)

git add korea_investment_stock/korea_investment_stock.py
git commit -m "[feat] #40 - Simplify main module (Private → Public, remove decorators)"

git add korea_investment_stock/__init__.py
git commit -m "[feat] #40 - Update package exports"

git add korea_investment_stock/tests
git commit -m "[feat] #40 - Update tests for simplified API"

git add examples
git commit -m "[feat] #40 - Update examples for simplified API"

git add README.md CLAUDE.md CHANGELOG.md
git commit -m "[feat] #40 - Update documentation"

git add pyproject.toml
git commit -m "[feat] #40 - Bump version to 0.6.0"

# 3. Push & PR
git push origin feat/issue-40-simplify
gh pr create --title "[feat] #40 - Simplify library to pure API wrapper" \
  --body "$(cat docs/issue-40/prd.md)"
```

### 2. 검증 스크립트

```bash
#!/bin/bash
# verify_simplification.sh

echo "=== Verification Script ==="

# 1. 삭제된 모듈 확인
echo "1. Checking deleted modules..."
! test -d korea_investment_stock/rate_limiting && echo "✓ rate_limiting deleted"
! test -d korea_investment_stock/caching && echo "✓ caching deleted"
! test -d korea_investment_stock/visualization && echo "✓ visualization deleted"
! test -d korea_investment_stock/batch_processing && echo "✓ batch_processing deleted"
! test -d korea_investment_stock/monitoring && echo "✓ monitoring deleted"
! test -d korea_investment_stock/error_handling && echo "✓ error_handling deleted"

# 2. 라인 수 확인
echo -e "\n2. Checking line count..."
lines=$(wc -l < korea_investment_stock/korea_investment_stock.py)
echo "Main module: $lines lines (target: ~800 lines)"
if [ $lines -lt 1000 ]; then
    echo "✓ Line count acceptable"
else
    echo "✗ Line count too high"
fi

# 3. Public 메서드 확인
echo -e "\n3. Checking public methods..."
grep -c "^    def fetch_price(" korea_investment_stock/korea_investment_stock.py > /dev/null && echo "✓ fetch_price() is public"
grep -c "^    def fetch_domestic_price(" korea_investment_stock/korea_investment_stock.py > /dev/null && echo "✓ fetch_domestic_price() is public"

# 4. 데코레이터 제거 확인
echo -e "\n4. Checking decorator removal..."
! grep -q "@retry_on_rate_limit" korea_investment_stock/korea_investment_stock.py && echo "✓ @retry_on_rate_limit removed"
! grep -q "@cacheable" korea_investment_stock/korea_investment_stock.py && echo "✓ @cacheable removed"

# 5. 테스트 실행
echo -e "\n5. Running tests..."
pytest korea_investment_stock/tests/ -v

echo -e "\n=== Verification Complete ==="
```

### 3. Before/After 비교

```bash
# Before (0.5.0)
$ find korea_investment_stock -name "*.py" | wc -l
32 files

$ wc -l korea_investment_stock/korea_investment_stock.py
1941 lines

# After (0.6.0)
$ find korea_investment_stock -name "*.py" | wc -l
8 files  # ~75% reduction

$ wc -l korea_investment_stock/korea_investment_stock.py
800 lines  # ~60% reduction
```

---

**작성**: Claude Code  
**검토**: (To be reviewed)  
**승인**: (To be approved)
