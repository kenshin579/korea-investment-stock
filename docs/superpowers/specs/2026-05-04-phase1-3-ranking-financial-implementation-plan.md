# Phase 1.3 — 국내 순위 + 재무 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** `korea-investment-stock` Go 라이브러리에 국내주식 순위 (4 메서드) + 재무제표 (5 메서드) 총 9 메서드 추가 (`v1.1.0` release).

**Architecture:** Phase 1.2 의 인프라 + 패턴 재사용. `domestic/ranking.go` (4 메서드) + `domestic/financial.go` (5 메서드) 추가. 한투 API path 1:1 매핑 (Style A — endpoint path 의 마지막 segment 를 PascalCase). 새 internal package 불필요 (REST only). TDD 흐름: testdata fixture (한투 docs 응답 필드 정의 → 합성 JSON) → 실패 테스트 → struct + 메서드 구현 → 통과 → commit.

**Tech Stack:** Go 1.25+, `github.com/jarcoal/httpmock` (Phase 1.1 부터), `github.com/shopspring/decimal`, `github.com/stretchr/testify`. 새 dependency 없음.

**참고 spec:**
- Phase 1 design spec (Phase 1.3 amendment 적용): `docs/superpowers/specs/2026-05-03-phase1-api-coverage-design.md` (commit `e3ed393`)
- Phase 1.2 plan (참조 패턴): `docs/superpowers/specs/2026-05-03-phase1-2-domestic-quotes-implementation-plan.md`
- 한투 API docs: `docs/api/국내주식/{거래량순위.md, 국내주식_등락률_순위.md, 국내주식_시가총액_상위.md, 국내주식_배당률_상위.md, 국내주식_재무비율.md, 국내주식_손익계산서.md, 국내주식_대차대조표.md, 국내주식_수익성비율.md, 국내주식_성장성비율.md}`

---

## 사전 정보

| 항목 | 값 |
|------|---|
| 작업 브랜치 | `docs/phase1-3-spec` (이미 생성됨) |
| 시작 HEAD | `e3ed393` (Phase 1 design spec amendment commit) |
| Release 목표 | `v1.1.0` (PR merge 후 태그) |
| PR 베이스 | `main` |
| 현재 main 상태 | v1.0.0 publish 완료 (Phase 1.1 + 1.2 통합) |

---

## 메서드 → 한투 API 매핑

| Go 메서드 | 한투 path | TR_ID | docs |
|-----------|----------|-------|------|
| `Domestic.InquireVolumeRank(ctx, params)` | `/uapi/domestic-stock/v1/quotations/volume-rank` | FHPST01710000 | 거래량순위.md |
| `Domestic.InquireFluctuation(ctx, params)` | `/uapi/domestic-stock/v1/ranking/fluctuation` | FHPST01700000 | 국내주식_등락률_순위.md |
| `Domestic.InquireMarketCap(ctx, params)` | `/uapi/domestic-stock/v1/ranking/market-cap` | FHPST01740000 | 국내주식_시가총액_상위.md |
| `Domestic.InquireDividendRate(ctx, params)` | `/uapi/domestic-stock/v1/ranking/dividend-rate` | HHKDB13470100 | 국내주식_배당률_상위.md |
| `Domestic.InquireFinancialRatio(ctx, params)` | `/uapi/domestic-stock/v1/finance/financial-ratio` | FHKST66430300 | 국내주식_재무비율.md |
| `Domestic.InquireIncomeStatement(ctx, params)` | `/uapi/domestic-stock/v1/finance/income-statement` | FHKST66430200 | 국내주식_손익계산서.md |
| `Domestic.InquireBalanceSheet(ctx, params)` | `/uapi/domestic-stock/v1/finance/balance-sheet` | FHKST66430100 | 국내주식_대차대조표.md |
| `Domestic.InquireProfitRatio(ctx, params)` | `/uapi/domestic-stock/v1/finance/profit-ratio` | FHKST66430400 | 국내주식_수익성비율.md |
| `Domestic.InquireGrowthRatio(ctx, params)` | `/uapi/domestic-stock/v1/finance/growth-ratio` | FHKST66430800 | 국내주식_성장성비율.md |

**참고**: 거래량순위만 path 가 `quotations/` (Phase 1.2 의 `inquire-price` 와 동일 group), 나머지 ranking 3개는 `ranking/` 디렉터리. financial 5개는 `finance/`. KIS API 의 query parameter naming 도 inconsistent (거래량순위만 대문자 `FID_*`, 나머지 ranking 은 소문자 `fid_*`). KIS docs 그대로 노출.

---

## 파일 구조

### 신규 (domestic)

- `domestic/ranking.go` — 4 ranking 메서드 + Result struct + Params struct
- `domestic/ranking_test.go`
- `domestic/financial.go` — 5 financial 메서드 + Result struct + Params struct
- `domestic/financial_test.go`
- `domestic/testdata/volume_rank_success.json`
- `domestic/testdata/fluctuation_success.json`
- `domestic/testdata/market_cap_success.json`
- `domestic/testdata/dividend_rate_success.json`
- `domestic/testdata/financial_ratio_success.json`
- `domestic/testdata/income_statement_success.json`
- `domestic/testdata/balance_sheet_success.json`
- `domestic/testdata/profit_ratio_success.json`
- `domestic/testdata/growth_ratio_success.json`

### 신규 (examples)

- `examples/domestic_ranking/main.go` — VolumeRank + Fluctuation + MarketCap 사용 예
- `examples/domestic_financial/main.go` — FinancialRatio + IncomeStatement + BalanceSheet 사용 예

### 수정 (root)

- `CLAUDE.md` — Phase 1.3 메서드 안내
- `README.md` — Available Methods 표 갱신
- `CHANGELOG.md` — `[1.1.0]` entry
- `domestic/doc.go` — Phase 1.3 메서드 안내 갱신

---

## 타입 매핑 규칙 (Phase 1.2 와 동일)

- **가격/주당 가치 → `decimal.Decimal` (bare tag)**: `stck_prpr`, `prdy_vrss`, `stck_hgpr`, `stck_lwpr`, `oprc_vrss_prpr`, `prd_rsfl`, `eps`, `bps`, `sps`, `per_sto_divi_amt`
- **수량/금액 → `int64,string`**: `acml_vol`, `prdy_vol`, `lstn_stcn`, `avrg_vol`, `avrg_tr_pbmn`, `acml_tr_pbmn`, `stck_avls`, `cnnt_ascn_dynu`, `cnnt_down_dynu`, 손익계산서 amount, 대차대조표 amount
- **비율/percentage → `float64,string`**: `prdy_ctrt`, `vol_inrt`, `vol_tnrt`, `tr_pbmn_tnrt`, `lwpr_vrss_prpr_rate`, `hgpr_vrss_prpr_rate`, `dsgt_date_clpr_vrss_prpr_rate`, `oprc_vrss_prpr_rate`, `prd_rsfl_rate`, `mrkt_whol_avls_rlim`, `divi_rate`, `grs`, `bsop_prfi_inrt`, `ntin_inrt`, `roe_val`, `equt_inrt`, `totl_aset_inrt`, `cptl_ntin_rate`, `self_cptl_ntin_inrt`, `sale_ntin_rate`, `sale_totl_rate`, `rsrv_rate`, `lblt_rate`, `n_befr_clpr_vrss_prpr_rate`, `nday_vol_tnrt`, `nday_tr_pbmn_tnrt`
- **코드/이름/날짜/Y-N/부호 → 평문 `string`**: `mksc_shrn_iscd`, `stck_shrn_iscd`, `hts_kor_isnm`, `sht_cd`, `isin_name`, `data_rank`, `rank`, `prdy_vrss_sign`, `oprc_vrss_prpr_sign`, `hgpr_hour`, `lwpr_hour`, `acml_hgpr_date`, `acml_lwpr_date`, `record_date`, `stac_yymm`, `divi_kind`

---

## Task 1: testdata fixtures (9 개 합성 JSON)

**Files:**
- Create: `domestic/testdata/volume_rank_success.json`
- Create: `domestic/testdata/fluctuation_success.json`
- Create: `domestic/testdata/market_cap_success.json`
- Create: `domestic/testdata/dividend_rate_success.json`
- Create: `domestic/testdata/financial_ratio_success.json`
- Create: `domestic/testdata/income_statement_success.json`
- Create: `domestic/testdata/balance_sheet_success.json`
- Create: `domestic/testdata/profit_ratio_success.json`
- Create: `domestic/testdata/growth_ratio_success.json`

> 9 testdata 모두 한 task 로 묶음. KIS docs 의 응답 필드 list 기반 합성 JSON. 값은 가능한 합리적 (실제 시세 아니지만 stretch 하지 않은 정상 범위).

- [ ] **Step 1: `domestic/testdata/volume_rank_success.json`**

```json
{
  "rt_cd": "0",
  "msg_cd": "MCA00000",
  "msg1": "정상처리 되었습니다.",
  "Output": [
    {
      "hts_kor_isnm": "삼성전자",
      "mksc_shrn_iscd": "005930",
      "data_rank": "1",
      "stck_prpr": "75800",
      "prdy_vrss_sign": "5",
      "prdy_vrss": "-200",
      "prdy_ctrt": "-0.26",
      "acml_vol": "12345678",
      "prdy_vol": "11000000",
      "lstn_stcn": "5969782550",
      "avrg_vol": "10500000",
      "n_befr_clpr_vrss_prpr_rate": "1.20",
      "vol_inrt": "12.23",
      "vol_tnrt": "0.21",
      "nday_vol_tnrt": "1.05",
      "avrg_tr_pbmn": "800000000000",
      "tr_pbmn_tnrt": "15.42",
      "nday_tr_pbmn_tnrt": "1.20",
      "acml_tr_pbmn": "938223456000"
    },
    {
      "hts_kor_isnm": "SK하이닉스",
      "mksc_shrn_iscd": "000660",
      "data_rank": "2",
      "stck_prpr": "210000",
      "prdy_vrss_sign": "2",
      "prdy_vrss": "3000",
      "prdy_ctrt": "1.45",
      "acml_vol": "8500000",
      "prdy_vol": "7800000",
      "lstn_stcn": "728002365",
      "avrg_vol": "7000000",
      "n_befr_clpr_vrss_prpr_rate": "2.10",
      "vol_inrt": "8.97",
      "vol_tnrt": "1.17",
      "nday_vol_tnrt": "1.20",
      "avrg_tr_pbmn": "1500000000000",
      "tr_pbmn_tnrt": "8.50",
      "nday_tr_pbmn_tnrt": "1.10",
      "acml_tr_pbmn": "1785000000000"
    }
  ]
}
```

