# Phase 1.5 — 해외주식 (Python parity 완성) Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** `korea-investment-stock` Go 라이브러리에 해외주식 5 메서드 + `FetchOverseasSymbols(market)` 1 통합 메서드 = 6 메서드 추가, Python parity 완성 (`v1.3.0` release).

**Architecture:** Phase 1.2 의 인프라 + 패턴 재사용. `overseas/{price,search,chart,ranking,symbols}.go` 추가, `internal/overseasmaster/` 신규 (KRX 와 분리). 한투 API path 1:1 매핑 (Style A). `overseas.New(http, master)` 시그니처 확장 (Phase 1.2 의 domestic 패턴). TDD: testdata fixture (한투 docs 응답 필드 정의 → 합성 JSON, 해외 마스터는 sample 첫 3행 가능 시) → 실패 테스트 → struct + 메서드 구현 → 통과 → commit.

**Tech Stack:** Go 1.25+, `github.com/jarcoal/httpmock`, `github.com/shopspring/decimal`, `github.com/stretchr/testify`. 새 dependency 없음 (cp949 디코딩 의 `golang.org/x/text` 는 Phase 1.2 부터 추가됨).

**참고 spec:**
- Phase 1 design spec (Phase 1.5 amendment): `docs/superpowers/specs/2026-05-03-phase1-api-coverage-design.md` (commit `b824d64`)
- Phase 1.2 plan (krxmaster 패턴 reference): `docs/superpowers/specs/2026-05-03-phase1-2-domestic-quotes-implementation-plan.md`
- Phase 1.4 plan (가장 최근): `docs/superpowers/specs/2026-05-05-phase1-4-investor-industry-ipo-implementation-plan.md`
- 한투 API docs: `docs/api/해외주식/{해외주식_현재가상세, 해외주식_상품기본정보, 해외주식_기간별시세, 해외주식_종목_지수_환율기간별시세(일_주_월_년), 해외주식_상승율_하락율}.md`

---

## 사전 정보

| 항목 | 값 |
|------|---|
| 작업 브랜치 | `docs/phase1-5-spec` (이미 생성됨) |
| 시작 HEAD | `b824d64` (Phase 1 design spec amendment commit) |
| Release 목표 | `v1.3.0` (PR merge 후 태그) |
| PR 베이스 | `main` |
| 현재 main 상태 | v1.2.0 publish 완료 (Phase 1.2+1.3+1.4 통합, 22 메서드) |

---

## 메서드 → 한투 API 매핑

| Go 메서드 | 한투 path | TR_ID | docs |
|-----------|----------|-------|------|
| `Overseas.InquirePriceDetail(ctx, params)` | `/uapi/overseas-price/v1/quotations/price-detail` | HHDFS76200200 | 해외주식_현재가상세.md |
| `Overseas.SearchInfo(ctx, params)` | `/uapi/overseas-price/v1/quotations/search-info` | CTPF1702R | 해외주식_상품기본정보.md |
| `Overseas.InquireDailyPrice(ctx, params)` | `/uapi/overseas-price/v1/quotations/dailyprice` | HHDFS76240000 | 해외주식_기간별시세.md |
| `Overseas.InquireDailyChartPrice(ctx, params)` | `/uapi/overseas-price/v1/quotations/inquire-daily-chartprice` | FHKST03030100 | 해외주식_종목_지수_환율기간별시세(일_주_월_년).md |
| `Overseas.InquireUpdownRate(ctx, params)` | `/uapi/overseas-stock/v1/ranking/updown-rate` | HHDFS76290000 | 해외주식_상승율_하락율.md |
| `Overseas.FetchOverseasSymbols(ctx, market)` | `https://new.real.download.dws.co.kr/common/master/<market>mst.cod.zip` (KRX 공개) | — | (KIS API 아님) |

---

## 파일 구조

### 신규 (internal)
- `internal/overseasmaster/doc.go`
- `internal/overseasmaster/overseasmaster.go` — `Symbol` struct + `Parse(market, zipBytes)` + `MarketURL(market)` 상수
- `internal/overseasmaster/overseasmaster_test.go`
- `internal/overseasmaster/testdata/<market>_code_sample.cod.zip` (실제 KIS 다운로드 첫 3행)
- `internal/overseasmaster/testdata/README.md`

### 신규 (overseas)
- `overseas/price.go` — `InquirePriceDetail` + `PriceDetail` struct
- `overseas/price_test.go`
- `overseas/search.go` — `SearchInfo` + `OverseasProductInfo` struct (이름 충돌 회피)
- `overseas/search_test.go`
- `overseas/chart.go` — `InquireDailyPrice` + `InquireDailyChartPrice` + structs
- `overseas/chart_test.go`
- `overseas/ranking.go` — `InquireUpdownRate` + `UpdownRate` struct
- `overseas/ranking_test.go`
- `overseas/symbols.go` — `FetchOverseasSymbols` + type alias re-export + `downloadURL` helper
- `overseas/symbols_test.go`
- `overseas/testhelper_test.go` — `newTestClient`, `loadFixture`, `stubTokenManager` (overseas 패키지 전용)
- `overseas/testdata/price_detail_success.json`
- `overseas/testdata/search_info_success.json`
- `overseas/testdata/daily_price_success.json`
- `overseas/testdata/daily_chart_price_success.json`
- `overseas/testdata/updown_rate_success.json`
- `overseas/testdata/<market>_code_sample.cod.zip` (overseasmaster 와 동일 sample 복사)
- `overseas/testdata/README.md`

### 신규 (examples)
- `examples/overseas_price/main.go` — InquirePriceDetail + SearchInfo
- `examples/overseas_chart/main.go` — InquireDailyPrice + InquireDailyChartPrice
- `examples/overseas_symbols/main.go` — FetchOverseasSymbols("nas")

### 수정 (root)
- `client.go` — `wireInfra` 의 `c.Overseas = overseas.New(c.httpClient)` 를 `overseas.New(c.httpClient, c.masterC)` 로 변경
- `CLAUDE.md` — Phase 1.5 안내, "Out of Scope" section 업데이트
- `README.md` — Available Methods 표 갱신 (22 → 28 메서드)
- `CHANGELOG.md` — `[1.3.0]` entry

### 수정 (sub-packages)
- `overseas/client.go` — `Client` struct 에 `master *mastercache.Cache` 추가, `New(http, master)` 시그니처
- `overseas/doc.go` — Phase 1.5 메서드 안내

---

## 타입 매핑 규칙

- **가격/지수 → `decimal.Decimal` (bare tag)**: `last`, `base`, `open`, `high`, `low`, `clos`, `diff`, `h52p`, `l52p`, `uplp`, `dnlp`, `pbid`, `pask`, `n_base`, `n_diff`, `t_xprc`, `t_xdif`, `p_xprc`, `p_xdif`, `e_parp`, `ovrs_papr`, `ovrs_nmix_prpr`, `ovrs_nmix_prdy_vrss`, `ovrs_nmix_prdy_clpr`, `ovrs_prod_oprc`, `ovrs_prod_hgpr`, `ovrs_prod_lwpr`, `ovrs_nmix_oprc`, `ovrs_nmix_hgpr`, `ovrs_nmix_lwpr`, `ovrs_now_pric1`
- **수량/금액 → `int64,string`**: `pvol`, `tvol`, `tamt`, `pamt`, `tomv`, `mcap`, `shar`, `vbid`, `vask`, `acml_vol`, `prdy_vol`, `lstg_stck_num`, `crec`, `trec`, `nrec`, `rank` (KIS 가 string 이지만 정수)
- **비율/percentage → `float64,string`**: `rate`, `prdy_ctrt`, `n_rate`, `t_xrat`, `p_xrat`, `t_rate`, `p_rate`, `perx`, `pbrx`
- **주당 가치 → `decimal.Decimal`**: `epsx`, `bpsx`
- **코드/이름/날짜/Y-N/부호 → 평문 `string`**: `rsym`, `sign`, `curr`, `zdiv`, `vnit`, `t_xsgn`, `p_xsng`, `e_ordyn`, `e_hogau`, `e_icod`, `etyp_nm`, `h52d`, `l52d`, `xymd`, `mod_yn`, `gubn`, `excd`, `symb`, `name`, `ename`, `stat`, `stck_bsop_date`, `stck_shrn_iscd`, `hts_kor_isnm`, `prdy_vrss_sign`, 그리고 SearchInfo 의 거의 모든 메타 필드 (std_pdno, prdt_eng_name, natn_cd, natn_name, tr_mket_cd, tr_mket_name, ovrs_excg_cd, ovrs_excg_name, tr_crcy_cd, crcy_name, ovrs_stck_dvsn_cd, prdt_clsf_cd, prdt_clsf_name, lstg_dt, lstg_abol_item_yn, lstg_yn, tax_levy_yn, ovrs_stck_erlm_rosn_cd, ovrs_stck_hist_rght_dvsn_cd, chng_bf_pdno, prdt_type_cd_2, ovrs_item_name, sedol_no, blbg_tckr_text, ovrs_stck_etf_risk_drtp_cd, etp_chas_erng_rt_dbnb, istt_usge_isin_cd, mint_svc_yn, mint_svc_yn_chng_dt, prdt_name, lei_cd, ovrs_stck_stop_rson_cd, lstg_abol_dt, mini_stk_tr_stat_dvsn_cd, mint_frst_svc_erlm_dt, mint_dcpt_trad_psbl_yn, mint_fnum_trad_psbl_yn, mint_cblc_cvsn_ipsb_yn, ptp_item_yn, ptp_item_trfx_exmt_yn, ptp_item_trfx_exmt_strt_dt, ptp_item_trfx_exmt_end_dt, dtm_tr_psbl_yn, sdrf_stop_ecls_yn, sdrf_stop_ecls_erlm_dt, memo_text1, last_rcvg_dtime), 그리고 SearchInfo 의 `sll_unit_qty`, `buy_unit_qty`, `tr_unit_amt`, `ovrs_stck_tr_stop_dvsn_cd`, `ovrs_stck_prdt_grp_no` (KIS 가 string 으로 줌)

---

## Task 1: testdata fixtures (5 KIS API JSON + 1 해외 마스터 sample)

**Files:**
- Create: `overseas/testdata/{price_detail,search_info,daily_price,daily_chart_price,updown_rate}_success.json` (5)
- Create: `overseas/testdata/README.md`
- Create (Step 6+): `internal/overseasmaster/testdata/<market>_code_sample.cod.zip` + README + `overseas/testdata/<market>_code_sample.cod.zip` (sample 복사)

> 5 KIS API testdata 는 docs 응답 필드 정의 기반 합성. 해외 마스터 sample 은 KIS 공개 다운로드 (`https://new.real.download.dws.co.kr/common/master/nasmst.cod.zip` 등) 의 첫 3 행. KRX 패턴과 동일 — 만약 해당 URL 이 동작하지 않으면 STOP 하고 보고 (Phase 1.5 의 핵심 가정 검증).

- [ ] **Step 1: price_detail_success.json**

