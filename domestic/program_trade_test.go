// File: domestic/program_trade_test.go
package domestic_test

import (
	"context"
	"net/http"
	"net/url"
	"os"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kenshin579/korea-investment-stock/domestic"
)

func TestClient_InquireProgramTradeByStockDaily(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/quotations/program-trade-by-stock-daily`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "program_trade_by_stock_daily_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireProgramTradeByStockDaily(context.Background(), domestic.InquireProgramTradeByStockDailyParams{
		Symbol:   "005930",
		BaseDate: "0020260505",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, "J", capturedQuery.Get("FID_COND_MRKT_DIV_CODE"))
	assert.Equal(t, "005930", capturedQuery.Get("FID_INPUT_ISCD"))
	assert.Equal(t, "0020260505", capturedQuery.Get("FID_INPUT_DATE_1"))

	require.Len(t, res.Output, 2)
	assert.Equal(t, "20260505", res.Output[0].StckBsopDate)
	assert.Equal(t, decimal.NewFromInt(75800), res.Output[0].StckClpr)
	assert.Equal(t, int64(150000), res.Output[0].WholSmtnNtbyQty)
	assert.Equal(t, int64(10000), res.Output[0].WholNtbyVolIcdc)
	assert.Equal(t, int64(500000000), res.Output[0].WholNtbyTrPbmnIcdc2)
}

func TestClient_InquireProgramTradeByStock(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/quotations/program-trade-by-stock`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "program_trade_by_stock_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireProgramTradeByStock(context.Background(), domestic.InquireProgramTradeByStockParams{
		Symbol: "005930",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, "J", capturedQuery.Get("FID_COND_MRKT_DIV_CODE"))
	assert.Equal(t, "005930", capturedQuery.Get("FID_INPUT_ISCD"))

	require.Len(t, res.Output, 2)
	assert.Equal(t, "090100", res.Output[0].BsopHour)
	assert.Equal(t, decimal.NewFromInt(75800), res.Output[0].StckPrpr)
	assert.Equal(t, int64(1500), res.Output[0].WholSmtnNtbyQty)
	assert.Equal(t, int64(200), res.Output[0].WholNtbyVolIcdc)
	assert.Equal(t, int64(15160000), res.Output[0].WholNtbyTrPbmnIcdc)
}

func TestClient_InquireCompProgramTradeToday(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/quotations/comp-program-trade-today`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "comp_program_trade_today_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireCompProgramTradeToday(context.Background(), domestic.InquireCompProgramTradeTodayParams{
		MarketCode:  "J",
		MrktClsCode: "K",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, "J", capturedQuery.Get("FID_COND_MRKT_DIV_CODE"))
	assert.Equal(t, "K", capturedQuery.Get("FID_MRKT_CLS_CODE"))
	assert.Equal(t, "", capturedQuery.Get("FID_SCTN_CLS_CODE"))
	assert.Equal(t, "", capturedQuery.Get("FID_INPUT_ISCD"))

	require.Len(t, res.Output1, 2)
	assert.Equal(t, "090100", res.Output1[0].BsopHour)
	assert.Equal(t, int64(12000000000), res.Output1[0].ArbtSmtnSelnTrPbmn)
	assert.Equal(t, 35.50, res.Output1[0].ArbtSmtmSelnTrPbmnRate)
	assert.Equal(t, 36.10, res.Output1[0].ArbtSmtmShunTrPbmnRate)
	assert.True(t, decimal.NewFromFloat(2750.50).Equal(res.Output1[0].BstpNmixPrpr))
}

func TestClient_InquireCompProgramTradeDaily(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	raw, err := os.ReadFile("testdata/comp_program_trade_daily_success.json")
	require.NoError(t, err)
	httpmock.RegisterResponder("GET",
		"=~.*comp-program-trade-daily.*",
		httpmock.NewBytesResponder(200, raw))

	c := newTestClient(t)
	ctx := context.Background()
	res, err := c.InquireCompProgramTradeDaily(ctx, domestic.InquireCompProgramTradeDailyParams{
		MarketCode:  "J",
		MrktClsCode: "K",
		StartDate:   "20260101",
		EndDate:     "20260505",
	})
	require.NoError(t, err)
	require.NotEmpty(t, res.Output)

	item := res.Output[0]
	assert.NotEmpty(t, item.StckBsopDate)                          // string
	assert.GreaterOrEqual(t, item.NabtEntmSelnTrPbmn, int64(0))    // int64
	assert.GreaterOrEqual(t, item.NabtEntmSelnVolRate, float64(0)) // float64
	assert.GreaterOrEqual(t, item.ArbtSmtmShunVolRate, float64(0)) // float64 (shun typo)
}

func TestClient_InquireInvestorProgramTradeToday(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	raw, err := os.ReadFile("testdata/investor_program_trade_today_success.json")
	require.NoError(t, err)
	httpmock.RegisterResponder("GET",
		"=~.*investor-program-trade-today.*",
		httpmock.NewBytesResponder(200, raw))

	c := newTestClient(t)
	ctx := context.Background()
	res, err := c.InquireInvestorProgramTradeToday(ctx, domestic.InquireInvestorProgramTradeTodayParams{
		ExchDivClsCode: "J",
		MrktDivClsCode: "1",
	})
	require.NoError(t, err)
	require.NotEmpty(t, res.Output1)

	item := res.Output1[0]
	assert.NotEmpty(t, item.InvrClsCode)
	assert.NotEmpty(t, item.InvrClsName)
	assert.GreaterOrEqual(t, item.AllSelnQty, int64(0))
	assert.GreaterOrEqual(t, item.ArbtNtbyAmt, int64(0))
	assert.GreaterOrEqual(t, item.NabtNtbyAmt, int64(0))
}
