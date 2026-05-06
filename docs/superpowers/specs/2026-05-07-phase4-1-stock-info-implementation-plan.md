# Phase 4.1 — 종목정보/분석 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** korea-investment-stock Go 라이브러리에 국내주식 종목정보/분석 10 메서드 추가 (`v1.12.0` release). 투자의견 3개(`opinion.go` 신규) + 체결분석/순위 7개(`extended.go` append). TDD: testdata fixture → 실패 테스트 → struct + 메서드 구현 → 통과 → commit.

**Architecture:** Phase 1+2 인프라 + 패턴 재사용. `domestic/opinion.go` 신규 생성 (3 메서드) + 기존 `domestic/extended.go` APPEND (7 메서드). 새 internal package 불필요.

**Tech Stack:** Go 1.25+, `github.com/jarcoal/httpmock`, `github.com/stretchr/testify`, `github.com/shopspring/decimal`. 새 dependency 없음.

**참고 spec:**
- Phase 4 design spec: `docs/superpowers/specs/2026-05-07-phase4-stock-info-design.md` (§Phase 4.1)
- Phase 2.7 plan (참조 패턴 — compact task structure): `docs/superpowers/specs/2026-05-05-phase2-7-industry-implementation-plan.md`
- Phase 3.1 plan (bonds 패턴 참조): `docs/superpowers/specs/2026-05-05-phase3-1-bonds-quote-implementation-plan.md`
- 기존 구현 참조: `domestic/extended.go`, `domestic/info.go`

---

## 사전 정보

| 항목 | 값 |
|------|---|
| 작업 브랜치 | `feat/phase4-1-stock-info` |
| 시작 HEAD | Phase 3.1 구현 완료 commit (v1.11.0) |
| Release 목표 | `v1.12.0` |
| PR 베이스 | `main` |
| 현재 main 상태 | v1.11.0 publish 완료 (Phase 3.1 통합, 79 메서드) |

> **누적 메서드 카운트:** 79 → **89** (10 신규).

---

## 메서드 매핑

| # | Go 메서드 | path (last segment) | TR_ID | output key | fields | anomalies |
|---|---|---|---|---|---|---|
| EP1 | `InquireInvestOpinion` | `quotations/invest-opinion` | FHKST663300C0 | `output []` | 12 | UPPERCASE FID_, date-range pagination |
| EP2 | `InquireInvestOpbysec` | `quotations/invest-opbysec` | FHKST663400C0 | `output []` | 16 | `FID_INPUT_ISCD` = 증권사코드 (NOT 종목코드) |
| EP3 | `InquireEstimatePerform` | `quotations/estimate-perform` | HHKST668300C0 | `output1{}+output2[]+output3[]+output4[]` | 8+5+5+1=19 | **QUAD OUTPUT**, non-FID param `SHT_CD`, doc labels mislabeled (use Python dataclass field names), HH prefix TR_ID |
| EP4 | `InquireVolumePower` | `ranking/volume-power` | FHPST01680000 | `output []` | 11 | **lowercase fid_*** all 9 params |
| EP5 | `InquireBulkTransNum` | `ranking/bulk-trans-num` | FHKST190900C0 | `output []` | 11 | **lowercase fid_*** all 12 params, **`mksc_shrn_iscd`** (not stck_shrn_iscd) |
| EP6 | `InquireTradprtByamt` | `quotations/tradprt-byamt` | FHKST111900C0 | `output []` | 11 | UPPERCASE FID_, **`whol_shun_vol_rate` typo preserved** |
| EP7 | `InquireHtsTopView` | `ranking/hts-top-view` | HHMCM000100C0 | `output1 {}` | 2 | **Zero query params**, `output1` (not `output`), HH prefix |
| EP8 | `InquirePbarTraRatio` | `quotations/pbar-tratio` | FHPST01130000 | `output1{}+output2[]` | 11+4 | UPPERCASE FID_, dual output |
| EP9 | `InquireExpPriceTrend` | `quotations/exp-price-trend` | FHPST01810000 | `output1{}+output2[]` | 7+7 | **lowercase fid_*** params |
| EP10 | `InquireExpTransUpdown` | `ranking/exp-trans-updown` | FHPST01820000 | `output []` | 15 | **lowercase fid_*** all 10 params |

Default `FID_COND_MRKT_DIV_CODE` = `"J"` (KRX standard).

---

## 파일 구조

### 신규 (Go source)
- `domestic/opinion.go` — EP1+EP2+EP3 (3 메서드 + structs + Params)
- `domestic/opinion_test.go` — EP1~EP3 테스트
- `examples/domestic_stock_info/main.go` — 사용 예제

### 수정 (APPEND)
- `domestic/extended.go` — EP4~EP10 (7 메서드 + structs + Params) append
- `domestic/extended_test.go` — EP4~EP10 테스트 append
- `CLAUDE.md` — banner Phase 3.1 → Phase 4.1, plan link 추가
- `README.md` — Available Methods 표 갱신 (79 → 89 메서드)
- `CHANGELOG.md` — `[1.12.0]` entry ABOVE `[1.11.0]`
- `domestic/doc.go` — Phase 4.1 section 추가

### 신규 (testdata — 10 files)
- `domestic/testdata/invest_opinion_success.json`
- `domestic/testdata/invest_opbysec_success.json`
- `domestic/testdata/estimate_perform_success.json`
- `domestic/testdata/volume_power_success.json`
- `domestic/testdata/bulk_trans_num_success.json`
- `domestic/testdata/tradprt_byamt_success.json`
- `domestic/testdata/hts_top_view_success.json`
- `domestic/testdata/pbar_tratio_success.json`
- `domestic/testdata/exp_price_trend_success.json`
- `domestic/testdata/exp_trans_updown_success.json`

---

## 타입 매핑

Phase 2 표준 타입 매핑 — Phase 4.1 종목정보/분석 특화.

| 카테고리 | Go 타입 | json tag suffix | 예시 필드 |
|---|---|---|---|
| 가격 | `decimal.Decimal` | (bare) | `stck_prpr`, `sdpr`, `prdy_vrss`, `hts_goal_prc`, `stck_prdy_clpr`, `stck_nday_esdg`, `stft_esdg`, `antc_cnpr`, `antc_cntg_vrss`, `askp`, `bidp`, `smtn_avrg_prpr`, `wghn_avrg_stck_prc` |
| 수량/금액 | `int64` | `,string` | `acml_vol`, `prdy_vol`, `cntg_vol`, `seln_cnqn_smtn`, `shnu_cnqn_smtn`, `seln_cntg_csnu`, `shnu_cntg_csnu`, `ntby_cnqn`, `ntby_cntg_csnu`, `lstn_stcn`, `seln_rsqn`, `shnu_rsqn`, `total_askp_rsqn`, `total_bidp_rsqn`, `antc_vol`, `antc_tr_pbmn` |
| 비율 | `float64` | `,string` | `prdy_ctrt`, `nday_dprt`, `dprt`, `tday_rltv`, `whol_ntby_qty_rate`, `whol_seln_vol_rate`, `whol_shun_vol_rate`, `acml_vol_rlim`, `antc_cntg_prdy_ctrt` |
| 코드/이름/날짜/Y-N | `string` | (bare) | `stck_bsop_date`, `stck_cntg_hour`, `prdy_vrss_sign`, `mbcr_name`, `hts_kor_isnm`, `stck_shrn_iscd`, `mksc_shrn_iscd`, `mrkt_div_cls_code`, `data_rank`, `invt_opnn`, `invt_opnn_cls_code`, `rgbf_invt_opnn`, `rgbf_invt_opnn_cls_code`, `rprs_mrkt_kor_name`, `prpr_name`, `antc_cntg_vrss_sign` |

> **EP3 anomaly:** `output1` KIS docs body labels MISLABELED. Python dataclass field names 사용: `sht_cd`, `item_kor_nm`, `name1`, `name2`, `estdate`, `rcmd_name`, `capital`, `forn_item_lmtrt` — 모두 `string` (opaque values).
>
> **EP5 anomaly:** 첫 번째 식별 필드는 `mksc_shrn_iscd` (시장구분 포함 단축 종목코드) — `stck_shrn_iscd` 아님.
>
> **EP6 anomaly:** `whol_shun_vol_rate` 는 KIS wire format 그대로 유지 (typo 의도적 보존).

---

## Tasks (13 total)

| # | 내용 | Files |
|---|---|---|
| Task 1 | testdata fixtures (10 합성 JSON) | `domestic/testdata/*.json` |
| Task 2 | EP1 `InquireInvestOpinion` — `opinion.go` 신규 | CREATE `opinion.go` + `opinion_test.go` |
| Task 3 | EP2 `InquireInvestOpbysec` | APPEND `opinion.go` / `opinion_test.go` |
| Task 4 | EP3 `InquireEstimatePerform` (quad-output anomaly) | APPEND `opinion.go` / `opinion_test.go` |
| Task 5 | EP4 `InquireVolumePower` (lowercase params) | APPEND `extended.go` / `extended_test.go` |
| Task 6 | EP5 `InquireBulkTransNum` (lowercase + mksc_shrn_iscd) | APPEND `extended.go` / `extended_test.go` |
| Task 7 | EP6 `InquireTradprtByamt` (whol_shun_vol_rate typo) | APPEND `extended.go` / `extended_test.go` |
| Task 8 | EP7 `InquireHtsTopView` (zero query params, output1) | APPEND `extended.go` / `extended_test.go` |
| Task 9 | EP8 `InquirePbarTraRatio` (dual output) | APPEND `extended.go` / `extended_test.go` |
| Task 10 | EP9 `InquireExpPriceTrend` (lowercase + dual output) | APPEND `extended.go` / `extended_test.go` |
| Task 11 | EP10 `InquireExpTransUpdown` (lowercase params) | APPEND `extended.go` / `extended_test.go` |
| Task 12 | examples `examples/domestic_stock_info/main.go` | CREATE |
| Task 13 | `domestic/doc.go` + `CHANGELOG.md` + `CLAUDE.md` + `README.md` 갱신 + `go test ./...` + tag | MODIFY |

---

## Task 1: testdata fixtures (10 합성 JSON)

- [ ] Step 1: `domestic/testdata/invest_opinion_success.json` (EP1 — `output []`, 12 fields)
- [ ] Step 2: `domestic/testdata/invest_opbysec_success.json` (EP2 — `output []`, 16 fields)
- [ ] Step 3: `domestic/testdata/estimate_perform_success.json` (EP3 — quad output: `output1{}` 8 + `output2[]` 5 + `output3[]` 5 + `output4[]` 1)
- [ ] Step 4: `domestic/testdata/volume_power_success.json` (EP4 — `output []`, 11 fields, lowercase params fixture)
- [ ] Step 5: `domestic/testdata/bulk_trans_num_success.json` (EP5 — `output []`, 11 fields, mksc_shrn_iscd)
- [ ] Step 6: `domestic/testdata/tradprt_byamt_success.json` (EP6 — `output []`, 11 fields, whol_shun_vol_rate)
- [ ] Step 7: `domestic/testdata/hts_top_view_success.json` (EP7 — `output1{}`, 2 fields)
- [ ] Step 8: `domestic/testdata/pbar_tratio_success.json` (EP8 — `output1{}` 11 + `output2[]` 4)
- [ ] Step 9: `domestic/testdata/exp_price_trend_success.json` (EP9 — `output1{}` 7 + `output2[]` 7)
- [ ] Step 10: `domestic/testdata/exp_trans_updown_success.json` (EP10 — `output []`, 15 fields)
- [ ] Step 11: validation

```bash
for f in \
  domestic/testdata/invest_opinion_success.json \
  domestic/testdata/invest_opbysec_success.json \
  domestic/testdata/estimate_perform_success.json \
  domestic/testdata/volume_power_success.json \
  domestic/testdata/bulk_trans_num_success.json \
  domestic/testdata/tradprt_byamt_success.json \
  domestic/testdata/hts_top_view_success.json \
  domestic/testdata/pbar_tratio_success.json \
  domestic/testdata/exp_price_trend_success.json \
  domestic/testdata/exp_trans_updown_success.json; do
  python3 -c "import json; json.load(open('$f'))" && echo "$f OK"
done
# Expected: 10 OK lines
```

- [ ] Step 12: commit

