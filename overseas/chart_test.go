package overseas_test

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kenshin579/korea-investment-stock/overseas"
)

func TestClient_InquireDailyPrice(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/quotations/dailyprice`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "daily_price_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireDailyPrice(context.Background(), overseas.InquireDailyPriceParams{
		Excd: "NAS",
		Symb: "AAPL",
		Bymd: "20260505",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, "NAS", capturedQuery.Get("EXCD"))
	assert.Equal(t, "AAPL", capturedQuery.Get("SYMB"))
	assert.Equal(t, "0", capturedQuery.Get("GUBN")) // default 일
	assert.Equal(t, "0", capturedQuery.Get("MODP")) // default 미반영

	assert.Equal(t, "DNASAAPL", res.Output1.Rsym)
	require.Len(t, res.Output2, 2)
	assert.Equal(t, "20260505", res.Output2[0].Xymd)
	d, _ := decimal.NewFromString("181.45")
	assert.True(t, d.Equal(res.Output2[0].Clos))
	assert.Equal(t, int64(85000000), res.Output2[0].Tvol)
}

func TestClient_InquireDailyChartPrice(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/quotations/inquire-daily-chartprice`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "daily_chart_price_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireDailyChartPrice(context.Background(), overseas.InquireDailyChartPriceParams{
		MarketCode: "N",
		Symbol:     "SPX",
		FromDate:   "20260101",
		ToDate:     "20260505",
		Period:     "D",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, "N", capturedQuery.Get("FID_COND_MRKT_DIV_CODE"))
	assert.Equal(t, "SPX", capturedQuery.Get("FID_INPUT_ISCD"))
	assert.Equal(t, "20260101", capturedQuery.Get("FID_INPUT_DATE_1"))
	assert.Equal(t, "20260505", capturedQuery.Get("FID_INPUT_DATE_2"))
	assert.Equal(t, "D", capturedQuery.Get("FID_PERIOD_DIV_CODE"))

	assert.Equal(t, "S&P 500", res.Output1.HtsKorIsnm)
	require.Len(t, res.Output2, 2)
	assert.Equal(t, "20260505", res.Output2[0].StckBsopDate)
}
