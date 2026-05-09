# Phase 11.1 — 국내선물옵션 시세/조회 11 EP Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 신규 `futures/` sub-package 도입 + 국내선물옵션 시세/조회 11 EP REST 추가 (`v1.21.0` release). bonds/ (Phase 3.1) 패턴 참조.

**Architecture:** Phase 1+ 인프라 + 패턴 재사용. 신규 `futures/` sub-package 신설 (Client, doc.go). Root `client.go` 에 `Futures *futures.Client` 필드 + `wireInfra` 에서 `bonds.New(...)` 옆에 `futures.New(...)` 추가. TDD: docs analyzer → schemas.md → fixture → 실패 테스트 → struct + 메서드 → 통과 → commit.

**Tech Stack:** Go 1.25+, `github.com/jarcoal/httpmock`, `github.com/stretchr/testify`, `github.com/shopspring/decimal`. 새 dependency 없음.

**참고 spec:**
- Phase 11.1 design spec: `docs/superpowers/specs/2026-05-09-phase11-1-futures-quote-design.md`
- bonds plan (참조 패턴 — 신규 sub-package): `docs/superpowers/specs/2026-05-05-phase3-1-bonds-quote-implementation-plan.md`

---

## 사전 정보

| 항목 | 값 |
|------|---|
| 작업 브랜치 | `feat/phase11-1-futures-quote` (이미 생성, design spec 2 commit 보유) |
| 시작 HEAD | v1.20.0 (Phase 10 완료, 121 REST + 17 WS = 138 endpoints) |
| Release 목표 | `v1.21.0` |
| PR 베이스 | `main` |
| 누적 (완료 후) | 132 REST + 17 WS = **149 endpoints** |

> **IMPORTANT**: #4 (`InquireCcnlBstime`) + #7 (`InquireDailyAmountFee`) 두 EP 는 path 가 `trading/...` 으로 시작 (조회 성격이지만 path 카테고리는 trading). docs analyzer (Task 1) 단계에서 계좌 정보 (CANO/ACNT_PRDT_CD) 가 query param 으로 필요한지 확인 필수. 만약 계좌 정보 필요 시 본 phase 에서 빼고 Phase 11.4 (Trading) 로 미룰 것 — 사용자 확인 받음 후 결정.

## 메서드 매핑 (docs grep 결과, Task 1 에서 정밀 검증)

base path = `/uapi/domestic-futureoption/v1/`

| # | docs (한글 파일명) | path (last segment) | TR_ID | Method (path 기반 PascalCase) | 분류 | 모의 | 비고 |
|---|---|---|---|---|---|---|---|
| 1 | 선물옵션_시세 | `quotations/inquire-price` | FHMIF10000000 | `InquirePrice` | quote.go | ? | 현재가 |
| 2 | 선물옵션_시세호가 | `quotations/inquire-asking-price` | FHMIF10010000 | `InquireAskingPrice` | quote.go | ? | 시세 + 호가 |
| 3 | 선물옵션_분봉조회 | `quotations/inquire-time-fuopchartprice` | FHKIF03020200 | `InquireTimeFuopchartprice` | chart.go | ? | 분봉 |
| 4 | 선물옵션_기준일체결내역 | `trading/inquire-ccnl-bstime` | CTFO5139R | `InquireCcnlBstime` | conclusion.go | 미지원 | trading path — 계좌 정보 검증 |
| 5 | 선물옵션_일중예상체결추이 | `quotations/exp-price-trend` | FHPIF05110100 | `ExpPriceTrend` | conclusion.go | ? | 일중 예상체결 |
| 6 | 선물옵션기간별시세(일/주/월/년) | `quotations/inquire-daily-fuopchartprice` | FHKIF03020100 | `InquireDailyFuopchartprice` | chart.go | ? | 일/주/월/년 차트 |
| 7 | 선물옵션기간약정수수료일별 | `trading/inquire-daily-amount-fee` | CTFO6119R | `InquireDailyAmountFee` | conclusion.go | 미지원 | trading path — 계좌 정보 검증 |
| 8 | 국내선물_기초자산_시세 | `quotations/display-board-top` | FHPIF05030000 | `DisplayBoardTop` | board.go | ? | 전광판 top (제목 vs path 차이) |
| 9 | 국내옵션전광판_선물 | `quotations/display-board-futures` | FHPIF05030200 | `DisplayBoardFutures` | board.go | ? | 전광판 선물 |
| 10 | 국내옵션전광판_옵션월물리스트 | `quotations/display-board-option-list` | FHPIO056104C0 | `DisplayBoardOptionList` | board.go | ? | 월물 리스트 |
| 11 | 국내옵션전광판_콜풋 | `quotations/display-board-callput` | FHPIF05030100 | `DisplayBoardCallput` | board.go | ? | 전광판 콜풋 |

