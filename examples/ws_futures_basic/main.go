// examples/ws_futures_basic/main.go — Phase 11.2 국내선물옵션 실시간 시세 시연.
//
// 11 EP 중 4종 시연:
//   - SubscribeStockFuturesTrade  (H0ZFCNT0, 주식 선물 체결가)
//   - SubscribeStockFuturesAsk    (H0ZFASP0, 주식 선물 호가 — 10단계)
//   - SubscribeStockOptionTrade   (H0ZOCNT0, 주식 옵션 체결가 — 그릭스 포함)
//   - SubscribeKrxNightFuturesTrade (H0MFCNT0, KRX 야간 선물 체결가)
//
// 종목코드: 9자리 alphanumeric (예: 101W3000 선물, 201X3300 옵션, KAK0F 야간선물).
//
// 모든 EP 모의 미지원 — 실전 환경에서만 동작.
//
// Run: KIS 환경변수 설정 후 go run ./examples/ws_futures_basic
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
		fmt.Println("=== WebSocket 연결됨 (선물옵션 실시간 시연) ===")
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

	// 1) 주식 선물 체결가
	client.WS.OnStockFuturesTrade(func(ev websocket.StockFuturesTradeEvent) {
		fmt.Printf("[주식선물 체결] %s %s @ %s원 (acc_vol=%d)\n",
			ev.Symbol, ev.Time, ev.Price.String(), ev.AccumVolume)
	})

	// 2) 주식 선물 호가 (10단계)
	client.WS.OnStockFuturesAsk(func(ev websocket.StockFuturesAskEvent) {
		fmt.Printf("[주식선물 호가] %s %s 매도1=%s(%d) 매수1=%s(%d)\n",
			ev.Symbol, ev.Time,
			ev.Ask[0].String(), ev.AskSize[0],
			ev.Bid[0].String(), ev.BidSize[0])
	})

	// 3) 주식 옵션 체결가 (그릭스 포함)
	client.WS.OnStockOptionTrade(func(ev websocket.StockOptionTradeEvent) {
		fmt.Printf("[주식옵션 체결] %s %s @ %s | IV=%.4f Delta=%.4f Gamma=%.4f\n",
			ev.Symbol, ev.Time, ev.Price.String(),
			ev.ImpliedVolatility, ev.Delta, ev.Gamma)
	})

	// 4) KRX 야간 선물 체결가
	client.WS.OnKrxNightFuturesTrade(func(ev websocket.KrxNightFuturesTradeEvent) {
		fmt.Printf("[KRX야간선물 체결] %s %s @ %s원\n",
			ev.Symbol, ev.Time, ev.Price.String())
	})

	const stockFutures = "101W3000" // 예시 — 실제 활성 종목 코드로 교체
	const stockOption = "201X3300"  // 예시
	const krxNightFutures = "KAK0F" // 예시

	if err := client.WS.SubscribeStockFuturesTrade(stockFutures); err != nil {
		log.Fatalf("SubscribeStockFuturesTrade: %v", err)
	}
	if err := client.WS.SubscribeStockFuturesAsk(stockFutures); err != nil {
		log.Fatalf("SubscribeStockFuturesAsk: %v", err)
	}
	if err := client.WS.SubscribeStockOptionTrade(stockOption); err != nil {
		log.Fatalf("SubscribeStockOptionTrade: %v", err)
	}
	if err := client.WS.SubscribeKrxNightFuturesTrade(krxNightFutures); err != nil {
		log.Fatalf("SubscribeKrxNightFuturesTrade: %v", err)
	}

	fmt.Println("Ctrl+C 로 종료...")
	if err := client.WS.Run(ctx); err != nil && !errors.Is(err, context.Canceled) {
		log.Printf("Run: %v", err)
	}
	fmt.Println("종료")
}
