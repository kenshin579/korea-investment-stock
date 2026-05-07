# Phase 4.2 — 시장운영/특수상태 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** korea-investment-stock Go 라이브러리에 국내주식 시장운영/특수상태 4 메서드 추가 (`v1.13.0` release). `domestic/market_op.go` 신규 생성 (4 메서드). TDD: testdata fixture → 실패 테스트 → struct + 메서드 구현 → 통과 → commit.

**Architecture:** Phase 1+2 인프라 + 패턴 재사용. `domestic/market_op.go` 신규 생성 (4 메서드). 새 internal package 불필요.

**Tech Stack:** Go 1.25+, `github.com/jarcoal/httpmock`, `github.com/stretchr/testify`, `github.com/shopspring/decimal`. 새 dependency 없음.

**WebSocket scope deviation (CRITICAL):**
Phase 4 design spec §Phase 4.2 에는 7 메서드가 나열되어 있었으나, 그 중 3개 (`InquireMarketOpInfoKrx` / `InquireMarketOpInfoNxt` / `InquireMarketOpInfoTotal` — TR_ID: H0STMKO0 / H0NXMKO0 / H0UNMKO0) 는 **REST GET이 아닌 WebSocket push API**로 확인되었다. 이 라이브러리는 read-only REST 전용 (CLAUDE.md: "실시간 WebSocket" 제외)이므로 Phase 4.2 실제 범위는 **4 REST 메서드**만이다. 해당 3개 WebSocket endpoint 는 잠재적 Phase 5 (WebSocket) 로 이연.

- Phase 4 전체 메서드 수: 30 → **27** (WebSocket 3개 제외)
- Phase 4.2 후 누적: 89 → **93**

**참고 spec:**
- Phase 4 design spec: `docs/superpowers/specs/2026-05-07-phase4-stock-info-design.md` (§Phase 4.2)
- Phase 4.1 plan (compact task structure): `docs/superpowers/specs/2026-05-07-phase4-1-stock-info-implementation-plan.md`
- 기존 구현 참조: `domestic/extended.go`, `domestic/opinion.go`

---

## 사전 정보

| 항목 | 값 |
|------|---|
| 작업 브랜치 | `feat/phase4-2-market-op` |
| 시작 HEAD | Phase 4.1 구현 완료 (v1.12.0) |
| Release 목표 | `v1.13.0` |
| PR 베이스 | `main` |
| 현재 main 상태 | v1.12.0 publish 완료 (Phase 4.1 통합, 89 메서드) |

> **누적 메서드 카운트:** 89 → **93** (4 신규).

---

## 메서드 매핑

| # | Go 메서드 | Path | TR_ID | output key | fields | anomalies |
|---|---|---|---|---|---|---|
| EP4 | `InquireExpClosingPrice` | `quotations/exp-closing-price` | FHKST117300C0 | `output1 []` | 9 | `output1` (not `output`), `FID_COND_SCR_DIV_CODE="11173"` hardcoded, `FID_INPUT_ISCD` = 시장구분코드 (NOT 종목코드) |
| EP5 | `InquireChkHoliday` | `quotations/chk-holiday` | CTCA0903R | `output {}` | 6 | **non-FID UPPERCASE params** (BASS_DT / CTX_AREA_NK / CTX_AREA_FK), CTCA prefix TR_ID, **rate-limit 1/day 권장** |
| EP6 | `InquireViStatus` | `quotations/inquire-vi-status` | FHPST01390000 | `output {}` | 13 | Doc declares `output: {}` Object (단일), but "30건" 문구 있음 — struct는 single object로 구현, 주석에 runtime 배열 가능성 명기 |
| EP7 | `InquireCaptureUplowprice` | `quotations/capture-uplowprice` | FHKST130000C0 | `output []` | 17 | `FID_COND_SCR_DIV_CODE="11300"` hardcoded |

Default `FID_COND_MRKT_DIV_CODE` = `"J"`.

---

## 파일 구조

### 신규 (Go source)
- `domestic/market_op.go` — EP4~EP7 (4 메서드 + structs + Params) 신규 생성
- `domestic/market_op_test.go` — EP4~EP7 테스트
- `examples/domestic_market_op/main.go` — 사용 예제

### 수정 (APPEND / MODIFY)
- `CLAUDE.md` — banner v1.12.0 → v1.13.0, Phase 4.2 plan link 추가
- `README.md` — Available Methods 표 갱신 (89 → 93 메서드), Phase 4.2 섹션 추가
- `CHANGELOG.md` — `[1.13.0]` entry ABOVE `[1.12.0]`, WebSocket 제외 사유 명기
- `domestic/doc.go` — Phase 4.2 section 추가

### 신규 (testdata — 4 files)
- `domestic/testdata/exp_closing_price_success.json`
- `domestic/testdata/chk_holiday_success.json`
- `domestic/testdata/vi_status_success.json`
- `domestic/testdata/capture_uplowprice_success.json`

---

## 타입 매핑

Phase 2 표준 타입 매핑 — Phase 4.2 시장운영/특수상태 특화.

| 카테고리 | Go 타입 | json tag suffix | 예시 필드 |
|---|---|---|---|
| 가격 | `decimal.Decimal` | (bare) | `stck_prpr`, `prdy_vrss`, `sdpr_vrss_prpr`, `vi_prc`, `vi_stnd_prc`, `vi_dmc_stnd_prc`, `stck_llam`, `stck_mxpr` |
| 수량/금액 | `int64` | `,string` | `cntg_vol`, `acml_vol`, `prdy_vol`, `vi_count`, `total_askp_rsqn`, `total_bidp_rsqn`, `askp_rsqn1`, `bidp_rsqn1`, `seln_cnqn`, `shnu_cnqn` |
| 비율 | `float64` | `,string` | `prdy_ctrt`, `sdpr_vrss_prpr_rate`, `vi_dprt`, `vi_dmc_dprt`, `prdy_vrss_vol_rate` |
| 코드/이름/날짜/Y-N | `string` | (bare) | `stck_shrn_iscd`, `mksc_shrn_iscd`, `hts_kor_isnm`, `prdy_vrss_sign`, `vi_cls_code`, `vi_kind_code`, `bsop_date`, `cntg_vi_hour`, `vi_cncl_hour`, `bass_dt`, `wday_dvsn_cd`, `bzdy_yn`, `tr_day_yn`, `opnd_yn`, `sttl_day_yn` |

---

## Tasks (9 total)

| # | 내용 | Files |
|---|---|---|
| Task 1 | testdata fixtures (4 합성 JSON) | `domestic/testdata/*.json` |
| Task 2 | EP4 `InquireExpClosingPrice` — `market_op.go` 신규, output1 array | CREATE `market_op.go` + `market_op_test.go` |
| Task 3 | EP5 `InquireChkHoliday` (non-FID params, 1/day 권장) | APPEND `market_op.go` / `market_op_test.go` |
| Task 4 | EP6 `InquireViStatus` (single object, array 가능성 주석) | APPEND `market_op.go` / `market_op_test.go` |
| Task 5 | EP7 `InquireCaptureUplowprice` (output array 17 fields) | APPEND `market_op.go` / `market_op_test.go` |
| Task 6 | examples `examples/domestic_market_op/main.go` | CREATE |
| Task 7 | 문서 갱신 (`CLAUDE.md` / `README.md` / `CHANGELOG.md` / `domestic/doc.go`) | MODIFY |
| Task 8 | 최종 점검 (`gofmt` / `build` / `vet` / `race` / coverage ≥ 80%) | — |
| Task 9 | PR 생성 (사용자 승인 후) | push + `gh pr create` + merge + tag `v1.13.0` + GitHub Release |

