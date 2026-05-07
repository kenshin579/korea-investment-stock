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

func TestClient_InquireEtfPrice(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/etfetn/v1/quotations/inquire-price`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "etf_price_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireEtfPrice(context.Background(), domestic.InquireEtfPriceParams{
		Symbol: "069500",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	// query param 검증
	assert.Equal(t, "J", capturedQuery.Get("fid_cond_mrkt_div_code"))
	assert.Equal(t, "069500", capturedQuery.Get("fid_input_iscd"))

	// output 필드 검증
	d, _ := decimal.NewFromString("28350")
	assert.True(t, d.Equal(res.Output.StckPrpr))
	assert.Equal(t, "5", res.Output.PrdyVrssSign)
	dVrss, _ := decimal.NewFromString("-150")
	assert.True(t, dVrss.Equal(res.Output.PrdyVrss))
	assert.InDelta(t, -0.53, res.Output.PrdyCtrt, 0.001)
	assert.Equal(t, int64(12345678), res.Output.AcmlVol)
	assert.Equal(t, int64(10234567), res.Output.PrdyVol)

	dNav, _ := decimal.NewFromString("28342")
	assert.True(t, dNav.Equal(res.Output.Nav))
	dNavVrss, _ := decimal.NewFromString("-138")
	assert.True(t, dNavVrss.Equal(res.Output.NavPrdyVrss))
	assert.InDelta(t, -0.48, res.Output.NavPrdyCtrt, 0.001)
	assert.InDelta(t, 0.02, res.Output.TrcErrt, 0.001)

	assert.Equal(t, int64(50000000), res.Output.EtfCrclStcn)
	assert.Equal(t, int64(1417100000000), res.Output.EtfNtasTtam)
	assert.Equal(t, int64(50000), res.Output.EtfCuUnitScrtCnt)
	assert.Equal(t, int64(93), res.Output.EtfCnfgIssuCnt)

	assert.InDelta(t, -0.03, res.Output.Dprt, 0.001)
	assert.Equal(t, "삼성자산운용", res.Output.MbcrName)
	assert.InDelta(t, 0.58, res.Output.LpHldnRate, 0.001)
	assert.Equal(t, int64(290000), res.Output.LpHldnVol)
}

func TestClient_InquireComponentStockPrice(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/etfetn/v1/quotations/inquire-component-stock-price`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "component_stock_price_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireComponentStockPrice(context.Background(), domestic.InquireComponentStockPriceParams{
		Symbol: "069500",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	// query param 검증
	assert.Equal(t, "J", capturedQuery.Get("FID_COND_MRKT_DIV_CODE"))
	assert.Equal(t, "069500", capturedQuery.Get("FID_INPUT_ISCD"))
	assert.Equal(t, "11216", capturedQuery.Get("FID_COND_SCR_DIV_CODE"))

	// output1 필드 검증
	d1, _ := decimal.NewFromString("28350")
	assert.True(t, d1.Equal(res.Output1.StckPrpr))
	dNav, _ := decimal.NewFromString("28342")
	assert.True(t, dNav.Equal(res.Output1.Nav))
	assert.Equal(t, int64(1417100000000), res.Output1.EtfNtasTtam)
	assert.Equal(t, int64(93), res.Output1.EtfCnfgIssuCnt)
	assert.Equal(t, int64(50000), res.Output1.EtfCuUnitScrtCnt)

	// output2 배열 검증
	require.Len(t, res.Output2, 2)
	assert.Equal(t, "005930", res.Output2[0].StckShrnIscd)
	assert.Equal(t, "삼성전자", res.Output2[0].HtsKorIsnm)
	d2, _ := decimal.NewFromString("57800")
	assert.True(t, d2.Equal(res.Output2[0].StckPrpr))
	assert.Equal(t, int64(18234567), res.Output2[0].AcmlVol)
	assert.InDelta(t, 19.99, res.Output2[0].EtfCnfgIssuRlim, 0.001)
	assert.Equal(t, "000660", res.Output2[1].StckShrnIscd)
}

