# CHANGELOG

## [1.1.0] - 2026-05-04

### Added — Phase 1.3 (국내주식 순위/재무)

- `Domestic.InquireVolumeRank` — 거래량순위 (FHPST01710000)
- `Domestic.InquireFluctuation` — 등락률 순위 (FHPST01700000)
- `Domestic.InquireMarketCap` — 시가총액 상위 (FHPST01740000)
- `Domestic.InquireDividendRate` — 배당률 상위 (HHKDB13470100)
- `Domestic.InquireFinancialRatio` — 재무비율 (FHKST66430300)
- `Domestic.InquireIncomeStatement` — 손익계산서 (FHKST66430200)
- `Domestic.InquireBalanceSheet` — 대차대조표 (FHKST66430100)
- `Domestic.InquireProfitRatio` — 수익성비율 (FHKST66430400)
- `Domestic.InquireGrowthRatio` — 성장성비율 (FHKST66430800)
- examples: `domestic_ranking`, `domestic_financial`

### Notes

- ranking 메서드의 query parameter naming 이 inconsistent (거래량순위만 대문자 `FID_*`, 나머지 소문자 `fid_*`) — KIS docs 그대로 노출
- 거래량순위 응답의 최상위 키가 대문자 `Output` (KIS docs 명시), 다른 ranking/finance 는 소문자 `output`/`output1`
- 손익계산서 / 대차대조표 의 일부 필드 (감가상각비, 영업외 수익/비용 등) 는 출력되지 않을 시 `"99.99"` 반환 — string 필드로 노출, caller 가 처리

## [1.0.0] - 2026-05-04

> Go 라이브러리 첫 stable release. Phase 1.1 (인프라+Config) + Phase 1.2 (국내 시세/심볼/차트) 통합.
>
> **Namespace transition**: Python 시대 (`v0.6.0` ~ `v0.19.0`) 와 명확한 분리를 위해 Go 라이브러리는 `v1.0.0` 부터 publish. 이전 Go pre-release 태그 (`v0.1.0`, `v0.2.0`) 는 삭제됨.

### Added — Phase 1.1 (인프라 + Config)

- **3 진입점**: `kis.NewClient(apiKey, apiSecret, accountNo, ...opts)`, `kis.NewClientFromEnv()`, `kis.NewClientFromYAML(path)`
- **10 functional options**: `WithBaseURL`, `WithRetries`, `WithRateLimit`, `WithHTTPClient`, `WithTokenStorage`, `WithMasterCacheDir`, `WithLogger`, `WithTimeout`, `WithUserAgent`, `WithRedisURL`
- **`internal/httpclient`**: `go-resty/resty/v2` wrapper. tr_id 헤더, 토큰 자동 재발급 (`EGW00123` 만료 감지), 5xx/429 retry with exponential backoff, 응답 정규화 (`rt_cd`/`msg_cd`/`msg1` + `output`/`output1`/`output2`)
- **`internal/ratelimit`**: token bucket rate limiter (default 15 calls/sec)
- **`internal/token`**: token storage abstraction. `FileStorage` (JSON file in `~/Library/Caches/kis/`), `RedisStorage` (TTL-aware), 자동 만료 감지 + refresh
- **`internal/mastercache`**: 디스크 file cache (default TTL 7일). atomic write (temp + rename), stale fallback on fetch error
- examples: `basic_example`, `env_config_example`, `yaml_config_example`

### Added — Phase 1.2 (국내주식 시세/심볼/차트)

- `Domestic.InquirePrice` — 주식현재가 시세 (FHKST01010100)
- `Domestic.SearchInfo` — 상품기본조회 (CTPF1604R)
- `Domestic.SearchStockInfo` — 주식기본조회 (CTPF1002R)
- `Domestic.InquireDailyItemChartPrice` — 국내주식기간별시세 일/주/월/년 (FHKST03010100)
- `Domestic.InquireTimeItemChartPrice` — 주식당일분봉조회 (FHKST03010200)
- `Domestic.FetchKospiSymbols` / `FetchKosdaqSymbols` — KRX 종목 마스터 (cp949+fwf 파서, mastercache 디스크 캐시)
- `internal/krxmaster` 패키지 — KRX 마스터 파일 파싱
- examples: `domestic_price`, `domestic_chart`, `kospi_symbols`