> **메서드명 변경 사유 (spec §2 의 임시명 vs path 기반)**: spec 의 임시명 일부 (예: `InquireUnderlyingPrice`, `OptionBoardFuture`, `OptionMonthlyList`, `OptionBoardCallPut`) 가 path 와 안 맞아 정정. Style A (path 기반 PascalCase) 일관성. 본 plan 의 메서드명이 source of truth.

---

## 파일 구조

### 신규 (futures package)

- `futures/client.go` — `Client` struct + `New(http *httpclient.Client) *Client`
- `futures/doc.go` — package doc (11 메서드 목록)
- `futures/quote.go` — 2 메서드 (InquirePrice / InquireAskingPrice)
- `futures/quote_test.go`
- `futures/chart.go` — 2 메서드 (InquireTimeFuopchartprice / InquireDailyFuopchartprice)
- `futures/chart_test.go`
- `futures/conclusion.go` — 3 메서드 (InquireCcnlBstime / ExpPriceTrend / InquireDailyAmountFee)
- `futures/conclusion_test.go`
- `futures/board.go` — 4 메서드 (DisplayBoardTop / DisplayBoardFutures / DisplayBoardOptionList / DisplayBoardCallput)
- `futures/board_test.go`
- `futures/testhelper_test.go` — `newTestClient(t)` + `loadFixtureString(t, name)` (bonds 패턴)

### 신규 (testdata)

- `futures/testdata/_schemas.md` — Task 1 의 11 EP 정확한 schema reference
- `futures/testdata/inquire_price_success.json`
- `futures/testdata/inquire_asking_price_success.json`
- `futures/testdata/inquire_time_fuopchartprice_success.json`
- `futures/testdata/inquire_ccnl_bstime_success.json`
- `futures/testdata/exp_price_trend_success.json`
- `futures/testdata/inquire_daily_fuopchartprice_success.json`
- `futures/testdata/inquire_daily_amount_fee_success.json`
- `futures/testdata/display_board_top_success.json`
- `futures/testdata/display_board_futures_success.json`
- `futures/testdata/display_board_option_list_success.json`
- `futures/testdata/display_board_callput_success.json`

### 신규 (examples)

- `examples/futures_basic/main.go` — 4-5 메서드 통합 시연

### 수정

- 기존 `client.go` (root) — `Futures *futures.Client` 필드 + `wireInfra` 에 `c.Futures = futures.New(c.httpClient)` 추가
- `CLAUDE.md` — banner Phase 10 → Phase 11.1, plan/spec link 추가
- `README.md` — Available Methods 표 갱신, futures section 추가
- `CHANGELOG.md` — `[1.21.0]` entry ABOVE `[1.20.0]`

---

## 핵심 컨벤션

- **type 매핑** (Phase 1 부터 일관):
  - 가격/체결가/액면가/호가 → `decimal.Decimal`
  - 거래량/거래대금/수량 → `int64,string`
  - 비율/등락율/체결강도 → `float64,string`
  - 코드/Y-N/날짜/시간 → `string`
- **JSON 태그**: 한투 API 원본 그대로 (`stck_prpr` 등)
- **메서드 시그니처 패턴** (단순):
  ```go
  func (c *Client) InquirePrice(ctx context.Context, code string) (*InquirePriceData, error)
  ```
- **메서드 시그니처 패턴** (Params struct):
  ```go
  func (c *Client) InquireDailyFuopchartprice(ctx context.Context, params InquireDailyFuopchartpriceParams) (*InquireDailyFuopchartpriceData, error)
  ```
- **Struct 명**: `<Method>Data` (응답) / `<Method>Params` (요청 옵션) — bonds 패턴 일관
- **Default `FID_COND_MRKT_DIV_CODE`**: Task 1 docs analyzer 에서 EP 별 확정 (`F`/`JF`/`O`/`U` 가능성)

---

## Task 1: docs analyzer + `futures/testdata/_schemas.md` 작성

**Files:**
- Create: `futures/testdata/_schemas.md`

11 EP 의 docs (`docs/api/국내선물옵션/<docs 파일>.md`) 응답 body 표를 직접 읽어 정확한 매핑 추출. 본 task 의 출력물이 후속 task 의 source of truth.

- [ ] **Step 1: 각 EP 별 docs 응답 body 표 직접 분석**

