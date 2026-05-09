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
	Endpoint      string // ws://ops.koreainvestment.com:21000 (실전) / :31000 (모의)
	BaseURL       string // /oauth2/Approval base URL (예: https://openapi.koreainvestment.com:9443)
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

// TR_ID 상수 — Task 2 schema 분석으로 확정 (plan 의 추정 일부 정정).
const (
	// Phase 8 — KRX
	trIDKrxTrade           = "H0STCNT0" // 실시간체결가 KRX
	trIDKrxAsk             = "H0STASP0" // 실시간호가 KRX
	trIDKrxExpectTrade     = "H0STANC0" // 실시간예상체결 KRX
	trIDKrxOvernightTrade  = "H0STOUP0" // 시간외 실시간체결가 KRX
	trIDKrxOvernightExpect = "H0STOAC0" // 시간외 실시간예상체결 KRX

	// Phase 9 — NXT
	trIDNxtTrade        = "H0NXCNT0" // 실시간체결가 NXT
	trIDNxtAsk          = "H0NXASP0" // 실시간호가 NXT
	trIDNxtExpectTrade  = "H0NXANC0" // 실시간예상체결 NXT
	trIDNxtProgramTrade = "H0NXPGM0" // 실시간프로그램매매 NXT
	trIDNxtMember       = "H0NXMBC0" // 실시간회원사 NXT

	// Phase 9 — 통합
	trIDUnifiedTrade        = "H0UNCNT0" // 실시간체결가 통합
	trIDUnifiedAsk          = "H0UNASP0" // 실시간호가 통합
	trIDUnifiedExpectTrade  = "H0UNANC0" // 실시간예상체결 통합
	trIDUnifiedProgramTrade = "H0UNPGM0" // 실시간프로그램매매 통합
	trIDUnifiedMember       = "H0UNMBC0" // 실시간회원사 통합
)

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

// NewClient 는 Options 로 WS Client 생성.
func NewClient(opts Options) *Client {
	opts.defaults()
	return &Client{
		opts:       opts,
		approval:   newApprovalKeyManager(opts.HTTPClient, opts.BaseURL, opts.AppKey, opts.AppSecret, opts.ApprovalTTL),
		conn:       newConnManager(opts.Endpoint),
		sub:        newSubscriber(),
		dispatcher: newDispatcher(),
		reconnect: newReconnectController(reconnectOpts{
			Min:         opts.ReconnectMin,
			Max:         opts.ReconnectMax,
			MaxAttempts: opts.MaxReconnects,
		}),
	}
}

// === Subscribe ===

func (c *Client) SubscribeKrxTrade(symbols ...string) error {
	return c.subscribe(trIDKrxTrade, symbols)
}
func (c *Client) SubscribeKrxAsk(symbols ...string) error {
	return c.subscribe(trIDKrxAsk, symbols)
}
func (c *Client) SubscribeKrxExpectTrade(symbols ...string) error {
	return c.subscribe(trIDKrxExpectTrade, symbols)
}
func (c *Client) SubscribeKrxOvernightTrade(symbols ...string) error {
	return c.subscribe(trIDKrxOvernightTrade, symbols)
}
func (c *Client) SubscribeKrxOvernightExpect(symbols ...string) error {
	return c.subscribe(trIDKrxOvernightExpect, symbols)
}

// Phase 9 — NXT Subscribe
func (c *Client) SubscribeNxtTrade(symbols ...string) error {
	return c.subscribe(trIDNxtTrade, symbols)
}
func (c *Client) SubscribeNxtAsk(symbols ...string) error {
	return c.subscribe(trIDNxtAsk, symbols)
}
func (c *Client) SubscribeNxtExpectTrade(symbols ...string) error {
	return c.subscribe(trIDNxtExpectTrade, symbols)
}
func (c *Client) SubscribeNxtProgramTrade(symbols ...string) error {
	return c.subscribe(trIDNxtProgramTrade, symbols)
}
func (c *Client) SubscribeNxtMember(symbols ...string) error {
	return c.subscribe(trIDNxtMember, symbols)
}

// Phase 9 — 통합 Subscribe
func (c *Client) SubscribeUnifiedTrade(symbols ...string) error {
	return c.subscribe(trIDUnifiedTrade, symbols)
}
func (c *Client) SubscribeUnifiedAsk(symbols ...string) error {
	return c.subscribe(trIDUnifiedAsk, symbols)
}
func (c *Client) SubscribeUnifiedExpectTrade(symbols ...string) error {
	return c.subscribe(trIDUnifiedExpectTrade, symbols)
}
func (c *Client) SubscribeUnifiedProgramTrade(symbols ...string) error {
	return c.subscribe(trIDUnifiedProgramTrade, symbols)
}
func (c *Client) SubscribeUnifiedMember(symbols ...string) error {
	return c.subscribe(trIDUnifiedMember, symbols)
}

