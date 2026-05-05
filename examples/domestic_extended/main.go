// domestic_extended example: InquireNearNewHighlow + InquireOvertimePrice +
// InquireOvertimeAskingPrice + InquireOvertimeVolume + InquireOvertimeFluctuation.
//
// Run: KIS credentials env vars 후 go run ./examples/domestic_extended
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
	symbol := "005930"

	// 1. 신고/신저근접종목 상위
	nh, err := client.Domestic.InquireNearNewHighlow(ctx, domestic.InquireNearNewHighlowParams{
		InputISCD:  "0000",
		PrcClsCode: "0", // 0:신고근접
	})
	if err != nil {
		log.Fatalf("InquireNearNewHighlow: %v", err)
	}
	fmt.Printf("[신고근접 상위 %d 건]\n", len(nh.Output))
	for i, item := range nh.Output {
		if i >= 5 {
			fmt.Println("  ... (이하 생략)")
			break
		}
		fmt.Printf("  %s (%s): %s원 (신고가 %s, 근접율 %v%%)\n",
			item.HtsKorIsnm, item.MkscShrnIscd, item.StckPrpr, item.NewHgpr, item.HprcNearRate)
	}

	// 2. 시간외현재가
	op, err := client.Domestic.InquireOvertimePrice(ctx, domestic.InquireOvertimePriceParams{
		Symbol: symbol,
	})
	if err != nil {
		log.Fatalf("InquireOvertimePrice: %v", err)
	}
	fmt.Printf("\n[%s] 시간외현재가\n", symbol)
	fmt.Printf("  현재가: %s원 (전일대비 %s, %v%%)\n",
		op.Output.OvtmUntpPrpr, op.Output.OvtmUntpPrdyVrss, op.Output.OvtmUntpPrdyCtrt)
	fmt.Printf("  거래량: %d주, 예상체결: %s원\n",
		op.Output.OvtmUntpVol, op.Output.OvtmUntpAntcCnpr)

	// 3. 시간외호가
	oa, err := client.Domestic.InquireOvertimeAskingPrice(ctx, domestic.InquireOvertimeAskingPriceParams{
		Symbol: symbol,
	})
	if err != nil {
		log.Fatalf("InquireOvertimeAskingPrice: %v", err)
	}
	fmt.Printf("\n[%s] 시간외호가 (최종시간 %s)\n", symbol, oa.Output1.OvtmUntpLastHour)
	fmt.Printf("  매도1: %s @ %d주, 매수1: %s @ %d주\n",
		oa.Output1.OvtmUntpAskp1, oa.Output1.OvtmUntpAskpRsqn1,
		oa.Output1.OvtmUntpBidp1, oa.Output1.OvtmUntpBidpRsqn1)
	fmt.Printf("  시간외 총매도잔량: %d, 총매수잔량: %d\n",
		oa.Output1.OvtmUntpTotalAskpRsqn, oa.Output1.OvtmUntpTotalBidpRsqn)

	// 4. 시간외거래량순위
	ov, err := client.Domestic.InquireOvertimeVolume(ctx, domestic.InquireOvertimeVolumeParams{
		InputISCD: "0000",
	})
	if err != nil {
		log.Fatalf("InquireOvertimeVolume: %v", err)
	}
	fmt.Printf("\n[시간외거래량순위 상위 %d 건]\n", len(ov.Output2))
	fmt.Printf("  거래소 합계: 거래량 %d주, 거래대금 %d원\n",
		ov.Output1.OvtmUntpExchVol, ov.Output1.OvtmUntpExchTrPbmn)
	for i, item := range ov.Output2 {
		if i >= 5 {
			fmt.Println("  ... (이하 생략)")
			break
		}
		fmt.Printf("  %s (%s): 시간외 %d주 (%v%%)\n",
			item.HtsKorIsnm, item.StckShrnIscd, item.OvtmUntpVol, item.OvtmVrssAcmlVolRlim)
	}

	// 5. 시간외등락율순위
	of, err := client.Domestic.InquireOvertimeFluctuation(ctx, domestic.InquireOvertimeFluctuationParams{
		InputISCD:  "0000",
		DivClsCode: "2", // 2:상승률
	})
	if err != nil {
		log.Fatalf("InquireOvertimeFluctuation: %v", err)
	}
	fmt.Printf("\n[시간외등락율순위 (상승률) 상위 %d 건]\n", len(of.Output2))
	fmt.Printf("  상한: %d, 상승: %d, 보합: %d, 하락: %d, 하한: %d\n",
		of.Output1.OvtmUntpUplmIssuCnt, of.Output1.OvtmUntpAscnIssuCnt,
		of.Output1.OvtmUntpStnrIssuCnt, of.Output1.OvtmUntpDownIssuCnt,
		of.Output1.OvtmUntpLslmIssuCnt)
	for i, item := range of.Output2 {
		if i >= 5 {
			fmt.Println("  ... (이하 생략)")
			break
		}
		fmt.Printf("  %s (%s): %s원 (시간외 %v%%)\n",
			item.HtsKorIsnm, item.MkscShrnIscd, item.OvtmUntpPrpr, item.OvtmUntpPrdyCtrt)
	}
}
