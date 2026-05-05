# Phase 2.2 — 국내 신고저가 / 시간외 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** `korea-investment-stock` Go 라이브러리에 국내주식 신고저 근접 + 시간외 5 메서드 추가 (`v1.5.0` release).

**Architecture:** Phase 2.1 의 인프라 + 패턴 재사용. `domestic/extended.go` (5 메서드 한 file) 추가. 한투 API path 1:1 매핑 (Style A). 새 internal package 불필요. TDD: testdata fixture (KIS docs 응답 필드 기반 합성 JSON) → 실패 테스트 → struct + 메서드 구현 → 통과 → commit.

**Tech Stack:** Go 1.25+, `github.com/jarcoal/httpmock`, `github.com/shopspring/decimal`, `github.com/stretchr/testify`. 새 dependency 없음.

**참고 spec:**
- Phase 2 design spec: `docs/superpowers/specs/2026-05-05-phase2-readonly-extension-design.md`
- Phase 2.1 plan (참조 패턴): `docs/superpowers/specs/2026-05-05-phase2-1-domestic-quote-implementation-plan.md`
- 한투 API docs: `docs/api/국내주식/{국내주식_신고_신저근접종목_상위.md, 국내주식_시간외현재가.md, 국내주식_시간외호가.md, 국내주식_시간외거래량순위.md, 국내주식_시간외등락율순위.md}`

---

## 사전 정보

| 항목 | 값 |
|------|---|
| 작업 브랜치 | `feat/phase2-2-extended` |
| 시작 HEAD | Phase 2.1 구현 완료 commit (v1.4.0) |
| Release 목표 | `v1.5.0` |
| PR 베이스 | `main` |
| 현재 main 상태 | v1.4.0 publish 완료 (Phase 2.1 통합, 31 메서드) |

---

## 메서드 → 한투 API 매핑

| Go 메서드 | 한투 path | TR_ID | 응답 구조 |
|---|---|---|---|
| `Domestic.InquireNearNewHighlow(ctx, params)` | `/uapi/domestic-stock/v1/ranking/near-new-highlow` | FHPST01870000 | `output: []` array (16 fields/item) |
| `Domestic.InquireOvertimePrice(ctx, params)` | `/uapi/domestic-stock/v1/quotations/inquire-overtime-price` | FHPST02300000 | `output: {}` object (35 fields) |
| `Domestic.InquireOvertimeAskingPrice(ctx, params)` | `/uapi/domestic-stock/v1/quotations/inquire-overtime-asking-price` | FHPST02300400 | `output1: {}` (74 fields) |
| `Domestic.InquireOvertimeVolume(ctx, params)` | `/uapi/domestic-stock/v1/ranking/overtime-volume` | FHPST02350000 | `output1: {}` (4 fields) + `output2: []` (14 fields/item) |
| `Domestic.InquireOvertimeFluctuation(ctx, params)` | `/uapi/domestic-stock/v1/ranking/overtime-fluctuation` | FHPST02340000 | `output1: {}` (11 fields) + `output2: []` (16 fields/item) |

> **명명 노트:** Phase 2 design spec §2 는 Phase 2.2 의 메서드 명을 `InquireNewHighLowProximity`, `InquireAfterHourPrice` 등 의미 기반으로 기술했으나, §3 에서 "Style A — path last segment PascalCase 1:1 매핑" 을 원칙으로 확정. 실제 path last segment (`near-new-highlow`, `inquire-overtime-price` 등) PascalCase 변환을 따름. 결과: `InquireNearNewHighlow`, `InquireOvertimePrice`, `InquireOvertimeAskingPrice`, `InquireOvertimeVolume`, `InquireOvertimeFluctuation`.

---

## 파일 구조

### 신규 (domestic)
- `domestic/extended.go` — 5 메서드 + structs + Params (single file)
- `domestic/extended_test.go`
- `domestic/testdata/near_new_highlow_success.json`
- `domestic/testdata/overtime_price_success.json`
- `domestic/testdata/overtime_asking_price_success.json`
- `domestic/testdata/overtime_volume_success.json`
- `domestic/testdata/overtime_fluctuation_success.json`

### 신규 (examples)
- `examples/domestic_extended/main.go` — 5 메서드 통합 예시

### 수정 (root)
- `CLAUDE.md` — banner 갱신 (v1.4.0 → v1.5.0, "Phase 2.1" → "Phase 2.2"), Phase 2.2 plan link bullet 추가
- `README.md` — Available Methods 표 갱신 (31 → 36 메서드, heading 갱신)
- `CHANGELOG.md` — `[1.5.0]` entry (above `[1.4.0]`)
- `domestic/doc.go` — Phase 2.2 section (after Phase 2.1 section)

---

## 타입 매핑 (Phase 1 동일)

- **가격류 → `decimal.Decimal` (bare tag)**: `stck_prpr`, `prdy_vrss`, `new_hgpr`, `new_lwpr`, `stck_sdpr`, `ovtm_untp_prpr`, `ovtm_untp_prdy_vrss`, `ovtm_untp_mxpr`, `ovtm_untp_llam`, `ovtm_untp_oprc`, `ovtm_untp_hgpr`, `ovtm_untp_lwpr`, `ovtm_untp_sdpr`, `ovtm_untp_antc_cnpr`, `ovtm_untp_antc_cntg_vrss`, `askp`, `bidp`, `ovtm_untp_askp1..10`, `ovtm_untp_bidp1..10`
- **수량/금액 → `int64,string`**: `acml_vol`, `ovtm_untp_vol`, `ovtm_untp_tr_pbmn`, `ovtm_untp_antc_cnqn`, `ovtm_untp_seln_rsqn`, `ovtm_untp_shnu_rsqn`, `ovtm_untp_exch_vol`, `ovtm_untp_exch_tr_pbmn`, `ovtm_untp_kosdaq_vol`, `ovtm_untp_kosdaq_tr_pbmn`, `ovtm_untp_acml_vol`, `ovtm_untp_acml_tr_pbmn`, `askp_rsqn1`, `bidp_rsqn1`, `ovtm_untp_askp_rsqn1..10`, `ovtm_untp_bidp_rsqn1..10`, `ovtm_untp_askp_icdc1..10`, `ovtm_untp_bidp_icdc1..10`, `ovtm_untp_total_askp_rsqn`, `ovtm_untp_total_bidp_rsqn`, `ovtm_untp_total_askp_rsqn_icdc`, `ovtm_untp_total_bidp_rsqn_icdc`, `ovtm_untp_ntby_bidp_rsqn`, `total_askp_rsqn`, `total_bidp_rsqn`, `total_askp_rsqn_icdc`, `total_bidp_rsqn_icdc`, `ovtm_total_askp_rsqn`, `ovtm_total_bidp_rsqn`, `ovtm_total_askp_icdc`, `ovtm_total_bidp_icdc`, `ovtm_untp_uplm_issu_cnt`, `ovtm_untp_ascn_issu_cnt`, `ovtm_untp_stnr_issu_cnt`, `ovtm_untp_lslm_issu_cnt`, `ovtm_untp_down_issu_cnt`
- **비율 → `float64,string`**: `prdy_ctrt`, `ovtm_untp_prdy_ctrt`, `ovtm_untp_antc_cntg_ctrt`, `hprc_near_rate`, `lwpr_near_rate`, `marg_rate`, `ovtm_vrss_acml_vol_rlim`
- **코드/시간/Y-N/이름/날짜 → 평문 `string`**: `hts_kor_isnm`, `mksc_shrn_iscd`, `stck_shrn_iscd`, `bstp_kor_isnm`, `prdy_vrss_sign`, `ovtm_untp_prdy_vrss_sign`, `ovtm_untp_antc_cntg_vrss_sign`, `mang_issu_cls_name`, `crdt_able_yn`, `new_lstn_cls_name`, `sltr_yn`, `mang_issu_yn`, `mrkt_warn_cls_code`, `mrkt_warn_cls_name`, `trht_yn`, `vlnt_deal_cls_name`, `revl_issu_reas_name`, `insn_pbnt_yn`, `flng_cls_name`, `rprs_mrkt_kor_name`, `ovtm_vi_cls_code`, `ovtm_untp_last_hour`

---

## Task 1: testdata fixtures (5 합성 JSON)

**Files (Create):**
- `domestic/testdata/near_new_highlow_success.json`
- `domestic/testdata/overtime_price_success.json`
- `domestic/testdata/overtime_asking_price_success.json`
- `domestic/testdata/overtime_volume_success.json`
- `domestic/testdata/overtime_fluctuation_success.json`

> KIS docs 응답 필드 기반 합성. OvertimeAskingPrice 의 호가 1~10 + 잔량 1~10 + 증감 1~10 (60 fields) 모두 채움 (testdata 가 cover 하지 않으면 Go json zero-value default 적용).

- [ ] **Step 1: near_new_highlow_success.json**

```json
{
  "rt_cd": "0",
  "msg_cd": "MCA00000",
  "msg1": "정상처리 되었습니다.",
  "output": [
    {
      "hts_kor_isnm": "삼성전자",
      "mksc_shrn_iscd": "005930",
      "stck_prpr": "75800",
      "prdy_vrss_sign": "5",
      "prdy_vrss": "-200",
      "prdy_ctrt": "-0.26",
      "askp": "75900",
      "askp_rsqn1": "1500",
      "bidp": "75800",
      "bidp_rsqn1": "2500",
      "acml_vol": "12345678",
      "new_hgpr": "76500",
      "hprc_near_rate": "1.24",
      "new_lwpr": "74000",
      "lwpr_near_rate": "2.43",
      "stck_sdpr": "76000"
    },
    {
      "hts_kor_isnm": "SK하이닉스",
      "mksc_shrn_iscd": "000660",
      "stck_prpr": "180000",
      "prdy_vrss_sign": "2",
      "prdy_vrss": "1000",
      "prdy_ctrt": "0.56",
      "askp": "180100",
      "askp_rsqn1": "300",
      "bidp": "180000",
      "bidp_rsqn1": "500",
      "acml_vol": "3456789",
      "new_hgpr": "182000",
      "hprc_near_rate": "1.10",
      "new_lwpr": "175000",
      "lwpr_near_rate": "2.86",
      "stck_sdpr": "179000"
    }
  ]
}
```

- [ ] **Step 2: overtime_price_success.json**

