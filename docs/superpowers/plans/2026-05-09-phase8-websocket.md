# Phase 8 WebSocket Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** KIS WebSocket 도메인을 처음 도입. 인프라 (인증/연결/구독/재연결/decoding) + 국내주식 KRX 시세 5 endpoint 구현. 누적 121 → 126 endpoints, v1.18.0 release.

**Architecture:** 신규 top-level `websocket/` 패키지에 `coder/websocket` (구 nhooyr.io) 기반 callback handler 패턴 구현. `kis.Client.WS` 가 진입점. ApprovalKeyManager (23h TTL) + ConnectionManager + Dispatcher + Subscriber + ReconnectController 모듈로 분리. 자동 재연결 + 구독 자동 복원 (exp backoff, max 10).

**Tech Stack:** Go 1.25+, `github.com/coder/websocket`, `github.com/shopspring/decimal`, `github.com/jarcoal/httpmock` (test), `github.com/stretchr/testify` (test). 기존 `internal/httpclient` 재사용.

**Spec:** `docs/superpowers/specs/2026-05-09-phase8-websocket-design.md`

---

## File Structure

```
go.mod                                   # +github.com/coder/websocket dep
client.go                                # +WS field, NewClient* 가 ws.Client 주입
websocket/                               # 신규 top-level 패키지
  doc.go                                 # 패키지 doc
  errors.go                              # typed errors + WSServerError
  events_krx.go                          # 5 Event struct (KrxTradeEvent 등)
  subscriber.go                          # 활성 구독 map + Add/Remove/RestoreAll
  subscriber_test.go
  reconnect.go                           # ReconnectController (exp backoff)
  reconnect_test.go
  frame.go                               # 첫 글자 분기 (JSON / 0|/ 1|) + chunk 분리
  frame_test.go
  decode_krx.go                          # caret-separated → Event (5 EP)
  decode_krx_test.go
  approval.go                            # ApprovalKeyManager (POST /oauth2/Approval)
  approval_test.go
  conn.go                                # ConnectionManager (coder/websocket wrapper)
  dispatcher.go                          # TR_ID → handler 라우팅 + recover
  dispatcher_test.go
  client.go                              # public Client + Options + Run + Subscribe* + On*
  client_test.go                         # public API smoke
  testdata/                              # caret payload fixtures (5 EP)
  internal/wsmock/
    server.go                            # local mock KIS WebSocket server
examples/ws_krx_basic/main.go            # 실 KIS smoke 시연
README.md                                # +Phase 8 section
CLAUDE.md                                # phase 7 → phase 8
CHANGELOG.md                             # [1.18.0] entry
```

---

## Task 1: Setup — feature branch, deps, scaffold

**Files:**
- Modify: `go.mod`, `go.sum`
- Create: `websocket/doc.go`, `websocket/testdata/.gitkeep`

- [ ] **Step 1: feature 브랜치 생성**

```bash
git checkout main
git pull origin main
git checkout -b feat/phase8-websocket
```

Expected: `Switched to a new branch 'feat/phase8-websocket'`

- [ ] **Step 2: WebSocket 라이브러리 추가**

```bash
cd /Users/user/src/workspace_moneyflow/korea-investment-stock
go get github.com/coder/websocket@latest
go mod tidy
```

Expected: `go.mod` 에 `github.com/coder/websocket vX.Y.Z` 추가됨. `go.sum` 갱신.

- [ ] **Step 3: 패키지 디렉토리 + doc.go 생성**

`websocket/doc.go`:
```go
// Package websocket 은 KIS 실시간 (WebSocket) 도메인.
//
// 한투 docs: docs/api/기타/실시간_(웹소켓)_접속키_발급.md (인증)
// docs/api/국내주식/국내주식_실시간*.md (KRX 5 EP)
//
// 디자인: docs/superpowers/specs/2026-05-09-phase8-websocket-design.md
//
// Phase 8 — 인프라 + 국내주식 KRX 시세 5 endpoint:
//
//	H0STCNT0  실시간체결가 (KRX)        SubscribeKrxTrade / OnKrxTrade
//	H0STASP0  실시간호가 (KRX)          SubscribeKrxAsk / OnKrxAsk
//	H0STANC0  실시간예상체결 (KRX)      SubscribeKrxExpectTrade / OnKrxExpectTrade
//	H0STOAC0  시간외 실시간체결가 (KRX) SubscribeKrxOvernightTrade / OnKrxOvernightTrade
//	H0STOAA0  시간외 실시간예상체결 (KRX) SubscribeKrxOvernightExpect / OnKrxOvernightExpect
//
// 사용자는 root kis.Client 의 WS 필드로 접근.
package websocket
```

`websocket/testdata/.gitkeep`: 빈 파일 (디렉토리 git 추적용).

- [ ] **Step 4: build 검증**

```bash
go build ./...
```

Expected: 출력 없음 (clean build). `websocket/` 빈 패키지 빌드 OK.

- [ ] **Step 5: commit**

```bash
git add go.mod go.sum websocket/
git commit -m "[chore] Phase 8 — websocket 패키지 scaffold + coder/websocket dep"
```

---

## Task 2: docs analyzer 5 EP — TR_ID + field schema 추출

**Files:**
- Read-only: 5 docs files
- Create: `websocket/testdata/h0stcnt0_success.txt`, `h0stasp0_success.txt`, `h0stanc0_success.txt`, `h0stoac0_success.txt`, `h0stoaa0_success.txt`
- Create: `websocket/testdata/_schemas.md` (reference for next tasks)

> 이 task 는 5 docs analyzer 를 병렬 dispatch 해 정확한 TR_ID 와 field schema 를 추출하는 단계. 결과로 다음 task 들의 type / decoder / fixture 가 결정됨.

- [ ] **Step 1: 5 subagent (Explore type) 동시 dispatch**

각 subagent prompt 패턴:
```
korea-investment-stock 의 KIS WebSocket docs 파일을 분석.
대상: docs/api/국내주식/<file>.md

추출:
1. 정확한 TR_ID (실전/모의)
2. tr_key 형식 (종목코드 6자리 등)
3. Response output 의 caret-separated 필드 list:
   - 필드명 (대문자, KIS docs 그대로)
   - 한글 설명
   - 타입 (String/Number — Number 는 정수 vs 소수 추정)
   - 길이
4. 필드 총 개수
5. count > 1 페이징 사용 여부
6. 특이 사항 (encrypted=1 가능 여부, 등)

출력: markdown table.
```

5 docs:
- `docs/api/국내주식/국내주식_실시간체결가_(KRX).md` → H0STCNT0
- `docs/api/국내주식/국내주식_실시간호가_(KRX).md` → H0STASP0 (예상)
- `docs/api/국내주식/국내주식_실시간예상체결_(KRX).md` → H0STANC0 (예상)
- `docs/api/국내주식/국내주식_시간외_실시간체결가_(KRX).md` → H0STOAC0 (예상)
- `docs/api/국내주식/국내주식_시간외_실시간예상체결_(KRX).md` → H0STOAA0 (예상)

- [ ] **Step 2: 결과 정리해 `websocket/testdata/_schemas.md` 작성**

5 EP 의 정확한 TR_ID + 필드 list 를 markdown table 로 정리. 다음 task 들이 reference.

- [ ] **Step 3: 5 fixture 파일 생성 (caret 샘플)**

각 fixture 는 `0|<TR_ID>|<count>|<caret-separated payload>` 형식. 실제 KIS 데이터 형태.

예: `websocket/testdata/h0stcnt0_success.txt` (체결가 — count=1, 005930):
```
0|H0STCNT0|001|005930^123929^73100^2^1500^2.09^72850^72500^73200^72400^73100^73000^150^123456^987654000000^...
```

총 필드 수는 schema 결과에 맞춤. 길이가 다양해 fixture 마다 60~100자 정도.

추가 fixture: `h0stcnt0_paging.txt` (count=2, 같은 종목의 2 frame 페이징 검증용).

- [ ] **Step 4: schema 검토 — Phase 8 spec 의 Event struct 명명과 일치 확인**

`_schemas.md` 의 필드명 ↔ spec §5-3 의 Event field 가 일관한지. 불일치 있으면 spec 의 의도 (Go-friendly 이름) 우선.

- [ ] **Step 5: commit**

```bash
git add websocket/testdata/
git commit -m "[chore] Phase 8 — 5 EP schema 추출 + caret fixture"
```

---

## Task 3: errors.go — typed errors

**Files:**
- Create: `websocket/errors.go`
- Test: `websocket/errors_test.go`

- [ ] **Step 1: 실패 테스트 작성**

`websocket/errors_test.go`:
```go
package websocket_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kenshin579/korea-investment-stock/websocket"
)

func TestWSServerError_Error(t *testing.T) {
	e := &websocket.WSServerError{TrID: "H0STCNT0", MsgCd: "OPSP0001", Msg: "ALREADY IN SUBSCRIBE"}
	assert.Contains(t, e.Error(), "H0STCNT0")
	assert.Contains(t, e.Error(), "OPSP0001")
	assert.Contains(t, e.Error(), "ALREADY IN SUBSCRIBE")
}

func TestSentinelErrors(t *testing.T) {
	assert.NotNil(t, websocket.ErrWSGiveUp)
	assert.NotNil(t, websocket.ErrWSApprovalFailed)
	assert.NotNil(t, websocket.ErrWSInvalidFrame)
	assert.NotNil(t, websocket.ErrWSNotConnected)
	assert.NotNil(t, websocket.ErrWSDuplicateSub)
	assert.NotNil(t, websocket.ErrWSEncryptedNotSupported)
	// errors.Is 동작 검증
	wrapped := errors.Join(websocket.ErrWSGiveUp, errors.New("net err"))
	assert.True(t, errors.Is(wrapped, websocket.ErrWSGiveUp))
}
```

- [ ] **Step 2: 테스트 실패 확인**

```bash
go test ./websocket/ -run TestWSServerError -v
```

Expected: FAIL — `undefined: websocket.WSServerError`

- [ ] **Step 3: 구현**

`websocket/errors.go`:
```go
package websocket

import (
	"errors"
	"fmt"
)

var (
	ErrWSNotConnected          = errors.New("kis ws: not connected")
	ErrWSGiveUp                = errors.New("kis ws: give up after max reconnect attempts")
	ErrWSApprovalFailed        = errors.New("kis ws: approval key issuance failed")
	ErrWSInvalidFrame          = errors.New("kis ws: invalid frame format")
	ErrWSDuplicateSub          = errors.New("kis ws: duplicate subscription")
	ErrWSEncryptedNotSupported = errors.New("kis ws: encrypted frames not supported in Phase 8")
)

// WSServerError 는 KIS 서버가 등록/해제 응답에서 반환한 에러.
type WSServerError struct {
	TrID  string // H0STCNT0
	MsgCd string // OPSP0001 등
	Msg   string // "ALREADY IN SUBSCRIBE" 등
}

func (e *WSServerError) Error() string {
	return fmt.Sprintf("kis ws server error: tr_id=%s msg_cd=%s msg=%s", e.TrID, e.MsgCd, e.Msg)
}
```

- [ ] **Step 4: 테스트 통과 확인**

```bash
go test ./websocket/ -v
```

