package domestic_test

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_InquireMarketTime(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	var capturedTrID string
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/quotations/market-time`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			capturedTrID = req.Header.Get("tr_id")
			return httpmock.NewStringResponse(200, loadFixtureString(t, "market_time_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireMarketTime(context.Background())
	require.NoError(t, err)
	require.NotNil(t, res)

	// 파라미터 없음 — query 가 비어 있어야 함
	assert.Empty(t, capturedQuery)
	// tr_id 가 정확히 들어갔는지
	assert.Equal(t, "HHMCM000002C0", capturedTrID)

	require.Len(t, res.Output1, 1)
	assert.Equal(t, "20260504", res.Output1[0].Date1)
	assert.Equal(t, "20260508", res.Output1[0].Date3) // 영업일 당일
	assert.Equal(t, "20260512", res.Output1[0].Date5)
	assert.Equal(t, "20260508", res.Output1[0].Today)
	assert.Equal(t, "143005", res.Output1[0].Time)
	assert.Equal(t, "090000", res.Output1[0].STime)
	assert.Equal(t, "153000", res.Output1[0].ETime)
}

func TestClient_InquireMarketTime_InvalidJSON(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// envelope 은 valid 하지만 output1 이 array 가 아닌 string → unmarshal 실패
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/quotations/market-time`,
		httpmock.NewStringResponder(200, `{"rt_cd":"0","msg_cd":"X","msg1":"x","output1":"not-array"}`),
	)

	c := newTestClient(t)
	_, err := c.InquireMarketTime(context.Background())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "MarketTime")
}
