// examples/overseas_futures_basic/main.go — Phase 11.5 해외선물 시세 시연.
//
// 4 메서드 통합:
//   - OverseasFutures.InquirePrice (현재가)
//   - OverseasFutures.StockDetail (종목 상세)
//   - OverseasFutures.InquireAskingPrice (호가)
//   - OverseasFutures.DailyCcnl (일간 체결추이)
//
// 종목코드: 해외선물 종목 코드 형식 (예: ES, CL, GC — 거래소 + 만기 별도 query 인자).
//
// 모든 EP 모의 미지원 — 실전 환경에서만 동작.
//
// Run: KIS 환경변수 설정 후 go run ./examples/overseas_futures_basic
//
//	KOREA_INVESTMENT_APP_KEY=...
//	KOREA_INVESTMENT_APP_SECRET=...
//	KOREA_INVESTMENT_ACCOUNT_NO=...
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
		log.Fatalf("NewClientFromEnv: %v", err)
	}
	ctx := context.Background()

	const code = "ESM24" // 예시: E-mini S&P500 2024년 6월물 (실제 호출 시 활성 만기물로 교체)

	// 1) 현재가
	if price, err := client.OverseasFutures.InquirePrice(ctx, code); err != nil {
		log.Printf("InquirePrice: %v", err)
	} else {
		fmt.Printf("[현재가] %s last=%s high=%s low=%s vol=%d\n",
			code,
			price.Output1.LastPrice.String(),
			price.Output1.HighPrice.String(),
			price.Output1.LowPrice.String(),
			price.Output1.Vol)
	}

	// 2) 종목 상세
	if detail, err := client.OverseasFutures.StockDetail(ctx, code); err != nil {
		log.Printf("StockDetail: %v", err)
	} else {
		fmt.Printf("[상세] %s output1=%+v\n", code, detail.Output1)
	}

	// 3) 호가
	if ask, err := client.OverseasFutures.InquireAskingPrice(ctx, code); err != nil {
		log.Printf("InquireAskingPrice: %v", err)
	} else {
		fmt.Printf("[호가] %s output1=%+v output2 size=%d\n", code, ask.Output1, len(ask.Output2))
	}

	// 4) 일간 체결추이 (Params 인자 — _schemas.md EP5 참조)
	// 실제 사용 시 Params 의 정확한 필드명/Required 값은 overseasfutures/chart.go 의
	// DailyCcnlParams struct 참조. 본 예시는 컴파일 위주 시연.
	fmt.Printf("[일간 체결추이] DailyCcnl 호출 예시 — overseasfutures/chart.go 의 DailyCcnlParams struct 참조\n")
}