Expected: PASS — TestWSServerError_Error, TestSentinelErrors.

- [ ] **Step 5: commit**

```bash
git add websocket/errors.go websocket/errors_test.go
git commit -m "[feat] Phase 8 — websocket typed errors + WSServerError"
```

---

## Task 4: events_krx.go — Event structs (5 EP)

**Files:**
- Create: `websocket/events_krx.go`

> Task 2 의 schema 결과를 토대로 5 Event struct 정의. 시간외 (H0STOAC0/H0STOAA0) 가 KRX 본장 (H0STCNT0/H0STANC0) 와 동일 schema 면 alias type 으로 재사용.

- [ ] **Step 1: 5 Event struct 정의**

`websocket/events_krx.go`:
```go
package websocket

import "github.com/shopspring/decimal"

// KrxTradeEvent — H0STCNT0 (체결가, KRX) + H0STOAC0 (시간외 체결가, KRX).
//
// 두 EP 의 caret-separated schema 는 동일 (Task 2 schema 검증 결과).
type KrxTradeEvent struct {
	Symbol         string          // MKSC_SHRN_ISCD
	Time           string          // STCK_CNTG_HOUR (HHMMSS)
	Price          decimal.Decimal // STCK_PRPR
	PrevDiffSign   string          // PRDY_VRSS_SIGN (1~5)
	PrevDiff       decimal.Decimal // PRDY_VRSS
	PrevChangeRate float64         // PRDY_CTRT
	WeightedAvg    decimal.Decimal // WGHN_AVRG_STCK_PRC
	Open           decimal.Decimal // STCK_OPRC
	High           decimal.Decimal // STCK_HGPR
	Low            decimal.Decimal // STCK_LWPR
	Ask1           decimal.Decimal // ASKP1
	Bid1           decimal.Decimal // BIDP1
	TradeVolume    int64           // CNTG_VOL
	AccumVolume    int64           // ACML_VOL
	AccumValue     int64           // ACML_TR_PBMN
	AskCount       int64           // SELN_CNTG_CSNU
	BidCount       int64           // SHNU_CNTG_CSNU
	NetCount       int64           // NTBY_CNTG_CSNU
	TradeStrength  float64         // CTTR
	TotalAskVolume int64           // SELN_CNTG_SMTN
	TotalBidVolume int64           // SHNU_CNTG_SMTN
	TradeKind      string          // CCLD_DVSN (1=매수, 3=장전, 5=매도)
	BidRate        float64         // SHNU_RATE
	PrevVolRate    float64         // PRDY_VOL_VRSS_ACML_VOL_RATE
	OpenTime       string          // OPRC_HOUR
	OpenDiffSign   string          // OPRC_VRSS_PRPR_SIGN
	OpenDiff       decimal.Decimal // OPRC_VRSS_PRPR
	HighTime       string          // HGPR_HOUR
	HighDiffSign   string          // HGPR_VRSS_PRPR_SIGN
	HighDiff       decimal.Decimal // HGPR_VRSS_PRPR
	// ... Task 2 schema 결과의 나머지 필드 (총 ~46 필드)
	Raw []string // caret 분리 원본 (escape hatch)
}

// KrxAskEvent — H0STASP0 (호가, KRX).
type KrxAskEvent struct {
	Symbol  string          // MKSC_SHRN_ISCD
	Time    string          // BSOP_HOUR (HHMMSS)
	Hour    string          // HOUR_CLS_CODE
	Ask     [10]decimal.Decimal // ASKP1..10
	Bid     [10]decimal.Decimal // BIDP1..10
	AskSize [10]int64       // ASKP_RSQN1..10
	BidSize [10]int64       // BIDP_RSQN1..10
	// ... Task 2 schema 결과의 나머지 필드
	Raw []string
}

// KrxExpectTradeEvent — H0STANC0 (예상체결, KRX) + H0STOAA0 (시간외 예상체결, KRX).
//
// 두 EP 의 schema 는 동일 (Task 2 schema 검증 결과).
type KrxExpectTradeEvent struct {
	Symbol     string          // MKSC_SHRN_ISCD
	Time       string          // STCK_CNTG_HOUR
	ExpectPrice decimal.Decimal // ANTC_CNQN_PRC 등 (Task 2 결과 반영)
	// ... 나머지 필드
	Raw []string
}
```

> Task 2 의 결과로 위 일부 필드명이 변경될 수 있음. 일관성 우선: spec §5-3 의 Go-friendly 이름 (camelCase 영어) 사용.

- [ ] **Step 2: build 검증**

```bash
go build ./websocket/
```

Expected: 출력 없음 (clean).

- [ ] **Step 3: gofmt 검증**

```bash
gofmt -l websocket/events_krx.go
```

Expected: 출력 없음.

- [ ] **Step 4: doc 검증**

`go doc ./websocket KrxTradeEvent` 가 한글 설명 포함해서 출력되는지.

- [ ] **Step 5: commit**

```bash
git add websocket/events_krx.go
git commit -m "[feat] Phase 8 — 5 EP Event struct (KrxTradeEvent / KrxAskEvent / KrxExpectTradeEvent)"
```

---

## Task 5: subscriber.go — 활성 구독 map

**Files:**
- Create: `websocket/subscriber.go`
- Test: `websocket/subscriber_test.go`

- [ ] **Step 1: 실패 테스트 작성**

`websocket/subscriber_test.go`:
```go
package websocket

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSubscriber_Add(t *testing.T) {
	s := newSubscriber()
	added, err := s.Add("H0STCNT0", "005930")
	assert.NoError(t, err)
	assert.True(t, added)

	// 중복
	added, err = s.Add("H0STCNT0", "005930")
	assert.NoError(t, err)
	assert.False(t, added) // 이미 존재 → false
}

func TestSubscriber_Remove(t *testing.T) {
	s := newSubscriber()
	s.Add("H0STCNT0", "005930")
	removed := s.Remove("H0STCNT0", "005930")
	assert.True(t, removed)

	removed = s.Remove("H0STCNT0", "005930")
	assert.False(t, removed) // 이미 없음
}

func TestSubscriber_RestoreAll(t *testing.T) {
	s := newSubscriber()
	s.Add("H0STCNT0", "005930")
	s.Add("H0STCNT0", "000660")
	s.Add("H0STASP0", "005930")

	all := s.All()
	assert.Len(t, all, 3)

	// 정렬되지 않은 list — set 비교
	keys := map[string]bool{}
	for _, sub := range all {
		keys[sub.TrID+":"+sub.TrKey] = true
	}
	assert.True(t, keys["H0STCNT0:005930"])
	assert.True(t, keys["H0STCNT0:000660"])
	assert.True(t, keys["H0STASP0:005930"])
}

func TestSubscriber_Concurrent(t *testing.T) {
	s := newSubscriber()
	done := make(chan struct{}, 100)
	for i := 0; i < 100; i++ {
		go func(i int) {
			s.Add("H0STCNT0", "00593" + string(rune('0'+i%10)))
			done <- struct{}{}
		}(i)
	}
	for i := 0; i < 100; i++ {
		<-done
	}
	// race detector 가 -race 모드에서 panic 없이 통과하면 OK
	assert.NotEmpty(t, s.All())
}
```

- [ ] **Step 2: 테스트 실패 확인**

```bash
go test ./websocket/ -run TestSubscriber -v
```

Expected: FAIL — `undefined: newSubscriber`.

- [ ] **Step 3: 구현**

`websocket/subscriber.go`:
```go
package websocket

import "sync"

// subKey 는 한 구독을 (tr_id, tr_key) 로 식별.
type subKey struct {
	TrID  string
	TrKey string
}

// subscriber 는 활성 구독을 thread-safe 하게 추적.
type subscriber struct {
	mu sync.RWMutex
	m  map[subKey]struct{}
}

func newSubscriber() *subscriber {
	return &subscriber{m: make(map[subKey]struct{})}
}

// Add 는 (tr_id, tr_key) 를 등록. 새로 추가했으면 true, 이미 존재했으면 false.
func (s *subscriber) Add(trID, trKey string) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	k := subKey{TrID: trID, TrKey: trKey}
	if _, exists := s.m[k]; exists {
		return false, nil
	}
	s.m[k] = struct{}{}
	return true, nil
}

// Remove 는 (tr_id, tr_key) 를 해제. 존재했으면 true.
func (s *subscriber) Remove(trID, trKey string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	k := subKey{TrID: trID, TrKey: trKey}
	if _, exists := s.m[k]; !exists {
		return false
	}
	delete(s.m, k)
	return true
}

// All 은 현재 활성 구독 list 를 snapshot 으로 반환 (reconnect 시 RestoreAll 용).
func (s *subscriber) All() []subKey {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]subKey, 0, len(s.m))
	for k := range s.m {
		out = append(out, k)
	}
	return out
}
```

- [ ] **Step 4: 테스트 통과 확인 (race 포함)**

```bash
go test -race ./websocket/ -run TestSubscriber -v
```

Expected: PASS — 4 테스트 모두.

- [ ] **Step 5: commit**

```bash
git add websocket/subscriber.go websocket/subscriber_test.go
git commit -m "[feat] Phase 8 — Subscriber (활성 구독 map, thread-safe)"
```

---

## Task 6: reconnect.go — ReconnectController

**Files:**
- Create: `websocket/reconnect.go`
- Test: `websocket/reconnect_test.go`

- [ ] **Step 1: 실패 테스트 작성**

`websocket/reconnect_test.go`:
```go
package websocket

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestReconnect_BackoffSequence(t *testing.T) {
	r := newReconnectController(reconnectOpts{
		Min: 1 * time.Second,
		Max: 30 * time.Second,
		MaxAttempts: 10,
	})

	// 1, 2, 4, 8, 16, 30, 30, 30, 30, 30 (cap)
	expected := []time.Duration{1, 2, 4, 8, 16, 30, 30, 30, 30, 30}
	for i, want := range expected {
		got, _ := r.NextBackoff()
		assert.Equalf(t, want*time.Second, got, "attempt %d", i+1)
	}
}

func TestReconnect_GiveUp(t *testing.T) {
	r := newReconnectController(reconnectOpts{
		Min: 1 * time.Second,
		Max: 30 * time.Second,
		MaxAttempts: 3,
	})

	r.NextBackoff()
	r.NextBackoff()
	r.NextBackoff()
	_, err := r.NextBackoff()
	assert.True(t, errors.Is(err, ErrWSGiveUp))
}

func TestReconnect_Reset(t *testing.T) {
	r := newReconnectController(reconnectOpts{
		Min: 1 * time.Second,
		Max: 30 * time.Second,
		MaxAttempts: 10,
	})
	r.NextBackoff()
	r.NextBackoff()
	r.Reset()
	got, _ := r.NextBackoff()
	assert.Equal(t, 1*time.Second, got)
}
```

- [ ] **Step 2: 테스트 실패 확인**

```bash
go test ./websocket/ -run TestReconnect -v
```

Expected: FAIL.

- [ ] **Step 3: 구현**

