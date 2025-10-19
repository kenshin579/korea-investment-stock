# PRD: Korea Investment Stock 단순화 프로젝트 (Issue #40)

> **관련 문서**:
> - [Implementation Guide](implementation.md) - 구현 상세 가이드
> - [TODO Checklist](todo.md) - 구현 체크리스트

---

## 📋 Executive Summary

### 프로젝트 목표
Korea Investment Stock 라이브러리를 **순수한 API Wrapper**로 단순화하여 유지보수성을 향상시키고, 사용자가 필요에 따라 커스터마이징 가능한 구조로 개선합니다.

### 핵심 변경사항
- **제거**: Rate limiting, Caching, Visualization, Batch processing, Monitoring, Error recovery 시스템
- **단순화**: Private 메서드를 Public으로 변경하여 사용자가 직접 제어 가능하도록 변경
- **목표**: 한국투자증권 OpenAPI의 얇은 wrapper 역할만 수행

### 기대효과
- 코드 라인 수 ~60% 감소 (1,941 → ~800 lines)
- 의존성 최소화 (requests, pandas, websockets, pycryptodome만 유지)
- 사용자가 필요한 기능을 직접 구현 가능한 유연성 제공

---

## 🎯 Background & Context

### 현재 아키텍처 문제점

**1. 과도한 기능 집중**
- 단순 API Wrapper에 Rate Limiter, Cache, Monitoring 등 너무 많은 기능 포함
- 라이브러리 복잡도 증가 → 유지보수 부담 증가
- 사용자가 원하지 않는 기능도 강제 포함

**2. 사용자 제어권 부족**
- Private 메서드(`__fetch_price`, `__fetch_stock_info` 등)로 숨겨진 핵심 기능
- 배치 처리, 캐싱 정책을 사용자가 커스터마이징 불가
- Rate limiting 로직이 내부에 고정되어 있어 조정 어려움

**3. 불필요한 복잡성**
- List 기반 메서드 (`fetch_price_list`, `fetch_stock_info_list`)가 ThreadPoolExecutor 사용
- Dynamic batch controller가 자동으로 배치 사이즈/딜레이 조정
- 사용자가 간단한 단일 조회만 원하는 경우에도 복잡한 인프라 로딩

### 목표 상태 (Target State)

```python
# BEFORE: 복잡한 사용 패턴
with KoreaInvestment(api_key, secret, acc_no) as broker:
    # 내부적으로 rate limiter, cache, thread pool 등 초기화
    prices = broker.fetch_price_list(stock_list)  # 자동 배치, 캐싱, 재시도
    broker.save_monitoring_dashboard("dashboard.html")

# AFTER: 단순하고 명확한 패턴
broker = KoreaInvestment(api_key, secret, acc_no)
price = broker.fetch_price("005930", "KR")  # 단일 조회, 직접 제어
stock_info = broker.fetch_stock_info("AAPL", "US")

# 사용자가 필요하면 직접 구현
for symbol in stock_list:
    price = broker.fetch_price(symbol, "KR")
    time.sleep(0.1)  # 사용자가 rate limiting 제어
```

---

## 📝 Detailed Requirements

### R1: Rate Limiting 시스템 제거

**제거 대상**:
- `rate_limiting/enhanced_rate_limiter.py` (~400 lines)
- `rate_limiting/enhanced_backoff_strategy.py` (~300 lines)
- `rate_limiting/enhanced_retry_decorator.py` (~200 lines)
- `rate_limiting/__init__.py` (~50 lines)

**총 제거**: ~950 lines

**제거 대상 데코레이터**: `@retry_on_rate_limit`, `@retry_on_network_error` (13개 메서드에 적용됨)

**영향**: 사용자가 직접 rate limiting 제어 (time.sleep, semaphore 등)

---

### R2: List 기반 메서드 제거

