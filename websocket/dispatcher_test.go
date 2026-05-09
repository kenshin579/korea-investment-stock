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

func TestDispatcher_Phase11_2Routers(t *testing.T) {
	// Phase 11.2 — 국내선물옵션 11 핸들러 라우팅 검증.
	d := newDispatcher()
	var (
		krxNightFutTrade   atomic.Int32
		krxNightFutAsk     atomic.Int32
		krxNightOptTrade   atomic.Int32
		krxNightOptAsk     atomic.Int32
		krxNightOptExpTrade atomic.Int32
		stockFutTrade      atomic.Int32
		stockFutAsk        atomic.Int32
		stockFutExpTrade   atomic.Int32
		stockOptTrade      atomic.Int32
		stockOptAsk        atomic.Int32
		stockOptExpTrade   atomic.Int32
	)

	d.OnKrxNightFuturesTrade(func(KrxNightFuturesTradeEvent) { krxNightFutTrade.Add(1) })
	d.OnKrxNightFuturesAsk(func(KrxNightFuturesAskEvent) { krxNightFutAsk.Add(1) })
	d.OnKrxNightOptionTrade(func(KrxNightOptionTradeEvent) { krxNightOptTrade.Add(1) })
	d.OnKrxNightOptionAsk(func(KrxNightOptionAskEvent) { krxNightOptAsk.Add(1) })
	d.OnKrxNightOptionExpectTrade(func(KrxNightOptionExpectTradeEvent) { krxNightOptExpTrade.Add(1) })
	d.OnStockFuturesTrade(func(StockFuturesTradeEvent) { stockFutTrade.Add(1) })
	d.OnStockFuturesAsk(func(StockFuturesAskEvent) { stockFutAsk.Add(1) })
	d.OnStockFuturesExpectTrade(func(StockFuturesExpectTradeEvent) { stockFutExpTrade.Add(1) })
	d.OnStockOptionTrade(func(StockOptionTradeEvent) { stockOptTrade.Add(1) })
	d.OnStockOptionAsk(func(StockOptionAskEvent) { stockOptAsk.Add(1) })
	d.OnStockOptionExpectTrade(func(StockOptionExpectTradeEvent) { stockOptExpTrade.Add(1) })

	d.RouteKrxNightFuturesTrade(KrxNightFuturesTradeEvent{})
	d.RouteKrxNightFuturesAsk(KrxNightFuturesAskEvent{})
	d.RouteKrxNightOptionTrade(KrxNightOptionTradeEvent{})
	d.RouteKrxNightOptionAsk(KrxNightOptionAskEvent{})
	d.RouteKrxNightOptionExpectTrade(KrxNightOptionExpectTradeEvent{})
	d.RouteStockFuturesTrade(StockFuturesTradeEvent{})
	d.RouteStockFuturesAsk(StockFuturesAskEvent{})
	d.RouteStockFuturesExpectTrade(StockFuturesExpectTradeEvent{})
	d.RouteStockOptionTrade(StockOptionTradeEvent{})
	d.RouteStockOptionAsk(StockOptionAskEvent{})
	d.RouteStockOptionExpectTrade(StockOptionExpectTradeEvent{})

	// 11개 핸들러 모두 정확히 1번 호출되어야 함
	assert.Equal(t, int32(1), krxNightFutTrade.Load())
	assert.Equal(t, int32(1), krxNightFutAsk.Load())
	assert.Equal(t, int32(1), krxNightOptTrade.Load())
	assert.Equal(t, int32(1), krxNightOptAsk.Load())
	assert.Equal(t, int32(1), krxNightOptExpTrade.Load())
	assert.Equal(t, int32(1), stockFutTrade.Load())
	assert.Equal(t, int32(1), stockFutAsk.Load())
	assert.Equal(t, int32(1), stockFutExpTrade.Load())
	assert.Equal(t, int32(1), stockOptTrade.Load())
	assert.Equal(t, int32(1), stockOptAsk.Load())
	assert.Equal(t, int32(1), stockOptExpTrade.Load())
}