`websocket/reconnect.go`:
```go
package websocket

import "time"

type reconnectOpts struct {
	Min         time.Duration // 1s
	Max         time.Duration // 30s
	MaxAttempts int           // 10
}

type reconnectController struct {
	opts     reconnectOpts
	attempts int
}

func newReconnectController(opts reconnectOpts) *reconnectController {
	return &reconnectController{opts: opts}
}

// NextBackoff 는 다음 sleep duration 을 반환. attempts > MaxAttempts 면 ErrWSGiveUp.
func (r *reconnectController) NextBackoff() (time.Duration, error) {
	r.attempts++
	if r.attempts > r.opts.MaxAttempts {
		return 0, ErrWSGiveUp
	}
	// 1 * 2^(attempts-1), capped at Max
	d := r.opts.Min << (r.attempts - 1)
	if d > r.opts.Max {
		d = r.opts.Max
	}
	return d, nil
}

// Reset — 재연결 성공 시 호출.
func (r *reconnectController) Reset() {
	r.attempts = 0
}
```

- [ ] **Step 4: 테스트 통과 확인**

```bash
go test ./websocket/ -run TestReconnect -v
```

Expected: PASS — 3 테스트.

- [ ] **Step 5: commit**

```bash
git add websocket/reconnect.go websocket/reconnect_test.go
git commit -m "[feat] Phase 8 — ReconnectController (exp backoff, max attempts)"
```

---

## Task 7: frame.go — 프레임 파서 (첫 글자 분기)

**Files:**
- Create: `websocket/frame.go`
- Test: `websocket/frame_test.go`

- [ ] **Step 1: 실패 테스트 작성**

`websocket/frame_test.go`:
```go
package websocket

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseFrame_RealtimeData(t *testing.T) {
	raw := "0|H0STCNT0|001|005930^123929^73100^2"
	f, err := parseFrame(raw)
	require.NoError(t, err)
	assert.Equal(t, frameKindRealtime, f.Kind)
	assert.False(t, f.Encrypted)
	assert.Equal(t, "H0STCNT0", f.TrID)
	assert.Equal(t, 1, f.Count)
	assert.Equal(t, []string{"005930", "123929", "73100", "2"}, f.Fields)
}

func TestParseFrame_Encrypted(t *testing.T) {
	raw := "1|H0STCNI0|001|encrypted-payload"
	f, err := parseFrame(raw)
	require.NoError(t, err)
	assert.Equal(t, frameKindRealtime, f.Kind)
	assert.True(t, f.Encrypted)
}

func TestParseFrame_RealtimePaging(t *testing.T) {
	// count=2 인 페이징
	raw := "0|H0STCNT0|002|f1^f2^f3^f4^f1b^f2b^f3b^f4b"
	f, err := parseFrame(raw)
	require.NoError(t, err)
	assert.Equal(t, 2, f.Count)
	assert.Len(t, f.Fields, 8)
}

func TestParseFrame_JSON_SubscribeSuccess(t *testing.T) {
	raw := `{"header":{"tr_id":"H0STCNT0"},"body":{"rt_cd":"0","msg_cd":"OPSP0000","msg1":"SUBSCRIBE SUCCESS","output":{"iv":"abc","key":"def"}}}`
	f, err := parseFrame(raw)
	require.NoError(t, err)
	assert.Equal(t, frameKindJSON, f.Kind)
	assert.Equal(t, "0", f.JSON.RtCd)
	assert.Equal(t, "OPSP0000", f.JSON.MsgCd)
	assert.Contains(t, f.JSON.Msg1, "SUBSCRIBE")
}

func TestParseFrame_JSON_PingPong(t *testing.T) {
	raw := `{"header":{"tr_id":"PINGPONG"}}`
	f, err := parseFrame(raw)
	require.NoError(t, err)
	assert.Equal(t, frameKindPingPong, f.Kind)
}

func TestParseFrame_Invalid(t *testing.T) {
	_, err := parseFrame("garbage-data-no-pipe-no-brace")
	assert.True(t, errors.Is(err, ErrWSInvalidFrame))
}

func TestChunkFields(t *testing.T) {
	// 1 frame = 2 chunks (count=2), 4 fields/chunk
	chunks := chunkFields([]string{"a", "b", "c", "d", "e", "f", "g", "h"}, 2, 4)
	assert.Equal(t, [][]string{{"a", "b", "c", "d"}, {"e", "f", "g", "h"}}, chunks)
}

func TestChunkFields_MismatchedLength(t *testing.T) {
	// 7 fields, count=2, fieldsPerChunk=4 → mismatch
	_, err := chunkFieldsErr([]string{"a", "b", "c", "d", "e", "f", "g"}, 2, 4)
	assert.True(t, errors.Is(err, ErrWSInvalidFrame))
}
```

- [ ] **Step 2: 테스트 실패 확인**

```bash
go test ./websocket/ -run TestParseFrame -v
```

Expected: FAIL.

- [ ] **Step 3: 구현**

`websocket/frame.go`:
```go
package websocket

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type frameKind int

const (
	frameKindUnknown frameKind = iota
	frameKindRealtime
	frameKindJSON
	frameKindPingPong
)

// frame 은 파싱된 raw WebSocket 메시지.
type frame struct {
	Kind      frameKind
	Encrypted bool     // 0/1 flag (realtime 만 의미)
	TrID      string   // realtime 만
	Count     int      // realtime 만
	Fields    []string // realtime caret-separated payload
	JSON      jsonFrame
}

type jsonFrame struct {
	Header struct {
		TrID string `json:"tr_id"`
	} `json:"header"`
	Body struct {
		RtCd  string `json:"rt_cd"`
		MsgCd string `json:"msg_cd"`
		Msg1  string `json:"msg1"`
	} `json:"body"`
	// 호환성 위해 top-level 필드도 받아둠
	RtCd  string `json:"rt_cd,omitempty"`
	MsgCd string `json:"msg_cd,omitempty"`
	Msg1  string `json:"msg1,omitempty"`
}

// parseFrame 은 raw text 를 frame 으로 파싱. 첫 글자로 종류 분기.
func parseFrame(raw string) (frame, error) {
	if len(raw) == 0 {
		return frame{}, fmt.Errorf("%w: empty", ErrWSInvalidFrame)
	}
	switch raw[0] {
	case '{':
		return parseJSONFrame(raw)
	case '0', '1':
		return parseRealtimeFrame(raw)
	default:
		return frame{}, fmt.Errorf("%w: unknown leader %q", ErrWSInvalidFrame, raw[0])
	}
}

func parseRealtimeFrame(raw string) (frame, error) {
	parts := strings.SplitN(raw, "|", 4)
	if len(parts) != 4 {
		return frame{}, fmt.Errorf("%w: realtime frame requires 4 pipe-parts", ErrWSInvalidFrame)
	}
	encrypted := parts[0] == "1"
	count, err := strconv.Atoi(parts[2])
	if err != nil {
		return frame{}, fmt.Errorf("%w: bad count %q", ErrWSInvalidFrame, parts[2])
	}
	fields := strings.Split(parts[3], "^")
	return frame{
		Kind:      frameKindRealtime,
		Encrypted: encrypted,
		TrID:      parts[1],
		Count:     count,
		Fields:    fields,
	}, nil
}

func parseJSONFrame(raw string) (frame, error) {
	var jf jsonFrame
	if err := json.Unmarshal([]byte(raw), &jf); err != nil {
		return frame{}, fmt.Errorf("%w: json: %v", ErrWSInvalidFrame, err)
	}
	// PINGPONG 분기
	if jf.Header.TrID == "PINGPONG" {
		return frame{Kind: frameKindPingPong, JSON: jf}, nil
	}
	// body 우선, top-level fallback
	if jf.Body.RtCd == "" && jf.RtCd != "" {
		jf.Body.RtCd = jf.RtCd
		jf.Body.MsgCd = jf.MsgCd
		jf.Body.Msg1 = jf.Msg1
	}
	return frame{Kind: frameKindJSON, JSON: jf}, nil
}

// chunkFields 는 fields 를 count 개의 chunk 로 분리. 길이 mismatch 면 에러.
func chunkFields(fields []string, count, fieldsPerChunk int) [][]string {
	chunks, _ := chunkFieldsErr(fields, count, fieldsPerChunk)
	return chunks
}

func chunkFieldsErr(fields []string, count, fieldsPerChunk int) ([][]string, error) {
	if count*fieldsPerChunk != len(fields) {
		return nil, fmt.Errorf("%w: chunk mismatch: expect %d*%d=%d, got %d",
			ErrWSInvalidFrame, count, fieldsPerChunk, count*fieldsPerChunk, len(fields))
	}
	out := make([][]string, count)
	for i := 0; i < count; i++ {
		out[i] = fields[i*fieldsPerChunk : (i+1)*fieldsPerChunk]
	}
	return out, nil
}

// 보조: jsonFrame.SubscribeSuccess 헬퍼
func (j jsonFrame) IsSubscribeSuccess() bool {
	return j.Body.RtCd == "0" && strings.Contains(j.Body.Msg1, "SUBSCRIBE SUCCESS")
}
```

- [ ] **Step 4: 테스트 통과 확인**

```bash
go test ./websocket/ -run "TestParseFrame|TestChunk" -v
```

Expected: PASS — 8 테스트 (paging, encrypted, JSON, pingpong, invalid, chunk x2 등).

- [ ] **Step 5: commit**

```bash
git add websocket/frame.go websocket/frame_test.go
git commit -m "[feat] Phase 8 — frame parser (realtime/JSON/PINGPONG 분기) + chunkFields"
```

---

## Task 8: decode_krx.go — 5 EP decoder

**Files:**
- Create: `websocket/decode_krx.go`
- Test: `websocket/decode_krx_test.go`

- [ ] **Step 1: 실패 테스트 작성 (EP1 체결가만 작성, 나머지 4개 동일 패턴)**

`websocket/decode_krx_test.go`:
```go
package websocket

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func loadFixture(t *testing.T, name string) string {
	t.Helper()
	b, err := os.ReadFile(filepath.Join("testdata", name))
	require.NoError(t, err)
	return strings.TrimRight(string(b), "\n")
}

func TestDecodeKrxTrade_Single(t *testing.T) {
	raw := loadFixture(t, "h0stcnt0_success.txt")
	f, err := parseFrame(raw)
	require.NoError(t, err)
	require.Equal(t, "H0STCNT0", f.TrID)

	events, err := decodeKrxTrade(f)
	require.NoError(t, err)
	require.Len(t, events, 1)

	ev := events[0]
	assert.Equal(t, "005930", ev.Symbol)
	assert.Equal(t, "123929", ev.Time)
	assert.True(t, decimal.NewFromInt(73100).Equal(ev.Price))
	assert.Equal(t, "2", ev.PrevDiffSign)
	assert.Equal(t, int64(150), ev.TradeVolume)
	assert.NotEmpty(t, ev.Raw)
}

func TestDecodeKrxTrade_Paging(t *testing.T) {
	raw := loadFixture(t, "h0stcnt0_paging.txt") // count=2
	f, err := parseFrame(raw)
	require.NoError(t, err)

	events, err := decodeKrxTrade(f)
	require.NoError(t, err)
	require.Len(t, events, 2) // 페이징 2건
}

func TestDecodeKrxTrade_BadNumeric(t *testing.T) {
	// 일부러 숫자 필드를 "abc" 로 바꿔도 zero value, 에러 안 남.
	raw := strings.Replace(loadFixture(t, "h0stcnt0_success.txt"), "73100", "abc", 1)
	f, err := parseFrame(raw)
	require.NoError(t, err)

	events, err := decodeKrxTrade(f)
	require.NoError(t, err)
	require.Len(t, events, 1)
	assert.True(t, events[0].Price.IsZero()) // 파싱 실패 → zero
}

// 동일 패턴 4 개 더: TestDecodeKrxAsk / TestDecodeKrxExpectTrade /
// TestDecodeKrxOvernightTrade / TestDecodeKrxOvernightExpect
```