- [ ] **Step 2: `domestic/testdata/fluctuation_success.json`**

```json
{
  "rt_cd": "0",
  "msg_cd": "MCA00000",
  "msg1": "정상처리 되었습니다.",
  "output": [
    {
      "stck_shrn_iscd": "005930",
      "data_rank": "1",
      "hts_kor_isnm": "삼성전자",
      "stck_prpr": "75800",
      "prdy_vrss": "-200",
      "prdy_vrss_sign": "5",
      "prdy_ctrt": "-0.26",
      "acml_vol": "12345678",
      "stck_hgpr": "76200",
      "hgpr_hour": "131542",
      "acml_hgpr_date": "20260502",
      "stck_lwpr": "75500",
      "lwpr_hour": "094230",
      "acml_lwpr_date": "20260502",
      "lwpr_vrss_prpr_rate": "0.40",
      "dsgt_date_clpr_vrss_prpr_rate": "-0.26",
      "cnnt_ascn_dynu": "0",
      "hgpr_vrss_prpr_rate": "-0.52",
      "cnnt_down_dynu": "1",
      "oprc_vrss_prpr_sign": "5",
      "oprc_vrss_prpr": "-200",
      "oprc_vrss_prpr_rate": "-0.26",
      "prd_rsfl": "-200",
      "prd_rsfl_rate": "-0.26"
    }
  ]
}
```

- [ ] **Step 3: `domestic/testdata/market_cap_success.json`**

```json
{
  "rt_cd": "0",
  "msg_cd": "MCA00000",
  "msg1": "정상처리 되었습니다.",
  "output": [
    {
      "mksc_shrn_iscd": "005930",
      "data_rank": "1",
      "hts_kor_isnm": "삼성전자",
      "stck_prpr": "75800",
      "prdy_vrss": "-200",
      "prdy_vrss_sign": "5",
      "prdy_ctrt": "-0.26",
      "acml_vol": "12345678",
      "lstn_stcn": "5969782550",
      "stck_avls": "452329543",
      "mrkt_whol_avls_rlim": "20.45"
    },
    {
      "mksc_shrn_iscd": "000660",
      "data_rank": "2",
      "hts_kor_isnm": "SK하이닉스",
      "stck_prpr": "210000",
      "prdy_vrss": "3000",
      "prdy_vrss_sign": "2",
      "prdy_ctrt": "1.45",
      "acml_vol": "8500000",
      "lstn_stcn": "728002365",
      "stck_avls": "152880497",
      "mrkt_whol_avls_rlim": "6.91"
    }
  ]
}
```

- [ ] **Step 4: `domestic/testdata/dividend_rate_success.json`**

```json
{
  "rt_cd": "0",
  "msg_cd": "MCA00000",
  "msg1": "정상처리 되었습니다.",
  "output1": [
    {
      "rank": "1",
      "sht_cd": "005930",
      "isin_name": "삼성전자",
      "record_date": "20251231",
      "per_sto_divi_amt": "1444",
      "divi_rate": "1.91",
      "divi_kind": "현금배당"
    },
    {
      "rank": "2",
      "sht_cd": "000660",
      "isin_name": "SK하이닉스",
      "record_date": "20251231",
      "per_sto_divi_amt": "1500",
      "divi_rate": "0.71",
      "divi_kind": "현금배당"
    }
  ]
}
```

- [ ] **Step 5: `domestic/testdata/financial_ratio_success.json`**

```json
{
  "rt_cd": "0",
  "msg_cd": "MCA00000",
  "msg1": "정상처리 되었습니다.",
  "output": [
    {
      "stac_yymm": "202412",
      "grs": "5.42",
      "bsop_prfi_inrt": "12.30",
      "ntin_inrt": "8.75",
      "roe_val": "11.50",
      "eps": "6638",
      "sps": "55432",
      "bps": "57420",
      "rsrv_rate": "32.45",
      "lblt_rate": "28.10"
    },
    {
      "stac_yymm": "202312",
      "grs": "3.21",
      "bsop_prfi_inrt": "-2.50",
      "ntin_inrt": "-5.10",
      "roe_val": "9.80",
      "eps": "6105",
      "sps": "52600",
      "bps": "53000",
      "rsrv_rate": "30.10",
      "lblt_rate": "29.50"
    }
  ]
}
```

- [ ] **Step 6: `domestic/testdata/income_statement_success.json`**

```json
{
  "rt_cd": "0",
  "msg_cd": "MCA00000",
  "msg1": "정상처리 되었습니다.",
  "output": [
    {
      "stac_yymm": "202412",
      "sale_account": "279600000",
      "sale_cost": "176000000",
      "sale_totl_prfi": "103600000",
      "depr_cost": "99.99",
      "sell_mang": "99.99",
      "bsop_prti": "32830000",
      "bsop_non_ernn": "99.99",
      "bsop_non_expn": "99.99",
      "op_prfi": "30450000",
      "spec_prfi": "1200000",
      "spec_loss": "800000",
      "thtr_ntin": "23456000"
    }
  ]
}
```

- [ ] **Step 7: `domestic/testdata/balance_sheet_success.json`**

```json
{
  "rt_cd": "0",
  "msg_cd": "MCA00000",
  "msg1": "정상처리 되었습니다.",
  "output": [
    {
      "stac_yymm": "202412",
      "cras": "189000000",
      "fxas": "245000000",
      "total_aset": "434000000",
      "flow_lblt": "62000000",
      "fix_lblt": "32000000",
      "total_lblt": "94000000",
      "cpfn": "778038",
      "cfp_surp": "99.99",
      "prfi_surp": "99.99",
      "total_cptl": "340000000"
    }
  ]
}
```

- [ ] **Step 8: `domestic/testdata/profit_ratio_success.json`**

```json
{
  "rt_cd": "0",
  "msg_cd": "MCA00000",
  "msg1": "정상처리 되었습니다.",
  "output": [
    {
      "stac_yymm": "202412",
      "cptl_ntin_rate": "8.45",
      "self_cptl_ntin_inrt": "11.50",
      "sale_ntin_rate": "12.30",
      "sale_totl_rate": "37.05"
    }
  ]
}
```

- [ ] **Step 9: `domestic/testdata/growth_ratio_success.json`**

```json
{
  "rt_cd": "0",
  "msg_cd": "MCA00000",
  "msg1": "정상처리 되었습니다.",
  "output": [
    {
      "stac_yymm": "202412",
      "grs": "5.42",
      "bsop_prfi_inrt": "12.30",
      "equt_inrt": "8.50",
      "totl_aset_inrt": "10.20"
    }
  ]
}
```

- [ ] **Step 10: 검증**

Run:
```bash
for f in domestic/testdata/{volume_rank,fluctuation,market_cap,dividend_rate,financial_ratio,income_statement,balance_sheet,profit_ratio,growth_ratio}_success.json; do
  python3 -c "import json; json.load(open('$f'))" && echo "$f OK" || echo "$f BROKEN"
done
```

Expected: 9 줄 모두 `OK`.

- [ ] **Step 11: Commit**

