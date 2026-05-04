package domestic_test

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kenshin579/korea-investment-stock/domestic"
)

func TestClient_InquireInvestorTradeByStockDaily(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/quotations/investor-trade-by-stock-daily`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "investor_trade_by_stock_daily_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireInvestorTradeByStockDaily(context.Background(), domestic.InquireInvestorTradeByStockDailyParams{
		Symbol:   "005930",
		BaseDate: "20260505",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, "J", capturedQuery.Get("FID_COND_MRKT_DIV_CODE"))
	assert.Equal(t, "005930", capturedQuery.Get("FID_INPUT_ISCD"))
	assert.Equal(t, "20260505", capturedQuery.Get("FID_INPUT_DATE_1"))

	// output1 (요약)
	assert.Equal(t, decimal.NewFromInt(75800), res.Output1.StckPrpr)
	assert.Equal(t, "KOSPI200", res.Output1.RprsMrktKorName)
	assert.Equal(t, int64(12345678), res.Output1.AcmlVol)

	// output2 (Array, 일별 거래)
	require.Len(t, res.Output2, 2)
	assert.Equal(t, "20260505", res.Output2[0].StckBsopDate)
	assert.Equal(t, decimal.NewFromInt(75800), res.Output2[0].StckClpr)
	assert.Equal(t, int64(12345678), res.Output2[0].AcmlVol)
	assert.Equal(t, int64(-123456), res.Output2[0].FrgnNtbyQty)
	assert.Equal(t, int64(234567), res.Output2[0].PrsnNtbyQty)
	assert.Equal(t, int64(-100000), res.Output2[0].OrgnNtbyQty)
}
