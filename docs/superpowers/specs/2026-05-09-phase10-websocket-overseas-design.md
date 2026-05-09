# Phase 10 — WebSocket 해외주식 실시간 Design (Lightweight)

**Status:** Active design (2026-05-09)
**Goal:** Phase 8 (KRX 5) + Phase 9 (NXT/통합 10) 위에 해외주식 실시간 시세 2 EP 추가. 누적 121 REST + 17 WS = **138 endpoints**.
**Out of Scope:** 해외주식 체결통보 (H0GSCNI0, AES256 별도 phase), 선물옵션 실시간, ELW/지수 실시간.

> **Lightweight spec**: Phase 8/9 인프라 그대로 재사용 (`websocket/` 패키지). 새 도메인 (해외) 만 추가.

---

## §1. EP 매핑

| EP | TR_ID | Body fields | 모의 | 비고 |
|---|---|---|---|---|
| 체결가 (지연) | HDFSCNT0 | 26 | 미지원 | 미국 0분지연 무료 / 아시아 15분지연 |
| 호가 | HDFSASP0 | 17 | 미지원 | 미국 1호가 무료 |

전체 schema reference: `websocket/testdata/_schemas_phase10.md`.

---

## §2. 사용자 합의 결정

| 항목 | 결정 | 근거 |
|---|---|---|
| Phase scope | 시세 2 EP (체결가 + 호가) | 사용자 선택 — 체결통보(CNI0) 는 AES256 인프라 필요로 별도 phase |
| Event types | 2 distinct struct (`OverseasTradeEvent`, `OverseasAskEvent`) | 해외는 단일 도메인, alias 불필요 |
| 진행 절차 | 경량 spec + plan skip + 다이렉트 구현 | Phase 6/7/9 패턴 재사용 |
| 1호가 vs 10호가 | 해외는 docs 명시 1호가만 | KIS docs HDFSASP0 응답 표 직접 검증 |

---

## §3. 핵심 deviation (Phase 8/9 spec 대비)

1. **Symbol 형식**: KRX/NXT/통합은 6자리 종목코드 (예: `005930`). 해외는 `D`/`R` + 시장구분(3자리) + 종목코드 (예: `DNASAAPL`).
2. **모든 응답 필드를 String 으로 docs 표기**: KRX 와 다르게 KIS docs 가 모든 type 을 String 으로 명시. 그러나 매핑은 KRX 패턴 따라 가격→decimal, 수량→int64, 비율→float64.
3. **시간 포맷 차이**: 해외는 `XHMS` (현지시간) + `KHMS` (한국시간) 두 가지. 일자도 `TYMD`/`XYMD`/`KYMD` 등 분리.
4. **호가 1단계만**: KRX 10단계 (ASKP1..10) 와 다르게 해외는 1단계 (PBID1/PASK1) 만.
5. **모든 EP 모의 미지원**: KRX 만 모의 지원 (Phase 8). NXT/통합/해외 모두 실전 only.

---

## §4. 구조

### 4-1. Event 정의 (`websocket/events_phase10.go`)

```go
type OverseasTradeEvent struct { /* 26 fields */ }
type OverseasAskEvent   struct { /* 17 fields */ }
```

### 4-2. Decoder (`websocket/decode_phase10.go`)

```go
func decodeOverseasTrade(f frame) ([]OverseasTradeEvent, error)
func decodeOverseasAsk(f frame) ([]OverseasAskEvent, error)
```

### 4-3. Client 확장

```go
const (
    trIDOverseasTrade = "HDFSCNT0"
    trIDOverseasAsk   = "HDFSASP0"
)

func (c *Client) SubscribeOverseasTrade(symbols ...string) error
func (c *Client) SubscribeOverseasAsk(symbols ...string) error
func (c *Client) UnsubscribeOverseasTrade(symbols ...string) error
func (c *Client) UnsubscribeOverseasAsk(symbols ...string) error
func (c *Client) OnOverseasTrade(h func(OverseasTradeEvent))
func (c *Client) OnOverseasAsk(h func(OverseasAskEvent))
```

`subscribe()`/`unsubscribe()` 헬퍼는 그대로 재사용. `routeRealtime` 에 2 case 추가.

### 4-4. Dispatcher 확장

`onOverseasTrade` / `onOverseasAsk` 2 슬롯 + On/Route 메서드.

---

## §5. Testing

- Decoder unit tests (2 fixture × 1 each + FieldCountMismatch + BadNumeric)
- Dispatcher tests (2 새 handler 라우팅)
- Integration test 1 시나리오 (OverseasTrade subscribe + 메시지 수신)

목표 coverage: ≥70% (Phase 8/9 동일).

---

## §6. 진입/종료 조건

- **진입**: main HEAD = v1.19.0 (Phase 9 완료, 121 REST + 15 WS).
- **종료**: PR merge, v1.20.0 tag, GitHub Release.
- **누적**: 121 REST + 17 WS = **138 endpoints**.

---

## §7. 진행 절차

Phase 6/7/9 패턴: plan 작성 skip + 직접 batch 구현.

Tasks:
1. _schemas_phase10.md + 2 fixtures + design spec ✅ (이 commit)
2. events_phase10.go (2 struct)
3. decode_phase10.go (2 decoder)
4. decode_phase10_test.go
5. dispatcher / client 확장
6. dispatcher_test.go / integration_test.go 추가
7. 문서 갱신 (CLAUDE.md / README.md / CHANGELOG.md / doc.go)
8. example 추가 (`examples/ws_overseas_basic/`)
9. 최종 점검 (gofmt/vet/build/race/coverage ≥70%)
10. PR + merge + tag v1.20.0 + GitHub Release
