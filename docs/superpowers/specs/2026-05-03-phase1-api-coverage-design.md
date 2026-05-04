# Phase 1 API Coverage 설계 — Python parity Go 라이브러리

- **상태**: Draft
- **작성일**: 2026-05-03
- **작성자**: kenshin579
- **연관 spec**: Phase 0 (`docs/superpowers/specs/2026-05-03-korea-investment-go-migration-design.md`)
- **후속 plan**: 5개 implementation plan (Phase 1.1 ~ Phase 1.5)

---

## §1. 목적과 결정 요약

> **Amendment (2026-05-03, Phase 1.2 brainstorming)**: 본 spec 의 "Python parity" framing 은 메서드 *범위* 기준 (Python 의 28 fetch + 9 IPO helpers 와 같은 도메인 커버리지) 으로 한정. 메서드 *시그니처* 와 *동작* 은 한투 API 문서를 source of truth 로 함 — Python wrapper 의 자체 편의 동작 (country_code 통합 진입점, ETF tr_id 분기, 다국가 fallback 루프 등) 은 미구현. 메서드 명명 스타일은 한투 endpoint path 1:1 (Style A) — §2.2 표 참조.

### Phase 1 목적

`korea-investment-stock` Go 라이브러리를 외부 사용자가 import 해서 쓸 수 있는 **첫 사용 가능 release (`v1.0.0`)** 으로 만든다. Python `v0.18.x` 와 동일한 한투 API 커버리지를 한투 docs 충실 디자인으로 재구현.

### 핵심 결정

| 항목 | 결정 |
|------|------|
| Phase 1 의 framing | 외부 공개 한투 API Go 클라이언트의 첫 사용 가능 release. batch / moneyflow 사용처는 Phase 1 결정에 영향 없음 |
| 메서드 범위 | Python `v0.18.x` 의 28개 fetch 메서드 + IPO helpers 9개 |
| 인프라 범위 | rate limiter (token bucket), token manager (auto refresh, file/redis storage), HTTP client (resty + 재시도), hashkey, master file cache (KOSPI/KOSDAQ ZIP 등) |
| **제외** | Memory cache (사용자가 외부 라이브러리 활용), 선물옵션, 장내채권, 실시간 WebSocket, 주식 주문/잔고 |
| Sub-plan 분해 | 5개 묶음 (각 PR 단위) |
| 호출 스타일 | `client.Domestic.InquirePrice(ctx, "005930")` 형태 (한투 endpoint path 의 마지막 segment 를 PascalCase). **메서드는 sub-package 에 정의** (`domestic.Client.InquirePrice`) — Phase 0 spec §3 일부 refine |
| Config 진입점 | 3개: `NewClient(...)`, `NewClientFromEnv(...)`, `NewClientFromYAML(path, ...)` |
| Master cache | default `os.UserCacheDir()/kis/`, `WithMasterCacheDir(path)` 로 override |
| 응답 typed struct | 한투 API 약어를 PascalCase 로 변환해 그대로 사용 (`StckPrpr`, `PrdyVrss`). 인라인 한국어 코멘트 |
| 테스트 | TDD + `httpmock` + integration build tag + `examples/` + coverage ≥ 80% |
| Release 패턴 | 각 sub-plan PR merge 후 `v0.x.0` tag. 5개 PR 모두 끝난 시점에 `v1.0.0` |

---

## §2. Sub-plan 분해

각 sub-plan 은 한 PR 단위. PR merge 후 release tag.

### Phase 1.1 — 인프라 + Config (`v0.1.0`)

**메서드 없음, 인프라만**. 이 PR 후엔 라이브러리 동작 가능 토대 마련됨.

| Deliverable | 위치 |
|-------------|------|
| Rate limiter (token bucket, default 15 req/sec, thread-safe, 통계) | `internal/ratelimit/` |
| Token manager (OAuth 발급/검증/자동 갱신, 만료 5분 전 선제 발급) | `internal/token/` |
| Token storage (FileStorage, RedisStorage 인터페이스 + 구현) | `internal/token/storage.go` |
| HTTP client (resty wrap, 토큰 자동 주입, 429/5xx 재시도, 토큰 만료 자동 감지/재시도) | `internal/httpclient/` |
| Hashkey 발급 helper | `internal/httpclient/hashkey.go` |
| Master file cache (KOSPI/KOSDAQ ZIP 다운로드 + 디스크 캐시, default `os.UserCacheDir()/kis/`) | `internal/mastercache/` |
| `auth.go` (root) — 토큰 발급 외부 노출 함수 | `auth.go` |
| Functional options (`WithBaseURL`, `WithRetries`, `WithRateLimit`, ...) | `options.go` |
| `errors.go` — `*APIError`, sentinel 에러 (`ErrTokenExpired`, `ErrRateLimited`, `ErrNotFound`) | `errors.go` |
| `NewClient`, `NewClientFromEnv`, `NewClientFromYAML` 3개 진입점 | `client.go` |
| Config 파일 (YAML) 파서 | `config.go` |

