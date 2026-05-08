# Phase 7 — Helpers (Group D) Design

**Status:** Active design (2026-05-08)
**Goal:** Phase 1~6 (117 methods) 후 보조/헬퍼 성격 endpoint 5 후보 audit → 4 메서드 추가 (1 SKIP).
**Out of Scope:** WebSocket, 주문, 잔고, 선물옵션.

---

## §1. 목적

읽기 전용 보조 API (영업일/금리/순위/신용 정보) 4 메서드 마이그레이션. 5 후보 중 1 (상품기본조회) 은 Phase 1.1 의 `SearchInfo` 와 path/TR_ID 동일 → 중복으로 SKIP.

---

## §2. 메서드 매핑

| EP | Method (Style A) | Path | TR_ID | File | Output | Fields |
|---|---|---|---|---|---|---|
| 1 | `InquireMarketTime` | `/uapi/domestic-stock/v1/quotations/market-time` | HHMCM000002C0 | `domestic/market.go` (신규) | `output1 []` | 9 |
| 2 | `InquireCompInterest` | `/uapi/domestic-stock/v1/quotations/comp-interest` | FHPST07020000 | `domestic/interest.go` (신규) | `output1` single + `output2 []` | 7+7 |
| 3 | `InquireTradedByCompany` | `/uapi/domestic-stock/v1/ranking/traded-by-company` | FHPST01860000 | `domestic/ranking.go` (append) | `output []` | 12 |
| 4 | `InquireCreditByCompany` | `/uapi/domestic-stock/v1/quotations/credit-by-company` | FHPST04770000 | `domestic/credit.go` (신규) | `output []` | 3 |

**제외**: 상품기본조회 (CTPF1604R) — Phase 1.1 `domestic/info.go:47 SearchInfo` 와 동일 (path/TR/Query/Output 모두 일치).

---

## §3. Anomalies (구현 시 주의)

1. **EP1 파라미터 없음** — `Query: nil` 또는 빈 map. Path 만으로 호출. `InquireMarketTimeParams` struct 도 빈 struct 또는 omit.
2. **EP1 output1 array** — 단일 객체가 아닌 array 형태. `Output1 []MarketTimeItem`. `time`/`s_time`/`e_time` 은 6자리 HHmmss string. `date1~date5`/`today` 는 8자리 YYYYMMDD string.
3. **EP2 4 query 모두 UPPERCASE + 거의 hardcoded**: `FID_COND_MRKT_DIV_CODE="I"`, `FID_COND_SCR_DIV_CODE="20702"`, `FID_DIV_CLS_CODE="1"`(해외금리지표), `FID_DIV_CLS_CODE1=""`(공백=전체). 사용자 입력 없음 → `Params` struct 비움 (`InquireCompInterestParams{}`). `inquireFinanceQuery` 같은 helper 사용 안 함.
4. **EP2 dual output**: `output1` (single object, 7 fields — 종합 metadata) + `output2` (array, 7 fields — 개별 금리 항목). 두 output 모두 같은 7 필드 거의 일치하지만 `output1.prdy_ctrt` ↔ `output2.bstp_nmix_prdy_ctrt` 차이 있음.
5. **EP3 12 query 중 4 hardcoded**: `fid_trgt_exls_cls_code="0"`, `fid_cond_scr_div_code="20186"`, `fid_trgt_cls_code="0"`, `fid_aply_rang_vol="0"`. 사용자 입력 8개 → `MarketCode`/`DivCode`/`SortCode`/`InputDate1`/`InputDate2`/`InputISCD`/`PriceFrom`/`PriceTo`. lowercase fid (Phase 1.3 ranking 패턴).
6. **EP3 default values**: `MarketCode` 비면 `"J"`(KRX), `DivCode` 비면 `"0"`(전체), `SortCode` 비면 `"0"`(매도상위), `InputISCD` 비면 `"0000"`(전체). 가격 범위 빈 값은 그대로 전달.
7. **EP4 5 query 중 2 hardcoded**: `fid_cond_scr_div_code="20477"`, `fid_cond_mrkt_div_code="J"`. 사용자 입력 3개 → `SortCode`/`SelectYN`/`InputISCD`. lowercase fid.
8. **EP4 default values**: `SortCode` 비면 `"0"`(코드순), `SelectYN` 비면 `"0"`(신용주문가능), `InputISCD` 비면 `"0000"`(전체).
9. **타입 매핑 (모든 EP 공통)**:
   - 가격 (`stck_prpr`/`prdy_vrss`) → `decimal.Decimal`
   - 누적 거래량/거래대금/순위/체결합계 → `int64,string`
   - 비율 (`prdy_ctrt`/`crdt_rate`) → `float64,string`
   - 채권금리 (`bond_mnrt_prpr`/`bond_mnrt_prdy_vrss`) → `decimal.Decimal` (소수 있을 수 있음)
   - 기타 (이름/코드/날짜/시간) → plain `string`
10. **모든 EP 모의투자 미지원** — 실전 only. 별도 처리 없음 (httpclient envelope 가 alarm 처리).

---

## §4. 인프라 변경

**신규 파일 3개**:
- `domestic/market.go` (EP1)
- `domestic/interest.go` (EP2)
- `domestic/credit.go` (EP4)

**기존 파일 1개 append**:
- `domestic/ranking.go` (EP3)

**testdata fixtures 4개 신규**:
- `testdata/market_time_success.json`
- `testdata/comp_interest_success.json`
- `testdata/traded_by_company_success.json`
- `testdata/credit_by_company_success.json`

---

## §5. 진입/종료 조건

- 진입: main HEAD = v1.16.0 (Phase 6 완료, 누적 117)
- 종료: PR merge, v1.17.0 tag, GitHub Release
- 누적: 117 → **121** 메서드

---

## §6. 진행 절차

Phase 6 패턴 동일: **plan 작성 skip + 직접 batch 구현**. 4 메서드라 main agent 가 직접 구현.

Tasks:
1. testdata 4 fixtures (성공 응답)
2. `domestic/market.go` 신규 + `InquireMarketTime` (EP1, 파라미터 없음)
3. `domestic/interest.go` 신규 + `InquireCompInterest` (EP2, dual output, 4 hardcoded)
4. `domestic/ranking.go` append `InquireTradedByCompany` (EP3, 12 fields, 4 hardcoded)
5. `domestic/credit.go` 신규 + `InquireCreditByCompany` (EP4, 3 fields, 2 hardcoded)
6. 4 메서드 unit tests (httpmock + InvalidJSON 패턴 envelope-valid + output-shape-invalid)
7. examples 시연 (`domestic_market` / `domestic_ranking` append / `domestic_credit` / `domestic_interest`)
8. 문서 갱신 (CLAUDE.md / README.md / CHANGELOG.md / `domestic/doc.go`)
9. 최종 점검 (gofmt/vet/build/race/coverage ≥80%)
10. PR 생성 (사용자 승인 후)
