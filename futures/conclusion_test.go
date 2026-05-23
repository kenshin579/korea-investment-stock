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

// ─── EP5: ExpPriceTrend ──────────────────────────────────────────────────────

func TestClient_ExpPriceTrend(t *testing.T) {
	client, transport := newTestClient(t)
	transport.RegisterResponder(
		http.MethodGet,
		"=~/uapi/domestic-futureoption/v1/quotations/exp-price-trend",
		httpmock.NewStringResponder(http.StatusOK, loadFixtureString(t, "exp_price_trend_success.json")),
	)

	got, err := client.ExpPriceTrend(context.Background(), futures.ExpPriceTrendParams{
		Code:       "101S06",
		MarketCode: "F",
	})
	require.NoError(t, err)
	require.NotNil(t, got)

	// output1 assertions
	assert.Equal(t, "KOSPI200선물 2506", got.Output1.HtsKorIsnm)
	antcCnpr, _ := decimal.NewFromString("363.75")
	assert.True(t, antcCnpr.Equal(got.Output1.FutsAntcCnpr))
	assert.Equal(t, "2", got.Output1.AntcCntgVrssSign)
	antcVrss, _ := decimal.NewFromString("2.50")
	assert.True(t, antcVrss.Equal(got.Output1.FutsAntcCntgVrss))
	assert.InDelta(t, 0.70, float64(got.Output1.AntcCntgPrdyCtrt), 0.001)
	sdpr, _ := decimal.NewFromString("361.25")
	assert.True(t, sdpr.Equal(got.Output1.FutsSdpr))

	// output2 assertions
	require.Len(t, got.Output2, 3)
	assert.Equal(t, "090000", got.Output2[0].StckCntgHour)
	antcCnpr0, _ := decimal.NewFromString("363.50")
	assert.True(t, antcCnpr0.Equal(got.Output2[0].FutsAntcCnpr))

	assert.Equal(t, "093000", got.Output2[1].StckCntgHour)
	antcCnpr1, _ := decimal.NewFromString("363.75")
	assert.True(t, antcCnpr1.Equal(got.Output2[1].FutsAntcCnpr))

	assert.Equal(t, "100000", got.Output2[2].StckCntgHour)
	antcCnpr2, _ := decimal.NewFromString("364.00")
	assert.True(t, antcCnpr2.Equal(got.Output2[2].FutsAntcCnpr))
}

func TestClient_ExpPriceTrend_InvalidJSON(t *testing.T) {
	client, transport := newTestClient(t)
	transport.RegisterResponder(
		http.MethodGet,
		"=~/uapi/domestic-futureoption/v1/quotations/exp-price-trend",
		httpmock.NewStringResponder(http.StatusOK, "invalid json"),
	)

	got, err := client.ExpPriceTrend(context.Background(), futures.ExpPriceTrendParams{
		Code:       "101S06",
		MarketCode: "F",
	})
	require.Error(t, err)
	require.Nil(t, got)
}
