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

func TestClient_InquireIndexPrice(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/quotations/inquire-index-price`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "index_price_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireIndexPrice(context.Background(), domestic.InquireIndexPriceParams{
		Symbol: "0001",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, "U", capturedQuery.Get("FID_COND_MRKT_DIV_CODE"))
	assert.Equal(t, "0001", capturedQuery.Get("FID_INPUT_ISCD"))

	d, _ := decimal.NewFromString("2650.45")
	assert.True(t, d.Equal(res.Output.BstpNmixPrpr))
	assert.InDelta(t, -0.46, res.Output.BstpNmixPrdyCtrt, 0.001)
	assert.Equal(t, int64(350000000), res.Output.AcmlVol)
	assert.Equal(t, "315", res.Output.AscnIssuCnt)
	assert.Equal(t, "450", res.Output.DownIssuCnt)
}

func TestClient_InquireIndexCategoryPrice(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/quotations/inquire-index-category-price`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "index_category_price_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireIndexCategoryPrice(context.Background(), domestic.InquireIndexCategoryPriceParams{
		Symbol:    "0001",
		MarketCls: "K",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, "U", capturedQuery.Get("FID_COND_MRKT_DIV_CODE"))
	assert.Equal(t, "0001", capturedQuery.Get("FID_INPUT_ISCD"))
	assert.Equal(t, "20214", capturedQuery.Get("FID_COND_SCR_DIV_CODE"))
	assert.Equal(t, "K", capturedQuery.Get("FID_MRKT_CLS_CODE"))
	assert.Equal(t, "0", capturedQuery.Get("FID_BLNG_CLS_CODE"))

	d, _ := decimal.NewFromString("2650.45")
	assert.True(t, d.Equal(res.Output1.BstpNmixPrpr))
	assert.Equal(t, int64(350000000), res.Output1.AcmlVol)

	require.Len(t, res.Output2, 2)
	assert.Equal(t, "0001", res.Output2[0].BstpClsCode)
	assert.Equal(t, "코스피", res.Output2[0].HtsKorIsnm)
	assert.InDelta(t, 100.0, res.Output2[0].AcmlVolRlim, 0.01)
}

func TestClient_InquireIndexDailyPrice(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/quotations/inquire-index-daily-price`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "index_daily_price_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireIndexDailyPrice(context.Background(), domestic.InquireIndexDailyPriceParams{
		Symbol:        "0001",
		PeriodDivCode: "D",
		InputDate1:    "20260505",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, "U", capturedQuery.Get("FID_COND_MRKT_DIV_CODE"))
	assert.Equal(t, "0001", capturedQuery.Get("FID_INPUT_ISCD"))
	assert.Equal(t, "D", capturedQuery.Get("FID_PERIOD_DIV_CODE"))
	assert.Equal(t, "20260505", capturedQuery.Get("FID_INPUT_DATE_1"))

	d, _ := decimal.NewFromString("2650.45")
	assert.True(t, d.Equal(res.Output1.BstpNmixPrpr))
	assert.Equal(t, "315", res.Output1.AscnIssuCnt)
	assert.InDelta(t, -0.46, res.Output1.BstpNmixPrdyCtrt, 0.001)

	require.Len(t, res.Output2, 2)
	assert.Equal(t, "20260505", res.Output2[0].StckBsopDate)
	d2, _ := decimal.NewFromString("2650.45")
	assert.True(t, d2.Equal(res.Output2[0].BstpNmixPrpr))
	assert.InDelta(t, 100.00, res.Output2[0].AcmlVolRlim, 0.01)
}

func TestClient_InquireIndexTimeprice(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/quotations/inquire-index-timeprice`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "index_timeprice_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireIndexTimeprice(context.Background(), domestic.InquireIndexTimepriceParams{
		InputHour1: "60",
		Symbol:     "0001",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, "60", capturedQuery.Get("FID_INPUT_HOUR_1"))
	assert.Equal(t, "0001", capturedQuery.Get("FID_INPUT_ISCD"))
	assert.Equal(t, "U", capturedQuery.Get("FID_COND_MRKT_DIV_CODE"))

	require.Len(t, res.Output, 2)
	assert.Equal(t, "100000", res.Output[0].BsopHour)
	d, _ := decimal.NewFromString("2652.10")
	assert.True(t, d.Equal(res.Output[0].BstpNmixPrpr))
	assert.Equal(t, int64(800000), res.Output[0].CntgVol)
}

func TestClient_InquireIndexTickprice(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/quotations/inquire-index-tickprice`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "index_tickprice_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireIndexTickprice(context.Background(), domestic.InquireIndexTickpriceParams{
		Symbol: "0001",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, "0001", capturedQuery.Get("FID_INPUT_ISCD"))
	assert.Equal(t, "U", capturedQuery.Get("FID_COND_MRKT_DIV_CODE"))

	require.Len(t, res.Output, 2)
	assert.Equal(t, "100015", res.Output[0].StckCntgHour)
	d, _ := decimal.NewFromString("2651.35")
	assert.True(t, d.Equal(res.Output[0].BstpNmixPrpr))
	assert.Equal(t, int64(50000), res.Output[0].CntgVol)
}

func TestClient_InquireDailyIndexchartprice(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/quotations/inquire-daily-indexchartprice`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "daily_indexchartprice_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireDailyIndexchartprice(context.Background(), domestic.InquireDailyIndexchartpriceParams{
		Symbol:        "0001",
		InputDate1:    "20260401",
		InputDate2:    "20260505",
		PeriodDivCode: "D",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, "U", capturedQuery.Get("FID_COND_MRKT_DIV_CODE"))
	assert.Equal(t, "0001", capturedQuery.Get("FID_INPUT_ISCD"))
	assert.Equal(t, "20260401", capturedQuery.Get("FID_INPUT_DATE_1"))
	assert.Equal(t, "20260505", capturedQuery.Get("FID_INPUT_DATE_2"))
	assert.Equal(t, "D", capturedQuery.Get("FID_PERIOD_DIV_CODE"))

	assert.Equal(t, "코스피", res.Output1.HtsKorIsnm)
	assert.Equal(t, "0001", res.Output1.BstpClsCode)
	futs, _ := decimal.NewFromString("355.50")
	assert.True(t, futs.Equal(res.Output1.FutsPrdyOprc))

	require.Len(t, res.Output2, 2)
	assert.Equal(t, "20260505", res.Output2[0].StckBsopDate)
	assert.Equal(t, "N", res.Output2[0].ModYn)
	assert.Equal(t, int64(350000000), res.Output2[0].AcmlVol)
}

