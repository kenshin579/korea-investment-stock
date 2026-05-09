# Phase 9 — WebSocket NXT/통합 변형 Design (Lightweight)

**Status:** Active design (2026-05-09)
**Goal:** Phase 8 KRX 5 endpoint 후 NXT/통합 변형 10 endpoint 추가. 누적 121 REST + 5 WS (Phase 8) + 10 WS (Phase 9) = **136 endpoints**.
**Out of Scope:** ELW/지수 실시간, 해외주식 실시간, 선물옵션 실시간, 체결통보 (AES256), Trading.

> **Lightweight spec**: Phase 8 패턴 재사용이 대부분. 이 spec 은 Phase 8 design spec (`2026-05-09-phase8-websocket-design.md`) 의 확장. 차이점만 기록.

---

## §1. EP 매핑

| EP | NXT TR_ID | 통합 TR_ID | Body fields | KRX 차이 |
|---|---|---|---|---|
| 체결가 | H0NXCNT0 | H0UNCNT0 | 46 | 22번 CNTG_CLS_CODE (KRX 는 CCLD_DVSN, 의미 동일) |
| 호가 | H0NXASP0 | H0UNASP0 | 65 | KRX 59 + KMID/NMID 6 (중간가) |
| 예상체결 | H0NXANC0 | H0UNANC0 | 46 | KRX 45 + VI_STND_PRC |
| 프로그램매매 | H0NXPGM0 | H0UNPGM0 | 11 | 신규 EP (KRX H0STPGM0 도 Phase 8 OoS) |
| 회원사 | H0NXMBC0 | H0UNMBC0 | 78 | 신규 EP (KRX H0STMBC0 도 Phase 8 OoS) |

전체 schema reference: `websocket/testdata/_schemas_phase9.md` (10 fixtures 와 함께 commit 됨).

---

## §2. 사용자 합의 결정

| 항목 | 결정 | 근거 |
|---|---|---|
| Phase scope | NXT 5 + 통합 5 = 10 EP | 사용자 명시 ("KRX 5 EP의 다른 시장 변형") |
| Event types | 5 base struct + 10 type alias | NXT/통합 schema 동일 → DRY + 시장 구분 명확 |
| 진행 절차 | 경량 spec + plan skip + 다이렉트 구현 | Phase 8 패턴 재사용이 대부분이라 plan overhead 불필요 |
| 시장별 명명 | NxtTradeEvent / UnifiedTradeEvent (alias) | base = AltMarketTradeEvent (또는 SharedTradeEvent), type alias 로 시장 명시 |

---

## §3. 핵심 deviation (Phase 8 spec 대비)

1. **Schema 재사용 부정**: NXT/통합은 KRX 와 schema 가 항상 다름 (호가 +6, 예상체결 +1, 프로그램매매/회원사 신규). KRX Event 재사용 X.
2. **NXT vs 통합은 schema 동일**: 5 base struct 로 처리 + type alias. → 사용자 입장 Subscribe/On 메서드는 시장별 (10개), 내부 decoder/struct 는 공유 (5개).
3. **22번 필드명 차이 (체결가)**: KRX = `CCLD_DVSN`, NXT/통합 = `CNTG_CLS_CODE`. 의미 동일 (체결구분). Go 필드명 `TradeKind` 통일.
4. **모든 모의 미지원**: KRX 만 모의 지원. NXT/통합은 실전 only. 별도 처리 없음 (Options.Endpoint 가 RealEnv 시 자동 동작).

---

## §4. 구조

### 4-1. Event 정의 (`websocket/events_phase9.go`)

```go
// 5 base struct (NXT/통합 공유)
type AltMarketTradeEvent struct { /* 46 fields */ }
type AltMarketAskEvent struct { /* 65 fields, KrxMidPrice/NxtMidPrice 등 추가 */ }
type AltMarketExpectTradeEvent struct { /* 46 fields, ViStandardPrice 추가 */ }
type ProgramTradeEvent struct { /* 11 fields */ }
type MemberEvent struct { /* 78 fields, 5단계 매도/매수 + 외국계 + 영문회원사명 */ }

// 10 type alias
type NxtTradeEvent = AltMarketTradeEvent
type UnifiedTradeEvent = AltMarketTradeEvent
type NxtAskEvent = AltMarketAskEvent
type UnifiedAskEvent = AltMarketAskEvent
type NxtExpectTradeEvent = AltMarketExpectTradeEvent
type UnifiedExpectTradeEvent = AltMarketExpectTradeEvent
type NxtProgramTradeEvent = ProgramTradeEvent
type UnifiedProgramTradeEvent = ProgramTradeEvent
type NxtMemberEvent = MemberEvent
type UnifiedMemberEvent = MemberEvent
```

### 4-2. Decoder (`websocket/decode_phase9.go`)