```bash
git commit -m "$(cat <<'EOF'
[chore] testdata — 10 stock-info fixture JSON (Phase 4.1)

합성 JSON fixtures (2 records each where applicable):
- invest_opinion_success.json (output[] 12 fields, date-range)
- invest_opbysec_success.json (output[] 16 fields, 증권사코드)
- estimate_perform_success.json (quad output1/2/3/4, SHT_CD param)
- volume_power_success.json (output[] 11 fields, lowercase fid_* params)
- bulk_trans_num_success.json (output[] 11 fields, mksc_shrn_iscd)
- tradprt_byamt_success.json (output[] 11 fields, whol_shun_vol_rate typo)
- hts_top_view_success.json (output1{} 2 fields, zero query params)
- pbar_tratio_success.json (output1{} 11 + output2[] 4 fields)
- exp_price_trend_success.json (output1{} 7 + output2[] 7 fields)
- exp_trans_updown_success.json (output[] 15 fields, lowercase fid_*)

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

### Fixture content

**Step 1 — `invest_opinion_success.json`** (EP1: output[] 12 fields, date-range pagination)

```json
{
  "rt_cd": "0",
  "msg_cd": "MCA00000",
  "msg1": "정상처리 되었습니다.",
  "output": [
    {
      "stck_bsop_date": "20260506",
      "invt_opnn": "매수",
      "invt_opnn_cls_code": "1",
      "rgbf_invt_opnn": "매수",
      "rgbf_invt_opnn_cls_code": "1",
      "mbcr_name": "삼성증권",
      "hts_goal_prc": "95000",
      "stck_prdy_clpr": "82500",
      "stck_nday_esdg": "93000",
      "nday_dprt": "12.73",
      "stft_esdg": "90000",
      "dprt": "9.09"
    },
    {
      "stck_bsop_date": "20260430",
      "invt_opnn": "매수",
      "invt_opnn_cls_code": "1",
      "rgbf_invt_opnn": "중립",
      "rgbf_invt_opnn_cls_code": "3",
      "mbcr_name": "미래에셋증권",
      "hts_goal_prc": "92000",
      "stck_prdy_clpr": "81000",
      "stck_nday_esdg": "90000",
      "nday_dprt": "11.11",
      "stft_esdg": "88000",
      "dprt": "8.64"
    }
  ]
}
```

**Step 2 — `invest_opbysec_success.json`** (EP2: output[] 16 fields, 증권사코드 조회)

```json
{
  "rt_cd": "0",
  "msg_cd": "MCA00000",
  "msg1": "정상처리 되었습니다.",
  "output": [
    {
      "stck_bsop_date": "20260506",
      "stck_shrn_iscd": "005930",
      "hts_kor_isnm": "삼성전자",
      "invt_opnn": "매수",
      "invt_opnn_cls_code": "1",
      "rgbf_invt_opnn": "매수",
      "rgbf_invt_opnn_cls_code": "1",
      "mbcr_name": "삼성증권",
      "stck_prpr": "82500",
      "prdy_vrss": "500",
      "prdy_vrss_sign": "2",
      "prdy_ctrt": "0.61",
      "hts_goal_prc": "95000",
      "stck_prdy_clpr": "82000",
      "stft_esdg": "90000",
      "dprt": "9.76"
    },
    {
      "stck_bsop_date": "20260505",
      "stck_shrn_iscd": "000660",
      "hts_kor_isnm": "SK하이닉스",
      "invt_opnn": "매수",
      "invt_opnn_cls_code": "1",
      "rgbf_invt_opnn": "매수",
      "rgbf_invt_opnn_cls_code": "1",
      "mbcr_name": "삼성증권",
      "stck_prpr": "185000",
      "prdy_vrss": "2000",
      "prdy_vrss_sign": "2",
      "prdy_ctrt": "1.09",
      "hts_goal_prc": "220000",
      "stck_prdy_clpr": "183000",
      "stft_esdg": "210000",
      "dprt": "14.75"
    }
  ]
}
```

**Step 3 — `estimate_perform_success.json`** (EP3: quad output — output1{} 8 + output2[] 5 per row + output3[] 5 per row + output4[] 1 field)

```json
{
  "rt_cd": "0",
  "msg_cd": "MCA00000",
  "msg1": "정상처리 되었습니다.",
  "output1": {
    "sht_cd": "005930",
    "item_kor_nm": "삼성전자",
    "name1": "82500",
    "name2": "500",
    "estdate": "20251231",
    "rcmd_name": "매수",
    "capital": "7780",
    "forn_item_lmtrt": "52.00"
  },
  "output2": [
    {
      "data1": "매출액",
      "data2": "305000",
      "data3": "320000",
      "data4": "335000",
      "data5": "350000"
    },
    {
      "data1": "매출액증감율",
      "data2": "5.20",
      "data3": "4.92",
      "data4": "4.69",
      "data5": "4.48"
    },
    {
      "data1": "영업이익",
      "data2": "35000",
      "data3": "40000",
      "data4": "45000",
      "data5": "50000"
    },
    {
      "data1": "영업이익증감율",
      "data2": "12.50",
      "data3": "14.29",
      "data4": "12.50",
      "data5": "11.11"
    },
    {
      "data1": "순이익",
      "data2": "28000",
      "data3": "32000",
      "data4": "36000",
      "data5": "40000"
    },
    {
      "data1": "순이익증감율",
      "data2": "10.00",
      "data3": "14.29",
      "data4": "12.50",
      "data5": "11.11"
    }
  ],
  "output3": [
    {
      "data1": "EBITDA",
      "data2": "55000",
      "data3": "62000",
      "data4": "68000",
      "data5": "74000"
    },
    {
      "data1": "EPS",
      "data2": "4800",
      "data3": "5500",
      "data4": "6200",
      "data5": "6900"
    },
    {
      "data1": "EPS증감율",
      "data2": "8.00",
      "data3": "14.58",
      "data4": "12.73",
      "data5": "11.29"
    },
    {
      "data1": "PER",
      "data2": "17.19",
      "data3": "15.00",
      "data4": "13.31",
      "data5": "11.96"
    },
    {
      "data1": "EV-EBITDA",
      "data2": "8.50",
      "data3": "7.80",
      "data4": "7.20",
      "data5": "6.70"
    },
    {
      "data1": "ROE",
      "data2": "15.20",
      "data3": "16.50",
      "data4": "17.80",
      "data5": "19.00"
    },
    {
      "data1": "부채비율",
      "data2": "35.00",
      "data3": "33.00",
      "data4": "31.00",
      "data5": "29.00"
    },
    {
      "data1": "이자보상배율",
      "data2": "42.00",
      "data3": "48.00",
      "data4": "55.00",
      "data5": "62.00"
    }
  ],
  "output4": [
    {"dt": "202412"},
    {"dt": "202512"},
    {"dt": "202612E"},
    {"dt": "202712E"},
    {"dt": "202812E"}
  ]
}
```

**Step 4 — `volume_power_success.json`** (EP4: output[] 11 fields, lowercase fid_* wire params)

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
      "stck_prpr": "82500",
      "prdy_vrss": "500",
      "prdy_vrss_sign": "2",
      "prdy_ctrt": "0.61",
      "acml_vol": "12500000",
      "tday_rltv": "125.30",
      "seln_cnqn_smtn": "6800000",
      "shnu_cnqn_smtn": "7200000"
    },
    {
      "stck_shrn_iscd": "000660",
      "data_rank": "2",
      "hts_kor_isnm": "SK하이닉스",
      "stck_prpr": "185000",
      "prdy_vrss": "2000",
      "prdy_vrss_sign": "2",
      "prdy_ctrt": "1.09",
      "acml_vol": "8300000",
      "tday_rltv": "118.70",
      "seln_cnqn_smtn": "4200000",
      "shnu_cnqn_smtn": "4500000"
    }
  ]
}
```

**Step 5 — `bulk_trans_num_success.json`** (EP5: output[] 11 fields, mksc_shrn_iscd)

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
      "stck_prpr": "82500",
      "prdy_vrss_sign": "2",
      "prdy_vrss": "500",
      "prdy_ctrt": "0.61",
      "acml_vol": "12500000",
      "shnu_cntg_csnu": "3200",
      "seln_cntg_csnu": "2800",
      "ntby_cnqn": "400000"
    },
    {
      "mksc_shrn_iscd": "000660",
      "data_rank": "2",
      "hts_kor_isnm": "SK하이닉스",
      "stck_prpr": "185000",
      "prdy_vrss_sign": "2",
      "prdy_vrss": "2000",
      "prdy_ctrt": "1.09",
      "acml_vol": "8300000",
      "shnu_cntg_csnu": "1900",
      "seln_cntg_csnu": "1750",
      "ntby_cnqn": "150000"
    }
  ]
}
```

**Step 6 — `tradprt_byamt_success.json`** (EP6: output[] 11 fields, whol_shun_vol_rate typo preserved)

```json
{
  "rt_cd": "0",
  "msg_cd": "MCA00000",
  "msg1": "정상처리 되었습니다.",
  "output": [
    {
      "prpr_name": "1억원 이상",
      "smtn_avrg_prpr": "150000000",
      "acml_vol": "8500",
      "whol_ntby_qty_rate": "12.50",
      "ntby_cntg_csnu": "320",
      "seln_cnqn_smtn": "620000000",
      "whol_seln_vol_rate": "42.30",
      "seln_cntg_csnu": "1800",
      "shnu_cnqn_smtn": "740000000",
      "whol_shun_vol_rate": "45.20",
      "shnu_cntg_csnu": "2120"
    },
    {
      "prpr_name": "5천만원 이상",
      "smtn_avrg_prpr": "65000000",
      "acml_vol": "15200",
      "whol_ntby_qty_rate": "8.30",
      "ntby_cntg_csnu": "180",
      "seln_cnqn_smtn": "480000000",
      "whol_seln_vol_rate": "32.80",
      "seln_cntg_csnu": "3200",
      "shnu_cnqn_smtn": "492000000",
      "whol_shun_vol_rate": "33.60",
      "shnu_cntg_csnu": "3380"
    }
  ]
}
```

**Step 7 — `hts_top_view_success.json`** (EP7: output1{} 2 fields, zero query params, HH prefix TR_ID)

```json
{
  "rt_cd": "0",
  "msg_cd": "MCA00000",
  "msg1": "정상처리 되었습니다.",
  "output1": {
    "stck_shrn_iscd": "005930",
    "hts_kor_isnm": "삼성전자"
  }
}
```

**Step 8 — `pbar_tratio_success.json`** (EP8: output1{} 11 + output2[] 4 fields)

```json
{
  "rt_cd": "0",
  "msg_cd": "MCA00000",
  "msg1": "정상처리 되었습니다.",
  "output1": {
    "stck_prpr": "82500",
    "prdy_vrss": "500",
    "prdy_vrss_sign": "2",
    "prdy_ctrt": "0.61",
    "acml_vol": "12500000",
    "total_askp_rsqn": "3500000",
    "total_bidp_rsqn": "4200000",
    "askp": "82600",
    "bidp": "82400",
    "antc_cnpr": "82500",
    "antc_cntg_vrss": "500"
  },
  "output2": [
    {
      "stck_bsop_date": "20260506",
      "acml_vol": "12500000",
      "seln_rsqn": "3500000",
      "shnu_rsqn": "4200000"
    },
    {
      "stck_bsop_date": "20260505",
      "acml_vol": "11800000",
      "seln_rsqn": "3200000",
      "shnu_rsqn": "4000000"
    }
  ]
}
```

**Step 9 — `exp_price_trend_success.json`** (EP9: output1{} 7 + output2[] 7 fields, lowercase params)

```json
{
  "rt_cd": "0",
  "msg_cd": "MCA00000",
  "msg1": "정상처리 되었습니다.",
  "output1": {
    "antc_cnpr": "82700",
    "antc_cntg_vrss": "200",
    "antc_cntg_vrss_sign": "2",
    "antc_cntg_prdy_ctrt": "0.24",
    "antc_vol": "850000",
    "antc_tr_pbmn": "70297500000",
    "rprs_mrkt_kor_name": "코스피"
  },
  "output2": [
    {
      "stck_cntg_hour": "153000",
      "antc_cnpr": "82700",
      "antc_cntg_vrss": "200",
      "antc_cntg_vrss_sign": "2",
      "antc_cntg_prdy_ctrt": "0.24",
      "antc_vol": "850000",
      "acml_vol_rlim": "6.80"
    },
    {
      "stck_cntg_hour": "152900",
      "antc_cnpr": "82600",
      "antc_cntg_vrss": "100",
      "antc_cntg_vrss_sign": "2",
      "antc_cntg_prdy_ctrt": "0.12",
      "antc_vol": "820000",
      "acml_vol_rlim": "6.56"
    }
  ]
}
```

**Step 10 — `exp_trans_updown_success.json`** (EP10: output[] 15 fields, lowercase params)

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
      "stck_prpr": "82500",
      "prdy_vrss": "500",
      "prdy_vrss_sign": "2",
      "prdy_ctrt": "0.61",
      "antc_cnpr": "82700",
      "antc_cntg_vrss": "200",
      "antc_cntg_vrss_sign": "2",
      "antc_cntg_prdy_ctrt": "0.24",
      "antc_vol": "850000",
      "antc_tr_pbmn": "70297500000",
      "lstn_stcn": "5969782550",
      "mrkt_div_cls_code": "0"
    },
    {
      "stck_shrn_iscd": "000660",
      "data_rank": "2",
      "hts_kor_isnm": "SK하이닉스",
      "stck_prpr": "185000",
      "prdy_vrss": "2000",
      "prdy_vrss_sign": "2",
      "prdy_ctrt": "1.09",
      "antc_cnpr": "186500",
      "antc_cntg_vrss": "1500",
      "antc_cntg_vrss_sign": "2",
      "antc_cntg_prdy_ctrt": "0.81",
      "antc_vol": "420000",
      "antc_tr_pbmn": "78330000000",
      "lstn_stcn": "728002365",
      "mrkt_div_cls_code": "0"
    }
  ]
}
```

---

## Task 2: InquireInvestOpinion (EP1)

**Files:** CREATE `domestic/opinion.go` + CREATE `domestic/opinion_test.go`

종목투자의견 조회. `output []` (12 fields, 날짜 범위 페이지네이션).

- [ ] Step 1: CREATE `domestic/opinion_test.go` (패키지 선언 + import + 첫 테스트 함수)
- [ ] Step 2: Verify FAIL — `go test ./domestic/... -run TestClient_InquireInvestOpinion -v` (compile error expected)
- [ ] Step 3: CREATE `domestic/opinion.go` (패키지 선언 + import + struct + Params + 메서드)
- [ ] Step 4: Verify PASS — `go test ./domestic/... -run TestClient_InquireInvestOpinion -v`
- [ ] Step 5: `gofmt -w domestic/opinion.go domestic/opinion_test.go && go vet ./domestic/...`
- [ ] Step 6: commit

