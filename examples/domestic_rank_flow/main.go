// domestic_rank_flow example: Phase 4.3 — 13 ranking/flow 메서드 시연.
//
// Run: KIS credentials env vars 후 go run ./examples/domestic_rank_flow
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
	today := time.Now().Format("20060102")
	weekAgo := time.Now().AddDate(0, 0, -7).Format("20060102")

	symbol := "005930" // 삼성전자

	// 1. InquireShortSale — 공매도 상위 (FHPST04820000)
	ss, err := client.Domestic.InquireShortSale(ctx, domestic.InquireShortSaleParams{
		Symbol:        symbol,
		PeriodDivCode: "D",
	})
	if err != nil {
		log.Printf("InquireShortSale 오류: %v", err)
	} else {
		fmt.Printf("[1] 공매도 상위 %d건\n", len(ss.Output))
		if len(ss.Output) > 0 {
			fmt.Printf("    %s (%s) 현재가=%s 공매도거래량=%d 비중=%.2f%%\n",
				ss.Output[0].HtsKorIsnm, ss.Output[0].MkscShrnIscd,
				ss.Output[0].StckPrpr, ss.Output[0].SstsCntgQty, ss.Output[0].SstsVolRlim)
		}
	}

	// 2. InquireDailyShortSale — 공매도 일별추이 (FHPST04830000)
	dss, err := client.Domestic.InquireDailyShortSale(ctx, domestic.InquireDailyShortSaleParams{
		Symbol:     symbol,
		InputDate1: weekAgo,
		InputDate2: today,
	})
	if err != nil {
		log.Printf("InquireDailyShortSale 오류: %v", err)
	} else {
		fmt.Printf("[2] 공매도 일별추이: 현재가=%s 거래량=%d 일별%d건\n",
			dss.Output1.StckPrpr, dss.Output1.AcmlVol, len(dss.Output2))
		if len(dss.Output2) > 0 {
			fmt.Printf("    최근일: %s 공매도비중=%.2f%%\n",
				dss.Output2[0].StckBsopDate, dss.Output2[0].AcmlSstsCntgQtyRlim)
		}
	}

	// 3. InquireCreditBalance — 신용잔고 상위 (FHKST17010000)
	cb, err := client.Domestic.InquireCreditBalance(ctx, domestic.InquireCreditBalanceParams{
		Symbol:       symbol,
		Option:       "5",
		RankSortCode: "0",
	})
	if err != nil {
		log.Printf("InquireCreditBalance 오류: %v", err)
	} else {
		fmt.Printf("[3] 신용잔고 상위 헤더%d건, 종목%d건\n", len(cb.Output1), len(cb.Output2))
		if len(cb.Output2) > 0 {
			fmt.Printf("    %s 신용융자잔고=%d주 잔고금액=%d\n",
				cb.Output2[0].HtsKorIsnm, cb.Output2[0].WholLoanRmndStcn, cb.Output2[0].WholLoanRmndAmt)
		}
	}

	// 4. InquireDailyCreditBalance — 신용잔고 일별추이 (FHPST04760000)
	dcb, err := client.Domestic.InquireDailyCreditBalance(ctx, domestic.InquireDailyCreditBalanceParams{
		Symbol:     symbol,
		InputDate1: today,
	})
	if err != nil {
		log.Printf("InquireDailyCreditBalance 오류: %v", err)
	} else {
		fmt.Printf("[4] 신용잔고 일별추이 %d건\n", len(dcb.Output))
		if len(dcb.Output) > 0 {
			fmt.Printf("    %s 신용융자잔고=%d주 담보비율=%.2f%%\n",
				dcb.Output[0].DealDate, dcb.Output[0].WholLoanRmndStcn, dcb.Output[0].WholLoanGvrt)
		}
	}

	// 5. InquireLendableByCompany — 당사 대주가능 (CTSC2702R)
	lbc, err := client.Domestic.InquireLendableByCompany(ctx, domestic.InquireLendableByCompanyParams{
		ExcgDvsnCd:     "02",
		Pdno:           symbol,
		ThcoStlnPsblYn: "Y",
		InqrDvsn1:      "1",
	})
	if err != nil {
		log.Printf("InquireLendableByCompany 오류: %v", err)
	} else {
		fmt.Printf("[5] 당사 대주가능 종목%d건 총한도=%d주 가능수량=%d주\n",
			len(lbc.Output1), lbc.Output2.TotStupLmtQty, lbc.Output2.RqstPsblQty)
		if len(lbc.Output1) > 0 {
			fmt.Printf("    %s (%s) 가능여부=%s\n",
				lbc.Output1[0].PrdtName, lbc.Output1[0].Pdno, lbc.Output1[0].PsblYnName)
		}
	}

	// 6. InquireQuoteBalance — 호가잔량 순위 (FHPST01720000)
	qb, err := client.Domestic.InquireQuoteBalance(ctx, domestic.InquireQuoteBalanceParams{
		Symbol:       "0000",
		VolCnt:       "30",
		RankSortCode: "0",
		DivClsCode:   "0",
		TrgtClsCode:  "111111111",
		TrgtExlsCode: "000000000",
	})
	if err != nil {
		log.Printf("InquireQuoteBalance 오류: %v", err)
	} else {
		fmt.Printf("[6] 호가잔량 순위 %d건\n", len(qb.Output))
		if len(qb.Output) > 0 {
			fmt.Printf("    %s. %s 매도잔량=%d 매수잔량=%d 매수비율=%.2f%%\n",
				qb.Output[0].DataRank, qb.Output[0].HtsKorIsnm,
				qb.Output[0].TotalAskpRsqn, qb.Output[0].TotalBidpRsqn, qb.Output[0].ShnuRsqnRate)
		}
	}

	// 7. InquireAfterHourBalance — 시간외잔량 순위 (FHPST01760000)
	ahb, err := client.Domestic.InquireAfterHourBalance(ctx, domestic.InquireAfterHourBalanceParams{
		Symbol:       "0000",
		RankSortCode: "0",
		DivClsCode:   "0",
	})
	if err != nil {
		log.Printf("InquireAfterHourBalance 오류: %v", err)
	} else {
		fmt.Printf("[7] 시간외잔량 순위 %d건\n", len(ahb.Output))
		if len(ahb.Output) > 0 {
			fmt.Printf("    %s. %s 시간외매도잔량=%d 시간외매수잔량=%d\n",
				ahb.Output[0].DataRank, ahb.Output[0].HtsKorIsnm,
				ahb.Output[0].OvtmTotalAskpRsqn, ahb.Output[0].OvtmTotalBidpRsqn)
		}
	}

	// 8. InquireOvertimeExpTransFluct — 시간외 예상체결 등락률 (FHKST11860000)
	// output 은 단일 객체 (배열 아님)
	oetf, err := client.Domestic.InquireOvertimeExpTransFluct(ctx, domestic.InquireOvertimeExpTransFluctParams{
		Symbol:       "0000",
		RankSortCode: "0",
		DivClsCode:   "0",
	})
	if err != nil {
		log.Printf("InquireOvertimeExpTransFluct 오류: %v", err)
	} else {
		fmt.Printf("[8] 시간외 예상체결 등락률 (단일객체): %s (%s) 예상체결가=%s 등락율=%.2f%%\n",
			oetf.Output.HtsKorIsnm, oetf.Output.StckShrnIscd,
			oetf.Output.OvtmUntpAntcCnpr, oetf.Output.OvtmUntpAntcCntgCtrt)
	}

	// 9. InquireMarketValue — 시장가치 순위 (FHPST01790000)
	mv, err := client.Domestic.InquireMarketValue(ctx, domestic.InquireMarketValueParams{
		Symbol:       "0000",
		RankSortCode: "0",
		DivClsCode:   "0",
		TrgtClsCode:  "111111111",
		TrgtExlsCode: "000000000",
		BlngClsCode:  "0",
		InputOption1: "2025",
		InputOption2: "1",
	})
	if err != nil {
		log.Printf("InquireMarketValue 오류: %v", err)
	} else {
		fmt.Printf("[9] 시장가치 순위 %d건\n", len(mv.Output))
		if len(mv.Output) > 0 {
			fmt.Printf("    %s. %s PER=%.2f PBR=%.2f PCR=%.2f PSR=%.2f\n",
				mv.Output[0].DataRank, mv.Output[0].HtsKorIsnm,
				mv.Output[0].Per, mv.Output[0].Pbr, mv.Output[0].Pcr, mv.Output[0].Psr)
		}
	}

	// 10. InquireDisparity — 이격도 순위 (FHPST01780000)
	disp, err := client.Domestic.InquireDisparity(ctx, domestic.InquireDisparityParams{
		Symbol:       "0000",
		HourClsCode:  "20",
		DivClsCode:   "0",
		RankSortCode: "0",
		TrgtClsCode:  "111111111",
		TrgtExlsCode: "000000000",
	})
	if err != nil {
		log.Printf("InquireDisparity 오류: %v", err)
	} else {
		fmt.Printf("[10] 이격도 순위 %d건\n", len(disp.Output))
		if len(disp.Output) > 0 {
			fmt.Printf("     %s. %s 5일=%.2f%% 20일=%.2f%% 60일=%.2f%%\n",
				disp.Output[0].DataRank, disp.Output[0].HtsKorIsnm,
				disp.Output[0].D5Dsrt, disp.Output[0].D20Dsrt, disp.Output[0].D60Dsrt)
		}
	}

	// 11. InquirePreferDisparateRatio — 우선주 괴리율 (FHPST01770000)
	pdr, err := client.Domestic.InquirePreferDisparateRatio(ctx, domestic.InquirePreferDisparateRatioParams{
		Symbol:       "0000",
		DivClsCode:   "0",
		TrgtClsCode:  "111111111",
		TrgtExlsCode: "000000000",
	})
	if err != nil {
		log.Printf("InquirePreferDisparateRatio 오류: %v", err)
	} else {
		fmt.Printf("[11] 우선주 괴리율 %d건\n", len(pdr.Output))
		if len(pdr.Output) > 0 {
			fmt.Printf("     %s (%s) ↔ %s (%s) 괴리율=%.2f%%\n",
				pdr.Output[0].HtsKorIsnm, pdr.Output[0].MkscShrnIscd,
				pdr.Output[0].PrstKorIsnm, pdr.Output[0].PrstIscd, pdr.Output[0].Dprt)
		}
	}

	// 12. InquireProfitAssetIndex — 수익자산지표 순위 (FHPST01730000)
	pai, err := client.Domestic.InquireProfitAssetIndex(ctx, domestic.InquireProfitAssetIndexParams{
		Symbol:       "0000",
		DivClsCode:   "0",
		TrgtClsCode:  "111111111",
		TrgtExlsCode: "000000000",
		RankSortCode: "0",
		BlngClsCode:  "0",
		InputOption1: "2025",
		InputOption2: "1",
	})
	if err != nil {
		log.Printf("InquireProfitAssetIndex 오류: %v", err)
	} else {
		fmt.Printf("[12] 수익자산지표 순위 %d건\n", len(pai.Output))
		if len(pai.Output) > 0 {
			fmt.Printf("     %s. %s 영업이익=%d 순이익=%d 자본총계=%d\n",
				pai.Output[0].DataRank, pai.Output[0].HtsKorIsnm,
				pai.Output[0].BsopPrti, pai.Output[0].ThtrNtin, pai.Output[0].TotalCptl)
		}
	}

	// 13. InquireMktfunds — 증시자금 종합 (FHKST649100C0)
	mf, err := client.Domestic.InquireMktfunds(ctx, domestic.InquireMktfundsParams{
		InputDate1: today,
	})
	if err != nil {
		log.Printf("InquireMktfunds 오류: %v", err)
	} else {
		fmt.Printf("[13] 증시자금 종합 %d건\n", len(mf.Output))
		if len(mf.Output) > 0 {
			fmt.Printf("     %s 지수=%s 고객예탁금=%d억원 신용융자=%d억원 MMF=%d억원\n",
				mf.Output[0].BsopDate, mf.Output[0].BstpNmixPrpr,
				mf.Output[0].CustDpmnAmt, mf.Output[0].CrdtLoanRmnd, mf.Output[0].MmfAmt)
		}
	}
}