---

## Task 1: testdata fixtures (4 합성 JSON)

- [ ] Step 1: `domestic/testdata/exp_closing_price_success.json` (EP4 — `output1 []`, 9 fields)
- [ ] Step 2: `domestic/testdata/chk_holiday_success.json` (EP5 — `output {}` 6 string fields)
- [ ] Step 3: `domestic/testdata/vi_status_success.json` (EP6 — `output {}` 13 fields)
- [ ] Step 4: `domestic/testdata/capture_uplowprice_success.json` (EP7 — `output []`, 17 fields)
- [ ] Step 5: validation

```bash
for f in \
  domestic/testdata/exp_closing_price_success.json \
  domestic/testdata/chk_holiday_success.json \
  domestic/testdata/vi_status_success.json \
  domestic/testdata/capture_uplowprice_success.json; do
  python3 -c "import json; json.load(open('$f'))" && echo "$f OK"
done
# Expected: 4 OK lines
```

- [ ] Step 6: commit

```bash
git commit -m "$(cat <<'EOF'
[chore] testdata — 4 market-op fixture JSON (Phase 4.2)

합성 JSON fixtures:
- exp_closing_price_success.json (output1[] 9 fields, FID_COND_SCR_DIV_CODE=11173)
- chk_holiday_success.json (output{} 6 string fields, non-FID UPPERCASE params)
- vi_status_success.json (output{} 13 fields, VI 발동/해제)
- capture_uplowprice_success.json (output[] 17 fields, 상하한가 포착)

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

### Fixture content

**Step 1 — `exp_closing_price_success.json`** (EP4: output1[] 9 fields, 장마감 예상체결가)

```json
{
  "rt_cd": "0",
  "msg_cd": "MCA00000",
  "msg1": "정상처리 되었습니다.",
  "output1": [
    {
      "stck_shrn_iscd": "005930",
      "hts_kor_isnm": "삼성전자",
      "stck_prpr": "82500",
      "prdy_vrss": "500",
      "prdy_vrss_sign": "2",
      "prdy_ctrt": "0.61",
      "sdpr_vrss_prpr": "1000",
      "sdpr_vrss_prpr_rate": "1.23",
      "cntg_vol": "125000"
    },
    {
      "stck_shrn_iscd": "000660",
      "hts_kor_isnm": "SK하이닉스",
      "stck_prpr": "185000",
      "prdy_vrss": "2000",
      "prdy_vrss_sign": "2",
      "prdy_ctrt": "1.09",
      "sdpr_vrss_prpr": "3000",
      "sdpr_vrss_prpr_rate": "1.65",
      "cntg_vol": "87500"
    }
  ]
}
```

**Step 2 — `chk_holiday_success.json`** (EP5: output{} 6 string fields, 휴장일 조회)

```json
{
  "rt_cd": "0",
  "msg_cd": "MCA00000",
  "msg1": "정상처리 되었습니다.",
  "output": {
    "bass_dt": "20260507",
    "wday_dvsn_cd": "04",
    "bzdy_yn": "Y",
    "tr_day_yn": "Y",
    "opnd_yn": "Y",
    "sttl_day_yn": "Y"
  }
}
```

**Step 3 — `vi_status_success.json`** (EP6: output{} 13 fields, 변동성완화장치 현황)

```json
{
  "rt_cd": "0",
  "msg_cd": "MCA00000",
  "msg1": "정상처리 되었습니다.",
  "output": {
    "hts_kor_isnm": "삼성전자",
    "mksc_shrn_iscd": "005930",
    "vi_cls_code": "Y",
    "bsop_date": "20260507",
    "cntg_vi_hour": "100530",
    "vi_cncl_hour": "100830",
    "vi_kind_code": "1",
    "vi_prc": "82500",
    "vi_stnd_prc": "81000",
    "vi_dprt": "1.85",
    "vi_dmc_stnd_prc": "80500",
    "vi_dmc_dprt": "2.48",
    "vi_count": "3"
  }
}
```

**Step 4 — `capture_uplowprice_success.json`** (EP7: output[] 17 fields/item, 상하한가 포착)

```json
{
  "rt_cd": "0",
  "msg_cd": "MCA00000",
  "msg1": "정상처리 되었습니다.",
  "output": [
    {
      "mksc_shrn_iscd": "005930",
      "hts_kor_isnm": "삼성전자",
      "stck_prpr": "82500",
      "prdy_vrss_sign": "2",
      "prdy_vrss": "500",
      "prdy_ctrt": "0.61",
      "acml_vol": "8750000",
      "total_askp_rsqn": "250000",
      "total_bidp_rsqn": "310000",
      "askp_rsqn1": "45000",
      "bidp_rsqn1": "62000",
      "prdy_vol": "9200000",
      "seln_cnqn": "4200000",
      "shnu_cnqn": "4550000",
      "stck_llam": "57750",
      "stck_mxpr": "107250",
      "prdy_vrss_vol_rate": "-4.89"
    },
    {
      "mksc_shrn_iscd": "000660",
      "hts_kor_isnm": "SK하이닉스",
      "stck_prpr": "185000",
      "prdy_vrss_sign": "2",
      "prdy_vrss": "2000",
      "prdy_ctrt": "1.09",
      "acml_vol": "3100000",
      "total_askp_rsqn": "95000",
      "total_bidp_rsqn": "120000",
      "askp_rsqn1": "18000",
      "bidp_rsqn1": "22000",
      "prdy_vol": "3400000",
      "seln_cnqn": "1500000",
      "shnu_cnqn": "1600000",
      "stck_llam": "129500",
      "stck_mxpr": "240500",
      "prdy_vrss_vol_rate": "-8.82"
    }
  ]
}
```

---

## Task 2: EP4 `InquireExpClosingPrice` — `market_op.go` 신규 생성

### Step 1: CREATE `domestic/market_op_test.go` with EP4 test

```go
package domestic_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/kenshin579/korea-investment-stock/domestic"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInquireExpClosingPrice(t *testing.T) {
	cl, transport := domestic.NewTestClient(t)
	transport.RegisterResponder(
		http.MethodGet,
		"https://openapi.koreainvestment.com:9443/uapi/domestic-stock/v1/quotations/exp-closing-price",
		httpmock.NewJsonResponderOrPanic(200, httpmock.File("testdata/exp_closing_price_success.json")),
	)

	res, err := cl.InquireExpClosingPrice(context.Background(), domestic.InquireExpClosingPriceParams{
		RankSortClsCode: "0",
		Symbol:          "0000",
		BlngClsCode:     "0",
	})
	require.NoError(t, err)
	require.Len(t, res.Output1, 2)

	item := res.Output1[0]
	assert.Equal(t, "005930", item.StckShrnIscd)
	assert.Equal(t, "삼성전자", item.HtsKorIsnm)
	assert.Equal(t, "82500", item.StckPrpr.String())
	assert.Equal(t, "500", item.PrdyVrss.String())
	assert.Equal(t, "2", item.PrdyVrssSign)
	assert.InDelta(t, 0.61, item.PrdyCtrt, 0.001)
	assert.Equal(t, "1000", item.SdprVrssPrpr.String())
	assert.InDelta(t, 1.23, item.SdprVrssPrprRate, 0.001)
	assert.Equal(t, int64(125000), item.CntgVol)
}
```

### Step 2: Verify FAIL

```bash
go test ./domestic/... -run TestInquireExpClosingPrice -v 2>&1 | tail -5
# Expected: FAIL (undefined: domestic.InquireExpClosingPrice)
```

### Step 3: CREATE `domestic/market_op.go`

```go
// Package domestic — market_op.go
// Phase 4.2: 시장운영/특수상태 4 메서드 (EP4~EP7)
//
// EP4  InquireExpClosingPrice    — 장마감 예상체결가   FHKST117300C0
// EP5  InquireChkHoliday         — 휴장일 조회         CTCA0903R
// EP6  InquireViStatus           — 변동성완화장치 현황 FHPST01390000
// EP7  InquireCaptureUplowprice  — 상하한가 포착       FHKST130000C0
//
// WebSocket 제외: 장운영정보 KRX(H0STMKO0) / NXT(H0NXMKO0) / 통합(H0UNMKO0) 은
// REST GET 이 아닌 WebSocket push API — Phase 4.2 범위 외. Phase 5 (WebSocket) 에서 처리 예정.
package domestic

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kenshin579/korea-investment-stock/internal/httpclient"
	"github.com/shopspring/decimal"
)

