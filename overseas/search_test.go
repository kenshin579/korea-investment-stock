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

func TestClient_SearchInfo(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/quotations/search-info`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "search_info_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.SearchInfo(context.Background(), overseas.SearchInfoParams{
		PrdtTypeCD: "512",
		Pdno:       "AAPL",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, "512", capturedQuery.Get("PRDT_TYPE_CD"))
	assert.Equal(t, "AAPL", capturedQuery.Get("PDNO"))

	assert.Equal(t, "APPLE INC", res.Output.PrdtEngName)
	assert.Equal(t, "미국", res.Output.NatnName)
	assert.Equal(t, "NASD", res.Output.OvrsExcgCd)
	assert.Equal(t, "USD", res.Output.TrCrcyCd)
	assert.Equal(t, "STK", res.Output.PrdtClsfCd)
	assert.Equal(t, "Y", res.Output.LstgYn)
}
