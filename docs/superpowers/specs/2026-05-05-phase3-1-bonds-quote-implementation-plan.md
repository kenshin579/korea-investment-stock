# Phase 3.1 — 장내채권 시세 8 메서드 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** korea-investment-stock Go 라이브러리에 신규 `bonds/` sub-package 도입 + 장내채권 시세 8 메서드 추가 (`v1.11.0` release). 새 패키지 scaffolding (Client, doc.go) + 8 메서드 TDD 구현.

**Architecture:** Phase 2 인프라 + 패턴 재사용. 신규 `bonds/` sub-package 신설. Root `client.go` 에 `Bonds *bonds.Client` 필드 추가. TDD: testdata fixture → 실패 테스트 → struct + 메서드 구현 → 통과 → commit.

**Tech Stack:** Go 1.25+, `github.com/jarcoal/httpmock`, `github.com/stretchr/testify`, `github.com/shopspring/decimal`. 새 dependency 없음.

**참고 spec:**
- Phase 3 bonds design spec: `docs/superpowers/specs/2026-05-05-phase3-bonds-design.md`
- Phase 2.7 plan (참조 패턴 — compact task structure): `docs/superpowers/specs/2026-05-05-phase2-7-industry-implementation-plan.md`
- Phase 2.4 KSD plan (all-string struct 패턴 참조): `docs/superpowers/specs/2026-05-05-phase2-4-ksd-implementation-plan.md`

---

## 사전 정보

| 항목 | 값 |
|------|---|
| 작업 브랜치 | `feat/phase3-1-bonds-quote` |
| 시작 HEAD | Phase 2.7 구현 완료 commit (v1.10.0) |
| Release 목표 | `v1.11.0` |
| PR 베이스 | `main` |
| 현재 main 상태 | v1.10.0 publish 완료 (Phase 2.7 통합, 71 메서드) |

> **CRITICAL NOTE:** Phase 3.1 은 완전 신규 `bonds/` sub-package 도입. EP1(`SearchBondInfo`) + EP2(`InquireIssueInfo`) 는 KSD-like all-string 70/69 필드 (Phase 2 standard 타입 매핑 **미적용**). EP3~EP8 은 Phase 2 standard typed 매핑 적용. 누적: 71 → **79** 메서드.

---

## 메서드 매핑

| # | Go 메서드 | path (last segment) | TR_ID | output key | fields | notes |
|---|---|---|---|---|---|---|
| EP1 | `SearchBondInfo` | `search-bond-info` | CTPF1114R | `output{}` | 70 | all-string KSD-like, `Inquire` prefix 없음 |
| EP2 | `InquireIssueInfo` | `issue-info` | CTPF1101R | `output{}` | 69 | all-string |
| EP3 | `InquirePrice` | `inquire-price` | FHKBJ773400C0 | `output{}` | 17 | typed |
| EP4 | `InquireCcnl` | `inquire-ccnl` | FHKBJ773403C0 | `output{}` | 7 | typed, single snapshot |
| EP5 | `InquireAskingPrice` | `inquire-asking-price` | FHKBJ773401C0 | `output{}` | 34 | typed, 5단계 호가 |
| EP6 | `InquireDailyPrice` | `inquire-daily-price` | FHKBJ773404C0 | `output{}` | 9 | typed |
| EP7 | `InquireDailyItemchartprice` | `inquire-daily-itemchartprice` | FHKBJ773701C0 | `output[]` | 6/item | typed array |
| EP8 | `InquireAvgUnit` | `avg-unit` | CTPF2005R | `output1{}+output2[]{}+output3[]{}` | 23+10+16 | typed, CTX_AREA pagination |

Default `FID_COND_MRKT_DIV_CODE` = `"B"` (Bond market) for EP3~EP6.

---

## 파일 구조

### 신규 (bonds package)

- `bonds/client.go` — Client struct + `New(http *httpclient.Client) *Client`
- `bonds/doc.go` — package doc (8 메서드 목록)
- `bonds/quote.go` — 8 메서드 + structs + Params
- `bonds/quote_test.go` — 8 테스트 함수
- `bonds/testhelper_test.go` — `newTestClient(t)` + `loadFixtureString(t, name)` helpers

### 신규 (testdata — 8 fixtures)

- `bonds/testdata/search_bond_info_success.json`
- `bonds/testdata/issue_info_success.json`
- `bonds/testdata/bond_price_success.json`
- `bonds/testdata/bond_ccnl_success.json`
- `bonds/testdata/bond_asking_price_success.json`
- `bonds/testdata/bond_daily_price_success.json`
- `bonds/testdata/bond_daily_itemchartprice_success.json`
- `bonds/testdata/avg_unit_success.json`

### 신규 (examples)

- `examples/bonds_quote/main.go` — 8 메서드 integration example

### 수정

- root `client.go` — `Bonds *bonds.Client` 필드 추가 + `wireInfra` 에서 `c.Bonds = bonds.New(c.httpClient)` 추가
- `CLAUDE.md` — banner Phase 2.7 → Phase 3.1, plan link 추가
- `README.md` — Available Methods 표 갱신 (71 → 79 메서드), bonds section 추가
- `CHANGELOG.md` — `[1.11.0]` entry ABOVE `[1.10.0]`

---

## 타입 매핑

### EP1 + EP2 (KSD-like all-string)

KIS docs 가 70/69 fields 를 모두 `String` 타입으로 명시. 모든 필드 plain `string`. Phase 2 standard typed mapping **미적용**.

### EP3 ~ EP8 (Phase 2 standard)

| 카테고리 | Go 타입 | json tag suffix | 예시 필드 |
|---|---|---|---|
| 채권 가격 | `decimal.Decimal` | (bare) | `bond_prpr`, `bond_oprc`, `bond_hgpr`, `bond_lwpr`, `bond_prdy_vrss`, `bond_mxpr`, `bond_llam`, `bond_askp1`~`bond_askp5`, `bond_bidp1`~`bond_bidp5`, `kis_unpr`, `kbp_unpr`, `nice_evlu_unpr`, `fnp_unpr`, `avg_evlu_unpr`, `*_rf_unpr`, `*_evlu_unit_pric`, `*_evlu_pric` |
| 평가금액 | `int64` | `,string` | `*_evlu_amt` |
| 거래량/잔량 | `int64` | `,string` | `acml_vol`, `cntg_vol`, `askp_rsqn1`~`askp_rsqn5`, `bidp_rsqn1`~`bidp_rsqn5`, `total_askp_rsqn`, `total_bidp_rsqn`, `ntby_aspr_rsqn` |
| 비율/수익률/등락률 | `float64` | `,string` | `prdy_ctrt`, `ernn_rate`, `oprc_ert`, `hgpr_ert`, `lwpr_ert`, `seln_ernn_rate1`~`seln_ernn_rate5`, `shnu_ernn_rate1`~`shnu_ernn_rate5`, `kis_evlu_erng_rt`, `kbp_evlu_erng_rt`, `nice_evlu_erng_rt`, `fnp_evlu_erng_rt`, `avg_evlu_erng_rt` |
| 코드/이름/날짜/시간/통화/Y-N/등급텍스트 | `string` | (bare) | `stnd_iscd`, `hts_kor_isnm`, `prdy_vrss_sign`, `stck_cntg_hour`, `stck_bsop_date`, `aspr_acpt_hour`, `evlu_dt`, `pdno`, `prdt_type_cd`, `prdt_name`, `crcy_cd`, `chng_yn`, `*_crdt_grad_text` |

---

## Tasks (13 total) — overview

| # | 내용 | Files |
|---|---|---|
| Task 1 | testdata fixtures (8 합성 JSON) | `bonds/testdata/*.json` |
| Task 2 | bonds 패키지 scaffolding + `SearchBondInfo` (EP1) | CREATE `bonds/client.go`, `bonds/doc.go`, `bonds/quote.go`, `bonds/quote_test.go`, `bonds/testhelper_test.go`; MODIFY root `client.go` |
| Task 3 | `InquireIssueInfo` (EP2) | APPEND `bonds/quote.go` / `bonds/quote_test.go` |
| Task 4 | `InquirePrice` (EP3) | APPEND `bonds/quote.go` / `bonds/quote_test.go` |
| Task 5 | `InquireCcnl` (EP4) | APPEND `bonds/quote.go` / `bonds/quote_test.go` |
| Task 6 | `InquireAskingPrice` (EP5) | APPEND `bonds/quote.go` / `bonds/quote_test.go` |
| Task 7 | `InquireDailyPrice` (EP6) | APPEND `bonds/quote.go` / `bonds/quote_test.go` |
| Task 8 | `InquireDailyItemchartprice` (EP7) | APPEND `bonds/quote.go` / `bonds/quote_test.go` |
| Task 9 | `InquireAvgUnit` (EP8, CTX_AREA pagination) | APPEND `bonds/quote.go` / `bonds/quote_test.go` |
| Task 10 | examples `examples/bonds_quote/main.go` | CREATE |
| Task 11 | `bonds/doc.go` Phase 3.1 section (finalize) | MODIFY |
| Task 12 | `CHANGELOG.md` `[1.11.0]` entry | MODIFY |
| Task 13 | `CLAUDE.md` + `README.md` 갱신 + `go test ./...` 전체 + tag | MODIFY |

---

## Task 1: testdata fixtures (8 합성 JSON)

- [ ] Step 1: `bonds/testdata/search_bond_info_success.json` (EP1 — `output{}`, 70 fields all-string)
- [ ] Step 2: `bonds/testdata/issue_info_success.json` (EP2 — `output{}`, 69 fields all-string)
- [ ] Step 3: `bonds/testdata/bond_price_success.json` (EP3 — `output{}`, 17 fields typed)
- [ ] Step 4: `bonds/testdata/bond_ccnl_success.json` (EP4 — `output{}`, 7 fields typed)
- [ ] Step 5: `bonds/testdata/bond_asking_price_success.json` (EP5 — `output{}`, 34 fields typed)
- [ ] Step 6: `bonds/testdata/bond_daily_price_success.json` (EP6 — `output{}`, 9 fields typed)
- [ ] Step 7: `bonds/testdata/bond_daily_itemchartprice_success.json` (EP7 — `output[]`, 6 fields/item)
- [ ] Step 8: `bonds/testdata/avg_unit_success.json` (EP8 — `output1{}+output2[]+output3[]`, 23+10+16 fields)
- [ ] Step 9: validation

```bash
for f in \
  bonds/testdata/search_bond_info_success.json \
  bonds/testdata/issue_info_success.json \
  bonds/testdata/bond_price_success.json \
  bonds/testdata/bond_ccnl_success.json \
  bonds/testdata/bond_asking_price_success.json \
  bonds/testdata/bond_daily_price_success.json \
  bonds/testdata/bond_daily_itemchartprice_success.json \
  bonds/testdata/avg_unit_success.json; do
  python3 -c "import json; json.load(open('$f'))" && echo "$f OK"
done
# Expected: 8 OK lines
```

- [ ] Step 10: commit