### 메서드 매핑
- Path: `/uapi/domestic-stock/v1/quotations/invest-opinion`
- TR_ID: `FHKST663300C0`
- Params (5): `MarketCode` (`FID_COND_MRKT_DIV_CODE` default `"J"`), `CondScrDivCode` (`FID_COND_SCR_DIV_CODE` default `"16633"`), `Symbol` (`FID_INPUT_ISCD` 필수), `StartDate` (`FID_INPUT_DATE_1` YYYYMMDD), `EndDate` (`FID_INPUT_DATE_2` YYYYMMDD)

### output struct 필드 (12 fields)

| Go field | json tag | Go type | 설명 |
|---|---|---|---|
| `StckBsopDate` | `stck_bsop_date` | `string` | 영업 일자 |
| `InvtOpnn` | `invt_opnn` | `string` | 투자의견 |
| `InvtOpnnClsCode` | `invt_opnn_cls_code` | `string` | 투자의견 구분 코드 |
| `RgbfInvtOpnn` | `rgbf_invt_opnn` | `string` | 직전 투자의견 |
| `RgbfInvtOpnnClsCode` | `rgbf_invt_opnn_cls_code` | `string` | 직전 투자의견 구분 코드 |
| `MbcrName` | `mbcr_name` | `string` | 회원사명 (증권사) |
| `HtsGoalPrc` | `hts_goal_prc` | `decimal.Decimal` | HTS 목표 가격 |
| `StckPrdyClpr` | `stck_prdy_clpr` | `decimal.Decimal` | 주식 전일 종가 |
| `StckNdayEsdg` | `stck_nday_esdg` | `decimal.Decimal` | 주식 N일 추정 단가 |
| `NdayDprt` | `nday_dprt` | `float64,string` | N일 이격도 |
| `StftEsdg` | `stft_esdg` | `decimal.Decimal` | 직전 추정 단가 |
| `Dprt` | `dprt` | `float64,string` | 이격도 |

### 구현 코드

```go
// opinion.go

package domestic

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kenshin579/korea-investment-stock/internal/httpclient"
	"github.com/shopspring/decimal"
)

// InvestOpinion 은 종목투자의견 (FHKST663300C0) 응답.
//
// 한투 docs: docs/api/국내주식/종목투자의견.md
// path: /uapi/domestic-stock/v1/quotations/invest-opinion
type InvestOpinion struct {
	Output []InvestOpinionItem `json:"output"`
}

// InvestOpinionItem 은 응답의 output 한 행 (12 fields).
type InvestOpinionItem struct {
	StckBsopDate        string          `json:"stck_bsop_date"`           // 영업 일자
	InvtOpnn            string          `json:"invt_opnn"`                // 투자의견
	InvtOpnnClsCode     string          `json:"invt_opnn_cls_code"`       // 투자의견 구분 코드
	RgbfInvtOpnn        string          `json:"rgbf_invt_opnn"`           // 직전 투자의견
	RgbfInvtOpnnClsCode string          `json:"rgbf_invt_opnn_cls_code"`  // 직전 투자의견 구분 코드
	MbcrName            string          `json:"mbcr_name"`                // 회원사명
	HtsGoalPrc          decimal.Decimal `json:"hts_goal_prc"`             // HTS 목표 가격
	StckPrdyClpr        decimal.Decimal `json:"stck_prdy_clpr"`           // 주식 전일 종가
	StckNdayEsdg        decimal.Decimal `json:"stck_nday_esdg"`           // 주식 N일 추정 단가
	NdayDprt            float64         `json:"nday_dprt,string"`         // N일 이격도
	StftEsdg            decimal.Decimal `json:"stft_esdg"`                // 직전 추정 단가
	Dprt                float64         `json:"dprt,string"`              // 이격도
}

// InquireInvestOpinionParams 는 종목투자의견 조회 파라미터.
type InquireInvestOpinionParams struct {
	MarketCode     string // FID_COND_MRKT_DIV_CODE — 빈 값=>"J"
	CondScrDivCode string // FID_COND_SCR_DIV_CODE — 빈 값=>"16633"
	Symbol         string // FID_INPUT_ISCD — 필수, 단축 종목코드 (예 "005930")
	StartDate      string // FID_INPUT_DATE_1 — 조회 시작일 YYYYMMDD
	EndDate        string // FID_INPUT_DATE_2 — 조회 종료일 YYYYMMDD
}

// InquireInvestOpinion 은 종목투자의견 호출.
//
// 한투 docs: docs/api/국내주식/종목투자의견.md
// path: /uapi/domestic-stock/v1/quotations/invest-opinion (FHKST663300C0)
func (c *Client) InquireInvestOpinion(ctx context.Context, params InquireInvestOpinionParams) (*InvestOpinion, error) {
	market := params.MarketCode
	if market == "" {
		market = "J"
	}
	scrDiv := params.CondScrDivCode
	if scrDiv == "" {
		scrDiv = "16633"
	}
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/domestic-stock/v1/quotations/invest-opinion",
		TrID:   "FHKST663300C0",
		Query: map[string]string{
			"FID_COND_MRKT_DIV_CODE": market,
			"FID_COND_SCR_DIV_CODE":  scrDiv,
			"FID_INPUT_ISCD":         params.Symbol,
			"FID_INPUT_DATE_1":       params.StartDate,
			"FID_INPUT_DATE_2":       params.EndDate,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}
	var res InvestOpinion
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse InvestOpinion: %w", err)
	}
	return &res, nil
}
```

### 테스트 코드

```go
// opinion_test.go

package domestic_test

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/kenshin579/korea-investment-stock/domestic"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_InquireInvestOpinion(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/quotations/invest-opinion`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "invest_opinion_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireInvestOpinion(context.Background(), domestic.InquireInvestOpinionParams{
		Symbol:    "005930",
		StartDate: "20260401",
		EndDate:   "20260506",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, "J", capturedQuery.Get("FID_COND_MRKT_DIV_CODE"))
	assert.Equal(t, "16633", capturedQuery.Get("FID_COND_SCR_DIV_CODE"))
	assert.Equal(t, "005930", capturedQuery.Get("FID_INPUT_ISCD"))
	assert.Equal(t, "20260401", capturedQuery.Get("FID_INPUT_DATE_1"))
	assert.Equal(t, "20260506", capturedQuery.Get("FID_INPUT_DATE_2"))

	require.Len(t, res.Output, 2)
	assert.Equal(t, "20260506", res.Output[0].StckBsopDate)
	assert.Equal(t, "삼성증권", res.Output[0].MbcrName)
	assert.Equal(t, "매수", res.Output[0].InvtOpnn)

	wantGoal, _ := decimal.NewFromString("95000")
	assert.True(t, wantGoal.Equal(res.Output[0].HtsGoalPrc))
	assert.InDelta(t, 12.73, res.Output[0].NdayDprt, 0.001)
	assert.InDelta(t, 9.09, res.Output[0].Dprt, 0.001)
}
```

### commit

```bash
git commit -m "$(cat <<'EOF'
[feat] domestic/opinion.go — InquireInvestOpinion (EP1, Phase 4.1)

종목투자의견 조회 (FHKST663300C0). output[] 12 fields.
날짜 범위 페이지네이션 파라미터 지원 (FID_INPUT_DATE_1/2).

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 3: InquireInvestOpbysec (EP2)

**Files:** APPEND to `domestic/opinion.go` and `domestic/opinion_test.go`

증권사별 투자의견 조회. `output []` (16 fields). `FID_INPUT_ISCD` = 증권사코드 (종목코드 아님 — 주의).

- [ ] Step 1: APPEND test code to `domestic/opinion_test.go`
- [ ] Step 2: Verify FAIL — `go test ./domestic/... -run TestClient_InquireInvestOpbysec -v`
- [ ] Step 3: APPEND struct + Params + method to `domestic/opinion.go`
- [ ] Step 4: Verify PASS — `go test ./domestic/... -run TestClient_InquireInvestOpbysec -v`
- [ ] Step 5: `gofmt -w domestic/opinion.go domestic/opinion_test.go && go vet ./domestic/...`
- [ ] Step 6: commit

### 메서드 매핑
- Path: `/uapi/domestic-stock/v1/quotations/invest-opbysec`
- TR_ID: `FHKST663400C0`
- Params (6): `MarketCode` (`FID_COND_MRKT_DIV_CODE` default `"J"`), `CondScrDivCode` (`FID_COND_SCR_DIV_CODE` default `"16634"`), `SecBrokerCode` (`FID_INPUT_ISCD` **증권사코드** 필수), `DivClsCode` (`FID_DIV_CLS_CODE` 0=전체), `StartDate` (`FID_INPUT_DATE_1` YYYYMMDD), `EndDate` (`FID_INPUT_DATE_2` YYYYMMDD)

### output struct 필드 (16 fields)

| Go field | json tag | Go type | 설명 |
|---|---|---|---|
| `StckBsopDate` | `stck_bsop_date` | `string` | 영업 일자 |
| `StckShrnIscd` | `stck_shrn_iscd` | `string` | 주식 단축 종목코드 |
| `HtsKorIsnm` | `hts_kor_isnm` | `string` | HTS 한글 종목명 |
| `InvtOpnn` | `invt_opnn` | `string` | 투자의견 |
| `InvtOpnnClsCode` | `invt_opnn_cls_code` | `string` | 투자의견 구분 코드 |
| `RgbfInvtOpnn` | `rgbf_invt_opnn` | `string` | 직전 투자의견 |
| `RgbfInvtOpnnClsCode` | `rgbf_invt_opnn_cls_code` | `string` | 직전 투자의견 구분 코드 |
| `MbcrName` | `mbcr_name` | `string` | 회원사명 |
| `StckPrpr` | `stck_prpr` | `decimal.Decimal` | 주식 현재가 |
| `PrdyVrss` | `prdy_vrss` | `decimal.Decimal` | 전일 대비 |
| `PrdyVrssSign` | `prdy_vrss_sign` | `string` | 전일 대비 부호 |
| `PrdyCtrt` | `prdy_ctrt` | `float64,string` | 전일 대비율 |
| `HtsGoalPrc` | `hts_goal_prc` | `decimal.Decimal` | HTS 목표 가격 |
| `StckPrdyClpr` | `stck_prdy_clpr` | `decimal.Decimal` | 주식 전일 종가 |
| `StftEsdg` | `stft_esdg` | `decimal.Decimal` | 직전 추정 단가 |
| `Dprt` | `dprt` | `float64,string` | 이격도 |

### 구현 코드

```go
// InvestOpbysec 은 증권사별 투자의견 (FHKST663400C0) 응답.
//
// 한투 docs: docs/api/국내주식/증권사별투자의견.md
// path: /uapi/domestic-stock/v1/quotations/invest-opbysec
type InvestOpbysec struct {
	Output []InvestOpbysecItem `json:"output"`
}

// InvestOpbysecItem 은 응답의 output 한 행 (16 fields).
type InvestOpbysecItem struct {
	StckBsopDate        string          `json:"stck_bsop_date"`           // 영업 일자
	StckShrnIscd        string          `json:"stck_shrn_iscd"`           // 주식 단축 종목코드
	HtsKorIsnm          string          `json:"hts_kor_isnm"`             // HTS 한글 종목명
	InvtOpnn            string          `json:"invt_opnn"`                // 투자의견
	InvtOpnnClsCode     string          `json:"invt_opnn_cls_code"`       // 투자의견 구분 코드
	RgbfInvtOpnn        string          `json:"rgbf_invt_opnn"`           // 직전 투자의견
	RgbfInvtOpnnClsCode string          `json:"rgbf_invt_opnn_cls_code"`  // 직전 투자의견 구분 코드
	MbcrName            string          `json:"mbcr_name"`                // 회원사명
	StckPrpr            decimal.Decimal `json:"stck_prpr"`                // 주식 현재가
	PrdyVrss            decimal.Decimal `json:"prdy_vrss"`                // 전일 대비
	PrdyVrssSign        string          `json:"prdy_vrss_sign"`           // 전일 대비 부호
	PrdyCtrt            float64         `json:"prdy_ctrt,string"`         // 전일 대비율
	HtsGoalPrc          decimal.Decimal `json:"hts_goal_prc"`             // HTS 목표 가격
	StckPrdyClpr        decimal.Decimal `json:"stck_prdy_clpr"`           // 주식 전일 종가
	StftEsdg            decimal.Decimal `json:"stft_esdg"`                // 직전 추정 단가
	Dprt                float64         `json:"dprt,string"`              // 이격도
}

// InquireInvestOpbysecParams 는 증권사별 투자의견 조회 파라미터.
type InquireInvestOpbysecParams struct {
	MarketCode     string // FID_COND_MRKT_DIV_CODE — 빈 값=>"J"
	CondScrDivCode string // FID_COND_SCR_DIV_CODE — 빈 값=>"16634"
	SecBrokerCode  string // FID_INPUT_ISCD — 필수, 증권사코드 (종목코드 아님!)
	DivClsCode     string // FID_DIV_CLS_CODE — 0=전체
	StartDate      string // FID_INPUT_DATE_1 — 조회 시작일 YYYYMMDD
	EndDate        string // FID_INPUT_DATE_2 — 조회 종료일 YYYYMMDD
}

// InquireInvestOpbysec 은 증권사별 투자의견 호출.
//
// 한투 docs: docs/api/국내주식/증권사별투자의견.md
// path: /uapi/domestic-stock/v1/quotations/invest-opbysec (FHKST663400C0)
//
// 주의: FID_INPUT_ISCD 는 종목코드가 아닌 증권사코드를 입력한다.
func (c *Client) InquireInvestOpbysec(ctx context.Context, params InquireInvestOpbysecParams) (*InvestOpbysec, error) {
	market := params.MarketCode
	if market == "" {
		market = "J"
	}
	scrDiv := params.CondScrDivCode
	if scrDiv == "" {
		scrDiv = "16634"
	}
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/domestic-stock/v1/quotations/invest-opbysec",
		TrID:   "FHKST663400C0",
		Query: map[string]string{
			"FID_COND_MRKT_DIV_CODE": market,
			"FID_COND_SCR_DIV_CODE":  scrDiv,
			"FID_INPUT_ISCD":         params.SecBrokerCode,
			"FID_DIV_CLS_CODE":       params.DivClsCode,
			"FID_INPUT_DATE_1":       params.StartDate,
			"FID_INPUT_DATE_2":       params.EndDate,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}
	var res InvestOpbysec
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse InvestOpbysec: %w", err)
	}
	return &res, nil
}
```

