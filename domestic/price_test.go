package domestic_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_InquirePrice(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/quotations/inquire-price`,
		httpmock.NewStringResponder(200, loadFixtureString(t, "price_success.json")),
	)

	c := newTestClient(t)
	price, err := c.InquirePrice(context.Background(), "005930")
	require.NoError(t, err)
	require.NotNil(t, price)

	assert.Equal(t, decimal.NewFromInt(75800), price.StckPrpr)
	assert.Equal(t, decimal.NewFromInt(-200), price.PrdyVrss)
	assert.Equal(t, "5", price.PrdyVrssSign)
	assert.InDelta(t, -0.26, price.PrdyCtrt, 0.001)
	assert.Equal(t, int64(12345678), price.AcmlVol)
	assert.Equal(t, int64(938223456000), price.AcmlTrPbmn)
	assert.Equal(t, decimal.NewFromInt(76000), price.StckOprc)
	assert.Equal(t, decimal.NewFromInt(76200), price.StckHgpr)
	assert.Equal(t, decimal.NewFromInt(75500), price.StckLwpr)
	assert.Equal(t, "005930", price.StckShrnIscd)
	assert.Equal(t, "Y", price.SstsYn)
	assert.Equal(t, "N", price.MangIssuClsCode)
	assert.InDelta(t, 11.42, price.Per, 0.001)
	assert.InDelta(t, 1.32, price.Pbr, 0.001)
}

func TestClient_InquirePrice_APIError(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/quotations/inquire-price`,
		httpmock.NewStringResponder(200, `{"rt_cd":"1","msg_cd":"MCA00001","msg1":"잘못된 요청","output":null}`),
	)

	c := newTestClient(t)
	_, err := c.InquirePrice(context.Background(), "005930")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "MCA00001")
}
