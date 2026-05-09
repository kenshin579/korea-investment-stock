package websocket_test

import (
	"context"
	"net/http"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kenshin579/korea-investment-stock/websocket"
	"github.com/kenshin579/korea-investment-stock/websocket/internal/wsmock"
)

// approvalClient 는 httpmock 전용 HTTP 클라이언트.
// http.DefaultClient 를 오염시키지 않아야 wslib.Dial 이 실제 TCP로 연결 가능.
var approvalClient = &http.Client{}

func setupApprovalMock(t *testing.T) {
	t.Helper()
	// ActivateNonDefault 로 approvalClient 만 intercept → DefaultClient(ws) 는 그대로.
	httpmock.ActivateNonDefault(approvalClient)
	t.Cleanup(func() { httpmock.DeactivateAndReset() })
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
		HTTPClient:    approvalClient, // httpmock 가 가로챔 (DefaultClient 는 건드리지 않음)
	})
}

// samplePayload46 — h0stcnt0 fixture 와 동일 layout, 46 fields.
func samplePayload46(symbol string) string {
	return symbol + "^123929^73100^2^1500^2.09^72850^72500^73200^72400^73100^73000^150^123456^987654000000^5000^7345^2345^120.5^65000^85000^1^53.4^102.3^090030^2^600^102345^2^700^090015^5^200^20260509^11^N^15000^25000^150000^180000^65.5^110000^99.8^0^^72500"
}

func TestIntegration_HappyPath(t *testing.T) {
	setupApprovalMock(t)

	srv := wsmock.New(t)
	defer srv.Close()

	c := newClient(t, srv.URL())

	var received atomic.Int32
	c.OnKrxTrade(func(ev websocket.KrxTradeEvent) {
		received.Add(1)
		assert.Equal(t, "005930", ev.Symbol)
	})

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
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
	require.NoError(t, srv.SendRealtime(ctx, "H0STCNT0", samplePayload46("005930")))

	require.Eventually(t, func() bool { return received.Load() > 0 }, 2*time.Second, 10*time.Millisecond)
}

func TestIntegration_Reconnect(t *testing.T) {
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
	select {
	case <-srv.Received():
	case <-ctx.Done():
		t.Fatal("first subscribe not received")
	}

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
	setupApprovalMock(t)

	srv := wsmock.New(t)
	defer srv.Close()

	c := newClient(t, srv.URL())

	var got atomic.Value // error
	c.OnError(func(err error) { got.Store(err) })

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	go c.Run(ctx)
	require.NoError(t, srv.WaitConnected(ctx))

	require.NoError(t, srv.SendText(ctx, `{"header":{"tr_id":"H0STCNT0"},"body":{"rt_cd":"1","msg_cd":"OPSP0001","msg1":"ALREADY IN SUBSCRIBE"}}`))

	require.Eventually(t, func() bool { return got.Load() != nil }, 2*time.Second, 10*time.Millisecond)
	err, _ := got.Load().(error)
	require.NotNil(t, err)
	wsErr, ok := err.(*websocket.WSServerError)
	require.True(t, ok, "expected *WSServerError, got %T", err)
	assert.Equal(t, "OPSP0001", wsErr.MsgCd)
}

func TestIntegration_PingPong(t *testing.T) {
	setupApprovalMock(t)

	srv := wsmock.New(t)
	defer srv.Close()

	c := newClient(t, srv.URL())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	go c.Run(ctx)
	require.NoError(t, srv.WaitConnected(ctx))

	pingMsg := `{"header":{"tr_id":"PINGPONG"}}`
	require.NoError(t, srv.SendText(ctx, pingMsg))

	// 클라이언트가 echo 응답
	select {
	case msg := <-srv.Received():
		// PING 메시지를 그대로 echo (SendText 가 동일 메시지 송신)
		assert.True(t, strings.Contains(msg, "PINGPONG"))
	case <-time.After(2 * time.Second):
		t.Fatal("did not receive PONG echo")
	}
}

// samplePayloadProgramTrade — H0NXPGM0 / H0UNPGM0 fixture 와 동일 layout, 11 fields.
func samplePayloadProgramTrade(symbol string) string {
	return symbol + "^123929^15000^1100000000^25000^1850000000^10000^750000000^5000^8000^3000"
}

