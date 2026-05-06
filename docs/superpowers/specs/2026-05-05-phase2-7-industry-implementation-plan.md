# Phase 2.7 — 국내업종 확장 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** korea-investment-stock Go 라이브러리에 국내업종 7 메서드 추가 (`v1.10.0` release). EP1(`InquireIndexPrice`) + EP2(`InquireIndexCategoryPrice`)는 Phase 1.4에서 이미 구현됨 — Phase 2.7은 EP3~EP9 신규 7 메서드.

**Architecture:** Phase 1 인프라 + 패턴 재사용. `domestic/industry.go` 에 APPEND (기존 파일 수정). 새 internal package 불필요. TDD: testdata fixture → 실패 테스트 → struct + 메서드 구현 → 통과 → commit.

**Tech Stack:** Go 1.25+, `github.com/jarcoal/httpmock`, `github.com/stretchr/testify`, `github.com/shopspring/decimal`. 새 dependency 없음.

**참고 spec:**
- Phase 2.5+ design spec: `docs/superpowers/specs/2026-05-05-phase2-5plus-extension-design.md` (§Phase 2.7)
- Phase 2.6 plan (참조 패턴 — compact task structure): `docs/superpowers/specs/2026-05-05-phase2-6-overseas-info-implementation-plan.md`
- Phase 1.4 구현 참조 (same file, EP1/EP2 already there): `domestic/industry.go`

---

## 사전 정보

| 항목 | 값 |
|------|---|
| 작업 브랜치 | `feat/phase2-7-industry` |
| 시작 HEAD | Phase 2.6 구현 완료 commit (v1.9.0) |
| Release 목표 | `v1.10.0` |
| PR 베이스 | `main` |
| 현재 main 상태 | v1.9.0 publish 완료 (Phase 2.6 통합, 64 메서드) |

> **CRITICAL NOTE:** design spec은 9 메서드를 기술하나, EP1(`InquireIndexPrice`, FHPUP02100000) + EP2(`InquireIndexCategoryPrice`, FHPUP02140000)는 **Phase 1.4에서 이미 구현** (`domestic/industry.go` 상단 참고). Phase 2.7 = 7 NEW 메서드(EP3~EP9). 누적: 64 → **71** 메서드.

---

## 메서드 매핑

| # | Go 메서드 | path (last segment) | TR_ID | output key | fields | anomalies |
|---|---|---|---|---|---|---|
| EP3 | `InquireIndexDailyPrice` | `inquire-index-daily-price` | FHPUP02120000 | `output1{}+output2[]` | 20+13 | none |
| EP4 | `InquireIndexTimeprice` | `inquire-index-timeprice` | FHPUP02110200 | `output []` | 8 | `bsop_hour` timestamp param |
| EP5 | `InquireIndexTickprice` | `inquire-index-tickprice` | FHPUP02110100 | `output []` | 8 | `stck_cntg_hour` timestamp field |
| EP6 | `InquireDailyIndexchartprice` | `inquire-daily-indexchartprice` | FHKUP03500100 | `output1{}+output2[]` | 15+8 | `futs_prdy_*` embedded in output1 |
| EP7 | `InquireTimeIndexchartprice` | `inquire-time-indexchartprice` | FHKUP03500200 | `output1{}+output2[]` | 16+8 | EP6 유사, intraday |
| EP8 | `ExpTotalIndex` | `exp-total-index` | FHKUP11750000 | `output1{}+output2[]` | 9+10 | **LOWERCASE query params** (`fid_*`), `prdy_ctrt` (short form), `fid_cond_scr_div_code="11175"` hardcoded |
| EP9 | `ExpIndexTrend` | `exp-index-trend` | FHPST01840000 | `output []` | 7 | doc Korean labels scrambled (field names correct), `prdy_ctrt` short form |

Default `FID_COND_MRKT_DIV_CODE` = `"U"` for all index endpoints.

---

## 파일 구조

### 수정 (APPEND)
- `domestic/industry.go` — EP3~EP9 struct + Params + 7 메서드 추가
- `domestic/industry_test.go` — 7 테스트 함수 추가
- `CLAUDE.md` — banner Phase 2.6 → Phase 2.7, plan link 추가
- `README.md` — Available Methods 표 갱신 (64 → 71 메서드)
- `CHANGELOG.md` — `[1.10.0]` entry ABOVE `[1.9.0]`
- `domestic/doc.go` — Phase 2.7 section 추가

### 신규 (testdata)
- `domestic/testdata/index_daily_price_success.json`
- `domestic/testdata/index_timeprice_success.json`
- `domestic/testdata/index_tickprice_success.json`
- `domestic/testdata/daily_indexchartprice_success.json`
- `domestic/testdata/time_indexchartprice_success.json`
- `domestic/testdata/exp_total_index_success.json`
- `domestic/testdata/exp_index_trend_success.json`

### 신규 (examples)
- `examples/industry_extended/main.go`

---

## 타입 매핑

Phase 2 표준 타입 매핑 — 업종 지수 특화 주석 포함.

| 카테고리 | Go 타입 | json tag suffix | 예시 필드 |
|---|---|---|---|
| 지수/가격 | `decimal.Decimal` | (bare) | `bstp_nmix_prpr`, `bstp_nmix_oprc/hgpr/lwpr`, `prdy_nmix`, `bstp_nmix_prdy_vrss`, `nmix_sdpr`, `futs_prdy_*`, `dryy_bstp_nmix_*`, `invt_new_psdg`, `d20_dsrt` |
| 거래량/거래대금 | `int64` | `,string` | `acml_vol`, `acml_tr_pbmn`, `cntg_vol`, `prdy_vol`, `prdy_tr_pbmn` |
| 비율 | `float64` | `,string` | `bstp_nmix_prdy_ctrt`, `prdy_ctrt`, `acml_vol_rlim` |
| 코드/날짜/Y-N/수 | `string` | (bare) | `prdy_vrss_sign`, `bsop_hour`, `stck_cntg_hour`, `stck_bsop_date`, `bstp_cls_code`, `hts_kor_isnm`, `mod_yn`, `*_issu_cnt`, `*_date` |

> **EP8 anomaly:** `prdy_ctrt` (short form) — `bstp_nmix_prdy_ctrt` 아님. `float64,string` 타입 동일.
> **EP9 anomaly:** `prdy_ctrt` 동일 short form.

---

## Tasks (12 total)

| # | 내용 | Files |
|---|---|---|
| Task 1 | testdata fixtures (7 합성 JSON) | `domestic/testdata/*.json` |
| Task 2 | EP3 `InquireIndexDailyPrice` | APPEND `industry.go` / `industry_test.go` |
| Task 3 | EP4 `InquireIndexTimeprice` | APPEND `industry.go` / `industry_test.go` |
| Task 4 | EP5 `InquireIndexTickprice` | APPEND `industry.go` / `industry_test.go` |
| Task 5 | EP6 `InquireDailyIndexchartprice` | APPEND `industry.go` / `industry_test.go` |
| Task 6 | EP7 `InquireTimeIndexchartprice` | APPEND `industry.go` / `industry_test.go` |
| Task 7 | EP8 `ExpTotalIndex` (lowercase params anomaly) | APPEND `industry.go` / `industry_test.go` |
| Task 8 | EP9 `ExpIndexTrend` (prdy_ctrt short form) | APPEND `industry.go` / `industry_test.go` |
| Task 9 | examples `examples/industry_extended/main.go` | CREATE |
| Task 10 | `domestic/doc.go` Phase 2.7 section 추가 | MODIFY |
| Task 11 | `CHANGELOG.md` `[1.10.0]` entry | MODIFY |
| Task 12 | `CLAUDE.md` + `README.md` 갱신 + `go test ./...` 전체 + tag | MODIFY |

---

## Task 1: testdata fixtures (7 합성 JSON)

- [ ] Step 1: `domestic/testdata/index_daily_price_success.json` (EP3 — `output1{}+output2[]`, 20+13 fields)
- [ ] Step 2: `domestic/testdata/index_timeprice_success.json` (EP4 — `output []`, 8 fields)
- [ ] Step 3: `domestic/testdata/index_tickprice_success.json` (EP5 — `output []`, 8 fields)
- [ ] Step 4: `domestic/testdata/daily_indexchartprice_success.json` (EP6 — `output1{}+output2[]`, 15+8 fields)
- [ ] Step 5: `domestic/testdata/time_indexchartprice_success.json` (EP7 — `output1{}+output2[]`, 16+8 fields)
- [ ] Step 6: `domestic/testdata/exp_total_index_success.json` (EP8 — `output1{}+output2[]`, 9+10 fields)
- [ ] Step 7: `domestic/testdata/exp_index_trend_success.json` (EP9 — `output []`, 7 fields)
- [ ] Step 8: validation

```bash
for f in \
  domestic/testdata/index_daily_price_success.json \
  domestic/testdata/index_timeprice_success.json \
  domestic/testdata/index_tickprice_success.json \
  domestic/testdata/daily_indexchartprice_success.json \
  domestic/testdata/time_indexchartprice_success.json \
  domestic/testdata/exp_total_index_success.json \
  domestic/testdata/exp_index_trend_success.json; do
  python3 -c "import json; json.load(open('$f'))" && echo "$f OK"
done
# Expected: 7 OK lines
```

- [ ] Step 9: commit

```bash
git commit -m "$(cat <<'EOF'
[chore] testdata — 7 industry extended fixture JSON (Phase 2.7)

합성 JSON fixtures (2 records each where applicable):
- index_daily_price_success.json (output1 20 fields + output2[] 13 fields)
- index_timeprice_success.json (output[] 8 fields, bsop_hour timestamp)
- index_tickprice_success.json (output[] 8 fields, stck_cntg_hour timestamp)
- daily_indexchartprice_success.json (output1 15 fields + output2[] 8 fields)
- time_indexchartprice_success.json (output1 16 fields + output2[] 8 fields)
- exp_total_index_success.json (output1 9 fields + output2[] 10 fields, lowercase params)
- exp_index_trend_success.json (output[] 7 fields, prdy_ctrt short form)

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

### Fixture content

**Step 1 — `index_daily_price_success.json`** (EP3: output1 20 fields + output2[] 13 fields)

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
      "stck_bsop_date": "20260505",
      "bstp_nmix_prpr": "2650.45",
      "prdy_vrss_sign": "5",
      "bstp_nmix_prdy_vrss": "-12.30",
      "bstp_nmix_prdy_ctrt": "-0.46",
      "bstp_nmix_oprc": "2660.00",
      "bstp_nmix_hgpr": "2665.20",
      "bstp_nmix_lwpr": "2645.10",
      "acml_vol_rlim": "100.00",
      "acml_vol": "350000000",
      "acml_tr_pbmn": "9500000",
      "invt_new_psdg": "0.45",
      "d20_dsrt": "52.30"
    },
    {
      "stck_bsop_date": "20260504",
      "bstp_nmix_prpr": "2662.75",
      "prdy_vrss_sign": "2",
      "bstp_nmix_prdy_vrss": "5.50",
      "bstp_nmix_prdy_ctrt": "0.21",
      "bstp_nmix_oprc": "2658.00",
      "bstp_nmix_hgpr": "2668.40",
      "bstp_nmix_lwpr": "2655.30",
      "acml_vol_rlim": "98.50",
      "acml_vol": "345000000",
      "acml_tr_pbmn": "9200000",
      "invt_new_psdg": "0.38",
      "d20_dsrt": "51.80"
    }
  ]
}
```

**Step 2 — `index_timeprice_success.json`** (EP4: output[] 8 fields, bsop_hour timestamp)

```json
{
  "rt_cd": "0",
  "msg_cd": "MCA00000",
  "msg1": "정상처리 되었습니다.",
  "output": [
    {
      "bsop_hour": "100000",
      "bstp_nmix_prpr": "2652.10",
      "bstp_nmix_prdy_vrss": "-10.65",
      "prdy_vrss_sign": "5",
      "bstp_nmix_prdy_ctrt": "-0.40",
      "acml_tr_pbmn": "3200000",
      "acml_vol": "120000000",
      "cntg_vol": "800000"
    },
    {
      "bsop_hour": "093000",
      "bstp_nmix_prpr": "2658.80",
      "bstp_nmix_prdy_vrss": "-3.95",
      "prdy_vrss_sign": "5",
      "bstp_nmix_prdy_ctrt": "-0.15",
      "acml_tr_pbmn": "1500000",
      "acml_vol": "55000000",
      "cntg_vol": "600000"
    }
  ]
}
```

**Step 3 — `index_tickprice_success.json`** (EP5: output[] 8 fields, stck_cntg_hour timestamp)

```json
{
  "rt_cd": "0",
  "msg_cd": "MCA00000",
  "msg1": "정상처리 되었습니다.",
  "output": [
    {
      "stck_cntg_hour": "100015",
      "bstp_nmix_prpr": "2651.35",
      "bstp_nmix_prdy_vrss": "-11.40",
      "prdy_vrss_sign": "5",
      "bstp_nmix_prdy_ctrt": "-0.43",
      "acml_tr_pbmn": "3350000",
      "acml_vol": "123000000",
      "cntg_vol": "50000"
    },
    {
      "stck_cntg_hour": "100014",
      "bstp_nmix_prpr": "2651.80",
      "bstp_nmix_prdy_vrss": "-10.95",
      "prdy_vrss_sign": "5",
      "bstp_nmix_prdy_ctrt": "-0.41",
      "acml_tr_pbmn": "3300000",
      "acml_vol": "122500000",
      "cntg_vol": "45000"
    }
  ]
}
```

**Step 4 — `daily_indexchartprice_success.json`** (EP6: output1 15 fields + output2[] 8 fields)