**제거 대상 메서드 (6개)**:
- `fetch_price_list()` - 주식 리스트 가격 조회
- `fetch_price_list_with_batch()` - 배치 파라미터 지정 조회
- `fetch_price_list_with_dynamic_batch()` - Dynamic batch controller 사용
- `fetch_stock_info_list()` - 주식 정보 리스트 조회
- `fetch_search_stock_info_list()` - 주식 검색 정보 리스트 조회 (중복 정의 2개)
- `fetch_price_detail_oversea_list()` - 해외 주식 리스트 조회

**제거 대상 내부 메서드 (2개)**:
- `__execute_concurrent_requests()` (~150 lines)
- `__execute_concurrent_requests_with_cache()` (~80 lines)

**총 제거**: ~400 lines

**영향**: 사용자가 for loop로 직접 배치 조회 구현

---

### R3: Private 메서드를 Public으로 변경

**변경 대상 (8개)**:

| 현재 (Private) | 변경 후 (Public) | 설명 |
|---------------|-----------------|------|
| `__fetch_price()` | `fetch_price()` | 단일 주식 가격 조회 (국내/해외 자동 판단) |
| `__get_symbol_type()` | `get_symbol_type()` | 심볼 타입 판단 (주식/ETF) |
| `__fetch_etf_domestic_price()` | `fetch_etf_domestic_price()` | 국내 ETF 가격 조회 |
| `__fetch_domestic_price()` | `fetch_domestic_price()` | 국내 주식 가격 조회 |
| `__fetch_price_detail_oversea()` | `fetch_price_detail_oversea()` | 해외 주식 상세 가격 조회 |
| `__fetch_stock_info()` | `fetch_stock_info()` | 단일 주식 정보 조회 |
| `__fetch_search_stock_info()` | `fetch_search_stock_info()` | 단일 주식 검색 정보 조회 |
| `__handle_rate_limit_error()` | ~~삭제~~ | DEPRECATED |

**작업**: 메서드명 변경 + Docstring 추가 + 데코레이터 제거

---

### R4: Cache 시스템 제거

**제거 대상 모듈**:
- `caching/ttl_cache.py` (~500 lines)
- `caching/market_hours.py` (~100 lines)
- `caching/__init__.py` (~50 lines)

**제거 대상 메서드 (5개)**:
- `clear_cache()`, `get_cache_stats()`, `set_cache_enabled()`, `preload_cache()`

**총 제거**: ~730 lines

**영향**: 사용자가 dict, redis 등으로 직접 캐싱 구현

---

### R5: Visualization 시스템 제거

**제거 대상 모듈**:
- `visualization/plotly_visualizer.py` (~400 lines)
- `visualization/dashboard.py` (~350 lines)
- `visualization/charts.py` (~250 lines)
- `visualization/__init__.py` (~50 lines)

**제거 대상 메서드 (7개)**:
- `create_monitoring_dashboard()`, `save_monitoring_dashboard()`, `create_stats_report()`, 등

**총 제거**: ~1,200 lines

**영향**: 사용자가 Prometheus, Grafana 등 외부 도구 사용

---

### R6: Batch Processing 시스템 제거

**제거 대상 모듈**:
- `batch_processing/dynamic_batch_controller.py` (~300 lines)
- `batch_processing/__init__.py` (~30 lines)

**총 제거**: ~330 lines

---

### R7: Monitoring & Error Handling 시스템 제거

**제거 대상 모듈**:
- `monitoring/stats_manager.py` (~600 lines)
- `error_handling/error_recovery_system.py` (~500 lines)
- 각각의 `__init__.py` (~60 lines)

**총 제거**: ~1,160 lines

---

### R8: Threading 시스템 제거

**제거 대상 코드**:
- `ThreadPoolExecutor` 초기화 및 shutdown
- `Semaphore` 초기화
- Background cleanup thread (Cache 관련)

**총 제거**: ~100 lines

---

## 🔄 API Surface Changes

### Before (현재 Public API - 30+ 메서드)

