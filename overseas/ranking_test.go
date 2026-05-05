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

func TestClient_InquireUpdownRate(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/ranking/updown-rate`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "updown_rate_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireUpdownRate(context.Background(), overseas.InquireUpdownRateParams{
		Excd: "NAS",
		Gubn: "1", // 상승율
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, "NAS", capturedQuery.Get("EXCD"))
	assert.Equal(t, "1", capturedQuery.Get("GUBN"))
	assert.Equal(t, "0", capturedQuery.Get("NDAY"))
	assert.Equal(t, "0", capturedQuery.Get("VOL_RANG"))

	require.Len(t, res.Output2, 2)
	assert.Equal(t, "NVDA", res.Output2[0].Symb)
	assert.Equal(t, "엔비디아", res.Output2[0].Name)
	d, _ := decimal.NewFromString("920.45")
	assert.True(t, d.Equal(res.Output2[0].Last))
	assert.InDelta(t, 5.16, res.Output2[0].Rate, 0.001)
}

func TestClient_InquireMarketCap(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/ranking/market-cap`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "market_cap_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireMarketCap(context.Background(), overseas.InquireMarketCapParams{
		ExcdCode: "NAS",
		VolRang:  "0",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, "NAS", capturedQuery.Get("EXCD"))
	assert.Equal(t, "0", capturedQuery.Get("VOL_RANG"))

	// output1 검증
	assert.Equal(t, "2", res.Output1.Zdiv)
	assert.Equal(t, int64(2), res.Output1.Crec)
	assert.Equal(t, int64(500), res.Output1.Trec)
	assert.Equal(t, int64(30), res.Output1.Nrec)

	// output2[0] 검증
	require.Len(t, res.Output2, 2)
	assert.Equal(t, "AAPL", res.Output2[0].Symb)
	assert.Equal(t, "애플", res.Output2[0].Name)
	d, _ := decimal.NewFromString("189.30")
	assert.True(t, d.Equal(res.Output2[0].Last))
	assert.InDelta(t, 1.34, res.Output2[0].Rate, 0.001)
	assert.Equal(t, int64(55000000), res.Output2[0].Tvol)
	assert.Equal(t, int64(15634232000), res.Output2[0].Shar)
	tomv, _ := decimal.NewFromString("2958652560000")
	assert.True(t, tomv.Equal(res.Output2[0].Tomv))
	assert.InDelta(t, 6.85, res.Output2[0].Grav, 0.001)
	assert.Equal(t, int64(1), res.Output2[0].Rank)
	assert.Equal(t, "APPLE INC", res.Output2[0].Ename)
}
