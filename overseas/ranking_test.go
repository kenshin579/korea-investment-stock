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
