# Phase 8 — WebSocket Real-time (Domestic KRX) Design

**Status:** Active design (2026-05-09)
**Goal:** Phase 1~7 (121 REST methods) 후 첫 WebSocket 도입. 인프라 + 국내주식 KRX 실시간 시세 5 endpoint.
**Out of Scope:** NXT/통합 변형, ELW, 지수 실시간, 해외주식 실시간, 선물옵션 실시간, 체결통보 (암호화), Trading.

---

## §1. 목적

KIS 의 실시간 (WebSocket) 도메인을 처음 도입. 시세 모니터링 시나리오 (체결가/호가/예상체결 + 시간외) 5 endpoint 를 통해 인프라 (인증, 연결, 구독, 재연결, decoding) 를 정착시키고, Phase 9+ 에서 NXT/통합/해외/선물옵션 실시간으로 확장.

**Phase 8 = 인프라 + 국내주식 KRX 시세 5 endpoint.**

---

## §2. Endpoint 매핑

| EP | TR_ID | 한글명 | docs |
|---|---|---|---|
| 1 | H0STCNT0 | 국내주식 실시간체결가 (KRX) | `docs/api/국내주식/국내주식_실시간체결가_(KRX).md` |
| 2 | H0STASP0 | 국내주식 실시간호가 (KRX) | `docs/api/국내주식/국내주식_실시간호가_(KRX).md` |
| 3 | H0STANC0 | 국내주식 실시간예상체결 (KRX) | `docs/api/국내주식/국내주식_실시간예상체결_(KRX).md` |
| 4 | H0STOAC0 | 국내주식 시간외 실시간체결가 (KRX) | `docs/api/국내주식/국내주식_시간외_실시간체결가_(KRX).md` |
| 5 | H0STOAA0 | 국내주식 시간외 실시간예상체결 (KRX) | `docs/api/국내주식/국내주식_시간외_실시간예상체결_(KRX).md` |

> 위 TR_ID 는 docs 파일명 기준. 구현 단계에서 docs analyzer 가 정확한 TR_ID/필드 schema/필드 개수 추출.

추가 인프라 endpoint:
- **`/oauth2/Approval`** (REST POST) — WebSocket 접속키 (approval_key) 발급. WS 패키지 내부에서만 호출.

---

## §3. 사용자 합의 결정

| 항목 | 결정 | 근거 |
|---|---|---|
| Phase scope | 인프라 + 국내주식 KRX 시세 5 endpoint | 58 endpoints 한 phase 에 못 넣음. 가장 자주 쓰이는 KRX 시세부터. |
| API style | Callback handler (`OnXxx(handler)`) | 한 connection 에 multi tr_id 가 fan-in 되는 KIS 모델에 fit. Java/Python WebSocket SDK 패턴과 일관. |
| Reconnect | 자동 재연결 + 구독 자동 복원 (exp backoff) | KIS 는 장 시작/마감 일시 끊김 빈번. 사용자 책임으로 두면 boilerplate. |
| Architecture | Approach A — 신규 top-level `websocket/` 패키지 | 기존 `client.Domestic`/`client.Overseas` 패턴과 일관. cross-domain reuse 자연스러움. |
| WS 라이브러리 | `coder/websocket` (구 nhooyr.io/websocket) | modern, context-aware, std-lib only deps. |

---

## §4. Architecture

### 4-1. Top-level 통합

```
kis.Client
  ├── Domestic   (기존)
  ├── Overseas   (기존)
  ├── Bonds      (기존)
  └── WS         ← 신규 *websocket.Client
```

### 4-2. WS Client 내부 구조

```
WS Client
  ├── ApprovalKeyManager   /oauth2/Approval 호출, 24h 만료 자동 재발급 (보수적 23h TTL)
  ├── ConnectionManager    coder/websocket dial, ping/pong, lifecycle
  ├── Dispatcher           TR_ID → 등록된 handler 라우팅
  ├── Subscriber           활성 구독 (tr_id, tr_key) 추적, reconnect 시 자동 복원
  ├── ReconnectController  exp backoff (initial 1s, max 30s, max 10 attempts)
  └── Decoders             TR_ID 별 caret-separated → typed Event
```