```json
{
  "rt_cd": "0",
  "msg_cd": "MCA00000",
  "msg1": "정상처리 되었습니다.",
  "output": {
    "rsym": "DNASAAPL",
    "pvol": "120000000",
    "open": "180.50",
    "high": "182.30",
    "low": "179.80",
    "last": "181.45",
    "base": "180.20",
    "tomv": "2800000000000",
    "pamt": "21800000000",
    "uplp": "0",
    "dnlp": "0",
    "h52p": "199.62",
    "h52d": "20251220",
    "l52p": "164.08",
    "l52d": "20260105",
    "perx": "29.45",
    "pbrx": "45.20",
    "epsx": "6.15",
    "bpsx": "4.01",
    "shar": "15400000000",
    "mcap": "73450000000",
    "curr": "USD",
    "zdiv": "2",
    "vnit": "1",
    "t_xprc": "246123",
    "t_xdif": "1696",
    "t_xrat": "0.69",
    "p_xprc": "244427",
    "p_xdif": "0",
    "p_xrat": "0",
    "t_rate": "1356.50",
    "p_rate": "1356.20",
    "t_xsgn": "2",
    "p_xsng": "3",
    "e_ordyn": "매매가능",
    "e_hogau": "0.01",
    "e_icod": "Technology",
    "e_parp": "0.0001",
    "tvol": "85000000",
    "tamt": "15400000000",
    "etyp_nm": ""
  }
}
```

- [ ] **Step 2: search_info_success.json**

```json
{
  "rt_cd": "0",
  "msg_cd": "MCA00000",
  "msg1": "정상처리 되었습니다.",
  "output": {
    "std_pdno": "US0378331005",
    "prdt_eng_name": "APPLE INC",
    "natn_cd": "840",
    "natn_name": "미국",
    "tr_mket_cd": "01",
    "tr_mket_name": "나스닥",
    "ovrs_excg_cd": "NASD",
    "ovrs_excg_name": "NASDAQ",
    "tr_crcy_cd": "USD",
    "ovrs_papr": "0.00001",
    "crcy_name": "미국 달러",
    "ovrs_stck_dvsn_cd": "01",
    "prdt_clsf_cd": "STK",
    "prdt_clsf_name": "주권",
    "sll_unit_qty": "1",
    "buy_unit_qty": "1",
    "tr_unit_amt": "0.01",
    "lstg_stck_num": "15400000000",
    "lstg_dt": "19801212",
    "ovrs_stck_tr_stop_dvsn_cd": "01",
    "lstg_abol_item_yn": "N",
    "ovrs_stck_prdt_grp_no": "",
    "lstg_yn": "Y",
    "tax_levy_yn": "Y",
    "ovrs_stck_erlm_rosn_cd": "01",
    "ovrs_stck_hist_rght_dvsn_cd": "",
    "chng_bf_pdno": "",
    "prdt_type_cd_2": "512",
    "ovrs_item_name": "APPLE INC",
    "sedol_no": "2046251",
    "blbg_tckr_text": "AAPL US Equity",
    "ovrs_stck_etf_risk_drtp_cd": "",
    "etp_chas_erng_rt_dbnb": "",
    "istt_usge_isin_cd": "",
    "mint_svc_yn": "Y",
    "mint_svc_yn_chng_dt": "20210301",
    "prdt_name": "애플",
    "lei_cd": "HWUPKR0MPOU8FGXBT394",
    "ovrs_stck_stop_rson_cd": "",
    "lstg_abol_dt": "",
    "mini_stk_tr_stat_dvsn_cd": "01",
    "mint_frst_svc_erlm_dt": "20210301",
    "mint_dcpt_trad_psbl_yn": "Y",
    "mint_fnum_trad_psbl_yn": "Y",
    "mint_cblc_cvsn_ipsb_yn": "N",
    "ptp_item_yn": "N",
    "ptp_item_trfx_exmt_yn": "N",
    "ptp_item_trfx_exmt_strt_dt": "",
    "ptp_item_trfx_exmt_end_dt": "",
    "dtm_tr_psbl_yn": "Y",
    "sdrf_stop_ecls_yn": "N",
    "sdrf_stop_ecls_erlm_dt": "",
    "memo_text1": "",
    "ovrs_now_pric1": "181.45",
    "last_rcvg_dtime": "20260505131500"
  }
}
```

- [ ] **Step 3: daily_price_success.json**

```json
{
  "rt_cd": "0",
  "msg_cd": "MCA00000",
  "msg1": "정상처리 되었습니다.",
  "output1": {
    "rsym": "DNASAAPL",
    "zdiv": "2",
    "nrec": "180.20"
  },
  "output2": [
    {
      "xymd": "20260505",
      "clos": "181.45",
      "sign": "2",
      "diff": "1.25",
      "rate": "0.69",
      "open": "180.50",
      "high": "182.30",
      "low": "179.80",
      "tvol": "85000000",
      "tamt": "15400000000",
      "pbid": "181.40",
      "vbid": "1500",
      "pask": "181.50",
      "vask": "2300"
    },
    {
      "xymd": "20260502",
      "clos": "180.20",
      "sign": "5",
      "diff": "-0.85",
      "rate": "-0.47",
      "open": "181.00",
      "high": "181.80",
      "low": "179.50",
      "tvol": "75000000",
      "tamt": "13500000000",
      "pbid": "180.15",
      "vbid": "1200",
      "pask": "180.25",
      "vask": "1800"
    }
  ]
}
```

- [ ] **Step 4: daily_chart_price_success.json**

```json
{
  "rt_cd": "0",
  "msg_cd": "MCA00000",
  "msg1": "정상처리 되었습니다.",
  "output1": {
    "ovrs_nmix_prdy_vrss": "12.4500",
    "prdy_vrss_sign": "2",
    "prdy_ctrt": "0.25",
    "ovrs_nmix_prdy_clpr": "5012.3400",
    "acml_vol": "1234567",
    "hts_kor_isnm": "S&P 500",
    "ovrs_nmix_prpr": "5024.7900",
    "stck_shrn_iscd": "SPX",
    "prdy_vol": "1100000",
    "ovrs_prod_oprc": "5015.0000",
    "ovrs_prod_hgpr": "5028.5000",
    "ovrs_prod_lwpr": "5010.2000"
  },
  "output2": [
    {
      "stck_bsop_date": "20260505",
      "ovrs_nmix_prpr": "5024.7900",
      "ovrs_nmix_oprc": "5015.0000",
      "ovrs_nmix_hgpr": "5028.5000",
      "ovrs_nmix_lwpr": "5010.2000",
      "acml_vol": "1234567",
      "mod_yn": "N"
    },
    {
      "stck_bsop_date": "20260502",
      "ovrs_nmix_prpr": "5012.3400",
      "ovrs_nmix_oprc": "5008.0000",
      "ovrs_nmix_hgpr": "5018.7000",
      "ovrs_nmix_lwpr": "5005.1000",
      "acml_vol": "1100000",
      "mod_yn": "N"
    }
  ]
}
```

- [ ] **Step 5: updown_rate_success.json**

```json
{
  "rt_cd": "0",
  "msg_cd": "MCA00000",
  "msg1": "정상처리 되었습니다.",
  "output1": {
    "zdiv": "2",
    "stat": "정상",
    "crec": "2",
    "trec": "100",
    "nrec": "30"
  },
  "output2": [
    {
      "rsym": "DNASNVDA",
      "excd": "NAS",
      "symb": "NVDA",
      "name": "엔비디아",
      "last": "920.45",
      "sign": "2",
      "diff": "45.20",
      "rate": "5.16",
      "tvol": "120000000",
      "pask": "920.50",
      "pbid": "920.40",
      "n_base": "875.25",
      "n_diff": "45.20",
      "n_rate": "5.16",
      "rank": "1",
      "ename": "NVIDIA CORP",
      "e_ordyn": "Y"
    },
    {
      "rsym": "DNASTSLA",
      "excd": "NAS",
      "symb": "TSLA",
      "name": "테슬라",
      "last": "245.80",
      "sign": "2",
      "diff": "10.50",
      "rate": "4.46",
      "tvol": "85000000",
      "pask": "245.85",
      "pbid": "245.75",
      "n_base": "235.30",
      "n_diff": "10.50",
      "n_rate": "4.46",
      "rank": "2",
      "ename": "TESLA INC",
      "e_ordyn": "Y"
    }
  ]
}
```

- [ ] **Step 6: overseas/testdata/README.md**

```markdown
# overseas testdata

각 한투 API 메서드의 단위 테스트 fixture.

## REST API 응답 (합성 JSON)

- `price_detail_success.json` — 해외주식_현재가상세 (HHDFS76200200) 정상 응답
- `search_info_success.json` — 해외주식_상품기본정보 (CTPF1702R) 정상 응답
- `daily_price_success.json` — 해외주식_기간별시세 (HHDFS76240000) 정상 응답
- `daily_chart_price_success.json` — 해외주식 종목/지수/환율 기간별시세 (FHKST03030100) 정상 응답
- `updown_rate_success.json` — 해외주식_상승율_하락율 (HHDFS76290000) 정상 응답

각 JSON 의 필드는 `docs/api/해외주식/<API>.md` 의 응답 필드 정의에 1:1 매핑. 값은 합성 (실제 시세 아님).

## 해외 마스터 sample

- `<market>_code_sample.cod.zip` — 출처 + 재생성 방법은 `internal/overseasmaster/testdata/README.md` 참조
```

- [ ] **Step 7: 검증**

```bash
for f in overseas/testdata/{price_detail,search_info,daily_price,daily_chart_price,updown_rate}_success.json; do
  python3 -c "import json; json.load(open('$f'))" && echo "$f OK" || echo "$f BROKEN"
done
```
Expected: 5 lines all `OK`.

- [ ] **Step 8: Commit**

```bash
git add overseas/testdata/{price_detail,search_info,daily_price,daily_chart_price,updown_rate}_success.json overseas/testdata/README.md
git commit -m "$(cat <<'EOF'
[chore] Phase 1.5 testdata — 5 KIS API 합성 JSON fixtures

해외주식 5 메서드 (price_detail, search_info, daily_price,
daily_chart_price, updown_rate) testdata. 한투 docs (docs/api/해외주식/)
응답 필드 정의 기반 합성. AAPL/SPX/NVDA/TSLA 샘플.

해외 마스터 sample (Task 2 에서 추가 예정).

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 2: internal/overseasmaster — 해외 마스터 파서

**Files:**
- Create: `internal/overseasmaster/doc.go`
- Create: `internal/overseasmaster/overseasmaster.go`
- Create: `internal/overseasmaster/overseasmaster_test.go`
- Create: `internal/overseasmaster/testdata/<market>_code_sample.cod.zip` (실제 KIS sample)
- Create: `internal/overseasmaster/testdata/README.md`
- Create: `overseas/testdata/<market>_code_sample.cod.zip` (위 파일 복사)

> 해외 마스터 파일 형식이 KIS docs 에 명시 안 됨 — implementation 시점에 sample 다운로드 + 분석 후 결정. 가설: KRX 와 같은 cp949 + fwf (또는 변형). 파일명 패턴: `<market>mst.cod.zip` (예: `nasmst.cod.zip`, `nysmst.cod.zip`).

- [ ] **Step 1: KIS sample 다운로드 + 형식 분석**

```bash
mkdir -p /tmp/overseas_master && cd /tmp/overseas_master
for m in nas nys ams shs shi szs szi tse hks hnx hsx; do
  curl -sSL -o "${m}mst.cod.zip" "https://new.real.download.dws.co.kr/common/master/${m}mst.cod.zip"
  ls -la "${m}mst.cod.zip" 2>&1
done
```

If any URL returns 404 or empty, STOP and report — Phase 1.5 의 마스터 다운로드 가정이 깨진 것. 다른 패턴 (e.g., `nasdaq.mst.zip`) 시도 또는 `https://new.real.download.dws.co.kr/common/master/` index 페이지 직접 확인.