```bash
git commit -m "$(cat <<'EOF'
[chore] testdata — 8 bonds quote fixture JSON (Phase 3.1)

합성 JSON fixtures:
- search_bond_info_success.json (output{} 70 fields all-string)
- issue_info_success.json (output{} 69 fields all-string)
- bond_price_success.json (output{} 17 fields typed)
- bond_ccnl_success.json (output{} 7 fields typed)
- bond_asking_price_success.json (output{} 34 fields typed)
- bond_daily_price_success.json (output{} 9 fields typed)
- bond_daily_itemchartprice_success.json (output[] 6 fields/item)
- avg_unit_success.json (output1{} 23 + output2[] 10 + output3[] 16)

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

### Fixture content

**Step 1 — `search_bond_info_success.json`** (EP1: output{} 70 fields all-string)

```json
{
  "rt_cd": "0",
  "msg_cd": "MCA00000",
  "msg1": "정상처리 되었습니다.",
  "output": {
    "pdno": "KR103501GCC7",
    "prdt_type_cd": "300",
    "ksd_bond_item_name": "국고채권 03500-5103(21-3)",
    "ksd_bond_item_eng_name": "KTB 03500-5103(21-3)",
    "ksd_bond_lstg_type_cd": "1",
    "ksd_ofrg_dvsn_cd": "1",
    "ksd_bond_int_dfrm_dvsn_cd": "4",
    "issu_dt": "20210310",
    "rdpt_dt": "20510310",
    "rvnu_dt": "20510310",
    "iso_crcy_cd": "KRW",
    "mdwy_rdpt_dt": "",
    "ksd_rcvg_bond_dsct_rt": "0.00",
    "ksd_rcvg_bond_srfc_inrt": "3.50",
    "bond_expd_rdpt_rt": "0.00",
    "ksd_prca_rdpt_mthd_cd": "1",
    "int_caltm_mcnt": "6",
    "ksd_int_calc_unit_cd": "1",
    "uval_cut_dvsn_cd": "1",
    "uval_cut_dcpt_dgit": "2",
    "ksd_dydv_caltm_aply_dvsn_cd": "1",
    "dydv_calc_dcnt": "0",
    "bond_expd_asrc_erng_rt": "3.50",
    "padf_plac_hdof_name": "한국은행",
    "lstg_dt": "20210310",
    "lstg_abol_dt": "",
    "ksd_bond_issu_mthd_cd": "1",
    "laps_indf_yn": "N",
    "ksd_lhdy_pnia_dfrm_mthd_cd": "1",
    "frst_int_dfrm_dt": "20210910",
    "ksd_prcm_lnkg_gvbd_yn": "N",
    "dpsi_end_dt": "",
    "dpsi_strt_dt": "",
    "dpsi_psbl_yn": "Y",
    "atyp_rdpt_bond_erlm_yn": "N",
    "dshn_occr_yn": "N",
    "expd_exts_yn": "N",
    "pclr_ptcr_text": "",
    "dpsi_psbl_excp_stat_cd": "",
    "expd_exts_srdp_rcnt": "0",
    "expd_exts_srdp_rt": "0.00",
    "expd_rdpt_rt": "0.00",
    "expd_asrc_erng_rt": "3.50",
    "bond_int_dfrm_mthd_cd": "1",
    "int_dfrm_day_type_cd": "1",
    "prca_dfmt_term_mcnt": "0",
    "splt_rdpt_rcnt": "0",
    "rgbf_int_dfrm_dt": "20240910",
    "nxtm_int_dfrm_dt": "20250310",
    "sprx_psbl_yn": "N",
    "ictx_rt_dvsn_cd": "1",
    "bond_clsf_cd": "1",
    "bond_clsf_kor_name": "국채",
    "int_mned_dvsn_cd": "1",
    "pnia_int_calc_unpr": "10000",
    "frn_intr": "0.00",
    "aply_day_prcm_idx_lnkg_cefc": "0.00",
    "ksd_expd_dydv_calc_bass_cd": "1",
    "expd_dydv_calc_dcnt": "0",
    "ksd_cbbw_dvsn_cd": "0",
    "crfd_item_yn": "N",
    "pnia_bank_ofdy_dfrm_mthd_cd": "1",
    "qib_yn": "N",
    "qib_cclc_dt": "",
    "csbd_yn": "N",
    "csbd_cclc_dt": "",
    "ksd_opcb_yn": "N",
    "ksd_sodn_yn": "N",
    "ksd_rqdi_scty_yn": "N",
    "elec_scty_yn": "Y",
    "rght_ecis_mbdy_dvsn_cd": "0",
    "int_rkng_mthd_dvsn_cd": "1",
    "ofrg_dvsn_cd": "1",
    "ksd_tot_issu_amt": "22600000000000"
  }
}
```

**Step 2 — `issue_info_success.json`** (EP2: output{} 69 fields all-string)

```json
{
  "rt_cd": "0",
  "msg_cd": "MCA00000",
  "msg1": "정상처리 되었습니다.",
  "output": {
    "pdno": "KR103501GCC7",
    "prdt_type_cd": "300",
    "prdt_name": "국고채권 03500-5103(21-3)",
    "prdt_eng_name": "KTB 03500-5103(21-3)",
    "ivst_heed_prdt_yn": "N",
    "exts_yn": "N",
    "bond_clsf_cd": "1",
    "bond_clsf_kor_name": "국채",
    "papr": "10000",
    "int_mned_dvsn_cd": "1",
    "rvnu_shap_cd": "1",
    "issu_amt": "22600000000000",
    "lstg_rmnd": "22600000000000",
    "int_dfrm_mcnt": "6",
    "bond_int_dfrm_mthd_cd": "1",
    "splt_rdpt_rcnt": "0",
    "prca_dfmt_term_mcnt": "0",
    "int_anap_dvsn_cd": "1",
    "bond_rght_dvsn_cd": "0",
    "prdt_pclc_text": "",
    "prdt_abrv_name": "국고03500-5103",
    "prdt_eng_abrv_name": "KTB03500-5103",
    "sprx_psbl_yn": "N",
    "pbff_pplc_ofrg_mthd_cd": "1",
    "cmco_cd": "",
    "issu_istt_cd": "B001",
    "issu_istt_name": "기획재정부",
    "pnia_dfrm_agcy_istt_cd": "B001",
    "dsct_ec_rt": "0.00",
    "srfc_inrt": "3.50",
    "expd_rdpt_rt": "0.00",
    "expd_asrc_erng_rt": "3.50",
    "bond_grte_istt_name": "",
    "int_dfrm_day_type_cd": "1",
    "ksd_int_calc_unit_cd": "1",
    "int_wunt_uder_prcs_dvsn_cd": "1",
    "rvnu_dt": "20510310",
    "issu_dt": "20210310",
    "lstg_dt": "20210310",
    "expd_dt": "20510310",
    "rdpt_dt": "20510310",
    "sbst_pric": "10000",
    "rgbf_int_dfrm_dt": "20240910",
    "nxtm_int_dfrm_dt": "20250310",
    "frst_int_dfrm_dt": "20210910",
    "ecis_pric": "",
    "rght_stck_std_pdno": "",
    "ecis_opng_dt": "",
    "ecis_end_dt": "",
    "bond_rvnu_mthd_cd": "1",
    "oprt_stfno": "",
    "oprt_stff_name": "",
    "rgbf_int_dfrm_wday": "4",
    "nxtm_int_dfrm_wday": "1",
    "kis_crdt_grad_text": "AAA",
    "kbp_crdt_grad_text": "AAA",
    "nice_crdt_grad_text": "AAA",
    "fnp_crdt_grad_text": "AAA",
    "dpsi_psbl_yn": "Y",
    "pnia_int_calc_unpr": "10000",
    "prcm_idx_bond_yn": "N",
    "expd_exts_srdp_rcnt": "0",
    "expd_exts_srdp_rt": "0.00",
    "loan_psbl_yn": "Y",
    "grte_dvsn_cd": "0",
    "fnrr_rank_dvsn_cd": "1",
    "krx_lstg_abol_dvsn_cd": "0",
    "asst_rqdi_dvsn_cd": "0",
    "opcb_dvsn_cd": "0",
    "crfd_item_yn": "N",
    "crfd_item_rstc_cclc_dt": "",
    "bond_nmpr_unit_pric": "10000",
    "ivst_heed_bond_dvsn_name": "",
    "add_erng_rt": "0.00",
    "add_erng_rt_aply_dt": "",
    "bond_tr_stop_dvsn_cd": "0",
    "ivst_heed_bond_dvsn_cd": "0",
    "pclr_cndt_text": "",
    "hbbd_yn": "N",
    "cdtl_cptl_scty_type_cd": "",
    "elec_scty_yn": "Y",
    "sq1_clop_ecis_opng_dt": "",
    "frst_erlm_stfno": "",
    "frst_erlm_dt": "20210310",
    "frst_erlm_tmd": "090000",
    "tlg_rcvg_dtl_dtime": "20260505090001123456"
  }
}
```

**Step 3 — `bond_price_success.json`** (EP3: output{} 17 fields typed)

```json
{
  "rt_cd": "0",
  "msg_cd": "MCA00000",
  "msg1": "정상처리 되었습니다.",
  "output": {
    "stnd_iscd": "KR103501GCC7",
    "hts_kor_isnm": "국고채권 03500-5103(21-3)",
    "bond_prpr": "9823.50",
    "prdy_vrss_sign": "5",
    "bond_prdy_vrss": "-12.50",
    "prdy_ctrt": "-0.13",
    "acml_vol": "152000",
    "bond_prdy_clpr": "9836.00",
    "bond_oprc": "9825.00",
    "bond_hgpr": "9830.00",
    "bond_lwpr": "9820.00",
    "ernn_rate": "3.62",
    "oprc_ert": "3.59",
    "hgpr_ert": "3.56",
    "lwpr_ert": "3.65",
    "bond_mxpr": "9950.00",
    "bond_llam": "9700.00"
  }
}
```

**Step 4 — `bond_ccnl_success.json`** (EP4: output{} 7 fields typed, single snapshot)

```json
{
  "rt_cd": "0",
  "msg_cd": "MCA00000",
  "msg1": "정상처리 되었습니다.",
  "output": {
    "stck_cntg_hour": "141523",
    "bond_prpr": "9823.50",
    "bond_prdy_vrss": "-12.50",
    "prdy_vrss_sign": "5",
    "prdy_ctrt": "-0.13",
    "cntg_vol": "500",
    "acml_vol": "152000"
  }
}
```

**Step 5 — `bond_asking_price_success.json`** (EP5: output{} 34 fields typed, 5단계 호가)

```json
{
  "rt_cd": "0",
  "msg_cd": "MCA00000",
  "msg1": "정상처리 되었습니다.",
  "output": {
    "aspr_acpt_hour": "141600",
    "bond_askp1": "9824.00",
    "bond_askp2": "9825.00",
    "bond_askp3": "9826.00",
    "bond_askp4": "9827.00",
    "bond_askp5": "9828.00",
    "bond_bidp1": "9823.00",
    "bond_bidp2": "9822.00",
    "bond_bidp3": "9821.00",
    "bond_bidp4": "9820.00",
    "bond_bidp5": "9819.00",
    "askp_rsqn1": "10000",
    "askp_rsqn2": "8000",
    "askp_rsqn3": "5000",
    "askp_rsqn4": "3000",
    "askp_rsqn5": "2000",
    "bidp_rsqn1": "9000",
    "bidp_rsqn2": "7000",
    "bidp_rsqn3": "4500",
    "bidp_rsqn4": "2500",
    "bidp_rsqn5": "1500",
    "total_askp_rsqn": "28000",
    "total_bidp_rsqn": "24500",
    "ntby_aspr_rsqn": "3500",
    "seln_ernn_rate1": "3.61",
    "seln_ernn_rate2": "3.60",
    "seln_ernn_rate3": "3.59",
    "seln_ernn_rate4": "3.58",
    "seln_ernn_rate5": "3.57",
    "shnu_ernn_rate1": "3.62",
    "shnu_ernn_rate2": "3.63",
    "shnu_ernn_rate3": "3.64",
    "shnu_ernn_rate4": "3.65",
    "shnu_ernn_rate5": "3.66"
  }
}
```

**Step 6 — `bond_daily_price_success.json`** (EP6: output{} 9 fields typed)

```json
{
  "rt_cd": "0",
  "msg_cd": "MCA00000",
  "msg1": "정상처리 되었습니다.",
  "output": {
    "stck_bsop_date": "20260505",
    "bond_prpr": "9823.50",
    "bond_prdy_vrss": "-12.50",
    "prdy_vrss_sign": "5",
    "prdy_ctrt": "-0.13",
    "acml_vol": "152000",
    "bond_oprc": "9825.00",
    "bond_hgpr": "9830.00",
    "bond_lwpr": "9820.00"
  }
}
```

**Step 7 — `bond_daily_itemchartprice_success.json`** (EP7: output[] 6 fields/item, 2 records)

```json
{
  "rt_cd": "0",
  "msg_cd": "MCA00000",
  "msg1": "정상처리 되었습니다.",
  "output": [
    {
      "stck_bsop_date": "20260505",
      "bond_oprc": "9825.00",
      "bond_hgpr": "9830.00",
      "bond_lwpr": "9820.00",
      "bond_prpr": "9823.50",
      "acml_vol": "152000"
    },
    {
      "stck_bsop_date": "20260502",
      "bond_oprc": "9830.00",
      "bond_hgpr": "9840.00",
      "bond_lwpr": "9825.00",
      "bond_prpr": "9836.00",
      "acml_vol": "98000"
    }
  ]
}
```

**Step 8 — `avg_unit_success.json`** (EP8: output1{} 23 fields + output2[]{} 10 fields + output3[]{} 16 fields)

```json
{
  "rt_cd": "0",
  "msg_cd": "MCA00000",
  "msg1": "정상처리 되었습니다.",
  "output1": {
    "pdno": "KR103501GCC7",
    "prdt_type_cd": "300",
    "prdt_name": "국고채권 03500-5103(21-3)",
    "bond_prpr": "9823.50",
    "ernn_rate": "3.62",
    "bond_hgpr": "9830.00",
    "bond_lwpr": "9820.00",
    "kis_unpr": "9820.00",
    "kbp_unpr": "9821.50",
    "nice_evlu_unpr": "9822.00",
    "fnp_unpr": "9823.00",
    "avg_evlu_unpr": "9821.63",
    "kis_evlu_erng_rt": "3.65",
    "kbp_evlu_erng_rt": "3.64",
    "nice_evlu_erng_rt": "3.63",
    "fnp_evlu_erng_rt": "3.62",
    "avg_evlu_erng_rt": "3.64",
    "evlu_dt": "20260505",
    "crcy_cd": "KRW",
    "pdno_cnt": "1",
    "tot_evlu_amt": "9821630",
    "acml_vol": "1000",
    "chng_yn": "Y"
  },
  "output2": [
    {
      "pdno": "KR103501GCC7",
      "prdt_type_cd": "300",
      "prdt_name": "국고채권 03500-5103(21-3)",
      "hldg_qty": "1000",
      "avg_unpr": "9810.00",
      "evlu_pric": "9821.63",
      "evlu_amt": "9821630",
      "evlu_pfls_amt": "11630",
      "crcy_cd": "KRW",
      "chng_yn": "Y"
    }
  ],
  "output3": [
    {
      "pdno": "KR103501GCC7",
      "prdt_type_cd": "300",
      "prdt_name": "국고채권 03500-5103(21-3)",
      "hldg_qty": "1000",
      "avg_unpr": "9810.00",
      "kis_rf_unpr": "9820.00",
      "kbp_rf_unpr": "9821.50",
      "nice_rf_unpr": "9822.00",
      "fnp_rf_unpr": "9823.00",
      "avg_rf_unpr": "9821.63",
      "kis_evlu_unit_pric": "9820.00",
      "kbp_evlu_unit_pric": "9821.50",
      "nice_evlu_unit_pric": "9822.00",
      "fnp_evlu_unit_pric": "9823.00",
      "avg_evlu_unit_pric": "9821.63",
      "evlu_amt": "9821630"
    }
  ]
}
```

---

## Task 2: bonds 패키지 scaffolding + SearchBondInfo (EP1)

**Files:**
- CREATE `bonds/client.go`
- CREATE `bonds/doc.go`
- CREATE `bonds/testhelper_test.go`
- CREATE `bonds/quote.go` (첫 번째 메서드 + struct)
- CREATE `bonds/quote_test.go` (첫 번째 테스트)
- MODIFY root `client.go` (Bonds 필드 + wireInfra)

- [ ] Step 1: CREATE `bonds/client.go`

```go
package bonds

import (
	"github.com/kenshin579/korea-investment-stock/internal/httpclient"
)

// Client 는 장내채권 (Korean bond) API 클라이언트.
type Client struct {
	http *httpclient.Client
}

// New 는 채권 Client 생성. http 는 root client 의 internal httpclient.
func New(http *httpclient.Client) *Client {
	return &Client{http: http}
}
```

- [ ] Step 2: CREATE `bonds/doc.go`

```go
// Package bonds 는 장내채권 (Korean bond) API 클라이언트.
//
// Phase 3.1 메서드 (8):
//
//   - SearchBondInfo             — 채권 기본조회 (CTPF1114R, 70 fields)
//   - InquireIssueInfo           — 발행정보 (CTPF1101R, 69 fields)
//   - InquirePrice               — 현재가 시세 (FHKBJ773400C0)
//   - InquireCcnl                — 현재가 체결 (FHKBJ773403C0)
//   - InquireAskingPrice         — 현재가 호가 (FHKBJ773401C0, 5단계)
//   - InquireDailyPrice          — 현재가 일별 (FHKBJ773404C0)
//   - InquireDailyItemchartprice — 기간별 시세 (FHKBJ773701C0)
//   - InquireAvgUnit             — 평균단가조회 (CTPF2005R)
//
// 사용자는 root kis.Client.Bonds 로 접근.
package bonds
```

- [ ] Step 3: CREATE `bonds/testhelper_test.go`

```go
package bonds_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/kenshin579/korea-investment-stock/bonds"
	"github.com/kenshin579/korea-investment-stock/internal/httpclient"
)

