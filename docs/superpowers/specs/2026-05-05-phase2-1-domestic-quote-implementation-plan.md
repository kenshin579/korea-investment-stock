# Phase 2.1 — 국내 호가/체결 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** `korea-investment-stock` Go 라이브러리에 국내주식 호가/체결 디테일 3 메서드 추가 (`v1.4.0` release).

**Architecture:** Phase 1 의 인프라 + 패턴 재사용. `domestic/quote.go` (3 메서드 한 file) 추가. 한투 API path 1:1 매핑 (Style A). 새 internal package 불필요. TDD: testdata fixture (KIS docs 응답 필드 기반 합성 JSON) → 실패 테스트 → struct + 메서드 구현 → 통과 → commit.

**Tech Stack:** Go 1.25+, `github.com/jarcoal/httpmock`, `github.com/shopspring/decimal`, `github.com/stretchr/testify`. 새 dependency 없음.

**참고 spec:**
- Phase 2 design spec: `docs/superpowers/specs/2026-05-05-phase2-readonly-extension-design.md` (commit `527ac96`)
- Phase 1.4 plan (참조 패턴): `docs/superpowers/specs/2026-05-05-phase1-4-investor-industry-ipo-implementation-plan.md`
- 한투 API docs: `docs/api/국내주식/{주식현재가_호가_예상체결.md, 주식현재가_체결.md, 주식현재가_일자별.md}`

---

## 사전 정보

| 항목 | 값 |
|------|---|
| 작업 브랜치 | `docs/phase2-spec` (Phase 2 spec 작성 시 생성됨) |
| 시작 HEAD | `527ac96` (Phase 2 design spec commit) |
| Release 목표 | `v1.4.0` |
| PR 베이스 | `main` |
| 현재 main 상태 | v1.3.0 publish 완료 (Phase 1.5 통합, 28 메서드) |

---

## 메서드 → 한투 API 매핑

| Go 메서드 | 한투 path | TR_ID |
|---|---|---|
| `Domestic.InquireAskingPriceExpCcn(ctx, params)` | `/uapi/domestic-stock/v1/quotations/inquire-asking-price-exp-ccn` | FHKST01010200 |
| `Domestic.InquireCcnl(ctx, params)` | `/uapi/domestic-stock/v1/quotations/inquire-ccnl` | FHKST01010300 |
| `Domestic.InquireDailyPrice(ctx, params)` | `/uapi/domestic-stock/v1/quotations/inquire-daily-price` | FHKST01010400 |

---

## 파일 구조

### 신규 (domestic)
- `domestic/quote.go` — 3 메서드 + structs + Params
- `domestic/quote_test.go`
- `domestic/testdata/asking_price_exp_ccn_success.json`
- `domestic/testdata/ccnl_success.json`
- `domestic/testdata/daily_price_success.json`

### 신규 (examples)
- `examples/domestic_quote/main.go` — 3 메서드 통합 예시

### 수정 (root)
- `CLAUDE.md` — banner 갱신 (v1.3.0 → v1.4.0, Phase 1 → Phase 2.1)
- `README.md` — Available Methods 표 갱신 (28 → 31 메서드)
- `CHANGELOG.md` — `[1.4.0]` entry
- `domestic/doc.go` — Phase 2.1 메서드 안내 추가

---

## 타입 매핑 (Phase 1 동일)

- **가격류 → `decimal.Decimal` (bare tag)**: `stck_prpr`, `stck_oprc`, `stck_hgpr`, `stck_lwpr`, `stck_sdpr`, `stck_clpr`, `prdy_vrss`, `askp1~10`, `bidp1~10`, `antc_cnpr`, `antc_cntg_vrss`
- **수량/금액 → `int64,string`**: `cntg_vol`, `acml_vol`, `antc_vol`, `frgn_ntby_qty`, `askp_rsqn1~10`, `bidp_rsqn1~10`, `askp_rsqn_icdc1~10`, `bidp_rsqn_icdc1~10`, `total_askp_rsqn`, `total_bidp_rsqn`, `total_askp_rsqn_icdc`, `total_bidp_rsqn_icdc`, `ovtm_total_askp_icdc`, `ovtm_total_bidp_icdc`, `ovtm_total_askp_rsqn`, `ovtm_total_bidp_rsqn`, `ntby_aspr_rsqn`
- **비율 → `float64,string`**: `prdy_ctrt`, `tday_rltv`, `antc_cntg_prdy_ctrt`, `prdy_vrss_vol_rate`, `hts_frgn_ehrt`, `acml_prtt_rate`
- **코드/시간/Y-N → 평문 `string`**: `aspr_acpt_hour`, `stck_cntg_hour`, `prdy_vrss_sign`, `antc_cntg_vrss_sign`, `new_mkop_cls_code`, `antc_mkop_cls_code`, `vi_cls_code`, `stck_shrn_iscd`, `stck_bsop_date`, `flng_cls_code`

---

## Task 1: testdata fixtures (3 합성 JSON)

**Files (Create):**
- `domestic/testdata/asking_price_exp_ccn_success.json`
- `domestic/testdata/ccnl_success.json`
- `domestic/testdata/daily_price_success.json`

> KIS docs 응답 필드 기반 합성. AskingPriceExpCcn 의 호가 1~10 + 잔량 1~10 + 증감 1~10 (60 fields) 모두 채움 (testdata 가 cover 하지 않으면 Go json zero-value default 적용).

- [ ] **Step 1: asking_price_exp_ccn_success.json**