성공 시: nas (nasdaq) sample 추출 + 형식 분석:
```bash
cd /tmp/overseas_master
unzip -o nasmst.cod.zip
file nasmst.cod  # 형식 확인 (text vs binary, encoding)
head -c 500 nasmst.cod | iconv -f cp949 -t utf-8 2>&1 | head -3   # cp949 시도
# 또는 head -c 500 nasmst.cod | iconv -f utf-8 -t utf-8  (UTF-8 시도)
# 또는 head -c 500 nasmst.cod | hexdump -C | head -20  (binary 확인)
```

분석 결과 형식 (cp949+fwf 인지, csv 인지, 다른 형식인지) 기반으로 파서 결정.

- [ ] **Step 2: doc.go**

```go
// Package overseasmaster 는 KIS 공개 다운로드의 해외 거래소 마스터 파일
// (NASDAQ/NYSE/AMEX/홍콩/일본 등 11 거래소) 의 디코딩/파싱 로직.
//
// 한투 API 가 아니라 KIS 가 공개 다운로드로 제공하는 .cod.zip 파일을 처리.
// KRX 마스터와 형식 다를 수 있음 (Step 1 분석 결과에 따라 코덱 결정).
//
// 사용자에게 노출되지 않는 internal 패키지. overseas 패키지의 FetchOverseasSymbols 가 호출.
package overseasmaster
```

- [ ] **Step 3: 테스트 작성** — `internal/overseasmaster/overseasmaster_test.go`

```go
package overseasmaster

import (
	"os"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse_NAS(t *testing.T) {
	zipBytes, err := os.ReadFile("testdata/nas_code_sample.cod.zip")
	require.NoError(t, err)

	syms, err := Parse("nas", zipBytes)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(syms), 1, "nas sample 은 최소 1 행 포함")

	// 종목 코드는 영문 대문자 + 숫자 (NASDAQ 형식)
	symRe := regexp.MustCompile(`^[A-Z0-9.]+$`)
	for i, s := range syms {
		assert.True(t, symRe.MatchString(s.Symbol), "row %d: Symbol %q 는 영문/숫자/점", i, s.Symbol)
		assert.NotEmpty(t, s.EnglishName, "row %d: EnglishName 비어있음", i)
		assert.NotNil(t, s.Raw, "row %d: Raw map 비어있지 않음", i)
	}
}

func TestParse_InvalidMarket(t *testing.T) {
	_, err := Parse("invalid", []byte{})
	assert.Error(t, err)
}

func TestParse_EmptyZip(t *testing.T) {
	_, err := Parse("nas", nil)
	assert.Error(t, err)
}
```

- [ ] **Step 4: 테스트 실행 → FAIL**

`go test ./internal/overseasmaster/... -v` — 컴파일 실패.

- [ ] **Step 5: 구현 — `internal/overseasmaster/overseasmaster.go`**

Step 1 의 형식 분석 결과를 반영. KRX 와 동일 cp949+fwf 라면 `internal/krxmaster` 의 helper (decodeCP949, parseFwf) 재사용 가능 — 단, internal package 끼리 import 안 되니 코드 중복 OR `internal/encoding/krxfwf/` 같은 공통 패키지로 추출. 가장 단순: 코드 중복 (helper 5 줄 정도).

다음은 cp949+fwf 가정한 skeleton (Step 1 분석 결과에 따라 fwf widths/columns 조정):

```go
package overseasmaster

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"strings"

	"golang.org/x/text/encoding/korean"
	"golang.org/x/text/transform"
)

// MarketURLs 는 KIS 공개 다운로드의 11 거래소 마스터 파일 URL.
var MarketURLs = map[string]string{
	"nas": "https://new.real.download.dws.co.kr/common/master/nasmst.cod.zip",
	"nys": "https://new.real.download.dws.co.kr/common/master/nysmst.cod.zip",
	"ams": "https://new.real.download.dws.co.kr/common/master/amsmst.cod.zip",
	"shs": "https://new.real.download.dws.co.kr/common/master/shsmst.cod.zip",
	"shi": "https://new.real.download.dws.co.kr/common/master/shimst.cod.zip",
	"szs": "https://new.real.download.dws.co.kr/common/master/szsmst.cod.zip",
	"szi": "https://new.real.download.dws.co.kr/common/master/szimst.cod.zip",
	"tse": "https://new.real.download.dws.co.kr/common/master/tsemst.cod.zip",
	"hks": "https://new.real.download.dws.co.kr/common/master/hksmst.cod.zip",
	"hnx": "https://new.real.download.dws.co.kr/common/master/hnxmst.cod.zip",
	"hsx": "https://new.real.download.dws.co.kr/common/master/hsxmst.cod.zip",
}

// Symbol 은 해외 거래소 마스터 한 행 (typed 핵심 필드 + Raw map fallback).
//
// Step 1 분석 결과에 따라 typed field 와 fwf widths 결정. 다음은 일반적 후보 — 실제 sample 확인 후 조정.
type Symbol struct {
	Symbol      string // 종목코드 (예 "AAPL")
	EnglishName string // 영문 종목명
	KoreanName  string // 한글 종목명 (KIS 가 한글 컬럼 제공 시)
	ISIN        string // 국제증권식별번호 (있는 경우)
	Currency    string // 거래 통화 (USD, HKD, JPY, ...)

	Raw map[string]string // 모든 컬럼 (한글 컬럼명 → 값)
}

// Parse 는 market ("nas"/"nys"/"ams" 등) 별 마스터 ZIP byte 를 파싱.
//
// Step 1 분석 결과에 따라 코덱 (cp949/utf-8/csv) 분기 또는 통합.
func Parse(market string, zipBytes []byte) ([]Symbol, error) {
	if _, ok := MarketURLs[market]; !ok {
		return nil, fmt.Errorf("overseasmaster: unknown market %q", market)
	}
	mst, err := openCodFromZip(zipBytes)
	if err != nil {
		return nil, fmt.Errorf("overseasmaster: %s: %w", market, err)
	}
	decoded, err := decodeCP949(mst)
	if err != nil {
		return nil, fmt.Errorf("overseasmaster: %s: cp949: %w", market, err)
	}

	var out []Symbol
	for _, line := range strings.Split(decoded, "\n") {
		line = strings.TrimRight(line, "\r")
		if line == "" {
			continue
		}
		// Step 1 분석 결과 기반: fwf widths 또는 csv split.
		// 다음은 placeholder — implementation 시 정확한 widths/separator 적용.
		raw := make(map[string]string)
		// raw["..."] = ...
		out = append(out, Symbol{
			Symbol:      "", // raw["종목코드"] 등
			EnglishName: "",
			Raw:         raw,
		})
	}
	return out, nil
}

func openCodFromZip(zipBytes []byte) ([]byte, error) {
	zr, err := zip.NewReader(bytes.NewReader(zipBytes), int64(len(zipBytes)))
	if err != nil {
		return nil, fmt.Errorf("zip open: %w", err)
	}
	for _, f := range zr.File {
		if strings.HasSuffix(f.Name, ".cod") || strings.HasSuffix(f.Name, ".mst") {
			rc, err := f.Open()
			if err != nil {
				return nil, fmt.Errorf("zip read %s: %w", f.Name, err)
			}
			defer rc.Close()
			return io.ReadAll(rc)
		}
	}
	return nil, fmt.Errorf(".cod or .mst file not found in zip")
}

func decodeCP949(b []byte) (string, error) {
	decoded, _, err := transform.Bytes(korean.EUCKR.NewDecoder(), b)
	if err != nil {
		return "", err
	}
	return string(decoded), nil
}
```

> **NOTE for implementer**: Step 1 분석 결과에 따라 위 skeleton 의 `Parse` 함수 내부 (raw map 채우기 + Symbol fields 매핑) 를 정확히 작성. KIS 가 csv 형식이라면 fwf 대신 csv parser 사용. 만약 형식이 KRX 와 동일하면 `internal/krxmaster` 의 `parseFwf` 같은 helper 재사용 (내부 패키지라 import 불가 — 코드 중복 OR internal common package).

- [ ] **Step 6: 테스트 실행 → PASS**

`go test ./internal/overseasmaster/... -v` — 3 cases PASS.

- [ ] **Step 7: testdata 디렉터리 + sample 파일 복사**

```bash
mkdir -p internal/overseasmaster/testdata overseas/testdata
# nas (nasdaq) 첫 3 행 sample 만 commit
LC_ALL=C head -n 3 /tmp/overseas_master/nasmst.cod > /tmp/nas_code_sample.cod
zip /tmp/overseas_master/nas_code_sample.cod.zip /tmp/nas_code_sample.cod
cp /tmp/overseas_master/nas_code_sample.cod.zip internal/overseasmaster/testdata/
cp /tmp/overseas_master/nas_code_sample.cod.zip overseas/testdata/
```

- [ ] **Step 8: testdata README**

`internal/overseasmaster/testdata/README.md`:

```markdown
# overseasmaster testdata

`nas_code_sample.cod.zip` 는 KIS 공개 다운로드의 NASDAQ 마스터 파일 첫 3 행 sample.

## 출처

- NASDAQ 마스터: https://new.real.download.dws.co.kr/common/master/nasmst.cod.zip
- 다른 거래소: https://new.real.download.dws.co.kr/common/master/{nys,ams,shs,shi,szs,szi,tse,hks,hnx,hsx}mst.cod.zip

KIS 가 공개 다운로드로 제공. `internal/overseasmaster` 의 파서가 실제 KIS byte 와 호환되는지 검증하기 위한 단위 테스트 sample.

## 재생성 방법

```bash
cd /tmp && rm -rf overseas_master && mkdir overseas_master && cd overseas_master
curl -sSL -o nasmst.cod.zip https://new.real.download.dws.co.kr/common/master/nasmst.cod.zip
unzip -o nasmst.cod.zip
LC_ALL=C head -n 3 nasmst.cod > nas_code_sample.cod
zip nas_code_sample.cod.zip nas_code_sample.cod
```

## 라이선스

KIS 가 무료 공개 다운로드로 배포. 본 sample 은 학습/테스트 용도.
```

- [ ] **Step 9: Commit**

```bash
git add internal/overseasmaster overseas/testdata/nas_code_sample.cod.zip
git commit -m "$(cat <<'EOF'
[feat] internal/overseasmaster — 해외 거래소 마스터 파서 (NASDAQ 등 11 거래소)

Symbol struct + Parse(market, zipBytes). 11 거래소 (nas/nys/ams/shs/shi/
szs/szi/tse/hks/hnx/hsx) 통합 마스터 다운로드 + 파싱. KIS 공개 다운로드
URL (https://new.real.download.dws.co.kr/common/master/<market>mst.cod.zip).
실제 sample (NASDAQ 첫 3행) 으로 cp949 디코딩 검증.

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 3: overseas.Client 시그니처 확장 + root wireInfra 갱신

**Files:**
- Modify: `overseas/client.go`
- Modify: `client.go` (root)

- [ ] **Step 1: overseas/client.go — Client 에 master 필드 추가**

기존 `overseas/client.go` 를 다음으로 교체:

```go
package overseas

import (
	"github.com/kenshin579/korea-investment-stock/internal/httpclient"
	"github.com/kenshin579/korea-investment-stock/internal/mastercache"
)

// Client 는 해외주식 API sub-client.
//
// 사용자는 직접 생성하지 않고 kis.Client.Overseas 으로 접근.
type Client struct {
	http   *httpclient.Client
	master *mastercache.Cache // 해외 마스터 파일 디스크 캐시 (FetchOverseasSymbols 가 사용)
}

