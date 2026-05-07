# Phase 6 — Financial Extension Design

**Status:** Active design (2026-05-08)
**Goal:** Phase 1.3 의 5 재무 메서드 외 누락된 read-only 재무/순위 endpoint 2 메서드 추가.
**Out of Scope:** WebSocket, 주문, 잔고, 선물옵션.

---

## §1. 목적

Phase 1~5 (115 methods) 후 재무 도메인 잔존 누락 audit 결과 3 개 docs (대차대조표 / 기타주요비율 / 재무비율_순위) 발견. **대차대조표는 Phase 1.3 의 `InquireBalanceSheet` 와 path/TR_ID 동일 → 중복으로 SKIP**. 실제 신규 = 2 메서드.

---

## §2. 메서드 매핑

| EP | Method (Style A) | Path | TR_ID | File | Output | Fields |
|---|---|---|---|---|---|---|
| 1 | `InquireOtherMajorRatios` | `/uapi/domestic-stock/v1/finance/other-major-ratios` | FHKST66430500 | `domestic/financial.go` | `output []` | 5 |
| 2 | `InquireFinanceRatioRanking` | `/uapi/domestic-stock/v1/ranking/finance-ratio` | FHPST01750000 | `domestic/ranking.go` | `output []` | 27 |

---

## §3. Anomalies (구현 시 주의)

1. **EP1 `fid_div_cls_code` lowercase** — Phase 1.3 의 `InquireGrowthRatio` (FHKST66430800) 와 동일. `inquireFinanceQuery` helper 사용 불가, inline query map.
2. **EP1 `payout_rate` 은 비정상 출력** — KIS docs 권고: "비정상 출력으로 무시". `string` 보존, 코멘트 명시.
3. **EP1 큰 금액 vs 비율**: `eva` (EVA), `ebitda` (EBITDA) → `decimal.Decimal` (금액). `ev_ebitda` → `float64,string` (배수).
4. **EP2 13 query 중 5 hardcoded**: `fid_trgt_cls_code="0"` / `fid_cond_scr_div_code="20175"` / `fid_div_cls_code="0"` / `fid_blng_cls_code="0"` / `fid_trgt_exls_cls_code="0"` — Params struct 에 노출하지 않음. `InquireMarketCap` 의 hardcoded screen code 패턴과 동일.
5. **EP2 페이지네이션 없음** — "최대 30 건, 다음 조회 불가". `tr_cont` 미사용.
6. **EP2 27 fields**: 가격(`stck_prpr`/`prdy_vrss`) → `decimal.Decimal`, 비율(`prdy_ctrt`/`bis`/`grs` 등 17개) → `float64,string`, 수량/순위(`acml_vol`/`data_rank`/`iqry_csnu`) → `int64,string`, 코드/이름 → `string`.
7. **EP1 dropped (대차대조표)**: 동일 path/TR_ID 가 Phase 1.3 `InquireBalanceSheet` 로 이미 출시됨. 재구현 X.

---

## §4. 인프라 변경 (없음)

신규 file 없음. 기존 `domestic/financial.go` 와 `domestic/ranking.go` 에 append.

---

## §5. 진입/종료 조건

- 진입: main HEAD = v1.15.0 (Phase 5 완료, 누적 115)
- 종료: PR merge, v1.16.0 tag, GitHub Release
- 누적: 115 → **117** 메서드

---

## §6. 진행 절차

Phase 4.3 / Phase 5 패턴: **plan 작성 skip + 직접 batch 구현**. 작은 phase (2 메서드) 라 main agent 가 직접 구현.

Tasks:
1. testdata 2 fixtures (성공 + 1 InvalidJSON 가능)
2. `domestic/financial.go` append `InquireOtherMajorRatios` (EP1, 5 fields, lowercase fid)
3. `domestic/ranking.go` append `InquireFinanceRatioRanking` (EP2, 27 fields, 5 hardcoded params)
4. 2 메서드 unit tests (httpmock + InvalidJSON)
5. examples 시연 (기존 `examples/domestic_financial/` 또는 `domestic_ranking/` 에 추가)
6. 문서 갱신 (CLAUDE.md / README.md / CHANGELOG.md / `domestic/doc.go`)
7. 최종 점검 (gofmt/vet/build/race/coverage ≥80%)
8. PR 생성 (사용자 승인 후)
