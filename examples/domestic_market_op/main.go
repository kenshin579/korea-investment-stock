// examples/domestic_market_op/main.go — Phase 4.2 시장운영/특수상태 사용 예제
//
// EP4: InquireExpClosingPrice  — 장마감 예상체결가    (FHKST117300C0)
// EP5: InquireChkHoliday       — 휴장일 조회          (CTCA0903R)
// EP6: InquireViStatus         — 변동성완화장치(VI) 현황 (FHPST01390000)
// EP7: InquireCaptureUplowprice — 상하한가 포착        (FHKST130000C0)
//
// Run: KIS 환경변수 설정 후 go run ./examples/domestic_market_op
//
//	KOREA_INVESTMENT_APP_KEY=...
//	KOREA_INVESTMENT_APP_SECRET=...
//	KOREA_INVESTMENT_ACCOUNT_NO=...
//	KOREA_INVESTMENT_MOCK=true  (모의투자 시)
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

	// ── EP4: 장마감 예상체결가 (FHKST117300C0) ──────────────────────────────
	// FID_INPUT_ISCD 는 종목코드가 아닌 시장구분코드:
	//   0000=전체 / 0001=코스피 / 1001=코스닥 / 2001=코스피200 / 4001=KRX100
	expClose, err := client.Domestic.InquireExpClosingPrice(ctx, domestic.InquireExpClosingPriceParams{
		RankSortClsCode: "0",    // 0=전체 / 1=상한가마감 / 2=하한가마감 / 3=상승률상위 / 4=하락률상위
		Symbol:          "0000", // 전체 시장
		BlngClsCode:     "0",    // 0=전체 / 1=종가범위연장
	})
	if err != nil {
		log.Printf("[EP4] InquireExpClosingPrice error: %v", err)
	} else {
		fmt.Printf("=== [EP4] 장마감 예상체결가 (상위 %d건) ===\n", len(expClose.Output1))
		for i, item := range expClose.Output1 {
			if i >= 3 {
				fmt.Println("  ... (이하 생략)")
				break
			}
			fmt.Printf("  [%d] %s (%s) 예상가=%s원 등락률=%.2f%% 체결량=%d\n",
				i+1, item.HtsKorIsnm, item.StckShrnIscd,
				item.StckPrpr.String(), item.PrdyCtrt, item.CntgVol)
		}
	}

	// ── EP5: 휴장일 조회 (CTCA0903R) ────────────────────────────────────────
	// 주의: 단시간 다수 호출 자제 (KIS docs 권장 1일 1회).
	// 파라미터명이 FID_ 접두어 없는 비표준 UPPERCASE 형식 (BASS_DT/CTX_AREA_NK/CTX_AREA_FK).
	holiday, err := client.Domestic.InquireChkHoliday(ctx, domestic.InquireChkHolidayParams{
		BassDt:    "20260507", // 조회기준일 YYYYMMDD
		CtxAreaNk: "",
		CtxAreaFk: "",
	})
	if err != nil {
		log.Printf("[EP5] InquireChkHoliday error: %v", err)
	} else if holiday.Output != nil {
		out := holiday.Output
		fmt.Printf("\n=== [EP5] 휴장일 조회 (%s) ===\n", out.Bassdt)
		fmt.Printf("  요일구분=%s 영업일=%s 거래일=%s 개장일=%s 결제일=%s\n",
			out.WdayDvsnCd, out.BzdyYn, out.TrDayYn, out.OpndYn, out.SttlDayYn)
	}

	// ── EP6: 변동성완화장치(VI) 현황 (FHPST01390000) ─────────────────────────
	// KIS 문서는 output 을 단일 Object로 선언.
	// 실 API 에서 배열 반환 시 []ViStatusOutput 로 struct 변경 필요.
	vi, err := client.Domestic.InquireViStatus(ctx, domestic.InquireViStatusParams{
		DivClsCode:      "0", // 0=전체 / 1=상승 / 2=하락
		MrktClsCode:     "0", // 0=전체 / K=거래소 / Q=코스닥
		Symbol:          "",  // 공란 = 전체
		RankSortClsCode: "0", // 0=전체 / 1=정적 / 2=동적 / 3=정적&동적
		InputDate1:      "20260507",
		TrgtClsCode:     "",
		TrgtExlsCode:    "",
	})
	if err != nil {
		log.Printf("[EP6] InquireViStatus error: %v", err)
	} else if vi.Output != nil {
		out := vi.Output
		fmt.Printf("\n=== [EP6] VI 현황 — %s (%s) ===\n", out.HtsKorIsnm, out.MkscShrnIscd)
		fmt.Printf("  VI발동=%s 종류=%s 발동가=%s 기준가=%s 발동시=%s 해제시=%s 횟수=%d\n",
			out.ViClsCode, out.ViKindCode,
			out.ViPrc.String(), out.ViStndPrc.String(),
			out.CntgViHour, out.ViCnclHour, out.ViCount)
	}

	// ── EP7: 상하한가 포착 (FHKST130000C0) ──────────────────────────────────
	uplow, err := client.Domestic.InquireCaptureUplowprice(ctx, domestic.InquireCaptureUplowpriceParams{
		PrcClsCode: "0",    // 0=상한가 / 1=하한가
		DivClsCode: "0",    // 0=상하한가 / 6=8%근접 / 5=10%근접 / 1=15%근접 / 2=20%근접 / 3=25%근접
		Symbol:     "0000", // 0000=전체 / 0001=코스피 / 1001=코스닥
	})
	if err != nil {
		log.Printf("[EP7] InquireCaptureUplowprice error: %v", err)
	} else {
		fmt.Printf("\n=== [EP7] 상한가 포착 (%d건) ===\n", len(uplow.Output))
		for i, item := range uplow.Output {
			if i >= 3 {
				fmt.Println("  ... (이하 생략)")
				break
			}
			fmt.Printf("  [%d] %s (%s) 현재가=%s원 등락률=%.2f%% 누적거래량=%d\n",
				i+1, item.HtsKorIsnm, item.MkscShrnIscd,
				item.StckPrpr.String(), item.PrdyCtrt, item.AcmlVol)
		}
	}
}
