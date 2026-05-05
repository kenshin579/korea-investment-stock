// overseas_chart example: InquireDailyPrice + InquireDailyChartPrice.
//
// Run: KIS credentials env vars 후 go run ./examples/overseas_chart
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	kis "github.com/kenshin579/korea-investment-stock"
	"github.com/kenshin579/korea-investment-stock/overseas"
)

func main() {
	client, err := kis.NewClientFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	today := time.Now().Format("20060102")

	// AAPL 일봉
	daily, err := client.Overseas.InquireDailyPrice(ctx, overseas.InquireDailyPriceParams{
		Excd: "NAS",
		Symb: "AAPL",
		Bymd: today,
	})
	if err != nil {
		log.Fatalf("InquireDailyPrice: %v", err)
	}
	fmt.Printf("[AAPL] 일봉 %d 캔들\n", len(daily.Output2))
	for i, c := range daily.Output2 {
		if i >= 5 {
			fmt.Println("  ... (이하 생략)")
			break
		}
		fmt.Printf("  %s: O=%s H=%s L=%s C=%s V=%d\n",
			c.Xymd, c.Open, c.High, c.Low, c.Clos, c.Tvol)
	}

	// S&P500 지수 (chart-price endpoint)
	from := time.Now().AddDate(0, 0, -30).Format("20060102")
	idx, err := client.Overseas.InquireDailyChartPrice(ctx, overseas.InquireDailyChartPriceParams{
		MarketCode: "N", // 해외지수
		Symbol:     "SPX",
		FromDate:   from,
		ToDate:     today,
		Period:     "D",
	})
	if err != nil {
		log.Fatalf("InquireDailyChartPrice: %v", err)
	}
	fmt.Printf("\n[%s %s] 30일 일봉 %d 개\n",
		idx.Output1.StckShrnIscd, idx.Output1.HtsKorIsnm, len(idx.Output2))
	for i, c := range idx.Output2 {
		if i >= 5 {
			fmt.Println("  ... (이하 생략)")
			break
		}
		fmt.Printf("  %s: O=%s H=%s L=%s C=%s V=%d\n",
			c.StckBsopDate, c.OvrsNmixOprc, c.OvrsNmixHgpr, c.OvrsNmixLwpr, c.OvrsNmixPrpr, c.AcmlVol)
	}
}