// ─── EP4: InquireExpClosingPrice ────────────────────────────────────────────

// InquireExpClosingPriceParams 는 장마감 예상체결가 조회 파라미터.
// FID_INPUT_ISCD 는 종목코드가 아닌 시장 구분코드: 0000(전체)/0001(코스피)/1001(코스닥)/2001(코스피200)/4001(KRX100).
type InquireExpClosingPriceParams struct {
	RankSortClsCode string // FID_RANK_SORT_CLS_CODE: 0=전체/1=상한가마감/2=하한가마감/3=상승률상위/4=하락률상위
	MarketCode      string // FID_COND_MRKT_DIV_CODE: 기본 "J"
	CondScrDivCode  string // FID_COND_SCR_DIV_CODE: 기본 "11173" (hardcoded)
	Symbol          string // FID_INPUT_ISCD: 시장구분코드 0000/0001/1001/2001/4001
	BlngClsCode     string // FID_BLNG_CLS_CODE: 0=전체/1=종가범위연장
}

// ExpClosingPriceItem 은 장마감 예상체결가 종목별 데이터.
type ExpClosingPriceItem struct {
	StckShrnIscd    string          `json:"stck_shrn_iscd"`
	HtsKorIsnm      string          `json:"hts_kor_isnm"`
	StckPrpr        decimal.Decimal `json:"stck_prpr"`
	PrdyVrss        decimal.Decimal `json:"prdy_vrss"`
	PrdyVrssSign    string          `json:"prdy_vrss_sign"`
	PrdyCtrt        float64         `json:"prdy_ctrt,string"`
	SdprVrssPrpr    decimal.Decimal `json:"sdpr_vrss_prpr"`
	SdprVrssPrprRate float64        `json:"sdpr_vrss_prpr_rate,string"`
	CntgVol         int64           `json:"cntg_vol,string"`
}

// InquireExpClosingPriceResponse 는 장마감 예상체결가 응답.
type InquireExpClosingPriceResponse struct {
	RtCd   string                `json:"rt_cd"`
	MsgCd  string                `json:"msg_cd"`
	Msg1   string                `json:"msg1"`
	Output1 []ExpClosingPriceItem `json:"output1"`
}

// InquireExpClosingPrice 는 장마감 예상체결가를 조회한다 (FHKST117300C0).
func (c *Client) InquireExpClosingPrice(ctx context.Context, params InquireExpClosingPriceParams) (*InquireExpClosingPriceResponse, error) {
	market := params.MarketCode
	if market == "" {
		market = "J"
	}
	scrDiv := params.CondScrDivCode
	if scrDiv == "" {
		scrDiv = "11173"
	}
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/domestic-stock/v1/quotations/exp-closing-price",
		TrID:   "FHKST117300C0",
		Query: map[string]string{
			"FID_COND_MRKT_DIV_CODE": market,
			"FID_COND_SCR_DIV_CODE":  scrDiv,
			"FID_RANK_SORT_CLS_CODE": params.RankSortClsCode,
			"FID_INPUT_ISCD":         params.Symbol,
			"FID_BLNG_CLS_CODE":      params.BlngClsCode,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}
	var res InquireExpClosingPriceResponse
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse InquireExpClosingPriceResponse: %w", err)
	}
	return &res, nil
}
```

### Step 4: Verify PASS

```bash
go test ./domestic/... -run TestInquireExpClosingPrice -v 2>&1 | tail -5
# Expected: PASS
```

### Step 5: gofmt / vet

```bash
gofmt -w domestic/market_op.go domestic/market_op_test.go
go vet ./domestic/...
```

### Step 6: commit

```bash
git commit -m "$(cat <<'EOF'
[feat] EP4 InquireExpClosingPrice — 장마감 예상체결가 (Phase 4.2)

domestic/market_op.go 신규 생성. FHKST117300C0.
output1 array (9 fields). FID_COND_SCR_DIV_CODE="11173" hardcoded.
FID_INPUT_ISCD = 시장구분코드 (종목코드 아님: 0000/0001/1001/2001/4001).
WebSocket 3 endpoints (H0STMKO0/H0NXMKO0/H0UNMKO0) 제외 사유 파일 상단 주석에 명기.

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 3: EP5 `InquireChkHoliday`

### Step 1: APPEND test to `domestic/market_op_test.go`

```go
func TestInquireChkHoliday(t *testing.T) {
	cl, transport := domestic.NewTestClient(t)
	transport.RegisterResponder(
		http.MethodGet,
		"https://openapi.koreainvestment.com:9443/uapi/domestic-stock/v1/quotations/chk-holiday",
		httpmock.NewJsonResponderOrPanic(200, httpmock.File("testdata/chk_holiday_success.json")),
	)

	res, err := cl.InquireChkHoliday(context.Background(), domestic.InquireChkHolidayParams{
		BassDt:    "20260507",
		CtxAreaNk: "",
		CtxAreaFk: "",
	})
	require.NoError(t, err)
	require.NotNil(t, res.Output)

	out := res.Output
	assert.Equal(t, "20260507", out.Bassdt)
	assert.Equal(t, "04", out.WdayDvsnCd)
	assert.Equal(t, "Y", out.BzdyYn)
	assert.Equal(t, "Y", out.TrDayYn)
	assert.Equal(t, "Y", out.OpndYn)
	assert.Equal(t, "Y", out.SttlDayYn)
}
```

### Step 2: Verify FAIL

```bash
go test ./domestic/... -run TestInquireChkHoliday -v 2>&1 | tail -5
# Expected: FAIL (undefined: domestic.InquireChkHoliday)
```

### Step 3: APPEND to `domestic/market_op.go`

