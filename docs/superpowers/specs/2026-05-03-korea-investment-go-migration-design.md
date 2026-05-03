# korea-investment-stock Go 마이그레이션 설계 (Phase 0)

- **상태**: Draft
- **작성일**: 2026-05-03
- **작성자**: kenshin579
- **대상 repo**: `korea-investment-stock`, `stock-data-batch`, `moneyflow`
- **연관 후속 spec**: Phase 1 (필요 API 식별/우선순위), `moneyflow 새 개발 spec`

---

## §1. 목적과 결정 요약

### Phase 0 목적

`korea-investment-stock` 과 `stock-data-batch` 의 언어/플랫폼 전환 결정 및 마이그레이션 절차 확정. 구체적인 API 추가 계획은 Phase 1 에서 다룬다.

### 핵심 결정

| 항목 | 결정 |
|------|------|
| 언어/플랫폼 | Python → **Go** 로 전환 |
| 동기 | ① moneyflow 백엔드(Go) 와 스택 통합 ② 동시성/성능 ③ 단일 바이너리 배포 ④ 타입 안정성 |
| `korea-investment-stock` 위치 | **현 repo 그대로 사용**, Go module 로 공개 (외부 사용자도 import 가능) |
| `korea-investment-stock` Python 코드 | main 에서 즉시 제거. **v1.x 태그·옛 commit 보존**. PyPI v1.x 는 그대로 두되 **deprecation notice** 부착, 신규 추가 없음 |
| `stock-data-batch` repo | **archive** 처리 (별도 repo 자체 폐기) |
| `stock-data-batch` 기능 | **moneyflow 에 통합**, 별도 CLI 로 실행 가능 |
| `moneyflow` 자체 | 신규 개발 (작업 시작 전 repo 비우기). **본 spec 의 직접 scope 아님** — 별도 spec 에서 다룸 |
| 선물옵션 / 장내채권 | 모두 **scope 영구 제외** |

---

## §2. 전환 절차/타임라인

전환은 외부 사용자에게 충격이 가지 않도록 단계적으로 진행한다.

### Step 1. `korea-investment-stock` 사전 정리

| 작업 | 비고 |
|------|------|
| Python 코드 freeze (이 시점부터 Python 신규 PR 받지 않음) | |
| `v1.x-final` 태그 (또는 정확한 마지막 버전 태그) push | 외부 사용자가 태그로 옛 코드 참조 가능 |
| PyPI 패키지에 deprecation notice 부착 | `pyproject.toml` description, `README.md` 상단, 다음 release note |
| `README.md` 에 마이그레이션 안내 추가 | "Python v1.x 사용자는 태그 참조, 신규 사용자는 Go 사용" |
| main 에서 Python 코드 제거 | git history 에는 그대로 남음 |
| `go.mod` 초기화 + 빈 패키지 스켈레톤 | `client.go`, `domestic/`, `overseas/`, `internal/` 디렉터리만 존재 |

### Step 2. Go 라이브러리 개발 (Phase 1 이후)

| 작업 | 비고 |
|------|------|
| Phase 1 에서 추려낸 필요 API 부터 구현 | Phase 1 spec 의 결과물 |
| 외부 사용자도 import 가능한 형태로 공개 | Go module proxy 자동 등록 |

### Step 3. `moneyflow` 신규 개발 + `stock-data-batch` 흡수

| 작업 | 비고 |
|------|------|
| moneyflow repo 비우고 새로 개발 시작 | **별도 spec — 본 spec scope 외** |
| Go `korea-investment-stock` 을 import 해 batch 기능 통합 | moneyflow 새 개발 spec 영역 |
| 별도 CLI 로 실행 가능하게 설계 | moneyflow 새 개발 spec 영역 |

### Step 4. `stock-data-batch` repo archive

| 작업 | 비고 |
|------|------|
| moneyflow 통합 완료 + 운영 안정화 확인 후 | |
| repo 를 GitHub archive 처리 (read-only) | git history 보존 |
| `README.md` 에 "moneyflow 로 통합되었음" 안내 | |

### 의존성 그래프

```
Step 1 (Python 정리) ──────┐
                          ↓
                    Step 2 (Go 라이브러리 개발) ─────┐
                                                  ↓
                                            Step 3 (moneyflow 신규 개발 + 통합)
                                                  ↓
                                            Step 4 (stock-data-batch archive)
```

Step 1 과 Step 2 는 순차 진행 (Python freeze 후 Go 작업 시작 — 같은 repo 에서 충돌 방지). Step 3 는 별도 spec 에서 다루며, Step 4 는 Step 3 완료 후.

---

## §3. Go 라이브러리 구조와 명명 규칙

### Module path / versioning

| 항목 | 결정 |
|------|------|
| Module path | `github.com/kenshin579/korea-investment-stock` (repo 경로 그대로) |
| 패키지 이름 | `kis` |
| 시작 버전 | `v0.1.0` 부터 → 안정화 시 `v1.0.0`. **Python v1.x 와 별개**. README 상단에 명시 |

### 디렉터리 구조