### Conventions

- **호출 스타일**: `client.Domestic.InquirePrice(ctx, "005930")` — 한투 API path 의 마지막 segment 를 PascalCase 로 1:1 매핑 (Style A)
- **응답 typed struct**: 한투 API 약어 그대로 PascalCase 변환 (`stck_prpr` → `StckPrpr`), 인라인 한국어 코멘트, JSON 태그 한투 원본 preserve
- **타입 매핑**: 가격/액면가 = `decimal.Decimal` (bare tag), 수량/백만원 단위 = `int64,string`, 비율/PER/PBR = `float64,string`, 코드/Y-N/날짜 = `string`
- **Params struct**: 차트류 메서드는 `XxxParams` struct (zero-value default — `Period=""→"D"`, `OriginalPrice false→수정주가`)
- **Output1+Output2**: 차트는 KIS 키 verbatim 노출

### Removed

- `kis.APIError` 타입 + sentinel errors (`ErrTokenExpired`, `ErrRateLimited`, `ErrNotFound`, `ErrUnauthorized`) — 미구현 dead code 정리. 에러는 `error.Error()` 메시지의 `msg_cd`/`msg1` 로 구분 (typed error 는 추후 사용자 demand 시 재도입 검토)

### Notes

- KRX 마스터의 `fwfLen` plan 값 (228 / 222) 실제 (227 / 221) 로 수정 — 첫 행 fund-record 회피 위해 일반 주권 6자리 코드 grep 필터 testdata 사용
- `DailyChartSummary` 에 `itewhol_loan_rmnd_ratem` (전체 융자 잔고 비율) 필드 추가

## [Unreleased]

> 본 repo 는 Go 로 마이그레이션 중입니다. Python 신규 기능 추가는 중단되었습니다.

## [v0.19.0] — 2026-05-03

### Deprecation Notice

- **이 버전이 Python 라이브러리의 마지막 기능 release 입니다.** 이후 Go 모듈로 대체됩니다.
- 마지막 Python 커밋은 `python-final` 태그로 영구 보존됩니다.
- PyPI 패키지 자체는 archive 하지 않으며, critical security fix 만 v0.19.x patch 로 받을 수 있습니다.
- 신규 사용자는 Go 모듈 (`github.com/kenshin579/korea-investment-stock`) 을 사용해주세요.

상세 내용: [Phase 0 design spec](docs/superpowers/specs/2026-05-03-korea-investment-go-migration-design.md)

### Added (이 버전에 포함된 기능 — 기존 [Unreleased] 항목 유지)

#### API 확장 Phase 1: 15개 GET API 추가 (#124)

차트, 시세 순위, 재무제표, 배당/업종 4개 카테고리에 걸쳐 15개 API를 추가합니다.

**차트 데이터 API (3개)**:
- `fetch_domestic_chart()` - 국내주식 기간별시세 (일/주/월/년봉)
- `fetch_domestic_minute_chart()` - 주식당일분봉조회
- `fetch_overseas_chart()` - 해외주식 기간별시세

**시세 순위 API (4개)**:
- `fetch_volume_ranking()` - 거래량순위
- `fetch_change_rate_ranking()` - 등락률 순위
- `fetch_market_cap_ranking()` - 시가총액 상위
- `fetch_overseas_change_rate_ranking()` - 해외주식 상승율/하락율

**재무제표 API (5개)**:
- `fetch_financial_ratio()` - 재무비율 (ROE, EPS, BPS 등)
- `fetch_income_statement()` - 손익계산서
- `fetch_balance_sheet()` - 대차대조표
- `fetch_profitability_ratio()` - 수익성비율
- `fetch_growth_ratio()` - 성장성비율

**배당 + 업종 API (3개)**:
- `fetch_dividend_ranking()` - 배당률 상위
- `fetch_industry_index()` - 업종 현재지수
- `fetch_industry_category_price()` - 업종 구분별전체시세

모든 API에 Cache 래퍼와 Rate Limit 래퍼가 포함됩니다.

