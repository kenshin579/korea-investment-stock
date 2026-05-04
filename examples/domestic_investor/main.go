// domestic_investor example: InquireInvestorTradeByStockDaily +
// InquireInvestorDailyByMarket + InquireIndexPrice + InquirePubOffer.
//
// Run: KIS credentials env vars 후 go run ./examples/domestic_investor
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
	yesterday := time.Now().AddDate(0, 0, -1).Format("20060102")
	symbol := "005930"

	// 1. 종목별 투자자매매동향 (일별, 어제 기준 N일치)
	stockDaily, err := client.Domestic.InquireInvestorTradeByStockDaily(ctx, domestic.InquireInvestorTradeByStockDailyParams{
		Symbol:   symbol,
		BaseDate: yesterday,
	})
	if err != nil {
		log.Fatalf("InquireInvestorTradeByStockDaily: %v", err)
	}
	fmt.Printf("[%s] 종목별 투자자매매동향 (어제 기준 %d일치)\n", symbol, len(stockDaily.Output2))
	for i, item := range stockDaily.Output2 {
		if i >= 5 {
			fmt.Println("  ... (이하 생략)")
			break
		}
		fmt.Printf("  %s: 외국인=%d, 개인=%d, 기관=%d (주)\n",
			item.StckBsopDate, item.FrgnNtbyQty, item.PrsnNtbyQty, item.OrgnNtbyQty)
	}

	// 2. 시장별 투자자매매동향 (일별, 코스피 종합)
	marketDaily, err := client.Domestic.InquireInvestorDailyByMarket(ctx, domestic.InquireInvestorDailyByMarketParams{
		Symbol:    "0001",
		BaseDate:  yesterday,
		Market:    "KSP",
		BaseDate2: yesterday,
		SubCode:   "0001",
	})
	if err != nil {
		log.Fatalf("InquireInvestorDailyByMarket: %v", err)
	}
	fmt.Printf("\n시장별 (코스피 종합) 투자자매매동향 %d 행\n", len(marketDaily.Output))
	for i, item := range marketDaily.Output {
		if i >= 3 {
			fmt.Println("  ... (이하 생략)")
			break
		}
		fmt.Printf("  %s: 지수=%s, 외국인=%d, 개인=%d (주)\n",
			item.StckBsopDate, item.BstpNmixPrpr, item.FrgnNtbyQty, item.PrsnNtbyQty)
	}

	// 3. 코스피 현재 지수
	idx, err := client.Domestic.InquireIndexPrice(ctx, domestic.InquireIndexPriceParams{
		Symbol: "0001",
	})
	if err != nil {
		log.Fatalf("InquireIndexPrice: %v", err)
	}
	fmt.Printf("\n코스피 현재 지수: %s (%v%%)\n", idx.Output.BstpNmixPrpr, idx.Output.BstpNmixPrdyCtrt)
	fmt.Printf("  상승=%s 하락=%s 보합=%s\n",
		idx.Output.AscnIssuCnt, idx.Output.DownIssuCnt, idx.Output.StnrIssuCnt)

	// 4. 공모주 청약일정 (다음 달)
	from := time.Now().Format("20060102")
	to := time.Now().AddDate(0, 1, 0).Format("20060102")
	ipo, err := client.Domestic.InquirePubOffer(ctx, domestic.InquirePubOfferParams{
		FromDate: from,
		ToDate:   to,
	})
	if err != nil {
		log.Fatalf("InquirePubOffer: %v", err)
	}
	fmt.Printf("\n공모주 청약일정 (%s ~ %s) %d 건\n", from, to, len(ipo.Output1))
	for i, item := range ipo.Output1 {
		if i >= 5 {
			fmt.Println("  ... (이하 생략)")
			break
		}
		fmt.Printf("  %s (%s) 공모가=%s, 청약일=%s, 주간사=%s\n",
			item.IsinName, item.ShtCd, item.FixSubscrPri, item.SubscrDt, item.LeadMgr)
	}
}
