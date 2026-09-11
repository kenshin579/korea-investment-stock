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
// 성공 시 개수·tr_cont 만 로그한다(금액·종목 없음). 실패 시 SDK 에러에 응답 본문이 포함될 수 있다.
func TestIntegration_InquireBalance(t *testing.T) {
	if _, err := kis.LoadConfigFromEnv(); err != nil {
		t.Skipf("skip: %v", err)
	}
	c, err := kis.NewClientFromEnv()
	require.NoError(t, err)
	ctx := context.Background()

	dom, err := c.Domestic.InquireBalanceAll(ctx, domestic.InquireBalanceParams{})
	require.NoError(t, err)
	require.Len(t, dom.Output2, 1, "국내 요약 1행")
	t.Logf("domestic: %d holdings, last tr_cont=%q", len(dom.Output1), dom.TrCont)

	ovs, err := c.Overseas.InquireBalanceAll(ctx, overseas.InquireBalanceParams{OvrsExcgCd: "NASD", TrCrcyCd: "USD"})
	require.NoError(t, err)
	t.Logf("overseas(NASD/USD): %d holdings, last tr_cont=%q", len(ovs.Output1), ovs.TrCont)
}