// New 는 internal 용도. root kis.NewClient 가 호출.
func New(http *httpclient.Client, master *mastercache.Cache) *Client {
	return &Client{http: http, master: master}
}
```

- [ ] **Step 2: client.go (root) — wireInfra 갱신**

`client.go` (root) 의 `c.Overseas = overseas.New(c.httpClient)` 라인을 다음으로 교체:

```go
	c.Overseas = overseas.New(c.httpClient, c.masterC)
```

- [ ] **Step 3: 빌드/테스트 회귀 검증**

```bash
go build ./... && go test ./... -count=1
```
Expected: 모든 패키지 PASS.

- [ ] **Step 4: Commit**

```bash
git add overseas/client.go client.go
git commit -m "$(cat <<'EOF'
[refactor] overseas.Client 에 master *mastercache.Cache 필드 추가

Phase 1.5 의 FetchOverseasSymbols 가 마스터 캐시에 접근하기 위한 인프라.
overseas.New(http, master) 시그니처 변경 — internal use 만이라 BC 영향 없음.
root client.go 의 wireInfra 가 c.masterC 주입 (Phase 1.2 의 domestic 패턴).

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 4: overseas/testhelper_test.go

**Files:**
- Create: `overseas/testhelper_test.go`

> `domestic/testhelper_test.go` 와 동일 패턴, package 만 `overseas_test`.

- [ ] **Step 1: testhelper 작성**

```go
package overseas_test

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"

	"github.com/kenshin579/korea-investment-stock/internal/httpclient"
	"github.com/kenshin579/korea-investment-stock/internal/mastercache"
	"github.com/kenshin579/korea-investment-stock/internal/ratelimit"
	"github.com/kenshin579/korea-investment-stock/overseas"
)

const testBaseURL = "https://openapi.koreainvestment.com:9443"

func loadFixture(t *testing.T, name string) []byte {
	t.Helper()
	b, err := os.ReadFile(filepath.Join("testdata", name))
	require.NoError(t, err)
	return b
}

func loadFixtureString(t *testing.T, name string) string {
	return string(loadFixture(t, name))
}

type stubTokenManager struct{}

func (stubTokenManager) Get(ctx context.Context) (string, error)     { return "Bearer test", nil }
func (stubTokenManager) Refresh(ctx context.Context) (string, error) { return "Bearer test", nil }

func newTestClient(t *testing.T) *overseas.Client {
	t.Helper()
	httpClient := &http.Client{Transport: httpmock.DefaultTransport}
	httpcli := httpclient.New(httpclient.Config{
		BaseURL:    testBaseURL,
		AppKey:     "test-key",
		AppSecret:  "test-secret",
		AccountNo:  "00000000-00",
		Limiter:    ratelimit.New(1000),
		TokenMgr:   stubTokenManager{},
		Retries:    0,
		Timeout:    5 * time.Second,
		HTTPClient: httpClient,
	})
	master := mastercache.New(t.TempDir(), time.Hour)
	return overseas.New(httpcli, master)
}
```

- [ ] **Step 2: 컴파일 검증**

`go test ./overseas/... -run NoSuchTest -v` — 컴파일 OK (0 tests).

- [ ] **Step 3: Commit**

```bash
git add overseas/testhelper_test.go
git commit -m "$(cat <<'EOF'
[test] overseas — 공통 테스트 helper

newTestClient + loadFixture + stubTokenManager (overseas 패키지 전용,
domestic/testhelper_test.go 와 분리). Phase 1.5 의 메서드 테스트가 공통 사용.

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 5: overseas/symbols.go — FetchOverseasSymbols(market)

**Files:**
- Create: `overseas/symbols.go`
- Create: `overseas/symbols_test.go`

- [ ] **Step 1: 테스트 작성**

```go
package overseas_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kenshin579/korea-investment-stock/internal/overseasmaster"
)

func TestClient_FetchOverseasSymbols_NAS(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	zipBytes := loadFixture(t, "nas_code_sample.cod.zip")
	httpmock.RegisterResponder(http.MethodGet, overseasmaster.MarketURLs["nas"],
		httpmock.NewBytesResponder(200, zipBytes))

	c := newTestClient(t)
	syms, err := c.FetchOverseasSymbols(context.Background(), "nas")
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(syms), 1)
	for _, s := range syms {
		assert.NotEmpty(t, s.Symbol)
	}
}

func TestClient_FetchOverseasSymbols_UnknownMarket(t *testing.T) {
	c := newTestClient(t)
	_, err := c.FetchOverseasSymbols(context.Background(), "invalid")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unknown market")
}

func TestClient_FetchOverseasSymbols_DownloadError(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(http.MethodGet, overseasmaster.MarketURLs["nas"],
		httpmock.NewStringResponder(500, "internal error"))

	c := newTestClient(t)
	_, err := c.FetchOverseasSymbols(context.Background(), "nas")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "HTTP 500")
}
```

- [ ] **Step 2: FAIL**

- [ ] **Step 3: 구현 — `overseas/symbols.go`**

```go
package overseas

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/kenshin579/korea-investment-stock/internal/overseasmaster"
)

// OverseasSymbol 은 internal/overseasmaster 의 type alias (외부 사용자 노출).
type OverseasSymbol = overseasmaster.Symbol

// FetchOverseasSymbols 는 KIS 공개 마스터 (`<market>mst.cod.zip`) 를 다운로드/캐시 후 파싱.
//
// 한투 REST API 가 아니라 KIS 가 공개 다운로드로 제공. 토큰 인증 불필요.
// 마스터 파일은 mastercache 에 디스크 캐시 (default TTL 7일).
//
// market: "nas"(NASDAQ)/"nys"(NYSE)/"ams"(AMEX)/"shs"(상해)/"shi"(상해지수)/
// "szs"(심천)/"szi"(심천지수)/"tse"(도쿄)/"hks"(홍콩)/"hnx"(하노이)/"hsx"(호치민)
func (c *Client) FetchOverseasSymbols(ctx context.Context, market string) ([]OverseasSymbol, error) {
	url, ok := overseasmaster.MarketURLs[market]
	if !ok {
		return nil, fmt.Errorf("overseas: unknown market %q (valid: nas/nys/ams/shs/shi/szs/szi/tse/hks/hnx/hsx)", market)
	}
	cacheName := market + "mst.cod.zip"
	raw, err := c.master.Get(ctx, cacheName, func(ctx context.Context) ([]byte, error) {
		return downloadURL(ctx, url)
	})
	if err != nil {
		return nil, err
	}
	return overseasmaster.Parse(market, raw)
}

// downloadURL 은 KIS 공개 마스터 파일 단순 GET. 한투 API transport 와 분리 의도로
// http.DefaultClient 사용 (KIS 마스터 도메인은 한투 API 가 아니므로 토큰/proxy 정책 다름).
func downloadURL(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("overseas master new request: %w", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("overseas master %s: %w", url, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("overseas master %s: HTTP %d", url, resp.StatusCode)
	}
	return io.ReadAll(resp.Body)
}
```

- [ ] **Step 4: PASS**

`go test ./overseas/... -run FetchOverseas -v` — 3 cases PASS.

- [ ] **Step 5: Commit**

```bash
git add overseas/symbols.go overseas/symbols_test.go
git commit -m "$(cat <<'EOF'
[feat] overseas — FetchOverseasSymbols(market)

11 거래소 (NASDAQ/NYSE/AMEX/상해/심천/도쿄/홍콩/하노이/호치민/지수)
통합 메서드. mastercache 디스크 캐시 (default TTL 7일) + overseasmaster
파서. http.DefaultClient 로 KIS 마스터 도메인 직접 호출. unknown market /
HTTP error 검증 테스트 포함.

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 6: overseas/price.go — InquirePriceDetail

**Files:**
- Create: `overseas/price.go`
- Create: `overseas/price_test.go`

- [ ] **Step 1: 테스트 작성**

```go
package overseas_test

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kenshin579/korea-investment-stock/overseas"
)

func TestClient_InquirePriceDetail(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/quotations/price-detail`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "price_detail_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquirePriceDetail(context.Background(), overseas.InquirePriceDetailParams{
		Excd: "NAS",
		Symb: "AAPL",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, "", capturedQuery.Get("AUTH"))
	assert.Equal(t, "NAS", capturedQuery.Get("EXCD"))
	assert.Equal(t, "AAPL", capturedQuery.Get("SYMB"))

	d, _ := decimal.NewFromString("181.45")
	assert.True(t, d.Equal(res.Output.Last))
	assert.Equal(t, "USD", res.Output.Curr)
	assert.Equal(t, int64(85000000), res.Output.Tvol)
	assert.InDelta(t, 29.45, res.Output.Perx, 0.001)
}
```

- [ ] **Step 2: FAIL**

- [ ] **Step 3: 구현 — `overseas/price.go`**

```go
package overseas

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/shopspring/decimal"

	"github.com/kenshin579/korea-investment-stock/internal/httpclient"
)

// PriceDetail 은 해외주식_현재가상세 (HHDFS76200200) 응답.
//
// 한투 docs: docs/api/해외주식/해외주식_현재가상세.md
// path: /uapi/overseas-price/v1/quotations/price-detail
type PriceDetail struct {
	Output PriceDetailSnapshot `json:"output"`
}

// PriceDetailSnapshot 은 응답의 output (단일 객체, ~46 fields).
type PriceDetailSnapshot struct {
	Rsym       string          `json:"rsym"`             // 실시간조회종목코드
	Pvol       int64           `json:"pvol,string"`      // 전일거래량
	Open       decimal.Decimal `json:"open"`             // 시가
	High       decimal.Decimal `json:"high"`             // 고가
	Low        decimal.Decimal `json:"low"`              // 저가
	Last       decimal.Decimal `json:"last"`             // 현재가
	Base       decimal.Decimal `json:"base"`             // 전일종가
	Tomv       int64           `json:"tomv,string"`      // 시가총액
	Pamt       int64           `json:"pamt,string"`      // 전일거래대금
	Uplp       decimal.Decimal `json:"uplp"`             // 상한가
	Dnlp       decimal.Decimal `json:"dnlp"`             // 하한가
	H52p       decimal.Decimal `json:"h52p"`             // 52주최고가
	H52d       string          `json:"h52d"`             // 52주최고일자
	L52p       decimal.Decimal `json:"l52p"`             // 52주최저가
	L52d       string          `json:"l52d"`             // 52주최저일자
	Perx       float64         `json:"perx,string"`      // PER
	Pbrx       float64         `json:"pbrx,string"`      // PBR
	Epsx       decimal.Decimal `json:"epsx"`             // EPS
	Bpsx       decimal.Decimal `json:"bpsx"`             // BPS
	Shar       int64           `json:"shar,string"`      // 상장주수
	Mcap       int64           `json:"mcap,string"`      // 자본금
	Curr       string          `json:"curr"`             // 통화
	Zdiv       string          `json:"zdiv"`             // 소수점자리수
	Vnit       string          `json:"vnit"`             // 매매단위
	TXprc      decimal.Decimal `json:"t_xprc"`           // 원환산당일가격
	TXdif      decimal.Decimal `json:"t_xdif"`           // 원환산당일대비
	TXrat      float64         `json:"t_xrat,string"`    // 원환산당일등락
	PXprc      decimal.Decimal `json:"p_xprc"`           // 원환산전일가격
	PXdif      decimal.Decimal `json:"p_xdif"`           // 원환산전일대비
	PXrat      float64         `json:"p_xrat,string"`    // 원환산전일등락
	TRate      float64         `json:"t_rate,string"`    // 당일환율
	PRate      float64         `json:"p_rate,string"`    // 전일환율
	TXsgn      string          `json:"t_xsgn"`           // 원환산당일기호
	PXsng      string          `json:"p_xsng"`           // 원환산전일기호
	EOrdyn     string          `json:"e_ordyn"`          // 거래가능여부
	EHogau     string          `json:"e_hogau"`          // 호가단위
	EIcod      string          `json:"e_icod"`           // 업종(섹터)
	EParp      decimal.Decimal `json:"e_parp"`           // 액면가
	Tvol       int64           `json:"tvol,string"`      // 거래량
	Tamt       int64           `json:"tamt,string"`      // 거래대금
	EtypNm     string          `json:"etyp_nm"`          // ETP 분류명
}

