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

	symbol := "0001" // KOSPI 지수

	// 1. InquireIndexDailyPrice — 국내업종 일자별지수 (FHPUP02120000)
	dailyResp, err := client.Domestic.InquireIndexDailyPrice(ctx, domestic.InquireIndexDailyPriceParams{
		Symbol:        symbol,
		MarketCode:    "U",
		PeriodDivCode: "D",
		InputDate1:    "20260101",
	})
	if err != nil {
		log.Printf("InquireIndexDailyPrice error: %v", err)
	} else {
		fmt.Printf("[EP3] InquireIndexDailyPrice: BstpNmixPrpr=%s, items=%d\n",
			dailyResp.Output1.BstpNmixPrpr, len(dailyResp.Output2))
	}

	// 2. InquireIndexTimeprice — 국내업종 시간별지수 분 (FHPUP02110200)
	timeResp, err := client.Domestic.InquireIndexTimeprice(ctx, domestic.InquireIndexTimepriceParams{
		Symbol:     symbol,
		MarketCode: "U",
	})
	if err != nil {
		log.Printf("InquireIndexTimeprice error: %v", err)
	} else {
		fmt.Printf("[EP4] InquireIndexTimeprice: items=%d\n", len(timeResp.Output))
	}

	// 3. InquireIndexTickprice — 국내업종 시간별지수 초 (FHPUP02110100)
	tickResp, err := client.Domestic.InquireIndexTickprice(ctx, domestic.InquireIndexTickpriceParams{
		Symbol:     symbol,
		MarketCode: "U",
	})
	if err != nil {
		log.Printf("InquireIndexTickprice error: %v", err)
	} else {
		fmt.Printf("[EP5] InquireIndexTickprice: items=%d\n", len(tickResp.Output))
	}

	// 4. InquireDailyIndexchartprice — 국내주식업종기간별시세 (FHKUP03500100)
	dailyChartResp, err := client.Domestic.InquireDailyIndexchartprice(ctx, domestic.InquireDailyIndexchartpriceParams{
		Symbol:        symbol,
		MarketCode:    "U",
		PeriodDivCode: "D",
		InputDate1:    "20260101",
		InputDate2:    "20260430",
	})
	if err != nil {
		log.Printf("InquireDailyIndexchartprice error: %v", err)
	} else {
		fmt.Printf("[EP6] InquireDailyIndexchartprice: BstpNmixPrpr=%s, items=%d\n",
			dailyChartResp.Output1.BstpNmixPrpr, len(dailyChartResp.Output2))
	}

	// 5. InquireTimeIndexchartprice — 업종 분봉조회 (FHKUP03500200)
	timeChartResp, err := client.Domestic.InquireTimeIndexchartprice(ctx, domestic.InquireTimeIndexchartpriceParams{
		Symbol:     symbol,
		MarketCode: "U",
	})
	if err != nil {
		log.Printf("InquireTimeIndexchartprice error: %v", err)
	} else {
		fmt.Printf("[EP7] InquireTimeIndexchartprice: BstpNmixPrpr=%s, items=%d\n",
			timeChartResp.Output1.BstpNmixPrpr, len(timeChartResp.Output2))
	}

	// 6. ExpTotalIndex — 예상체결 전체지수 (FHKUP11750000)
	expTotalResp, err := client.Domestic.ExpTotalIndex(ctx, domestic.ExpTotalIndexParams{
		MrktClsCode: "0",
		MarketCode:  "U",
		Symbol:      symbol,
		MkopClsCode: "0",
	})
	if err != nil {
		log.Printf("ExpTotalIndex error: %v", err)
	} else {
		fmt.Printf("[EP8] ExpTotalIndex: BstpNmixPrpr=%s, output2 items=%d\n",
			expTotalResp.Output1.BstpNmixPrpr, len(expTotalResp.Output2))
	}

	// 7. ExpIndexTrend — 예상체결지수 추이 (FHPST01840000)
	expTrendResp, err := client.Domestic.ExpIndexTrend(ctx, domestic.ExpIndexTrendParams{
		MkopClsCode: "1",
		InputHour1:  "10",
		Symbol:      symbol,
		MarketCode:  "U",
	})
	if err != nil {
		log.Printf("ExpIndexTrend error: %v", err)
	} else if len(expTrendResp.Output) > 0 {
		item := expTrendResp.Output[0]
		fmt.Printf("[EP9] ExpIndexTrend: StckCntgHour=%s, BstpNmixPrpr=%s, AcmlVol=%d\n",
			item.StckCntgHour, item.BstpNmixPrpr, item.AcmlVol)
	}
}
