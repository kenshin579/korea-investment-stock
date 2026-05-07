// domestic_ranking example: InquireVolumeRank + InquireFluctuation + InquireMarketCap + InquireFinanceRatioRanking.
//
// Run: KIS credentials env vars 후 go run ./examples/domestic_ranking
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

	// 1. 거래량 상위 30
	vol, err := client.Domestic.InquireVolumeRank(ctx, domestic.InquireVolumeRankParams{
		InputISCD: "0000",
	})
	if err != nil {
		log.Fatalf("InquireVolumeRank: %v", err)
	}
	fmt.Printf("거래량 상위 %d 개\n", len(vol.Output))
	for i, item := range vol.Output {
		if i >= 5 {
			fmt.Println("  ... (이하 생략)")
			break
		}
		fmt.Printf("  %s. %s (%s) 현재가=%s 거래량=%d\n",
			item.DataRank, item.HtsKorIsnm, item.MkscShrnIscd,
			item.StckPrpr, item.AcmlVol)
	}

	// 2. 등락률 상위 (상승율순)
	flux, err := client.Domestic.InquireFluctuation(ctx, domestic.InquireFluctuationParams{
		InputISCD: "0000",
		SortCode:  "0", // 0=상승율순
	})
	if err != nil {
		log.Fatalf("InquireFluctuation: %v", err)
	}
	fmt.Printf("\n등락률 상위 (상승) %d 개\n", len(flux.Output))
	for i, item := range flux.Output {
		if i >= 5 {
			fmt.Println("  ... (이하 생략)")
			break
		}
		fmt.Printf("  %s. %s (%s) %s원 (%v%%)\n",
			item.DataRank, item.HtsKorIsnm, item.StckShrnIscd,
			item.StckPrpr, item.PrdyCtrt)
	}

	// 3. 시가총액 상위
	cap, err := client.Domestic.InquireMarketCap(ctx, domestic.InquireMarketCapParams{
		InputISCD: "0000",
	})
	if err != nil {
		log.Fatalf("InquireMarketCap: %v", err)
	}
	fmt.Printf("\n시가총액 상위 %d 개\n", len(cap.Output))
	for i, item := range cap.Output {
		if i >= 5 {
			fmt.Println("  ... (이하 생략)")
			break
		}
		fmt.Printf("  %s. %s (%s) 시총=%d백만원 비중=%v%%\n",
			item.DataRank, item.HtsKorIsnm, item.MkscShrnIscd,
			item.StckAvls, item.MrktWholAvlsRlim)
	}

	// 4. 재무비율 순위 (Phase 6) — 안정성 (BIS, 부채비율 등) 기준 상위 30
	finRank, err := client.Domestic.InquireFinanceRatioRanking(ctx, domestic.InquireFinanceRatioRankingParams{
		Year:     "2024",
		Period:   "3",  // 결산
		RankSort: "11", // 안정성
	})
	if err != nil {
		log.Fatalf("InquireFinanceRatioRanking: %v", err)
	}
	fmt.Printf("\n재무비율 순위 (안정성) %d 개\n", len(finRank.Output))
	for i, item := range finRank.Output {
		if i >= 5 {
			fmt.Println("  ... (이하 생략)")
			break
		}
		fmt.Printf("  %d. %s (%s) BIS=%v%% 부채비율=%v%% ROE=%v%%\n",
			item.DataRank, item.HtsKorIsnm, item.MkscShrnIscd,
			item.Bis, item.LbltRate, item.CptlNtinRate)
	}
}
