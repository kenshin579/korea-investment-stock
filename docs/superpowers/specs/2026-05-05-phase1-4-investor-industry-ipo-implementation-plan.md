# Phase 1.4 — 국내 투자자 + 업종 + IPO Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** `korea-investment-stock` Go 라이브러리에 국내주식 투자자 매매동향 (3) + 업종 (2) + IPO (1) 총 6 메서드 추가 (`v1.2.0` release).

**Architecture:** Phase 1.2/1.3 의 인프라 + 패턴 재사용. `domestic/investor.go` (3 메서드) + `domestic/industry.go` (2 메서드) + `domestic/ipo.go` (1 메서드) 추가. 한투 API path 1:1 매핑 (Style A — endpoint path 의 마지막 segment 를 PascalCase). 새 internal package 불필요. IPO helpers 9개 미반영 (Phase 1.2 amendment 의 "Python wrapper convenience 미반영" 정책 일관). TDD: testdata fixture (한투 docs 응답 필드 정의 → 합성 JSON) → 실패 테스트 → struct + 메서드 구현 → 통과 → commit.

**Tech Stack:** Go 1.25+, `github.com/jarcoal/httpmock`, `github.com/shopspring/decimal`, `github.com/stretchr/testify`. 새 dependency 없음.

**참고 spec:**
- Phase 1 design spec (Phase 1.4 amendment 적용): `docs/superpowers/specs/2026-05-03-phase1-api-coverage-design.md` (commit `80b6ed3`)
- Phase 1.3 plan (참조 패턴): `docs/superpowers/specs/2026-05-04-phase1-3-ranking-financial-implementation-plan.md`
- 한투 API docs: `docs/api/국내주식/{종목별_투자자매매동향(일별).md, 시장별_투자자매매동향(일별).md, 시장별_투자자매매동향(시세).md, 국내업종_현재지수.md, 국내업종_구분별전체시세.md, 예탁원정보(공모주청약일정).md}`

---

## 사전 정보

| 항목 | 값 |
|------|---|
| 작업 브랜치 | `docs/phase1-4-spec` (이미 생성됨) |
| 시작 HEAD | `80b6ed3` (Phase 1 design spec amendment commit) |
| Release 목표 | `v1.2.0` (PR merge 후 태그) |
| PR 베이스 | `main` |
| 현재 main 상태 | v1.1.0 publish 완료 (Phase 1.1+1.2+1.3 통합, 16 메서드) |

---

## 메서드 → 한투 API 매핑

| Go 메서드 | 한투 path | TR_ID | docs |
|-----------|----------|-------|------|
| `Domestic.InquireInvestorTradeByStockDaily(ctx, params)` | `/uapi/domestic-stock/v1/quotations/investor-trade-by-stock-daily` | FHPTJ04160001 | 종목별_투자자매매동향(일별).md |
| `Domestic.InquireInvestorDailyByMarket(ctx, params)` | `/uapi/domestic-stock/v1/quotations/inquire-investor-daily-by-market` | FHPTJ04040000 | 시장별_투자자매매동향(일별).md |
| `Domestic.InquireInvestorTimeByMarket(ctx, params)` | `/uapi/domestic-stock/v1/quotations/inquire-investor-time-by-market` | FHPTJ04030000 | 시장별_투자자매매동향(시세).md |
| `Domestic.InquireIndexPrice(ctx, params)` | `/uapi/domestic-stock/v1/quotations/inquire-index-price` | FHPUP02100000 | 국내업종_현재지수.md |
| `Domestic.InquireIndexCategoryPrice(ctx, params)` | `/uapi/domestic-stock/v1/quotations/inquire-index-category-price` | FHPUP02140000 | 국내업종_구분별전체시세.md |
| `Domestic.InquirePubOffer(ctx, params)` | `/uapi/domestic-stock/v1/ksdinfo/pub-offer` | HHKDB669108C0 | 예탁원정보(공모주청약일정).md |

**참고**: `InquirePubOffer` 만 `ksdinfo/` path (예탁원 정보), 다른 5개는 `quotations/`. PubOffer 의 query 키도 다른 형식 (`SHT_CD`, `CTS`, `F_DT`, `T_DT` — 대문자 + 한글식). KIS docs 그대로 노출.

---

## 파일 구조

### 신규 (domestic)
- `domestic/investor.go` — 3 investor 메서드 + Result/Item structs + Params structs
- `domestic/investor_test.go`
- `domestic/industry.go` — 2 industry 메서드 + Result/Item structs + Params structs
- `domestic/industry_test.go`
- `domestic/ipo.go` — InquirePubOffer + PubOffer + PubOfferItem + InquirePubOfferParams
- `domestic/ipo_test.go`
- `domestic/testdata/investor_trade_by_stock_daily_success.json`
- `domestic/testdata/investor_daily_by_market_success.json`
- `domestic/testdata/investor_time_by_market_success.json`
- `domestic/testdata/index_price_success.json`
- `domestic/testdata/index_category_price_success.json`
- `domestic/testdata/pub_offer_success.json`

### 신규 (examples)
- `examples/domestic_investor/main.go` — InvestorTradeByStockDaily + InvestorDailyByMarket 예제

### 수정 (root)
- `CLAUDE.md` — Phase 1.4 메서드 안내
- `README.md` — Available Methods 표 갱신 (16 → 22 메서드)
- `CHANGELOG.md` — `[1.2.0]` entry
- `domestic/doc.go` — Phase 1.4 메서드 안내 갱신

---

## 타입 매핑 규칙 (Phase 1.2/1.3 와 동일)

- **가격/지수/주당 → `decimal.Decimal` (bare tag)**: `stck_prpr`, `stck_clpr`, `stck_oprc`, `stck_hgpr`, `stck_lwpr`, `prdy_vrss`, `bstp_nmix_prpr`, `bstp_nmix_oprc`, `bstp_nmix_hgpr`, `bstp_nmix_lwpr`, `bstp_nmix_prdy_vrss`, `prdy_nmix_vrss_nmix_*`, `dryy_bstp_nmix_*pr`, `prdy_clpr_vrss_lwpr`, `fix_subscr_pri`, `face_value`
- **수량/금액 → `int64,string`**: `acml_vol`, `prdy_vol`, `acml_tr_pbmn`, `prdy_tr_pbmn`, `*_ntby_qty`, `*_ntby_tr_pbmn`, `*_ntby_vol`, `*_seln_vol`, `*_shnu_vol`, `*_seln_tr_pbmn`, `*_shnu_tr_pbmn`, `*_askp_qty`, `*_bidp_qty`, `*_askp_pbmn`, `*_bidp_pbmn`, `total_askp_rsqn`, `total_bidp_rsqn`, `ntby_rsqn`, `ascn_issu_cnt`, `down_issu_cnt`, `stnr_issu_cnt`, `uplm_issu_cnt`, `lslm_issu_cnt`, `pub_bf_cap`, `pub_af_cap`, `assign_stk_qty`
- **비율 → `float64,string`**: `prdy_ctrt`, `bstp_nmix_prdy_ctrt`, `bstp_nmix_*_prdy_ctrt`, `prdy_clpr_vrss_lwpr_rate`, `dryy_*_rate`, `seln_rsqn_rate`, `shnu_rsqn_rate`, `acml_vol_rlim`, `acml_tr_pbmn_rlim`
- **코드/이름/날짜/Y-N/부호 → 평문 `string`**: `prdy_vrss_sign`, `oprc_vrss_prpr_sign`, `hgpr_vrss_prpr_sign`, `lwpr_vrss_prpr_sign`, `stck_bsop_date`, `dryy_bstp_nmix_*_date`, `bstp_cls_code`, `hts_kor_isnm`, `rprs_mrkt_kor_name`, `record_date`, `sht_cd`, `isin_name`, `subscr_dt`, `pay_dt`, `refund_dt`, `list_dt`, `lead_mgr`

---

## Task 1: testdata fixtures (6 합성 JSON)

**Files (Create):**
- `domestic/testdata/investor_trade_by_stock_daily_success.json`
- `domestic/testdata/investor_daily_by_market_success.json`
- `domestic/testdata/investor_time_by_market_success.json`
- `domestic/testdata/index_price_success.json`
- `domestic/testdata/index_category_price_success.json`
- `domestic/testdata/pub_offer_success.json`

> 각 fixture 는 KIS docs 응답 필드 정의 기반 합성. **투자자 매매동향 응답은 80~100 필드 수준으로 매우 크므로** testdata 에는 핵심 필드 (가격류 + 외국인/개인/기관 ntby_qty + 합계 거래량) 만 포함. 누락된 KIS field 는 Go json unmarshal 시 zero value 자동 처리. struct 정의에는 모든 KIS field 포함.

- [ ] **Step 1: investor_trade_by_stock_daily_success.json**

```json
{
  "rt_cd": "0",
  "msg_cd": "MCA00000",
  "msg1": "정상처리 되었습니다.",
  "output1": {
    "stck_prpr": "75800",
    "prdy_vrss": "-200",
    "prdy_vrss_sign": "5",
    "prdy_ctrt": "-0.26",
    "acml_vol": "12345678",
    "prdy_vol": "11000000",
    "rprs_mrkt_kor_name": "KOSPI200"
  },
  "output2": [
    {
      "stck_bsop_date": "20260505",
      "stck_clpr": "75800",
      "stck_oprc": "76000",
      "stck_hgpr": "76200",
      "stck_lwpr": "75500",
      "prdy_vrss": "-200",
      "prdy_vrss_sign": "5",
      "prdy_ctrt": "-0.26",
      "acml_vol": "12345678",
      "acml_tr_pbmn": "938223456",
      "frgn_ntby_qty": "-123456",
      "frgn_ntby_tr_pbmn": "-9357",
      "prsn_ntby_qty": "234567",
      "prsn_ntby_tr_pbmn": "17783",
      "orgn_ntby_qty": "-100000",
      "orgn_ntby_tr_pbmn": "-7580"
    },
    {
      "stck_bsop_date": "20260504",
      "stck_clpr": "76000",
      "stck_oprc": "76500",
      "stck_hgpr": "76800",
      "stck_lwpr": "75900",
      "prdy_vrss": "100",
      "prdy_vrss_sign": "2",
      "prdy_ctrt": "0.13",
      "acml_vol": "11000000",
      "acml_tr_pbmn": "836000000",
      "frgn_ntby_qty": "50000",
      "frgn_ntby_tr_pbmn": "3800",
      "prsn_ntby_qty": "-30000",
      "prsn_ntby_tr_pbmn": "-2280",
      "orgn_ntby_qty": "-15000",
      "orgn_ntby_tr_pbmn": "-1140"
    }
  ]
}
```

- [ ] **Step 2: investor_daily_by_market_success.json**

```json
{
  "rt_cd": "0",
  "msg_cd": "MCA00000",
  "msg1": "정상처리 되었습니다.",
  "output": [
    {
      "stck_bsop_date": "20260505",
      "bstp_nmix_prpr": "2650.45",
      "bstp_nmix_prdy_vrss": "-12.30",
      "prdy_vrss_sign": "5",
      "bstp_nmix_prdy_ctrt": "-0.46",
      "bstp_nmix_oprc": "2660.00",
      "bstp_nmix_hgpr": "2665.20",
      "bstp_nmix_lwpr": "2645.10",
      "stck_prdy_clpr": "2662.75",
      "frgn_ntby_qty": "-123456",
      "prsn_ntby_qty": "234567",
      "orgn_ntby_qty": "-100000",
      "frgn_ntby_tr_pbmn": "-9357",
      "prsn_ntby_tr_pbmn": "17783",
      "orgn_ntby_tr_pbmn": "-7580"
    }
  ]
}
```

- [ ] **Step 3: investor_time_by_market_success.json**