// InquirePriceDetailParams 는 해외주식_현재가상세 조회 파라미터.
type InquirePriceDetailParams struct {
	Auth string // AUTH — 사용자권한정보. 빈 값 default
	Excd string // EXCD — 거래소명 (HKS/NYS/NAS/AMS/TSE/SHS/SZS/SHI/SZI/HSX/HNX/BAY/BAQ/BAA)
	Symb string // SYMB — 종목코드
}

// InquirePriceDetail 은 해외주식_현재가상세 호출.
//
// 한투 docs: docs/api/해외주식/해외주식_현재가상세.md
// path: /uapi/overseas-price/v1/quotations/price-detail (HHDFS76200200)
func (c *Client) InquirePriceDetail(ctx context.Context, params InquirePriceDetailParams) (*PriceDetail, error) {
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/overseas-price/v1/quotations/price-detail",
		TrID:   "HHDFS76200200",
		Query: map[string]string{
			"AUTH": params.Auth,
			"EXCD": params.Excd,
			"SYMB": params.Symb,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}

	var res PriceDetail
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse PriceDetail: %w", err)
	}
	return &res, nil
}
```

- [ ] **Step 4: PASS**

- [ ] **Step 5: Commit**

```bash
git add overseas/price.go overseas/price_test.go
git commit -m "$(cat <<'EOF'
[feat] overseas — InquirePriceDetail (해외주식 현재가상세, HHDFS76200200)

PriceDetail + PriceDetailSnapshot (~40 필드) + InquirePriceDetailParams
(AUTH/EXCD/SYMB). 가격/52주/PER/PBR/원환산/환율 모두 포함. EXCD 11 거래소
지원 (NAS/NYS/AMS/TSE/SHS/SZS/SHI/SZI/HSX/HNX/HKS).

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 7: overseas/search.go — SearchInfo

**Files:**
- Create: `overseas/search.go`
- Create: `overseas/search_test.go`

- [ ] **Step 1: 테스트 작성**

```go
package overseas_test

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kenshin579/korea-investment-stock/overseas"
)

func TestClient_SearchInfo(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/quotations/search-info`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "search_info_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.SearchInfo(context.Background(), overseas.SearchInfoParams{
		PrdtTypeCD: "512",
		Pdno:       "AAPL",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, "512", capturedQuery.Get("PRDT_TYPE_CD"))
	assert.Equal(t, "AAPL", capturedQuery.Get("PDNO"))

	assert.Equal(t, "APPLE INC", res.Output.PrdtEngName)
	assert.Equal(t, "미국", res.Output.NatnName)
	assert.Equal(t, "NASD", res.Output.OvrsExcgCd)
	assert.Equal(t, "USD", res.Output.TrCrcyCd)
	assert.Equal(t, "STK", res.Output.PrdtClsfCd)
	assert.Equal(t, "Y", res.Output.LstgYn)
}
```

- [ ] **Step 2: FAIL**

- [ ] **Step 3: 구현 — `overseas/search.go`**

```go
package overseas

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kenshin579/korea-investment-stock/internal/httpclient"
)

// OverseasProductInfo 는 해외주식_상품기본정보 (CTPF1702R) 응답.
//
// 한투 docs: docs/api/해외주식/해외주식_상품기본정보.md
// path: /uapi/overseas-price/v1/quotations/search-info
//
// domestic.ProductInfo 와 다른 패키지/타입 — 해외 거래소 메타정보 풍부.
type OverseasProductInfo struct {
	Output OverseasProductInfoOutput `json:"output"`
}

// OverseasProductInfoOutput 은 응답의 output (단일 객체, ~50 fields).
//
// KIS docs 의 모든 필드 1:1 매핑.
type OverseasProductInfoOutput struct {
	StdPdno                 string `json:"std_pdno"`                    // 표준상품번호 (ISIN)
	PrdtEngName             string `json:"prdt_eng_name"`               // 상품영문명
	NatnCd                  string `json:"natn_cd"`                     // 국가코드
	NatnName                string `json:"natn_name"`                   // 국가명
	TrMketCd                string `json:"tr_mket_cd"`                  // 거래시장코드
	TrMketName              string `json:"tr_mket_name"`                // 거래시장명
	OvrsExcgCd              string `json:"ovrs_excg_cd"`                // 해외거래소코드
	OvrsExcgName            string `json:"ovrs_excg_name"`              // 해외거래소명
	TrCrcyCd                string `json:"tr_crcy_cd"`                  // 거래통화코드
	OvrsPapr                string `json:"ovrs_papr"`                   // 해외액면가 (string — 형식 다양)
	CrcyName                string `json:"crcy_name"`                   // 통화명
	OvrsStckDvsnCd          string `json:"ovrs_stck_dvsn_cd"`           // 해외주식구분코드 (01.주식 02.WARRANT 03.ETF 04.우선주)
	PrdtClsfCd              string `json:"prdt_clsf_cd"`                // 상품분류코드
	PrdtClsfName            string `json:"prdt_clsf_name"`              // 상품분류명
	SllUnitQty              string `json:"sll_unit_qty"`                // 매도단위수량
	BuyUnitQty              string `json:"buy_unit_qty"`                // 매수단위수량
	TrUnitAmt               string `json:"tr_unit_amt"`                 // 거래단위금액
	LstgStckNum             int64  `json:"lstg_stck_num,string"`        // 상장주식수
	LstgDt                  string `json:"lstg_dt"`                     // 상장일자
	OvrsStckTrStopDvsnCd    string `json:"ovrs_stck_tr_stop_dvsn_cd"`   // 해외주식거래정지구분코드
	LstgAbolItemYn          string `json:"lstg_abol_item_yn"`           // 상장폐지종목여부
	OvrsStckPrdtGrpNo       string `json:"ovrs_stck_prdt_grp_no"`       // 해외주식상품그룹번호
	LstgYn                  string `json:"lstg_yn"`                     // 상장여부
	TaxLevyYn               string `json:"tax_levy_yn"`                 // 세금징수여부
	OvrsStckErlmRosnCd      string `json:"ovrs_stck_erlm_rosn_cd"`      // 해외주식등록사유코드
	OvrsStckHistRghtDvsnCd  string `json:"ovrs_stck_hist_rght_dvsn_cd"` // 해외주식이력권리구분코드
	ChngBfPdno              string `json:"chng_bf_pdno"`                // 변경전상품번호
	PrdtTypeCd2             string `json:"prdt_type_cd_2"`              // 상품유형코드2
	OvrsItemName            string `json:"ovrs_item_name"`              // 해외종목명
	SedolNo                 string `json:"sedol_no"`                    // SEDOL번호
	BlbgTckrText            string `json:"blbg_tckr_text"`              // 블름버그티커내용
	OvrsStckEtfRiskDrtpCd   string `json:"ovrs_stck_etf_risk_drtp_cd"`  // ETF위험지표코드
	EtpChasErngRtDbnb       string `json:"etp_chas_erng_rt_dbnb"`       // ETP추적수익율배수
	IsttUsgeIsinCd          string `json:"istt_usge_isin_cd"`           // 기관용도ISIN코드
	MintSvcYn               string `json:"mint_svc_yn"`                 // MINT서비스여부
	MintSvcYnChngDt         string `json:"mint_svc_yn_chng_dt"`         // MINT서비스여부변경일자
	PrdtName                string `json:"prdt_name"`                   // 상품명 (한글)
	LeiCd                   string `json:"lei_cd"`                      // LEI코드
	OvrsStckStopRsonCd      string `json:"ovrs_stck_stop_rson_cd"`      // 해외주식정지사유코드
	LstgAbolDt              string `json:"lstg_abol_dt"`                // 상장폐지일자
	MiniStkTrStatDvsnCd     string `json:"mini_stk_tr_stat_dvsn_cd"`    // 미니스탁거래상태구분코드
	MintFrstSvcErlmDt       string `json:"mint_frst_svc_erlm_dt"`       // MINT최초서비스등록일자
	MintDcptTradPsblYn      string `json:"mint_dcpt_trad_psbl_yn"`      // MINT소수점매매가능여부
	MintFnumTradPsblYn      string `json:"mint_fnum_trad_psbl_yn"`      // MINT정수매매가능여부
	MintCblcCvsnIpsbYn      string `json:"mint_cblc_cvsn_ipsb_yn"`      // MINT잔고전환불가여부
	PtpItemYn               string `json:"ptp_item_yn"`                 // PTP종목여부
	PtpItemTrfxExmtYn       string `json:"ptp_item_trfx_exmt_yn"`       // PTP종목양도세면제여부
	PtpItemTrfxExmtStrtDt   string `json:"ptp_item_trfx_exmt_strt_dt"`  // PTP양도세면제시작일자
	PtpItemTrfxExmtEndDt    string `json:"ptp_item_trfx_exmt_end_dt"`   // PTP양도세면제종료일자
	DtmTrPsblYn             string `json:"dtm_tr_psbl_yn"`              // 주간거래가능여부
	SdrfStopEclsYn          string `json:"sdrf_stop_ecls_yn"`           // 급등락정지제외여부
	SdrfStopEclsErlmDt      string `json:"sdrf_stop_ecls_erlm_dt"`      // 급등락정지제외등록일자
	MemoText1               string `json:"memo_text1"`                  // 메모내용1
	OvrsNowPric1            string `json:"ovrs_now_pric1"`              // 해외현재가격1 (string — 23.5 형식)
	LastRcvgDtime           string `json:"last_rcvg_dtime"`             // 최종수신일시
}

// SearchInfoParams 는 해외주식_상품기본정보 조회 파라미터.
type SearchInfoParams struct {
	PrdtTypeCD string // PRDT_TYPE_CD — 512(나스닥)/513(뉴욕)/529(아멕스)/515(일본)/501(홍콩)/543(홍콩CNY)/558(홍콩USD)/507(베트남 하노이)/508(베트남 호치민)/551(중국 상해A)/552(중국 심천A)
	Pdno       string // PDNO — 상품번호 (예 "AAPL")
}

// SearchInfo 는 해외주식_상품기본정보 호출.
//
// 한투 docs: docs/api/해외주식/해외주식_상품기본정보.md
// path: /uapi/overseas-price/v1/quotations/search-info (CTPF1702R)
func (c *Client) SearchInfo(ctx context.Context, params SearchInfoParams) (*OverseasProductInfo, error) {
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/overseas-price/v1/quotations/search-info",
		TrID:   "CTPF1702R",
		Query: map[string]string{
			"PRDT_TYPE_CD": params.PrdtTypeCD,
			"PDNO":         params.Pdno,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}

	var res OverseasProductInfo
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse OverseasProductInfo: %w", err)
	}
	return &res, nil
}
```

- [ ] **Step 4: PASS**

- [ ] **Step 5: Commit**

```bash
git add overseas/search.go overseas/search_test.go
git commit -m "$(cat <<'EOF'
[feat] overseas — SearchInfo (해외주식 상품기본정보, CTPF1702R)