func newTestClient(t *testing.T) (*bonds.Client, *httpmock.MockTransport) {
	t.Helper()
	transport := httpmock.NewMockTransport()
	hc := httpclient.NewForTest(transport)
	return bonds.New(hc), transport
}

func loadFixtureString(t *testing.T, name string) string {
	t.Helper()
	b, err := os.ReadFile(filepath.Join("testdata", name))
	if err != nil {
		t.Fatalf("loadFixtureString: %v", err)
	}
	return string(b)
}
```

- [ ] Step 4: APPEND failing test to CREATE `bonds/quote_test.go`

```go
package bonds_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_SearchBondInfo(t *testing.T) {
	client, transport := newTestClient(t)
	transport.RegisterResponder(
		http.MethodGet,
		"=~/uapi/domestic-bond/v1/quotations/search-bond-info",
		httpmock.NewStringResponder(http.StatusOK, loadFixtureString(t, "search_bond_info_success.json")),
	)

	got, err := client.SearchBondInfo(context.Background(), SearchBondInfoParams{
		Pdno:       "KR103501GCC7",
		PrdtTypeCd: "300",
	})
	require.NoError(t, err)
	require.NotNil(t, got)

	assert.Equal(t, "KR103501GCC7", got.Pdno)
	assert.Equal(t, "국고채권 03500-5103(21-3)", got.KsdBondItemName)
	assert.Equal(t, "3.50", got.KsdRcvgBondSrfcInrt)
	assert.Equal(t, "국채", got.BondClsfKorName)
	assert.Equal(t, "Y", got.ElecSctyYn)
}
```

- [ ] Step 5: Verify FAIL

```bash
cd /Users/user/src/workspace_moneyflow/korea-investment-stock && \
  go test ./bonds/... 2>&1 | head -20
# Expected: compile error — SearchBondInfo undefined
```

- [ ] Step 6: CREATE `bonds/quote.go` with SearchBondInfo struct + Params + method

```go
package bonds

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kenshin579/korea-investment-stock/internal/httpclient"
)

// ─── EP1: SearchBondInfo ──────────────────────────────────────────────────────

// SearchBondInfoParams 는 채권 기본조회 요청 파라미터.
type SearchBondInfoParams struct {
	Pdno       string // PDNO: 채권 종목 번호 (필수)
	PrdtTypeCd string // PRDT_TYPE_CD: 상품유형코드 (필수, 예: "300")
}

// SearchBondInfo 는 채권 기본조회 결과. CTPF1114R — all-string 70 fields.
type SearchBondInfo struct {
	Pdno                       string `json:"pdno"`
	PrdtTypeCd                 string `json:"prdt_type_cd"`
	KsdBondItemName            string `json:"ksd_bond_item_name"`
	KsdBondItemEngName         string `json:"ksd_bond_item_eng_name"`
	KsdBondLstgTypeCd          string `json:"ksd_bond_lstg_type_cd"`
	KsdOfrgDvsnCd              string `json:"ksd_ofrg_dvsn_cd"`
	KsdBondIntDfrmDvsnCd       string `json:"ksd_bond_int_dfrm_dvsn_cd"`
	IssuDt                     string `json:"issu_dt"`
	RdptDt                     string `json:"rdpt_dt"`
	RvnuDt                     string `json:"rvnu_dt"`
	IsoCrcyCd                  string `json:"iso_crcy_cd"`
	MdwyRdptDt                 string `json:"mdwy_rdpt_dt"`
	KsdRcvgBondDsctRt          string `json:"ksd_rcvg_bond_dsct_rt"`
	KsdRcvgBondSrfcInrt        string `json:"ksd_rcvg_bond_srfc_inrt"`
	BondExpdRdptRt             string `json:"bond_expd_rdpt_rt"`
	KsdPrcaRdptMthdCd          string `json:"ksd_prca_rdpt_mthd_cd"`
	IntCaltmMcnt               string `json:"int_caltm_mcnt"`
	KsdIntCalcUnitCd           string `json:"ksd_int_calc_unit_cd"`
	UvalCutDvsnCd              string `json:"uval_cut_dvsn_cd"`
	UvalCutDcptDgit            string `json:"uval_cut_dcpt_dgit"`
	KsdDydvCaltmAplyDvsnCd     string `json:"ksd_dydv_caltm_aply_dvsn_cd"`
	DydvCalcDcnt               string `json:"dydv_calc_dcnt"`
	BondExpdAsrcErngRt         string `json:"bond_expd_asrc_erng_rt"`
	PadfPlacHdofName           string `json:"padf_plac_hdof_name"`
	LstgDt                     string `json:"lstg_dt"`
	LstgAbolDt                 string `json:"lstg_abol_dt"`
	KsdBondIssuMthdCd          string `json:"ksd_bond_issu_mthd_cd"`
	LapsIndfYn                 string `json:"laps_indf_yn"`
	KsdLhdyPniaDfrmMthdCd      string `json:"ksd_lhdy_pnia_dfrm_mthd_cd"`
	FrstIntDfrmDt              string `json:"frst_int_dfrm_dt"`
	KsdPrcmLnkgGvbdYn          string `json:"ksd_prcm_lnkg_gvbd_yn"`
	DpsiEndDt                  string `json:"dpsi_end_dt"`
	DpsiStrtDt                 string `json:"dpsi_strt_dt"`
	DpsiPsblYn                 string `json:"dpsi_psbl_yn"`
	AtypRdptBondErlmYn         string `json:"atyp_rdpt_bond_erlm_yn"`
	DshnOccrYn                 string `json:"dshn_occr_yn"`
	ExpdExtsYn                 string `json:"expd_exts_yn"`
	PclrPtcrText               string `json:"pclr_ptcr_text"`
	DpsiPsblExcpStatCd         string `json:"dpsi_psbl_excp_stat_cd"`
	ExpdExtsSrdpRcnt           string `json:"expd_exts_srdp_rcnt"`
	ExpdExtsSrdpRt             string `json:"expd_exts_srdp_rt"`
	ExpdRdptRt                 string `json:"expd_rdpt_rt"`
	ExpdAsrcErngRt             string `json:"expd_asrc_erng_rt"`
	BondIntDfrmMthdCd          string `json:"bond_int_dfrm_mthd_cd"`
	IntDfrmDayTypeCd           string `json:"int_dfrm_day_type_cd"`
	PrcaDfmtTermMcnt           string `json:"prca_dfmt_term_mcnt"`
	SpltRdptRcnt               string `json:"splt_rdpt_rcnt"`
	RgbfIntDfrmDt              string `json:"rgbf_int_dfrm_dt"`
	NxtmIntDfrmDt              string `json:"nxtm_int_dfrm_dt"`
	SprxPsblYn                 string `json:"sprx_psbl_yn"`
	IctxRtDvsnCd               string `json:"ictx_rt_dvsn_cd"`
	BondClsfCd                 string `json:"bond_clsf_cd"`
	BondClsfKorName            string `json:"bond_clsf_kor_name"`
	IntMnedDvsnCd              string `json:"int_mned_dvsn_cd"`
	PniaIntCalcUnpr            string `json:"pnia_int_calc_unpr"`
	FrnIntr                    string `json:"frn_intr"`
	AplyDayPrcmIdxLnkgCefc     string `json:"aply_day_prcm_idx_lnkg_cefc"`
	KsdExpdDydvCalcBassCd      string `json:"ksd_expd_dydv_calc_bass_cd"`
	ExpdDydvCalcDcnt           string `json:"expd_dydv_calc_dcnt"`
	KsdCbbwDvsnCd              string `json:"ksd_cbbw_dvsn_cd"`
	CrfdItemYn                 string `json:"crfd_item_yn"`
	PniaBankOfdyDfrmMthdCd     string `json:"pnia_bank_ofdy_dfrm_mthd_cd"`
	QibYn                      string `json:"qib_yn"`
	QibCclcDt                  string `json:"qib_cclc_dt"`
	CsbdYn                     string `json:"csbd_yn"`
	CsbdCclcDt                 string `json:"csbd_cclc_dt"`
	KsdOpcbYn                  string `json:"ksd_opcb_yn"`
	KsdSodnYn                  string `json:"ksd_sodn_yn"`
	KsdRqdiSctyYn              string `json:"ksd_rqdi_scty_yn"`
	ElecSctyYn                 string `json:"elec_scty_yn"`
	RghtEcisMbdyDvsnCd         string `json:"rght_ecis_mbdy_dvsn_cd"`
	IntRkngMthdDvsnCd          string `json:"int_rkng_mthd_dvsn_cd"`
	OfrgDvsnCd                 string `json:"ofrg_dvsn_cd"`
	KsdTotIssuAmt              string `json:"ksd_tot_issu_amt"`
}

type searchBondInfoResponse struct {
	RtCd  string         `json:"rt_cd"`
	MsgCd string         `json:"msg_cd"`
	Msg1  string         `json:"msg1"`
	Output SearchBondInfo `json:"output"`
}

// SearchBondInfo 는 채권 기본조회 (CTPF1114R).
//
// KIS API: GET /uapi/domestic-bond/v1/quotations/search-bond-info
func (c *Client) SearchBondInfo(ctx context.Context, params SearchBondInfoParams) (*SearchBondInfo, error) {
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method:   http.MethodGet,
		Path:     "/uapi/domestic-bond/v1/quotations/search-bond-info",
		TrID:     "CTPF1114R",
		CustType: "P",
		Query: map[string]string{
			"PDNO":         params.Pdno,
			"PRDT_TYPE_CD": params.PrdtTypeCd,
		},
	})
	if err != nil {
		return nil, err
	}
	var res searchBondInfoResponse
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse SearchBondInfo: %w", err)
	}
	return &res.Output, nil
}
```

- [ ] Step 7: Verify PASS

```bash
cd /Users/user/src/workspace_moneyflow/korea-investment-stock && \
  go test ./bonds/... -run TestClient_SearchBondInfo -v 2>&1
# Expected: PASS
```

- [ ] Step 8: MODIFY root `client.go` — add Bonds field + wireInfra

  Add `"github.com/kenshin579/korea-investment-stock/bonds"` to import block.

  In `type Client struct { ... }`, add after `Overseas *overseas.Client`:
  ```go
  Bonds    *bonds.Client
  ```

  In `NewClient` function body, add after `c.Overseas = overseas.New(c.httpClient, c.masterC)`:
  ```go
  c.Bonds = bonds.New(c.httpClient)
  ```

- [ ] Step 9: gofmt + vet

```bash
cd /Users/user/src/workspace_moneyflow/korea-investment-stock && \
  gofmt -w bonds/client.go bonds/doc.go bonds/quote.go bonds/quote_test.go bonds/testhelper_test.go client.go && \
  go vet ./... 2>&1
# Expected: no output (clean)
```

- [ ] Step 10: commit

```bash
git commit -m "$(cat <<'EOF'
[feat] bonds: scaffold package + SearchBondInfo EP1 (Phase 3.1)

신규 bonds/ sub-package 도입:
- bonds/client.go: Client struct + New constructor
- bonds/doc.go: package doc (8 메서드 목록)
- bonds/testhelper_test.go: newTestClient + loadFixtureString helpers
- bonds/quote.go: SearchBondInfo (CTPF1114R, 70 fields all-string)
- bonds/quote_test.go: TestClient_SearchBondInfo

Root client.go: Bonds *bonds.Client 필드 + wireInfra 주입.

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 3: InquireIssueInfo (EP2) — APPEND

**Files:** APPEND to `bonds/quote.go` and `bonds/quote_test.go`

- [ ] Step 1: APPEND test to `bonds/quote_test.go`

```go
func TestClient_InquireIssueInfo(t *testing.T) {
	client, transport := newTestClient(t)
	transport.RegisterResponder(
		http.MethodGet,
		"=~/uapi/domestic-bond/v1/quotations/issue-info",
		httpmock.NewStringResponder(http.StatusOK, loadFixtureString(t, "issue_info_success.json")),
	)

	got, err := client.InquireIssueInfo(context.Background(), InquireIssueInfoParams{
		Pdno:       "KR103501GCC7",
		PrdtTypeCd: "300",
	})
	require.NoError(t, err)
	require.NotNil(t, got)

	assert.Equal(t, "KR103501GCC7", got.Pdno)
	assert.Equal(t, "국고채권 03500-5103(21-3)", got.PrdtName)
	assert.Equal(t, "3.50", got.SrfcInrt)
	assert.Equal(t, "AAA", got.KisCrdtGradText)
	assert.Equal(t, "Y", got.ElecSctyYn)
}
```

- [ ] Step 2: Verify FAIL

```bash
cd /Users/user/src/workspace_moneyflow/korea-investment-stock && \
  go test ./bonds/... -run TestClient_InquireIssueInfo -v 2>&1 | head -20
# Expected: compile error — InquireIssueInfo undefined
```

- [ ] Step 3: APPEND struct + Params + method to `bonds/quote.go`