```json
{
  "rt_cd": "0",
  "msg_cd": "MCA00000",
  "msg1": "정상처리 되었습니다.",
  "output1": {
    "aspr_acpt_hour": "131542",
    "askp1": "75900", "askp2": "75950", "askp3": "76000", "askp4": "76050", "askp5": "76100",
    "askp6": "76150", "askp7": "76200", "askp8": "76250", "askp9": "76300", "askp10": "76350",
    "bidp1": "75800", "bidp2": "75750", "bidp3": "75700", "bidp4": "75650", "bidp5": "75600",
    "bidp6": "75550", "bidp7": "75500", "bidp8": "75450", "bidp9": "75400", "bidp10": "75350",
    "askp_rsqn1": "1500", "askp_rsqn2": "2300", "askp_rsqn3": "1800", "askp_rsqn4": "1200", "askp_rsqn5": "900",
    "askp_rsqn6": "1100", "askp_rsqn7": "800", "askp_rsqn8": "600", "askp_rsqn9": "500", "askp_rsqn10": "400",
    "bidp_rsqn1": "2500", "bidp_rsqn2": "1900", "bidp_rsqn3": "2100", "bidp_rsqn4": "1700", "bidp_rsqn5": "1400",
    "bidp_rsqn6": "1100", "bidp_rsqn7": "900", "bidp_rsqn8": "700", "bidp_rsqn9": "600", "bidp_rsqn10": "500",
    "askp_rsqn_icdc1": "100", "askp_rsqn_icdc2": "-50", "askp_rsqn_icdc3": "200", "askp_rsqn_icdc4": "-30", "askp_rsqn_icdc5": "10",
    "askp_rsqn_icdc6": "0", "askp_rsqn_icdc7": "-20", "askp_rsqn_icdc8": "5", "askp_rsqn_icdc9": "0", "askp_rsqn_icdc10": "0",
    "bidp_rsqn_icdc1": "200", "bidp_rsqn_icdc2": "-100", "bidp_rsqn_icdc3": "150", "bidp_rsqn_icdc4": "-80", "bidp_rsqn_icdc5": "20",
    "bidp_rsqn_icdc6": "0", "bidp_rsqn_icdc7": "-10", "bidp_rsqn_icdc8": "0", "bidp_rsqn_icdc9": "0", "bidp_rsqn_icdc10": "0",
    "total_askp_rsqn": "11100",
    "total_bidp_rsqn": "13400",
    "total_askp_rsqn_icdc": "215",
    "total_bidp_rsqn_icdc": "180",
    "ovtm_total_askp_icdc": "0",
    "ovtm_total_bidp_icdc": "0",
    "ovtm_total_askp_rsqn": "0",
    "ovtm_total_bidp_rsqn": "0",
    "ntby_aspr_rsqn": "2300",
    "new_mkop_cls_code": "20"
  },
  "output2": {
    "antc_mkop_cls_code": "030",
    "stck_prpr": "75800",
    "stck_oprc": "76000",
    "stck_hgpr": "76200",
    "stck_lwpr": "75500",
    "stck_sdpr": "76000",
    "antc_cnpr": "75800",
    "antc_cntg_vrss_sign": "5",
    "antc_cntg_vrss": "-200",
    "antc_cntg_prdy_ctrt": "-0.26",
    "antc_vol": "12345678",
    "stck_shrn_iscd": "005930",
    "vi_cls_code": "N"
  }
}
```

- [ ] **Step 2: ccnl_success.json**

```json
{
  "rt_cd": "0",
  "msg_cd": "MCA00000",
  "msg1": "정상처리 되었습니다.",
  "output": [
    {
      "stck_cntg_hour": "131542",
      "stck_prpr": "75800",
      "prdy_vrss": "-200",
      "prdy_vrss_sign": "5",
      "cntg_vol": "12345",
      "tday_rltv": "104.32",
      "prdy_ctrt": "-0.26"
    },
    {
      "stck_cntg_hour": "131541",
      "stck_prpr": "75900",
      "prdy_vrss": "-100",
      "prdy_vrss_sign": "5",
      "cntg_vol": "8000",
      "tday_rltv": "103.50",
      "prdy_ctrt": "-0.13"
    }
  ]
}
```

- [ ] **Step 3: daily_price_success.json**

```json
{
  "rt_cd": "0",
  "msg_cd": "MCA00000",
  "msg1": "정상처리 되었습니다.",
  "output": [
    {
      "stck_bsop_date": "20260505",
      "stck_oprc": "76000",
      "stck_hgpr": "76200",
      "stck_lwpr": "75500",
      "stck_clpr": "75800",
      "acml_vol": "12345678",
      "prdy_vrss_vol_rate": "112.23",
      "prdy_vrss": "-200",
      "prdy_vrss_sign": "5",
      "prdy_ctrt": "-0.26",
      "hts_frgn_ehrt": "53.42",
      "frgn_ntby_qty": "-123456",
      "flng_cls_code": "00",
      "acml_prtt_rate": "0.00"
    },
    {
      "stck_bsop_date": "20260502",
      "stck_oprc": "76500",
      "stck_hgpr": "76800",
      "stck_lwpr": "75900",
      "stck_clpr": "76000",
      "acml_vol": "11000000",
      "prdy_vrss_vol_rate": "98.50",
      "prdy_vrss": "100",
      "prdy_vrss_sign": "2",
      "prdy_ctrt": "0.13",
      "hts_frgn_ehrt": "53.40",
      "frgn_ntby_qty": "50000",
      "flng_cls_code": "00",
      "acml_prtt_rate": "0.00"
    }
  ]
}
```

- [ ] **Step 4: 검증**

```bash
for f in domestic/testdata/{asking_price_exp_ccn,ccnl,daily_price}_success.json; do
  python3 -c "import json; json.load(open('$f'))" && echo "$f OK" || echo "$f BROKEN"
done
```
Expected: 3 lines `OK`.

- [ ] **Step 5: Commit**

