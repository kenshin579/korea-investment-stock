// overseas_info example: InquireNewsTitle + InquireBrknewsTitle
// + InquireRightsByIce + InquirePeriodRights.
//
// Run: KIS credentials env vars 후 go run ./examples/overseas_info
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
		log.Fatalf("client init: %v", err)
	}
	ctx := context.Background()

	// EP1: InquireNewsTitle — 해외뉴스종합(제목)
	// ANOMALY: 응답 key outblock1
	news, err := client.Overseas.InquireNewsTitle(ctx, overseas.InquireNewsTitleParams{
		NationCd:   "US",
		ExchangeCd: "NAS",
		DataDt:     "20260505",
	})
	if err != nil {
		log.Printf("InquireNewsTitle error: %v", err)
	} else {
		fmt.Printf("=== InquireNewsTitle (outblock1 key) ===\n")
		for _, item := range news.Outblock1 {
			fmt.Printf("  [%s %s] %s — %s\n", item.DataDt, item.DataTm, item.Symb, item.Title)
		}
	}

	// EP2: InquireBrknewsTitle — 해외속보(제목)
	// ANOMALY: FID_ prefix params, iscd1-10/kor_isnm1-10 flat, hardcoded FID_COND_SCR_DIV_CODE
	brknews, err := client.Overseas.InquireBrknewsTitle(ctx, overseas.InquireBrknewsTitleParams{
		InputDate1: "20260505",
	})
	if err != nil {
		log.Printf("InquireBrknewsTitle error: %v", err)
	} else {
		fmt.Printf("\n=== InquireBrknewsTitle (FID_ prefix params) ===\n")
		for _, item := range brknews.Output {
			fmt.Printf("  [%s %s] %s (iscd1=%s, kor_isnm1=%s)\n",
				item.DataDt, item.DataTm, item.HtsPbntTitlCntt, item.Iscd1, item.KorIsnm1)
		}
	}

	// EP3: InquireRightsByIce — 해외주식_권리종합
	// ANOMALY: output1 only (no output2)
	rights, err := client.Overseas.InquireRightsByIce(ctx, overseas.InquireRightsByIceParams{
		NCod:  "US",
		Symb:  "AAPL",
		StYmd: "20260401",
		EdYmd: "20260430",
	})
	if err != nil {
		log.Printf("InquireRightsByIce error: %v", err)
	} else {
		fmt.Printf("\n=== InquireRightsByIce (output1 only) ===\n")
		for _, item := range rights.Output1 {
			fmt.Printf("  [%s] %s — pay_dt=%s\n", item.AnnoDt, item.CaTitle, item.PayDt)
		}
	}

	// EP4: InquirePeriodRights — 해외주식_기간별권리조회
	// ANOMALY: TR_ID C prefix, CTX_AREA_NK50/FK50 cursor pagination
	periodRights, err := client.Overseas.InquirePeriodRights(ctx, overseas.InquirePeriodRightsParams{
		InqrStrtDt: "20260401",
		InqrEndDt:  "20260430",
	})
	if err != nil {
		log.Printf("InquirePeriodRights error: %v", err)
	} else {
		fmt.Printf("\n=== InquirePeriodRights (CTRGT011R, CTX cursor) ===\n")
		for _, item := range periodRights.Output {
			fmt.Printf("  [%s] %s (%s) — cash_alct_rt=%s dfnt_yn=%s\n",
				item.BassDt, item.Pdno, item.PrdtName, item.CashAlctRt, item.DfntYn)
		}
	}
}
