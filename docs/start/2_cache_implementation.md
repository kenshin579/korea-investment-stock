# 캐싱 기능 구현 가이드

## 🎯 구현 아키텍처

**Option B: Wrapper 클래스 패턴**

```
KoreaInvestment (기존)
    ↓
CachedKoreaInvestment (래퍼)
    ↓
CacheManager (캐시 엔진)
```

---

## 📦 구현 파일 구조

```
korea_investment_stock/
├── cache_manager.py              # 캐시 매니저
├── cached_korea_investment.py    # 래퍼 클래스
├── __init__.py                   # 모듈 export 추가
└── tests/
    ├── test_cache_manager.py     # 단위 테스트
    ├── test_cached_integration.py # 통합 테스트
    └── test_cache_performance.py  # 성능 테스트
```

---

## 💻 1. CacheManager 구현

### korea_investment_stock/cache_manager.py

```python
from typing import Dict, Any, Optional
from datetime import datetime, timedelta
import threading

class CacheEntry:
    """캐시 엔트리"""
    def __init__(self, data: Any, ttl_seconds: int):
        self.data = data
        self.cached_at = datetime.now()
        self.expires_at = self.cached_at + timedelta(seconds=ttl_seconds)

    def is_expired(self) -> bool:
        """만료 여부 확인"""
        return datetime.now() > self.expires_at

    def age_seconds(self) -> float:
        """캐시 생성 후 경과 시간 (초)"""
        return (datetime.now() - self.cached_at).total_seconds()


class CacheManager:
    """메모리 기반 캐시 매니저 (Thread-safe)"""

    def __init__(self):
        self._cache: Dict[str, CacheEntry] = {}
        self._lock = threading.Lock()
        self._stats = {
            'hits': 0,
            'misses': 0,
            'evictions': 0
        }

    def get(self, key: str) -> Optional[Any]:
        """캐시에서 데이터 조회"""
        with self._lock:
            entry = self._cache.get(key)

            if entry is None:
                self._stats['misses'] += 1
                return None

            if entry.is_expired():
                del self._cache[key]
                self._stats['evictions'] += 1
                self._stats['misses'] += 1
                return None

            self._stats['hits'] += 1
            return entry.data

    def set(self, key: str, data: Any, ttl_seconds: int):
        """캐시에 데이터 저장"""
        with self._lock:
            self._cache[key] = CacheEntry(data, ttl_seconds)

    def invalidate(self, key: str):
        """특정 캐시 무효화"""
        with self._lock:
            if key in self._cache:
                del self._cache[key]
                self._stats['evictions'] += 1

    def clear(self):
        """전체 캐시 삭제"""
        with self._lock:
            count = len(self._cache)
            self._cache.clear()
            self._stats['evictions'] += count

    def get_stats(self) -> Dict[str, Any]:
        """캐시 통계 반환"""
        with self._lock:
            total_requests = self._stats['hits'] + self._stats['misses']
            hit_rate = (self._stats['hits'] / total_requests * 100
                       if total_requests > 0 else 0)

            return {
                'cache_size': len(self._cache),
                'hits': self._stats['hits'],
                'misses': self._stats['misses'],
                'evictions': self._stats['evictions'],
                'hit_rate': f"{hit_rate:.2f}%"
            }

    def get_cache_info(self, key: str) -> Optional[Dict[str, Any]]:
        """특정 캐시 엔트리 정보 반환"""
        with self._lock:
            entry = self._cache.get(key)
            if entry is None:
                return None

            return {
                'cached_at': entry.cached_at.isoformat(),
                'expires_at': entry.expires_at.isoformat(),
                'age_seconds': entry.age_seconds(),
                'is_expired': entry.is_expired()
            }
```

---

## 🎁 2. CachedKoreaInvestment 래퍼 구현

### korea_investment_stock/cached_korea_investment.py