**Functional options (예상)**:
- `WithBaseURL(url)` (default `RealEnv`), `WithPaperEnv()`
- `WithRetries(n)` (default 3)
- `WithRateLimit(rps)` (default 15)
- `WithHTTPClient(*http.Client)`
- `WithTokenStorage(storage)`
- `WithMasterCacheDir(path)`
- `WithLogger(*slog.Logger)`
- `WithTimeout(duration)`
- `WithUserAgent(string)`

### Phase 1.2 — 국내 시세 + 심볼 + 차트 (`v0.2.0`)

> **Amendment (2026-05-03, brainstorming 결과)**: 본 sub-plan 은 한투 API 문서를 source of truth 로 재정렬. Python parity wrapper 인 root `FetchPrice (KR/US 통합)` 제거, 메서드명 스타일은 한투 endpoint path 의 마지막 segment 를 PascalCase 로 1:1 매핑 (Style A). KRX 공개 마스터 파일은 한투 API 가 아니므로 `Fetch` prefix 로 구분.

| 메서드 | 위치 | 한투 path | TR_ID |
|-------|------|----------|-------|
| `Domestic.InquirePrice` | `domestic/price.go` | `inquire-price` | FHKST01010100 |
| `Domestic.SearchInfo` | `domestic/info.go` | `search-info` | CTPF1604R |
| `Domestic.SearchStockInfo` | `domestic/info.go` | `search-stock-info` | CTPF1002R |
| `Domestic.InquireDailyItemChartPrice` | `domestic/chart.go` | `inquire-daily-itemchartprice` | FHKST03010100 |
| `Domestic.InquireTimeItemChartPrice` | `domestic/chart.go` | `inquire-time-itemchartprice` | FHKST03010200 |
| `Domestic.FetchKospiSymbols` | `domestic/symbols.go` | KRX 마스터 (한투 API 아님) | — |
| `Domestic.FetchKosdaqSymbols` | `domestic/symbols.go` | KRX 마스터 (한투 API 아님) | — |

총 **7 메서드** (root `FetchPrice` 통합 진입점 제거 — 한투 API 에 없는 wrapper)

**응답 typed struct 명**: `Price`, `ProductInfo` (search-info), `StockInfo` (search-stock-info, 국내 디테일), `DailyChart`, `MinuteChart`. Params struct 는 `<MethodName>Params` (예 `InquireDailyItemChartPriceParams`).

**한투 spec 충실 원칙**: Python `fetch_price` 의 ETF tr_id 분기, `fetch_stock_info` 의 country_code fallback 루프, `fetch_search_stock_info` 의 "KR" 검사 등 Python 자체 wrapper 동작은 모두 미구현. 한투 spec 의 query param 그대로 노출.

### Phase 1.3 — 국내 순위 + 재무 (`v1.1.0`)

> **Amendment (2026-05-04, Phase 1.3 brainstorming)**: 메서드명을 Phase 1.2 와 동일한 Style A (한투 endpoint path 의 마지막 segment 를 PascalCase 로 1:1 매핑) 로 갱신. release tag 는 v1.0.0 publish 이후 Python 시대 태그와 namespace 분리됨에 따라 `v0.3.0` → `v1.1.0`. Python parity wrapper 인 `FetchVolumeRanking`/`FetchChangeRateRanking` 등은 한투 API 의 path 직접 노출로 변경.

| 메서드 | 위치 | 한투 path | TR_ID |
|-------|------|----------|-------|
| `Domestic.InquireVolumeRank` | `domestic/ranking.go` | `quotations/volume-rank` | FHPST01710000 |
| `Domestic.InquireFluctuation` | `domestic/ranking.go` | `ranking/fluctuation` | FHPST01700000 |
| `Domestic.InquireMarketCap` | `domestic/ranking.go` | `ranking/market-cap` | FHPST01740000 |
| `Domestic.InquireDividendRate` | `domestic/ranking.go` | `ranking/dividend-rate` | HHKDB13470100 |
| `Domestic.InquireFinancialRatio` | `domestic/financial.go` | `finance/financial-ratio` | FHKST66430300 |
| `Domestic.InquireIncomeStatement` | `domestic/financial.go` | `finance/income-statement` | FHKST66430200 |
| `Domestic.InquireBalanceSheet` | `domestic/financial.go` | `finance/balance-sheet` | FHKST66430100 |
| `Domestic.InquireProfitRatio` | `domestic/financial.go` | `finance/profit-ratio` | FHKST66430400 |
| `Domestic.InquireGrowthRatio` | `domestic/financial.go` | `finance/growth-ratio` | FHKST66430800 |

총 **9 메서드**

**응답 typed struct 명**: `VolumeRank`, `Fluctuation`, `MarketCap`, `DividendRate`, `FinancialRatio`, `IncomeStatement`, `BalanceSheet`, `ProfitRatio`, `GrowthRatio`. Params struct 는 `<MethodName>Params` (예 `InquireVolumeRankParams`).

**한투 spec 충실 원칙** (Phase 1.2 와 동일): query param/응답 필드를 한투 docs (`docs/api/국내주식/<API>.md`) 그대로 노출. Python wrapper convenience 미반영.

### Phase 1.4 — 국내 투자자 + 업종 + IPO (`v1.2.0`)