func TestClient_InquireTimeIndexchartprice(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/quotations/inquire-time-indexchartprice`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "time_indexchartprice_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireTimeIndexchartprice(context.Background(), domestic.InquireTimeIndexchartpriceParams{
		EtcClsCode:   "0",
		Symbol:       "0001",
		InputHour1:   "60",
		PwDataIncuYn: "Y",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, "U", capturedQuery.Get("FID_COND_MRKT_DIV_CODE"))
	assert.Equal(t, "0", capturedQuery.Get("FID_ETC_CLS_CODE"))
	assert.Equal(t, "0001", capturedQuery.Get("FID_INPUT_ISCD"))
	assert.Equal(t, "60", capturedQuery.Get("FID_INPUT_HOUR_1"))
	assert.Equal(t, "Y", capturedQuery.Get("FID_PW_DATA_INCU_YN"))

	assert.Equal(t, "코스피", res.Output1.HtsKorIsnm)
	futs, _ := decimal.NewFromString("355.50")
	assert.True(t, futs.Equal(res.Output1.FutsPrdyOprc))
	vrss, _ := decimal.NewFromString("-12.30")
	assert.True(t, vrss.Equal(res.Output1.BstpNmixPrdyVrss))

	require.Len(t, res.Output2, 2)
	assert.Equal(t, "20260505", res.Output2[0].StckBsopDate)
	assert.Equal(t, "100000", res.Output2[0].StckCntgHour)
	assert.Equal(t, int64(800000), res.Output2[0].CntgVol)
}

func TestClient_ExpTotalIndex(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/quotations/exp-total-index`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "exp_total_index_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.ExpTotalIndex(context.Background(), domestic.ExpTotalIndexParams{
		MrktClsCode: "K",
		Symbol:      "0001",
		MkopClsCode: "1",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	// ANOMALY: lowercase query param keys
	assert.Equal(t, "K", capturedQuery.Get("fid_mrkt_cls_code"))
	assert.Equal(t, "U", capturedQuery.Get("fid_cond_mrkt_div_code"))
	assert.Equal(t, "11175", capturedQuery.Get("fid_cond_scr_div_code"))
	assert.Equal(t, "0001", capturedQuery.Get("fid_input_iscd"))
	assert.Equal(t, "1", capturedQuery.Get("fid_mkop_cls_code"))

	d, _ := decimal.NewFromString("2650.45")
	assert.True(t, d.Equal(res.Output1.BstpNmixPrpr))
	// prdy_ctrt (short form) — bstp_nmix_prdy_ctrt 아님
	assert.InDelta(t, -0.46, res.Output1.PrdyCtrt, 0.001)
	assert.Equal(t, "315", res.Output1.AscnIssuCnt)
	assert.Equal(t, "0001", res.Output1.BstpClsCode)

	require.Len(t, res.Output2, 2)
	assert.Equal(t, "코스피", res.Output2[0].HtsKorIsnm)
	sdpr, _ := decimal.NewFromString("2662.75")
	assert.True(t, sdpr.Equal(res.Output2[0].NmixSdpr))
	assert.InDelta(t, -0.46, res.Output2[0].BstpNmixPrdyCtrt, 0.001)
}

func TestClient_ExpIndexTrend(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/quotations/exp-index-trend`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "exp_index_trend_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.ExpIndexTrend(context.Background(), domestic.ExpIndexTrendParams{
		MkopClsCode: "1",
		InputHour1:  "10",
		Symbol:      "0001",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, "1", capturedQuery.Get("FID_MKOP_CLS_CODE"))
	assert.Equal(t, "10", capturedQuery.Get("FID_INPUT_HOUR_1"))
	assert.Equal(t, "0001", capturedQuery.Get("FID_INPUT_ISCD"))
	assert.Equal(t, "U", capturedQuery.Get("FID_COND_MRKT_DIV_CODE"))

	require.Len(t, res.Output, 2)
	assert.Equal(t, "090000", res.Output[0].StckCntgHour)
	prpr, _ := decimal.NewFromString("2550.25")
	assert.True(t, prpr.Equal(res.Output[0].BstpNmixPrpr))
	assert.Equal(t, "2", res.Output[0].PrdyVrssSign)
	assert.InDelta(t, 0.49, res.Output[0].PrdyCtrt, 0.001)
	assert.Equal(t, int64(123456789), res.Output[0].AcmlVol)
	assert.Equal(t, int64(987654321000), res.Output[0].AcmlTrPbmn)
}
