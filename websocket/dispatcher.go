package websocket

import (
	"fmt"
	"sync"
)

// dispatcher 는 TR_ID 별 handler 등록/라우팅.
//
// NOTE: handler 내부에서 dispatcher 의 On* 등록 메서드를 재호출하면
// RLock 상태에서 Lock 시도로 deadlock 이 발생한다.
// Phase 8 에서는 handler 가 dispatcher 메서드를 재호출하지 않는다고 가정한다.
type dispatcher struct {
	mu sync.RWMutex

	onKrxTrade           func(KrxTradeEvent)
	onKrxAsk             func(KrxAskEvent)
	onKrxExpectTrade     func(KrxExpectTradeEvent)
	onKrxOvernightTrade  func(KrxOvernightTradeEvent)
	onKrxOvernightExpect func(KrxOvernightExpectEvent)

	// Phase 9 — NXT/통합 (5 base type, 10 distinct handler 슬롯)
	onNxtTrade            func(NxtTradeEvent)
	onUnifiedTrade        func(UnifiedTradeEvent)
	onNxtAsk              func(NxtAskEvent)
	onUnifiedAsk          func(UnifiedAskEvent)
	onNxtExpectTrade      func(NxtExpectTradeEvent)
	onUnifiedExpectTrade  func(UnifiedExpectTradeEvent)
	onNxtProgramTrade     func(NxtProgramTradeEvent)
	onUnifiedProgramTrade func(UnifiedProgramTradeEvent)
	onNxtMember           func(NxtMemberEvent)
	onUnifiedMember       func(UnifiedMemberEvent)

	onConnected  func()
	onReconnect  func(attempt int)
	onDisconnect func(error)
	onError      func(error)
}

func newDispatcher() *dispatcher { return &dispatcher{} }

// === 등록 메서드 ===

func (d *dispatcher) OnKrxTrade(h func(KrxTradeEvent)) {
	d.mu.Lock()
	d.onKrxTrade = h
	d.mu.Unlock()
}
func (d *dispatcher) OnKrxAsk(h func(KrxAskEvent)) {
	d.mu.Lock()
	d.onKrxAsk = h
	d.mu.Unlock()
}
func (d *dispatcher) OnKrxExpectTrade(h func(KrxExpectTradeEvent)) {
	d.mu.Lock()
	d.onKrxExpectTrade = h
	d.mu.Unlock()
}
func (d *dispatcher) OnKrxOvernightTrade(h func(KrxOvernightTradeEvent)) {
	d.mu.Lock()
	d.onKrxOvernightTrade = h
	d.mu.Unlock()
}
func (d *dispatcher) OnKrxOvernightExpect(h func(KrxOvernightExpectEvent)) {
	d.mu.Lock()
	d.onKrxOvernightExpect = h
	d.mu.Unlock()
}

// Phase 9 — NXT/통합 등록 메서드 (10)

func (d *dispatcher) OnNxtTrade(h func(NxtTradeEvent)) {
	d.mu.Lock()
	d.onNxtTrade = h
	d.mu.Unlock()
}
func (d *dispatcher) OnUnifiedTrade(h func(UnifiedTradeEvent)) {
	d.mu.Lock()
	d.onUnifiedTrade = h
	d.mu.Unlock()
}
func (d *dispatcher) OnNxtAsk(h func(NxtAskEvent)) {
	d.mu.Lock()
	d.onNxtAsk = h
	d.mu.Unlock()
}
func (d *dispatcher) OnUnifiedAsk(h func(UnifiedAskEvent)) {
	d.mu.Lock()
	d.onUnifiedAsk = h
	d.mu.Unlock()
}
func (d *dispatcher) OnNxtExpectTrade(h func(NxtExpectTradeEvent)) {
	d.mu.Lock()
	d.onNxtExpectTrade = h
	d.mu.Unlock()
}
func (d *dispatcher) OnUnifiedExpectTrade(h func(UnifiedExpectTradeEvent)) {
	d.mu.Lock()
	d.onUnifiedExpectTrade = h
	d.mu.Unlock()
}
func (d *dispatcher) OnNxtProgramTrade(h func(NxtProgramTradeEvent)) {
	d.mu.Lock()
	d.onNxtProgramTrade = h
	d.mu.Unlock()
}
func (d *dispatcher) OnUnifiedProgramTrade(h func(UnifiedProgramTradeEvent)) {
	d.mu.Lock()
	d.onUnifiedProgramTrade = h
	d.mu.Unlock()
}
func (d *dispatcher) OnNxtMember(h func(NxtMemberEvent)) {
	d.mu.Lock()
	d.onNxtMember = h
	d.mu.Unlock()
}
func (d *dispatcher) OnUnifiedMember(h func(UnifiedMemberEvent)) {
	d.mu.Lock()
	d.onUnifiedMember = h
	d.mu.Unlock()
}

func (d *dispatcher) OnConnected(h func()) {
	d.mu.Lock()
	d.onConnected = h
	d.mu.Unlock()
}
func (d *dispatcher) OnReconnect(h func(int)) {
	d.mu.Lock()
	d.onReconnect = h
	d.mu.Unlock()
}
func (d *dispatcher) OnDisconnect(h func(error)) {
	d.mu.Lock()
	d.onDisconnect = h
	d.mu.Unlock()
}
func (d *dispatcher) OnError(h func(error)) {
	d.mu.Lock()
	d.onError = h
	d.mu.Unlock()
}