> **Amendment (2026-05-05, Phase 1.4 brainstorming)**: 메서드명을 Phase 1.2/1.3 와 동일한 Style A (path 의 last segment PascalCase) 로 갱신. 시장별 투자자매매동향 (시세) 1개 추가 (총 6 메서드). IPO helpers 9개 제거 — Phase 1.2 amendment 의 "Python wrapper convenience 미반영" 정책과 일관 (helpers 는 client-side data 가공이라 caller 가 직접 처리). release tag `v0.4.0` → `v1.2.0`.

| 메서드 | 위치 | 한투 path | TR_ID |
|-------|------|----------|-------|
| `Domestic.InquireInvestorTradeByStockDaily` | `domestic/investor.go` | `quotations/investor-trade-by-stock-daily` | FHPTJ04160001 |
| `Domestic.InquireInvestorDailyByMarket` | `domestic/investor.go` | `quotations/inquire-investor-daily-by-market` | FHPTJ04040000 |
| `Domestic.InquireInvestorTimeByMarket` | `domestic/investor.go` | `quotations/inquire-investor-time-by-market` | FHPTJ04030000 |
| `Domestic.InquireIndexPrice` | `domestic/industry.go` | `quotations/inquire-index-price` | FHPUP02100000 |
| `Domestic.InquireIndexCategoryPrice` | `domestic/industry.go` | `quotations/inquire-index-category-price` | FHPUP02140000 |
| `Domestic.InquirePubOffer` | `domestic/ipo.go` | `ksdinfo/pub-offer` | HHKDB669108C0 |

총 **6 메서드** (helpers 제외)

**응답 typed struct 명**: `InvestorTradeByStockDaily`, `InvestorDailyByMarket`, `InvestorTimeByMarket`, `IndexPrice`, `IndexCategoryPrice`, `PubOffer`. Params struct 는 `<MethodName>Params`.

**IPO 메서드명 결정**: `InquirePubOffer` 는 Style A 룰 (path last segment `pub-offer` → `PubOffer`) 정확히 따름. godoc 으로 "공모주청약일정 (IPO Schedule) 조회" 명시.

**한투 spec 충실 원칙** (Phase 1.2/1.3 와 동일): query param/응답 필드를 한투 docs (`docs/api/국내주식/<API>.md`) 그대로 노출. Python wrapper convenience 미반영.

### Phase 1.5 — 해외 전체 (`v1.3.0`, Python parity 완성)

> **Amendment (2026-05-05, Phase 1.5 brainstorming)**: 메서드명을 Phase 1.2~1.4 와 동일한 Style A (한투 endpoint path 의 마지막 segment 를 PascalCase 로 1:1 매핑) 로 갱신. NASDAQ/NYSE/AMEX 별 메서드는 `FetchOverseasSymbols(market)` 로 통합 (Python wrapper convenience 미반영 정책 일관). 차트는 `dailyprice` (단일 종목) + `inquire-daily-chartprice` (종목/지수/환율 일/주/월/년) 두 endpoint 모두 별도 메서드로 노출. `해외주식_상품기본정보` 추가 (Python `fetch_search_stock_info` 파리티). 총 6 메서드.

| 메서드 | 위치 | 한투 path | TR_ID |
|-------|------|----------|-------|
| `Overseas.InquirePriceDetail` | `overseas/price.go` | `overseas-price/v1/quotations/price-detail` | HHDFS76200200 |
| `Overseas.SearchInfo` | `overseas/search.go` | `overseas-price/v1/quotations/search-info` | CTPF1702R |
| `Overseas.InquireDailyPrice` | `overseas/chart.go` | `overseas-price/v1/quotations/dailyprice` | HHDFS76240000 |
| `Overseas.InquireDailyChartPrice` | `overseas/chart.go` | `overseas-price/v1/quotations/inquire-daily-chartprice` | FHKST03030100 |
| `Overseas.InquireUpdownRate` | `overseas/ranking.go` | `overseas-stock/v1/ranking/updown-rate` | HHDFS76290000 |
| `Overseas.FetchOverseasSymbols(market)` | `overseas/symbols.go` | (외부 다운로드 — 11 거래소) | — |

총 **6 메서드**

**응답 typed struct 명**: `PriceDetail`, `OverseasSearchInfo` (또는 `OverseasProductInfo` — domestic 의 `ProductInfo` 와 구분), `DailyPrice`, `DailyChartPrice`, `UpdownRate`, `OverseasSymbol[]`. Params struct 는 `<MethodName>Params`.

**해외 마스터 파일**: KRX 와 동일 패턴 — 외부 다운로드 (`https://new.real.download.dws.co.kr/common/master/<market>mst.cod.zip`). 11 거래소 (`nas`/`nys`/`ams`/`shs`/`shi`/`szs`/`szi`/`tse`/`hks`/`hnx`/`hsx`). 새 internal package: `internal/overseasmaster/` (KRX 와 분리, cp949 + fwf 형식 가능 — 실제 형식 확인 후 codec 재사용 여부 결정).

**`overseas.Client` 시그니처 확장**: `overseas.New(http)` → `overseas.New(http, master)` (Phase 1.2 의 domestic 패턴). root `client.go` 의 `wireInfra` 가 `c.masterC` 주입.