```bash
git add domestic/testdata/{volume_rank,fluctuation,market_cap,dividend_rate,financial_ratio,income_statement,balance_sheet,profit_ratio,growth_ratio}_success.json
git commit -m "$(cat <<'EOF'
[chore] Phase 1.3 testdata — 9 합성 JSON fixtures

ranking 4 (volume_rank, fluctuation, market_cap, dividend_rate) +
financial 5 (financial_ratio, income_statement, balance_sheet,
profit_ratio, growth_ratio). 한투 docs (docs/api/국내주식/<API>.md) 의
응답 필드 정의 기반 합성. 값은 합리적 정상 범위.

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 2: domestic/ranking.go — InquireVolumeRank + 공통 base

**Files:**
- Create: `domestic/ranking.go`
- Create: `domestic/ranking_test.go`

> 첫 ranking 메서드 — file 도 같이 생성. 거래량순위 응답 키가 대문자 `Output` (KIS docs 표기). 다른 ranking 은 소문자 `output` 사용.

- [ ] **Step 1: 테스트 작성** — `domestic/ranking_test.go`

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

func TestClient_InquireVolumeRank(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/quotations/volume-rank`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "volume_rank_success.json")), nil
		},
	)

	c := newTestClient(t)
	rank, err := c.InquireVolumeRank(context.Background(), domestic.InquireVolumeRankParams{
		InputISCD: "0000",
	})
	require.NoError(t, err)
	require.NotNil(t, rank)

	// 필수 query 기본값 검증
	assert.Equal(t, "J", capturedQuery.Get("FID_COND_MRKT_DIV_CODE"))
	assert.Equal(t, "20171", capturedQuery.Get("FID_COND_SCR_DIV_CODE"))
	assert.Equal(t, "0000", capturedQuery.Get("FID_INPUT_ISCD"))
	assert.Equal(t, "0", capturedQuery.Get("FID_DIV_CLS_CODE"))
	assert.Equal(t, "0", capturedQuery.Get("FID_BLNG_CLS_CODE"))
	assert.Equal(t, "111111111", capturedQuery.Get("FID_TRGT_CLS_CODE"))
	assert.Equal(t, "0000000000", capturedQuery.Get("FID_TRGT_EXLS_CLS_CODE"))

	// 응답 검증
	require.Len(t, rank.Output, 2)
	assert.Equal(t, "삼성전자", rank.Output[0].HtsKorIsnm)
	assert.Equal(t, "005930", rank.Output[0].MkscShrnIscd)
	assert.Equal(t, "1", rank.Output[0].DataRank)
	assert.Equal(t, decimal.NewFromInt(75800), rank.Output[0].StckPrpr)
	assert.Equal(t, int64(12345678), rank.Output[0].AcmlVol)
	assert.InDelta(t, 0.21, rank.Output[0].VolTnrt, 0.001)
	assert.Equal(t, int64(938223456000), rank.Output[0].AcmlTrPbmn)
}
```

- [ ] **Step 2: 테스트 실행 → FAIL**

Run: `go test ./domestic/... -run InquireVolumeRank -v`
Expected: 컴파일 실패 (`InquireVolumeRank`, `VolumeRank`, `InquireVolumeRankParams` 미정의).

- [ ] **Step 3: 구현** — `domestic/ranking.go`

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

// VolumeRank 는 거래량순위 (FHPST01710000) 응답.
//
// 한투 docs: docs/api/국내주식/거래량순위.md
// path: /uapi/domestic-stock/v1/quotations/volume-rank
//
// 최대 30건 확인 가능, 다음 조회 불가.
type VolumeRank struct {
	Output []VolumeRankItem `json:"Output"` // KIS docs 가 대문자 'O' 표기
}

// VolumeRankItem 은 거래량순위 응답의 한 행.
type VolumeRankItem struct {
	HtsKorIsnm             string          `json:"hts_kor_isnm"`                    // HTS 한글 종목명
	MkscShrnIscd           string          `json:"mksc_shrn_iscd"`                  // 유가증권 단축 종목코드
	DataRank               string          `json:"data_rank"`                       // 데이터 순위
	StckPrpr               decimal.Decimal `json:"stck_prpr"`                       // 주식 현재가
	PrdyVrssSign           string          `json:"prdy_vrss_sign"`                  // 전일 대비 부호
	PrdyVrss               decimal.Decimal `json:"prdy_vrss"`                       // 전일 대비
	PrdyCtrt               float64         `json:"prdy_ctrt,string"`                // 전일 대비율
	AcmlVol                int64           `json:"acml_vol,string"`                 // 누적 거래량
	PrdyVol                int64           `json:"prdy_vol,string"`                 // 전일 거래량
	LstnStcn               int64           `json:"lstn_stcn,string"`                // 상장 주수
	AvrgVol                int64           `json:"avrg_vol,string"`                 // 평균 거래량
	NBefrClprVrssPrprRate  float64         `json:"n_befr_clpr_vrss_prpr_rate,string"` // N일전종가대비현재가대비율
	VolInrt                float64         `json:"vol_inrt,string"`                 // 거래량 증가율
	VolTnrt                float64         `json:"vol_tnrt,string"`                 // 거래량 회전율
	NdayVolTnrt            float64         `json:"nday_vol_tnrt,string"`            // N일 거래량 회전율
	AvrgTrPbmn             int64           `json:"avrg_tr_pbmn,string"`             // 평균 거래 대금
	TrPbmnTnrt             float64         `json:"tr_pbmn_tnrt,string"`             // 거래대금 회전율
	NdayTrPbmnTnrt         float64         `json:"nday_tr_pbmn_tnrt,string"`        // N일 거래대금 회전율
	AcmlTrPbmn             int64           `json:"acml_tr_pbmn,string"`             // 누적 거래 대금
}

// InquireVolumeRankParams 는 거래량순위 조회 파라미터.
//
// 필수: InputISCD (종목코드 또는 "0000" 전체).
// 나머지는 zero-value 시 sensible default 사용.
type InquireVolumeRankParams struct {
	MarketCode    string // FID_COND_MRKT_DIV_CODE — "J":KRX, "NX":NXT. 빈 값=>"J"
	ScreenCode    string // FID_COND_SCR_DIV_CODE — Unique key. 빈 값=>"20171"
	InputISCD     string // FID_INPUT_ISCD — 필수, "0000"(전체) 또는 업종코드
	DivCode       string // FID_DIV_CLS_CODE — "0":전체, "1":보통주, "2":우선주. 빈 값=>"0"
	BelongCode    string // FID_BLNG_CLS_CODE — "0":평균거래량, "1":거래증가율, "2":평균거래회전율, "3":거래금액순, "4":평균거래금액회전율. 빈 값=>"0"
	TargetCode    string // FID_TRGT_CLS_CODE — 9자리 (증거금/신용보증금 비율). 빈 값=>"111111111"
	TargetExclude string // FID_TRGT_EXLS_CLS_CODE — 10자리 (제외 항목). 빈 값=>"0000000000"
	InputPrice1   string // FID_INPUT_PRICE_1 — 가격 ~. 빈 값 OK
	InputPrice2   string // FID_INPUT_PRICE_2 — ~ 가격. 빈 값 OK
	VolCount      string // FID_VOL_CNT — 거래량 ~. 빈 값 OK
	InputDate1    string // FID_INPUT_DATE_1 — 빈 값 OK
}

// InquireVolumeRank 는 거래량순위 호출.
//
// 한투 docs: docs/api/국내주식/거래량순위.md
// path: /uapi/domestic-stock/v1/quotations/volume-rank (FHPST01710000)
func (c *Client) InquireVolumeRank(ctx context.Context, params InquireVolumeRankParams) (*VolumeRank, error) {
	market := params.MarketCode
	if market == "" {
		market = "J"
	}
	scr := params.ScreenCode
	if scr == "" {
		scr = "20171"
	}
	div := params.DivCode
	if div == "" {
		div = "0"
	}
	belong := params.BelongCode
	if belong == "" {
		belong = "0"
	}
	tgt := params.TargetCode
	if tgt == "" {
		tgt = "111111111"
	}
	tgtExcl := params.TargetExclude
	if tgtExcl == "" {
		tgtExcl = "0000000000"
	}

	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/domestic-stock/v1/quotations/volume-rank",
		TrID:   "FHPST01710000",
		Query: map[string]string{
			"FID_COND_MRKT_DIV_CODE": market,
			"FID_COND_SCR_DIV_CODE":  scr,
			"FID_INPUT_ISCD":         params.InputISCD,
			"FID_DIV_CLS_CODE":       div,
			"FID_BLNG_CLS_CODE":      belong,
			"FID_TRGT_CLS_CODE":      tgt,
			"FID_TRGT_EXLS_CLS_CODE": tgtExcl,
			"FID_INPUT_PRICE_1":      params.InputPrice1,
			"FID_INPUT_PRICE_2":      params.InputPrice2,
			"FID_VOL_CNT":            params.VolCount,
			"FID_INPUT_DATE_1":       params.InputDate1,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}

	// 'Output' 키가 대문자 (KIS docs 명시) — resp.Raw 로 unmarshal.
	var rank VolumeRank
	if err := json.Unmarshal(resp.Raw, &rank); err != nil {
		return nil, fmt.Errorf("kis: parse VolumeRank: %w", err)
	}
	return &rank, nil
}
```

- [ ] **Step 4: 테스트 실행 → PASS**

Run: `go test ./domestic/... -run InquireVolumeRank -v`
Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add domestic/ranking.go domestic/ranking_test.go
git commit -m "$(cat <<'EOF'
[feat] domestic — InquireVolumeRank (거래량순위, FHPST01710000)

VolumeRank + VolumeRankItem struct (19 필드) + InquireVolumeRankParams
(11 query, zero-value default). Output 키가 대문자라 resp.Raw 직접 unmarshal.

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 3: domestic/ranking.go — InquireFluctuation

**Files:**
- Modify: `domestic/ranking.go`
- Modify: `domestic/ranking_test.go`

- [ ] **Step 1: 테스트 추가** — `domestic/ranking_test.go` 끝에 함수 추가

```go
func TestClient_InquireFluctuation(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/ranking/fluctuation`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "fluctuation_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireFluctuation(context.Background(), domestic.InquireFluctuationParams{
		InputISCD: "0000",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, "J", capturedQuery.Get("fid_cond_mrkt_div_code"))
	assert.Equal(t, "20170", capturedQuery.Get("fid_cond_scr_div_code"))
	assert.Equal(t, "0000", capturedQuery.Get("fid_input_iscd"))
	assert.Equal(t, "0", capturedQuery.Get("fid_rank_sort_cls_code"))

	require.Len(t, res.Output, 1)
	assert.Equal(t, "005930", res.Output[0].StckShrnIscd)
	assert.Equal(t, "삼성전자", res.Output[0].HtsKorIsnm)
	assert.Equal(t, decimal.NewFromInt(75800), res.Output[0].StckPrpr)
	assert.Equal(t, decimal.NewFromInt(76200), res.Output[0].StckHgpr)
	assert.Equal(t, decimal.NewFromInt(75500), res.Output[0].StckLwpr)
	assert.Equal(t, "131542", res.Output[0].HgprHour)
}
```

- [ ] **Step 2: 테스트 실행 → FAIL**

Run: `go test ./domestic/... -run InquireFluctuation -v`
Expected: 컴파일 실패.

- [ ] **Step 3: 구현 추가** — `domestic/ranking.go` 끝에 추가

```go
// Fluctuation 은 등락률 순위 (FHPST01700000) 응답.
//
// 한투 docs: docs/api/국내주식/국내주식_등락률_순위.md
// path: /uapi/domestic-stock/v1/ranking/fluctuation
type Fluctuation struct {
	Output []FluctuationItem `json:"output"`
}

// FluctuationItem 은 등락률 순위 응답의 한 행.
type FluctuationItem struct {
	StckShrnIscd            string          `json:"stck_shrn_iscd"`                       // 주식 단축 종목코드
	DataRank                string          `json:"data_rank"`                            // 데이터 순위
	HtsKorIsnm              string          `json:"hts_kor_isnm"`                         // HTS 한글 종목명
	StckPrpr                decimal.Decimal `json:"stck_prpr"`                            // 주식 현재가
	PrdyVrss                decimal.Decimal `json:"prdy_vrss"`                            // 전일 대비
	PrdyVrssSign            string          `json:"prdy_vrss_sign"`                       // 전일 대비 부호
	PrdyCtrt                float64         `json:"prdy_ctrt,string"`                     // 전일 대비율
	AcmlVol                 int64           `json:"acml_vol,string"`                      // 누적 거래량
	StckHgpr                decimal.Decimal `json:"stck_hgpr"`                            // 주식 최고가
	HgprHour                string          `json:"hgpr_hour"`                            // 최고가 시간 (HHMMSS)
	AcmlHgprDate            string          `json:"acml_hgpr_date"`                       // 누적 최고가 일자
	StckLwpr                decimal.Decimal `json:"stck_lwpr"`                            // 주식 최저가
	LwprHour                string          `json:"lwpr_hour"`                            // 최저가 시간
	AcmlLwprDate            string          `json:"acml_lwpr_date"`                       // 누적 최저가 일자
	LwprVrssPrprRate        float64         `json:"lwpr_vrss_prpr_rate,string"`           // 최저가 대비 현재가 비율
	DsgtDateClprVrssPrprRate float64        `json:"dsgt_date_clpr_vrss_prpr_rate,string"` // 지정 일자 종가 대비 현재가 비율
	CnntAscnDynu            int64           `json:"cnnt_ascn_dynu,string"`                // 연속 상승 일수
	HgprVrssPrprRate        float64         `json:"hgpr_vrss_prpr_rate,string"`           // 최고가 대비 현재가 비율
	CnntDownDynu            int64           `json:"cnnt_down_dynu,string"`                // 연속 하락 일수
	OprcVrssPrprSign        string          `json:"oprc_vrss_prpr_sign"`                  // 시가 대비 현재가 부호
	OprcVrssPrpr            decimal.Decimal `json:"oprc_vrss_prpr"`                       // 시가 대비 현재가
	OprcVrssPrprRate        float64         `json:"oprc_vrss_prpr_rate,string"`           // 시가 대비 현재가 비율
	PrdRsfl                 decimal.Decimal `json:"prd_rsfl"`                             // 기간 등락
	PrdRsflRate             float64         `json:"prd_rsfl_rate,string"`                 // 기간 등락 비율
}

