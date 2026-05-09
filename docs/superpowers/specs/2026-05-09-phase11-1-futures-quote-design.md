# Phase 11.1 — 국내선물옵션 시세/조회 Design (Lightweight)

**Status:** Active design (2026-05-09)
**Goal:** 신규 도메인 `futures/` 도입. 국내선물옵션 시세/조회 9 EP REST 구현. 누적 121 → **130 REST + 17 WS = 147 endpoints**.

> **Scope 정정 (2026-05-09)**: 원래 11 EP 였으나 docs analyzer (Plan Task 1) 결과 EP4 `InquireCcnlBstime` (CTFO5139R) + EP7 `InquireDailyAmountFee` (CTFO6119R) 가 query 에 CANO/ACNT_PRDT_CD 를 요구해서 본 phase 의 시세/조회 scope 와 충돌. 두 EP 는 **Phase 11.4 Trading** 으로 미룸. 사용자 confirm 받음.
**Out of Scope:** 국내선물옵션 실시간 (Phase 11.2~11.3, 20 EP), Trading (Phase 11.4, HIGH RISK 12 EP), 해외선물옵션 (Phase 11.5~).

> **Lightweight spec + 형식 plan**: 신규 sub-package `futures/` 스캐폴딩이 Phase 3.1 (bonds/) 패턴 그대로 재사용. bonds 도 plan 작성했던 전례라 본 phase 도 plan 작성 (writing-plans skill 호출).

---

## §1. 사용자 합의 결정

| 항목 | 결정 | 근거 |
|---|---|---|
| 다음 phase driver | 한투 API coverage 완성도 | 외부 사용자 대상 라이브러리 |
| Coverage focus | 시세/조회 영역 (선물옵션/ELW/지수) | 미커버 분량 가장 큼 |
| 1순위 도메인 | 선물옵션 (큰 도메인 본격 시작) | ELW 보다 분량 크지만 우선순위 높음 |
| 시작 sub-phase | 국내 시세 (Phase 11.1) | 해외 (11.5+) 와 실시간 (11.2~11.3) 은 후속 |
| Sub-package 명 | `futures/` | bonds/ 패턴 일관, 짧고 단순. 옵션 포함 (yfinance/alpaca 등 관례) |
| Phase 11.1 scope | 9 EP (계좌정보 필요한 EP4/EP7 제외) | docs analyzer 결과 후 정정 |
| 진행 절차 | lightweight spec + 형식 plan 작성 | bonds/ 신규 sub-package 도입은 Phase 3.1 처럼 plan 가치 있음 |

---

## §2. EP Matrix (11)

| # | docs (한글 파일명) | 임시 메서드명 | 비고 |
|---|---|---|---|
| 1 | 선물옵션_시세 | `Futures.InquirePrice` | 현재가 |
| 2 | 선물옵션_시세호가 | `Futures.InquireAskingPrice` | 시세 + 호가 |
| 3 | 선물옵션_분봉조회 | `Futures.InquireTimeItemchartprice` | 분봉 차트 |
| ~~4~~ | ~~선물옵션_기준일체결내역~~ | ~~`Futures.InquireCcnlBstime`~~ | **Phase 11.4 미룸 (CANO/ACNT_PRDT_CD 필요)** |
| 5 | 선물옵션_일중예상체결추이 | `Futures.InquireExpectFluctTrend` | 일중 예상체결 |
| 6 | 선물옵션기간별시세(일/주/월/년) | `Futures.InquireDailyPrice` | 일/주/월/년 차트 |
| ~~7~~ | ~~선물옵션기간약정수수료일별~~ | ~~`Futures.InquireDailyAmountFee`~~ | **Phase 11.4 미룸 (CANO/ACNT_PRDT_CD 필요)** |
| 8 | 국내선물_기초자산_시세 | `Futures.InquireUnderlyingPrice` | 기초자산 가격 |
| 9 | 국내옵션전광판_선물 | `Futures.OptionBoardFuture` | 옵션 전광판 — 선물 |
| 10 | 국내옵션전광판_옵션월물리스트 | `Futures.OptionMonthlyList` | 월물 리스트 |
| 11 | 국내옵션전광판_콜풋 | `Futures.OptionBoardCallPut` | 콜/풋 전광판 |

**메서드명 확정 절차**: docs analyzer 단계에서 각 EP 의 한투 API path 마지막 segment 를 PascalCase 로 1:1 매핑 (Style A — Phase 1 부터 적용된 규칙). 위는 제안이며 path 확인 후 정정 가능.

**TR_ID/path/output schema**: docs analyzer 단계에서 각 docs 직접 분석 후 확정. 본 spec 은 high-level scope 만 명시.

---

## §3. Sub-package 구조

bonds/ 패턴 (Phase 3.1) 재사용:

```
futures/
  client.go            # Client struct + New(http) + wireInfra (master cache 미주입 — futures 마스터파일은 후속 phase)
  doc.go               # 패키지 doc
  testhelper_test.go   # bonds 패턴 그대로 — httpclient.NewForTest 사용
  quote.go             # InquirePrice / InquireAskingPrice / InquireDailyPrice / InquireUnderlyingPrice
  quote_test.go
  chart.go             # InquireTimeItemchartprice
  chart_test.go
  conclusion.go        # InquireConclusion / InquireExpectFluctTrend / InquireDailyCcldFee
  conclusion_test.go
  board.go             # OptionBoardFuture / OptionMonthlyList / OptionBoardCallPut
  board_test.go
```

root `client.go` 에 `Futures *futures.Client` 필드 추가 (`Bonds *bonds.Client` 옆에).

파일 분할은 의미적 그룹화 — quote (시세), chart (분봉/일별), conclusion (체결/예상), board (전광판). 메서드 수에 따라 implementation 단계에서 통합/분할 조정 가능.

---

## §4. 핵심 deviation

1. **종목코드 형식**: KRX 6자리 (`005930`) 와 다른 9자리 alphanumeric 선물옵션 코드 (예: `101W3000` 선물 / `201X3300` 옵션). 메서드 시그니처는 `string` 그대로 받고 caller 가 정확한 코드 입력. 마스터파일 cache 는 후속 phase 에서.
2. **MarketCode default**: 선물옵션 endpoint 의 `FID_COND_MRKT_DIV_CODE` default 가 `"F"` (선물) 또는 `"JF"` (선물옵션 통합) 또는 `"O"` (옵션) 가능성 — docs analyzer 단계에서 EP 별 정확히 확정. KRX `"J"`, 채권 `"B"`, 해외 `"X"`/`"N"` 와 다름.
3. **Multi-output 가능성**: 옵션전광판 (#9~11) 같은 endpoint 가 output1+output2 dual 응답일 수 있음 (콜/풋 분리 등). docs analyzer 시 정확히 확인.
4. **수수료 EP (#7)**: 거래수수료 정보. 시세 영역으로 분류했지만 trading-related — 본 phase 에 포함 (사용자 confirm). 응답 schema 가 일반 시세와 다를 수 있음.
5. **type 매핑 정책 (Phase 1 부터 일관)**:
   - 가격/체결가/액면가 → `decimal.Decimal`
   - 거래량/거래대금/수량 → `int64,string`
   - 비율/등락율/체결강도 → `float64,string`
   - 코드/Y-N/날짜/시간 → `string`

---

## §5. Testing

- TDD + httpmock + KIS docs sample response (Phase 1 부터 일관)
- 각 메서드: happy path + 1 invalid-JSON test (Phase 3.1 lesson — 신규 sub-package 도입 시 ~67-71% coverage 가 80% 충족 못 해 보강 필요)
- Coverage 목표:
  - 전체 ≥ 80%
  - `futures/` package ≥ 80%
- testdata fixture 작성 시 envelope (`rt_cd`/`msg_cd`/`msg1`) 포함 강제 (Phase 2.4 lesson)

---

## §6. 진행 절차

Phase 3.1 (bonds/) 패턴 — 신규 sub-package 도입은 형식 plan 작성:

1. ✅ design doc commit (이 단계)
2. **writing-plans skill 호출 → implementation plan 작성** (각 task verbatim code 포함)
3. docs analyzer 단계: 11 EP 각각의 한투 docs 직접 분석 → path / TR_ID / 정확한 메서드명 / Query params / Output schema 확정 (plan 의 task 0 또는 분리)
4. testdata fixtures 작성 (11 fixtures)
5. 신규 `futures/` package 스캐폴딩 — `client.go`, `doc.go`, `testhelper_test.go` (bonds 패턴 복사)
6. 메서드 구현 + 테스트 (4 그룹: quote / chart / conclusion / board) — plan 의 task 단위로 진행
7. example `examples/futures_basic/` (4-5 메서드 통합 시연)
8. 문서 갱신 (CLAUDE.md / README.md / CHANGELOG.md / `futures/doc.go`)
9. 최종 점검 (gofmt / vet / build / race / coverage ≥ 80%)
10. PR + merge + tag v1.21.0 + GitHub Release

---

## §7. 진입/종료 조건

- **진입**: main HEAD = v1.20.0 (Phase 10 완료, 121 REST + 17 WS = 138 endpoints)
- **종료**: PR merge to main, v1.21.0 tag, GitHub Release
- **누적**: 121 + 11 REST + 17 WS = **149 endpoints**

## §8. 다음 sub-phase (Phase 11.x roadmap)

- Phase 11.2: KRX 야간/주식/지수 선물옵션 실시간 (~13 EP)
- Phase 11.3: 상품선물 + 지수선물옵션 + 선물옵션 실시간체결통보 잔여 (~7 EP)
- Phase 11.4: 국내선물옵션 Trading 12 EP — HIGH RISK, 별도 design + safety review
- Phase 11.5+: 해외선물옵션 35 docs (시세/실시간/trading 분할)

본 spec 은 11.1 만 cover. 후속 phase 는 별도 design spec 작성.