```go
// ─── EP5: InquireChkHoliday ──────────────────────────────────────────────────

// InquireChkHolidayParams 는 휴장일 조회 파라미터.
// 주의: 파라미터명이 FID_ 접두어 없는 비표준 UPPERCASE 형식 (BASS_DT / CTX_AREA_NK / CTX_AREA_FK).
// 주의: 단시간 다수 호출 자제 (KIS docs 권장 1일 1회).
type InquireChkHolidayParams struct {
	BassDt    string // BASS_DT (Y): 조회기준일 YYYYMMDD
	CtxAreaNk string // CTX_AREA_NK (Y): 연속조회검색조건 (공란 가능)
	CtxAreaFk string // CTX_AREA_FK (Y): 연속조회키 (공란 가능)
}

// ChkHolidayOutput 는 휴장일 조회 단일 응답 객체.
type ChkHolidayOutput struct {
	BassDateTime string `json:"bass_dt"` // 기준일자 YYYYMMDD
	WdayDvsnCd   string `json:"wday_dvsn_cd"` // 요일구분코드 01(일)~07(토)
	BzdyYn       string `json:"bzdy_yn"`      // 영업일여부 Y/N
	TrDayYn      string `json:"tr_day_yn"`    // 거래일여부 Y/N
	OpndYn       string `json:"opnd_yn"`      // 개장일여부 Y/N
	SttlDayYn    string `json:"sttl_day_yn"`  // 결제일여부 Y/N
}

// ChkHolidayOutputAlias 는 json bass_dt 필드를 BassDateTime 이름 충돌 없이 매핑하기 위한 내부 alias.
// KIS wire key "bass_dt" → Go 필드 Bassdt.
type ChkHolidayItem struct {
	Bassdt     string `json:"bass_dt"`
	WdayDvsnCd string `json:"wday_dvsn_cd"`
	BzdyYn     string `json:"bzdy_yn"`
	TrDayYn    string `json:"tr_day_yn"`
	OpndYn     string `json:"opnd_yn"`
	SttlDayYn  string `json:"sttl_day_yn"`
}

// InquireChkHolidayResponse 는 휴장일 조회 응답.
type InquireChkHolidayResponse struct {
	RtCd   string          `json:"rt_cd"`
	MsgCd  string          `json:"msg_cd"`
	Msg1   string          `json:"msg1"`
	Output *ChkHolidayItem `json:"output"`
}

// InquireChkHoliday 는 휴장일을 조회한다 (CTCA0903R).
//
// 주의: 단시간 다수 호출 자제 (KIS docs 권장 1일 1회).
// 파라미터명이 FID_ 접두어 없는 비표준 UPPERCASE 형식임에 유의.
func (c *Client) InquireChkHoliday(ctx context.Context, params InquireChkHolidayParams) (*InquireChkHolidayResponse, error) {
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/domestic-stock/v1/quotations/chk-holiday",
		TrID:   "CTCA0903R",
		Query: map[string]string{
			"BASS_DT":    params.BassDt,
			"CTX_AREA_NK": params.CtxAreaNk,
			"CTX_AREA_FK": params.CtxAreaFk,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}
	var res InquireChkHolidayResponse
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse InquireChkHolidayResponse: %w", err)
	}
	return &res, nil
}
```

### Step 4: Verify PASS

```bash
go test ./domestic/... -run TestInquireChkHoliday -v 2>&1 | tail -5
# Expected: PASS
```

### Step 5: gofmt / vet

```bash
gofmt -w domestic/market_op.go domestic/market_op_test.go
go vet ./domestic/...
```

### Step 6: commit

```bash
git commit -m "$(cat <<'EOF'
[feat] EP5 InquireChkHoliday — 휴장일 조회 (Phase 4.2)

CTCA0903R. non-FID UPPERCASE params (BASS_DT/CTX_AREA_NK/CTX_AREA_FK).
output single object 6 string fields.
1일 1회 호출 권장 주석 명기.

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 4: EP6 `InquireViStatus`

### Step 1: APPEND test to `domestic/market_op_test.go`

```go
func TestInquireViStatus(t *testing.T) {
	cl, transport := domestic.NewTestClient(t)
	transport.RegisterResponder(
		http.MethodGet,
		"https://openapi.koreainvestment.com:9443/uapi/domestic-stock/v1/quotations/inquire-vi-status",
		httpmock.NewJsonResponderOrPanic(200, httpmock.File("testdata/vi_status_success.json")),
	)

	res, err := cl.InquireViStatus(context.Background(), domestic.InquireViStatusParams{
		DivClsCode:       "0",
		MrktClsCode:      "0",
		Symbol:           "",
		RankSortClsCode:  "0",
		InputDate1:       "20260507",
		TrgtClsCode:      "",
		TrgtExlsCode:     "",
	})
	require.NoError(t, err)
	require.NotNil(t, res.Output)

	out := res.Output
	assert.Equal(t, "삼성전자", out.HtsKorIsnm)
	assert.Equal(t, "005930", out.MkscShrnIscd)
	assert.Equal(t, "Y", out.ViClsCode)
	assert.Equal(t, "20260507", out.BsopDate)
	assert.Equal(t, "100530", out.CntgViHour)
	assert.Equal(t, "100830", out.ViCnclHour)
	assert.Equal(t, "1", out.ViKindCode)
	assert.Equal(t, "82500", out.ViPrc.String())
	assert.Equal(t, "81000", out.ViStndPrc.String())
	assert.InDelta(t, 1.85, out.ViDprt, 0.001)
	assert.Equal(t, "80500", out.ViDmcStndPrc.String())
	assert.InDelta(t, 2.48, out.ViDmcDprt, 0.001)
	assert.Equal(t, int64(3), out.ViCount)
}
```

### Step 2: Verify FAIL

```bash
go test ./domestic/... -run TestInquireViStatus -v 2>&1 | tail -5
# Expected: FAIL (undefined: domestic.InquireViStatus)
```

### Step 3: APPEND to `domestic/market_op.go`

```go
// ─── EP6: InquireViStatus ────────────────────────────────────────────────────

// InquireViStatusParams 는 변동성완화장치(VI) 현황 조회 파라미터.
type InquireViStatusParams struct {
	DivClsCode      string // FID_DIV_CLS_CODE (Y): 0=전체/1=상승/2=하락
	CondScrDivCode  string // FID_COND_SCR_DIV_CODE: 기본 "20139" (hardcoded)
	MrktClsCode     string // FID_MRKT_CLS_CODE (Y): 0=전체/K=거래소/Q=코스닥
	Symbol          string // FID_INPUT_ISCD (Y, 공란 가능)
	RankSortClsCode string // FID_RANK_SORT_CLS_CODE (Y): 0=전체/1=정적/2=동적/3=정적&동적
	InputDate1      string // FID_INPUT_DATE_1 (Y): YYYYMMDD
	TrgtClsCode     string // FID_TRGT_CLS_CODE (Y, 공란 가능)
	TrgtExlsCode    string // FID_TRGT_EXLS_CLS_CODE (Y, 공란 가능)
}

