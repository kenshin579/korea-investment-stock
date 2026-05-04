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

func TestClient_InquireDailyItemChartPrice(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/quotations/inquire-daily-itemchartprice`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "daily_chart_success.json")), nil
		},
	)

	c := newTestClient(t)
	chart, err := c.InquireDailyItemChartPrice(context.Background(), domestic.InquireDailyItemChartPriceParams{
		Symbol:   "005930",
		FromDate: "20260430",
		ToDate:   "20260502",
	})
	require.NoError(t, err)
	require.NotNil(t, chart)

	// query default 검증 (zero-value → "D" / "J" / 수정주가)
	assert.Equal(t, "J", capturedQuery.Get("FID_COND_MRKT_DIV_CODE"))
	assert.Equal(t, "005930", capturedQuery.Get("FID_INPUT_ISCD"))
	assert.Equal(t, "20260430", capturedQuery.Get("FID_INPUT_DATE_1"))
	assert.Equal(t, "20260502", capturedQuery.Get("FID_INPUT_DATE_2"))
	assert.Equal(t, "D", capturedQuery.Get("FID_PERIOD_DIV_CODE"))
	assert.Equal(t, "0", capturedQuery.Get("FID_ORG_ADJ_PRC")) // 0 = 수정주가

	// output1 검증
	assert.Equal(t, "삼성전자", chart.Output1.HtsKorIsnm)
	assert.Equal(t, decimal.NewFromInt(75800), chart.Output1.StckPrpr)
	assert.Equal(t, "005930", chart.Output1.StckShrnIscd)

	// output2 검증
	require.Len(t, chart.Output2, 3)
	assert.Equal(t, "20260502", chart.Output2[0].StckBsopDate)
	assert.Equal(t, decimal.NewFromInt(75800), chart.Output2[0].StckClpr)
	assert.Equal(t, decimal.NewFromInt(76000), chart.Output2[0].StckOprc)
	assert.Equal(t, int64(12345678), chart.Output2[0].AcmlVol)
}

func TestClient_InquireDailyItemChartPrice_OriginalPrice(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/quotations/inquire-daily-itemchartprice`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "daily_chart_success.json")), nil
		},
	)

	c := newTestClient(t)
	_, err := c.InquireDailyItemChartPrice(context.Background(), domestic.InquireDailyItemChartPriceParams{
		Symbol:        "005930",
		Period:        "W",
		FromDate:      "20260101",
		ToDate:        "20260502",
		OriginalPrice: true,
		MarketCode:    "NX",
	})
	require.NoError(t, err)
	assert.Equal(t, "NX", capturedQuery.Get("FID_COND_MRKT_DIV_CODE"))
	assert.Equal(t, "W", capturedQuery.Get("FID_PERIOD_DIV_CODE"))
	assert.Equal(t, "1", capturedQuery.Get("FID_ORG_ADJ_PRC")) // 1 = 원주가
}

func TestClient_InquireTimeItemChartPrice(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/quotations/inquire-time-itemchartprice`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "minute_chart_success.json")), nil
		},
	)

	c := newTestClient(t)
	chart, err := c.InquireTimeItemChartPrice(context.Background(), domestic.InquireTimeItemChartPriceParams{
		Symbol:   "005930",
		TimeFrom: "150000",
	})
	require.NoError(t, err)
	require.NotNil(t, chart)

	assert.Equal(t, "J", capturedQuery.Get("FID_COND_MRKT_DIV_CODE"))
	assert.Equal(t, "005930", capturedQuery.Get("FID_INPUT_ISCD"))
	assert.Equal(t, "150000", capturedQuery.Get("FID_INPUT_HOUR_1"))
	assert.Equal(t, "N", capturedQuery.Get("FID_PW_DATA_INCU_YN"))

	require.Len(t, chart.Output2, 3)
	assert.Equal(t, "150000", chart.Output2[0].StckCntgHour)
	assert.Equal(t, decimal.NewFromInt(75800), chart.Output2[0].StckPrpr)
	assert.Equal(t, int64(12345), chart.Output2[0].CntgVol)
}

func TestClient_InquireTimeItemChartPrice_PastDataInclude(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/quotations/inquire-time-itemchartprice`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "minute_chart_success.json")), nil
		},
	)

	c := newTestClient(t)
	_, err := c.InquireTimeItemChartPrice(context.Background(), domestic.InquireTimeItemChartPriceParams{
		Symbol:          "005930",
		TimeFrom:        "150000",
		PastDataInclude: true,
		MarketCode:      "UN",
	})
	require.NoError(t, err)
	assert.Equal(t, "Y", capturedQuery.Get("FID_PW_DATA_INCU_YN"))
	assert.Equal(t, "UN", capturedQuery.Get("FID_COND_MRKT_DIV_CODE"))
}