```json
{
  "rt_cd": "0",
  "msg_cd": "MCA00000",
  "msg1": "정상처리 되었습니다.",
  "output1": {
    "prdy_vrss_sign": "5",
    "bstp_nmix_prdy_ctrt": "-0.46",
    "prdy_nmix": "2662.75",
    "acml_vol": "350000000",
    "acml_tr_pbmn": "9500000",
    "hts_kor_isnm": "코스피",
    "bstp_nmix_prpr": "2650.45",
    "bstp_cls_code": "0001",
    "prdy_vol": "330000000",
    "bstp_nmix_oprc": "2660.00",
    "bstp_nmix_hgpr": "2665.20",
    "bstp_nmix_lwpr": "2645.10",
    "futs_prdy_oprc": "355.50",
    "futs_prdy_hgpr": "358.00",
    "futs_prdy_lwpr": "353.20"
  },
  "output2": [
    {
      "stck_bsop_date": "20260505",
      "bstp_nmix_prpr": "2650.45",
      "bstp_nmix_oprc": "2660.00",
      "bstp_nmix_hgpr": "2665.20",
      "bstp_nmix_lwpr": "2645.10",
      "acml_vol": "350000000",
      "acml_tr_pbmn": "9500000",
      "mod_yn": "N"
    },
    {
      "stck_bsop_date": "20260504",
      "bstp_nmix_prpr": "2662.75",
      "bstp_nmix_oprc": "2658.00",
      "bstp_nmix_hgpr": "2668.40",
      "bstp_nmix_lwpr": "2655.30",
      "acml_vol": "345000000",
      "acml_tr_pbmn": "9200000",
      "mod_yn": "N"
    }
  ]
}
```

**Step 5 — `time_indexchartprice_success.json`** (EP7: output1 16 fields + output2[] 8 fields)

```json
{
  "rt_cd": "0",
  "msg_cd": "MCA00000",
  "msg1": "정상처리 되었습니다.",
  "output1": {
    "bstp_nmix_prdy_vrss": "-12.30",
    "prdy_vrss_sign": "5",
    "bstp_nmix_prdy_ctrt": "-0.46",
    "prdy_nmix": "2662.75",
    "acml_vol": "350000000",
    "acml_tr_pbmn": "9500000",
    "hts_kor_isnm": "코스피",
    "bstp_nmix_prpr": "2650.45",
    "bstp_cls_code": "0001",
    "prdy_vol": "330000000",
    "bstp_nmix_oprc": "2660.00",
    "bstp_nmix_hgpr": "2665.20",
    "bstp_nmix_lwpr": "2645.10",
    "futs_prdy_oprc": "355.50",
    "futs_prdy_hgpr": "358.00",
    "futs_prdy_lwpr": "353.20"
  },
  "output2": [
    {
      "stck_bsop_date": "20260505",
      "stck_cntg_hour": "100000",
      "bstp_nmix_prpr": "2652.10",
      "bstp_nmix_oprc": "2660.00",
      "bstp_nmix_hgpr": "2665.20",
      "bstp_nmix_lwpr": "2648.50",
      "cntg_vol": "800000",
      "acml_tr_pbmn": "3200000"
    },
    {
      "stck_bsop_date": "20260505",
      "stck_cntg_hour": "093000",
      "bstp_nmix_prpr": "2658.80",
      "bstp_nmix_oprc": "2660.00",
      "bstp_nmix_hgpr": "2661.00",
      "bstp_nmix_lwpr": "2656.20",
      "cntg_vol": "600000",
      "acml_tr_pbmn": "1500000"
    }
  ]
}
```

**Step 6 — `exp_total_index_success.json`** (EP8: output1 9 fields + output2[] 10 fields, prdy_ctrt short form)

```json
{
  "rt_cd": "0",
  "msg_cd": "MCA00000",
  "msg1": "정상처리 되었습니다.",
  "output1": {
    "bstp_nmix_prpr": "2650.45",
    "bstp_nmix_prdy_vrss": "-12.30",
    "prdy_vrss_sign": "5",
    "prdy_ctrt": "-0.46",
    "acml_vol": "350000000",
    "ascn_issu_cnt": "315",
    "down_issu_cnt": "450",
    "stnr_issu_cnt": "120",
    "bstp_cls_code": "0001"
  },
  "output2": [
    {
      "hts_kor_isnm": "코스피",
      "bstp_nmix_prpr": "2650.45",
      "bstp_nmix_prdy_vrss": "-12.30",
      "prdy_vrss_sign": "5",
      "bstp_nmix_prdy_ctrt": "-0.46",
      "acml_vol": "350000000",
      "nmix_sdpr": "2662.75",
      "ascn_issu_cnt": "315",
      "stnr_issu_cnt": "120",
      "down_issu_cnt": "450"
    },
    {
      "hts_kor_isnm": "코스닥",
      "bstp_nmix_prpr": "870.20",
      "bstp_nmix_prdy_vrss": "-3.50",
      "prdy_vrss_sign": "5",
      "bstp_nmix_prdy_ctrt": "-0.40",
      "acml_vol": "680000000",
      "nmix_sdpr": "873.70",
      "ascn_issu_cnt": "580",
      "stnr_issu_cnt": "210",
      "down_issu_cnt": "750"
    }
  ]
}
```

**Step 7 — `exp_index_trend_success.json`** (EP9: output[] 7 fields, prdy_ctrt short form)

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
      "prdy_ctrt": "-0.46",
      "acml_vol": "350000000",
      "acml_tr_pbmn": "9500000"
    },
    {
      "stck_bsop_date": "20260504",
      "bstp_nmix_prpr": "2662.75",
      "bstp_nmix_prdy_vrss": "5.50",
      "prdy_vrss_sign": "2",
      "prdy_ctrt": "0.21",
      "acml_vol": "345000000",
      "acml_tr_pbmn": "9200000"
    }
  ]
}
```

---

## Task 2: InquireIndexDailyPrice (EP3)

**Files:** APPEND to `domestic/industry.go` and `domestic/industry_test.go`

국내업종 일별 지수 조회. `output1{}` (대표 스냅샷 20 fields) + `output2[]` (일별 배열 13 fields).

- [ ] Step 1: APPEND test code to `domestic/industry_test.go`
- [ ] Step 2: Verify FAIL — `go test ./domestic/... -run TestClient_InquireIndexDailyPrice -v` (compile error expected)
- [ ] Step 3: APPEND struct + Params + method to `domestic/industry.go`
- [ ] Step 4: Verify PASS — `go test ./domestic/... -run TestClient_InquireIndexDailyPrice -v`
- [ ] Step 5: `gofmt -w domestic/industry.go domestic/industry_test.go && go vet ./domestic/...`
- [ ] Step 6: commit

### 메서드 매핑
- Path: `/uapi/domestic-stock/v1/quotations/inquire-index-daily-price`
- TR_ID: `FHPUP02120000`
- Params (4): `MarketCode` (`FID_COND_MRKT_DIV_CODE` default `"U"`), `Symbol` (`FID_INPUT_ISCD`), `PeriodDivCode` (`FID_PERIOD_DIV_CODE`), `InputDate1` (`FID_INPUT_DATE_1`)

### output1 struct 필드 (20 fields)

| Go field | json tag | Go type | 설명 |
|---|---|---|---|
| `BstpNmixPrpr` | `bstp_nmix_prpr` | `decimal.Decimal` | 업종 지수 현재가 |
| `BstpNmixPrdyVrss` | `bstp_nmix_prdy_vrss` | `decimal.Decimal` | 업종 지수 전일 대비 |
| `PrdyVrssSign` | `prdy_vrss_sign` | `string` | 전일 대비 부호 |
| `BstpNmixPrdyCtrt` | `bstp_nmix_prdy_ctrt` | `float64,string` | 업종 지수 전일 대비율 |
| `AcmlVol` | `acml_vol` | `int64,string` | 누적 거래량 |
| `AcmlTrPbmn` | `acml_tr_pbmn` | `int64,string` | 누적 거래 대금 |
| `BstpNmixOprc` | `bstp_nmix_oprc` | `decimal.Decimal` | 업종 지수 시가 |
| `BstpNmixHgpr` | `bstp_nmix_hgpr` | `decimal.Decimal` | 업종 지수 최고가 |
| `BstpNmixLwpr` | `bstp_nmix_lwpr` | `decimal.Decimal` | 업종 지수 최저가 |
| `PrdyVol` | `prdy_vol` | `int64,string` | 전일 거래량 |
| `AscnIssuCnt` | `ascn_issu_cnt` | `string` | 상승 종목 수 |
| `DownIssuCnt` | `down_issu_cnt` | `string` | 하락 종목 수 |
| `StnrIssuCnt` | `stnr_issu_cnt` | `string` | 보합 종목 수 |
| `UplmIssuCnt` | `uplm_issu_cnt` | `string` | 상한 종목 수 |
| `LslmIssuCnt` | `lslm_issu_cnt` | `string` | 하한 종목 수 |
| `PrdyTrPbmn` | `prdy_tr_pbmn` | `int64,string` | 전일 거래 대금 |
| `DryyBstpNmixHgprDate` | `dryy_bstp_nmix_hgpr_date` | `string` | 연중 최고가 일자 |
| `DryyBstpNmixHgpr` | `dryy_bstp_nmix_hgpr` | `decimal.Decimal` | 연중 업종 지수 최고가 |
| `DryyBstpNmixLwpr` | `dryy_bstp_nmix_lwpr` | `decimal.Decimal` | 연중 업종 지수 최저가 |
| `DryyBstpNmixLwprDate` | `dryy_bstp_nmix_lwpr_date` | `string` | 연중 최저가 일자 |

### output2 struct 필드 (13 fields)

| Go field | json tag | Go type | 설명 |
|---|---|---|---|
| `StckBsopDate` | `stck_bsop_date` | `string` | 영업 일자 |
| `BstpNmixPrpr` | `bstp_nmix_prpr` | `decimal.Decimal` | 업종 지수 현재가 |
| `PrdyVrssSign` | `prdy_vrss_sign` | `string` | 전일 대비 부호 |
| `BstpNmixPrdyVrss` | `bstp_nmix_prdy_vrss` | `decimal.Decimal` | 업종 지수 전일 대비 |
| `BstpNmixPrdyCtrt` | `bstp_nmix_prdy_ctrt` | `float64,string` | 업종 지수 전일 대비율 |
| `BstpNmixOprc` | `bstp_nmix_oprc` | `decimal.Decimal` | 업종 지수 시가 |
| `BstpNmixHgpr` | `bstp_nmix_hgpr` | `decimal.Decimal` | 업종 지수 최고가 |
| `BstpNmixLwpr` | `bstp_nmix_lwpr` | `decimal.Decimal` | 업종 지수 최저가 |
| `AcmlVolRlim` | `acml_vol_rlim` | `float64,string` | 누적 거래량 비중 |
| `AcmlVol` | `acml_vol` | `int64,string` | 누적 거래량 |
| `AcmlTrPbmn` | `acml_tr_pbmn` | `int64,string` | 누적 거래 대금 |
| `InvtNewPsdg` | `invt_new_psdg` | `decimal.Decimal` | 투자자 순매수 주도 |
| `D20Dsrt` | `d20_dsrt` | `decimal.Decimal` | 20일 이격도 |

### 구현 코드

```go
// IndexDailyPrice 는 국내업종 일별지수 (FHPUP02120000) 응답.
//
// 한투 docs: docs/api/국내주식/국내업종_일별지수.md
// path: /uapi/domestic-stock/v1/quotations/inquire-index-daily-price
type IndexDailyPrice struct {
	Output1 IndexDailyPriceSummary `json:"output1"`
	Output2 []IndexDailyPriceItem  `json:"output2"`
}

// IndexDailyPriceSummary 는 응답의 output1 (대표 스냅샷, 20 fields).
type IndexDailyPriceSummary struct {
	BstpNmixPrpr         decimal.Decimal `json:"bstp_nmix_prpr"`             // 업종 지수 현재가
	BstpNmixPrdyVrss     decimal.Decimal `json:"bstp_nmix_prdy_vrss"`        // 업종 지수 전일 대비
	PrdyVrssSign         string          `json:"prdy_vrss_sign"`             // 전일 대비 부호
	BstpNmixPrdyCtrt     float64         `json:"bstp_nmix_prdy_ctrt,string"` // 업종 지수 전일 대비율
	AcmlVol              int64           `json:"acml_vol,string"`            // 누적 거래량
	AcmlTrPbmn           int64           `json:"acml_tr_pbmn,string"`        // 누적 거래 대금
	BstpNmixOprc         decimal.Decimal `json:"bstp_nmix_oprc"`             // 업종 지수 시가
	BstpNmixHgpr         decimal.Decimal `json:"bstp_nmix_hgpr"`             // 업종 지수 최고가
	BstpNmixLwpr         decimal.Decimal `json:"bstp_nmix_lwpr"`             // 업종 지수 최저가
	PrdyVol              int64           `json:"prdy_vol,string"`            // 전일 거래량
	AscnIssuCnt          string          `json:"ascn_issu_cnt"`              // 상승 종목 수
	DownIssuCnt          string          `json:"down_issu_cnt"`              // 하락 종목 수
	StnrIssuCnt          string          `json:"stnr_issu_cnt"`              // 보합 종목 수
	UplmIssuCnt          string          `json:"uplm_issu_cnt"`              // 상한 종목 수
	LslmIssuCnt          string          `json:"lslm_issu_cnt"`              // 하한 종목 수
	PrdyTrPbmn           int64           `json:"prdy_tr_pbmn,string"`        // 전일 거래 대금
	DryyBstpNmixHgprDate string          `json:"dryy_bstp_nmix_hgpr_date"`   // 연중 최고가 일자
	DryyBstpNmixHgpr     decimal.Decimal `json:"dryy_bstp_nmix_hgpr"`        // 연중 업종 지수 최고가
	DryyBstpNmixLwpr     decimal.Decimal `json:"dryy_bstp_nmix_lwpr"`        // 연중 업종 지수 최저가
	DryyBstpNmixLwprDate string          `json:"dryy_bstp_nmix_lwpr_date"`   // 연중 최저가 일자
}