**주식 조회 (List 기반) - 제거 예정**:
- ❌ `fetch_price_list(stock_list)`
- ❌ `fetch_price_list_with_batch(...)`
- ❌ `fetch_price_list_with_dynamic_batch(...)`
- ❌ `fetch_stock_info_list(stock_market_list)`
- ❌ `fetch_search_stock_info_list(stock_market_list)`
- ❌ `fetch_price_detail_oversea_list(stock_market_list)`

**Cache 관리 - 제거 예정**:
- ❌ `clear_cache(pattern)`
- ❌ `get_cache_stats()`
- ❌ `set_cache_enabled(enabled)`
- ❌ `preload_cache(symbols, market)`

**Monitoring & Visualization - 제거 예정**:
- ❌ `create_monitoring_dashboard(...)`
- ❌ `save_monitoring_dashboard(filename)`
- ❌ `create_stats_report(save_as)`
- ❌ `get_system_health_chart()`
- ❌ `get_api_usage_chart(hours)`
- ❌ `show_monitoring_dashboard()`

---

### After (단순화된 Public API - 18개 메서드)

**인증 & 설정 (5개)**:
- `issue_access_token()`
- `check_access_token()`
- `load_access_token()`
- `issue_hashkey(data)`
- `set_base_url(mock)`

**단일 주식 조회 (7개) - 🆕 Public 전환**:
- `fetch_price(symbol, market)` ← `__fetch_price`
- `fetch_domestic_price(market_code, symbol)` ← `__fetch_domestic_price`
- `fetch_etf_domestic_price(market_code, symbol)` ← `__fetch_etf_domestic_price`
- `fetch_price_detail_oversea(symbol, market)` ← `__fetch_price_detail_oversea`
- `fetch_stock_info(symbol, market)` ← `__fetch_stock_info`
- `fetch_search_stock_info(symbol, market)` ← `__fetch_search_stock_info`
- `get_symbol_type(symbol_info)` ← `__get_symbol_type`

**심볼 조회 (6개)**:
- `fetch_kospi_symbols()`, `fetch_kosdaq_symbols()`, `fetch_symbols()`
- `download_master_file(...)`, `parse_kospi_master(base_dir)`, `parse_kosdaq_master(base_dir)`

**IPO 조회 (6개)**:
- `fetch_ipo_schedule(from_date, to_date, symbol)`
- Static 메서드 5개: `parse_ipo_date_range()`, `format_ipo_date()`, `calculate_ipo_d_day()`, 등

**리소스 관리 (1개)**:
- `shutdown()` (간소화)

---

## 📚 Migration Guide

### Breaking Changes Summary

**버전**: 0.5.0 → 0.6.0  
**변경 범위**: Major breaking changes (하위 호환성 없음)

---

### Code Migration Examples

#### 1. 단일 주식 조회

```python
# ❌ BEFORE (0.5.0)
broker = KoreaInvestment(api_key, secret, acc_no)
# __fetch_price는 private이라 직접 호출 불가
prices = broker.fetch_price_list([("005930", "KR")])
price = prices[0]

# ✅ AFTER (0.6.0)
broker = KoreaInvestment(api_key, secret, acc_no)
price = broker.fetch_price("005930", "KR")  # 직접 호출 가능
```

---

#### 2. 배치 조회 (여러 주식)