```go
// ─── EP2: InquireIssueInfo ────────────────────────────────────────────────────

// InquireIssueInfoParams 는 발행정보 조회 요청 파라미터.
type InquireIssueInfoParams struct {
	Pdno       string // PDNO: 채권 종목 번호 (필수)
	PrdtTypeCd string // PRDT_TYPE_CD: 상품유형코드 (필수)
}

// IssueInfo 는 채권 발행정보. CTPF1101R — all-string 69 fields.
type IssueInfo struct {
	Pdno                  string `json:"pdno"`
	PrdtTypeCd            string `json:"prdt_type_cd"`
	PrdtName              string `json:"prdt_name"`
	PrdtEngName           string `json:"prdt_eng_name"`
	IvstHeedPrdtYn        string `json:"ivst_heed_prdt_yn"`
	ExtsYn                string `json:"exts_yn"`
	BondClsfCd            string `json:"bond_clsf_cd"`
	BondClsfKorName       string `json:"bond_clsf_kor_name"`
	Papr                  string `json:"papr"`
	IntMnedDvsnCd         string `json:"int_mned_dvsn_cd"`
	RvnuShapCd            string `json:"rvnu_shap_cd"`
	IssuAmt               string `json:"issu_amt"`
	LstgRmnd              string `json:"lstg_rmnd"`
	IntDfrmMcnt           string `json:"int_dfrm_mcnt"`
	BondIntDfrmMthdCd     string `json:"bond_int_dfrm_mthd_cd"`
	SpltRdptRcnt          string `json:"splt_rdpt_rcnt"`
	PrcaDfmtTermMcnt      string `json:"prca_dfmt_term_mcnt"`
	IntAnapDvsnCd         string `json:"int_anap_dvsn_cd"`
	BondRghtDvsnCd        string `json:"bond_rght_dvsn_cd"`
	PrdtPclcText          string `json:"prdt_pclc_text"`
	PrdtAbrvName          string `json:"prdt_abrv_name"`
	PrdtEngAbrvName       string `json:"prdt_eng_abrv_name"`
	SprxPsblYn            string `json:"sprx_psbl_yn"`
	PbffPplcOfrgMthdCd    string `json:"pbff_pplc_ofrg_mthd_cd"`
	CmcoCd                string `json:"cmco_cd"`
	IssuIsttCd            string `json:"issu_istt_cd"`
	IssuIsttName          string `json:"issu_istt_name"`
	PniaDfrmAgcyIsttCd    string `json:"pnia_dfrm_agcy_istt_cd"`
	DsctEcRt              string `json:"dsct_ec_rt"`
	SrfcInrt              string `json:"srfc_inrt"`
	ExpdRdptRt            string `json:"expd_rdpt_rt"`
	ExpdAsrcErngRt        string `json:"expd_asrc_erng_rt"`
	BondGrteIsttName      string `json:"bond_grte_istt_name"`
	IntDfrmDayTypeCd      string `json:"int_dfrm_day_type_cd"`
	KsdIntCalcUnitCd      string `json:"ksd_int_calc_unit_cd"`
	IntWuntUderPrcsDvsnCd string `json:"int_wunt_uder_prcs_dvsn_cd"`
	RvnuDt                string `json:"rvnu_dt"`
	IssuDt                string `json:"issu_dt"`
	LstgDt                string `json:"lstg_dt"`
	ExpdDt                string `json:"expd_dt"`
	RdptDt                string `json:"rdpt_dt"`
	SbstPric              string `json:"sbst_pric"`
	RgbfIntDfrmDt         string `json:"rgbf_int_dfrm_dt"`
	NxtmIntDfrmDt         string `json:"nxtm_int_dfrm_dt"`
	FrstIntDfrmDt         string `json:"frst_int_dfrm_dt"`
	EcisPric              string `json:"ecis_pric"`
	RghtStckStdPdno       string `json:"rght_stck_std_pdno"`
	EcisOpngDt            string `json:"ecis_opng_dt"`
	EcisEndDt             string `json:"ecis_end_dt"`
	BondRvnuMthdCd        string `json:"bond_rvnu_mthd_cd"`
	OprtStfno             string `json:"oprt_stfno"`
	OprtStffName          string `json:"oprt_stff_name"`
	RgbfIntDfrmWday       string `json:"rgbf_int_dfrm_wday"`
	NxtmIntDfrmWday       string `json:"nxtm_int_dfrm_wday"`
	KisCrdtGradText       string `json:"kis_crdt_grad_text"`
	KbpCrdtGradText       string `json:"kbp_crdt_grad_text"`
	NiceCrdtGradText      string `json:"nice_crdt_grad_text"`
	FnpCrdtGradText       string `json:"fnp_crdt_grad_text"`
	DpsiPsblYn            string `json:"dpsi_psbl_yn"`
	PniaIntCalcUnpr       string `json:"pnia_int_calc_unpr"`
	PrcmIdxBondYn         string `json:"prcm_idx_bond_yn"`
	ExpdExtsSrdpRcnt      string `json:"expd_exts_srdp_rcnt"`
	ExpdExtsSrdpRt        string `json:"expd_exts_srdp_rt"`
	LoanPsblYn            string `json:"loan_psbl_yn"`
	GrteDvsnCd            string `json:"grte_dvsn_cd"`
	FnrrRankDvsnCd        string `json:"fnrr_rank_dvsn_cd"`
	KrxLstgAbolDvsnCd     string `json:"krx_lstg_abol_dvsn_cd"`
	AsstRqdiDvsnCd        string `json:"asst_rqdi_dvsn_cd"`
	OpcbDvsnCd            string `json:"opcb_dvsn_cd"`
	CrfdItemYn            string `json:"crfd_item_yn"`
	CrfdItemRstcCclcDt    string `json:"crfd_item_rstc_cclc_dt"`
	BondNmprUnitPric      string `json:"bond_nmpr_unit_pric"`
	IvstHeedBondDvsnName  string `json:"ivst_heed_bond_dvsn_name"`
	AddErngRt             string `json:"add_erng_rt"`
	AddErngRtAplyDt       string `json:"add_erng_rt_aply_dt"`
	BondTrStopDvsnCd      string `json:"bond_tr_stop_dvsn_cd"`
	IvstHeedBondDvsnCd    string `json:"ivst_heed_bond_dvsn_cd"`
	PclrCndtText          string `json:"pclr_cndt_text"`
	HbbdYn                string `json:"hbbd_yn"`
	CdtlCptlSctyTypeCd    string `json:"cdtl_cptl_scty_type_cd"`
	ElecSctyYn            string `json:"elec_scty_yn"`
	Sq1ClopEcisOpngDt     string `json:"sq1_clop_ecis_opng_dt"`
	FrstErlmStfno         string `json:"frst_erlm_stfno"`
	FrstErlmDt            string `json:"frst_erlm_dt"`
	FrstErlmTmd           string `json:"frst_erlm_tmd"`
	TlgRcvgDtlDtime       string `json:"tlg_rcvg_dtl_dtime"`
}

type inquireIssueInfoResponse struct {
	RtCd   string    `json:"rt_cd"`
	MsgCd  string    `json:"msg_cd"`
	Msg1   string    `json:"msg1"`
	Output IssueInfo `json:"output"`
}

// InquireIssueInfo 는 채권 발행정보 조회 (CTPF1101R).
//
// KIS API: GET /uapi/domestic-bond/v1/quotations/issue-info
func (c *Client) InquireIssueInfo(ctx context.Context, params InquireIssueInfoParams) (*IssueInfo, error) {
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method:   http.MethodGet,
		Path:     "/uapi/domestic-bond/v1/quotations/issue-info",
		TrID:     "CTPF1101R",
		CustType: "P",
		Query: map[string]string{
			"PDNO":         params.Pdno,
			"PRDT_TYPE_CD": params.PrdtTypeCd,
		},
	})
	if err != nil {
		return nil, err
	}
	var res inquireIssueInfoResponse
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse InquireIssueInfo: %w", err)
	}
	return &res.Output, nil
}
```

- [ ] Step 4: Verify PASS

```bash
cd /Users/user/src/workspace_moneyflow/korea-investment-stock && \
  go test ./bonds/... -run TestClient_InquireIssueInfo -v 2>&1
# Expected: PASS
```

- [ ] Step 5: gofmt/vet

```bash
cd /Users/user/src/workspace_moneyflow/korea-investment-stock && \
  gofmt -w bonds/quote.go bonds/quote_test.go && go vet ./bonds/... 2>&1
```

- [ ] Step 6: commit

```bash
git commit -m "$(cat <<'EOF'
[feat] bonds: InquireIssueInfo EP2 (CTPF1101R, 69 fields)

발행정보 조회 — all-string 69 필드 struct + method + test.

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 4: InquirePrice (EP3) — APPEND

**Files:** APPEND to `bonds/quote.go` and `bonds/quote_test.go`

- [ ] Step 1: APPEND test to `bonds/quote_test.go`

```go
func TestClient_InquirePrice(t *testing.T) {
	client, transport := newTestClient(t)
	transport.RegisterResponder(
		http.MethodGet,
		"=~/uapi/domestic-bond/v1/quotations/inquire-price",
		httpmock.NewStringResponder(http.StatusOK, loadFixtureString(t, "bond_price_success.json")),
	)

	got, err := client.InquirePrice(context.Background(), InquirePriceParams{
		MarketCode: "B",
		Symbol:     "KR103501GCC7",
	})
	require.NoError(t, err)
	require.NotNil(t, got)

	assert.Equal(t, "KR103501GCC7", got.StndIscd)
	assert.Equal(t, "9823.50", got.BondPrpr.String())
	assert.Equal(t, "-12.50", got.BondPrdyVrss.String())
	assert.InDelta(t, -0.13, got.PrdyCtrt, 0.001)
	assert.Equal(t, int64(152000), got.AcmlVol)
}
```

- [ ] Step 2: Verify FAIL

```bash
cd /Users/user/src/workspace_moneyflow/korea-investment-stock && \
  go test ./bonds/... -run TestClient_InquirePrice -v 2>&1 | head -20
# Expected: compile error — InquirePrice undefined
```

- [ ] Step 3: APPEND struct + Params + method to `bonds/quote.go`

```go
// ─── EP3: InquirePrice ────────────────────────────────────────────────────────

// InquirePriceParams 는 채권 현재가 시세 요청 파라미터.
type InquirePriceParams struct {
	MarketCode string // FID_COND_MRKT_DIV_CODE: 기본 "B"
	Symbol     string // FID_INPUT_ISCD: 채권 단축 종목코드 (필수)
}

// BondPrice 는 채권 현재가 시세. FHKBJ773400C0 — 17 fields typed.
type BondPrice struct {
	StndIscd     string          `json:"stnd_iscd"`
	HtsKorIsnm   string          `json:"hts_kor_isnm"`
	BondPrpr     decimal.Decimal `json:"bond_prpr"`
	PrdyVrssSign string          `json:"prdy_vrss_sign"`
	BondPrdyVrss decimal.Decimal `json:"bond_prdy_vrss"`
	PrdyCtrt     float64         `json:"prdy_ctrt,string"`
	AcmlVol      int64           `json:"acml_vol,string"`
	BondPrdyClpr decimal.Decimal `json:"bond_prdy_clpr"`
	BondOprc     decimal.Decimal `json:"bond_oprc"`
	BondHgpr     decimal.Decimal `json:"bond_hgpr"`
	BondLwpr     decimal.Decimal `json:"bond_lwpr"`
	ErnnRate     float64         `json:"ernn_rate,string"`
	OprcErt      float64         `json:"oprc_ert,string"`
	HgprErt      float64         `json:"hgpr_ert,string"`
	LwprErt      float64         `json:"lwpr_ert,string"`
	BondMxpr     decimal.Decimal `json:"bond_mxpr"`
	BondLlam     decimal.Decimal `json:"bond_llam"`
}

type inquirePriceResponse struct {
	RtCd   string    `json:"rt_cd"`
	MsgCd  string    `json:"msg_cd"`
	Msg1   string    `json:"msg1"`
	Output BondPrice `json:"output"`
}

// InquirePrice 는 채권 현재가 시세 (FHKBJ773400C0).
//
// KIS API: GET /uapi/domestic-bond/v1/quotations/inquire-price
func (c *Client) InquirePrice(ctx context.Context, params InquirePriceParams) (*BondPrice, error) {
	if params.MarketCode == "" {
		params.MarketCode = "B"
	}
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method:   http.MethodGet,
		Path:     "/uapi/domestic-bond/v1/quotations/inquire-price",
		TrID:     "FHKBJ773400C0",
		CustType: "P",
		Query: map[string]string{
			"FID_COND_MRKT_DIV_CODE": params.MarketCode,
			"FID_INPUT_ISCD":         params.Symbol,
		},
	})
	if err != nil {
		return nil, err
	}
	var res inquirePriceResponse
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse InquirePrice: %w", err)
	}
	return &res.Output, nil
}
```

> **Note:** Add `"github.com/shopspring/decimal"` to import block in `bonds/quote.go`.

- [ ] Step 4: Verify PASS

```bash
cd /Users/user/src/workspace_moneyflow/korea-investment-stock && \
  go test ./bonds/... -run TestClient_InquirePrice -v 2>&1
# Expected: PASS
```

- [ ] Step 5: gofmt/vet

```bash
cd /Users/user/src/workspace_moneyflow/korea-investment-stock && \
  gofmt -w bonds/quote.go bonds/quote_test.go && go vet ./bonds/... 2>&1
```

- [ ] Step 6: commit

```bash
git commit -m "$(cat <<'EOF'
[feat] bonds: InquirePrice EP3 (FHKBJ773400C0, 17 fields typed)

채권 현재가 시세 — decimal/float64/int64 typed struct + method + test.
Default FID_COND_MRKT_DIV_CODE="B".

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 5: InquireCcnl (EP4) — APPEND

**Files:** APPEND to `bonds/quote.go` and `bonds/quote_test.go`

- [ ] Step 1: APPEND test to `bonds/quote_test.go`

```go
func TestClient_InquireCcnl(t *testing.T) {
	client, transport := newTestClient(t)
	transport.RegisterResponder(
		http.MethodGet,
		"=~/uapi/domestic-bond/v1/quotations/inquire-ccnl",
		httpmock.NewStringResponder(http.StatusOK, loadFixtureString(t, "bond_ccnl_success.json")),
	)

	got, err := client.InquireCcnl(context.Background(), InquireCcnlParams{
		MarketCode: "B",
		Symbol:     "KR103501GCC7",
	})
	require.NoError(t, err)
	require.NotNil(t, got)

	assert.Equal(t, "141523", got.StckCntgHour)
	assert.Equal(t, "9823.50", got.BondPrpr.String())
	assert.Equal(t, int64(500), got.CntgVol)
	assert.Equal(t, int64(152000), got.AcmlVol)
	assert.InDelta(t, -0.13, got.PrdyCtrt, 0.001)
}
```

- [ ] Step 2: Verify FAIL

```bash
cd /Users/user/src/workspace_moneyflow/korea-investment-stock && \
  go test ./bonds/... -run TestClient_InquireCcnl -v 2>&1 | head -20
# Expected: compile error — InquireCcnl undefined
```

- [ ] Step 3: APPEND struct + Params + method to `bonds/quote.go`