**한투 spec 충실 원칙** (Phase 1.2~1.4 와 동일): query param/응답 필드를 한투 docs (`docs/api/해외주식/<API>.md`) 그대로 노출. Python wrapper 의 `country_code` fallback 루프, NASDAQ/NYSE/AMEX 별 편의 메서드 등은 미구현.

### 합계

- **인프라**: 1 PR
- **메서드**: 7 (1.2) + 9 (1.3) + 6 (1.4) + 6 (1.5) = **28 메서드** (Phase 1.5 amendment: NASDAQ/NYSE/AMEX 통합으로 7→6, 해외 SearchInfo 추가)
- **Release tags**: ~~v0.1.0/v0.2.0~~ (Python era namespace 충돌로 삭제) → **`v1.0.0`** (Phase 1.1+1.2 통합), `v1.1.0` (Phase 1.3), `v1.2.0` (Phase 1.4), `v1.3.0` (Phase 1.5 — 해외 전체)

---

## §3. 인프라 디자인 (Phase 1.1)

### Rate limiter (`internal/ratelimit/`)

```go
package ratelimit

type Limiter struct {
    callsPerSec    float64
    minInterval    time.Duration
    lastCall       time.Time
    mu             sync.Mutex
    totalCalls     int64
    throttledCalls int64
    totalWait      time.Duration
}

func New(callsPerSec float64) *Limiter
func (l *Limiter) Wait(ctx context.Context) error  // ctx.Done() 시 sleep 인터럽트
func (l *Limiter) Stats() Stats
```

- 토큰 버킷 (Python 동일), thread-safe
- `Wait` 가 `context.Context` 받아 취소/타임아웃 지원 (Python 에는 없음 — 개선)
- Default 15 req/sec

### Token manager (`internal/token/`)

```go
package token

type Storage interface {
    Save(ctx context.Context, token *AccessToken) error
    Load(ctx context.Context) (*AccessToken, error)
    Clear(ctx context.Context) error
}

type AccessToken struct {
    Value     string
    TokenType string  // "Bearer"
    ExpiresAt time.Time
}

type Manager struct {
    storage  Storage
    apiKey   string
    secret   string
    baseURL  string
    httpDo   func(*http.Request) (*http.Response, error)
    mu       sync.Mutex
    cache    *AccessToken
}

func NewManager(...) *Manager
func (m *Manager) Get(ctx context.Context) (string, error)  // 만료 5분전 자동 갱신
func (m *Manager) Refresh(ctx context.Context) (string, error)
```

- `FileStorage` (`~/.cache/kis/token.json`), `RedisStorage` 두 구현
- 만료 5분 전 선제 발급 (Python 동일)
- 동시 호출 시 한 번만 발급 (singleflight 패턴)

### HTTP client (`internal/httpclient/`)

```go
package httpclient

type Client struct {
    resty       *resty.Client
    rateLimiter *ratelimit.Limiter
    tokenMgr    *token.Manager
}

func New(opts Config) *Client
func (c *Client) Do(ctx context.Context, req *Request) (*Response, error)
```

동작 순서:
1. `rateLimiter.Wait(ctx)` — 호출 빈도 제어
2. `tokenMgr.Get(ctx)` — 유효 토큰 확보
3. resty 로 호출
4. 응답 검사:
   - `rt_cd != "0"` → `*kis.APIError` 리턴
   - 토큰 만료 메시지 ("기간이 만료된 token") → `tokenMgr.Refresh()` 후 1회 재시도
   - 5xx → `WithRetries(n)` 만큼 exponential backoff
   - 429 → 대기 후 재시도
5. JSON unmarshal → typed struct

### Hashkey (`internal/httpclient/hashkey.go`)

```go
func (c *Client) Hashkey(ctx context.Context, body any) (string, error)
```

- 주문 등 일부 한투 API 가 요구하는 hashkey 생성
- Phase 1 에서 주문 API 는 안 만들지만 hashkey 함수는 미리 만들어둠

### Master file cache (`internal/mastercache/`)

```go
package mastercache

type Cache struct {
    dir string  // default: os.UserCacheDir()/kis/
    ttl time.Duration  // default: 168h (1주)
    mu  sync.Mutex
}

func New(dir string, ttl time.Duration) *Cache
func (c *Cache) Get(ctx context.Context, name string, fetch func() (io.Reader, error)) ([]byte, error)
```

- KOSPI ZIP, KOSDAQ ZIP 다운로드 + 디스크 캐시
- TTL 만료 시 자동 재다운로드
- `WithMasterCacheDir(path)` 로 사용자 override 가능

### Errors (`errors.go` root)

```go
package kis

type APIError struct {
    RtCode  string  // rt_cd
    MsgCode string  // msg_cd
    Message string  // msg1
    TrID    string  // 디버깅용
}

func (e *APIError) Error() string

var (
    ErrTokenExpired = errors.New("kis: token expired")
    ErrRateLimited  = errors.New("kis: rate limited")
    ErrNotFound     = errors.New("kis: resource not found")
    ErrUnauthorized = errors.New("kis: unauthorized")
)
```

`errors.As(err, &apiErr)` 로 한투 응답 에러 분기, `errors.Is(err, kis.ErrTokenExpired)` 로 sentinel 확인.