// IndexDailyPriceItem 은 응답의 output2 한 행 (일별, 13 fields).
type IndexDailyPriceItem struct {
	StckBsopDate     string          `json:"stck_bsop_date"`             // 영업 일자
	BstpNmixPrpr     decimal.Decimal `json:"bstp_nmix_prpr"`             // 업종 지수 현재가
	PrdyVrssSign     string          `json:"prdy_vrss_sign"`             // 전일 대비 부호
	BstpNmixPrdyVrss decimal.Decimal `json:"bstp_nmix_prdy_vrss"`        // 업종 지수 전일 대비
	BstpNmixPrdyCtrt float64         `json:"bstp_nmix_prdy_ctrt,string"` // 업종 지수 전일 대비율
	BstpNmixOprc     decimal.Decimal `json:"bstp_nmix_oprc"`             // 업종 지수 시가
	BstpNmixHgpr     decimal.Decimal `json:"bstp_nmix_hgpr"`             // 업종 지수 최고가
	BstpNmixLwpr     decimal.Decimal `json:"bstp_nmix_lwpr"`             // 업종 지수 최저가
	AcmlVolRlim      float64         `json:"acml_vol_rlim,string"`       // 누적 거래량 비중
	AcmlVol          int64           `json:"acml_vol,string"`            // 누적 거래량
	AcmlTrPbmn       int64           `json:"acml_tr_pbmn,string"`        // 누적 거래 대금
	InvtNewPsdg      decimal.Decimal `json:"invt_new_psdg"`              // 투자자 순매수 주도
	D20Dsrt          decimal.Decimal `json:"d20_dsrt"`                   // 20일 이격도
}

// InquireIndexDailyPriceParams 는 국내업종 일별지수 조회 파라미터.
type InquireIndexDailyPriceParams struct {
	MarketCode    string // FID_COND_MRKT_DIV_CODE — 빈 값=>"U" (업종)
	Symbol        string // FID_INPUT_ISCD — 필수, 업종 코드 (예 "0001":코스피)
	PeriodDivCode string // FID_PERIOD_DIV_CODE — D:일 W:주 M:월 Y:년
	InputDate1    string // FID_INPUT_DATE_1 — 조회 시작일 YYYYMMDD
}