```bash
git add domestic/testdata/{asking_price_exp_ccn,ccnl,daily_price}_success.json
git commit -m "$(cat <<'EOF'
[chore] Phase 2.1 testdata — 3 합성 JSON fixtures

호가/예상체결 (asking_price_exp_ccn, FHKST01010200) + 체결 (ccnl, FHKST01010300) +
일자별 (daily_price, FHKST01010400). KIS docs 의 응답 필드 정의 기반 합성.
삼성전자 (005930) 샘플.

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 2: domestic/quote.go — InquireAskingPriceExpCcn (file base)

**Files:**
- Create: `domestic/quote.go`
- Create: `domestic/quote_test.go`

> 가장 큰 struct (output1 ~58 fields, output2 13 fields). 호가 1~10 + 잔량 1~10 + 증감 1~10 indexed 필드. file 의 base 도 같이 setup.

- [ ] **Step 1: 테스트 작성** — `domestic/quote_test.go`

```go
package domestic_test

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kenshin579/korea-investment-stock/domestic"
)

func TestClient_InquireAskingPriceExpCcn(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/quotations/inquire-asking-price-exp-ccn`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "asking_price_exp_ccn_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireAskingPriceExpCcn(context.Background(), domestic.InquireAskingPriceExpCcnParams{
		Symbol: "005930",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, "J", capturedQuery.Get("FID_COND_MRKT_DIV_CODE"))
	assert.Equal(t, "005930", capturedQuery.Get("FID_INPUT_ISCD"))

	// output1 검증
	assert.Equal(t, "131542", res.Output1.AsprAcptHour)
	d, _ := decimal.NewFromString("75900")
	assert.True(t, d.Equal(res.Output1.Askp1))
	d, _ = decimal.NewFromString("75800")
	assert.True(t, d.Equal(res.Output1.Bidp1))
	assert.Equal(t, int64(1500), res.Output1.AskpRsqn1)
	assert.Equal(t, int64(2500), res.Output1.BidpRsqn1)
	assert.Equal(t, int64(11100), res.Output1.TotalAskpRsqn)
	assert.Equal(t, int64(13400), res.Output1.TotalBidpRsqn)

	// output2 검증
	assert.Equal(t, "030", res.Output2.AntcMkopClsCode)
	assert.True(t, d.Equal(res.Output2.StckPrpr))
	assert.Equal(t, int64(12345678), res.Output2.AntcVol)
	assert.Equal(t, "005930", res.Output2.StckShrnIscd)
}
```

- [ ] **Step 2: FAIL**

`go test ./domestic/... -run InquireAskingPriceExpCcn -v` — 컴파일 실패.

- [ ] **Step 3: 구현 — `domestic/quote.go`**

```go
package domestic

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/shopspring/decimal"

	"github.com/kenshin579/korea-investment-stock/internal/httpclient"
)

// AskingPriceExpCcn 은 주식현재가 호가/예상체결 (FHKST01010200) 응답.
//
// 한투 docs: docs/api/국내주식/주식현재가_호가_예상체결.md
// path: /uapi/domestic-stock/v1/quotations/inquire-asking-price-exp-ccn
type AskingPriceExpCcn struct {
	Output1 AskingPriceExpCcnOrderbook `json:"output1"`
	Output2 AskingPriceExpCcnExpected  `json:"output2"`
}

// AskingPriceExpCcnOrderbook 은 응답의 output1 — 10 단계 호가, 잔량, 증감.
type AskingPriceExpCcnOrderbook struct {
	AsprAcptHour string `json:"aspr_acpt_hour"` // 호가 접수 시간 (HHMMSS)

	Askp1  decimal.Decimal `json:"askp1"`  // 매도호가 1
	Askp2  decimal.Decimal `json:"askp2"`
	Askp3  decimal.Decimal `json:"askp3"`
	Askp4  decimal.Decimal `json:"askp4"`
	Askp5  decimal.Decimal `json:"askp5"`
	Askp6  decimal.Decimal `json:"askp6"`
	Askp7  decimal.Decimal `json:"askp7"`
	Askp8  decimal.Decimal `json:"askp8"`
	Askp9  decimal.Decimal `json:"askp9"`
	Askp10 decimal.Decimal `json:"askp10"`

	Bidp1  decimal.Decimal `json:"bidp1"`  // 매수호가 1
	Bidp2  decimal.Decimal `json:"bidp2"`
	Bidp3  decimal.Decimal `json:"bidp3"`
	Bidp4  decimal.Decimal `json:"bidp4"`
	Bidp5  decimal.Decimal `json:"bidp5"`
	Bidp6  decimal.Decimal `json:"bidp6"`
	Bidp7  decimal.Decimal `json:"bidp7"`
	Bidp8  decimal.Decimal `json:"bidp8"`
	Bidp9  decimal.Decimal `json:"bidp9"`
	Bidp10 decimal.Decimal `json:"bidp10"`

	AskpRsqn1  int64 `json:"askp_rsqn1,string"`  // 매도호가 잔량 1
	AskpRsqn2  int64 `json:"askp_rsqn2,string"`
	AskpRsqn3  int64 `json:"askp_rsqn3,string"`
	AskpRsqn4  int64 `json:"askp_rsqn4,string"`
	AskpRsqn5  int64 `json:"askp_rsqn5,string"`
	AskpRsqn6  int64 `json:"askp_rsqn6,string"`
	AskpRsqn7  int64 `json:"askp_rsqn7,string"`
	AskpRsqn8  int64 `json:"askp_rsqn8,string"`
	AskpRsqn9  int64 `json:"askp_rsqn9,string"`
	AskpRsqn10 int64 `json:"askp_rsqn10,string"`

	BidpRsqn1  int64 `json:"bidp_rsqn1,string"`  // 매수호가 잔량 1
	BidpRsqn2  int64 `json:"bidp_rsqn2,string"`
	BidpRsqn3  int64 `json:"bidp_rsqn3,string"`
	BidpRsqn4  int64 `json:"bidp_rsqn4,string"`
	BidpRsqn5  int64 `json:"bidp_rsqn5,string"`
	BidpRsqn6  int64 `json:"bidp_rsqn6,string"`
	BidpRsqn7  int64 `json:"bidp_rsqn7,string"`
	BidpRsqn8  int64 `json:"bidp_rsqn8,string"`
	BidpRsqn9  int64 `json:"bidp_rsqn9,string"`
	BidpRsqn10 int64 `json:"bidp_rsqn10,string"`

	AskpRsqnIcdc1  int64 `json:"askp_rsqn_icdc1,string"`  // 매도호가 잔량 증감 1
	AskpRsqnIcdc2  int64 `json:"askp_rsqn_icdc2,string"`
	AskpRsqnIcdc3  int64 `json:"askp_rsqn_icdc3,string"`
	AskpRsqnIcdc4  int64 `json:"askp_rsqn_icdc4,string"`
	AskpRsqnIcdc5  int64 `json:"askp_rsqn_icdc5,string"`
	AskpRsqnIcdc6  int64 `json:"askp_rsqn_icdc6,string"`
	AskpRsqnIcdc7  int64 `json:"askp_rsqn_icdc7,string"`
	AskpRsqnIcdc8  int64 `json:"askp_rsqn_icdc8,string"`
	AskpRsqnIcdc9  int64 `json:"askp_rsqn_icdc9,string"`
	AskpRsqnIcdc10 int64 `json:"askp_rsqn_icdc10,string"`

	BidpRsqnIcdc1  int64 `json:"bidp_rsqn_icdc1,string"`
	BidpRsqnIcdc2  int64 `json:"bidp_rsqn_icdc2,string"`
	BidpRsqnIcdc3  int64 `json:"bidp_rsqn_icdc3,string"`
	BidpRsqnIcdc4  int64 `json:"bidp_rsqn_icdc4,string"`
	BidpRsqnIcdc5  int64 `json:"bidp_rsqn_icdc5,string"`
	BidpRsqnIcdc6  int64 `json:"bidp_rsqn_icdc6,string"`
	BidpRsqnIcdc7  int64 `json:"bidp_rsqn_icdc7,string"`
	BidpRsqnIcdc8  int64 `json:"bidp_rsqn_icdc8,string"`
	BidpRsqnIcdc9  int64 `json:"bidp_rsqn_icdc9,string"`
	BidpRsqnIcdc10 int64 `json:"bidp_rsqn_icdc10,string"`

	TotalAskpRsqn      int64  `json:"total_askp_rsqn,string"`      // 총 매도호가 잔량
	TotalBidpRsqn      int64  `json:"total_bidp_rsqn,string"`      // 총 매수호가 잔량
	TotalAskpRsqnIcdc  int64  `json:"total_askp_rsqn_icdc,string"` // 총 매도호가 잔량 증감
	TotalBidpRsqnIcdc  int64  `json:"total_bidp_rsqn_icdc,string"` // 총 매수호가 잔량 증감
	OvtmTotalAskpIcdc  int64  `json:"ovtm_total_askp_icdc,string"` // 시간외 총 매도호가 증감
	OvtmTotalBidpIcdc  int64  `json:"ovtm_total_bidp_icdc,string"` // 시간외 총 매수호가 증감
	OvtmTotalAskpRsqn  int64  `json:"ovtm_total_askp_rsqn,string"` // 시간외 총 매도호가 잔량
	OvtmTotalBidpRsqn  int64  `json:"ovtm_total_bidp_rsqn,string"` // 시간외 총 매수호가 잔량
	NtbyAsprRsqn       int64  `json:"ntby_aspr_rsqn,string"`       // 순매수 호가 잔량
	NewMkopClsCode     string `json:"new_mkop_cls_code"`           // 신 장운영 구분 코드
}

// AskingPriceExpCcnExpected 은 응답의 output2 — 예상체결 + 시세.
type AskingPriceExpCcnExpected struct {
	AntcMkopClsCode    string          `json:"antc_mkop_cls_code"`         // 예상 장운영 구분 코드
	StckPrpr           decimal.Decimal `json:"stck_prpr"`                  // 주식 현재가
	StckOprc           decimal.Decimal `json:"stck_oprc"`                  // 주식 시가
	StckHgpr           decimal.Decimal `json:"stck_hgpr"`                  // 주식 최고가
	StckLwpr           decimal.Decimal `json:"stck_lwpr"`                  // 주식 최저가
	StckSdpr           decimal.Decimal `json:"stck_sdpr"`                  // 주식 기준가
	AntcCnpr           decimal.Decimal `json:"antc_cnpr"`                  // 예상 체결가
	AntcCntgVrssSign   string          `json:"antc_cntg_vrss_sign"`        // 예상 체결 대비 부호
	AntcCntgVrss       decimal.Decimal `json:"antc_cntg_vrss"`             // 예상 체결 대비
	AntcCntgPrdyCtrt   float64         `json:"antc_cntg_prdy_ctrt,string"` // 예상 체결 전일 대비율
	AntcVol            int64           `json:"antc_vol,string"`            // 예상 거래량
	StckShrnIscd       string          `json:"stck_shrn_iscd"`             // 주식 단축 종목코드
	ViClsCode          string          `json:"vi_cls_code"`                // VI 적용 구분 코드
}

// InquireAskingPriceExpCcnParams 는 호가/예상체결 조회 파라미터.
type InquireAskingPriceExpCcnParams struct {
	MarketCode string // FID_COND_MRKT_DIV_CODE — "J":KRX/"NX":NXT/"UN":통합. 빈 값=>"J"
	Symbol     string // FID_INPUT_ISCD — 종목코드 (예 "005930")
}

// InquireAskingPriceExpCcn 은 주식현재가 호가/예상체결 호출.
//
// 한투 docs: docs/api/국내주식/주식현재가_호가_예상체결.md
// path: /uapi/domestic-stock/v1/quotations/inquire-asking-price-exp-ccn (FHKST01010200)
//
// output1: 10 단계 호가/잔량/증감 + 시간외 + VI 등.
// output2: 예상체결가 + 시세.
func (c *Client) InquireAskingPriceExpCcn(ctx context.Context, params InquireAskingPriceExpCcnParams) (*AskingPriceExpCcn, error) {
	market := params.MarketCode
	if market == "" {
		market = "J"
	}

	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/domestic-stock/v1/quotations/inquire-asking-price-exp-ccn",
		TrID:   "FHKST01010200",
		Query: map[string]string{
			"FID_COND_MRKT_DIV_CODE": market,
			"FID_INPUT_ISCD":         params.Symbol,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}

	var res AskingPriceExpCcn
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse AskingPriceExpCcn: %w", err)
	}
	return &res, nil
}
```

- [ ] **Step 4: PASS**

`go test ./domestic/... -run InquireAskingPriceExpCcn -v`

- [ ] **Step 5: Commit**

```bash
git add domestic/quote.go domestic/quote_test.go
git commit -m "$(cat <<'EOF'
[feat] domestic — InquireAskingPriceExpCcn (호가/예상체결, FHKST01010200)

AskingPriceExpCcn (Output1 호가/잔량/증감 10 단계 + 시간외 + VI = ~60 필드,
Output2 예상체결 + 시세 = 13 필드) + InquireAskingPriceExpCcnParams.
MarketCode default "J" (KRX).

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 3: InquireCcnl

**Files:**
- Modify: `domestic/quote.go` (append)
- Modify: `domestic/quote_test.go` (append)

- [ ] **Step 1: 테스트 추가**

```go
func TestClient_InquireCcnl(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/quotations/inquire-ccnl`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "ccnl_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireCcnl(context.Background(), domestic.InquireCcnlParams{
		Symbol: "005930",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, "J", capturedQuery.Get("FID_COND_MRKT_DIV_CODE"))
	assert.Equal(t, "005930", capturedQuery.Get("FID_INPUT_ISCD"))

	require.Len(t, res.Output, 2)
	assert.Equal(t, "131542", res.Output[0].StckCntgHour)
	d, _ := decimal.NewFromString("75800")
	assert.True(t, d.Equal(res.Output[0].StckPrpr))
	assert.Equal(t, int64(12345), res.Output[0].CntgVol)
	assert.InDelta(t, 104.32, res.Output[0].TdayRltv, 0.01)
}
```

- [ ] **Step 2: FAIL**

- [ ] **Step 3: 구현 추가 — APPEND to `domestic/quote.go`**

```go
// Ccnl 은 주식현재가 체결 (FHKST01010300) 응답.
//
// 한투 docs: docs/api/국내주식/주식현재가_체결.md
// path: /uapi/domestic-stock/v1/quotations/inquire-ccnl
//
// 최근 체결 list (~30건). 체결강도 (tday_rltv) 포함.
type Ccnl struct {
	Output []CcnlItem `json:"output"`
}

// CcnlItem 은 체결 한 건.
type CcnlItem struct {
	StckCntgHour string          `json:"stck_cntg_hour"`     // 체결 시간 (HHMMSS)
	StckPrpr     decimal.Decimal `json:"stck_prpr"`          // 현재가
	PrdyVrss     decimal.Decimal `json:"prdy_vrss"`          // 전일 대비
	PrdyVrssSign string          `json:"prdy_vrss_sign"`     // 전일 대비 부호
	CntgVol      int64           `json:"cntg_vol,string"`    // 체결 거래량
	TdayRltv     float64         `json:"tday_rltv,string"`   // 당일 체결강도
	PrdyCtrt     float64         `json:"prdy_ctrt,string"`   // 전일 대비율
}

// InquireCcnlParams 는 체결 조회 파라미터.
type InquireCcnlParams struct {
	MarketCode string // FID_COND_MRKT_DIV_CODE — 빈 값=>"J"
	Symbol     string // FID_INPUT_ISCD
}

// InquireCcnl 은 주식현재가 체결 호출.
//
// 한투 docs: docs/api/국내주식/주식현재가_체결.md
// path: /uapi/domestic-stock/v1/quotations/inquire-ccnl (FHKST01010300)
func (c *Client) InquireCcnl(ctx context.Context, params InquireCcnlParams) (*Ccnl, error) {
	market := params.MarketCode
	if market == "" {
		market = "J"
	}

	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/domestic-stock/v1/quotations/inquire-ccnl",
		TrID:   "FHKST01010300",
		Query: map[string]string{
			"FID_COND_MRKT_DIV_CODE": market,
			"FID_INPUT_ISCD":         params.Symbol,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}

	var res Ccnl
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse Ccnl: %w", err)
	}
	return &res, nil
}
```

- [ ] **Step 4: PASS**

- [ ] **Step 5: Commit**

```bash
git add domestic/quote.go domestic/quote_test.go
git commit -m "$(cat <<'EOF'
[feat] domestic — InquireCcnl (체결, FHKST01010300)

Ccnl + CcnlItem (7 필드) + InquireCcnlParams. 최근 체결 list 와 체결강도.

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 4: InquireDailyPrice

**Files:**
- Modify: `domestic/quote.go` (append)
- Modify: `domestic/quote_test.go` (append)

> **NOTE**: domestic 패키지에 이미 `InquireDailyItemChartPrice` (Phase 1.2) 가 있음. `InquireDailyPrice` 는 이름 충돌 X — 다른 메서드. KIS API path 도 다름 (`inquire-daily-price` vs `inquire-daily-itemchartprice`).

- [ ] **Step 1: 테스트 추가**

```go
func TestClient_InquireDailyPrice(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/quotations/inquire-daily-price`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "daily_price_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireDailyPrice(context.Background(), domestic.InquireDailyPriceParams{
		Symbol: "005930",
		Period: "D",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, "J", capturedQuery.Get("FID_COND_MRKT_DIV_CODE"))
	assert.Equal(t, "005930", capturedQuery.Get("FID_INPUT_ISCD"))
	assert.Equal(t, "D", capturedQuery.Get("FID_PERIOD_DIV_CODE"))
	assert.Equal(t, "0", capturedQuery.Get("FID_ORG_ADJ_PRC")) // default 미반영

	require.Len(t, res.Output, 2)
	assert.Equal(t, "20260505", res.Output[0].StckBsopDate)
	d, _ := decimal.NewFromString("75800")
	assert.True(t, d.Equal(res.Output[0].StckClpr))
	assert.Equal(t, int64(12345678), res.Output[0].AcmlVol)
	assert.InDelta(t, 53.42, res.Output[0].HtsFrgnEhrt, 0.01)
}

func TestClient_InquireDailyPrice_OriginalPrice(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/quotations/inquire-daily-price`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "daily_price_success.json")), nil
		},
	)

	c := newTestClient(t)
	_, err := c.InquireDailyPrice(context.Background(), domestic.InquireDailyPriceParams{
		Symbol:        "005930",
		Period:        "W",
		OriginalPrice: true,
	})
	require.NoError(t, err)
	assert.Equal(t, "W", capturedQuery.Get("FID_PERIOD_DIV_CODE"))
	assert.Equal(t, "1", capturedQuery.Get("FID_ORG_ADJ_PRC")) // 1 = 원주가
}
```

- [ ] **Step 2: FAIL**

- [ ] **Step 3: 구현 추가**

```go
// DailyPrice 는 주식현재가 일자별 (FHKST01010400) 응답.
//
// 한투 docs: docs/api/국내주식/주식현재가_일자별.md
// path: /uapi/domestic-stock/v1/quotations/inquire-daily-price
//
// 최근 30 거래일 (D) / 30 주 (W) / 30 개월 (M) 일별 시세 + 외국인 순매수 + 락 구분.
// 단순 일자별 — Phase 1.2 의 InquireDailyItemChartPrice (chart) 와는 다른 endpoint.
type DailyPrice struct {
	Output []DailyPriceItem `json:"output"`
}