---

## §4. Config 시스템 디자인 (Phase 1.1)

### 3 가지 진입점

```go
// 1. 직접 전달
client, err := kis.NewClient(apiKey, apiSecret, accountNo, opts...)

// 2. env vars 자동 감지 (KOREA_INVESTMENT_*)
client, err := kis.NewClientFromEnv(opts...)

// 3. YAML 파일
client, err := kis.NewClientFromYAML("~/.config/kis/config.yaml", opts...)
```

세 함수 모두 내부적으로 `NewClient` 의 thin wrapper. `opts` 는 functional options (모든 진입점 공통).

### 환경변수 (Python 동일)

| 변수 | 의미 | 필수 |
|------|------|------|
| `KOREA_INVESTMENT_API_KEY` | API key | ✓ |
| `KOREA_INVESTMENT_API_SECRET` | API secret | ✓ |
| `KOREA_INVESTMENT_ACCOUNT_NO` | 계좌번호 (8-2 형식) | ✓ |
| `KOREA_INVESTMENT_TOKEN_STORAGE` | `file` 또는 `redis` (default `file`) | |
| `KOREA_INVESTMENT_TOKEN_FILE` | 토큰 파일 경로 (default `~/.cache/kis/token.json`) | |
| `KOREA_INVESTMENT_REDIS_URL` | Redis URL | (storage=redis 시) |
| `KOREA_INVESTMENT_REDIS_PASSWORD` | Redis password | |
| `KOREA_INVESTMENT_BASE_URL` | API base URL (default RealEnv) | |

### YAML 형식 (Python 동일)

```yaml
api_key: your-api-key
api_secret: your-api-secret
acc_no: "12345678-01"

# Optional
base_url: https://openapi.koreainvestment.com:9443
token_storage_type: file
token_file: ~/.cache/kis/token.json
master_cache_dir: ~/.cache/kis/
rate_limit: 15
retries: 3

# Redis (token_storage_type: redis 시)
redis_url: redis://localhost:6379/0
redis_password: ""
```

### 우선순위 (한 진입점 안에서)

`NewClientFromEnv` 또는 `NewClientFromYAML` 은 base config 를 로드한 후, **functional options 가 마지막에 적용** (= override 가능).

```go
client, _ := kis.NewClientFromEnv(
    kis.WithRateLimit(20),  // env 의 rate_limit 무시하고 20 적용
)
```

### Config 구조체 (외부 노출)

```go
package kis

type Config struct {
    APIKey         string
    APISecret      string
    AccountNo      string
    BaseURL        string
    TokenStorage   string  // "file" | "redis"
    TokenFile      string
    RedisURL       string
    RedisPassword  string
    MasterCacheDir string
    RateLimit      float64
    Retries        int
}

func LoadConfigFromEnv() (*Config, error)
func LoadConfigFromYAML(path string) (*Config, error)
```

`Config` 자체를 외부 노출하면 사용자가 직접 만들어 수정도 가능.

---

## §5. 메서드 디자인

> **Amendment (2026-05-03, Phase 1.2 brainstorming)**: 본 §5 의 메서드명 (`FetchPrice` 등) 과 디렉터리 구조 (`info.go` 등) 코드 예시는 Phase 0 시점의 일반 패턴 시각화이며, **실제 sub-plan 의 메서드명/응답 typed struct 명은 각 sub-plan 표 (§2.2~§2.5) 의 한투 endpoint path 1:1 매핑 (Style A) 을 source of truth 로 함**. 본 §5 의 코드 예시는 godoc/typed struct/에러처리/JSON tag 등 *구조적 패턴* 만 참조하고, 식별자는 sub-plan 표 기준으로 적용.


### 디렉터리 구조 (Phase 0 spec §3 refine — 메서드 sub-package 이동)

```
korea-investment-stock/
├── client.go              # Client struct + sub-client wiring (Domestic *domestic.Client, Overseas *overseas.Client)
├── auth.go                # 토큰 발급 외부 노출 (root)
├── errors.go              # APIError + sentinel
├── options.go             # functional options
├── config.go              # Config struct + LoadConfigFromEnv / LoadConfigFromYAML
├── domestic/
│   ├── client.go          # domestic.Client struct + 생성자 (internal use)
│   ├── price.go           # FetchDomesticPrice + Price 타입
│   ├── info.go            # FetchStockInfo, FetchSearchStockInfo + StockInfo 타입
│   ├── chart.go           # FetchDomesticChart, FetchDomesticMinuteChart + Chart 타입들
│   ├── ranking.go         # 4개 ranking 메서드 + Ranking 타입들
│   ├── financial.go       # 5개 재무 메서드 + Financial 타입들
│   ├── investor.go        # 2개 투자자 메서드 + 타입
│   ├── industry.go        # 2개 업종 메서드 + 타입
│   ├── ipo.go             # FetchIPOSchedule + 9 helpers + IPO 타입들
│   └── symbols.go         # FetchKospiSymbols, FetchKosdaqSymbols + Symbol 타입
├── overseas/
│   ├── client.go          # overseas.Client struct + 생성자
│   ├── price.go           # FetchPriceDetailOverseas + 타입
│   ├── chart.go           # FetchOverseasChart + 타입
│   ├── ranking.go         # FetchOverseasChangeRateRanking + 타입
│   └── symbols.go         # NASDAQ/NYSE/AMEX/Overseas + 타입
└── internal/
    ├── httpclient/
    ├── ratelimit/
    ├── token/
    └── mastercache/
```