// ViStatusOutput 는 변동성완화장치(VI) 현황 응답 단일 객체.
//
// KIS 공식 문서는 output 을 단일 Object({})로 선언하나, 실제 응답에서 배열([])을 반환할 수 있음.
// ("30건" 등 복수 건 문구 포함). 실 API 호출 시 배열 반환 확인 시 []ViStatusOutput 로 전환 필요.
type ViStatusOutput struct {
	HtsKorIsnm   string          `json:"hts_kor_isnm"`
	MkscShrnIscd string          `json:"mksc_shrn_iscd"`
	ViClsCode    string          `json:"vi_cls_code"`    // Y=발동/N=해제
	BsopDate     string          `json:"bsop_date"`      // YYYYMMDD
	CntgViHour   string          `json:"cntg_vi_hour"`   // HHMMSS
	ViCnclHour   string          `json:"vi_cncl_hour"`   // HHMMSS
	ViKindCode   string          `json:"vi_kind_code"`   // 1=정적/2=동적/3=정적&동적
	ViPrc        decimal.Decimal `json:"vi_prc"`
	ViStndPrc    decimal.Decimal `json:"vi_stnd_prc"`
	ViDprt       float64         `json:"vi_dprt,string"`
	ViDmcStndPrc decimal.Decimal `json:"vi_dmc_stnd_prc"`
	ViDmcDprt    float64         `json:"vi_dmc_dprt,string"`
	ViCount      int64           `json:"vi_count,string"`
}

// InquireViStatusResponse 는 변동성완화장치(VI) 현황 응답.
type InquireViStatusResponse struct {
	RtCd   string          `json:"rt_cd"`
	MsgCd  string          `json:"msg_cd"`
	Msg1   string          `json:"msg1"`
	Output *ViStatusOutput `json:"output"`
}

// InquireViStatus 는 변동성완화장치(VI) 현황을 조회한다 (FHPST01390000).
//
// KIS 문서가 output 을 단일 Object로 선언. 실 API 에서 배열 반환 시 struct 변경 필요.
// FID_COND_SCR_DIV_CODE 는 "20139" 로 hardcoded.
func (c *Client) InquireViStatus(ctx context.Context, params InquireViStatusParams) (*InquireViStatusResponse, error) {
	scrDiv := params.CondScrDivCode
	if scrDiv == "" {
		scrDiv = "20139"
	}
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/domestic-stock/v1/quotations/inquire-vi-status",
		TrID:   "FHPST01390000",
		Query: map[string]string{
			"FID_DIV_CLS_CODE":        params.DivClsCode,
			"FID_COND_SCR_DIV_CODE":   scrDiv,
			"FID_MRKT_CLS_CODE":       params.MrktClsCode,
			"FID_INPUT_ISCD":          params.Symbol,
			"FID_RANK_SORT_CLS_CODE":  params.RankSortClsCode,
			"FID_INPUT_DATE_1":        params.InputDate1,
			"FID_TRGT_CLS_CODE":       params.TrgtClsCode,
			"FID_TRGT_EXLS_CLS_CODE":  params.TrgtExlsCode,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}
	var res InquireViStatusResponse
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse InquireViStatusResponse: %w", err)
	}
	return &res, nil
}
```

### Step 4: Verify PASS

```bash
go test ./domestic/... -run TestInquireViStatus -v 2>&1 | tail -5
# Expected: PASS
```

### Step 5: gofmt / vet

```bash
gofmt -w domestic/market_op.go domestic/market_op_test.go
go vet ./domestic/...
```

### Step 6: commit

```bash
git commit -m "$(cat <<'EOF'
[feat] EP6 InquireViStatus — 변동성완화장치 현황 (Phase 4.2)

FHPST01390000. FID_COND_SCR_DIV_CODE="20139" hardcoded.
output 단일 Object (13 fields). KIS 문서 "30건" 언급으로 runtime 배열 가능성
struct 주석에 명기. decimal/float64/int64 typed.

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 5: EP7 `InquireCaptureUplowprice`

### Step 1: APPEND test to `domestic/market_op_test.go`

```go
func TestInquireCaptureUplowprice(t *testing.T) {
	cl, transport := domestic.NewTestClient(t)
	transport.RegisterResponder(
		http.MethodGet,
		"https://openapi.koreainvestment.com:9443/uapi/domestic-stock/v1/quotations/capture-uplowprice",
		httpmock.NewJsonResponderOrPanic(200, httpmock.File("testdata/capture_uplowprice_success.json")),
	)

	res, err := cl.InquireCaptureUplowprice(context.Background(), domestic.InquireCaptureUplowpriceParams{
		PrcClsCode: "0",
		DivClsCode: "0",
		Symbol:     "0000",
	})
	require.NoError(t, err)
	require.Len(t, res.Output, 2)

	item := res.Output[0]
	assert.Equal(t, "005930", item.MkscShrnIscd)
	assert.Equal(t, "삼성전자", item.HtsKorIsnm)
	assert.Equal(t, "82500", item.StckPrpr.String())
	assert.Equal(t, "2", item.PrdyVrssSign)
	assert.Equal(t, "500", item.PrdyVrss.String())
	assert.InDelta(t, 0.61, item.PrdyCtrt, 0.001)
	assert.Equal(t, int64(8750000), item.AcmlVol)
	assert.Equal(t, int64(250000), item.TotalAskpRsqn)
	assert.Equal(t, int64(310000), item.TotalBidpRsqn)
	assert.Equal(t, int64(45000), item.AskpRsqn1)
	assert.Equal(t, int64(62000), item.BidpRsqn1)
	assert.Equal(t, int64(9200000), item.PrdyVol)
	assert.Equal(t, int64(4200000), item.SelnCnqn)
	assert.Equal(t, int64(4550000), item.ShnuCnqn)
	assert.Equal(t, "57750", item.StckLlam.String())
	assert.Equal(t, "107250", item.StckMxpr.String())
	assert.InDelta(t, -4.89, item.PrdyVrssVolRate, 0.001)
}
```

### Step 2: Verify FAIL

```bash
go test ./domestic/... -run TestInquireCaptureUplowprice -v 2>&1 | tail -5
# Expected: FAIL (undefined: domestic.InquireCaptureUplowprice)
```

### Step 3: APPEND to `domestic/market_op.go`