```json
{
  "rt_cd": "0",
  "msg_cd": "MCA00000",
  "msg1": "정상처리 되었습니다.",
  "output": {
    "frgn_seln_vol": "5000000",
    "frgn_shnu_vol": "4876544",
    "frgn_ntby_qty": "-123456",
    "frgn_seln_tr_pbmn": "379500",
    "frgn_shnu_tr_pbmn": "370143",
    "frgn_ntby_tr_pbmn": "-9357",
    "prsn_seln_vol": "10000000",
    "prsn_shnu_vol": "10234567",
    "prsn_ntby_qty": "234567",
    "prsn_seln_tr_pbmn": "759000",
    "prsn_shnu_tr_pbmn": "776783",
    "prsn_ntby_tr_pbmn": "17783",
    "orgn_seln_vol": "3000000",
    "orgn_shnu_vol": "2900000",
    "orgn_ntby_qty": "-100000",
    "orgn_seln_tr_pbmn": "227700",
    "orgn_shnu_tr_pbmn": "220120",
    "orgn_ntby_tr_pbmn": "-7580"
  }
}
```

- [ ] **Step 4: index_price_success.json**

```json
{
  "rt_cd": "0",
  "msg_cd": "MCA00000",
  "msg1": "정상처리 되었습니다.",
  "output": {
    "bstp_nmix_prpr": "2650.45",
    "bstp_nmix_prdy_vrss": "-12.30",
    "prdy_vrss_sign": "5",
    "bstp_nmix_prdy_ctrt": "-0.46",
    "acml_vol": "350000000",
    "prdy_vol": "330000000",
    "acml_tr_pbmn": "9500000",
    "prdy_tr_pbmn": "9000000",
    "bstp_nmix_oprc": "2660.00",
    "prdy_nmix_vrss_nmix_oprc": "-2.75",
    "oprc_vrss_prpr_sign": "5",
    "bstp_nmix_oprc_prdy_ctrt": "-0.10",
    "bstp_nmix_hgpr": "2665.20",
    "prdy_nmix_vrss_nmix_hgpr": "2.45",
    "hgpr_vrss_prpr_sign": "5",
    "bstp_nmix_hgpr_prdy_ctrt": "0.09",
    "bstp_nmix_lwpr": "2645.10",
    "prdy_clpr_vrss_lwpr": "-17.65",
    "lwpr_vrss_prpr_sign": "2",
    "prdy_clpr_vrss_lwpr_rate": "-0.66",
    "ascn_issu_cnt": "315",
    "uplm_issu_cnt": "5",
    "stnr_issu_cnt": "120",
    "down_issu_cnt": "450",
    "lslm_issu_cnt": "2",
    "dryy_bstp_nmix_hgpr": "2780.00",
    "dryy_hgpr_vrss_prpr_rate": "-4.66",
    "dryy_bstp_nmix_hgpr_date": "20260301",
    "dryy_bstp_nmix_lwpr": "2480.50",
    "dryy_lwpr_vrss_prpr_rate": "6.85",
    "dryy_bstp_nmix_lwpr_date": "20260115",
    "total_askp_rsqn": "12345678",
    "total_bidp_rsqn": "10234567",
    "seln_rsqn_rate": "54.68",
    "shnu_rsqn_rate": "45.32",
    "ntby_rsqn": "-2111111"
  }
}
```

- [ ] **Step 5: index_category_price_success.json**

```json
{
  "rt_cd": "0",
  "msg_cd": "MCA00000",
  "msg1": "정상처리 되었습니다.",
  "output1": {
    "bstp_nmix_prpr": "2650.45",
    "bstp_nmix_prdy_vrss": "-12.30",
    "prdy_vrss_sign": "5",
    "bstp_nmix_prdy_ctrt": "-0.46",
    "acml_vol": "350000000",
    "acml_tr_pbmn": "9500000",
    "bstp_nmix_oprc": "2660.00",
    "bstp_nmix_hgpr": "2665.20",
    "bstp_nmix_lwpr": "2645.10",
    "prdy_vol": "330000000",
    "ascn_issu_cnt": "315",
    "down_issu_cnt": "450",
    "stnr_issu_cnt": "120",
    "uplm_issu_cnt": "5",
    "lslm_issu_cnt": "2",
    "prdy_tr_pbmn": "9000000",
    "dryy_bstp_nmix_hgpr_date": "20260301",
    "dryy_bstp_nmix_hgpr": "2780.00",
    "dryy_bstp_nmix_lwpr": "2480.50",
    "dryy_bstp_nmix_lwpr_date": "20260115"
  },
  "output2": [
    {
      "bstp_cls_code": "0001",
      "hts_kor_isnm": "코스피",
      "bstp_nmix_prpr": "2650.45",
      "bstp_nmix_prdy_vrss": "-12.30",
      "prdy_vrss_sign": "5",
      "bstp_nmix_prdy_ctrt": "-0.46",
      "acml_vol": "350000000",
      "acml_tr_pbmn": "9500000",
      "acml_vol_rlim": "100.00",
      "acml_tr_pbmn_rlim": "100.00"
    },
    {
      "bstp_cls_code": "0002",
      "hts_kor_isnm": "대형주",
      "bstp_nmix_prpr": "2700.10",
      "bstp_nmix_prdy_vrss": "-15.20",
      "prdy_vrss_sign": "5",
      "bstp_nmix_prdy_ctrt": "-0.56",
      "acml_vol": "150000000",
      "acml_tr_pbmn": "5500000",
      "acml_vol_rlim": "42.86",
      "acml_tr_pbmn_rlim": "57.89"
    }
  ]
}
```

- [ ] **Step 6: pub_offer_success.json**

```json
{
  "rt_cd": "0",
  "msg_cd": "MCA00000",
  "msg1": "정상처리 되었습니다.",
  "output1": [
    {
      "record_date": "20260505",
      "sht_cd": "999998",
      "isin_name": "샘플바이오",
      "fix_subscr_pri": "30000",
      "face_value": "100",
      "subscr_dt": "20260505 ~ 20260506",
      "pay_dt": "20260508",
      "refund_dt": "20260510",
      "list_dt": "20260520",
      "lead_mgr": "한국투자증권",
      "pub_bf_cap": "5000",
      "pub_af_cap": "8000",
      "assign_stk_qty": "100000"
    },
    {
      "record_date": "20260510",
      "sht_cd": "999997",
      "isin_name": "샘플테크",
      "fix_subscr_pri": "15000",
      "face_value": "500",
      "subscr_dt": "20260510 ~ 20260511",
      "pay_dt": "20260513",
      "refund_dt": "20260515",
      "list_dt": "20260525",
      "lead_mgr": "한국투자증권",
      "pub_bf_cap": "3000",
      "pub_af_cap": "4500",
      "assign_stk_qty": "75000"
    }
  ]
}
```

- [ ] **Step 7: 검증**

```bash
for f in domestic/testdata/{investor_trade_by_stock_daily,investor_daily_by_market,investor_time_by_market,index_price,index_category_price,pub_offer}_success.json; do
  python3 -c "import json; json.load(open('$f'))" && echo "$f OK" || echo "$f BROKEN"
done
```
Expected: 6 줄 모두 `OK`.

- [ ] **Step 8: Commit**