// DailyPriceItem 은 한 일자/주/월 봉.
type DailyPriceItem struct {
	StckBsopDate     string          `json:"stck_bsop_date"`            // 영업 일자
	StckOprc         decimal.Decimal `json:"stck_oprc"`                 // 시가
	StckHgpr         decimal.Decimal `json:"stck_hgpr"`                 // 최고가
	StckLwpr         decimal.Decimal `json:"stck_lwpr"`                 // 최저가
	StckClpr         decimal.Decimal `json:"stck_clpr"`                 // 종가
	AcmlVol          int64           `json:"acml_vol,string"`           // 누적 거래량
	PrdyVrssVolRate  float64         `json:"prdy_vrss_vol_rate,string"` // 전일 대비 거래량 비율
	PrdyVrss         decimal.Decimal `json:"prdy_vrss"`                 // 전일 대비
	PrdyVrssSign     string          `json:"prdy_vrss_sign"`            // 전일 대비 부호
	PrdyCtrt         float64         `json:"prdy_ctrt,string"`          // 전일 대비율
	HtsFrgnEhrt      float64         `json:"hts_frgn_ehrt,string"`      // HTS 외국인 소진율
	FrgnNtbyQty      int64           `json:"frgn_ntby_qty,string"`      // 외국인 순매수 수량
	FlngClsCode      string          `json:"flng_cls_code"`             // 락 구분 (00=일반, 01=권리, 02=배당, ...)
	AcmlPrttRate     float64         `json:"acml_prtt_rate,string"`     // 누적 분할 비율
}