```python
# ❌ BEFORE (0.5.0)
with KoreaInvestment(api_key, secret, acc_no) as broker:
    stock_list = [("005930", "KR"), ("035420", "KR"), ("000660", "KR")]
    prices = broker.fetch_price_list(stock_list)  # 자동 배치, 캐싱, Rate limiting

# ✅ AFTER (0.6.0) - Option 1: 직접 제어
broker = KoreaInvestment(api_key, secret, acc_no)
stock_list = ["005930", "035420", "000660"]
prices = []

for symbol in stock_list:
    try:
        price = broker.fetch_price(symbol, "KR")
        prices.append(price)
        time.sleep(0.1)  # Rate limiting 직접 제어
    except Exception as e:
        print(f"Error fetching {symbol}: {e}")

# ✅ AFTER (0.6.0) - Option 2: 병렬 처리 (사용자 구현)
from concurrent.futures import ThreadPoolExecutor, as_completed

def fetch_with_retry(symbol, retries=3):
    for i in range(retries):
        try:
            return broker.fetch_price(symbol, "KR")
        except Exception as e:
            if i == retries - 1:
                raise
            time.sleep(2 ** i)  # Exponential backoff

prices = []
with ThreadPoolExecutor(max_workers=3) as executor:
    futures = {executor.submit(fetch_with_retry, symbol): symbol
               for symbol in stock_list}
    
    for future in as_completed(futures):
        try:
            price = future.result()
            prices.append(price)
        except Exception as e:
            print(f"Failed: {e}")
```

---

#### 3. Cache 사용 (직접 구현)

```python
# ❌ BEFORE (0.5.0)
with KoreaInvestment(api_key, secret, acc_no) as broker:
    # 자동 캐싱 (TTL 5분)
    price1 = broker.fetch_price_list([("005930", "KR")])[0]
    price2 = broker.fetch_price_list([("005930", "KR")])[0]  # Cache hit

# ✅ AFTER (0.6.0) - 직접 캐싱 구현
from datetime import datetime, timedelta

class CachedBroker:
    def __init__(self, api_key, secret, acc_no):
        self.broker = KoreaInvestment(api_key, secret, acc_no)
        self.cache = {}
        self.cache_ttl = timedelta(minutes=5)
    
    def fetch_price(self, symbol, market="KR"):
        cache_key = f"{symbol}:{market}"
        now = datetime.now()
        
        # Cache hit check
        if cache_key in self.cache:
            cached_time, cached_price = self.cache[cache_key]
            if now - cached_time < self.cache_ttl:
                return cached_price
        
        # Cache miss
        price = self.broker.fetch_price(symbol, market)
        self.cache[cache_key] = (now, price)
        return price

cached_broker = CachedBroker(api_key, secret, acc_no)
price1 = cached_broker.fetch_price("005930", "KR")
price2 = cached_broker.fetch_price("005930", "KR")  # Cache hit
```

---

#### 4. Monitoring & Visualization

```python
# ❌ BEFORE (0.5.0)
with KoreaInvestment(api_key, secret, acc_no) as broker:
    prices = broker.fetch_price_list(stock_list)
    
    # 내장 모니터링
    broker.save_monitoring_dashboard("dashboard.html")
    stats = broker.create_stats_report()

# ✅ AFTER (0.6.0) - 직접 구현 or 외부 도구
import time
import json

broker = KoreaInvestment(api_key, secret, acc_no)
stats = {
    "total_requests": 0,
    "errors": 0,
    "start_time": time.time()
}

for symbol in stock_list:
    try:
        price = broker.fetch_price(symbol, "KR")
        stats["total_requests"] += 1
    except Exception as e:
        stats["errors"] += 1

stats["duration"] = time.time() - stats["start_time"]
stats["success_rate"] = (stats["total_requests"] - stats["errors"]) / stats["total_requests"]

with open("stats.json", "w") as f:
    json.dump(stats, f, indent=2)
```

---

#### 5. IPO 조회 (변경 없음)

```python
# ✅ BEFORE & AFTER (동일)
broker = KoreaInvestment(api_key, secret, acc_no)

# 전체 IPO 일정 조회
ipo_data = broker.fetch_ipo_schedule()

# 기간 지정 조회
ipo_data = broker.fetch_ipo_schedule(
    from_date="20250101",
    to_date="20250131"
)

# Static 메서드 사용
d_day = KoreaInvestment.calculate_ipo_d_day("20250120")
status = KoreaInvestment.get_ipo_status("20250120")
```