// InquireFluctuationParams 는 등락률 순위 조회 파라미터.
type InquireFluctuationParams struct {
	RsflRate2     string // fid_rsfl_rate2 — 등락 비율2 (~ 비율). 빈 값 OK
	MarketCode    string // fid_cond_mrkt_div_code — "J":KRX, "NX":NXT. 빈 값=>"J"
	ScreenCode    string // fid_cond_scr_div_code — Unique key. 빈 값=>"20170"
	InputISCD     string // fid_input_iscd — 필수, "0000"(전체)/"0001"(코스피)/"1001"(코스닥)/"2001"(코스피200)
	SortCode      string // fid_rank_sort_cls_code — "0":상승율순, "1":하락율순, "2":시가대비상승, "3":시가대비하락, "4":변동율. 빈 값=>"0"
	InputCnt1     string // fid_input_cnt_1 — "0":전체, 또는 누적일수. 빈 값=>"0"
	PriceCode     string // fid_prc_cls_code — 가격 구분. 빈 값=>"0"
	InputPrice1   string // fid_input_price_1
	InputPrice2   string // fid_input_price_2
	VolCount      string // fid_vol_cnt
	TargetCode    string // fid_trgt_cls_code. 빈 값=>"0"
	TargetExclude string // fid_trgt_exls_cls_code. 빈 값=>"0"
	DivCode       string // fid_div_cls_code. 빈 값=>"0"
	RsflRate1     string // fid_rsfl_rate1 — 등락 비율1 (비율 ~). 빈 값 OK
}

// InquireFluctuation 는 등락률 순위 호출.
//
// 한투 docs: docs/api/국내주식/국내주식_등락률_순위.md
// path: /uapi/domestic-stock/v1/ranking/fluctuation (FHPST01700000)
func (c *Client) InquireFluctuation(ctx context.Context, params InquireFluctuationParams) (*Fluctuation, error) {
	market := params.MarketCode
	if market == "" {
		market = "J"
	}
	scr := params.ScreenCode
	if scr == "" {
		scr = "20170"
	}
	sort := params.SortCode
	if sort == "" {
		sort = "0"
	}
	cnt := params.InputCnt1
	if cnt == "" {
		cnt = "0"
	}
	prc := params.PriceCode
	if prc == "" {
		prc = "0"
	}
	tgt := params.TargetCode
	if tgt == "" {
		tgt = "0"
	}
	tgtExcl := params.TargetExclude
	if tgtExcl == "" {
		tgtExcl = "0"
	}
	div := params.DivCode
	if div == "" {
		div = "0"
	}

	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/domestic-stock/v1/ranking/fluctuation",
		TrID:   "FHPST01700000",
		Query: map[string]string{
			"fid_rsfl_rate2":         params.RsflRate2,
			"fid_cond_mrkt_div_code": market,
			"fid_cond_scr_div_code":  scr,
			"fid_input_iscd":         params.InputISCD,
			"fid_rank_sort_cls_code": sort,
			"fid_input_cnt_1":        cnt,
			"fid_prc_cls_code":       prc,
			"fid_input_price_1":      params.InputPrice1,
			"fid_input_price_2":      params.InputPrice2,
			"fid_vol_cnt":            params.VolCount,
			"fid_trgt_cls_code":      tgt,
			"fid_trgt_exls_cls_code": tgtExcl,
			"fid_div_cls_code":       div,
			"fid_rsfl_rate1":         params.RsflRate1,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}

	var res Fluctuation
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse Fluctuation: %w", err)
	}
	return &res, nil
}
```

- [ ] **Step 4: 테스트 실행 → PASS**

Run: `go test ./domestic/... -run InquireFluctuation -v`

- [ ] **Step 5: Commit**

```bash
git add domestic/ranking.go domestic/ranking_test.go
git commit -m "$(cat <<'EOF'
[feat] domestic — InquireFluctuation (등락률 순위, FHPST01700000)

Fluctuation + FluctuationItem (24 필드) + InquireFluctuationParams (14 query).
ranking 카테고리의 첫 메서드 — query naming 소문자 fid_*.

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 4: domestic/ranking.go — InquireMarketCap

**Files:**
- Modify: `domestic/ranking.go`
- Modify: `domestic/ranking_test.go`

- [ ] **Step 1: 테스트 추가**

```go
func TestClient_InquireMarketCap(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/ranking/market-cap`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "market_cap_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireMarketCap(context.Background(), domestic.InquireMarketCapParams{
		InputISCD: "0000",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, "J", capturedQuery.Get("fid_cond_mrkt_div_code"))
	assert.Equal(t, "20174", capturedQuery.Get("fid_cond_scr_div_code"))
	assert.Equal(t, "0000", capturedQuery.Get("fid_input_iscd"))

	require.Len(t, res.Output, 2)
	assert.Equal(t, "삼성전자", res.Output[0].HtsKorIsnm)
	assert.Equal(t, decimal.NewFromInt(75800), res.Output[0].StckPrpr)
	assert.Equal(t, int64(452329543), res.Output[0].StckAvls)
	assert.InDelta(t, 20.45, res.Output[0].MrktWholAvlsRlim, 0.001)
}
```

- [ ] **Step 2: 테스트 실행 → FAIL**

- [ ] **Step 3: 구현 추가** — `domestic/ranking.go` 끝에

```go
// MarketCap 은 시가총액 상위 (FHPST01740000) 응답.
//
// 한투 docs: docs/api/국내주식/국내주식_시가총액_상위.md
// path: /uapi/domestic-stock/v1/ranking/market-cap
type MarketCap struct {
	Output []MarketCapItem `json:"output"`
}

// MarketCapItem 은 시가총액 상위 응답의 한 행.
type MarketCapItem struct {
	MkscShrnIscd     string          `json:"mksc_shrn_iscd"`            // 유가증권 단축 종목코드
	DataRank         string          `json:"data_rank"`                 // 데이터 순위
	HtsKorIsnm       string          `json:"hts_kor_isnm"`              // HTS 한글 종목명
	StckPrpr         decimal.Decimal `json:"stck_prpr"`                 // 주식 현재가
	PrdyVrss         decimal.Decimal `json:"prdy_vrss"`                 // 전일 대비
	PrdyVrssSign     string          `json:"prdy_vrss_sign"`            // 전일 대비 부호
	PrdyCtrt         float64         `json:"prdy_ctrt,string"`          // 전일 대비율
	AcmlVol          int64           `json:"acml_vol,string"`           // 누적 거래량
	LstnStcn         int64           `json:"lstn_stcn,string"`          // 상장 주수
	StckAvls         int64           `json:"stck_avls,string"`          // 시가 총액
	MrktWholAvlsRlim float64         `json:"mrkt_whol_avls_rlim,string"` // 시장 전체 시가총액 비중
}

// InquireMarketCapParams 는 시가총액 상위 조회 파라미터.
type InquireMarketCapParams struct {
	InputPrice2   string // fid_input_price_2
	MarketCode    string // fid_cond_mrkt_div_code — 빈 값=>"J"
	ScreenCode    string // fid_cond_scr_div_code — 빈 값=>"20174"
	DivCode       string // fid_div_cls_code — "0":전체, "1":보통주, "2":우선주. 빈 값=>"0"
	InputISCD     string // fid_input_iscd — 필수, "0000"(전체)/"0001"(거래소)/"1001"(코스닥)/"2001"(코스피200)
	TargetCode    string // fid_trgt_cls_code — 빈 값=>"0"
	TargetExclude string // fid_trgt_exls_cls_code — 빈 값=>"0"
	InputPrice1   string // fid_input_price_1
	VolCount      string // fid_vol_cnt
}

// InquireMarketCap 은 시가총액 상위 호출.
//
// 한투 docs: docs/api/국내주식/국내주식_시가총액_상위.md
// path: /uapi/domestic-stock/v1/ranking/market-cap (FHPST01740000)
func (c *Client) InquireMarketCap(ctx context.Context, params InquireMarketCapParams) (*MarketCap, error) {
	market := params.MarketCode
	if market == "" {
		market = "J"
	}
	scr := params.ScreenCode
	if scr == "" {
		scr = "20174"
	}
	div := params.DivCode
	if div == "" {
		div = "0"
	}
	tgt := params.TargetCode
	if tgt == "" {
		tgt = "0"
	}
	tgtExcl := params.TargetExclude
	if tgtExcl == "" {
		tgtExcl = "0"
	}

	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/domestic-stock/v1/ranking/market-cap",
		TrID:   "FHPST01740000",
		Query: map[string]string{
			"fid_input_price_2":      params.InputPrice2,
			"fid_cond_mrkt_div_code": market,
			"fid_cond_scr_div_code":  scr,
			"fid_div_cls_code":       div,
			"fid_input_iscd":         params.InputISCD,
			"fid_trgt_cls_code":      tgt,
			"fid_trgt_exls_cls_code": tgtExcl,
			"fid_input_price_1":      params.InputPrice1,
			"fid_vol_cnt":            params.VolCount,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}

	var res MarketCap
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse MarketCap: %w", err)
	}
	return &res, nil
}
```

- [ ] **Step 4: 테스트 실행 → PASS**

- [ ] **Step 5: Commit**

```bash
git add domestic/ranking.go domestic/ranking_test.go
git commit -m "$(cat <<'EOF'
[feat] domestic — InquireMarketCap (시가총액 상위, FHPST01740000)