### 4-3. 스레드 모델

- `Run(ctx)` 가 blocking. 내부에 1개 reader goroutine.
- Reader 가 raw frame 수신 → Decoder type 변환 → Dispatcher 가 사용자 handler 호출.
- Handler 는 reader goroutine 에서 동기 실행 (사용자가 무거운 작업이면 본인 channel 로 fan-out).
- Subscribe/Unsubscribe 는 caller goroutine 에서 즉시 frame 송신 (mutex 보호).

### 4-4. Connection model

- Single connection per WS Client 인스턴스. KIS 정책상 single approval_key = single session.
- Multi-connection (예: 2 인스턴스 = 2 connection) 은 사용자 책임. Phase 9+ 에서 결정.

---

## §5. Components

### 5-1. ApprovalKeyManager (`websocket/approval.go`)

```go
type ApprovalKeyManager struct {
    httpClient HTTPClient
    appKey     string
    appSecret  string
    mu         sync.Mutex
    cached     string
    expiry     time.Time
}

func (m *ApprovalKeyManager) Get(ctx context.Context) (string, error)
// 캐시된 key 가 valid 하면 즉시 반환
// expired/없음 이면 POST /oauth2/Approval → 새 key 저장 → 반환
// 23h TTL (KIS 24h 보다 보수적)
```

기존 `internal/httpclient` 재사용 (POST /oauth2/Approval).

### 5-2. WS Client 공개 API

```go
package websocket

type Client struct { /* unexported */ }

// Lifecycle — Run 은 blocking, ctx 끝나면 graceful close.
func (c *Client) Run(ctx context.Context) error

// 구독 (5 endpoint)
func (c *Client) SubscribeKrxTrade(symbols ...string) error           // H0STCNT0
func (c *Client) SubscribeKrxAsk(symbols ...string) error             // H0STASP0
func (c *Client) SubscribeKrxExpectTrade(symbols ...string) error     // H0STANC0
func (c *Client) SubscribeKrxOvernightTrade(symbols ...string) error  // H0STOAC0
func (c *Client) SubscribeKrxOvernightExpect(symbols ...string) error // H0STOAA0

// 해제 (대칭)
func (c *Client) UnsubscribeKrxTrade(symbols ...string) error
func (c *Client) UnsubscribeKrxAsk(symbols ...string) error
func (c *Client) UnsubscribeKrxExpectTrade(symbols ...string) error
func (c *Client) UnsubscribeKrxOvernightTrade(symbols ...string) error
func (c *Client) UnsubscribeKrxOvernightExpect(symbols ...string) error

// Handler 등록 (later 등록 OK — Run() 전이든 후든)
func (c *Client) OnKrxTrade(h func(KrxTradeEvent))
func (c *Client) OnKrxAsk(h func(KrxAskEvent))
func (c *Client) OnKrxExpectTrade(h func(KrxExpectTradeEvent))
func (c *Client) OnKrxOvernightTrade(h func(KrxTradeEvent))      // 동일 schema 재사용
func (c *Client) OnKrxOvernightExpect(h func(KrxExpectTradeEvent))

// 인프라 이벤트
func (c *Client) OnConnected(h func())
func (c *Client) OnReconnect(h func(attempt int))
func (c *Client) OnDisconnect(h func(err error))
func (c *Client) OnError(h func(err error))
```

옵션:
```go
type Options struct {
    Endpoint        string         // default 실전 ws://ops.koreainvestment.com:21000, 모의 :31000
    MaxReconnects   int            // default 10
    ReconnectMin    time.Duration  // default 1s
    ReconnectMax    time.Duration  // default 30s
    ApprovalTTL     time.Duration  // default 23h
    Logger          *slog.Logger   // default discard
    CustType        string         // default "P" (개인). 법인 "B".
}
```