- [ ] **Step 2: 테스트 실패 확인**

```bash
go test ./websocket/ -run TestDecodeKrx -v
```

Expected: FAIL.

- [ ] **Step 3: 구현 (EP1 패턴, 나머지 동일)**

`websocket/decode_krx.go`:
```go
package websocket

import (
	"strconv"

	"github.com/shopspring/decimal"
)

// EP1 H0STCNT0 + EP4 H0STOAC0 (시간외) 동일 schema.
// fieldsPerChunk = 46 (Task 2 schema 결과 기반 — 실제 값 반영 필요)
const krxTradeFieldCount = 46

func decodeKrxTrade(f frame) ([]KrxTradeEvent, error) {
	chunks, err := chunkFieldsErr(f.Fields, f.Count, krxTradeFieldCount)
	if err != nil {
		return nil, err
	}
	out := make([]KrxTradeEvent, 0, len(chunks))
	for _, c := range chunks {
		out = append(out, parseKrxTradeChunk(c))
	}
	return out, nil
}

func parseKrxTradeChunk(c []string) KrxTradeEvent {
	// 정확한 필드 인덱스는 Task 2 schema 결과로 확정.
	// 아래는 docs 기반 예시 매핑.
	return KrxTradeEvent{
		Symbol:         c[0],
		Time:           c[1],
		Price:          asDecimal(c[2]),
		PrevDiffSign:   c[3],
		PrevDiff:       asDecimal(c[4]),
		PrevChangeRate: asFloat(c[5]),
		WeightedAvg:    asDecimal(c[6]),
		Open:           asDecimal(c[7]),
		High:           asDecimal(c[8]),
		Low:            asDecimal(c[9]),
		Ask1:           asDecimal(c[10]),
		Bid1:           asDecimal(c[11]),
		TradeVolume:    asInt64(c[12]),
		AccumVolume:    asInt64(c[13]),
		AccumValue:     asInt64(c[14]),
		// ... 나머지 필드 (Task 2 schema 기반)
		Raw: c,
	}
}

// asDecimal — 파싱 실패 시 zero. error 안 던짐.
func asDecimal(s string) decimal.Decimal {
	d, err := decimal.NewFromString(s)
	if err != nil {
		return decimal.Zero
	}
	return d
}

func asInt64(s string) int64 {
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0
	}
	return v
}

func asFloat(s string) float64 {
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return v
}

// === EP2 H0STASP0 (호가) ===

const krxAskFieldCount = 59 // Task 2 결과 반영

func decodeKrxAsk(f frame) ([]KrxAskEvent, error) {
	chunks, err := chunkFieldsErr(f.Fields, f.Count, krxAskFieldCount)
	if err != nil {
		return nil, err
	}
	out := make([]KrxAskEvent, 0, len(chunks))
	for _, c := range chunks {
		out = append(out, parseKrxAskChunk(c))
	}
	return out, nil
}

func parseKrxAskChunk(c []string) KrxAskEvent {
	ev := KrxAskEvent{
		Symbol: c[0],
		Time:   c[1],
		Hour:   c[2],
		Raw:    c,
	}
	// ASKP1..10 / BIDP1..10 / ASKP_RSQN1..10 / BIDP_RSQN1..10 — Task 2 schema 결과 인덱스
	for i := 0; i < 10; i++ {
		ev.Ask[i] = asDecimal(c[3+i])
		ev.Bid[i] = asDecimal(c[13+i])
		ev.AskSize[i] = asInt64(c[23+i])
		ev.BidSize[i] = asInt64(c[33+i])
	}
	return ev
}

// === EP3 H0STANC0 (예상체결) + EP5 H0STOAA0 (시간외 예상체결) ===

const krxExpectTradeFieldCount = 30 // Task 2 결과 반영

func decodeKrxExpectTrade(f frame) ([]KrxExpectTradeEvent, error) {
	chunks, err := chunkFieldsErr(f.Fields, f.Count, krxExpectTradeFieldCount)
	if err != nil {
		return nil, err
	}
	out := make([]KrxExpectTradeEvent, 0, len(chunks))
	for _, c := range chunks {
		out = append(out, parseKrxExpectTradeChunk(c))
	}
	return out, nil
}

func parseKrxExpectTradeChunk(c []string) KrxExpectTradeEvent {
	return KrxExpectTradeEvent{
		Symbol:      c[0],
		Time:        c[1],
		ExpectPrice: asDecimal(c[2]),
		// ... Task 2 결과 기반 나머지 필드
		Raw: c,
	}
}
```

- [ ] **Step 4: 테스트 통과 확인 (5 EP 모두)**

```bash
go test ./websocket/ -run TestDecodeKrx -v
```

Expected: PASS — 5 EP × (Single + Paging + BadNumeric) ≈ 15 테스트.

- [ ] **Step 5: commit**

```bash
git add websocket/decode_krx.go websocket/decode_krx_test.go
git commit -m "[feat] Phase 8 — decode_krx (5 EP, paging + bad numeric tolerance)"
```

---

## Task 9: dispatcher.go — TR_ID → handler 라우팅

**Files:**
- Create: `websocket/dispatcher.go`
- Test: `websocket/dispatcher_test.go`

- [ ] **Step 1: 실패 테스트 작성**

`websocket/dispatcher_test.go`:
```go
package websocket

import (
	"errors"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDispatcher_RouteKrxTrade(t *testing.T) {
	d := newDispatcher()
	var calls atomic.Int32
	d.OnKrxTrade(func(ev KrxTradeEvent) {
		calls.Add(1)
		assert.Equal(t, "005930", ev.Symbol)
	})

	d.RouteKrxTrade(KrxTradeEvent{Symbol: "005930"})
	assert.Equal(t, int32(1), calls.Load())
}

func TestDispatcher_HandlerPanic(t *testing.T) {
	d := newDispatcher()
	var errors []error
	d.OnError(func(err error) {
		errors = append(errors, err)
	})
	d.OnKrxTrade(func(ev KrxTradeEvent) {
		panic("boom")
	})

	d.RouteKrxTrade(KrxTradeEvent{Symbol: "005930"})
	assert.Len(t, errors, 1)
	assert.Contains(t, errors[0].Error(), "panic")
}

func TestDispatcher_NoHandler(t *testing.T) {
	d := newDispatcher()
	// handler 미등록 — silent ignore (panic 없어야 함)
	assert.NotPanics(t, func() {
		d.RouteKrxTrade(KrxTradeEvent{Symbol: "005930"})
	})
}

func TestDispatcher_RouteError(t *testing.T) {
	d := newDispatcher()
	var got error
	d.OnError(func(err error) { got = err })
	d.RouteError(ErrWSInvalidFrame)
	assert.True(t, errors.Is(got, ErrWSInvalidFrame))
}
```

- [ ] **Step 2: 테스트 실패 확인**

```bash
go test ./websocket/ -run TestDispatcher -v
```

Expected: FAIL.

- [ ] **Step 3: 구현**

`websocket/dispatcher.go`:
```go
package websocket

import (
	"fmt"
	"sync"
)

// dispatcher 는 TR_ID 별 handler 등록/라우팅.
type dispatcher struct {
	mu sync.RWMutex

	onKrxTrade           func(KrxTradeEvent)
	onKrxAsk             func(KrxAskEvent)
	onKrxExpectTrade     func(KrxExpectTradeEvent)
	onKrxOvernightTrade  func(KrxTradeEvent)
	onKrxOvernightExpect func(KrxExpectTradeEvent)

	onConnected  func()
	onReconnect  func(attempt int)
	onDisconnect func(error)
	onError      func(error)
}

func newDispatcher() *dispatcher { return &dispatcher{} }

// === 등록 메서드 ===

func (d *dispatcher) OnKrxTrade(h func(KrxTradeEvent))                  { d.mu.Lock(); d.onKrxTrade = h; d.mu.Unlock() }
func (d *dispatcher) OnKrxAsk(h func(KrxAskEvent))                       { d.mu.Lock(); d.onKrxAsk = h; d.mu.Unlock() }
func (d *dispatcher) OnKrxExpectTrade(h func(KrxExpectTradeEvent))       { d.mu.Lock(); d.onKrxExpectTrade = h; d.mu.Unlock() }
func (d *dispatcher) OnKrxOvernightTrade(h func(KrxTradeEvent))          { d.mu.Lock(); d.onKrxOvernightTrade = h; d.mu.Unlock() }
func (d *dispatcher) OnKrxOvernightExpect(h func(KrxExpectTradeEvent))   { d.mu.Lock(); d.onKrxOvernightExpect = h; d.mu.Unlock() }
func (d *dispatcher) OnConnected(h func())                                { d.mu.Lock(); d.onConnected = h; d.mu.Unlock() }
func (d *dispatcher) OnReconnect(h func(int))                             { d.mu.Lock(); d.onReconnect = h; d.mu.Unlock() }
func (d *dispatcher) OnDisconnect(h func(error))                          { d.mu.Lock(); d.onDisconnect = h; d.mu.Unlock() }
func (d *dispatcher) OnError(h func(error))                               { d.mu.Lock(); d.onError = h; d.mu.Unlock() }

// === 라우팅 메서드 ===

func (d *dispatcher) RouteKrxTrade(ev KrxTradeEvent)             { d.safeCall(func() { if d.onKrxTrade != nil { d.onKrxTrade(ev) } }) }
func (d *dispatcher) RouteKrxAsk(ev KrxAskEvent)                  { d.safeCall(func() { if d.onKrxAsk != nil { d.onKrxAsk(ev) } }) }
func (d *dispatcher) RouteKrxExpectTrade(ev KrxExpectTradeEvent)  { d.safeCall(func() { if d.onKrxExpectTrade != nil { d.onKrxExpectTrade(ev) } }) }
func (d *dispatcher) RouteKrxOvernightTrade(ev KrxTradeEvent)     { d.safeCall(func() { if d.onKrxOvernightTrade != nil { d.onKrxOvernightTrade(ev) } }) }
func (d *dispatcher) RouteKrxOvernightExpect(ev KrxExpectTradeEvent) { d.safeCall(func() { if d.onKrxOvernightExpect != nil { d.onKrxOvernightExpect(ev) } }) }

func (d *dispatcher) RouteConnected()         { d.safeCall(func() { if d.onConnected != nil { d.onConnected() } }) }
func (d *dispatcher) RouteReconnect(att int)  { d.safeCall(func() { if d.onReconnect != nil { d.onReconnect(att) } }) }
func (d *dispatcher) RouteDisconnect(e error) { d.safeCall(func() { if d.onDisconnect != nil { d.onDisconnect(e) } }) }
func (d *dispatcher) RouteError(e error)      { d.safeCall(func() { if d.onError != nil { d.onError(e) } }) }

// safeCall — handler panic 을 OnError 로 라우팅.
func (d *dispatcher) safeCall(fn func()) {
	defer func() {
		if r := recover(); r != nil {
			d.mu.RLock()
			h := d.onError
			d.mu.RUnlock()
			if h != nil {
				h(fmt.Errorf("kis ws: handler panic: %v", r))
			}
		}
	}()
	d.mu.RLock()
	defer d.mu.RUnlock()
	fn()
}
```

