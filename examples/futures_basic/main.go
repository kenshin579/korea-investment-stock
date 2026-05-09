// futures_basic example: Phase 11.1 국내선물 시세 4 메서드 (EP1~EP4, 선택) 시연.
//
// Run:
//
//	export KOREA_INVESTMENT_API_KEY=...
//	export KOREA_INVESTMENT_API_SECRET=...
//	export KOREA_INVESTMENT_ACCOUNT_NO=...
//	go run ./examples/futures_basic
package main

import (
	"context"
	"fmt"
	"log"

	kis "github.com/kenshin579/korea-investment-stock"
	"github.com/kenshin579/korea-investment-stock/futures"
)

func main() {
	client, err := kis.NewClientFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()

	// 종목코드: 선물 KOSPI200 근월물 예시 코드
	// 실제 호출 시 caller가 정확한 활성 종목 코드 입력 (예: 101W3000, 101V06 등)
	code := "101W3000"
	marketCode := "F" // F:지수선물

	// 1. InquirePrice — 선물 현재가 시세 (EP1, FHMIF10000000)
	price, err := client.Futures.InquirePrice(ctx, marketCode, code)
	if err != nil {
		log.Printf("InquirePrice error: %v", err)
	} else {
		fmt.Printf("[EP1] InquirePrice: code=%s name=%s prpr=%s prdy_vrss=%s prdy_ctrt=%.2f%%\n",
			code, price.Output1.HtsKorIsnm, price.Output1.FutsPrpr.String(),
			price.Output1.FutsPrdyVrss.String(), price.Output1.FutsPrdyCtrt)
		fmt.Printf("      hi=%s lo=%s vol=%d basis=%s\n",
			price.Output1.FutsHgpr.String(), price.Output1.FutsLwpr.String(),
			price.Output1.AcmlVol, price.Output1.Basis.String())
	}

	// 2. InquireAskingPrice — 선물 호가 5단계 (EP2, FHMIF10010000)
	asking, err := client.Futures.InquireAskingPrice(ctx, marketCode, code)
	if err != nil {
		log.Printf("InquireAskingPrice error: %v", err)
	} else {
		fmt.Printf("[EP2] InquireAskingPrice: prpr=%s vol=%d\n",
			asking.Output1.FutsPrpr.String(), asking.Output1.AcmlVol)
		if len(asking.Output2) > 0 {
			ask := asking.Output2[0]
			fmt.Printf("      askp1=%s bidp1=%s total_ask=%d total_bid=%d\n",
				ask.FutsAskp1.String(), ask.FutsBidp1.String(),
				ask.TotalAskpRsqn, ask.TotalBidpRsqn)
		}
	}

	// 3. InquireDailyFuopchartprice — 기간별 시세(일봉) (EP3, FHKIF03020100)
	chartParams := futures.InquireDailyFuopchartpriceParams{
		MarketCode: marketCode,
		Code:       code,
		FromDate:   "20260401", // 조회 시작일
		ToDate:     "20260508", // 조회 종료일
		Period:     "D",        // D:일봉 (기본값)
	}
	chart, err := client.Futures.InquireDailyFuopchartprice(ctx, chartParams)
	if err != nil {
		log.Printf("InquireDailyFuopchartprice error: %v", err)
	} else {
		fmt.Printf("[EP3] InquireDailyFuopchartprice: %d records\n", len(chart.Output2))
		if len(chart.Output2) > 0 {
			first := chart.Output2[0]
			fmt.Printf("      first: date=%s open=%s high=%s low=%s close=%s vol=%d\n",
				first.StckBsopDate, first.FutsOprc.String(),
				first.FutsHgpr.String(), first.FutsLwpr.String(),
				first.FutsPrpr.String(), first.AcmlVol)
		}
	}

	// 4. DisplayBoardTop — 기초자산 전광판 조회 (EP4, FHPIF05030000)
	boardParams := futures.DisplayBoardTopParams{
		MarketCode: "F", // F:선물 (기본값)
		Code:       code,
	}
	board, err := client.Futures.DisplayBoardTop(ctx, boardParams)
	if err != nil {
		log.Printf("DisplayBoardTop error: %v", err)
	} else {
		fmt.Printf("[EP4] DisplayBoardTop: unas_prpr=%s futs_prpr=%s\n",
			board.Output1.UnasPrpr.String(), board.Output1.FutsPrpr.String())
		fmt.Printf("      unas_vrss=%s futs_vrss=%s maturity=%d\n",
			board.Output1.UnasPrdyVrss.String(), board.Output1.FutsPrdyVrss.String(),
			len(board.Output2))
	}
}