### 테스트 코드

```go
func TestClient_InquireInvestOpbysec(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/quotations/invest-opbysec`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "invest_opbysec_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireInvestOpbysec(context.Background(), domestic.InquireInvestOpbysecParams{
		SecBrokerCode: "240",  // 삼성증권 코드 예시
		DivClsCode:    "0",
		StartDate:     "20260401",
		EndDate:       "20260506",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	// 증권사코드가 FID_INPUT_ISCD 로 전송되는지 확인
	assert.Equal(t, "J", capturedQuery.Get("FID_COND_MRKT_DIV_CODE"))
	assert.Equal(t, "16634", capturedQuery.Get("FID_COND_SCR_DIV_CODE"))
	assert.Equal(t, "240", capturedQuery.Get("FID_INPUT_ISCD"))
	assert.Equal(t, "0", capturedQuery.Get("FID_DIV_CLS_CODE"))

	require.Len(t, res.Output, 2)
	assert.Equal(t, "005930", res.Output[0].StckShrnIscd)
	assert.Equal(t, "삼성전자", res.Output[0].HtsKorIsnm)
	assert.Equal(t, "삼성증권", res.Output[0].MbcrName)

	wantPrpr, _ := decimal.NewFromString("82500")
	assert.True(t, wantPrpr.Equal(res.Output[0].StckPrpr))
	assert.InDelta(t, 0.61, res.Output[0].PrdyCtrt, 0.001)
	assert.InDelta(t, 9.76, res.Output[0].Dprt, 0.001)
}
```

### commit

```bash
git commit -m "$(cat <<'EOF'
[feat] domestic/opinion.go — InquireInvestOpbysec (EP2, Phase 4.1)

증권사별 투자의견 조회 (FHKST663400C0). output[] 16 fields.
FID_INPUT_ISCD = 증권사코드 (종목코드 아님) 주의사항 문서화.

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 4: InquireEstimatePerform (EP3)

**Files:** APPEND to `domestic/opinion.go` and `domestic/opinion_test.go`

종목 추정실적 조회. **QUAD OUTPUT**: `output1{}` (8 fields) + `output2[]` (5 fields × 6 rows 추정손익계산서) + `output3[]` (5 fields × 8 rows 투자지표) + `output4[]` (1 field × 5 rows 결산년월). TR_ID HH prefix (`HHKST668300C0`). Param 이름 `SHT_CD` (non-FID, 비표준).

> **CRITICAL:** KIS docs body table은 output1 labels이 MISLABELED. Python dataclass field names 그대로 사용: `sht_cd`, `item_kor_nm`, `name1`, `name2`, `estdate`, `rcmd_name`, `capital`, `forn_item_lmtrt`. 모두 `string` (opaque — 실제 의미: 단축종목코드/한글종목명/현재가/전일대비/결산일자/투자의견명/자본금or거래량/외국인한도비율).

- [ ] Step 1: APPEND test code to `domestic/opinion_test.go`
- [ ] Step 2: Verify FAIL — `go test ./domestic/... -run TestClient_InquireEstimatePerform -v`
- [ ] Step 3: APPEND structs + Params + method to `domestic/opinion.go`
- [ ] Step 4: Verify PASS — `go test ./domestic/... -run TestClient_InquireEstimatePerform -v`
- [ ] Step 5: `gofmt -w domestic/opinion.go domestic/opinion_test.go && go vet ./domestic/...`
- [ ] Step 6: commit

### 메서드 매핑
- Path: `/uapi/domestic-stock/v1/quotations/estimate-perform`
- TR_ID: `HHKST668300C0`
- Params (1): `Symbol` (`SHT_CD` 필수, 6자리 단축 종목코드 — **non-FID 파라미터명**)

### output1 struct 필드 (8 fields — EstimatePerformSummary)

| Go field | json tag | Go type | 실제 의미 (Python dataclass 기준) |
|---|---|---|---|
| `ShtCd` | `sht_cd` | `string` | 단축 종목코드 |
| `ItemKorNm` | `item_kor_nm` | `string` | 한글 종목명 |
| `Name1` | `name1` | `string` | 현재가 (opaque string) |
| `Name2` | `name2` | `string` | 전일 대비 (opaque string) |
| `Estdate` | `estdate` | `string` | 결산 일자 |
| `RcmdName` | `rcmd_name` | `string` | 투자의견명 |
| `Capital` | `capital` | `string` | 자본금 또는 거래량 (opaque) |
| `FornItemLmtrt` | `forn_item_lmtrt` | `string` | 외국인 한도 비율 (opaque) |

### output2 / output3 struct 필드 (5 fields per row — EstimatePerformRow)

output2 (추정손익계산서 6 rows) 와 output3 (투자지표 8 rows) 는 동일 shape:

| Go field | json tag | Go type | 설명 |
|---|---|---|---|
| `Data1` | `data1` | `string` | 항목명 (레이블) |
| `Data2` | `data2` | `string` | 기간1 값 |
| `Data3` | `data3` | `string` | 기간2 값 |
| `Data4` | `data4` | `string` | 기간3 값 |
| `Data5` | `data5` | `string` | 기간4 값 |

### output4 struct 필드 (1 field per row — EstimatePerformPeriod)

| Go field | json tag | Go type | 설명 |
|---|---|---|---|
| `Dt` | `dt` | `string` | 결산 년월 (예 "202412", "202612E") |

### 구현 코드

```go
// EstimatePerform 은 종목 추정실적 (HHKST668300C0) 응답.
//
// 한투 docs: docs/api/국내주식/종목추정실적.md
// path: /uapi/domestic-stock/v1/quotations/estimate-perform
//
// QUAD OUTPUT: output1(요약) + output2(추정손익계산서 6행) + output3(투자지표 8행) + output4(결산년월 5행).
// output1 KIS docs body labels 오표기 — Python dataclass field names 사용.
type EstimatePerform struct {
	Output1 EstimatePerformSummary    `json:"output1"`
	Output2 []EstimatePerformRow      `json:"output2"`
	Output3 []EstimatePerformRow      `json:"output3"`
	Output4 []EstimatePerformPeriod   `json:"output4"`
}

// EstimatePerformSummary 는 응답의 output1 (종목 요약, 8 fields).
// KIS docs body table 오표기 — Python dataclass 기준 field names 사용.
type EstimatePerformSummary struct {
	ShtCd         string `json:"sht_cd"`           // 단축 종목코드
	ItemKorNm     string `json:"item_kor_nm"`       // 한글 종목명
	Name1         string `json:"name1"`             // 현재가 (opaque)
	Name2         string `json:"name2"`             // 전일 대비 (opaque)
	Estdate       string `json:"estdate"`           // 결산 일자
	RcmdName      string `json:"rcmd_name"`         // 투자의견명
	Capital       string `json:"capital"`           // 자본금/거래량 (opaque)
	FornItemLmtrt string `json:"forn_item_lmtrt"`   // 외국인 한도 비율 (opaque)
}

// EstimatePerformRow 는 output2(추정손익계산서) / output3(투자지표) 한 행 (5 fields).
type EstimatePerformRow struct {
	Data1 string `json:"data1"` // 항목명
	Data2 string `json:"data2"` // 기간1 값
	Data3 string `json:"data3"` // 기간2 값
	Data4 string `json:"data4"` // 기간3 값
	Data5 string `json:"data5"` // 기간4 값
}

// EstimatePerformPeriod 는 output4 한 행 (결산년월, 1 field).
type EstimatePerformPeriod struct {
	Dt string `json:"dt"` // 결산 년월 (예 "202412", "202612E")
}

// InquireEstimatePerformParams 는 종목 추정실적 조회 파라미터.
type InquireEstimatePerformParams struct {
	Symbol string // SHT_CD — 필수, 6자리 단축 종목코드 (비표준 param명, FID_ 접두어 없음)
}

// InquireEstimatePerform 은 종목 추정실적 호출.
//
// 한투 docs: docs/api/국내주식/종목추정실적.md
// path: /uapi/domestic-stock/v1/quotations/estimate-perform (HHKST668300C0)
//
// 주의: query param 이름이 SHT_CD (FID_ 접두어 없음).
func (c *Client) InquireEstimatePerform(ctx context.Context, params InquireEstimatePerformParams) (*EstimatePerform, error) {
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/domestic-stock/v1/quotations/estimate-perform",
		TrID:   "HHKST668300C0",
		Query: map[string]string{
			"SHT_CD": params.Symbol,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}
	var res EstimatePerform
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse EstimatePerform: %w", err)
	}
	return &res, nil
}
```

### 테스트 코드

```go
func TestClient_InquireEstimatePerform(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/quotations/estimate-perform`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "estimate_perform_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireEstimatePerform(context.Background(), domestic.InquireEstimatePerformParams{
		Symbol: "005930",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	// non-FID param 이름 확인
	assert.Equal(t, "005930", capturedQuery.Get("SHT_CD"))
	assert.Empty(t, capturedQuery.Get("FID_INPUT_ISCD"), "SHT_CD 파라미터만 사용해야 함")

	// output1 확인
	assert.Equal(t, "005930", res.Output1.ShtCd)
	assert.Equal(t, "삼성전자", res.Output1.ItemKorNm)
	assert.Equal(t, "매수", res.Output1.RcmdName)

	// output2 (추정손익계산서 6행) 확인
	require.Len(t, res.Output2, 6)
	assert.Equal(t, "매출액", res.Output2[0].Data1)
	assert.Equal(t, "305000", res.Output2[0].Data2)

	// output3 (투자지표 8행) 확인
	require.Len(t, res.Output3, 8)
	assert.Equal(t, "EBITDA", res.Output3[0].Data1)

	// output4 (결산년월 5행) 확인
	require.Len(t, res.Output4, 5)
	assert.Equal(t, "202412", res.Output4[0].Dt)
	assert.Equal(t, "202612E", res.Output4[2].Dt)
}
```

### commit

```bash
git commit -m "$(cat <<'EOF'
[feat] domestic/opinion.go — InquireEstimatePerform (EP3, Phase 4.1)

종목 추정실적 조회 (HHKST668300C0). quad output 구조:
output1(요약 8fields) + output2(추정손익계산서 6행×5fields) +
output3(투자지표 8행×5fields) + output4(결산년월 5행×1field).
SHT_CD 비표준 param명, KIS docs mislabeled output1 주의사항 문서화.

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 5: InquireVolumePower (EP4)

**Files:** APPEND to `domestic/extended.go` and `domestic/extended_test.go`

체결강도 상위 조회 (ranking). `output []` (11 fields). **모든 9개 query params가 lowercase `fid_*`** (UPPERCASE FID_ 아님 — anomaly).

- [ ] Step 1: APPEND test code to `domestic/extended_test.go`
- [ ] Step 2: Verify FAIL — `go test ./domestic/... -run TestClient_InquireVolumePower -v`
- [ ] Step 3: APPEND struct + Params + method to `domestic/extended.go`
- [ ] Step 4: Verify PASS — `go test ./domestic/... -run TestClient_InquireVolumePower -v`
- [ ] Step 5: `gofmt -w domestic/extended.go domestic/extended_test.go && go vet ./domestic/...`
- [ ] Step 6: commit

### 메서드 매핑
- Path: `/uapi/domestic-stock/v1/ranking/volume-power`
- TR_ID: `FHPST01680000`
- Params (9, **모두 lowercase wire keys**):
  - `MarketCode` → `fid_cond_mrkt_div_code` (default `"J"`)
  - `CondScrDivCode` → `fid_cond_scr_div_code` (default `"20168"`)
  - `Symbol` → `fid_input_iscd` (0000:전체 / 0001:코스피 / 1001:코스닥)
  - `DivClsCode` → `fid_div_cls_code`
  - `Price1` → `fid_input_price_1`
  - `Price2` → `fid_input_price_2`
  - `VolCnt` → `fid_vol_cnt`
  - `TrgtClsCode` → `fid_trgt_cls_code`
  - `TrgtExlsCode` → `fid_trgt_exls_cls_code`

### output struct 필드 (11 fields)