```
korea-investment-stock/
├── go.mod
├── go.sum
├── README.md (마이그레이션 안내 + Go 사용법)
├── CHANGELOG.md
├── LICENSE
├── client.go              # Client 생성/설정 (functional options)
├── auth.go                # 토큰 발급/폐기
├── errors.go              # KIS 에러 타입
├── domestic/              # 국내주식 API
│   ├── price.go
│   ├── chart.go
│   ├── ranking.go
│   ├── financial.go
│   ├── investor.go
│   ├── ipo.go
│   └── symbols.go
├── overseas/              # 해외주식 API
│   ├── price.go
│   ├── chart.go
│   └── ranking.go
└── internal/              # 외부 노출 X
    ├── httpclient/        # resty 기반, 재시도 + 토큰 자동 갱신
    ├── ratelimit/         # token bucket
    └── token/             # 토큰 저장소 (file/redis)
```

### 호출 스타일

go-github / stripe-go 의 **1단계 서비스 그룹화** 채택.

```go
client := kis.NewClient(appKey, appSecret, accountNo,
    kis.WithBaseURL(kis.RealEnv),
    kis.WithRetries(3),
)

// 국내
price, err := client.Domestic.FetchPrice(ctx, "005930")
chart, err := client.Domestic.FetchDailyChart(ctx, "005930", from, to)
ranking, err := client.Domestic.FetchVolumeRanking(ctx, kis.MarketKospi)

// 해외
usPrice, err := client.Overseas.FetchPrice(ctx, "AAPL", kis.US)
```

도메인당 메서드가 30+ 로 늘어나면 **2단계 그룹화** (`client.Domestic.Price.Fetch(...)`) 로 진화. 결정은 Phase 1 의 API 매핑 시점에.

### Functional options

```go
kis.NewClient(appKey, appSecret, accountNo,
    kis.WithBaseURL(kis.RealEnv),       // or kis.PaperEnv
    kis.WithRetries(3),
    kis.WithRateLimit(15),               // 호출/초
    kis.WithTokenStorage(fileStorage),
    kis.WithHTTPClient(customClient),    // *http.Client (표준)
    kis.WithLogger(slog.Default()),
)
```

옵션 추가가 SemVer breaking 이 아님 → 장기 유지보수 우위.

### Context.Context 필수

모든 API 호출 메서드는 첫 인자로 `context.Context` 받음.

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
price, err := client.Domestic.FetchPrice(ctx, "005930")
```

### 응답 타입

- typed struct 사용
- 한투 API 응답의 한글 약어 필드(`stck_prpr`, `prdy_vrss` 등)는 JSON 태그로 매핑, 사용자에겐 명료한 영문 필드명으로 노출
- 가격은 `shopspring/decimal` 사용 (금융 정밀도 보호)

```go
type DomesticPrice struct {
    Current        decimal.Decimal `json:"stck_prpr"`
    PreviousChange decimal.Decimal `json:"prdy_vrss"`
    ChangeRate     float64         `json:"prdy_ctrt,string"`
    Open           decimal.Decimal `json:"stck_oprc"`
    High           decimal.Decimal `json:"stck_hgpr"`
    Low            decimal.Decimal `json:"stck_lwpr"`
    Volume         int64           `json:"acml_vol,string"`
    // ...
}
```

### 에러 처리 — typed error + sentinel

```go
type APIError struct {
    RtCode  string  // rt_cd
    MsgCode string  // msg_cd
    Message string  // msg1
    TrID    string
}

var (
    ErrTokenExpired = errors.New("kis: token expired")
    ErrRateLimited  = errors.New("kis: rate limited")
    ErrNotFound     = errors.New("kis: resource not found")
)
```

사용자는 `errors.As(err, &apiErr)` 또는 `errors.Is(err, kis.ErrTokenExpired)` 로 분기.

### Rate limit / 토큰 자동 처리

- `internal/ratelimit/`: token bucket, 기본 15 req/sec (goroutine 동시 호출 안전)
- `internal/token/`: 토큰 자동 갱신 (만료 5분 전 선제 발급), file/redis 저장소
- `internal/httpclient/`: resty 기반, 429/5xx 자동 재시도, 토큰 만료 자동 감지/재시도

### HTTP 클라이언트 — resty

| 항목 | 결정 |
|------|------|
| 내부 구현 | `github.com/go-resty/resty/v2` |
| 외부 노출 인터페이스 | 표준 `*http.Client` (`WithHTTPClient(*http.Client)`) |

원칙: **사용자는 resty 를 몰라도 된다.** 사용자가 `WithHTTPClient(customHTTP)` 로 표준 `*http.Client` 를 넘기면, 내부적으로 resty 가 그 transport 를 감싸서 사용. resty 는 구현 디테일.

### 의존성 정책

| 영역 | 선택 |
|------|------|
| Go 버전 | **1.23+** |
| HTTP (internal) | `go-resty/resty/v2` |
| JSON | 표준 `encoding/json` |
| 로깅 | 표준 `log/slog` |
| 테스트 | `stretchr/testify` + `httpmock` |
| 금융 정밀도 | `shopspring/decimal` |
| Redis (선택) | `redis/go-redis/v9` |
| YAML 설정 (선택) | `gopkg.in/yaml.v3` |

원칙: **표준 라이브러리 우선, 외부 의존성 최소**.

### 명명 규칙

- Python `fetch_price` → Go `FetchPrice` (Go 컨벤션)
- KIS 고유 transaction ID 등은 unexported 상수로 모듈 내부에 보관

---

## §4. stock-data-batch → moneyflow 통합 방향성

본 spec 의 직접 scope 는 **korea-investment-stock 의 Python → Go 전환** 이며, 통합 디테일은 별도 spec(`moneyflow 신규 개발`) 에서 다룬다. 여기서는 의존성과 큰 방향만 박아둔다.

### 의존성 방향

```
moneyflow (신규 Go 백엔드, MySQL)
    └── imports → github.com/kenshin579/korea-investment-stock (Go 라이브러리)