// InquireDailyPriceParams 는 일자별 조회 파라미터.
type InquireDailyPriceParams struct {
	MarketCode    string // FID_COND_MRKT_DIV_CODE — 빈 값=>"J"
	Symbol        string // FID_INPUT_ISCD
	Period        string // FID_PERIOD_DIV_CODE — "D":일/"W":주/"M":월. 빈 값=>"D"
	OriginalPrice bool   // FID_ORG_ADJ_PRC — false=>"0"(수정주가 미반영, default), true=>"1"(원주가)
}

// InquireDailyPrice 는 주식현재가 일자별 호출.
//
// 한투 docs: docs/api/국내주식/주식현재가_일자별.md
// path: /uapi/domestic-stock/v1/quotations/inquire-daily-price (FHKST01010400)
//
// 최근 30 거래일 (D) / 30 주 (W) / 30 개월 (M). Phase 1.2 의
// InquireDailyItemChartPrice (chart) 와는 다른 endpoint — 외국인 소진율/락 구분 등 추가.
func (c *Client) InquireDailyPrice(ctx context.Context, params InquireDailyPriceParams) (*DailyPrice, error) {
	market := params.MarketCode
	if market == "" {
		market = "J"
	}
	period := params.Period
	if period == "" {
		period = "D"
	}
	adjPrc := "0" // 수정주가 미반영 default
	if params.OriginalPrice {
		adjPrc = "1"
	}

	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/domestic-stock/v1/quotations/inquire-daily-price",
		TrID:   "FHKST01010400",
		Query: map[string]string{
			"FID_COND_MRKT_DIV_CODE": market,
			"FID_INPUT_ISCD":         params.Symbol,
			"FID_PERIOD_DIV_CODE":    period,
			"FID_ORG_ADJ_PRC":        adjPrc,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}

	var res DailyPrice
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse DailyPrice: %w", err)
	}
	return &res, nil
}
```

- [ ] **Step 4: PASS**

`go test ./domestic/... -run InquireDailyPrice -v` — 2 cases PASS.

- [ ] **Step 5: 전체 회귀 테스트**

`go test ./... -count=1` — all PASS.

- [ ] **Step 6: Commit**

```bash
git add domestic/quote.go domestic/quote_test.go
git commit -m "$(cat <<'EOF'
[feat] domestic — InquireDailyPrice (주식현재가 일자별, FHKST01010400)