// === 라우팅 메서드 ===

func (d *dispatcher) RouteKrxTrade(ev KrxTradeEvent) {
	d.safeCall(func(h *dispatcher) {
		if h.onKrxTrade != nil {
			h.onKrxTrade(ev)
		}
	})
}
func (d *dispatcher) RouteKrxAsk(ev KrxAskEvent) {
	d.safeCall(func(h *dispatcher) {
		if h.onKrxAsk != nil {
			h.onKrxAsk(ev)
		}
	})
}
func (d *dispatcher) RouteKrxExpectTrade(ev KrxExpectTradeEvent) {
	d.safeCall(func(h *dispatcher) {
		if h.onKrxExpectTrade != nil {
			h.onKrxExpectTrade(ev)
		}
	})
}
func (d *dispatcher) RouteKrxOvernightTrade(ev KrxOvernightTradeEvent) {
	d.safeCall(func(h *dispatcher) {
		if h.onKrxOvernightTrade != nil {
			h.onKrxOvernightTrade(ev)
		}
	})
}
func (d *dispatcher) RouteKrxOvernightExpect(ev KrxOvernightExpectEvent) {
	d.safeCall(func(h *dispatcher) {
		if h.onKrxOvernightExpect != nil {
			h.onKrxOvernightExpect(ev)
		}
	})
}

// Phase 9 — NXT/통합 라우팅 메서드 (10)

func (d *dispatcher) RouteNxtTrade(ev NxtTradeEvent) {
	d.safeCall(func(h *dispatcher) {
		if h.onNxtTrade != nil {
			h.onNxtTrade(ev)
		}
	})
}
func (d *dispatcher) RouteUnifiedTrade(ev UnifiedTradeEvent) {
	d.safeCall(func(h *dispatcher) {
		if h.onUnifiedTrade != nil {
			h.onUnifiedTrade(ev)
		}
	})
}
func (d *dispatcher) RouteNxtAsk(ev NxtAskEvent) {
	d.safeCall(func(h *dispatcher) {
		if h.onNxtAsk != nil {
			h.onNxtAsk(ev)
		}
	})
}
func (d *dispatcher) RouteUnifiedAsk(ev UnifiedAskEvent) {
	d.safeCall(func(h *dispatcher) {
		if h.onUnifiedAsk != nil {
			h.onUnifiedAsk(ev)
		}
	})
}
func (d *dispatcher) RouteNxtExpectTrade(ev NxtExpectTradeEvent) {
	d.safeCall(func(h *dispatcher) {
		if h.onNxtExpectTrade != nil {
			h.onNxtExpectTrade(ev)
		}
	})
}
func (d *dispatcher) RouteUnifiedExpectTrade(ev UnifiedExpectTradeEvent) {
	d.safeCall(func(h *dispatcher) {
		if h.onUnifiedExpectTrade != nil {
			h.onUnifiedExpectTrade(ev)
		}
	})
}
func (d *dispatcher) RouteNxtProgramTrade(ev NxtProgramTradeEvent) {
	d.safeCall(func(h *dispatcher) {
		if h.onNxtProgramTrade != nil {
			h.onNxtProgramTrade(ev)
		}
	})
}
func (d *dispatcher) RouteUnifiedProgramTrade(ev UnifiedProgramTradeEvent) {
	d.safeCall(func(h *dispatcher) {
		if h.onUnifiedProgramTrade != nil {
			h.onUnifiedProgramTrade(ev)
		}
	})
}
func (d *dispatcher) RouteNxtMember(ev NxtMemberEvent) {
	d.safeCall(func(h *dispatcher) {
		if h.onNxtMember != nil {
			h.onNxtMember(ev)
		}
	})
}
func (d *dispatcher) RouteUnifiedMember(ev UnifiedMemberEvent) {
	d.safeCall(func(h *dispatcher) {
		if h.onUnifiedMember != nil {
			h.onUnifiedMember(ev)
		}
	})
}

func (d *dispatcher) RouteConnected() {
	d.safeCall(func(h *dispatcher) {
		if h.onConnected != nil {
			h.onConnected()
		}
	})
}
func (d *dispatcher) RouteReconnect(att int) {
	d.safeCall(func(h *dispatcher) {
		if h.onReconnect != nil {
			h.onReconnect(att)
		}
	})
}
func (d *dispatcher) RouteDisconnect(e error) {
	d.safeCall(func(h *dispatcher) {
		if h.onDisconnect != nil {
			h.onDisconnect(e)
		}
	})
}

// RouteError 는 panic recover 후에도 직접 호출되므로 safeCall 없이 직접 호출.
func (d *dispatcher) RouteError(e error) {
	d.mu.RLock()
	h := d.onError
	d.mu.RUnlock()
	if h != nil {
		h(e)
	}
}

// safeCall — handler panic 을 OnError 로 라우팅.
// pattern: RLock 획득 → handler 실행 (with recover) → panic 시 onError 호출.
func (d *dispatcher) safeCall(fn func(h *dispatcher)) {
	defer func() {
		if r := recover(); r != nil {
			d.mu.RLock()
			eh := d.onError
			d.mu.RUnlock()
			if eh != nil {
				eh(fmt.Errorf("kis ws: handler panic: %v", r))
			}
		}
	}()
	d.mu.RLock()
	defer d.mu.RUnlock()
	fn(d)
}