```python
from typing import Optional, Dict, Any
from .korea_investment_stock import KoreaInvestment
from .cache_manager import CacheManager

class CachedKoreaInvestment:
    """캐싱 기능이 추가된 KoreaInvestment 래퍼"""

    DEFAULT_TTL = {
        'price': 5,           # 실시간 가격: 5초
        'stock_info': 300,    # 종목 정보: 5분
        'symbols': 3600,      # 종목 리스트: 1시간
        'ipo': 1800           # IPO 일정: 30분
    }

    def __init__(
        self,
        broker: KoreaInvestment,
        enable_cache: bool = True,
        price_ttl: Optional[int] = None,
        stock_info_ttl: Optional[int] = None,
        symbols_ttl: Optional[int] = None,
        ipo_ttl: Optional[int] = None
    ):
        """
        Args:
            broker: KoreaInvestment 인스턴스
            enable_cache: 캐싱 활성화 여부
            price_ttl: 실시간 가격 TTL (초)
            stock_info_ttl: 종목정보 TTL (초)
            symbols_ttl: 종목리스트 TTL (초)
            ipo_ttl: IPO 일정 TTL (초)
        """
        self.broker = broker
        self.enable_cache = enable_cache
        self.cache = CacheManager() if enable_cache else None

        # TTL 설정
        self.ttl = {
            'price': price_ttl or self.DEFAULT_TTL['price'],
            'stock_info': stock_info_ttl or self.DEFAULT_TTL['stock_info'],
            'symbols': symbols_ttl or self.DEFAULT_TTL['symbols'],
            'ipo': ipo_ttl or self.DEFAULT_TTL['ipo']
        }

    def _make_cache_key(self, method: str, *args, **kwargs) -> str:
        """캐시 키 생성"""
        args_str = "_".join(str(arg) for arg in args)
        kwargs_str = "_".join(f"{k}={v}" for k, v in sorted(kwargs.items()))
        return f"{method}:{args_str}:{kwargs_str}"

    def fetch_price(self, symbol: str, market: str = "KR") -> dict:
        """가격 조회 (캐싱 지원)"""
        if not self.enable_cache:
            return self.broker.fetch_price(symbol, market)

        cache_key = self._make_cache_key("fetch_price", symbol, market)
        cached_data = self.cache.get(cache_key)

        if cached_data is not None:
            return cached_data

        result = self.broker.fetch_price(symbol, market)

        if result.get('rt_cd') == '0':
            self.cache.set(cache_key, result, self.ttl['price'])

        return result

    def fetch_domestic_price(self, market_code: str, symbol: str) -> dict:
        """국내 주식 가격 조회 (캐싱 지원)"""
        if not self.enable_cache:
            return self.broker.fetch_domestic_price(market_code, symbol)

        cache_key = self._make_cache_key("fetch_domestic_price", market_code, symbol)
        cached_data = self.cache.get(cache_key)

        if cached_data is not None:
            return cached_data

        result = self.broker.fetch_domestic_price(market_code, symbol)

        if result.get('rt_cd') == '0':
            self.cache.set(cache_key, result, self.ttl['price'])

        return result

    def fetch_etf_domestic_price(self, market_code: str, symbol: str) -> dict:
        """ETF 가격 조회 (캐싱 지원)"""
        if not self.enable_cache:
            return self.broker.fetch_etf_domestic_price(market_code, symbol)

        cache_key = self._make_cache_key("fetch_etf_domestic_price", market_code, symbol)
        cached_data = self.cache.get(cache_key)

        if cached_data is not None:
            return cached_data

        result = self.broker.fetch_etf_domestic_price(market_code, symbol)

        if result.get('rt_cd') == '0':
            self.cache.set(cache_key, result, self.ttl['price'])

        return result

    def fetch_price_detail_oversea(self, symbol: str, market: str = "KR") -> dict:
        """해외 주식 가격 조회 (캐싱 지원)"""
        if not self.enable_cache:
            return self.broker.fetch_price_detail_oversea(symbol, market)

        cache_key = self._make_cache_key("fetch_price_detail_oversea", symbol, market)
        cached_data = self.cache.get(cache_key)

        if cached_data is not None:
            return cached_data

        result = self.broker.fetch_price_detail_oversea(symbol, market)

        if result.get('rt_cd') == '0':
            self.cache.set(cache_key, result, self.ttl['price'])

        return result

    def fetch_stock_info(self, symbol: str, market: str = "KR") -> dict:
        """종목 정보 조회 (캐싱 지원)"""
        if not self.enable_cache:
            return self.broker.fetch_stock_info(symbol, market)

        cache_key = self._make_cache_key("fetch_stock_info", symbol, market)
        cached_data = self.cache.get(cache_key)

        if cached_data is not None:
            return cached_data

        result = self.broker.fetch_stock_info(symbol, market)

        if result.get('rt_cd') == '0':
            self.cache.set(cache_key, result, self.ttl['stock_info'])

        return result

    def fetch_search_stock_info(self, symbol: str, market: str = "KR") -> dict:
        """종목 검색 (캐싱 지원)"""
        if not self.enable_cache:
            return self.broker.fetch_search_stock_info(symbol, market)

        cache_key = self._make_cache_key("fetch_search_stock_info", symbol, market)
        cached_data = self.cache.get(cache_key)

        if cached_data is not None:
            return cached_data

        result = self.broker.fetch_search_stock_info(symbol, market)

        if result.get('rt_cd') == '0':
            self.cache.set(cache_key, result, self.ttl['stock_info'])

        return result

    def fetch_kospi_symbols(self) -> dict:
        """KOSPI 종목 리스트 (캐싱 지원)"""
        if not self.enable_cache:
            return self.broker.fetch_kospi_symbols()

        cache_key = "fetch_kospi_symbols"
        cached_data = self.cache.get(cache_key)

        if cached_data is not None:
            return cached_data

        result = self.broker.fetch_kospi_symbols()
        self.cache.set(cache_key, result, self.ttl['symbols'])

        return result

    def fetch_kosdaq_symbols(self) -> dict:
        """KOSDAQ 종목 리스트 (캐싱 지원)"""
        if not self.enable_cache:
            return self.broker.fetch_kosdaq_symbols()

        cache_key = "fetch_kosdaq_symbols"
        cached_data = self.cache.get(cache_key)

        if cached_data is not None:
            return cached_data

        result = self.broker.fetch_kosdaq_symbols()
        self.cache.set(cache_key, result, self.ttl['symbols'])

        return result

    def fetch_ipo_schedule(
        self,
        from_date: Optional[str] = None,
        to_date: Optional[str] = None,
        symbol: str = ""
    ) -> dict:
        """IPO 일정 조회 (캐싱 지원)"""
        if not self.enable_cache:
            return self.broker.fetch_ipo_schedule(from_date, to_date, symbol)

        cache_key = self._make_cache_key("fetch_ipo_schedule", from_date, to_date, symbol)
        cached_data = self.cache.get(cache_key)

        if cached_data is not None:
            return cached_data

        result = self.broker.fetch_ipo_schedule(from_date, to_date, symbol)

        if result.get('rt_cd') == '0':
            self.cache.set(cache_key, result, self.ttl['ipo'])

        return result

    def invalidate_cache(self, method: Optional[str] = None):
        """캐시 무효화"""
        if not self.enable_cache:
            return

        self.cache.clear()

    def get_cache_stats(self) -> Dict[str, Any]:
        """캐시 통계 반환"""
        if not self.enable_cache:
            return {'cache_enabled': False}

        stats = self.cache.get_stats()
        stats['cache_enabled'] = True
        stats['ttl_config'] = self.ttl
        return stats

    def __enter__(self):
        """컨텍스트 매니저 진입"""
        return self

    def __exit__(self, exc_type, exc_val, exc_tb):
        """컨텍스트 매니저 종료"""
        if self.enable_cache:
            self.cache.clear()
        return False
```

