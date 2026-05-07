# Phase 4 — 국내주식 시장정보/순위 Read-Only Design

**Status:** Active design (2026-05-07)
**Goal:** 국내주식의 추가 read-only 시장정보/순위/시장운영 30 메서드. Phase 1+2+2.5+ 에서 누락된 niche/utility 영역.
**Out of Scope:** 선물옵션, 잔고/조회 (계좌 인증), ELW, 실시간 (WebSocket), Trading.

---

## §1. 목적과 결정 요약

### 목적

Phase 1+2+2.5+ 에서 핵심 시세/차트/순위/재무/투자자/업종/IPO/KSD 영역 cover. Phase 3.1 에서 채권 시세. 그러나 국내주식 영역에 여전히 niche 영역 (투자의견, 공매도, 신용잔고, 장운영정보, 체결분석, 추가 ranking 등) 30 메서드 미커버.

이 영역은 moneyflow / stock-data-batch 의 잠재 use case (애널리스트 의견, 공매도 트렌드, 시장 운영 상태 알림 등). 위험도 LOW — 모두 OAuth token only.

### 핵심 결정

| 항목 | 결정 |
|---|---|
| 기존 패키지 활용 | `domestic/` 패키지 재사용. 새 sub-package 도입 안 함. 도메인별 file 분리. |
| Sub-phase 분해 | 3 sub-phase (4.1 ~ 4.3), 각 별도 PR + minor release. PR size 7-13 메서드. |
| Style A 명명 | KIS path last segment PascalCase 1:1 매핑. `Inquire*` prefix 일관. |
| 타입 매핑 | Phase 2 standard. KSD-like all-string 미적용 (시세/순위 endpoints). |
| Release tags | v1.12.0 (4.1) ~ v1.14.0 (4.3). breaking change 없음. |

---

## §2. Sub-phase 분해

### Phase 4.1 — 종목정보/분석 (`v1.12.0`)

종목 단위 분석 + 시장 종합 정보. 10 메서드.

| 메서드 | 위치 (예상) | docs |
|---|---|---|
| `Domestic.InquireInvestOpinion` | `domestic/opinion.go` (신규) | 종목투자의견 |
| `Domestic.InquireSecBrokerOpinion` | `domestic/opinion.go` | 증권사별 투자의견 |
| `Domestic.InquireEstimatePerform` | `domestic/opinion.go` | 종목 추정실적 |
| `Domestic.InquireCcnlStrengthRank` | `domestic/extended.go` (append) | 체결강도 상위 |
| `Domestic.InquireBulkTradeNumberRank` | `domestic/extended.go` | 대량체결건수 상위 |
| `Domestic.InquireTrAmountTradeRatio` | `domestic/extended.go` | 체결금액별 매매비중 |
| `Domestic.InquireHtsTop20` | `domestic/extended.go` | HTS조회상위20 |
| `Domestic.InquireVolumeProfile` | `domestic/extended.go` | 매물대 거래비중 |
| `Domestic.InquireExpectedCcnlPriceTrend` | `domestic/extended.go` | 예상체결가 추이 |
| `Domestic.InquireExpectedUpDownRank` | `domestic/extended.go` | 예상체결 상승/하락 상위 |

**파일**: 새 `domestic/opinion.go` (3 투자의견) + 기존 `domestic/extended.go` 에 7 append.

### Phase 4.2 — 시장운영/특수상태 (`v1.13.0`)

장 운영 정보 + 시장 특수 상태 (VI/상하한가). 7 메서드.

| 메서드 | 위치 (예상) | docs |
|---|---|---|
| `Domestic.InquireMarketOpInfoKrx` | `domestic/market_op.go` (신규) | 장운영정보 (KRX) |
| `Domestic.InquireMarketOpInfoNxt` | `domestic/market_op.go` | 장운영정보 (NXT) |
| `Domestic.InquireMarketOpInfoTotal` | `domestic/market_op.go` | 장운영정보 (통합) |
| `Domestic.InquireExpectedClosingPrice` | `domestic/market_op.go` | 장마감 예상체결가 |
| `Domestic.InquireMarketHoliday` | `domestic/market_op.go` | 휴장일 조회 |
| `Domestic.InquireViStatus` | `domestic/market_op.go` | 변동성완화장치(VI) 현황 |
| `Domestic.InquireUpDownLimitCatch` | `domestic/market_op.go` | 상하한가 포착 |

**파일**: 새 `domestic/market_op.go` (7 메서드).

### Phase 4.3 — 추가 ranking/시장흐름 (`v1.14.0`)

추가 ranking + 시장 자금/공매도/신용 흐름. 13 메서드.