---

### 권장 Migration 전략

**Phase 1: 의존성 업데이트**
```bash
pip install korea-investment-stock==0.6.0
```

**Phase 2: 코드 수정**
1. `fetch_price_list()` → `fetch_price()` loop로 변경
2. Context manager 제거 (선택사항)
3. 필요시 캐싱/배치 처리 직접 구현

**Phase 3: 테스트**
1. 단일 조회 기능 테스트
2. 배치 조회 로직 검증 (사용자 구현)
3. 에러 처리 검증

**Phase 4: 배포**
- 점진적 배포 (Canary deployment 권장)
- 모니터링 강화 (외부 APM 도구 사용)

---

## ⚠️ Risk Assessment

### High Risk Areas

**1. Breaking Changes (심각도: HIGH)**
- 모든 기존 사용자 코드가 영향받음
- `fetch_price_list()` 제거 → 사용자가 직접 loop 구현 필요

**완화 전략**:
- 명확한 migration guide 제공
- CHANGELOG에 Breaking changes 명시
- Major version bump (0.5 → 0.6)

**2. 성능 저하 가능성 (심각도: MEDIUM)**
- Rate limiting 제거 → API 서버 과부하 위험
- Cache 제거 → API 호출 증가
- ThreadPool 제거 → 병렬 처리 성능 저하

**완화 전략**:
- Documentation에 rate limiting 권장사항 명시
- 사용자 구현 예시 제공 (retry, cache, batch)

**3. 기능 손실 (심각도: MEDIUM)**
- Monitoring/Visualization 제거 → 문제 진단 어려움

**완화 전략**:
- 외부 도구 사용 가이드 제공 (Prometheus, Grafana)

---

## 🎯 Success Criteria

### 정량적 지표

1. **코드 라인 수**: 60% 감소 (1,941 → ~800 lines)
2. **의존성 수**: Plotly 제거 (선택 의존성 0개)
3. **테스트 통과율**: 100% (유지된 테스트 전체 통과)
4. **Public API 수**: 18개 메서드 (명확한 API surface)

### 정성적 지표

5. **단순성**: `__init__()`이 최소한의 초기화만 수행
6. **명확성**: 모든 Public 메서드에 docstring 존재
7. **유연성**: 사용자가 rate limiting, caching 직접 제어 가능

---

## 📅 Timeline Estimate

**총 예상 시간**: 2-3일 (개발자 1명 기준)

- **Day 1**: 모듈 삭제 + 메인 모듈 수정 (7시간)
- **Day 2**: 테스트 + Example 수정 (6시간)
- **Day 3**: 문서화 + 검증 + 배포 (7시간)

---

## 📌 Notes

### 설계 원칙

1. **KISS (Keep It Simple, Stupid)**
   - 한국투자증권 API의 얇은 wrapper만 제공
   - 복잡한 기능은 사���자가 필요시 직접 구현

2. **Separation of Concerns**
   - API 통신: 라이브러리 담당
   - Rate limiting, Caching, Monitoring: 사용자 담당

3. **Principle of Least Surprise**
   - 메서드 이름이 동작을 명확히 표현
   - Private → Public 전환으로 숨겨진 기능 없음

### 참고 자료

- Issue: https://github.com/kenshin579/korea-investment-stock/issues/40
- 한국투자증권 OpenAPI 문서: https://apiportal.koreainvestment.com/
- Python API Wrapper Best Practices: https://realpython.com/api-integration-in-python/

---

## ✍️ Document History

| 버전 | 날짜 | 작성자 | 변경사항 |
|-----|------|--------|---------|
| 1.0 | 2025-10-19 | Claude Code | 초안 작성 |
| 1.1 | 2025-10-19 | Claude Code | 문서 분리 (prd.md, implementation.md, todo.md) |

---

**작성**: Claude Code  
**검토**: (To be reviewed)  
**승인**: (To be approved)
