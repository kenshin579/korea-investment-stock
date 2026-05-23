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

func TestClient_InquirePriceDetail(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/quotations/price-detail`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "price_detail_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquirePriceDetail(context.Background(), overseas.InquirePriceDetailParams{
		Excd: "NAS",
		Symb: "AAPL",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, "", capturedQuery.Get("AUTH"))
	assert.Equal(t, "NAS", capturedQuery.Get("EXCD"))
	assert.Equal(t, "AAPL", capturedQuery.Get("SYMB"))

	d, _ := decimal.NewFromString("181.45")
	assert.True(t, d.Equal(res.Output.Last))
	assert.Equal(t, "USD", res.Output.Curr)
	assert.Equal(t, int64(85000000), res.Output.Tvol)
	assert.InDelta(t, 29.45, float64(res.Output.Perx), 0.001)
}
