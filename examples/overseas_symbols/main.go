// overseas_symbols example: FetchOverseasSymbols("nas") — NASDAQ 마스터 다운로드.
//
// Run: KIS credentials env vars 후 go run ./examples/overseas_symbols
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

	syms, err := client.Overseas.FetchOverseasSymbols(ctx, "nas")
	if err != nil {
		log.Fatalf("FetchOverseasSymbols: %v", err)
	}
	fmt.Printf("NASDAQ 종목 %d 개\n", len(syms))
	for i, s := range syms {
		if i >= 10 {
			fmt.Println("  ... (이하 생략)")
			break
		}
		fmt.Printf("  %s: %s\n", s.Symbol, s.EnglishName)
	}
}
