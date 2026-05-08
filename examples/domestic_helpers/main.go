// examples/domestic_helpers/main.go — Phase 7 헬퍼 (Group D) 사용 예제
//
// EP1: InquireMarketTime       — 국내선물 영업일조회      (HHMCM000002C0)
// EP2: InquireCompInterest     — 금리 종합 (국내채권 금리) (FHPST07020000)
// EP3: InquireTradedByCompany  — 당사매매종목 상위         (FHPST01860000)
// EP4: InquireCreditByCompany  — 당사 신용가능종목         (FHPST04770000)
//
// Run: KIS 환경변수 설정 후 go run ./examples/domestic_helpers
//
//	KOREA_INVESTMENT_APP_KEY=...
//	KOREA_INVESTMENT_APP_SECRET=...
//	KOREA_INVESTMENT_ACCOUNT_NO=...
//
// 모의투자 미지원 — 4개 EP 모두 실전 환경 only.
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
		log.Fatalf("NewClientFromEnv: %v", err)
	}
	ctx := context.Background()

	// ── EP1: 영업일조회 (HHMCM000002C0) ──────────────────────────────────────
	// 파라미터 없음. date1~date5(영업일) + today + time + s_time/e_time(장 시간).
	mt, err := client.Domestic.InquireMarketTime(ctx)
	if err != nil {
		log.Printf("[EP1] InquireMarketTime error: %v", err)
	} else if len(mt.Output1) > 0 {
		out := mt.Output1[0]
		fmt.Printf("=== [EP1] 영업일조회 ===\n")
		fmt.Printf("  영업일: %s %s [%s] %s %s (D±2, [당일])\n",
			out.Date1, out.Date2, out.Date3, out.Date4, out.Date5)
		fmt.Printf("  현재시간=%s 장시작=%s 장마감=%s\n", out.Time, out.STime, out.ETime)
	}

	// ── EP2: 금리 종합 (FHPST07020000) ───────────────────────────────────────
	// 파라미터 모두 hardcoded — Params 빈 struct.
	// output1=대표 metadata, output2=개별 금리 항목 array.
	ci, err := client.Domestic.InquireCompInterest(ctx, domestic.InquireCompInterestParams{})
	if err != nil {
		log.Printf("[EP2] InquireCompInterest error: %v", err)
	} else {
		fmt.Printf("\n=== [EP2] 금리 종합 — 대표: %s (%s) ===\n", ci.Output1.HtsKorIsnm, ci.Output1.BcdtCode)
		fmt.Printf("  현재금리=%s %% (전일대비 %s, %.2f%%)\n",
			ci.Output1.BondMnrtPrpr.String(),
			ci.Output1.BondMnrtPrdyVrss.String(),
			ci.Output1.PrdyCtrt)
		fmt.Printf("  개별 항목 %d 종\n", len(ci.Output2))
		for i, item := range ci.Output2 {
			if i >= 5 {
				fmt.Println("  ... (이하 생략)")
				break
			}
			fmt.Printf("    [%d] %s: %s%% (전일대비율 %.2f%%)\n",
				i+1, item.HtsKorIsnm, item.BondMnrtPrpr.String(), item.BstpNmixPrdyCtrt)
		}
	}

	// ── EP3: 당사매매종목 상위 (FHPST01860000) ───────────────────────────────
	// 매도/매수 누적 합계 + 순매수 (음수 가능). 최대 30 건.
	tbc, err := client.Domestic.InquireTradedByCompany(ctx, domestic.InquireTradedByCompanyParams{
		InputDate1: "20260501", // 기간 시작 YYYYMMDD
		InputDate2: "20260508", // 기간 종료 YYYYMMDD
		SortCode:   "1",        // 1=매수상위
	})
	if err != nil {
		log.Printf("[EP3] InquireTradedByCompany error: %v", err)
	} else {
		fmt.Printf("\n=== [EP3] 당사매매종목 상위 (매수상위 %d건) ===\n", len(tbc.Output))
		for i, item := range tbc.Output {
			if i >= 5 {
				fmt.Println("  ... (이하 생략)")
				break
			}
			fmt.Printf("  [%d] %s (%s) 현재가=%s 매수합=%d 순매수=%d\n",
				item.DataRank, item.HtsKorIsnm, item.MkscShrnIscd,
				item.StckPrpr.String(), item.ShnuCnqnSmtn, item.NtbyCnqn)
		}
	}

	// ── EP4: 당사 신용가능종목 (FHPST04770000) ───────────────────────────────
	// 신용주문 가능 종목 list + 신용 비율 (%). 최대 100 건.
	cbc, err := client.Domestic.InquireCreditByCompany(ctx, domestic.InquireCreditByCompanyParams{
		// 모두 default (코드순/신용주문가능/전체)
	})
	if err != nil {
		log.Printf("[EP4] InquireCreditByCompany error: %v", err)
	} else {
		fmt.Printf("\n=== [EP4] 당사 신용가능종목 (%d건) ===\n", len(cbc.Output))
		for i, item := range cbc.Output {
			if i >= 5 {
				fmt.Println("  ... (이하 생략)")
				break
			}
			fmt.Printf("  [%d] %s (%s) 신용비율=%.2f%%\n",
				i+1, item.HtsKorIsnm, item.StckShrnIscd, item.CrdtRate)
		}
	}
}