// TestIntegration_NxtTrade_HappyPath — Phase 9 NXT 체결가 시나리오.
// SubscribeNxtTrade → wsmock 이 H0NXCNT0 realtime frame 송신 → OnNxtTrade handler 호출 검증.
func TestIntegration_NxtTrade_HappyPath(t *testing.T) {
	setupApprovalMock(t)

	srv := wsmock.New(t)
	defer srv.Close()

	c := newClient(t, srv.URL())

	var nxtCalls atomic.Int32
	var unCalls atomic.Int32
	c.OnNxtTrade(func(ev websocket.NxtTradeEvent) {
		nxtCalls.Add(1)
		assert.Equal(t, "005930", ev.Symbol)
	})
	c.OnUnifiedTrade(func(ev websocket.UnifiedTradeEvent) {
		// 검증용 — Nxt 와 Unified 가 별도 슬롯인지 확인
		unCalls.Add(1)
	})

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	go c.Run(ctx)
	require.NoError(t, srv.WaitConnected(ctx))

	require.NoError(t, c.SubscribeNxtTrade("005930"))

	select {
	case msg := <-srv.Received():
		assert.Contains(t, msg, "H0NXCNT0")
		assert.Contains(t, msg, "005930")
	case <-ctx.Done():
		t.Fatal("did not receive NXT subscribe frame")
	}

	// 동일 페이로드 layout (KRX 와 schema 호환). 22번 필드만 의미적으로 NXT 의 CNTG_CLS_CODE.
	require.NoError(t, srv.SendRealtime(ctx, "H0NXCNT0", samplePayload46("005930")))

	require.Eventually(t, func() bool { return nxtCalls.Load() > 0 }, 2*time.Second, 10*time.Millisecond)
	// Unified handler 는 호출되지 않아야 함 (별도 슬롯)
	assert.Equal(t, int32(0), unCalls.Load())
}

// TestIntegration_UnifiedProgramTrade — Phase 9 통합 프로그램매매 신규 EP 시나리오.
func TestIntegration_UnifiedProgramTrade(t *testing.T) {
	setupApprovalMock(t)

	srv := wsmock.New(t)
	defer srv.Close()

	c := newClient(t, srv.URL())

	var calls atomic.Int32
	c.OnUnifiedProgramTrade(func(ev websocket.UnifiedProgramTradeEvent) {
		calls.Add(1)
		assert.Equal(t, "005930", ev.Symbol)
		assert.Equal(t, int64(15000), ev.AskQuantity)
		assert.Equal(t, int64(3000), ev.TotalNetQuantity)
	})

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	go c.Run(ctx)
	require.NoError(t, srv.WaitConnected(ctx))

	require.NoError(t, c.SubscribeUnifiedProgramTrade("005930"))

	select {
	case msg := <-srv.Received():
		assert.Contains(t, msg, "H0UNPGM0")
	case <-ctx.Done():
		t.Fatal("did not receive 통합 program-trade subscribe frame")
	}

	require.NoError(t, srv.SendRealtime(ctx, "H0UNPGM0", samplePayloadProgramTrade("005930")))
	require.Eventually(t, func() bool { return calls.Load() > 0 }, 2*time.Second, 10*time.Millisecond)
}

// samplePayloadOverseasTrade — HDFSCNT0 fixture 와 동일 layout, 26 fields.
func samplePayloadOverseasTrade(symbol string) string {
	return symbol + "^AAPL^4^20260509^260509^093015^260509^223015^195.32^196.50^194.80^195.85^2^0.53^0.27^195.83^195.86^1500^1200^150^25000000^4900000000^120000^150000^102.5^1"
}

// samplePayloadStockFuturesTrade — H0ZFCNT0 fixture 와 동일 layout, 49 fields.
func samplePayloadStockFuturesTrade(symbol string) string {
	return symbol + "^133500^73400^2^1500^2.09^73000^74500^72800^50^23456^171500000^73350^-50^-0.07^73500^73300^100^8900^30^133000^2^73350^133200^2^73400^133100^5^73200^52.30^115.80^-0.20^-10^1.50^73500^73300^120^180^300^450^-150^3500^5000^35000^45000^1.85^75000^71000^N"
}