#### 시장별 투자자매매동향(시세) API 추가 (#120)

시장별 투자자 유형(외국인, 개인, 기관 등)의 매매 현황을 시간대별로 조회합니다.
한국투자 HTS [0403] 시장별 시간동향 화면과 동일한 기능입니다.

```python
from korea_investment_stock import KoreaInvestment

broker = KoreaInvestment()

# 코스피 종합 투자자 매매동향
result = broker.fetch_investor_trend_by_market("KSP", "0001")

# 코스닥 종합 투자자 매매동향
result = broker.fetch_investor_trend_by_market("KSQ", "1001")

# ETF 전체 투자자 매매동향
result = broker.fetch_investor_trend_by_market("ETF", "T000")

if result['rt_cd'] == '0':
    for item in result['output']:
        print(f"외국인 순매수: {item['frgn_ntby_qty']}주")
        print(f"기관 순매수: {item['orgn_ntby_qty']}주")
        print(f"개인 순매수: {item['prsn_ntby_qty']}주")
```

**주요 기능**:
- 시장별(코스피, 코스닥, ETF 등) 투자자 매매동향 조회
- 투자자 유형별(외국인, 개인, 기관, 증권, 투신, 사모펀드, 은행, 보험, 기금 등) 순매수 수량/금액 조회
- 자동 토큰 재발급 지원

**시장 코드 상수 추가**:
- `MARKET_INVESTOR_TREND_CODE`: 시장 코드 (KSP, KSQ, ETF 등)
- `SECTOR_CODE`: 업종 코드 (0001, 1001, T000 등)

#### 종목별 투자자매매동향(일별) API 추가 (#114)

특정 종목의 날짜별 외국인/기관/개인 매수매도 현황을 조회합니다.
한국투자 HTS [0416] 종목별 일별동향 화면과 동일한 기능입니다.

```python
from korea_investment_stock import KoreaInvestment

broker = KoreaInvestment()

# 삼성전자 어제 투자자 매매동향
from datetime import datetime, timedelta
yesterday = (datetime.now() - timedelta(days=1)).strftime("%Y%m%d")

result = broker.fetch_investor_trading_by_stock_daily("005930", yesterday)

if result['rt_cd'] == '0':
    for day in result['output2']:
        print(f"날짜: {day['stck_bsop_date']}")
        print(f"외국인 순매수: {day['frgn_ntby_qty']}주 ({day['frgn_ntby_tr_pbmn']}백만원)")
        print(f"기관 순매수: {day['orgn_ntby_qty']}주 ({day['orgn_ntby_tr_pbmn']}백만원)")
        print(f"개인 순매수: {day['prsn_ntby_qty']}주")
```

**주요 기능**:
- 외국인/기관/개인 순매수 수량 및 금액 조회
- 기관 세부 분류 (증권, 투자신탁, 사모펀드, 은행, 보험 등)
- 캐시 및 Rate Limit 래퍼 지원
- 자동 토큰 재발급 지원

**캐시 전략**:
- 과거 날짜 데이터: 1시간 캐시 (확정된 데이터)
- 당일 데이터: 5초 캐시 (장중 실시간 변동)

#### API 호출 중 토큰 만료 시 자동 재발급 기능 (#109)

장시간 실행되는 배치 작업 중 토큰이 만료되어도 자동으로 재발급되어 중단 없이 처리됩니다.

**동작 방식**:
- API 응답에서 토큰 만료 에러 감지 (`"기간이 만료된 token 입니다"`)
- 자동으로 `issue_access_token(force=True)` 호출 후 재시도
- 사용자 코드 수정 불필요 (투명한 처리)

**적용된 API 메서드**:
- `fetch_domestic_price()`
- `fetch_price_detail_oversea()`
- `fetch_stock_info()`
- `fetch_search_stock_info()`
- `fetch_ipo_schedule()`

**새로운 기능**:
- `issue_access_token(force=True)` - 저장소 상태와 무관하게 강제 토큰 재발급

**로깅**:
토큰 재발급 이벤트는 INFO 레벨로 로깅됩니다:
```python
import logging
logging.basicConfig(level=logging.INFO)
# LOG: 토큰 만료 감지, 재발급 시도...
```