```go
// ─── EP4: InquireCcnl ─────────────────────────────────────────────────────────

// InquireCcnlParams 는 채권 현재가 체결 요청 파라미터.
type InquireCcnlParams struct {
	MarketCode string // FID_COND_MRKT_DIV_CODE: 기본 "B"
	Symbol     string // FID_INPUT_ISCD: 채권 단축 종목코드 (필수)
}

// BondCcnl 는 채권 현재가 체결 (single snapshot). FHKBJ773403C0 — 7 fields typed.
type BondCcnl struct {
	StckCntgHour string          `json:"stck_cntg_hour"`
	BondPrpr     decimal.Decimal `json:"bond_prpr"`
	BondPrdyVrss decimal.Decimal `json:"bond_prdy_vrss"`
	PrdyVrssSign string          `json:"prdy_vrss_sign"`
	PrdyCtrt     float64         `json:"prdy_ctrt,string"`
	CntgVol      int64           `json:"cntg_vol,string"`
	AcmlVol      int64           `json:"acml_vol,string"`
}

type inquireCcnlResponse struct {
	RtCd   string   `json:"rt_cd"`
	MsgCd  string   `json:"msg_cd"`
	Msg1   string   `json:"msg1"`
	Output BondCcnl `json:"output"`
}

// InquireCcnl 는 채권 현재가 체결 (FHKBJ773403C0).
//
// KIS API: GET /uapi/domestic-bond/v1/quotations/inquire-ccnl
func (c *Client) InquireCcnl(ctx context.Context, params InquireCcnlParams) (*BondCcnl, error) {
	if params.MarketCode == "" {
		params.MarketCode = "B"
	}
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method:   http.MethodGet,
		Path:     "/uapi/domestic-bond/v1/quotations/inquire-ccnl",
		TrID:     "FHKBJ773403C0",
		CustType: "P",
		Query: map[string]string{
			"FID_COND_MRKT_DIV_CODE": params.MarketCode,
			"FID_INPUT_ISCD":         params.Symbol,
		},
	})
	if err != nil {
		return nil, err
	}
	var res inquireCcnlResponse
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse InquireCcnl: %w", err)
	}
	return &res.Output, nil
}
```

- [ ] Step 4: Verify PASS

```bash
cd /Users/user/src/workspace_moneyflow/korea-investment-stock && \
  go test ./bonds/... -run TestClient_InquireCcnl -v 2>&1
# Expected: PASS
```

- [ ] Step 5: gofmt/vet

```bash
cd /Users/user/src/workspace_moneyflow/korea-investment-stock && \
  gofmt -w bonds/quote.go bonds/quote_test.go && go vet ./bonds/... 2>&1
```

- [ ] Step 6: commit

```bash
git commit -m "$(cat <<'EOF'
[feat] bonds: InquireCcnl EP4 (FHKBJ773403C0, 7 fields typed)

채권 현재가 체결 단일 스냅샷 — typed struct + method + test.

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 6: InquireAskingPrice (EP5) — APPEND

**Files:** APPEND to `bonds/quote.go` and `bonds/quote_test.go`

- [ ] Step 1: APPEND test to `bonds/quote_test.go`

```go
func TestClient_InquireAskingPrice(t *testing.T) {
	client, transport := newTestClient(t)
	transport.RegisterResponder(
		http.MethodGet,
		"=~/uapi/domestic-bond/v1/quotations/inquire-asking-price",
		httpmock.NewStringResponder(http.StatusOK, loadFixtureString(t, "bond_asking_price_success.json")),
	)

	got, err := client.InquireAskingPrice(context.Background(), InquireAskingPriceParams{
		MarketCode: "B",
		Symbol:     "KR103501GCC7",
	})
	require.NoError(t, err)
	require.NotNil(t, got)

	assert.Equal(t, "141600", got.AsprAcptHour)
	assert.Equal(t, "9824.00", got.BondAskp1.String())
	assert.Equal(t, "9823.00", got.BondBidp1.String())
	assert.Equal(t, int64(10000), got.AskpRsqn1)
	assert.Equal(t, int64(28000), got.TotalAskpRsqn)
	assert.InDelta(t, 3.61, got.SelnErnnRate1, 0.001)
	assert.InDelta(t, 3.62, got.ShnuErnnRate1, 0.001)
}
```

- [ ] Step 2: Verify FAIL

```bash
cd /Users/user/src/workspace_moneyflow/korea-investment-stock && \
  go test ./bonds/... -run TestClient_InquireAskingPrice -v 2>&1 | head -20
# Expected: compile error — InquireAskingPrice undefined
```

- [ ] Step 3: APPEND struct + Params + method to `bonds/quote.go`

```go
// ─── EP5: InquireAskingPrice ──────────────────────────────────────────────────

// InquireAskingPriceParams 는 채권 현재가 호가 요청 파라미터.
type InquireAskingPriceParams struct {
	MarketCode string // FID_COND_MRKT_DIV_CODE: 기본 "B"
	Symbol     string // FID_INPUT_ISCD: 채권 단축 종목코드 (필수)
}

// BondAskingPrice 는 채권 현재가 호가 (5단계). FHKBJ773401C0 — 34 fields typed.
type BondAskingPrice struct {
	AsprAcptHour  string          `json:"aspr_acpt_hour"`
	BondAskp1     decimal.Decimal `json:"bond_askp1"`
	BondAskp2     decimal.Decimal `json:"bond_askp2"`
	BondAskp3     decimal.Decimal `json:"bond_askp3"`
	BondAskp4     decimal.Decimal `json:"bond_askp4"`
	BondAskp5     decimal.Decimal `json:"bond_askp5"`
	BondBidp1     decimal.Decimal `json:"bond_bidp1"`
	BondBidp2     decimal.Decimal `json:"bond_bidp2"`
	BondBidp3     decimal.Decimal `json:"bond_bidp3"`
	BondBidp4     decimal.Decimal `json:"bond_bidp4"`
	BondBidp5     decimal.Decimal `json:"bond_bidp5"`
	AskpRsqn1     int64           `json:"askp_rsqn1,string"`
	AskpRsqn2     int64           `json:"askp_rsqn2,string"`
	AskpRsqn3     int64           `json:"askp_rsqn3,string"`
	AskpRsqn4     int64           `json:"askp_rsqn4,string"`
	AskpRsqn5     int64           `json:"askp_rsqn5,string"`
	BidpRsqn1     int64           `json:"bidp_rsqn1,string"`
	BidpRsqn2     int64           `json:"bidp_rsqn2,string"`
	BidpRsqn3     int64           `json:"bidp_rsqn3,string"`
	BidpRsqn4     int64           `json:"bidp_rsqn4,string"`
	BidpRsqn5     int64           `json:"bidp_rsqn5,string"`
	TotalAskpRsqn int64           `json:"total_askp_rsqn,string"`
	TotalBidpRsqn int64           `json:"total_bidp_rsqn,string"`
	NtbyAsprRsqn  int64           `json:"ntby_aspr_rsqn,string"`
	SelnErnnRate1 float64         `json:"seln_ernn_rate1,string"`
	SelnErnnRate2 float64         `json:"seln_ernn_rate2,string"`
	SelnErnnRate3 float64         `json:"seln_ernn_rate3,string"`
	SelnErnnRate4 float64         `json:"seln_ernn_rate4,string"`
	SelnErnnRate5 float64         `json:"seln_ernn_rate5,string"`
	ShnuErnnRate1 float64         `json:"shnu_ernn_rate1,string"`
	ShnuErnnRate2 float64         `json:"shnu_ernn_rate2,string"`
	ShnuErnnRate3 float64         `json:"shnu_ernn_rate3,string"`
	ShnuErnnRate4 float64         `json:"shnu_ernn_rate4,string"`
	ShnuErnnRate5 float64         `json:"shnu_ernn_rate5,string"`
}

type inquireAskingPriceResponse struct {
	RtCd   string          `json:"rt_cd"`
	MsgCd  string          `json:"msg_cd"`
	Msg1   string          `json:"msg1"`
	Output BondAskingPrice `json:"output"`
}

// InquireAskingPrice 는 채권 현재가 호가 5단계 (FHKBJ773401C0).
//
// KIS API: GET /uapi/domestic-bond/v1/quotations/inquire-asking-price
func (c *Client) InquireAskingPrice(ctx context.Context, params InquireAskingPriceParams) (*BondAskingPrice, error) {
	if params.MarketCode == "" {
		params.MarketCode = "B"
	}
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method:   http.MethodGet,
		Path:     "/uapi/domestic-bond/v1/quotations/inquire-asking-price",
		TrID:     "FHKBJ773401C0",
		CustType: "P",
		Query: map[string]string{
			"FID_COND_MRKT_DIV_CODE": params.MarketCode,
			"FID_INPUT_ISCD":         params.Symbol,
		},
	})
	if err != nil {
		return nil, err
	}
	var res inquireAskingPriceResponse
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse InquireAskingPrice: %w", err)
	}
	return &res.Output, nil
}
```

- [ ] Step 4: Verify PASS

```bash
cd /Users/user/src/workspace_moneyflow/korea-investment-stock && \
  go test ./bonds/... -run TestClient_InquireAskingPrice -v 2>&1
# Expected: PASS
```

- [ ] Step 5: gofmt/vet

```bash
cd /Users/user/src/workspace_moneyflow/korea-investment-stock && \
  gofmt -w bonds/quote.go bonds/quote_test.go && go vet ./bonds/... 2>&1
```

- [ ] Step 6: commit

```bash
git commit -m "$(cat <<'EOF'
[feat] bonds: InquireAskingPrice EP5 (FHKBJ773401C0, 34 fields, 5단계 호가)

채권 5단계 호가 — decimal askp/bidp, int64 잔량, float64 수익률 typed struct + method + test.

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 7: InquireDailyPrice (EP6) — APPEND

**Files:** APPEND to `bonds/quote.go` and `bonds/quote_test.go`

- [ ] Step 1: APPEND test to `bonds/quote_test.go`

```go
func TestClient_InquireDailyPrice(t *testing.T) {
	client, transport := newTestClient(t)
	transport.RegisterResponder(
		http.MethodGet,
		"=~/uapi/domestic-bond/v1/quotations/inquire-daily-price",
		httpmock.NewStringResponder(http.StatusOK, loadFixtureString(t, "bond_daily_price_success.json")),
	)

	got, err := client.InquireDailyPrice(context.Background(), InquireDailyPriceParams{
		MarketCode: "B",
		Symbol:     "KR103501GCC7",
	})
	require.NoError(t, err)
	require.NotNil(t, got)

	assert.Equal(t, "20260505", got.StckBsopDate)
	assert.Equal(t, "9823.50", got.BondPrpr.String())
	assert.Equal(t, int64(152000), got.AcmlVol)
	assert.InDelta(t, -0.13, got.PrdyCtrt, 0.001)
	assert.Equal(t, "9830.00", got.BondHgpr.String())
}
```

- [ ] Step 2: Verify FAIL

```bash
cd /Users/user/src/workspace_moneyflow/korea-investment-stock && \
  go test ./bonds/... -run TestClient_InquireDailyPrice -v 2>&1 | head -20
# Expected: compile error — InquireDailyPrice undefined
```

- [ ] Step 3: APPEND struct + Params + method to `bonds/quote.go`

```go
// ─── EP6: InquireDailyPrice ───────────────────────────────────────────────────

// InquireDailyPriceParams 는 채권 현재가 일별 요청 파라미터.
type InquireDailyPriceParams struct {
	MarketCode string // FID_COND_MRKT_DIV_CODE: 기본 "B"
	Symbol     string // FID_INPUT_ISCD: 채권 단축 종목코드 (필수)
}

// BondDailyPrice 는 채권 현재가 일별 시세. FHKBJ773404C0 — 9 fields typed.
//
// Note: KIS docs 는 output{} (object) 로 명시. 실제 API 가 array 를 반환하면 patch 에서 수정.
type BondDailyPrice struct {
	StckBsopDate string          `json:"stck_bsop_date"`
	BondPrpr     decimal.Decimal `json:"bond_prpr"`
	BondPrdyVrss decimal.Decimal `json:"bond_prdy_vrss"`
	PrdyVrssSign string          `json:"prdy_vrss_sign"`
	PrdyCtrt     float64         `json:"prdy_ctrt,string"`
	AcmlVol      int64           `json:"acml_vol,string"`
	BondOprc     decimal.Decimal `json:"bond_oprc"`
	BondHgpr     decimal.Decimal `json:"bond_hgpr"`
	BondLwpr     decimal.Decimal `json:"bond_lwpr"`
}

type inquireDailyPriceResponse struct {
	RtCd   string         `json:"rt_cd"`
	MsgCd  string         `json:"msg_cd"`
	Msg1   string         `json:"msg1"`
	Output BondDailyPrice `json:"output"`
}

// InquireDailyPrice 는 채권 현재가 일별 시세 (FHKBJ773404C0).
//
// KIS API: GET /uapi/domestic-bond/v1/quotations/inquire-daily-price
func (c *Client) InquireDailyPrice(ctx context.Context, params InquireDailyPriceParams) (*BondDailyPrice, error) {
	if params.MarketCode == "" {
		params.MarketCode = "B"
	}
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method:   http.MethodGet,
		Path:     "/uapi/domestic-bond/v1/quotations/inquire-daily-price",
		TrID:     "FHKBJ773404C0",
		CustType: "P",
		Query: map[string]string{
			"FID_COND_MRKT_DIV_CODE": params.MarketCode,
			"FID_INPUT_ISCD":         params.Symbol,
		},
	})
	if err != nil {
		return nil, err
	}
	var res inquireDailyPriceResponse
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse InquireDailyPrice: %w", err)
	}
	return &res.Output, nil
}
```

- [ ] Step 4: Verify PASS

```bash
cd /Users/user/src/workspace_moneyflow/korea-investment-stock && \
  go test ./bonds/... -run TestClient_InquireDailyPrice -v 2>&1
# Expected: PASS
```

- [ ] Step 5: gofmt/vet

```bash
cd /Users/user/src/workspace_moneyflow/korea-investment-stock && \
  gofmt -w bonds/quote.go bonds/quote_test.go && go vet ./bonds/... 2>&1
```

- [ ] Step 6: commit

