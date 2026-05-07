# Phase 5 — ETF/NAV/관심종목 Read-Only Design

**Status:** Active design (2026-05-07)
**Goal:** ETF/NAV/관심종목 도메인 9 REST 메서드. Phase 4 (시장정보/순위 27 methods) 완료 후 Group B 영역 cover.
**Out of Scope:** WebSocket (EP3 H0STNAV0 — 별도 phase), 잔고/주문조회 (계좌 인증).

---

## §1. 목적과 결정 요약

KIS API 의 read-only 영역에서 Phase 1+2+2.5+ + Phase 3 + Phase 4 (총 106 methods) 가 cover 한 외에도 ETF/NAV/관심종목 영역이 남아있음. moneyflow 의 펀드 추적/portfolio 다종목 시세 조회 등에 유용.

### 핵심 결정

| 항목 | 결정 |
|---|---|
| 기존 패키지 활용 | `domestic/` 패키지 재사용. 새 sub-package 도입 안 함. 도메인별 file 분리. |
| 파일 구성 | `domestic/etf.go` (5: EP1+EP2+EP4+EP5+EP6 ETF/NAV) + `domestic/watchlist.go` (4: EP7+EP8+EP9+EP10 관심종목) |
| Sub-phase 분해 | 단일 phase 진행 (9 methods). 분해 불필요. |
| Style A 명명 | 메서드명 충돌 회피 위해 prefix 변형 (예: `InquireEtfPrice` for `inquire-price` ETF endpoint) |
| 타입 매핑 | Phase 2 standard. 가격 decimal, 수량 int64,string, 비율 float64,string, 코드/날짜 string. |

---

## §2. 메서드 매핑

| EP | Method (Style A) | Path | TR_ID | Output | Fields |
|---|---|---|---|---|---|
| 1 | `InquireEtfPrice` | `/uapi/etfetn/v1/quotations/inquire-price` | FHPST02400000 | `output {}` | 54 |
| 2 | `InquireComponentStockPrice` | `/uapi/etfetn/v1/quotations/inquire-component-stock-price` | FHKST121600C0 | `output1 {} + output2 []` | 16+15 |
| 4 | `InquireNavComparisonTimeTrend` | `/uapi/etfetn/v1/quotations/nav-comparison-time-trend` | FHPST02440100 | `output []` | 13 |
| 5 | `InquireNavComparisonDailyTrend` | `/uapi/etfetn/v1/quotations/nav-comparison-daily-trend` | FHPST02440200 | `output []` | 13 |
| 6 | `InquireNavComparisonTrend` | `/uapi/etfetn/v1/quotations/nav-comparison-trend` | FHPST02440000 | `output1 {} + output2 {}` | 12+8 |
| 7 | `InquireIntstockMultprice` | `/uapi/domestic-stock/v1/quotations/intstock-multprice` | FHKST11300006 | `output {}` (single obj despite 30 stocks!) | 29 |
| 8 | `InquireIntstockStocklistByGroup` | `/uapi/domestic-stock/v1/quotations/intstock-stocklist-by-group` | HHKCM113004C6 | `output1 {} + output2 []` | 2+10 |
| 9 | `InquireIntstockGrouplist` | `/uapi/domestic-stock/v1/quotations/intstock-grouplist` | HHKCM113004C7 | `output2 {}` (no output1!) | 6 |
| 10 | `InquireTopInterestStock` | `/uapi/domestic-stock/v1/ranking/top-interest-stock` | FHPST01800000 | `output []` | 13 |

**EP3 (`H0STNAV0`)는 WebSocket → 제외**. Phase 5 = 9 REST methods.

---

## §3. Anomalies (구현 시 주의)

1. **EP1 name collision**: 기존 `domestic.InquirePrice` (Phase 1.2) 와 path `inquire-price` 충돌 → `InquireEtfPrice` 로 rename.
2. **ETF base path**: EP1, EP2, EP4, EP5, EP6 은 `/uapi/etfetn/v1/quotations/` (NOT `/uapi/domestic-stock/v1/`).
3. **Mixed FID_ casing**: EP1/EP4/EP5/EP10 = lowercase fid_*; EP2/EP6 = UPPERCASE FID_; EP7/EP8/EP9 = non-FID UPPERCASE (TYPE/USER_ID/INTER_GRP_CODE 등).
4. **EP2 output1 docs label scramble**: KIS docs 의 Korean labels 가 다른 endpoint 에서 잘못 복사됨. Field 명만 source of truth.
5. **EP7 single object despite batch**: 30 (market_code, stock_code) pairs 를 받지만 응답은 single `output {}` (not array) — KIS docs 명시. 실 API 동작 검증 필요.
6. **EP8/EP9 USER_ID 필요**: 사용자의 HTS_ID 입력 필요 (계좌 인증 X, but user-identity required). Params struct 코멘트로 명시.
7. **EP9 output key is `output2` (no `output1`)**: 응답 key 가 `output2` 만 — `output1` 없음.
8. **EP10 hardcoded params**: `fid_cond_scr_div_code="20180"`, `fid_input_iscd_2="000000"` hardcode.
9. **EP2 hardcoded `FID_COND_SCR_DIV_CODE="11216"`**.

---

## §4. 인프라 변경 (없음)

기존 `domestic/` 패키지 재사용. 새 internal package 불필요.

신규 file:
- `domestic/etf.go` — ETF/NAV 5 methods (EP1, EP2, EP4, EP5, EP6)
- `domestic/watchlist.go` — 관심종목 4 methods (EP7, EP8, EP9, EP10)

---

## §5. 진입/종료 조건

- 진입: main HEAD = v1.14.0 (Phase 4 완료, 누적 106)
- 종료: PR merge, v1.15.0 tag, GitHub Release
- 누적: 106 → 115 메서드

---

## §6. 진행 절차

WebSocket EP3 1개 제외하면 9 REST. Phase 4.3 처럼 **plan 작성 skip + 직접 batch dispatch** 진행 (docs analyzer 결과를 spec 으로 사용).

Tasks:
1. testdata 9 fixtures
2. `domestic/etf.go` base + InquireEtfPrice (EP1, 54 fields)
3. InquireComponentStockPrice (EP2)
4. InquireNavComparisonTimeTrend (EP4)
5. InquireNavComparisonDailyTrend (EP5)
6. InquireNavComparisonTrend (EP6)
7. `domestic/watchlist.go` base + InquireIntstockMultprice (EP7)
8. InquireIntstockStocklistByGroup (EP8, USER_ID)
9. InquireIntstockGrouplist (EP9, USER_ID)
10. InquireTopInterestStock (EP10)
11. examples/domestic_etf_watchlist/main.go
12. 문서 갱신 (CLAUDE/README/CHANGELOG/domestic/doc.go)
13. 최종 점검 (gofmt/build/vet/race/coverage ≥80%)
14. PR 생성 (사용자 승인 후)
