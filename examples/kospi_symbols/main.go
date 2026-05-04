// kospi_symbols example: FetchKospiSymbols — KRX KOSPI 마스터 다운로드 + 통계.
//
// Run: KIS credentials env vars 후 go run ./examples/kospi_symbols
//
// 첫 실행 시 ~수 MB ZIP 다운로드 (디스크 캐시, default TTL 7일).
package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	kis "github.com/kenshin579/korea-investment-stock"
)

func main() {
	client, err := kis.NewClientFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()

	syms, err := client.Domestic.FetchKospiSymbols(ctx)
	if err != nil {
		log.Fatalf("FetchKospiSymbols: %v", err)
	}
	fmt.Printf("KOSPI 종목 %d 개\n", len(syms))

	// 그룹코드별 분포
	byGroup := make(map[string]int)
	for _, s := range syms {
		byGroup[s.GroupCode]++
	}
	fmt.Println("그룹별 분포:")
	for k, v := range byGroup {
		fmt.Printf("  %s: %d\n", k, v)
	}

	// 우선주만 추리기
	var preferred []string
	for _, s := range syms {
		if s.PreferredStock == "Y" {
			preferred = append(preferred, fmt.Sprintf("%s:%s", s.ShortCode, s.KoreanName))
		}
	}
	fmt.Printf("\n우선주 %d 개 (앞 10):\n  %s\n",
		len(preferred), strings.Join(preferred[:min(10, len(preferred))], ", "))

	// KOSPI200 편입 종목
	var kospi200 int
	for _, s := range syms {
		if s.KOSPI200 == "Y" {
			kospi200++
		}
	}
	fmt.Printf("\nKOSPI200 편입: %d 개\n", kospi200)
}
