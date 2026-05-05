// examples/domestic_ksd/main.go — 예탁원 정보 11 메서드 사용 예시
package main

import (
	"context"
	"fmt"
	"log"

	kis "github.com/kenshin579/korea-investment-stock"
	"github.com/kenshin579/korea-investment-stock/domestic"
)

func main() {
	c, err := kis.NewClientFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()

	fromDate := "20260101"
	toDate := "20260505"
	symbol := "005930"

	// 1. 배당일정
	div, err := c.Domestic.InquireKsdDividend(ctx, domestic.InquireKsdDividendParams{
		FromDate: fromDate, ToDate: toDate,
	})
	if err != nil {
		log.Printf("InquireKsdDividend: %v", err)
	} else {
		fmt.Printf("InquireKsdDividend: %d rows\n", len(div.Output1))
		for i, item := range div.Output1 {
			if i >= 3 {
				break
			}
			fmt.Printf("  [%d] %s %s 배당금=%s\n", i, item.RecordDate, item.IsinName, item.PerStoDiviAmt)
		}
	}

	// 2. 무상증자
	bonus, err := c.Domestic.InquireKsdBonusIssue(ctx, domestic.InquireKsdBonusIssueParams{
		FromDate: fromDate, ToDate: toDate,
	})
	if err != nil {
		log.Printf("InquireKsdBonusIssue: %v", err)
	} else {
		fmt.Printf("InquireKsdBonusIssue: %d rows\n", len(bonus.Output1))
	}

	// 3. 유상증자
	paid, err := c.Domestic.InquireKsdPaidinCapin(ctx, domestic.InquireKsdPaidinCapinParams{
		FromDate: fromDate, ToDate: toDate,
	})
	if err != nil {
		log.Printf("InquireKsdPaidinCapin: %v", err)
	} else {
		fmt.Printf("InquireKsdPaidinCapin: %d rows\n", len(paid.Output)) // output (not output1)
	}

	// 4. 주주총회
	meet, err := c.Domestic.InquireKsdSharehldMeet(ctx, domestic.InquireKsdSharehldMeetParams{
		FromDate: fromDate, ToDate: toDate,
	})
	if err != nil {
		log.Printf("InquireKsdSharehldMeet: %v", err)
	} else {
		fmt.Printf("InquireKsdSharehldMeet: %d rows\n", len(meet.Output1))
	}

	// 5. 합병/분할
	merge, err := c.Domestic.InquireKsdMergerSplit(ctx, domestic.InquireKsdMergerSplitParams{
		FromDate: fromDate, ToDate: toDate,
	})
	if err != nil {
		log.Printf("InquireKsdMergerSplit: %v", err)
	} else {
		fmt.Printf("InquireKsdMergerSplit: %d rows\n", len(merge.Output1))
		for i, item := range merge.Output1 {
			if i >= 3 {
				break
			}
			fmt.Printf("  [%d] %s → %s (%s)\n", i, item.OppCustNm, item.CustNm, item.MergeType)
		}
	}

	// 6. 액면변경
	rev, err := c.Domestic.InquireKsdRevSplit(ctx, domestic.InquireKsdRevSplitParams{
		FromDate: fromDate, ToDate: toDate,
	})
	if err != nil {
		log.Printf("InquireKsdRevSplit: %v", err)
	} else {
		fmt.Printf("InquireKsdRevSplit: %d rows\n", len(rev.Output1))
	}

	// 7. 실권주청약
	forf, err := c.Domestic.InquireKsdForfeit(ctx, domestic.InquireKsdForfeitParams{
		Symbol: symbol, FromDate: fromDate, ToDate: toDate,
	})
	if err != nil {
		log.Printf("InquireKsdForfeit: %v", err)
	} else {
		fmt.Printf("InquireKsdForfeit: %d rows\n", len(forf.Output1))
	}

	// 8. 의무보호예수
	dep, err := c.Domestic.InquireKsdMandDeposit(ctx, domestic.InquireKsdMandDepositParams{
		Symbol: symbol, FromDate: fromDate, ToDate: toDate,
	})
	if err != nil {
		log.Printf("InquireKsdMandDeposit: %v", err)
	} else {
		fmt.Printf("InquireKsdMandDeposit: %d rows\n", len(dep.Output1))
		for i, item := range dep.Output1 {
			if i >= 3 {
				break
			}
			fmt.Printf("  [%d] depo_date=%s qty=%s reason=%s\n", i, item.DepoDate, item.StkQty, item.DepoReason)
		}
	}

	// 9. 감자
	cap, err := c.Domestic.InquireKsdCapDcrs(ctx, domestic.InquireKsdCapDcrsParams{
		FromDate: fromDate, ToDate: toDate,
	})
	if err != nil {
		log.Printf("InquireKsdCapDcrs: %v", err)
	} else {
		fmt.Printf("InquireKsdCapDcrs: %d rows\n", len(cap.Output1))
	}

	// 10. 주식매수청구
	pur, err := c.Domestic.InquireKsdPurreq(ctx, domestic.InquireKsdPurreqParams{
		FromDate: fromDate, ToDate: toDate,
	})
	if err != nil {
		log.Printf("InquireKsdPurreq: %v", err)
	} else {
		fmt.Printf("InquireKsdPurreq: %d rows\n", len(pur.Output1))
	}

	// 11. 주식상장정보
	lst, err := c.Domestic.InquireKsdListInfo(ctx, domestic.InquireKsdListInfoParams{
		Symbol: symbol, FromDate: fromDate, ToDate: toDate,
	})
	if err != nil {
		log.Printf("InquireKsdListInfo: %v", err)
	} else {
		fmt.Printf("InquireKsdListInfo: %d rows\n", len(lst.Output1))
		for i, item := range lst.Output1 {
			if i >= 3 {
				break
			}
			fmt.Printf("  [%d] list_dt=%s %s 발행가=%s\n", i, item.ListDt, item.IsinName, item.IssuePrice)
		}
	}
}