```go
// ─── EP7: InquireCaptureUplowprice ──────────────────────────────────────────

// InquireCaptureUplowpriceParams 는 상하한가 포착 조회 파라미터.
type InquireCaptureUplowpriceParams struct {
	MarketCode     string // FID_COND_MRKT_DIV_CODE: 기본 "J"
	CondScrDivCode string // FID_COND_SCR_DIV_CODE: 기본 "11300" (hardcoded)
	PrcClsCode     string // FID_PRC_CLS_CODE: 0=상한가/1=하한가
	DivClsCode     string // FID_DIV_CLS_CODE: 0=상하한가/6=8%근접/5=10%근접/1=15%근접/2=20%근접/3=25%근접
	Symbol         string // FID_INPUT_ISCD: 0000=전체/0001=코스피/1001=코스닥
	TrgtClsCode    string // FID_TRGT_CLS_CODE (공란 가능)
	TrgtExlsCode   string // FID_TRGT_EXLS_CLS_CODE (공란 가능)
	InputPrice1    string // FID_INPUT_PRICE_1 (공란 가능)
	InputPrice2    string // FID_INPUT_PRICE_2 (공란 가능)
	VolCnt         string // FID_VOL_CNT (공란 가능)
}

// CaptureUplowpriceItem 은 상하한가 포착 종목별 데이터.
type CaptureUplowpriceItem struct {
	MkscShrnIscd   string          `json:"mksc_shrn_iscd"`
	HtsKorIsnm     string          `json:"hts_kor_isnm"`
	StckPrpr       decimal.Decimal `json:"stck_prpr"`
	PrdyVrssSign   string          `json:"prdy_vrss_sign"`
	PrdyVrss       decimal.Decimal `json:"prdy_vrss"`
	PrdyCtrt       float64         `json:"prdy_ctrt,string"`
	AcmlVol        int64           `json:"acml_vol,string"`
	TotalAskpRsqn  int64           `json:"total_askp_rsqn,string"`
	TotalBidpRsqn  int64           `json:"total_bidp_rsqn,string"`
	AskpRsqn1      int64           `json:"askp_rsqn1,string"`
	BidpRsqn1      int64           `json:"bidp_rsqn1,string"`
	PrdyVol        int64           `json:"prdy_vol,string"`
	SelnCnqn       int64           `json:"seln_cnqn,string"`
	ShnuCnqn       int64           `json:"shnu_cnqn,string"`
	StckLlam       decimal.Decimal `json:"stck_llam"`
	StckMxpr       decimal.Decimal `json:"stck_mxpr"`
	PrdyVrssVolRate float64        `json:"prdy_vrss_vol_rate,string"`
}

// InquireCaptureUplowpriceResponse 는 상하한가 포착 응답.
type InquireCaptureUplowpriceResponse struct {
	RtCd   string                  `json:"rt_cd"`
	MsgCd  string                  `json:"msg_cd"`
	Msg1   string                  `json:"msg1"`
	Output []CaptureUplowpriceItem `json:"output"`
}

// InquireCaptureUplowprice 는 상하한가 포착 종목을 조회한다 (FHKST130000C0).
// FID_COND_SCR_DIV_CODE 는 "11300" 으로 hardcoded.
func (c *Client) InquireCaptureUplowprice(ctx context.Context, params InquireCaptureUplowpriceParams) (*InquireCaptureUplowpriceResponse, error) {
	market := params.MarketCode
	if market == "" {
		market = "J"
	}
	scrDiv := params.CondScrDivCode
	if scrDiv == "" {
		scrDiv = "11300"
	}
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/domestic-stock/v1/quotations/capture-uplowprice",
		TrID:   "FHKST130000C0",
		Query: map[string]string{
			"FID_COND_MRKT_DIV_CODE": market,
			"FID_COND_SCR_DIV_CODE":  scrDiv,
			"FID_PRC_CLS_CODE":       params.PrcClsCode,
			"FID_DIV_CLS_CODE":       params.DivClsCode,
			"FID_INPUT_ISCD":         params.Symbol,
			"FID_TRGT_CLS_CODE":      params.TrgtClsCode,
			"FID_TRGT_EXLS_CLS_CODE": params.TrgtExlsCode,
			"FID_INPUT_PRICE_1":      params.InputPrice1,
			"FID_INPUT_PRICE_2":      params.InputPrice2,
			"FID_VOL_CNT":            params.VolCnt,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}
	var res InquireCaptureUplowpriceResponse
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse InquireCaptureUplowpriceResponse: %w", err)
	}
	return &res, nil
}
```

### Step 4: Verify PASS

```bash
go test ./domestic/... -run TestInquireCaptureUplowprice -v 2>&1 | tail -5
# Expected: PASS
```

### Step 5: gofmt / vet

```bash
gofmt -w domestic/market_op.go domestic/market_op_test.go
go vet ./domestic/...
```

### Step 6: commit

```bash
git commit -m "$(cat <<'EOF'
[feat] EP7 InquireCaptureUplowprice — 상하한가 포착 (Phase 4.2)

FHKST130000C0. FID_COND_SCR_DIV_CODE="11300" hardcoded.
output array 17 fields/item. decimal/int64/float64 typed.

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 6: examples `examples/domestic_market_op/main.go`

### Step 1: CREATE `examples/domestic_market_op/main.go`

```go
// examples/domestic_market_op/main.go — Phase 4.2 시장운영/특수상태 사용 예제
//
// 환경 변수 필요:
//   KIS_APP_KEY, KIS_APP_SECRET, KIS_ACCOUNT_NO, KIS_MOCK (true/false)
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	kis "github.com/kenshin579/korea-investment-stock"
	"github.com/kenshin579/korea-investment-stock/domestic"
)

func main() {
	client, err := kis.NewClient(kis.Config{
		AppKey:    os.Getenv("KIS_APP_KEY"),
		AppSecret: os.Getenv("KIS_APP_SECRET"),
		AccountNo: os.Getenv("KIS_ACCOUNT_NO"),
		Mock:      os.Getenv("KIS_MOCK") == "true",
	})
	if err != nil {
		log.Fatalf("NewClient: %v", err)
	}
	ctx := context.Background()

	// EP4: 장마감 예상체결가 — 전체 종목 상위
	expClose, err := client.Domestic.InquireExpClosingPrice(ctx, domestic.InquireExpClosingPriceParams{
		RankSortClsCode: "0",  // 전체
		Symbol:          "0000", // 전체 시장
		BlngClsCode:     "0",  // 전체
	})
	if err != nil {
		log.Printf("InquireExpClosingPrice error: %v", err)
	} else {
		fmt.Printf("=== 장마감 예상체결가 (상위 %d건) ===\n", len(expClose.Output1))
		for i, item := range expClose.Output1 {
			if i >= 3 {
				fmt.Println("  ...")
				break
			}
			fmt.Printf("  [%d] %s (%s) 예상가=%s 등락률=%.2f%%\n",
				i+1, item.HtsKorIsnm, item.StckShrnIscd,
				item.StckPrpr.String(), item.PrdyCtrt)
		}
	}

	// EP5: 휴장일 조회 (1일 1회 권장)
	holiday, err := client.Domestic.InquireChkHoliday(ctx, domestic.InquireChkHolidayParams{
		BassDt:    "20260507",
		CtxAreaNk: "",
		CtxAreaFk: "",
	})
	if err != nil {
		log.Printf("InquireChkHoliday error: %v", err)
	} else if holiday.Output != nil {
		out := holiday.Output
		fmt.Printf("\n=== 휴장일 조회 (%s) ===\n", out.Bassdt)
		fmt.Printf("  영업일=%s 거래일=%s 개장일=%s 결제일=%s\n",
			out.BzdyYn, out.TrDayYn, out.OpndYn, out.SttlDayYn)
	}

	// EP6: 변동성완화장치(VI) 현황
	vi, err := client.Domestic.InquireViStatus(ctx, domestic.InquireViStatusParams{
		DivClsCode:      "0", // 전체
		MrktClsCode:     "0", // 전체
		Symbol:          "",
		RankSortClsCode: "0", // 전체
		InputDate1:      "20260507",
		TrgtClsCode:     "",
		TrgtExlsCode:    "",
	})
	if err != nil {
		log.Printf("InquireViStatus error: %v", err)
	} else if vi.Output != nil {
		out := vi.Output
		fmt.Printf("\n=== VI 현황 — %s (%s) ===\n", out.HtsKorIsnm, out.MkscShrnIscd)
		fmt.Printf("  발동여부=%s 종류=%s 발동가=%s 발동시=%s\n",
			out.ViClsCode, out.ViKindCode, out.ViPrc.String(), out.CntgViHour)
	}

	// EP7: 상하한가 포착
	uplow, err := client.Domestic.InquireCaptureUplowprice(ctx, domestic.InquireCaptureUplowpriceParams{
		PrcClsCode: "0", // 상한가
		DivClsCode: "0", // 상하한가
		Symbol:     "0000", // 전체
	})
	if err != nil {
		log.Printf("InquireCaptureUplowprice error: %v", err)
	} else {
		fmt.Printf("\n=== 상한가 포착 (%d건) ===\n", len(uplow.Output))
		for i, item := range uplow.Output {
			if i >= 3 {
				fmt.Println("  ...")
				break
			}
			fmt.Printf("  [%d] %s (%s) 현재가=%s 등락률=%.2f%%\n",
				i+1, item.HtsKorIsnm, item.MkscShrnIscd,
				item.StckPrpr.String(), item.PrdyCtrt)
		}
	}
}
```

### Step 2: Build check

```bash
go build ./examples/domestic_market_op/...
# Expected: no error
```

### Step 3: commit

```bash
git commit -m "$(cat <<'EOF'
[feat] examples/domestic_market_op — Phase 4.2 시장운영 사용 예제

