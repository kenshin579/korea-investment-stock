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