#### 해외 주식 마스터 파일 다운로드 기능 (#102)

**해외 11개 거래소 종목 코드 다운로드 지원**:

```python
from korea_investment_stock import KoreaInvestment, OVERSEAS_MARKETS

broker = KoreaInvestment(api_key, api_secret, acc_no)

# 나스닥 종목 조회
nasdaq = broker.fetch_nasdaq_symbols()

# 뉴욕증권거래소 종목 조회
nyse = broker.fetch_nyse_symbols()

# 홍콩 종목 조회
hk = broker.fetch_overseas_symbols("hks")

# 지원 시장 확인
print(OVERSEAS_MARKETS)
# {'nas': '나스닥', 'nys': '뉴욕', 'ams': '아멕스', 'shs': '상해', ...}
```

**지원 거래소 (11개)**:
| 코드 | 거래소 |
|------|--------|
| `nas` | 나스닥 (NASDAQ) |
| `nys` | 뉴욕 (NYSE) |
| `ams` | 아멕스 (AMEX) |
| `shs` | 상해 |
| `shi` | 상해지수 |
| `szs` | 심천 |
| `szi` | 심천지수 |
| `tse` | 도쿄 |
| `hks` | 홍콩 |
| `hnx` | 하노이 |
| `hsx` | 호치민 |

**새로운 메서드**:
- `fetch_overseas_symbols(market)` - 해외 종목 코드 조회
- `fetch_nasdaq_symbols()` - 나스닥 편의 메서드
- `fetch_nyse_symbols()` - 뉴욕 편의 메서드
- `fetch_amex_symbols()` - 아멕스 편의 메서드

**새로운 상수**:
- `OVERSEAS_MARKETS` - 지원 시장 코드 (11개)
- `OVERSEAS_COLUMNS` - 컬럼명 목록 (24개)

**Wrapper 호환**:
- `CachedKoreaInvestment` 지원
- `RateLimitedKoreaInvestment` 지원

#### Testcontainers 도입 - Redis 통합 테스트 (#92)

**실제 Docker 컨테이너 기반 통합 테스트 환경 구축**:

- `testcontainers>=4.0.0` 의존성 추가
- pytest marker로 테스트 유형 구분 (`unit`, `integration`)
- Redis 통합 테스트 7개 추가:
  - 토큰 저장/로드/삭제
  - 다중 스레드 연결 풀
  - 실제 TTL 만료 확인
  - 다중 데이터베이스 격리

**테스트 실행**:
```bash
# 단위 테스트만 (Docker 불필요)
pytest -m "not integration"

# 통합 테스트만 (Docker 필요)
pytest -m integration

# 전체 테스트
pytest
```

**fakeredis와의 공존**:
- 기존 fakeredis 단위 테스트 유지 (빠른 피드백)
- testcontainers 통합 테스트 추가 (실제 환경 검증)
- Docker 미설치 시 통합 테스트 자동 스킵

#### Hybrid Configuration System (v1.1.0) (#76)

**5단계 설정 우선순위 시스템**:

1. 생성자 파라미터 (최고 우선순위)
2. `config` 객체
3. `config_file` 파라미터
4. 환경 변수
5. 기본 config 파일 (`~/.config/kis/config.yaml`)

**새로운 파라미터**:
```python
broker = KoreaInvestment(
    config=Config.from_yaml("config.yaml"),  # Config 객체 주입
    config_file="./my_config.yaml",          # YAML 파일 경로
)
```

**기본 config 파일 자동 탐색**:
```yaml
# ~/.config/kis/config.yaml
api_key: your-api-key
api_secret: your-api-secret
acc_no: "12345678-01"
```

**혼합 사용 (부분 override)**:
```python
config = Config.from_yaml("~/.config/kis/config.yaml")
broker = KoreaInvestment(
    config=config,
    api_key="override-key"  # config보다 우선
)
```

**하위 호환성**: 기존 코드 100% 호환
```python
# 기존 방식 모두 동작
broker = KoreaInvestment(api_key, api_secret, acc_no)  # 생성자 파라미터
broker = KoreaInvestment()  # 환경 변수 자동 감지
```

### Changed

