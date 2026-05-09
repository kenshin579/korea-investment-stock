package websocket

import (
	"strings"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestRouteRealtime_Phase11_3_All — 6 routeRealtime case 직접 호출 (coverage 보강).
func TestRouteRealtime_Phase11_3_All(t *testing.T) {
	c := NewClient(Options{Endpoint: "ws://localhost:1", BaseURL: "http://x"})

	var (
		ifTrade, ifAsk, ioTrade, ioAsk, cfTrade, cfAsk atomic.Int32
	)
	c.OnIndexFuturesTrade(func(IndexFuturesTradeEvent) { ifTrade.Add(1) })
	c.OnIndexFuturesAsk(func(IndexFuturesAskEvent) { ifAsk.Add(1) })
	c.OnIndexOptionTrade(func(IndexOptionTradeEvent) { ioTrade.Add(1) })
	c.OnIndexOptionAsk(func(IndexOptionAskEvent) { ioAsk.Add(1) })
	c.OnCommodityFuturesTrade(func(CommodityFuturesTradeEvent) { cfTrade.Add(1) })
	c.OnCommodityFuturesAsk(func(CommodityFuturesAskEvent) { cfAsk.Add(1) })

	loadIntoFrame := func(t *testing.T, name, trID string) frame {
		raw := loadFixture(t, name)
		// raw = "0|TR_ID|001|payload"
		parts := strings.SplitN(raw, "|", 4)
		return frame{Kind: frameKindRealtime, TrID: trID, Count: 1, Fields: strings.Split(parts[3], "^")}
	}

	c.routeRealtime(loadIntoFrame(t, "h0ifcnt0_success.txt", trIDIndexFuturesTrade))
	c.routeRealtime(loadIntoFrame(t, "h0ifasp0_success.txt", trIDIndexFuturesAsk))
	c.routeRealtime(loadIntoFrame(t, "h0iocnt0_success.txt", trIDIndexOptionTrade))
	c.routeRealtime(loadIntoFrame(t, "h0ioasp0_success.txt", trIDIndexOptionAsk))
	c.routeRealtime(loadIntoFrame(t, "h0cfcnt0_success.txt", trIDCommodityFuturesTrade))
	c.routeRealtime(loadIntoFrame(t, "h0cfasp0_success.txt", trIDCommodityFuturesAsk))

	assert.Equal(t, int32(1), ifTrade.Load())
	assert.Equal(t, int32(1), ifAsk.Load())
	assert.Equal(t, int32(1), ioTrade.Load())
	assert.Equal(t, int32(1), ioAsk.Load())
	assert.Equal(t, int32(1), cfTrade.Load())
	assert.Equal(t, int32(1), cfAsk.Load())
}
