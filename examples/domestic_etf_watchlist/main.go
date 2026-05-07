// domestic_etf_watchlist example: Phase 5 — 9 ETF/NAV/관심종목 메서드 시연.
//
// Run: KIS credentials env vars 후 go run ./examples/domestic_etf_watchlist
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	kis "github.com/kenshin579/korea-investment-stock"
	"github.com/kenshin579/korea-investment-stock/domestic"
)

func main() {
	client, err := kis.NewClientFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()

	etfSymbol := "069500" // KODEX 200

	today := time.Now().Format("20060102")
	weekAgo := time.Now().AddDate(0, 0, -7).Format("20060102")

	// 1. InquireEtfPrice — ETF/ETN 현재가 (FHPST02400000)
	ep, err := client.Domestic.InquireEtfPrice(ctx, domestic.InquireEtfPriceParams{
		Symbol: etfSymbol,
	})
	if err != nil {
		log.Printf("InquireEtfPrice 오류: %v", err)
	} else {
		fmt.Printf("[1] ETF 현재가 (KODEX 200): 현재가=%s NAV=%s 추적오차율=%.4f%% 괴리율=%.4f%%\n",
			ep.Output.StckPrpr, ep.Output.Nav, ep.Output.TrcErrt, ep.Output.Dprt)
	}

	// 2. InquireComponentStockPrice — ETF 구성종목 시세 (FHKST121600C0)
	csp, err := client.Domestic.InquireComponentStockPrice(ctx, domestic.InquireComponentStockPriceParams{
		Symbol: etfSymbol,
	})
	if err != nil {
		log.Printf("InquireComponentStockPrice 오류: %v", err)
	} else {
		fmt.Printf("[2] ETF 구성종목 시세: ETF현재가=%s NAV=%s 구성종목%d개\n",
			csp.Output1.StckPrpr, csp.Output1.Nav, len(csp.Output2))
		if len(csp.Output2) > 0 {
			fmt.Printf("    구성1: %s (%s) 현재가=%s 비중=%.4f%%\n",
				csp.Output2[0].HtsKorIsnm, csp.Output2[0].StckShrnIscd,
				csp.Output2[0].StckPrpr, csp.Output2[0].EtfCnfgIssuRlim)
		}
	}

	// 3. InquireNavComparisonTimeTrend — NAV 비교 시간 추이 (FHPST02440100)
	// fid_hour_cls_code: "60" = 60분 간격
	ntt, err := client.Domestic.InquireNavComparisonTimeTrend(ctx, domestic.InquireNavComparisonTimeTrendParams{
		Symbol:      etfSymbol,
		HourClsCode: "60",
	})
	if err != nil {
		log.Printf("InquireNavComparisonTimeTrend 오류: %v", err)
	} else {
		fmt.Printf("[3] NAV 비교 시간 추이 (60분 간격): %d건\n", len(ntt.Output))
		if len(ntt.Output) > 0 {
			fmt.Printf("    최근: %s NAV=%s 현재가=%s 괴리율=%.4f%%\n",
				ntt.Output[0].BsopHour, ntt.Output[0].Nav,
				ntt.Output[0].StckPrpr, ntt.Output[0].Dprt)
		}
	}

	// 4. InquireNavComparisonDailyTrend — NAV 비교 일별 추이 (FHPST02440200)
	ndt, err := client.Domestic.InquireNavComparisonDailyTrend(ctx, domestic.InquireNavComparisonDailyTrendParams{
		Symbol:     etfSymbol,
		InputDate1: weekAgo,
		InputDate2: today,
	})
	if err != nil {
		log.Printf("InquireNavComparisonDailyTrend 오류: %v", err)
	} else {
		fmt.Printf("[4] NAV 비교 일별 추이 (%s~%s): %d건\n", weekAgo, today, len(ndt.Output))
		if len(ndt.Output) > 0 {
			fmt.Printf("    최근일: %s 종가=%s NAV=%s 괴리율=%.4f%%\n",
				ndt.Output[0].StckBsopDate, ndt.Output[0].StckClpr,
				ndt.Output[0].Nav, ndt.Output[0].Dprt)
		}
	}

	// 5. InquireNavComparisonTrend — NAV 비교 추이 (FHPST02440000)
	nt, err := client.Domestic.InquireNavComparisonTrend(ctx, domestic.InquireNavComparisonTrendParams{
		Symbol: etfSymbol,
	})
	if err != nil {
		log.Printf("InquireNavComparisonTrend 오류: %v", err)
	} else {
		fmt.Printf("[5] NAV 비교 추이: 현재가=%s 시가=%s 고가=%s 저가=%s / NAV=%s 전일종가NAV=%s\n",
			nt.Output1.StckPrpr, nt.Output1.StckOprc, nt.Output1.StckHgpr, nt.Output1.StckLwpr,
			nt.Output2.Nav, nt.Output2.PrdyClprNav)
	}

	// 6. InquireIntstockMultprice — 관심종목 멀티 시세 (FHKST11300006)
	// 최대 30종목 batch; 여기서는 삼성전자+SK하이닉스 2종목
	mtp, err := client.Domestic.InquireIntstockMultprice(ctx, domestic.InquireIntstockMultpriceParams{
		MarketCodes: []string{"J", "J"},
		Symbols:     []string{"005930", "000660"}, // 삼성전자, SK하이닉스
	})
	if err != nil {
		log.Printf("InquireIntstockMultprice 오류: %v", err)
	} else {
		fmt.Printf("[6] 관심종목 멀티 시세 (배치): %s (%s) 현재가=%s 전일대비=%.2f%%\n",
			mtp.Output.InterKorIsnm, mtp.Output.InterShrnIscd,
			mtp.Output.Inter2Prpr, mtp.Output.PrdyCtrt)
	}

	// 7. InquireIntstockStocklistByGroup — 관심종목 그룹별 종목조회 (HHKCM113004C6)
	// USER_ID = HTS 로그인 ID (사용자 본인의 HTS 계정)
	userID := "YOUR_HTS_ID"
	grp, err := client.Domestic.InquireIntstockStocklistByGroup(ctx, domestic.InquireIntstockStocklistByGroupParams{
		UserID:       userID,
		InterGrpCode: "0001", // 관심 그룹 코드 (HTS 에서 확인)
	})
	if err != nil {
		log.Printf("InquireIntstockStocklistByGroup 오류 (HTS_ID 필요): %v", err)
	} else {
		fmt.Printf("[7] 관심종목 그룹별 종목조회: 그룹='%s' 종목%d개\n",
			grp.Output1.InterGrpName, len(grp.Output2))
		if len(grp.Output2) > 0 {
			fmt.Printf("    1위: %s (%s) 체결단가=%s\n",
				grp.Output2[0].HtsKorIsnm, grp.Output2[0].JongCode,
				grp.Output2[0].CntgUnpr)
		}
	}

	// 8. InquireIntstockGrouplist — 관심종목 그룹조회 (HHKCM113004C7)
	// USER_ID = HTS 로그인 ID (사용자 본인의 HTS 계정)
	gl, err := client.Domestic.InquireIntstockGrouplist(ctx, domestic.InquireIntstockGrouplistParams{
		UserID: userID,
	})
	if err != nil {
		log.Printf("InquireIntstockGrouplist 오류 (HTS_ID 필요): %v", err)
	} else {
		// NOTE: output2 만 존재 (output1 없음)
		fmt.Printf("[8] 관심종목 그룹조회 (output2 only): 그룹코드=%s 그룹명='%s' 종목수=%s\n",
			gl.Output2.InterGrpCode, gl.Output2.InterGrpName, gl.Output2.AskCnt)
	}

	// 9. InquireTopInterestStock — 관심종목등록 상위 (FHPST01800000)
	tis, err := client.Domestic.InquireTopInterestStock(ctx, domestic.InquireTopInterestStockParams{
		Symbol:     "0000", // 전체 조회
		DivClsCode: "0",
		InputCnt1:  "1",
	})
	if err != nil {
		log.Printf("InquireTopInterestStock 오류: %v", err)
	} else {
		fmt.Printf("[9] 관심종목등록 상위 %d건\n", len(tis.Output))
		if len(tis.Output) > 0 {
			fmt.Printf("    %s위: %s (%s) 현재가=%s 등록고객수=%d\n",
				tis.Output[0].DataRank, tis.Output[0].HtsKorIsnm,
				tis.Output[0].MkscShrnIscd, tis.Output[0].StckPrpr,
				tis.Output[0].InterIssuRegCsnu)
		}
	}
}