| Go field | json tag | Go type | 설명 |
|---|---|---|---|
| `StckShrnIscd` | `stck_shrn_iscd` | `string` | 주식 단축 종목코드 |
| `DataRank` | `data_rank` | `string` | 데이터 순위 |
| `HtsKorIsnm` | `hts_kor_isnm` | `string` | HTS 한글 종목명 |
| `StckPrpr` | `stck_prpr` | `decimal.Decimal` | 주식 현재가 |
| `PrdyVrss` | `prdy_vrss` | `decimal.Decimal` | 전일 대비 |
| `PrdyVrssSign` | `prdy_vrss_sign` | `string` | 전일 대비 부호 |
| `PrdyCtrt` | `prdy_ctrt` | `float64,string` | 전일 대비율 |
| `AcmlVol` | `acml_vol` | `int64,string` | 누적 거래량 |
| `TdayRltv` | `tday_rltv` | `float64,string` | 당일 체결강도 |
| `SelnCnqnSmtn` | `seln_cnqn_smtn` | `int64,string` | 매도 체결량 합계 |
| `ShnuCnqnSmtn` | `shnu_cnqn_smtn` | `int64,string` | 매수 체결량 합계 |

### 구현 코드

```go
// VolumePower 는 체결강도 상위 (FHPST01680000) 응답.
//
// 한투 docs: docs/api/국내주식/체결강도상위.md
// path: /uapi/domestic-stock/v1/ranking/volume-power
//
// 주의: 모든 query 파라미터가 lowercase fid_* (대문자 FID_ 아님).
type VolumePower struct {
	Output []VolumePowerItem `json:"output"`
}

// VolumePowerItem 은 응답의 output 한 행 (11 fields).
type VolumePowerItem struct {
	StckShrnIscd string          `json:"stck_shrn_iscd"`    // 주식 단축 종목코드
	DataRank     string          `json:"data_rank"`         // 데이터 순위
	HtsKorIsnm   string          `json:"hts_kor_isnm"`      // HTS 한글 종목명
	StckPrpr     decimal.Decimal `json:"stck_prpr"`         // 주식 현재가
	PrdyVrss     decimal.Decimal `json:"prdy_vrss"`         // 전일 대비
	PrdyVrssSign string          `json:"prdy_vrss_sign"`    // 전일 대비 부호
	PrdyCtrt     float64         `json:"prdy_ctrt,string"`  // 전일 대비율
	AcmlVol      int64           `json:"acml_vol,string"`   // 누적 거래량
	TdayRltv     float64         `json:"tday_rltv,string"`  // 당일 체결강도
	SelnCnqnSmtn int64           `json:"seln_cnqn_smtn,string"` // 매도 체결량 합계
	ShnuCnqnSmtn int64           `json:"shnu_cnqn_smtn,string"` // 매수 체결량 합계
}

// InquireVolumePowerParams 는 체결강도 상위 조회 파라미터.
type InquireVolumePowerParams struct {
	MarketCode     string // fid_cond_mrkt_div_code — 빈 값=>"J" (lowercase wire key)
	CondScrDivCode string // fid_cond_scr_div_code — 빈 값=>"20168"
	Symbol         string // fid_input_iscd — 0000:전체/0001:코스피/1001:코스닥
	DivClsCode     string // fid_div_cls_code
	Price1         string // fid_input_price_1
	Price2         string // fid_input_price_2
	VolCnt         string // fid_vol_cnt
	TrgtClsCode    string // fid_trgt_cls_code
	TrgtExlsCode   string // fid_trgt_exls_cls_code
}

// InquireVolumePower 는 체결강도 상위 호출.
//
// 한투 docs: docs/api/국내주식/체결강도상위.md
// path: /uapi/domestic-stock/v1/ranking/volume-power (FHPST01680000)
func (c *Client) InquireVolumePower(ctx context.Context, params InquireVolumePowerParams) (*VolumePower, error) {
	market := params.MarketCode
	if market == "" {
		market = "J"
	}
	scrDiv := params.CondScrDivCode
	if scrDiv == "" {
		scrDiv = "20168"
	}
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/domestic-stock/v1/ranking/volume-power",
		TrID:   "FHPST01680000",
		Query: map[string]string{
			"fid_cond_mrkt_div_code":  market,
			"fid_cond_scr_div_code":   scrDiv,
			"fid_input_iscd":          params.Symbol,
			"fid_div_cls_code":        params.DivClsCode,
			"fid_input_price_1":       params.Price1,
			"fid_input_price_2":       params.Price2,
			"fid_vol_cnt":             params.VolCnt,
			"fid_trgt_cls_code":       params.TrgtClsCode,
			"fid_trgt_exls_cls_code":  params.TrgtExlsCode,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}
	var res VolumePower
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse VolumePower: %w", err)
	}
	return &res, nil
}
```

### 테스트 코드

```go
func TestClient_InquireVolumePower(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/ranking/volume-power`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "volume_power_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireVolumePower(context.Background(), domestic.InquireVolumePowerParams{
		Symbol:     "0001",
		DivClsCode: "0",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	// lowercase wire keys 확인 (UPPERCASE FID_ 아님)
	assert.Equal(t, "J", capturedQuery.Get("fid_cond_mrkt_div_code"))
	assert.Equal(t, "20168", capturedQuery.Get("fid_cond_scr_div_code"))
	assert.Equal(t, "0001", capturedQuery.Get("fid_input_iscd"))
	assert.Empty(t, capturedQuery.Get("FID_COND_MRKT_DIV_CODE"), "lowercase 사용 확인")

	require.Len(t, res.Output, 2)
	assert.Equal(t, "005930", res.Output[0].StckShrnIscd)
	assert.Equal(t, "1", res.Output[0].DataRank)

	wantPrpr, _ := decimal.NewFromString("82500")
	assert.True(t, wantPrpr.Equal(res.Output[0].StckPrpr))
	assert.InDelta(t, 125.30, res.Output[0].TdayRltv, 0.001)
	assert.Equal(t, int64(6800000), res.Output[0].SelnCnqnSmtn)
	assert.Equal(t, int64(7200000), res.Output[0].ShnuCnqnSmtn)
}
```

### commit

```bash
git commit -m "$(cat <<'EOF'
[feat] domestic/extended.go — InquireVolumePower (EP4, Phase 4.1)

체결강도 상위 조회 (FHPST01680000). output[] 11 fields.
모든 9개 query params lowercase fid_* anomaly 반영.

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 6: InquireBulkTransNum (EP5)

**Files:** APPEND to `domestic/extended.go` and `domestic/extended_test.go`

대량체결건수 상위 조회 (ranking). `output []` (11 fields). **모든 12개 query params lowercase `fid_*`**. 첫 번째 종목코드 필드는 **`mksc_shrn_iscd`** (시장구분 포함, `stck_shrn_iscd` 아님).

- [ ] Step 1: APPEND test code to `domestic/extended_test.go`
- [ ] Step 2: Verify FAIL — `go test ./domestic/... -run TestClient_InquireBulkTransNum -v`
- [ ] Step 3: APPEND struct + Params + method to `domestic/extended.go`
- [ ] Step 4: Verify PASS — `go test ./domestic/... -run TestClient_InquireBulkTransNum -v`
- [ ] Step 5: `gofmt -w domestic/extended.go domestic/extended_test.go && go vet ./domestic/...`
- [ ] Step 6: commit

### 메서드 매핑
- Path: `/uapi/domestic-stock/v1/ranking/bulk-trans-num`
- TR_ID: `FHKST190900C0`
- Params (12, **모두 lowercase wire keys**):
  - `MarketCode` → `fid_cond_mrkt_div_code` (default `"J"`)
  - `CondScrDivCode` → `fid_cond_scr_div_code` (default `"11909"`)
  - `Symbol` → `fid_input_iscd`
  - `DivClsCode` → `fid_div_cls_code`
  - `RankSortCode` → `fid_rank_sort_cls_code`
  - `BlngClsCode` → `fid_blng_cls_code`
  - `TrgtClsCode` → `fid_trgt_cls_code`
  - `TrgtExlsCode` → `fid_trgt_exls_cls_code`
  - `InputPrice1` → `fid_input_price_1`
  - `InputPrice2` → `fid_input_price_2`
  - `VolCnt` → `fid_vol_cnt`
  - `InputDate1` → `fid_input_date_1`

### output struct 필드 (11 fields)

| Go field | json tag | Go type | 설명 |
|---|---|---|---|
| `MkscShrnIscd` | `mksc_shrn_iscd` | `string` | 시장구분+단축 종목코드 (**stck_shrn_iscd 아님**) |
| `DataRank` | `data_rank` | `string` | 데이터 순위 |
| `HtsKorIsnm` | `hts_kor_isnm` | `string` | HTS 한글 종목명 |
| `StckPrpr` | `stck_prpr` | `decimal.Decimal` | 주식 현재가 |
| `PrdyVrssSign` | `prdy_vrss_sign` | `string` | 전일 대비 부호 |
| `PrdyVrss` | `prdy_vrss` | `decimal.Decimal` | 전일 대비 |
| `PrdyCtrt` | `prdy_ctrt` | `float64,string` | 전일 대비율 |
| `AcmlVol` | `acml_vol` | `int64,string` | 누적 거래량 |
| `ShnuCntgCsnu` | `shnu_cntg_csnu` | `int64,string` | 매수 체결 건수 |
| `SelnCntgCsnu` | `seln_cntg_csnu` | `int64,string` | 매도 체결 건수 |
| `NtbyCnqn` | `ntby_cnqn` | `int64,string` | 순매수 체결량 |

### 구현 코드

```go
// BulkTransNum 은 대량체결건수 상위 (FHKST190900C0) 응답.
//
// 한투 docs: docs/api/국내주식/대량체결건수상위.md
// path: /uapi/domestic-stock/v1/ranking/bulk-trans-num
//
// 주의 1: 모든 query 파라미터가 lowercase fid_* (대문자 FID_ 아님).
// 주의 2: 종목코드 필드는 mksc_shrn_iscd (시장구분 포함, stck_shrn_iscd 아님).
type BulkTransNum struct {
	Output []BulkTransNumItem `json:"output"`
}

// BulkTransNumItem 은 응답의 output 한 행 (11 fields).
type BulkTransNumItem struct {
	MkscShrnIscd string          `json:"mksc_shrn_iscd"`      // 시장구분+단축 종목코드
	DataRank     string          `json:"data_rank"`           // 데이터 순위
	HtsKorIsnm   string          `json:"hts_kor_isnm"`        // HTS 한글 종목명
	StckPrpr     decimal.Decimal `json:"stck_prpr"`           // 주식 현재가
	PrdyVrssSign string          `json:"prdy_vrss_sign"`      // 전일 대비 부호
	PrdyVrss     decimal.Decimal `json:"prdy_vrss"`           // 전일 대비
	PrdyCtrt     float64         `json:"prdy_ctrt,string"`    // 전일 대비율
	AcmlVol      int64           `json:"acml_vol,string"`     // 누적 거래량
	ShnuCntgCsnu int64           `json:"shnu_cntg_csnu,string"` // 매수 체결 건수
	SelnCntgCsnu int64           `json:"seln_cntg_csnu,string"` // 매도 체결 건수
	NtbyCnqn     int64           `json:"ntby_cnqn,string"`    // 순매수 체결량
}

// InquireBulkTransNumParams 는 대량체결건수 상위 조회 파라미터.
type InquireBulkTransNumParams struct {
	MarketCode     string // fid_cond_mrkt_div_code — 빈 값=>"J"
	CondScrDivCode string // fid_cond_scr_div_code — 빈 값=>"11909"
	Symbol         string // fid_input_iscd
	DivClsCode     string // fid_div_cls_code
	RankSortCode   string // fid_rank_sort_cls_code
	BlngClsCode    string // fid_blng_cls_code
	TrgtClsCode    string // fid_trgt_cls_code
	TrgtExlsCode   string // fid_trgt_exls_cls_code
	InputPrice1    string // fid_input_price_1
	InputPrice2    string // fid_input_price_2
	VolCnt         string // fid_vol_cnt
	InputDate1     string // fid_input_date_1
}

// InquireBulkTransNum 은 대량체결건수 상위 호출.
//
// 한투 docs: docs/api/국내주식/대량체결건수상위.md
// path: /uapi/domestic-stock/v1/ranking/bulk-trans-num (FHKST190900C0)
func (c *Client) InquireBulkTransNum(ctx context.Context, params InquireBulkTransNumParams) (*BulkTransNum, error) {
	market := params.MarketCode
	if market == "" {
		market = "J"
	}
	scrDiv := params.CondScrDivCode
	if scrDiv == "" {
		scrDiv = "11909"
	}
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/domestic-stock/v1/ranking/bulk-trans-num",
		TrID:   "FHKST190900C0",
		Query: map[string]string{
			"fid_cond_mrkt_div_code":  market,
			"fid_cond_scr_div_code":   scrDiv,
			"fid_input_iscd":          params.Symbol,
			"fid_div_cls_code":        params.DivClsCode,
			"fid_rank_sort_cls_code":  params.RankSortCode,
			"fid_blng_cls_code":       params.BlngClsCode,
			"fid_trgt_cls_code":       params.TrgtClsCode,
			"fid_trgt_exls_cls_code":  params.TrgtExlsCode,
			"fid_input_price_1":       params.InputPrice1,
			"fid_input_price_2":       params.InputPrice2,
			"fid_vol_cnt":             params.VolCnt,
			"fid_input_date_1":        params.InputDate1,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}
	var res BulkTransNum
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse BulkTransNum: %w", err)
	}
	return &res, nil
}
```

### 테스트 코드