// InquireIndexDailyPrice 는 국내업종 일별지수 호출.
//
// 한투 docs: docs/api/국내주식/국내업종_일별지수.md
// path: /uapi/domestic-stock/v1/quotations/inquire-index-daily-price (FHPUP02120000)
func (c *Client) InquireIndexDailyPrice(ctx context.Context, params InquireIndexDailyPriceParams) (*IndexDailyPrice, error) {
	market := params.MarketCode
	if market == "" {
		market = "U"
	}
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/domestic-stock/v1/quotations/inquire-index-daily-price",
		TrID:   "FHPUP02120000",
		Query: map[string]string{
			"FID_COND_MRKT_DIV_CODE": market,
			"FID_INPUT_ISCD":         params.Symbol,
			"FID_PERIOD_DIV_CODE":    params.PeriodDivCode,
			"FID_INPUT_DATE_1":       params.InputDate1,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}
	var res IndexDailyPrice
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse IndexDailyPrice: %w", err)
	}
	return &res, nil
}
```

### 테스트 코드

```go
func TestClient_InquireIndexDailyPrice(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/quotations/inquire-index-daily-price`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "index_daily_price_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireIndexDailyPrice(context.Background(), domestic.InquireIndexDailyPriceParams{
		Symbol:        "0001",
		PeriodDivCode: "D",
		InputDate1:    "20260505",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, "U", capturedQuery.Get("FID_COND_MRKT_DIV_CODE"))
	assert.Equal(t, "0001", capturedQuery.Get("FID_INPUT_ISCD"))
	assert.Equal(t, "D", capturedQuery.Get("FID_PERIOD_DIV_CODE"))
	assert.Equal(t, "20260505", capturedQuery.Get("FID_INPUT_DATE_1"))

	d, _ := decimal.NewFromString("2650.45")
	assert.True(t, d.Equal(res.Output1.BstpNmixPrpr))
	assert.Equal(t, "315", res.Output1.AscnIssuCnt)
	assert.InDelta(t, -0.46, res.Output1.BstpNmixPrdyCtrt, 0.001)

	require.Len(t, res.Output2, 2)
	assert.Equal(t, "20260505", res.Output2[0].StckBsopDate)
	d2, _ := decimal.NewFromString("2650.45")
	assert.True(t, d2.Equal(res.Output2[0].BstpNmixPrpr))
	assert.InDelta(t, 100.00, res.Output2[0].AcmlVolRlim, 0.01)
}
```

### Commit

```bash
git commit -m "$(cat <<'EOF'
[feat] domestic — InquireIndexDailyPrice (국내업종 일별지수, FHPUP02120000)

- IndexDailyPrice / IndexDailyPriceSummary (20 fields) / IndexDailyPriceItem (13 fields)
- InquireIndexDailyPriceParams (MarketCode/Symbol/PeriodDivCode/InputDate1)
- output1 snapshot + output2[] 일별 배열 (invt_new_psdg, d20_dsrt decimal)
- TestClient_InquireIndexDailyPrice — fixture index_daily_price_success.json

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 3: InquireIndexTimeprice (EP4)

**Files:** APPEND to `domestic/industry.go` and `domestic/industry_test.go`

국내업종 시간별 지수 조회. `output []` 배열 8 fields. `bsop_hour` (HHMMSS) 타임스탬프 필드 포함. `FID_INPUT_HOUR_1` 파라미터로 집계 단위(60/300/600초) 설정.

- [ ] Step 1: APPEND test code to `domestic/industry_test.go`
- [ ] Step 2: Verify FAIL — `go test ./domestic/... -run TestClient_InquireIndexTimeprice -v`
- [ ] Step 3: APPEND struct + Params + method to `domestic/industry.go`
- [ ] Step 4: Verify PASS — `go test ./domestic/... -run TestClient_InquireIndexTimeprice -v`
- [ ] Step 5: `gofmt -w domestic/industry.go domestic/industry_test.go && go vet ./domestic/...`
- [ ] Step 6: commit

### 메서드 매핑
- Path: `/uapi/domestic-stock/v1/quotations/inquire-index-timeprice`
- TR_ID: `FHPUP02110200`
- Params (3): `InputHour1` (`FID_INPUT_HOUR_1` — 60/300/600), `Symbol` (`FID_INPUT_ISCD`), `MarketCode` (`FID_COND_MRKT_DIV_CODE` default `"U"`)

### output struct 필드 (8 fields)

| Go field | json tag | Go type | 설명 |
|---|---|---|---|
| `BsopHour` | `bsop_hour` | `string` | 영업 시간 (HHMMSS) |
| `BstpNmixPrpr` | `bstp_nmix_prpr` | `decimal.Decimal` | 업종 지수 현재가 |
| `BstpNmixPrdyVrss` | `bstp_nmix_prdy_vrss` | `decimal.Decimal` | 업종 지수 전일 대비 |
| `PrdyVrssSign` | `prdy_vrss_sign` | `string` | 전일 대비 부호 |
| `BstpNmixPrdyCtrt` | `bstp_nmix_prdy_ctrt` | `float64,string` | 업종 지수 전일 대비율 |
| `AcmlTrPbmn` | `acml_tr_pbmn` | `int64,string` | 누적 거래 대금 |
| `AcmlVol` | `acml_vol` | `int64,string` | 누적 거래량 |
| `CntgVol` | `cntg_vol` | `int64,string` | 체결 거래량 |

### 구현 코드

```go
// IndexTimeprice 는 국내업종 시간별 지수 (FHPUP02110200) 응답.
//
// 한투 docs: docs/api/국내주식/국내업종_시간별지수.md
// path: /uapi/domestic-stock/v1/quotations/inquire-index-timeprice
//
// BsopHour 는 HHMMSS 형식 타임스탬프. FID_INPUT_HOUR_1 파라미터로 집계 단위 설정 (60/300/600초).
type IndexTimeprice struct {
	Output []IndexTimepriceItem `json:"output"`
}

// IndexTimepriceItem 은 응답의 output 한 행 (시간별, 8 fields).
type IndexTimepriceItem struct {
	BsopHour         string          `json:"bsop_hour"`                  // 영업 시간 HHMMSS
	BstpNmixPrpr     decimal.Decimal `json:"bstp_nmix_prpr"`             // 업종 지수 현재가
	BstpNmixPrdyVrss decimal.Decimal `json:"bstp_nmix_prdy_vrss"`        // 업종 지수 전일 대비
	PrdyVrssSign     string          `json:"prdy_vrss_sign"`             // 전일 대비 부호
	BstpNmixPrdyCtrt float64         `json:"bstp_nmix_prdy_ctrt,string"` // 업종 지수 전일 대비율
	AcmlTrPbmn       int64           `json:"acml_tr_pbmn,string"`        // 누적 거래 대금
	AcmlVol          int64           `json:"acml_vol,string"`            // 누적 거래량
	CntgVol          int64           `json:"cntg_vol,string"`            // 체결 거래량
}

// InquireIndexTimepriceParams 는 국내업종 시간별 지수 조회 파라미터.
type InquireIndexTimepriceParams struct {
	InputHour1 string // FID_INPUT_HOUR_1 — 집계 단위: "60"(1분)/"300"(5분)/"600"(10분)
	Symbol     string // FID_INPUT_ISCD — 필수, 업종 코드
	MarketCode string // FID_COND_MRKT_DIV_CODE — 빈 값=>"U" (업종)
}

// InquireIndexTimeprice 는 국내업종 시간별 지수 호출.
//
// 한투 docs: docs/api/국내주식/국내업종_시간별지수.md
// path: /uapi/domestic-stock/v1/quotations/inquire-index-timeprice (FHPUP02110200)
func (c *Client) InquireIndexTimeprice(ctx context.Context, params InquireIndexTimepriceParams) (*IndexTimeprice, error) {
	market := params.MarketCode
	if market == "" {
		market = "U"
	}
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/domestic-stock/v1/quotations/inquire-index-timeprice",
		TrID:   "FHPUP02110200",
		Query: map[string]string{
			"FID_INPUT_HOUR_1":       params.InputHour1,
			"FID_INPUT_ISCD":         params.Symbol,
			"FID_COND_MRKT_DIV_CODE": market,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}
	var res IndexTimeprice
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse IndexTimeprice: %w", err)
	}
	return &res, nil
}
```

### 테스트 코드

```go
func TestClient_InquireIndexTimeprice(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/quotations/inquire-index-timeprice`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "index_timeprice_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireIndexTimeprice(context.Background(), domestic.InquireIndexTimepriceParams{
		InputHour1: "60",
		Symbol:     "0001",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, "60", capturedQuery.Get("FID_INPUT_HOUR_1"))
	assert.Equal(t, "0001", capturedQuery.Get("FID_INPUT_ISCD"))
	assert.Equal(t, "U", capturedQuery.Get("FID_COND_MRKT_DIV_CODE"))

	require.Len(t, res.Output, 2)
	assert.Equal(t, "100000", res.Output[0].BsopHour)
	d, _ := decimal.NewFromString("2652.10")
	assert.True(t, d.Equal(res.Output[0].BstpNmixPrpr))
	assert.Equal(t, int64(800000), res.Output[0].CntgVol)
}
```

### Commit

```bash
git commit -m "$(cat <<'EOF'
[feat] domestic — InquireIndexTimeprice (국내업종 시간별지수, FHPUP02110200)

- IndexTimeprice / IndexTimepriceItem (8 fields, bsop_hour HHMMSS timestamp)
- InquireIndexTimepriceParams (InputHour1 60/300/600 / Symbol / MarketCode)
- output[] 배열; FID_INPUT_HOUR_1 집계 단위 파라미터
- TestClient_InquireIndexTimeprice — fixture index_timeprice_success.json

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 4: InquireIndexTickprice (EP5)

**Files:** APPEND to `domestic/industry.go` and `domestic/industry_test.go`

국내업종 틱별 지수 조회. `output []` 배열 8 fields. `stck_cntg_hour` (HHMMSS) 틱 타임스탬프 필드 포함. 파라미터 2개만 — Symbol + MarketCode.

- [ ] Step 1: APPEND test code to `domestic/industry_test.go`
- [ ] Step 2: Verify FAIL — `go test ./domestic/... -run TestClient_InquireIndexTickprice -v`
- [ ] Step 3: APPEND struct + Params + method to `domestic/industry.go`
- [ ] Step 4: Verify PASS — `go test ./domestic/... -run TestClient_InquireIndexTickprice -v`
- [ ] Step 5: `gofmt -w domestic/industry.go domestic/industry_test.go && go vet ./domestic/...`
- [ ] Step 6: commit

### 메서드 매핑
- Path: `/uapi/domestic-stock/v1/quotations/inquire-index-tickprice`
- TR_ID: `FHPUP02110100`
- Params (2): `Symbol` (`FID_INPUT_ISCD`), `MarketCode` (`FID_COND_MRKT_DIV_CODE` default `"U"`)

### output struct 필드 (8 fields)

| Go field | json tag | Go type | 설명 |
|---|---|---|---|
| `StckCntgHour` | `stck_cntg_hour` | `string` | 주식 체결 시간 (HHMMSS) |
| `BstpNmixPrpr` | `bstp_nmix_prpr` | `decimal.Decimal` | 업종 지수 현재가 |
| `BstpNmixPrdyVrss` | `bstp_nmix_prdy_vrss` | `decimal.Decimal` | 업종 지수 전일 대비 |
| `PrdyVrssSign` | `prdy_vrss_sign` | `string` | 전일 대비 부호 |
| `BstpNmixPrdyCtrt` | `bstp_nmix_prdy_ctrt` | `float64,string` | 업종 지수 전일 대비율 |
| `AcmlTrPbmn` | `acml_tr_pbmn` | `int64,string` | 누적 거래 대금 |
| `AcmlVol` | `acml_vol` | `int64,string` | 누적 거래량 |
| `CntgVol` | `cntg_vol` | `int64,string` | 체결 거래량 |

### 구현 코드

```go
// IndexTickprice 는 국내업종 틱별 지수 (FHPUP02110100) 응답.
//
// 한투 docs: docs/api/국내주식/국내업종_틱별지수.md
// path: /uapi/domestic-stock/v1/quotations/inquire-index-tickprice
//
// StckCntgHour 는 HHMMSS 형식 틱 타임스탬프.
type IndexTickprice struct {
	Output []IndexTickpriceItem `json:"output"`
}

// IndexTickpriceItem 은 응답의 output 한 행 (틱별, 8 fields).
type IndexTickpriceItem struct {
	StckCntgHour     string          `json:"stck_cntg_hour"`             // 주식 체결 시간 HHMMSS
	BstpNmixPrpr     decimal.Decimal `json:"bstp_nmix_prpr"`             // 업종 지수 현재가
	BstpNmixPrdyVrss decimal.Decimal `json:"bstp_nmix_prdy_vrss"`        // 업종 지수 전일 대비
	PrdyVrssSign     string          `json:"prdy_vrss_sign"`             // 전일 대비 부호
	BstpNmixPrdyCtrt float64         `json:"bstp_nmix_prdy_ctrt,string"` // 업종 지수 전일 대비율
	AcmlTrPbmn       int64           `json:"acml_tr_pbmn,string"`        // 누적 거래 대금
	AcmlVol          int64           `json:"acml_vol,string"`            // 누적 거래량
	CntgVol          int64           `json:"cntg_vol,string"`            // 체결 거래량
}

// InquireIndexTickpriceParams 는 국내업종 틱별 지수 조회 파라미터.
type InquireIndexTickpriceParams struct {
	Symbol     string // FID_INPUT_ISCD — 필수, 업종 코드
	MarketCode string // FID_COND_MRKT_DIV_CODE — 빈 값=>"U" (업종)
}

// InquireIndexTickprice 는 국내업종 틱별 지수 호출.
//
// 한투 docs: docs/api/국내주식/국내업종_틱별지수.md
// path: /uapi/domestic-stock/v1/quotations/inquire-index-tickprice (FHPUP02110100)
func (c *Client) InquireIndexTickprice(ctx context.Context, params InquireIndexTickpriceParams) (*IndexTickprice, error) {
	market := params.MarketCode
	if market == "" {
		market = "U"
	}
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/domestic-stock/v1/quotations/inquire-index-tickprice",
		TrID:   "FHPUP02110100",
		Query: map[string]string{
			"FID_INPUT_ISCD":         params.Symbol,
			"FID_COND_MRKT_DIV_CODE": market,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}
	var res IndexTickprice
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse IndexTickprice: %w", err)
	}
	return &res, nil
}
```

### 테스트 코드

```go
func TestClient_InquireIndexTickprice(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/quotations/inquire-index-tickprice`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "index_tickprice_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireIndexTickprice(context.Background(), domestic.InquireIndexTickpriceParams{
		Symbol: "0001",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, "0001", capturedQuery.Get("FID_INPUT_ISCD"))
	assert.Equal(t, "U", capturedQuery.Get("FID_COND_MRKT_DIV_CODE"))

	require.Len(t, res.Output, 2)
	assert.Equal(t, "100015", res.Output[0].StckCntgHour)
	d, _ := decimal.NewFromString("2651.35")
	assert.True(t, d.Equal(res.Output[0].BstpNmixPrpr))
	assert.Equal(t, int64(50000), res.Output[0].CntgVol)
}
```

### Commit

```bash
git commit -m "$(cat <<'EOF'
[feat] domestic — InquireIndexTickprice (국내업종 틱별지수, FHPUP02110100)

- IndexTickprice / IndexTickpriceItem (8 fields, stck_cntg_hour HHMMSS tick timestamp)
- InquireIndexTickpriceParams (Symbol / MarketCode — 2 params only)
- output[] 배열; EP4(Timeprice)와 유사 구조, 파라미터 최소화
- TestClient_InquireIndexTickprice — fixture index_tickprice_success.json

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 5: InquireDailyIndexchartprice (EP6)

**Files:** APPEND to `domestic/industry.go` and `domestic/industry_test.go`

국내업종 일봉 차트 조회. `output1{}` 15 fields + `output2[]` 8 fields. `futs_prdy_*` (선물 전일 시가/고가/저가) 3 필드가 output1 에 포함 — 특이 패턴 주석 필요.

- [ ] Step 1: APPEND test code to `domestic/industry_test.go`
- [ ] Step 2: Verify FAIL — `go test ./domestic/... -run TestClient_InquireDailyIndexchartprice -v`
- [ ] Step 3: APPEND struct + Params + method to `domestic/industry.go`
- [ ] Step 4: Verify PASS — `go test ./domestic/... -run TestClient_InquireDailyIndexchartprice -v`
- [ ] Step 5: `gofmt -w domestic/industry.go domestic/industry_test.go && go vet ./domestic/...`
- [ ] Step 6: commit

### 메서드 매핑
- Path: `/uapi/domestic-stock/v1/quotations/inquire-daily-indexchartprice`
- TR_ID: `FHKUP03500100`
- Params (5): `MarketCode` (`FID_COND_MRKT_DIV_CODE` default `"U"`), `Symbol` (`FID_INPUT_ISCD`), `InputDate1` (`FID_INPUT_DATE_1`), `InputDate2` (`FID_INPUT_DATE_2`), `PeriodDivCode` (`FID_PERIOD_DIV_CODE` D/W/M/Y)

### output1 struct 필드 (15 fields)

| Go field | json tag | Go type | 설명 |
|---|---|---|---|
| `PrdyVrssSign` | `prdy_vrss_sign` | `string` | 전일 대비 부호 |
| `BstpNmixPrdyCtrt` | `bstp_nmix_prdy_ctrt` | `float64,string` | 업종 지수 전일 대비율 |
| `PrdyNmix` | `prdy_nmix` | `decimal.Decimal` | 전일 지수 |
| `AcmlVol` | `acml_vol` | `int64,string` | 누적 거래량 |
| `AcmlTrPbmn` | `acml_tr_pbmn` | `int64,string` | 누적 거래 대금 |
| `HtsKorIsnm` | `hts_kor_isnm` | `string` | HTS 한글 종목명 |
| `BstpNmixPrpr` | `bstp_nmix_prpr` | `decimal.Decimal` | 업종 지수 현재가 |
| `BstpClsCode` | `bstp_cls_code` | `string` | 업종 구분 코드 |
| `PrdyVol` | `prdy_vol` | `int64,string` | 전일 거래량 |
| `BstpNmixOprc` | `bstp_nmix_oprc` | `decimal.Decimal` | 업종 지수 시가 |
| `BstpNmixHgpr` | `bstp_nmix_hgpr` | `decimal.Decimal` | 업종 지수 최고가 |
| `BstpNmixLwpr` | `bstp_nmix_lwpr` | `decimal.Decimal` | 업종 지수 최저가 |
| `FutsPrdyOprc` | `futs_prdy_oprc` | `decimal.Decimal` | 선물 전일 시가 |
| `FutsPrdyHgpr` | `futs_prdy_hgpr` | `decimal.Decimal` | 선물 전일 고가 |
| `FutsPrdyLwpr` | `futs_prdy_lwpr` | `decimal.Decimal` | 선물 전일 저가 |

### output2 struct 필드 (8 fields)

| Go field | json tag | Go type | 설명 |
|---|---|---|---|
| `StckBsopDate` | `stck_bsop_date` | `string` | 영업 일자 |
| `BstpNmixPrpr` | `bstp_nmix_prpr` | `decimal.Decimal` | 업종 지수 현재가 |
| `BstpNmixOprc` | `bstp_nmix_oprc` | `decimal.Decimal` | 업종 지수 시가 |
| `BstpNmixHgpr` | `bstp_nmix_hgpr` | `decimal.Decimal` | 업종 지수 최고가 |
| `BstpNmixLwpr` | `bstp_nmix_lwpr` | `decimal.Decimal` | 업종 지수 최저가 |
| `AcmlVol` | `acml_vol` | `int64,string` | 누적 거래량 |
| `AcmlTrPbmn` | `acml_tr_pbmn` | `int64,string` | 누적 거래 대금 |
| `ModYn` | `mod_yn` | `string` | 수정 여부 (Y/N) |

### 구현 코드

```go
// DailyIndexchartprice 는 국내업종 일봉 차트 (FHKUP03500100) 응답.
//
// 한투 docs: docs/api/국내주식/국내업종_일봉차트.md
// path: /uapi/domestic-stock/v1/quotations/inquire-daily-indexchartprice
//
// output1 에 futs_prdy_* (선물 전일 시가/고가/저가) 3 필드 포함 — 업종+선물 복합 스냅샷.
type DailyIndexchartprice struct {
	Output1 DailyIndexchartpriceSummary `json:"output1"`
	Output2 []DailyIndexchartpriceItem  `json:"output2"`
}

// DailyIndexchartpriceSummary 는 응답의 output1 (현재 스냅샷 + 선물 전일 OHLC, 15 fields).
type DailyIndexchartpriceSummary struct {
	PrdyVrssSign     string          `json:"prdy_vrss_sign"`             // 전일 대비 부호
	BstpNmixPrdyCtrt float64         `json:"bstp_nmix_prdy_ctrt,string"` // 업종 지수 전일 대비율
	PrdyNmix         decimal.Decimal `json:"prdy_nmix"`                  // 전일 지수
	AcmlVol          int64           `json:"acml_vol,string"`            // 누적 거래량
	AcmlTrPbmn       int64           `json:"acml_tr_pbmn,string"`        // 누적 거래 대금
	HtsKorIsnm       string          `json:"hts_kor_isnm"`               // HTS 한글 종목명
	BstpNmixPrpr     decimal.Decimal `json:"bstp_nmix_prpr"`             // 업종 지수 현재가
	BstpClsCode      string          `json:"bstp_cls_code"`              // 업종 구분 코드
	PrdyVol          int64           `json:"prdy_vol,string"`            // 전일 거래량
	BstpNmixOprc     decimal.Decimal `json:"bstp_nmix_oprc"`             // 업종 지수 시가
	BstpNmixHgpr     decimal.Decimal `json:"bstp_nmix_hgpr"`             // 업종 지수 최고가
	BstpNmixLwpr     decimal.Decimal `json:"bstp_nmix_lwpr"`             // 업종 지수 최저가
	FutsPrdyOprc     decimal.Decimal `json:"futs_prdy_oprc"`             // 선물 전일 시가
	FutsPrdyHgpr     decimal.Decimal `json:"futs_prdy_hgpr"`             // 선물 전일 고가
	FutsPrdyLwpr     decimal.Decimal `json:"futs_prdy_lwpr"`             // 선물 전일 저가
}

// DailyIndexchartpriceItem 은 응답의 output2 한 행 (일봉, 8 fields).
type DailyIndexchartpriceItem struct {
	StckBsopDate string          `json:"stck_bsop_date"` // 영업 일자
	BstpNmixPrpr decimal.Decimal `json:"bstp_nmix_prpr"` // 업종 지수 현재가
	BstpNmixOprc decimal.Decimal `json:"bstp_nmix_oprc"` // 업종 지수 시가
	BstpNmixHgpr decimal.Decimal `json:"bstp_nmix_hgpr"` // 업종 지수 최고가
	BstpNmixLwpr decimal.Decimal `json:"bstp_nmix_lwpr"` // 업종 지수 최저가
	AcmlVol      int64           `json:"acml_vol,string"` // 누적 거래량
	AcmlTrPbmn   int64           `json:"acml_tr_pbmn,string"` // 누적 거래 대금
	ModYn        string          `json:"mod_yn"`          // 수정 여부
}

// InquireDailyIndexchartpriceParams 는 국내업종 일봉 차트 조회 파라미터.
type InquireDailyIndexchartpriceParams struct {
	MarketCode    string // FID_COND_MRKT_DIV_CODE — 빈 값=>"U" (업종)
	Symbol        string // FID_INPUT_ISCD — 필수, 업종 코드
	InputDate1    string // FID_INPUT_DATE_1 — 조회 시작일 YYYYMMDD
	InputDate2    string // FID_INPUT_DATE_2 — 조회 종료일 YYYYMMDD
	PeriodDivCode string // FID_PERIOD_DIV_CODE — D:일 W:주 M:월 Y:년
}

// InquireDailyIndexchartprice 는 국내업종 일봉 차트 호출.
//
// 한투 docs: docs/api/국내주식/국내업종_일봉차트.md
// path: /uapi/domestic-stock/v1/quotations/inquire-daily-indexchartprice (FHKUP03500100)
func (c *Client) InquireDailyIndexchartprice(ctx context.Context, params InquireDailyIndexchartpriceParams) (*DailyIndexchartprice, error) {
	market := params.MarketCode
	if market == "" {
		market = "U"
	}
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/domestic-stock/v1/quotations/inquire-daily-indexchartprice",
		TrID:   "FHKUP03500100",
		Query: map[string]string{
			"FID_COND_MRKT_DIV_CODE": market,
			"FID_INPUT_ISCD":         params.Symbol,
			"FID_INPUT_DATE_1":       params.InputDate1,
			"FID_INPUT_DATE_2":       params.InputDate2,
			"FID_PERIOD_DIV_CODE":    params.PeriodDivCode,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}
	var res DailyIndexchartprice
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse DailyIndexchartprice: %w", err)
	}
	return &res, nil
}
```

### 테스트 코드

```go
func TestClient_InquireDailyIndexchartprice(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/quotations/inquire-daily-indexchartprice`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "daily_indexchartprice_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireDailyIndexchartprice(context.Background(), domestic.InquireDailyIndexchartpriceParams{
		Symbol:        "0001",
		InputDate1:    "20260401",
		InputDate2:    "20260505",
		PeriodDivCode: "D",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, "U", capturedQuery.Get("FID_COND_MRKT_DIV_CODE"))
	assert.Equal(t, "0001", capturedQuery.Get("FID_INPUT_ISCD"))
	assert.Equal(t, "20260401", capturedQuery.Get("FID_INPUT_DATE_1"))
	assert.Equal(t, "20260505", capturedQuery.Get("FID_INPUT_DATE_2"))
	assert.Equal(t, "D", capturedQuery.Get("FID_PERIOD_DIV_CODE"))

	assert.Equal(t, "코스피", res.Output1.HtsKorIsnm)
	assert.Equal(t, "0001", res.Output1.BstpClsCode)
	futs, _ := decimal.NewFromString("355.50")
	assert.True(t, futs.Equal(res.Output1.FutsPrdyOprc))

	require.Len(t, res.Output2, 2)
	assert.Equal(t, "20260505", res.Output2[0].StckBsopDate)
	assert.Equal(t, "N", res.Output2[0].ModYn)
	assert.Equal(t, int64(350000000), res.Output2[0].AcmlVol)
}
```

### Commit

```bash
git commit -m "$(cat <<'EOF'
[feat] domestic — InquireDailyIndexchartprice (국내업종 일봉차트, FHKUP03500100)

- DailyIndexchartprice / DailyIndexchartpriceSummary (15 fields incl. futs_prdy_* 3 fields) / DailyIndexchartpriceItem (8 fields, mod_yn)
- InquireDailyIndexchartpriceParams (MarketCode/Symbol/InputDate1/InputDate2/PeriodDivCode)
- futs_prdy_oprc/hgpr/lwpr (선물 전일 OHLC) decimal 매핑 in output1
- TestClient_InquireDailyIndexchartprice — fixture daily_indexchartprice_success.json

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 6: InquireTimeIndexchartprice (EP7)

**Files:** APPEND to `domestic/industry.go` and `domestic/industry_test.go`

국내업종 분봉 차트 조회 (intraday). `output1{}` 16 fields + `output2[]` 8 fields. EP6 일봉 차트와 유사하나 output1에 `bstp_nmix_prdy_vrss` 필드 추가(16 fields), output2에 `stck_cntg_hour` 타임스탬프 + `cntg_vol` 포함(EP6 output2의 `acml_vol`/`mod_yn` 대신).

- [ ] Step 1: APPEND test code to `domestic/industry_test.go`
- [ ] Step 2: Verify FAIL — `go test ./domestic/... -run TestClient_InquireTimeIndexchartprice -v`
- [ ] Step 3: APPEND struct + Params + method to `domestic/industry.go`
- [ ] Step 4: Verify PASS — `go test ./domestic/... -run TestClient_InquireTimeIndexchartprice -v`
- [ ] Step 5: `gofmt -w domestic/industry.go domestic/industry_test.go && go vet ./domestic/...`
- [ ] Step 6: commit

### 메서드 매핑
- Path: `/uapi/domestic-stock/v1/quotations/inquire-time-indexchartprice`
- TR_ID: `FHKUP03500200`
- Params (5): `MarketCode` (`FID_COND_MRKT_DIV_CODE` default `"U"`), `EtcClsCode` (`FID_ETC_CLS_CODE` 0/1), `Symbol` (`FID_INPUT_ISCD`), `InputHour1` (`FID_INPUT_HOUR_1` 30/60/600/3600), `PwDataIncuYn` (`FID_PW_DATA_INCU_YN` Y/N)

### output1 struct 필드 (16 fields)

| Go field | json tag | Go type | 설명 |
|---|---|---|---|
| `BstpNmixPrdyVrss` | `bstp_nmix_prdy_vrss` | `decimal.Decimal` | 업종 지수 전일 대비 |
| `PrdyVrssSign` | `prdy_vrss_sign` | `string` | 전일 대비 부호 |
| `BstpNmixPrdyCtrt` | `bstp_nmix_prdy_ctrt` | `float64,string` | 업종 지수 전일 대비율 |
| `PrdyNmix` | `prdy_nmix` | `decimal.Decimal` | 전일 지수 |
| `AcmlVol` | `acml_vol` | `int64,string` | 누적 거래량 |
| `AcmlTrPbmn` | `acml_tr_pbmn` | `int64,string` | 누적 거래 대금 |
| `HtsKorIsnm` | `hts_kor_isnm` | `string` | HTS 한글 종목명 |
| `BstpNmixPrpr` | `bstp_nmix_prpr` | `decimal.Decimal` | 업종 지수 현재가 |
| `BstpClsCode` | `bstp_cls_code` | `string` | 업종 구분 코드 |
| `PrdyVol` | `prdy_vol` | `int64,string` | 전일 거래량 |
| `BstpNmixOprc` | `bstp_nmix_oprc` | `decimal.Decimal` | 업종 지수 시가 |
| `BstpNmixHgpr` | `bstp_nmix_hgpr` | `decimal.Decimal` | 업종 지수 최고가 |
| `BstpNmixLwpr` | `bstp_nmix_lwpr` | `decimal.Decimal` | 업종 지수 최저가 |
| `FutsPrdyOprc` | `futs_prdy_oprc` | `decimal.Decimal` | 선물 전일 시가 |
| `FutsPrdyHgpr` | `futs_prdy_hgpr` | `decimal.Decimal` | 선물 전일 고가 |
| `FutsPrdyLwpr` | `futs_prdy_lwpr` | `decimal.Decimal` | 선물 전일 저가 |

### output2 struct 필드 (8 fields)

| Go field | json tag | Go type | 설명 |
|---|---|---|---|
| `StckBsopDate` | `stck_bsop_date` | `string` | 영업 일자 |
| `StckCntgHour` | `stck_cntg_hour` | `string` | 주식 체결 시간 (HHMMSS) |
| `BstpNmixPrpr` | `bstp_nmix_prpr` | `decimal.Decimal` | 업종 지수 현재가 |
| `BstpNmixOprc` | `bstp_nmix_oprc` | `decimal.Decimal` | 업종 지수 시가 |
| `BstpNmixHgpr` | `bstp_nmix_hgpr` | `decimal.Decimal` | 업종 지수 최고가 |
| `BstpNmixLwpr` | `bstp_nmix_lwpr` | `decimal.Decimal` | 업종 지수 최저가 |
| `CntgVol` | `cntg_vol` | `int64,string` | 체결 거래량 |
| `AcmlTrPbmn` | `acml_tr_pbmn` | `int64,string` | 누적 거래 대금 |

### 구현 코드

```go
// TimeIndexchartprice 는 국내업종 분봉 차트 (FHKUP03500200) 응답.
//
// 한투 docs: docs/api/국내주식/국내업종_분봉차트.md
// path: /uapi/domestic-stock/v1/quotations/inquire-time-indexchartprice
//
// EP6(DailyIndexchartprice) 유사 구조 — output1 에 bstp_nmix_prdy_vrss 추가(16 fields).
// output2 에 stck_cntg_hour 타임스탬프 포함.
type TimeIndexchartprice struct {
	Output1 TimeIndexchartpriceSummary `json:"output1"`
	Output2 []TimeIndexchartpriceItem  `json:"output2"`
}

// TimeIndexchartpriceSummary 는 응답의 output1 (현재 스냅샷 + 선물 전일 OHLC, 16 fields).
type TimeIndexchartpriceSummary struct {
	BstpNmixPrdyVrss decimal.Decimal `json:"bstp_nmix_prdy_vrss"`        // 업종 지수 전일 대비
	PrdyVrssSign     string          `json:"prdy_vrss_sign"`             // 전일 대비 부호
	BstpNmixPrdyCtrt float64         `json:"bstp_nmix_prdy_ctrt,string"` // 업종 지수 전일 대비율
	PrdyNmix         decimal.Decimal `json:"prdy_nmix"`                  // 전일 지수
	AcmlVol          int64           `json:"acml_vol,string"`            // 누적 거래량
	AcmlTrPbmn       int64           `json:"acml_tr_pbmn,string"`        // 누적 거래 대금
	HtsKorIsnm       string          `json:"hts_kor_isnm"`               // HTS 한글 종목명
	BstpNmixPrpr     decimal.Decimal `json:"bstp_nmix_prpr"`             // 업종 지수 현재가
	BstpClsCode      string          `json:"bstp_cls_code"`              // 업종 구분 코드
	PrdyVol          int64           `json:"prdy_vol,string"`            // 전일 거래량
	BstpNmixOprc     decimal.Decimal `json:"bstp_nmix_oprc"`             // 업종 지수 시가
	BstpNmixHgpr     decimal.Decimal `json:"bstp_nmix_hgpr"`             // 업종 지수 최고가
	BstpNmixLwpr     decimal.Decimal `json:"bstp_nmix_lwpr"`             // 업종 지수 최저가
	FutsPrdyOprc     decimal.Decimal `json:"futs_prdy_oprc"`             // 선물 전일 시가
	FutsPrdyHgpr     decimal.Decimal `json:"futs_prdy_hgpr"`             // 선물 전일 고가
	FutsPrdyLwpr     decimal.Decimal `json:"futs_prdy_lwpr"`             // 선물 전일 저가
}

// TimeIndexchartpriceItem 은 응답의 output2 한 행 (분봉, 8 fields).
type TimeIndexchartpriceItem struct {
	StckBsopDate string          `json:"stck_bsop_date"`      // 영업 일자
	StckCntgHour string          `json:"stck_cntg_hour"`      // 주식 체결 시간 HHMMSS
	BstpNmixPrpr decimal.Decimal `json:"bstp_nmix_prpr"`      // 업종 지수 현재가
	BstpNmixOprc decimal.Decimal `json:"bstp_nmix_oprc"`      // 업종 지수 시가
	BstpNmixHgpr decimal.Decimal `json:"bstp_nmix_hgpr"`      // 업종 지수 최고가
	BstpNmixLwpr decimal.Decimal `json:"bstp_nmix_lwpr"`      // 업종 지수 최저가
	CntgVol      int64           `json:"cntg_vol,string"`     // 체결 거래량
	AcmlTrPbmn   int64           `json:"acml_tr_pbmn,string"` // 누적 거래 대금
}

// InquireTimeIndexchartpriceParams 는 국내업종 분봉 차트 조회 파라미터.
type InquireTimeIndexchartpriceParams struct {
	MarketCode   string // FID_COND_MRKT_DIV_CODE — 빈 값=>"U" (업종)
	EtcClsCode   string // FID_ETC_CLS_CODE — 0:현물 1:선물
	Symbol       string // FID_INPUT_ISCD — 필수, 업종 코드
	InputHour1   string // FID_INPUT_HOUR_1 — 집계 단위: "30"/"60"/"600"/"3600"
	PwDataIncuYn string // FID_PW_DATA_INCU_YN — Y:과거데이터포함 N:미포함
}

// InquireTimeIndexchartprice 는 국내업종 분봉 차트 호출.
//
// 한투 docs: docs/api/국내주식/국내업종_분봉차트.md
// path: /uapi/domestic-stock/v1/quotations/inquire-time-indexchartprice (FHKUP03500200)
func (c *Client) InquireTimeIndexchartprice(ctx context.Context, params InquireTimeIndexchartpriceParams) (*TimeIndexchartprice, error) {
	market := params.MarketCode
	if market == "" {
		market = "U"
	}
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/domestic-stock/v1/quotations/inquire-time-indexchartprice",
		TrID:   "FHKUP03500200",
		Query: map[string]string{
			"FID_COND_MRKT_DIV_CODE": market,
			"FID_ETC_CLS_CODE":       params.EtcClsCode,
			"FID_INPUT_ISCD":         params.Symbol,
			"FID_INPUT_HOUR_1":       params.InputHour1,
			"FID_PW_DATA_INCU_YN":    params.PwDataIncuYn,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}
	var res TimeIndexchartprice
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse TimeIndexchartprice: %w", err)
	}
	return &res, nil
}
```

### 테스트 코드

```go
func TestClient_InquireTimeIndexchartprice(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/quotations/inquire-time-indexchartprice`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "time_indexchartprice_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireTimeIndexchartprice(context.Background(), domestic.InquireTimeIndexchartpriceParams{
		EtcClsCode:   "0",
		Symbol:       "0001",
		InputHour1:   "60",
		PwDataIncuYn: "Y",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, "U", capturedQuery.Get("FID_COND_MRKT_DIV_CODE"))
	assert.Equal(t, "0", capturedQuery.Get("FID_ETC_CLS_CODE"))
	assert.Equal(t, "0001", capturedQuery.Get("FID_INPUT_ISCD"))
	assert.Equal(t, "60", capturedQuery.Get("FID_INPUT_HOUR_1"))
	assert.Equal(t, "Y", capturedQuery.Get("FID_PW_DATA_INCU_YN"))

	assert.Equal(t, "코스피", res.Output1.HtsKorIsnm)
	futs, _ := decimal.NewFromString("355.50")
	assert.True(t, futs.Equal(res.Output1.FutsPrdyOprc))
	vrss, _ := decimal.NewFromString("-12.30")
	assert.True(t, vrss.Equal(res.Output1.BstpNmixPrdyVrss))

	require.Len(t, res.Output2, 2)
	assert.Equal(t, "20260505", res.Output2[0].StckBsopDate)
	assert.Equal(t, "100000", res.Output2[0].StckCntgHour)
	assert.Equal(t, int64(800000), res.Output2[0].CntgVol)
}
```

### Commit

```bash
git commit -m "$(cat <<'EOF'
[feat] domestic — InquireTimeIndexchartprice (국내업종 분봉차트, FHKUP03500200)

- TimeIndexchartprice / TimeIndexchartpriceSummary (16 fields, bstp_nmix_prdy_vrss 추가 vs EP6) / TimeIndexchartpriceItem (8 fields, stck_cntg_hour + cntg_vol)
- InquireTimeIndexchartpriceParams (MarketCode/EtcClsCode/Symbol/InputHour1/PwDataIncuYn)
- futs_prdy_* 선물 전일 OHLC in output1 (EP6 패턴 동일)
- TestClient_InquireTimeIndexchartprice — fixture time_indexchartprice_success.json

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 7: ExpTotalIndex (EP8)

**Files:** APPEND to `domestic/industry.go` and `domestic/industry_test.go`

예상 전체 지수 조회. **ANOMALY: 쿼리 파라미터 키가 모두 소문자** (`fid_*`) — 다른 모든 KIS 엔드포인트는 `FID_*` 대문자 사용. `output1{}` 9 fields + `output2[]` 10 fields. output1 에서 `prdy_ctrt` (short form) 사용 — `bstp_nmix_prdy_ctrt` 아님. `fid_cond_scr_div_code="11175"` Params 미노출 hardcoded.

- [ ] Step 1: APPEND test code to `domestic/industry_test.go`
- [ ] Step 2: Verify FAIL — `go test ./domestic/... -run TestClient_ExpTotalIndex -v`
- [ ] Step 3: APPEND struct + Params + method to `domestic/industry.go`
- [ ] Step 4: Verify PASS — `go test ./domestic/... -run TestClient_ExpTotalIndex -v`
- [ ] Step 5: `gofmt -w domestic/industry.go domestic/industry_test.go && go vet ./domestic/...`
- [ ] Step 6: commit

### 메서드 매핑
- Path: `/uapi/domestic-stock/v1/quotations/exp-total-index`
- TR_ID: `FHKUP11750000`
- Params (5, **ALL lowercase wire keys**):
  - `MrktClsCode` → `fid_mrkt_cls_code`
  - `MarketCode` → `fid_cond_mrkt_div_code` (default `"U"`)
  - `CondScrDivCode` → `fid_cond_scr_div_code` (default `"11175"` — **hardcoded in Query map, Params 미노출**)
  - `Symbol` → `fid_input_iscd`
  - `MkopClsCode` → `fid_mkop_cls_code` (1/2)

### output1 struct 필드 (9 fields)

| Go field | json tag | Go type | 설명 |
|---|---|---|---|
| `BstpNmixPrpr` | `bstp_nmix_prpr` | `decimal.Decimal` | 업종 지수 현재가 |
| `BstpNmixPrdyVrss` | `bstp_nmix_prdy_vrss` | `decimal.Decimal` | 업종 지수 전일 대비 |
| `PrdyVrssSign` | `prdy_vrss_sign` | `string` | 전일 대비 부호 |
| `PrdyCtrt` | `prdy_ctrt` | `float64,string` | 전일 대비율 (**short form** — `bstp_nmix_prdy_ctrt` 아님) |
| `AcmlVol` | `acml_vol` | `int64,string` | 누적 거래량 |
| `AscnIssuCnt` | `ascn_issu_cnt` | `string` | 상승 종목 수 |
| `DownIssuCnt` | `down_issu_cnt` | `string` | 하락 종목 수 |
| `StnrIssuCnt` | `stnr_issu_cnt` | `string` | 보합 종목 수 |
| `BstpClsCode` | `bstp_cls_code` | `string` | 업종 구분 코드 |

### output2 struct 필드 (10 fields)

| Go field | json tag | Go type | 설명 |
|---|---|---|---|
| `HtsKorIsnm` | `hts_kor_isnm` | `string` | HTS 한글 종목명 |
| `BstpNmixPrpr` | `bstp_nmix_prpr` | `decimal.Decimal` | 업종 지수 현재가 |
| `BstpNmixPrdyVrss` | `bstp_nmix_prdy_vrss` | `decimal.Decimal` | 업종 지수 전일 대비 |
| `PrdyVrssSign` | `prdy_vrss_sign` | `string` | 전일 대비 부호 |
| `BstpNmixPrdyCtrt` | `bstp_nmix_prdy_ctrt` | `float64,string` | 업종 지수 전일 대비율 |
| `AcmlVol` | `acml_vol` | `int64,string` | 누적 거래량 |
| `NmixSdpr` | `nmix_sdpr` | `decimal.Decimal` | 지수 기준가 |
| `AscnIssuCnt` | `ascn_issu_cnt` | `string` | 상승 종목 수 |
| `StnrIssuCnt` | `stnr_issu_cnt` | `string` | 보합 종목 수 |
| `DownIssuCnt` | `down_issu_cnt` | `string` | 하락 종목 수 |

### 구현 코드

```go
// ExpTotalIndex 는 예상 전체 지수 (FHKUP11750000) 응답.
//
// 한투 docs: docs/api/국내주식/예상전체지수.md
// path: /uapi/domestic-stock/v1/quotations/exp-total-index
//
// ANOMALY: 쿼리 파라미터 키가 모두 소문자 fid_* (다른 엔드포인트는 FID_* 대문자).
// ANOMALY: output1 의 비율 필드명이 prdy_ctrt (short form) — bstp_nmix_prdy_ctrt 아님.
// fid_cond_scr_div_code="11175" 는 Query map 에 하드코딩 (Params 미노출).
type ExpTotalIndex struct {
	Output1 ExpTotalIndexSummary `json:"output1"`
	Output2 []ExpTotalIndexItem  `json:"output2"`
}

// ExpTotalIndexSummary 는 응답의 output1 (9 fields).
//
// PrdyCtrt 는 prdy_ctrt (short form) — bstp_nmix_prdy_ctrt 와 다른 필드명 주의.
type ExpTotalIndexSummary struct {
	BstpNmixPrpr     decimal.Decimal `json:"bstp_nmix_prpr"`      // 업종 지수 현재가
	BstpNmixPrdyVrss decimal.Decimal `json:"bstp_nmix_prdy_vrss"` // 업종 지수 전일 대비
	PrdyVrssSign     string          `json:"prdy_vrss_sign"`      // 전일 대비 부호
	PrdyCtrt         float64         `json:"prdy_ctrt,string"`    // 전일 대비율 (short form)
	AcmlVol          int64           `json:"acml_vol,string"`     // 누적 거래량
	AscnIssuCnt      string          `json:"ascn_issu_cnt"`       // 상승 종목 수
	DownIssuCnt      string          `json:"down_issu_cnt"`       // 하락 종목 수
	StnrIssuCnt      string          `json:"stnr_issu_cnt"`       // 보합 종목 수
	BstpClsCode      string          `json:"bstp_cls_code"`       // 업종 구분 코드
}

// ExpTotalIndexItem 은 응답의 output2 한 행 (10 fields).
type ExpTotalIndexItem struct {
	HtsKorIsnm       string          `json:"hts_kor_isnm"`               // HTS 한글 종목명
	BstpNmixPrpr     decimal.Decimal `json:"bstp_nmix_prpr"`             // 업종 지수 현재가
	BstpNmixPrdyVrss decimal.Decimal `json:"bstp_nmix_prdy_vrss"`        // 업종 지수 전일 대비
	PrdyVrssSign     string          `json:"prdy_vrss_sign"`             // 전일 대비 부호
	BstpNmixPrdyCtrt float64         `json:"bstp_nmix_prdy_ctrt,string"` // 업종 지수 전일 대비율
	AcmlVol          int64           `json:"acml_vol,string"`            // 누적 거래량
	NmixSdpr         decimal.Decimal `json:"nmix_sdpr"`                  // 지수 기준가
	AscnIssuCnt      string          `json:"ascn_issu_cnt"`              // 상승 종목 수
	StnrIssuCnt      string          `json:"stnr_issu_cnt"`              // 보합 종목 수
	DownIssuCnt      string          `json:"down_issu_cnt"`              // 하락 종목 수
}

// ExpTotalIndexParams 는 예상 전체 지수 조회 파라미터.
//
// 주의: wire key 가 모두 소문자 fid_* (KIS anomaly).
// fid_cond_scr_div_code="11175" 는 내부 hardcoded (노출 안함).
type ExpTotalIndexParams struct {
	MrktClsCode string // fid_mrkt_cls_code — 시장 구분 코드
	MarketCode  string // fid_cond_mrkt_div_code — 빈 값=>"U" (업종)
	Symbol      string // fid_input_iscd — 필수, 업종 코드
	MkopClsCode string // fid_mkop_cls_code — 1:매수 2:매도
}

// ExpTotalIndex 는 예상 전체 지수 호출.
//
// 한투 docs: docs/api/국내주식/예상전체지수.md
// path: /uapi/domestic-stock/v1/quotations/exp-total-index (FHKUP11750000)
//
// ANOMALY: 쿼리 파라미터 키 소문자 fid_* 사용. fid_cond_scr_div_code="11175" hardcoded.
func (c *Client) ExpTotalIndex(ctx context.Context, params ExpTotalIndexParams) (*ExpTotalIndex, error) {
	market := params.MarketCode
	if market == "" {
		market = "U"
	}
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/domestic-stock/v1/quotations/exp-total-index",
		TrID:   "FHKUP11750000",
		Query: map[string]string{
			"fid_mrkt_cls_code":      params.MrktClsCode,
			"fid_cond_mrkt_div_code": market,
			"fid_cond_scr_div_code":  "11175",
			"fid_input_iscd":         params.Symbol,
			"fid_mkop_cls_code":      params.MkopClsCode,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}
	var res ExpTotalIndex
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse ExpTotalIndex: %w", err)
	}
	return &res, nil
}
```

### 테스트 코드

```go
func TestClient_ExpTotalIndex(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/quotations/exp-total-index`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "exp_total_index_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.ExpTotalIndex(context.Background(), domestic.ExpTotalIndexParams{
		MrktClsCode: "K",
		Symbol:      "0001",
		MkopClsCode: "1",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	// ANOMALY: lowercase query param keys
	assert.Equal(t, "K", capturedQuery.Get("fid_mrkt_cls_code"))
	assert.Equal(t, "U", capturedQuery.Get("fid_cond_mrkt_div_code"))
	assert.Equal(t, "11175", capturedQuery.Get("fid_cond_scr_div_code"))
	assert.Equal(t, "0001", capturedQuery.Get("fid_input_iscd"))
	assert.Equal(t, "1", capturedQuery.Get("fid_mkop_cls_code"))

	d, _ := decimal.NewFromString("2650.45")
	assert.True(t, d.Equal(res.Output1.BstpNmixPrpr))
	// prdy_ctrt (short form) — bstp_nmix_prdy_ctrt 아님
	assert.InDelta(t, -0.46, res.Output1.PrdyCtrt, 0.001)
	assert.Equal(t, "315", res.Output1.AscnIssuCnt)
	assert.Equal(t, "0001", res.Output1.BstpClsCode)

	require.Len(t, res.Output2, 2)
	assert.Equal(t, "코스피", res.Output2[0].HtsKorIsnm)
	sdpr, _ := decimal.NewFromString("2662.75")
	assert.True(t, sdpr.Equal(res.Output2[0].NmixSdpr))
	assert.InDelta(t, -0.46, res.Output2[0].BstpNmixPrdyCtrt, 0.001)
}
```

### Commit

```bash
git commit -m "$(cat <<'EOF'
[feat] domestic — ExpTotalIndex (예상전체지수, FHKUP11750000)

- ExpTotalIndex / ExpTotalIndexSummary (9 fields) / ExpTotalIndexItem (10 fields)
- ExpTotalIndexParams (MrktClsCode/MarketCode/Symbol/MkopClsCode)
- ANOMALY: 쿼리 파라미터 키 소문자 fid_* (KIS 유일 예외 — 다른 EP는 FID_* 대문자)
- ANOMALY: output1 PrdyCtrt = prdy_ctrt (short form, bstp_nmix_prdy_ctrt 아님)
- fid_cond_scr_div_code="11175" hardcoded in Query map (Params 미노출)
- NmixSdpr (지수 기준가) decimal in output2
- TestClient_ExpTotalIndex — fixture exp_total_index_success.json

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

---

## Task 8: ExpIndexTrend (EP9) — output array (7 fields)

**파일**: APPEND to `domestic/industry.go` and `domestic/industry_test.go`

- Path: `/uapi/domestic-stock/v1/quotations/exp-index-trend`
- TR_ID: `FHPST01840000`
- Output: `output []ExpIndexTrendItem` (7 fields)

> **ANOMALY — KIS 문서 Korean 라벨 scrambled**: `stck_cntg_hour` 가 "주식 단축 종목코드" 로, `bstp_nmix_prpr` 가 "HTS 한글 종목명" 으로 잘못 라벨링되어 있음. **Field 명(영문)은 정확** — Korean description column 만 무시하고 field name 을 source of truth 로 사용할 것.

### Params (4)

| Go field | Wire key | 예시 값 |
|---|---|---|
| MkopClsCode | FID_MKOP_CLS_CODE | `1` or `2` |
| InputHour1 | FID_INPUT_HOUR_1 | `10` / `30` / `60` / `600` |
| Symbol | FID_INPUT_ISCD | e.g. `"0001"` |
| MarketCode | FID_COND_MRKT_DIV_CODE | `"U"` (default) |

### Struct Fields (7)

| JSON key | Go field | Go type |
|---|---|---|
| stck_cntg_hour | StckCntgHour | string |
| bstp_nmix_prpr | BstpNmixPrpr | decimal.Decimal |
| prdy_vrss_sign | PrdyVrssSign | string |
| bstp_nmix_prdy_vrss | BstpNmixPrdyVrss | decimal.Decimal |
| prdy_ctrt | PrdyCtrt | float64 (string in JSON) |
| acml_vol | AcmlVol | int64 (string in JSON) |
| acml_tr_pbmn | AcmlTrPbmn | int64 (string in JSON) |

> `prdy_ctrt` — short form (NOT `bstp_nmix_prdy_ctrt`). EP8 (ExpTotalIndex) 와 동일 패턴.

### Step 1 — 테스트 먼저 작성 (RED)

`domestic/industry_test.go` 에 아래 테스트 APPEND:

```go
func TestClient_ExpIndexTrend(t *testing.T) {
    c, mock := newTestClient(t)
    defer mock.DeactivateAndReset()

    fixture := loadFixture(t, "testdata/exp_index_trend_success.json")
    mock.RegisterResponder(http.MethodGet,
        "https://openapi.koreainvestment.com:9443/uapi/domestic-stock/v1/quotations/exp-index-trend",
        httpmock.NewBytesResponder(http.StatusOK, fixture))

    params := domestic.ExpIndexTrendParams{
        MkopClsCode: "1",
        InputHour1:  "10",
        Symbol:      "0001",
        MarketCode:  "U",
    }
    resp, err := c.Domestic.ExpIndexTrend(context.Background(), params)
    require.NoError(t, err)
    require.NotEmpty(t, resp.Output)

    item := resp.Output[0]
    assert.NotEmpty(t, item.StckCntgHour)
    assert.False(t, item.BstpNmixPrpr.IsZero())
    assert.NotEmpty(t, item.PrdyVrssSign)
    assert.NotZero(t, item.AcmlVol)
}

func TestClient_ExpIndexTrend_InvalidJSON(t *testing.T) {
    c, mock := newTestClient(t)
    defer mock.DeactivateAndReset()

    mock.RegisterResponder(http.MethodGet,
        "https://openapi.koreainvestment.com:9443/uapi/domestic-stock/v1/quotations/exp-index-trend",
        httpmock.NewStringResponder(http.StatusOK, `{invalid}`))

    _, err := c.Domestic.ExpIndexTrend(context.Background(), domestic.ExpIndexTrendParams{
        MkopClsCode: "1",
        InputHour1:  "10",
        Symbol:      "0001",
        MarketCode:  "U",
    })
    require.Error(t, err)
}
```

### Step 2 — `go test ./domestic/... -run TestClient_ExpIndexTrend` → FAIL (컴파일 에러)

### Step 3 — 구현 (`domestic/industry.go` APPEND)

```go
// ExpIndexTrendItem — 예상체결지수 추이 행
// ANOMALY: KIS 문서의 Korean field 라벨이 scrambled 되어 있음.
//   e.g., stck_cntg_hour → "주식 단축 종목코드" (잘못된 라벨)
//         bstp_nmix_prpr → "HTS 한글 종목명" (잘못된 라벨)
// Field name(영문)은 정확하므로 field name 기준으로 구현.
type ExpIndexTrendItem struct {
    StckCntgHour    string          `json:"stck_cntg_hour"`    // 체결 시간
    BstpNmixPrpr    decimal.Decimal `json:"bstp_nmix_prpr"`    // 업종 지수 현재가
    PrdyVrssSign    string          `json:"prdy_vrss_sign"`    // 전일 대비 부호
    BstpNmixPrdyVrss decimal.Decimal `json:"bstp_nmix_prdy_vrss"` // 업종 지수 전일 대비
    PrdyCtrt        FloatString     `json:"prdy_ctrt"`          // 전일 대비율 (short form)
    AcmlVol         Int64String     `json:"acml_vol"`           // 누적 거래량
    AcmlTrPbmn      Int64String     `json:"acml_tr_pbmn"`       // 누적 거래 대금
}

// ExpIndexTrendResponse — FHPST01840000 응답
type ExpIndexTrendResponse struct {
    RtCd   string             `json:"rt_cd"`
    MsgCd  string             `json:"msg_cd"`
    Msg1   string             `json:"msg1"`
    Output []ExpIndexTrendItem `json:"output"`
}

// ExpIndexTrendParams — 예상체결지수 추이 요청 파라미터
type ExpIndexTrendParams struct {
    MkopClsCode string // FID_MKOP_CLS_CODE: 1 or 2
    InputHour1  string // FID_INPUT_HOUR_1: 10/30/60/600
    Symbol      string // FID_INPUT_ISCD: 종목코드
    MarketCode  string // FID_COND_MRKT_DIV_CODE: default "U"
}

// ExpIndexTrend — 예상체결지수 추이 (FHPST01840000)
// ANOMALY: KIS 문서 Korean field 라벨이 scrambled — field name 기준으로 구현.
// path: /uapi/domestic-stock/v1/quotations/exp-index-trend
func (s *DomesticService) ExpIndexTrend(ctx context.Context, params ExpIndexTrendParams) (*ExpIndexTrendResponse, error) {
    if params.MarketCode == "" {
        params.MarketCode = "U"
    }
    var result ExpIndexTrendResponse
    resp, err := s.client.R().
        SetContext(ctx).
        SetHeader("tr_id", "FHPST01840000").
        SetQueryParams(map[string]string{
            "FID_MKOP_CLS_CODE":      params.MkopClsCode,
            "FID_INPUT_HOUR_1":       params.InputHour1,
            "FID_INPUT_ISCD":         params.Symbol,
            "FID_COND_MRKT_DIV_CODE": params.MarketCode,
        }).
        SetResult(&result).
        Get("/uapi/domestic-stock/v1/quotations/exp-index-trend")
    if err != nil {
        return nil, err
    }
    if resp.IsError() {
        return nil, fmt.Errorf("ExpIndexTrend: %s %s", result.MsgCd, result.Msg1)
    }
    return &result, nil
}
```

**Fixture** `domestic/testdata/exp_index_trend_success.json`:

```json
{
  "rt_cd": "0",
  "msg_cd": "OPSP0000",
  "msg1": "정상처리 되었습니다.",
  "output": [
    {
      "stck_cntg_hour": "090000",
      "bstp_nmix_prpr": "2550.25",
      "prdy_vrss_sign": "2",
      "bstp_nmix_prdy_vrss": "12.50",
      "prdy_ctrt": "0.49",
      "acml_vol": "123456789",
      "acml_tr_pbmn": "987654321000"
    },
    {
      "stck_cntg_hour": "090010",
      "bstp_nmix_prpr": "2551.00",
      "prdy_vrss_sign": "2",
      "bstp_nmix_prdy_vrss": "13.25",
      "prdy_ctrt": "0.52",
      "acml_vol": "234567890",
      "acml_tr_pbmn": "1098765432000"
    }
  ]
}
```

### Step 4 — `go test ./domestic/... -run TestClient_ExpIndexTrend` → PASS

### Step 5 — `gofmt -w domestic/industry.go domestic/industry_test.go`

### Step 6 — Commit

```bash
git add domestic/industry.go domestic/industry_test.go domestic/testdata/exp_index_trend_success.json
git commit -m "$(cat <<'EOF'
[feat] domestic — InquireExpIndexTrend (예상체결지수 추이, FHPST01840000)

- ExpIndexTrendItem (7 fields) / ExpIndexTrendResponse / ExpIndexTrendParams
- ANOMALY: KIS 문서 Korean field 라벨이 scrambled (e.g., stck_cntg_hour → "주식 단축 종목코드")
  Field name(영문)은 정확 — Korean 라벨만 무시
- prdy_ctrt = short form (NOT bstp_nmix_prdy_ctrt) — EP8 (ExpTotalIndex) 와 동일 패턴
- MarketCode 기본값 "U" (빈 문자열 시 자동 설정)
- TestClient_ExpIndexTrend + TestClient_ExpIndexTrend_InvalidJSON
- fixture: domestic/testdata/exp_index_trend_success.json

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 9: examples/domestic_industry/main.go

Phase 2.7 에서 구현한 7개 NEW 메서드(EP3~EP9)를 시연하는 예제 파일 작성.

> **참고**: EP1 (`InquireIndexPrice`) + EP2 (`InquireIndexCategoryPrice`) 는 Phase 1.4 에서 이미 구현 — 예제에서 제외.

샘플 심볼: `"0001"` (코스피 지수), MarketCode: `"U"` (기본값)

### Step 1 — 예제 파일 작성

**파일**: `examples/domestic_industry/main.go`

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

	symbol := "0001" // KOSPI 지수

	// 1. InquireIndexDailyPrice — 국내업종 일자별지수 (FHPUP02120000)
	dailyResp, err := client.Domestic.InquireIndexDailyPrice(ctx, domestic.InquireIndexDailyPriceParams{
		Symbol:        symbol,
		MarketCode:    "U",
		PeriodDivCode: "D",
		StartDate:     "20260101",
		EndDate:       "20260430",
	})
	if err != nil {
		log.Printf("InquireIndexDailyPrice error: %v", err)
	} else {
		fmt.Printf("[EP3] InquireIndexDailyPrice: BstpNmixPrpr=%s, items=%d\n",
			dailyResp.Output1.BstpNmixPrpr, len(dailyResp.Output2))
	}

	// 2. InquireIndexTimeprice — 국내업종 시간별지수 분 (FHPUP02110200)
	timeResp, err := client.Domestic.InquireIndexTimeprice(ctx, domestic.InquireIndexTimepriceParams{
		Symbol:     symbol,
		MarketCode: "U",
	})
	if err != nil {
		log.Printf("InquireIndexTimeprice error: %v", err)
	} else {
		fmt.Printf("[EP4] InquireIndexTimeprice: items=%d\n", len(timeResp.Output))
	}

	// 3. InquireIndexTickprice — 국내업종 시간별지수 초 (FHPUP02110100)
	tickResp, err := client.Domestic.InquireIndexTickprice(ctx, domestic.InquireIndexTickpriceParams{
		Symbol:     symbol,
		MarketCode: "U",
	})
	if err != nil {
		log.Printf("InquireIndexTickprice error: %v", err)
	} else {
		fmt.Printf("[EP5] InquireIndexTickprice: items=%d\n", len(tickResp.Output))
	}

	// 4. InquireDailyIndexchartprice — 국내주식업종기간별시세 (FHKUP03500100)
	dailyChartResp, err := client.Domestic.InquireDailyIndexchartprice(ctx, domestic.InquireDailyIndexchartpriceParams{
		Symbol:        symbol,
		MarketCode:    "U",
		PeriodDivCode: "D",
		StartDate:     "20260101",
		EndDate:       "20260430",
	})
	if err != nil {
		log.Printf("InquireDailyIndexchartprice error: %v", err)
	} else {
		fmt.Printf("[EP6] InquireDailyIndexchartprice: BstpNmixPrpr=%s, items=%d\n",
			dailyChartResp.Output1.BstpNmixPrpr, len(dailyChartResp.Output2))
	}

	// 5. InquireTimeIndexchartprice — 업종 분봉조회 (FHKUP03500200)
	timeChartResp, err := client.Domestic.InquireTimeIndexchartprice(ctx, domestic.InquireTimeIndexchartpriceParams{
		Symbol:     symbol,
		MarketCode: "U",
	})
	if err != nil {
		log.Printf("InquireTimeIndexchartprice error: %v", err)
	} else {
		fmt.Printf("[EP7] InquireTimeIndexchartprice: BstpNmixPrpr=%s, items=%d\n",
			timeChartResp.Output1.BstpNmixPrpr, len(timeChartResp.Output2))
	}

	// 6. ExpTotalIndex — 예상체결 전체지수 (FHKUP11750000)
	expTotalResp, err := client.Domestic.ExpTotalIndex(ctx, domestic.ExpTotalIndexParams{
		MrktClsCode: "0",
		MarketCode:  "U",
		Symbol:      symbol,
		MkopClsCode: "0",
	})
	if err != nil {
		log.Printf("ExpTotalIndex error: %v", err)
	} else {
		fmt.Printf("[EP8] ExpTotalIndex: output1 items=%d, output2 items=%d\n",
			len(expTotalResp.Output1), len(expTotalResp.Output2))
	}

	// 7. ExpIndexTrend — 예상체결지수 추이 (FHPST01840000)
	expTrendResp, err := client.Domestic.ExpIndexTrend(ctx, domestic.ExpIndexTrendParams{
		MkopClsCode: "1",
		InputHour1:  "10",
		Symbol:      symbol,
		MarketCode:  "U",
	})
	if err != nil {
		log.Printf("ExpIndexTrend error: %v", err)
	} else if len(expTrendResp.Output) > 0 {
		item := expTrendResp.Output[0]
		fmt.Printf("[EP9] ExpIndexTrend: StckCntgHour=%s, BstpNmixPrpr=%s, AcmlVol=%d\n",
			item.StckCntgHour, item.BstpNmixPrpr, item.AcmlVol.Int64())
	}
}
```

### Step 2 — 빌드 확인

```bash
cd /Users/user/src/workspace_moneyflow/korea-investment-stock && go build ./examples/domestic_industry && echo OK
```

출력: `OK`

### Step 3 — 누출 바이너리 정리

```bash
rm -f /Users/user/src/workspace_moneyflow/korea-investment-stock/domestic_industry
```

### Step 4 — gofmt

```bash
gofmt -w examples/domestic_industry/main.go
```

### Step 5 — Commit

```bash
git add examples/domestic_industry/main.go
git commit -m "$(cat <<'EOF'
[feat] examples — domestic_industry (Phase 2.7 7 메서드 시연)