### 호출 흐름

```go
// 1. root 의 NewClient 가 sub-client 를 wiring
type Client struct {
    apiKey, apiSecret, accountNo string
    opts                         clientOptions
    httpClient                   *httpclient.Client  // shared
    Domestic                     *domestic.Client
    Overseas                     *overseas.Client
}

// 2. sub-package 의 Client 가 shared httpclient 를 받아 메서드 구현
package domestic

type Client struct {
    http *httpclient.Client  // root 가 주입
}

func New(http *httpclient.Client) *Client { return &Client{http: http} }

func (c *Client) FetchPrice(ctx context.Context, symbol string) (*Price, error) {
    // ... build request, call c.http.Do, unmarshal ...
}
```

→ root 의 `kis.Client` 가 `httpclient` 를 만들어 sub-package 에 주입. sub-package 끼리 공유.

### Typed struct 컨벤션 — 한투 API 약어 + 인라인 코멘트

```go
package domestic

// Price 는 주식현재가_시세 응답.
// API 문서: docs/api/국내주식/주식현재가_시세.md
type Price struct {
    StckPrpr     decimal.Decimal `json:"stck_prpr"`           // 주식 현재가
    PrdyVrss     decimal.Decimal `json:"prdy_vrss"`           // 전일 대비
    PrdyVrssSign string          `json:"prdy_vrss_sign"`      // 전일 대비 부호 (1: 상한, 2: 상승, 3: 보합, 4: 하한, 5: 하락)
    PrdyCtrt     float64         `json:"prdy_ctrt,string"`    // 전일 대비율 (등락률 %)
    StckOprc     decimal.Decimal `json:"stck_oprc"`           // 시가
    StckHgpr     decimal.Decimal `json:"stck_hgpr"`           // 고가
    StckLwpr     decimal.Decimal `json:"stck_lwpr"`           // 저가
    AcmlVol      int64           `json:"acml_vol,string"`     // 누적 거래량
    AcmlTrPbmn   int64           `json:"acml_tr_pbmn,string"` // 누적 거래대금
    HtsAvls      decimal.Decimal `json:"hts_avls"`            // HTS 시가총액 (억원)
    TrhtYn       string          `json:"trht_yn"`             // 거래 정지 여부 (Y/N)
    AskPrice1    decimal.Decimal `json:"askp1"`               // 매도호가1
    BidPrice1    decimal.Decimal `json:"bidp1"`               // 매수호가1
}
```

### 규칙

| 항목 | 규칙 |
|------|------|
| **필드명** | 한투 API 문서의 약어를 PascalCase 로 변환 (`stck_prpr` → `StckPrpr`) |
| **JSON 태그** | 한투 원본 키 그대로. 숫자가 문자열로 오는 경우 `,string` 옵션 |
| **타입** | 가격 = `decimal.Decimal` / 수량·거래량 = `int64` / 백분율 = `float64` / 시각·날짜 = `time.Time` 또는 `civil.Date` / boolean 플래그 = `string` (Y/N) 또는 `bool` (변환 helper 제공) |
| **타입 doc 코멘트** | 타입 위 (godoc 표준). 어느 API 응답인지 + docs/api 경로 명시 |
| **필드 코멘트** | **오른쪽 인라인** (`// 한국어 설명`). enum 값이 있으면 가능한 값들 명시 |
| **메서드명** | 사용자가 발견하기 쉽게 의미 있는 영문 (`FetchPrice`, `FetchKospiSymbols`) |

### 긴 enum 의 경우 const 그룹

```go
type Price struct {
    PrdyVrssSign string `json:"prdy_vrss_sign"` // 전일 대비 부호 (PrdyVrssSign* 상수 참조)
}

const (
    PrdyVrssSignUpperLimit = "1"  // 상한
    PrdyVrssSignUp         = "2"  // 상승
    PrdyVrssSignFlat       = "3"  // 보합
    PrdyVrssSignLowerLimit = "4"  // 하한
    PrdyVrssSignDown       = "5"  // 하락
)
```

(enum 값이 정말 많거나 사용자가 자주 비교할 때만 추가)

### IPO helpers — 9개 (Python parity)

```go
package domestic

type IPOSchedule struct {
    Output1 []IPOItem `json:"output1"`  // 데이터 항목 리스트
}

type IPOItem struct {
    ShCode   string `json:"sh_code"`     // 종목 코드
    IsinName string `json:"isin_name"`   // 종목명
    SubsDt   string `json:"subs_dt"`     // 청약 시작일 (YYYYMMDD)
    // ... (한투 약어 유지, 인라인 코멘트)
}

// helpers — 의미 있는 영문 메서드명, 응답 필드는 한투 약어
func (s *IPOSchedule) FilterByPeriod(from, to time.Time) []IPOItem
func (s *IPOSchedule) UpcomingSubscriptions(within time.Duration) []IPOItem
// ... 7 more
```