MarketCap + MarketCapItem (11 필드) + InquireMarketCapParams (9 query).

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 5: domestic/ranking.go — InquireDividendRate

**Files:**
- Modify: `domestic/ranking.go`
- Modify: `domestic/ranking_test.go`

> 배당률 순위는 다른 ranking 과 query 형식이 완전히 다름 (CTS_AREA, GB1, UPJONG, F_DT 등). TR_ID 도 `HHKDB13470100` 으로 다른 group. 응답은 `output1` (출력1).

- [ ] **Step 1: 테스트 추가**

```go
func TestClient_InquireDividendRate(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/ranking/dividend-rate`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "dividend_rate_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireDividendRate(context.Background(), domestic.InquireDividendRateParams{
		Sector:   "0001",
		FromDate: "20250101",
		ToDate:   "20251231",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	// 기본 query 검증
	assert.Equal(t, "0", capturedQuery.Get("GB1"))
	assert.Equal(t, "0001", capturedQuery.Get("UPJONG"))
	assert.Equal(t, "0", capturedQuery.Get("GB2"))
	assert.Equal(t, "1", capturedQuery.Get("GB3"))
	assert.Equal(t, "20250101", capturedQuery.Get("F_DT"))
	assert.Equal(t, "20251231", capturedQuery.Get("T_DT"))
	assert.Equal(t, "0", capturedQuery.Get("GB4"))

	require.Len(t, res.Output1, 2)
	assert.Equal(t, "1", res.Output1[0].Rank)
	assert.Equal(t, "005930", res.Output1[0].ShtCd)
	assert.Equal(t, "삼성전자", res.Output1[0].IsinName)
	assert.Equal(t, "20251231", res.Output1[0].RecordDate)
	assert.Equal(t, decimal.NewFromInt(1444), res.Output1[0].PerStoDiviAmt)
	assert.InDelta(t, 1.91, res.Output1[0].DiviRate, 0.001)
	assert.Equal(t, "현금배당", res.Output1[0].DiviKind)
}
```

- [ ] **Step 2: 테스트 실행 → FAIL**

- [ ] **Step 3: 구현 추가** — `domestic/ranking.go` 끝에

```go
// DividendRate 는 배당률 상위 (HHKDB13470100) 응답.
//
// 한투 docs: docs/api/국내주식/국내주식_배당률_상위.md
// path: /uapi/domestic-stock/v1/ranking/dividend-rate
type DividendRate struct {
	Output1 []DividendRateItem `json:"output1"` // 응답상세 (output1)
}

// DividendRateItem 은 배당률 상위 응답의 한 행.
type DividendRateItem struct {
	Rank          string          `json:"rank"`             // 순위
	ShtCd         string          `json:"sht_cd"`           // 종목코드
	IsinName      string          `json:"isin_name"`        // 종목명
	RecordDate    string          `json:"record_date"`      // 기준일 (YYYYMMDD)
	PerStoDiviAmt decimal.Decimal `json:"per_sto_divi_amt"` // 현금/주식배당금
	DiviRate      float64         `json:"divi_rate,string"` // 현금/주식배당률 (%)
	DiviKind      string          `json:"divi_kind"`        // 배당종류
}

// InquireDividendRateParams 는 배당률 상위 조회 파라미터.
//
// 다른 ranking 과 query 형식이 다름. KIS docs 의 query 키 (CTS_AREA, GB1~GB4, UPJONG, F_DT, T_DT) 그대로 노출.
type InquireDividendRateParams struct {
	CtsArea     string // CTS_AREA — 빈 값(공백) default
	Market      string // GB1 — KOSPI 구분: "0":전체, "1":코스피, "2":코스피200, "3":코스닥. 빈 값=>"0"
	Sector      string // UPJONG — 업종구분 (필수). 예: "0001"(코스피 종합), "1001"(코스닥 종합)
	StockType   string // GB2 — 종목선택: "0":전체, "6":보통주, "7":우선주. 빈 값=>"0"
	DividendCls string // GB3 — 배당구분: "1":주식배당, "2":현금배당. 빈 값=>"1"
	FromDate    string // F_DT — 기준일 From (필수, YYYYMMDD)
	ToDate      string // T_DT — 기준일 To (필수, YYYYMMDD)
	YearCls     string // GB4 — 결산/중간배당: "0":전체, "1":결산배당, "2":중간배당. 빈 값=>"0"
}

// InquireDividendRate 는 배당률 상위 호출.
//
// 한투 docs: docs/api/국내주식/국내주식_배당률_상위.md
// path: /uapi/domestic-stock/v1/ranking/dividend-rate (HHKDB13470100)
func (c *Client) InquireDividendRate(ctx context.Context, params InquireDividendRateParams) (*DividendRate, error) {
	gb1 := params.Market
	if gb1 == "" {
		gb1 = "0"
	}
	gb2 := params.StockType
	if gb2 == "" {
		gb2 = "0"
	}
	gb3 := params.DividendCls
	if gb3 == "" {
		gb3 = "1"
	}
	gb4 := params.YearCls
	if gb4 == "" {
		gb4 = "0"
	}

	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/domestic-stock/v1/ranking/dividend-rate",
		TrID:   "HHKDB13470100",
		Query: map[string]string{
			"CTS_AREA": params.CtsArea,
			"GB1":      gb1,
			"UPJONG":   params.Sector,
			"GB2":      gb2,
			"GB3":      gb3,
			"F_DT":     params.FromDate,
			"T_DT":     params.ToDate,
			"GB4":      gb4,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}

	var res DividendRate
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse DividendRate: %w", err)
	}
	return &res, nil
}
```

- [ ] **Step 4: 테스트 실행 → PASS**

- [ ] **Step 5: Commit**

```bash
git add domestic/ranking.go domestic/ranking_test.go
git commit -m "$(cat <<'EOF'
[feat] domestic — InquireDividendRate (배당률 상위, HHKDB13470100)

DividendRate + DividendRateItem (7 필드) + InquireDividendRateParams (8 query).
다른 ranking 과 다르게 query 키 (CTS_AREA, GB1~GB4, UPJONG, F_DT, T_DT) 그대로 노출.

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 6: domestic/financial.go — InquireFinancialRatio + 공통 base

**Files:**
- Create: `domestic/financial.go`
- Create: `domestic/financial_test.go`

> 5 financial 메서드 모두 query 가 거의 동일 (`FID_DIV_CLS_CODE`, `fid_cond_mrkt_div_code`, `fid_input_iscd`). 첫 메서드에 file 생성, 나머지는 추가만.

- [ ] **Step 1: 테스트 작성** — `domestic/financial_test.go`

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

func TestClient_InquireFinancialRatio(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/finance/financial-ratio`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "financial_ratio_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireFinancialRatio(context.Background(), domestic.InquireFinancialRatioParams{
		Symbol: "005930",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, "0", capturedQuery.Get("FID_DIV_CLS_CODE"))
	assert.Equal(t, "J", capturedQuery.Get("fid_cond_mrkt_div_code"))
	assert.Equal(t, "005930", capturedQuery.Get("fid_input_iscd"))

	require.Len(t, res.Output, 2)
	assert.Equal(t, "202412", res.Output[0].StacYymm)
	assert.InDelta(t, 5.42, res.Output[0].Grs, 0.001)
	assert.InDelta(t, 11.50, res.Output[0].RoeVal, 0.001)
	assert.Equal(t, decimal.NewFromInt(6638), res.Output[0].Eps)
	assert.Equal(t, decimal.NewFromInt(57420), res.Output[0].Bps)
}

func TestClient_InquireFinancialRatio_Quarter(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/finance/financial-ratio`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "financial_ratio_success.json")), nil
		},
	)

	c := newTestClient(t)
	_, err := c.InquireFinancialRatio(context.Background(), domestic.InquireFinancialRatioParams{
		Symbol:  "005930",
		Quarter: true,
	})
	require.NoError(t, err)
	assert.Equal(t, "1", capturedQuery.Get("FID_DIV_CLS_CODE")) // 1=분기
}
```

- [ ] **Step 2: 테스트 실행 → FAIL**

Run: `go test ./domestic/... -run InquireFinancialRatio -v`
Expected: 컴파일 실패.

- [ ] **Step 3: 구현** — `domestic/financial.go`

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

// 5 financial 메서드 모두 동일 query (Symbol + 분기/연도) — 공통 helper.

// inquireFinanceQuery 는 finance category 의 공통 query 생성.
//
// 모든 finance API 가 fid_cond_mrkt_div_code="J" 고정, FID_DIV_CLS_CODE 는 0(년)/1(분기).
func inquireFinanceQuery(symbol string, quarter bool) map[string]string {
	div := "0"
	if quarter {
		div = "1"
	}
	return map[string]string{
		"FID_DIV_CLS_CODE":       div,
		"fid_cond_mrkt_div_code": "J",
		"fid_input_iscd":         symbol,
	}
}

// FinancialRatio 는 재무비율 (FHKST66430300) 응답.
//
// 한투 docs: docs/api/국내주식/국내주식_재무비율.md
// path: /uapi/domestic-stock/v1/finance/financial-ratio
type FinancialRatio struct {
	Output []FinancialRatioItem `json:"output"`
}

// FinancialRatioItem 은 재무비율 응답의 한 행 (분기/년).
type FinancialRatioItem struct {
	StacYymm     string          `json:"stac_yymm"`              // 결산 년월 (YYYYMM)
	Grs          float64         `json:"grs,string"`             // 매출액 증가율
	BsopPrfiInrt float64         `json:"bsop_prfi_inrt,string"`  // 영업 이익 증가율 (적자지속/흑자전환/적자전환은 0)
	NtinInrt     float64         `json:"ntin_inrt,string"`       // 순이익 증가율
	RoeVal       float64         `json:"roe_val,string"`         // ROE 값
	Eps          decimal.Decimal `json:"eps"`                    // EPS
	Sps          decimal.Decimal `json:"sps"`                    // 주당매출액
	Bps          decimal.Decimal `json:"bps"`                    // BPS
	RsrvRate     float64         `json:"rsrv_rate,string"`       // 유보 비율
	LbltRate     float64         `json:"lblt_rate,string"`       // 부채 비율
}

// InquireFinancialRatioParams 는 재무비율 조회 파라미터.
//
// 5 financial 메서드 공통: Symbol (필수) + Quarter (false=년 default, true=분기).
type InquireFinancialRatioParams struct {
	Symbol  string // fid_input_iscd (필수, 종목코드)
	Quarter bool   // FID_DIV_CLS_CODE — false=>"0"(년 default), true=>"1"(분기)
}

// InquireFinancialRatio 는 재무비율 호출.
//
// 한투 docs: docs/api/국내주식/국내주식_재무비율.md
// path: /uapi/domestic-stock/v1/finance/financial-ratio (FHKST66430300)
func (c *Client) InquireFinancialRatio(ctx context.Context, params InquireFinancialRatioParams) (*FinancialRatio, error) {
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method:   http.MethodGet,
		Path:     "/uapi/domestic-stock/v1/finance/financial-ratio",
		TrID:     "FHKST66430300",
		Query:    inquireFinanceQuery(params.Symbol, params.Quarter),
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}

	var res FinancialRatio
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse FinancialRatio: %w", err)
	}
	return &res, nil
}
```