5 base decoder (NXT/통합 공유):
```go
func decodeAltMarketTrade(f frame) ([]AltMarketTradeEvent, error)
func decodeAltMarketAsk(f frame) ([]AltMarketAskEvent, error)
func decodeAltMarketExpectTrade(f frame) ([]AltMarketExpectTradeEvent, error)
func decodeProgramTrade(f frame) ([]ProgramTradeEvent, error)
func decodeMember(f frame) ([]MemberEvent, error)
```

각 decoder 의 field index → struct field 매핑은 `_schemas_phase9.md` 참고.

### 4-3. Client (확장 — `websocket/client.go`)

추가 const + 메서드:
```go
const (
    trIDNxtTrade           = "H0NXCNT0"
    trIDNxtAsk             = "H0NXASP0"
    trIDNxtExpectTrade     = "H0NXANC0"
    trIDNxtProgramTrade    = "H0NXPGM0"
    trIDNxtMember          = "H0NXMBC0"
    trIDUnifiedTrade        = "H0UNCNT0"
    trIDUnifiedAsk          = "H0UNASP0"
    trIDUnifiedExpectTrade  = "H0UNANC0"
    trIDUnifiedProgramTrade = "H0UNPGM0"
    trIDUnifiedMember       = "H0UNMBC0"
)

// Subscribe (10 추가)
func (c *Client) SubscribeNxtTrade(symbols ...string) error
func (c *Client) SubscribeUnifiedTrade(symbols ...string) error
// ... 8 more

// On (10 추가) — alias 라 type 은 base 와 동일
func (c *Client) OnNxtTrade(h func(NxtTradeEvent))
func (c *Client) OnUnifiedTrade(h func(UnifiedTradeEvent))
// ... 8 more
```

### 4-4. Dispatcher (확장 — `websocket/dispatcher.go`)

10 distinct handler fields + Route methods. NXT 와 통합은 base type 같지만 handler 다름:
```go
type dispatcher struct {
    // ... Phase 8 fields
    onNxtTrade           func(AltMarketTradeEvent)
    onUnifiedTrade       func(AltMarketTradeEvent)
    onNxtAsk             func(AltMarketAskEvent)
    onUnifiedAsk         func(AltMarketAskEvent)
    // ... 6 more
}
```

### 4-5. routeRealtime 확장 (client.go)

10 case 추가 (TR_ID → decoder + market-specific Route):
```go
case trIDNxtTrade:
    evs, err := decodeAltMarketTrade(f); ...
    for _, ev := range evs { c.dispatcher.RouteNxtTrade(ev) }
case trIDUnifiedTrade:
    evs, err := decodeAltMarketTrade(f); ...
    for _, ev := range evs { c.dispatcher.RouteUnifiedTrade(ev) }
// ... 8 more
```

---

## §5. Testing

- Decoder unit tests (5 base × Single/Paging/InvalidJSON ≈ 15 tests)
- Dispatcher tests (10 새 handler 라우팅 검증, ≈ 5 tests)
- Integration test 추가: NXT 종목 subscribe → mock 데이터 → handler 호출 (1-2 시나리오, 기존 wsmock 재사용)

목표 coverage: ≥70% (Phase 8 와 동일).

---

## §6. 진입/종료 조건

- **진입**: main HEAD = v1.18.0 (Phase 8 완료, 121 REST + 5 WS = 126 endpoints).
- **종료**: PR merge, v1.19.0 tag, GitHub Release.
- **누적**: 121 REST + 15 WS = **136 endpoints**.

---

## §7. 진행 절차

Phase 6/7 패턴: **plan 작성 skip + 직접 batch 구현**. Phase 8 와 다르게 새 인프라가 없음 (모두 재사용).

Tasks:
1. ~~_schemas_phase9.md + 10 fixtures + design spec~~ ✅ (이 commit)
2. Event types (events_phase9.go) — 5 base + 10 alias
3. Decoder (decode_phase9.go) — 5 base
4. Decoder tests (decode_phase9_test.go) — 5+ tests
5. Dispatcher 확장 (10 새 handler 필드 + 10 Route 메서드)
6. Dispatcher tests 추가
7. Client 확장 (10 Subscribe + 10 Unsubscribe + 10 On + routeRealtime 10 case)
8. Integration test 추가 (NXT 1-2 시나리오)
9. example 추가 또는 기존 ws_krx_basic 확장
10. 문서 갱신 (CLAUDE.md / README.md / CHANGELOG.md / domestic/doc.go)
11. 최종 점검 (gofmt/vet/build/race/coverage ≥70%)
12. PR + merge + tag v1.19.0 + GitHub Release

회원사 EP body field count = **78** (NXT/통합 동일, 2026-05-09 docs 응답 표 직접 검증).