### 에러 처리 흐름 (모든 메서드 공통)

```go
func (c *Client) FetchPrice(ctx context.Context, symbol string) (*Price, error) {
    resp, err := c.http.Do(ctx, &Request{
        Method: http.MethodGet,
        Path:   "/uapi/domestic-stock/v1/quotations/inquire-price",
        TrID:   "FHKST01010100",
        Query:  map[string]string{...},
    })
    if err != nil {
        return nil, err  // *APIError, ErrTokenExpired, ErrRateLimited 등
    }
    var price Price
    if err := json.Unmarshal(resp.Output1, &price); err != nil {
        return nil, fmt.Errorf("kis: parse price: %w", err)
    }
    return &price, nil
}
```

`httpclient.Do` 가 이미 rate limit / token / 재시도 처리. 메서드 코드는 단순.

---

## §6. 테스트 전략

### TDD 흐름 (각 메서드)

각 fetch 메서드 구현은 다음 4단계:

1. **테스트 작성** (`*_test.go`)
   - `docs/api/<카테고리>/<API>.md` 의 응답 샘플을 읽어와서 mock JSON 으로 사용
   - 예상 typed struct 비교
2. **실행 → 실패** (메서드 미구현)
3. **메서드 구현**
4. **실행 → 성공**

### httpmock 패턴 — `domestic/price_test.go` 예시

```go
package domestic

import (
    "context"
    "net/http"
    "testing"

    "github.com/jarcoal/httpmock"
    "github.com/shopspring/decimal"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"

    "github.com/kenshin579/korea-investment-stock/internal/httpclient"
)

func TestClient_FetchPrice(t *testing.T) {
    httpmock.Activate()
    defer httpmock.DeactivateAndReset()

    httpmock.RegisterResponder(
        http.MethodGet,
        "=~^https://openapi\\.koreainvestment\\.com.*/inquire-price",
        httpmock.NewStringResponder(200, `{
            "rt_cd": "0",
            "msg_cd": "MCA00000",
            "msg1": "정상처리 되었습니다.",
            "output": {
                "stck_prpr": "75800",
                "prdy_vrss": "-200",
                "prdy_vrss_sign": "5",
                "prdy_ctrt": "-0.26",
                "stck_oprc": "76000",
                "stck_hgpr": "76200",
                "stck_lwpr": "75500",
                "acml_vol": "12345678",
                "hts_avls": "452312"
            }
        }`),
    )

    http := httpclient.NewForTest(...)  // resty + mocked transport
    c := New(http)

    price, err := c.FetchPrice(context.Background(), "005930")
    require.NoError(t, err)
    assert.Equal(t, decimal.NewFromInt(75800), price.StckPrpr)
    assert.Equal(t, decimal.NewFromInt(-200), price.PrdyVrss)
    assert.Equal(t, "5", price.PrdyVrssSign)
    assert.Equal(t, -0.26, price.PrdyCtrt)
    assert.Equal(t, int64(12345678), price.AcmlVol)
}
```

mock JSON 은 가능하면 `docs/api/국내주식/주식현재가_시세.md` 의 example 을 그대로 복사. AI 가 docs 와 mock 을 동시에 참조하면 자동 생성 친화적.

### Mock data 위치

각 sub-package 에 `testdata/` 디렉터리:

```
domestic/
├── price.go
├── price_test.go
└── testdata/
    ├── price_success.json
    ├── price_error_token_expired.json
    └── ...
```

긴 mock JSON 은 `testdata/*.json` 로 분리, `os.ReadFile` 로 로드.

### Integration test (build tag)

```go
//go:build integration
// +build integration

package domestic

func TestClient_FetchPrice_Integration(t *testing.T) {
    if os.Getenv("KOREA_INVESTMENT_API_KEY") == "" {
        t.Skip("KOREA_INVESTMENT_API_KEY not set")
    }
    // 실제 API 호출
}
```

```bash
# Unit only (default, CI)
go test ./...

# Integration 포함
go test -tags=integration ./...
```

### Examples 디렉터리

```
examples/
├── basic/
│   └── main.go              # NewClient + 단일 호출
├── env_config/
│   └── main.go              # NewClientFromEnv
├── yaml_config/
│   ├── main.go
│   └── config.yaml
├── domestic_price/
│   └── main.go              # FetchPrice + 응답 출력
├── domestic_chart/
│   └── main.go
├── overseas_price/
│   └── main.go
└── batch_with_rate_limit/
    └── main.go              # 250 종목 fetch + rate limit 동작 확인
```

각 example 은 `go run examples/<dir>/main.go` 로 실행. 실제 KIS credentials 필요.

### Coverage

| 영역 | 목표 |
|------|------|
| `internal/ratelimit` | ≥ 90% (단순 logic, 외부 의존 없음) |
| `internal/token` | ≥ 85% (storage 포함) |
| `internal/httpclient` | ≥ 80% (httpmock 으로 transport 단위) |
| `internal/mastercache` | ≥ 80% |
| `domestic/`, `overseas/` 메서드 | ≥ 80% (mock 단위 테스트) |
| 전체 | ≥ 80% |