OverseasProductInfo + OverseasProductInfoOutput (~52 필드) + SearchInfoParams
(PRDT_TYPE_CD/PDNO). 거래소/통화/SEDOL/Bloomberg/MINT/PTP 등 메타데이터.
domestic.ProductInfo 와 다른 패키지라 이름 충돌 없음.

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 8: overseas/chart.go — InquireDailyPrice + InquireDailyChartPrice

**Files:**
- Create: `overseas/chart.go`
- Create: `overseas/chart_test.go`

> 두 메서드 한 file 에 (Phase 1.2 의 domestic/chart.go 와 같은 패턴 — 일봉 + 분봉).

- [ ] **Step 1: 테스트 작성**

```go
package overseas_test

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kenshin579/korea-investment-stock/overseas"
)

func TestClient_InquireDailyPrice(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/quotations/dailyprice`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "daily_price_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireDailyPrice(context.Background(), overseas.InquireDailyPriceParams{
		Excd: "NAS",
		Symb: "AAPL",
		Bymd: "20260505",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, "NAS", capturedQuery.Get("EXCD"))
	assert.Equal(t, "AAPL", capturedQuery.Get("SYMB"))
	assert.Equal(t, "0", capturedQuery.Get("GUBN")) // default 일
	assert.Equal(t, "0", capturedQuery.Get("MODP")) // default 미반영

	assert.Equal(t, "DNASAAPL", res.Output1.Rsym)
	require.Len(t, res.Output2, 2)
	assert.Equal(t, "20260505", res.Output2[0].Xymd)
	d, _ := decimal.NewFromString("181.45")
	assert.True(t, d.Equal(res.Output2[0].Clos))
	assert.Equal(t, int64(85000000), res.Output2[0].Tvol)
}

func TestClient_InquireDailyChartPrice(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/quotations/inquire-daily-chartprice`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "daily_chart_price_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireDailyChartPrice(context.Background(), overseas.InquireDailyChartPriceParams{
		MarketCode: "N",
		Symbol:     "SPX",
		FromDate:   "20260101",
		ToDate:     "20260505",
		Period:     "D",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, "N", capturedQuery.Get("FID_COND_MRKT_DIV_CODE"))
	assert.Equal(t, "SPX", capturedQuery.Get("FID_INPUT_ISCD"))
	assert.Equal(t, "20260101", capturedQuery.Get("FID_INPUT_DATE_1"))
	assert.Equal(t, "20260505", capturedQuery.Get("FID_INPUT_DATE_2"))
	assert.Equal(t, "D", capturedQuery.Get("FID_PERIOD_DIV_CODE"))

	assert.Equal(t, "S&P 500", res.Output1.HtsKorIsnm)
	require.Len(t, res.Output2, 2)
	assert.Equal(t, "20260505", res.Output2[0].StckBsopDate)
}
```

- [ ] **Step 2: FAIL**

- [ ] **Step 3: 구현 — `overseas/chart.go`**

```go
package overseas

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/shopspring/decimal"

	"github.com/kenshin579/korea-investment-stock/internal/httpclient"
)

// DailyPrice 는 해외주식_기간별시세 (HHDFS76240000) 응답.
type DailyPrice struct {
	Output1 DailyPriceSummary `json:"output1"`
	Output2 []DailyPriceCandle `json:"output2"`
}

// DailyPriceSummary 는 응답의 output1 (단일 객체).
type DailyPriceSummary struct {
	Rsym string `json:"rsym"` // 실시간조회종목코드
	Zdiv string `json:"zdiv"` // 소수점자리수
	Nrec string `json:"nrec"` // 전일종가 (KIS 가 string 으로 줌 — 형식 다양)
}

// DailyPriceCandle 은 응답의 output2 한 행 (한 일자).
type DailyPriceCandle struct {
	Xymd string          `json:"xymd"`         // 일자(YYYYMMDD)
	Clos decimal.Decimal `json:"clos"`         // 종가
	Sign string          `json:"sign"`         // 대비기호 (1상한/2상승/3보합/4하한/5하락)
	Diff decimal.Decimal `json:"diff"`         // 대비
	Rate float64         `json:"rate,string"`  // 등락율
	Open decimal.Decimal `json:"open"`         // 시가
	High decimal.Decimal `json:"high"`         // 고가
	Low  decimal.Decimal `json:"low"`          // 저가
	Tvol int64           `json:"tvol,string"`  // 거래량
	Tamt int64           `json:"tamt,string"`  // 거래대금
	Pbid decimal.Decimal `json:"pbid"`         // 매수호가
	Vbid int64           `json:"vbid,string"`  // 매수호가잔량
	Pask decimal.Decimal `json:"pask"`         // 매도호가
	Vask int64           `json:"vask,string"`  // 매도호가잔량
}

// InquireDailyPriceParams 는 해외주식_기간별시세 조회 파라미터.
type InquireDailyPriceParams struct {
	Auth string // AUTH — 빈 값 default
	Excd string // EXCD — 거래소코드 (HKS/NYS/NAS/AMS/TSE/SHS/SZS/SHI/SZI/HSX/HNX)
	Symb string // SYMB — 종목코드
	Gubn string // GUBN — "0":일, "1":주, "2":월. 빈 값=>"0"
	Bymd string // BYMD — 조회기준일자 (YYYYMMDD). 빈 값=>오늘
	Modp string // MODP — "0":미반영, "1":반영(수정주가). 빈 값=>"0"
	Keyb string // KEYB — NEXT KEY BUFF (다음 조회 시 사용). 빈 값 default
}

// InquireDailyPrice 는 해외주식_기간별시세 호출.
//
// 한투 docs: docs/api/해외주식/해외주식_기간별시세.md
// path: /uapi/overseas-price/v1/quotations/dailyprice (HHDFS76240000)
//
// 한 번의 호출에 최대 100건. 미국은 0분지연, 홍콩/베트남/중국/일본은 15분지연시세.
func (c *Client) InquireDailyPrice(ctx context.Context, params InquireDailyPriceParams) (*DailyPrice, error) {
	gubn := params.Gubn
	if gubn == "" {
		gubn = "0"
	}
	modp := params.Modp
	if modp == "" {
		modp = "0"
	}

	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/overseas-price/v1/quotations/dailyprice",
		TrID:   "HHDFS76240000",
		Query: map[string]string{
			"AUTH": params.Auth,
			"EXCD": params.Excd,
			"SYMB": params.Symb,
			"GUBN": gubn,
			"BYMD": params.Bymd,
			"MODP": modp,
			"KEYB": params.Keyb,
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

// DailyChartPrice 는 해외주식 종목/지수/환율 기간별시세 (FHKST03030100) 응답.
//
// 한투 docs: docs/api/해외주식/해외주식_종목_지수_환율기간별시세(일_주_월_년).md
// path: /uapi/overseas-price/v1/quotations/inquire-daily-chartprice
//
// 미국 주식은 다우30/나스닥100/S&P500 만 조회 가능 (다른 미국 종목은 InquireDailyPrice).
type DailyChartPrice struct {
	Output1 DailyChartPriceSummary `json:"output1"`
	Output2 []DailyChartPriceCandle `json:"output2"`
}

// DailyChartPriceSummary 는 응답의 output1 (단일 객체, 기본정보).
type DailyChartPriceSummary struct {
	OvrsNmixPrdyVrss decimal.Decimal `json:"ovrs_nmix_prdy_vrss"`        // 전일 대비
	PrdyVrssSign     string          `json:"prdy_vrss_sign"`             // 전일 대비 부호
	PrdyCtrt         float64         `json:"prdy_ctrt,string"`           // 전일 대비율
	OvrsNmixPrdyClpr decimal.Decimal `json:"ovrs_nmix_prdy_clpr"`        // 전일 종가
	AcmlVol          int64           `json:"acml_vol,string"`            // 누적 거래량
	HtsKorIsnm       string          `json:"hts_kor_isnm"`               // HTS 한글 종목명
	OvrsNmixPrpr     decimal.Decimal `json:"ovrs_nmix_prpr"`             // 현재가
	StckShrnIscd     string          `json:"stck_shrn_iscd"`             // 단축 종목코드
	PrdyVol          int64           `json:"prdy_vol,string"`            // 전일 거래량
	OvrsProdOprc     decimal.Decimal `json:"ovrs_prod_oprc"`             // 시가
	OvrsProdHgpr     decimal.Decimal `json:"ovrs_prod_hgpr"`             // 최고가
	OvrsProdLwpr     decimal.Decimal `json:"ovrs_prod_lwpr"`             // 최저가
}

// DailyChartPriceCandle 은 응답의 output2 한 행 (한 일자/주/월/년 봉).
type DailyChartPriceCandle struct {
	StckBsopDate  string          `json:"stck_bsop_date"`        // 영업 일자
	OvrsNmixPrpr  decimal.Decimal `json:"ovrs_nmix_prpr"`        // 현재가
	OvrsNmixOprc  decimal.Decimal `json:"ovrs_nmix_oprc"`        // 시가
	OvrsNmixHgpr  decimal.Decimal `json:"ovrs_nmix_hgpr"`        // 최고가
	OvrsNmixLwpr  decimal.Decimal `json:"ovrs_nmix_lwpr"`        // 최저가
	AcmlVol       int64           `json:"acml_vol,string"`       // 누적 거래량
	ModYn         string          `json:"mod_yn"`                // 변경 여부
}

// InquireDailyChartPriceParams 는 해외주식 종목/지수/환율 기간별시세 조회 파라미터.
type InquireDailyChartPriceParams struct {
	MarketCode string // FID_COND_MRKT_DIV_CODE — "N":해외지수, "X":환율, "I":국채, "S":금선물
	Symbol     string // FID_INPUT_ISCD — 종목코드 (해외주식 마스터 코드 참조)
	FromDate   string // FID_INPUT_DATE_1 — 시작일자(YYYYMMDD)
	ToDate     string // FID_INPUT_DATE_2 — 종료일자(YYYYMMDD)
	Period     string // FID_PERIOD_DIV_CODE — "D":일/"W":주/"M":월/"Y":년
}

// InquireDailyChartPrice 는 해외주식 종목/지수/환율 기간별시세 호출.
//
// 한투 docs: docs/api/해외주식/해외주식_종목_지수_환율기간별시세(일_주_월_년).md
// path: /uapi/overseas-price/v1/quotations/inquire-daily-chartprice (FHKST03030100)
//
// ※ 미국 주식 조회 시 다우30, 나스닥100, S&P500 종목만 가능. 다른 미국 종목은 InquireDailyPrice 사용.
func (c *Client) InquireDailyChartPrice(ctx context.Context, params InquireDailyChartPriceParams) (*DailyChartPrice, error) {
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/overseas-price/v1/quotations/inquire-daily-chartprice",
		TrID:   "FHKST03030100",
		Query: map[string]string{
			"FID_COND_MRKT_DIV_CODE": params.MarketCode,
			"FID_INPUT_ISCD":         params.Symbol,
			"FID_INPUT_DATE_1":       params.FromDate,
			"FID_INPUT_DATE_2":       params.ToDate,
			"FID_PERIOD_DIV_CODE":    params.Period,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}

	var res DailyChartPrice
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse DailyChartPrice: %w", err)
	}
	return &res, nil
}
```

- [ ] **Step 4: PASS**

- [ ] **Step 5: Commit**

```bash
git add overseas/chart.go overseas/chart_test.go
git commit -m "$(cat <<'EOF'
[feat] overseas — InquireDailyPrice + InquireDailyChartPrice (chart 2 메서드)

DailyPrice (HHDFS76240000) — 단일 종목 일/주/월별, 11 거래소, 최대 100건.
DailyChartPrice (FHKST03030100) — 종목/지수/환율 통합 일/주/월/년. 미국은
다우30/나스닥100/S&P500 한정.

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 9: overseas/ranking.go — InquireUpdownRate

**Files:**
- Create: `overseas/ranking.go`
- Create: `overseas/ranking_test.go`

- [ ] **Step 1: 테스트 작성**

```go
package overseas_test

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kenshin579/korea-investment-stock/overseas"
)

