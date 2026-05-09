// examples/ws_index_futures_basic/main.go — Phase 11.3 지수선물옵션 + 상품선물 실시간 시연.
//
// 6 EP 중 4종 시연:
//   - SubscribeIndexFuturesTrade   (H0IFCNT0, 지수선물 체결)
//   - SubscribeIndexOptionTrade    (H0IOCNT0, 지수옵션 체결, 그릭스)
//   - SubscribeIndexOptionAsk      (H0IOASP0, 지수옵션 호가 5단계)
//   - SubscribeCommodityFuturesTrade (H0CFCNT0, 상품선물 체결, 지수선물과 alias)
//
// 4 base struct + 2 alias 패턴:
//
//	IndexFuturesTradeEvent ≡ CommodityFuturesTradeEvent
//	IndexFuturesAskEvent   ≡ CommodityFuturesAskEvent
//
// 사용자 facing API 는 시장 구분을 위해 별도 메서드, 내부 decoder 는 base 공유.
//
// 모든 EP 모의 미지원 — 실전 환경에서만 동작.
//
// Run: KIS 환경변수 설정 후 go run ./examples/ws_index_futures_basic
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
		fmt.Println("=== WebSocket 연결됨 (지수/상품 실시간 시연) ===")
	})
	client.WS.OnError(func(err error) {
		fmt.Printf(">>> ERROR: %v\n", err)
	})

	// 1) 지수선물 체결 (KOSPI200 등)
	client.WS.OnIndexFuturesTrade(func(ev websocket.IndexFuturesTradeEvent) {
		fmt.Printf("[지수선물 체결] %s %s @ %s (acc_vol=%d)\n",
			ev.Symbol, ev.Time, ev.Price.String(), ev.AccumVolume)
	})

	// 2) 지수옵션 체결 (그릭스 + IV/HV/AVRG_VLTL)
	client.WS.OnIndexOptionTrade(func(ev websocket.IndexOptionTradeEvent) {
		fmt.Printf("[지수옵션 체결] %s %s @ %s | IV=%.4f Delta=%.4f Gamma=%.4f AvgVol=%.4f\n",
			ev.Symbol, ev.Time, ev.Price.String(),
			ev.ImpliedVolatility, ev.Delta, ev.Gamma, ev.AvgVolatility)
	})

	// 3) 지수옵션 호가 (5단계)
	client.WS.OnIndexOptionAsk(func(ev websocket.IndexOptionAskEvent) {
		fmt.Printf("[지수옵션 호가] %s %s 매도1=%s(%d) 매수1=%s(%d)\n",
			ev.Symbol, ev.Time,
			ev.Ask[0].String(), ev.AskSize[0],
			ev.Bid[0].String(), ev.BidSize[0])
	})

	// 4) 상품선물 체결 (지수선물과 alias — 같은 struct, 다른 핸들러 슬롯)
	client.WS.OnCommodityFuturesTrade(func(ev websocket.CommodityFuturesTradeEvent) {
		fmt.Printf("[상품선물 체결] %s %s @ %s\n",
			ev.Symbol, ev.Time, ev.Price.String())
	})

	const indexFutures = "101S12"     // 예시 — 활성 종목 코드로 교체
	const indexOption = "201T12300"   // 예시
	const commodityFutures = "175S12" // 예시 (금, 원유 등)

	if err := client.WS.SubscribeIndexFuturesTrade(indexFutures); err != nil {
		log.Fatalf("SubscribeIndexFuturesTrade: %v", err)
	}
	if err := client.WS.SubscribeIndexOptionTrade(indexOption); err != nil {
		log.Fatalf("SubscribeIndexOptionTrade: %v", err)
	}
	if err := client.WS.SubscribeIndexOptionAsk(indexOption); err != nil {
		log.Fatalf("SubscribeIndexOptionAsk: %v", err)
	}
	if err := client.WS.SubscribeCommodityFuturesTrade(commodityFutures); err != nil {
		log.Fatalf("SubscribeCommodityFuturesTrade: %v", err)
	}

	fmt.Println("Ctrl+C 로 종료...")
	if err := client.WS.Run(ctx); err != nil && !errors.Is(err, context.Canceled) {
		log.Printf("Run: %v", err)
	}
	fmt.Println("종료")
}
