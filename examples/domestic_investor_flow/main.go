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

	// 1. 투자자 추정 (외인기관 가집계) — EP1 HHPTJ04160200
	fmt.Println("=== EP1: 투자자 매매 추정 가집계 ===")
	est, err := client.Domestic.InquireInvestorTrendEstimate(ctx, domestic.InquireInvestorTrendEstimateParams{
		Symbol: "005930",
	})
	if err != nil {
		log.Printf("EP1 error: %v", err)
	} else if len(est.Output2) > 0 {
		item := est.Output2[0]
		fmt.Printf("  외인 순매수(가집계): %d\n", item.FrgnFakeNtbyQty)
		fmt.Printf("  기관 순매수(가집계): %d\n", item.OrgnFakeNtbyQty)
		fmt.Printf("  합산 순매수(가집계): %d\n", item.SumFakeNtbyQty)
	}

	// 2. 외국인/기관 매매종목 집계 — EP2 FHPTJ04400000
	fmt.Println("\n=== EP2: 외인기관 매매종목가 집계 ===")
	fi, err := client.Domestic.InquireForeignInstitutionTotal(ctx, domestic.InquireForeignInstitutionTotalParams{
		Symbol:       "0001",
		DivClsCode:   "0",
		EtcClsCode:   "0",
		RankSortCode: "0",
	})
	if err != nil {
		log.Printf("EP2 error: %v", err)
	} else if len(fi.Output) > 0 {
		item := fi.Output[0]
		fmt.Printf("  종목명: %s\n", item.HtsKorIsnm)
		fmt.Printf("  외인 순매수: %d\n", item.FrgnNtbyQty)
		fmt.Printf("  기관 순매수: %d\n", item.OrgnNtbyQty)
	}

	// 3. 종목별 프로그램매매 추이(일별) — EP3 FHPPG04650201
	fmt.Println("\n=== EP3: 종목별 프로그램매매 추이(일별) ===")
	ptd, err := client.Domestic.InquireProgramTradeByStockDaily(ctx, domestic.InquireProgramTradeByStockDailyParams{
		MarketCode: "J",
		Symbol:     "005930",
		BaseDate:   "0020260505", // 002 prefix (KIS docs 예시)
	})
	if err != nil {
		log.Printf("EP3 error: %v", err)
	} else if len(ptd.Output) > 0 {
		item := ptd.Output[0]
		fmt.Printf("  영업일: %s\n", item.StckBsopDate)
		fmt.Printf("  전체 합산 순매수: %d\n", item.WholSmtnNtbyQty)
		fmt.Printf("  전체 합산 순매수 거래대금: %d\n", item.WholSmtnNtbyTrPbmn)
	}

	// 4. 종목별 프로그램매매 추이(체결) — EP4 FHPPG04650101
	fmt.Println("\n=== EP4: 종목별 프로그램매매 추이(체결) ===")
	pt, err := client.Domestic.InquireProgramTradeByStock(ctx, domestic.InquireProgramTradeByStockParams{
		MarketCode: "J",
		Symbol:     "005930",
	})
	if err != nil {
		log.Printf("EP4 error: %v", err)
	} else if len(pt.Output) > 0 {
		item := pt.Output[0]
		fmt.Printf("  영업시간: %s\n", item.BsopHour)
		fmt.Printf("  전체 합산 순매수: %d\n", item.WholSmtnNtbyQty)
		fmt.Printf("  전체 합산 매도 수량: %d\n", item.WholSmtnSelnVol)
	}

	// 5. 프로그램매매 종합현황(시간) — EP5 FHPPG04600101
	fmt.Println("\n=== EP5: 프로그램매매 종합현황(시간) ===")
	cpt, err := client.Domestic.InquireCompProgramTradeToday(ctx, domestic.InquireCompProgramTradeTodayParams{
		MarketCode:  "J",
		MrktClsCode: "K",
	})
	if err != nil {
		log.Printf("EP5 error: %v", err)
	} else if len(cpt.Output1) > 0 {
		item := cpt.Output1[0]
		fmt.Printf("  시간: %s\n", item.BsopHour)
		fmt.Printf("  차익 합산 순매수 거래대금: %d\n", item.ArbtSmtnNtbyTrPbmn)
		fmt.Printf("  비차익 합산 순매수 거래대금: %d\n", item.NabtSmtnNtbyTrPbmn)
	}

	// 6. 프로그램매매 종합현황(일별) — EP6 FHPPG04600001
	fmt.Println("\n=== EP6: 프로그램매매 종합현황(일별) ===")
	cpd, err := client.Domestic.InquireCompProgramTradeDaily(ctx, domestic.InquireCompProgramTradeDailyParams{
		MarketCode:  "J",
		MrktClsCode: "K",
		StartDate:   "20260101",
		EndDate:     "20260505",
	})
	if err != nil {
		log.Printf("EP6 error: %v", err)
	} else if len(cpd.Output) > 0 {
		item := cpd.Output[0]
		fmt.Printf("  영업일: %s\n", item.StckBsopDate)
		fmt.Printf("  비차익 위탁 매도 거래대금: %d\n", item.NabtEntmSelnTrPbmn)
		fmt.Printf("  차익 합계 매수 거래량: %d\n", item.ArbtSmtnShnuVol)
	}

	// 7. 당일 투자자별 프로그램매매 동향 — EP7 HHPPG046600C1
	fmt.Println("\n=== EP7: 당일 투자자별 프로그램매매 동향 ===")
	ipt, err := client.Domestic.InquireInvestorProgramTradeToday(ctx, domestic.InquireInvestorProgramTradeTodayParams{
		ExchDivClsCode: "J",
		MrktDivClsCode: "1",
	})
	if err != nil {
		log.Printf("EP7 error: %v", err)
	} else if len(ipt.Output1) > 0 {
		item := ipt.Output1[0]
		fmt.Printf("  투자자: %s (%s)\n", item.InvrClsName, item.InvrClsCode)
		fmt.Printf("  전체 순매수 금액: %d\n", item.AllNtbyAmt)
		fmt.Printf("  차익 순매수 금액: %d\n", item.ArbtNtbyAmt)
		fmt.Printf("  비차익 순매수 금액: %d\n", item.NabtNtbyAmt)
	}
}