// samplePayloadIndexFuturesTrade — H0IFCNT0 fixture 와 동일 layout, 50 fields.
func samplePayloadIndexFuturesTrade(symbol string) string {
	return symbol + "^090130^5.00^2^0.15^3395.00^3380.00^3400.00^3370.00^150^78234^265234560000^3392.50^-2.50^-0.07^3393.00^3391.00^2.00^87654^200^090000^2^15.00^085520^2^5.00^090010^5^-25.00^54.30^118.50^-0.50^-100^1.20^3396.00^3394.00^320^480^850^1200^350^12500^18000^45000^62000^2.35^0^3500.00^3290.00^N"
}

// TestIntegration_IndexFuturesTrade — Phase 11.3 지수선물 체결가 시나리오 (routeRealtime 보강).
func TestIntegration_IndexFuturesTrade(t *testing.T) {
	setupApprovalMock(t)

	srv := wsmock.New(t)
	defer srv.Close()

	c := newClient(t, srv.URL())

	var calls atomic.Int32
	c.OnIndexFuturesTrade(func(ev websocket.IndexFuturesTradeEvent) {
		calls.Add(1)
		assert.Equal(t, "101S12", ev.Symbol)
	})

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	go c.Run(ctx)
	require.NoError(t, srv.WaitConnected(ctx))

	require.NoError(t, c.SubscribeIndexFuturesTrade("101S12"))

	select {
	case msg := <-srv.Received():
		assert.Contains(t, msg, "H0IFCNT0")
	case <-ctx.Done():
		t.Fatal("did not receive subscribe frame")
	}

	require.NoError(t, srv.SendRealtime(ctx, "H0IFCNT0", samplePayloadIndexFuturesTrade("101S12")))
	require.Eventually(t, func() bool { return calls.Load() > 0 }, 2*time.Second, 10*time.Millisecond)
}

// TestIntegration_OverseasTrade — Phase 10 해외주식 체결가 시나리오.
func TestIntegration_OverseasTrade(t *testing.T) {
	setupApprovalMock(t)

	srv := wsmock.New(t)
	defer srv.Close()

	c := newClient(t, srv.URL())

	var calls atomic.Int32
	c.OnOverseasTrade(func(ev websocket.OverseasTradeEvent) {
		calls.Add(1)
		assert.Equal(t, "DNASAAPL", ev.Symbol)
		assert.Equal(t, "AAPL", ev.SymbolCode)
	})

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	go c.Run(ctx)
	require.NoError(t, srv.WaitConnected(ctx))

	require.NoError(t, c.SubscribeOverseasTrade("DNASAAPL"))

	select {
	case msg := <-srv.Received():
		assert.Contains(t, msg, "HDFSCNT0")
		assert.Contains(t, msg, "DNASAAPL")
	case <-ctx.Done():
		t.Fatal("did not receive overseas subscribe frame")
	}

	require.NoError(t, srv.SendRealtime(ctx, "HDFSCNT0", samplePayloadOverseasTrade("DNASAAPL")))
	require.Eventually(t, func() bool { return calls.Load() > 0 }, 2*time.Second, 10*time.Millisecond)
}

// TestIntegration_StockFuturesTrade — Phase 11.2 주식선물 체결가 시나리오.
func TestIntegration_StockFuturesTrade(t *testing.T) {
	setupApprovalMock(t)

	srv := wsmock.New(t)
	defer srv.Close()

	c := newClient(t, srv.URL())

	var calls atomic.Int32
	c.OnStockFuturesTrade(func(ev websocket.StockFuturesTradeEvent) {
		calls.Add(1)
		assert.Equal(t, "KAK0F", ev.Symbol)
		assert.Equal(t, "73400", ev.Price.String())
	})

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	go c.Run(ctx)
	require.NoError(t, srv.WaitConnected(ctx))

	require.NoError(t, c.SubscribeStockFuturesTrade("KAK0F"))

	select {
	case msg := <-srv.Received():
		assert.Contains(t, msg, "H0ZFCNT0")
		assert.Contains(t, msg, "KAK0F")
	case <-ctx.Done():
		t.Fatal("did not receive stock futures subscribe frame")
	}

	require.NoError(t, srv.SendRealtime(ctx, "H0ZFCNT0", samplePayloadStockFuturesTrade("KAK0F")))
	require.Eventually(t, func() bool { return calls.Load() > 0 }, 2*time.Second, 10*time.Millisecond)
}

func TestIntegration_GracefulShutdown(t *testing.T) {
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