- EP3 InquireIndexDailyPrice (FHPUP02120000)
- EP4 InquireIndexTimeprice (FHPUP02110200)
- EP5 InquireIndexTickprice (FHPUP02110100)
- EP6 InquireDailyIndexchartprice (FHKUP03500100)
- EP7 InquireTimeIndexchartprice (FHKUP03500200)
- EP8 ExpTotalIndex (FHKUP11750000)
- EP9 ExpIndexTrend (FHPST01840000)
- 샘플 심볼 "0001" (KOSPI 지수), MarketCode "U"

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 10: 문서 갱신

4개 파일 업데이트.

### Step 1 — 기존 파일 읽기

```bash
# 현재 상태 확인
head -10 /Users/user/src/workspace_moneyflow/korea-investment-stock/CLAUDE.md
grep -n "Phase 2\." /Users/user/src/workspace_moneyflow/korea-investment-stock/README.md | head -20
grep -n "\[1\." /Users/user/src/workspace_moneyflow/korea-investment-stock/CHANGELOG.md | head -10
grep -n "Phase 2\." /Users/user/src/workspace_moneyflow/korea-investment-stock/domestic/doc.go
```

### Step 2 — 파일 편집

#### 2-a. CLAUDE.md

Replace banner:

```
> **Phase 2.6 — 해외 정보 4 메서드 (v1.9.0).** Phase 2.5+ design spec 및 plan 참고.
```