```bash
git add domestic/testdata/{investor_trade_by_stock_daily,investor_daily_by_market,investor_time_by_market,index_price,index_category_price,pub_offer}_success.json
git commit -m "$(cat <<'EOF'
[chore] Phase 1.4 testdata — 6 합성 JSON fixtures

investor 3 (trade_by_stock_daily, daily_by_market, time_by_market) +
industry 2 (index_price, index_category_price) + ipo 1 (pub_offer).
한투 docs (docs/api/국내주식/<API>.md) 의 응답 필드 정의 기반 합성.
거대 응답 (investor 80~100 필드) 은 핵심 필드만 포함, 나머지는 Go json
unmarshal 의 zero-value default 활용.

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 2: domestic/investor.go — InquireInvestorTradeByStockDaily + 공통 base

**Files:**
- Create: `domestic/investor.go`
- Create: `domestic/investor_test.go`

> **NOTE**: 종목별 투자자매매동향(일별) 응답이 가장 큼. output1 (요약 7 fields) + output2 (Array, 한 행에 ~95 fields). KIS docs (`docs/api/국내주식/종목별_투자자매매동향(일별).md`) line 86~180 의 모든 필드를 struct 에 1:1 매핑. testdata 는 핵심만 (Task 1 Step 1) — Go json unmarshal 이 누락 필드는 zero value 처리.

- [ ] **Step 1: 테스트 작성** — `domestic/investor_test.go`

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

func TestClient_InquireInvestorTradeByStockDaily(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/quotations/investor-trade-by-stock-daily`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "investor_trade_by_stock_daily_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireInvestorTradeByStockDaily(context.Background(), domestic.InquireInvestorTradeByStockDailyParams{
		Symbol:    "005930",
		BaseDate:  "20260505",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, "J", capturedQuery.Get("FID_COND_MRKT_DIV_CODE"))
	assert.Equal(t, "005930", capturedQuery.Get("FID_INPUT_ISCD"))
	assert.Equal(t, "20260505", capturedQuery.Get("FID_INPUT_DATE_1"))

	// output1 (요약)
	assert.Equal(t, decimal.NewFromInt(75800), res.Output1.StckPrpr)
	assert.Equal(t, "KOSPI200", res.Output1.RprsMrktKorName)
	assert.Equal(t, int64(12345678), res.Output1.AcmlVol)

	// output2 (Array, 일별 거래)
	require.Len(t, res.Output2, 2)
	assert.Equal(t, "20260505", res.Output2[0].StckBsopDate)
	assert.Equal(t, decimal.NewFromInt(75800), res.Output2[0].StckClpr)
	assert.Equal(t, int64(12345678), res.Output2[0].AcmlVol)
	assert.Equal(t, int64(-123456), res.Output2[0].FrgnNtbyQty)
	assert.Equal(t, int64(234567), res.Output2[0].PrsnNtbyQty)
	assert.Equal(t, int64(-100000), res.Output2[0].OrgnNtbyQty)
}
```

- [ ] **Step 2: 테스트 실행 → FAIL**

Run: `go test ./domestic/... -run InquireInvestorTradeByStockDaily -v`
Expected: 컴파일 실패 (`InquireInvestorTradeByStockDaily`, struct 미정의).

- [ ] **Step 3: 구현** — `domestic/investor.go`

KIS docs (`docs/api/국내주식/종목별_투자자매매동향(일별).md`) 의 line 78~180+ 응답 필드 모두 PascalCase 1:1 매핑. **모든 필드** 포함 (한 응답 ~95 fields). 다음은 핵심 skeleton — 누락 필드는 KIS docs 의 line 78~180+ 모든 row 를 struct field 로 추가:

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

// InvestorTradeByStockDaily 는 종목별 투자자매매동향(일별) (FHPTJ04160001) 응답.
//
// 한투 docs: docs/api/국내주식/종목별_투자자매매동향(일별).md
// path: /uapi/domestic-stock/v1/quotations/investor-trade-by-stock-daily
//
// output1 (요약) + output2 (일별 Array). 각 일별 행에 외국인/개인/기관 등
// 13개 투자자 type 의 매수/매도/순매수 수량 + 거래대금 (~95 필드).
type InvestorTradeByStockDaily struct {
	Output1 InvestorTradeByStockDailySummary `json:"output1"`
	Output2 []InvestorTradeByStockDailyItem  `json:"output2"`
}

// InvestorTradeByStockDailySummary 는 응답의 output1 (단일 객체, 요약).
type InvestorTradeByStockDailySummary struct {
	StckPrpr        decimal.Decimal `json:"stck_prpr"`         // 주식 현재가
	PrdyVrss        decimal.Decimal `json:"prdy_vrss"`         // 전일 대비
	PrdyVrssSign    string          `json:"prdy_vrss_sign"`    // 전일 대비 부호
	PrdyCtrt        float64         `json:"prdy_ctrt,string"`  // 전일 대비율
	AcmlVol         int64           `json:"acml_vol,string"`   // 누적 거래량
	PrdyVol         int64           `json:"prdy_vol,string"`   // 전일 거래량
	RprsMrktKorName string          `json:"rprs_mrkt_kor_name"` // 대표 시장 한글명
}

// InvestorTradeByStockDailyItem 은 응답의 output2 한 행 (한 일자).
//
// KIS docs 의 line 86~180+ 모든 필드 1:1 매핑. 13 투자자 type
// (외국인/개인/기관계/증권/투자신탁/사모펀드/은행/보험/종금/기금/기타/기타법인/기타단체)
// 각각 ntby_qty + seln_vol + shnu_vol + seln_tr_pbmn + shnu_tr_pbmn + ntby_tr_pbmn = 6 fields.
type InvestorTradeByStockDailyItem struct {
	// 일자 + 시세 (10 fields)
	StckBsopDate string          `json:"stck_bsop_date"`     // 주식 영업 일자
	StckClpr     decimal.Decimal `json:"stck_clpr"`          // 주식 종가
	PrdyVrss     decimal.Decimal `json:"prdy_vrss"`          // 전일 대비
	PrdyVrssSign string          `json:"prdy_vrss_sign"`     // 전일 대비 부호
	PrdyCtrt     float64         `json:"prdy_ctrt,string"`   // 전일 대비율
	AcmlVol      int64           `json:"acml_vol,string"`    // 누적 거래량 (주)
	AcmlTrPbmn   int64           `json:"acml_tr_pbmn,string"` // 누적 거래 대금 (백만원)
	StckOprc     decimal.Decimal `json:"stck_oprc"`          // 시가
	StckHgpr     decimal.Decimal `json:"stck_hgpr"`          // 최고가
	StckLwpr     decimal.Decimal `json:"stck_lwpr"`          // 최저가

	// 외국인 (10 fields)
	FrgnNtbyQty       int64 `json:"frgn_ntby_qty,string"`       // 외국인 순매수 수량
	FrgnRegNtbyQty    int64 `json:"frgn_reg_ntby_qty,string"`   // 외국인 등록 순매수 수량
	FrgnNregNtbyQty   int64 `json:"frgn_nreg_ntby_qty,string"`  // 외국인 비등록 순매수 수량
	FrgnRegNtbyPbmn   int64 `json:"frgn_reg_ntby_pbmn,string"`  // 외국인 등록 순매수 대금
	FrgnNtbyTrPbmn    int64 `json:"frgn_ntby_tr_pbmn,string"`   // 외국인 순매수 거래 대금
	FrgnNregNtbyPbmn  int64 `json:"frgn_nreg_ntby_pbmn,string"` // 외국인 비등록 순매수 대금
	FrgnSelnVol       int64 `json:"frgn_seln_vol,string"`       // 외국인 매도 거래량
	FrgnShnuVol       int64 `json:"frgn_shnu_vol,string"`       // 외국인 매수 거래량
	FrgnSelnTrPbmn    int64 `json:"frgn_seln_tr_pbmn,string"`   // 외국인 매도 거래 대금
	FrgnShnuTrPbmn    int64 `json:"frgn_shnu_tr_pbmn,string"`   // 외국인 매수 거래 대금

	// 외국인 등록/비등록 매도매수 (8 fields)
	FrgnRegAskpQty   int64 `json:"frgn_reg_askp_qty,string"`   // 외국인 등록 매도 수량
	FrgnRegBidpQty   int64 `json:"frgn_reg_bidp_qty,string"`   // 외국인 등록 매수 수량
	FrgnRegAskpPbmn  int64 `json:"frgn_reg_askp_pbmn,string"`  // 외국인 등록 매도 대금
	FrgnRegBidpPbmn  int64 `json:"frgn_reg_bidp_pbmn,string"`  // 외국인 등록 매수 대금
	FrgnNregAskpQty  int64 `json:"frgn_nreg_askp_qty,string"`  // 외국인 비등록 매도 수량
	FrgnNregBidpQty  int64 `json:"frgn_nreg_bidp_qty,string"`  // 외국인 비등록 매수 수량
	FrgnNregAskpPbmn int64 `json:"frgn_nreg_askp_pbmn,string"` // 외국인 비등록 매도 대금
	FrgnNregBidpPbmn int64 `json:"frgn_nreg_bidp_pbmn,string"` // 외국인 비등록 매수 대금

	// 개인 (5 fields)
	PrsnNtbyQty    int64 `json:"prsn_ntby_qty,string"`     // 개인 순매수 수량
	PrsnNtbyTrPbmn int64 `json:"prsn_ntby_tr_pbmn,string"` // 개인 순매수 거래 대금
	PrsnSelnVol    int64 `json:"prsn_seln_vol,string"`     // 개인 매도 거래량
	PrsnShnuVol    int64 `json:"prsn_shnu_vol,string"`     // 개인 매수 거래량
	PrsnSelnTrPbmn int64 `json:"prsn_seln_tr_pbmn,string"` // 개인 매도 거래 대금
	PrsnShnuTrPbmn int64 `json:"prsn_shnu_tr_pbmn,string"` // 개인 매수 거래 대금

	// 기관계 (6 fields)
	OrgnNtbyQty    int64 `json:"orgn_ntby_qty,string"`     // 기관계 순매수 수량
	OrgnNtbyTrPbmn int64 `json:"orgn_ntby_tr_pbmn,string"` // 기관계 순매수 거래 대금
	OrgnSelnVol    int64 `json:"orgn_seln_vol,string"`     // 기관계 매도 거래량
	OrgnShnuVol    int64 `json:"orgn_shnu_vol,string"`     // 기관계 매수 거래량
	OrgnSelnTrPbmn int64 `json:"orgn_seln_tr_pbmn,string"` // 기관계 매도 거래 대금
	OrgnShnuTrPbmn int64 `json:"orgn_shnu_tr_pbmn,string"` // 기관계 매수 거래 대금

	// 증권 (6 fields)
	ScrtNtbyQty    int64 `json:"scrt_ntby_qty,string"`
	ScrtNtbyTrPbmn int64 `json:"scrt_ntby_tr_pbmn,string"`
	ScrtSelnVol    int64 `json:"scrt_seln_vol,string"`
	ScrtShnuVol    int64 `json:"scrt_shnu_vol,string"`
	ScrtSelnTrPbmn int64 `json:"scrt_seln_tr_pbmn,string"`
	ScrtShnuTrPbmn int64 `json:"scrt_shnu_tr_pbmn,string"`

	// 투자신탁 (6 fields)
	IvtrNtbyQty    int64 `json:"ivtr_ntby_qty,string"`
	IvtrNtbyTrPbmn int64 `json:"ivtr_ntby_tr_pbmn,string"`
	IvtrSelnVol    int64 `json:"ivtr_seln_vol,string"`
	IvtrShnuVol    int64 `json:"ivtr_shnu_vol,string"`
	IvtrSelnTrPbmn int64 `json:"ivtr_seln_tr_pbmn,string"`
	IvtrShnuTrPbmn int64 `json:"ivtr_shnu_tr_pbmn,string"`

	// 사모펀드 (6 fields, KIS docs 가 vol/qty 혼용 — vol 사용)
	PeFundNtbyVol    int64 `json:"pe_fund_ntby_vol,string"`
	PeFundNtbyTrPbmn int64 `json:"pe_fund_ntby_tr_pbmn,string"`
	PeFundSelnVol    int64 `json:"pe_fund_seln_vol,string"`
	PeFundShnuVol    int64 `json:"pe_fund_shnu_vol,string"`
	PeFundSelnTrPbmn int64 `json:"pe_fund_seln_tr_pbmn,string"`
	PeFundShnuTrPbmn int64 `json:"pe_fund_shnu_tr_pbmn,string"`

	// 은행 (6 fields)
	BankNtbyQty    int64 `json:"bank_ntby_qty,string"`
	BankNtbyTrPbmn int64 `json:"bank_ntby_tr_pbmn,string"`
	BankSelnVol    int64 `json:"bank_seln_vol,string"`
	BankShnuVol    int64 `json:"bank_shnu_vol,string"`
	BankSelnTrPbmn int64 `json:"bank_seln_tr_pbmn,string"`
	BankShnuTrPbmn int64 `json:"bank_shnu_tr_pbmn,string"`

	// 보험 (6 fields)
	InsuNtbyQty    int64 `json:"insu_ntby_qty,string"`
	InsuNtbyTrPbmn int64 `json:"insu_ntby_tr_pbmn,string"`
	InsuSelnVol    int64 `json:"insu_seln_vol,string"`
	InsuShnuVol    int64 `json:"insu_shnu_vol,string"`
	InsuSelnTrPbmn int64 `json:"insu_seln_tr_pbmn,string"`
	InsuShnuTrPbmn int64 `json:"insu_shnu_tr_pbmn,string"`

	// 종금 (6 fields)
	MrbnNtbyQty    int64 `json:"mrbn_ntby_qty,string"`
	MrbnNtbyTrPbmn int64 `json:"mrbn_ntby_tr_pbmn,string"`
	MrbnSelnVol    int64 `json:"mrbn_seln_vol,string"`
	MrbnShnuVol    int64 `json:"mrbn_shnu_vol,string"`
	MrbnSelnTrPbmn int64 `json:"mrbn_seln_tr_pbmn,string"`
	MrbnShnuTrPbmn int64 `json:"mrbn_shnu_tr_pbmn,string"`

	// 기금 (6 fields)
	FundNtbyQty    int64 `json:"fund_ntby_qty,string"`
	FundNtbyTrPbmn int64 `json:"fund_ntby_tr_pbmn,string"`
	FundSelnVol    int64 `json:"fund_seln_vol,string"`
	FundShnuVol    int64 `json:"fund_shnu_vol,string"`
	FundSelnTrPbmn int64 `json:"fund_seln_tr_pbmn,string"`
	FundShnuTrPbmn int64 `json:"fund_shnu_tr_pbmn,string"`

	// 기타 (6 fields)
	EtcNtbyQty    int64 `json:"etc_ntby_qty,string"`
	EtcNtbyTrPbmn int64 `json:"etc_ntby_tr_pbmn,string"`
	EtcSelnVol    int64 `json:"etc_seln_vol,string"`
	EtcShnuVol    int64 `json:"etc_shnu_vol,string"`
	EtcSelnTrPbmn int64 `json:"etc_seln_tr_pbmn,string"`
	EtcShnuTrPbmn int64 `json:"etc_shnu_tr_pbmn,string"`

	// 기타 법인 (3 fields, KIS docs 가 vol 사용)
	EtcCorpNtbyVol    int64 `json:"etc_corp_ntby_vol,string"`
	EtcCorpNtbyTrPbmn int64 `json:"etc_corp_ntby_tr_pbmn,string"`

	// 기타 단체 (6 fields, vol)
	EtcOrgtNtbyVol    int64 `json:"etc_orgt_ntby_vol,string"`
	EtcOrgtNtbyTrPbmn int64 `json:"etc_orgt_ntby_tr_pbmn,string"`
	EtcOrgtSelnVol    int64 `json:"etc_orgt_seln_vol,string"`
	EtcOrgtShnuVol    int64 `json:"etc_orgt_shnu_vol,string"`
	EtcOrgtSelnTrPbmn int64 `json:"etc_orgt_seln_tr_pbmn,string"`
	EtcOrgtShnuTrPbmn int64 `json:"etc_orgt_shnu_tr_pbmn,string"`
}

// InquireInvestorTradeByStockDailyParams 는 종목별 투자자매매동향(일별) 조회 파라미터.
type InquireInvestorTradeByStockDailyParams struct {
	MarketCode string // FID_COND_MRKT_DIV_CODE — "J":KRX, "NX":NXT, "UN":통합. 빈 값=>"J"
	Symbol     string // FID_INPUT_ISCD — 필수, 종목코드 (6자리)
	BaseDate   string // FID_INPUT_DATE_1 — 필수, YYYYMMDD (해당일 조회는 장 종료 후 가능)
	OrgAdjPrc  string // FID_ORG_ADJ_PRC — 빈 값(공란) default
	EtcClsCode string // FID_ETC_CLS_CODE — 빈 값(공란) default
}

// InquireInvestorTradeByStockDaily 는 종목별 투자자매매동향(일별) 호출.
//
// 한투 docs: docs/api/국내주식/종목별_투자자매매동향(일별).md
// path: /uapi/domestic-stock/v1/quotations/investor-trade-by-stock-daily (FHPTJ04160001)
func (c *Client) InquireInvestorTradeByStockDaily(ctx context.Context, params InquireInvestorTradeByStockDailyParams) (*InvestorTradeByStockDaily, error) {
	market := params.MarketCode
	if market == "" {
		market = "J"
	}

	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/domestic-stock/v1/quotations/investor-trade-by-stock-daily",
		TrID:   "FHPTJ04160001",
		Query: map[string]string{
			"FID_COND_MRKT_DIV_CODE": market,
			"FID_INPUT_ISCD":         params.Symbol,
			"FID_INPUT_DATE_1":       params.BaseDate,
			"FID_ORG_ADJ_PRC":        params.OrgAdjPrc,
			"FID_ETC_CLS_CODE":       params.EtcClsCode,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}

	var res InvestorTradeByStockDaily
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse InvestorTradeByStockDaily: %w", err)
	}
	return &res, nil
}
```