### 5-3. Event Structs (체결가 예시)

```go
type KrxTradeEvent struct {
    Symbol         string          // MKSC_SHRN_ISCD
    Time           string          // STCK_CNTG_HOUR (HHMMSS) — string 보존
    Price          decimal.Decimal // STCK_PRPR
    PrevDiffSign   string          // PRDY_VRSS_SIGN (1=상한, 2=상승, 3=보합, 4=하한, 5=하락)
    PrevDiff       decimal.Decimal
    PrevChangeRate float64
    Open           decimal.Decimal
    High           decimal.Decimal
    Low            decimal.Decimal
    Ask1           decimal.Decimal
    Bid1           decimal.Decimal
    TradeVolume    int64           // CNTG_VOL
    AccumVolume    int64           // ACML_VOL
    AccumValue     int64           // ACML_TR_PBMN
    // ... 전체 필드 (구현 단계 docs analyzer 결과 기반)
    Raw            []string        // caret 분리 원본 (escape hatch)
}
```

다른 Event 들 (`KrxAskEvent`, `KrxExpectTradeEvent`) 도 동일 패턴.

### 5-4. 내부 구조 (구현 단계에서 결정될 file 분리)

- `websocket/client.go` — public Client + Options + lifecycle
- `websocket/approval.go` — ApprovalKeyManager
- `websocket/conn.go` — ConnectionManager (coder/websocket wrapper)
- `websocket/subscriber.go` — Subscriber (활성 구독 map)
- `websocket/dispatcher.go` — TR_ID → handler 라우팅
- `websocket/reconnect.go` — ReconnectController (backoff)
- `websocket/decode_krx.go` — 5 endpoint decoder
- `websocket/events_krx.go` — Event struct
- `websocket/errors.go` — typed errors
- `websocket/internal/wsmock/server.go` — test mock server
- `websocket/testdata/` — caret payload fixtures

---

## §6. Data Flow

### 6-1. Subscribe 흐름

```
사용자: client.WS.SubscribeKrxTrade("005930", "000660")
    │
    ▼
[Subscriber.Add] — map 에 (H0STCNT0, "005930") + (H0STCNT0, "000660") 등록
    │
    ▼ 연결 상태 분기
    ├─ 연결 X: map 저장 후 return nil (Run() 시 dial 후 일괄 송신)
    └─ 연결 O: 즉시 send subscribe frame (각 종목당 1 frame)
            │
            ▼
[Encoder] tr_type="1" (등록) → JSON
    {
      "header": {"approval_key":"...", "custtype":"P", "tr_type":"1", "content-type":"utf-8"},
      "body":   {"input": {"tr_id":"H0STCNT0", "tr_key":"005930"}}
    }
    ▼
[Conn.Write] WebSocket text frame
```

### 6-2. 수신 흐름 (reader goroutine)

```
[Conn.Read] WebSocket text frame
    │
    ▼ 첫 글자 분기
    ├─ "{" → JSON 메시지
    │     ├─ rt_cd="0" + msg1="SUBSCRIBE SUCCESS" → log (또는 OnConnected)
    │     ├─ tr_id="PINGPONG" → 즉시 PINGPONG 응답
    │     └─ msg_cd ≠ "OPSP0000" → OnError(WSServerError)
    │
    └─ "0|" or "1|" → 실시간 데이터
          │
          ▼ split("|") → [encryptFlag, trID, count, payload]
          │
          ▼ encryptFlag="1" → Phase 8 에서는 OnError(ErrWSEncryptedNotSupported), skip
          │
          ▼ payload.split("^") → []string fields
          │
          ▼ count > 1 이면 fields 를 endpoint 의 fieldCount 단위로 chunk
          │
          ▼ [Decoder.For(trID)] 각 chunk → typed Event
          │
          ▼ [Dispatcher.Route(trID, event)] 등록된 handler 호출 (동기, reader goroutine 내)
```