```go
func TestClient_InquireBulkTransNum(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/ranking/bulk-trans-num`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "bulk_trans_num_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireBulkTransNum(context.Background(), domestic.InquireBulkTransNumParams{
		Symbol:       "0000",
		DivClsCode:   "0",
		RankSortCode: "0",
		BlngClsCode:  "0",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	// lowercase wire keys 확인
	assert.Equal(t, "J", capturedQuery.Get("fid_cond_mrkt_div_code"))
	assert.Equal(t, "11909", capturedQuery.Get("fid_cond_scr_div_code"))
	assert.Equal(t, "0000", capturedQuery.Get("fid_input_iscd"))

	require.Len(t, res.Output, 2)
	// mksc_shrn_iscd (stck_shrn_iscd 아님) 확인
	assert.Equal(t, "005930", res.Output[0].MkscShrnIscd)
	assert.Equal(t, "1", res.Output[0].DataRank)
	assert.Equal(t, int64(3200), res.Output[0].ShnuCntgCsnu)
	assert.Equal(t, int64(2800), res.Output[0].SelnCntgCsnu)
	assert.Equal(t, int64(400000), res.Output[0].NtbyCnqn)
}
```

### commit

```bash
git commit -m "$(cat <<'EOF'
[feat] domestic/extended.go — InquireBulkTransNum (EP5, Phase 4.1)

대량체결건수 상위 조회 (FHKST190900C0). output[] 11 fields.
lowercase fid_* anomaly + mksc_shrn_iscd (not stck_shrn_iscd) 반영.

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 7: InquireTradprtByamt (EP6)

**Files:** APPEND to `domestic/extended.go` and `domestic/extended_test.go`

체결금액별 매매비중 조회. `output []` (11 fields). UPPERCASE FID_ params. **`whol_shun_vol_rate` 필드명은 KIS wire format typo — 그대로 보존** (수정 금지).

- [ ] Step 1: APPEND test code to `domestic/extended_test.go`
- [ ] Step 2: Verify FAIL — `go test ./domestic/... -run TestClient_InquireTradprtByamt -v`
- [ ] Step 3: APPEND struct + Params + method to `domestic/extended.go`
- [ ] Step 4: Verify PASS — `go test ./domestic/... -run TestClient_InquireTradprtByamt -v`
- [ ] Step 5: `gofmt -w domestic/extended.go domestic/extended_test.go && go vet ./domestic/...`
- [ ] Step 6: commit

### 메서드 매핑
- Path: `/uapi/domestic-stock/v1/quotations/tradprt-byamt`
- TR_ID: `FHKST111900C0`
- Params (3, UPPERCASE): `MarketCode` (`FID_COND_MRKT_DIV_CODE` default `"J"`), `CondScrDivCode` (`FID_COND_SCR_DIV_CODE` default `"11119"`), `Symbol` (`FID_INPUT_ISCD`)

### output struct 필드 (11 fields)

| Go field | json tag | Go type | 설명 |
|---|---|---|---|
| `PrprName` | `prpr_name` | `string` | 체결금액 구간명 |
| `SmtnAvrgPrpr` | `smtn_avrg_prpr` | `decimal.Decimal` | 합산 평균 가격 |
| `AcmlVol` | `acml_vol` | `int64,string` | 누적 거래량 |
| `WholNtbyQtyRate` | `whol_ntby_qty_rate` | `float64,string` | 전체 순매수 수량 비율 |
| `NtbyCntgCsnu` | `ntby_cntg_csnu` | `int64,string` | 순매수 체결 건수 |
| `SelnCnqnSmtn` | `seln_cnqn_smtn` | `int64,string` | 매도 체결량 합계 |
| `WholSelnVolRate` | `whol_seln_vol_rate` | `float64,string` | 전체 매도 거래량 비율 |
| `SelnCntgCsnu` | `seln_cntg_csnu` | `int64,string` | 매도 체결 건수 |
| `ShnuCnqnSmtn` | `shnu_cnqn_smtn` | `int64,string` | 매수 체결량 합계 |
| `WholShunVolRate` | `whol_shun_vol_rate` | `float64,string` | 전체 매수 거래량 비율 (**typo 보존**: shun not shnu) |
| `ShnuCntgCsnu` | `shnu_cntg_csnu` | `int64,string` | 매수 체결 건수 |

### 구현 코드

```go
// TradprtByamt 는 체결금액별 매매비중 (FHKST111900C0) 응답.
//
// 한투 docs: docs/api/국내주식/체결금액별매매비중.md
// path: /uapi/domestic-stock/v1/quotations/tradprt-byamt
//
// 주의: whol_shun_vol_rate 필드명은 KIS wire format typo (shun ≠ shnu).
// 실제 의미는 전체 매수 거래량 비율이나 KIS wire key 그대로 보존.
type TradprtByamt struct {
	Output []TradprtByamtItem `json:"output"`
}

// TradprtByamtItem 은 응답의 output 한 행 (11 fields).
type TradprtByamtItem struct {
	PrprName        string          `json:"prpr_name"`              // 체결금액 구간명
	SmtnAvrgPrpr    decimal.Decimal `json:"smtn_avrg_prpr"`         // 합산 평균 가격
	AcmlVol         int64           `json:"acml_vol,string"`        // 누적 거래량
	WholNtbyQtyRate float64         `json:"whol_ntby_qty_rate,string"` // 전체 순매수 수량 비율
	NtbyCntgCsnu    int64           `json:"ntby_cntg_csnu,string"`  // 순매수 체결 건수
	SelnCnqnSmtn    int64           `json:"seln_cnqn_smtn,string"`  // 매도 체결량 합계
	WholSelnVolRate float64         `json:"whol_seln_vol_rate,string"` // 전체 매도 거래량 비율
	SelnCntgCsnu    int64           `json:"seln_cntg_csnu,string"`  // 매도 체결 건수
	ShnuCnqnSmtn    int64           `json:"shnu_cnqn_smtn,string"`  // 매수 체결량 합계
	WholShunVolRate float64         `json:"whol_shun_vol_rate,string"` // 전체 매수 거래량 비율 (KIS typo 보존)
	ShnuCntgCsnu    int64           `json:"shnu_cntg_csnu,string"`  // 매수 체결 건수
}

// InquireTradprtByamtParams 는 체결금액별 매매비중 조회 파라미터.
type InquireTradprtByamtParams struct {
	MarketCode     string // FID_COND_MRKT_DIV_CODE — 빈 값=>"J"
	CondScrDivCode string // FID_COND_SCR_DIV_CODE — 빈 값=>"11119"
	Symbol         string // FID_INPUT_ISCD — 필수, 단축 종목코드
}

// InquireTradprtByamt 는 체결금액별 매매비중 호출.
//
// 한투 docs: docs/api/국내주식/체결금액별매매비중.md
// path: /uapi/domestic-stock/v1/quotations/tradprt-byamt (FHKST111900C0)
func (c *Client) InquireTradprtByamt(ctx context.Context, params InquireTradprtByamtParams) (*TradprtByamt, error) {
	market := params.MarketCode
	if market == "" {
		market = "J"
	}
	scrDiv := params.CondScrDivCode
	if scrDiv == "" {
		scrDiv = "11119"
	}
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/domestic-stock/v1/quotations/tradprt-byamt",
		TrID:   "FHKST111900C0",
		Query: map[string]string{
			"FID_COND_MRKT_DIV_CODE": market,
			"FID_COND_SCR_DIV_CODE":  scrDiv,
			"FID_INPUT_ISCD":         params.Symbol,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}
	var res TradprtByamt
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse TradprtByamt: %w", err)
	}
	return &res, nil
}
```

### 테스트 코드

```go
func TestClient_InquireTradprtByamt(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/quotations/tradprt-byamt`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "tradprt_byamt_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireTradprtByamt(context.Background(), domestic.InquireTradprtByamtParams{
		Symbol: "005930",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, "J", capturedQuery.Get("FID_COND_MRKT_DIV_CODE"))
	assert.Equal(t, "11119", capturedQuery.Get("FID_COND_SCR_DIV_CODE"))
	assert.Equal(t, "005930", capturedQuery.Get("FID_INPUT_ISCD"))

	require.Len(t, res.Output, 2)
	assert.Equal(t, "1억원 이상", res.Output[0].PrprName)
	assert.Equal(t, int64(8500), res.Output[0].AcmlVol)
	assert.InDelta(t, 12.50, res.Output[0].WholNtbyQtyRate, 0.001)
	// whol_shun_vol_rate typo 필드 확인
	assert.InDelta(t, 45.20, res.Output[0].WholShunVolRate, 0.001)
	assert.InDelta(t, 42.30, res.Output[0].WholSelnVolRate, 0.001)

	wantAvrgPrpr, _ := decimal.NewFromString("150000000")
	assert.True(t, wantAvrgPrpr.Equal(res.Output[0].SmtnAvrgPrpr))
}
```

### commit

```bash
git commit -m "$(cat <<'EOF'
[feat] domestic/extended.go — InquireTradprtByamt (EP6, Phase 4.1)

체결금액별 매매비중 조회 (FHKST111900C0). output[] 11 fields.
whol_shun_vol_rate KIS wire format typo 의도적 보존.

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

> **Note**: 총 태스크 수는 **15개** (초기 개요의 13에서 변경). Tasks 14-15는 최종 QA + PR 생성 (Phase 2.x 패턴과 동일).

---

## Task 8: InquireHtsTopView (EP7)

**파일**: `domestic/extended.go` (append)
**엔드포인트**: `GET /uapi/domestic-stock/v1/ranking/hts-top-view`
**TR_ID**: `HHMCM000100C0`
**특이사항**: 쿼리 파라미터 없음 (zero params), 응답 key가 `output1` (not `output`)

### 타입 정의

```go
type HtsTopView struct {
    Output1 HtsTopViewItem `json:"output1"`
}

type HtsTopViewItem struct {
    MrktDivClsCode string `json:"mrkt_div_cls_code"` // 시장구분 (J:코스피, Q:코스닥)
    MkscShrnIscd   string `json:"mksc_shrn_iscd"`    // 종목코드
}

type InquireHtsTopViewParams struct {
    // No fields — endpoint takes no query parameters
}
```

### 메서드 구현

```go
func (c *Client) InquireHtsTopView(ctx context.Context, _ InquireHtsTopViewParams) (*HtsTopView, error) {
    resp, err := c.http.Do(ctx, &httpclient.Request{
        Method:   http.MethodGet,
        Path:     "/uapi/domestic-stock/v1/ranking/hts-top-view",
        TrID:     "HHMCM000100C0",
        Query:    map[string]string{},
        CustType: "P",
    })
    if err != nil {
        return nil, err
    }
    var res HtsTopView
    if err := json.Unmarshal(resp.Raw, &res); err != nil {
        return nil, fmt.Errorf("kis: parse HtsTopView: %w", err)
    }
    return &res, nil
}
```

### 테스트

- `testdata/domestic_hts_top_view.json` fixture 사용
- 쿼리 파라미터가 전송되지 않음을 assert
- `output1.mrkt_div_cls_code`, `output1.mksc_shrn_iscd` 필드 값 검증

### 스텝

1. `testdata/domestic_hts_top_view.json` fixture 작성 (output1, 2 fields)
2. `domestic/extended.go` 에 타입 + 메서드 append
3. `domestic/extended_test.go` 에 TestInquireHtsTopView 추가 (zero-params + output1 anomaly 검증)
4. `go build ./... && go vet ./...` — silent
5. `go test ./domestic/... -run TestInquireHtsTopView -v` — PASS
6. `gofmt -w domestic/extended.go domestic/extended_test.go`

### commit

```bash
git commit -m "$(cat <<'EOF'
[feat] domestic — InquireHtsTopView (HTS조회상위20종목, HHMCM000100C0)

zero query params + output1 (not output) 응답 구조 anomaly 처리.
HHMCM prefix TR_ID (HH 시리즈).

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 9: InquirePbarTraRatio (EP8)

**파일**: `domestic/extended.go` (append)
**엔드포인트**: `GET /uapi/domestic-stock/v1/quotations/pbar-tratio`
**TR_ID**: `FHPST01130000`
**출력**: `output1 PbarTraRatioSummary` (11 fields) + `output2 []PbarTraRatioItem` (4 fields)

### 파라미터 (4개)

| Go 필드 | 와이어 키 | 기본값 | 설명 |
|---------|-----------|--------|------|
| MarketCode | `FID_COND_MRKT_DIV_CODE` | `"J"` | 시장구분 |
| CondScrDivCode | `FID_COND_SCR_DIV_CODE` | `"11130"` | 조건화면분류코드 |
| Symbol | `FID_INPUT_ISCD` | — | 종목코드 |
| InputHour1 | `FID_INPUT_HOUR_1` | — | 입력시간1 |

### 타입 정의

```go
// output1
type PbarTraRatioSummary struct {
    RprsMrktKorName string          `json:"rprs_mrkt_kor_name"`
    StckShrnIscd    string          `json:"stck_shrn_iscd"`
    HtsKorIsnm      string          `json:"hts_kor_isnm"`
    StckPrpr        decimal.Decimal `json:"stck_prpr"`
    PrdyVrssSign    string          `json:"prdy_vrss_sign"`
    PrdyVrss        decimal.Decimal `json:"prdy_vrss"`
    PrdyCtrt        string          `json:"prdy_ctrt"`   // float64 as string
    AcmlVol         string          `json:"acml_vol"`    // int64 as string
    PrdyVol         string          `json:"prdy_vol"`    // int64 as string
    WghnAvrgStckPrc decimal.Decimal `json:"wghn_avrg_stck_prc"` // 가중평균주식가격
    LstnStcn        string          `json:"lstn_stcn"`   // 상장주수, int64 as string
}