→

```
> **Phase 2.7 — 업종/지수 7 메서드 (v1.10.0). Phase 2.5+ 완료.**
```

Phase 2.6 plan link 줄 아래에 Phase 2.7 plan link bullet 추가:

```markdown
- Phase 2.7 implementation plan: [`docs/superpowers/specs/2026-05-05-phase2-7-industry-implementation-plan.md`](docs/superpowers/specs/2026-05-05-phase2-7-industry-implementation-plan.md)
```

#### 2-b. README.md

Heading 변경:

```
Available Methods (Phase 1.2 ~ 2.6)
```

→

```
Available Methods (Phase 1.2 ~ 2.7)
```

메서드 표 맨 끝에 7개 행 추가 (기존 행 아래):

```markdown
| `Domestic.InquireIndexDailyPrice` | `quotations/inquire-index-daily-price` | FHPUP02120000 |
| `Domestic.InquireIndexTimeprice` | `quotations/inquire-index-timeprice` | FHPUP02110200 |
| `Domestic.InquireIndexTickprice` | `quotations/inquire-index-tickprice` | FHPUP02110100 |
| `Domestic.InquireDailyIndexchartprice` | `quotations/inquire-daily-indexchartprice` | FHKUP03500100 |
| `Domestic.InquireTimeIndexchartprice` | `quotations/inquire-time-indexchartprice` | FHKUP03500200 |
| `Domestic.ExpTotalIndex` | `quotations/exp-total-index` | FHKUP11750000 |
| `Domestic.ExpIndexTrend` | `quotations/exp-index-trend` | FHPST01840000 |
```

