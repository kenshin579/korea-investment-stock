// examples/ws_overseas_futures_basic/main.go — Phase 11.7 해외선물옵션 실시간 시연.
//
// 2 EP 모두 시연:
//   - SubscribeOverseasFuturesTrade (HDFFF020, 체결가)
//   - SubscribeOverseasFuturesAsk   (HDFFF010, 호가, BID/ASK 5단계 교차 배열)
//
// 선물/옵션 통합 EP — 단일 TR_ID 로 선물/옵션 모두 수신. 그릭스 미포함.
//
// 종목코드 (SERIES_CD): 해외선물옵션 종목 series 코드 형식.
// Endpoint: 기존 인프라 (`ws://ops.koreainvestment.com:21000`) 그대로.
//
// 모든 EP 모의 미지원 — 실전 환경에서만 동작.
//
// Run: KIS 환경변수 설정 후 go run ./examples/ws_overseas_futures_basic
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
		fmt.Println("=== WebSocket 연결됨 (해외선물옵션 실시간 시연) ===")
	})
	client.WS.OnError(func(err error) {
		fmt.Printf(">>> ERROR: %v\n", err)
	})

	// 1) 체결가
	client.WS.OnOverseasFuturesTrade(func(ev websocket.OverseasFuturesTradeEvent) {
		fmt.Printf("[해외선물옵션 체결] %s %s%s @ %s (vol=%d, sign=%s)\n",
			ev.Symbol, ev.RecvDate, ev.RecvTime,
			ev.LastPrice.String(), ev.Vol, ev.QuotSign)
	})

	// 2) 호가 (BID/ASK 5단계)
	client.WS.OnOverseasFuturesAsk(func(ev websocket.OverseasFuturesAskEvent) {
		fmt.Printf("[해외선물옵션 호가] %s %s%s 매수1=%s(qty=%d) 매도1=%s(qty=%d)\n",
			ev.Symbol, ev.RecvDate, ev.RecvTime,
			ev.BidPrice[0].String(), ev.BidQntt[0],
			ev.AskPrice[0].String(), ev.AskQntt[0])
	})

	const series = "ESM24" // 예시: E-mini S&P500 6월물 (실제 호출 시 활성 series 코드)

	if err := client.WS.SubscribeOverseasFuturesTrade(series); err != nil {
		log.Fatalf("SubscribeOverseasFuturesTrade: %v", err)
	}
	if err := client.WS.SubscribeOverseasFuturesAsk(series); err != nil {
		log.Fatalf("SubscribeOverseasFuturesAsk: %v", err)
	}

	fmt.Println("Ctrl+C 로 종료...")
	if err := client.WS.Run(ctx); err != nil && !errors.Is(err, context.Canceled) {
		log.Printf("Run: %v", err)
	}
	fmt.Println("종료")
}