```json
{
  "rt_cd": "0",
  "msg_cd": "MCA00000",
  "msg1": "정상처리 되었습니다.",
  "output": {
    "bstp_kor_isnm": "전기전자",
    "mang_issu_cls_name": "",
    "ovtm_untp_prpr": "75700",
    "ovtm_untp_prdy_vrss": "-300",
    "ovtm_untp_prdy_vrss_sign": "5",
    "ovtm_untp_prdy_ctrt": "-0.39",
    "ovtm_untp_vol": "234567",
    "ovtm_untp_tr_pbmn": "17756100000",
    "ovtm_untp_mxpr": "83600",
    "ovtm_untp_llam": "68300",
    "ovtm_untp_oprc": "75900",
    "ovtm_untp_hgpr": "76000",
    "ovtm_untp_lwpr": "75600",
    "marg_rate": "20.00",
    "ovtm_untp_antc_cnpr": "75700",
    "ovtm_untp_antc_cntg_vrss": "-100",
    "ovtm_untp_antc_cntg_vrss_sign": "5",
    "ovtm_untp_antc_cntg_ctrt": "-0.13",
    "ovtm_untp_antc_cnqn": "12345",
    "crdt_able_yn": "Y",
    "new_lstn_cls_name": "",
    "sltr_yn": "N",
    "mang_issu_yn": "N",
    "mrkt_warn_cls_code": "00",
    "trht_yn": "N",
    "vlnt_deal_cls_name": "",
    "ovtm_untp_sdpr": "76000",
    "mrkt_warn_cls_name": "정상",
    "revl_issu_reas_name": "",
    "insn_pbnt_yn": "N",
    "flng_cls_name": "",
    "rprs_mrkt_kor_name": "KOSPI",
    "ovtm_vi_cls_code": "N",
    "bidp": "75700",
    "askp": "75750"
  }
}
```

- [ ] **Step 3: overtime_asking_price_success.json**

```json
{
  "rt_cd": "0",
  "msg_cd": "MCA00000",
  "msg1": "정상처리 되었습니다.",
  "output1": {
    "ovtm_untp_last_hour": "180542",
    "ovtm_untp_askp1": "75750", "ovtm_untp_askp2": "75800", "ovtm_untp_askp3": "75850",
    "ovtm_untp_askp4": "75900", "ovtm_untp_askp5": "75950", "ovtm_untp_askp6": "76000",
    "ovtm_untp_askp7": "76050", "ovtm_untp_askp8": "76100", "ovtm_untp_askp9": "76150",
    "ovtm_untp_askp10": "76200",
    "ovtm_untp_bidp1": "75700", "ovtm_untp_bidp2": "75650", "ovtm_untp_bidp3": "75600",
    "ovtm_untp_bidp4": "75550", "ovtm_untp_bidp5": "75500", "ovtm_untp_bidp6": "75450",
    "ovtm_untp_bidp7": "75400", "ovtm_untp_bidp8": "75350", "ovtm_untp_bidp9": "75300",
    "ovtm_untp_bidp10": "75250",
    "ovtm_untp_askp_icdc1": "100", "ovtm_untp_askp_icdc2": "-50", "ovtm_untp_askp_icdc3": "200",
    "ovtm_untp_askp_icdc4": "0", "ovtm_untp_askp_icdc5": "-30", "ovtm_untp_askp_icdc6": "10",
    "ovtm_untp_askp_icdc7": "0", "ovtm_untp_askp_icdc8": "-20", "ovtm_untp_askp_icdc9": "5",
    "ovtm_untp_askp_icdc10": "0",
    "ovtm_untp_bidp_icdc1": "200", "ovtm_untp_bidp_icdc2": "-100", "ovtm_untp_bidp_icdc3": "150",
    "ovtm_untp_bidp_icdc4": "-80", "ovtm_untp_bidp_icdc5": "20", "ovtm_untp_bidp_icdc6": "0",
    "ovtm_untp_bidp_icdc7": "-10", "ovtm_untp_bidp_icdc8": "0", "ovtm_untp_bidp_icdc9": "0",
    "ovtm_untp_bidp_icdc10": "0",
    "ovtm_untp_askp_rsqn1": "1200", "ovtm_untp_askp_rsqn2": "900", "ovtm_untp_askp_rsqn3": "1500",
    "ovtm_untp_askp_rsqn4": "800", "ovtm_untp_askp_rsqn5": "600", "ovtm_untp_askp_rsqn6": "400",
    "ovtm_untp_askp_rsqn7": "300", "ovtm_untp_askp_rsqn8": "200", "ovtm_untp_askp_rsqn9": "150",
    "ovtm_untp_askp_rsqn10": "100",
    "ovtm_untp_bidp_rsqn1": "2000", "ovtm_untp_bidp_rsqn2": "1500", "ovtm_untp_bidp_rsqn3": "1200",
    "ovtm_untp_bidp_rsqn4": "1000", "ovtm_untp_bidp_rsqn5": "800", "ovtm_untp_bidp_rsqn6": "600",
    "ovtm_untp_bidp_rsqn7": "400", "ovtm_untp_bidp_rsqn8": "300", "ovtm_untp_bidp_rsqn9": "200",
    "ovtm_untp_bidp_rsqn10": "100",
    "ovtm_untp_total_askp_rsqn": "6150",
    "ovtm_untp_total_bidp_rsqn": "8100",
    "ovtm_untp_total_askp_rsqn_icdc": "215",
    "ovtm_untp_total_bidp_rsqn_icdc": "180",
    "ovtm_untp_ntby_bidp_rsqn": "1950",
    "total_askp_rsqn": "11100",
    "total_bidp_rsqn": "13400",
    "total_askp_rsqn_icdc": "215",
    "total_bidp_rsqn_icdc": "180",
    "ovtm_total_askp_rsqn": "6150",
    "ovtm_total_bidp_rsqn": "8100",
    "ovtm_total_askp_icdc": "50",
    "ovtm_total_bidp_icdc": "30"
  }
}
```

- [ ] **Step 4: overtime_volume_success.json**

```json
{
  "rt_cd": "0",
  "msg_cd": "MCA00000",
  "msg1": "정상처리 되었습니다.",
  "output1": {
    "ovtm_untp_exch_vol": "12345678",
    "ovtm_untp_exch_tr_pbmn": "987654321000",
    "ovtm_untp_kosdaq_vol": "9876543",
    "ovtm_untp_kosdaq_tr_pbmn": "543210000000"
  },
  "output2": [
    {
      "stck_shrn_iscd": "005930",
      "hts_kor_isnm": "삼성전자",
      "ovtm_untp_prpr": "75700",
      "ovtm_untp_prdy_vrss": "-300",
      "ovtm_untp_prdy_vrss_sign": "5",
      "ovtm_untp_prdy_ctrt": "-0.39",
      "ovtm_untp_seln_rsqn": "3200",
      "ovtm_untp_shnu_rsqn": "4500",
      "ovtm_untp_vol": "234567",
      "ovtm_vrss_acml_vol_rlim": "1.90",
      "stck_prpr": "75800",
      "acml_vol": "12345678",
      "bidp": "75700",
      "askp": "75750"
    },
    {
      "stck_shrn_iscd": "000660",
      "hts_kor_isnm": "SK하이닉스",
      "ovtm_untp_prpr": "180200",
      "ovtm_untp_prdy_vrss": "1200",
      "ovtm_untp_prdy_vrss_sign": "2",
      "ovtm_untp_prdy_ctrt": "0.67",
      "ovtm_untp_seln_rsqn": "800",
      "ovtm_untp_shnu_rsqn": "1200",
      "ovtm_untp_vol": "98765",
      "ovtm_vrss_acml_vol_rlim": "2.85",
      "stck_prpr": "180000",
      "acml_vol": "3456789",
      "bidp": "180200",
      "askp": "180300"
    }
  ]
}
```

- [ ] **Step 5: overtime_fluctuation_success.json**

```json
{
  "rt_cd": "0",
  "msg_cd": "MCA00000",
  "msg1": "정상처리 되었습니다.",
  "output1": {
    "ovtm_untp_uplm_issu_cnt": "5",
    "ovtm_untp_ascn_issu_cnt": "312",
    "ovtm_untp_stnr_issu_cnt": "187",
    "ovtm_untp_lslm_issu_cnt": "2",
    "ovtm_untp_down_issu_cnt": "154",
    "ovtm_untp_acml_vol": "22345678",
    "ovtm_untp_acml_tr_pbmn": "1987654321000",
    "ovtm_untp_exch_vol": "12345678",
    "ovtm_untp_exch_tr_pbmn": "987654321000",
    "ovtm_untp_kosdaq_vol": "9876543",
    "ovtm_untp_kosdaq_tr_pbmn": "543210000000"
  },
  "output2": [
    {
      "mksc_shrn_iscd": "005930",
      "hts_kor_isnm": "삼성전자",
      "ovtm_untp_prpr": "75700",
      "ovtm_untp_prdy_vrss": "-300",
      "ovtm_untp_prdy_vrss_sign": "5",
      "ovtm_untp_prdy_ctrt": "-0.39",
      "ovtm_untp_askp1": "75750",
      "ovtm_untp_seln_rsqn": "3200",
      "ovtm_untp_bidp1": "75700",
      "ovtm_untp_shnu_rsqn": "4500",
      "ovtm_untp_vol": "234567",
      "ovtm_vrss_acml_vol_rlim": "1.90",
      "stck_prpr": "75800",
      "acml_vol": "12345678",
      "bidp": "75700",
      "askp": "75750"
    },
    {
      "mksc_shrn_iscd": "000660",
      "hts_kor_isnm": "SK하이닉스",
      "ovtm_untp_prpr": "180200",
      "ovtm_untp_prdy_vrss": "1200",
      "ovtm_untp_prdy_vrss_sign": "2",
      "ovtm_untp_prdy_ctrt": "0.67",
      "ovtm_untp_askp1": "180300",
      "ovtm_untp_seln_rsqn": "800",
      "ovtm_untp_bidp1": "180200",
      "ovtm_untp_shnu_rsqn": "1200",
      "ovtm_untp_vol": "98765",
      "ovtm_vrss_acml_vol_rlim": "2.85",
      "stck_prpr": "180000",
      "acml_vol": "3456789",
      "bidp": "180200",
      "askp": "180300"
    }
  ]
}
```

- [ ] **Step 6: 검증**

```bash
for f in domestic/testdata/{near_new_highlow,overtime_price,overtime_asking_price,overtime_volume,overtime_fluctuation}_success.json; do
  python3 -c "import json; json.load(open('$f'))" && echo "$f OK" || echo "$f BROKEN"
done
```
Expected: 5 lines `OK`.

- [ ] **Step 7: Commit**

