//go:build integration

package kis_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	kis "github.com/kenshin579/korea-investment-stock"
	"github.com/kenshin579/korea-investment-stock/domestic"
	"github.com/kenshin579/korea-investment-stock/overseas"
)

// TestIntegration_InquireBalance 는 실계좌 잔고 실호출. KOREA_INVESTMENT_* env 필요.
// 실행: go test -tags integration -run TestIntegration_InquireBalance -v .
// 금액·종목은 로그에 남기지 않는다(개수만).
func TestIntegration_InquireBalance(t *testing.T) {
	c, err := kis.NewClientFromEnv()
	if err != nil {
		t.Skipf("skip: %v", err)
	}
	ctx := context.Background()

	dom, err := c.Domestic.InquireBalanceAll(ctx, domestic.InquireBalanceParams{})
	require.NoError(t, err)
	require.Len(t, dom.Output2, 1, "국내 요약 1행")
	t.Logf("domestic: %d holdings, last tr_cont=%q", len(dom.Output1), dom.TrCont)

	ovs, err := c.Overseas.InquireBalanceAll(ctx, overseas.InquireBalanceParams{OvrsExcgCd: "NASD", TrCrcyCd: "USD"})
	require.NoError(t, err)
	t.Logf("overseas(NASD/USD): %d holdings, last tr_cont=%q", len(ovs.Output1), ovs.TrCont)
}