DailyPrice + DailyPriceItem (14 필드) + InquireDailyPriceParams (Period D/W/M
+ OriginalPrice bool). 최근 30 거래일/주/월 + 외국인 소진율/락 구분 코드.
Phase 1.2 의 InquireDailyItemChartPrice 와 다른 endpoint.

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 5: examples/domestic_quote/main.go

**Files:**
- Create: `examples/domestic_quote/main.go`

- [ ] **Step 1: example 작성**

```go
// domestic_quote example: InquireAskingPriceExpCcn + InquireCcnl + InquireDailyPrice.
//
// Run: KIS credentials env vars 후 go run ./examples/domestic_quote
package main

import (
	"context"
	"fmt"
	"log"

	kis "github.com/kenshin579/korea-investment-stock"
	"github.com/kenshin579/korea-investment-stock/domestic"
)

func main() {
	client, err := kis.NewClientFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	symbol := "005930"

	// 1. 호가/예상체결
	ob, err := client.Domestic.InquireAskingPriceExpCcn(ctx, domestic.InquireAskingPriceExpCcnParams{
		Symbol: symbol,
	})
	if err != nil {
		log.Fatalf("InquireAskingPriceExpCcn: %v", err)
	}
	fmt.Printf("[%s] 호가 (접수시간 %s)\n", symbol, ob.Output1.AsprAcptHour)
	fmt.Printf("  매도1: %s @ %d, 매수1: %s @ %d\n",
		ob.Output1.Askp1, ob.Output1.AskpRsqn1, ob.Output1.Bidp1, ob.Output1.BidpRsqn1)
	fmt.Printf("  총 매도잔량 %d, 총 매수잔량 %d\n",
		ob.Output1.TotalAskpRsqn, ob.Output1.TotalBidpRsqn)
	fmt.Printf("  예상체결: %s (전일대비 %s, %v%%)\n",
		ob.Output2.AntcCnpr, ob.Output2.AntcCntgVrss, ob.Output2.AntcCntgPrdyCtrt)

	// 2. 최근 체결
	cc, err := client.Domestic.InquireCcnl(ctx, domestic.InquireCcnlParams{
		Symbol: symbol,
	})
	if err != nil {
		log.Fatalf("InquireCcnl: %v", err)
	}
	fmt.Printf("\n[%s] 최근 체결 %d 건\n", symbol, len(cc.Output))
	for i, item := range cc.Output {
		if i >= 5 {
			fmt.Println("  ... (이하 생략)")
			break
		}
		fmt.Printf("  %s: %s원 @ %d주 (체결강도 %v)\n",
			item.StckCntgHour, item.StckPrpr, item.CntgVol, item.TdayRltv)
	}

	// 3. 일자별
	dp, err := client.Domestic.InquireDailyPrice(ctx, domestic.InquireDailyPriceParams{
		Symbol: symbol,
		Period: "D",
	})
	if err != nil {
		log.Fatalf("InquireDailyPrice: %v", err)
	}
	fmt.Printf("\n[%s] 최근 일자별 %d 일\n", symbol, len(dp.Output))
	for i, item := range dp.Output {
		if i >= 5 {
			fmt.Println("  ... (이하 생략)")
			break
		}
		fmt.Printf("  %s: O=%s H=%s L=%s C=%s V=%d 외인소진=%v%%\n",
			item.StckBsopDate, item.StckOprc, item.StckHgpr, item.StckLwpr,
			item.StckClpr, item.AcmlVol, item.HtsFrgnEhrt)
	}
}
```

