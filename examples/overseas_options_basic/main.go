// examples/overseas_options_basic/main.go — Phase 11.6 해외옵션 시세 + 장운영시간 시연.
//
// 4 메서드 통합:
//   - OverseasFutures.OptPrice (옵션 현재가)
//   - OverseasFutures.OptDetail (옵션 종목 상세)
//   - OverseasFutures.OptAskingPrice (옵션 호가)
//   - OverseasFutures.MarketTime (해외선물옵션 공통 장운영시간)
//
// 옵션 종목코드: 거래소/만기/행사가 형식 (예: ES 콜/풋 행사가별).
//
// 모든 EP 모의 미지원 — 실전 환경에서만 동작.
//
// Run: KIS 환경변수 설정 후 go run ./examples/overseas_options_basic
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

	const code = "ESM24C5000" // 예시: E-mini S&P500 6월 콜 5000 (실제 호출 시 활성 종목)

	// 1) 옵션 현재가
	if price, err := client.OverseasFutures.OptPrice(ctx, code); err != nil {
		log.Printf("OptPrice: %v", err)
	} else {
		fmt.Printf("[옵션 현재가] %s last=%s high=%s low=%s vol=%d\n",
			code,
			price.Output1.LastPrice.String(),
			price.Output1.HighPrice.String(),
			price.Output1.LowPrice.String(),
			price.Output1.Vol)
	}

	// 2) 옵션 종목 상세
	if detail, err := client.OverseasFutures.OptDetail(ctx, code); err != nil {
		log.Printf("OptDetail: %v", err)
	} else {
		fmt.Printf("[옵션 상세] %s output1=%+v\n", code, detail.Output1)
	}

	// 3) 옵션 호가
	if ask, err := client.OverseasFutures.OptAskingPrice(ctx, code); err != nil {
		log.Printf("OptAskingPrice: %v", err)
	} else {
		fmt.Printf("[옵션 호가] %s output1=%+v output2 size=%d\n", code, ask.Output1, len(ask.Output2))
	}

	// 4) 해외선물옵션 공통 장운영시간 (실제 사용 시 Params 정확 입력)
	// MarketTimeParams 의 인자 구조는 overseasfutures/market_time.go 의 struct 참조.
	fmt.Printf("[장운영시간] MarketTime 호출 예시 — overseasfutures/market_time.go 의 MarketTimeParams 참조\n")
}