```

- moneyflow 가 Go 라이브러리를 단방향 의존. 라이브러리는 moneyflow 를 모름 (재사용성 보호)
- `stock-data-batch` 의 기존 책임(시세 수집, 분류, DB 저장) 은 moneyflow 내부 모듈로 흡수

### DB

- moneyflow 신규 백엔드도 MySQL 사용 → 기존 stock-data-batch (MySQL) 와 정합 이슈 없음
- 스키마 마이그레이션 도구는 moneyflow 새 개발 spec 에서 결정 (현 워크스페이스에서 Liquibase 사용 중)

### 통합 형태 (방향성만)

- moneyflow 백엔드의 별도 CLI 진입점으로 실행 가능 (예: `moneyflow batch sync-stocks`)
- 스케줄링 (cron / k8s CronJob 등) 은 moneyflow 새 개발 spec 영역
- 동일 binary 내 sub-command 또는 별도 cmd 디렉터리 (`cmd/batch/main.go`) 결정도 moneyflow 새 개발 spec 영역

### 본 spec 에서 보장할 것

- Go 라이브러리가 batch 시나리오에 충분한 API 커버리지 (Phase 1 에서 우선순위 산정 시 stock-data-batch 가 사용하는 메서드 매핑 포함)
- 라이브러리 인터페이스가 batch + 단발성 호출 양쪽 모두에 자연스럽도록 설계 (context, rate limit 옵션 등)

---

## §5. 성공 기준 / Phase 1 진입 조건

### Phase 0 의 Deliverable

1. **본 설계 문서** (`docs/superpowers/specs/2026-05-03-korea-investment-go-migration-design.md`) — git 커밋
2. **`korea-investment-stock` repo 의 Step 1 실행 결과**:
   - `v1.x-final` 태그 (또는 정확한 마지막 Python 버전 태그) push
   - PyPI 패키지에 deprecation notice (`pyproject.toml` description, README 상단, 다음 release note)
   - `README.md` 에 마이그레이션 안내
   - main 에서 Python 코드 제거 (git history 보존)
   - `go.mod` 초기화 + 빈 패키지 스켈레톤

### Phase 1 진입 조건

- Phase 0 spec 승인 완료
- Step 1 (Python freeze + tag + PyPI notice + main 정리 + Go 스켈레톤) 완료
- Phase 1 (필요 API 식별 + 우선순위 산정) 시작 의지 확인

### Phase 0 성공 판정

| 항목 | 판정 기준 |
|------|----------|
| 외부 사용자가 v1.x 사용을 계속할 수 있는가 | 태그/PyPI 보존 |
| 외부 사용자가 deprecation 사실을 알 수 있는가 | README + PyPI description + release note |
| 다음 단계(Go 라이브러리 개발) 시작할 수 있는 상태인가 | main 비고 + Go 스켈레톤 준비 |
| Go 라이브러리의 핵심 설계 원칙이 합의되었는가 | 본 spec |

### Non-goals (의도적으로 다루지 않은 것)

- 어떤 API 를 어떤 우선순위로 옮길지 → **Phase 1**
- moneyflow 의 디렉터리/CLI/스케줄링/스키마 디테일 → **moneyflow 신규 개발 spec**
- 실시간 WebSocket 시세 → **추후 별도 spec** (batch 시나리오에 불필요)
- 선물옵션 / 장내채권 → **본 spec 에서 영구 제외**
- 주식 주문/잔고/예약주문 (트레이딩 기능) → **본 spec 에서 다루지 않음**

### 위험 요소와 대응

| 위험 | 대응 |
|------|------|
| 외부 Python 사용자가 라이브러리가 바뀐 줄 모르고 혼란 | README + PyPI description deprecation notice + GitHub Releases 안내 |
| Go 가 처음에는 API 커버리지 0 인 기간 | Phase 1 에서 stock-data-batch 가 실제 쓰는 메서드 우선 구현 → 갭 최소화 |
| Python `v1.x` 의 미해결 버그 / 보안 이슈 | 신규 추가는 안 하지만 critical security fix 는 v1.x patch 허용 (정책 명시) |
| Go 의존성 (resty, decimal 등) 향후 breaking change 발생 | go.sum lock + 필요 시 thin wrapper 로 격리 |