4 메서드 (InquireExpClosingPrice/InquireChkHoliday/InquireViStatus/
InquireCaptureUplowprice) 통합 예제.

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 7: 문서 갱신

### Step 1: `CLAUDE.md` 갱신

`CLAUDE.md` 상단 banner 를 아래와 같이 수정:

```
> **Phase 4.2 — 국내주식 시장운영/특수상태 4 메서드 (v1.13.0). Phase 4 design spec 및 plan 참고.**
```

Phase 4.2 plan 링크 추가 (Phase 4.1 링크 아래):

```
- Phase 4.2 implementation plan: [`docs/superpowers/specs/2026-05-07-phase4-2-market-op-implementation-plan.md`](docs/superpowers/specs/2026-05-07-phase4-2-market-op-implementation-plan.md)
```

### Step 2: `README.md` 갱신

"Available Methods" 섹션에 Phase 4.2 소제목과 4개 행 추가:

```markdown
#### Phase 4.2 — 시장운영/특수상태 (v1.13.0)

| Method | TR_ID | Description |
|--------|-------|-------------|
| `Domestic.InquireExpClosingPrice` | FHKST117300C0 | 장마감 예상체결가 |
| `Domestic.InquireChkHoliday` | CTCA0903R | 휴장일 조회 (1일 1회 권장) |
| `Domestic.InquireViStatus` | FHPST01390000 | 변동성완화장치(VI) 현황 |
| `Domestic.InquireCaptureUplowprice` | FHKST130000C0 | 상하한가 포착 |
```

총 메서드 수 카운터: `89 → 93`.

### Step 3: `CHANGELOG.md` 갱신

`[1.12.0]` 섹션 **바로 위**에 `[1.13.0]` 추가:

```markdown
## [1.13.0] - 2026-05-07

### Added — Phase 4.2 (국내주식 시장운영/특수상태)

- `Domestic.InquireExpClosingPrice` — 장마감 예상체결가 (FHKST117300C0) — output1 array 9 fields; FID_INPUT_ISCD=시장구분코드
- `Domestic.InquireChkHoliday` — 휴장일 조회 (CTCA0903R) — output single object 6 string fields; 1일 1회 호출 권장
- `Domestic.InquireViStatus` — 변동성완화장치(VI) 현황 (FHPST01390000) — output single object 13 fields (runtime 배열 가능성 있음)
- `Domestic.InquireCaptureUplowprice` — 상하한가 포착 (FHKST130000C0) — output array 17 fields/item
- examples: `domestic_market_op`

### Notes

- Phase 4 design spec §Phase 4.2 는 7 메서드를 나열했으나, 3개 (장운영정보 KRX/NXT/통합 — TR_ID: H0STMKO0/H0NXMKO0/H0UNMKO0) 는 WebSocket push API 로 확인되어 Phase 4.2 범위에서 제외. 잠재적 Phase 5 (WebSocket) 로 이연.
- Phase 4 전체 메서드 수: 30 → 27 (WebSocket 3개 제외).
- EP5 (`InquireChkHoliday`) 파라미터명: FID_ 접두어 없는 비표준 UPPERCASE (BASS_DT/CTX_AREA_NK/CTX_AREA_FK).
- EP4 `FID_COND_SCR_DIV_CODE="11173"`, EP6 `FID_COND_SCR_DIV_CODE="20139"`, EP7 `FID_COND_SCR_DIV_CODE="11300"` hardcoded.
- EP6 (`InquireViStatus`) KIS 문서 output 단일 Object 선언 — 실 API 배열 반환 시 struct 변경 필요 (ViStatusOutput 주석 참조).
- 누적 89 → 93 메서드.
```

### Step 4: `domestic/doc.go` Phase 4.2 섹션 추가

파일 마지막 `package domestic` 선언 바로 앞의 `Anomalies (Phase 4.1)` 블록 뒤에 추가:

```go
// Phase 4.2 — 시장운영/특수상태 (v1.13.0)
//
//	EP4  InquireExpClosingPrice    — 장마감 예상체결가   FHKST117300C0
//	EP5  InquireChkHoliday         — 휴장일 조회         CTCA0903R
//	EP6  InquireViStatus           — 변동성완화장치 현황 FHPST01390000
//	EP7  InquireCaptureUplowprice  — 상하한가 포착       FHKST130000C0
//
// Anomalies (Phase 4.2):
//
//	EP4 output1 (not output) array, FID_INPUT_ISCD=시장구분코드 (종목코드 아님)
//	EP4 FID_COND_SCR_DIV_CODE="11173" hardcoded
//	EP5 non-FID UPPERCASE params (BASS_DT/CTX_AREA_NK/CTX_AREA_FK), CTCA prefix TR_ID
//	EP5 단시간 다수 호출 자제 (1일 1회 권장)
//	EP6 FID_COND_SCR_DIV_CODE="20139" hardcoded, output {} 단일 Object (runtime 배열 가능)
//	EP7 FID_COND_SCR_DIV_CODE="11300" hardcoded
//	WebSocket 제외: H0STMKO0/H0NXMKO0/H0UNMKO0 (장운영정보 KRX/NXT/통합) → Phase 5 이연
```

### Step 5: commit

```bash
git commit -m "$(cat <<'EOF'
[docs] Phase 4.2 문서 갱신 — v1.13.0 (market_op 4 메서드)

- CLAUDE.md: banner v1.12.0→v1.13.0, Phase 4.2 plan 링크 추가
- README.md: Phase 4.2 섹션 추가, 89→93 메서드
- CHANGELOG.md: [1.13.0] 추가 (WebSocket 3 endpoints 제외 사유 명기)
- domestic/doc.go: Phase 4.2 섹션 + Anomalies 추가

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 8: 최종 점검

### Step 1: 전체 빌드 / vet

```bash
go build ./...
go vet ./...
```

### Step 2: gofmt

```bash
gofmt -l ./domestic/market_op.go ./domestic/market_op_test.go
# Expected: 출력 없음 (이미 정형화됨)
```

### Step 3: race detector + 전체 테스트

```bash
go test -race ./...
# Expected: ok (no FAIL, no DATA RACE)
```

### Step 4: coverage 확인

```bash
go test -coverprofile=coverage.out ./domestic/...
go tool cover -func=coverage.out | grep -E "(market_op|total)"
# domestic coverage ≥ 80%