// output2 item
type PbarTraRatioItem struct {
    DataRank  string          `json:"data_rank"`
    StckPrpr  decimal.Decimal `json:"stck_prpr"`
    CntgVol   string          `json:"cntg_vol"`    // int64 as string
    AcmlVolRlim string        `json:"acml_vol_rlim"` // float64 as string
}

type PbarTraRatio struct {
    Output1 PbarTraRatioSummary `json:"output1"`
    Output2 []PbarTraRatioItem  `json:"output2"`
}

type InquirePbarTraRatioParams struct {
    MarketCode     string // FID_COND_MRKT_DIV_CODE, default "J"
    CondScrDivCode string // FID_COND_SCR_DIV_CODE, default "11130"
    Symbol         string // FID_INPUT_ISCD
    InputHour1     string // FID_INPUT_HOUR_1
}
```

### 스텝

1. `testdata/domestic_pbar_tra_ratio.json` fixture 작성 (output1 + output2 배열)
2. `domestic/extended.go` 에 타입 + 메서드 append
3. `domestic/extended_test.go` 에 TestInquirePbarTraRatio 추가 (dual-output 검증)
4. `go build ./... && go vet ./...` — silent
5. `go test ./domestic/... -run TestInquirePbarTraRatio -v` — PASS
6. `gofmt -w domestic/extended.go domestic/extended_test.go`

### commit

```bash
git commit -m "$(cat <<'EOF'
[feat] domestic — InquirePbarTraRatio (체결금액별 매매비중, FHPST01130000)

dual output (output1 PbarTraRatioSummary 11 fields + output2 []PbarTraRatioItem 4 fields).
wghn_avrg_stck_prc (가중평균주식가격) + lstn_stcn (상장주수) 포함.

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 10: InquireExpPriceTrend (EP9)

**파일**: `domestic/extended.go` (append)
**엔드포인트**: `GET /uapi/domestic-stock/v1/quotations/exp-price-trend`
**TR_ID**: `FHPST01810000`
**특이사항**: 쿼리 파라미터 와이어 키가 **소문자 fid_***
**출력**: `output1 ExpPriceTrendSummary` (7 fields) + `output2 []ExpPriceTrendItem` (7 fields)

### 파라미터 (3개, 소문자 와이어 키)

| Go 필드 | 와이어 키 (소문자) | 기본값 | 설명 |
|---------|-------------------|--------|------|
| MarketCode | `fid_cond_mrkt_div_code` | `"J"` | 시장구분 |
| CondScrDivCode | `fid_cond_scr_div_code` | `"11810"` | 조건화면분류코드 |
| Symbol | `fid_input_iscd` | — | 종목코드 |

### 타입 정의

```go
// output1
type ExpPriceTrendSummary struct {
    RprsMrktKorName    string          `json:"rprs_mrkt_kor_name"`
    AntcCnpr           decimal.Decimal `json:"antc_cnpr"`           // 예상체결가
    AntcCntgVrssSign   string          `json:"antc_cntg_vrss_sign"`
    AntcCntgVrss       decimal.Decimal `json:"antc_cntg_vrss"`
    AntcCntgPrdyCtrt   string          `json:"antc_cntg_prdy_ctrt"` // float64 as string
    AntcVol            string          `json:"antc_vol"`             // int64 as string
    AntcTrPbmn         string          `json:"antc_tr_pbmn"`         // int64 as string
}

// output2 item
type ExpPriceTrendItem struct {
    StckBsopDate  string          `json:"stck_bsop_date"`
    StckCntgHour  string          `json:"stck_cntg_hour"`
    StckPrpr      decimal.Decimal `json:"stck_prpr"`
    PrdyVrssSign  string          `json:"prdy_vrss_sign"`
    PrdyVrss      decimal.Decimal `json:"prdy_vrss"`
    PrdyCtrt      string          `json:"prdy_ctrt"`  // float64 as string
    AcmlVol       string          `json:"acml_vol"`   // int64 as string
}

type ExpPriceTrend struct {
    Output1 ExpPriceTrendSummary `json:"output1"`
    Output2 []ExpPriceTrendItem  `json:"output2"`
}

type InquireExpPriceTrendParams struct {
    MarketCode     string // fid_cond_mrkt_div_code (소문자), default "J"
    CondScrDivCode string // fid_cond_scr_div_code (소문자), default "11810"
    Symbol         string // fid_input_iscd (소문자)
}
```

### 스텝

1. `testdata/domestic_exp_price_trend.json` fixture 작성 (output1 + output2 배열)
2. `domestic/extended.go` 에 타입 + 메서드 append (소문자 쿼리 키 주의)
3. `domestic/extended_test.go` 에 TestInquireExpPriceTrend 추가
4. `go build ./... && go vet ./...` — silent
5. `go test ./domestic/... -run TestInquireExpPriceTrend -v` — PASS
6. `gofmt -w domestic/extended.go domestic/extended_test.go`

### commit

```bash
git commit -m "$(cat <<'EOF'
[feat] domestic — InquireExpPriceTrend (예상체결가 추이, FHPST01810000)

lowercase fid_* 파라미터 와이어 키 (EP4/5와 동일 패턴).
dual output (output1 ExpPriceTrendSummary + output2 []ExpPriceTrendItem).
antc_cnpr (예상체결가) 포함.

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 11: InquireExpTransUpdown (EP10)

**파일**: `domestic/extended.go` (append)
**엔드포인트**: `GET /uapi/domestic-stock/v1/ranking/exp-trans-updown`
**TR_ID**: `FHPST01820000`
**특이사항**: 쿼리 파라미터 와이어 키 **전부 소문자 fid_***, ranking 엔드포인트
**출력**: `output []ExpTransUpdownItem` (15 fields per item)

### 파라미터 (10개, 모두 소문자 와이어 키)

| Go 필드 | 와이어 키 (소문자) | 기본값 | 설명 |
|---------|-------------------|--------|------|
| MarketCode | `fid_cond_mrkt_div_code` | `"J"` | 시장구분 |
| CondScrDivCode | `fid_cond_scr_div_code` | `"11820"` | 조건화면분류코드 |
| Symbol | `fid_input_iscd` | — | 종목코드 |
| DivClsCode | `fid_div_cls_code` | — | 분류구분코드 |
| RankSortCode | `fid_rank_sort_cls_code` | — | 순위정렬구분코드 |
| InputPrice1 | `fid_input_price_1` | — | 입력가격1 |
| InputPrice2 | `fid_input_price_2` | — | 입력가격2 |
| VolCnt | `fid_vol_cnt` | — | 거래량수 |
| TrgtClsCode | `fid_trgt_cls_code` | — | 대상구분코드 |
| TrgtExlsCode | `fid_trgt_exls_cls_code` | — | 대상제외구분코드 |

### 타입 정의

```go
type ExpTransUpdownItem struct {
    StckShrnIscd   string          `json:"stck_shrn_iscd"`
    HtsKorIsnm     string          `json:"hts_kor_isnm"`
    StckPrpr       decimal.Decimal `json:"stck_prpr"`
    PrdyVrss       decimal.Decimal `json:"prdy_vrss"`
    PrdyVrssSign   string          `json:"prdy_vrss_sign"`
    PrdyCtrt       string          `json:"prdy_ctrt"`       // float64 as string
    StckSdpr       decimal.Decimal `json:"stck_sdpr"`
    SelnRsqn       string          `json:"seln_rsqn"`       // int64 as string
    Askp           decimal.Decimal `json:"askp"`
    Bidp           decimal.Decimal `json:"bidp"`
    ShnuRsqn       string          `json:"shnu_rsqn"`       // int64 as string
    CntgVol        string          `json:"cntg_vol"`        // int64 as string
    AntcTrPbmn     string          `json:"antc_tr_pbmn"`    // int64 as string
    TotalAskpRsqn  string          `json:"total_askp_rsqn"` // int64 as string
    TotalBidpRsqn  string          `json:"total_bidp_rsqn"` // int64 as string
}

type ExpTransUpdown struct {
    Output []ExpTransUpdownItem `json:"output"`
}

type InquireExpTransUpdownParams struct {
    MarketCode     string // fid_cond_mrkt_div_code, default "J"
    CondScrDivCode string // fid_cond_scr_div_code, default "11820"
    Symbol         string // fid_input_iscd
    DivClsCode     string // fid_div_cls_code
    RankSortCode   string // fid_rank_sort_cls_code
    InputPrice1    string // fid_input_price_1
    InputPrice2    string // fid_input_price_2
    VolCnt         string // fid_vol_cnt
    TrgtClsCode    string // fid_trgt_cls_code
    TrgtExlsCode   string // fid_trgt_exls_cls_code
}
```

### 스텝

1. `testdata/domestic_exp_trans_updown.json` fixture 작성 (output 배열, 15 fields/item)
2. `domestic/extended.go` 에 타입 + 메서드 append (소문자 쿼리 키 주의)
3. `domestic/extended_test.go` 에 TestInquireExpTransUpdown 추가 (10 파라미터 + 15 fields 검증)
4. `go build ./... && go vet ./...` — silent
5. `go test ./domestic/... -run TestInquireExpTransUpdown -v` — PASS
6. `gofmt -w domestic/extended.go domestic/extended_test.go`

### commit

```bash
git commit -m "$(cat <<'EOF'
[feat] domestic — InquireExpTransUpdown (예상체결 상승/하락 상위, FHPST01820000)

lowercase fid_* 파라미터 10개 (ranking endpoint 패턴).
output []ExpTransUpdownItem 15 fields.
antc_tr_pbmn (예상거래대금) + total_askp/bidp_rsqn (총호가잔량) 포함.

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 12: examples/domestic_stock_info/main.go

Phase 4.1 전체 10 메서드를 시연하는 예제 프로그램.

**파일**: `examples/domestic_stock_info/main.go` (신규 생성)

### 구조

```go
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

    // 1. InquireInvestOpinion — 투자의견
    opinion, err := client.Domestic.InquireInvestOpinion(ctx, domestic.InquireInvestOpinionParams{
        Symbol: symbol,
    })
    if err != nil {
        log.Printf("InquireInvestOpinion: %v", err)
    } else {
        fmt.Printf("[EP1] InquireInvestOpinion: output 건수=%d\n", len(opinion.Output))
    }

    // 2. InquireInvestOpbysec — 증권사별 투자의견 (증권사코드 예: "62")
    opbysec, err := client.Domestic.InquireInvestOpbysec(ctx, domestic.InquireInvestOpbysecParams{
        Symbol:   "62",
        InputDate: "20260101",
    })
    if err != nil {
        log.Printf("InquireInvestOpbysec: %v", err)
    } else {
        fmt.Printf("[EP2] InquireInvestOpbysec: output 건수=%d\n", len(opbysec.Output))
    }

    // 3. InquireEstimatePerform — 추정실적
    estPerf, err := client.Domestic.InquireEstimatePerform(ctx, domestic.InquireEstimatePerformParams{
        Symbol: symbol,
    })
    if err != nil {
        log.Printf("InquireEstimatePerform: %v", err)
    } else {
        fmt.Printf("[EP3] InquireEstimatePerform: output1 건수=%d\n", len(estPerf.Output1))
    }

    // 4. InquireVolumePower — 체결강도
    volPow, err := client.Domestic.InquireVolumePower(ctx, domestic.InquireVolumePowerParams{
        Symbol: symbol,
    })
    if err != nil {
        log.Printf("InquireVolumePower: %v", err)
    } else {
        fmt.Printf("[EP4] InquireVolumePower: output 건수=%d\n", len(volPow.Output))
    }

    // 5. InquireBulkTransNum — 대량체결건수
    bulkTrans, err := client.Domestic.InquireBulkTransNum(ctx, domestic.InquireBulkTransNumParams{
        Symbol: symbol,
    })
    if err != nil {
        log.Printf("InquireBulkTransNum: %v", err)
    } else {
        fmt.Printf("[EP5] InquireBulkTransNum: output 건수=%d\n", len(bulkTrans.Output))
    }

    // 6. InquireTradprtByamt — 체결금액별 거래비중
    tradprt, err := client.Domestic.InquireTradprtByamt(ctx, domestic.InquireTradprtByamtParams{
        Symbol: symbol,
    })
    if err != nil {
        log.Printf("InquireTradprtByamt: %v", err)
    } else {
        fmt.Printf("[EP6] InquireTradprtByamt: output 건수=%d\n", len(tradprt.Output))
    }

    // 7. InquireHtsTopView — HTS 조회 상위 20 종목 (zero params)
    htsTop, err := client.Domestic.InquireHtsTopView(ctx, domestic.InquireHtsTopViewParams{})
    if err != nil {
        log.Printf("InquireHtsTopView: %v", err)
    } else {
        fmt.Printf("[EP7] InquireHtsTopView: mrkt_div_cls_code=%s mksc_shrn_iscd=%s\n",
            htsTop.Output1.MrktDivClsCode, htsTop.Output1.MkscShrnIscd)
    }

    // 8. InquirePbarTraRatio — 체결금액별 매매비중
    pbarRatio, err := client.Domestic.InquirePbarTraRatio(ctx, domestic.InquirePbarTraRatioParams{
        Symbol:     symbol,
        InputHour1: "090000",
    })
    if err != nil {
        log.Printf("InquirePbarTraRatio: %v", err)
    } else {
        fmt.Printf("[EP8] InquirePbarTraRatio: output2 건수=%d, wghn_avrg_stck_prc=%s\n",
            len(pbarRatio.Output2), pbarRatio.Output1.WghnAvrgStckPrc)
    }

    // 9. InquireExpPriceTrend — 예상체결가 추이
    expTrend, err := client.Domestic.InquireExpPriceTrend(ctx, domestic.InquireExpPriceTrendParams{
        Symbol: symbol,
    })
    if err != nil {
        log.Printf("InquireExpPriceTrend: %v", err)
    } else {
        fmt.Printf("[EP9] InquireExpPriceTrend: antc_cnpr=%s output2 건수=%d\n",
            expTrend.Output1.AntcCnpr, len(expTrend.Output2))
    }

    // 10. InquireExpTransUpdown — 예상체결 상승/하락 상위
    expUpdown, err := client.Domestic.InquireExpTransUpdown(ctx, domestic.InquireExpTransUpdownParams{
        Symbol: symbol,
    })
    if err != nil {
        log.Printf("InquireExpTransUpdown: %v", err)
    } else {
        fmt.Printf("[EP10] InquireExpTransUpdown: output 건수=%d\n", len(expUpdown.Output))
    }
}
```