- [ ] **Step 4: 테스트 통과 확인**

```bash
go test -race ./websocket/ -run TestDispatcher -v
```

Expected: PASS — 4 테스트.

- [ ] **Step 5: commit**

```bash
git add websocket/dispatcher.go websocket/dispatcher_test.go
git commit -m "[feat] Phase 8 — Dispatcher (TR_ID 라우팅 + recover)"
```

---

## Task 10: approval.go — ApprovalKeyManager

**Files:**
- Create: `websocket/approval.go`
- Test: `websocket/approval_test.go`

- [ ] **Step 1: 실패 테스트 작성**

`websocket/approval_test.go`:
```go
package websocket

import (
	"context"
	"net/http"
	"sync/atomic"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestApproval_FetchAndCache(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var calls atomic.Int32
	httpmock.RegisterResponder(http.MethodPost, `=~/oauth2/Approval`,
		func(req *http.Request) (*http.Response, error) {
			calls.Add(1)
			return httpmock.NewStringResponse(200, `{"approval_key":"key-12345"}`), nil
		},
	)

	m := newApprovalKeyManager(http.DefaultClient, "https://api.example", "appkey", "appsecret", 23*time.Hour)

	k1, err := m.Get(context.Background())
	require.NoError(t, err)
	assert.Equal(t, "key-12345", k1)

	// 캐시 동작 — 두 번째 호출은 HTTP 안 침
	k2, err := m.Get(context.Background())
	require.NoError(t, err)
	assert.Equal(t, "key-12345", k2)
	assert.Equal(t, int32(1), calls.Load())
}

func TestApproval_Expiry(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	keys := []string{"key-A", "key-B"}
	idx := 0
	httpmock.RegisterResponder(http.MethodPost, `=~/oauth2/Approval`,
		func(req *http.Request) (*http.Response, error) {
			defer func() { idx++ }()
			return httpmock.NewStringResponse(200, `{"approval_key":"`+keys[idx]+`"}`), nil
		},
	)

	// TTL 0 → 매번 갱신
	m := newApprovalKeyManager(http.DefaultClient, "https://api.example", "appkey", "appsecret", 0)
	k1, _ := m.Get(context.Background())
	k2, _ := m.Get(context.Background())
	assert.NotEqual(t, k1, k2)
}

func TestApproval_Failure(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(http.MethodPost, `=~/oauth2/Approval`,
		httpmock.NewStringResponder(500, `{"error":"server"}`),
	)

	m := newApprovalKeyManager(http.DefaultClient, "https://api.example", "appkey", "appsecret", 23*time.Hour)
	_, err := m.Get(context.Background())
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrWSApprovalFailed)
}
```

- [ ] **Step 2: 테스트 실패 확인**

```bash
go test ./websocket/ -run TestApproval -v
```

Expected: FAIL.

- [ ] **Step 3: 구현**

`websocket/approval.go`:
```go
package websocket

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

type approvalKeyManager struct {
	httpClient *http.Client
	baseURL    string
	appKey     string
	appSecret  string
	ttl        time.Duration

	mu     sync.Mutex
	cached string
	expiry time.Time
}

func newApprovalKeyManager(c *http.Client, baseURL, appKey, appSecret string, ttl time.Duration) *approvalKeyManager {
	if c == nil {
		c = http.DefaultClient
	}
	return &approvalKeyManager{
		httpClient: c,
		baseURL:    baseURL,
		appKey:     appKey,
		appSecret:  appSecret,
		ttl:        ttl,
	}
}

func (m *approvalKeyManager) Get(ctx context.Context) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.cached != "" && time.Now().Before(m.expiry) {
		return m.cached, nil
	}

	body, _ := json.Marshal(map[string]string{
		"grant_type": "client_credentials",
		"appkey":     m.appKey,
		"secretkey":  m.appSecret,
	})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, m.baseURL+"/oauth2/Approval", bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("%w: build req: %v", ErrWSApprovalFailed, err)
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	resp, err := m.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("%w: http: %v", ErrWSApprovalFailed, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		raw, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("%w: status %d body %s", ErrWSApprovalFailed, resp.StatusCode, string(raw))
	}

	var out struct {
		ApprovalKey string `json:"approval_key"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", fmt.Errorf("%w: decode: %v", ErrWSApprovalFailed, err)
	}
	if out.ApprovalKey == "" {
		return "", fmt.Errorf("%w: empty approval_key", ErrWSApprovalFailed)
	}

	m.cached = out.ApprovalKey
	m.expiry = time.Now().Add(m.ttl)
	return m.cached, nil
}
```

- [ ] **Step 4: 테스트 통과 확인**

```bash
go test ./websocket/ -run TestApproval -v
```

Expected: PASS — 3 테스트.

- [ ] **Step 5: commit**

```bash
git add websocket/approval.go websocket/approval_test.go
git commit -m "[feat] Phase 8 — ApprovalKeyManager (POST /oauth2/Approval, TTL 캐시)"
```

---

## Task 11: conn.go — ConnectionManager

**Files:**
- Create: `websocket/conn.go`

> 이 task 는 `coder/websocket` wrapper. unit test 어렵고 (실제 WS 필요) integration test (Task 14) 에서 검증.

- [ ] **Step 1: 구현**

`websocket/conn.go`:
```go
package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/coder/websocket"
)

type connManager struct {
	mu       sync.Mutex
	endpoint string
	conn     *websocket.Conn
}

func newConnManager(endpoint string) *connManager {
	return &connManager{endpoint: endpoint}
}

// Dial 은 WebSocket 연결.
func (cm *connManager) Dial(ctx context.Context) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	if cm.conn != nil {
		_ = cm.conn.Close(websocket.StatusNormalClosure, "redial")
		cm.conn = nil
	}
	c, _, err := websocket.Dial(ctx, cm.endpoint, &websocket.DialOptions{
		Subprotocols: nil,
	})
	if err != nil {
		return fmt.Errorf("kis ws: dial %s: %w", cm.endpoint, err)
	}
	c.SetReadLimit(1 << 20) // 1 MiB
	cm.conn = c
	return nil
}

// SendSubscribe 는 subscribe/unsubscribe frame 송신.
func (cm *connManager) SendSubscribe(ctx context.Context, approvalKey, custType, trType, trID, trKey string) error {
	cm.mu.Lock()
	c := cm.conn
	cm.mu.Unlock()
	if c == nil {
		return ErrWSNotConnected
	}
	msg := map[string]any{
		"header": map[string]string{
			"approval_key": approvalKey,
			"custtype":     custType,
			"tr_type":      trType,
			"content-type": "utf-8",
		},
		"body": map[string]any{
			"input": map[string]string{
				"tr_id":  trID,
				"tr_key": trKey,
			},
		},
	}
	raw, _ := json.Marshal(msg)
	return c.Write(ctx, websocket.MessageText, raw)
}

// Read 는 다음 text frame 을 받음 (blocking).
func (cm *connManager) Read(ctx context.Context) (string, error) {
	cm.mu.Lock()
	c := cm.conn
	cm.mu.Unlock()
	if c == nil {
		return "", ErrWSNotConnected
	}
	_, raw, err := c.Read(ctx)
	if err != nil {
		return "", err
	}
	return string(raw), nil
}

// Pong — PINGPONG 응답 (KIS 가 보낸 PING 메시지 그대로 echo).
func (cm *connManager) Pong(ctx context.Context, raw string) error {
	cm.mu.Lock()
	c := cm.conn
	cm.mu.Unlock()
	if c == nil {
		return ErrWSNotConnected
	}
	return c.Write(ctx, websocket.MessageText, []byte(raw))
}

// Close — graceful close.
func (cm *connManager) Close() error {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	if cm.conn == nil {
		return nil
	}
	err := cm.conn.Close(websocket.StatusNormalClosure, "client shutdown")
	cm.conn = nil
	return err
}
```

- [ ] **Step 2: build 검증**

```bash
go build ./websocket/
```

Expected: 출력 없음.

- [ ] **Step 3: vet 검증**

```bash
go vet ./websocket/
```

Expected: 출력 없음.

- [ ] **Step 4: 단위 테스트 생략 — Task 14 integration tests 에서 wsmock 으로 e2e 검증**

이 task 는 wrapper layer 이므로 실제 ws 가 없으면 의미있는 test 어려움. Task 14 가 wsmock 으로 검증.

- [ ] **Step 5: commit**

```bash
git add websocket/conn.go
git commit -m "[feat] Phase 8 — ConnectionManager (coder/websocket wrapper)"
```

---

## Task 12: client.go — public Client + Run + Subscribe + Handler

**Files:**
- Create: `websocket/client.go`

> 이 task 는 모든 모듈을 결합하는 facade layer. 단위 test 는 dispatcher/subscriber/reconnect 가 이미 cover. Task 14 integration test 가 e2e 검증.

- [ ] **Step 1: Options 정의**

`websocket/client.go`:
```go
package websocket

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"sync"
	"time"
)

// Options 는 WS Client 생성 옵션.
type Options struct {
	Endpoint      string        // ws://ops.koreainvestment.com:21000 (실전) / :31000 (모의)
	BaseURL       string        // /oauth2/Approval 호출 base URL (예: https://openapi.koreainvestment.com:9443)
	AppKey        string
	AppSecret     string
	CustType      string        // "P" (개인, default) / "B" (법인)
	MaxReconnects int           // default 10
	ReconnectMin  time.Duration // default 1s
	ReconnectMax  time.Duration // default 30s
	ApprovalTTL   time.Duration // default 23h
	HTTPClient    *http.Client  // default http.DefaultClient
	Logger        *slog.Logger  // default discard
}

func (o *Options) defaults() {
	if o.MaxReconnects == 0 {
		o.MaxReconnects = 10
	}
	if o.ReconnectMin == 0 {
		o.ReconnectMin = 1 * time.Second
	}
	if o.ReconnectMax == 0 {
		o.ReconnectMax = 30 * time.Second
	}
	if o.ApprovalTTL == 0 {
		o.ApprovalTTL = 23 * time.Hour
	}
	if o.CustType == "" {
		o.CustType = "P"
	}
	if o.Logger == nil {
		o.Logger = slog.New(slog.NewTextHandler(io.Discard, nil))
	}
}
```

- [ ] **Step 2: Client struct + NewClient**

```go
// Client 는 KIS WebSocket 진입점. kis.Client.WS 로 접근.
type Client struct {
	opts Options

	approval   *approvalKeyManager
	conn       *connManager
	sub        *subscriber
	dispatcher *dispatcher
	reconnect  *reconnectController

	mu        sync.Mutex
	connected bool // 현재 dial 된 상태 — Subscribe 시 즉시 송신 vs 보류 결정
}

