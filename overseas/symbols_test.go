package overseas_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kenshin579/korea-investment-stock/internal/overseasmaster"
)

func TestClient_FetchOverseasSymbols_NAS(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	zipBytes := loadFixture(t, "nas_code_sample.cod.zip")
	httpmock.RegisterResponder(http.MethodGet, overseasmaster.MarketURLs["nas"],
		httpmock.NewBytesResponder(200, zipBytes))

	c := newTestClient(t)
	syms, err := c.FetchOverseasSymbols(context.Background(), "nas")
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(syms), 1)
	for _, s := range syms {
		assert.NotEmpty(t, s.Symbol)
	}
}

func TestClient_FetchOverseasSymbols_UnknownMarket(t *testing.T) {
	c := newTestClient(t)
	_, err := c.FetchOverseasSymbols(context.Background(), "invalid")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unknown market")
}

func TestClient_FetchOverseasSymbols_DownloadError(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(http.MethodGet, overseasmaster.MarketURLs["nas"],
		httpmock.NewStringResponder(500, "internal error"))

	c := newTestClient(t)
	_, err := c.FetchOverseasSymbols(context.Background(), "nas")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "HTTP 500")
}