func TestClient_InquireNavComparisonTimeTrend(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/etfetn/v1/quotations/nav-comparison-time-trend`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "nav_comparison_time_trend_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireNavComparisonTimeTrend(context.Background(), domestic.InquireNavComparisonTimeTrendParams{
		HourClsCode: "60",
		Symbol:      "069500",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	// query param 검증
	assert.Equal(t, "E", capturedQuery.Get("fid_cond_mrkt_div_code"))
	assert.Equal(t, "60", capturedQuery.Get("fid_hour_cls_code"))
	assert.Equal(t, "069500", capturedQuery.Get("fid_input_iscd"))

	// output 배열 검증
	require.Len(t, res.Output, 2)
	assert.Equal(t, "130000", res.Output[0].BsopHour)
	dNav, _ := decimal.NewFromString("28342")
	assert.True(t, dNav.Equal(res.Output[0].Nav))
	assert.InDelta(t, -0.48, res.Output[0].NavPrdyCtrt, 0.001)
	assert.InDelta(t, -0.03, res.Output[0].Dprt, 0.001)
	dPrpr, _ := decimal.NewFromString("28350")
	assert.True(t, dPrpr.Equal(res.Output[0].StckPrpr))
	assert.Equal(t, int64(9876543), res.Output[0].AcmlVol)
	assert.Equal(t, int64(1234), res.Output[0].CntgVol)
	assert.Equal(t, "131500", res.Output[1].BsopHour)
}

func TestClient_InquireNavComparisonDailyTrend(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/etfetn/v1/quotations/nav-comparison-daily-trend`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "nav_comparison_daily_trend_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireNavComparisonDailyTrend(context.Background(), domestic.InquireNavComparisonDailyTrendParams{
		Symbol:     "069500",
		InputDate1: "20250401",
		InputDate2: "20250506",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	// query param 검증
	assert.Equal(t, "J", capturedQuery.Get("fid_cond_mrkt_div_code"))
	assert.Equal(t, "069500", capturedQuery.Get("fid_input_iscd"))
	assert.Equal(t, "20250401", capturedQuery.Get("fid_input_date_1"))
	assert.Equal(t, "20250506", capturedQuery.Get("fid_input_date_2"))

	// output 배열 검증
	require.Len(t, res.Output, 2)
	assert.Equal(t, "20250506", res.Output[0].StckBsopDate)
	dClpr, _ := decimal.NewFromString("28350")
	assert.True(t, dClpr.Equal(res.Output[0].StckClpr))
	assert.Equal(t, int64(12345678), res.Output[0].AcmlVol)
	assert.InDelta(t, -0.03, res.Output[0].Dprt, 0.001)
	dNav, _ := decimal.NewFromString("28342")
	assert.True(t, dNav.Equal(res.Output[0].Nav))
	assert.InDelta(t, -0.48, res.Output[0].NavPrdyCtrt, 0.001)
	assert.Equal(t, "20250502", res.Output[1].StckBsopDate)
}

func TestClient_InquireNavComparisonTrend(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/etfetn/v1/quotations/nav-comparison-trend`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "nav_comparison_trend_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireNavComparisonTrend(context.Background(), domestic.InquireNavComparisonTrendParams{
		Symbol: "069500",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	// query param 검증
	assert.Equal(t, "J", capturedQuery.Get("FID_COND_MRKT_DIV_CODE"))
	assert.Equal(t, "069500", capturedQuery.Get("FID_INPUT_ISCD"))

	// output1 필드 검증
	d1, _ := decimal.NewFromString("28350")
	assert.True(t, d1.Equal(res.Output1.StckPrpr))
	dVrss, _ := decimal.NewFromString("-150")
	assert.True(t, dVrss.Equal(res.Output1.PrdyVrss))
	assert.InDelta(t, -0.53, res.Output1.PrdyCtrt, 0.001)
	assert.Equal(t, int64(12345678), res.Output1.AcmlVol)
	assert.Equal(t, int64(350123456789), res.Output1.AcmlTrPbmn)
	dMxpr, _ := decimal.NewFromString("36850")
	assert.True(t, dMxpr.Equal(res.Output1.StckMxpr))
	dLlam, _ := decimal.NewFromString("19850")
	assert.True(t, dLlam.Equal(res.Output1.StckLlam))

	// output2 필드 검증
	dNav, _ := decimal.NewFromString("28342")
	assert.True(t, dNav.Equal(res.Output2.Nav))
	assert.Equal(t, "5", res.Output2.NavPrdyVrssSign)
	dNavVrss, _ := decimal.NewFromString("-138")
	assert.True(t, dNavVrss.Equal(res.Output2.NavPrdyVrss))
	assert.InDelta(t, -0.48, res.Output2.NavPrdyCtrt, 0.001)
	dHprcNav, _ := decimal.NewFromString("28560")
	assert.True(t, dHprcNav.Equal(res.Output2.HprcNav))
}
