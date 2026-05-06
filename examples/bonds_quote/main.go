// bonds_quote example: Phase 3.1 장내채권 시세 8 메서드 (EP1~EP8) 시연.
//
// Run:
//
//	export KOREA_INVESTMENT_API_KEY=...
//	export KOREA_INVESTMENT_API_SECRET=...
//	export KOREA_INVESTMENT_ACCOUNT_NO=...
//	go run ./examples/bonds_quote
package main

import (
	"context"
	"fmt"
	"log"

	kis "github.com/kenshin579/korea-investment-stock"
	"github.com/kenshin579/korea-investment-stock/bonds"
)

func main() {
	client, err := kis.NewClientFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()

	pdno := "KR103501GCC7" // 국고채권 03500-5103(21-3) 단축 종목코드

	// 1. SearchBondInfo — 채권 기본정보 조회 (EP1, CTPF1114R)
	info, err := client.Bonds.SearchBondInfo(ctx, bonds.SearchBondInfoParams{
		Pdno:       pdno,
		PrdtTypeCd: "300",
	})
	if err != nil {
		log.Printf("SearchBondInfo error: %v", err)
	} else {
		fmt.Printf("[EP1] SearchBondInfo: pdno=%s name=%s srfc_inrt=%s\n",
			info.Pdno, info.KsdBondItemName, info.KsdRcvgBondSrfcInrt)
	}

	// 2. InquireIssueInfo — 발행정보 조회 (EP2, CTPF1101R)
	issue, err := client.Bonds.InquireIssueInfo(ctx, bonds.InquireIssueInfoParams{
		Pdno:       pdno,
		PrdtTypeCd: "300",
	})
	if err != nil {
		log.Printf("InquireIssueInfo error: %v", err)
	} else {
		fmt.Printf("[EP2] InquireIssueInfo: pdno=%s prdt_name=%s issu_dt=%s kis_crdt=%s\n",
			issue.Pdno, issue.PrdtName, issue.IssuDt, issue.KisCrdtGradText)
	}

	// 3. InquirePrice — 현재가 시세 (EP3, FHKBJ773400C0)
	price, err := client.Bonds.InquirePrice(ctx, bonds.InquirePriceParams{
		MarketCode: "B",
		Symbol:     pdno,
	})
	if err != nil {
		log.Printf("InquirePrice error: %v", err)
	} else {
		fmt.Printf("[EP3] InquirePrice: bond_prpr=%s prdy_vrss=%s ernn_rate=%f\n",
			price.BondPrpr.String(), price.BondPrdyVrss.String(), price.ErnnRate)
	}

	// 4. InquireCcnl — 현재가 체결 (EP4, FHKBJ773403C0)
	ccnl, err := client.Bonds.InquireCcnl(ctx, bonds.InquireCcnlParams{
		MarketCode: "B",
		Symbol:     pdno,
	})
	if err != nil {
		log.Printf("InquireCcnl error: %v", err)
	} else {
		fmt.Printf("[EP4] InquireCcnl: stck_cntg_hour=%s bond_prpr=%s cntg_vol=%d\n",
			ccnl.StckCntgHour, ccnl.BondPrpr.String(), ccnl.CntgVol)
	}

	// 5. InquireAskingPrice — 현재가 호가 5단계 (EP5, FHKBJ773401C0)
	ask, err := client.Bonds.InquireAskingPrice(ctx, bonds.InquireAskingPriceParams{
		MarketCode: "B",
		Symbol:     pdno,
	})
	if err != nil {
		log.Printf("InquireAskingPrice error: %v", err)
	} else {
		fmt.Printf("[EP5] InquireAskingPrice: askp1=%s bidp1=%s total_ask=%d total_bid=%d\n",
			ask.BondAskp1.String(), ask.BondBidp1.String(),
			ask.TotalAskpRsqn, ask.TotalBidpRsqn)
	}

	// 6. InquireDailyPrice — 현재가 일별 (EP6, FHKBJ773404C0)
	daily, err := client.Bonds.InquireDailyPrice(ctx, bonds.InquireDailyPriceParams{
		MarketCode: "B",
		Symbol:     pdno,
	})
	if err != nil {
		log.Printf("InquireDailyPrice error: %v", err)
	} else {
		fmt.Printf("[EP6] InquireDailyPrice: stck_bsop_date=%s bond_prpr=%s acml_vol=%d\n",
			daily.StckBsopDate, daily.BondPrpr.String(), daily.AcmlVol)
	}

	// 7. InquireDailyItemchartprice — 기간별 시세 배열 (EP7, FHKBJ773701C0)
	chart, err := client.Bonds.InquireDailyItemchartprice(ctx, bonds.InquireDailyItemchartpriceParams{
		MarketCode: "B",
		Symbol:     pdno,
	})
	if err != nil {
		log.Printf("InquireDailyItemchartprice error: %v", err)
	} else {
		fmt.Printf("[EP7] InquireDailyItemchartprice: %d records\n", len(chart.Output))
		if len(chart.Output) > 0 {
			first := chart.Output[0]
			fmt.Printf("       first: date=%s open=%s high=%s low=%s close=%s vol=%d\n",
				first.StckBsopDate, first.BondOprc.String(),
				first.BondHgpr.String(), first.BondLwpr.String(),
				first.BondPrpr.String(), first.AcmlVol)
		}
	}

	// 8. InquireAvgUnit — 평균단가조회 output1+output2+output3 (EP8, CTPF2005R)
	avg, err := client.Bonds.InquireAvgUnit(ctx, bonds.InquireAvgUnitParams{
		InqrStrtDt:   "20260401",
		InqrEndDt:    "20260505",
		Pdno:         pdno,
		PrdtTypeCd:   "300",
		VrfcKindCd:   "01",
		CtxAreaNk30:  "",
		CtxAreaFk100: "",
	})
	if err != nil {
		log.Printf("InquireAvgUnit error: %v", err)
	} else {
		fmt.Printf("[EP8] InquireAvgUnit: output1=%d output2=%d output3=%d records\n",
			len(avg.Output1), len(avg.Output2), len(avg.Output3))
		if len(avg.Output1) > 0 {
			u := avg.Output1[0]
			fmt.Printf("       avg_evlu_unpr=%s avg_evlu_erng_rt=%f\n",
				u.AvgEvluUnpr.String(), u.AvgEvluErngRt)
		}
	}
}
