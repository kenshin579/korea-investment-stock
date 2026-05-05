// overseas_price example: InquirePriceDetail + SearchInfo for AAPL.
//
// Run: KIS credentials env vars 후 go run ./examples/overseas_price
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

	price, err := client.Overseas.InquirePriceDetail(ctx, overseas.InquirePriceDetailParams{
		Excd: "NAS",
		Symb: "AAPL",
	})
	if err != nil {
		log.Fatalf("InquirePriceDetail: %v", err)
	}
	fmt.Printf("[AAPL @ NAS] 현재가 %s %s, PER=%v, PBR=%v\n",
		price.Output.Last, price.Output.Curr, price.Output.Perx, price.Output.Pbrx)
	fmt.Printf("  52주 최고/최저: %s / %s\n", price.Output.H52p, price.Output.L52p)
	fmt.Printf("  거래량: %d, 시가총액: %d\n", price.Output.Tvol, price.Output.Tomv)

	info, err := client.Overseas.SearchInfo(ctx, overseas.SearchInfoParams{
		PrdtTypeCD: "512", // NASDAQ
		Pdno:       "AAPL",
	})
	if err != nil {
		log.Fatalf("SearchInfo: %v", err)
	}
	fmt.Printf("\n상품정보: %s (%s)\n", info.Output.PrdtEngName, info.Output.PrdtName)
	fmt.Printf("  거래소: %s (%s), 통화: %s\n",
		info.Output.OvrsExcgName, info.Output.OvrsExcgCd, info.Output.TrCrcyCd)
	fmt.Printf("  ISIN: %s, SEDOL: %s, Bloomberg: %s\n",
		info.Output.StdPdno, info.Output.SedolNo, info.Output.BlbgTckrText)
}
