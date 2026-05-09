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