func NewClient(opts Options) *Client {
	opts.defaults()
	return &Client{
		opts:     opts,
		approval: newApprovalKeyManager(opts.HTTPClient, opts.BaseURL, opts.AppKey, opts.AppSecret, opts.ApprovalTTL),
		conn:     newConnManager(opts.Endpoint),
		sub:      newSubscriber(),
		dispatcher: newDispatcher(),
		reconnect: newReconnectController(reconnectOpts{
			Min:         opts.ReconnectMin,
			Max:         opts.ReconnectMax,
			MaxAttempts: opts.MaxReconnects,
		}),
	}
}
```

- [ ] **Step 3: Subscribe / Unsubscribe / Handler 위임**

```go
// === Subscribe (5 EP) ===

const (
	trIDKrxTrade           = "H0STCNT0"
	trIDKrxAsk             = "H0STASP0"
	trIDKrxExpectTrade     = "H0STANC0"
	trIDKrxOvernightTrade  = "H0STOAC0"
	trIDKrxOvernightExpect = "H0STOAA0"
)

func (c *Client) SubscribeKrxTrade(symbols ...string) error           { return c.subscribe(trIDKrxTrade, symbols) }
func (c *Client) SubscribeKrxAsk(symbols ...string) error             { return c.subscribe(trIDKrxAsk, symbols) }
func (c *Client) SubscribeKrxExpectTrade(symbols ...string) error     { return c.subscribe(trIDKrxExpectTrade, symbols) }
func (c *Client) SubscribeKrxOvernightTrade(symbols ...string) error  { return c.subscribe(trIDKrxOvernightTrade, symbols) }
func (c *Client) SubscribeKrxOvernightExpect(symbols ...string) error { return c.subscribe(trIDKrxOvernightExpect, symbols) }

func (c *Client) UnsubscribeKrxTrade(symbols ...string) error           { return c.unsubscribe(trIDKrxTrade, symbols) }
// ... 나머지 4개 동일 패턴