---

## 📝 3. __init__.py 업데이트

### korea_investment_stock/__init__.py

```python
from .korea_investment_stock import KoreaInvestment
from .cache_manager import CacheManager, CacheEntry
from .cached_korea_investment import CachedKoreaInvestment

__all__ = [
    'KoreaInvestment',
    'CacheManager',
    'CacheEntry',
    'CachedKoreaInvestment'
]
```

---

## 🧪 4. 테스트 코드

### 4.1 단위 테스트: test_cache_manager.py

```python
import pytest
import time
from korea_investment_stock.cache_manager import CacheManager, CacheEntry

class TestCacheEntry:
    def test_cache_entry_creation(self):
        data = {"key": "value"}
        entry = CacheEntry(data, ttl_seconds=5)

        assert entry.data == data
        assert not entry.is_expired()
        assert entry.age_seconds() < 1

    def test_cache_entry_expiration(self):
        entry = CacheEntry("test", ttl_seconds=1)
        assert not entry.is_expired()

        time.sleep(1.1)
        assert entry.is_expired()


class TestCacheManager:
    def test_cache_set_get(self):
        cache = CacheManager()
        cache.set("key1", "value1", ttl_seconds=10)
        assert cache.get("key1") == "value1"

    def test_cache_miss(self):
        cache = CacheManager()
        assert cache.get("nonexistent") is None

    def test_cache_expiration(self):
        cache = CacheManager()
        cache.set("key1", "value1", ttl_seconds=1)
        assert cache.get("key1") == "value1"

        time.sleep(1.1)
        assert cache.get("key1") is None

    def test_cache_stats(self):
        cache = CacheManager()
        cache.get("key1")  # miss
        cache.set("key1", "value1", ttl_seconds=10)
        cache.get("key1")  # hit

        stats = cache.get_stats()
        assert stats['hits'] == 1
        assert stats['misses'] == 1
```

### 4.2 통합 테스트: test_cached_integration.py