go test -coverprofile=coverage_root.out ./...
go tool cover -func=coverage_root.out | grep "total:"
# root coverage ≥ 80%
```

### Step 5: 최종 확인 체크리스트

```bash
# 4 메서드 exports 확인
grep -n "^func (c \*Client) Inquire" domestic/market_op.go
# Expected: 4 lines (InquireExpClosingPrice, InquireChkHoliday, InquireViStatus, InquireCaptureUplowprice)

# testdata 4 파일 존재 확인
ls domestic/testdata/ | grep -E "(exp_closing|chk_holiday|vi_status|capture_uplowprice)"
# Expected: 4 files

# examples build 확인
go build ./examples/domestic_market_op/...

# placeholder 없음 확인
grep -rn "TODO\|TBD\|FIXME\|placeholder" domestic/market_op.go domestic/market_op_test.go || echo "clean"
```

### Step 6: commit (if any formatting fix)

```bash
git commit -m "$(cat <<'EOF'
[chore] Phase 4.2 최종 점검 — gofmt/vet/race/coverage

모든 점검 통과: build/vet clean, race detector clean, domestic ≥80%.

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 9: PR 생성 (사용자 승인 후)

> **주의:** 이 Task 는 사용자의 명시적 승인 이후에만 실행한다.

### Step 1: push

```bash
git push -u origin feat/phase4-2-market-op
```

### Step 2: PR 생성 (HEREDOC)

```bash
gh pr create \
  --title "feat: Phase 4.2 — 국내주식 시장운영/특수상태 4 메서드 (v1.13.0)" \
  --base main \
  --reviewer kenshin579 \
  --body "$(cat <<'EOF'
## Summary

Phase 4.2: 국내주식 시장운영/특수상태 4 REST 메서드 추가 (`domestic/market_op.go` 신규).

### 추가 메서드 (89 → 93 누적)

| Method | TR_ID | Description |
|--------|-------|-------------|
| `Domestic.InquireExpClosingPrice` | FHKST117300C0 | 장마감 예상체결가 (output1 array) |
| `Domestic.InquireChkHoliday` | CTCA0903R | 휴장일 조회 (1일 1회 권장) |
| `Domestic.InquireViStatus` | FHPST01390000 | 변동성완화장치(VI) 현황 |
| `Domestic.InquireCaptureUplowprice` | FHKST130000C0 | 상하한가 포착 |

### WebSocket scope deviation

Phase 4 design spec §4.2 는 7 메서드를 나열했으나, 장운영정보 KRX/NXT/통합
(H0STMKO0 / H0NXMKO0 / H0UNMKO0) 3개는 WebSocket push API 로 확인.
이 라이브러리는 read-only REST 전용이므로 Phase 4.2 실제 범위는 **4 메서드**.
해당 3개 endpoint 는 잠재적 Phase 5 (WebSocket) 로 이연.
Phase 4 전체: 30 → 27 메서드.

### 주요 anomalies

- EP4: `output1` (not `output`), `FID_INPUT_ISCD` = 시장구분코드 (종목코드 아님), scr_div=11173 hardcoded
- EP5: non-FID UPPERCASE params (BASS_DT/CTX_AREA_NK/CTX_AREA_FK), CTCA prefix TR_ID, 1일 1회 권장
- EP6: scr_div=20139 hardcoded, output 단일 Object (runtime 배열 가능성 struct 주석 명기)
- EP7: scr_div=11300 hardcoded

## Test plan

- [ ] `go test -race ./domestic/... -run "TestInquireExpClosingPrice|TestInquireChkHoliday|TestInquireViStatus|TestInquireCaptureUplowprice"` PASS
- [ ] `go test -race ./...` clean (no FAIL, no DATA RACE)
- [ ] domestic coverage ≥ 80%
- [ ] `go build ./examples/domestic_market_op/...` success
- [ ] CHANGELOG [1.13.0] WebSocket 제외 사유 명기 확인
- [ ] `domestic/doc.go` Phase 4.2 섹션 + Anomalies 존재 확인

🤖 Generated with [Claude Code](https://claude.com/claude-code)
EOF
)"
```

### Step 3: merge 후 tag + GitHub Release

```bash
# PR merge 후
git checkout main && git pull origin main

git tag v1.13.0
git push origin v1.13.0

gh release create v1.13.0 \
  --title "v1.13.0 — Phase 4.2 시장운영/특수상태 4 메서드" \
  --notes "$(cat <<'EOF'
## Phase 4.2 — 국내주식 시장운영/특수상태 (v1.13.0)

4 REST 메서드 추가. 누적 93 메서드.

### 추가

- `Domestic.InquireExpClosingPrice` (FHKST117300C0) — 장마감 예상체결가
- `Domestic.InquireChkHoliday` (CTCA0903R) — 휴장일 조회 (1일 1회 권장)
- `Domestic.InquireViStatus` (FHPST01390000) — 변동성완화장치(VI) 현황
- `Domestic.InquireCaptureUplowprice` (FHKST130000C0) — 상하한가 포착

### WebSocket 범위 조정

장운영정보 KRX/NXT/통합 (H0STMKO0/H0NXMKO0/H0UNMKO0) 은 WebSocket push API 로
REST 라이브러리 범위 외. Phase 5 (WebSocket) 에서 처리 예정.
Phase 4 전체: 30 → 27 메서드.
EOF
)"
```

### Step 4: memory 갱신

사용자에게 MEMORY.md 갱신 요청:

```
[korea-investment-stock Go 마이그레이션 진행](kis_go_migration.md) — Phase 4.2 완료 ✅ (v1.13.0, 누적 93 메서드, 4 REST 메서드: 장마감예상체결가/휴장일조회/VI현황/상하한가포착). Phase 4 실제 범위: 27 메서드 (WebSocket 3 제외). 다음: Phase 4.3 (추가 ranking/시장흐름 13 메서드, v1.14.0).
```

---

## Self-review checklist

- [ ] 0 placeholders (TBD/TODO/Similar) — 없음
- [ ] 4 메서드 모두 커버 (7 아님 — WebSocket 3 제외 문서화됨)
- [ ] EP5 non-FID UPPERCASE params (BASS_DT/CTX_AREA_NK/CTX_AREA_FK) 검증
- [ ] EP4 scr_div="11173" / EP6 scr_div="20139" / EP7 scr_div="11300" hardcoded 확인
- [ ] HEREDOC commit messages — 모든 6개 commit
- [ ] Task 9 "사용자 승인 후" 명시
- [ ] CHANGELOG `[1.13.0]` — WebSocket scope deviation 명기
- [ ] EP6 ViStatusOutput 구조체 주석에 runtime 배열 가능성 명기
- [ ] EP5 1일 1회 권장 주석 (`InquireChkHoliday` 함수 및 Params 주석)