`go test -coverprofile=coverage.out ./...` 로 측정. CI 에 coverage threshold gate 추가 (Phase 1.1).

### Test helpers

`internal/testutil/` 에 공통 helper:

- `LoadFixture(t, "price_success.json")` — testdata 로드
- `NewMockClient(t)` — http client 모킹된 인스턴스 생성
- `AssertAPIError(t, err, "MCA00001")` — APIError 검증 helper

---

## §7. 성공 기준 / 진입/종료 조건

### Phase 1 의 Deliverable

5 sub-plan 모두 완료 시 다음 상태:

1. **본 설계 문서** — `docs/superpowers/specs/2026-05-03-phase1-api-coverage-design.md` git commit
2. **5개 implementation plan** — `docs/superpowers/specs/2026-05-03-phase1-{1..5}-*.md`
3. **5개 PR merge** — 각각 release tag 박힘
   - Phase 1.1 → `v0.1.0` (인프라+Config)
   - Phase 1.2 → `v0.2.0` (국내 시세+심볼+차트)
   - Phase 1.3 → `v0.3.0` (국내 순위+재무)
   - Phase 1.4 → `v0.4.0` (국내 투자자+업종+IPO)
   - Phase 1.5 → **`v1.0.0`** (해외 전체 = Python parity 완성)
4. **`v1.0.0` Go module proxy 등록** — `go get github.com/kenshin579/korea-investment-stock@v1.0.0` 가능
5. **README 업데이트** — Quick Start 의 코드가 실제로 동작 (Phase 0 의 forward reference 가 사실이 됨)

### 각 sub-plan 의 success criteria (공통)

| 항목 | 기준 |
|------|------|
| Build | `go build ./...` 통과 |
| Vet | `go vet ./...` 출력 없음 |
| Format | `gofmt -l .` 출력 없음 |
| Unit test | 새로 추가된 메서드들 모두 mock 기반 테스트 통과 |
| Coverage | sub-plan 의 추가 코드 영역 ≥ 80% |
| Integration test (선택) | `go test -tags=integration` 가 KIS credentials 있을 때 통과 (PR merge 의 strict 조건 아님 — 사용자가 manual 검증) |
| Examples | 해당 카테고리 example 추가 + 사용자가 실행해서 동작 확인 |
| Doc | 새 메서드의 godoc + README 의 Quick Start 갱신 |
| Commit/PR 정책 | Phase 0 와 동일 — feature branch, `[chore]`/`[feat]` prefix, Co-Authored-By, kenshin579 reviewer |

### Phase 1 종료 후 (= Phase 2 진입 조건)

Phase 1 = Python parity. 그 다음:

- **Phase 2** — Python 에 없던 추가 한투 API (시간외, 공매도, 종목조건검색, 관심종목, 휴장일, ETF NAV 추이 등). docs/api 의 미구현 240여 개 중 가치 큰 것 추려냄. 별도 spec
- **Phase 3 (선택)** — 실시간 WebSocket
- **Phase 4 (선택)** — 주식 주문/잔고/예약주문 (트레이딩)

### Non-goals (Phase 1 에서 의도적으로 제외)

- ❌ Memory cache 라이브러리 빌트인 (사용자가 외부 라이브러리 활용)
- ❌ 선물옵션 / 장내채권 (Phase 0 영구 제외)
- ❌ 실시간 WebSocket
- ❌ 주식 주문 / 잔고 / 예약주문 (Phase 4 후보)
- ❌ Python 의 `fetch_etf_domestic_price` — Python 에 있긴 한데 batch 가 안 쓰고 일반 `fetch_domestic_price` 로 ETF 도 처리됨. Go 에선 단일 진입점 통합

### 위험 요소와 대응

| 위험 | 대응 |
|------|------|
| 한투 API 응답 스키마 변경 | docs/api 와 mock 응답 1:1 매칭 → 변경 시 mock 만 업데이트해 빠르게 대응 |
| 토큰 만료/재발급 race condition | singleflight 패턴 (한 토큰 만료 시 한 번만 재발급, 다른 호출은 결과 공유) |
| Rate limit 초과 폭주 시 race | `Wait(ctx)` 가 atomic mutex 사용, ctx.Done() 시 sleep 인터럽트 |
| Master file (KOSPI ZIP) 다운로드 실패 | 캐시된 옛 파일 fallback, 명시적 `force` 플래그로 강제 재시도 |
| Decimal 변환 오류 (한투가 빈 문자열 보낼 때) | helper 로 통합, `""` → `decimal.Zero` 명시 |
| 5 sub-plan 중간에 의존성 충돌 (Phase 1.2 가 Phase 1.1 의 인프라 변경 요구) | sub-plan 작성 전 Phase 1.1 PR merge 후 시작. 인프라 추가 변경 필요 시 별도 fix PR |

### 본 spec 의 Non-goals (의도적으로 다루지 않은 것)

- 각 메서드의 정확한 query 파라미터/요청 헤더 디테일 → 각 sub-plan implementation plan 에서 다룸
- 모든 응답 필드의 정확한 매핑 표 → docs/api/<API>.md 를 source of truth 로 사용
- moneyflow / stock-data-batch 의 데이터 흐름 → 별도 spec