- [ ] **Step 2: 컴파일 검증**

`go build ./examples/domestic_quote && echo OK`

- [ ] **Step 3: Commit**

```bash
git add examples/domestic_quote
git commit -m "$(cat <<'EOF'
[feat] examples/domestic_quote — 호가 + 체결 + 일자별 통합 예시

삼성전자 (005930) 의 10단계 호가/잔량, 최근 체결, 일별 시세 출력.

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 6: 문서 갱신

**Files:**
- Modify: `CLAUDE.md`
- Modify: `README.md`
- Modify: `CHANGELOG.md`
- Modify: `domestic/doc.go`

- [ ] **Step 1: README.md — Available Methods 표 갱신**

Find existing `## Available Methods (Phase 1.2 + 1.3 + 1.4 + 1.5)` heading. Update heading to `## Available Methods (Phase 1.2 ~ 2.1)` and APPEND 3 rows AT THE END:

```markdown
| `Domestic.InquireAskingPriceExpCcn` | `quotations/inquire-asking-price-exp-ccn` | FHKST01010200 |
| `Domestic.InquireCcnl` | `quotations/inquire-ccnl` | FHKST01010300 |
| `Domestic.InquireDailyPrice` | `quotations/inquire-daily-price` | FHKST01010400 |
```

- [ ] **Step 2: CLAUDE.md — banner 갱신**

Replace `> **Phase 1.5 — 해외주식 (v1.3.0, Python parity 완성).** Phase 2+ 는 추후 결정.` with:
```
> **Phase 2.1 — 국내 호가/체결 (v1.4.0).** Phase 2.2+ 는 추후 sub-plan 으로.
```

ADD spec link bullet:
```markdown
- Phase 2 design spec: [`docs/superpowers/specs/2026-05-05-phase2-readonly-extension-design.md`](docs/superpowers/specs/2026-05-05-phase2-readonly-extension-design.md)
- Phase 2.1 implementation plan: [`docs/superpowers/specs/2026-05-05-phase2-1-domestic-quote-implementation-plan.md`](docs/superpowers/specs/2026-05-05-phase2-1-domestic-quote-implementation-plan.md)
```