// === Unsubscribe (대칭) ===

func (c *Client) UnsubscribeKrxTrade(symbols ...string) error {
	return c.unsubscribe(trIDKrxTrade, symbols)
}
func (c *Client) UnsubscribeKrxAsk(symbols ...string) error {
	return c.unsubscribe(trIDKrxAsk, symbols)
}
func (c *Client) UnsubscribeKrxExpectTrade(symbols ...string) error {
	return c.unsubscribe(trIDKrxExpectTrade, symbols)
}
func (c *Client) UnsubscribeKrxOvernightTrade(symbols ...string) error {
	return c.unsubscribe(trIDKrxOvernightTrade, symbols)
}
func (c *Client) UnsubscribeKrxOvernightExpect(symbols ...string) error {
	return c.unsubscribe(trIDKrxOvernightExpect, symbols)
}

// Phase 9 — NXT Unsubscribe
func (c *Client) UnsubscribeNxtTrade(symbols ...string) error {
	return c.unsubscribe(trIDNxtTrade, symbols)
}
func (c *Client) UnsubscribeNxtAsk(symbols ...string) error {
	return c.unsubscribe(trIDNxtAsk, symbols)
}
func (c *Client) UnsubscribeNxtExpectTrade(symbols ...string) error {
	return c.unsubscribe(trIDNxtExpectTrade, symbols)
}
func (c *Client) UnsubscribeNxtProgramTrade(symbols ...string) error {
	return c.unsubscribe(trIDNxtProgramTrade, symbols)
}
func (c *Client) UnsubscribeNxtMember(symbols ...string) error {
	return c.unsubscribe(trIDNxtMember, symbols)
}

// Phase 9 — 통합 Unsubscribe
func (c *Client) UnsubscribeUnifiedTrade(symbols ...string) error {
	return c.unsubscribe(trIDUnifiedTrade, symbols)
}
func (c *Client) UnsubscribeUnifiedAsk(symbols ...string) error {
	return c.unsubscribe(trIDUnifiedAsk, symbols)
}
func (c *Client) UnsubscribeUnifiedExpectTrade(symbols ...string) error {
	return c.unsubscribe(trIDUnifiedExpectTrade, symbols)
}
func (c *Client) UnsubscribeUnifiedProgramTrade(symbols ...string) error {
	return c.unsubscribe(trIDUnifiedProgramTrade, symbols)
}
func (c *Client) UnsubscribeUnifiedMember(symbols ...string) error {
	return c.unsubscribe(trIDUnifiedMember, symbols)
}

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

// === Handler 위임 (5 distinct types) ===

func (c *Client) OnKrxTrade(h func(KrxTradeEvent))             { c.dispatcher.OnKrxTrade(h) }
func (c *Client) OnKrxAsk(h func(KrxAskEvent))                 { c.dispatcher.OnKrxAsk(h) }
func (c *Client) OnKrxExpectTrade(h func(KrxExpectTradeEvent)) { c.dispatcher.OnKrxExpectTrade(h) }
func (c *Client) OnKrxOvernightTrade(h func(KrxOvernightTradeEvent)) {
	c.dispatcher.OnKrxOvernightTrade(h)
}
func (c *Client) OnKrxOvernightExpect(h func(KrxOvernightExpectEvent)) {
	c.dispatcher.OnKrxOvernightExpect(h)
}

// Phase 9 — NXT/통합 Handler 위임 (10)
func (c *Client) OnNxtTrade(h func(NxtTradeEvent))             { c.dispatcher.OnNxtTrade(h) }
func (c *Client) OnUnifiedTrade(h func(UnifiedTradeEvent))     { c.dispatcher.OnUnifiedTrade(h) }
func (c *Client) OnNxtAsk(h func(NxtAskEvent))                 { c.dispatcher.OnNxtAsk(h) }
func (c *Client) OnUnifiedAsk(h func(UnifiedAskEvent))         { c.dispatcher.OnUnifiedAsk(h) }
func (c *Client) OnNxtExpectTrade(h func(NxtExpectTradeEvent)) { c.dispatcher.OnNxtExpectTrade(h) }
func (c *Client) OnUnifiedExpectTrade(h func(UnifiedExpectTradeEvent)) {
	c.dispatcher.OnUnifiedExpectTrade(h)
}
func (c *Client) OnNxtProgramTrade(h func(NxtProgramTradeEvent)) {
	c.dispatcher.OnNxtProgramTrade(h)
}
func (c *Client) OnUnifiedProgramTrade(h func(UnifiedProgramTradeEvent)) {
	c.dispatcher.OnUnifiedProgramTrade(h)
}
func (c *Client) OnNxtMember(h func(NxtMemberEvent))         { c.dispatcher.OnNxtMember(h) }
func (c *Client) OnUnifiedMember(h func(UnifiedMemberEvent)) { c.dispatcher.OnUnifiedMember(h) }

