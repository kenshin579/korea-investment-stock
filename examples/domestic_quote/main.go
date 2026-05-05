// domestic_quote example: InquireAskingPriceExpCcn + InquireCcnl + InquireDailyPrice.
//
// Run: KIS credentials env vars 후 go run ./examples/domestic_quote
package main

import (
	"context"
	"fmt"
	"log"

	kis "github.com/kenshin579/korea-investment-stock"
	"github.com/kenshin579/korea-investment-stock/domestic"
)

func main() {
	client, err := kis.NewClientFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	symbol := "005930"

	// 1. 호가/예상체결
	ob, err := client.Domestic.InquireAskingPriceExpCcn(ctx, domestic.InquireAskingPriceExpCcnParams{
		Symbol: symbol,
	})
	if err != nil {
		log.Fatalf("InquireAskingPriceExpCcn: %v", err)
	}
	fmt.Printf("[%s] 호가 (접수시간 %s)\n", symbol, ob.Output1.AsprAcptHour)
	fmt.Printf("  매도1: %s @ %d, 매수1: %s @ %d\n",
		ob.Output1.Askp1, ob.Output1.AskpRsqn1, ob.Output1.Bidp1, ob.Output1.BidpRsqn1)
	fmt.Printf("  총 매도잔량 %d, 총 매수잔량 %d\n",
		ob.Output1.TotalAskpRsqn, ob.Output1.TotalBidpRsqn)
	fmt.Printf("  예상체결: %s (전일대비 %s, %v%%)\n",
		ob.Output2.AntcCnpr, ob.Output2.AntcCntgVrss, ob.Output2.AntcCntgPrdyCtrt)

	// 2. 최근 체결
	cc, err := client.Domestic.InquireCcnl(ctx, domestic.InquireCcnlParams{
		Symbol: symbol,
	})
	if err != nil {
		log.Fatalf("InquireCcnl: %v", err)
	}
	fmt.Printf("\n[%s] 최근 체결 %d 건\n", symbol, len(cc.Output))
	for i, item := range cc.Output {
		if i >= 5 {
			fmt.Println("  ... (이하 생략)")
			break
		}
		fmt.Printf("  %s: %s원 @ %d주 (체결강도 %v)\n",
			item.StckCntgHour, item.StckPrpr, item.CntgVol, item.TdayRltv)
	}

	// 3. 일자별
	dp, err := client.Domestic.InquireDailyPrice(ctx, domestic.InquireDailyPriceParams{
		Symbol: symbol,
		Period: "D",
	})
	if err != nil {
		log.Fatalf("InquireDailyPrice: %v", err)
	}
	fmt.Printf("\n[%s] 최근 일자별 %d 일\n", symbol, len(dp.Output))
	for i, item := range dp.Output {
		if i >= 5 {
			fmt.Println("  ... (이하 생략)")
			break
		}
		fmt.Printf("  %s: O=%s H=%s L=%s C=%s V=%d 외인소진=%v%%\n",
			item.StckBsopDate, item.StckOprc, item.StckHgpr, item.StckLwpr,
			item.StckClpr, item.AcmlVol, item.HtsFrgnEhrt)
	}
}
