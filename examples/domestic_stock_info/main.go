// domestic_stock_info example: Phase 4.1 — 국내주식 종목정보/분석 10 메서드.
//
// EP1: InquireInvestOpinion  — 종목투자의견 (FHKST663300C0)
// EP2: InquireInvestOpbysec  — 증권사별투자의견 (FHKST663400C0)
// EP3: InquireEstimatePerform — 종목추정실적 (HHKST668300C0) quad-output
// EP4: InquireVolumePower    — 체결강도상위 (FHPST01680000)
// EP5: InquireBulkTransNum   — 대량체결건수상위 (FHKST190900C0)
// EP6: InquireTradprtByamt   — 체결금액별매매비중 (FHKST111900C0)
// EP7: InquireHtsTopView     — HTS조회상위20종목 (HHMCM000100C0)
// EP8: InquirePbarTraRatio   — 매물대거래비중 (FHPST01130000)
// EP9: InquireExpPriceTrend  — 예상체결가추이 (FHPST01810000)
// EP10: InquireExpTransUpdown — 예상체결상승/하락상위 (FHPST01820000)
//
// Run: KIS credentials env vars 후 go run ./examples/domestic_stock_info
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

	// ── EP1: 종목투자의견 (FHKST663300C0) ───────────────────────────────────
	io, err := client.Domestic.InquireInvestOpinion(ctx, domestic.InquireInvestOpinionParams{
		Symbol:    symbol,
		StartDate: "20250101",
		EndDate:   "20251231",
	})
	if err != nil {
		log.Fatalf("InquireInvestOpinion: %v", err)
	}
	fmt.Printf("[EP1] 종목투자의견 %d 건\n", len(io.Output))
	for i, item := range io.Output {
		if i >= 3 {
			fmt.Println("  ... (이하 생략)")
			break
		}
		fmt.Printf("  %s | %s | 목표: %s원 | 회원사: %s\n",
			item.StckBsopDate, item.InvtOpnn, item.HtsGoalPrc, item.MbcrName)
	}

	// ── EP2: 증권사별투자의견 (FHKST663400C0) ───────────────────────────────
	ob, err := client.Domestic.InquireInvestOpbysec(ctx, domestic.InquireInvestOpbysecParams{
		SecBrokerCode: "0", // 전체
		DivClsCode:    "0", // 전체
		StartDate:     "20250101",
		EndDate:       "20251231",
	})
	if err != nil {
		log.Fatalf("InquireInvestOpbysec: %v", err)
	}
	fmt.Printf("\n[EP2] 증권사별투자의견 %d 건\n", len(ob.Output))
	for i, item := range ob.Output {
		if i >= 3 {
			fmt.Println("  ... (이하 생략)")
			break
		}
		fmt.Printf("  %s | %s (%s) | 의견: %s | 목표: %s원\n",
			item.StckBsopDate, item.HtsKorIsnm, item.StckShrnIscd, item.InvtOpnn, item.HtsGoalPrc)
	}

	// ── EP3: 종목추정실적 (HHKST668300C0) quad-output ───────────────────────
	ep, err := client.Domestic.InquireEstimatePerform(ctx, domestic.InquireEstimatePerformParams{
		Symbol: symbol,
	})
	if err != nil {
		log.Fatalf("InquireEstimatePerform: %v", err)
	}
	fmt.Printf("\n[EP3] 종목추정실적 — %s (%s)\n", ep.Output1.ItemKorNm, ep.Output1.ShtCd)
	fmt.Printf("  투자의견: %s | 결산일: %s\n", ep.Output1.RcmdName, ep.Output1.Estdate)
	fmt.Printf("  output2(손익계산서) %d 행, output3(투자지표) %d 행, output4(결산년월) %d 행\n",
		len(ep.Output2), len(ep.Output3), len(ep.Output4))
	if len(ep.Output4) > 0 {
		fmt.Printf("  결산년월: ")
		for i, p := range ep.Output4 {
			if i > 0 {
				fmt.Printf(", ")
			}
			fmt.Printf("%s", p.Dt)
		}
		fmt.Println()
	}

	// ── EP4: 체결강도상위 (FHPST01680000) ───────────────────────────────────
	vp, err := client.Domestic.InquireVolumePower(ctx, domestic.InquireVolumePowerParams{
		Symbol: "0000", // 전체
	})
	if err != nil {
		log.Fatalf("InquireVolumePower: %v", err)
	}
	fmt.Printf("\n[EP4] 체결강도상위 %d 건\n", len(vp.Output))
	for i, item := range vp.Output {
		if i >= 5 {
			fmt.Println("  ... (이하 생략)")
			break
		}
		fmt.Printf("  %s (%s): %s원 체결강도 %v\n",
			item.HtsKorIsnm, item.StckShrnIscd, item.StckPrpr, item.TdayRltv)
	}

	// ── EP5: 대량체결건수상위 (FHKST190900C0) ───────────────────────────────
	bt, err := client.Domestic.InquireBulkTransNum(ctx, domestic.InquireBulkTransNumParams{
		Symbol: "0000", // 전체
	})
	if err != nil {
		log.Fatalf("InquireBulkTransNum: %v", err)
	}
	fmt.Printf("\n[EP5] 대량체결건수상위 %d 건\n", len(bt.Output))
	for i, item := range bt.Output {
		if i >= 5 {
			fmt.Println("  ... (이하 생략)")
			break
		}
		fmt.Printf("  %s (%s): %s원 매수건수 %d 매도건수 %d\n",
			item.HtsKorIsnm, item.MkscShrnIscd, item.StckPrpr, item.ShnuCntgCsnu, item.SelnCntgCsnu)
	}

	// ── EP6: 체결금액별매매비중 (FHKST111900C0) ──────────────────────────────
	tb, err := client.Domestic.InquireTradprtByamt(ctx, domestic.InquireTradprtByamtParams{
		Symbol: symbol,
	})
	if err != nil {
		log.Fatalf("InquireTradprtByamt: %v", err)
	}
	fmt.Printf("\n[EP6] 체결금액별매매비중 %d 구간 — %s\n", len(tb.Output), symbol)
	for i, item := range tb.Output {
		if i >= 5 {
			fmt.Println("  ... (이하 생략)")
			break
		}
		fmt.Printf("  %s: 매수 %d건 매도 %d건 순매수비중 %v%%\n",
			item.PrprName, item.ShnuCntgCsnu, item.SelnCntgCsnu, item.WholNtbyQtyRate)
	}

	// ── EP7: HTS조회상위20종목 (HHMCM000100C0) zero params ──────────────────
	ht, err := client.Domestic.InquireHtsTopView(ctx, domestic.InquireHtsTopViewParams{})
	if err != nil {
		log.Fatalf("InquireHtsTopView: %v", err)
	}
	fmt.Printf("\n[EP7] HTS조회상위20종목\n")
	fmt.Printf("  시장구분: %s | 종목코드: %s\n",
		ht.Output1.MrktDivClsCode, ht.Output1.MkscShrnIscd)

	// ── EP8: 매물대거래비중 (FHPST01130000) ─────────────────────────────────
	pr, err := client.Domestic.InquirePbarTraRatio(ctx, domestic.InquirePbarTraRatioParams{
		Symbol: symbol,
	})
	if err != nil {
		log.Fatalf("InquirePbarTraRatio: %v", err)
	}
	fmt.Printf("\n[EP8] 매물대거래비중 — %s (%s)\n",
		pr.Output1.HtsKorIsnm, pr.Output1.StckShrnIscd)
	fmt.Printf("  현재가: %s원 | 가중평균: %s원 | 가격대 %d 건\n",
		pr.Output1.StckPrpr, pr.Output1.WghnAvrgStckPrc, len(pr.Output2))
	for i, item := range pr.Output2 {
		if i >= 3 {
			fmt.Println("  ... (이하 생략)")
			break
		}
		fmt.Printf("  rank %s | %s원 | %d주 | %v%%\n",
			item.DataRank, item.StckPrpr, item.CntgVol, item.AcmlVolRlim)
	}

	// ── EP9: 예상체결가추이 (FHPST01810000) ─────────────────────────────────
	et, err := client.Domestic.InquireExpPriceTrend(ctx, domestic.InquireExpPriceTrendParams{
		Symbol: symbol,
	})
	if err != nil {
		log.Fatalf("InquireExpPriceTrend: %v", err)
	}
	fmt.Printf("\n[EP9] 예상체결가추이 — %s\n", symbol)
	fmt.Printf("  예상체결가: %s원 (대비 %s, %v%%) 예상거래량: %d\n",
		et.Output1.AntcCnpr, et.Output1.AntcCntgVrss, et.Output1.AntcCntgPrdyCtrt, et.Output1.AntcVol)
	fmt.Printf("  시간별 추이 %d 건\n", len(et.Output2))
	for i, item := range et.Output2 {
		if i >= 3 {
			fmt.Println("  ... (이하 생략)")
			break
		}
		fmt.Printf("  %s %s | %s원 (%v%%)\n",
			item.StckBsopDate, item.StckCntgHour, item.StckPrpr, item.PrdyCtrt)
	}

	// ── EP10: 예상체결상승/하락상위 (FHPST01820000) ──────────────────────────
	eu, err := client.Domestic.InquireExpTransUpdown(ctx, domestic.InquireExpTransUpdownParams{
		Symbol:       "0000", // 전체
		RankSortCode: "0",    // 0:상승률
	})
	if err != nil {
		log.Fatalf("InquireExpTransUpdown: %v", err)
	}
	fmt.Printf("\n[EP10] 예상체결상승상위 %d 건\n", len(eu.Output))
	for i, item := range eu.Output {
		if i >= 5 {
			fmt.Println("  ... (이하 생략)")
			break
		}
		fmt.Printf("  %s (%s): %s원 (%v%%) 예상거래대금 %d\n",
			item.HtsKorIsnm, item.StckShrnIscd, item.StckPrpr, item.PrdyCtrt, item.AntcTrPbmn)
	}
}