- [ ] **Step 4: 테스트 실행 → PASS**

Run: `go test ./domestic/... -run InquireInvestorTradeByStockDaily -v`
Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add domestic/investor.go domestic/investor_test.go
git commit -m "$(cat <<'EOF'
[feat] domestic — InquireInvestorTradeByStockDaily (종목별 투자자매매동향 일별, FHPTJ04160001)

InvestorTradeByStockDaily (Output1 요약 7 + Output2 일별 array 95 필드) +
InquireInvestorTradeByStockDailyParams (5 query). 13 투자자 type 의
ntby_qty/seln_vol/shnu_vol/*_tr_pbmn 모두 PascalCase 1:1 매핑.

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 3: domestic/investor.go — InquireInvestorDailyByMarket

**Files:**
- Modify: `domestic/investor.go` (append)
- Modify: `domestic/investor_test.go` (append)

> 시장별 투자자매매동향(일별). output (Array) 한 행에 ~40 fields. KIS docs (`docs/api/국내주식/시장별_투자자매매동향(일별).md`) line 75~113 응답.

- [ ] **Step 1: 테스트 추가**

```go
func TestClient_InquireInvestorDailyByMarket(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/quotations/inquire-investor-daily-by-market`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "investor_daily_by_market_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireInvestorDailyByMarket(context.Background(), domestic.InquireInvestorDailyByMarketParams{
		Symbol:     "0001",   // 코스피 종합
		BaseDate:   "20260505",
		Market:     "KSP",
		BaseDate2:  "20260505",
		SubCode:    "0001",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, "0001", capturedQuery.Get("FID_INPUT_ISCD"))
	assert.Equal(t, "20260505", capturedQuery.Get("FID_INPUT_DATE_1"))
	assert.Equal(t, "KSP", capturedQuery.Get("FID_INPUT_ISCD_1"))

	require.Len(t, res.Output, 1)
	assert.Equal(t, "20260505", res.Output[0].StckBsopDate)
	assert.Equal(t, "5", res.Output[0].PrdyVrssSign)
	assert.Equal(t, int64(-123456), res.Output[0].FrgnNtbyQty)
	assert.Equal(t, int64(234567), res.Output[0].PrsnNtbyQty)
}
```

- [ ] **Step 2: FAIL**

- [ ] **Step 3: 구현 추가** — `domestic/investor.go` 끝에

```go
// InvestorDailyByMarket 은 시장별 투자자매매동향(일별) (FHPTJ04040000) 응답.
//
// 한투 docs: docs/api/국내주식/시장별_투자자매매동향(일별).md
// path: /uapi/domestic-stock/v1/quotations/inquire-investor-daily-by-market
type InvestorDailyByMarket struct {
	Output []InvestorDailyByMarketItem `json:"output"`
}

// InvestorDailyByMarketItem 은 응답 한 행 (한 일자).
type InvestorDailyByMarketItem struct {
	// 일자 + 지수 (9 fields)
	StckBsopDate       string          `json:"stck_bsop_date"`              // 주식 영업 일자
	BstpNmixPrpr       decimal.Decimal `json:"bstp_nmix_prpr"`              // 업종 지수 현재가
	BstpNmixPrdyVrss   decimal.Decimal `json:"bstp_nmix_prdy_vrss"`         // 업종 지수 전일 대비
	PrdyVrssSign       string          `json:"prdy_vrss_sign"`              // 전일 대비 부호
	BstpNmixPrdyCtrt   float64         `json:"bstp_nmix_prdy_ctrt,string"`  // 업종 지수 전일 대비율
	BstpNmixOprc       decimal.Decimal `json:"bstp_nmix_oprc"`              // 업종 지수 시가
	BstpNmixHgpr       decimal.Decimal `json:"bstp_nmix_hgpr"`              // 업종 지수 최고가
	BstpNmixLwpr       decimal.Decimal `json:"bstp_nmix_lwpr"`              // 업종 지수 최저가
	StckPrdyClpr       decimal.Decimal `json:"stck_prdy_clpr"`              // 전일 종가

	// 13 type ntby_qty (수량)
	FrgnNtbyQty     int64 `json:"frgn_ntby_qty,string"`     // 외국인 순매수 수량
	FrgnRegNtbyQty  int64 `json:"frgn_reg_ntby_qty,string"` // 외국인 등록 순매수 수량
	FrgnNregNtbyQty int64 `json:"frgn_nreg_ntby_qty,string"` // 외국인 비등록 순매수 수량
	PrsnNtbyQty     int64 `json:"prsn_ntby_qty,string"`     // 개인 순매수 수량
	OrgnNtbyQty     int64 `json:"orgn_ntby_qty,string"`     // 기관계 순매수 수량
	ScrtNtbyQty     int64 `json:"scrt_ntby_qty,string"`     // 증권 순매수 수량
	IvtrNtbyQty     int64 `json:"ivtr_ntby_qty,string"`     // 투자신탁 순매수 수량
	PeFundNtbyVol   int64 `json:"pe_fund_ntby_vol,string"`  // 사모 펀드 순매수 거래량 (vol)
	BankNtbyQty     int64 `json:"bank_ntby_qty,string"`     // 은행 순매수 수량
	InsuNtbyQty     int64 `json:"insu_ntby_qty,string"`     // 보험 순매수 수량
	MrbnNtbyQty     int64 `json:"mrbn_ntby_qty,string"`     // 종금 순매수 수량
	FundNtbyQty     int64 `json:"fund_ntby_qty,string"`     // 기금 순매수 수량
	EtcNtbyQty      int64 `json:"etc_ntby_qty,string"`      // 기타 순매수 수량
	EtcOrgtNtbyVol  int64 `json:"etc_orgt_ntby_vol,string"` // 기타 단체 순매수 거래량
	EtcCorpNtbyVol  int64 `json:"etc_corp_ntby_vol,string"` // 기타 법인 순매수 거래량

	// 14 type ntby_tr_pbmn (거래대금)
	FrgnNtbyTrPbmn    int64 `json:"frgn_ntby_tr_pbmn,string"`     // 외국인 순매수 거래 대금
	FrgnRegNtbyPbmn   int64 `json:"frgn_reg_ntby_pbmn,string"`    // 외국인 등록 순매수 대금
	FrgnNregNtbyPbmn  int64 `json:"frgn_nreg_ntby_pbmn,string"`   // 외국인 비등록 순매수 대금
	PrsnNtbyTrPbmn    int64 `json:"prsn_ntby_tr_pbmn,string"`     // 개인 순매수 거래 대금
	OrgnNtbyTrPbmn    int64 `json:"orgn_ntby_tr_pbmn,string"`     // 기관계 순매수 거래 대금
	ScrtNtbyTrPbmn    int64 `json:"scrt_ntby_tr_pbmn,string"`     // 증권
	IvtrNtbyTrPbmn    int64 `json:"ivtr_ntby_tr_pbmn,string"`     // 투자신탁
	PeFundNtbyTrPbmn  int64 `json:"pe_fund_ntby_tr_pbmn,string"`  // 사모 펀드
	BankNtbyTrPbmn    int64 `json:"bank_ntby_tr_pbmn,string"`     // 은행
	InsuNtbyTrPbmn    int64 `json:"insu_ntby_tr_pbmn,string"`     // 보험
	MrbnNtbyTrPbmn    int64 `json:"mrbn_ntby_tr_pbmn,string"`     // 종금
	FundNtbyTrPbmn    int64 `json:"fund_ntby_tr_pbmn,string"`     // 기금
	EtcNtbyTrPbmn     int64 `json:"etc_ntby_tr_pbmn,string"`      // 기타
	EtcOrgtNtbyTrPbmn int64 `json:"etc_orgt_ntby_tr_pbmn,string"` // 기타 단체
	EtcCorpNtbyTrPbmn int64 `json:"etc_corp_ntby_tr_pbmn,string"` // 기타 법인
}

// InquireInvestorDailyByMarketParams 는 시장별 투자자매매동향(일별) 조회 파라미터.
//
// KIS docs 의 query 키 그대로 노출. FID_INPUT_ISCD = 업종분류코드, FID_INPUT_ISCD_1 = 시장 (KSP/KSQ), FID_INPUT_ISCD_2 = 하위 분류.
type InquireInvestorDailyByMarketParams struct {
	MarketCode string // FID_COND_MRKT_DIV_CODE — 빈 값=>"J"
	Symbol     string // FID_INPUT_ISCD — 업종분류코드 (예 "0001":코스피 종합)
	BaseDate   string // FID_INPUT_DATE_1 — YYYYMMDD
	Market     string // FID_INPUT_ISCD_1 — "KSP"(코스피) 또는 "KSQ"(코스닥)
	BaseDate2  string // FID_INPUT_DATE_2 — BaseDate 와 동일 일자
	SubCode    string // FID_INPUT_ISCD_2 — 하위 분류코드 (업종분류코드)
}

// InquireInvestorDailyByMarket 은 시장별 투자자매매동향(일별) 호출.
//
// 한투 docs: docs/api/국내주식/시장별_투자자매매동향(일별).md
// path: /uapi/domestic-stock/v1/quotations/inquire-investor-daily-by-market (FHPTJ04040000)
func (c *Client) InquireInvestorDailyByMarket(ctx context.Context, params InquireInvestorDailyByMarketParams) (*InvestorDailyByMarket, error) {
	market := params.MarketCode
	if market == "" {
		market = "J"
	}

	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/domestic-stock/v1/quotations/inquire-investor-daily-by-market",
		TrID:   "FHPTJ04040000",
		Query: map[string]string{
			"FID_COND_MRKT_DIV_CODE": market,
			"FID_INPUT_ISCD":         params.Symbol,
			"FID_INPUT_DATE_1":       params.BaseDate,
			"FID_INPUT_ISCD_1":       params.Market,
			"FID_INPUT_DATE_2":       params.BaseDate2,
			"FID_INPUT_ISCD_2":       params.SubCode,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}

	var res InvestorDailyByMarket
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse InvestorDailyByMarket: %w", err)
	}
	return &res, nil
}
```

- [ ] **Step 4: PASS**

- [ ] **Step 5: Commit**

```bash
git add domestic/investor.go domestic/investor_test.go
git commit -m "$(cat <<'EOF'
[feat] domestic — InquireInvestorDailyByMarket (시장별 투자자매매동향 일별, FHPTJ04040000)

InvestorDailyByMarket + InvestorDailyByMarketItem (40+ 필드) +
InquireInvestorDailyByMarketParams (6 query). KSP/KSQ 시장 구분 + 업종분류코드.

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 4: domestic/investor.go — InquireInvestorTimeByMarket

**Files:**
- Modify: `domestic/investor.go` (append)
- Modify: `domestic/investor_test.go` (append)

> 시장별 투자자매매동향(시세) — 다양한 product (코스피/코스닥/선물/옵션/ETF/ELW/ETN 등) 지원. output (단일 객체, ~108 fields). KIS docs (`docs/api/국내주식/시장별_투자자매매동향(시세).md`) line 113~180+ 응답.

- [ ] **Step 1: 테스트 추가**

```go
func TestClient_InquireInvestorTimeByMarket(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/quotations/inquire-investor-time-by-market`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "investor_time_by_market_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireInvestorTimeByMarket(context.Background(), domestic.InquireInvestorTimeByMarketParams{
		Market:  "KSP",
		SubCode: "0001",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, "KSP", capturedQuery.Get("fid_input_iscd"))
	assert.Equal(t, "0001", capturedQuery.Get("fid_input_iscd_2"))

	assert.Equal(t, int64(5000000), res.Output.FrgnSelnVol)
	assert.Equal(t, int64(-123456), res.Output.FrgnNtbyQty)
	assert.Equal(t, int64(234567), res.Output.PrsnNtbyQty)
}
```

- [ ] **Step 2: FAIL**

- [ ] **Step 3: 구현 추가**

```go
// InvestorTimeByMarket 은 시장별 투자자매매동향(시세) (FHPTJ04030000) 응답.
//
// 한투 docs: docs/api/국내주식/시장별_투자자매매동향(시세).md
// path: /uapi/domestic-stock/v1/quotations/inquire-investor-time-by-market
type InvestorTimeByMarket struct {
	Output InvestorTimeByMarketSnapshot `json:"output"`
}

// InvestorTimeByMarketSnapshot 은 응답의 output (단일 객체, 시세).
//
// 13 type 의 (ntby_qty/seln_vol/shnu_vol/seln_tr_pbmn/shnu_tr_pbmn/ntby_tr_pbmn) = 6 fields × 13 = 78 fields.
// 일부 type 은 vol 표기. KIS docs 의 line 113~180+ 모든 필드 1:1 매핑.
type InvestorTimeByMarketSnapshot struct {
	// 외국인 (6 fields)
	FrgnSelnVol    int64 `json:"frgn_seln_vol,string"`
	FrgnShnuVol    int64 `json:"frgn_shnu_vol,string"`
	FrgnNtbyQty    int64 `json:"frgn_ntby_qty,string"`
	FrgnSelnTrPbmn int64 `json:"frgn_seln_tr_pbmn,string"`
	FrgnShnuTrPbmn int64 `json:"frgn_shnu_tr_pbmn,string"`
	FrgnNtbyTrPbmn int64 `json:"frgn_ntby_tr_pbmn,string"`

	// 개인 (6 fields)
	PrsnSelnVol    int64 `json:"prsn_seln_vol,string"`
	PrsnShnuVol    int64 `json:"prsn_shnu_vol,string"`
	PrsnNtbyQty    int64 `json:"prsn_ntby_qty,string"`
	PrsnSelnTrPbmn int64 `json:"prsn_seln_tr_pbmn,string"`
	PrsnShnuTrPbmn int64 `json:"prsn_shnu_tr_pbmn,string"`
	PrsnNtbyTrPbmn int64 `json:"prsn_ntby_tr_pbmn,string"`

	// 기관계 (6 fields)
	OrgnSelnVol    int64 `json:"orgn_seln_vol,string"`
	OrgnShnuVol    int64 `json:"orgn_shnu_vol,string"`
	OrgnNtbyQty    int64 `json:"orgn_ntby_qty,string"`
	OrgnSelnTrPbmn int64 `json:"orgn_seln_tr_pbmn,string"`
	OrgnShnuTrPbmn int64 `json:"orgn_shnu_tr_pbmn,string"`
	OrgnNtbyTrPbmn int64 `json:"orgn_ntby_tr_pbmn,string"`

	// 증권 (6 fields)
	ScrtSelnVol    int64 `json:"scrt_seln_vol,string"`
	ScrtShnuVol    int64 `json:"scrt_shnu_vol,string"`
	ScrtNtbyQty    int64 `json:"scrt_ntby_qty,string"`
	ScrtSelnTrPbmn int64 `json:"scrt_seln_tr_pbmn,string"`
	ScrtShnuTrPbmn int64 `json:"scrt_shnu_tr_pbmn,string"`
	ScrtNtbyTrPbmn int64 `json:"scrt_ntby_tr_pbmn,string"`

	// 투자신탁 (6 fields)
	IvtrSelnVol    int64 `json:"ivtr_seln_vol,string"`
	IvtrShnuVol    int64 `json:"ivtr_shnu_vol,string"`
	IvtrNtbyQty    int64 `json:"ivtr_ntby_qty,string"`
	IvtrSelnTrPbmn int64 `json:"ivtr_seln_tr_pbmn,string"`
	IvtrShnuTrPbmn int64 `json:"ivtr_shnu_tr_pbmn,string"`
	IvtrNtbyTrPbmn int64 `json:"ivtr_ntby_tr_pbmn,string"`

	// 사모 펀드 (6 fields, vol)
	PeFundSelnVol    int64 `json:"pe_fund_seln_vol,string"`
	PeFundShnuVol    int64 `json:"pe_fund_shnu_vol,string"`
	PeFundNtbyVol    int64 `json:"pe_fund_ntby_vol,string"`
	PeFundSelnTrPbmn int64 `json:"pe_fund_seln_tr_pbmn,string"`
	PeFundShnuTrPbmn int64 `json:"pe_fund_shnu_tr_pbmn,string"`
	PeFundNtbyTrPbmn int64 `json:"pe_fund_ntby_tr_pbmn,string"`

	// 은행 (6 fields)
	BankSelnVol    int64 `json:"bank_seln_vol,string"`
	BankShnuVol    int64 `json:"bank_shnu_vol,string"`
	BankNtbyQty    int64 `json:"bank_ntby_qty,string"`
	BankSelnTrPbmn int64 `json:"bank_seln_tr_pbmn,string"`
	BankShnuTrPbmn int64 `json:"bank_shnu_tr_pbmn,string"`
	BankNtbyTrPbmn int64 `json:"bank_ntby_tr_pbmn,string"`

	// 보험 (6 fields)
	InsuSelnVol    int64 `json:"insu_seln_vol,string"`
	InsuShnuVol    int64 `json:"insu_shnu_vol,string"`
	InsuNtbyQty    int64 `json:"insu_ntby_qty,string"`
	InsuSelnTrPbmn int64 `json:"insu_seln_tr_pbmn,string"`
	InsuShnuTrPbmn int64 `json:"insu_shnu_tr_pbmn,string"`
	InsuNtbyTrPbmn int64 `json:"insu_ntby_tr_pbmn,string"`

	// 종금 (6 fields)
	MrbnSelnVol    int64 `json:"mrbn_seln_vol,string"`
	MrbnShnuVol    int64 `json:"mrbn_shnu_vol,string"`
	MrbnNtbyQty    int64 `json:"mrbn_ntby_qty,string"`
	MrbnSelnTrPbmn int64 `json:"mrbn_seln_tr_pbmn,string"`
	MrbnShnuTrPbmn int64 `json:"mrbn_shnu_tr_pbmn,string"`
	MrbnNtbyTrPbmn int64 `json:"mrbn_ntby_tr_pbmn,string"`

	// 기금 (6 fields)
	FundSelnVol    int64 `json:"fund_seln_vol,string"`
	FundShnuVol    int64 `json:"fund_shnu_vol,string"`
	FundNtbyQty    int64 `json:"fund_ntby_qty,string"`
	FundSelnTrPbmn int64 `json:"fund_seln_tr_pbmn,string"`
	FundShnuTrPbmn int64 `json:"fund_shnu_tr_pbmn,string"`
	FundNtbyTrPbmn int64 `json:"fund_ntby_tr_pbmn,string"`

	// 기타 단체 (6 fields, vol)
	EtcOrgtSelnVol    int64 `json:"etc_orgt_seln_vol,string"`
	EtcOrgtShnuVol    int64 `json:"etc_orgt_shnu_vol,string"`
	EtcOrgtNtbyVol    int64 `json:"etc_orgt_ntby_vol,string"`
	EtcOrgtSelnTrPbmn int64 `json:"etc_orgt_seln_tr_pbmn,string"`
	EtcOrgtShnuTrPbmn int64 `json:"etc_orgt_shnu_tr_pbmn,string"`
	EtcOrgtNtbyTrPbmn int64 `json:"etc_orgt_ntby_tr_pbmn,string"`
}

// InquireInvestorTimeByMarketParams 는 시장별 투자자매매동향(시세) 조회 파라미터.
//
// fid_input_iscd 는 product 코드 (KSP/KSQ/K2I/999/ETF/ELW/ETN/MKI/WKM/WKI/KQI). fid_input_iscd_2 는 product 별 하위 분류.
type InquireInvestorTimeByMarketParams struct {
	Market  string // fid_input_iscd — KSP(코스피)/KSQ(코스닥)/K2I(선물옵션)/등
	SubCode string // fid_input_iscd_2 — Market 에 따른 하위 분류 (코스피 0001 종합 등)
}

// InquireInvestorTimeByMarket 은 시장별 투자자매매동향(시세) 호출.
//
// 한투 docs: docs/api/국내주식/시장별_투자자매매동향(시세).md
// path: /uapi/domestic-stock/v1/quotations/inquire-investor-time-by-market (FHPTJ04030000)
//
// 다양한 product (KSP/KSQ/선물옵션/ETF/ELW/ETN/등) 지원. SubCode 의 의미는 Market 에 따라 달라짐 — KIS docs 의 line 50~93 참조.
func (c *Client) InquireInvestorTimeByMarket(ctx context.Context, params InquireInvestorTimeByMarketParams) (*InvestorTimeByMarket, error) {
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/domestic-stock/v1/quotations/inquire-investor-time-by-market",
		TrID:   "FHPTJ04030000",
		Query: map[string]string{
			"fid_input_iscd":   params.Market,
			"fid_input_iscd_2": params.SubCode,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}

	var res InvestorTimeByMarket
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse InvestorTimeByMarket: %w", err)
	}
	return &res, nil
}
```

- [ ] **Step 4: PASS**

- [ ] **Step 5: Commit**

```bash
git add domestic/investor.go domestic/investor_test.go
git commit -m "$(cat <<'EOF'
[feat] domestic — InquireInvestorTimeByMarket (시장별 투자자매매동향 시세, FHPTJ04030000)

InvestorTimeByMarket + InvestorTimeByMarketSnapshot (78 필드 — 13 투자자
type × 6 fields) + InquireInvestorTimeByMarketParams (2 query, 소문자 fid_*).
KSP/KSQ/선물/ETF/ELW/ETN 등 다양한 product 지원.

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 5: domestic/industry.go — InquireIndexPrice + 공통 base

**Files:**
- Create: `domestic/industry.go`
- Create: `domestic/industry_test.go`

> 국내업종 현재지수. output (단일 객체, ~36 fields). KIS docs (`docs/api/국내주식/국내업종_현재지수.md`).

- [ ] **Step 1: 테스트 작성**

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

func TestClient_InquireIndexPrice(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/quotations/inquire-index-price`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "index_price_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireIndexPrice(context.Background(), domestic.InquireIndexPriceParams{
		Symbol: "0001", // 코스피
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, "U", capturedQuery.Get("FID_COND_MRKT_DIV_CODE"))
	assert.Equal(t, "0001", capturedQuery.Get("FID_INPUT_ISCD"))

	d, _ := decimal.NewFromString("2650.45")
	assert.True(t, d.Equal(res.Output.BstpNmixPrpr))
	assert.InDelta(t, -0.46, res.Output.BstpNmixPrdyCtrt, 0.001)
	assert.Equal(t, int64(350000000), res.Output.AcmlVol)
	assert.Equal(t, "315", res.Output.AscnIssuCnt)
	assert.Equal(t, "450", res.Output.DownIssuCnt)
}
```

- [ ] **Step 2: FAIL**

- [ ] **Step 3: 구현** — `domestic/industry.go`

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

// IndexPrice 는 국내업종 현재지수 (FHPUP02100000) 응답.
//
// 한투 docs: docs/api/국내주식/국내업종_현재지수.md
// path: /uapi/domestic-stock/v1/quotations/inquire-index-price
type IndexPrice struct {
	Output IndexPriceSnapshot `json:"output"`
}

// IndexPriceSnapshot 은 응답의 output (단일 객체).
//
// KIS docs 의 line 73~108 모든 필드 1:1 매핑 (~36 fields).
type IndexPriceSnapshot struct {
	// 지수 + 변동
	BstpNmixPrpr     decimal.Decimal `json:"bstp_nmix_prpr"`              // 업종 지수 현재가
	BstpNmixPrdyVrss decimal.Decimal `json:"bstp_nmix_prdy_vrss"`         // 업종 지수 전일 대비
	PrdyVrssSign     string          `json:"prdy_vrss_sign"`              // 전일 대비 부호
	BstpNmixPrdyCtrt float64         `json:"bstp_nmix_prdy_ctrt,string"`  // 업종 지수 전일 대비율

	// 거래량/거래대금
	AcmlVol    int64 `json:"acml_vol,string"`     // 누적 거래량
	PrdyVol    int64 `json:"prdy_vol,string"`     // 전일 거래량
	AcmlTrPbmn int64 `json:"acml_tr_pbmn,string"` // 누적 거래 대금
	PrdyTrPbmn int64 `json:"prdy_tr_pbmn,string"` // 전일 거래 대금

	// 시가 + 시가대비
	BstpNmixOprc            decimal.Decimal `json:"bstp_nmix_oprc"`                       // 업종 지수 시가
	PrdyNmixVrssNmixOprc    decimal.Decimal `json:"prdy_nmix_vrss_nmix_oprc"`             // 전일 지수 대비 지수 시가
	OprcVrssPrprSign        string          `json:"oprc_vrss_prpr_sign"`                  // 시가 대비 현재가 부호
	BstpNmixOprcPrdyCtrt    float64         `json:"bstp_nmix_oprc_prdy_ctrt,string"`      // 업종 지수 시가 전일 대비율

	// 최고가
	BstpNmixHgpr            decimal.Decimal `json:"bstp_nmix_hgpr"`                       // 업종 지수 최고가
	PrdyNmixVrssNmixHgpr    decimal.Decimal `json:"prdy_nmix_vrss_nmix_hgpr"`             // 전일 지수 대비 지수 최고가
	HgprVrssPrprSign        string          `json:"hgpr_vrss_prpr_sign"`                  // 최고가 대비 현재가 부호
	BstpNmixHgprPrdyCtrt    float64         `json:"bstp_nmix_hgpr_prdy_ctrt,string"`      // 업종 지수 최고가 전일 대비율

	// 최저가
	BstpNmixLwpr         decimal.Decimal `json:"bstp_nmix_lwpr"`                  // 업종 지수 최저가
	PrdyClprVrssLwpr     decimal.Decimal `json:"prdy_clpr_vrss_lwpr"`             // 전일 종가 대비 최저가
	LwprVrssPrprSign     string          `json:"lwpr_vrss_prpr_sign"`             // 최저가 대비 현재가 부호
	PrdyClprVrssLwprRate float64         `json:"prdy_clpr_vrss_lwpr_rate,string"` // 전일 종가 대비 최저가 비율

	// 종목수 (5 fields, KIS 가 string 으로 줌 — 작은 정수)
	AscnIssuCnt string `json:"ascn_issu_cnt"` // 상승 종목 수
	UplmIssuCnt string `json:"uplm_issu_cnt"` // 상한 종목 수
	StnrIssuCnt string `json:"stnr_issu_cnt"` // 보합 종목 수
	DownIssuCnt string `json:"down_issu_cnt"` // 하락 종목 수
	LslmIssuCnt string `json:"lslm_issu_cnt"` // 하한 종목 수

	// 연중 (6 fields)
	DryyBstpNmixHgpr      decimal.Decimal `json:"dryy_bstp_nmix_hgpr"`             // 연중 업종 지수 최고가
	DryyHgprVrssPrprRate  float64         `json:"dryy_hgpr_vrss_prpr_rate,string"` // 연중 최고가 대비 현재가 비율
	DryyBstpNmixHgprDate  string          `json:"dryy_bstp_nmix_hgpr_date"`        // 연중 업종 지수 최고가 일자
	DryyBstpNmixLwpr      decimal.Decimal `json:"dryy_bstp_nmix_lwpr"`             // 연중 업종 지수 최저가
	DryyLwprVrssPrprRate  float64         `json:"dryy_lwpr_vrss_prpr_rate,string"` // 연중 최저가 대비 현재가 비율
	DryyBstpNmixLwprDate  string          `json:"dryy_bstp_nmix_lwpr_date"`        // 연중 업종 지수 최저가 일자

	// 호가 잔량 (5 fields)
	TotalAskpRsqn int64   `json:"total_askp_rsqn,string"` // 총 매도호가 잔량
	TotalBidpRsqn int64   `json:"total_bidp_rsqn,string"` // 총 매수호가 잔량
	SelnRsqnRate  float64 `json:"seln_rsqn_rate,string"`  // 매도 잔량 비율
	ShnuRsqnRate  float64 `json:"shnu_rsqn_rate,string"`  // 매수 잔량 비율
	NtbyRsqn      int64   `json:"ntby_rsqn,string"`       // 순매수 잔량
}

// InquireIndexPriceParams 는 국내업종 현재지수 조회 파라미터.
type InquireIndexPriceParams struct {
	MarketCode string // FID_COND_MRKT_DIV_CODE — 빈 값=>"U" (업종)
	Symbol     string // FID_INPUT_ISCD — 필수, 업종 코드 (예 "0001":코스피, "1001":코스닥, "2001":코스피200)
}

// InquireIndexPrice 는 국내업종 현재지수 호출.
//
// 한투 docs: docs/api/국내주식/국내업종_현재지수.md
// path: /uapi/domestic-stock/v1/quotations/inquire-index-price (FHPUP02100000)
func (c *Client) InquireIndexPrice(ctx context.Context, params InquireIndexPriceParams) (*IndexPrice, error) {
	market := params.MarketCode
	if market == "" {
		market = "U"
	}

	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/domestic-stock/v1/quotations/inquire-index-price",
		TrID:   "FHPUP02100000",
		Query: map[string]string{
			"FID_COND_MRKT_DIV_CODE": market,
			"FID_INPUT_ISCD":         params.Symbol,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}

	var res IndexPrice
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse IndexPrice: %w", err)
	}
	return &res, nil
}
```

- [ ] **Step 4: PASS**

- [ ] **Step 5: Commit**

```bash
git add domestic/industry.go domestic/industry_test.go
git commit -m "$(cat <<'EOF'
[feat] domestic — InquireIndexPrice (국내업종 현재지수, FHPUP02100000)

IndexPrice + IndexPriceSnapshot (~36 필드) + InquireIndexPriceParams.
지수/시가/최고/최저 + 종목수 (상승/하한/보합/하락) + 연중 + 호가 잔량.
MarketCode default "U" (업종).

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 6: domestic/industry.go — InquireIndexCategoryPrice

**Files:**
- Modify: `domestic/industry.go` (append)
- Modify: `domestic/industry_test.go` (append)

> 국내업종 구분별 전체시세. output1 (단일, 20 fields) + output2 (Array, 10 fields per row). KIS docs (`docs/api/국내주식/국내업종_구분별전체시세.md`).

- [ ] **Step 1: 테스트 추가**

```go
func TestClient_InquireIndexCategoryPrice(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/quotations/inquire-index-category-price`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "index_category_price_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireIndexCategoryPrice(context.Background(), domestic.InquireIndexCategoryPriceParams{
		Symbol:    "0001",
		MarketCls: "K",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, "U", capturedQuery.Get("FID_COND_MRKT_DIV_CODE"))
	assert.Equal(t, "0001", capturedQuery.Get("FID_INPUT_ISCD"))
	assert.Equal(t, "20214", capturedQuery.Get("FID_COND_SCR_DIV_CODE"))
	assert.Equal(t, "K", capturedQuery.Get("FID_MRKT_CLS_CODE"))
	assert.Equal(t, "0", capturedQuery.Get("FID_BLNG_CLS_CODE"))

	// output1 (요약)
	d, _ := decimal.NewFromString("2650.45")
	assert.True(t, d.Equal(res.Output1.BstpNmixPrpr))
	assert.Equal(t, int64(350000000), res.Output1.AcmlVol)

	// output2 (Array)
	require.Len(t, res.Output2, 2)
	assert.Equal(t, "0001", res.Output2[0].BstpClsCode)
	assert.Equal(t, "코스피", res.Output2[0].HtsKorIsnm)
	assert.InDelta(t, 100.0, res.Output2[0].AcmlVolRlim, 0.01)
}
```

- [ ] **Step 2: FAIL**

- [ ] **Step 3: 구현 추가**

```go
// IndexCategoryPrice 는 국내업종 구분별 전체시세 (FHPUP02140000) 응답.
//
// 한투 docs: docs/api/국내주식/국내업종_구분별전체시세.md
// path: /uapi/domestic-stock/v1/quotations/inquire-index-category-price
type IndexCategoryPrice struct {
	Output1 IndexCategoryPriceSummary `json:"output1"`
	Output2 []IndexCategoryPriceItem  `json:"output2"`
}

// IndexCategoryPriceSummary 는 응답의 output1 (대표 업종 지수).
type IndexCategoryPriceSummary struct {
	BstpNmixPrpr         decimal.Decimal `json:"bstp_nmix_prpr"`
	BstpNmixPrdyVrss     decimal.Decimal `json:"bstp_nmix_prdy_vrss"`
	PrdyVrssSign         string          `json:"prdy_vrss_sign"`
	BstpNmixPrdyCtrt     float64         `json:"bstp_nmix_prdy_ctrt,string"`
	AcmlVol              int64           `json:"acml_vol,string"`
	AcmlTrPbmn           int64           `json:"acml_tr_pbmn,string"`
	BstpNmixOprc         decimal.Decimal `json:"bstp_nmix_oprc"`
	BstpNmixHgpr         decimal.Decimal `json:"bstp_nmix_hgpr"`
	BstpNmixLwpr         decimal.Decimal `json:"bstp_nmix_lwpr"`
	PrdyVol              int64           `json:"prdy_vol,string"`
	AscnIssuCnt          string          `json:"ascn_issu_cnt"`
	DownIssuCnt          string          `json:"down_issu_cnt"`
	StnrIssuCnt          string          `json:"stnr_issu_cnt"`
	UplmIssuCnt          string          `json:"uplm_issu_cnt"`
	LslmIssuCnt          string          `json:"lslm_issu_cnt"`
	PrdyTrPbmn           int64           `json:"prdy_tr_pbmn,string"`
	DryyBstpNmixHgprDate string          `json:"dryy_bstp_nmix_hgpr_date"`
	DryyBstpNmixHgpr     decimal.Decimal `json:"dryy_bstp_nmix_hgpr"`
	DryyBstpNmixLwpr     decimal.Decimal `json:"dryy_bstp_nmix_lwpr"`
	DryyBstpNmixLwprDate string          `json:"dryy_bstp_nmix_lwpr_date"`
}

// IndexCategoryPriceItem 은 응답의 output2 한 행 (구분별 업종).
type IndexCategoryPriceItem struct {
	BstpClsCode      string          `json:"bstp_cls_code"`              // 업종 구분 코드
	HtsKorIsnm       string          `json:"hts_kor_isnm"`               // HTS 한글 종목명
	BstpNmixPrpr     decimal.Decimal `json:"bstp_nmix_prpr"`             // 업종 지수 현재가
	BstpNmixPrdyVrss decimal.Decimal `json:"bstp_nmix_prdy_vrss"`        // 업종 지수 전일 대비
	PrdyVrssSign     string          `json:"prdy_vrss_sign"`             // 전일 대비 부호
	BstpNmixPrdyCtrt float64         `json:"bstp_nmix_prdy_ctrt,string"` // 업종 지수 전일 대비율
	AcmlVol          int64           `json:"acml_vol,string"`            // 누적 거래량
	AcmlTrPbmn       int64           `json:"acml_tr_pbmn,string"`        // 누적 거래 대금
	AcmlVolRlim      float64         `json:"acml_vol_rlim,string"`       // 누적 거래량 비중
	AcmlTrPbmnRlim   float64         `json:"acml_tr_pbmn_rlim,string"`   // 누적 거래 대금 비중
}

// InquireIndexCategoryPriceParams 는 국내업종 구분별 전체시세 조회 파라미터.
type InquireIndexCategoryPriceParams struct {
	MarketCode string // FID_COND_MRKT_DIV_CODE — 빈 값=>"U" (업종)
	Symbol     string // FID_INPUT_ISCD — 필수, 업종 코드 (코스피 0001 등)
	ScreenCode string // FID_COND_SCR_DIV_CODE — 빈 값=>"20214"
	MarketCls  string // FID_MRKT_CLS_CODE — "K":거래소, "Q":코스닥, "K2":코스피200
	BelongCls  string // FID_BLNG_CLS_CODE — 빈 값=>"0" (전업종)
}

// InquireIndexCategoryPrice 는 국내업종 구분별 전체시세 호출.
//
// 한투 docs: docs/api/국내주식/국내업종_구분별전체시세.md
// path: /uapi/domestic-stock/v1/quotations/inquire-index-category-price (FHPUP02140000)
func (c *Client) InquireIndexCategoryPrice(ctx context.Context, params InquireIndexCategoryPriceParams) (*IndexCategoryPrice, error) {
	market := params.MarketCode
	if market == "" {
		market = "U"
	}
	scr := params.ScreenCode
	if scr == "" {
		scr = "20214"
	}
	belong := params.BelongCls
	if belong == "" {
		belong = "0"
	}

	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/domestic-stock/v1/quotations/inquire-index-category-price",
		TrID:   "FHPUP02140000",
		Query: map[string]string{
			"FID_COND_MRKT_DIV_CODE": market,
			"FID_INPUT_ISCD":         params.Symbol,
			"FID_COND_SCR_DIV_CODE":  scr,
			"FID_MRKT_CLS_CODE":      params.MarketCls,
			"FID_BLNG_CLS_CODE":      belong,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}

	var res IndexCategoryPrice
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse IndexCategoryPrice: %w", err)
	}
	return &res, nil
}
```

- [ ] **Step 4: PASS**

- [ ] **Step 5: Commit**

```bash
git add domestic/industry.go domestic/industry_test.go
git commit -m "$(cat <<'EOF'
[feat] domestic — InquireIndexCategoryPrice (국내업종 구분별 전체시세, FHPUP02140000)

IndexCategoryPrice (Output1 요약 20 + Output2 구분별 array 10 필드) +
InquireIndexCategoryPriceParams (5 query). MarketCls (K:거래소/Q:코스닥/K2)
+ BelongCls (소속 구분).

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 7: domestic/ipo.go — InquirePubOffer

**Files:**
- Create: `domestic/ipo.go`
- Create: `domestic/ipo_test.go`

> 예탁원정보(공모주청약일정). output1 (Array, 13 fields per row). KIS docs (`docs/api/국내주식/예탁원정보(공모주청약일정).md`). query 키가 다른 형식 (`SHT_CD`, `CTS`, `F_DT`, `T_DT`).

- [ ] **Step 1: 테스트 작성**

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

func TestClient_InquirePubOffer(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/ksdinfo/pub-offer`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "pub_offer_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquirePubOffer(context.Background(), domestic.InquirePubOfferParams{
		FromDate: "20260501",
		ToDate:   "20260531",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	// query 키가 대문자 + 한글식
	assert.Equal(t, "", capturedQuery.Get("SHT_CD"))
	assert.Equal(t, "", capturedQuery.Get("CTS"))
	assert.Equal(t, "20260501", capturedQuery.Get("F_DT"))
	assert.Equal(t, "20260531", capturedQuery.Get("T_DT"))

	require.Len(t, res.Output1, 2)
	assert.Equal(t, "999998", res.Output1[0].ShtCd)
	assert.Equal(t, "샘플바이오", res.Output1[0].IsinName)
	assert.Equal(t, decimal.NewFromInt(30000), res.Output1[0].FixSubscrPri)
	assert.Equal(t, decimal.NewFromInt(100), res.Output1[0].FaceValue)
	assert.Equal(t, "20260505 ~ 20260506", res.Output1[0].SubscrDt)
	assert.Equal(t, "한국투자증권", res.Output1[0].LeadMgr)
}
```

- [ ] **Step 2: FAIL**

- [ ] **Step 3: 구현** — `domestic/ipo.go`

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

// PubOffer 는 예탁원정보(공모주청약일정) (HHKDB669108C0) 응답.
//
// 한투 docs: docs/api/국내주식/예탁원정보(공모주청약일정).md
// path: /uapi/domestic-stock/v1/ksdinfo/pub-offer
//
// IPO 청약일정 list. output1 (Array). 다른 ranking/financial 과 query 키 형식이 다름 (대문자+한글식).
type PubOffer struct {
	Output1 []PubOfferItem `json:"output1"`
}

// PubOfferItem 은 한 IPO 청약 일정 항목.
type PubOfferItem struct {
	RecordDate   string          `json:"record_date"`     // 기준일 (YYYYMMDD)
	ShtCd        string          `json:"sht_cd"`          // 종목코드
	IsinName     string          `json:"isin_name"`       // 종목명
	FixSubscrPri decimal.Decimal `json:"fix_subscr_pri"`  // 공모가
	FaceValue    decimal.Decimal `json:"face_value"`      // 액면가
	SubscrDt     string          `json:"subscr_dt"`       // 청약기간 (예 "20260505 ~ 20260506")
	PayDt        string          `json:"pay_dt"`          // 납입일 (YYYYMMDD)
	RefundDt     string          `json:"refund_dt"`       // 환불일
	ListDt       string          `json:"list_dt"`         // 상장/등록일
	LeadMgr      string          `json:"lead_mgr"`        // 주간사
	PubBfCap     int64           `json:"pub_bf_cap,string"`     // 공모전 자본금
	PubAfCap     int64           `json:"pub_af_cap,string"`     // 공모후 자본금
	AssignStkQty int64           `json:"assign_stk_qty,string"` // 당사 배정물량
}

// InquirePubOfferParams 는 공모주청약일정 조회 파라미터.
//
// 다른 ranking 과 query 키 형식이 다름 — KIS docs 그대로 노출 (SHT_CD, CTS, F_DT, T_DT).
type InquirePubOfferParams struct {
	Symbol   string // SHT_CD — 종목코드. 빈 값(공백) = 전체
	Cts      string // CTS — 빈 값(공백) default
	FromDate string // F_DT — 조회일자 From (YYYYMMDD)
	ToDate   string // T_DT — 조회일자 To (YYYYMMDD)
}

// InquirePubOffer 는 예탁원정보(공모주청약일정) 호출.
//
// 한투 docs: docs/api/국내주식/예탁원정보(공모주청약일정).md
// path: /uapi/domestic-stock/v1/ksdinfo/pub-offer (HHKDB669108C0)
//
// 공모주(IPO) 청약일정 list 조회. Symbol 빈 값 시 전체.
func (c *Client) InquirePubOffer(ctx context.Context, params InquirePubOfferParams) (*PubOffer, error) {
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/domestic-stock/v1/ksdinfo/pub-offer",
		TrID:   "HHKDB669108C0",
		Query: map[string]string{
			"SHT_CD": params.Symbol,
			"CTS":    params.Cts,
			"F_DT":   params.FromDate,
			"T_DT":   params.ToDate,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}

	var res PubOffer
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse PubOffer: %w", err)
	}
	return &res, nil
}
```

- [ ] **Step 4: PASS**

- [ ] **Step 5: 전체 회귀 테스트**

```bash
go test ./... -count=1
```
Expected: 모든 패키지 PASS.

- [ ] **Step 6: Commit**

```bash
git add domestic/ipo.go domestic/ipo_test.go
git commit -m "$(cat <<'EOF'
[feat] domestic — InquirePubOffer (예탁원정보 공모주청약일정, HHKDB669108C0)

PubOffer + PubOfferItem (13 필드) + InquirePubOfferParams (4 query).
ksdinfo/ path (다른 quotations/ 와 다름). query 키 대문자+한글식
(SHT_CD/CTS/F_DT/T_DT) 그대로 노출.

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 8: examples/domestic_investor/main.go

**Files:**
- Create: `examples/domestic_investor/main.go`

- [ ] **Step 1: example 작성**

```go
// domestic_investor example: InquireInvestorTradeByStockDaily +
// InquireInvestorDailyByMarket + InquireIndexPrice + InquirePubOffer.
//
// Run: KIS credentials env vars 후 go run ./examples/domestic_investor
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	kis "github.com/kenshin579/korea-investment-stock"
	"github.com/kenshin579/korea-investment-stock/domestic"
)

func main() {
	client, err := kis.NewClientFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	yesterday := time.Now().AddDate(0, 0, -1).Format("20060102")
	symbol := "005930"

	// 1. 종목별 투자자매매동향 (일별, 어제 기준 N일치)
	stockDaily, err := client.Domestic.InquireInvestorTradeByStockDaily(ctx, domestic.InquireInvestorTradeByStockDailyParams{
		Symbol:   symbol,
		BaseDate: yesterday,
	})
	if err != nil {
		log.Fatalf("InquireInvestorTradeByStockDaily: %v", err)
	}
	fmt.Printf("[%s] 종목별 투자자매매동향 (어제 기준 %d일치)\n", symbol, len(stockDaily.Output2))
	for i, item := range stockDaily.Output2 {
		if i >= 5 {
			fmt.Println("  ... (이하 생략)")
			break
		}
		fmt.Printf("  %s: 외국인=%d, 개인=%d, 기관=%d (주)\n",
			item.StckBsopDate, item.FrgnNtbyQty, item.PrsnNtbyQty, item.OrgnNtbyQty)
	}

	// 2. 시장별 투자자매매동향 (일별, 코스피 종합)
	marketDaily, err := client.Domestic.InquireInvestorDailyByMarket(ctx, domestic.InquireInvestorDailyByMarketParams{
		Symbol:    "0001",
		BaseDate:  yesterday,
		Market:    "KSP",
		BaseDate2: yesterday,
		SubCode:   "0001",
	})
	if err != nil {
		log.Fatalf("InquireInvestorDailyByMarket: %v", err)
	}
	fmt.Printf("\n시장별 (코스피 종합) 투자자매매동향 %d 행\n", len(marketDaily.Output))
	for i, item := range marketDaily.Output {
		if i >= 3 {
			fmt.Println("  ... (이하 생략)")
			break
		}
		fmt.Printf("  %s: 지수=%s, 외국인=%d, 개인=%d (주)\n",
			item.StckBsopDate, item.BstpNmixPrpr, item.FrgnNtbyQty, item.PrsnNtbyQty)
	}

	// 3. 코스피 현재 지수
	idx, err := client.Domestic.InquireIndexPrice(ctx, domestic.InquireIndexPriceParams{
		Symbol: "0001",
	})
	if err != nil {
		log.Fatalf("InquireIndexPrice: %v", err)
	}
	fmt.Printf("\n코스피 현재 지수: %s (%v%%)\n", idx.Output.BstpNmixPrpr, idx.Output.BstpNmixPrdyCtrt)
	fmt.Printf("  상승=%s 하락=%s 보합=%s\n",
		idx.Output.AscnIssuCnt, idx.Output.DownIssuCnt, idx.Output.StnrIssuCnt)

	// 4. 공모주 청약일정 (다음 달)
	from := time.Now().Format("20060102")
	to := time.Now().AddDate(0, 1, 0).Format("20060102")
	ipo, err := client.Domestic.InquirePubOffer(ctx, domestic.InquirePubOfferParams{
		FromDate: from,
		ToDate:   to,
	})
	if err != nil {
		log.Fatalf("InquirePubOffer: %v", err)
	}
	fmt.Printf("\n공모주 청약일정 (%s ~ %s) %d 건\n", from, to, len(ipo.Output1))
	for i, item := range ipo.Output1 {
		if i >= 5 {
			fmt.Println("  ... (이하 생략)")
			break
		}
		fmt.Printf("  %s (%s) 공모가=%s, 청약일=%s, 주간사=%s\n",
			item.IsinName, item.ShtCd, item.FixSubscrPri, item.SubscrDt, item.LeadMgr)
	}
}
```

- [ ] **Step 2: 컴파일 검증**

Run: `go build ./examples/domestic_investor && echo OK`
Expected: `OK`.

- [ ] **Step 3: Commit**

```bash
git add examples/domestic_investor
git commit -m "$(cat <<'EOF'
[feat] examples/domestic_investor — Investor + Index + IPO 통합 예시

종목별/시장별 투자자매매동향 + 코스피 현재 지수 + 공모주 청약일정 출력.

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 9: 문서 갱신 (CLAUDE.md, README.md, CHANGELOG.md, domestic/doc.go)

**Files:**
- Modify: `CLAUDE.md`
- Modify: `README.md`
- Modify: `CHANGELOG.md`
- Modify: `domestic/doc.go`

- [ ] **Step 1: README.md — Available Methods 표 갱신**

기존 `## Available Methods (Phase 1.2 + 1.3)` 헤딩과 16 메서드 표를 다음으로 REPLACE:

```markdown
## Available Methods (Phase 1.2 + 1.3 + 1.4)

| Method | 한투 path | TR_ID |
|--------|----------|-------|
| `Domestic.InquirePrice` | `inquire-price` | FHKST01010100 |
| `Domestic.SearchInfo` | `search-info` | CTPF1604R |
| `Domestic.SearchStockInfo` | `search-stock-info` | CTPF1002R |
| `Domestic.InquireDailyItemChartPrice` | `inquire-daily-itemchartprice` | FHKST03010100 |
| `Domestic.InquireTimeItemChartPrice` | `inquire-time-itemchartprice` | FHKST03010200 |
| `Domestic.FetchKospiSymbols` | (KRX 공개 마스터) | — |
| `Domestic.FetchKosdaqSymbols` | (KRX 공개 마스터) | — |
| `Domestic.InquireVolumeRank` | `quotations/volume-rank` | FHPST01710000 |
| `Domestic.InquireFluctuation` | `ranking/fluctuation` | FHPST01700000 |
| `Domestic.InquireMarketCap` | `ranking/market-cap` | FHPST01740000 |
| `Domestic.InquireDividendRate` | `ranking/dividend-rate` | HHKDB13470100 |
| `Domestic.InquireFinancialRatio` | `finance/financial-ratio` | FHKST66430300 |
| `Domestic.InquireIncomeStatement` | `finance/income-statement` | FHKST66430200 |
| `Domestic.InquireBalanceSheet` | `finance/balance-sheet` | FHKST66430100 |
| `Domestic.InquireProfitRatio` | `finance/profit-ratio` | FHKST66430400 |
| `Domestic.InquireGrowthRatio` | `finance/growth-ratio` | FHKST66430800 |
| `Domestic.InquireInvestorTradeByStockDaily` | `quotations/investor-trade-by-stock-daily` | FHPTJ04160001 |
| `Domestic.InquireInvestorDailyByMarket` | `quotations/inquire-investor-daily-by-market` | FHPTJ04040000 |
| `Domestic.InquireInvestorTimeByMarket` | `quotations/inquire-investor-time-by-market` | FHPTJ04030000 |
| `Domestic.InquireIndexPrice` | `quotations/inquire-index-price` | FHPUP02100000 |
| `Domestic.InquireIndexCategoryPrice` | `quotations/inquire-index-category-price` | FHPUP02140000 |
| `Domestic.InquirePubOffer` | `ksdinfo/pub-offer` | HHKDB669108C0 |
```

- [ ] **Step 2: CLAUDE.md — banner 갱신**

기존 `> **Phase 1.3 — domestic 순위/재무 (v1.1.0)...`** 를 다음으로 REPLACE:

```
> **Phase 1.4 — domestic 투자자/업종/IPO (v1.2.0).** Phase 1.5+ 메서드는 추후 sub-plan 으로.
```

또한 spec 링크 section 에 Phase 1.4 plan 추가:
```markdown
- Phase 1.4 implementation plan: [`docs/superpowers/specs/2026-05-05-phase1-4-investor-industry-ipo-implementation-plan.md`](docs/superpowers/specs/2026-05-05-phase1-4-investor-industry-ipo-implementation-plan.md)
```

- [ ] **Step 3: CHANGELOG.md — `[1.2.0]` entry 추가**

`[1.1.0]` 위에 새 section 추가:

```markdown
## [1.2.0] - 2026-05-05

### Added — Phase 1.4 (국내주식 투자자/업종/IPO)

- `Domestic.InquireInvestorTradeByStockDaily` — 종목별 투자자매매동향 일별 (FHPTJ04160001)
- `Domestic.InquireInvestorDailyByMarket` — 시장별 투자자매매동향 일별 (FHPTJ04040000)
- `Domestic.InquireInvestorTimeByMarket` — 시장별 투자자매매동향 시세 (FHPTJ04030000)
- `Domestic.InquireIndexPrice` — 국내업종 현재지수 (FHPUP02100000)
- `Domestic.InquireIndexCategoryPrice` — 국내업종 구분별 전체시세 (FHPUP02140000)
- `Domestic.InquirePubOffer` — 예탁원정보 공모주청약일정 (HHKDB669108C0)
- examples: `domestic_investor`

### Notes

- IPO helpers 9개 omit — Phase 1.2 amendment 의 "Python wrapper convenience 미반영" 정책 일관 (client-side data 가공이라 caller 가 직접 처리)
- 투자자 매매동향 응답이 매우 큼 (종목별 일별: 95+ 필드, 시세: 78 필드) — KIS docs 1:1 매핑, struct field 모두 포함
- `InquireInvestorTimeByMarket` 의 query 키가 소문자 `fid_input_iscd` (다른 quotations/ 메서드와 다름) — KIS docs 그대로 노출
- `InquirePubOffer` 의 query 키가 대문자+한글식 (`SHT_CD`, `CTS`, `F_DT`, `T_DT`) + path 가 `ksdinfo/` (다른 메서드의 `quotations/`/`ranking/`/`finance/` 와 다름)
```

- [ ] **Step 4: domestic/doc.go 갱신**

기존 16 메서드 + 6 추가 = 22 메서드 list 로 갱신:

```go
// Package domestic 은 한국투자증권 OpenAPI 의 국내주식 카테고리 메서드.
//
// Phase 1.2 메서드 (7):
//
//   - InquirePrice                 — 주식현재가 시세 (FHKST01010100)
//   - SearchInfo                   — 상품기본조회 (CTPF1604R)
//   - SearchStockInfo              — 주식기본조회 (CTPF1002R)
//   - InquireDailyItemChartPrice   — 국내주식기간별시세 일/주/월/년 (FHKST03010100)
//   - InquireTimeItemChartPrice    — 주식당일분봉조회 (FHKST03010200)
//   - FetchKospiSymbols            — KRX KOSPI 마스터 (한투 API 가 아닌 KRX 공개 다운로드)
//   - FetchKosdaqSymbols           — KRX KOSDAQ 마스터
//
// Phase 1.3 메서드 (9):
//
//   - InquireVolumeRank            — 거래량순위 (FHPST01710000)
//   - InquireFluctuation           — 등락률 순위 (FHPST01700000)
//   - InquireMarketCap             — 시가총액 상위 (FHPST01740000)
//   - InquireDividendRate          — 배당률 상위 (HHKDB13470100)
//   - InquireFinancialRatio        — 재무비율 (FHKST66430300)
//   - InquireIncomeStatement       — 손익계산서 (FHKST66430200)
//   - InquireBalanceSheet          — 대차대조표 (FHKST66430100)
//   - InquireProfitRatio           — 수익성비율 (FHKST66430400)
//   - InquireGrowthRatio           — 성장성비율 (FHKST66430800)
//
// Phase 1.4 메서드 (6):
//
//   - InquireInvestorTradeByStockDaily — 종목별 투자자매매동향 일별 (FHPTJ04160001)
//   - InquireInvestorDailyByMarket    — 시장별 투자자매매동향 일별 (FHPTJ04040000)
//   - InquireInvestorTimeByMarket     — 시장별 투자자매매동향 시세 (FHPTJ04030000)
//   - InquireIndexPrice               — 국내업종 현재지수 (FHPUP02100000)
//   - InquireIndexCategoryPrice       — 국내업종 구분별 전체시세 (FHPUP02140000)
//   - InquirePubOffer                 — 예탁원정보 공모주청약일정 (HHKDB669108C0)
//
// 사용자는 root kis.Client 의 Domestic 필드로 접근.
package domestic
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
[doc] Phase 1.4 메서드 문서 갱신 — CLAUDE.md, README.md, CHANGELOG.md, domestic/doc.go

Phase 1.4 의 6 메서드 (investor 3 + industry 2 + ipo 1) 목록 + CHANGELOG
[1.2.0] entry. CLAUDE.md banner 갱신 (Phase 1.3 → 1.4, v1.1.0 → v1.2.0).
domestic/doc.go 패키지 doc 에 Phase 1.4 section 추가.

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 10: 최종 점검

- [ ] **Step 1: gofmt cleanup (필요 시)**

```bash
gofmt -w domestic/investor.go domestic/investor_test.go domestic/industry.go domestic/industry_test.go domestic/ipo.go domestic/ipo_test.go
gofmt -l .
```
Expected: empty.

If diff exists, commit as `[chore] Phase 1.4 — gofmt cleanup`.

- [ ] **Step 2: 빌드/vet**

```bash
go build ./... && go vet ./...
```
Expected: silent.

- [ ] **Step 3: 모든 테스트 + race**

```bash
go test ./... -race -count=1
```
Expected: all PASS.

- [ ] **Step 4: Coverage**

```bash
go test ./... -coverprofile=/tmp/cov.out -covermode=atomic
go tool cover -func=/tmp/cov.out | tail -10
```

Expected:
- domestic/ ≥ 80%
- root kis ≥ 80%

If short, ADD targeted tests (don't lower thresholds).

- [ ] **Step 5: 디렉터리 구조 확인**

```bash
ls -la domestic/{investor,industry,ipo}.go domestic/{investor,industry,ipo}_test.go domestic/testdata/{investor_trade_by_stock_daily,investor_daily_by_market,investor_time_by_market,index_price,index_category_price,pub_offer}_success.json examples/domestic_investor/main.go 2>&1 | wc -l
```
Expected: 13 (6 testdata + 6 .go files + 1 example).

- [ ] **Step 6: Commit history**

```bash
git log main..HEAD --oneline | wc -l
```
Expected: ~10-12 commits.

---

## Task 11: PR 생성 (사용자 승인 후)

> Claude 는 push / PR 생성을 사용자 명시적 승인 후에만 실행 (글로벌 정책).

- [ ] **Step 1: 사용자 승인 요청**

작업 진행 보고 + PR 생성 가능 여부 confirm.

- [ ] **Step 2: Push branch**

```bash
git push -u origin docs/phase1-4-spec
```

- [ ] **Step 3: PR 생성**

```bash
gh pr create --title "Phase 1.4 — domestic 투자자/업종/IPO (v1.2.0)" --reviewer kenshin579 --base main --head docs/phase1-4-spec --body "$(cat <<'EOF'
## Summary

- 국내주식 투자자 매매동향 3 + 업종 2 + IPO 1 = 6 메서드 추가
- Phase 1.2/1.3 패턴 재사용 (Style A 메서드명, Params struct, KIS docs 1:1)
- 새 internal package 불필요
- IPO helpers 9개 미반영 (Phase 1.2 amendment 정책 일관)
- v1.2.0 release 대상

## 메서드 → 한투 API 매핑

| Go 메서드 | path | TR_ID |
|-----------|------|-------|
| InquireInvestorTradeByStockDaily | quotations/investor-trade-by-stock-daily | FHPTJ04160001 |
| InquireInvestorDailyByMarket | quotations/inquire-investor-daily-by-market | FHPTJ04040000 |
| InquireInvestorTimeByMarket | quotations/inquire-investor-time-by-market | FHPTJ04030000 |
| InquireIndexPrice | quotations/inquire-index-price | FHPUP02100000 |
| InquireIndexCategoryPrice | quotations/inquire-index-category-price | FHPUP02140000 |
| InquirePubOffer | ksdinfo/pub-offer | HHKDB669108C0 |

## Test Plan

- [x] go build/vet/fmt clean
- [x] go test ./... -race -count=1 모든 패키지 PASS
- [x] Coverage domestic/ ≥ 80%
- [x] httpmock 단위 테스트 (각 메서드)

## Breaking Changes

없음 — 신규 메서드 추가만.

🤖 Generated with [Claude Code](https://claude.com/claude-code)
EOF
)"
```

- [ ] **Step 4: Merge (사용자 승인 후)**

`gh pr merge <PR#> --merge`

- [ ] **Step 5: 후속 작업 (사용자 승인 후)**

```bash
git checkout main && git pull
git tag -a v1.2.0 -m "v1.2.0: Phase 1.4 — domestic 투자자/업종/IPO 6 메서드"
git push origin v1.2.0
gh release create v1.2.0 --title "v1.2.0 — Phase 1.4: domestic 투자자/업종/IPO" --notes-file <(awk '/^## \[/{c++} c==1' CHANGELOG.md)
```
