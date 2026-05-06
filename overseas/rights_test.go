package overseas_test

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kenshin579/korea-investment-stock/overseas"
)

func TestClient_InquireRightsByIce(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/quotations/rights-by-ice`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "rights_by_ice_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireRightsByIce(context.Background(), overseas.InquireRightsByIceParams{
		NCod:  "US",
		Symb:  "AAPL",
		StYmd: "20260401",
		EdYmd: "20260430",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, "US", capturedQuery.Get("NCOD"))
	assert.Equal(t, "AAPL", capturedQuery.Get("SYMB"))
	assert.Equal(t, "20260401", capturedQuery.Get("ST_YMD"))
	assert.Equal(t, "20260430", capturedQuery.Get("ED_YMD"))

	// output1 only — output2 없음 anomaly 검증
	require.Len(t, res.Output1, 2)
	assert.Equal(t, "20260401", res.Output1[0].AnnoDt)
	assert.Equal(t, "주식배당", res.Output1[0].CaTitle)
	assert.Equal(t, "20260430", res.Output1[0].PayDt)
}
