// domestic_price example: InquirePrice + SearchInfo + SearchStockInfo.
//
// Run:
//   export KOREA_INVESTMENT_API_KEY=...
//   export KOREA_INVESTMENT_API_SECRET=...
//   export KOREA_INVESTMENT_ACCOUNT_NO=...
//   go run ./examples/domestic_price
package main

import (
	"context"
	"fmt"
	"log"

	kis "github.com/kenshin579/korea-investment-stock"
)

func main() {
	client, err := kis.NewClientFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	symbol := "005930" // 삼성전자

	price, err := client.Domestic.InquirePrice(ctx, symbol)
	if err != nil {
		log.Fatalf("InquirePrice: %v", err)
	}
	fmt.Printf("[%s] 현재가 %s, 전일대비 %s (%v%%)\n",
		symbol, price.StckPrpr.String(), price.PrdyVrss.String(), price.PrdyCtrt)
	fmt.Printf("  시가/고가/저가: %s / %s / %s\n",
		price.StckOprc.String(), price.StckHgpr.String(), price.StckLwpr.String())
	fmt.Printf("  거래량: %d, PER: %v, PBR: %v\n", price.AcmlVol, price.Per, price.Pbr)

	info, err := client.Domestic.SearchInfo(ctx, symbol, "300")
	if err != nil {
		log.Fatalf("SearchInfo: %v", err)
	}
	fmt.Printf("  상품명: %s (%s, %s)\n", info.PrdtName, info.PrdtClsfName, info.StdPdno)

	stockInfo, err := client.Domestic.SearchStockInfo(ctx, symbol, "300")
	if err != nil {
		log.Fatalf("SearchStockInfo: %v", err)
	}
	fmt.Printf("  시장: %s (%s), 업종: %s, KOSPI200=%s\n",
		stockInfo.MketIdCd, stockInfo.ScrtGrpIdCd,
		stockInfo.IdxBztpLclsCdName, stockInfo.Kospi200ItemYn)
}