```bash
git commit -m "$(cat <<'EOF'
[feat] bonds: InquireDailyPrice EP6 (FHKBJ773404C0, 9 fields typed)

채권 현재가 일별 시세 — typed struct + method + test.
KIS docs output{} 준수; 실제 API array 시 patch 예정.

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 8: InquireDailyItemchartprice (EP7) — output array (6 fields/item)

**File**: APPEND to `bonds/quote.go` and `bonds/quote_test.go`

| 항목 | 값 |
|---|---|
| Path | `/uapi/domestic-bond/v1/quotations/inquire-daily-itemchartprice` |
| TR_ID | `FHKBJ773701C0` |
| Output | `output []DailyItemchartpriceItem` (6 fields/item, max 30 records) |

### Params (2)

| Go 필드 | KIS 파라미터 | Required | 기본값 |
|---|---|---|---|
| MarketCode | `FID_COND_MRKT_DIV_CODE` | Y | `"B"` |
| Symbol | `FID_INPUT_ISCD` | Y | — |

### Output fields (6 per item)

| Go 필드 | KIS 필드 | 타입 |
|---|---|---|
| StckBsopDate | `stck_bsop_date` | string |
| BondOprc | `bond_oprc` | decimal |
| BondHgpr | `bond_hgpr` | decimal |
| BondLwpr | `bond_lwpr` | decimal |
| BondPrpr | `bond_prpr` | decimal |
| AcmlVol | `acml_vol` | int64 (string) |

### Step-by-step

- [ ] Step 1: test 작성 (`bonds/quote_test.go` APPEND)

```go
// TestInquireDailyItemchartprice — httpmock, testdata/bond_daily_itemchartprice_success.json
func TestInquireDailyItemchartprice(t *testing.T) {
    c, teardown := newMockBondsClient(t,
        "GET",
        "/uapi/domestic-bond/v1/quotations/inquire-daily-itemchartprice",
        "testdata/bond_daily_itemchartprice_success.json",
    )
    defer teardown()

    params := bonds.InquireDailyItemchartpriceParams{
        MarketCode: "B",
        Symbol:     "KR1010012345",
    }
    got, err := c.InquireDailyItemchartprice(context.Background(), params)
    require.NoError(t, err)

    // query param assertions
    lastReq := getLastRequest(t)
    assert.Equal(t, "B",            lastReq.URL.Query().Get("FID_COND_MRKT_DIV_CODE"))
    assert.Equal(t, "KR1010012345", lastReq.URL.Query().Get("FID_INPUT_ISCD"))

    // array shape
    require.GreaterOrEqual(t, len(got.Output), 2)

    // spot-check first item fields
    item := got.Output[0]
    assert.NotEmpty(t, item.StckBsopDate)
    assert.False(t, item.BondPrpr.IsZero())
    assert.GreaterOrEqual(t, item.AcmlVol, int64(0))
}
```

Testdata (`bonds/testdata/bond_daily_itemchartprice_success.json`):

```json
{
  "rt_cd": "0",
  "msg_cd": "MCA00000",
  "msg1": "정상처리 되었습니다.",
  "output": [
    {
      "stck_bsop_date": "20260505",
      "bond_oprc": "10000.00",
      "bond_hgpr": "10050.00",
      "bond_lwpr": "9980.00",
      "bond_prpr": "10020.00",
      "acml_vol": "3000"
    },
    {
      "stck_bsop_date": "20260504",
      "bond_oprc": "9990.00",
      "bond_hgpr": "10010.00",
      "bond_lwpr": "9960.00",
      "bond_prpr": "9995.00",
      "acml_vol": "2500"
    }
  ]
}
```

- [ ] Step 2: `go test ./bonds/... -run TestInquireDailyItemchartprice` — **FAIL** (미구현)

- [ ] Step 3: 구현 (`bonds/quote.go` APPEND)

```go
// DailyItemchartpriceItem — 장내채권 기간별 시세 항목 (EP7, FHKBJ773701C0)
type DailyItemchartpriceItem struct {
    StckBsopDate string          `json:"stck_bsop_date"` // 주식 영업 일자
    BondOprc     decimal.Decimal `json:"bond_oprc"`      // 채권 시가2
    BondHgpr     decimal.Decimal `json:"bond_hgpr"`      // 채권 고가
    BondLwpr     decimal.Decimal `json:"bond_lwpr"`      // 채권 저가
    BondPrpr     decimal.Decimal `json:"bond_prpr"`      // 채권 현재가
    AcmlVol      int64           `json:"acml_vol,string"` // 누적 거래량
}

// DailyItemchartprice — InquireDailyItemchartprice 응답
type DailyItemchartprice struct {
    Output []DailyItemchartpriceItem `json:"output"`
}

// InquireDailyItemchartpriceParams — EP7 요청 파라미터
type InquireDailyItemchartpriceParams struct {
    MarketCode string // FID_COND_MRKT_DIV_CODE (default "B")
    Symbol     string // FID_INPUT_ISCD — 채권 단축 종목코드
}

// InquireDailyItemchartprice — 장내채권 기간별 시세 조회 (EP7)
// TR_ID: FHKBJ773701C0
// Path:  /uapi/domestic-bond/v1/quotations/inquire-daily-itemchartprice
func (c *Client) InquireDailyItemchartprice(ctx context.Context, params InquireDailyItemchartpriceParams) (*DailyItemchartprice, error) {
    if params.MarketCode == "" {
        params.MarketCode = "B"
    }
    var result DailyItemchartprice
    resp, err := c.http.R().
        SetContext(ctx).
        SetHeader("tr_id", "FHKBJ773701C0").
        SetQueryParams(map[string]string{
            "FID_COND_MRKT_DIV_CODE": params.MarketCode,
            "FID_INPUT_ISCD":         params.Symbol,
        }).
        SetResult(&result).
        Get("/uapi/domestic-bond/v1/quotations/inquire-daily-itemchartprice")
    if err != nil {
        return nil, err
    }
    if err := checkResponse(resp); err != nil {
        return nil, err
    }
    return &result, nil
}
```

- [ ] Step 4: `go test ./bonds/... -run TestInquireDailyItemchartprice` — **PASS**

- [ ] Step 5: format + vet

```bash
cd /Users/user/src/workspace_moneyflow/korea-investment-stock && \
  gofmt -w bonds/quote.go bonds/quote_test.go && go vet ./bonds/... 2>&1
```

- [ ] Step 6: commit

```bash
git commit -m "$(cat <<'EOF'
[feat] bonds — InquireDailyItemchartprice (장내채권 기간별 시세, FHKBJ773701C0)

EP7: []DailyItemchartpriceItem array (6 fields/item, max 30 records).
Mirrors InquireDailyPrice pattern with array output shape.

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 9: InquireAvgUnit (EP8) — output1+output2+output3 (23+10+16 fields)

**File**: APPEND to `bonds/quote.go` and `bonds/quote_test.go`

| 항목 | 값 |
|---|---|
| Path | `/uapi/domestic-bond/v1/quotations/avg-unit` |
| TR_ID | `CTPF2005R` |
| Output | `output1 []AvgUnitEvalUnit` + `output2 []AvgUnitEvalAmt` + `output3 []AvgUnitPrice` |

### Params (7)

| Go 필드 | KIS 파라미터 | Required | 비고 |
|---|---|---|---|
| InqrStrtDt | `INQR_STRT_DT` | Y | 조회 시작일 (YYYYMMDD) |
| InqrEndDt | `INQR_END_DT` | Y | 조회 종료일 (YYYYMMDD) |
| Pdno | `PDNO` | Y | 채권 종목코드 |
| PrdtTypeCd | `PRDT_TYPE_CD` | Y | 상품유형코드 |
| VrfcKindCd | `VRFC_KIND_CD` | Y | 검증종류코드 |
| CtxAreaNk30 | `CTX_AREA_NK30` | Y | 연속조회검색조건 (blank 허용) |
| CtxAreaFk100 | `CTX_AREA_FK100` | Y | 연속조회키 (blank 허용) |

### output1 (23 fields) — `AvgUnitEvalUnit`

| Go 필드 | KIS 필드 | 타입 |
|---|---|---|
| EvluDt | `evlu_dt` | string |
| Pdno | `pdno` | string |
| PrdtTypeCd | `prdt_type_cd` | string |
| PrdtName | `prdt_name` | string |
| KisUnpr | `kis_unpr` | decimal |
| KbpUnpr | `kbp_unpr` | decimal |
| NiceEvluUnpr | `nice_evlu_unpr` | decimal |
| FnpUnpr | `fnp_unpr` | decimal |
| AvgEvluUnpr | `avg_evlu_unpr` | decimal |
| KisCrdtGradText | `kis_crdt_grad_text` | string |
| KbpCrdtGradText | `kbp_crdt_grad_text` | string |
| NiceCrdtGradText | `nice_crdt_grad_text` | string |
| FnpCrdtGradText | `fnp_crdt_grad_text` | string |
| ChngYn | `chng_yn` | string |
| KisErngRt | `kis_erng_rt` | float64 (string) |
| KbpErngRt | `kbp_erng_rt` | float64 (string) |
| NiceEvluErngRt | `nice_evlu_erng_rt` | float64 (string) |
| FnpErngRt | `fnp_erng_rt` | float64 (string) |
| AvgEvluErngRt | `avg_evlu_erng_rt` | float64 (string) |
| KisRfUnpr | `kis_rf_unpr` | decimal |
| KbpRfUnpr | `kbp_rf_unpr` | decimal |
| NiceEvluRfUnpr | `nice_evlu_rf_unpr` | decimal |
| AvgEvluRfUnpr | `avg_evlu_rf_unpr` | decimal |

### output2 (10 fields) — `AvgUnitEvalAmt`

| Go 필드 | KIS 필드 | 타입 |
|---|---|---|
| EvluDt | `evlu_dt` | string |
| Pdno | `pdno` | string |
| PrdtTypeCd | `prdt_type_cd` | string |
| PrdtName | `prdt_name` | string |
| KisEvluAmt | `kis_evlu_amt` | int64 (string) |
| KbpEvluAmt | `kbp_evlu_amt` | int64 (string) |
| NiceEvluAmt | `nice_evlu_amt` | int64 (string) |
| FnpEvluAmt | `fnp_evlu_amt` | int64 (string) |
| AvgEvluAmt | `avg_evlu_amt` | int64 (string) |
| ChngYn | `chng_yn` | string |

### output3 (16 fields) — `AvgUnitPrice`

| Go 필드 | KIS 필드 | 타입 |
|---|---|---|
| EvluDt | `evlu_dt` | string |
| Pdno | `pdno` | string |
| PrdtTypeCd | `prdt_type_cd` | string |
| PrdtName | `prdt_name` | string |
| KisCrcyCd | `kis_crcy_cd` | string |
| KisEvluUnitPric | `kis_evlu_unit_pric` | decimal |
| KisEvluPric | `kis_evlu_pric` | decimal |
| KbpCrcyCd | `kbp_crcy_cd` | string |
| KbpEvluUnitPric | `kbp_evlu_unit_pric` | decimal |
| KbpEvluPric | `kbp_evlu_pric` | decimal |
| NiceCrcyCd | `nice_crcy_cd` | string |
| NiceEvluUnitPric | `nice_evlu_unit_pric` | decimal |
| NiceEvluPric | `nice_evlu_pric` | decimal |
| AvgEvluUnitPric | `avg_evlu_unit_pric` | decimal |
| AvgEvluPric | `avg_evlu_pric` | decimal |
| ChngYn | `chng_yn` | string |

### Step-by-step

- [ ] Step 1: test 작성 (`bonds/quote_test.go` APPEND)

```go
// TestInquireAvgUnit — httpmock, testdata/avg_unit_success.json
func TestInquireAvgUnit(t *testing.T) {
    c, teardown := newMockBondsClient(t,
        "GET",
        "/uapi/domestic-bond/v1/quotations/avg-unit",
        "testdata/avg_unit_success.json",
    )
    defer teardown()

    params := bonds.InquireAvgUnitParams{
        InqrStrtDt:  "20260401",
        InqrEndDt:   "20260505",
        Pdno:        "KR1010012345",
        PrdtTypeCd:  "301",
        VrfcKindCd:  "01",
        CtxAreaNk30: "",
        CtxAreaFk100: "",
    }
    got, err := c.InquireAvgUnit(context.Background(), params)
    require.NoError(t, err)

    // query param assertions
    lastReq := getLastRequest(t)
    assert.Equal(t, "20260401",     lastReq.URL.Query().Get("INQR_STRT_DT"))
    assert.Equal(t, "20260505",     lastReq.URL.Query().Get("INQR_END_DT"))
    assert.Equal(t, "KR1010012345", lastReq.URL.Query().Get("PDNO"))
    assert.Equal(t, "301",          lastReq.URL.Query().Get("PRDT_TYPE_CD"))
    assert.Equal(t, "01",           lastReq.URL.Query().Get("VRFC_KIND_CD"))

    // output1 array
    require.GreaterOrEqual(t, len(got.Output1), 1)
    u := got.Output1[0]
    assert.NotEmpty(t, u.EvluDt)
    assert.NotEmpty(t, u.Pdno)
    assert.False(t, u.AvgEvluUnpr.IsZero())

    // output2 array
    require.GreaterOrEqual(t, len(got.Output2), 1)
    a := got.Output2[0]
    assert.NotEmpty(t, a.EvluDt)
    assert.GreaterOrEqual(t, a.AvgEvluAmt, int64(0))

    // output3 array
    require.GreaterOrEqual(t, len(got.Output3), 1)
    p := got.Output3[0]
    assert.NotEmpty(t, p.EvluDt)
    assert.False(t, p.AvgEvluUnitPric.IsZero())
}
```

Testdata (`bonds/testdata/avg_unit_success.json`):

```json
{
  "rt_cd": "0",
  "msg_cd": "MCA00000",
  "msg1": "정상처리 되었습니다.",
  "output1": [
    {
      "evlu_dt": "20260505",
      "pdno": "KR1010012345",
      "prdt_type_cd": "301",
      "prdt_name": "국고채권",
      "kis_unpr": "9850.00",
      "kbp_unpr": "9855.00",
      "nice_evlu_unpr": "9848.00",
      "fnp_unpr": "9852.00",
      "avg_evlu_unpr": "9851.25",
      "kis_crdt_grad_text": "AAA",
      "kbp_crdt_grad_text": "AAA",
      "nice_crdt_grad_text": "AAA",
      "fnp_crdt_grad_text": "AAA",
      "chng_yn": "N",
      "kis_erng_rt": "3.500",
      "kbp_erng_rt": "3.495",
      "nice_evlu_erng_rt": "3.505",
      "fnp_erng_rt": "3.498",
      "avg_evlu_erng_rt": "3.4995",
      "kis_rf_unpr": "9840.00",
      "kbp_rf_unpr": "9842.00",
      "nice_evlu_rf_unpr": "9838.00",
      "avg_evlu_rf_unpr": "9840.00"
    }
  ],
  "output2": [
    {
      "evlu_dt": "20260505",
      "pdno": "KR1010012345",
      "prdt_type_cd": "301",
      "prdt_name": "국고채권",
      "kis_evlu_amt": "9850000",
      "kbp_evlu_amt": "9855000",
      "nice_evlu_amt": "9848000",
      "fnp_evlu_amt": "9852000",
      "avg_evlu_amt": "9851250",
      "chng_yn": "N"
    }
  ],
  "output3": [
    {
      "evlu_dt": "20260505",
      "pdno": "KR1010012345",
      "prdt_type_cd": "301",
      "prdt_name": "국고채권",
      "kis_crcy_cd": "KRW",
      "kis_evlu_unit_pric": "9850.00",
      "kis_evlu_pric": "9850000.00",
      "kbp_crcy_cd": "KRW",
      "kbp_evlu_unit_pric": "9855.00",
      "kbp_evlu_pric": "9855000.00",
      "nice_crcy_cd": "KRW",
      "nice_evlu_unit_pric": "9848.00",
      "nice_evlu_pric": "9848000.00",
      "avg_evlu_unit_pric": "9851.25",
      "avg_evlu_pric": "9851250.00",
      "chng_yn": "N"
    }
  ]
}
```

