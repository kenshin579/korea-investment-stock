// examples/ws_nxt_basic/main.go — Phase 9 NXT/통합 WebSocket 시세 시연.
//
// 신규 5종 EP × NXT/통합 = 10 EP. 본 예제는 그 중 4종을 시연:
//   - SubscribeNxtTrade           (H0NXCNT0, 체결가)
//   - SubscribeUnifiedAsk         (H0UNASP0, 호가, KMID/NMID 중간가 포함)
//   - SubscribeNxtProgramTrade    (H0NXPGM0, 신규 EP — 프로그램매매)
//   - SubscribeUnifiedMember      (H0UNMBC0, 신규 EP — 회원사 5단계)
//
// NXT 와 통합은 schema 동일 (5 base struct + 10 type alias). 따라서 핸들러 시그니처는
// 시장 구분 없이 base struct 의 alias 로 받음 — `NxtTradeEvent ≡ AltMarketTradeEvent`.
//
// 모든 EP **모의 미지원** — 실전 환경에서만 동작.
//
// Run: KIS 환경변수 설정 후 go run ./examples/ws_nxt_basic
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
		fmt.Println("=== WebSocket 연결됨 (NXT/통합 시연) ===")
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

	// 1) NXT 체결가 (KRX 와 같은 46 fields, 22번 = CNTG_CLS_CODE)
	client.WS.OnNxtTrade(func(ev websocket.NxtTradeEvent) {
		fmt.Printf("[NXT 체결] %s %s @ %s원 (vol=%d, kind=%s)\n",
			ev.Symbol, ev.Time, ev.Price.String(), ev.TradeVolume, ev.TradeKind)
	})

	// 2) 통합 호가 (65 fields, KRX/NXT 중간가 6 fields 포함)
	client.WS.OnUnifiedAsk(func(ev websocket.UnifiedAskEvent) {
		fmt.Printf("[통합 호가] %s %s 매도1=%s 매수1=%s | KMID=%s NMID=%s\n",
			ev.Symbol, ev.Time,
			ev.Ask[0].String(), ev.Bid[0].String(),
			ev.KrxMidPrice.String(), ev.NxtMidPrice.String())
	})

	// 3) NXT 프로그램매매 (신규 EP, 11 fields)
	client.WS.OnNxtProgramTrade(func(ev websocket.NxtProgramTradeEvent) {
		fmt.Printf("[NXT 프로그램] %s %s 매수=%d 매도=%d 순매수=%d\n",
			ev.Symbol, ev.Time, ev.BidQuantity, ev.AskQuantity, ev.NetQuantity)
	})

	// 4) 통합 회원사 (신규 EP, 78 fields — 5단계 매도/매수 + 외국계 + 영문)
	client.WS.OnUnifiedMember(func(ev websocket.UnifiedMemberEvent) {
		fmt.Printf("[통합 회원사] %s 매도1위=%s(%d주) 매수1위=%s(%d주) | 외국계 순매수=%d\n",
			ev.Symbol,
			ev.SellBrokerNames[0], ev.TotalSellQty[0],
			ev.BuyBrokerNames[0], ev.TotalBuyQty[0],
			ev.GlobalNetBuyQty)
	})

	const symbol = "005930" // 삼성전자

	if err := client.WS.SubscribeNxtTrade(symbol); err != nil {
		log.Fatalf("SubscribeNxtTrade: %v", err)
	}
	if err := client.WS.SubscribeUnifiedAsk(symbol); err != nil {
		log.Fatalf("SubscribeUnifiedAsk: %v", err)
	}
	if err := client.WS.SubscribeNxtProgramTrade(symbol); err != nil {
		log.Fatalf("SubscribeNxtProgramTrade: %v", err)
	}
	if err := client.WS.SubscribeUnifiedMember(symbol); err != nil {
		log.Fatalf("SubscribeUnifiedMember: %v", err)
	}

	fmt.Println("Ctrl+C 로 종료...")
	if err := client.WS.Run(ctx); err != nil && !errors.Is(err, context.Canceled) {
		log.Printf("Run: %v", err)
	}
	fmt.Println("종료")
}
