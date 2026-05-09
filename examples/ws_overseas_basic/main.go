// examples/ws_overseas_basic/main.go — Phase 10 해외주식 실시간 시세 시연.
//
// EP1: SubscribeOverseasTrade — HDFSCNT0 (체결가 지연, 26 fields)
// EP2: SubscribeOverseasAsk   — HDFSASP0 (호가, 17 fields, 1호가만)
//
// tr_key 형식:
//   - "D"+시장구분(3자리)+종목코드 (무료시세, 미국 0분지연 / 아시아 15분지연)
//   - "R"+시장구분(3자리)+종목코드 (유료시세 + 미국주간거래)
// 시장구분: NAS(나스닥), NYS(뉴욕), AMS(아멕스), TSE(도쿄), HKS(홍콩),
//          SHS(상해), SZS(심천), HSX(호치민), HNX(하노이),
//          BAY/BAQ/BAA(미국 주간 — 뉴욕/나스닥/아멕스).
//
// 모든 EP **모의 미지원** — 실전 환경에서만 동작.
//
// Run: KIS 환경변수 설정 후 go run ./examples/ws_overseas_basic
//
//	KOREA_INVESTMENT_APP_KEY=...
//	KOREA_INVESTMENT_APP_SECRET=...
//	KOREA_INVESTMENT_ACCOUNT_NO=...
//
// Ctrl+C 로 종료.
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	kis "github.com/kenshin579/korea-investment-stock"
	"github.com/kenshin579/korea-investment-stock/websocket"
)

func main() {
	client, err := kis.NewClientFromEnv()
	if err != nil {
		log.Fatalf("NewClientFromEnv: %v", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	client.WS.OnConnected(func() {
		fmt.Println("=== WebSocket 연결됨 (해외주식 시연) ===")
	})
	client.WS.OnReconnect(func(attempt int) {
		fmt.Printf(">>> 재연결 #%d\n", attempt)
	})
	client.WS.OnDisconnect(func(err error) {
		fmt.Printf(">>> 연결 끊김: %v\n", err)
	})
	client.WS.OnError(func(err error) {
		fmt.Printf(">>> ERROR: %v\n", err)
	})

	client.WS.OnOverseasTrade(func(ev websocket.OverseasTradeEvent) {
		fmt.Printf("[체결] %s %s 현재가=%s 등락율=%.2f%% 거래량=%d (체결강도=%.1f)\n",
			ev.Symbol, ev.LocalTime, ev.Last.String(), ev.ChangeRate,
			ev.TradeVolume, ev.TradeStrength)
	})

	client.WS.OnOverseasAsk(func(ev websocket.OverseasAskEvent) {
		fmt.Printf("[호가] %s %s 매수1=%s(%d) 매도1=%s(%d) | 총잔량 매수=%d 매도=%d\n",
			ev.Symbol, ev.LocalTime,
			ev.Bid1.String(), ev.Bid1Size,
			ev.Ask1.String(), ev.Ask1Size,
			ev.TotalBidSize, ev.TotalAskSize)
	})

	const symbol = "DNASAAPL" // 무료시세, NASDAQ AAPL

	if err := client.WS.SubscribeOverseasTrade(symbol); err != nil {
		log.Fatalf("SubscribeOverseasTrade: %v", err)
	}
	if err := client.WS.SubscribeOverseasAsk(symbol); err != nil {
		log.Fatalf("SubscribeOverseasAsk: %v", err)
	}

	fmt.Println("Ctrl+C 로 종료...")
	if err := client.WS.Run(ctx); err != nil && !errors.Is(err, context.Canceled) {
		log.Printf("Run: %v", err)
	}
	fmt.Println("종료")
}