- [ ] Step 2: `go test ./bonds/... -run TestInquireAvgUnit` — **FAIL** (미구현)

- [ ] Step 3: 구현 (`bonds/quote.go` APPEND)

```go
// AvgUnitEvalUnit — 평균단가조회 output1 항목 (EP8, CTPF2005R)
type AvgUnitEvalUnit struct {
    EvluDt           string          `json:"evlu_dt"`            // 평가 일자
    Pdno             string          `json:"pdno"`               // 상품번호
    PrdtTypeCd       string          `json:"prdt_type_cd"`       // 상품유형코드
    PrdtName         string          `json:"prdt_name"`          // 상품명
    KisUnpr          decimal.Decimal `json:"kis_unpr"`           // KIS 단가
    KbpUnpr          decimal.Decimal `json:"kbp_unpr"`           // KBP 단가
    NiceEvluUnpr     decimal.Decimal `json:"nice_evlu_unpr"`     // NICE 평가단가
    FnpUnpr          decimal.Decimal `json:"fnp_unpr"`           // FnP 단가
    AvgEvluUnpr      decimal.Decimal `json:"avg_evlu_unpr"`      // 평균 평가단가
    KisCrdtGradText  string          `json:"kis_crdt_grad_text"` // KIS 신용등급
    KbpCrdtGradText  string          `json:"kbp_crdt_grad_text"` // KBP 신용등급
    NiceCrdtGradText string          `json:"nice_crdt_grad_text"`// NICE 신용등급
    FnpCrdtGradText  string          `json:"fnp_crdt_grad_text"` // FnP 신용등급
    ChngYn           string          `json:"chng_yn"`            // 변경여부
    KisErngRt        float64         `json:"kis_erng_rt,string"` // KIS 수익률
    KbpErngRt        float64         `json:"kbp_erng_rt,string"` // KBP 수익률
    NiceEvluErngRt   float64         `json:"nice_evlu_erng_rt,string"` // NICE 평가수익률
    FnpErngRt        float64         `json:"fnp_erng_rt,string"` // FnP 수익률
    AvgEvluErngRt    float64         `json:"avg_evlu_erng_rt,string"`  // 평균 평가수익률
    KisRfUnpr        decimal.Decimal `json:"kis_rf_unpr"`        // KIS 기준단가
    KbpRfUnpr        decimal.Decimal `json:"kbp_rf_unpr"`        // KBP 기준단가
    NiceEvluRfUnpr   decimal.Decimal `json:"nice_evlu_rf_unpr"`  // NICE 평가기준단가
    AvgEvluRfUnpr    decimal.Decimal `json:"avg_evlu_rf_unpr"`   // 평균 평가기준단가
}

// AvgUnitEvalAmt — 평균단가조회 output2 항목 (EP8, CTPF2005R)
type AvgUnitEvalAmt struct {
    EvluDt      string `json:"evlu_dt"`           // 평가 일자
    Pdno        string `json:"pdno"`              // 상품번호
    PrdtTypeCd  string `json:"prdt_type_cd"`      // 상품유형코드
    PrdtName    string `json:"prdt_name"`         // 상품명
    KisEvluAmt  int64  `json:"kis_evlu_amt,string"`  // KIS 평가금액
    KbpEvluAmt  int64  `json:"kbp_evlu_amt,string"`  // KBP 평가금액
    NiceEvluAmt int64  `json:"nice_evlu_amt,string"` // NICE 평가금액
    FnpEvluAmt  int64  `json:"fnp_evlu_amt,string"`  // FnP 평가금액
    AvgEvluAmt  int64  `json:"avg_evlu_amt,string"`  // 평균 평가금액
    ChngYn      string `json:"chng_yn"`           // 변경여부
}

// AvgUnitPrice — 평균단가조회 output3 항목 (EP8, CTPF2005R)
type AvgUnitPrice struct {
    EvluDt          string          `json:"evlu_dt"`             // 평가 일자
    Pdno            string          `json:"pdno"`                // 상품번호
    PrdtTypeCd      string          `json:"prdt_type_cd"`        // 상품유형코드
    PrdtName        string          `json:"prdt_name"`           // 상품명
    KisCrcyCd       string          `json:"kis_crcy_cd"`         // KIS 통화코드
    KisEvluUnitPric decimal.Decimal `json:"kis_evlu_unit_pric"`  // KIS 평가단위가격
    KisEvluPric     decimal.Decimal `json:"kis_evlu_pric"`       // KIS 평가가격
    KbpCrcyCd       string          `json:"kbp_crcy_cd"`         // KBP 통화코드
    KbpEvluUnitPric decimal.Decimal `json:"kbp_evlu_unit_pric"`  // KBP 평가단위가격
    KbpEvluPric     decimal.Decimal `json:"kbp_evlu_pric"`       // KBP 평가가격
    NiceCrcyCd      string          `json:"nice_crcy_cd"`        // NICE 통화코드
    NiceEvluUnitPric decimal.Decimal `json:"nice_evlu_unit_pric"` // NICE 평가단위가격
    NiceEvluPric    decimal.Decimal `json:"nice_evlu_pric"`      // NICE 평가가격
    AvgEvluUnitPric decimal.Decimal `json:"avg_evlu_unit_pric"`  // 평균 평가단위가격
    AvgEvluPric     decimal.Decimal `json:"avg_evlu_pric"`       // 평균 평가가격
    ChngYn          string          `json:"chng_yn"`             // 변경여부
}

// AvgUnit — InquireAvgUnit 응답 (output1 + output2 + output3)
type AvgUnit struct {
    Output1 []AvgUnitEvalUnit `json:"output1"`
    Output2 []AvgUnitEvalAmt  `json:"output2"`
    Output3 []AvgUnitPrice    `json:"output3"`
}

// InquireAvgUnitParams — EP8 요청 파라미터
type InquireAvgUnitParams struct {
    InqrStrtDt   string // INQR_STRT_DT — 조회 시작일 (YYYYMMDD)
    InqrEndDt    string // INQR_END_DT — 조회 종료일 (YYYYMMDD)
    Pdno         string // PDNO — 채권 종목코드
    PrdtTypeCd   string // PRDT_TYPE_CD — 상품유형코드
    VrfcKindCd   string // VRFC_KIND_CD — 검증종류코드
    CtxAreaNk30  string // CTX_AREA_NK30 — 연속조회검색조건 (blank 허용)
    CtxAreaFk100 string // CTX_AREA_FK100 — 연속조회키 (blank 허용)
}

// InquireAvgUnit — 장내채권 평균단가 조회 (EP8)
// TR_ID: CTPF2005R
// Path:  /uapi/domestic-bond/v1/quotations/avg-unit
func (c *Client) InquireAvgUnit(ctx context.Context, params InquireAvgUnitParams) (*AvgUnit, error) {
    var result AvgUnit
    resp, err := c.http.R().
        SetContext(ctx).
        SetHeader("tr_id", "CTPF2005R").
        SetQueryParams(map[string]string{
            "INQR_STRT_DT":   params.InqrStrtDt,
            "INQR_END_DT":    params.InqrEndDt,
            "PDNO":           params.Pdno,
            "PRDT_TYPE_CD":   params.PrdtTypeCd,
            "VRFC_KIND_CD":   params.VrfcKindCd,
            "CTX_AREA_NK30":  params.CtxAreaNk30,
            "CTX_AREA_FK100": params.CtxAreaFk100,
        }).
        SetResult(&result).
        Get("/uapi/domestic-bond/v1/quotations/avg-unit")
    if err != nil {
        return nil, err
    }
    if err := checkResponse(resp); err != nil {
        return nil, err
    }
    return &result, nil
}
```

- [ ] Step 4: `go test ./bonds/... -run TestInquireAvgUnit` — **PASS**

- [ ] Step 5: format + vet

```bash
cd /Users/user/src/workspace_moneyflow/korea-investment-stock && \
  gofmt -w bonds/quote.go bonds/quote_test.go && go vet ./bonds/... 2>&1
```

- [ ] Step 6: commit

```bash
git commit -m "$(cat <<'EOF'
[feat] bonds — InquireAvgUnit (장내채권 평균단가조회, CTPF2005R)

EP8: multi-output (output1/output2/output3) 23+10+16 fields.
CTX_AREA_NK30/CTX_AREA_FK100 cursor pagination params 노출.

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 10: examples/bonds_quote/main.go

**Goal**: Phase 3.1 전체 8 메서드(EP1-EP8) 시연 예제 파일 생성.

**File**: `examples/bonds_quote/main.go` (신규 생성)

### Step-by-step

- [ ] Step 1: 예제 파일 작성

```go
package main

import (
    "context"
    "fmt"
    "log"

    kis "github.com/kenshin579/korea-investment-stock"
    "github.com/kenshin579/korea-investment-stock/bonds"
)

func main() {
    client, err := kis.NewClientFromEnv()
    if err != nil {
        log.Fatal(err)
    }
    ctx := context.Background()

    pdno := "KR1010012345" // example bond code (단축 종목코드)

    // 1. SearchBondInfo — 채권 기본정보 조회 (EP1, CTPF1114R)
    info, err := client.Bonds.SearchBondInfo(ctx, bonds.SearchBondInfoParams{
        Pdno: pdno, PrdtTypeCd: "301",
    })
    if err != nil {
        log.Printf("SearchBondInfo error: %v", err)
    } else {
        fmt.Printf("[EP1] SearchBondInfo: pdno=%s prdt_name=%s\n",
            info.Output.Pdno, info.Output.PrdtName)
    }

    // 2. InquireIssueInfo — 발행정보 조회 (EP2, CTPF1101R)
    issue, err := client.Bonds.InquireIssueInfo(ctx, bonds.InquireIssueInfoParams{
        Pdno: pdno, PrdtTypeCd: "301",
    })
    if err != nil {
        log.Printf("InquireIssueInfo error: %v", err)
    } else {
        fmt.Printf("[EP2] InquireIssueInfo: pdno=%s bond_issu_dt=%s\n",
            issue.Output.Pdno, issue.Output.BondIssuDt)
    }

    // 3. InquirePrice — 현재가 시세 (EP3, FHKBJ773400C0)
    price, err := client.Bonds.InquirePrice(ctx, bonds.InquirePriceParams{
        MarketCode: "B", Symbol: pdno,
    })
    if err != nil {
        log.Printf("InquirePrice error: %v", err)
    } else {
        fmt.Printf("[EP3] InquirePrice: bond_prpr=%s erng_rt=%f\n",
            price.Output.BondPrpr.String(), price.Output.ErngRt)
    }

    // 4. InquireCcnl — 현재가 체결 (EP4, FHKBJ773403C0)
    ccnl, err := client.Bonds.InquireCcnl(ctx, bonds.InquireCcnlParams{
        MarketCode: "B", Symbol: pdno,
    })
    if err != nil {
        log.Printf("InquireCcnl error: %v", err)
    } else {
        fmt.Printf("[EP4] InquireCcnl: stck_cntg_hour=%s bond_prpr=%s\n",
            ccnl.Output.StckCntgHour, ccnl.Output.BondPrpr.String())
    }

    // 5. InquireAskingPrice — 현재가 호가 (EP5, FHKBJ773401C0)
    ask, err := client.Bonds.InquireAskingPrice(ctx, bonds.InquireAskingPriceParams{
        MarketCode: "B", Symbol: pdno,
    })
    if err != nil {
        log.Printf("InquireAskingPrice error: %v", err)
    } else {
        fmt.Printf("[EP5] InquireAskingPrice: askp1=%s bidp1=%s\n",
            ask.Output.BondAskp1.String(), ask.Output.BondBidp1.String())
    }

    // 6. InquireDailyPrice — 현재가 일별 (EP6, FHKBJ773404C0)
    daily, err := client.Bonds.InquireDailyPrice(ctx, bonds.InquireDailyPriceParams{
        MarketCode: "B", Symbol: pdno,
    })
    if err != nil {
        log.Printf("InquireDailyPrice error: %v", err)
    } else {
        fmt.Printf("[EP6] InquireDailyPrice: stck_bsop_date=%s bond_prpr=%s\n",
            daily.Output.StckBsopDate, daily.Output.BondPrpr.String())
    }

    // 7. InquireDailyItemchartprice — 기간별 시세 (EP7, FHKBJ773701C0)
    chart, err := client.Bonds.InquireDailyItemchartprice(ctx, bonds.InquireDailyItemchartpriceParams{
        MarketCode: "B", Symbol: pdno,
    })
    if err != nil {
        log.Printf("InquireDailyItemchartprice error: %v", err)
    } else {
        fmt.Printf("[EP7] InquireDailyItemchartprice: %d records, first date=%s price=%s\n",
            len(chart.Output),
            func() string {
                if len(chart.Output) > 0 {
                    return chart.Output[0].StckBsopDate
                }
                return "N/A"
            }(),
            func() string {
                if len(chart.Output) > 0 {
                    return chart.Output[0].BondPrpr.String()
                }
                return "N/A"
            }(),
        )
    }

    // 8. InquireAvgUnit — 평균단가조회 (EP8, CTPF2005R)
    avg, err := client.Bonds.InquireAvgUnit(ctx, bonds.InquireAvgUnitParams{
        InqrStrtDt:   "20260401",
        InqrEndDt:    "20260505",
        Pdno:         pdno,
        PrdtTypeCd:   "301",
        VrfcKindCd:   "01",
        CtxAreaNk30:  "",
        CtxAreaFk100: "",
    })
    if err != nil {
        log.Printf("InquireAvgUnit error: %v", err)
    } else {
        fmt.Printf("[EP8] InquireAvgUnit: output1=%d output2=%d output3=%d records\n",
            len(avg.Output1), len(avg.Output2), len(avg.Output3))
        if len(avg.Output1) > 0 {
            fmt.Printf("       avg_evlu_unpr=%s avg_evlu_erng_rt=%f\n",
                avg.Output1[0].AvgEvluUnpr.String(), avg.Output1[0].AvgEvluErngRt)
        }
    }
}
```

- [ ] Step 2: 빌드 확인

```bash
cd /Users/user/src/workspace_moneyflow/korea-investment-stock && \
  go build ./examples/bonds_quote && echo OK