#### fetch_stock_info, fetch_search_stock_info 개선 (#94)

**Breaking Change: 인자 변경**

```python
# 변경 전
broker.fetch_stock_info("005930", market="KR")
broker.fetch_search_stock_info("005930", market="KR")

# 변경 후
broker.fetch_stock_info("005930", country_code="KR")
broker.fetch_search_stock_info("005930", country_code="KR")  # KR만 지원, 그 외 ValueError
```

**주요 변경 내용**:

- `fetch_stock_info` 인자: `market` → `country_code`
- `fetch_search_stock_info` 인자: `market` → `country_code` (KR만 지원, 그 외 ValueError)
- API 문서 기반 상세 docstring 추가
- 반환 타입 힌트 `-> dict` 추가

**상수 변경**:

- `MARKET_TYPE_MAP` → `PRDT_TYPE_CD_BY_COUNTRY`로 이름 변경
- `PRDT_TYPE_CD` 상수 참조 사용으로 코드 품질 향상
- `OVRS_EXCG_CD` 키 형태 변경 (NASD:NASD 패턴)

**호환성 노트**:

- `fetch_stock_info`: 위치 인자 사용 시 호환 (예: `broker.fetch_stock_info("005930", "KR")`)
- `fetch_stock_info`: 키워드 인자 `market=` 사용 시 `country_code=`로 변경 필요
- `fetch_search_stock_info`: 키워드 인자 `market=` 사용 시 `country_code=`로 변경 필요
- `fetch_search_stock_info`: KR 외 country_code 사용 시 ValueError 발생

#### fetch_price_detail_oversea 리팩토링 (#90)

**인자명 변경**: `market` → `country_code`

```python
# v1.0.x (Before)
broker.fetch_price_detail_oversea("AAPL", market="US")

# v1.1.0 (After)
broker.fetch_price_detail_oversea("AAPL")  # 기본값 "US"
broker.fetch_price_detail_oversea("AAPL", country_code="US")
broker.fetch_price_detail_oversea("9988", country_code="HK")  # 홍콩 알리바바
broker.fetch_price_detail_oversea("7203", country_code="JP")  # 일본 토요타
```

**지원 국가**:
- `"US"`: 미국 (NYSE, NASDAQ, AMEX + 주간거래)
- `"HK"`: 홍콩
- `"JP"`: 일본
- `"CN"`: 중국 (상하이, 심천)
- `"VN"`: 베트남 (호치민, 하노이)

**상수 변경**:
- `EXCD` 키 변경: `"NYSE"` → `"NYS"`, `"NASDAQ"` → `"NAS"` 등
- `EXCD_BY_COUNTRY` 신규 추가: 국가별 거래소 코드 매핑