func (c *Client) subscribe(trID string, symbols []string) error {
	for _, sym := range symbols {
		added, err := c.sub.Add(trID, sym)
		if err != nil {
			return err
		}
		c.mu.Lock()
		conn := c.connected
		c.mu.Unlock()
		if added && conn {
			ak, err := c.approval.Get(context.Background())
			if err != nil {
				return err
			}
			if err := c.conn.SendSubscribe(context.Background(), ak, c.opts.CustType, "1", trID, sym); err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *Client) unsubscribe(trID string, symbols []string) error {
	for _, sym := range symbols {
		removed := c.sub.Remove(trID, sym)
		c.mu.Lock()
		conn := c.connected
		c.mu.Unlock()
		if removed && conn {
			ak, err := c.approval.Get(context.Background())
			if err != nil {
				return err
			}
			if err := c.conn.SendSubscribe(context.Background(), ak, c.opts.CustType, "2", trID, sym); err != nil {
				return err
			}
		}
	}
	return nil
}

// === Handler 위임 ===

func (c *Client) OnKrxTrade(h func(KrxTradeEvent))                  { c.dispatcher.OnKrxTrade(h) }
func (c *Client) OnKrxAsk(h func(KrxAskEvent))                       { c.dispatcher.OnKrxAsk(h) }
func (c *Client) OnKrxExpectTrade(h func(KrxExpectTradeEvent))       { c.dispatcher.OnKrxExpectTrade(h) }
func (c *Client) OnKrxOvernightTrade(h func(KrxTradeEvent))          { c.dispatcher.OnKrxOvernightTrade(h) }
func (c *Client) OnKrxOvernightExpect(h func(KrxExpectTradeEvent))   { c.dispatcher.OnKrxOvernightExpect(h) }
func (c *Client) OnConnected(h func())                                { c.dispatcher.OnConnected(h) }
func (c *Client) OnReconnect(h func(attempt int))                     { c.dispatcher.OnReconnect(h) }
func (c *Client) OnDisconnect(h func(error))                          { c.dispatcher.OnDisconnect(h) }
func (c *Client) OnError(h func(error))                               { c.dispatcher.OnError(h) }
```

- [ ] **Step 4: Run loop**

```go
// Run 은 blocking — ctx 끝나면 graceful close.
// 자동 재연결 + 구독 자동 복원 포함.
func (c *Client) Run(ctx context.Context) error {
	for {
		// dial
		if err := c.dial(ctx); err != nil {
			d, gErr := c.reconnect.NextBackoff()
			if errors.Is(gErr, ErrWSGiveUp) {
				c.dispatcher.RouteError(ErrWSGiveUp)
				return ErrWSGiveUp
			}
			c.opts.Logger.Warn("ws dial failed", "err", err, "backoff", d)
			select {
			case <-time.After(d):
			case <-ctx.Done():
				return ctx.Err()
			}
			continue
		}

		// 연결 성공 — 기존 구독 복원
		if err := c.restoreSubs(ctx); err != nil {
			c.opts.Logger.Warn("ws restore subs failed", "err", err)
		}
		c.dispatcher.RouteConnected()
		c.reconnect.Reset()

		// 메시지 read loop
		err := c.readLoop(ctx)
		if errors.Is(err, ctx.Err()) {
			_ = c.conn.Close()
			return ctx.Err()
		}
		c.dispatcher.RouteDisconnect(err)
		c.mu.Lock()
		c.connected = false
		c.mu.Unlock()

		// 재연결 backoff
		d, gErr := c.reconnect.NextBackoff()
		if errors.Is(gErr, ErrWSGiveUp) {
			c.dispatcher.RouteError(ErrWSGiveUp)
			return ErrWSGiveUp
		}
		c.dispatcher.RouteReconnect(c.reconnect.attempts)
		select {
		case <-time.After(d):
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (c *Client) dial(ctx context.Context) error {
	if err := c.conn.Dial(ctx); err != nil {
		return err
	}
	c.mu.Lock()
	c.connected = true
	c.mu.Unlock()
	return nil
}

func (c *Client) restoreSubs(ctx context.Context) error {
	subs := c.sub.All()
	if len(subs) == 0 {
		return nil
	}
	ak, err := c.approval.Get(ctx)
	if err != nil {
		return err
	}
	for _, k := range subs {
		if err := c.conn.SendSubscribe(ctx, ak, c.opts.CustType, "1", k.TrID, k.TrKey); err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) readLoop(ctx context.Context) error {
	for {
		raw, err := c.conn.Read(ctx)
		if err != nil {
			return err
		}
		f, perr := parseFrame(raw)
		if perr != nil {
			c.dispatcher.RouteError(perr)
			continue
		}
		c.handleFrame(ctx, raw, f)
	}
}

func (c *Client) handleFrame(ctx context.Context, raw string, f frame) {
	switch f.Kind {
	case frameKindPingPong:
		_ = c.conn.Pong(ctx, raw)
	case frameKindJSON:
		if f.JSON.Body.RtCd != "0" {
			c.dispatcher.RouteError(&WSServerError{
				TrID:  f.JSON.Header.TrID,
				MsgCd: f.JSON.Body.MsgCd,
				Msg:   f.JSON.Body.Msg1,
			})
		}
		// 등록 성공 응답 등은 silent (logger 만)
		c.opts.Logger.Debug("ws json frame", "tr_id", f.JSON.Header.TrID, "msg_cd", f.JSON.Body.MsgCd)
	case frameKindRealtime:
		if f.Encrypted {
			c.dispatcher.RouteError(ErrWSEncryptedNotSupported)
			return
		}
		c.routeRealtime(f)
	}
}

func (c *Client) routeRealtime(f frame) {
	switch f.TrID {
	case trIDKrxTrade:
		evs, err := decodeKrxTrade(f)
		if err != nil { c.dispatcher.RouteError(err); return }
		for _, ev := range evs { c.dispatcher.RouteKrxTrade(ev) }
	case trIDKrxAsk:
		evs, err := decodeKrxAsk(f)
		if err != nil { c.dispatcher.RouteError(err); return }
		for _, ev := range evs { c.dispatcher.RouteKrxAsk(ev) }
	case trIDKrxExpectTrade:
		evs, err := decodeKrxExpectTrade(f)
		if err != nil { c.dispatcher.RouteError(err); return }
		for _, ev := range evs { c.dispatcher.RouteKrxExpectTrade(ev) }
	case trIDKrxOvernightTrade:
		evs, err := decodeKrxTrade(f) // 동일 schema
		if err != nil { c.dispatcher.RouteError(err); return }
		for _, ev := range evs { c.dispatcher.RouteKrxOvernightTrade(ev) }
	case trIDKrxOvernightExpect:
		evs, err := decodeKrxExpectTrade(f) // 동일 schema
		if err != nil { c.dispatcher.RouteError(err); return }
		for _, ev := range evs { c.dispatcher.RouteKrxOvernightExpect(ev) }
	default:
		c.dispatcher.RouteError(ErrWSInvalidFrame)
	}
}
```

- [ ] **Step 5: build + commit**

```bash
go build ./...
git add websocket/client.go
git commit -m "[feat] Phase 8 — public Client (Run + Subscribe + Handler 위임)"
```

Expected build: 출력 없음.

---

## Task 13: kis.Client integration

**Files:**
- Modify: `client.go` (root)

- [ ] **Step 1: 기존 client.go 에 WS 필드 + 초기화 추가**

`client.go` 의 변경 (modify, 핵심 부분):
```go
import (
	// ... 기존 imports
	"github.com/kenshin579/korea-investment-stock/websocket"
)

type Client struct {
	// ... 기존 필드
	Domestic *domestic.Client
	Overseas *overseas.Client
	Bonds    *bonds.Client
	WS       *websocket.Client  // 신규
	// ...
}

// NewClient (또는 NewClientFromEnv 등) 내부:
func NewClient(...) (*Client, error) {
	// ... 기존 로직
	c.Domestic = domestic.New(c.httpClient, c.masterC)
	c.Overseas = overseas.New(c.httpClient, c.masterC)
	c.Bonds = bonds.New(c.httpClient)

	// WebSocket endpoint 자동 결정 (real vs paper)
	wsEndpoint := "ws://ops.koreainvestment.com:21000"
	if c.opts.env == PaperEnv {
		wsEndpoint = "ws://ops.koreainvestment.com:31000"
	}
	c.WS = websocket.NewClient(websocket.Options{
		Endpoint:  wsEndpoint,
		BaseURL:   c.opts.env, // 실전/모의 base URL
		AppKey:    c.apiKey,
		AppSecret: c.apiSecret,
	})

	return c, nil
}
```

> 정확한 변경 지점은 기존 `NewClient*` 코드를 read 해서 결정. `c.opts.env` 정확한 필드명은 기존 코드 패턴 따라감.

- [ ] **Step 2: 기존 모든 테스트 PASS 검증**

```bash
go test ./...
```

Expected: 모든 패키지 PASS (기존 121 메서드 회귀 없음).

- [ ] **Step 3: 새 사용 시나리오 manual import 검증**

```go
// 임시 _scratch.go 작성:
import kis "github.com/kenshin579/korea-investment-stock"
client, _ := kis.NewClientFromEnv()
client.WS.SubscribeKrxTrade("005930") // 컴파일만 검증
```

```bash
go build ./...
```

Expected: 출력 없음. 임시 파일 삭제.

- [ ] **Step 4: vet**

```bash
go vet ./...
```

Expected: 출력 없음.

- [ ] **Step 5: commit**

```bash
git add client.go
git commit -m "[feat] Phase 8 — kis.Client.WS 필드 추가 (websocket.Client 주입)"
```

---

## Task 14: wsmock — local KIS WebSocket 서버

**Files:**
- Create: `websocket/internal/wsmock/server.go`

- [ ] **Step 1: wsmock server 구현**

`websocket/internal/wsmock/server.go`:
```go
// Package wsmock 은 KIS WebSocket 서버를 모방하는 local mock.
//
// 사용:
//
//	srv := wsmock.New(t)
//	defer srv.Close()
//	// srv.URL() 을 websocket.Options.Endpoint 로 사용
//
//	srv.SendRealtime("H0STCNT0", "005930^123929^73100^...")
//	srv.SendJSON(`{"header":{"tr_id":"H0STCNT0"},"body":{"rt_cd":"0","msg_cd":"OPSP0000","msg1":"SUBSCRIBE SUCCESS"}}`)
package wsmock

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/coder/websocket"
)

type Server struct {
	t       *testing.T
	hs      *httptest.Server
	mu      sync.Mutex
	conn    *websocket.Conn
	connCh  chan struct{} // 클라이언트 connect 알림
	receive chan string    // 클라이언트 → mock 송신 메시지
}

func New(t *testing.T) *Server {
	t.Helper()
	s := &Server{
		t:       t,
		connCh:  make(chan struct{}, 1),
		receive: make(chan string, 100),
	}
	s.hs = httptest.NewServer(http.HandlerFunc(s.handle))
	return s
}

func (s *Server) URL() string {
	// httptest 는 http:// 반환 — ws:// 로 변환
	u := s.hs.URL
	return "ws" + u[len("http"):]
}

func (s *Server) Close() {
	s.mu.Lock()
	if s.conn != nil {
		_ = s.conn.Close(websocket.StatusNormalClosure, "test end")
	}
	s.mu.Unlock()
	s.hs.Close()
}

func (s *Server) handle(w http.ResponseWriter, r *http.Request) {
	c, err := websocket.Accept(w, r, nil)
	if err != nil {
		s.t.Logf("wsmock accept: %v", err)
		return
	}
	s.mu.Lock()
	s.conn = c
	s.mu.Unlock()
	select {
	case s.connCh <- struct{}{}:
	default:
	}
	defer c.Close(websocket.StatusNormalClosure, "")

	// 클라이언트 → mock 메시지 read loop
	for {
		_, raw, err := c.Read(r.Context())
		if err != nil {
			return
		}
		select {
		case s.receive <- string(raw):
		default:
		}
	}
}

// WaitConnected 는 클라이언트 connect 까지 blocking.
func (s *Server) WaitConnected(ctx context.Context) error {
	select {
	case <-s.connCh:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// SendText 는 mock → 클라이언트 raw text frame 송신.
func (s *Server) SendText(ctx context.Context, msg string) error {
	s.mu.Lock()
	c := s.conn
	s.mu.Unlock()
	return c.Write(ctx, websocket.MessageText, []byte(msg))
}

// SendRealtime — 0|TR_ID|001|payload 형식 송신 helper.
func (s *Server) SendRealtime(ctx context.Context, trID, payload string) error {
	return s.SendText(ctx, "0|"+trID+"|001|"+payload)
}

// CloseConn — 클라이언트 측에서는 abnormal close 로 보임 (reconnect 시나리오 테스트용).
func (s *Server) CloseConn() {
	s.mu.Lock()
	c := s.conn
	s.conn = nil
	s.mu.Unlock()
	if c != nil {
		_ = c.Close(websocket.StatusAbnormalClosure, "test forced close")
	}
}

// Received 는 클라이언트 → mock 으로 송신된 메시지 channel.
func (s *Server) Received() <-chan string { return s.receive }
```

- [ ] **Step 2: build 검증**

```bash
go build ./websocket/internal/wsmock/
```

Expected: 출력 없음.

- [ ] **Step 3: vet**

```bash
go vet ./websocket/internal/wsmock/
```

Expected: 출력 없음.

- [ ] **Step 4: 단독 unit test 생략 — Task 15 가 사용자**

- [ ] **Step 5: commit**

```bash
git add websocket/internal/wsmock/
git commit -m "[chore] Phase 8 — wsmock local KIS WebSocket server"
```

---

## Task 15: integration tests (e2e via wsmock)

**Files:**
- Create: `websocket/integration_test.go`
- Create: `websocket/testdata/approval_success.json`

- [ ] **Step 1: approval mock fixture**

`websocket/testdata/approval_success.json`:
```json
{"approval_key":"test-approval-key-123"}
```

- [ ] **Step 2: integration test 작성**

`websocket/integration_test.go`:
```go
package websocket_test

import (
	"context"
	"net/http"
	"sync/atomic"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kenshin579/korea-investment-stock/websocket"
	"github.com/kenshin579/korea-investment-stock/websocket/internal/wsmock"
)

func setupApprovalMock(t *testing.T) {
	t.Helper()
	httpmock.RegisterResponder(http.MethodPost, `=~/oauth2/Approval`,
		httpmock.NewStringResponder(200, `{"approval_key":"test-approval-key-123"}`),
	)
}

func newClient(t *testing.T, endpoint string) *websocket.Client {
	t.Helper()
	return websocket.NewClient(websocket.Options{
		Endpoint:      endpoint,
		BaseURL:       "https://api.example",
		AppKey:        "appkey",
		AppSecret:     "appsecret",
		ReconnectMin:  10 * time.Millisecond,
		ReconnectMax:  100 * time.Millisecond,
		MaxReconnects: 5,
	})
}

func TestIntegration_HappyPath(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	setupApprovalMock(t)

	srv := wsmock.New(t)
	defer srv.Close()

	c := newClient(t, srv.URL())

	var received atomic.Int32
	c.OnKrxTrade(func(ev websocket.KrxTradeEvent) {
		received.Add(1)
		assert.Equal(t, "005930", ev.Symbol)
	})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	go c.Run(ctx)
	require.NoError(t, srv.WaitConnected(ctx))

	require.NoError(t, c.SubscribeKrxTrade("005930"))

	// 클라이언트가 보낸 subscribe frame 확인
	select {
	case msg := <-srv.Received():
		assert.Contains(t, msg, "H0STCNT0")
		assert.Contains(t, msg, "005930")
	case <-ctx.Done():
		t.Fatal("did not receive subscribe frame")
	}

	// mock 가 realtime frame 송신
	payload := samplePayload46Fields("005930") // helper, 46 fields caret-separated
	require.NoError(t, srv.SendRealtime(ctx, "H0STCNT0", payload))

	// handler 호출 대기
	require.Eventually(t, func() bool { return received.Load() > 0 }, 1*time.Second, 10*time.Millisecond)
}

func TestIntegration_Reconnect(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	setupApprovalMock(t)

	srv := wsmock.New(t)
	defer srv.Close()

	c := newClient(t, srv.URL())

	var reconnects atomic.Int32
	c.OnReconnect(func(att int) {
		reconnects.Add(1)
	})

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	go c.Run(ctx)
	require.NoError(t, srv.WaitConnected(ctx))
	require.NoError(t, c.SubscribeKrxTrade("005930"))

	// 첫 subscribe drain
	<-srv.Received()

	// mock 측에서 강제 close → SDK 재연결
	srv.CloseConn()
	require.NoError(t, srv.WaitConnected(ctx))

	// 재연결 후 기존 구독 자동 복원 frame 검증
	select {
	case msg := <-srv.Received():
		assert.Contains(t, msg, "H0STCNT0")
		assert.Contains(t, msg, "005930")
	case <-ctx.Done():
		t.Fatal("did not receive resubscribe frame")
	}

	require.Eventually(t, func() bool { return reconnects.Load() > 0 }, 2*time.Second, 10*time.Millisecond)
}

func TestIntegration_ServerError(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	setupApprovalMock(t)

	srv := wsmock.New(t)
	defer srv.Close()

	c := newClient(t, srv.URL())

	var got error
	c.OnError(func(err error) { got = err })

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	go c.Run(ctx)
	require.NoError(t, srv.WaitConnected(ctx))

	require.NoError(t, srv.SendText(ctx, `{"header":{"tr_id":"H0STCNT0"},"body":{"rt_cd":"1","msg_cd":"OPSP0001","msg1":"ALREADY IN SUBSCRIBE"}}`))

	require.Eventually(t, func() bool { return got != nil }, 1*time.Second, 10*time.Millisecond)
	wsErr, ok := got.(*websocket.WSServerError)
	require.True(t, ok)
	assert.Equal(t, "OPSP0001", wsErr.MsgCd)
}

func TestIntegration_PingPong(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	setupApprovalMock(t)

	srv := wsmock.New(t)
	defer srv.Close()

	c := newClient(t, srv.URL())

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	go c.Run(ctx)
	require.NoError(t, srv.WaitConnected(ctx))

	pingMsg := `{"header":{"tr_id":"PINGPONG"}}`
	require.NoError(t, srv.SendText(ctx, pingMsg))

	// 클라이언트가 echo 응답 — Received 에 동일 메시지
	select {
	case msg := <-srv.Received():
		assert.Equal(t, pingMsg, msg)
	case <-time.After(1 * time.Second):
		t.Fatal("did not receive PONG echo")
	}
}

func TestIntegration_GracefulShutdown(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	setupApprovalMock(t)

	srv := wsmock.New(t)
	defer srv.Close()

	c := newClient(t, srv.URL())

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- c.Run(ctx) }()

	require.NoError(t, srv.WaitConnected(context.Background()))

	cancel()
	select {
	case err := <-done:
		assert.ErrorIs(t, err, context.Canceled)
	case <-time.After(2 * time.Second):
		t.Fatal("Run did not return after ctx cancel")
	}
}

// helper — 46 필드 caret payload (실 schema 와 호환되는 minimal 값)
func samplePayload46Fields(symbol string) string {
	return symbol + "^123929^73100^2^1500^2.09^72850^72500^73200^72400^73100^73000^150^123456^987654000000^0^0^0^0^0^0^1^0^0^0^0^0^0^0^0^0^0^0^0^0^0^0^0^0^0^0^0^0^0^0^0"
}
```

- [ ] **Step 3: 테스트 실행 (race 포함)**

```bash
go test -race ./websocket/ -run TestIntegration -v
```

Expected: PASS — 5 시나리오 (Happy / Reconnect / ServerError / PingPong / Shutdown).

- [ ] **Step 4: coverage 확인**

```bash
go test -cover ./websocket/
```

Expected: ≥70% (network/timing 코드 한계 고려).

- [ ] **Step 5: commit**

```bash
git add websocket/integration_test.go websocket/testdata/approval_success.json
git commit -m "[test] Phase 8 — integration tests (wsmock e2e: happy/reconnect/error/pingpong/shutdown)"
```

---

## Task 16: example + docs + PR

**Files:**
- Create: `examples/ws_krx_basic/main.go`
- Modify: `README.md`, `CLAUDE.md`, `CHANGELOG.md`

- [ ] **Step 1: example 작성**

`examples/ws_krx_basic/main.go`:
```go
// examples/ws_krx_basic/main.go — Phase 8 KRX WebSocket 시세 시연.
//
// EP1: SubscribeKrxTrade — H0STCNT0 (체결가, KRX)
// EP2: SubscribeKrxAsk — H0STASP0 (호가, KRX)
//
// Run: KIS 환경변수 설정 후 go run ./examples/ws_krx_basic
//
//	KOREA_INVESTMENT_APP_KEY=...
//	KOREA_INVESTMENT_APP_SECRET=...
//	KOREA_INVESTMENT_ACCOUNT_NO=...
//
// Ctrl+C 로 종료. 자동 재연결 + 구독 자동 복원 동작.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	kis "github.com/kenshin579/korea-investment-stock"
)

func main() {
	client, err := kis.NewClientFromEnv()
	if err != nil {
		log.Fatalf("NewClientFromEnv: %v", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	client.WS.OnConnected(func() {
		fmt.Println("=== WebSocket 연결됨 ===")
	})
	client.WS.OnReconnect(func(attempt int) {
		fmt.Printf(">>> 재연결 #%d (구독 자동 복원)\n", attempt)
	})
	client.WS.OnDisconnect(func(err error) {
		fmt.Printf(">>> 연결 끊김: %v\n", err)
	})
	client.WS.OnError(func(err error) {
		fmt.Printf(">>> ERROR: %v\n", err)
	})

	client.WS.OnKrxTrade(func(ev kis.KrxTradeEvent) { // 또는 websocket.KrxTradeEvent
		fmt.Printf("[체결] %s %s @ %s원 (vol=%d, accum=%d)\n",
			ev.Symbol, ev.Time, ev.Price.String(), ev.TradeVolume, ev.AccumVolume)
	})
	client.WS.OnKrxAsk(func(ev kis.KrxAskEvent) {
		fmt.Printf("[호가] %s %s 매도1=%s 매수1=%s\n",
			ev.Symbol, ev.Time, ev.Ask[0].String(), ev.Bid[0].String())
	})

	// 구독 (Run 전 등록 OK — 연결되면 자동 송신)
	if err := client.WS.SubscribeKrxTrade("005930"); err != nil {
		log.Fatalf("SubscribeKrxTrade: %v", err)
	}
	if err := client.WS.SubscribeKrxAsk("005930"); err != nil {
		log.Fatalf("SubscribeKrxAsk: %v", err)
	}

	fmt.Println("Ctrl+C 로 종료...")
	if err := client.WS.Run(ctx); err != nil && err != context.Canceled {
		log.Printf("Run: %v", err)
	}
	fmt.Println("종료")
}
```

> 만약 `kis.KrxTradeEvent` re-export 안 했으면 `websocket.KrxTradeEvent` 로 import. 결정은 Task 13 정확한 client.go 변경 시점 확인.

- [ ] **Step 2: README.md 수정**

`README.md` 의 "Available Methods" 섹션 끝 (Bonds 다음) 에 추가:
```markdown
### WebSocket — Phase 8 (v1.18.0)

| Method | TR_ID | 설명 |
|--------|-------|------|
| `WS.SubscribeKrxTrade` / `OnKrxTrade` | H0STCNT0 | 국내주식 실시간체결가 (KRX) |
| `WS.SubscribeKrxAsk` / `OnKrxAsk` | H0STASP0 | 국내주식 실시간호가 (KRX) |
| `WS.SubscribeKrxExpectTrade` / `OnKrxExpectTrade` | H0STANC0 | 국내주식 실시간예상체결 (KRX) |
| `WS.SubscribeKrxOvernightTrade` / `OnKrxOvernightTrade` | H0STOAC0 | 국내주식 시간외 실시간체결가 (KRX) |
| `WS.SubscribeKrxOvernightExpect` / `OnKrxOvernightExpect` | H0STOAA0 | 국내주식 시간외 실시간예상체결 (KRX) |

자동 재연결 + 구독 자동 복원 (exp backoff). 사용 예: `examples/ws_krx_basic/`.
```

- [ ] **Step 3: CLAUDE.md 수정**

```markdown
> **Phase 8 — WebSocket KRX 5 endpoint (v1.18.0). 누적 121 REST + 5 WS = 126 endpoints.**
```

추가 spec 링크:
```markdown
- Phase 8 design spec: [`docs/superpowers/specs/2026-05-09-phase8-websocket-design.md`](docs/superpowers/specs/2026-05-09-phase8-websocket-design.md)
- Phase 8 implementation plan: [`docs/superpowers/plans/2026-05-09-phase8-websocket.md`](docs/superpowers/plans/2026-05-09-phase8-websocket.md)
```

- [ ] **Step 4: CHANGELOG.md 수정**

```markdown
## [1.18.0] - 2026-05-09

### Added — Phase 8 (WebSocket — KRX 시세 5 endpoint)

- `client.WS` — 신규 top-level WebSocket client (`websocket/` 패키지)
- `WS.SubscribeKrxTrade` / `OnKrxTrade` — 실시간체결가 KRX (H0STCNT0)
- `WS.SubscribeKrxAsk` / `OnKrxAsk` — 실시간호가 KRX (H0STASP0)
- `WS.SubscribeKrxExpectTrade` / `OnKrxExpectTrade` — 실시간예상체결 KRX (H0STANC0)
- `WS.SubscribeKrxOvernightTrade` / `OnKrxOvernightTrade` — 시간외 체결가 (H0STOAC0)
- `WS.SubscribeKrxOvernightExpect` / `OnKrxOvernightExpect` — 시간외 예상체결 (H0STOAA0)
- ApprovalKeyManager: `/oauth2/Approval` 23h TTL 캐시
- 자동 재연결 + 구독 자동 복원 (exp backoff, max 10 attempts)
- examples: `ws_krx_basic`

### Notes

- 첫 architecture 변경 (REST → WebSocket).
- WebSocket 라이브러리: `github.com/coder/websocket` (구 nhooyr.io).
- Phase 8 = KRX 시세 5 endpoint 만. NXT/통합/ELW/지수/해외/선물옵션 실시간 + 체결통보 (암호화) → Phase 9+.
- Single connection per WS Client. multi-connection 은 사용자 책임.
- Handler 는 reader goroutine 에서 동기 실행 — 무거운 작업은 사용자가 channel 로 fan-out.
- 누적 121 REST + 5 WS = 126 endpoints.
```

- [ ] **Step 5: 최종 점검 + PR**

```bash
gofmt -l .
go vet ./...
go build ./...
go test -race ./...
go test -cover ./websocket/
```

Expected: 모두 clean / PASS / coverage ≥70%.

```bash
git add examples/ws_krx_basic/ README.md CLAUDE.md CHANGELOG.md
git commit -m "[docs] Phase 8 — example + README + CLAUDE + CHANGELOG (v1.18.0)"
git push -u origin feat/phase8-websocket
```

PR 생성:
```bash
gh pr create --title "[feat] Phase 8 — WebSocket KRX 시세 5 endpoint (v1.18.0)" --body "$(cat <<'EOF'
## Summary

Phase 8 — 첫 WebSocket 도입. 인프라 (인증/연결/구독/재연결/decoding) + 국내주식 KRX 시세 5 endpoint.

| EP | TR_ID | 설명 |
|---|---|---|
| 1 | H0STCNT0 | 실시간체결가 (KRX) |
| 2 | H0STASP0 | 실시간호가 (KRX) |
| 3 | H0STANC0 | 실시간예상체결 (KRX) |
| 4 | H0STOAC0 | 시간외 실시간체결가 (KRX) |
| 5 | H0STOAA0 | 시간외 실시간예상체결 (KRX) |

**누적 121 REST + 5 WS = 126 endpoints.**

## 디자인

- Architecture: 신규 top-level `websocket/` 패키지 (Approach A)
- API style: callback handler (`client.WS.OnKrxTrade(handler)`)
- Reconnect: 자동 + 구독 자동 복원 (exp backoff, max 10)
- WS 라이브러리: `github.com/coder/websocket`
- Single connection / single approval_key / 23h TTL 캐시

Spec: `docs/superpowers/specs/2026-05-09-phase8-websocket-design.md`
Plan: `docs/superpowers/plans/2026-05-09-phase8-websocket.md`

## Test plan

- [x] `go build ./...` clean
- [x] `go vet ./...` clean
- [x] `gofmt -l .` clean
- [x] `go test -race ./...` PASS (기존 121 메서드 회귀 없음)
- [x] `go test -cover ./websocket/` ≥70%
- [x] integration tests (5 시나리오: happy/reconnect/error/pingpong/shutdown) via wsmock
- [ ] (선택) `examples/ws_krx_basic` 실 KIS smoke

## Out of Scope (Phase 8)

- NXT/통합 변형, ELW, 지수, 해외주식, 선물옵션 실시간 → Phase 9+
- 체결통보 (AES256 복호화) → 별도 phase
- Multi-connection / backpressure queue → 필요 시

🤖 Generated with [Claude Code](https://claude.com/claude-code)
EOF
)" --reviewer kenshin579
```

---

## Self-Review

**Spec coverage**:
- §2 5 EP → Task 2 (schema) + Task 4 (events) + Task 8 (decode) + Task 12 (subscribe/handler/route) ✓
- §3 결정 → Task 1 (deps), Task 12 (callback API), Task 6 (reconnect), Task 13 (kis.WS) ✓
- §4 architecture → Task 1 (scaffold), Task 12 (Run loop) ✓
- §5 components → Task 5 (Subscriber), 6 (Reconnect), 7 (frame), 9 (Dispatcher), 10 (Approval), 11 (Conn), 12 (Client) ✓
- §6 data flow → Task 12 (Run/readLoop/handleFrame/restoreSubs) ✓
- §7 error handling → Task 3 (errors), 7 (parse), 9 (panic recover), 12 (route) ✓
- §8 testing → Task 5/6/7/8/9/10 unit + Task 14 (wsmock) + Task 15 (integration) ✓
- §9 진입/종료 → Task 1 (branch), Task 16 (release) ✓

**Placeholder scan**: 일부 "Task 2 schema 결과 반영" / "fieldsPerChunk = 46 (실제 값 반영 필요)" 가 있음 — 이는 실제 docs analyzer 결과로 결정될 동적 값이므로 의도된 placeholder. Plan 자체의 task 1 (Setup) + Task 2 (schema 추출) 가 그 결과를 produce. 코드 골격은 완전 (실제 인덱스 / 필드 개수만 schema 결과로 swap).

**Type consistency**:
- `KrxTradeEvent` / `KrxAskEvent` / `KrxExpectTradeEvent` 명명 spec §5-3 와 plan Task 4/8/9/12 일관 ✓
- `subKey{TrID, TrKey}` Task 5 와 Task 12 (`restoreSubs` 의 `k.TrID`/`k.TrKey`) 일관 ✓
- `frame{Kind, Encrypted, TrID, Count, Fields, JSON}` Task 7 정의 ↔ Task 12 `routeRealtime` 사용 일관 ✓
- `Options` 필드명 Task 12 (`opts.Logger`/`opts.CustType` 등) ↔ NewClient 사용 일관 ✓
- `dispatcher` 의 `Route*` / `On*` 메서드 Task 9 정의 ↔ Task 12 `client.dispatcher.RouteXxx()` 일관 ✓

이슈 없음.