각 EP 의 docs 에서 다음 추출:
1. `## Body` (요청) 표 → Query params (FID_COND_*, 종목코드 필드명, 등)
2. `## 응답 → ## Body` 표 → output 키 (`output`/`output1`/`output2`+ array vs single) + 필드 (이름 / type / Length / 한글명)
3. 모의투자 지원 여부 (header 영역 명시)
4. CANO/ACNT_PRDT_CD 같은 계좌 정보 query 가 있는지 (#4, #7 검증 critical)

- [ ] **Step 2: `_schemas.md` 작성**

각 EP 별로 다음 형식:

```markdown
## EP1 — `InquirePrice` (선물옵션 시세, FHMIF10000000)

**path**: `/uapi/domestic-futureoption/v1/quotations/inquire-price`
**모의**: 지원/미지원
**계좌 정보 필요**: 없음/CANO+ACNT_PRDT_CD

### Query
| Key | Type | Required | Default | 한글명 |
|---|---|---|---|---|
| FID_COND_MRKT_DIV_CODE | string | Y | "F" | 시장구분 |
| FID_INPUT_ISCD | string | Y | - | 종목코드 |
| ... | | | | |

### Output (`output{}` single 또는 `output[]` array)
| # | Key | Go field | Type | 한글명 |
|---|---|---|---|---|
| 1 | stck_prpr | Price | decimal | 현재가 |
| 2 | ... | | | |

### 비고
- ...
```

11 EP 모두 위 형식.

- [ ] **Step 3: 계좌 정보 필요 EP 사용자 confirm**

#4 (`InquireCcnlBstime`) + #7 (`InquireDailyAmountFee`) 가 CANO/ACNT_PRDT_CD 를 query 로 받는지 확인. 만약 받는다면 실제 계좌 인증 필요 — 본 phase scope (시세/조회) 와 충돌. 사용자에게 보고:
- 옵션 A: 본 phase 에서 빼고 Phase 11.4 Trading 으로 미룸 (10 EP 로 축소)
- 옵션 B: 본 phase 에 포함하되 caller 가 계좌 정보 입력하도록 Params 명시

- [ ] **Step 4: commit**

```bash
git add futures/testdata/_schemas.md
git commit -m "[docs] Phase 11.1 — futures/testdata/_schemas.md (11 EP schema reference)"
```

---

## Task 2: futures/ scaffolding (Client + doc + testhelper)

**Files:**
- Create: `futures/client.go`
- Create: `futures/doc.go`
- Create: `futures/testhelper_test.go`
- Modify: `client.go` (root)

bonds/ 패턴 그대로 복사 + 명칭 치환.

- [ ] **Step 1: `futures/client.go` 작성**

```go
package futures

import "github.com/kenshin579/korea-investment-stock/internal/httpclient"

// Client 는 국내선물옵션 도메인 진입점.
//
// kis.NewClient(...) 가 자동으로 wireInfra 에서 New(httpClient) 호출하여
// kis.Client.Futures 필드에 주입한다. 직접 인스턴스화 불필요.
type Client struct {
	http *httpclient.Client
}

// New 는 새 Futures Client 를 생성한다.
func New(http *httpclient.Client) *Client {
	return &Client{http: http}
}
```

- [ ] **Step 2: `futures/doc.go` 작성**

```go
// Package futures 는 KIS 국내선물옵션 도메인.
//
// 한투 docs: docs/api/국내선물옵션/*.md
//
// Phase 11.1 — 시세/조회 11 endpoint:
//
//   FHMIF10000000  선물옵션 시세                  InquirePrice
//   FHMIF10010000  선물옵션 시세호가              InquireAskingPrice
//   FHKIF03020200  선물옵션 분봉                  InquireTimeFuopchartprice
//   CTFO5139R      선물옵션 기준일체결내역        InquireCcnlBstime
//   FHPIF05110100  선물옵션 일중예상체결추이      ExpPriceTrend
//   FHKIF03020100  선물옵션 기간별시세 (일/주/월/년) InquireDailyFuopchartprice
//   CTFO6119R      선물옵션 기간약정수수료 일별   InquireDailyAmountFee
//   FHPIF05030000  국내선물 전광판 top (기초자산) DisplayBoardTop
//   FHPIF05030200  옵션 전광판 선물               DisplayBoardFutures
//   FHPIO056104C0  옵션 전광판 옵션월물리스트     DisplayBoardOptionList
//   FHPIF05030100  옵션 전광판 콜풋               DisplayBoardCallput
//
// 사용자는 root kis.Client 의 Futures 필드로 접근.
package futures
```

- [ ] **Step 3: `futures/testhelper_test.go` 작성**

```go
package futures

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/kenshin579/korea-investment-stock/internal/httpclient"
	"github.com/stretchr/testify/require"
)

// newTestClient 는 httpmock 가 가로챈 transport 로 stub Client 를 생성한다.
func newTestClient(t *testing.T) *Client {
	t.Helper()
	hc := httpclient.NewForTest(nil)
	return New(hc)
}

// loadFixtureString 은 testdata 파일 내용을 string 으로 반환.
func loadFixtureString(t *testing.T, name string) string {
	t.Helper()
	b, err := os.ReadFile(filepath.Join("testdata", name))
	require.NoError(t, err)
	return string(b)
}
```

- [ ] **Step 4: root `client.go` 에 Futures 필드 추가**

기존 `client.go` 의 `Client` struct 에 `Bonds *bonds.Client` 옆에 `Futures *futures.Client` 추가. `wireInfra` (또는 NewClient) 에서 `c.Bonds = bonds.New(...)` 옆에 `c.Futures = futures.New(c.httpClient)` 추가.

```bash
grep -n "Bonds " client.go
# 결과 확인 후 정확한 위치에 Futures 라인 추가
```

import 추가:
```go
"github.com/kenshin579/korea-investment-stock/futures"
```

- [ ] **Step 5: 빌드 확인**

```bash
go build ./...
```

Expected: 빌드 성공 (메서드 0개 상태 OK).

- [ ] **Step 6: commit**

```bash
git add futures/client.go futures/doc.go futures/testhelper_test.go client.go
git commit -m "[chore] Phase 11.1 — futures/ scaffolding (Client, doc, testhelper) + root Client.Futures"
```

---

## Task 3: EP1 — `InquirePrice` (선물옵션 시세)

**Files:**
- Create: `futures/testdata/inquire_price_success.json`
- Create: `futures/quote.go`
- Create: `futures/quote_test.go`

**Source of truth**: `futures/testdata/_schemas.md` 의 EP1 섹션. 본 task 의 verbatim Go struct 필드 매핑은 schemas.md 참조.

- [ ] **Step 1: testdata fixture 작성**

`futures/testdata/inquire_price_success.json`. KIS 응답 envelope (`rt_cd`/`msg_cd`/`msg1`/`output`) 포함 + EP1 의 정확한 필드 (schemas.md 참조).

```json
{
  "rt_cd": "0",
  "msg_cd": "MCA00000",
  "msg1": "정상처리 되었습니다",
  "output": {
    "iscd_stat_cls_code": "55",
    ...
  }
}
```

- [ ] **Step 2: failing test 작성**

`futures/quote_test.go`:

```go
package futures

import (
	"context"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInquirePrice(t *testing.T) {
	c := newTestClient(t)
	httpmock.ActivateNonDefault(c.http.Resty().GetClient())
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(http.MethodGet,
		`=~/uapi/domestic-futureoption/v1/quotations/inquire-price`,
		httpmock.NewStringResponder(200, loadFixtureString(t, "inquire_price_success.json")))

	got, err := c.InquirePrice(context.Background(), "101W3000")
	require.NoError(t, err)
	require.NotNil(t, got)
	// schemas.md 의 정확한 필드 검증 (가격/거래량 등 1-2개 strong assertion)
	assert.False(t, got.Price.IsZero())
}
```

- [ ] **Step 3: 테스트 실행 (FAIL 확인)**

```bash
go test ./futures/ -run TestInquirePrice -v
```

Expected: FAIL — `c.InquirePrice undefined`.

- [ ] **Step 4: 메서드 구현**

`futures/quote.go`:

```go
package futures

import (
	"context"

	"github.com/shopspring/decimal"
)

// InquirePriceData 는 InquirePrice (FHMIF10000000) 응답.
// 정확한 필드는 docs/api/국내선물옵션/선물옵션_시세.md 응답 body 표 + futures/testdata/_schemas.md EP1 참조.
type InquirePriceData struct {
	IscdStatClsCode string          `json:"iscd_stat_cls_code"` // 종목 상태
	Price           decimal.Decimal `json:"stck_prpr"`          // 현재가 (예시 — 실제 필드명 schemas.md 검증)
	// ... schemas.md 의 EP1 output 필드 모두 매핑
}

// InquirePrice 는 국내선물옵션 시세 (FHMIF10000000) 를 조회한다.
//
// FID_COND_MRKT_DIV_CODE default = "F" (선물; 옵션은 "JF"/"O" 가능성 — schemas.md 검증).
// 종목코드 형식: 9자리 alphanumeric (예: "101W3000" 선물, "201X3300" 옵션).
func (c *Client) InquirePrice(ctx context.Context, code string) (*InquirePriceData, error) {
	var resp struct {
		Output InquirePriceData `json:"output"`
	}
	err := c.http.Get(ctx,
		"/uapi/domestic-futureoption/v1/quotations/inquire-price",
		"FHMIF10000000",
		map[string]string{
			"FID_COND_MRKT_DIV_CODE": "F", // schemas.md EP1 검증
			"FID_INPUT_ISCD":         code,
		},
		&resp)
	if err != nil {
		return nil, err
	}
	return &resp.Output, nil
}
```

- [ ] **Step 5: 테스트 실행 (PASS 확인)**

```bash
go test ./futures/ -run TestInquirePrice -v
```

Expected: PASS.

- [ ] **Step 6: commit**

```bash
git add futures/testdata/inquire_price_success.json futures/quote.go futures/quote_test.go
git commit -m "[feat] Phase 11.1 — Futures.InquirePrice (FHMIF10000000)"
```

---

## Task 4: EP2 — `InquireAskingPrice` (선물옵션 시세호가)

**Files:**
- Create: `futures/testdata/inquire_asking_price_success.json`
- Modify: `futures/quote.go` (append 메서드 + struct)
- Modify: `futures/quote_test.go` (append 테스트)

**Source of truth**: `_schemas.md` EP2.

- [ ] **Step 1: fixture 작성** — schemas.md EP2 의 정확한 output 필드 + envelope.
- [ ] **Step 2: failing test 작성** — `TestInquireAskingPrice` 함수 (Task 3 패턴 그대로 — path 만 `inquire-asking-price` 로).
- [ ] **Step 3: 실행 (FAIL)** — `go test ./futures/ -run TestInquireAskingPrice -v`.
- [ ] **Step 4: 구현** — `InquireAskingPriceData` struct + `InquireAskingPrice(ctx, code)` 메서드. TR_ID `FHMIF10010000`. 5단계 호가 (선물옵션) 또는 다른 단계 — schemas.md 검증.
- [ ] **Step 5: 실행 (PASS)** — `go test ./futures/ -run TestInquireAskingPrice -v`.
- [ ] **Step 6: commit** — `[feat] Phase 11.1 — Futures.InquireAskingPrice (FHMIF10010000)`.

---

## Task 5: EP3 — `InquireTimeFuopchartprice` (분봉)

**Files:**
- Create: `futures/testdata/inquire_time_fuopchartprice_success.json`
- Create: `futures/chart.go`
- Create: `futures/chart_test.go`

**Source of truth**: `_schemas.md` EP3.

- [ ] **Step 1: fixture 작성** — output 이 array (분봉) 일 가능성 높음. envelope + `output1{}` 헤더 + `output2[]` 분봉 array (Phase 1.2 chart 패턴).
- [ ] **Step 2: failing test (`TestInquireTimeFuopchartprice`)** — Params struct (시간/구분 등) + happy path 검증.
- [ ] **Step 3: 실행 (FAIL)**.
- [ ] **Step 4: 구현** — `InquireTimeFuopchartpriceParams` (Period, MarketCode, Code, Time, IncludePastData, etc — schemas.md 검증) + `InquireTimeFuopchartpriceData` struct (output1 + output2 분봉) + 메서드. TR_ID `FHKIF03020200`.
- [ ] **Step 5: 실행 (PASS)**.
- [ ] **Step 6: commit** — `[feat] Phase 11.1 — Futures.InquireTimeFuopchartprice (FHKIF03020200, 분봉)`.

---

## Task 6: EP4 — `InquireCcnlBstime` (기준일 체결내역)

**Files:**
- Create: `futures/testdata/inquire_ccnl_bstime_success.json`
- Create: `futures/conclusion.go`
- Create: `futures/conclusion_test.go`

**Source of truth**: `_schemas.md` EP4. **계좌 정보 필요 시 Task 1 step 3 의 사용자 결정 따름.**

- [ ] **Step 1: 계좌 정보 분기**
  - 계좌 정보 미필요 → 일반 patten (Task 3 처럼)
  - 계좌 정보 필요 + 사용자가 옵션 A 선택 → 본 task 와 Task 9 (#7) skip (10 EP 로 축소)
  - 옵션 B → Params 에 CANO/ACNT_PRDT_CD 추가 + Task 진행
- [ ] **Step 2: fixture 작성** — schemas.md EP4 의 정확한 output. CTFO5139R TR_ID 명시. 기준일 (BASS_DT) 파라미터.
- [ ] **Step 3: failing test (`TestInquireCcnlBstime`)**.
- [ ] **Step 4: 실행 (FAIL)**.
- [ ] **Step 5: 구현** — `InquireCcnlBstimeParams` (BASS_DT 기준일자 등) + `InquireCcnlBstimeData` + 메서드. TR_ID `CTFO5139R`. trading path 명시.
- [ ] **Step 6: 실행 (PASS)**.
- [ ] **Step 7: commit** — `[feat] Phase 11.1 — Futures.InquireCcnlBstime (CTFO5139R, 기준일체결)`.

---

## Task 7: EP5 — `ExpPriceTrend` (일중 예상체결추이)

**Files:**
- Create: `futures/testdata/exp_price_trend_success.json`
- Modify: `futures/conclusion.go`
- Modify: `futures/conclusion_test.go`

**Source of truth**: `_schemas.md` EP5.

- [ ] **Step 1: fixture 작성** — output array (시간대별 추이) 가능성. TR_ID `FHPIF05110100`.
- [ ] **Step 2: failing test (`TestExpPriceTrend`)**.
- [ ] **Step 3: 실행 (FAIL)**.
- [ ] **Step 4: 구현** — `ExpPriceTrendParams` + `ExpPriceTrendData` (output1 헤더 + output2 array 가능) + 메서드.
- [ ] **Step 5: 실행 (PASS)**.
- [ ] **Step 6: commit** — `[feat] Phase 11.1 — Futures.ExpPriceTrend (FHPIF05110100, 일중 예상체결)`.

---

## Task 8: EP6 — `InquireDailyFuopchartprice` (기간별 시세 일/주/월/년)

**Files:**
- Create: `futures/testdata/inquire_daily_fuopchartprice_success.json`
- Modify: `futures/chart.go`
- Modify: `futures/chart_test.go`

**Source of truth**: `_schemas.md` EP6.

- [ ] **Step 1: fixture** — Period 별 (일/주/월/년) 차트. output1 헤더 + output2 array (Phase 1.2 chart 패턴).
- [ ] **Step 2: failing test (`TestInquireDailyFuopchartprice`)** — Params (Period: D/W/M/Y, 시작일/종료일).
- [ ] **Step 3: 실행 (FAIL)**.
- [ ] **Step 4: 구현** — `InquireDailyFuopchartpriceParams` + `InquireDailyFuopchartpriceData` + 메서드. TR_ID `FHKIF03020100`.
- [ ] **Step 5: 실행 (PASS)**.
- [ ] **Step 6: commit** — `[feat] Phase 11.1 — Futures.InquireDailyFuopchartprice (FHKIF03020100, 일/주/월/년)`.

---

## Task 9: EP7 — `InquireDailyAmountFee` (수수료 일별)

**Files:**
- Create: `futures/testdata/inquire_daily_amount_fee_success.json`
- Modify: `futures/conclusion.go`
- Modify: `futures/conclusion_test.go`

**Source of truth**: `_schemas.md` EP7. **계좌 정보 필요 시 Task 1 step 3 결정 따름.**

- [ ] **Step 1: 계좌 정보 분기** (Task 6 동일).
- [ ] **Step 2: fixture** — TR_ID `CTFO6119R`. 일별 수수료 array.
- [ ] **Step 3: failing test (`TestInquireDailyAmountFee`)**.
- [ ] **Step 4: 실행 (FAIL)**.
- [ ] **Step 5: 구현** — Params + Data struct + 메서드.
- [ ] **Step 6: 실행 (PASS)**.
- [ ] **Step 7: commit** — `[feat] Phase 11.1 — Futures.InquireDailyAmountFee (CTFO6119R, 수수료 일별)`.

---

## Task 10: EP8 — `DisplayBoardTop` (전광판 top, 기초자산 시세)

**Files:**
- Create: `futures/testdata/display_board_top_success.json`
- Create: `futures/board.go`
- Create: `futures/board_test.go`

**Source of truth**: `_schemas.md` EP8. docs 제목 ("기초자산 시세") 와 path ("display-board-top") 차이 — schemas.md 의 실제 응답 필드를 따름.

- [ ] **Step 1: fixture** — TR_ID `FHPIF05030000`. 전광판 top 종목 array 가능성 (output array).
- [ ] **Step 2: failing test (`TestDisplayBoardTop`)**.
- [ ] **Step 3: 실행 (FAIL)**.
- [ ] **Step 4: 구현** — `DisplayBoardTopParams` + `DisplayBoardTopData` + 메서드.
- [ ] **Step 5: 실행 (PASS)**.
- [ ] **Step 6: commit** — `[feat] Phase 11.1 — Futures.DisplayBoardTop (FHPIF05030000, 전광판 top)`.

---

## Task 11: EP9 — `DisplayBoardFutures` (전광판 선물)

**Files:**
- Create: `futures/testdata/display_board_futures_success.json`
- Modify: `futures/board.go`
- Modify: `futures/board_test.go`

**Source of truth**: `_schemas.md` EP9.

- [ ] **Step 1: fixture** — TR_ID `FHPIF05030200`. 선물 종목 array.
- [ ] **Step 2: failing test (`TestDisplayBoardFutures`)**.
- [ ] **Step 3: 실행 (FAIL)**.
- [ ] **Step 4: 구현**.
- [ ] **Step 5: 실행 (PASS)**.
- [ ] **Step 6: commit** — `[feat] Phase 11.1 — Futures.DisplayBoardFutures (FHPIF05030200, 전광판 선물)`.

---

## Task 12: EP10 — `DisplayBoardOptionList` (옵션 월물 리스트)

**Files:**
- Create: `futures/testdata/display_board_option_list_success.json`
- Modify: `futures/board.go`
- Modify: `futures/board_test.go`

**Source of truth**: `_schemas.md` EP10.

- [ ] **Step 1: fixture** — TR_ID `FHPIO056104C0` (특이 — `O` prefix 보임. schemas.md 검증). 월물 array.
- [ ] **Step 2: failing test (`TestDisplayBoardOptionList`)**.
- [ ] **Step 3: 실행 (FAIL)**.
- [ ] **Step 4: 구현**.
- [ ] **Step 5: 실행 (PASS)**.
- [ ] **Step 6: commit** — `[feat] Phase 11.1 — Futures.DisplayBoardOptionList (FHPIO056104C0, 옵션 월물)`.

---

## Task 13: EP11 — `DisplayBoardCallput` (전광판 콜풋)

**Files:**
- Create: `futures/testdata/display_board_callput_success.json`
- Modify: `futures/board.go`
- Modify: `futures/board_test.go`

**Source of truth**: `_schemas.md` EP11.

- [ ] **Step 1: fixture** — TR_ID `FHPIF05030100`. 콜/풋 dual output 가능성 (output1=콜, output2=풋 또는 array).
- [ ] **Step 2: failing test (`TestDisplayBoardCallput`)**.
- [ ] **Step 3: 실행 (FAIL)**.
- [ ] **Step 4: 구현** — multi-output 가능성 schemas.md 정확 매핑.
- [ ] **Step 5: 실행 (PASS)**.
- [ ] **Step 6: commit** — `[feat] Phase 11.1 — Futures.DisplayBoardCallput (FHPIF05030100, 전광판 콜풋)`.

---

## Task 14: Invalid-JSON 보강 테스트 (coverage)

**Files:**
- Modify: `futures/quote_test.go`, `chart_test.go`, `conclusion_test.go`, `board_test.go`

Phase 3.1 lesson: 신규 sub-package 도입 시 happy path 만 ~67-71% coverage. 80% 충족 위해 InvalidJSON test 추가.

- [ ] **Step 1: 각 파일 마다 1 InvalidJSON 테스트 추가**

```go
func TestInquirePrice_InvalidJSON(t *testing.T) {
	c := newTestClient(t)
	httpmock.ActivateNonDefault(c.http.Resty().GetClient())
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(http.MethodGet,
		`=~/uapi/domestic-futureoption/v1/quotations/inquire-price`,
		httpmock.NewStringResponder(200, `{"rt_cd":"0","msg_cd":"MCA00000","output":"not-an-object"}`))

	_, err := c.InquirePrice(context.Background(), "101W3000")
	require.Error(t, err)
}
```

11 메서드 모두 (또는 file 별 1개) 동일 패턴. 4 파일 × 1 = 4 테스트 충분 (Phase 3.1 패턴).

- [ ] **Step 2: 실행**

```bash
go test ./futures/ -count=1
```

Expected: 모두 PASS.

- [ ] **Step 3: coverage 측정**

```bash
go test ./futures/ -coverprofile=/tmp/futures_cov.out -count=1
go tool cover -func=/tmp/futures_cov.out | tail -1
```

Expected: ≥ 80%.

- [ ] **Step 4: commit**

```bash
git add futures/quote_test.go futures/chart_test.go futures/conclusion_test.go futures/board_test.go
git commit -m "[test] Phase 11.1 — InvalidJSON 보강 (coverage ≥80%)"
```

---

## Task 15: example `examples/futures_basic/`

**Files:**
- Create: `examples/futures_basic/main.go`

4-5 메서드 통합 시연 (InquirePrice + InquireAskingPrice + InquireDailyFuopchartprice + DisplayBoardTop 정도).

- [ ] **Step 1: example 작성**

`examples/futures_basic/main.go`:

```go
// examples/futures_basic/main.go — Phase 11.1 국내선물옵션 시세 시연.
//
// 4 메서드 통합:
//   - Futures.InquirePrice (현재가)
//   - Futures.InquireAskingPrice (시세호가)
//   - Futures.InquireDailyFuopchartprice (일별 차트)
//   - Futures.DisplayBoardTop (전광판 top)
//
// 종목코드: 9자리 alphanumeric (예: "101W3000" 선물, "201X3300" 옵션).
//
// Run:
//   KOREA_INVESTMENT_APP_KEY=... KOREA_INVESTMENT_APP_SECRET=... go run ./examples/futures_basic
package main

import (
	"context"
	"fmt"
	"log"

	kis "github.com/kenshin579/korea-investment-stock"
)

func main() {
	client, err := kis.NewClientFromEnv()
	if err != nil {
		log.Fatalf("NewClientFromEnv: %v", err)
	}
	ctx := context.Background()

	// Task 1 의 schemas.md 결과로 예시 종목코드 결정 (선물 KOSPI200 근월물 등)
	const code = "101W3000" // 예시 — 실제 호출 시 caller 가 정확한 코드 입력

	price, err := client.Futures.InquirePrice(ctx, code)
	if err != nil {
		log.Fatalf("InquirePrice: %v", err)
	}
	fmt.Printf("[현재가] %s 현재가=%s\n", code, price.Price.String())

	// 추가 메서드도 동일 패턴
	// ...
}
```

- [ ] **Step 2: 빌드 확인**

```bash
go build ./examples/futures_basic/
rm -f futures_basic
```

- [ ] **Step 3: commit**

```bash
git add examples/futures_basic/main.go
git commit -m "[example] Phase 11.1 — examples/futures_basic (4 메서드 시연)"
```

---

## Task 16: 문서 갱신 + 최종 점검

**Files:**
- Modify: `CLAUDE.md`, `README.md`, `CHANGELOG.md`

- [ ] **Step 1: `CLAUDE.md` banner + spec/plan link**

```markdown
> **Phase 11.1 — 국내선물옵션 시세/조회 11 EP (v1.21.0). 누적 132 REST + 17 WS = 149 endpoints.**
```

list 에 추가:
```markdown
- Phase 11.1 design spec: [`docs/superpowers/specs/2026-05-09-phase11-1-futures-quote-design.md`](...)
- Phase 11.1 implementation plan: [`docs/superpowers/plans/2026-05-09-phase11-1-futures-quote.md`](...)
```

- [ ] **Step 2: `README.md` Available Methods 표 + futures section**

futures section 추가 (bonds 섹션 아래):

```markdown
### Futures (국내선물옵션) — Phase 11.1

| Go 메서드 | path | TR_ID |
|---|---|---|
| `Futures.InquirePrice` | `quotations/inquire-price` | FHMIF10000000 |
| `Futures.InquireAskingPrice` | `quotations/inquire-asking-price` | FHMIF10010000 |
| ... (11 EP 모두) |
```

- [ ] **Step 3: `CHANGELOG.md` `[1.21.0]` entry**

`[1.20.0]` 위에 새 entry:

```markdown
## [1.21.0] - 2026-05-09

### Added — Phase 11.1 (국내선물옵션 시세/조회 11 EP, 신규 `futures/` sub-package)

- `Futures.InquirePrice` (FHMIF10000000) — 선물옵션 시세
- `Futures.InquireAskingPrice` (FHMIF10010000) — 시세호가
- ... (11 EP 모두)
- examples: `futures_basic`

### Notes
- 신규 sub-package `futures/` (bonds 패턴). root `Client.Futures *futures.Client`.
- 종목코드 9자리 alphanumeric (KRX 6자리와 다름).
- 누적 121 → 132 REST + 17 WS = 149 endpoints.
```

- [ ] **Step 4: 최종 점검**

```bash
gofmt -w .
gofmt -l .  # empty
go vet ./...
go build ./...
go test ./... -race -count=1
go test ./futures/ -coverprofile=/tmp/cov.out -count=1
go tool cover -func=/tmp/cov.out | tail -1
```

Expected: 모두 clean. coverage ≥ 80%.

- [ ] **Step 5: commit**

```bash
git add CLAUDE.md README.md CHANGELOG.md
git commit -m "[docs] Phase 11.1 — CLAUDE/README/CHANGELOG (v1.21.0)"
```

---

## Task 17: PR + merge + tag v1.21.0 + GitHub Release

- [ ] **Step 1: push branch**

```bash
git push -u origin feat/phase11-1-futures-quote
```

- [ ] **Step 2: PR 생성** — `gh pr create` + HEREDOC body. Title: `Phase 11.1 — 국내선물옵션 시세/조회 11 EP (v1.21.0)`. Reviewer: kenshin579.

- [ ] **Step 3: merge** — `gh pr merge <num> --squash --delete-branch`.

- [ ] **Step 4: main sync + tag**

```bash
git checkout main && git pull origin main
git tag -a v1.21.0 -m "v1.21.0: Phase 11.1 — 국내선물옵션 시세/조회 11 EP"
git push origin v1.21.0
```

- [ ] **Step 5: GitHub Release**

```bash
awk '/^## \[/{c++} c==1' CHANGELOG.md > /tmp/release_v1.21.0.md
gh release create v1.21.0 --title "v1.21.0 — Phase 11.1 국내선물옵션 시세/조회 11 EP" --notes-file /tmp/release_v1.21.0.md
```

- [ ] **Step 6: Go proxy verify**

```bash
GOPROXY=direct go list -m github.com/kenshin579/korea-investment-stock@latest
```

Expected: `v1.21.0`.

- [ ] **Step 7: 메모리 갱신** — `~/.claude/projects/.../memory/kis_go_migration.md` 의 Phase 11.1 entry 추가.

---

## 진행 절차 요약

총 17 task. 분량:
- Task 1 (docs analyzer): plan 의 가장 무거운 task — 11 docs 직접 분석
- Task 2 (scaffolding): bonds 패턴 복사
- Task 3-13 (11 EP × TDD): 각 ~5-6 step, 합쳐서 60+ step
- Task 14 (coverage 보강), 15 (example), 16 (문서), 17 (release)

**예상 진행 시간**: 1-2 세션 (집중 시 1, 분산 시 2).
**위험 요인**:
- Task 1 step 3 의 계좌 정보 필요 EP 사용자 결정 — early checkpoint.
- multi-output endpoint (#10, #11) schemas 정확 매핑.
- 종목코드 형식 (선물 vs 옵션 9자리) example 작성 시 정확한 활성 종목 코드 입력.