func TestClient_InquireUpdownRate(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/ranking/updown-rate`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "updown_rate_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireUpdownRate(context.Background(), overseas.InquireUpdownRateParams{
		Excd: "NAS",
		Gubn: "1", // 상승율
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, "NAS", capturedQuery.Get("EXCD"))
	assert.Equal(t, "1", capturedQuery.Get("GUBN"))
	assert.Equal(t, "0", capturedQuery.Get("NDAY")) // default 당일
	assert.Equal(t, "0", capturedQuery.Get("VOL_RANG")) // default 전체

	require.Len(t, res.Output2, 2)
	assert.Equal(t, "NVDA", res.Output2[0].Symb)
	assert.Equal(t, "엔비디아", res.Output2[0].Name)
	d, _ := decimal.NewFromString("920.45")
	assert.True(t, d.Equal(res.Output2[0].Last))
	assert.InDelta(t, 5.16, res.Output2[0].Rate, 0.001)
}
```

- [ ] **Step 2: FAIL**

- [ ] **Step 3: 구현 — `overseas/ranking.go`**

```go
package overseas

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/shopspring/decimal"

	"github.com/kenshin579/korea-investment-stock/internal/httpclient"
)

// UpdownRate 은 해외주식_상승율_하락율 (HHDFS76290000) 응답.
//
// 한투 docs: docs/api/해외주식/해외주식_상승율_하락율.md
// path: /uapi/overseas-stock/v1/ranking/updown-rate
type UpdownRate struct {
	Output1 UpdownRateSummary `json:"output1"`
	Output2 []UpdownRateItem  `json:"output2"`
}

// UpdownRateSummary 는 응답의 output1 (단일 객체).
type UpdownRateSummary struct {
	Zdiv string `json:"zdiv"`        // 소수점자리수
	Stat string `json:"stat"`        // 거래상태정보
	Crec int64  `json:"crec,string"` // 현재 Count
	Trec int64  `json:"trec,string"` // 전체조회종목수
	Nrec int64  `json:"nrec,string"` // RecordCount
}

// UpdownRateItem 은 응답의 output2 한 행.
type UpdownRateItem struct {
	Rsym   string          `json:"rsym"`         // 실시간조회심볼
	Excd   string          `json:"excd"`         // 거래소코드
	Symb   string          `json:"symb"`         // 종목코드
	Name   string          `json:"name"`         // 종목명 (한글)
	Last   decimal.Decimal `json:"last"`         // 현재가
	Sign   string          `json:"sign"`         // 기호
	Diff   decimal.Decimal `json:"diff"`         // 대비
	Rate   float64         `json:"rate,string"`  // 등락율
	Tvol   int64           `json:"tvol,string"`  // 거래량
	Pask   decimal.Decimal `json:"pask"`         // 매도호가
	Pbid   decimal.Decimal `json:"pbid"`         // 매수호가
	NBase  decimal.Decimal `json:"n_base"`       // 기준가격
	NDiff  decimal.Decimal `json:"n_diff"`       // 기준가격대비
	NRate  float64         `json:"n_rate,string"` // 기준가격대비율
	Rank   int64           `json:"rank,string"`  // 순위
	Ename  string          `json:"ename"`        // 영문종목명
	EOrdyn string          `json:"e_ordyn"`      // 매매가능
}

// InquireUpdownRateParams 는 해외주식_상승율_하락율 조회 파라미터.
type InquireUpdownRateParams struct {
	Keyb    string // KEYB — NEXT KEY BUFF. 빈 값=>"" (공백)
	Auth    string // AUTH — 사용자권한정보. 빈 값=>"" (공백)
	Excd    string // EXCD — 거래소코드 (NYS/NAS/AMS/HKS/SHS/SZS/HSX/HNX/TSE)
	Gubn    string // GUBN — "0":하락율, "1":상승율
	Nday    string // NDAY — N일전: "0"(당일)/"1"(2일)/"2"(3일)/"3"(5일)/"4"(10일)/"5"(20일)/"6"(30일)/"7"(60일)/"8"(120일)/"9"(1년). 빈 값=>"0"
	VolRang string // VOL_RANG — 거래량조건: "0"(전체)/"1"(1백주이상)/"2"(1천주이상)/"3"(1만주이상)/"4"(10만주이상)/"5"(100만주이상)/"6"(1000만주이상). 빈 값=>"0"
}

// InquireUpdownRate 는 해외주식_상승율_하락율 호출.
//
// 한투 docs: docs/api/해외주식/해외주식_상승율_하락율.md
// path: /uapi/overseas-stock/v1/ranking/updown-rate (HHDFS76290000)
func (c *Client) InquireUpdownRate(ctx context.Context, params InquireUpdownRateParams) (*UpdownRate, error) {
	nday := params.Nday
	if nday == "" {
		nday = "0"
	}
	vol := params.VolRang
	if vol == "" {
		vol = "0"
	}

	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/overseas-stock/v1/ranking/updown-rate",
		TrID:   "HHDFS76290000",
		Query: map[string]string{
			"KEYB":     params.Keyb,
			"AUTH":     params.Auth,
			"EXCD":     params.Excd,
			"GUBN":     params.Gubn,
			"NDAY":     nday,
			"VOL_RANG": vol,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}

	var res UpdownRate
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse UpdownRate: %w", err)
	}
	return &res, nil
}
```

- [ ] **Step 4: PASS**

- [ ] **Step 5: Commit**

```bash
git add overseas/ranking.go overseas/ranking_test.go
git commit -m "$(cat <<'EOF'
[feat] overseas — InquireUpdownRate (해외주식 상승율_하락율, HHDFS76290000)

UpdownRate (Output1 5 + Output2 17 필드) + InquireUpdownRateParams (6 query).
Gubn (0하락/1상승) + Nday (당일~1년 전) + VolRang (거래량 조건). 9 거래소
(NYS/NAS/AMS/HKS/SHS/SZS/HSX/HNX/TSE).

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 10: examples/overseas_*

**Files:**
- Create: `examples/overseas_price/main.go`
- Create: `examples/overseas_chart/main.go`
- Create: `examples/overseas_symbols/main.go`

- [ ] **Step 1: examples/overseas_price/main.go**

```go
// overseas_price example: InquirePriceDetail + SearchInfo for AAPL.
//
// Run: KIS credentials env vars 후 go run ./examples/overseas_price
package main

import (
	"context"
	"fmt"
	"log"

	kis "github.com/kenshin579/korea-investment-stock"
	"github.com/kenshin579/korea-investment-stock/overseas"
)

func main() {
	client, err := kis.NewClientFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()

	price, err := client.Overseas.InquirePriceDetail(ctx, overseas.InquirePriceDetailParams{
		Excd: "NAS",
		Symb: "AAPL",
	})
	if err != nil {
		log.Fatalf("InquirePriceDetail: %v", err)
	}
	fmt.Printf("[AAPL @ NAS] 현재가 %s %s, PER=%v, PBR=%v\n",
		price.Output.Last, price.Output.Curr, price.Output.Perx, price.Output.Pbrx)
	fmt.Printf("  52주 최고/최저: %s / %s\n", price.Output.H52p, price.Output.L52p)
	fmt.Printf("  거래량: %d, 시가총액: %d\n", price.Output.Tvol, price.Output.Tomv)

	info, err := client.Overseas.SearchInfo(ctx, overseas.SearchInfoParams{
		PrdtTypeCD: "512", // NASDAQ
		Pdno:       "AAPL",
	})
	if err != nil {
		log.Fatalf("SearchInfo: %v", err)
	}
	fmt.Printf("\n상품정보: %s (%s)\n", info.Output.PrdtEngName, info.Output.PrdtName)
	fmt.Printf("  거래소: %s (%s), 통화: %s\n",
		info.Output.OvrsExcgName, info.Output.OvrsExcgCd, info.Output.TrCrcyCd)
	fmt.Printf("  ISIN: %s, SEDOL: %s, Bloomberg: %s\n",
		info.Output.StdPdno, info.Output.SedolNo, info.Output.BlbgTckrText)
}
```

- [ ] **Step 2: examples/overseas_chart/main.go**

```go
// overseas_chart example: InquireDailyPrice + InquireDailyChartPrice.
//
// Run: KIS credentials env vars 후 go run ./examples/overseas_chart
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	kis "github.com/kenshin579/korea-investment-stock"
	"github.com/kenshin579/korea-investment-stock/overseas"
)

func main() {
	client, err := kis.NewClientFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	today := time.Now().Format("20060102")

	// AAPL 일봉
	daily, err := client.Overseas.InquireDailyPrice(ctx, overseas.InquireDailyPriceParams{
		Excd: "NAS",
		Symb: "AAPL",
		Bymd: today,
	})
	if err != nil {
		log.Fatalf("InquireDailyPrice: %v", err)
	}
	fmt.Printf("[AAPL] 일봉 %d 캔들\n", len(daily.Output2))
	for i, c := range daily.Output2 {
		if i >= 5 {
			fmt.Println("  ... (이하 생략)")
			break
		}
		fmt.Printf("  %s: O=%s H=%s L=%s C=%s V=%d\n",
			c.Xymd, c.Open, c.High, c.Low, c.Clos, c.Tvol)
	}

	// S&P500 지수 (chart-price endpoint)
	from := time.Now().AddDate(0, 0, -30).Format("20060102")
	idx, err := client.Overseas.InquireDailyChartPrice(ctx, overseas.InquireDailyChartPriceParams{
		MarketCode: "N", // 해외지수
		Symbol:     "SPX",
		FromDate:   from,
		ToDate:     today,
		Period:     "D",
	})
	if err != nil {
		log.Fatalf("InquireDailyChartPrice: %v", err)
	}
	fmt.Printf("\n[%s %s] 30일 일봉 %d 개\n",
		idx.Output1.StckShrnIscd, idx.Output1.HtsKorIsnm, len(idx.Output2))
	for i, c := range idx.Output2 {
		if i >= 5 {
			fmt.Println("  ... (이하 생략)")
			break
		}
		fmt.Printf("  %s: O=%s H=%s L=%s C=%s V=%d\n",
			c.StckBsopDate, c.OvrsNmixOprc, c.OvrsNmixHgpr, c.OvrsNmixLwpr, c.OvrsNmixPrpr, c.AcmlVol)
	}
}
```

- [ ] **Step 3: examples/overseas_symbols/main.go**

```go
// overseas_symbols example: FetchOverseasSymbols("nas") — NASDAQ 마스터 다운로드.
//
// Run: KIS credentials env vars 후 go run ./examples/overseas_symbols
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
		log.Fatal(err)
	}
	ctx := context.Background()

	syms, err := client.Overseas.FetchOverseasSymbols(ctx, "nas")
	if err != nil {
		log.Fatalf("FetchOverseasSymbols: %v", err)
	}
	fmt.Printf("NASDAQ 종목 %d 개\n", len(syms))
	for i, s := range syms {
		if i >= 10 {
			fmt.Println("  ... (이하 생략)")
			break
		}
		fmt.Printf("  %s: %s\n", s.Symbol, s.EnglishName)
	}
}
```

- [ ] **Step 4: 컴파일 검증**

```bash
go build ./examples/overseas_price && \
go build ./examples/overseas_chart && \
go build ./examples/overseas_symbols && \
echo OK
```

- [ ] **Step 5: Commit**

```bash
git add examples/overseas_price examples/overseas_chart examples/overseas_symbols
git commit -m "$(cat <<'EOF'
[feat] examples/overseas_* — 3 예시 (price + chart + symbols)