- [ ] **Step 4: 테스트 실행 → PASS**

Run: `go test ./domestic/... -run InquireFinancialRatio -v`
Expected: 2 PASS (TestClient_InquireFinancialRatio + _Quarter).

- [ ] **Step 5: Commit**

```bash
git add domestic/financial.go domestic/financial_test.go
git commit -m "$(cat <<'EOF'
[feat] domestic — InquireFinancialRatio (재무비율, FHKST66430300)

FinancialRatio + FinancialRatioItem (10 필드) + InquireFinancialRatioParams.
inquireFinanceQuery helper 도 추가 — 5 financial 메서드 공통 query 생성.
Quarter bool 로 년/분기 구분 (false=년 default).

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 7: domestic/financial.go — InquireIncomeStatement

**Files:**
- Modify: `domestic/financial.go`
- Modify: `domestic/financial_test.go`

- [ ] **Step 1: 테스트 추가**

```go
func TestClient_InquireIncomeStatement(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/finance/income-statement`,
		httpmock.NewStringResponder(200, loadFixtureString(t, "income_statement_success.json")),
	)

	c := newTestClient(t)
	res, err := c.InquireIncomeStatement(context.Background(), domestic.InquireIncomeStatementParams{
		Symbol: "005930",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	require.Len(t, res.Output, 1)
	assert.Equal(t, "202412", res.Output[0].StacYymm)
	assert.Equal(t, int64(279600000), res.Output[0].SaleAccount)
	assert.Equal(t, int64(176000000), res.Output[0].SaleCost)
	assert.Equal(t, int64(32830000), res.Output[0].BsopPrti)
	assert.Equal(t, int64(23456000), res.Output[0].ThtrNtin)
}
```

- [ ] **Step 2: 테스트 실행 → FAIL**

- [ ] **Step 3: 구현 추가** — `domestic/financial.go` 끝에

```go
// IncomeStatement 는 손익계산서 (FHKST66430200) 응답.
//
// 한투 docs: docs/api/국내주식/국내주식_손익계산서.md
// path: /uapi/domestic-stock/v1/finance/income-statement
//
// 분기 데이터는 연단위 누적 합산.
type IncomeStatement struct {
	Output []IncomeStatementItem `json:"output"`
}

// IncomeStatementItem 은 손익계산서 응답의 한 행.
type IncomeStatementItem struct {
	StacYymm      string `json:"stac_yymm"`             // 결산 년월
	SaleAccount   int64  `json:"sale_account,string"`   // 매출액
	SaleCost      int64  `json:"sale_cost,string"`      // 매출 원가
	SaleTotlPrfi  int64  `json:"sale_totl_prfi,string"` // 매출 총 이익
	DeprCost      string `json:"depr_cost"`             // 감가상각비 (출력 안 되면 "99.99" — string 그대로)
	SellMang      string `json:"sell_mang"`             // 판매 및 관리비 (출력 안 되면 "99.99")
	BsopPrti      int64  `json:"bsop_prti,string"`      // 영업 이익
	BsopNonErnn   string `json:"bsop_non_ernn"`         // 영업 외 수익 (출력 안 되면 "99.99")
	BsopNonExpn   string `json:"bsop_non_expn"`         // 영업 외 비용 (출력 안 되면 "99.99")
	OpPrfi        int64  `json:"op_prfi,string"`        // 경상 이익
	SpecPrfi      int64  `json:"spec_prfi,string"`      // 특별 이익
	SpecLoss      int64  `json:"spec_loss,string"`      // 특별 손실
	ThtrNtin      int64  `json:"thtr_ntin,string"`      // 당기순이익
}

// InquireIncomeStatementParams 는 손익계산서 조회 파라미터.
type InquireIncomeStatementParams struct {
	Symbol  string // fid_input_iscd (필수)
	Quarter bool   // FID_DIV_CLS_CODE — false=>년, true=>분기 (분기는 누적 합산)
}

// InquireIncomeStatement 는 손익계산서 호출.
//
// 한투 docs: docs/api/국내주식/국내주식_손익계산서.md
// path: /uapi/domestic-stock/v1/finance/income-statement (FHKST66430200)
//
// ※ 분기 데이터는 연단위 누적 합산. depr_cost / sell_mang / bsop_non_ernn / bsop_non_expn
// 은 출력 안 되면 "99.99" — caller 는 string 으로 받아 "99.99" 검사 후 처리.
func (c *Client) InquireIncomeStatement(ctx context.Context, params InquireIncomeStatementParams) (*IncomeStatement, error) {
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method:   http.MethodGet,
		Path:     "/uapi/domestic-stock/v1/finance/income-statement",
		TrID:     "FHKST66430200",
		Query:    inquireFinanceQuery(params.Symbol, params.Quarter),
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}

	var res IncomeStatement
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse IncomeStatement: %w", err)
	}
	return &res, nil
}
```

- [ ] **Step 4: 테스트 실행 → PASS**

- [ ] **Step 5: Commit**

```bash
git add domestic/financial.go domestic/financial_test.go
git commit -m "$(cat <<'EOF'
[feat] domestic — InquireIncomeStatement (손익계산서, FHKST66430200)

IncomeStatement + IncomeStatementItem (13 필드) + InquireIncomeStatementParams.
depr_cost/sell_mang/bsop_non_*는 출력 안 될 시 "99.99" — string 으로 노출.

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 8: domestic/financial.go — InquireBalanceSheet

**Files:**
- Modify: `domestic/financial.go`
- Modify: `domestic/financial_test.go`

- [ ] **Step 1: 테스트 추가**

```go
func TestClient_InquireBalanceSheet(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/finance/balance-sheet`,
		httpmock.NewStringResponder(200, loadFixtureString(t, "balance_sheet_success.json")),
	)

	c := newTestClient(t)
	res, err := c.InquireBalanceSheet(context.Background(), domestic.InquireBalanceSheetParams{
		Symbol: "005930",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	require.Len(t, res.Output, 1)
	assert.Equal(t, "202412", res.Output[0].StacYymm)
	assert.Equal(t, int64(189000000), res.Output[0].Cras)
	assert.Equal(t, int64(434000000), res.Output[0].TotalAset)
	assert.Equal(t, int64(94000000), res.Output[0].TotalLblt)
	assert.Equal(t, int64(340000000), res.Output[0].TotalCptl)
}
```

- [ ] **Step 2: 테스트 실행 → FAIL**

- [ ] **Step 3: 구현 추가**

```go
// BalanceSheet 는 대차대조표 (FHKST66430100) 응답.
//
// 한투 docs: docs/api/국내주식/국내주식_대차대조표.md
// path: /uapi/domestic-stock/v1/finance/balance-sheet
type BalanceSheet struct {
	Output []BalanceSheetItem `json:"output"`
}

// BalanceSheetItem 은 대차대조표 응답의 한 행.
type BalanceSheetItem struct {
	StacYymm   string `json:"stac_yymm"`         // 결산 년월
	Cras       int64  `json:"cras,string"`       // 유동자산
	Fxas       int64  `json:"fxas,string"`       // 고정자산
	TotalAset  int64  `json:"total_aset,string"` // 자산총계
	FlowLblt   int64  `json:"flow_lblt,string"`  // 유동부채
	FixLblt    int64  `json:"fix_lblt,string"`   // 고정부채
	TotalLblt  int64  `json:"total_lblt,string"` // 부채총계
	Cpfn       int64  `json:"cpfn,string"`       // 자본금
	CfpSurp    string `json:"cfp_surp"`          // 자본 잉여금 (출력 안 되면 "99.99")
	PrfiSurp   string `json:"prfi_surp"`         // 이익 잉여금 (출력 안 되면 "99.99")
	TotalCptl  int64  `json:"total_cptl,string"` // 자본총계
}

// InquireBalanceSheetParams 는 대차대조표 조회 파라미터.
type InquireBalanceSheetParams struct {
	Symbol  string // fid_input_iscd (필수)
	Quarter bool   // FID_DIV_CLS_CODE
}

// InquireBalanceSheet 는 대차대조표 호출.
//
// 한투 docs: docs/api/국내주식/국내주식_대차대조표.md
// path: /uapi/domestic-stock/v1/finance/balance-sheet (FHKST66430100)
func (c *Client) InquireBalanceSheet(ctx context.Context, params InquireBalanceSheetParams) (*BalanceSheet, error) {
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method:   http.MethodGet,
		Path:     "/uapi/domestic-stock/v1/finance/balance-sheet",
		TrID:     "FHKST66430100",
		Query:    inquireFinanceQuery(params.Symbol, params.Quarter),
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}

	var res BalanceSheet
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse BalanceSheet: %w", err)
	}
	return &res, nil
}
```

- [ ] **Step 4: 테스트 실행 → PASS**

- [ ] **Step 5: Commit**

```bash
git add domestic/financial.go domestic/financial_test.go
git commit -m "$(cat <<'EOF'
[feat] domestic — InquireBalanceSheet (대차대조표, FHKST66430100)

BalanceSheet + BalanceSheetItem (11 필드) + InquireBalanceSheetParams.

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 9: domestic/financial.go — InquireProfitRatio

**Files:**
- Modify: `domestic/financial.go`
- Modify: `domestic/financial_test.go`

- [ ] **Step 1: 테스트 추가**

```go
func TestClient_InquireProfitRatio(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/finance/profit-ratio`,
		httpmock.NewStringResponder(200, loadFixtureString(t, "profit_ratio_success.json")),
	)

	c := newTestClient(t)
	res, err := c.InquireProfitRatio(context.Background(), domestic.InquireProfitRatioParams{
		Symbol: "005930",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	require.Len(t, res.Output, 1)
	assert.Equal(t, "202412", res.Output[0].StacYymm)
	assert.InDelta(t, 8.45, res.Output[0].CptlNtinRate, 0.001)
	assert.InDelta(t, 11.50, res.Output[0].SelfCptlNtinInrt, 0.001)
	assert.InDelta(t, 12.30, res.Output[0].SaleNtinRate, 0.001)
	assert.InDelta(t, 37.05, res.Output[0].SaleTotlRate, 0.001)
}
```

- [ ] **Step 2: 테스트 실행 → FAIL**

- [ ] **Step 3: 구현 추가**

```go
// ProfitRatio 는 수익성비율 (FHKST66430400) 응답.
//
// 한투 docs: docs/api/국내주식/국내주식_수익성비율.md
// path: /uapi/domestic-stock/v1/finance/profit-ratio
type ProfitRatio struct {
	Output []ProfitRatioItem `json:"output"`
}

// ProfitRatioItem 은 수익성비율 응답의 한 행.
type ProfitRatioItem struct {
	StacYymm         string  `json:"stac_yymm"`                    // 결산 년월
	CptlNtinRate     float64 `json:"cptl_ntin_rate,string"`        // 총자본 순이익율
	SelfCptlNtinInrt float64 `json:"self_cptl_ntin_inrt,string"`   // 자기자본 순이익율
	SaleNtinRate     float64 `json:"sale_ntin_rate,string"`        // 매출액 순이익율
	SaleTotlRate     float64 `json:"sale_totl_rate,string"`        // 매출액 총이익율
}

// InquireProfitRatioParams 는 수익성비율 조회 파라미터.
type InquireProfitRatioParams struct {
	Symbol  string // fid_input_iscd (필수)
	Quarter bool   // FID_DIV_CLS_CODE
}

// InquireProfitRatio 는 수익성비율 호출.
//
// 한투 docs: docs/api/국내주식/국내주식_수익성비율.md
// path: /uapi/domestic-stock/v1/finance/profit-ratio (FHKST66430400)
func (c *Client) InquireProfitRatio(ctx context.Context, params InquireProfitRatioParams) (*ProfitRatio, error) {
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method:   http.MethodGet,
		Path:     "/uapi/domestic-stock/v1/finance/profit-ratio",
		TrID:     "FHKST66430400",
		Query:    inquireFinanceQuery(params.Symbol, params.Quarter),
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}

	var res ProfitRatio
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse ProfitRatio: %w", err)
	}
	return &res, nil
}
```

- [ ] **Step 4: 테스트 실행 → PASS**

- [ ] **Step 5: Commit**

```bash
git add domestic/financial.go domestic/financial_test.go
git commit -m "$(cat <<'EOF'
[feat] domestic — InquireProfitRatio (수익성비율, FHKST66430400)

ProfitRatio + ProfitRatioItem (5 필드) + InquireProfitRatioParams.

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 10: domestic/financial.go — InquireGrowthRatio

**Files:**
- Modify: `domestic/financial.go`
- Modify: `domestic/financial_test.go`

- [ ] **Step 1: 테스트 추가**

```go
func TestClient_InquireGrowthRatio(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/finance/growth-ratio`,
		httpmock.NewStringResponder(200, loadFixtureString(t, "growth_ratio_success.json")),
	)

	c := newTestClient(t)
	res, err := c.InquireGrowthRatio(context.Background(), domestic.InquireGrowthRatioParams{
		Symbol: "005930",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	require.Len(t, res.Output, 1)
	assert.Equal(t, "202412", res.Output[0].StacYymm)
	assert.InDelta(t, 5.42, res.Output[0].Grs, 0.001)
	assert.InDelta(t, 12.30, res.Output[0].BsopPrfiInrt, 0.001)
	assert.InDelta(t, 8.50, res.Output[0].EqutInrt, 0.001)
	assert.InDelta(t, 10.20, res.Output[0].TotlAsetInrt, 0.001)
}
```

- [ ] **Step 2: 테스트 실행 → FAIL**

- [ ] **Step 3: 구현 추가**

```go
// GrowthRatio 는 성장성비율 (FHKST66430800) 응답.
//
// 한투 docs: docs/api/국내주식/국내주식_성장성비율.md
// path: /uapi/domestic-stock/v1/finance/growth-ratio
type GrowthRatio struct {
	Output []GrowthRatioItem `json:"output"`
}

// GrowthRatioItem 은 성장성비율 응답의 한 행.
type GrowthRatioItem struct {
	StacYymm     string  `json:"stac_yymm"`              // 결산 년월
	Grs          float64 `json:"grs,string"`             // 매출액 증가율
	BsopPrfiInrt float64 `json:"bsop_prfi_inrt,string"`  // 영업 이익 증가율
	EqutInrt     float64 `json:"equt_inrt,string"`       // 자기자본 증가율
	TotlAsetInrt float64 `json:"totl_aset_inrt,string"`  // 총자산 증가율
}

// InquireGrowthRatioParams 는 성장성비율 조회 파라미터.
type InquireGrowthRatioParams struct {
	Symbol  string // fid_input_iscd (필수)
	Quarter bool   // fid_div_cls_code (소문자 — KIS docs 그대로). false=>"0"(년), true=>"1"(분기)
}

// InquireGrowthRatio 는 성장성비율 호출.
//
// 한투 docs: docs/api/국내주식/국내주식_성장성비율.md
// path: /uapi/domestic-stock/v1/finance/growth-ratio (FHKST66430800)
//
// 다른 finance API 와 다르게 query 키가 fid_div_cls_code (소문자) — 다른 4 메서드는 FID_DIV_CLS_CODE (대문자).
// 그래서 inquireFinanceQuery helper 는 사용하지 않고 inline query.
func (c *Client) InquireGrowthRatio(ctx context.Context, params InquireGrowthRatioParams) (*GrowthRatio, error) {
	div := "0"
	if params.Quarter {
		div = "1"
	}
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/domestic-stock/v1/finance/growth-ratio",
		TrID:   "FHKST66430800",
		Query: map[string]string{
			"fid_input_iscd":         params.Symbol,
			"fid_div_cls_code":       div,
			"fid_cond_mrkt_div_code": "J",
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}

	var res GrowthRatio
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse GrowthRatio: %w", err)
	}
	return &res, nil
}
```

- [ ] **Step 4: 테스트 실행 → PASS**

- [ ] **Step 5: 전체 회귀 검증**

Run: `go test ./... -count=1`
Expected: 모든 패키지 PASS.

- [ ] **Step 6: Commit**

```bash
git add domestic/financial.go domestic/financial_test.go
git commit -m "$(cat <<'EOF'
[feat] domestic — InquireGrowthRatio (성장성비율, FHKST66430800)

GrowthRatio + GrowthRatioItem (5 필드) + InquireGrowthRatioParams. KIS docs
가 다른 finance API 와 다르게 fid_div_cls_code (소문자) 사용 — inline query.

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 11: examples/domestic_ranking/main.go

**Files:**
- Create: `examples/domestic_ranking/main.go`

- [ ] **Step 1: example 작성**

```go
// domestic_ranking example: InquireVolumeRank + InquireFluctuation + InquireMarketCap.
//
// Run: KIS credentials env vars 후 go run ./examples/domestic_ranking
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

	// 1. 거래량 상위 30
	vol, err := client.Domestic.InquireVolumeRank(ctx, domestic.InquireVolumeRankParams{
		InputISCD: "0000",
	})
	if err != nil {
		log.Fatalf("InquireVolumeRank: %v", err)
	}
	fmt.Printf("거래량 상위 %d 개\n", len(vol.Output))
	for i, item := range vol.Output {
		if i >= 5 {
			fmt.Println("  ... (이하 생략)")
			break
		}
		fmt.Printf("  %s. %s (%s) 현재가=%s 거래량=%d\n",
			item.DataRank, item.HtsKorIsnm, item.MkscShrnIscd,
			item.StckPrpr, item.AcmlVol)
	}

	// 2. 등락률 상위 (상승율순)
	flux, err := client.Domestic.InquireFluctuation(ctx, domestic.InquireFluctuationParams{
		InputISCD: "0000",
		SortCode:  "0", // 0=상승율순
	})
	if err != nil {
		log.Fatalf("InquireFluctuation: %v", err)
	}
	fmt.Printf("\n등락률 상위 (상승) %d 개\n", len(flux.Output))
	for i, item := range flux.Output {
		if i >= 5 {
			fmt.Println("  ... (이하 생략)")
			break
		}
		fmt.Printf("  %s. %s (%s) %s원 (%v%%)\n",
			item.DataRank, item.HtsKorIsnm, item.StckShrnIscd,
			item.StckPrpr, item.PrdyCtrt)
	}

	// 3. 시가총액 상위
	cap, err := client.Domestic.InquireMarketCap(ctx, domestic.InquireMarketCapParams{
		InputISCD: "0000",
	})
	if err != nil {
		log.Fatalf("InquireMarketCap: %v", err)
	}
	fmt.Printf("\n시가총액 상위 %d 개\n", len(cap.Output))
	for i, item := range cap.Output {
		if i >= 5 {
			fmt.Println("  ... (이하 생략)")
			break
		}
		fmt.Printf("  %s. %s (%s) 시총=%d백만원 비중=%v%%\n",
			item.DataRank, item.HtsKorIsnm, item.MkscShrnIscd,
			item.StckAvls, item.MrktWholAvlsRlim)
	}
}
```

- [ ] **Step 2: 컴파일 검증**

Run: `go build ./examples/domestic_ranking && echo OK`
Expected: `OK`.

- [ ] **Step 3: Commit**

```bash
git add examples/domestic_ranking
git commit -m "$(cat <<'EOF'
[feat] examples/domestic_ranking — VolumeRank + Fluctuation + MarketCap

코스피 종합 거래량/등락률/시가총액 상위 5종목 출력 예시.

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 12: examples/domestic_financial/main.go

**Files:**
- Create: `examples/domestic_financial/main.go`

- [ ] **Step 1: example 작성**

```go
// domestic_financial example: InquireFinancialRatio + InquireIncomeStatement + InquireBalanceSheet.
//
// Run: KIS credentials env vars 후 go run ./examples/domestic_financial
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
	symbol := "005930" // 삼성전자

	// 1. 재무비율 (연단위)
	ratio, err := client.Domestic.InquireFinancialRatio(ctx, domestic.InquireFinancialRatioParams{
		Symbol: symbol,
	})
	if err != nil {
		log.Fatalf("InquireFinancialRatio: %v", err)
	}
	fmt.Printf("[%s] 재무비율 %d 기간\n", symbol, len(ratio.Output))
	for i, item := range ratio.Output {
		if i >= 3 {
			fmt.Println("  ... (이하 생략)")
			break
		}
		fmt.Printf("  %s: 매출증가=%v%% 영업증가=%v%% ROE=%v%% EPS=%s BPS=%s\n",
			item.StacYymm, item.Grs, item.BsopPrfiInrt, item.RoeVal,
			item.Eps, item.Bps)
	}

	// 2. 손익계산서
	is, err := client.Domestic.InquireIncomeStatement(ctx, domestic.InquireIncomeStatementParams{
		Symbol: symbol,
	})
	if err != nil {
		log.Fatalf("InquireIncomeStatement: %v", err)
	}
	fmt.Printf("\n[%s] 손익계산서 %d 기간\n", symbol, len(is.Output))
	for i, item := range is.Output {
		if i >= 3 {
			fmt.Println("  ... (이하 생략)")
			break
		}
		fmt.Printf("  %s: 매출=%d 영업이익=%d 당기순이익=%d (백만원)\n",
			item.StacYymm, item.SaleAccount, item.BsopPrti, item.ThtrNtin)
	}

	// 3. 대차대조표
	bs, err := client.Domestic.InquireBalanceSheet(ctx, domestic.InquireBalanceSheetParams{
		Symbol: symbol,
	})
	if err != nil {
		log.Fatalf("InquireBalanceSheet: %v", err)
	}
	fmt.Printf("\n[%s] 대차대조표 %d 기간\n", symbol, len(bs.Output))
	for i, item := range bs.Output {
		if i >= 3 {
			fmt.Println("  ... (이하 생략)")
			break
		}
		fmt.Printf("  %s: 자산=%d 부채=%d 자본=%d (백만원)\n",
			item.StacYymm, item.TotalAset, item.TotalLblt, item.TotalCptl)
	}
}
```

- [ ] **Step 2: 컴파일 검증**

Run: `go build ./examples/domestic_financial && echo OK`

- [ ] **Step 3: Commit**

```bash
git add examples/domestic_financial
git commit -m "$(cat <<'EOF'
[feat] examples/domestic_financial — FinancialRatio + IncomeStatement + BalanceSheet

삼성전자 (005930) 의 재무비율 / 손익계산서 / 대차대조표 출력 예시.

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 13: 문서 갱신 (CLAUDE.md, README.md, CHANGELOG.md, domestic/doc.go)

**Files:**
- Modify: `CLAUDE.md`
- Modify: `README.md`
- Modify: `CHANGELOG.md`
- Modify: `domestic/doc.go`

- [ ] **Step 1: README.md 의 Available Methods 표 갱신**

기존 표 (Phase 1.2 의 7 메서드) 에 9 개 추가:

```markdown
## Available Methods (Phase 1.2 + 1.3)

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
```

- [ ] **Step 2: CLAUDE.md 의 banner 갱신**

Phase 1.2 banner:
```
> **Phase 1.2 — domestic 시세/심볼/차트 (v0.2.0).** Phase 1.3+ 메서드는 추후 sub-plan 으로.
```

→
```
> **Phase 1.3 — domestic 순위/재무 (v1.1.0).** Phase 1.4+ 메서드는 추후 sub-plan 으로.
```

또한 spec 링크 section 에 Phase 1.3 plan 추가:
```markdown
- Phase 1.3 implementation plan: [`docs/superpowers/specs/2026-05-04-phase1-3-ranking-financial-implementation-plan.md`](docs/superpowers/specs/2026-05-04-phase1-3-ranking-financial-implementation-plan.md)
```

- [ ] **Step 3: CHANGELOG.md `[1.1.0]` entry 추가**

CHANGELOG.md 의 `[1.0.0]` section 위에 새 section 추가:

```markdown
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
```

- [ ] **Step 4: domestic/doc.go 갱신**

기존 7 메서드 + 9 메서드 → 16 메서드 list 로 갱신:

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
[doc] Phase 1.3 메서드 문서 갱신 — CLAUDE.md, README.md, CHANGELOG.md, domestic/doc.go

Phase 1.3 의 9 메서드 (ranking 4 + financial 5) 목록 + CHANGELOG [1.1.0] entry.
CLAUDE.md banner 갱신 (Phase 1.2 → 1.3, v0.2.0 → v1.1.0). domestic/doc.go 의
패키지 doc 에 Phase 1.3 section 추가.

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 14: 최종 점검

- [ ] **Step 1: 빌드/vet/fmt**

Run:
```bash
go build ./... && go vet ./... && gofmt -l . | tee /tmp/fmt.out
```
Expected: 빌드/vet 출력 없음. `gofmt -l` 출력 빈 파일.

- [ ] **Step 2: 모든 테스트 통과 + race**

Run: `go test ./... -race -count=1`
Expected: 모든 패키지 PASS.

- [ ] **Step 3: Coverage 측정**

Run:
```bash
go test ./... -coverprofile=/tmp/cov.out -covermode=atomic
go tool cover -func=/tmp/cov.out | tail -10
```
Expected: 마지막 줄 `total: (statements) ...` ≥ 80%. `domestic/` ≥ 80%.

만약 부족하면 다음 영역 추가 테스트 후 fix commit:
- ranking/financial 메서드의 API error 경로 (`rt_cd != "0"` 케이스)
- Params struct 의 zero-value default 동작 검증

- [ ] **Step 4: 디렉터리 구조 확인**

Run:
```bash
ls -la domestic/{ranking,financial}.go domestic/{ranking,financial}_test.go domestic/testdata/{volume_rank,fluctuation,market_cap,dividend_rate,financial_ratio,income_statement,balance_sheet,profit_ratio,growth_ratio}_success.json examples/{domestic_ranking,domestic_financial}/main.go
```
Expected: 모두 존재.

- [ ] **Step 5: Commit history 확인**

Run: `git log main..HEAD --oneline`
Expected: ~14-16 commits (각 task 의 commit + amendment).

(이 task 는 fix 가 필요하면 수행 후 commit; 이상 없으면 commit 없이 다음 단계로.)

---

## Task 15: PR 생성 및 push (사용자 승인 후)

> Claude 는 push / PR 생성을 사용자 명시적 승인 후에만 실행 (글로벌 정책).

- [ ] **Step 1: 사용자 승인 요청**

작업 진행 보고 + 사용자에게 PR 생성/push 진행 가능 여부 confirm.

- [ ] **Step 2: Push branch**

```bash
git push -u origin docs/phase1-3-spec
```

- [ ] **Step 3: PR 생성**

```bash
gh pr create --title "Phase 1.3 — domestic 순위/재무 (v1.1.0)" --reviewer kenshin579 --base main --head docs/phase1-3-spec --body "$(cat <<'EOF'
## Summary

- 국내주식 순위 4 메서드 + 재무 5 메서드 = 9 메서드 추가
- Phase 1.2 패턴 그대로 재사용 (Style A 메서드명, Params struct, KIS docs 1:1 매핑)
- 새 internal package 불필요 (REST only)
- v1.1.0 release 대상 — merge 후 태그 push

## 메서드 → 한투 API 매핑

| Go 메서드 | path | TR_ID |
|-----------|------|-------|
| InquireVolumeRank | quotations/volume-rank | FHPST01710000 |
| InquireFluctuation | ranking/fluctuation | FHPST01700000 |
| InquireMarketCap | ranking/market-cap | FHPST01740000 |
| InquireDividendRate | ranking/dividend-rate | HHKDB13470100 |
| InquireFinancialRatio | finance/financial-ratio | FHKST66430300 |
| InquireIncomeStatement | finance/income-statement | FHKST66430200 |
| InquireBalanceSheet | finance/balance-sheet | FHKST66430100 |
| InquireProfitRatio | finance/profit-ratio | FHKST66430400 |
| InquireGrowthRatio | finance/growth-ratio | FHKST66430800 |

## 주요 설계 결정

- ranking 메서드의 query parameter naming inconsistent (거래량순위만 대문자) — KIS docs 그대로
- 5 financial 메서드 공통 query helper (`inquireFinanceQuery`) — 단, GrowthRatio 만 query 키가 소문자라 inline
- 손익계산서/대차대조표 의 미출력 필드는 KIS 가 "99.99" 반환 — string 그대로 노출

## Test Plan

- [x] go build/vet/fmt clean
- [x] go test ./... -race -count=1 모든 패키지 PASS
- [x] Coverage domestic/ ≥ 80%
- [x] httpmock 단위 테스트 (각 메서드 + Quarter 옵션 변형)

## Breaking Changes

없음 — 신규 메서드 추가만.

## 참고 문서

- Phase 1 design spec (Phase 1.3 amendment 적용): docs/superpowers/specs/2026-05-03-phase1-api-coverage-design.md
- Phase 1.3 implementation plan: docs/superpowers/specs/2026-05-04-phase1-3-ranking-financial-implementation-plan.md

🤖 Generated with [Claude Code](https://claude.com/claude-code)
EOF
)"
```

- [ ] **Step 4: PR 검증**

`gh pr view <PR#> --json mergeable,statusCheckRollup` — 결과 확인.

- [ ] **Step 5: Merge (사용자 승인 후)**

`gh pr merge <PR#> --merge` — Phase 1.2 와 동일 방식 (merge commit).

- [ ] **Step 6: 후속 작업 (사용자 승인 후)**

```bash
git checkout main && git pull
git tag -a v1.1.0 -m "v1.1.0: Phase 1.3 — domestic 순위 + 재무 9 메서드"
git push origin v1.1.0
gh release create v1.1.0 --title "v1.1.0 — Phase 1.3: domestic 순위/재무" --notes-file <(awk '/^## \[/{c++} c==1' CHANGELOG.md)
```
