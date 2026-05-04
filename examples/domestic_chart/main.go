// domestic_chart example: InquireDailyItemChartPrice + InquireTimeItemChartPrice.
//
// Run: KIS credentials env vars 후 go run ./examples/domestic_chart
package main

import (
	"context"
	"fmt"
	"log"
	"time"

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

	// 일봉 (최근 약 30일)
	to := time.Now().Format("20060102")
	from := time.Now().AddDate(0, 0, -30).Format("20060102")
	daily, err := client.Domestic.InquireDailyItemChartPrice(ctx, domestic.InquireDailyItemChartPriceParams{
		Symbol:   symbol,
		Period:   "D",
		FromDate: from,
		ToDate:   to,
	})
	if err != nil {
		log.Fatalf("InquireDailyItemChartPrice: %v", err)
	}
	fmt.Printf("[%s %s] 일봉 %d 캔들\n", symbol, daily.Output1.HtsKorIsnm, len(daily.Output2))
	for i, c := range daily.Output2 {
		if i >= 5 {
			fmt.Println("  ... (이하 생략)")
			break
		}
		fmt.Printf("  %s: O=%s H=%s L=%s C=%s V=%d\n",
			c.StckBsopDate, c.StckOprc, c.StckHgpr, c.StckLwpr, c.StckClpr, c.AcmlVol)
	}

	// 분봉 (장 마감 직전부터 30개)
	minute, err := client.Domestic.InquireTimeItemChartPrice(ctx, domestic.InquireTimeItemChartPriceParams{
		Symbol:   symbol,
		TimeFrom: "153000",
	})
	if err != nil {
		log.Fatalf("InquireTimeItemChartPrice: %v", err)
	}
	fmt.Printf("\n[%s] 당일 분봉 %d 개 (15:30 시작):\n", symbol, len(minute.Output2))
	for i, c := range minute.Output2 {
		if i >= 5 {
			fmt.Println("  ... (이하 생략)")
			break
		}
		fmt.Printf("  %s %s: O=%s H=%s L=%s C=%s V=%d\n",
			c.StckBsopDate, c.StckCntgHour,
			c.StckOprc, c.StckHgpr, c.StckLwpr, c.StckPrpr, c.CntgVol)
	}
}