### 6-3. Reconnect 흐름

```
Conn.Read 에러 (context.Canceled 가 아닌 net err)
    │
    ▼ OnDisconnect(err) 호출
    │
    ▼ [ReconnectController]
    │   attempts++, backoff = min(1s * 2^attempts, 30s)
    │   attempts > maxAttempts (default 10) 면 ErrWSGiveUp → OnError → Run() 종료
    │   time.Sleep(backoff)
    │
    ▼ ApprovalKeyManager.Get(ctx) — 23h 지났으면 자동 재발급
    │
    ▼ ConnectionManager.Dial(ctx) — 새 ws 세션
    │
    ▼ Subscriber.RestoreAll() — map 의 모든 (tr_id, tr_key) 일괄 등록 frame 송신
    │
    ▼ OnReconnect(attempt) 호출
    │
    ▼ attempts = 0 (성공) → reader goroutine 재개
```

### 6-4. Graceful shutdown

```
ctx.Done()
    │
    ▼ reader goroutine: 진행 중 read 취소
    ▼ Conn.Close(websocket.StatusNormalClosure)
    ▼ Run() return ctx.Err()
```

### 6-5. Field 매핑 정책

| 결정 | 선택 |
|---|---|
| 숫자 파싱 실패 | 해당 필드 zero value, OnError 호출 X (data corruption 은 KIS bug) |
| Time field | string `HHMMSS` 보존 (Go time 변환은 사용자 책임) |
| count > 1 페이징 | Decoder 자동 split → handler 가 N번 호출 (각 chunk 마다) |
| Raw escape hatch | `Event.Raw []string` 으로 caret 원본 보존 (debug + future-proof) |

---

## §7. Error Handling

### 7-1. Typed errors

```go
// websocket/errors.go
var (
    ErrWSNotConnected            = errors.New("kis ws: not connected")
    ErrWSGiveUp                  = errors.New("kis ws: give up after max reconnect attempts")
    ErrWSApprovalFailed          = errors.New("kis ws: approval key issuance failed")
    ErrWSInvalidFrame            = errors.New("kis ws: invalid frame format")
    ErrWSDuplicateSub            = errors.New("kis ws: duplicate subscription")
    ErrWSEncryptedNotSupported   = errors.New("kis ws: encrypted frames not supported in Phase 8")
)

type WSServerError struct {
    TrID  string
    MsgCd string
    Msg   string
}
func (e *WSServerError) Error() string { ... }
```

### 7-2. 에러 라우팅

| 발생 위치 | 처리 |
|---|---|
| Subscribe frame 송신 실패 | Subscribe 메서드 직접 return error (OnError 호출 X) |
| Read 에러 (net.OpError, EOF, abnormal close) | OnDisconnect → Reconnect → 실패 시 OnError(ErrWSGiveUp) |
| Frame 파싱 실패 (`ErrWSInvalidFrame`) | OnError 만 호출, 연결 유지 (1 frame 손실) |
| 숫자 파싱 실패 | 해당 필드 zero value, **silent** |
| 서버 등록 실패 (msg_cd ≠ "OPSP0000") | OnError(WSServerError) |
| Approval key 만료 + 재발급 실패 | Reconnect 카운트 → ErrWSGiveUp |
| Handler panic | recover → OnError 로 라우팅 (다른 frame 처리 계속) |

### 7-3. PINGPONG / Keepalive

- KIS 서버 송신 `tr_id="PINGPONG"` JSON → 자동 PONG 응답 (사용자 노출 X).
- `coder/websocket` 의 `Conn.Ping(ctx)` 30s interval 호출 (TCP-level keepalive).

### 7-4. Backpressure