- overseas_price: AAPL InquirePriceDetail + SearchInfo
- overseas_chart: AAPL daily + S&P500 30일 chart
- overseas_symbols: NASDAQ 종목 마스터 다운로드 (앞 10개)

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 11: 문서 갱신 (CLAUDE.md, README.md, CHANGELOG.md, overseas/doc.go)

**Files:**
- Modify: `CLAUDE.md`
- Modify: `README.md`
- Modify: `CHANGELOG.md`
- Modify: `overseas/doc.go`

- [ ] **Step 1: README.md — Available Methods 표 갱신**

Find existing `## Available Methods (Phase 1.2 + 1.3 + 1.4)` heading and 22-method table. REPLACE with `## Available Methods (Phase 1.2 + 1.3 + 1.4 + 1.5)` heading and 28-method table (existing 22 + 6 new):

```markdown
## Available Methods (Phase 1.2 + 1.3 + 1.4 + 1.5)

| Method | 한투 path | TR_ID |
|--------|----------|-------|
... (existing 22 rows) ...
| `Overseas.InquirePriceDetail` | `overseas-price/v1/quotations/price-detail` | HHDFS76200200 |
| `Overseas.SearchInfo` | `overseas-price/v1/quotations/search-info` | CTPF1702R |
| `Overseas.InquireDailyPrice` | `overseas-price/v1/quotations/dailyprice` | HHDFS76240000 |
| `Overseas.InquireDailyChartPrice` | `overseas-price/v1/quotations/inquire-daily-chartprice` | FHKST03030100 |
| `Overseas.InquireUpdownRate` | `overseas-stock/v1/ranking/updown-rate` | HHDFS76290000 |
| `Overseas.FetchOverseasSymbols(market)` | (KIS 공개 마스터 — 11 거래소) | — |
```

- [ ] **Step 2: CLAUDE.md — banner 갱신**

Replace `> **Phase 1.4 — domestic 투자자/업종/IPO (v1.2.0)...**` with:
```
> **Phase 1.5 — 해외주식 (v1.3.0, Python parity 완성).** Phase 2+ 는 추후 결정.
```

ADD spec link bullet:
```markdown
- Phase 1.5 implementation plan: [`docs/superpowers/specs/2026-05-05-phase1-5-overseas-implementation-plan.md`](docs/superpowers/specs/2026-05-05-phase1-5-overseas-implementation-plan.md)
```

- [ ] **Step 3: CHANGELOG.md — `[1.3.0]` entry**

ADD section AT THE TOP (above `[1.2.0]`):

```markdown
## [1.3.0] - 2026-05-05

### Added — Phase 1.5 (해외주식, Python parity 완성)

- `Overseas.InquirePriceDetail` — 해외주식 현재가상세 (HHDFS76200200)
- `Overseas.SearchInfo` — 해외주식 상품기본정보 (CTPF1702R)
- `Overseas.InquireDailyPrice` — 해외주식 기간별시세 (HHDFS76240000) — 11 거래소 단일 종목 일/주/월
- `Overseas.InquireDailyChartPrice` — 해외주식 종목/지수/환율 기간별시세 (FHKST03030100) — 일/주/월/년 (미국 주식은 다우30/나스닥100/S&P500 한정)
- `Overseas.InquireUpdownRate` — 해외주식 상승율/하락율 (HHDFS76290000)
- `Overseas.FetchOverseasSymbols(market)` — 11 거래소 통합 (KIS 공개 마스터)
- `internal/overseasmaster` 패키지 — 해외 마스터 파일 파싱
- examples: `overseas_price`, `overseas_chart`, `overseas_symbols`

### Changed

- `overseas.New(http, master)` 시그니처 — `*mastercache.Cache` 파라미터 추가 (internal API; BC-safe)

### Notes

- NASDAQ/NYSE/AMEX 별 메서드는 `FetchOverseasSymbols(market)` 로 통합 (Python wrapper convenience 미반영 정책 일관)
- `Overseas.SearchInfo` 의 응답 struct 명은 `OverseasProductInfo` (domestic 의 `ProductInfo` 와 다른 패키지지만 명시적으로 구분)
- 차트 endpoint 두 개 보완: `dailyprice` 는 단일 종목 (모든 미국 종목 지원), `inquire-daily-chartprice` 는 지수/환율 통합 (미국은 다우30/나스닥100/S&P500 한정)

### Phase 1 완성

이번 release 로 Python 라이브러리의 28 fetch 메서드 도메인 커버리지 완성:
- Phase 1.2: 7 메서드 (국내 시세/심볼/차트)
- Phase 1.3: 9 메서드 (국내 순위/재무)
- Phase 1.4: 6 메서드 (국내 투자자/업종/IPO)
- Phase 1.5: 6 메서드 (해외주식)
- 총 28 메서드 (Python 의 fetch 28개 + IPO helpers 9개 omit 의 카테고리 커버리지)
```

- [ ] **Step 4: overseas/doc.go 갱신**

Replace existing content with:

```go
// Package overseas 는 한국투자증권 OpenAPI 의 해외주식 카테고리 메서드.
//
// Phase 1.5 메서드 (6):
//
//   - InquirePriceDetail        — 해외주식 현재가상세 (HHDFS76200200)
//   - SearchInfo                — 해외주식 상품기본정보 (CTPF1702R)
//   - InquireDailyPrice         — 해외주식 기간별시세 (HHDFS76240000) — 단일 종목 11 거래소
//   - InquireDailyChartPrice    — 해외주식 종목/지수/환율 기간별시세 (FHKST03030100) — 미국은 지수 한정
//   - InquireUpdownRate         — 해외주식 상승율/하락율 (HHDFS76290000)
//   - FetchOverseasSymbols      — 11 거래소 통합 마스터 (KIS 공개 다운로드)
//
// 사용자는 root kis.Client 의 Overseas 필드로 접근.
package overseas
```

- [ ] **Step 5: 검증**

```bash
go build ./... && go vet ./... && gofmt -l .
```
Expected: silent.

- [ ] **Step 6: Commit**

```bash
git add CLAUDE.md README.md CHANGELOG.md overseas/doc.go
git commit -m "$(cat <<'EOF'
[doc] Phase 1.5 메서드 문서 갱신 — CLAUDE.md, README.md, CHANGELOG.md, overseas/doc.go

Phase 1.5 의 6 메서드 (overseas) 목록 + CHANGELOG [1.3.0] entry. Phase 1
Python parity 28 메서드 완성 표시. CLAUDE.md banner 갱신
(Phase 1.4 → 1.5, v1.2.0 → v1.3.0).

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 12: 최종 점검

- [ ] **Step 1: gofmt cleanup (필요 시)**

`gofmt -w overseas/*.go internal/overseasmaster/*.go` 후 `gofmt -l .` empty 확인.

If diff exists, commit `[chore] Phase 1.5 — gofmt cleanup`.

- [ ] **Step 2: 빌드/vet**

`go build ./... && go vet ./...` — silent.

- [ ] **Step 3: 모든 테스트 + race**

`go test ./... -race -count=1` — all PASS.

- [ ] **Step 4: Coverage**

```bash
go test ./... -coverprofile=/tmp/cov.out -covermode=atomic
go tool cover -func=/tmp/cov.out | tail -15
```

Expected: domestic/ ≥ 80%, root kis ≥ 80%, overseas/ ≥ 70% (해외 마스터 다운로드 코드 포함이라 약간 낮을 수 있음).

If short: API error 경로 (`rt_cd != "0"`) 테스트 추가.

- [ ] **Step 5: 디렉터리 구조 확인**

```bash
ls -la overseas/{price,search,chart,ranking,symbols,client,doc}.go overseas/{price,search,chart,ranking,symbols}_test.go overseas/testdata/{price_detail,search_info,daily_price,daily_chart_price,updown_rate}_success.json overseas/testdata/nas_code_sample.cod.zip examples/overseas_{price,chart,symbols}/main.go internal/overseasmaster/{doc,overseasmaster}.go internal/overseasmaster/overseasmaster_test.go internal/overseasmaster/testdata/nas_code_sample.cod.zip 2>&1 | wc -l
```
Expected: ~22 (5 testdata json + 1 cod.zip + 7 overseas .go + 5 _test.go + 3 examples + 3 overseasmaster files + 1 master sample).

- [ ] **Step 6: Commit history**

`git log main..HEAD --oneline | wc -l` — should be ~13-15 commits.

---

## Task 13: PR 생성 (사용자 승인 후)

> Claude 는 push / PR 생성을 사용자 명시적 승인 후에만 실행 (글로벌 정책).

- [ ] **Step 1: 사용자 승인 요청**

작업 진행 보고 + PR 생성 가능 여부 confirm.

- [ ] **Step 2: Push branch**

`git push -u origin docs/phase1-5-spec`

- [ ] **Step 3: PR 생성**

```bash
gh pr create --title "Phase 1.5 — 해외주식 (v1.3.0, Python parity 완성)" --reviewer kenshin579 --base main --head docs/phase1-5-spec --body "$(cat <<'EOF'
## Summary

- 해외주식 5 KIS REST 메서드 + `FetchOverseasSymbols(market)` 통합 = 6 메서드
- Python parity 완성 — 28 fetch 메서드 도메인 커버리지 완료
- `internal/overseasmaster` 패키지 신규 (해외 마스터 파일 파서)
- `overseas.New(http, master)` 시그니처 확장
- v1.3.0 release 대상

## 메서드 → 한투 API 매핑

| Go 메서드 | path | TR_ID |
|-----------|------|-------|
| InquirePriceDetail | overseas-price/v1/quotations/price-detail | HHDFS76200200 |
| SearchInfo | overseas-price/v1/quotations/search-info | CTPF1702R |
| InquireDailyPrice | overseas-price/v1/quotations/dailyprice | HHDFS76240000 |
| InquireDailyChartPrice | overseas-price/v1/quotations/inquire-daily-chartprice | FHKST03030100 |
| InquireUpdownRate | overseas-stock/v1/ranking/updown-rate | HHDFS76290000 |
| FetchOverseasSymbols(market) | (KIS 공개 마스터, 11 거래소) | — |

## Test Plan

- [x] go build/vet/fmt clean
- [x] go test ./... -race -count=1 모든 패키지 PASS
- [x] Coverage 임계 충족 (domestic/ overseas/ root)
- [x] httpmock 단위 테스트 (각 5 KIS API + 3 symbols)

## Breaking Changes

없음 — 신규 메서드 추가만. `overseas.New(http, master)` 시그니처 변경은 internal API.

🤖 Generated with [Claude Code](https://claude.com/claude-code)
EOF
)"
```

- [ ] **Step 4: Merge (사용자 승인 후)**

`gh pr merge <PR#> --merge`

- [ ] **Step 5: 후속 작업 (사용자 승인 후)**

```bash
git checkout main && git pull
git tag -a v1.3.0 -m "v1.3.0: Phase 1.5 — 해외주식 6 메서드, Python parity 완성"
git push origin v1.3.0
gh release create v1.3.0 --title "v1.3.0 — Phase 1.5: 해외주식 (Python parity 완성)" --notes-file <(awk '/^## \[/{c++} c==1' CHANGELOG.md)
```
