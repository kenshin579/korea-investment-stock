# Phase 2.5 — 투자자/매매 동향 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** korea-investment-stock Go 라이브러리에 투자자/프로그램매매 동향 7 메서드 추가 (`v1.8.0` release). 외인기관 가집계 2 (investor.go append) + 프로그램매매 5 (program_trade.go 신규).

**Architecture:** Phase 1/2 인프라 + 패턴 재사용. `domestic/investor.go` 에 2 메서드 append, `domestic/program_trade.go` 신규 파일에 5 메서드. Phase 2 standard type mapping (decimal/int64/float64/string). TDD: testdata fixture → 실패 테스트 → struct + 메서드 구현 → 통과 → commit.

**Tech Stack:** Go 1.25+, `github.com/jarcoal/httpmock`, `github.com/stretchr/testify`, `github.com/shopspring/decimal`. 새 dependency 없음.

**참고 spec:**
- Phase 2.5+ design spec: `docs/superpowers/specs/2026-05-05-phase2-5plus-extension-design.md` (§Phase 2.5)
- Phase 2.4 plan (참조 패턴): `docs/superpowers/specs/2026-05-05-phase2-4-ksd-implementation-plan.md`
- 한투 API docs: `docs/api/국내주식/주식현재가_투자자.md`, `국내기관_외국인_매매종목가집계.md`, `종목별_프로그램매매추이(일별).md`, `종목별_프로그램매매추이(체결).md`, `프로그램매매_종합현황(시간).md`, `프로그램매매_종합현황(일별).md`, `프로그램매매_투자자매매동향(당일).md`

---

## 사전 정보

| 항목 | 값 |
|------|---|
| 작업 브랜치 | `docs/phase2-5plus-spec` |
| 시작 HEAD | Phase 2.4 구현 완료 commit (v1.7.0) |
| Release 목표 | `v1.8.0` |
| PR 베이스 | `main` |
| 현재 main 상태 | v1.7.0 publish 완료 (Phase 2.4 통합, 53 메서드) |

---

## 메서드 매핑

| Go 메서드 | path (last segment) | TR_ID | output key | fields | anomalies |
|---|---|---|---|---|---|
| `InquireInvestorTrendEstimate` | `investor-trend-estimate` | HHPTJ04160200 | `output2 []` | 4 | non-FID query (`MKSC_SHRN_ISCD`), output2 only (no output1) |
| `InquireForeignInstitutionTotal` | `foreign-institution-total` | FHPTJ04400000 | `Output []` | 26 | json:"Output" (capital O), default MKT="V"/SCR="16449" |
| `InquireProgramTradeByStockDaily` | `program-trade-by-stock-daily` | FHPPG04650201 | `output []` | 15 | "002" date prefix doc'd in Params comment, last field `_icdc2` |
| `InquireProgramTradeByStock` | `program-trade-by-stock` | FHPPG04650101 | `output []` | 14 | last field `_icdc` (no trailing "2") |
| `InquireCompProgramTradeToday` | `comp-program-trade-today` | FHPPG04600101 | `output1 []` | 18 | 6 params (4 blank), `shun` typo preserved verbatim |
| `InquireCompProgramTradeDaily` | `comp-program-trade-daily` | FHPPG04600001 | `output []` | 97 | 8-month limit, multiple `shun`/`smtm` typos preserved |
| `InquireInvestorProgramTradeToday` | `investor-program-trade-today` | HHPPG046600C1 | `output1 []` | 20 | non-FID query (EXCH_DIV_CLS_CODE/MRKT_DIV_CLS_CODE), MRKT="1"/"4" not K/Q, `_amt` suffix |

---

## 파일 구조

### 신규
- `domestic/program_trade.go` — 5 메서드 + structs + Params (EP3~EP7)
- `domestic/program_trade_test.go` — 5 테스트 함수
- `domestic/testdata/investor_trend_estimate_success.json` — EP1 fixture
- `domestic/testdata/foreign_institution_total_success.json` — EP2 fixture
- `domestic/testdata/program_trade_by_stock_daily_success.json` — EP3 fixture
- `domestic/testdata/program_trade_by_stock_success.json` — EP4 fixture
- `domestic/testdata/comp_program_trade_today_success.json` — EP5 fixture
- `domestic/testdata/comp_program_trade_daily_success.json` — EP6 fixture
- `domestic/testdata/investor_program_trade_today_success.json` — EP7 fixture
- `examples/domestic_program_trade/main.go`

### 수정
- `domestic/investor.go` — EP1 (`InquireInvestorTrendEstimate`) + EP2 (`InquireForeignInstitutionTotal`) append
- `domestic/investor_test.go` — EP1 + EP2 테스트 함수 append
- `CLAUDE.md` — banner Phase 2.4 → Phase 2.5, plan link 추가
- `README.md` — Available Methods 표 갱신 (53 → 60 메서드)
- `CHANGELOG.md` — `[1.8.0]` entry
- `domestic/doc.go` — Phase 2.5 section 추가

---

## 타입 매핑

Phase 2 standard mapping (KSD all-string 과 다름):

| 필드 패턴 | Go 타입 | 예시 |
|---|---|---|
| 가격 (`stck_prpr`, `stck_clpr`, `prdy_vrss`, `bstp_nmix_prpr`, `bstp_nmix_prdy_vrss`) | `decimal.Decimal` (bare, no `,string`) | `json:"stck_prpr"` |
| 수량/거래량/거래대금 (`*_vol`, `*_qty`, `*_tr_pbmn`, `*_amt`, `acml_*`, `*_icdc`, `*_icdc2`) | `int64` + `,string` tag | `json:"acml_vol,string"` |
| 비율 (`prdy_ctrt`, `*_rate`) | `float64` + `,string` tag | `json:"prdy_ctrt,string"` |
| 코드/날짜/명칭/Y-N (`prdy_vrss_sign`, `bsop_hour*`, `stck_bsop_date`, `mksc_shrn_iscd`, `hts_kor_isnm`, `invr_cls_*`) | `string` (plain) | `json:"bsop_hour"` |

---

## Tasks (12 total)

Task 1: testdata fixtures 7개 | Task 2: investor.go append + EP1 | Task 3: EP2 ForeignInstitutionTotal | Task 4: program_trade.go base + EP3 | Task 5: EP4 ProgramTradeByStock | Task 6: EP5 CompProgramTradeToday | Task 7: EP6 CompProgramTradeDaily | Task 8: EP7 InvestorProgramTradeToday | Task 9: examples/domestic_program_trade/main.go | Task 10: CHANGELOG + README + CLAUDE.md 갱신 | Task 11: domestic/doc.go Phase 2.5 section | Task 12: 최종 검증 + PR

---

## Task 1: testdata fixtures (7 합성 JSON)

- [ ] Step 1: `domestic/testdata/investor_trend_estimate_success.json` (EP1, output2, 4 fields)
- [ ] Step 2: `domestic/testdata/foreign_institution_total_success.json` (EP2, Output capital-O, 26 fields)
- [ ] Step 3: `domestic/testdata/program_trade_by_stock_daily_success.json` (EP3, output, 15 fields)
- [ ] Step 4: `domestic/testdata/program_trade_by_stock_success.json` (EP4, output, 14 fields)
- [ ] Step 5: `domestic/testdata/comp_program_trade_today_success.json` (EP5, output1, 18 fields)
- [ ] Step 6: `domestic/testdata/comp_program_trade_daily_success.json` (EP6, output, 97 fields)
- [ ] Step 7: `domestic/testdata/investor_program_trade_today_success.json` (EP7, output1, 20 fields)
- [ ] Step 8: validation

```bash
for f in \
  domestic/testdata/investor_trend_estimate_success.json \
  domestic/testdata/foreign_institution_total_success.json \
  domestic/testdata/program_trade_by_stock_daily_success.json \
  domestic/testdata/program_trade_by_stock_success.json \
  domestic/testdata/comp_program_trade_today_success.json \
  domestic/testdata/comp_program_trade_daily_success.json \
  domestic/testdata/investor_program_trade_today_success.json; do
  python3 -c "import json; json.load(open('$f'))" && echo "$f OK"
done
# Expected: 7 OK lines
```

- [ ] Step 9: commit