```bash
git add domestic/testdata/{near_new_highlow,overtime_price,overtime_asking_price,overtime_volume,overtime_fluctuation}_success.json
git commit -m "$(cat <<'EOF'
[chore] Phase 2.2 testdata — 5 합성 JSON fixtures

신고저근접 (near_new_highlow, FHPST01870000) + 시간외현재가 (overtime_price,
FHPST02300000) + 시간외호가 (overtime_asking_price, FHPST02300400) +
시간외거래량순위 (overtime_volume, FHPST02350000) + 시간외등락율순위
(overtime_fluctuation, FHPST02340000). KIS docs 응답 필드 기반 합성.
삼성전자 (005930) / SK하이닉스 (000660) 샘플.

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 2: domestic/extended.go base + InquireNearNewHighlow

**Files:**
- Create: `domestic/extended.go`
- Create: `domestic/extended_test.go`

> 12 query params. FID_COND_SCR_DIV_CODE 는 "20187" 고정. output 은 array (16 fields/item). 가장 먼저 file base 를 setup.

- [ ] **Step 1: 테스트 작성** — `domestic/extended_test.go`

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

func TestClient_InquireNearNewHighlow(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/ranking/near-new-highlow`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "near_new_highlow_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireNearNewHighlow(context.Background(), domestic.InquireNearNewHighlowParams{
		InputISCD:  "0000",
		PrcClsCode: "0",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	// query param 검증
	assert.Equal(t, "J", capturedQuery.Get("fid_cond_mrkt_div_code"))
	assert.Equal(t, "20187", capturedQuery.Get("fid_cond_scr_div_code"))
	assert.Equal(t, "0000", capturedQuery.Get("fid_input_iscd"))
	assert.Equal(t, "0", capturedQuery.Get("fid_prc_cls_code"))

	require.Len(t, res.Output, 2)

	// output[0] 필드 검증
	assert.Equal(t, "삼성전자", res.Output[0].HtsKorIsnm)
	assert.Equal(t, "005930", res.Output[0].MkscShrnIscd)
	d, _ := decimal.NewFromString("75800")
	assert.True(t, d.Equal(res.Output[0].StckPrpr))
	assert.Equal(t, int64(12345678), res.Output[0].AcmlVol)
	d2, _ := decimal.NewFromString("76500")
	assert.True(t, d2.Equal(res.Output[0].NewHgpr))
	assert.InDelta(t, 1.24, res.Output[0].HprcNearRate, 0.01)
	d3, _ := decimal.NewFromString("74000")
	assert.True(t, d3.Equal(res.Output[0].NewLwpr))
	assert.InDelta(t, 2.43, res.Output[0].LwprNearRate, 0.01)
}
```

- [ ] **Step 2: FAIL**

`go test ./domestic/... -run InquireNearNewHighlow -v` — 컴파일 실패.

- [ ] **Step 3: 구현 — `domestic/extended.go`**

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

// NearNewHighlow 는 국내주식 신고/신저근접종목 상위 (FHPST01870000) 응답.
//
// 한투 docs: docs/api/국내주식/국내주식_신고_신저근접종목_상위.md
// path: /uapi/domestic-stock/v1/ranking/near-new-highlow
//
// 최대 30건 확인 가능. 신고 근접 (PrcClsCode="0") 또는 신저 근접 (PrcClsCode="1").
type NearNewHighlow struct {
	Output []NearNewHighlowItem `json:"output"`
}

// NearNewHighlowItem 은 신고/신저근접종목 상위 응답의 한 행.
type NearNewHighlowItem struct {
	HtsKorIsnm   string          `json:"hts_kor_isnm"`   // HTS 한글 종목명
	MkscShrnIscd string          `json:"mksc_shrn_iscd"` // 유가증권 단축 종목코드
	StckPrpr     decimal.Decimal `json:"stck_prpr"`      // 주식 현재가
	PrdyVrssSign string          `json:"prdy_vrss_sign"` // 전일 대비 부호
	PrdyVrss     decimal.Decimal `json:"prdy_vrss"`      // 전일 대비
	PrdyCtrt     float64         `json:"prdy_ctrt,string"`  // 전일 대비율
	Askp         decimal.Decimal `json:"askp"`              // 매도호가
	AskpRsqn1    int64           `json:"askp_rsqn1,string"` // 매도호가 잔량1
	Bidp         decimal.Decimal `json:"bidp"`              // 매수호가
	BidpRsqn1    int64           `json:"bidp_rsqn1,string"` // 매수호가 잔량1
	AcmlVol      int64           `json:"acml_vol,string"`   // 누적 거래량
	NewHgpr      decimal.Decimal `json:"new_hgpr"`          // 신 최고가
	HprcNearRate float64         `json:"hprc_near_rate,string"` // 고가 근접 비율
	NewLwpr      decimal.Decimal `json:"new_lwpr"`              // 신 최저가
	LwprNearRate float64         `json:"lwpr_near_rate,string"` // 저가 근접 비율
	StckSdpr     decimal.Decimal `json:"stck_sdpr"`             // 주식 기준가
}

// InquireNearNewHighlowParams 는 신고/신저근접종목 상위 조회 파라미터.
//
// FID_COND_SCR_DIV_CODE = "20187" 고정 (사용자 변경 불가).
type InquireNearNewHighlowParams struct {
	MarketCode    string // fid_cond_mrkt_div_code — "J":KRX. 빈 값=>"J"
	DivClsCode    string // fid_div_cls_code — 0:전체, 1:관리종목, 2:투자주의, 3:투자경고. 빈 값=>"0"
	InputCnt1     string // fid_input_cnt_1 — 괴리율 최소. 빈 값=>"0"
	InputCnt2     string // fid_input_cnt_2 — 괴리율 최대. 빈 값=>"100"
	PrcClsCode    string // fid_prc_cls_code — 0:신고근접, 1:신저근접. 빈 값=>"0"
	InputISCD     string // fid_input_iscd — 0000:전체, 0001:거래소, 1001:코스닥, 2001:코스피200, 4001:KRX100
	TrgtClsCode   string // fid_trgt_cls_code — 0:전체. 빈 값=>"0"
	TrgtExlsCode  string // fid_trgt_exls_cls_code — 0:전체. 빈 값=>"0"
	AplyRangVol   string // fid_aply_rang_vol — 0:전체, 100:100주 이상. 빈 값=>"0"
	AplyRangPrc1  string // fid_aply_rang_prc_1 — 가격 ~. 빈 값 OK
	AplyRangPrc2  string // fid_aply_rang_prc_2 — ~ 가격. 빈 값 OK
}

// InquireNearNewHighlow 는 국내주식 신고/신저근접종목 상위 호출.
//
// 한투 docs: docs/api/국내주식/국내주식_신고_신저근접종목_상위.md
// path: /uapi/domestic-stock/v1/ranking/near-new-highlow (FHPST01870000)
//
// PrcClsCode="0" 신고 근접 / "1" 신저 근접. 최대 30건.
func (c *Client) InquireNearNewHighlow(ctx context.Context, params InquireNearNewHighlowParams) (*NearNewHighlow, error) {
	market := params.MarketCode
	if market == "" {
		market = "J"
	}
	div := params.DivClsCode
	if div == "" {
		div = "0"
	}
	cnt1 := params.InputCnt1
	if cnt1 == "" {
		cnt1 = "0"
	}
	cnt2 := params.InputCnt2
	if cnt2 == "" {
		cnt2 = "100"
	}
	prc := params.PrcClsCode
	if prc == "" {
		prc = "0"
	}
	tgt := params.TrgtClsCode
	if tgt == "" {
		tgt = "0"
	}
	tgtExcl := params.TrgtExlsCode
	if tgtExcl == "" {
		tgtExcl = "0"
	}
	vol := params.AplyRangVol
	if vol == "" {
		vol = "0"
	}

	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/domestic-stock/v1/ranking/near-new-highlow",
		TrID:   "FHPST01870000",
		Query: map[string]string{
			"fid_cond_mrkt_div_code":   market,
			"fid_cond_scr_div_code":    "20187",
			"fid_div_cls_code":         div,
			"fid_input_cnt_1":          cnt1,
			"fid_input_cnt_2":          cnt2,
			"fid_prc_cls_code":         prc,
			"fid_input_iscd":           params.InputISCD,
			"fid_trgt_cls_code":        tgt,
			"fid_trgt_exls_cls_code":   tgtExcl,
			"fid_aply_rang_vol":        vol,
			"fid_aply_rang_prc_1":      params.AplyRangPrc1,
			"fid_aply_rang_prc_2":      params.AplyRangPrc2,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}

	var res NearNewHighlow
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse NearNewHighlow: %w", err)
	}
	return &res, nil
}
```

- [ ] **Step 4: PASS**

`go test ./domestic/... -run InquireNearNewHighlow -v`

- [ ] **Step 5: Commit**