func (c *Client) OnConnected(h func())            { c.dispatcher.OnConnected(h) }
func (c *Client) OnReconnect(h func(attempt int)) { c.dispatcher.OnReconnect(h) }
func (c *Client) OnDisconnect(h func(error))      { c.dispatcher.OnDisconnect(h) }
func (c *Client) OnError(h func(error))           { c.dispatcher.OnError(h) }

// === Run loop ===

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
		if f.JSON.Body.RtCd != "" && f.JSON.Body.RtCd != "0" {
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
		if err != nil {
			c.dispatcher.RouteError(err)
			return
		}
		for _, ev := range evs {
			c.dispatcher.RouteKrxTrade(ev)
		}
	case trIDKrxAsk:
		evs, err := decodeKrxAsk(f)
		if err != nil {
			c.dispatcher.RouteError(err)
			return
		}
		for _, ev := range evs {
			c.dispatcher.RouteKrxAsk(ev)
		}
	case trIDKrxExpectTrade:
		evs, err := decodeKrxExpectTrade(f)
		if err != nil {
			c.dispatcher.RouteError(err)
			return
		}
		for _, ev := range evs {
			c.dispatcher.RouteKrxExpectTrade(ev)
		}
	case trIDKrxOvernightTrade:
		evs, err := decodeKrxOvernightTrade(f)
		if err != nil {
			c.dispatcher.RouteError(err)
			return
		}
		for _, ev := range evs {
			c.dispatcher.RouteKrxOvernightTrade(ev)
		}
	case trIDKrxOvernightExpect:
		evs, err := decodeKrxOvernightExpect(f)
		if err != nil {
			c.dispatcher.RouteError(err)
			return
		}
		for _, ev := range evs {
			c.dispatcher.RouteKrxOvernightExpect(ev)
		}

	// Phase 9 — NXT
	case trIDNxtTrade:
		evs, err := decodeAltMarketTrade(f)
		if err != nil {
			c.dispatcher.RouteError(err)
			return
		}
		for _, ev := range evs {
			c.dispatcher.RouteNxtTrade(ev)
		}
	case trIDNxtAsk:
		evs, err := decodeAltMarketAsk(f)
		if err != nil {
			c.dispatcher.RouteError(err)
			return
		}
		for _, ev := range evs {
			c.dispatcher.RouteNxtAsk(ev)
		}
	case trIDNxtExpectTrade:
		evs, err := decodeAltMarketExpectTrade(f)
		if err != nil {
			c.dispatcher.RouteError(err)
			return
		}
		for _, ev := range evs {
			c.dispatcher.RouteNxtExpectTrade(ev)
		}
	case trIDNxtProgramTrade:
		evs, err := decodeProgramTrade(f)
		if err != nil {
			c.dispatcher.RouteError(err)
			return
		}
		for _, ev := range evs {
			c.dispatcher.RouteNxtProgramTrade(ev)
		}
	case trIDNxtMember:
		evs, err := decodeMember(f)
		if err != nil {
			c.dispatcher.RouteError(err)
			return
		}
		for _, ev := range evs {
			c.dispatcher.RouteNxtMember(ev)
		}

	// Phase 9 — 통합
	case trIDUnifiedTrade:
		evs, err := decodeAltMarketTrade(f)
		if err != nil {
			c.dispatcher.RouteError(err)
			return
		}
		for _, ev := range evs {
			c.dispatcher.RouteUnifiedTrade(ev)
		}
	case trIDUnifiedAsk:
		evs, err := decodeAltMarketAsk(f)
		if err != nil {
			c.dispatcher.RouteError(err)
			return
		}
		for _, ev := range evs {
			c.dispatcher.RouteUnifiedAsk(ev)
		}
	case trIDUnifiedExpectTrade:
		evs, err := decodeAltMarketExpectTrade(f)
		if err != nil {
			c.dispatcher.RouteError(err)
			return
		}
		for _, ev := range evs {
			c.dispatcher.RouteUnifiedExpectTrade(ev)
		}
	case trIDUnifiedProgramTrade:
		evs, err := decodeProgramTrade(f)
		if err != nil {
			c.dispatcher.RouteError(err)
			return
		}
		for _, ev := range evs {
			c.dispatcher.RouteUnifiedProgramTrade(ev)
		}
	case trIDUnifiedMember:
		evs, err := decodeMember(f)
		if err != nil {
			c.dispatcher.RouteError(err)
			return
		}
		for _, ev := range evs {
			c.dispatcher.RouteUnifiedMember(ev)
		}

	default:
		c.dispatcher.RouteError(ErrWSInvalidFrame)
	}
}