```python
import pytest
import os
import time
from korea_investment_stock import KoreaInvestment, CachedKoreaInvestment

@pytest.fixture
def broker():
    api_key = os.getenv('KOREA_INVESTMENT_API_KEY')
    api_secret = os.getenv('KOREA_INVESTMENT_API_SECRET')
    acc_no = os.getenv('KOREA_INVESTMENT_ACCOUNT_NO')

    if not all([api_key, api_secret, acc_no]):
        pytest.skip("API credentials not set")

    return KoreaInvestment(api_key, api_secret, acc_no, mock=True)


class TestCachedKoreaInvestment:
    def test_cached_fetch_price(self, broker):
        cached_broker = CachedKoreaInvestment(broker, price_ttl=5)

        # 첫 번째 호출 (캐시 미스)
        result1 = cached_broker.fetch_price("005930", "KR")
        assert result1['rt_cd'] == '0'

        # 두 번째 호출 (캐시 히트)
        result2 = cached_broker.fetch_price("005930", "KR")
        assert result2 == result1

        stats = cached_broker.get_cache_stats()
        assert stats['hits'] == 1
        assert stats['misses'] == 1

    def test_cache_disabled(self, broker):
        cached_broker = CachedKoreaInvestment(broker, enable_cache=False)

        result1 = cached_broker.fetch_price("005930", "KR")
        result2 = cached_broker.fetch_price("005930", "KR")

        stats = cached_broker.get_cache_stats()
        assert stats['cache_enabled'] is False
```

---

## 📖 5. 사용 예제

### 5.0 환경 설정

**중요**: Python 스크립트 실행 전 반드시 가상환경을 생성하고 활성화해야 합니다.

```bash
# 가상환경 생성 (.venv는 필수 이름)
python -m venv .venv

# 가상환경 활성화
source .venv/bin/activate  # macOS/Linux
# .venv\Scripts\activate   # Windows

# 패키지 설치 (editable 모드)
pip install -e .

# 개발 의존성 포함 설치
pip install -e ".[dev]"
```

**환경 변수 설정** (OS 환경변수 사용, .env 파일 사용 안 함):
```bash
# ~/.zshrc 또는 ~/.bashrc에 추가
export KOREA_INVESTMENT_API_KEY="your-api-key"
export KOREA_INVESTMENT_API_SECRET="your-api-secret"
export KOREA_INVESTMENT_ACCOUNT_NO="12345678-01"
```

### 5.1 기본 사용법

```python
from korea_investment_stock import KoreaInvestment, CachedKoreaInvestment
import os

api_key = os.getenv('KOREA_INVESTMENT_API_KEY')
api_secret = os.getenv('KOREA_INVESTMENT_API_SECRET')
acc_no = os.getenv('KOREA_INVESTMENT_ACCOUNT_NO')

# 기본 broker 생성
broker = KoreaInvestment(api_key, api_secret, acc_no, mock=True)

# 캐싱 래퍼 적용
cached_broker = CachedKoreaInvestment(broker)

# 사용 (기존과 동일)
result = cached_broker.fetch_price("005930", "KR")
print(f"삼성전자 현재가: {result['output1']['stck_prpr']}원")

# 캐시 통계
stats = cached_broker.get_cache_stats()
print(f"캐시 히트율: {stats['hit_rate']}")
```

### 5.2 TTL 커스터마이징

```python
# 실시간 트레이딩: 짧은 TTL
cached_broker = CachedKoreaInvestment(
    broker,
    price_ttl=1,        # 1초
    stock_info_ttl=60   # 1분
)

# 백테스팅/분석: 긴 TTL
cached_broker = CachedKoreaInvestment(
    broker,
    price_ttl=60,       # 1분
    stock_info_ttl=3600 # 1시간
)
```

### 5.3 컨텍스트 매니저

```python
with CachedKoreaInvestment(broker) as cached_broker:
    for symbol in ["005930", "000660", "035720"]:
        result = cached_broker.fetch_price(symbol, "KR")
        print(f"{symbol}: {result['output1']['stck_prpr']}원")
# with 블록 종료 시 캐시 자동 정리
```

### 5.4 캐시 제어

```python
cached_broker = CachedKoreaInvestment(broker)

# 가격 조회
result = cached_broker.fetch_price("005930", "KR")

# 캐시 무효화 (장 시작/마감 시)
cached_broker.invalidate_cache()

# 캐시 통계
stats = cached_broker.get_cache_stats()
print(f"""
캐시 크기: {stats['cache_size']}
히트: {stats['hits']}
미스: {stats['misses']}
히트율: {stats['hit_rate']}
""")
```

---

## ✅ 구현 체크리스트

- [ ] `CacheManager` 클래스 구현
- [ ] `CachedKoreaInvestment` 래퍼 클래스 구현
- [ ] `__init__.py` 업데이트
- [ ] 단위 테스트 작성
- [ ] 통합 테스트 작성
- [ ] 기존 테스트 통과 확인
- [ ] 사용 예제 작성
