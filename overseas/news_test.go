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

func TestClient_InquireNewsTitle(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/quotations/news-title`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "news_title_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireNewsTitle(context.Background(), overseas.InquireNewsTitleParams{
		NationCd:   "US",
		ExchangeCd: "NAS",
		DataDt:     "20260505",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, "US", capturedQuery.Get("NATION_CD"))
	assert.Equal(t, "NAS", capturedQuery.Get("EXCHANGE_CD"))
	assert.Equal(t, "20260505", capturedQuery.Get("DATA_DT"))

	require.Len(t, res.Outblock1, 2)
	assert.Equal(t, "AAPL", res.Outblock1[0].Symb)
	assert.Equal(t, "20260505120001", res.Outblock1[0].NewsKey)
	assert.Equal(t, "애플 2분기 실적 서프라이즈 발표", res.Outblock1[0].Title)
}
