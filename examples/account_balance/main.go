// account_balance example: 설정된 계좌의 국내·해외 보유 종목 전체 (InquireBalanceAll).
//
// Run: KOREA_INVESTMENT_API_KEY / API_SECRET / ACCOUNT_NO 설정 후 go run ./examples/account_balance
package main

import (
	"context"
	"fmt"
	"log"

	kis "github.com/kenshin579/korea-investment-stock"
	"github.com/kenshin579/korea-investment-stock/domestic"
	"github.com/kenshin579/korea-investment-stock/overseas"
)

func main() {
	client, err := kis.NewClientFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()

	// 1. 국내주식 잔고 — 계좌번호는 클라이언트 설정값, 나머지는 한투 기본값
	dom, err := client.Domestic.InquireBalanceAll(ctx, domestic.InquireBalanceParams{})
	if err != nil {
		log.Fatalf("Domestic.InquireBalanceAll: %v", err)
	}
	fmt.Printf("국내 보유 %d 종목\n", len(dom.Output1))
	for _, it := range dom.Output1 {
		fmt.Printf("  %s %-20s 수량=%d 평단=%.0f 평가=%d 손익=%d (%.2f%%)\n",
			it.Pdno, it.PrdtName, int64(it.HldgQty), float64(it.PchsAvgPric),
			int64(it.EvluAmt), int64(it.EvluPflsAmt), float64(it.EvluPflsRt))
	}
	if len(dom.Output2) > 0 {
		s := dom.Output2[0]
		fmt.Printf("  예수금=%d 유가평가=%d 총평가=%d\n", int64(s.DncaTotAmt), int64(s.SctsEvluAmt), int64(s.TotEvluAmt))
	}

	// 2. 해외주식 잔고 — 실전 미국 전체(NASD) + USD
	// 모의투자(WithPaperEnv)는 거래소 코드 체계가 다르다(NASD=나스닥, NAS 없음) — overseas.InquireBalanceParams 참고.
	ovs, err := client.Overseas.InquireBalanceAll(ctx, overseas.InquireBalanceParams{OvrsExcgCd: "NASD", TrCrcyCd: "USD"})
	if err != nil {
		log.Fatalf("Overseas.InquireBalanceAll: %v", err)
	}
	fmt.Printf("해외(미국) 보유 %d 종목\n", len(ovs.Output1))
	for _, it := range ovs.Output1 {
		fmt.Printf("  %-6s %-30s %s 수량=%.4f 평단=%.2f 평가=%.2f 손익=%.2f (%.2f%%)\n",
			it.OvrsPdno, it.OvrsItemName, it.OvrsExcgCd, float64(it.OvrsCblcQty), float64(it.PchsAvgPric),
			float64(it.OvrsStckEvluAmt), float64(it.FrcrEvluPflsAmt), float64(it.EvluPflsRt))
	}
	fmt.Printf("  외화매입합계=%.2f 총평가손익=%.2f (%.2f%%)\n",
		float64(ovs.Output2.FrcrPchsAmt1), float64(ovs.Output2.TotEvluPflsAmt), float64(ovs.Output2.TotPftrt))
}
