// domestic_financial example: InquireFinancialRatio + InquireIncomeStatement + InquireBalanceSheet + InquireOtherMajorRatios.
//
// Run: KIS credentials env vars 후 go run ./examples/domestic_financial
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
	symbol := "005930" // 삼성전자

	// 1. 재무비율 (연단위)
	ratio, err := client.Domestic.InquireFinancialRatio(ctx, domestic.InquireFinancialRatioParams{
		Symbol: symbol,
	})
	if err != nil {
		log.Fatalf("InquireFinancialRatio: %v", err)
	}
	fmt.Printf("[%s] 재무비율 %d 기간\n", symbol, len(ratio.Output))
	for i, item := range ratio.Output {
		if i >= 3 {
			fmt.Println("  ... (이하 생략)")
			break
		}
		fmt.Printf("  %s: 매출증가=%v%% 영업증가=%v%% ROE=%v%% EPS=%s BPS=%s\n",
			item.StacYymm, item.Grs, item.BsopPrfiInrt, item.RoeVal,
			item.Eps, item.Bps)
	}

	// 2. 손익계산서
	is, err := client.Domestic.InquireIncomeStatement(ctx, domestic.InquireIncomeStatementParams{
		Symbol: symbol,
	})
	if err != nil {
		log.Fatalf("InquireIncomeStatement: %v", err)
	}
	fmt.Printf("\n[%s] 손익계산서 %d 기간\n", symbol, len(is.Output))
	for i, item := range is.Output {
		if i >= 3 {
			fmt.Println("  ... (이하 생략)")
			break
		}
		fmt.Printf("  %s: 매출=%d 영업이익=%d 당기순이익=%d (백만원)\n",
			item.StacYymm, item.SaleAccount, item.BsopPrti, item.ThtrNtin)
	}

	// 3. 대차대조표
	bs, err := client.Domestic.InquireBalanceSheet(ctx, domestic.InquireBalanceSheetParams{
		Symbol: symbol,
	})
	if err != nil {
		log.Fatalf("InquireBalanceSheet: %v", err)
	}
	fmt.Printf("\n[%s] 대차대조표 %d 기간\n", symbol, len(bs.Output))
	for i, item := range bs.Output {
		if i >= 3 {
			fmt.Println("  ... (이하 생략)")
			break
		}
		fmt.Printf("  %s: 자산=%d 부채=%d 자본=%d (백만원)\n",
			item.StacYymm, item.TotalAset, item.TotalLblt, item.TotalCptl)
	}

	// 4. 기타주요비율 (Phase 6) — EVA / EBITDA / EV/EBITDA
	other, err := client.Domestic.InquireOtherMajorRatios(ctx, domestic.InquireOtherMajorRatiosParams{
		Symbol: symbol,
	})
	if err != nil {
		log.Fatalf("InquireOtherMajorRatios: %v", err)
	}
	fmt.Printf("\n[%s] 기타주요비율 %d 기간\n", symbol, len(other.Output))
	for i, item := range other.Output {
		if i >= 3 {
			fmt.Println("  ... (이하 생략)")
			break
		}
		// payout_rate 는 비정상 출력으로 무시
		fmt.Printf("  %s: EVA=%s EBITDA=%s EV/EBITDA=%v배\n",
			item.StacYymm, item.Eva, item.Ebitda, item.EvEbitda)
	}
}