- **Project Structure**: Reorganized package into feature-based modules (#52)
  - Created `cache/` module for caching functionality
  - Created `token_storage/` module for token storage implementations
  - Moved test files to co-locate with implementation files (co-located tests)
  - Removed `tests/` directory in favor of feature-specific test files
  - All existing import paths remain compatible (backward compatible)
  - Updated version to 0.7.0

## [0.8.0] - 2025-01-XX (Breaking Changes) ⚠️

### ⚠️ BREAKING CHANGES

#### Mock 모드 완전 제거 (#55)

**제거된 기능**: 모의투자 서버 지원 (`mock` 파라미터)

**변경 사항**:

1. **생성자 시그니처 변경**
```python
# v0.7.x (Before)
broker = KoreaInvestment(api_key, api_secret, acc_no, mock=True)

# v0.8.0 (After)
broker = KoreaInvestment(api_key, api_secret, acc_no)
```

2. **제거된 메서드**
- `set_base_url(mock: bool)` 메서드 제거
- 실전 서버 URL 고정: `https://openapi.koreainvestment.com:9443`

3. **제거된 검증**
- `fetch_ipo_schedule()`: 모의투자 검증 로직 제거

**마이그레이션 가이드**:
```python
# Before (v0.7.x)
broker = KoreaInvestment(
    api_key="YOUR_API_KEY",
    api_secret="YOUR_API_SECRET",
    acc_no="12345678-01",
    mock=True  # 또는 mock=False
)

# After (v0.8.0)
broker = KoreaInvestment(
    api_key="YOUR_API_KEY",
    api_secret="YOUR_API_SECRET",
    acc_no="12345678-01"
)
```

**주의사항**:
- ⚠️ v0.8.0부터는 **실전 계좌만 지원**됩니다
- ⚠️ 테스트 환경이 필요한 경우 `unittest.mock` 사용 권장

**단위 테스트 예제**:
```python
from unittest.mock import patch

@patch('korea_investment_stock.requests.get')
def test_fetch_price(mock_get):
    mock_get.return_value.json.return_value = {
        'rt_cd': '0',
        'output1': {'stck_prpr': '70000'}
    }
    broker = KoreaInvestment(api_key, api_secret, acc_no)
    result = broker.fetch_price("005930", "KR")
    assert result['output1']['stck_prpr'] == '70000'
```

### Added

#### API Rate Limiting (#67)

**New Feature**: Automatic rate limiting to manage Korea Investment API's 20 calls/second limit.

**Components**:
- `RateLimiter`: Thread-safe rate limiter using token bucket algorithm
- `RateLimitedKoreaInvestment`: Wrapper class for automatic rate limiting

**Usage**:
```python
from korea_investment_stock import KoreaInvestment, RateLimitedKoreaInvestment

# Create base broker
broker = KoreaInvestment(api_key, api_secret, acc_no)

# Wrap with rate limiting (15 calls/second - conservative)
rate_limited = RateLimitedKoreaInvestment(broker, calls_per_second=15)

# Use as normal - rate limiting applied automatically
result = rate_limited.fetch_price("005930", "KR")
```

**Features**:
- ✅ Thread-safe using `threading.Lock`
- ✅ Default: 15 calls/second (conservative margin)
- ✅ Dynamic rate adjustment at runtime
- ✅ Statistics tracking (total_calls, min_interval)
- ✅ Context manager support
- ✅ Zero changes to existing `KoreaInvestment` class
- ✅ Works with `CachedKoreaInvestment` (recommended combination)

**Benefits**:
- Prevents API rate limit errors
- `examples/stress_test.py` now achieves 100% success (500 API calls)
- Batch processing of stocks is safe and reliable
- Opt-in design: users choose when to enable

**See Also**:
- Implementation guide: `docs/start/1_api_limit_implementation.md`
- PRD: `docs/start/1_api_limit_prd.md`
- CLAUDE.md: "API Rate Limiting" section

### Changed
- 실전 서버로 통일되어 모든 API 일관되게 지원
- 코드베이스 간소화 (mock 관련 로직 제거)
- `examples/stress_test.py` updated to use `RateLimitedKoreaInvestment`

### Removed
- `mock` 파라미터 (Breaking)
- `set_base_url()` 메서드 (Breaking)
- `self.mock` 인스턴스 변수
- IPO Schedule API의 모의투자 검증 로직

## [0.6.0] - 2025-01-19 (Breaking Changes) ⚠️

### 🎯 Major Simplification (#40)
**Philosophy Change**: Transformed from feature-rich library to **pure API wrapper**

This version removes all advanced features to focus on being a thin, reliable wrapper around the Korea Investment Securities OpenAPI. Users who need rate limiting, caching, batch processing, or monitoring should implement these features themselves according to their specific needs.

### ⚠️ BREAKING CHANGES

#### Removed Features (~6,000+ lines of code removed)
- **Rate Limiting System**: Removed EnhancedRateLimiter, BackoffStrategy, Circuit Breaker
  - Users should implement their own rate limiting if needed
- **Caching System**: Removed TTL cache, cache decorators, cache statistics
  - Users should implement their own caching strategy
- **Batch Processing**: Removed batch methods and dynamic batch controller
  - Use loops with `fetch_price()` instead of `fetch_price_list()`
- **Monitoring & Visualization**: Removed stats collection, Plotly dashboards, HTML reports
  - Users should implement their own monitoring
- **Error Recovery**: Removed automatic retry decorators and error recovery system
  - Users should handle errors according to their needs
- **Legacy Module**: Removed deprecated code and unused features

#### API Changes
- **Removed Methods**:
  - `fetch_price_list()` → Use loop with `fetch_price(symbol, market)`
  - `fetch_stock_info_list()` → Use loop with `fetch_stock_info(symbol, market)`
  - `fetch_price_list_with_batch()` → Use loop with `fetch_price()`
  - `fetch_price_list_with_dynamic_batch()` → Use loop with `fetch_price()`
  - All batch processing methods
  - All caching-related methods
  - All statistics and monitoring methods

- **Private → Public Methods** (now part of public API):
  - `__fetch_price()` → `fetch_price(symbol, market)`
  - `__fetch_stock_info()` → `fetch_stock_info(symbol, market)`
  - `__fetch_domestic_price()` → `fetch_domestic_price(market_code, symbol)`
  - `__fetch_etf_domestic_price()` → `fetch_etf_domestic_price(market_code, symbol)`
  - `__fetch_price_detail_oversea()` → `fetch_price_detail_oversea(symbol, market)`

#### Simplified Dependencies
- **Removed**: `websockets`, `pycryptodome`, `crypto`
- **Kept**: `requests`, `pandas` (minimal dependencies)

### ✅ What Remains
- ✅ Stock price queries (domestic & US)
- ✅ Stock information queries
- ✅ IPO schedule queries
- ✅ Unified interface for KR/US stocks via `fetch_price(symbol, market)`
- ✅ Basic error responses from API
- ✅ Context manager support
- ✅ Thread pool executor (basic concurrency)

### 📦 Migration Guide

#### Before (v0.5.0):
```python
# Batch query with automatic rate limiting, caching, retry
stocks = [("005930", "KR"), ("AAPL", "US")]
results = broker.fetch_price_list(stocks)
```

#### After (v0.6.0):
```python
# Simple loop - implement your own rate limiting if needed
stocks = [("005930", "KR"), ("AAPL", "US")]
results = []
for symbol, market in stocks:
    result = broker.fetch_price(symbol, market)
    results.append(result)
    # Add your own rate limiting, caching, retry logic here if needed
```

### 📈 Code Reduction
- Main file: 1,941 → 1,011 lines (48% reduction)
- Total deletion: ~6,000+ lines
- Module count: 15 → 1 (core module only)
- Test files: 18 → 4 (only integration tests remain)

### 🎯 Why This Change?
- **Simplicity**: Focus on doing one thing well - wrapping the API
- **Flexibility**: Users implement features their way
- **Maintainability**: Less code = fewer bugs
- **Transparency**: Pure wrapper with no magic

### 📚 Documentation Updates
- Updated README.md to reflect simple API wrapper approach
- Updated CLAUDE.md to remove advanced architecture details
- Updated examples to show simple usage patterns
- Added `basic_example.py` for simple use cases

## [Unreleased] - 2025-01-14

### 🚀 추가된 기능

#### 미국 주식 통합 지원 (#33) ✨
- **통합 인터페이스**: `fetch_price_list()`로 국내/미국 주식 모두 조회 가능
  - 기존: 국내 주식만 지원
  - 개선: `[("005930", "KR"), ("AAPL", "US")]` 혼합 조회 가능
- **자동 거래소 검색**: NASDAQ, NYSE, AMEX 순으로 자동 탐색
- **추가 재무 정보**: 미국 주식의 경우 PER, PBR, EPS, BPS, 52주 최고/최저가 등 제공
- **향상된 에러 처리**: 거래소별 심볼 검색 실패 시 명확한 에러 메시지
- **캐시 통합**: 미국 주식도 5분 TTL 캐시 적용으로 성능 향상

### 🔧 개선사항

#### API 메서드 캡슐화
- `fetch_etf_domestic_price()` → `__fetch_etf_domestic_price()` (private)
- `fetch_domestic_price()` → `__fetch_domestic_price()` (private)
- 사용자는 통합 인터페이스 `fetch_price_list()` 사용 권장

### ⚠️ 주의사항
- 미국 주식은 **실전투자 계정에서만** 조회 가능 (모의투자 미지원)
- 미국 주식은 실시간 무료시세 제공 (나스닥 마켓센터 기준)

## [Unreleased] - 2024-12-28

### 🏗️ 구조 개선

#### 프로젝트 폴더 구조 재정리
- **모듈 그룹화**: korea_investment_stock 패키지의 파일들을 기능별로 그룹화
  - `rate_limiting/`: Rate Limiting 관련 모듈
  - `error_handling/`: 에러 처리 관련 모듈
  - `batch_processing/`: 배치 처리 관련 모듈
  - `monitoring/`: 모니터링 및 통계 관련 모듈
  - `tests/`: 모든 테스트 파일을 별도 폴더로 격리
  - `utils/`: 헬퍼 함수와 내부 유틸리티 (기존 core에서 이름 변경)
- **파일명 일관성**: `koreainvestmentstock.py` → `korea_investment_stock.py`로 변경
- **메인 모듈 위치 변경**: Python 표준에 맞게 `korea_investment_stock.py`를 패키지 루트로 이동
- **Import 구조 개선**: 각 모듈별 `__init__.py`에서 주요 클래스/함수 export
- **하위 호환성 유지**: 공개 API는 변경 없이 내부 구조만 개선

### 🚀 추가된 기능

#### Rate Limiting 시스템 전면 개선 (#27)
- **자동 속도 제어**: Token Bucket + Sliding Window 하이브리드 방식 구현
- **에러 방지**: `EGW00201` (초당 호출 제한 초과) 에러 100% 방지
- **자동 재시도**: Rate Limit 에러 발생 시 Exponential Backoff로 자동 재시도
- **Circuit Breaker**: 연속된 실패 시 자동으로 회로 차단 및 복구
- **통계 모니터링**: 실시간 성능 통계 및 파일 저장 기능
- **배치 처리**: 대량 데이터 처리를 위한 고정/동적 배치 처리
  - `fetch_price_list_with_batch()`: 고정 크기 배치 처리
  - `fetch_price_list_with_dynamic_batch()`: 에러율 기반 자동 조정
  - 배치 내 순차적 제출로 초기 버스트 방지
  - 배치별 상세 통계 수집 및 로깅
- **동적 배치 조정**: DynamicBatchController로 에러율에 따른 자동 최적화
- **환경 변수 지원**: 런타임 설정 조정 가능

### 🔧 개선사항

#### ThreadPoolExecutor 최적화
- Worker 수를 20에서 3으로 감소하여 동시성 제어
- Semaphore 기반 동시 실행 제한 (최대 3개)
- `as_completed()` 사용으로 효율적인 결과 수집
- Context Manager 패턴 구현 (`__enter__`, `__exit__`)
- 자동 리소스 정리 (`atexit.register`)

#### 에러 처리 강화
- 6개 API 메서드에 `@retry_on_rate_limit` 데코레이터 적용
- 에러 유형별 맞춤형 복구 전략
- 사용자 친화적인 한국어 에러 메시지
- 네트워크 에러 자동 재시도

### 📊 성능 개선
- **안정적인 처리량**: 10-12 TPS 유지 (API 한계의 60%)
- **에러율**: 0% 달성 (목표 <1%)
- **100개 종목 조회**: 8.35초, 0 에러
- **장시간 안정성**: 30초 테스트 313 호출, 0 에러

### 📚 문서화
- README.md에 Rate Limiting 섹션 추가
- 상세한 사용 예제 제공 (`examples/rate_limiting_example.py`)
- 모범 사례 및 권장 설정 안내

### 🔄 하위 호환성
- 기존 API 인터페이스 완전 유지
- 기본 동작은 변경 없음
- 새로운 기능은 옵트인 방식

### 🗑️ 제거된 기능
- WebSocket 관련 코드 제거 (더 이상 사용하지 않음)
- 불필요한 레거시 메서드 제거

### 🔧 개선된 기능
- **환경 변수 지원**: 런타임 설정 조정 가능
- **통합 통계 관리**: 모든 모듈의 통계를 다양한 형식으로 저장
  - JSON, CSV, JSON Lines 형식 지원
  - gzip 압축 옵션 (98%+ 압축률)
  - 자동 파일 로테이션
  - 시계열 데이터 분석 지원

## [이전 버전]

(이전 버전 기록은 향후 추가 예정) 