```

- [ ] Step 3: 빌드 아티팩트 삭제 (repo root 바이너리 방지)

```bash
rm -f /Users/user/src/workspace_moneyflow/korea-investment-stock/bonds_quote
```

- [ ] Step 4: format + vet

```bash
cd /Users/user/src/workspace_moneyflow/korea-investment-stock && \
  gofmt -w examples/bonds_quote/main.go && go vet ./examples/bonds_quote/... 2>&1
```

- [ ] Step 5: commit

```bash
git commit -m "$(cat <<'EOF'
[feat] examples — bonds_quote (Phase 3.1 8 메서드 시연)

examples/bonds_quote/main.go: EP1~EP8 전체 호출 시연.
KR1010012345 예시 채권코드 사용. build OK 확인.

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 11: 문서 갱신

**목표**: CLAUDE.md, README.md, CHANGELOG.md, bonds/doc.go 4개 파일 갱신.

### Step-by-step

- [ ] Step 1: 기존 파일 읽기 (Edit 전 필수)

```bash
# 각 파일 상단/하단 확인
head -20 /Users/user/src/workspace_moneyflow/korea-investment-stock/CLAUDE.md
head -5  /Users/user/src/workspace_moneyflow/korea-investment-stock/README.md
head -30 /Users/user/src/workspace_moneyflow/korea-investment-stock/CHANGELOG.md
```

- [ ] Step 2: **CLAUDE.md** 갱신

배너 교체:
```
> **Phase 2.7 — 업종/지수 7 메서드 (v1.10.0). Phase 2.5+ 완료.**
```
→
```
> **Phase 3.1 — 장내채권 시세 8 메서드 (v1.11.0). Phase 3 design spec 및 plan 참고.**
```

Phase 3 링크 bullets 추가 (Phase 2.7 plan 링크 아래):
```markdown
- Phase 3 design spec: [`docs/superpowers/specs/2026-05-05-phase3-bonds-design.md`](docs/superpowers/specs/2026-05-05-phase3-bonds-design.md)
- Phase 3.1 implementation plan: [`docs/superpowers/specs/2026-05-05-phase3-1-bonds-quote-implementation-plan.md`](docs/superpowers/specs/2026-05-05-phase3-1-bonds-quote-implementation-plan.md)
```

- [ ] Step 3: **README.md** 갱신

헤딩 교체:
```
Available Methods (Phase 1.2 ~ 2.7)
```
→
```
Available Methods (Phase 1.2 ~ 3.1)
```

메서드 수 교체: `71` → `79`

신규 "Bonds (장내채권) — Phase 3.1" 섹션 추가 (기존 방식 섹션 아래):

```markdown
### Bonds (장내채권) — Phase 3.1

| Go 메서드 | path | TR_ID |
|---|---|---|
| `Bonds.SearchBondInfo` | `quotations/search-bond-info` | CTPF1114R |
| `Bonds.InquireIssueInfo` | `quotations/issue-info` | CTPF1101R |
| `Bonds.InquirePrice` | `quotations/inquire-price` | FHKBJ773400C0 |
| `Bonds.InquireCcnl` | `quotations/inquire-ccnl` | FHKBJ773403C0 |
| `Bonds.InquireAskingPrice` | `quotations/inquire-asking-price` | FHKBJ773401C0 |
| `Bonds.InquireDailyPrice` | `quotations/inquire-daily-price` | FHKBJ773404C0 |
| `Bonds.InquireDailyItemchartprice` | `quotations/inquire-daily-itemchartprice` | FHKBJ773701C0 |
| `Bonds.InquireAvgUnit` | `quotations/avg-unit` | CTPF2005R |
```

- [ ] Step 4: **CHANGELOG.md** 갱신

`## [1.10.0]` 항목 바로 위에 삽입:

```markdown
## [1.11.0] - 2026-05-05

### Added — Phase 3.1 (장내채권 시세) — 신규 도메인

- 신규 sub-package `bonds/` 도입 (`client.Bonds.*`)
- `Bonds.SearchBondInfo` — 채권 기본조회 (CTPF1114R) — 70 fields all-string
- `Bonds.InquireIssueInfo` — 발행정보 (CTPF1101R) — 69 fields all-string
- `Bonds.InquirePrice` — 현재가 시세 (FHKBJ773400C0) — 17 fields typed
- `Bonds.InquireCcnl` — 현재가 체결 (FHKBJ773403C0) — 7 fields typed (single snapshot)
- `Bonds.InquireAskingPrice` — 현재가 호가 (FHKBJ773401C0) — 34 fields typed (5단계 호가)
- `Bonds.InquireDailyPrice` — 현재가 일별 (FHKBJ773404C0) — 9 fields typed
- `Bonds.InquireDailyItemchartprice` — 기간별 시세 (FHKBJ773701C0) — 6 fields/item array
- `Bonds.InquireAvgUnit` — 평균단가조회 (CTPF2005R) — output1/output2/output3 (23+10+16 fields)
- examples: `bonds_quote`

### Notes

- Phase 3 신규 도메인 시작 — 장내채권 (Korean bond) sub-package.
- EP1 (`SearchBondInfo`) 는 path "search-bond-info" 에 동사 포함 — `Inquire` prefix 강제 안 함 (Style A 변형).
- EP1+EP2 의 70/69 fields 는 KSD-style all-string mapping (KIS docs 가 모두 String 타입 명시).
- EP3-EP8 은 Phase 2 standard typed mapping (decimal/int64/float64/string).
- 채권 현재가/호가에서 `bond_oprc` (시가2 — KIS naming artifact), `stck_cntg_hour`/`stck_bsop_date` (cross-domain stock prefix) 등 KIS 명명 그대로 보존.
- EP8 `CTX_AREA_NK30`/`CTX_AREA_FK100` cursor pagination params 노출.
- Phase 3.2 (잔고/주문조회 4 메서드, 계좌 인증 필요) 는 Trading 도메인과 함께 추후 결정.
- 누적 71 → 79 메서드.

```

- [ ] Step 5: **bonds/doc.go** — Task 2 에서 이미 작성됨. 변경 불필요.

- [ ] Step 6: 빌드/린트 검증

```bash
cd /Users/user/src/workspace_moneyflow/korea-investment-stock && \
  go build ./... && go vet ./... && gofmt -l . 2>&1
```

> **CRITICAL**: `gofmt -l .` 출력이 없어야 함 (빈 출력 = OK). `go build` 가 repo root 에 바이너리를 남겼다면 삭제 후 커밋.

- [ ] Step 7: commit

```bash
git commit -m "$(cat <<'EOF'
[docs] Phase 3.1 문서 갱신 (v1.11.0, 신규 bonds 도메인)

CLAUDE.md 배너 갱신 + Phase 3 링크 추가.
README.md 헤딩/카운트 갱신 + Bonds 섹션 추가.
CHANGELOG.md v1.11.0 항목 추가.

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 12: 최종 점검

**목표**: Phase 3.1 전체 품질 최종 확인.

### Step-by-step

- [ ] Step 1: gofmt 검사

```bash
cd /Users/user/src/workspace_moneyflow/korea-investment-stock && \
  gofmt -l . | head
# 출력 없음이어야 함 (silent = OK)
```

- [ ] Step 2: 빌드 + vet

```bash
cd /Users/user/src/workspace_moneyflow/korea-investment-stock && \
  go build ./... && go vet ./...
# 출력 없음이어야 함 (silent = OK)
```

- [ ] Step 3: 전체 테스트 (race detector)

```bash
cd /Users/user/src/workspace_moneyflow/korea-investment-stock && \
  go test ./... -race -count=1
# 모든 패키지 PASS
```

- [ ] Step 4: 커버리지 측정

```bash
# bonds 패키지
cd /Users/user/src/workspace_moneyflow/korea-investment-stock && \
  go test ./bonds/... -coverprofile=/tmp/cov_b.out -covermode=atomic && \
  go tool cover -func=/tmp/cov_b.out | tail -2
# 기대: total ≥ 80%

# root kis 패키지
cd /Users/user/src/workspace_moneyflow/korea-investment-stock && \
  go test . -coverprofile=/tmp/cov_r.out -covermode=atomic && \
  go tool cover -func=/tmp/cov_r.out | tail -2
# 기대: total ≥ 80%
```

> 커버리지 < 80% 시: Phase 2.6 교훈 적용 — 각 메서드에 InvalidJSON error path 테스트 추가 후 재측정.

- [ ] Step 5: 파일 존재 확인 (13개 기대)

```bash
ls \
  bonds/client.go \
  bonds/quote.go \
  bonds/quote_test.go \
  bonds/doc.go \
  bonds/testhelper_test.go \
  bonds/testdata/search_bond_info_success.json \
  bonds/testdata/issue_info_success.json \
  bonds/testdata/bond_price_success.json \
  bonds/testdata/bond_ccnl_success.json \
  bonds/testdata/bond_asking_price_success.json \
  bonds/testdata/bond_daily_price_success.json \
  bonds/testdata/bond_daily_itemchartprice_success.json \
  bonds/testdata/avg_unit_success.json \
  examples/bonds_quote/main.go \
  2>&1 | wc -l
# 기대: 14 (14 lines = 14 files found)
```

- [ ] Step 6: 브랜치 커밋 수 확인

```bash
git log main..HEAD --oneline | wc -l
# Phase 3.1 커밋 수 확인 (Task 1~11 커밋 포함)
```

---

## Task 13: PR 생성 (사용자 승인 후)

> **Claude 는 push / PR 생성을 사용자 명시적 승인 후에만 실행 (글로벌 정책).**

### Step-by-step

- [ ] Step 1: **사용자 승인 요청** — 작업 진행 보고 후 PR 생성 가능 여부 confirm.

  보고 내용:
  - 완료된 Tasks (1-12) 요약
  - 테스트 결과 (PASS/커버리지)
  - 커밋 수
  - PR 생성 진행 여부 질문

- [ ] Step 2: **Push** (승인 후)

```bash
git push -u origin feat/phase3-1-bonds-quote
```

- [ ] Step 3: **PR 생성** (승인 후)

```bash
gh pr create \
  --title "Phase 3.1 — 장내채권 시세 (v1.11.0) [신규 bonds 도메인]" \
  --reviewer kenshin579 \
  --base main \
  --head feat/phase3-1-bonds-quote \
  --body "$(cat <<'EOF'
## Summary

- **Phase 3.1 implementation (8 NEW methods)** — 신규 `bonds/` sub-package 도입
- Phase 2 패턴 그대로 (Style A, Params struct, KIS docs 1:1)
- v1.11.0 release (누적 71 → 79)
- **신규 도메인** — 장내채권 (Korean bond) 시세 조회

## 메서드 → 한투 API 매핑 (8 NEW)

| Go 메서드 | path | TR_ID |
|---|---|---|
| Bonds.SearchBondInfo | quotations/search-bond-info | CTPF1114R |
| Bonds.InquireIssueInfo | quotations/issue-info | CTPF1101R |
| Bonds.InquirePrice | quotations/inquire-price | FHKBJ773400C0 |
| Bonds.InquireCcnl | quotations/inquire-ccnl | FHKBJ773403C0 |
| Bonds.InquireAskingPrice | quotations/inquire-asking-price | FHKBJ773401C0 |
| Bonds.InquireDailyPrice | quotations/inquire-daily-price | FHKBJ773404C0 |
| Bonds.InquireDailyItemchartprice | quotations/inquire-daily-itemchartprice | FHKBJ773701C0 |
| Bonds.InquireAvgUnit | quotations/avg-unit | CTPF2005R |

## 구현 이슈 및 대응

- 신규 sub-package bonds/ 도입 + root client 에 Bonds 필드 통합
- EP1 path "search-bond-info" 에 동사 포함 → SearchBondInfo (Inquire prefix 강제 안 함)
- EP1+EP2 (CTPF11* TR_IDs) 70/69 fields 모두 plain string (KSD-like)
- EP3-EP8 typed mapping (decimal/int64/float64/string)
- EP4 (InquireCcnl) 단일 snapshot object (체결 list 아님)
- EP6 (InquireDailyPrice) doc 의 Object shape 모순 가능성 — doc verbatim 으로 single object 구현
- EP8 multi-output (output1/output2/output3) + CTX_AREA cursor pagination
- KIS cross-domain 명명 보존: bond_oprc (시가2), stck_cntg_hour/stck_bsop_date (stock prefix)

## Test Plan

- [x] go build/vet/fmt clean
- [x] go test ./... -race -count=1 모든 패키지 PASS
- [x] Coverage bonds ≥80%, root kis ≥80%
- [x] httpmock 단위 테스트 (8 methods)
- [x] examples/bonds_quote build OK

## Breaking Changes

없음 — 신규 도메인/메서드 추가만.

## Notes

- Phase 3.2-3.4 (잔고/주문조회/주문/실시간) 는 미정. 사용자 결정 후 진행.
- 누적 79 메서드 (Phase 1: 28 + Phase 2: 25 + Phase 2.5+: 18 + Phase 3.1: 8).

## 참고 문서

- Phase 3 design spec: docs/superpowers/specs/2026-05-05-phase3-bonds-design.md
- Phase 3.1 plan: docs/superpowers/specs/2026-05-05-phase3-1-bonds-quote-implementation-plan.md

🤖 Generated with [Claude Code](https://claude.com/claude-code)
EOF
)"
```

- [ ] Step 4: **Merge** (사용자 최종 승인 후)

```bash
gh pr merge <PR#> --merge
```

- [ ] Step 5: **릴리즈 태깅** (merge 후)

```bash
git tag -a v1.11.0 -m "Phase 3.1 — 장내채권 시세 8 메서드 (신규 bonds 도메인)"
git push origin v1.11.0
gh release create v1.11.0 \
  --title "v1.11.0 — Phase 3.1 장내채권 시세 (bonds 신규 도메인)" \
  --notes-file <(awk '/^## \[1.11.0\]/{f=1;next} /^## \[/&&f{exit} f' CHANGELOG.md)
```