### 스텝

1. `examples/domestic_stock_info/` 디렉터리 생성 후 `main.go` 작성
2. `go build ./examples/domestic_stock_info && echo OK` — "OK" 출력 확인
3. `rm -f domestic_stock_info` — 유출 바이너리 정리
4. `gofmt -w examples/domestic_stock_info/main.go`
5. `go vet ./examples/domestic_stock_info/...` — silent

### commit

```bash
git commit -m "$(cat <<'EOF'
[feat] examples — domestic_stock_info (Phase 4.1 10 메서드 시연)

EP1~EP10 전체 호출 예제.
005930 (삼성전자) 기준, 증권사코드 "62" (InquireInvestOpbysec EP2).

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 13: 문서 갱신

4개 파일 업데이트.

### 13-1. CLAUDE.md

**Banner 교체**:
```
> **Phase 3.1 — 장내채권 시세 8 메서드 (v1.11.0). Phase 3 design spec 및 plan 참고.**
```
→
```
> **Phase 4.1 — 국내주식 종목정보/분석 10 메서드 (v1.12.0). Phase 4 design spec 및 plan 참고.**
```

**링크 추가** (Phase 3.1 plan 아래):
```markdown
- Phase 4 design spec: [`docs/superpowers/specs/2026-05-07-phase4-stock-info-design.md`](docs/superpowers/specs/2026-05-07-phase4-stock-info-design.md)
- Phase 4.1 implementation plan: [`docs/superpowers/specs/2026-05-07-phase4-1-stock-info-implementation-plan.md`](docs/superpowers/specs/2026-05-07-phase4-1-stock-info-implementation-plan.md)
```

### 13-2. README.md

- 섹션 제목 변경: `Available Methods (Phase 1.2 ~ 3.1)` → `Available Methods (Phase 1.2 ~ 4.1)`
- 메서드 표에 10행 추가 (EP1~EP10):

| 메서드 | TR_ID | 설명 |
|--------|-------|------|
| `InquireInvestOpinion` | FHKST663300C0 | 종목 투자의견 |
| `InquireInvestOpbysec` | FHKST663400C0 | 증권사별 투자의견 |
| `InquireEstimatePerform` | HHKST668300C0 | 추정실적 (quad-output) |
| `InquireVolumePower` | FHPST01680000 | 체결강도 |
| `InquireBulkTransNum` | FHKST190900C0 | 대량체결건수 |
| `InquireTradprtByamt` | FHKST111900C0 | 체결금액별 거래비중 |
| `InquireHtsTopView` | HHMCM000100C0 | HTS 조회 상위 20 종목 |
| `InquirePbarTraRatio` | FHPST01130000 | 체결금액별 매매비중 |
| `InquireExpPriceTrend` | FHPST01810000 | 예상체결가 추이 |
| `InquireExpTransUpdown` | FHPST01820000 | 예상체결 상승/하락 상위 |

- 메서드 총수 업데이트: `79 → 89`

### 13-3. CHANGELOG.md

`## [1.11.0]` 위에 아래 블록 추가:

```markdown
## [1.12.0] - 2026-05-07

### Added
- `InquireInvestOpinion` (FHKST663300C0) — 종목 투자의견
- `InquireInvestOpbysec` (FHKST663400C0) — 증권사별 투자의견
- `InquireEstimatePerform` (HHKST668300C0) — 추정실적 (quad-output: output1+output2+output3+output4)
- `InquireVolumePower` (FHPST01680000) — 체결강도
- `InquireBulkTransNum` (FHKST190900C0) — 대량체결건수
- `InquireTradprtByamt` (FHKST111900C0) — 체결금액별 거래비중
- `InquireHtsTopView` (HHMCM000100C0) — HTS 조회 상위 20 종목
- `InquirePbarTraRatio` (FHPST01130000) — 체결금액별 매매비중
- `InquireExpPriceTrend` (FHPST01810000) — 예상체결가 추이
- `InquireExpTransUpdown` (FHPST01820000) — 예상체결 상승/하락 상위
- `examples/domestic_stock_info/main.go` — Phase 4.1 10 메서드 시연 예제

### Notes
- EP3 (InquireEstimatePerform): quad-output 구조, non-FID SHT_CD 파라미터, HH prefix TR_ID
- EP7 (InquireHtsTopView): zero query params, output1 응답 key anomaly
- EP4/5/9/10: lowercase fid_* 와이어 키 (EP3 SHT_CD 제외)
- EP5 (InquireBulkTransNum): mksc_shrn_iscd (not stck_shrn_iscd) 종목코드 필드
- EP6 (InquireTradprtByamt): `whol_shun_vol_rate` wire typo 의도적 보존
- EP2 (InquireInvestOpbysec): FID_INPUT_ISCD = 증권사코드 (종목코드 아님)
```

### 13-4. domestic/doc.go

Phase 4.1 섹션 추가:

```go
// Phase 4.1 — 국내주식 종목정보/분석 (v1.12.0)
//
// InquireInvestOpinion    — 종목 투자의견 (FHKST663300C0)
// InquireInvestOpbysec    — 증권사별 투자의견 (FHKST663400C0)
// InquireEstimatePerform  — 추정실적, quad-output (HHKST668300C0)
// InquireVolumePower      — 체결강도 (FHPST01680000)
// InquireBulkTransNum     — 대량체결건수 (FHKST190900C0)
// InquireTradprtByamt     — 체결금액별 거래비중 (FHKST111900C0)
// InquireHtsTopView       — HTS 조회 상위 20 종목, zero params (HHMCM000100C0)
// InquirePbarTraRatio     — 체결금액별 매매비중, dual output (FHPST01130000)
// InquireExpPriceTrend    — 예상체결가 추이, dual output (FHPST01810000)
// InquireExpTransUpdown   — 예상체결 상승/하락 상위, lowercase fid_* (FHPST01820000)
```

### 스텝

1. 기존 파일 Read
2. 4개 파일 Edit 적용
3. 검증: `cd /Users/user/src/workspace_moneyflow/korea-investment-stock && go build ./... && go vet ./... && gofmt -l .` — silent
4. **CRITICAL**: `go build` 로 루트에 바이너리 유출 시 삭제 후 커밋
5. 커밋

### commit

```bash
git commit -m "$(cat <<'EOF'
[docs] Phase 4.1 문서 갱신 (v1.12.0)

CLAUDE.md banner + 링크, README 메서드 표 79→89, CHANGELOG v1.12.0,
domestic/doc.go Phase 4.1 섹션.

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 14: 최종 점검

**Note**: 초기 개요는 13 태스크로 기재되었으나, Phase 2.x 패턴과 동일하게 최종 점검(Task 14) + PR 생성(Task 15)을 별도 태스크로 분리.

### 스텝

1. `gofmt -l . | head` — 출력 없음 (clean)
2. `go build ./... && go vet ./...` — silent
3. `go test ./... -race -count=1` — 전 패키지 PASS
4. 커버리지 확인:
   ```bash
   go test ./domestic/... -coverprofile=cov.out && go tool cover -func=cov.out | grep total
   go test ./... -coverprofile=cov_root.out && go tool cover -func=cov_root.out | grep total
   ```
   - `domestic` 패키지: ≥ 80%
   - root `kis` 패키지: ≥ 80%
5. 파일 수 확인:
   - `domestic/opinion.go` (신규)
   - `domestic/opinion_test.go` (신규)
   - `domestic/extended.go` (EP4~EP10 append)
   - `domestic/extended_test.go` (EP4~EP10 테스트 추가)
   - `testdata/domestic_invest_opinion.json` … `domestic_exp_trans_updown.json` (10개)
   - `examples/domestic_stock_info/main.go` (신규)
   - 총 **13 파일** 변경/신규
6. 커밋 수 확인: `git log --oneline feat/phase4-1-stock-info ^main | wc -l` → 약 14~15 커밋

### 커버리지 < 80% 시 대응

Phase 2.6/3.1 교훈 적용: 각 신규 메서드에 InvalidJSON 실패 테스트 추가.

```go
func TestInquireXxx_InvalidJSON(t *testing.T) {
    srv := httpmock.New(t, "TRXXX")
    srv.Respond([]byte(`{invalid`))
    defer srv.Close()
    client := newTestClient(t, srv.URL())
    _, err := client.Domestic.InquireXxx(context.Background(), domestic.InquireXxxParams{...})
    require.Error(t, err)
    assert.Contains(t, err.Error(), "parse")
}
```

커버리지 목표 달성 후 재검증.

---

## Task 15: PR 생성 (사용자 승인 후)

> Claude 는 push / PR 생성을 **사용자 명시적 승인 후에만** 실행 (글로벌 정책).

### 스텝

1. **사용자 승인 요청** — 전체 작업 완료 보고 후 PR 생성 가능 여부 confirm 요청.

2. **Push**:
   ```bash
   git push -u origin feat/phase4-1-stock-info
   ```

3. **PR 생성**:
   ```bash
   gh pr create \
     --title "Phase 4.1 — 국내주식 종목정보/분석 (v1.12.0)" \
     --reviewer kenshin579 \
     --base main \
     --head feat/phase4-1-stock-info \
     --body "$(cat <<'EOF'
   ## Summary

   - Phase 4.1 구현 (신규 10 메서드) — 국내주식 종목정보/분석
   - 신규 파일 `domestic/opinion.go` (3 투자의견) + `domestic/extended.go` 에 7 append
   - Phase 2 표준 패턴 그대로 적용
   - v1.12.0 release (누적 79 → 89 메서드)

   ## 메서드 → 한투 API 매핑 (10 NEW)

   | Go 메서드 | path | TR_ID |
   |---|---|---|
   | InquireInvestOpinion | quotations/invest-opinion | FHKST663300C0 |
   | InquireInvestOpbysec | quotations/invest-opbysec | FHKST663400C0 |
   | InquireEstimatePerform | quotations/estimate-perform | HHKST668300C0 |
   | InquireVolumePower | ranking/volume-power | FHPST01680000 |
   | InquireBulkTransNum | ranking/bulk-trans-num | FHKST190900C0 |
   | InquireTradprtByamt | quotations/tradprt-byamt | FHKST111900C0 |
   | InquireHtsTopView | ranking/hts-top-view | HHMCM000100C0 |
   | InquirePbarTraRatio | quotations/pbar-tratio | FHPST01130000 |
   | InquireExpPriceTrend | quotations/exp-price-trend | FHPST01810000 |
   | InquireExpTransUpdown | ranking/exp-trans-updown | FHPST01820000 |

   ## Anomalies handled

   - EP3 (InquireEstimatePerform) **quad-output** (output1+output2+output3+output4)
   - EP3 문서 한글 레이블 오류 — Python dataclass 필드명을 source of truth로 사용
   - EP3 non-FID SHT_CD 파라미터 + HH prefix TR_ID
   - EP4/5/9/10 **lowercase fid_*** 와이어 키
   - EP5 `mksc_shrn_iscd` (not `stck_shrn_iscd`) 종목코드 필드
   - EP6 `whol_shun_vol_rate` wire typo 의도적 보존
   - EP7 zero query params + `output1` (not `output`) 응답 키
   - EP2 `FID_INPUT_ISCD` = 증권사코드 (종목코드 아님)

   ## Test Plan
   - [x] go build/vet/fmt clean
   - [x] go test ./... -race -count=1 전 패키지 PASS
   - [x] Coverage domestic ≥ 80%, root kis ≥ 80%
   - [x] httpmock 단위 테스트 (10 메서드)
   - [x] examples/domestic_stock_info build OK

   ## Breaking Changes
   없음.

   ## 참고 문서
   - Phase 4 design spec: `docs/superpowers/specs/2026-05-07-phase4-stock-info-design.md`
   - Phase 4.1 plan: `docs/superpowers/specs/2026-05-07-phase4-1-stock-info-implementation-plan.md`

   🤖 Generated with [Claude Code](https://claude.com/claude-code)
   EOF
   )"
   ```

4. **Merge** (사용자 승인 후):
   ```bash
   gh pr merge <PR#> --merge
   ```

5. **후속 작업** (사용자 승인 후):
   ```bash
   git tag -a v1.12.0 -m "Phase 4.1 — 국내주식 종목정보/분석 10 메서드"
   git push origin v1.12.0
   gh release create v1.12.0 --title "v1.12.0 — Phase 4.1 국내주식 종목정보/분석" --notes-from-tag
   ```
