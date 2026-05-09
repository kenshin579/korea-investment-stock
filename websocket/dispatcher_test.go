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
	var errs []error
	d.OnError(func(err error) {
		errs = append(errs, err)
	})
	d.OnKrxTrade(func(ev KrxTradeEvent) {
		panic("boom")
	})

	d.RouteKrxTrade(KrxTradeEvent{Symbol: "005930"})
	assert.Len(t, errs, 1)
	assert.Contains(t, errs[0].Error(), "panic")
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

func TestDispatcher_AllRouters(t *testing.T) {
	// 5 EP + lifecycle handlers 가 모두 라우팅되는지 검증
	d := newDispatcher()
	var (
		trade, ask, expectTrade atomic.Int32
		ovTrade, ovExpect       atomic.Int32
		connected, disconnected atomic.Int32
		reconnect               atomic.Int32
	)

	d.OnKrxTrade(func(KrxTradeEvent) { trade.Add(1) })
	d.OnKrxAsk(func(KrxAskEvent) { ask.Add(1) })
	d.OnKrxExpectTrade(func(KrxExpectTradeEvent) { expectTrade.Add(1) })
	d.OnKrxOvernightTrade(func(KrxOvernightTradeEvent) { ovTrade.Add(1) })
	d.OnKrxOvernightExpect(func(KrxOvernightExpectEvent) { ovExpect.Add(1) })
	d.OnConnected(func() { connected.Add(1) })
	d.OnDisconnect(func(error) { disconnected.Add(1) })
	d.OnReconnect(func(int) { reconnect.Add(1) })

	d.RouteKrxTrade(KrxTradeEvent{})
	d.RouteKrxAsk(KrxAskEvent{})
	d.RouteKrxExpectTrade(KrxExpectTradeEvent{})
	d.RouteKrxOvernightTrade(KrxOvernightTradeEvent{})
	d.RouteKrxOvernightExpect(KrxOvernightExpectEvent{})
	d.RouteConnected()
	d.RouteDisconnect(errors.New("test"))
	d.RouteReconnect(3)

	assert.Equal(t, int32(1), trade.Load())
	assert.Equal(t, int32(1), ask.Load())
	assert.Equal(t, int32(1), expectTrade.Load())
	assert.Equal(t, int32(1), ovTrade.Load())
	assert.Equal(t, int32(1), ovExpect.Load())
	assert.Equal(t, int32(1), connected.Load())
	assert.Equal(t, int32(1), disconnected.Load())
	assert.Equal(t, int32(1), reconnect.Load())
}

func TestDispatcher_Phase9Routers(t *testing.T) {
	// Phase 9 — NXT/통합 10 핸들러 라우팅 검증.
	// alias type 이라 Nxt 와 Unified 같은 base 지만 handler 슬롯은 별도여야 함.
	d := newDispatcher()
	var (
		nxtTrade, unTrade   atomic.Int32
		nxtAsk, unAsk       atomic.Int32
		nxtExp, unExp       atomic.Int32
		nxtPgm, unPgm       atomic.Int32
		nxtMember, unMember atomic.Int32
	)

	d.OnNxtTrade(func(NxtTradeEvent) { nxtTrade.Add(1) })
	d.OnUnifiedTrade(func(UnifiedTradeEvent) { unTrade.Add(1) })
	d.OnNxtAsk(func(NxtAskEvent) { nxtAsk.Add(1) })
	d.OnUnifiedAsk(func(UnifiedAskEvent) { unAsk.Add(1) })
	d.OnNxtExpectTrade(func(NxtExpectTradeEvent) { nxtExp.Add(1) })
	d.OnUnifiedExpectTrade(func(UnifiedExpectTradeEvent) { unExp.Add(1) })
	d.OnNxtProgramTrade(func(NxtProgramTradeEvent) { nxtPgm.Add(1) })
	d.OnUnifiedProgramTrade(func(UnifiedProgramTradeEvent) { unPgm.Add(1) })
	d.OnNxtMember(func(NxtMemberEvent) { nxtMember.Add(1) })
	d.OnUnifiedMember(func(UnifiedMemberEvent) { unMember.Add(1) })

	d.RouteNxtTrade(NxtTradeEvent{})
	d.RouteUnifiedTrade(UnifiedTradeEvent{})
	d.RouteNxtAsk(NxtAskEvent{})
	d.RouteUnifiedAsk(UnifiedAskEvent{})
	d.RouteNxtExpectTrade(NxtExpectTradeEvent{})
	d.RouteUnifiedExpectTrade(UnifiedExpectTradeEvent{})
	d.RouteNxtProgramTrade(NxtProgramTradeEvent{})
	d.RouteUnifiedProgramTrade(UnifiedProgramTradeEvent{})
	d.RouteNxtMember(NxtMemberEvent{})
	d.RouteUnifiedMember(UnifiedMemberEvent{})

	// Nxt 와 Unified 는 별도 슬롯 — 각각 1번씩만 호출
	assert.Equal(t, int32(1), nxtTrade.Load())
	assert.Equal(t, int32(1), unTrade.Load())
	assert.Equal(t, int32(1), nxtAsk.Load())
	assert.Equal(t, int32(1), unAsk.Load())
	assert.Equal(t, int32(1), nxtExp.Load())
	assert.Equal(t, int32(1), unExp.Load())
	assert.Equal(t, int32(1), nxtPgm.Load())
	assert.Equal(t, int32(1), unPgm.Load())
	assert.Equal(t, int32(1), nxtMember.Load())
	assert.Equal(t, int32(1), unMember.Load())
}

func TestDispatcher_Phase10Routers(t *testing.T) {
	// Phase 10 — 해외주식 2 핸들러 라우팅 검증.
	d := newDispatcher()
	var (
		ovTrade atomic.Int32
		ovAsk   atomic.Int32
	)

	d.OnOverseasTrade(func(OverseasTradeEvent) { ovTrade.Add(1) })
	d.OnOverseasAsk(func(OverseasAskEvent) { ovAsk.Add(1) })

	d.RouteOverseasTrade(OverseasTradeEvent{})
	d.RouteOverseasAsk(OverseasAskEvent{})

	assert.Equal(t, int32(1), ovTrade.Load())
	assert.Equal(t, int32(1), ovAsk.Load())
}