- [ ] **Step 3: CHANGELOG.md — `[1.4.0]` entry**

ADD AT THE TOP (above `## [1.3.0]`):

```markdown
## [1.4.0] - 2026-05-05

### Added — Phase 2.1 (국내 호가/체결)

- `Domestic.InquireAskingPriceExpCcn` — 주식현재가 호가/예상체결 (FHKST01010200) — 10단계 호가/잔량/증감 + 시간외 + 예상체결
- `Domestic.InquireCcnl` — 주식현재가 체결 (FHKST01010300) — 최근 체결 list + 체결강도
- `Domestic.InquireDailyPrice` — 주식현재가 일자별 (FHKST01010400) — 최근 30 거래일/주/월 + 외국인 소진율 + 락 구분
- examples: `domestic_quote`

### Notes

- `InquireDailyPrice` 는 Phase 1.2 의 `InquireDailyItemChartPrice` 와 다른 endpoint — 외국인 소진율, 락 구분 등 추가 필드 포함
- Phase 2 시작 — Python wrapper 가 cover 하지 않은 KIS read-only API 확장 (Phase 2.1~2.4 sub-phase)
```

- [ ] **Step 4: domestic/doc.go 갱신**

ADD Phase 2.1 section after Phase 1.4:

```go
// Phase 2.1 메서드 (3):
//
//   - InquireAskingPriceExpCcn  — 주식현재가 호가/예상체결 (FHKST01010200)
//   - InquireCcnl               — 주식현재가 체결 (FHKST01010300)
//   - InquireDailyPrice         — 주식현재가 일자별 (FHKST01010400)
```

- [ ] **Step 5: 검증**

```bash
go build ./... && go vet ./... && gofmt -l .
```
Expected: silent.

- [ ] **Step 6: Commit**

```bash
git add CLAUDE.md README.md CHANGELOG.md domestic/doc.go
git commit -m "$(cat <<'EOF'
[doc] Phase 2.1 메서드 문서 갱신 — CLAUDE/README/CHANGELOG/domestic/doc.go

Phase 2.1 의 3 메서드 (호가/체결/일자별) 목록 + CHANGELOG [1.4.0] entry.
CLAUDE.md banner 갱신 (Phase 1.5 → 2.1, v1.3.0 → v1.4.0). domestic/doc.go
패키지 doc 에 Phase 2.1 section 추가.

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 7: 최종 점검

- [ ] **Step 1: gofmt cleanup (필요 시)**

`gofmt -w domestic/*.go && gofmt -l .` empty.

- [ ] **Step 2: 빌드/vet/test**

```bash
go build ./... && go vet ./... && go test ./... -race -count=1
```
Expected: silent + all PASS.

- [ ] **Step 3: Coverage**

```bash
go test ./... -coverprofile=/tmp/cov.out -covermode=atomic
go tool cover -func=/tmp/cov.out | tail -10
```
Expected: domestic/ ≥ 80%, root kis ≥ 80%.

- [ ] **Step 4: 디렉터리 구조 확인**

```bash
ls -la domestic/quote.go domestic/quote_test.go domestic/testdata/{asking_price_exp_ccn,ccnl,daily_price}_success.json examples/domestic_quote/main.go 2>&1 | wc -l
```
Expected: 6 files.

- [ ] **Step 5: Commit history**

`git log main..HEAD --oneline | wc -l` — should be ~7-9.

---

## Task 8: PR 생성 (사용자 승인 후)

> Claude 는 push / PR 생성을 사용자 명시적 승인 후에만 실행 (글로벌 정책).

- [ ] **Step 1: 사용자 승인 요청**

작업 진행 보고 + PR 생성 가능 여부 confirm.

- [ ] **Step 2: Push branch**

`git push -u origin docs/phase2-spec`

- [ ] **Step 3: PR 생성**

```bash
gh pr create --title "Phase 2.1 — 국내 호가/체결 (v1.4.0)" --reviewer kenshin579 --base main --head docs/phase2-spec --body "$(cat <<'EOF'
## Summary

- 국내주식 호가/체결 디테일 3 메서드 추가 (Phase 2 의 첫 sub-phase)
- Phase 1 패턴 그대로 재사용 (Style A, Params struct, KIS docs 1:1)
- v1.4.0 release 대상

## 메서드 → 한투 API 매핑

| Go 메서드 | path | TR_ID |
|---|---|---|
| InquireAskingPriceExpCcn | quotations/inquire-asking-price-exp-ccn | FHKST01010200 |
| InquireCcnl | quotations/inquire-ccnl | FHKST01010300 |
| InquireDailyPrice | quotations/inquire-daily-price | FHKST01010400 |

## Test Plan

- [x] go build/vet/fmt clean
- [x] go test ./... -race -count=1 모든 패키지 PASS
- [x] Coverage domestic/ ≥ 80%
- [x] httpmock 단위 테스트 (3 메서드 + InquireDailyPrice variant)

## Breaking Changes

없음 — 신규 메서드 추가만.

## 참고 문서

- Phase 2 design spec: \`docs/superpowers/specs/2026-05-05-phase2-readonly-extension-design.md\`
- Phase 2.1 implementation plan: \`docs/superpowers/specs/2026-05-05-phase2-1-domestic-quote-implementation-plan.md\`

🤖 Generated with [Claude Code](https://claude.com/claude-code)
EOF
)"
```

- [ ] **Step 4: Merge (사용자 승인 후)** — `gh pr merge <PR#> --merge`

- [ ] **Step 5: 후속 작업** — `git tag -a v1.4.0 -m "..."`, `git push origin v1.4.0`, `gh release create v1.4.0 --title "..." --notes-file <(awk '/^## \[/{c++} c==1' CHANGELOG.md)`