메서드 총 수 카운트: `64` → `71`

#### 2-c. CHANGELOG.md

`## [1.9.0]` 바로 위에 아래 섹션 삽입:

```markdown
## [1.10.0] - 2026-05-05

### Added — Phase 2.7 (업종/지수)

- `Domestic.InquireIndexDailyPrice` — 국내업종 일자별지수 (FHPUP02120000) — output1 20 + output2 13 fields
- `Domestic.InquireIndexTimeprice` — 국내업종 시간별지수 분 (FHPUP02110200) — output 8 fields, bsop_hour timestamp
- `Domestic.InquireIndexTickprice` — 국내업종 시간별지수 초 (FHPUP02110100) — output 8 fields, stck_cntg_hour timestamp
- `Domestic.InquireDailyIndexchartprice` — 국내주식업종기간별시세 (FHKUP03500100) — output1 15 + output2 8, futs_prdy_* embedded
- `Domestic.InquireTimeIndexchartprice` — 업종 분봉조회 (FHKUP03500200) — output1 16 + output2 8
- `Domestic.ExpTotalIndex` — 예상체결 전체지수 (FHKUP11750000) — output1 9 + output2 10, LOWERCASE fid_* query params
- `Domestic.ExpIndexTrend` — 예상체결지수 추이 (FHPST01840000) — output 7 fields
- examples: `domestic_industry`

### Notes

- Phase 2.5+ design spec §Phase 2.7 listed 9 methods. EP1 (`InquireIndexPrice`) + EP2 (`InquireIndexCategoryPrice`) 는 Phase 1.4 에서 이미 구현됨 — Phase 2.7 = 7 NEW methods.
- EP8 (`ExpTotalIndex`) 의 query param wire keys 는 lowercase (`fid_*`) — 다른 endpoint 의 `FID_*` 와 다름. 코드에서 lowercase 그대로 보존.
- EP8/EP9 응답 struct 는 `prdy_ctrt` (short form, NOT `bstp_nmix_prdy_ctrt`) 사용.
- EP9 (`ExpIndexTrend`) KIS docs 의 Korean field labels 가 scrambled 되어있음 (e.g., `stck_cntg_hour` 가 "주식 단축 종목코드" 로 잘못 라벨링). Field 명은 정확 — 라벨만 무시.
- **Phase 2.5+ 완료** (2.5 + 2.6 + 2.7 = 18 NEW methods 누적). Phase 2 + Phase 2.5+ = 43 read-only 확장.
- 누적 71 메서드 (Phase 1: 28 + Phase 2: 25 + Phase 2.5+: 18).
```

