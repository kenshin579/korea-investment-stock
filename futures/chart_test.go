package futures_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jarcoal/httpmock"

	"github.com/kenshin579/korea-investment-stock/futures"
)

// ─── EP3: InquireTimeFuopchartprice ──────────────────────────────────────────

func TestClient_InquireTimeFuopchartprice(t *testing.T) {
	client, transport := newTestClient(t)
	transport.RegisterResponder(
		http.MethodGet,
		"=~/uapi/domestic-futureoption/v1/quotations/inquire-time-fuopchartprice",
		httpmock.NewStringResponder(http.StatusOK, loadFixtureString(t, "inquire_time_fuopchartprice_success.json")),
	)

	got, err := client.InquireTimeFuopchartprice(context.Background(), futures.InquireTimeFuopchartpriceParams{
		MarketCode:     "F",
		Code:           "101S06",
		HourClsCode:    "60",
		PwDataIncuYn:   "N",
		FakeTickIncuYn: "N",
		InputDate:      "20250509",
		InputHour:      "141500",
	})
	require.NoError(t, err)
	require.NotNil(t, got)

	// output1 assertions
	assert.Equal(t, "KOSPI200선물 2506", got.Output1.HtsKorIsnm)
	prpr, _ := decimal.NewFromString("362.50")
	assert.True(t, prpr.Equal(got.Output1.FutsPrpr))
	assert.Equal(t, int64(125430), got.Output1.AcmlVol)
	assert.Equal(t, int64(310250), got.Output1.HtsOtstStplQty)
	assert.Equal(t, "101S06", got.Output1.FutsShrnIscd)
	assert.InDelta(t, 0.35, float64(got.Output1.FutsPrdyCtrt), 0.001)
	assert.InDelta(t, 55.20, float64(got.Output1.TdayRltv), 0.001)

	// output2 array assertions
	require.GreaterOrEqual(t, len(got.Output2), 1)
	candle := got.Output2[0]
	assert.Equal(t, "20250509", candle.StckBsopDate)
	assert.Equal(t, "141500", candle.StckCntgHour)
	candlePrpr, _ := decimal.NewFromString("362.50")
	assert.True(t, candlePrpr.Equal(candle.FutsPrpr))
	assert.Equal(t, int64(320), candle.CntgVol)
	assert.Equal(t, int64(115840), candle.AcmlTrPbmn)
}

func TestClient_InquireTimeFuopchartprice_InvalidJSON(t *testing.T) {
	client, transport := newTestClient(t)
	transport.RegisterResponder(
		http.MethodGet,
		"=~/uapi/domestic-futureoption/v1/quotations/inquire-time-fuopchartprice",
		httpmock.NewStringResponder(http.StatusOK, `{"rt_cd":"0","msg_cd":"X","msg1":"x","output1":"not-an-object"}`),
	)
	_, err := client.InquireTimeFuopchartprice(context.Background(), futures.InquireTimeFuopchartpriceParams{
		MarketCode: "F",
		Code:       "101S06",
	})
	require.Error(t, err)
}

// ─── EP6: InquireDailyFuopchartprice ─────────────────────────────────────────

func TestClient_InquireDailyFuopchartprice(t *testing.T) {
	client, transport := newTestClient(t)
	transport.RegisterResponder(
		http.MethodGet,
		"=~/uapi/domestic-futureoption/v1/quotations/inquire-daily-fuopchartprice",
		httpmock.NewStringResponder(http.StatusOK, loadFixtureString(t, "inquire_daily_fuopchartprice_success.json")),
	)

	got, err := client.InquireDailyFuopchartprice(context.Background(), futures.InquireDailyFuopchartpriceParams{
		MarketCode: "F",
		Code:       "101S06",
		FromDate:   "20250401",
		ToDate:     "20250509",
		Period:     "D",
	})
	require.NoError(t, err)
	require.NotNil(t, got)

	// output1 assertions
	assert.Equal(t, "KOSPI200선물 2506", got.Output1.HtsKorIsnm)
	prpr, _ := decimal.NewFromString("362.50")
	assert.True(t, prpr.Equal(got.Output1.FutsPrpr))
	assert.Equal(t, int64(125430), got.Output1.AcmlVol)
	assert.Equal(t, int64(310250), got.Output1.HtsOtstStplQty)
	assert.InDelta(t, 0.35, float64(got.Output1.FutsPrdyCtrt), 0.001)

	// output2 array assertions
	require.GreaterOrEqual(t, len(got.Output2), 1)
	candle := got.Output2[0]
	assert.Equal(t, "20250509", candle.StckBsopDate)
	candlePrpr, _ := decimal.NewFromString("362.50")
	assert.True(t, candlePrpr.Equal(candle.FutsPrpr))
	assert.Equal(t, int64(125430), candle.AcmlVol)
	assert.Equal(t, int64(4530210), candle.AcmlTrPbmn)
	assert.Equal(t, "N", candle.ModYn)
}

func TestClient_InquireDailyFuopchartprice_InvalidJSON(t *testing.T) {
	client, transport := newTestClient(t)
	transport.RegisterResponder(
		http.MethodGet,
		"=~/uapi/domestic-futureoption/v1/quotations/inquire-daily-fuopchartprice",
		httpmock.NewStringResponder(http.StatusOK, `{"rt_cd":"0","msg_cd":"X","msg1":"x","output1":"not-an-object"}`),
	)
	_, err := client.InquireDailyFuopchartprice(context.Background(), futures.InquireDailyFuopchartpriceParams{
		MarketCode: "F",
		Code:       "101S06",
		FromDate:   "20250401",
		ToDate:     "20250509",
	})
	require.Error(t, err)
}
