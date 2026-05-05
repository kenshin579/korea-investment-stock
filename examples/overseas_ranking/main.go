// overseas_ranking example: InquireMarketCap + InquireTradeVol + InquireTradePbmn
// + InquireVolumeSurge + InquireVolumePower + InquireNewHighlow.
//
// Run: KIS credentials env vars 후 go run ./examples/overseas_ranking
package main

import (
	"context"
	"fmt"
	"log"

	kis "github.com/kenshin579/korea-investment-stock"
	"github.com/kenshin579/korea-investment-stock/overseas"
)

func main() {
	client, err := kis.NewClientFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	excd := "NAS"

	// 1. 시가총액순위
	mc, err := client.Overseas.InquireMarketCap(ctx, overseas.InquireMarketCapParams{
		ExcdCode: excd,
		VolRang:  "0",
	})
	if err != nil {
		log.Fatalf("InquireMarketCap: %v", err)
	}
	fmt.Printf("[%s 시가총액 상위 %d 건]\n", excd, len(mc.Output2))
	for i, item := range mc.Output2 {
		if i >= 5 {
			fmt.Println("  ... (이하 생략)")
			break
		}
		fmt.Printf("  #%d %s (%s): %s USD (시총 %s, 비중 %v%%)\n",
			item.Rank, item.Name, item.Symb, item.Last, item.Tomv, item.Grav)
	}

	// 2. 거래량순위
	tv, err := client.Overseas.InquireTradeVol(ctx, overseas.InquireTradeVolParams{
		ExcdCode: excd,
		NDay:     "0",
	})
	if err != nil {
		log.Fatalf("InquireTradeVol: %v", err)
	}
	fmt.Printf("\n[%s 거래량 상위 %d 건]\n", excd, len(tv.Output2))
	for i, item := range tv.Output2 {
		if i >= 5 {
			fmt.Println("  ... (이하 생략)")
			break
		}
		fmt.Printf("  #%d %s (%s): %d주 (평균 %d주)\n",
			item.Rank, item.Name, item.Symb, item.Tvol, item.ATvol)
	}

	// 3. 거래대금순위
	tp, err := client.Overseas.InquireTradePbmn(ctx, overseas.InquireTradePbmnParams{
		ExcdCode: excd,
		NDay:     "0",
	})
	if err != nil {
		log.Fatalf("InquireTradePbmn: %v", err)
	}
	fmt.Printf("\n[%s 거래대금 상위 %d 건]\n", excd, len(tp.Output2))
	for i, item := range tp.Output2 {
		if i >= 5 {
			fmt.Println("  ... (이하 생략)")
			break
		}
		fmt.Printf("  #%d %s (%s): %d USD (평균 %d USD)\n",
			item.Rank, item.Name, item.Symb, item.Tamt, item.ATamt)
	}

	// 4. 거래량급증
	vs, err := client.Overseas.InquireVolumeSurge(ctx, overseas.InquireVolumeSurgeParams{
		ExcdCode: excd,
		MixN:     "0", // 1분전
	})
	if err != nil {
		log.Fatalf("InquireVolumeSurge: %v", err)
	}
	fmt.Printf("\n[%s 거래량급증 %d 건]\n", excd, len(vs.Output2))
	for i, item := range vs.Output2 {
		if i >= 5 {
			fmt.Println("  ... (이하 생략)")
			break
		}
		fmt.Printf("  %s (%s): 현재 %d주 / 기준 %d주 (+%v, %v%%)\n",
			item.Knam, item.Symb, item.Tvol, item.NTvol, item.NDiff, item.NRate)
	}

	// 5. 매수체결강도상위
	vp, err := client.Overseas.InquireVolumePower(ctx, overseas.InquireVolumePowerParams{
		ExcdCode: excd,
		NDay:     "0", // 1분전 (wire name NDAY, 실제 분 단위)
	})
	if err != nil {
		log.Fatalf("InquireVolumePower: %v", err)
	}
	fmt.Printf("\n[%s 매수체결강도 상위 %d 건]\n", excd, len(vp.Output2))
	for i, item := range vp.Output2 {
		if i >= 5 {
			fmt.Println("  ... (이하 생략)")
			break
		}
		fmt.Printf("  %s (%s): 당일체결강도 %v / 체결강도 %v\n",
			item.Knam, item.Symb, item.Tpow, item.Powx)
	}

	// 6. 신고/신저가
	nh, err := client.Overseas.InquireNewHighlow(ctx, overseas.InquireNewHighlowParams{
		ExcdCode: excd,
		Gubn:     "1", // 신고(1)
		Gubn2:    "1", // 돌파유지(1)
		NDay:     "6", // 52주
	})
	if err != nil {
		log.Fatalf("InquireNewHighlow: %v", err)
	}
	fmt.Printf("\n[%s 52주 신고가 %d 건]\n", excd, len(nh.Output2))
	for i, item := range nh.Output2 {
		if i >= 5 {
			fmt.Println("  ... (이하 생략)")
			break
		}
		fmt.Printf("  %s (%s): %s USD (기준가 %s, 대비 %s, %v%%)\n",
			item.Name, item.Symb, item.Last, item.NBase, item.NDiff, item.NRate)
	}
}