```bash
git add domestic/extended.go domestic/extended_test.go
git commit -m "$(cat <<'EOF'
[feat] domestic — InquireNearNewHighlow (신고/신저근접종목 상위, FHPST01870000)

NearNewHighlow + NearNewHighlowItem (16 필드) + InquireNearNewHighlowParams
(12 query params, FID_COND_SCR_DIV_CODE="20187" 고정). 신고근접 (PrcClsCode="0")
/ 신저근접 (PrcClsCode="1") 선택. 최대 30건.

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 3: InquireOvertimePrice

**Files:**
- Modify: `domestic/extended.go` (append)
- Modify: `domestic/extended_test.go` (append)

> output 이 object `{}` — array 아님. 35 fields. 시간외 단일가 현재가 + 예상체결 + 호가 + 관리구분.

- [ ] **Step 1: 테스트 추가**

```go
func TestClient_InquireOvertimePrice(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/quotations/inquire-overtime-price`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "overtime_price_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireOvertimePrice(context.Background(), domestic.InquireOvertimePriceParams{
		Symbol: "005930",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, "J", capturedQuery.Get("FID_COND_MRKT_DIV_CODE"))
	assert.Equal(t, "005930", capturedQuery.Get("FID_INPUT_ISCD"))

	// output 필드 검증
	assert.Equal(t, "전기전자", res.Output.BstpKorIsnm)
	d, _ := decimal.NewFromString("75700")
	assert.True(t, d.Equal(res.Output.OvtmUntpPrpr))
	assert.Equal(t, int64(234567), res.Output.OvtmUntpVol)
	assert.InDelta(t, 20.00, res.Output.MargRate, 0.01)
	assert.Equal(t, "N", res.Output.TrhtYn)
	assert.Equal(t, "KOSPI", res.Output.RprsMrktKorName)
	d2, _ := decimal.NewFromString("75700")
	assert.True(t, d2.Equal(res.Output.Bidp))
}
```

- [ ] **Step 2: FAIL**

- [ ] **Step 3: 구현 추가 — APPEND to `domestic/extended.go`**

```go
// OvertimePrice 는 국내주식 시간외현재가 (FHPST02300000) 응답.
//
// 한투 docs: docs/api/국내주식/국내주식_시간외현재가.md
// path: /uapi/domestic-stock/v1/quotations/inquire-overtime-price
//
// 시간외 단일가 현재가 + 예상체결 + 상하한가 + 증거금비율 + 관리구분 등.
type OvertimePrice struct {
	Output OvertimePriceOutput `json:"output"`
}

// OvertimePriceOutput 은 시간외현재가 응답의 output object.
type OvertimePriceOutput struct {
	BstpKorIsnm               string          `json:"bstp_kor_isnm"`                       // 업종 한글 종목명
	MangIssuClsName           string          `json:"mang_issu_cls_name"`                  // 관리 종목 구분 명
	OvtmUntpPrpr              decimal.Decimal `json:"ovtm_untp_prpr"`                      // 시간외 단일가 현재가
	OvtmUntpPrdyVrss          decimal.Decimal `json:"ovtm_untp_prdy_vrss"`                 // 시간외 단일가 전일 대비
	OvtmUntpPrdyVrssSign      string          `json:"ovtm_untp_prdy_vrss_sign"`            // 시간외 단일가 전일 대비 부호
	OvtmUntpPrdyCtrt          float64         `json:"ovtm_untp_prdy_ctrt,string"`          // 시간외 단일가 전일 대비율
	OvtmUntpVol               int64           `json:"ovtm_untp_vol,string"`                // 시간외 단일가 거래량
	OvtmUntpTrPbmn            int64           `json:"ovtm_untp_tr_pbmn,string"`            // 시간외 단일가 거래 대금
	OvtmUntpMxpr              decimal.Decimal `json:"ovtm_untp_mxpr"`                      // 시간외 단일가 상한가
	OvtmUntpLlam              decimal.Decimal `json:"ovtm_untp_llam"`                      // 시간외 단일가 하한가
	OvtmUntpOprc              decimal.Decimal `json:"ovtm_untp_oprc"`                      // 시간외 단일가 시가2
	OvtmUntpHgpr              decimal.Decimal `json:"ovtm_untp_hgpr"`                      // 시간외 단일가 최고가
	OvtmUntpLwpr              decimal.Decimal `json:"ovtm_untp_lwpr"`                      // 시간외 단일가 최저가
	MargRate                  float64         `json:"marg_rate,string"`                    // 증거금 비율
	OvtmUntpAntcCnpr          decimal.Decimal `json:"ovtm_untp_antc_cnpr"`                 // 시간외 단일가 예상 체결가
	OvtmUntpAntcCntgVrss      decimal.Decimal `json:"ovtm_untp_antc_cntg_vrss"`            // 시간외 단일가 예상 체결 대비
	OvtmUntpAntcCntgVrssSign  string          `json:"ovtm_untp_antc_cntg_vrss_sign"`       // 시간외 단일가 예상 체결 대비 부호
	OvtmUntpAntcCntgCtrt      float64         `json:"ovtm_untp_antc_cntg_ctrt,string"`     // 시간외 단일가 예상 체결 대비율
	OvtmUntpAntcCnqn          int64           `json:"ovtm_untp_antc_cnqn,string"`          // 시간외 단일가 예상 체결량
	CrdtAbleYn                string          `json:"crdt_able_yn"`                        // 신용 가능 여부
	NewLstnClsName            string          `json:"new_lstn_cls_name"`                   // 신규 상장 구분 명
	SltrYn                    string          `json:"sltr_yn"`                             // 정리매매 여부
	MangIssuYn                string          `json:"mang_issu_yn"`                        // 관리 종목 여부
	MrktWarnClsCode           string          `json:"mrkt_warn_cls_code"`                  // 시장 경고 구분 코드
	TrhtYn                    string          `json:"trht_yn"`                             // 거래정지 여부
	VlntDealClsName           string          `json:"vlnt_deal_cls_name"`                  // 임의 매매 구분 명
	OvtmUntpSdpr              decimal.Decimal `json:"ovtm_untp_sdpr"`                      // 시간외 단일가 기준가
	MrktWarnClsName           string          `json:"mrkt_warn_cls_name"`                  // 시장 경고 구분 명
	RevlIssuReasName          string          `json:"revl_issu_reas_name"`                 // 재평가 종목 사유 명
	InsnPbntYn                string          `json:"insn_pbnt_yn"`                        // 불성실 공시 여부
	FlngClsName               string          `json:"flng_cls_name"`                       // 락 구분 이름
	RprsMrktKorName           string          `json:"rprs_mrkt_kor_name"`                  // 대표 시장 한글 명
	OvtmViClsCode             string          `json:"ovtm_vi_cls_code"`                    // 시간외단일가VI적용구분코드
	Bidp                      decimal.Decimal `json:"bidp"`                                // 매수호가
	Askp                      decimal.Decimal `json:"askp"`                                // 매도호가
}

// InquireOvertimePriceParams 는 시간외현재가 조회 파라미터.
type InquireOvertimePriceParams struct {
	MarketCode string // FID_COND_MRKT_DIV_CODE — "J":KRX. 빈 값=>"J"
	Symbol     string // FID_INPUT_ISCD — 종목코드 (예 "005930")
}

// InquireOvertimePrice 는 국내주식 시간외현재가 호출.
//
// 한투 docs: docs/api/국내주식/국내주식_시간외현재가.md
// path: /uapi/domestic-stock/v1/quotations/inquire-overtime-price (FHPST02300000)
func (c *Client) InquireOvertimePrice(ctx context.Context, params InquireOvertimePriceParams) (*OvertimePrice, error) {
	market := params.MarketCode
	if market == "" {
		market = "J"
	}

	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/domestic-stock/v1/quotations/inquire-overtime-price",
		TrID:   "FHPST02300000",
		Query: map[string]string{
			"FID_COND_MRKT_DIV_CODE": market,
			"FID_INPUT_ISCD":         params.Symbol,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}

	var res OvertimePrice
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse OvertimePrice: %w", err)
	}
	return &res, nil
}
```

- [ ] **Step 4: PASS**

`go test ./domestic/... -run InquireOvertimePrice -v`

- [ ] **Step 5: Commit**

```bash
git add domestic/extended.go domestic/extended_test.go
git commit -m "$(cat <<'EOF'
[feat] domestic — InquireOvertimePrice (시간외현재가, FHPST02300000)

OvertimePrice + OvertimePriceOutput (35 필드) + InquireOvertimePriceParams.
시간외 단일가 현재가/예상체결/상하한가/증거금비율/관리구분 등.

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 4: InquireOvertimeAskingPrice

**Files:**
- Modify: `domestic/extended.go` (append)
- Modify: `domestic/extended_test.go` (append)

> 가장 큰 struct (~74 fields). 10단계 매도/매수 호가 + 증감 + 잔량 (각 10 = 60 fields) + 합계 + 정규장 합계. Phase 2.1 의 AskingPriceExpCcn 와 동일한 indexed 패턴.

- [ ] **Step 1: 테스트 추가**

```go
func TestClient_InquireOvertimeAskingPrice(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/quotations/inquire-overtime-asking-price`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "overtime_asking_price_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireOvertimeAskingPrice(context.Background(), domestic.InquireOvertimeAskingPriceParams{
		Symbol: "005930",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, "J", capturedQuery.Get("FID_COND_MRKT_DIV_CODE"))
	assert.Equal(t, "005930", capturedQuery.Get("FID_INPUT_ISCD"))

	// output1 핵심 필드 검증
	assert.Equal(t, "180542", res.Output1.OvtmUntpLastHour)
	d, _ := decimal.NewFromString("75750")
	assert.True(t, d.Equal(res.Output1.OvtmUntpAskp1))
	d2, _ := decimal.NewFromString("75700")
	assert.True(t, d2.Equal(res.Output1.OvtmUntpBidp1))
	assert.Equal(t, int64(100), res.Output1.OvtmUntpAskpIcdc1)
	assert.Equal(t, int64(200), res.Output1.OvtmUntpBidpIcdc1)
	assert.Equal(t, int64(1200), res.Output1.OvtmUntpAskpRsqn1)
	assert.Equal(t, int64(2000), res.Output1.OvtmUntpBidpRsqn1)
	assert.Equal(t, int64(6150), res.Output1.OvtmUntpTotalAskpRsqn)
	assert.Equal(t, int64(8100), res.Output1.OvtmUntpTotalBidpRsqn)
	assert.Equal(t, int64(11100), res.Output1.TotalAskpRsqn)
	assert.Equal(t, int64(13400), res.Output1.TotalBidpRsqn)
}
```

- [ ] **Step 2: FAIL**

- [ ] **Step 3: 구현 추가 — APPEND to `domestic/extended.go`**

```go
// OvertimeAskingPrice 는 국내주식 시간외호가 (FHPST02300400) 응답.
//
// 한투 docs: docs/api/국내주식/국내주식_시간외호가.md
// path: /uapi/domestic-stock/v1/quotations/inquire-overtime-asking-price
//
// 시간외 단일가 10단계 호가/증감/잔량 + 정규장 총잔량. output1 만 존재.
type OvertimeAskingPrice struct {
	Output1 OvertimeAskingPriceOrderbook `json:"output1"`
}

// OvertimeAskingPriceOrderbook 은 시간외호가 응답 output1 — 10단계 호가+증감+잔량.
type OvertimeAskingPriceOrderbook struct {
	OvtmUntpLastHour string `json:"ovtm_untp_last_hour"` // 시간외 단일가 최종 시간 (HHMMSS)

	OvtmUntpAskp1  decimal.Decimal `json:"ovtm_untp_askp1"`  // 시간외 단일가 매도호가1
	OvtmUntpAskp2  decimal.Decimal `json:"ovtm_untp_askp2"`
	OvtmUntpAskp3  decimal.Decimal `json:"ovtm_untp_askp3"`
	OvtmUntpAskp4  decimal.Decimal `json:"ovtm_untp_askp4"`
	OvtmUntpAskp5  decimal.Decimal `json:"ovtm_untp_askp5"`
	OvtmUntpAskp6  decimal.Decimal `json:"ovtm_untp_askp6"`
	OvtmUntpAskp7  decimal.Decimal `json:"ovtm_untp_askp7"`
	OvtmUntpAskp8  decimal.Decimal `json:"ovtm_untp_askp8"`
	OvtmUntpAskp9  decimal.Decimal `json:"ovtm_untp_askp9"`
	OvtmUntpAskp10 decimal.Decimal `json:"ovtm_untp_askp10"`

	OvtmUntpBidp1  decimal.Decimal `json:"ovtm_untp_bidp1"`  // 시간외 단일가 매수호가1
	OvtmUntpBidp2  decimal.Decimal `json:"ovtm_untp_bidp2"`
	OvtmUntpBidp3  decimal.Decimal `json:"ovtm_untp_bidp3"`
	OvtmUntpBidp4  decimal.Decimal `json:"ovtm_untp_bidp4"`
	OvtmUntpBidp5  decimal.Decimal `json:"ovtm_untp_bidp5"`
	OvtmUntpBidp6  decimal.Decimal `json:"ovtm_untp_bidp6"`
	OvtmUntpBidp7  decimal.Decimal `json:"ovtm_untp_bidp7"`
	OvtmUntpBidp8  decimal.Decimal `json:"ovtm_untp_bidp8"`
	OvtmUntpBidp9  decimal.Decimal `json:"ovtm_untp_bidp9"`
	OvtmUntpBidp10 decimal.Decimal `json:"ovtm_untp_bidp10"`

	OvtmUntpAskpIcdc1  int64 `json:"ovtm_untp_askp_icdc1,string"`  // 시간외 단일가 매도호가 증감1
	OvtmUntpAskpIcdc2  int64 `json:"ovtm_untp_askp_icdc2,string"`
	OvtmUntpAskpIcdc3  int64 `json:"ovtm_untp_askp_icdc3,string"`
	OvtmUntpAskpIcdc4  int64 `json:"ovtm_untp_askp_icdc4,string"`
	OvtmUntpAskpIcdc5  int64 `json:"ovtm_untp_askp_icdc5,string"`
	OvtmUntpAskpIcdc6  int64 `json:"ovtm_untp_askp_icdc6,string"`
	OvtmUntpAskpIcdc7  int64 `json:"ovtm_untp_askp_icdc7,string"`
	OvtmUntpAskpIcdc8  int64 `json:"ovtm_untp_askp_icdc8,string"`
	OvtmUntpAskpIcdc9  int64 `json:"ovtm_untp_askp_icdc9,string"`
	OvtmUntpAskpIcdc10 int64 `json:"ovtm_untp_askp_icdc10,string"`

	OvtmUntpBidpIcdc1  int64 `json:"ovtm_untp_bidp_icdc1,string"`  // 시간외 단일가 매수호가 증감1
	OvtmUntpBidpIcdc2  int64 `json:"ovtm_untp_bidp_icdc2,string"`
	OvtmUntpBidpIcdc3  int64 `json:"ovtm_untp_bidp_icdc3,string"`
	OvtmUntpBidpIcdc4  int64 `json:"ovtm_untp_bidp_icdc4,string"`
	OvtmUntpBidpIcdc5  int64 `json:"ovtm_untp_bidp_icdc5,string"`
	OvtmUntpBidpIcdc6  int64 `json:"ovtm_untp_bidp_icdc6,string"`
	OvtmUntpBidpIcdc7  int64 `json:"ovtm_untp_bidp_icdc7,string"`
	OvtmUntpBidpIcdc8  int64 `json:"ovtm_untp_bidp_icdc8,string"`
	OvtmUntpBidpIcdc9  int64 `json:"ovtm_untp_bidp_icdc9,string"`
	OvtmUntpBidpIcdc10 int64 `json:"ovtm_untp_bidp_icdc10,string"`

	OvtmUntpAskpRsqn1  int64 `json:"ovtm_untp_askp_rsqn1,string"`  // 시간외 단일가 매도호가 잔량1
	OvtmUntpAskpRsqn2  int64 `json:"ovtm_untp_askp_rsqn2,string"`
	OvtmUntpAskpRsqn3  int64 `json:"ovtm_untp_askp_rsqn3,string"`
	OvtmUntpAskpRsqn4  int64 `json:"ovtm_untp_askp_rsqn4,string"`
	OvtmUntpAskpRsqn5  int64 `json:"ovtm_untp_askp_rsqn5,string"`
	OvtmUntpAskpRsqn6  int64 `json:"ovtm_untp_askp_rsqn6,string"`
	OvtmUntpAskpRsqn7  int64 `json:"ovtm_untp_askp_rsqn7,string"`
	OvtmUntpAskpRsqn8  int64 `json:"ovtm_untp_askp_rsqn8,string"`
	OvtmUntpAskpRsqn9  int64 `json:"ovtm_untp_askp_rsqn9,string"`
	OvtmUntpAskpRsqn10 int64 `json:"ovtm_untp_askp_rsqn10,string"`

	OvtmUntpBidpRsqn1  int64 `json:"ovtm_untp_bidp_rsqn1,string"`  // 시간외 단일가 매수호가 잔량1
	OvtmUntpBidpRsqn2  int64 `json:"ovtm_untp_bidp_rsqn2,string"`
	OvtmUntpBidpRsqn3  int64 `json:"ovtm_untp_bidp_rsqn3,string"`
	OvtmUntpBidpRsqn4  int64 `json:"ovtm_untp_bidp_rsqn4,string"`
	OvtmUntpBidpRsqn5  int64 `json:"ovtm_untp_bidp_rsqn5,string"`
	OvtmUntpBidpRsqn6  int64 `json:"ovtm_untp_bidp_rsqn6,string"`
	OvtmUntpBidpRsqn7  int64 `json:"ovtm_untp_bidp_rsqn7,string"`
	OvtmUntpBidpRsqn8  int64 `json:"ovtm_untp_bidp_rsqn8,string"`
	OvtmUntpBidpRsqn9  int64 `json:"ovtm_untp_bidp_rsqn9,string"`
	OvtmUntpBidpRsqn10 int64 `json:"ovtm_untp_bidp_rsqn10,string"`

	OvtmUntpTotalAskpRsqn     int64 `json:"ovtm_untp_total_askp_rsqn,string"`      // 시간외 단일가 총 매도호가 잔량
	OvtmUntpTotalBidpRsqn     int64 `json:"ovtm_untp_total_bidp_rsqn,string"`      // 시간외 단일가 총 매수호가 잔량
	OvtmUntpTotalAskpRsqnIcdc int64 `json:"ovtm_untp_total_askp_rsqn_icdc,string"` // 시간외 단일가 총 매도호가 잔량 증감
	OvtmUntpTotalBidpRsqnIcdc int64 `json:"ovtm_untp_total_bidp_rsqn_icdc,string"` // 시간외 단일가 총 매수호가 잔량 증감
	OvtmUntpNtbyBidpRsqn      int64 `json:"ovtm_untp_ntby_bidp_rsqn,string"`       // 시간외 단일가 순매수 호가 잔량
	TotalAskpRsqn             int64 `json:"total_askp_rsqn,string"`                // 총 매도호가 잔량 (정규장)
	TotalBidpRsqn             int64 `json:"total_bidp_rsqn,string"`                // 총 매수호가 잔량 (정규장)
	TotalAskpRsqnIcdc         int64 `json:"total_askp_rsqn_icdc,string"`           // 총 매도호가 잔량 증감
	TotalBidpRsqnIcdc         int64 `json:"total_bidp_rsqn_icdc,string"`           // 총 매수호가 잔량 증감
	OvtmTotalAskpRsqn         int64 `json:"ovtm_total_askp_rsqn,string"`           // 시간외 총 매도호가 잔량
	OvtmTotalBidpRsqn         int64 `json:"ovtm_total_bidp_rsqn,string"`           // 시간외 총 매수호가 잔량
	OvtmTotalAskpIcdc         int64 `json:"ovtm_total_askp_icdc,string"`           // 시간외 총 매도호가 증감
	OvtmTotalBidpIcdc         int64 `json:"ovtm_total_bidp_icdc,string"`           // 시간외 총 매수호가 증감
}

// InquireOvertimeAskingPriceParams 는 시간외호가 조회 파라미터.
type InquireOvertimeAskingPriceParams struct {
	Symbol     string // FID_INPUT_ISCD — 종목코드 (예 "005930")
	MarketCode string // FID_COND_MRKT_DIV_CODE — "J":KRX. 빈 값=>"J"
}

// InquireOvertimeAskingPrice 는 국내주식 시간외호가 호출.
//
// 한투 docs: docs/api/국내주식/국내주식_시간외호가.md
// path: /uapi/domestic-stock/v1/quotations/inquire-overtime-asking-price (FHPST02300400)
//
// 시간외 단일가 10단계 호가/증감/잔량 (총 60 fields) + 시간외/정규장 총잔량.
func (c *Client) InquireOvertimeAskingPrice(ctx context.Context, params InquireOvertimeAskingPriceParams) (*OvertimeAskingPrice, error) {
	market := params.MarketCode
	if market == "" {
		market = "J"
	}

	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/domestic-stock/v1/quotations/inquire-overtime-asking-price",
		TrID:   "FHPST02300400",
		Query: map[string]string{
			"FID_INPUT_ISCD":          params.Symbol,
			"FID_COND_MRKT_DIV_CODE": market,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}

	var res OvertimeAskingPrice
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse OvertimeAskingPrice: %w", err)
	}
	return &res, nil
}
```

- [ ] **Step 4: PASS**

`go test ./domestic/... -run InquireOvertimeAskingPrice -v`

- [ ] **Step 5: Commit**

```bash
git add domestic/extended.go domestic/extended_test.go
git commit -m "$(cat <<'EOF'
[feat] domestic — InquireOvertimeAskingPrice (시간외호가, FHPST02300400)

OvertimeAskingPrice + OvertimeAskingPriceOrderbook (~74 필드):
- 시간외 단일가 호가 1~10 (10 fields)
- 시간외 단일가 매수호가 1~10 (10 fields)
- 시간외 단일가 매도호가 증감 1~10 (10 fields)
- 시간외 단일가 매수호가 증감 1~10 (10 fields)
- 시간외 단일가 매도호가 잔량 1~10 (10 fields)
- 시간외 단일가 매수호가 잔량 1~10 (10 fields)
- 시간외/정규장 총잔량 + 순매수 호가 잔량 (14 fields)

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 5: InquireOvertimeVolume

**Files:**
- Modify: `domestic/extended.go` (append)
- Modify: `domestic/extended_test.go` (append)

> output1 (4 fields summary object) + output2 (array, 14 fields/item). FID_COND_SCR_DIV_CODE = "20235" 고정.

- [ ] **Step 1: 테스트 추가**

```go
func TestClient_InquireOvertimeVolume(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/ranking/overtime-volume`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "overtime_volume_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireOvertimeVolume(context.Background(), domestic.InquireOvertimeVolumeParams{
		InputISCD: "0000",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, "J", capturedQuery.Get("FID_COND_MRKT_DIV_CODE"))
	assert.Equal(t, "20235", capturedQuery.Get("FID_COND_SCR_DIV_CODE"))
	assert.Equal(t, "0000", capturedQuery.Get("FID_INPUT_ISCD"))

	// output1 검증
	assert.Equal(t, int64(12345678), res.Output1.OvtmUntpExchVol)
	assert.Equal(t, int64(9876543), res.Output1.OvtmUntpKosdaqVol)

	// output2 검증
	require.Len(t, res.Output2, 2)
	assert.Equal(t, "005930", res.Output2[0].StckShrnIscd)
	d, _ := decimal.NewFromString("75700")
	assert.True(t, d.Equal(res.Output2[0].OvtmUntpPrpr))
	assert.Equal(t, int64(234567), res.Output2[0].OvtmUntpVol)
	assert.InDelta(t, 1.90, res.Output2[0].OvtmVrssAcmlVolRlim, 0.01)
}
```

- [ ] **Step 2: FAIL**

- [ ] **Step 3: 구현 추가 — APPEND to `domestic/extended.go`**

```go
// OvertimeVolume 은 국내주식 시간외거래량순위 (FHPST02350000) 응답.
//
// 한투 docs: docs/api/국내주식/국내주식_시간외거래량순위.md
// path: /uapi/domestic-stock/v1/ranking/overtime-volume
//
// output1: 거래소/코스닥 합계 (4 fields). output2: 종목별 시간외 거래량 순위 array.
type OvertimeVolume struct {
	Output1 OvertimeVolumeSummary `json:"output1"`
	Output2 []OvertimeVolumeItem  `json:"output2"`
}

// OvertimeVolumeSummary 는 시간외거래량순위 output1 — 시장 전체 합계.
type OvertimeVolumeSummary struct {
	OvtmUntpExchVol      int64 `json:"ovtm_untp_exch_vol,string"`       // 시간외 단일가 거래소 거래량
	OvtmUntpExchTrPbmn   int64 `json:"ovtm_untp_exch_tr_pbmn,string"`   // 시간외 단일가 거래소 거래대금
	OvtmUntpKosdaqVol    int64 `json:"ovtm_untp_kosdaq_vol,string"`     // 시간외 단일가 KOSDAQ 거래량
	OvtmUntpKosdaqTrPbmn int64 `json:"ovtm_untp_kosdaq_tr_pbmn,string"` // 시간외 단일가 KOSDAQ 거래대금
}

// OvertimeVolumeItem 은 시간외거래량순위 output2 의 한 행.
type OvertimeVolumeItem struct {
	StckShrnIscd         string          `json:"stck_shrn_iscd"`                   // 주식 단축 종목코드
	HtsKorIsnm           string          `json:"hts_kor_isnm"`                     // HTS 한글 종목명
	OvtmUntpPrpr         decimal.Decimal `json:"ovtm_untp_prpr"`                   // 시간외 단일가 현재가
	OvtmUntpPrdyVrss     decimal.Decimal `json:"ovtm_untp_prdy_vrss"`              // 시간외 단일가 전일 대비
	OvtmUntpPrdyVrssSign string          `json:"ovtm_untp_prdy_vrss_sign"`         // 시간외 단일가 전일 대비 부호
	OvtmUntpPrdyCtrt     float64         `json:"ovtm_untp_prdy_ctrt,string"`       // 시간외 단일가 전일 대비율
	OvtmUntpSelnRsqn     int64           `json:"ovtm_untp_seln_rsqn,string"`       // 시간외 단일가 매도 잔량
	OvtmUntpShnuRsqn     int64           `json:"ovtm_untp_shnu_rsqn,string"`       // 시간외 단일가 매수 잔량
	OvtmUntpVol          int64           `json:"ovtm_untp_vol,string"`             // 시간외 단일가 거래량
	OvtmVrssAcmlVolRlim  float64         `json:"ovtm_vrss_acml_vol_rlim,string"`   // 시간외 대비 누적 거래량 비중
	StckPrpr             decimal.Decimal `json:"stck_prpr"`                        // 주식 현재가 (정규장)
	AcmlVol              int64           `json:"acml_vol,string"`                  // 누적 거래량 (정규장)
	Bidp                 decimal.Decimal `json:"bidp"`                             // 매수호가
	Askp                 decimal.Decimal `json:"askp"`                             // 매도호가
}

// InquireOvertimeVolumeParams 는 시간외거래량순위 조회 파라미터.
//
// FID_COND_SCR_DIV_CODE = "20235" 고정 (사용자 변경 불가).
type InquireOvertimeVolumeParams struct {
	MarketCode    string // FID_COND_MRKT_DIV_CODE — "J":KRX. 빈 값=>"J"
	InputISCD     string // FID_INPUT_ISCD — 0000:전체, 0001:코스피, 1001:코스닥
	RankSortCode  string // FID_RANK_SORT_CLS_CODE — 0:매수잔량, 1:매도잔량, 2:거래량. 빈 값=>"2"
	InputPrice1   string // FID_INPUT_PRICE_1 — 가격 ~. 빈 값 OK
	InputPrice2   string // FID_INPUT_PRICE_2 — ~ 가격. 빈 값 OK
	VolCount      string // FID_VOL_CNT — 거래량 ~. 빈 값 OK
	TrgtClsCode   string // FID_TRGT_CLS_CODE — 공백 입력
	TrgtExlsCode  string // FID_TRGT_EXLS_CLS_CODE — 공백 입력
}

// InquireOvertimeVolume 은 국내주식 시간외거래량순위 호출.
//
// 한투 docs: docs/api/국내주식/국내주식_시간외거래량순위.md
// path: /uapi/domestic-stock/v1/ranking/overtime-volume (FHPST02350000)
//
// output1: 거래소/코스닥 합계. output2: 최대 30건 종목별 시간외 거래량 순위.
func (c *Client) InquireOvertimeVolume(ctx context.Context, params InquireOvertimeVolumeParams) (*OvertimeVolume, error) {
	market := params.MarketCode
	if market == "" {
		market = "J"
	}
	sort := params.RankSortCode
	if sort == "" {
		sort = "2"
	}

	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/domestic-stock/v1/ranking/overtime-volume",
		TrID:   "FHPST02350000",
		Query: map[string]string{
			"FID_COND_MRKT_DIV_CODE":  market,
			"FID_COND_SCR_DIV_CODE":   "20235",
			"FID_INPUT_ISCD":          params.InputISCD,
			"FID_RANK_SORT_CLS_CODE":  sort,
			"FID_INPUT_PRICE_1":       params.InputPrice1,
			"FID_INPUT_PRICE_2":       params.InputPrice2,
			"FID_VOL_CNT":             params.VolCount,
			"FID_TRGT_CLS_CODE":       params.TrgtClsCode,
			"FID_TRGT_EXLS_CLS_CODE":  params.TrgtExlsCode,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}

	var res OvertimeVolume
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse OvertimeVolume: %w", err)
	}
	return &res, nil
}
```

- [ ] **Step 4: PASS**

`go test ./domestic/... -run InquireOvertimeVolume -v`

- [ ] **Step 5: Commit**

```bash
git add domestic/extended.go domestic/extended_test.go
git commit -m "$(cat <<'EOF'
[feat] domestic — InquireOvertimeVolume (시간외거래량순위, FHPST02350000)

OvertimeVolume + OvertimeVolumeSummary (output1, 4 필드) + OvertimeVolumeItem
(output2 array, 14 필드) + InquireOvertimeVolumeParams (9 query params,
FID_COND_SCR_DIV_CODE="20235" 고정). 최대 30건.

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 6: InquireOvertimeFluctuation

**Files:**
- Modify: `domestic/extended.go` (append)
- Modify: `domestic/extended_test.go` (append)

> output1 (11 fields summary — 상한/상승/보합/하한/하락 종목 수 + 거래량/대금) + output2 (array, 16 fields/item). FID_COND_SCR_DIV_CODE = "20234" 고정.

- [ ] **Step 1: 테스트 추가**

```go
func TestClient_InquireOvertimeFluctuation(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/ranking/overtime-fluctuation`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "overtime_fluctuation_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireOvertimeFluctuation(context.Background(), domestic.InquireOvertimeFluctuationParams{
		InputISCD:  "0000",
		DivClsCode: "2",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, "J", capturedQuery.Get("FID_COND_MRKT_DIV_CODE"))
	assert.Equal(t, "20234", capturedQuery.Get("FID_COND_SCR_DIV_CODE"))
	assert.Equal(t, "0000", capturedQuery.Get("FID_INPUT_ISCD"))
	assert.Equal(t, "2", capturedQuery.Get("FID_DIV_CLS_CODE"))

	// output1 검증
	assert.Equal(t, int64(5), res.Output1.OvtmUntpUplmIssuCnt)
	assert.Equal(t, int64(312), res.Output1.OvtmUntpAscnIssuCnt)
	assert.Equal(t, int64(22345678), res.Output1.OvtmUntpAcmlVol)

	// output2 검증
	require.Len(t, res.Output2, 2)
	assert.Equal(t, "005930", res.Output2[0].MkscShrnIscd)
	d, _ := decimal.NewFromString("75700")
	assert.True(t, d.Equal(res.Output2[0].OvtmUntpPrpr))
	assert.InDelta(t, -0.39, res.Output2[0].OvtmUntpPrdyCtrt, 0.01)
	assert.Equal(t, int64(234567), res.Output2[0].OvtmUntpVol)
	assert.InDelta(t, 1.90, res.Output2[0].OvtmVrssAcmlVolRlim, 0.01)
}
```

- [ ] **Step 2: FAIL**

- [ ] **Step 3: 구현 추가 — APPEND to `domestic/extended.go`**

```go
// OvertimeFluctuation 은 국내주식 시간외등락율순위 (FHPST02340000) 응답.
//
// 한투 docs: docs/api/국내주식/국내주식_시간외등락율순위.md
// path: /uapi/domestic-stock/v1/ranking/overtime-fluctuation
//
// output1: 상한/상승/보합/하한/하락 종목 수 + 거래량/대금 합계 (11 fields).
// output2: 종목별 시간외 등락율 순위 array (16 fields/item).
type OvertimeFluctuation struct {
	Output1 OvertimeFluctuationSummary `json:"output1"`
	Output2 []OvertimeFluctuationItem  `json:"output2"`
}

// OvertimeFluctuationSummary 는 시간외등락율순위 output1 — 시장 전체 통계.
type OvertimeFluctuationSummary struct {
	OvtmUntpUplmIssuCnt   int64 `json:"ovtm_untp_uplm_issu_cnt,string"`    // 시간외 단일가 상한 종목 수
	OvtmUntpAscnIssuCnt   int64 `json:"ovtm_untp_ascn_issu_cnt,string"`    // 시간외 단일가 상승 종목 수
	OvtmUntpStnrIssuCnt   int64 `json:"ovtm_untp_stnr_issu_cnt,string"`    // 시간외 단일가 보합 종목 수
	OvtmUntpLslmIssuCnt   int64 `json:"ovtm_untp_lslm_issu_cnt,string"`    // 시간외 단일가 하한 종목 수
	OvtmUntpDownIssuCnt   int64 `json:"ovtm_untp_down_issu_cnt,string"`    // 시간외 단일가 하락 종목 수
	OvtmUntpAcmlVol       int64 `json:"ovtm_untp_acml_vol,string"`         // 시간외 단일가 누적 거래량
	OvtmUntpAcmlTrPbmn    int64 `json:"ovtm_untp_acml_tr_pbmn,string"`     // 시간외 단일가 누적 거래대금
	OvtmUntpExchVol       int64 `json:"ovtm_untp_exch_vol,string"`         // 시간외 단일가 거래소 거래량
	OvtmUntpExchTrPbmn    int64 `json:"ovtm_untp_exch_tr_pbmn,string"`     // 시간외 단일가 거래소 거래대금
	OvtmUntpKosdaqVol     int64 `json:"ovtm_untp_kosdaq_vol,string"`       // 시간외 단일가 KOSDAQ 거래량
	OvtmUntpKosdaqTrPbmn  int64 `json:"ovtm_untp_kosdaq_tr_pbmn,string"`   // 시간외 단일가 KOSDAQ 거래대금
}

// OvertimeFluctuationItem 은 시간외등락율순위 output2 의 한 행.
type OvertimeFluctuationItem struct {
	MkscShrnIscd         string          `json:"mksc_shrn_iscd"`                   // 유가증권 단축 종목코드
	HtsKorIsnm           string          `json:"hts_kor_isnm"`                     // HTS 한글 종목명
	OvtmUntpPrpr         decimal.Decimal `json:"ovtm_untp_prpr"`                   // 시간외 단일가 현재가
	OvtmUntpPrdyVrss     decimal.Decimal `json:"ovtm_untp_prdy_vrss"`              // 시간외 단일가 전일 대비
	OvtmUntpPrdyVrssSign string          `json:"ovtm_untp_prdy_vrss_sign"`         // 시간외 단일가 전일 대비 부호
	OvtmUntpPrdyCtrt     float64         `json:"ovtm_untp_prdy_ctrt,string"`       // 시간외 단일가 전일 대비율
	OvtmUntpAskp1        decimal.Decimal `json:"ovtm_untp_askp1"`                  // 시간외 단일가 매도호가1
	OvtmUntpSelnRsqn     int64           `json:"ovtm_untp_seln_rsqn,string"`       // 시간외 단일가 매도 잔량
	OvtmUntpBidp1        decimal.Decimal `json:"ovtm_untp_bidp1"`                  // 시간외 단일가 매수호가1
	OvtmUntpShnuRsqn     int64           `json:"ovtm_untp_shnu_rsqn,string"`       // 시간외 단일가 매수 잔량
	OvtmUntpVol          int64           `json:"ovtm_untp_vol,string"`             // 시간외 단일가 거래량
	OvtmVrssAcmlVolRlim  float64         `json:"ovtm_vrss_acml_vol_rlim,string"`   // 시간외 대비 누적 거래량 비중
	StckPrpr             decimal.Decimal `json:"stck_prpr"`                        // 주식 현재가 (정규장)
	AcmlVol              int64           `json:"acml_vol,string"`                  // 누적 거래량 (정규장)
	Bidp                 decimal.Decimal `json:"bidp"`                             // 매수호가
	Askp                 decimal.Decimal `json:"askp"`                             // 매도호가
}

// InquireOvertimeFluctuationParams 는 시간외등락율순위 조회 파라미터.
//
// FID_COND_SCR_DIV_CODE = "20234" 고정 (사용자 변경 불가).
type InquireOvertimeFluctuationParams struct {
	MarketCode   string // FID_COND_MRKT_DIV_CODE — "J":KRX. 빈 값=>"J"
	MrktClsCode  string // FID_MRKT_CLS_CODE — 공백 입력
	InputISCD    string // FID_INPUT_ISCD — 0000:전체, 0001:코스피, 1001:코스닥
	DivClsCode   string // FID_DIV_CLS_CODE — 1:상한가, 2:상승률, 3:보합, 4:하한가, 5:하락률. 빈 값=>"2"
	InputPrice1  string // FID_INPUT_PRICE_1 — 가격 ~. 빈 값 OK
	InputPrice2  string // FID_INPUT_PRICE_2 — ~ 가격. 빈 값 OK
	VolCount     string // FID_VOL_CNT — 거래량 ~. 빈 값 OK
	TrgtClsCode  string // FID_TRGT_CLS_CODE — 공백 입력
	TrgtExlsCode string // FID_TRGT_EXLS_CLS_CODE — 공백 입력
}

// InquireOvertimeFluctuation 은 국내주식 시간외등락율순위 호출.
//
// 한투 docs: docs/api/국내주식/국내주식_시간외등락율순위.md
// path: /uapi/domestic-stock/v1/ranking/overtime-fluctuation (FHPST02340000)
//
// output1: 상한/상승/보합/하한/하락 종목 수 + 거래량/대금 합계.
// output2: 최대 30건 종목별 시간외 등락율 순위.
func (c *Client) InquireOvertimeFluctuation(ctx context.Context, params InquireOvertimeFluctuationParams) (*OvertimeFluctuation, error) {
	market := params.MarketCode
	if market == "" {
		market = "J"
	}
	div := params.DivClsCode
	if div == "" {
		div = "2"
	}

	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/domestic-stock/v1/ranking/overtime-fluctuation",
		TrID:   "FHPST02340000",
		Query: map[string]string{
			"FID_COND_MRKT_DIV_CODE":  market,
			"FID_MRKT_CLS_CODE":       params.MrktClsCode,
			"FID_COND_SCR_DIV_CODE":   "20234",
			"FID_INPUT_ISCD":          params.InputISCD,
			"FID_DIV_CLS_CODE":        div,
			"FID_INPUT_PRICE_1":       params.InputPrice1,
			"FID_INPUT_PRICE_2":       params.InputPrice2,
			"FID_VOL_CNT":             params.VolCount,
			"FID_TRGT_CLS_CODE":       params.TrgtClsCode,
			"FID_TRGT_EXLS_CLS_CODE":  params.TrgtExlsCode,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}

	var res OvertimeFluctuation
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse OvertimeFluctuation: %w", err)
	}
	return &res, nil
}
```

- [ ] **Step 4: PASS**

`go test ./domestic/... -run InquireOvertimeFluctuation -v`

- [ ] **Step 5: 전체 회귀 테스트**

`go test ./... -count=1` — all PASS.

- [ ] **Step 6: Commit**

```bash
git add domestic/extended.go domestic/extended_test.go
git commit -m "$(cat <<'EOF'
[feat] domestic — InquireOvertimeFluctuation (시간외등락율순위, FHPST02340000)

OvertimeFluctuation + OvertimeFluctuationSummary (output1, 11 필드: 상한/상승/
보합/하한/하락 종목 수 + 거래소/코스닥 거래량/대금) + OvertimeFluctuationItem
(output2 array, 16 필드) + InquireOvertimeFluctuationParams (10 query params,
FID_COND_SCR_DIV_CODE="20234" 고정). 최대 30건.

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 7: examples/domestic_extended/main.go

**Files:**
- Create: `examples/domestic_extended/main.go`

- [ ] **Step 1: example 작성**

```go
// domestic_extended example: InquireNearNewHighlow + InquireOvertimePrice +
// InquireOvertimeAskingPrice + InquireOvertimeVolume + InquireOvertimeFluctuation.
//
// Run: KIS credentials env vars 후 go run ./examples/domestic_extended
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

	// 1. 신고/신저근접종목 상위
	nh, err := client.Domestic.InquireNearNewHighlow(ctx, domestic.InquireNearNewHighlowParams{
		InputISCD:  "0000",
		PrcClsCode: "0", // 0:신고근접
	})
	if err != nil {
		log.Fatalf("InquireNearNewHighlow: %v", err)
	}
	fmt.Printf("[신고근접 상위 %d 건]\n", len(nh.Output))
	for i, item := range nh.Output {
		if i >= 5 {
			fmt.Println("  ... (이하 생략)")
			break
		}
		fmt.Printf("  %s (%s): %s원 (신고가 %s, 근접율 %v%%)\n",
			item.HtsKorIsnm, item.MkscShrnIscd, item.StckPrpr, item.NewHgpr, item.HprcNearRate)
	}

	// 2. 시간외현재가
	op, err := client.Domestic.InquireOvertimePrice(ctx, domestic.InquireOvertimePriceParams{
		Symbol: symbol,
	})
	if err != nil {
		log.Fatalf("InquireOvertimePrice: %v", err)
	}
	fmt.Printf("\n[%s] 시간외현재가\n", symbol)
	fmt.Printf("  현재가: %s원 (전일대비 %s, %v%%)\n",
		op.Output.OvtmUntpPrpr, op.Output.OvtmUntpPrdyVrss, op.Output.OvtmUntpPrdyCtrt)
	fmt.Printf("  거래량: %d주, 예상체결: %s원\n",
		op.Output.OvtmUntpVol, op.Output.OvtmUntpAntcCnpr)

	// 3. 시간외호가
	oa, err := client.Domestic.InquireOvertimeAskingPrice(ctx, domestic.InquireOvertimeAskingPriceParams{
		Symbol: symbol,
	})
	if err != nil {
		log.Fatalf("InquireOvertimeAskingPrice: %v", err)
	}
	fmt.Printf("\n[%s] 시간외호가 (최종시간 %s)\n", symbol, oa.Output1.OvtmUntpLastHour)
	fmt.Printf("  매도1: %s @ %d주, 매수1: %s @ %d주\n",
		oa.Output1.OvtmUntpAskp1, oa.Output1.OvtmUntpAskpRsqn1,
		oa.Output1.OvtmUntpBidp1, oa.Output1.OvtmUntpBidpRsqn1)
	fmt.Printf("  시간외 총매도잔량: %d, 총매수잔량: %d\n",
		oa.Output1.OvtmUntpTotalAskpRsqn, oa.Output1.OvtmUntpTotalBidpRsqn)

	// 4. 시간외거래량순위
	ov, err := client.Domestic.InquireOvertimeVolume(ctx, domestic.InquireOvertimeVolumeParams{
		InputISCD: "0000",
	})
	if err != nil {
		log.Fatalf("InquireOvertimeVolume: %v", err)
	}
	fmt.Printf("\n[시간외거래량순위 상위 %d 건]\n", len(ov.Output2))
	fmt.Printf("  거래소 합계: 거래량 %d주, 거래대금 %d원\n",
		ov.Output1.OvtmUntpExchVol, ov.Output1.OvtmUntpExchTrPbmn)
	for i, item := range ov.Output2 {
		if i >= 5 {
			fmt.Println("  ... (이하 생략)")
			break
		}
		fmt.Printf("  %s (%s): 시간외 %d주 (%v%%)\n",
			item.HtsKorIsnm, item.StckShrnIscd, item.OvtmUntpVol, item.OvtmVrssAcmlVolRlim)
	}

	// 5. 시간외등락율순위
	of, err := client.Domestic.InquireOvertimeFluctuation(ctx, domestic.InquireOvertimeFluctuationParams{
		InputISCD:  "0000",
		DivClsCode: "2", // 2:상승률
	})
	if err != nil {
		log.Fatalf("InquireOvertimeFluctuation: %v", err)
	}
	fmt.Printf("\n[시간외등락율순위 (상승률) 상위 %d 건]\n", len(of.Output2))
	fmt.Printf("  상한: %d, 상승: %d, 보합: %d, 하락: %d, 하한: %d\n",
		of.Output1.OvtmUntpUplmIssuCnt, of.Output1.OvtmUntpAscnIssuCnt,
		of.Output1.OvtmUntpStnrIssuCnt, of.Output1.OvtmUntpDownIssuCnt,
		of.Output1.OvtmUntpLslmIssuCnt)
	for i, item := range of.Output2 {
		if i >= 5 {
			fmt.Println("  ... (이하 생략)")
			break
		}
		fmt.Printf("  %s (%s): %s원 (시간외 %v%%)\n",
			item.HtsKorIsnm, item.MkscShrnIscd, item.OvtmUntpPrpr, item.OvtmUntpPrdyCtrt)
	}
}
```

- [ ] **Step 2: 컴파일 검증**

`go build ./examples/domestic_extended && echo OK`

- [ ] **Step 3: Commit**

```bash
git add examples/domestic_extended
git commit -m "$(cat <<'EOF'
[feat] examples/domestic_extended — 신고저근접 + 시간외 5 메서드 통합 예시

신고/신저근접 상위 + 시간외현재가 + 시간외호가 + 시간외거래량순위 +
시간외등락율순위 출력. 삼성전자 (005930) + 전체시장 (0000) 샘플.

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 8: 문서 갱신

**Files:**
- Modify: `CLAUDE.md`
- Modify: `README.md`
- Modify: `CHANGELOG.md`
- Modify: `domestic/doc.go`

- [ ] **Step 1: CLAUDE.md — banner 갱신**

Replace:
```
> **Phase 2.1 — 국내 호가/체결 (v1.4.0).** Phase 2.2+ 는 추후 sub-plan 으로.
```
With:
```
> **Phase 2.2 — 국내 신고저/시간외 (v1.5.0).** Phase 2.3+ 는 추후 sub-plan 으로.
```

ADD spec link bullet after Phase 2.1 plan link:
```markdown
- Phase 2.2 implementation plan: [`docs/superpowers/specs/2026-05-05-phase2-2-extended-implementation-plan.md`](docs/superpowers/specs/2026-05-05-phase2-2-extended-implementation-plan.md)
```

- [ ] **Step 2: README.md — Available Methods 표 갱신**

Find existing `## Available Methods (Phase 1.2 ~ 2.1)` heading. Update heading to `## Available Methods (Phase 1.2 ~ 2.2)` and APPEND 5 rows AT THE END:

```markdown
| `Domestic.InquireNearNewHighlow` | `ranking/near-new-highlow` | FHPST01870000 |
| `Domestic.InquireOvertimePrice` | `quotations/inquire-overtime-price` | FHPST02300000 |
| `Domestic.InquireOvertimeAskingPrice` | `quotations/inquire-overtime-asking-price` | FHPST02300400 |
| `Domestic.InquireOvertimeVolume` | `ranking/overtime-volume` | FHPST02350000 |
| `Domestic.InquireOvertimeFluctuation` | `ranking/overtime-fluctuation` | FHPST02340000 |
```

Also update the method count in the heading if present (31 → 36).

- [ ] **Step 3: CHANGELOG.md — `[1.5.0]` entry**

ADD AT THE TOP (above `## [1.4.0]`):

```markdown
## [1.5.0] - 2026-05-05

### Added — Phase 2.2 (국내 신고저가 / 시간외)

- `Domestic.InquireNearNewHighlow` — 국내주식 신고/신저근접종목 상위 (FHPST01870000) — 신고근접/신저근접 최대 30건
- `Domestic.InquireOvertimePrice` — 국내주식 시간외현재가 (FHPST02300000) — 시간외 단일가 현재가/예상체결/상하한가/관리구분
- `Domestic.InquireOvertimeAskingPrice` — 국내주식 시간외호가 (FHPST02300400) — 10단계 호가/증감/잔량 + 정규장 총잔량
- `Domestic.InquireOvertimeVolume` — 국내주식 시간외거래량순위 (FHPST02350000) — 거래소/코스닥 합계 + 종목별 최대 30건
- `Domestic.InquireOvertimeFluctuation` — 국내주식 시간외등락율순위 (FHPST02340000) — 상한/상승/보합/하한/하락 통계 + 종목별 최대 30건
- examples: `domestic_extended`

### Notes

- `InquireOvertimeAskingPrice` 응답 struct 는 74 필드 (10단계 × 6 배열 + 합계). 시간외 단일가 최종 시간 (`ovtm_untp_last_hour`) 포함
- Phase 2.2 완료 — 누적 36 메서드 (Phase 2.1: 31 → Phase 2.2: 36)
```

- [ ] **Step 4: domestic/doc.go 갱신**

ADD Phase 2.2 section after Phase 2.1 section (before closing `package domestic` comment):

```go
// Phase 2.2 메서드 (5):
//
//   - InquireNearNewHighlow      — 국내주식 신고/신저근접종목 상위 (FHPST01870000)
//   - InquireOvertimePrice       — 국내주식 시간외현재가 (FHPST02300000)
//   - InquireOvertimeAskingPrice — 국내주식 시간외호가 (FHPST02300400)
//   - InquireOvertimeVolume      — 국내주식 시간외거래량순위 (FHPST02350000)
//   - InquireOvertimeFluctuation — 국내주식 시간외등락율순위 (FHPST02340000)
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
[doc] Phase 2.2 메서드 문서 갱신 — CLAUDE/README/CHANGELOG/domestic/doc.go

Phase 2.2 의 5 메서드 (신고저근접 + 시간외 4종) 목록 + CHANGELOG [1.5.0] entry.
CLAUDE.md banner 갱신 (Phase 2.1 → 2.2, v1.4.0 → v1.5.0). domestic/doc.go
패키지 doc 에 Phase 2.2 section 추가. README Available Methods 31 → 36.

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 9: 최종 점검

- [ ] **Step 1: gofmt cleanup (필요 시)**

`gofmt -w domestic/*.go && gofmt -l .` — empty output.

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
ls -la \
  domestic/extended.go \
  domestic/extended_test.go \
  domestic/testdata/near_new_highlow_success.json \
  domestic/testdata/overtime_price_success.json \
  domestic/testdata/overtime_asking_price_success.json \
  domestic/testdata/overtime_volume_success.json \
  domestic/testdata/overtime_fluctuation_success.json \
  examples/domestic_extended/main.go \
  2>&1 | wc -l
```
Expected: 8 files.

- [ ] **Step 5: Commit history**

`git log main..HEAD --oneline | wc -l` — should be ~9-12.

---

## Task 10: PR 생성 (사용자 승인 후)

> Claude 는 push / PR 생성을 사용자 명시적 승인 후에만 실행 (글로벌 정책).

- [ ] **Step 1: 사용자 승인 요청**

작업 진행 보고 + PR 생성 가능 여부 confirm.

- [ ] **Step 2: Push branch**

`git push -u origin feat/phase2-2-extended`

- [ ] **Step 3: PR 생성**

```bash
gh pr create --title "Phase 2.2 — 국내 신고저/시간외 (v1.5.0)" --reviewer kenshin579 --base main --head feat/phase2-2-extended --body "$(cat <<'EOF'
## Summary

- 국내주식 신고저 근접 + 시간외 5 메서드 추가 (Phase 2 두 번째 sub-phase)
- Phase 1/2.1 패턴 그대로 재사용 (Style A, Params struct, KIS docs 1:1)
- v1.5.0 release 대상 (누적 31 → 36 메서드)

## 메서드 → 한투 API 매핑

| Go 메서드 | path | TR_ID |
|---|---|---|
| InquireNearNewHighlow | ranking/near-new-highlow | FHPST01870000 |
| InquireOvertimePrice | quotations/inquire-overtime-price | FHPST02300000 |
| InquireOvertimeAskingPrice | quotations/inquire-overtime-asking-price | FHPST02300400 |
| InquireOvertimeVolume | ranking/overtime-volume | FHPST02350000 |
| InquireOvertimeFluctuation | ranking/overtime-fluctuation | FHPST02340000 |

## Test Plan

- [x] go build/vet/fmt clean
- [x] go test ./... -race -count=1 모든 패키지 PASS
- [x] Coverage domestic/ >= 80%
- [x] httpmock 단위 테스트 (5 메서드)
- [x] OvertimeAskingPrice 74 필드 (10단계 indexed) 전체 커버

## Breaking Changes

없음 — 신규 메서드 추가만.

## 참고 문서

- Phase 2 design spec: `docs/superpowers/specs/2026-05-05-phase2-readonly-extension-design.md`
- Phase 2.2 implementation plan: `docs/superpowers/specs/2026-05-05-phase2-2-extended-implementation-plan.md`

🤖 Generated with [Claude Code](https://claude.com/claude-code)
EOF
)"
```

- [ ] **Step 4: Merge (사용자 승인 후)** — `gh pr merge <PR#> --merge`

- [ ] **Step 5: 후속 작업**

```bash
git tag -a v1.5.0 -m "Phase 2.2 — 국내 신고저/시간외 (5 메서드, 누적 36)"
git push origin v1.5.0
gh release create v1.5.0 --title "v1.5.0 — Phase 2.2 국내 신고저/시간외" \
  --notes-file <(awk '/^## \[1\.5\.0\]/{p=1} p && /^## \[1\.4\.0\]/{exit} p' CHANGELOG.md)
```