| 메서드 | 위치 (예상) | docs |
|---|---|---|
| `Domestic.InquireShortSaleRank` | `domestic/extended.go` (append) | 공매도 상위종목 |
| `Domestic.InquireShortSaleDailyTrend` | `domestic/extended.go` | 공매도 일별추이 |
| `Domestic.InquireCreditBalanceRank` | `domestic/extended.go` | 신용잔고 상위 |
| `Domestic.InquireCreditBalanceTrend` | `domestic/extended.go` | 신용잔고 일별추이 |
| `Domestic.InquireLendablePossible` | `domestic/extended.go` | 당사 대주가능 종목 |
| `Domestic.InquireAskingPriceRsqnRank` | `domestic/extended.go` | 호가잔량 순위 |
| `Domestic.InquireAfterHourRsqnRank` | `domestic/extended.go` | 시간외 잔량 순위 |
| `Domestic.InquireAfterHourExpFluctRate` | `domestic/extended.go` | 시간외 예상체결 등락률 |
| `Domestic.InquireMarketCapRank` | `domestic/extended.go` | 시장가치 순위 |
| `Domestic.InquireDisparityRank` | `domestic/extended.go` | 이격도 순위 |
| `Domestic.InquirePreferredStockGapRank` | `domestic/extended.go` | 우선주 괴리율 상위 |
| `Domestic.InquireProfitableAssetRank` | `domestic/extended.go` | 수익자산지표 순위 |
| `Domestic.InquireMarketFundsTotal` | `domestic/extended.go` | 증시자금 종합 |

**파일**: 기존 `domestic/extended.go` 에 13 append.

> 메서드 명명은 docs analyzer 단계에서 실 path 확인 후 확정. Style A path-based mapping.

### 합계 (Phase 4.1 ~ 4.3)

- **메서드**: 10 + 7 + 13 = **30 메서드** (누적 79 → 109)
- **Release tags**: `v1.12.0` (4.1), `v1.13.0` (4.2), `v1.14.0` (4.3)
- **PR**: 3 sub-phase = 3 PR

---

## §3. 명명 / 패턴 (Phase 2 그대로)

### Style A 메서드 명명

KIS path 의 last segment PascalCase 1:1. 충돌 시 패키지 prefix (e.g., `domestic.InquireMarketCap` vs `overseas.InquireMarketCap`) 로 회피.

### 응답 typed struct

KIS docs 1:1 매핑. Output1+Output2 verbatim.

### Params struct

`Inquire<X>Params` (zero-value default).

### 타입 매핑

- **가격류 → `decimal.Decimal` (bare tag)**
- **수량/금액 → `int64,string`**
- **비율 → `float64,string`**
- **코드/이름/날짜/Y-N → 평문 `string`**

KSD-like all-string 미적용 — Phase 4 endpoints 는 일반 시세/순위.

---

## §4. 인프라 변경 (없음)

기존 `domestic/` 패키지 재사용. 새 internal package 불필요. `domestic.New(http, master)` 시그니처 변경 없음.

신규 file 도입:
- `domestic/opinion.go` — 투자의견 3 메서드 (Phase 4.1)
- `domestic/market_op.go` — 장운영/특수상태 7 메서드 (Phase 4.2)

기존 `domestic/extended.go` (Phase 2.2 의 시간외 + 신고저근접 5 메서드) 에 Phase 4.1+4.3 의 추가 메서드 append.

---

## §5. 테스트 / 문서 / Release 흐름 (Phase 2 동일)

각 sub-phase 별:

1. testdata fixtures (rt_cd/msg_cd/msg1 envelope 필수)
2. TDD: failing test → struct + method → PASS → commit
3. examples (`examples/<sub-phase>/main.go`)
4. CLAUDE.md / README.md / CHANGELOG.md / `domestic/doc.go` 갱신
5. 최종 점검 (build/vet/fmt/race/coverage ≥ 80%)
6. PR (사용자 승인 후), merge, tag, GitHub Release

---

## §6. 진입/종료 조건

### 진입 조건

- main HEAD = v1.11.0 (Phase 3.1 publish 완료, 누적 79)
- Phase 4 design spec (이 문서) 사용자 승인

### 각 sub-phase 종료 조건

- PR merge, CI clean
- minor version tag push
- GitHub Release publish
- memory 갱신

### Phase 4 전체 종료 조건

- Phase 4.1~4.3 모두 publish (v1.14.0)
- 누적 109 메서드
- 사용자가 다음 영역 (Group B ETF/NAV, Group C 재무 추가, Group D 헬퍼, 또는 다른 도메인) 결정