```bash
git commit -m "$(cat <<'EOF'
[chore] testdata — 7 investor/program_trade fixture JSON (Phase 2.5)

합성 JSON fixtures (2 records each):
- investor_trend_estimate_success.json (4 fields, output2)
- foreign_institution_total_success.json (26 fields, Output capital-O)
- program_trade_by_stock_daily_success.json (15 fields, icdc2)
- program_trade_by_stock_success.json (14 fields, icdc no-2)
- comp_program_trade_today_success.json (18 fields, shun typo)
- comp_program_trade_daily_success.json (97 fields, 1 record)
- investor_program_trade_today_success.json (20 fields, _amt suffix)

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

### Fixture content

**Step 1 — `investor_trend_estimate_success.json`** (EP1, output2 only, 4 fields)
```json
{
  "output2": [
    {
      "bsop_hour_gb": "1",
      "frgn_fake_ntby_qty": "123456",
      "orgn_fake_ntby_qty": "-45678",
      "sum_fake_ntby_qty": "77778"
    },
    {
      "bsop_hour_gb": "2",
      "frgn_fake_ntby_qty": "234567",
      "orgn_fake_ntby_qty": "-56789",
      "sum_fake_ntby_qty": "177778"
    }
  ]
}
```

**Step 2 — `foreign_institution_total_success.json`** (EP2, `Output` capital O, 26 fields)
```json
{
  "Output": [
    {
      "hts_kor_isnm": "삼성전자",
      "mksc_shrn_iscd": "005930",
      "ntby_qty": "123456",
      "stck_prpr": "75800",
      "prdy_vrss_sign": "2",
      "prdy_vrss": "800",
      "prdy_ctrt": "1.07",
      "acml_vol": "12345678",
      "frgn_ntby_qty": "100000",
      "orgn_ntby_qty": "-50000",
      "ivtr_ntby_qty": "-20000",
      "bank_ntby_qty": "5000",
      "insu_ntby_qty": "3000",
      "mrbn_ntby_qty": "1000",
      "fund_ntby_qty": "-10000",
      "etc_orgt_ntby_vol": "2000",
      "etc_corp_ntby_vol": "500",
      "frgn_ntby_tr_pbmn": "7580000000",
      "orgn_ntby_tr_pbmn": "-3790000000",
      "ivtr_ntby_tr_pbmn": "-1516000000",
      "bank_ntby_tr_pbmn": "379000000",
      "insu_ntby_tr_pbmn": "227400000",
      "mrbn_ntby_tr_pbmn": "75800000",
      "fund_ntby_tr_pbmn": "-758000000",
      "etc_orgt_ntby_tr_pbmn": "151600000",
      "etc_corp_ntby_tr_pbmn": "37900000"
    },
    {
      "hts_kor_isnm": "SK하이닉스",
      "mksc_shrn_iscd": "000660",
      "ntby_qty": "56789",
      "stck_prpr": "185000",
      "prdy_vrss_sign": "2",
      "prdy_vrss": "2000",
      "prdy_ctrt": "1.09",
      "acml_vol": "3456789",
      "frgn_ntby_qty": "60000",
      "orgn_ntby_qty": "-20000",
      "ivtr_ntby_qty": "-8000",
      "bank_ntby_qty": "2000",
      "insu_ntby_qty": "1000",
      "mrbn_ntby_qty": "500",
      "fund_ntby_qty": "-3000",
      "etc_orgt_ntby_vol": "800",
      "etc_corp_ntby_vol": "200",
      "frgn_ntby_tr_pbmn": "11100000000",
      "orgn_ntby_tr_pbmn": "-3700000000",
      "ivtr_ntby_tr_pbmn": "-1480000000",
      "bank_ntby_tr_pbmn": "370000000",
      "insu_ntby_tr_pbmn": "185000000",
      "mrbn_ntby_tr_pbmn": "92500000",
      "fund_ntby_tr_pbmn": "-555000000",
      "etc_orgt_ntby_tr_pbmn": "148000000",
      "etc_corp_ntby_tr_pbmn": "37000000"
    }
  ]
}
```

**Step 3 — `program_trade_by_stock_daily_success.json`** (EP3, output, 15 fields, last `icdc2`)
```json
{
  "output": [
    {
      "stck_bsop_date": "20260505",
      "stck_clpr": "75800",
      "prdy_vrss": "800",
      "prdy_vrss_sign": "2",
      "prdy_ctrt": "1.07",
      "acml_vol": "12345678",
      "acml_tr_pbmn": "936203964",
      "whol_smtn_seln_vol": "500000",
      "whol_smtn_shnu_vol": "650000",
      "whol_smtn_ntby_qty": "150000",
      "whol_smtn_seln_tr_pbmn": "37900000000",
      "whol_smtn_shnu_tr_pbmn": "49270000000",
      "whol_smtn_ntby_tr_pbmn": "11370000000",
      "whol_ntby_vol_icdc": "10000",
      "whol_ntby_tr_pbmn_icdc2": "500000000"
    },
    {
      "stck_bsop_date": "20260504",
      "stck_clpr": "75000",
      "prdy_vrss": "-500",
      "prdy_vrss_sign": "5",
      "prdy_ctrt": "-0.66",
      "acml_vol": "9876543",
      "acml_tr_pbmn": "740740725",
      "whol_smtn_seln_vol": "480000",
      "whol_smtn_shnu_vol": "430000",
      "whol_smtn_ntby_qty": "-50000",
      "whol_smtn_seln_tr_pbmn": "36000000000",
      "whol_smtn_shnu_tr_pbmn": "32250000000",
      "whol_smtn_ntby_tr_pbmn": "-3750000000",
      "whol_ntby_vol_icdc": "-8000",
      "whol_ntby_tr_pbmn_icdc2": "-300000000"
    }
  ]
}
```

**Step 4 — `program_trade_by_stock_success.json`** (EP4, output, 14 fields, last `icdc` no-"2")
```json
{
  "output": [
    {
      "bsop_hour": "090100",
      "stck_prpr": "75800",
      "prdy_vrss": "800",
      "prdy_vrss_sign": "2",
      "prdy_ctrt": "1.07",
      "acml_vol": "123456",
      "whol_smtn_seln_vol": "5000",
      "whol_smtn_shnu_vol": "6500",
      "whol_smtn_ntby_qty": "1500",
      "whol_smtn_seln_tr_pbmn": "379000000",
      "whol_smtn_shnu_tr_pbmn": "492700000",
      "whol_smtn_ntby_tr_pbmn": "113700000",
      "whol_ntby_vol_icdc": "200",
      "whol_ntby_tr_pbmn_icdc": "15160000"
    },
    {
      "bsop_hour": "090200",
      "stck_prpr": "75900",
      "prdy_vrss": "900",
      "prdy_vrss_sign": "2",
      "prdy_ctrt": "1.20",
      "acml_vol": "234567",
      "whol_smtn_seln_vol": "8000",
      "whol_smtn_shnu_vol": "10000",
      "whol_smtn_ntby_qty": "2000",
      "whol_smtn_seln_tr_pbmn": "607200000",
      "whol_smtn_shnu_tr_pbmn": "759000000",
      "whol_smtn_ntby_tr_pbmn": "151800000",
      "whol_ntby_vol_icdc": "500",
      "whol_ntby_tr_pbmn_icdc": "37950000"
    }
  ]
}
```

**Step 5 — `comp_program_trade_today_success.json`** (EP5, output1, 18 fields, `shun`/`smtm` typos preserved)
```json
{
  "output1": [
    {
      "bsop_hour": "090100",
      "arbt_smtn_seln_tr_pbmn": "12000000000",
      "arbt_smtm_seln_tr_pbmn_rate": "35.50",
      "arbt_smtn_shnu_tr_pbmn": "13000000000",
      "arbt_smtm_shun_tr_pbmn_rate": "36.10",
      "nabt_smtn_seln_tr_pbmn": "21800000000",
      "nabt_smtm_seln_tr_pbmn_rate": "64.50",
      "nabt_smtn_shnu_tr_pbmn": "23000000000",
      "nabt_smtm_shun_tr_pbmn_rate": "63.90",
      "arbt_smtn_ntby_tr_pbmn": "1000000000",
      "arbt_smtm_ntby_tr_pbmn_rate": "33.33",
      "nabt_smtn_ntby_tr_pbmn": "1200000000",
      "nabt_smtm_ntby_tr_pbmn_rate": "40.00",
      "whol_smtn_ntby_tr_pbmn": "2200000000",
      "whol_ntby_tr_pbmn_rate": "100.00",
      "bstp_nmix_prpr": "2750.50",
      "bstp_nmix_prdy_vrss": "12.30",
      "prdy_vrss_sign": "2"
    },
    {
      "bsop_hour": "090200",
      "arbt_smtn_seln_tr_pbmn": "13500000000",
      "arbt_smtm_seln_tr_pbmn_rate": "35.80",
      "arbt_smtn_shnu_tr_pbmn": "14200000000",
      "arbt_smtm_shun_tr_pbmn_rate": "36.40",
      "nabt_smtn_seln_tr_pbmn": "24200000000",
      "nabt_smtm_seln_tr_pbmn_rate": "64.20",
      "nabt_smtn_shnu_tr_pbmn": "24800000000",
      "nabt_smtm_shun_tr_pbmn_rate": "63.60",
      "arbt_smtn_ntby_tr_pbmn": "700000000",
      "arbt_smtm_ntby_tr_pbmn_rate": "30.43",
      "nabt_smtn_ntby_tr_pbmn": "600000000",
      "nabt_smtm_ntby_tr_pbmn_rate": "26.09",
      "whol_smtn_ntby_tr_pbmn": "1300000000",
      "whol_ntby_tr_pbmn_rate": "100.00",
      "bstp_nmix_prpr": "2752.80",
      "bstp_nmix_prdy_vrss": "14.60",
      "prdy_vrss_sign": "2"
    }
  ]
}
```

**Step 6 — `comp_program_trade_daily_success.json`** (EP6, output, 97 fields, 1 record)

Note: EP6 has 97 fields. Below is one record with all fields populated with realistic non-zero values. Field names preserve all `shun`/`smtm` typos verbatim as they appear in the KIS API.

```json
{
  "output": [
    {
      "stck_bsop_date": "20260505",
      "bstp_nmix_prpr": "2750.50",
      "bstp_nmix_prdy_vrss": "12.30",
      "prdy_vrss_sign": "2",
      "bstp_nmix_prdy_ctrt": "0.45",
      "whol_smtn_seln_tr_pbmn": "33800000000",
      "whol_smtm_seln_tr_pbmn_rate": "100.00",
      "whol_smtn_shnu_tr_pbmn": "36200000000",
      "whol_smtm_shun_tr_pbmn_rate": "100.00",
      "whol_smtn_ntby_tr_pbmn": "2400000000",
      "whol_smtm_ntby_tr_pbmn_rate": "100.00",
      "arbt_smtn_seln_tr_pbmn": "12000000000",
      "arbt_smtm_seln_tr_pbmn_rate": "35.50",
      "arbt_smtn_shnu_tr_pbmn": "13000000000",
      "arbt_smtm_shun_tr_pbmn_rate": "35.91",
      "arbt_smtn_ntby_tr_pbmn": "1000000000",
      "arbt_smtm_ntby_tr_pbmn_rate": "41.67",
      "nabt_smtn_seln_tr_pbmn": "21800000000",
      "nabt_smtm_seln_tr_pbmn_rate": "64.50",
      "nabt_smtn_shnu_tr_pbmn": "23200000000",
      "nabt_smtm_shun_tr_pbmn_rate": "64.09",
      "nabt_smtn_ntby_tr_pbmn": "1400000000",
      "nabt_smtm_ntby_tr_pbmn_rate": "58.33",
      "whol_smtn_seln_vol": "450000",
      "whol_smtm_seln_vol_rate": "100.00",
      "whol_smtn_shnu_vol": "480000",
      "whol_smtm_shun_vol_rate": "100.00",
      "whol_smtn_ntby_qty": "30000",
      "whol_smtm_ntby_qty_rate": "100.00",
      "arbt_smtn_seln_vol": "160000",
      "arbt_smtm_seln_vol_rate": "35.56",
      "arbt_smtn_shnu_vol": "172000",
      "arbt_smtm_shun_vol_rate": "35.83",
      "arbt_smtn_ntby_qty": "12000",
      "arbt_smtm_ntby_qty_rate": "40.00",
      "nabt_smtn_seln_vol": "290000",
      "nabt_smtm_seln_vol_rate": "64.44",
      "nabt_smtn_shnu_vol": "308000",
      "nabt_smtm_shun_vol_rate": "64.17",
      "nabt_smtn_ntby_qty": "18000",
      "nabt_smtm_ntby_qty_rate": "60.00",
      "frgn_smtn_seln_tr_pbmn": "8000000000",
      "frgn_smtm_seln_tr_pbmn_rate": "23.67",
      "frgn_smtn_shnu_tr_pbmn": "9000000000",
      "frgn_smtm_shun_tr_pbmn_rate": "24.86",
      "frgn_smtn_ntby_tr_pbmn": "1000000000",
      "frgn_smtm_ntby_tr_pbmn_rate": "41.67",
      "orgn_smtn_seln_tr_pbmn": "7000000000",
      "orgn_smtm_seln_tr_pbmn_rate": "20.71",
      "orgn_smtn_shnu_tr_pbmn": "7500000000",
      "orgn_smtm_shun_tr_pbmn_rate": "20.72",
      "orgn_smtn_ntby_tr_pbmn": "500000000",
      "orgn_smtm_ntby_tr_pbmn_rate": "20.83",
      "ivtr_smtn_seln_tr_pbmn": "3000000000",
      "ivtr_smtm_seln_tr_pbmn_rate": "8.88",
      "ivtr_smtn_shnu_tr_pbmn": "3200000000",
      "ivtr_smtm_shun_tr_pbmn_rate": "8.84",
      "ivtr_smtn_ntby_tr_pbmn": "200000000",
      "ivtr_smtm_ntby_tr_pbmn_rate": "8.33",
      "bank_smtn_seln_tr_pbmn": "1500000000",
      "bank_smtm_seln_tr_pbmn_rate": "4.44",
      "bank_smtn_shnu_tr_pbmn": "1600000000",
      "bank_smtm_shun_tr_pbmn_rate": "4.42",
      "bank_smtn_ntby_tr_pbmn": "100000000",
      "bank_smtm_ntby_tr_pbmn_rate": "4.17",
      "insu_smtn_seln_tr_pbmn": "1200000000",
      "insu_smtm_seln_tr_pbmn_rate": "3.55",
      "insu_smtn_shnu_tr_pbmn": "1300000000",
      "insu_smtm_shun_tr_pbmn_rate": "3.59",
      "insu_smtn_ntby_tr_pbmn": "100000000",
      "insu_smtm_ntby_tr_pbmn_rate": "4.17",
      "mrbn_smtn_seln_tr_pbmn": "400000000",
      "mrbn_smtm_seln_tr_pbmn_rate": "1.18",
      "mrbn_smtn_shnu_tr_pbmn": "420000000",
      "mrbn_smtm_shun_tr_pbmn_rate": "1.16",
      "mrbn_smtn_ntby_tr_pbmn": "20000000",
      "mrbn_smtm_ntby_tr_pbmn_rate": "0.83",
      "fund_smtn_seln_tr_pbmn": "900000000",
      "fund_smtm_seln_tr_pbmn_rate": "2.66",
      "fund_smtn_shnu_tr_pbmn": "980000000",
      "fund_smtm_shun_tr_pbmn_rate": "2.71",
      "fund_smtn_ntby_tr_pbmn": "80000000",
      "fund_smtm_ntby_tr_pbmn_rate": "3.33",
      "etc_smtn_seln_tr_pbmn": "500000000",
      "etc_smtm_seln_tr_pbmn_rate": "1.48",
      "etc_smtn_shnu_tr_pbmn": "600000000",
      "etc_smtm_shun_tr_pbmn_rate": "1.66",
      "etc_smtn_ntby_tr_pbmn": "100000000",
      "etc_smtm_ntby_tr_pbmn_rate": "4.17",
      "prsn_smtn_seln_tr_pbmn": "10000000000",
      "prsn_smtm_seln_tr_pbmn_rate": "29.59",
      "prsn_smtn_shnu_tr_pbmn": "9800000000",
      "prsn_smtm_shun_tr_pbmn_rate": "27.07",
      "prsn_smtn_ntby_tr_pbmn": "-200000000",
      "prsn_smtm_ntby_tr_pbmn_rate": "-8.33",
      "whol_ntby_tr_pbmn_icdc": "100000000",
      "whol_ntby_vol_icdc": "5000"
    }
  ]
}
```

**Step 7 — `investor_program_trade_today_success.json`** (EP7, output1, 20 fields, `_amt` suffix, MRKT="1"/"4")
```json
{
  "output1": [
    {
      "bsop_hour": "090100",
      "invr_cls_code": "1",
      "invr_cls_name": "외국인",
      "arbt_seln_tr_pbmn": "1200000000",
      "arbt_shnu_tr_pbmn": "1300000000",
      "arbt_ntby_tr_pbmn": "100000000",
      "arbt_seln_vol_amt": "15000",
      "arbt_shnu_vol_amt": "16000",
      "arbt_ntby_qty_amt": "1000",
      "nabt_seln_tr_pbmn": "2200000000",
      "nabt_shnu_tr_pbmn": "2400000000",
      "nabt_ntby_tr_pbmn": "200000000",
      "nabt_seln_vol_amt": "28000",
      "nabt_shnu_vol_amt": "30000",
      "nabt_ntby_qty_amt": "2000",
      "whol_seln_tr_pbmn": "3400000000",
      "whol_shnu_tr_pbmn": "3700000000",
      "whol_ntby_tr_pbmn": "300000000",
      "whol_seln_vol_amt": "43000",
      "whol_shnu_vol_amt": "46000"
    },
    {
      "bsop_hour": "090100",
      "invr_cls_code": "2",
      "invr_cls_name": "기관계",
      "arbt_seln_tr_pbmn": "800000000",
      "arbt_shnu_tr_pbmn": "750000000",
      "arbt_ntby_tr_pbmn": "-50000000",
      "arbt_seln_vol_amt": "10000",
      "arbt_shnu_vol_amt": "9500",
      "arbt_ntby_qty_amt": "-500",
      "nabt_seln_tr_pbmn": "1500000000",
      "nabt_shnu_tr_pbmn": "1400000000",
      "nabt_ntby_tr_pbmn": "-100000000",
      "nabt_seln_vol_amt": "19000",
      "nabt_shnu_vol_amt": "18000",
      "nabt_ntby_qty_amt": "-1000",
      "whol_seln_tr_pbmn": "2300000000",
      "whol_shnu_tr_pbmn": "2150000000",
      "whol_ntby_tr_pbmn": "-150000000",
      "whol_seln_vol_amt": "29000",
      "whol_shnu_vol_amt": "27500"
    }
  ]
}
```

---

## Task 2: domestic/investor.go append + `InquireInvestorTrendEstimate` (EP1)

**Files**
- Modify: `domestic/investor.go` (append)
- Modify: `domestic/investor_test.go` (append)

- [ ] Step 1: APPEND test code to `domestic/investor_test.go`

```go
func TestClient_InquireInvestorTrendEstimate(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/quotations/investor-trend-estimate`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "investor_trend_estimate_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireInvestorTrendEstimate(context.Background(), domestic.InquireInvestorTrendEstimateParams{
		Symbol: "005930",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, "005930", capturedQuery.Get("MKSC_SHRN_ISCD"))

	require.Len(t, res.Output2, 2)
	assert.Equal(t, "1", res.Output2[0].BsopHourGb)
	assert.Equal(t, int64(123456), res.Output2[0].FrgnFakeNtbyQty)
	assert.Equal(t, int64(-45678), res.Output2[0].OrgnFakeNtbyQty)
	assert.Equal(t, int64(77778), res.Output2[0].SumFakeNtbyQty)
}
```

- [ ] Step 2: Verify FAIL

```bash
cd /Users/user/src/workspace_moneyflow/korea-investment-stock
go test ./domestic/... -run TestClient_InquireInvestorTrendEstimate -v 2>&1 | tail -5
# Expected: FAIL (undefined: domestic.InquireInvestorTrendEstimateParams)
```

- [ ] Step 3: APPEND struct + Params + method to `domestic/investor.go`

```go
// InvestorTrendEstimate 는 주식현재가 투자자 가집계 (HHPTJ04160200) 응답.
//
// 한투 docs: docs/api/국내주식/주식현재가_투자자.md
// path: /uapi/domestic-stock/v1/quotations/investor-trend-estimate
//
// output2 (시간대별 가집계 Array). output1 없음. query 키가 FID_ prefix 없는 MKSC_SHRN_ISCD.
type InvestorTrendEstimate struct {
	Output2 []InvestorTrendEstimateItem `json:"output2"`
}

// InvestorTrendEstimateItem 은 응답 output2 한 행 (시간대 가집계).
//
// bsop_hour_gb: 1=09:30, 2=10:00, 3=11:20, 4=13:20, 5=14:30
type InvestorTrendEstimateItem struct {
	BsopHourGb      string `json:"bsop_hour_gb"`       // 입력구분 (1~5)
	FrgnFakeNtbyQty int64  `json:"frgn_fake_ntby_qty,string"` // 외국인 순매수 수량(가집계)
	OrgnFakeNtbyQty int64  `json:"orgn_fake_ntby_qty,string"` // 기관 순매수 수량(가집계)
	SumFakeNtbyQty  int64  `json:"sum_fake_ntby_qty,string"`  // 합산 순매수 수량(가집계)
}

// InquireInvestorTrendEstimateParams 는 주식현재가 투자자 가집계 조회 파라미터.
//
// query 키 MKSC_SHRN_ISCD — FID_ prefix 없음 (KIS docs 그대로).
type InquireInvestorTrendEstimateParams struct {
	Symbol string // MKSC_SHRN_ISCD — 필수, 종목코드 (6자리)
}

// InquireInvestorTrendEstimate 는 주식현재가 투자자 가집계 호출.
//
// 한투 docs: docs/api/국내주식/주식현재가_투자자.md
// path: /uapi/domestic-stock/v1/quotations/investor-trend-estimate (HHPTJ04160200)
func (c *Client) InquireInvestorTrendEstimate(ctx context.Context, params InquireInvestorTrendEstimateParams) (*InvestorTrendEstimate, error) {
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/domestic-stock/v1/quotations/investor-trend-estimate",
		TrID:   "HHPTJ04160200",
		Query: map[string]string{
			"MKSC_SHRN_ISCD": params.Symbol,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}

	var res InvestorTrendEstimate
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse InvestorTrendEstimate: %w", err)
	}
	return &res, nil
}
```

- [ ] Step 4: Verify PASS

```bash
go test ./domestic/... -run TestClient_InquireInvestorTrendEstimate -v 2>&1 | tail -5
# Expected: PASS
```

- [ ] Step 5: gofmt/vet

```bash
gofmt -w domestic/investor.go domestic/investor_test.go
go vet ./domestic/...
```

- [ ] Step 6: Commit

```bash
git commit -m "$(cat <<'EOF'
[feat] domestic: InquireInvestorTrendEstimate (EP1, Phase 2.5)

주식현재가 투자자 가집계 (HHPTJ04160200). output2 only (4 fields).
Non-FID query key MKSC_SHRN_ISCD. 시간대별(1~5) 외인/기관/합산 순매수 수량.

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 3: `InquireForeignInstitutionTotal` (EP2)

**Files**
- Modify: `domestic/investor.go` (append)
- Modify: `domestic/investor_test.go` (append)

- [ ] Step 1: APPEND test code to `domestic/investor_test.go`

```go
func TestClient_InquireForeignInstitutionTotal(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/quotations/foreign-institution-total`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "foreign_institution_total_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireForeignInstitutionTotal(context.Background(), domestic.InquireForeignInstitutionTotalParams{
		Symbol:       "0001",
		DivClsCode:   "0",
		RankSortCode: "0",
		EtcClsCode:   "0",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, "V", capturedQuery.Get("FID_COND_MRKT_DIV_CODE"))
	assert.Equal(t, "16449", capturedQuery.Get("FID_COND_SCR_DIV_CODE"))
	assert.Equal(t, "0001", capturedQuery.Get("FID_INPUT_ISCD"))
	assert.Equal(t, "0", capturedQuery.Get("FID_DIV_CLS_CODE"))

	require.Len(t, res.Output, 2)
	assert.Equal(t, "005930", res.Output[0].MkscShrnIscd)
	assert.Equal(t, "삼성전자", res.Output[0].HtsKorIsnm)
	assert.Equal(t, decimal.NewFromInt(75800), res.Output[0].StckPrpr)
	assert.Equal(t, int64(123456), res.Output[0].NtbyQty)
	assert.Equal(t, int64(100000), res.Output[0].FrgnNtbyQty)
	assert.Equal(t, int64(-50000), res.Output[0].OrgnNtbyQty)
}
```

- [ ] Step 2: Verify FAIL

```bash
go test ./domestic/... -run TestClient_InquireForeignInstitutionTotal -v 2>&1 | tail -5
# Expected: FAIL (undefined: domestic.InquireForeignInstitutionTotalParams)
```

- [ ] Step 3: APPEND struct + Params + method to `domestic/investor.go`

```go
// ForeignInstitutionTotal 는 국내기관 외국인 매매종목 가집계 (FHPTJ04400000) 응답.
//
// 한투 docs: docs/api/국내주식/국내기관_외국인_매매종목가집계.md
// path: /uapi/domestic-stock/v1/quotations/foreign-institution-total
//
// Output (capital O) — KIS API 의 json key 가 "Output" (대문자). 26 fields.
type ForeignInstitutionTotal struct {
	Output []ForeignInstitutionTotalItem `json:"Output"`
}

// ForeignInstitutionTotalItem 은 응답 Output 한 행 (종목별 외국인/기관 집계).
type ForeignInstitutionTotalItem struct {
	HtsKorIsnm        string          `json:"hts_kor_isnm"`             // HTS 한글 종목명
	MkscShrnIscd      string          `json:"mksc_shrn_iscd"`           // 유가증권 단축 종목코드
	NtbyQty           int64           `json:"ntby_qty,string"`          // 순매수 수량
	StckPrpr          decimal.Decimal `json:"stck_prpr"`                // 주식 현재가
	PrdyVrssSign      string          `json:"prdy_vrss_sign"`           // 전일 대비 부호
	PrdyVrss          decimal.Decimal `json:"prdy_vrss"`                // 전일 대비
	PrdyCtrt          float64         `json:"prdy_ctrt,string"`         // 전일 대비율
	AcmlVol           int64           `json:"acml_vol,string"`          // 누적 거래량
	FrgnNtbyQty       int64           `json:"frgn_ntby_qty,string"`     // 외국인 순매수 수량
	OrgnNtbyQty       int64           `json:"orgn_ntby_qty,string"`     // 기관계 순매수 수량
	IvtrNtbyQty       int64           `json:"ivtr_ntby_qty,string"`     // 투자신탁 순매수 수량
	BankNtbyQty       int64           `json:"bank_ntby_qty,string"`     // 은행 순매수 수량
	InsuNtbyQty       int64           `json:"insu_ntby_qty,string"`     // 보험 순매수 수량
	MrbnNtbyQty       int64           `json:"mrbn_ntby_qty,string"`     // 종금 순매수 수량
	FundNtbyQty       int64           `json:"fund_ntby_qty,string"`     // 기금 순매수 수량
	EtcOrgtNtbyVol    int64           `json:"etc_orgt_ntby_vol,string"` // 기타단체 순매수 거래량
	EtcCorpNtbyVol    int64           `json:"etc_corp_ntby_vol,string"` // 기타법인 순매수 거래량
	FrgnNtbyTrPbmn    int64           `json:"frgn_ntby_tr_pbmn,string"`     // 외국인 순매수 거래대금
	OrgnNtbyTrPbmn    int64           `json:"orgn_ntby_tr_pbmn,string"`     // 기관계 순매수 거래대금
	IvtrNtbyTrPbmn    int64           `json:"ivtr_ntby_tr_pbmn,string"`     // 투자신탁 순매수 거래대금
	BankNtbyTrPbmn    int64           `json:"bank_ntby_tr_pbmn,string"`     // 은행 순매수 거래대금
	InsuNtbyTrPbmn    int64           `json:"insu_ntby_tr_pbmn,string"`     // 보험 순매수 거래대금
	MrbnNtbyTrPbmn    int64           `json:"mrbn_ntby_tr_pbmn,string"`     // 종금 순매수 거래대금
	FundNtbyTrPbmn    int64           `json:"fund_ntby_tr_pbmn,string"`     // 기금 순매수 거래대금
	EtcOrgtNtbyTrPbmn int64           `json:"etc_orgt_ntby_tr_pbmn,string"` // 기타단체 순매수 거래대금
	EtcCorpNtbyTrPbmn int64           `json:"etc_corp_ntby_tr_pbmn,string"` // 기타법인 순매수 거래대금
}

// InquireForeignInstitutionTotalParams 는 국내기관 외국인 매매종목 가집계 조회 파라미터.
//
// MarketCode 기본값 "V", ScrDivCode 기본값 "16449".
// Symbol: "0000"=전체, "0001"=KOSPI, "1001"=KOSDAQ.
// DivClsCode: 0=수량정렬, 1=금액정렬.
// RankSortCode: 0=순매수상위, 1=순매도상위.
// EtcClsCode: 0=전체, 1=외국인, 2=기관계, 3=기타.
type InquireForeignInstitutionTotalParams struct {
	MarketCode   string // FID_COND_MRKT_DIV_CODE — 빈 값=>"V"
	ScrDivCode   string // FID_COND_SCR_DIV_CODE — 빈 값=>"16449"
	Symbol       string // FID_INPUT_ISCD — 필수
	DivClsCode   string // FID_DIV_CLS_CODE — 0:수량, 1:금액
	RankSortCode string // FID_RANK_SORT_CLS_CODE — 0:순매수, 1:순매도
	EtcClsCode   string // FID_ETC_CLS_CODE — 0:전체, 1:외국인, 2:기관계, 3:기타
}

// InquireForeignInstitutionTotal 는 국내기관 외국인 매매종목 가집계 호출.
//
// 한투 docs: docs/api/국내주식/국내기관_외국인_매매종목가집계.md
// path: /uapi/domestic-stock/v1/quotations/foreign-institution-total (FHPTJ04400000)
func (c *Client) InquireForeignInstitutionTotal(ctx context.Context, params InquireForeignInstitutionTotalParams) (*ForeignInstitutionTotal, error) {
	mkt := params.MarketCode
	if mkt == "" {
		mkt = "V"
	}
	scr := params.ScrDivCode
	if scr == "" {
		scr = "16449"
	}

	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/domestic-stock/v1/quotations/foreign-institution-total",
		TrID:   "FHPTJ04400000",
		Query: map[string]string{
			"FID_COND_MRKT_DIV_CODE":  mkt,
			"FID_COND_SCR_DIV_CODE":   scr,
			"FID_INPUT_ISCD":          params.Symbol,
			"FID_DIV_CLS_CODE":        params.DivClsCode,
			"FID_RANK_SORT_CLS_CODE":  params.RankSortCode,
			"FID_ETC_CLS_CODE":        params.EtcClsCode,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}

	var res ForeignInstitutionTotal
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse ForeignInstitutionTotal: %w", err)
	}
	return &res, nil
}
```

- [ ] Step 4: Verify PASS

```bash
go test ./domestic/... -run TestClient_InquireForeignInstitutionTotal -v 2>&1 | tail -5
# Expected: PASS
```

- [ ] Step 5: gofmt/vet

```bash
gofmt -w domestic/investor.go domestic/investor_test.go
go vet ./domestic/...
```

- [ ] Step 6: Commit

```bash
git commit -m "$(cat <<'EOF'
[feat] domestic: InquireForeignInstitutionTotal (EP2, Phase 2.5)

국내기관 외국인 매매종목 가집계 (FHPTJ04400000). Output (capital-O json key)
26 fields. default MarketCode="V", ScrDivCode="16449".

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 4: `domestic/program_trade.go` base + `InquireProgramTradeByStockDaily` (EP3)

**Files**
- Create: `domestic/program_trade.go`
- Create: `domestic/program_trade_test.go`

- [ ] Step 1: CREATE `domestic/program_trade_test.go` with first test

```go
// File: domestic/program_trade_test.go
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

func TestClient_InquireProgramTradeByStockDaily(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/quotations/program-trade-by-stock-daily`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "program_trade_by_stock_daily_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireProgramTradeByStockDaily(context.Background(), domestic.InquireProgramTradeByStockDailyParams{
		Symbol:   "005930",
		BaseDate: "0020260505",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, "J", capturedQuery.Get("FID_COND_MRKT_DIV_CODE"))
	assert.Equal(t, "005930", capturedQuery.Get("FID_INPUT_ISCD"))
	assert.Equal(t, "0020260505", capturedQuery.Get("FID_INPUT_DATE_1"))

	require.Len(t, res.Output, 2)
	assert.Equal(t, "20260505", res.Output[0].StckBsopDate)
	assert.Equal(t, decimal.NewFromInt(75800), res.Output[0].StckClpr)
	assert.Equal(t, int64(150000), res.Output[0].WholSmtnNtbyQty)
	assert.Equal(t, int64(10000), res.Output[0].WholNtbyVolIcdc)
	assert.Equal(t, int64(500000000), res.Output[0].WholNtbyTrPbmnIcdc2)
}
```

- [ ] Step 2: Verify FAIL

```bash
go test ./domestic/... -run TestClient_InquireProgramTradeByStockDaily -v 2>&1 | tail -5
# Expected: FAIL (undefined: domestic.InquireProgramTradeByStockDailyParams)
```

- [ ] Step 3: CREATE `domestic/program_trade.go` with struct + Params + method

```go
// File: domestic/program_trade.go
package domestic

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/shopspring/decimal"

	"github.com/kenshin579/korea-investment-stock/internal/httpclient"
)

// ProgramTradeByStockDaily 는 종목별 프로그램매매추이(일별) (FHPPG04650201) 응답.
//
// 한투 docs: docs/api/국내주식/종목별_프로그램매매추이(일별).md
// path: /uapi/domestic-stock/v1/quotations/program-trade-by-stock-daily
type ProgramTradeByStockDaily struct {
	Output []ProgramTradeByStockDailyItem `json:"output"`
}

// ProgramTradeByStockDailyItem 은 응답 output 한 행 (일자별 프로그램 매매).
//
// WholNtbyTrPbmnIcdc2: 마지막 필드. 필드명 trailing "2" — EP4 의 WholNtbyTrPbmnIcdc 와 구분.
type ProgramTradeByStockDailyItem struct {
	StckBsopDate          string          `json:"stck_bsop_date"`                  // 주식 영업 일자
	StckClpr              decimal.Decimal `json:"stck_clpr"`                       // 주식 종가
	PrdyVrss              decimal.Decimal `json:"prdy_vrss"`                       // 전일 대비
	PrdyVrssSign          string          `json:"prdy_vrss_sign"`                  // 전일 대비 부호
	PrdyCtrt              float64         `json:"prdy_ctrt,string"`                // 전일 대비율
	AcmlVol               int64           `json:"acml_vol,string"`                 // 누적 거래량
	AcmlTrPbmn            int64           `json:"acml_tr_pbmn,string"`             // 누적 거래 대금
	WholSmtnSelnVol       int64           `json:"whol_smtn_seln_vol,string"`       // 전체 합산 매도 수량
	WholSmtnShnuVol       int64           `json:"whol_smtn_shnu_vol,string"`       // 전체 합산 매수 수량
	WholSmtnNtbyQty       int64           `json:"whol_smtn_ntby_qty,string"`       // 전체 합산 순매수 수량
	WholSmtnSelnTrPbmn    int64           `json:"whol_smtn_seln_tr_pbmn,string"`   // 전체 합산 매도 거래대금
	WholSmtnShnuTrPbmn    int64           `json:"whol_smtn_shnu_tr_pbmn,string"`   // 전체 합산 매수 거래대금
	WholSmtnNtbyTrPbmn    int64           `json:"whol_smtn_ntby_tr_pbmn,string"`   // 전체 합산 순매수 거래대금
	WholNtbyVolIcdc       int64           `json:"whol_ntby_vol_icdc,string"`       // 전체 순매수 거래량 증감
	WholNtbyTrPbmnIcdc2   int64           `json:"whol_ntby_tr_pbmn_icdc2,string"`  // 전체 순매수 거래대금 증감 (trailing "2")
}

// InquireProgramTradeByStockDailyParams 는 종목별 프로그램매매추이(일별) 조회 파라미터.
//
// BaseDate: KIS docs 예시가 "002" prefix 포함 ("0020240308") — 호출자가 raw string 그대로 전달.
// MarketCode: "J"(KRX), "NX"(NXT), "UN"(통합). 빈 값=>"J".
type InquireProgramTradeByStockDailyParams struct {
	MarketCode string // FID_COND_MRKT_DIV_CODE — 빈 값=>"J"
	Symbol     string // FID_INPUT_ISCD — 필수
	BaseDate   string // FID_INPUT_DATE_1 — 필수, KIS docs 예시: "0020240308" ("002" prefix)
}

// InquireProgramTradeByStockDaily 는 종목별 프로그램매매추이(일별) 호출.
//
// 한투 docs: docs/api/국내주식/종목별_프로그램매매추이(일별).md
// path: /uapi/domestic-stock/v1/quotations/program-trade-by-stock-daily (FHPPG04650201)
func (c *Client) InquireProgramTradeByStockDaily(ctx context.Context, params InquireProgramTradeByStockDailyParams) (*ProgramTradeByStockDaily, error) {
	mkt := params.MarketCode
	if mkt == "" {
		mkt = "J"
	}

	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/domestic-stock/v1/quotations/program-trade-by-stock-daily",
		TrID:   "FHPPG04650201",
		Query: map[string]string{
			"FID_COND_MRKT_DIV_CODE": mkt,
			"FID_INPUT_ISCD":         params.Symbol,
			"FID_INPUT_DATE_1":       params.BaseDate,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}

	var res ProgramTradeByStockDaily
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse ProgramTradeByStockDaily: %w", err)
	}
	return &res, nil
}
```

- [ ] Step 4: Verify PASS

```bash
go test ./domestic/... -run TestClient_InquireProgramTradeByStockDaily -v 2>&1 | tail -5
# Expected: PASS
```

- [ ] Step 5: gofmt/vet

```bash
gofmt -w domestic/program_trade.go domestic/program_trade_test.go
go vet ./domestic/...
```

- [ ] Step 6: Commit

```bash
git commit -m "$(cat <<'EOF'
[feat] domestic: program_trade.go base + InquireProgramTradeByStockDaily (EP3, Phase 2.5)

종목별 프로그램매매추이(일별) (FHPPG04650201). 15 fields.
BaseDate "002" prefix doc'd in Params comment. 마지막 필드 WholNtbyTrPbmnIcdc2 (trailing 2).

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 5: `InquireProgramTradeByStock` (EP4)

**Files**
- Modify: `domestic/program_trade.go` (append)
- Modify: `domestic/program_trade_test.go` (append)

- [ ] Step 1: APPEND test code to `domestic/program_trade_test.go`

```go
func TestClient_InquireProgramTradeByStock(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/quotations/program-trade-by-stock`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "program_trade_by_stock_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireProgramTradeByStock(context.Background(), domestic.InquireProgramTradeByStockParams{
		Symbol: "005930",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, "J", capturedQuery.Get("FID_COND_MRKT_DIV_CODE"))
	assert.Equal(t, "005930", capturedQuery.Get("FID_INPUT_ISCD"))

	require.Len(t, res.Output, 2)
	assert.Equal(t, "090100", res.Output[0].BsopHour)
	assert.Equal(t, decimal.NewFromInt(75800), res.Output[0].StckPrpr)
	assert.Equal(t, int64(1500), res.Output[0].WholSmtnNtbyQty)
	assert.Equal(t, int64(200), res.Output[0].WholNtbyVolIcdc)
	assert.Equal(t, int64(15160000), res.Output[0].WholNtbyTrPbmnIcdc)
}
```

- [ ] Step 2: Verify FAIL

```bash
go test ./domestic/... -run TestClient_InquireProgramTradeByStock -v 2>&1 | tail -5
# Expected: FAIL (undefined: domestic.InquireProgramTradeByStockParams)
```

- [ ] Step 3: APPEND struct + Params + method to `domestic/program_trade.go`

```go
// ProgramTradeByStock 는 종목별 프로그램매매추이(체결) (FHPPG04650101) 응답.
//
// 한투 docs: docs/api/국내주식/종목별_프로그램매매추이(체결).md
// path: /uapi/domestic-stock/v1/quotations/program-trade-by-stock
type ProgramTradeByStock struct {
	Output []ProgramTradeByStockItem `json:"output"`
}

// ProgramTradeByStockItem 은 응답 output 한 행 (시간대별 프로그램 체결).
//
// WholNtbyTrPbmnIcdc: 마지막 필드. trailing "2" 없음 — EP3 의 WholNtbyTrPbmnIcdc2 와 구분.
type ProgramTradeByStockItem struct {
	BsopHour           string          `json:"bsop_hour"`                       // 영업 시간
	StckPrpr           decimal.Decimal `json:"stck_prpr"`                       // 주식 현재가
	PrdyVrss           decimal.Decimal `json:"prdy_vrss"`                       // 전일 대비
	PrdyVrssSign       string          `json:"prdy_vrss_sign"`                  // 전일 대비 부호
	PrdyCtrt           float64         `json:"prdy_ctrt,string"`                // 전일 대비율
	AcmlVol            int64           `json:"acml_vol,string"`                 // 누적 거래량
	WholSmtnSelnVol    int64           `json:"whol_smtn_seln_vol,string"`       // 전체 합산 매도 수량
	WholSmtnShnuVol    int64           `json:"whol_smtn_shnu_vol,string"`       // 전체 합산 매수 수량
	WholSmtnNtbyQty    int64           `json:"whol_smtn_ntby_qty,string"`       // 전체 합산 순매수 수량
	WholSmtnSelnTrPbmn int64           `json:"whol_smtn_seln_tr_pbmn,string"`   // 전체 합산 매도 거래대금
	WholSmtnShnuTrPbmn int64           `json:"whol_smtn_shnu_tr_pbmn,string"`   // 전체 합산 매수 거래대금
	WholSmtnNtbyTrPbmn int64           `json:"whol_smtn_ntby_tr_pbmn,string"`   // 전체 합산 순매수 거래대금
	WholNtbyVolIcdc    int64           `json:"whol_ntby_vol_icdc,string"`       // 전체 순매수 거래량 증감
	WholNtbyTrPbmnIcdc int64           `json:"whol_ntby_tr_pbmn_icdc,string"`   // 전체 순매수 거래대금 증감 (no trailing "2")
}

// InquireProgramTradeByStockParams 는 종목별 프로그램매매추이(체결) 조회 파라미터.
type InquireProgramTradeByStockParams struct {
	MarketCode string // FID_COND_MRKT_DIV_CODE — 빈 값=>"J"
	Symbol     string // FID_INPUT_ISCD — 필수
}

// InquireProgramTradeByStock 는 종목별 프로그램매매추이(체결) 호출.
//
// 한투 docs: docs/api/국내주식/종목별_프로그램매매추이(체결).md
// path: /uapi/domestic-stock/v1/quotations/program-trade-by-stock (FHPPG04650101)
func (c *Client) InquireProgramTradeByStock(ctx context.Context, params InquireProgramTradeByStockParams) (*ProgramTradeByStock, error) {
	mkt := params.MarketCode
	if mkt == "" {
		mkt = "J"
	}

	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/domestic-stock/v1/quotations/program-trade-by-stock",
		TrID:   "FHPPG04650101",
		Query: map[string]string{
			"FID_COND_MRKT_DIV_CODE": mkt,
			"FID_INPUT_ISCD":         params.Symbol,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}

	var res ProgramTradeByStock
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse ProgramTradeByStock: %w", err)
	}
	return &res, nil
}
```

- [ ] Step 4: Verify PASS

```bash
go test ./domestic/... -run TestClient_InquireProgramTradeByStock -v 2>&1 | tail -5
# Expected: PASS
```

- [ ] Step 5: gofmt/vet

```bash
gofmt -w domestic/program_trade.go domestic/program_trade_test.go
go vet ./domestic/...
```

- [ ] Step 6: Commit

```bash
git commit -m "$(cat <<'EOF'
[feat] domestic: InquireProgramTradeByStock (EP4, Phase 2.5)

종목별 프로그램매매추이(체결) (FHPPG04650101). 14 fields.
마지막 필드 WholNtbyTrPbmnIcdc (no trailing "2" — EP3 대비 구분).

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 6: `InquireCompProgramTradeToday` (EP5)

**Files**
- Modify: `domestic/program_trade.go` (append)
- Modify: `domestic/program_trade_test.go` (append)

- [ ] Step 1: APPEND test code to `domestic/program_trade_test.go`

```go
func TestClient_InquireCompProgramTradeToday(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/quotations/comp-program-trade-today`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "comp_program_trade_today_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireCompProgramTradeToday(context.Background(), domestic.InquireCompProgramTradeTodayParams{
		MarketCode:  "J",
		MrktClsCode: "K",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, "J", capturedQuery.Get("FID_COND_MRKT_DIV_CODE"))
	assert.Equal(t, "K", capturedQuery.Get("FID_MRKT_CLS_CODE"))
	assert.Equal(t, "", capturedQuery.Get("FID_SCTN_CLS_CODE"))
	assert.Equal(t, "", capturedQuery.Get("FID_INPUT_ISCD"))

	require.Len(t, res.Output1, 2)
	assert.Equal(t, "090100", res.Output1[0].BsopHour)
	assert.Equal(t, int64(12000000000), res.Output1[0].ArbtSmtnSelnTrPbmn)
	assert.Equal(t, 35.50, res.Output1[0].ArbtSmtmSelnTrPbmnRate)
	assert.Equal(t, 36.10, res.Output1[0].ArbtSmtmShunTrPbmnRate)
	assert.Equal(t, decimal.NewFromFloat(2750.50), res.Output1[0].BstpNmixPrpr)
}
```

- [ ] Step 2: Verify FAIL

```bash
go test ./domestic/... -run TestClient_InquireCompProgramTradeToday -v 2>&1 | tail -5
# Expected: FAIL (undefined: domestic.InquireCompProgramTradeTodayParams)
```

- [ ] Step 3: APPEND struct + Params + method to `domestic/program_trade.go`

```go
// CompProgramTradeToday 는 프로그램매매 종합현황(시간) (FHPPG04600101) 응답.
//
// 한투 docs: docs/api/국내주식/프로그램매매_종합현황(시간).md
// path: /uapi/domestic-stock/v1/quotations/comp-program-trade-today
//
// output1 (시간별 Array). 18 fields.
// 비고: `smtm`(rate 계열)과 `smtn`(금액/수량 계열) 혼용 — KIS API 원문 그대로.
// `shun` typo 2개 (arbt_smtm_shun_tr_pbmn_rate, nabt_smtm_shun_tr_pbmn_rate) 원문 보존.
type CompProgramTradeToday struct {
	Output1 []CompProgramTradeTodayItem `json:"output1"`
}

// CompProgramTradeTodayItem 은 응답 output1 한 행 (시간대별 종합현황).
//
// 필드명 내 typo 주의:
// - ArbtSmtmShunTrPbmnRate: "shun" (KIS docs typo) — 매수 율이지만 필드명 오기
// - NabtSmtmShunTrPbmnRate: 동일 패턴
type CompProgramTradeTodayItem struct {
	BsopHour                  string          `json:"bsop_hour"`                        // 영업 시간
	ArbtSmtnSelnTrPbmn        int64           `json:"arbt_smtn_seln_tr_pbmn,string"`    // 차익 합산 매도 거래대금
	ArbtSmtmSelnTrPbmnRate    float64         `json:"arbt_smtm_seln_tr_pbmn_rate,string"` // 차익 합산 매도 거래대금 비율
	ArbtSmtnShnuTrPbmn        int64           `json:"arbt_smtn_shnu_tr_pbmn,string"`    // 차익 합산 매수 거래대금
	ArbtSmtmShunTrPbmnRate    float64         `json:"arbt_smtm_shun_tr_pbmn_rate,string"` // 차익 합산 매수 거래대금 비율 ("shun" typo)
	NabtSmtnSelnTrPbmn        int64           `json:"nabt_smtn_seln_tr_pbmn,string"`    // 비차익 합산 매도 거래대금
	NabtSmtmSelnTrPbmnRate    float64         `json:"nabt_smtm_seln_tr_pbmn_rate,string"` // 비차익 합산 매도 거래대금 비율
	NabtSmtnShnuTrPbmn        int64           `json:"nabt_smtn_shnu_tr_pbmn,string"`    // 비차익 합산 매수 거래대금
	NabtSmtmShunTrPbmnRate    float64         `json:"nabt_smtm_shun_tr_pbmn_rate,string"` // 비차익 합산 매수 거래대금 비율 ("shun" typo)
	ArbtSmtnNtbyTrPbmn        int64           `json:"arbt_smtn_ntby_tr_pbmn,string"`    // 차익 합산 순매수 거래대금
	ArbtSmtmNtbyTrPbmnRate    float64         `json:"arbt_smtm_ntby_tr_pbmn_rate,string"` // 차익 합산 순매수 거래대금 비율
	NabtSmtnNtbyTrPbmn        int64           `json:"nabt_smtn_ntby_tr_pbmn,string"`    // 비차익 합산 순매수 거래대금
	NabtSmtmNtbyTrPbmnRate    float64         `json:"nabt_smtm_ntby_tr_pbmn_rate,string"` // 비차익 합산 순매수 거래대금 비율
	WholSmtnNtbyTrPbmn        int64           `json:"whol_smtn_ntby_tr_pbmn,string"`    // 전체 합산 순매수 거래대금
	WholNtbyTrPbmnRate        float64         `json:"whol_ntby_tr_pbmn_rate,string"`    // 전체 순매수 거래대금 비율
	BstpNmixPrpr              decimal.Decimal `json:"bstp_nmix_prpr"`                   // 업종 지수 현재가
	BstpNmixPrdyVrss          decimal.Decimal `json:"bstp_nmix_prdy_vrss"`              // 업종 지수 전일 대비
	PrdyVrssSign              string          `json:"prdy_vrss_sign"`                   // 전일 대비 부호
}

// InquireCompProgramTradeTodayParams 는 프로그램매매 종합현황(시간) 조회 파라미터.
//
// 6개 query 파라미터 중 첫 2개만 의미있음. 나머지 4개는 빈 문자열 전송.
// MrktClsCode: "K"=코스피, "Q"=코스닥.
type InquireCompProgramTradeTodayParams struct {
	MarketCode  string // FID_COND_MRKT_DIV_CODE — 필수
	MrktClsCode string // FID_MRKT_CLS_CODE — K:코스피, Q:코스닥
	// 아래 4개는 항상 "" 전송 (KIS docs 명시)
	// SctnClsCode FID_SCTN_CLS_CODE
	// Symbol      FID_INPUT_ISCD
	// MarketCode1 FID_COND_MRKT_DIV_CODE1
	// InputHour1  FID_INPUT_HOUR_1
}

// InquireCompProgramTradeToday 는 프로그램매매 종합현황(시간) 호출.
//
// 한투 docs: docs/api/국내주식/프로그램매매_종합현황(시간).md
// path: /uapi/domestic-stock/v1/quotations/comp-program-trade-today (FHPPG04600101)
func (c *Client) InquireCompProgramTradeToday(ctx context.Context, params InquireCompProgramTradeTodayParams) (*CompProgramTradeToday, error) {
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/domestic-stock/v1/quotations/comp-program-trade-today",
		TrID:   "FHPPG04600101",
		Query: map[string]string{
			"FID_COND_MRKT_DIV_CODE":  params.MarketCode,
			"FID_MRKT_CLS_CODE":       params.MrktClsCode,
			"FID_SCTN_CLS_CODE":       "",
			"FID_INPUT_ISCD":          "",
			"FID_COND_MRKT_DIV_CODE1": "",
			"FID_INPUT_HOUR_1":        "",
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}

	var res CompProgramTradeToday
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse CompProgramTradeToday: %w", err)
	}
	return &res, nil
}
```

- [ ] Step 4: Verify PASS

```bash
go test ./domestic/... -run TestClient_InquireCompProgramTradeToday -v 2>&1 | tail -5
# Expected: PASS
```

- [ ] Step 5: gofmt/vet

```bash
gofmt -w domestic/program_trade.go domestic/program_trade_test.go
go vet ./domestic/...
```

- [ ] Step 6: Commit

```bash
git commit -m "$(cat <<'EOF'
[feat] domestic: InquireCompProgramTradeToday (EP5, Phase 2.5)

프로그램매매 종합현황(시간) (FHPPG04600101). output1 18 fields.
6 query params (4 blank). shun/smtm typos preserved verbatim per KIS API.

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 7: InquireCompProgramTradeDaily (EP6) — 97 fields, LARGEST

**Endpoint**

| 항목 | 값 |
|---|---|
| Path | `/uapi/domestic-stock/v1/quotations/comp-program-trade-daily` |
| TR_ID | `FHPPG04600001` |
| Output | `output []CompProgramTradeDailyItem` (97 fields per row) |
| File | APPEND to `domestic/program_trade.go` and `domestic/program_trade_test.go` |

**Query Params (4)**

| Go field | KIS key | 설명 |
|---|---|---|
| MarketCode | `FID_COND_MRKT_DIV_CODE` | 시장 구분 코드 |
| MrktClsCode | `FID_MRKT_CLS_CODE` | K:코스피, Q:코스닥 |
| StartDate | `FID_INPUT_DATE_1` | blank 또는 YYYYMMDD (8개월 lookback 한도) |
| EndDate | `FID_INPUT_DATE_2` | blank 또는 YYYYMMDD |

**CompProgramTradeDailyItem 필드 테이블 (97 fields)**

> 타입 규칙: `_rate` 로 끝나면 `float64` (JSON string); `stck_bsop_date` → `string`; 그 외 → `int64` (JSON string).
> Implementer 는 KIS docs 원문 순서·철자를 따른다. 아래 table 이 기준이며 count 는 KIS docs 원문 우선.
> `shun` 타이포(KIS docs 명시)를 그대로 보존.

| # | JSON key | 한국어 설명 | Go type |
|---|---|---|---|
| 1 | `stck_bsop_date` | 주식 영업 일자 | string |
| 2 | `nabt_entm_seln_tr_pbmn` | 비차익 위탁 매도 거래대금 | int64,string |
| 3 | `nabt_onsl_seln_vol` | 비차익 자기 매도 거래량 | int64,string |
| 4 | `whol_onsl_seln_tr_pbmn` | 전체 자기 매도 거래대금 | int64,string |
| 5 | `arbt_smtn_shnu_vol` | 차익 합계 매수 거래량 | int64,string |
| 6 | `nabt_smtn_shnu_tr_pbmn` | 비차익 합계 매수 거래대금 | int64,string |
| 7 | `arbt_entm_ntby_qty` | 차익 위탁 순매수 수량 | int64,string |
| 8 | `nabt_entm_ntby_tr_pbmn` | 비차익 위탁 순매수 거래대금 | int64,string |
| 9 | `arbt_entm_seln_vol` | 차익 위탁 매도 거래량 | int64,string |
| 10 | `nabt_entm_seln_vol_rate` | 비차익 위탁 매도 거래량 비율 | float64,string |
| 11 | `nabt_onsl_seln_vol_rate` | 비차익 자기 매도 거래량 비율 | float64,string |
| 12 | `whol_onsl_seln_tr_pbmn_rate` | 전체 자기 매도 거래대금 비율 | float64,string |
| 13 | `arbt_smtm_shun_vol_rate` | 차익 합계 매수 거래량 비율 (shun 타이포) | float64,string |
| 14 | `nabt_smtm_shun_tr_pbmn_rate` | 비차익 합계 매수 거래대금 비율 (shun 타이포) | float64,string |
| 15 | `arbt_entm_ntby_qty_rate` | 차익 위탁 순매수 수량 비율 | float64,string |
| 16 | `nabt_entm_ntby_tr_pbmn_rate` | 비차익 위탁 순매수 거래대금 비율 | float64,string |
| 17 | `arbt_entm_seln_vol_rate` | 차익 위탁 매도 거래량 비율 | float64,string |
| 18 | `nabt_entm_seln_tr_pbmn_rate` | 비차익 위탁 매도 거래대금 비율 | float64,string |
| 19 | `nabt_onsl_seln_tr_pbmn` | 비차익 자기 매도 거래대금 | int64,string |
| 20 | `whol_smtn_seln_vol` | 전체 합계 매도 거래량 | int64,string |
| 21 | `arbt_smtn_shnu_tr_pbmn` | 차익 합계 매수 거래대금 | int64,string |
| 22 | `whol_entm_shnu_vol` | 전체 위탁 매수 거래량 | int64,string |
| 23 | `arbt_entm_ntby_tr_pbmn` | 차익 위탁 순매수 거래대금 | int64,string |
| 24 | `nabt_onsl_ntby_qty` | 비차익 자기 순매수 수량 | int64,string |
| 25 | `arbt_entm_seln_tr_pbmn` | 차익 위탁 매도 거래대금 | int64,string |
| 26 | `nabt_onsl_seln_tr_pbmn_rate` | 비차익 자기 매도 거래대금 비율 | float64,string |
| 27 | `whol_seln_vol_rate` | 전체 매도 거래량 비율 | float64,string |
| 28 | `arbt_smtm_shun_tr_pbmn_rate` | 차익 합계 매수 거래대금 비율 (shun 타이포) | float64,string |
| 29 | `whol_entm_shnu_vol_rate` | 전체 위탁 매수 거래량 비율 | float64,string |
| 30 | `arbt_entm_ntby_tr_pbmn_rate` | 차익 위탁 순매수 거래대금 비율 | float64,string |
| 31 | `nabt_onsl_ntby_qty_rate` | 비차익 자기 순매수 수량 비율 | float64,string |
| 32 | `arbt_entm_seln_tr_pbmn_rate` | 차익 위탁 매도 거래대금 비율 | float64,string |
| 33 | `nabt_smtn_seln_vol` | 비차익 합계 매도 거래량 | int64,string |
| 34 | `whol_smtn_seln_tr_pbmn` | 전체 합계 매도 거래대금 | int64,string |
| 35 | `nabt_entm_shnu_vol` | 비차익 위탁 매수 거래량 | int64,string |
| 36 | `whol_entm_shnu_tr_pbmn` | 전체 위탁 매수 거래대금 | int64,string |
| 37 | `arbt_onsl_ntby_qty` | 차익 자기 순매수 수량 | int64,string |
| 38 | `nabt_onsl_ntby_tr_pbmn` | 비차익 자기 순매수 거래대금 | int64,string |
| 39 | `arbt_onsl_seln_tr_pbmn` | 차익 자기 매도 거래대금 | int64,string |
| 40 | `nabt_smtm_seln_vol_rate` | 비차익 합계 매도 거래량 비율 | float64,string |
| 41 | `whol_seln_tr_pbmn_rate` | 전체 매도 거래대금 비율 | float64,string |
| 42 | `nabt_entm_shnu_vol_rate` | 비차익 위탁 매수 거래량 비율 | float64,string |
| 43 | `whol_entm_shnu_tr_pbmn_rate` | 전체 위탁 매수 거래대금 비율 | float64,string |
| 44 | `arbt_onsl_ntby_qty_rate` | 차익 자기 순매수 수량 비율 | float64,string |
| 45 | `nabt_onsl_ntby_tr_pbmn_rate` | 비차익 자기 순매수 거래대금 비율 | float64,string |
| 46 | `arbt_onsl_seln_tr_pbmn_rate` | 차익 자기 매도 거래대금 비율 | float64,string |
| 47 | `nabt_smtn_seln_tr_pbmn` | 비차익 합계 매도 거래대금 | int64,string |
| 48 | `arbt_entm_shnu_vol` | 차익 위탁 매수 거래량 | int64,string |
| 49 | `nabt_entm_shnu_tr_pbmn` | 비차익 위탁 매수 거래대금 | int64,string |
| 50 | `whol_onsl_shnu_vol` | 전체 자기 매수 거래량 | int64,string |
| 51 | `arbt_onsl_ntby_tr_pbmn` | 차익 자기 순매수 거래대금 | int64,string |
| 52 | `nabt_smtn_ntby_qty` | 비차익 합계 순매수 수량 | int64,string |
| 53 | `arbt_onsl_seln_vol` | 차익 자기 매도 거래량 | int64,string |
| 54 | `nabt_smtm_seln_tr_pbmn_rate` | 비차익 합계 매도 거래대금 비율 | float64,string |
| 55 | `arbt_entm_shnu_vol_rate` | 차익 위탁 매수 거래량 비율 | float64,string |
| 56 | `nabt_entm_shnu_tr_pbmn_rate` | 비차익 위탁 매수 거래대금 비율 | float64,string |
| 57 | `whol_onsl_shnu_tr_pbmn` | 전체 자기 매수 거래대금 | int64,string |
| 58 | `arbt_onsl_ntby_tr_pbmn_rate` | 차익 자기 순매수 거래대금 비율 | float64,string |
| 59 | `nabt_smtm_ntby_qty_rate` | 비차익 합계 순매수 수량 비율 | float64,string |
| 60 | `arbt_onsl_seln_vol_rate` | 차익 자기 매도 거래량 비율 | float64,string |
| 61 | `whol_entm_seln_vol` | 전체 위탁 매도 거래량 | int64,string |
| 62 | `arbt_entm_shnu_tr_pbmn` | 차익 위탁 매수 거래대금 | int64,string |
| 63 | `nabt_onsl_shnu_vol` | 비차익 자기 매수 거래량 | int64,string |
| 64 | `whol_onsl_shnu_tr_pbmn_rate` | 전체 자기 매수 거래대금 비율 | float64,string |
| 65 | `arbt_smtn_ntby_qty` | 차익 합계 순매수 수량 | int64,string |
| 66 | `nabt_smtn_ntby_tr_pbmn` | 비차익 합계 순매수 거래대금 | int64,string |
| 67 | `arbt_smtn_seln_vol` | 차익 합계 매도 거래량 | int64,string |
| 68 | `whol_entm_seln_tr_pbmn` | 전체 위탁 매도 거래대금 | int64,string |
| 69 | `arbt_entm_shnu_tr_pbmn_rate` | 차익 위탁 매수 거래대금 비율 | float64,string |
| 70 | `nabt_onsl_shnu_vol_rate` | 비차익 자기 매수 거래량 비율 | float64,string |
| 71 | `whol_onsl_shnu_vol_rate` | 전체 자기 매수 거래량 비율 | float64,string |
| 72 | `arbt_smtm_ntby_qty_rate` | 차익 합계 순매수 수량 비율 | float64,string |
| 73 | `nabt_smtm_ntby_tr_pbmn_rate` | 비차익 합계 순매수 거래대금 비율 | float64,string |
| 74 | `arbt_smtm_seln_vol_rate` | 차익 합계 매도 거래량 비율 | float64,string |
| 75 | `whol_entm_seln_vol_rate` | 전체 위탁 매도 거래량 비율 | float64,string |
| 76 | `arbt_onsl_shnu_vol` | 차익 자기 매수 거래량 | int64,string |
| 77 | `nabt_onsl_shnu_tr_pbmn` | 비차익 자기 매수 거래대금 | int64,string |
| 78 | `whol_smtn_shnu_vol` | 전체 합계 매수 거래량 | int64,string |
| 79 | `arbt_smtn_ntby_tr_pbmn` | 차익 합계 순매수 거래대금 | int64,string |
| 80 | `whol_entm_ntby_qty` | 전체 위탁 순매수 수량 | int64,string |
| 81 | `arbt_smtn_seln_tr_pbmn` | 차익 합계 매도 거래대금 | int64,string |
| 82 | `whol_entm_seln_tr_pbmn_rate` | 전체 위탁 매도 거래대금 비율 | float64,string |
| 83 | `arbt_onsl_shnu_vol_rate` | 차익 자기 매수 거래량 비율 | float64,string |
| 84 | `nabt_onsl_shnu_tr_pbmn_rate` | 비차익 자기 매수 거래대금 비율 | float64,string |
| 85 | `whol_shun_vol_rate` | 전체 매수 거래량 비율 (shun 타이포) | float64,string |
| 86 | `arbt_smtm_ntby_tr_pbmn_rate` | 차익 합계 순매수 거래대금 비율 | float64,string |
| 87 | `whol_entm_ntby_qty_rate` | 전체 위탁 순매수 수량 비율 | float64,string |
| 88 | `arbt_smtm_seln_tr_pbmn_rate` | 차익 합계 매도 거래대금 비율 | float64,string |
| 89 | `whol_onsl_seln_vol` | 전체 자기 매도 거래량 | int64,string |
| 90 | `arbt_onsl_shnu_tr_pbmn` | 차익 자기 매수 거래대금 | int64,string |
| 91 | `nabt_smtn_shnu_vol` | 비차익 합계 매수 거래량 | int64,string |
| 92 | `whol_smtn_shnu_tr_pbmn` | 전체 합계 매수 거래대금 | int64,string |
| 93 | `nabt_entm_ntby_qty` | 비차익 위탁 순매수 수량 | int64,string |
| 94 | `whol_entm_ntby_tr_pbmn` | 전체 위탁 순매수 거래대금 | int64,string |
| 95 | `nabt_entm_seln_vol` | 비차익 위탁 매도 거래량 | int64,string |
| 96 | `whol_onsl_seln_vol_rate` | 전체 자기 매도 거래량 비율 | float64,string |
| 97 | `arbt_onsl_shnu_tr_pbmn_rate` | 차익 자기 매수 거래대금 비율 | float64,string |

> **참고**: KIS docs 원문의 정확한 필드 수가 권위 있는 기준. 위 테이블은 97개이나, 실제 API 응답과 대조하여 dedup/추가할 것.
> `nabt_smtm_shun_vol_rate`, `whol_shun_tr_pbmn_rate`, `nabt_entm_ntby_qty_rate` 등 추가 rate 필드가 KIS docs 원문에 있으면 포함.

**구현 단계 (6 steps)**

- [ ] Step 1: APPEND test to `domestic/program_trade_test.go`

  Fixture 파일: `domestic/testdata/comp_program_trade_daily_success.json`

  ```go
  func TestClient_InquireCompProgramTradeDaily(t *testing.T) {
      httpmock.Activate()
      defer httpmock.DeactivateAndReset()

      raw, err := os.ReadFile("testdata/comp_program_trade_daily_success.json")
      require.NoError(t, err)
      httpmock.RegisterResponder("GET",
          "=~.*comp-program-trade-daily.*",
          httpmock.NewBytesResponder(200, raw))

      client := newTestClient(t)
      ctx := context.Background()
      res, err := client.Domestic.InquireCompProgramTradeDaily(ctx, domestic.InquireCompProgramTradeDailyParams{
          MarketCode:  "J",
          MrktClsCode: "K",
          StartDate:   "20260101",
          EndDate:     "20260505",
      })
      require.NoError(t, err)
      require.NotEmpty(t, res.Output)

      item := res.Output[0]
      assert.NotEmpty(t, item.StckBsopDate)               // string
      assert.GreaterOrEqual(t, item.NabtEntmSelnTrPbmn, int64(0)) // int64
      assert.GreaterOrEqual(t, item.NabtEntmSelnVolRate, float64(0)) // float64
      assert.GreaterOrEqual(t, item.ArbtSmtmShunVolRate, float64(0)) // float64 (shun typo)
  }
  ```

- [ ] Step 2: Verify FAIL

  ```bash
  go test ./domestic/... -run TestClient_InquireCompProgramTradeDaily -v 2>&1 | tail -5
  # Expected: compile error or FAIL (struct not yet defined)
  ```

- [ ] Step 3: APPEND to `domestic/program_trade.go`

  구조체는 위 필드 테이블에서 타입 규칙을 적용하여 직접 생성:
  - JSON tag = 테이블의 JSON key 그대로 (snake_case)
  - Go field name = PascalCase 변환 (예: `stck_bsop_date` → `StckBsopDate`)
  - `_rate` 로 끝나면 `float64` + `` `json:"...,string"` ``
  - `stck_bsop_date` → `string` + `` `json:"stck_bsop_date"` ``
  - 그 외 → `int64` + `` `json:"...,string"` ``

  ```go
  // CompProgramTradeDailyItem 은 프로그램매매 종합현황(일별) 한 행을 나타낸다.
  // FHPPG04600001 output 배열 원소. 8개월 lookback 한도.
  // shun 타이포 필드(arbt_smtm_shun_vol_rate 등)는 KIS docs 원문 보존.
  type CompProgramTradeDailyItem struct {
      // ... 97 fields from table above ...
  }

  // CompProgramTradeDaily 는 FHPPG04600001 전체 응답이다.
  type CompProgramTradeDaily struct {
      RTCd   string                      `json:"rt_cd"`
      MsgCd  string                      `json:"msg_cd"`
      Msg1   string                      `json:"msg1"`
      Output []CompProgramTradeDailyItem `json:"output"`
  }

  // InquireCompProgramTradeDailyParams 는 EP6 쿼리 파라미터다.
  type InquireCompProgramTradeDailyParams struct {
      MarketCode  string // FID_COND_MRKT_DIV_CODE
      MrktClsCode string // FID_MRKT_CLS_CODE  K:코스피  Q:코스닥
      StartDate   string // FID_INPUT_DATE_1  blank or YYYYMMDD  8개월 한도
      EndDate     string // FID_INPUT_DATE_2  blank or YYYYMMDD
  }

  // InquireCompProgramTradeDaily 는 프로그램매매 종합현황(일별)을 조회한다 (FHPPG04600001).
  func (c *Client) InquireCompProgramTradeDaily(ctx context.Context, p InquireCompProgramTradeDailyParams) (*CompProgramTradeDaily, error) {
      resp, err := c.client.Get(ctx, "/uapi/domestic-stock/v1/quotations/comp-program-trade-daily", map[string]string{
          "FID_COND_MRKT_DIV_CODE": p.MarketCode,
          "FID_MRKT_CLS_CODE":      p.MrktClsCode,
          "FID_INPUT_DATE_1":       p.StartDate,
          "FID_INPUT_DATE_2":       p.EndDate,
      }, "FHPPG04600001")
      if err != nil {
          return nil, err
      }

      var res CompProgramTradeDaily
      if err := json.Unmarshal(resp.Raw, &res); err != nil {
          return nil, fmt.Errorf("kis: parse CompProgramTradeDaily: %w", err)
      }
      return &res, nil
  }
  ```

- [ ] Step 4: Verify PASS

  ```bash
  go test ./domestic/... -run TestClient_InquireCompProgramTradeDaily -v 2>&1 | tail -5
  # Expected: PASS
  ```

- [ ] Step 5: gofmt/vet

  ```bash
  gofmt -w domestic/program_trade.go domestic/program_trade_test.go
  go vet ./domestic/...
  ```

- [ ] Step 6: Commit

  ```bash
  git add domestic/program_trade.go domestic/program_trade_test.go \
          domestic/testdata/comp_program_trade_daily_success.json
  git commit -m "$(cat <<'EOF'
  [feat] domestic — InquireCompProgramTradeDaily (프로그램매매 종합현황 일별, FHPPG04600001)

  EP6, Phase 2.5. output[] 97 fields. 8개월 lookback 한도.
  shun 타이포 필드명 KIS docs 원문 보존 (arbt_smtm_shun_vol_rate 등).

  Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
  EOF
  )"
  ```

---

## Task 8: InquireInvestorProgramTradeToday (EP7)

**Endpoint**

| 항목 | 값 |
|---|---|
| Path | `/uapi/domestic-stock/v1/quotations/investor-program-trade-today` |
| TR_ID | `HHPPG046600C1` |
| Output | `output1 []InvestorProgramTradeTodayItem` (20 fields) |
| File | APPEND to `domestic/program_trade.go` and `domestic/program_trade_test.go` |

**Query Params (2 — 비-FID prefix)**

| Go field | KIS key | 설명 |
|---|---|---|
| ExchDivClsCode | `EXCH_DIV_CLS_CODE` | J / NX / UN |
| MrktDivClsCode | `MRKT_DIV_CLS_CODE` | 1:코스피, 4:코스닥 |

> 주의: 이 endpoint 의 query key 는 `FID_` prefix 없음. `_amt` suffix 사용 (다른 endpoint 의 `_tr_pbmn` 와 다름).

**InvestorProgramTradeTodayItem 필드 (20)**

| JSON key | 한국어 설명 | Go type |
|---|---|---|
| `invr_cls_code` | 투자자 구분 코드 | string |
| `invr_cls_name` | 투자자 구분 이름 | string |
| `all_seln_qty` | 전체 매도 수량 | int64,string |
| `all_seln_amt` | 전체 매도 금액 | int64,string |
| `all_shnu_qty` | 전체 매수 수량 | int64,string |
| `all_shnu_amt` | 전체 매수 금액 | int64,string |
| `all_ntby_qty` | 전체 순매수 수량 | int64,string |
| `all_ntby_amt` | 전체 순매수 금액 | int64,string |
| `arbt_seln_qty` | 차익 매도 수량 | int64,string |
| `arbt_seln_amt` | 차익 매도 금액 | int64,string |
| `arbt_shnu_qty` | 차익 매수 수량 | int64,string |
| `arbt_shnu_amt` | 차익 매수 금액 | int64,string |
| `arbt_ntby_qty` | 차익 순매수 수량 | int64,string |
| `arbt_ntby_amt` | 차익 순매수 금액 | int64,string |
| `nabt_seln_qty` | 비차익 매도 수량 | int64,string |
| `nabt_seln_amt` | 비차익 매도 금액 | int64,string |
| `nabt_shnu_qty` | 비차익 매수 수량 | int64,string |
| `nabt_shnu_amt` | 비차익 매수 금액 | int64,string |
| `nabt_ntby_qty` | 비차익 순매수 수량 | int64,string |
| `nabt_ntby_amt` | 비차익 순매수 금액 | int64,string |

**구현 단계 (6 steps)**

- [ ] Step 1: APPEND test to `domestic/program_trade_test.go`

  Fixture 파일: `domestic/testdata/investor_program_trade_today_success.json`

  ```go
  func TestClient_InquireInvestorProgramTradeToday(t *testing.T) {
      httpmock.Activate()
      defer httpmock.DeactivateAndReset()

      raw, err := os.ReadFile("testdata/investor_program_trade_today_success.json")
      require.NoError(t, err)
      httpmock.RegisterResponder("GET",
          "=~.*investor-program-trade-today.*",
          httpmock.NewBytesResponder(200, raw))

      client := newTestClient(t)
      ctx := context.Background()
      res, err := client.Domestic.InquireInvestorProgramTradeToday(ctx, domestic.InquireInvestorProgramTradeTodayParams{
          ExchDivClsCode: "J",
          MrktDivClsCode: "1",
      })
      require.NoError(t, err)
      require.NotEmpty(t, res.Output1)

      item := res.Output1[0]
      assert.NotEmpty(t, item.InvrClsCode)
      assert.NotEmpty(t, item.InvrClsName)
      assert.GreaterOrEqual(t, item.AllSelnQty, int64(0))
      assert.GreaterOrEqual(t, item.ArbtNtbyAmt, int64(0))
      assert.GreaterOrEqual(t, item.NabtNtbyAmt, int64(0))
  }
  ```

- [ ] Step 2: Verify FAIL

  ```bash
  go test ./domestic/... -run TestClient_InquireInvestorProgramTradeToday -v 2>&1 | tail -5
  # Expected: compile error or FAIL (struct not yet defined)
  ```

- [ ] Step 3: APPEND to `domestic/program_trade.go`

  ```go
  // InvestorProgramTradeTodayItem 은 당일 투자자별 프로그램매매 동향 한 행을 나타낸다.
  // HHPPG046600C1 output1 배열 원소. 차익/비차익 breakdown 포함.
  // _amt suffix 사용 (다른 endpoint 의 _tr_pbmn 와 다름).
  type InvestorProgramTradeTodayItem struct {
      InvrClsCode string `json:"invr_cls_code"`
      InvrClsName string `json:"invr_cls_name"`
      AllSelnQty  int64  `json:"all_seln_qty,string"`
      AllSelnAmt  int64  `json:"all_seln_amt,string"`
      AllShnuQty  int64  `json:"all_shnu_qty,string"`
      AllShnuAmt  int64  `json:"all_shnu_amt,string"`
      AllNtbyQty  int64  `json:"all_ntby_qty,string"`
      AllNtbyAmt  int64  `json:"all_ntby_amt,string"`
      ArbtSelnQty int64  `json:"arbt_seln_qty,string"`
      ArbtSelnAmt int64  `json:"arbt_seln_amt,string"`
      ArbtShnuQty int64  `json:"arbt_shnu_qty,string"`
      ArbtShnuAmt int64  `json:"arbt_shnu_amt,string"`
      ArbtNtbyQty int64  `json:"arbt_ntby_qty,string"`
      ArbtNtbyAmt int64  `json:"arbt_ntby_amt,string"`
      NabtSelnQty int64  `json:"nabt_seln_qty,string"`
      NabtSelnAmt int64  `json:"nabt_seln_amt,string"`
      NabtShnuQty int64  `json:"nabt_shnu_qty,string"`
      NabtShnuAmt int64  `json:"nabt_shnu_amt,string"`
      NabtNtbyQty int64  `json:"nabt_ntby_qty,string"`
      NabtNtbyAmt int64  `json:"nabt_ntby_amt,string"`
  }

  // InvestorProgramTradeToday 는 HHPPG046600C1 전체 응답이다.
  type InvestorProgramTradeToday struct {
      RTCd    string                           `json:"rt_cd"`
      MsgCd   string                           `json:"msg_cd"`
      Msg1    string                           `json:"msg1"`
      Output1 []InvestorProgramTradeTodayItem  `json:"output1"`
  }

  // InquireInvestorProgramTradeTodayParams 는 EP7 쿼리 파라미터다.
  // KIS docs 기준 비-FID prefix 사용.
  type InquireInvestorProgramTradeTodayParams struct {
      ExchDivClsCode string // EXCH_DIV_CLS_CODE  J / NX / UN
      MrktDivClsCode string // MRKT_DIV_CLS_CODE  1:코스피  4:코스닥
  }

  // InquireInvestorProgramTradeToday 는 당일 투자자별 프로그램매매 동향을 조회한다 (HHPPG046600C1).
  func (c *Client) InquireInvestorProgramTradeToday(ctx context.Context, p InquireInvestorProgramTradeTodayParams) (*InvestorProgramTradeToday, error) {
      resp, err := c.client.Get(ctx, "/uapi/domestic-stock/v1/quotations/investor-program-trade-today", map[string]string{
          "EXCH_DIV_CLS_CODE": p.ExchDivClsCode,
          "MRKT_DIV_CLS_CODE": p.MrktDivClsCode,
      }, "HHPPG046600C1")
      if err != nil {
          return nil, err
      }

      var res InvestorProgramTradeToday
      if err := json.Unmarshal(resp.Raw, &res); err != nil {
          return nil, fmt.Errorf("kis: parse InvestorProgramTradeToday: %w", err)
      }
      return &res, nil
  }
  ```

- [ ] Step 4: Verify PASS

  ```bash
  go test ./domestic/... -run TestClient_InquireInvestorProgramTradeToday -v 2>&1 | tail -5
  # Expected: PASS
  ```

- [ ] Step 5: gofmt/vet

  ```bash
  gofmt -w domestic/program_trade.go domestic/program_trade_test.go
  go vet ./domestic/...
  ```

- [ ] Step 6: Commit

  ```bash
  git add domestic/program_trade.go domestic/program_trade_test.go \
          domestic/testdata/investor_program_trade_today_success.json
  git commit -m "$(cat <<'EOF'
  [feat] domestic — InquireInvestorProgramTradeToday (당일 투자자별 프로그램매매 동향, HHPPG046600C1)

  EP7, Phase 2.5. output1[] 20 fields. 비-FID query params.
  _amt suffix (다른 EP 의 _tr_pbmn 와 구분). 차익/비차익 breakdown.

  Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
  EOF
  )"
  ```

---

## Task 9: examples/domestic_investor_flow/main.go

EP1-EP7 전체를 시연하는 예제 파일을 생성한다.

**파일 경로**: `examples/domestic_investor_flow/main.go`

**구현 단계**

- [ ] Step 1: 예제 파일 작성

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

      // 1. 투자자 추정 (외인기관 가집계) — EP1 HHPTJ04160200
      fmt.Println("=== EP1: 투자자 매매 추정 가집계 ===")
      est, err := client.Domestic.InquireInvestorTrendEstimate(ctx, domestic.InquireInvestorTrendEstimateParams{
          MarketCode: "J",
          StockCode:  "005930",
      })
      if err != nil {
          log.Printf("EP1 error: %v", err)
      } else {
          fmt.Printf("  외인 순매수: %d\n", est.Output.FrgnNtbyQty)
          fmt.Printf("  기관 순매수: %d\n", est.Output.OrgNtbyQty)
          fmt.Printf("  개인 순매수: %d\n", est.Output.PsndNtbyQty)
          fmt.Printf("  합산 순매수: %d\n", est.Output.SmtnNtbyQty)
      }

      // 2. 외국인/기관 매매종목 집계 — EP2 FHPTJ04400000
      fmt.Println("\n=== EP2: 외인기관 매매종목가 집계 ===")
      fi, err := client.Domestic.InquireForeignInstitutionTotal(ctx, domestic.InquireForeignInstitutionTotalParams{
          MarketCode:       "J",
          SortDivisionCode: "0",
          StockCode:        "005930",
          TrgtClsCode:      "111111111",
      })
      if err != nil {
          log.Printf("EP2 error: %v", err)
      } else if len(fi.Output) > 0 {
          item := fi.Output[0]
          fmt.Printf("  종목명: %s\n", item.HtsFidIsnmDvsnCd)
          fmt.Printf("  외인 순매수: %d\n", item.FrgnNtbyQty)
          fmt.Printf("  기관 순매수: %d\n", item.OrgNtbyQty)
      }

      // 3. 종목별 프로그램매매 추이(일별) — EP3 FHPPG04650201
      fmt.Println("\n=== EP3: 종목별 프로그램매매 추이(일별) ===")
      ptd, err := client.Domestic.InquireProgramTradeByStockDaily(ctx, domestic.InquireProgramTradeByStockDailyParams{
          MarketCode: "J",
          StockCode:  "005930",
          BaseDate:   "0020260505", // 002 prefix (KIS docs 예시)
      })
      if err != nil {
          log.Printf("EP3 error: %v", err)
      } else if len(ptd.Output) > 0 {
          item := ptd.Output[0]
          fmt.Printf("  영업일: %s\n", item.StckBsopDate)
          fmt.Printf("  전체 순매수: %d\n", item.WholNtbyQty)
          fmt.Printf("  차익 순매수: %d\n", item.ArbtNtbyQty)
      }

      // 4. 종목별 프로그램매매 추이(체결) — EP4 FHPPG04650101
      fmt.Println("\n=== EP4: 종목별 프로그램매매 추이(체결) ===")
      pt, err := client.Domestic.InquireProgramTradeByStock(ctx, domestic.InquireProgramTradeByStockParams{
          MarketCode: "J",
          StockCode:  "005930",
      })
      if err != nil {
          log.Printf("EP4 error: %v", err)
      } else if len(pt.Output) > 0 {
          item := pt.Output[0]
          fmt.Printf("  체결시간: %s\n", item.StckCntgHour)
          fmt.Printf("  전체 순매수: %d\n", item.WholNtbyQty)
          fmt.Printf("  비차익 순매수: %d\n", item.NabtNtbyQty)
      }

      // 5. 프로그램매매 종합현황(시간) — EP5 FHPPG04600101
      fmt.Println("\n=== EP5: 프로그램매매 종합현황(시간) ===")
      cpt, err := client.Domestic.InquireCompProgramTradeToday(ctx, domestic.InquireCompProgramTradeTodayParams{
          MarketCode:  "J",
          MrktClsCode: "K",
      })
      if err != nil {
          log.Printf("EP5 error: %v", err)
      } else if len(cpt.Output1) > 0 {
          item := cpt.Output1[0]
          fmt.Printf("  시간: %s\n", item.BsopHour)
          fmt.Printf("  차익 순매수 금액: %d\n", item.ArbtNtbyTrPbmn)
          fmt.Printf("  비차익 순매수 금액: %d\n", item.NabtNtbyTrPbmn)
      }

      // 6. 프로그램매매 종합현황(일별) — EP6 FHPPG04600001
      fmt.Println("\n=== EP6: 프로그램매매 종합현황(일별) ===")
      cpd, err := client.Domestic.InquireCompProgramTradeDaily(ctx, domestic.InquireCompProgramTradeDailyParams{
          MarketCode:  "J",
          MrktClsCode: "K",
          StartDate:   "20260101",
          EndDate:     "20260505",
      })
      if err != nil {
          log.Printf("EP6 error: %v", err)
      } else if len(cpd.Output) > 0 {
          item := cpd.Output[0]
          fmt.Printf("  영업일: %s\n", item.StckBsopDate)
          fmt.Printf("  비차익 위탁 매도 거래대금: %d\n", item.NabtEntmSelnTrPbmn)
          fmt.Printf("  차익 합계 매수 거래량: %d\n", item.ArbtSmtnShnuVol)
      }

      // 7. 당일 투자자별 프로그램매매 동향 — EP7 HHPPG046600C1
      fmt.Println("\n=== EP7: 당일 투자자별 프로그램매매 동향 ===")
      ipt, err := client.Domestic.InquireInvestorProgramTradeToday(ctx, domestic.InquireInvestorProgramTradeTodayParams{
          ExchDivClsCode: "J",
          MrktDivClsCode: "1",
      })
      if err != nil {
          log.Printf("EP7 error: %v", err)
      } else if len(ipt.Output1) > 0 {
          item := ipt.Output1[0]
          fmt.Printf("  투자자: %s (%s)\n", item.InvrClsName, item.InvrClsCode)
          fmt.Printf("  전체 순매수 금액: %d\n", item.AllNtbyAmt)
          fmt.Printf("  차익 순매수 금액: %d\n", item.ArbtNtbyAmt)
          fmt.Printf("  비차익 순매수 금액: %d\n", item.NabtNtbyAmt)
      }
  }
  ```

- [ ] Step 2: 빌드 확인

  ```bash
  go build ./examples/domestic_investor_flow && echo OK
  ```

- [ ] Step 3: 빌드 산출물 정리 (leaked binary)

  ```bash
  rm -f domestic_investor_flow
  ```

- [ ] Step 4: gofmt

  ```bash
  gofmt -w examples/domestic_investor_flow/main.go
  ```

- [ ] Step 5: Commit

  ```bash
  git add examples/domestic_investor_flow/main.go
  git commit -m "$(cat <<'EOF'
  [feat] examples — domestic_investor_flow (EP1-EP7 Phase 2.5 시연)

  투자자 추정(EP1), 외인기관 집계(EP2), 종목별 프로그램매매 일별/체결(EP3/EP4),
  프로그램매매 종합현황 시간/일별(EP5/EP6), 투자자별 당일 동향(EP7) 7개 메서드 시연.

  Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
  EOF
  )"
  ```

---

## Task 10: 문서 갱신

4개 파일을 업데이트한다.

**구현 단계**

- [ ] Step 1: `CLAUDE.md` 배너 갱신

  변경 내용:
  - `Phase 2.4 — 예탁원 정보 확장 11 메서드 (v1.7.0). Phase 2.5+ 는 추후 sub-plan 으로.` →
    `Phase 2.5 — 투자자/매매 동향 7 메서드 (v1.8.0). Phase 2.5+ design spec 및 plan 참고.`
  - Phase 2.5+ design spec 링크 추가: `docs/superpowers/specs/2026-05-05-phase2-5plus-extension-design.md`
  - Phase 2.5 plan 링크 추가: `docs/superpowers/specs/2026-05-05-phase2-5-investor-flow-implementation-plan.md`

- [ ] Step 2: `README.md` 메서드 섹션 갱신

  변경 내용:
  - 헤딩 `Available Methods (Phase 1.2 ~ 2.4)` → `Available Methods (Phase 1.2 ~ 2.5)`
  - 메서드 수 `53` → `60`
  - 아래 7개 행을 domestic 메서드 테이블에 추가:

  | Method | TR_ID | Description |
  |---|---|---|
  | `Domestic.InquireInvestorTrendEstimate` | `HHPTJ04160200` | 투자자 매매 추정 가집계 |
  | `Domestic.InquireForeignInstitutionTotal` | `FHPTJ04400000` | 외인기관 매매종목가 집계 |
  | `Domestic.InquireProgramTradeByStockDaily` | `FHPPG04650201` | 종목별 프로그램매매 추이(일별) |
  | `Domestic.InquireProgramTradeByStock` | `FHPPG04650101` | 종목별 프로그램매매 추이(체결) |
  | `Domestic.InquireCompProgramTradeToday` | `FHPPG04600101` | 프로그램매매 종합현황(시간) |
  | `Domestic.InquireCompProgramTradeDaily` | `FHPPG04600001` | 프로그램매매 종합현황(일별) |
  | `Domestic.InquireInvestorProgramTradeToday` | `HHPPG046600C1` | 당일 투자자별 프로그램매매 동향 |

- [ ] Step 3: `CHANGELOG.md` — `## [1.7.0]` 위에 추가

  ```markdown
  ## [1.8.0] - 2026-05-05

  ### Added — Phase 2.5 (투자자/매매 동향)

  - `Domestic.InquireInvestorTrendEstimate` — 투자자 매매 추정 가집계 (HHPTJ04160200) — 외국인/기관/합산 가집계 4 fields
  - `Domestic.InquireForeignInstitutionTotal` — 외인기관 매매종목가 집계 (FHPTJ04400000) — 26 fields, 8 투자자 종류 ntby
  - `Domestic.InquireProgramTradeByStockDaily` — 종목별 프로그램매매추이(일별) (FHPPG04650201) — 15 fields
  - `Domestic.InquireProgramTradeByStock` — 종목별 프로그램매매추이(체결) (FHPPG04650101) — 14 fields
  - `Domestic.InquireCompProgramTradeToday` — 프로그램매매 종합현황(시간) (FHPPG04600101) — 18 fields
  - `Domestic.InquireCompProgramTradeDaily` — 프로그램매매 종합현황(일별) (FHPPG04600001) — 97 fields (largest in Phase 2.5)
  - `Domestic.InquireInvestorProgramTradeToday` — 당일 투자자별 프로그램매매 동향 (HHPPG046600C1) — 20 fields, 차익/비차익 breakdown
  - examples: `domestic_investor_flow`

  ### Notes

  - EP2 응답 키는 `Output` (대문자 O) — KIS docs 명시. `json:"Output"` 사용.
  - EP3 의 `FID_INPUT_DATE_1` 은 KIS docs 예시에서 "002" prefix 사용 (e.g., "0020240308"). 호출자가 raw string 전달.
  - EP3 vs EP4: 마지막 필드 `whol_ntby_tr_pbmn_icdc2` (EP3) vs `whol_ntby_tr_pbmn_icdc` (EP4) — 변경 시 주의.
  - EP5/EP6 의 일부 rate field 명에 `shun` 타이포 (KIS docs 명시) 보존: `arbt_smtm_shun_tr_pbmn_rate`, `nabt_smtm_shun_tr_pbmn_rate`, `whol_shun_vol_rate`, `whol_shun_tr_pbmn_rate` 등.
  - EP6 응답 struct 97 필드 (Phase 2.5 최대). 8개월 lookback 한도.
  - EP7 query param 은 비-FID prefix (`EXCH_DIV_CLS_CODE`, `MRKT_DIV_CLS_CODE`) + MRKT 값 "1"/"4" (코스피/코스닥). 필드 suffix 는 `_amt` (다른 endpoint 의 `_tr_pbmn` 와 다름).
  ```

- [ ] Step 4: `domestic/doc.go` — Phase 2.4 섹션 다음에 Phase 2.5 섹션 추가

  ```go
  // Phase 2.5 — 투자자/매매 동향 (v1.8.0)
  //
  //   InquireInvestorTrendEstimate       HHPTJ04160200  투자자 매매 추정 가집계
  //   InquireForeignInstitutionTotal     FHPTJ04400000  외인기관 매매종목가 집계
  //   InquireProgramTradeByStockDaily    FHPPG04650201  종목별 프로그램매매 추이(일별)
  //   InquireProgramTradeByStock         FHPPG04650101  종목별 프로그램매매 추이(체결)
  //   InquireCompProgramTradeToday       FHPPG04600101  프로그램매매 종합현황(시간)
  //   InquireCompProgramTradeDaily       FHPPG04600001  프로그램매매 종합현황(일별)
  //   InquireInvestorProgramTradeToday   HHPPG046600C1  당일 투자자별 프로그램매매 동향
  ```

- [ ] Step 5: 빌드/vet 확인

  ```bash
  go build ./... && go vet ./...
  ```

- [ ] Step 6: Commit

  ```bash
  git add CLAUDE.md README.md CHANGELOG.md domestic/doc.go
  git commit -m "$(cat <<'EOF'
  [docs] Phase 2.5 문서 갱신 (v1.8.0)

  CLAUDE.md 배너, README.md 메서드 목록 53→60, CHANGELOG.md v1.8.0 항목,
  domestic/doc.go Phase 2.5 섹션 추가.

  Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
  EOF
  )"
  ```

---

## Task 11: 최종 점검

**구현 단계**

- [ ] Step 1: gofmt 검사

  ```bash
  gofmt -l . | head
  # Expected: 출력 없음 (포맷 이상 없음)
  ```

- [ ] Step 2: 빌드 및 vet

  ```bash
  go build ./... && go vet ./...
  # Expected: 오류 없음
  ```

- [ ] Step 3: 전체 테스트 (race detector)

  ```bash
  go test ./... -race -count=1
  # Expected: 모든 패키지 PASS
  ```

- [ ] Step 4: 커버리지 측정

  ```bash
  # domestic 패키지
  go test ./domestic/... -coverprofile=/tmp/cov_d.out -covermode=atomic
  go tool cover -func=/tmp/cov_d.out | tail -2
  # Expected: total ≥ 80%

  # root kis 패키지
  go test . -coverprofile=/tmp/cov_r.out -covermode=atomic
  go tool cover -func=/tmp/cov_r.out | tail -2
  # Expected: total ≥ 80%
  ```

  커버리지 < 80% 인 경우: 필수 파라미터 검증, 기본값 체크 등 추가 테스트 작성 후 재측정.

- [ ] Step 5: 파일 존재 확인

  ```bash
  ls domestic/program_trade.go \
     domestic/program_trade_test.go \
     domestic/testdata/investor_trend_estimate_success.json \
     domestic/testdata/foreign_institution_total_success.json \
     domestic/testdata/program_trade_by_stock_daily_success.json \
     domestic/testdata/program_trade_by_stock_success.json \
     domestic/testdata/comp_program_trade_today_success.json \
     domestic/testdata/comp_program_trade_daily_success.json \
     domestic/testdata/investor_program_trade_today_success.json \
     examples/domestic_investor_flow/main.go \
     2>&1 | wc -l
  # Expected: 10
  ```

- [ ] Step 6: 브랜치 커밋 수 확인

  ```bash
  git log main..HEAD --oneline | wc -l
  # 예상: Task 1-10 의 개별 commit 수와 일치 (약 10-12개)
  ```

---

## Task 12: PR 생성 (사용자 승인 후)

> **Claude 는 push / PR 생성을 사용자 명시적 승인 후에만 실행 (글로벌 정책).**

**구현 단계**

- [ ] Step 1: **사용자 승인 요청**

  Task 11 결과 보고 (빌드 OK, 테스트 PASS, 커버리지 수치, 파일 수, 커밋 수) 후 PR 생성 진행 여부를 사용자에게 확인한다.

- [ ] Step 2: Push (사용자 승인 후)

  ```bash
  git push -u origin docs/phase2-5plus-spec
  ```

- [ ] Step 3: PR 생성

  ```bash
  gh pr create \
    --title "Phase 2.5 — 투자자/매매 동향 (v1.8.0) + Phase 2.5+ design" \
    --reviewer kenshin579 \
    --base main \
    --head docs/phase2-5plus-spec \
    --body "$(cat <<'EOF'
  ## Summary

  - Phase 2.5+ extension design spec + Phase 2.5 implementation (7 methods)
  - Phase 2 패턴 그대로 (Style A, Params struct, KIS docs 1:1)
  - v1.8.0 release (누적 53 → 60 메서드)

  ## 메서드 → 한투 API 매핑

  | Go 메서드 | path | TR_ID |
  |---|---|---|
  | InquireInvestorTrendEstimate | quotations/investor-trend-estimate | HHPTJ04160200 |
  | InquireForeignInstitutionTotal | quotations/foreign-institution-total | FHPTJ04400000 |
  | InquireProgramTradeByStockDaily | quotations/program-trade-by-stock-daily | FHPPG04650201 |
  | InquireProgramTradeByStock | quotations/program-trade-by-stock | FHPPG04650101 |
  | InquireCompProgramTradeToday | quotations/comp-program-trade-today | FHPPG04600101 |
  | InquireCompProgramTradeDaily | quotations/comp-program-trade-daily | FHPPG04600001 |
  | InquireInvestorProgramTradeToday | quotations/investor-program-trade-today | HHPPG046600C1 |

  ## Anomalies handled

  - EP2 capital-O `Output` key
  - EP3 "002" date prefix doc'd in Params comment
  - EP3 vs EP4 trailing-2 distinction
  - EP5/EP6 `shun` typo preserved verbatim
  - EP6 97-field struct + 8-month lookback note
  - EP7 non-FID query + `_amt` vs `_tr_pbmn` suffix

  ## Test Plan
  - [x] go build/vet/fmt clean
  - [x] go test ./... -race -count=1 모든 패키지 PASS
  - [x] Coverage domestic ≥80%, root kis ≥80%
  - [x] httpmock 단위 테스트 (7 methods)
  - [x] examples/domestic_investor_flow build OK

  ## Breaking Changes
  없음.

  ## 참고 문서
  - Phase 2.5+ design spec: docs/superpowers/specs/2026-05-05-phase2-5plus-extension-design.md
  - Phase 2.5 plan: docs/superpowers/specs/2026-05-05-phase2-5-investor-flow-implementation-plan.md

  🤖 Generated with [Claude Code](https://claude.com/claude-code)
  EOF
  )"
  ```

- [ ] Step 4: Merge (사용자 승인 후)

  ```bash
  gh pr merge <PR#> --merge
  ```

- [ ] Step 5: 태그 및 릴리스 (사용자 승인 후)

  ```bash
  git tag -a v1.8.0 -m "Phase 2.5 — 투자자/매매 동향 (7 메서드, 누적 60)"
  git push origin v1.8.0
  gh release create v1.8.0 \
    --title "v1.8.0 — Phase 2.5 투자자/매매 동향" \
    --notes-file <(awk '/^## \[1\.8\.0\]/{f=1;next} /^## \[/&&f{exit} f' CHANGELOG.md)
  ```