- Handler 동기 실행 → 느린 handler 가 read 차단 가능.
- 문서 명시: "handler 는 빠르게 return 해야 함. 무거운 작업은 사용자가 channel 로 fan-out".
- Phase 8 에서 buffered queue 도입 X (YAGNI). Phase 9+ 옵션.

### 7-5. Logging

- `*slog.Logger` 받음 (`Options.Logger`).
- 등록 frame 송신 / reconnect attempt / parse 실패 / give-up 시 로그.
- Default: silent (`io.Discard`).

---

## §8. Testing

### 8-1. Unit tests

| 대상 | 테스트 |
|---|---|
| Decoder | 5 endpoint 의 caret 샘플 → Event struct (testdata fixtures). count=1 / count>1 / 빈 필드. |
| Subscriber | Add/Remove map 상태 / 중복 등록 거부 / RestoreAll frame list. |
| ReconnectController | backoff 계산 (1s/2s/4s/.../30s) / maxAttempts 도달 → ErrWSGiveUp. |
| Frame parser | `0|...` / JSON / PINGPONG / `1|` 암호화 (OnError + skip). |

### 8-2. Integration tests — local mock server

- `websocket/internal/wsmock/server.go` — `httptest.NewServer` + `coder/websocket.Accept`.
- 시나리오:
  1. Happy path: Subscribe → mock 송신 → handler 호출.
  2. Multi-symbol fan-out: 2 종목 → 분리 라우팅.
  3. Reconnect: mock close → 재연결 → 구독 복원 → OnReconnect.
  4. PINGPONG: mock PING → SDK PONG.
  5. Server error: mock 가 `OPSP0001` → OnError(WSServerError).
  6. Graceful shutdown: ctx.Done → Run return + close frame 송신.

### 8-3. Approval key

- `httpmock` 으로 `POST /oauth2/Approval` 응답 mock.
- TTL 만료 직전/직후 동작.

### 8-4. Coverage

- `websocket/` 패키지 ≥ 70% (network/timing 코드 한계 인정 — 기존 패키지 80% 보다 낮은 목표).
- Decoder/Subscriber 단독 ≥ 90%.

### 8-5. Manual smoke test

- `examples/ws_krx_basic/main.go` — 실제 KIS 환경 KOSPI 종목 1개 체결가 구독 시연.

---

## §9. 진입/종료 조건

- **진입**: main HEAD = v1.17.0 (Phase 7 완료, 121 메서드).
- **종료**: PR merge, v1.18.0 tag, GitHub Release.
- **누적**: 121 REST + 5 WebSocket = **126 endpoints**.

---

## §10. Out of Scope (Phase 8)

- NXT / 통합 변형 (`H0NXCNT0` / `H0UNCNT0` 등) → Phase 9.
- ELW 실시간, 국내지수 실시간 → Phase 9+.
- 해외주식 실시간 (HDFSCNT0 등) → Phase 10+.
- 선물옵션 실시간 21 endpoints → Phase 11+.
- **체결통보 (H0STCNI0/H0STCNI9)** — AES256 복호화 필요. 별도 phase.
- 회원사/프로그램매매 실시간 → 후속 phase 결정.
- Multi-connection (2+ WS 인스턴스) → Phase 9+ 에서 필요 시.
- Backpressure buffered queue → 사용자 요구 발생 시.

---

## §11. 진행 절차

Phase 8 은 첫 architecture 변경 (REST → WebSocket) 으로 **plan 작성 필수** (Phase 4.3/5/6/7 처럼 직접 구현 X).

Step:
1. Brainstorming spec (이 문서) 작성 + commit ✅
2. **`writing-plans` skill 호출** → 상세 implementation plan 작성 (sub-tasks decomposition).
3. Plan 승인 후 구현 시작 (가능 시 subagent dispatch 로 병렬화 — decoder 5 endpoint 등).
4. 구현 → unit + integration 테스트 → manual smoke → PR → v1.18.0 release.