#### 2-d. domestic/doc.go

기존 Phase 2.5 (또는 최신) 섹션 이후에 Phase 2.7 섹션 추가:

```go
// Phase 2.7 — 업종/지수 (v1.10.0)  [Phase 2.5+ 마지막 sub-phase]
//
//   EP3  InquireIndexDailyPrice      — 국내업종 일자별지수       FHPUP02120000
//   EP4  InquireIndexTimeprice       — 국내업종 시간별지수 분    FHPUP02110200
//   EP5  InquireIndexTickprice       — 국내업종 시간별지수 초    FHPUP02110100
//   EP6  InquireDailyIndexchartprice — 국내주식업종기간별시세    FHKUP03500100
//   EP7  InquireTimeIndexchartprice  — 업종 분봉조회             FHKUP03500200
//   EP8  ExpTotalIndex               — 예상체결 전체지수         FHKUP11750000
//   EP9  ExpIndexTrend               — 예상체결지수 추이         FHPST01840000
//
// Anomalies:
//   EP1+EP2 already in Phase 1.4 → Phase 2.7 = 7 NEW (not 9)
//   EP8 lowercase fid_* query params (KIS 유일 예외)
//   EP8/EP9 prdy_ctrt short form (NOT bstp_nmix_prdy_ctrt)
//   EP9 KIS docs Korean labels scrambled — field names are correct
```

### Step 3 — 빌드/벳/포맷 검증

```bash
cd /Users/user/src/workspace_moneyflow/korea-investment-stock && go build ./... && go vet ./... && gofmt -l .
```

출력: 모두 silent (변경 없음).

### Step 4 — 누출 바이너리 확인 및 정리

```bash
ls /Users/user/src/workspace_moneyflow/korea-investment-stock/ | grep -v '/'
```

`.go`, `.mod`, `.sum`, `Makefile`, `README.md`, `CHANGELOG.md` 등 소스 파일만 존재해야 함. 바이너리 발견 시 즉시 `rm -f <binary>`.

### Step 5 — Commit

```bash
git add CLAUDE.md README.md CHANGELOG.md domestic/doc.go
git commit -m "$(cat <<'EOF'
[docs] Phase 2.7 문서 갱신 (v1.10.0, Phase 2.5+ 완료)

- CLAUDE.md: banner → Phase 2.7, Phase 2.7 plan link 추가
- README.md: heading Phase 1.2~2.7, 메서드 64→71, 7개 신규 행
- CHANGELOG.md: [1.10.0] 섹션 추가 (EP1+EP2 Phase 1.4 구현 주석 포함)
- domestic/doc.go: Phase 2.7 섹션 추가

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 11: 최종 점검

### Step 1 — gofmt

```bash
cd /Users/user/src/workspace_moneyflow/korea-investment-stock && gofmt -l . | head
```

출력: silent (변경 없음).

### Step 2 — Build + Vet

```bash
go build ./... && go vet ./...
```

출력: silent.

### Step 3 — Full Test Suite (race detector)

```bash
go test ./... -race -count=1
```

출력: 모든 패키지 `ok`.

### Step 4 — Coverage

```bash
# domestic 패키지
go test ./domestic/... -coverprofile=/tmp/cov_d.out -covermode=atomic
go tool cover -func=/tmp/cov_d.out | tail -2

# root kis 패키지
go test . -coverprofile=/tmp/cov_r.out -covermode=atomic
go tool cover -func=/tmp/cov_r.out | tail -2
```

기대값: domestic ≥ 80%, root ≥ 80%.

> **커버리지 < 80% 시**: 각 신규 메서드에 InvalidJSON 에러 경로 테스트 추가 (Phase 2.6 레슨 — `TestClient_*_InvalidJSON` 패턴).

### Step 5 — 파일 수 확인

```bash
ls \
  domestic/testdata/index_daily_price_success.json \
  domestic/testdata/index_timeprice_success.json \
  domestic/testdata/index_tickprice_success.json \
  domestic/testdata/daily_indexchartprice_success.json \
  domestic/testdata/time_indexchartprice_success.json \
  domestic/testdata/exp_total_index_success.json \
  domestic/testdata/exp_index_trend_success.json \
  examples/domestic_industry/main.go \
  2>&1 | wc -l
```

기대값: `8`

### Step 6 — 커밋 수 확인

```bash
git log main..HEAD --oneline | wc -l
```

기대값: Task 1~10 완료 시 10개 내외 (각 Task 당 1 commit).

### 최종 점검 통과 기준

| 항목 | 기준 |
|---|---|
| gofmt | silent |
| go build/vet | silent |
| go test -race | all PASS |
| domestic coverage | ≥ 80% |
| root coverage | ≥ 80% |
| testdata fixtures | 7개 모두 존재 |
| example file | 1개 존재 |

---

## Task 12: PR 생성 (사용자 승인 후)

> **Claude 는 push / PR 생성을 사용자 명시적 승인 후에만 실행 (글로벌 정책).**

### Step 1 — 사용자 승인 요청

작업 진행 보고 (Task 1~11 완료 요약) + PR 생성 가능 여부 confirm.

보고 내용 포함 사항:
- 구현된 메서드 수 (7 NEW)
- 테스트 결과 (all PASS, coverage %)
- 커밋 수
- 누적 메서드 수 (71)

### Step 2 — Push

```bash
git push -u origin feat/phase2-7-industry
```

### Step 3 — PR 생성

```bash
gh pr create \
  --title "Phase 2.7 — 업종/지수 (v1.10.0) [Phase 2.5+ 완료]" \
  --reviewer kenshin579 \
  --base main \
  --head feat/phase2-7-industry \
  --body "$(cat <<'EOF'
## Summary

- Phase 2.7 implementation (7 NEW methods) — **Phase 2.5+ 의 마지막 sub-phase**
- Phase 2 패턴 그대로 (Style A, Params struct, 일반 타입 매핑)
- v1.10.0 release (누적 64 → 71)
- Phase 2.5+ 전체 완료 (2.5: 7 + 2.6: 4 + 2.7: 7 = 18 NEW methods)

## 메서드 → 한투 API 매핑 (7 NEW)

| Go 메서드 | path | TR_ID |
|---|---|---|
| InquireIndexDailyPrice | quotations/inquire-index-daily-price | FHPUP02120000 |
| InquireIndexTimeprice | quotations/inquire-index-timeprice | FHPUP02110200 |
| InquireIndexTickprice | quotations/inquire-index-tickprice | FHPUP02110100 |
| InquireDailyIndexchartprice | quotations/inquire-daily-indexchartprice | FHKUP03500100 |
| InquireTimeIndexchartprice | quotations/inquire-time-indexchartprice | FHKUP03500200 |
| ExpTotalIndex | quotations/exp-total-index | FHKUP11750000 |
| ExpIndexTrend | quotations/exp-index-trend | FHPST01840000 |

## Anomalies

- Phase 2.5+ design spec listed 9 methods; EP1 + EP2 already implemented in Phase 1.4 — Phase 2.7 = 7 NEW.
- EP8 (ExpTotalIndex) lowercase fid_* query params (vs FID_* elsewhere)
- EP8/EP9 use prdy_ctrt short form (NOT bstp_nmix_prdy_ctrt)
- EP9 KIS docs Korean labels scrambled — field names correct, labels ignored
- EP6/EP7 embed futs_prdy_oprc/hgpr/lwpr (선물 전일 OHLC)

## Test Plan
- [x] go build/vet/fmt clean
- [x] go test ./... -race -count=1 모든 패키지 PASS
- [x] Coverage domestic ≥80%, root kis ≥80%
- [x] httpmock 단위 테스트 (7 methods + InvalidJSON tests if needed)
- [x] examples/domestic_industry build OK

## Breaking Changes
없음.

## Notes
- **Phase 2.5+ 완료** — 18 NEW read-only methods 추가.
- 다음: Trading/WebSocket/v1.x 유지보수 결정.

## 참고 문서
- Phase 2.5+ design spec: docs/superpowers/specs/2026-05-05-phase2-5plus-extension-design.md
- Phase 2.7 plan: docs/superpowers/specs/2026-05-05-phase2-7-industry-implementation-plan.md

🤖 Generated with [Claude Code](https://claude.com/claude-code)
EOF
)"
```

### Step 4 — Merge (사용자 승인 후)

```bash
gh pr merge <PR#> --merge
```

### Step 5 — 후속 작업

```bash
# 태그 생성
git tag -a v1.10.0 -m "Phase 2.7 — 업종/지수 7 메서드 (Phase 2.5+ 완료)"
git push origin v1.10.0

# GitHub Release 생성
gh release create v1.10.0 \
  --title "v1.10.0 — Phase 2.7 업종/지수 (Phase 2.5+ 완료)" \
  --notes-file <(awk '/^## \[1\.10\.0\]/{f=1;next} /^## \[/&&f{exit} f' CHANGELOG.md)
```
