// examples/ws_krx_basic/main.go — Phase 8 KRX WebSocket 시세 시연.
//
// EP1: SubscribeKrxTrade — H0STCNT0 (체결가, KRX)
// EP2: SubscribeKrxAsk   — H0STASP0 (호가, KRX)
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

	client.WS.OnKrxTrade(func(ev websocket.KrxTradeEvent) {
		fmt.Printf("[체결] %s %s @ %s원 (vol=%d, accum=%d)\n",
			ev.Symbol, ev.Time, ev.Price.String(), ev.TradeVolume, ev.AccumVolume)
	})
	client.WS.OnKrxAsk(func(ev websocket.KrxAskEvent) {
		fmt.Printf("[호가] %s %s 매도1=%s 매수1=%s\n",
			ev.Symbol, ev.Time, ev.Ask[0].String(), ev.Bid[0].String())
	})

	if err := client.WS.SubscribeKrxTrade("005930"); err != nil {
		log.Fatalf("SubscribeKrxTrade: %v", err)
	}
	if err := client.WS.SubscribeKrxAsk("005930"); err != nil {
		log.Fatalf("SubscribeKrxAsk: %v", err)
	}

	fmt.Println("Ctrl+C 로 종료...")
	if err := client.WS.Run(ctx); err != nil && !errors.Is(err, context.Canceled) {
		log.Printf("Run: %v", err)
	}
	fmt.Println("종료")
}
