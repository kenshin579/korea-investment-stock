package websocket

import (
	"strings"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestRouteRealtime_Phase11_7_All — 2 routeRealtime case 직접 호출 (coverage 보강).
func TestRouteRealtime_Phase11_7_All(t *testing.T) {
	c := NewClient(Options{Endpoint: "ws://localhost:1", BaseURL: "http://x"})

	var (
		trade atomic.Int32
		ask   atomic.Int32
	)
	c.OnOverseasFuturesTrade(func(OverseasFuturesTradeEvent) { trade.Add(1) })
	c.OnOverseasFuturesAsk(func(OverseasFuturesAskEvent) { ask.Add(1) })

	loadIntoFrame := func(t *testing.T, name, trID string) frame {
		raw := loadFixture(t, name)
		// raw = "0|TR_ID|001|payload"
		parts := strings.SplitN(raw, "|", 4)
		return frame{Kind: frameKindRealtime, TrID: trID, Count: 1, Fields: strings.Split(parts[3], "^")}
	}

	c.routeRealtime(loadIntoFrame(t, "hdfff020_success.txt", trIDOverseasFuturesTrade))
	c.routeRealtime(loadIntoFrame(t, "hdfff010_success.txt", trIDOverseasFuturesAsk))

	assert.Equal(t, int32(1), trade.Load())
	assert.Equal(t, int32(1), ask.Load())
}
