package domestic_test

import (
	"context"
	"net/http"
	"regexp"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kenshin579/korea-investment-stock/internal/krxmaster"
)

func TestClient_FetchKospiSymbols(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	zipBytes := loadFixture(t, "kospi_code_sample.mst.zip")
	httpmock.RegisterResponder(http.MethodGet, krxmaster.KospiURL,
		httpmock.NewBytesResponder(200, zipBytes))

	c := newTestClient(t)
	syms, err := c.FetchKospiSymbols(context.Background())
	require.NoError(t, err)
	require.Len(t, syms, 3)

	hangulRe := regexp.MustCompile(`[\x{AC00}-\x{D7A3}]`)
	for i, s := range syms {
		assert.Regexp(t, `^[0-9A-Z]{6}$`, s.ShortCode, "row %d ShortCode", i)
		assert.True(t, hangulRe.MatchString(s.KoreanName), "row %d KoreanName 한글", i)
	}
}

func TestClient_FetchKosdaqSymbols(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	zipBytes := loadFixture(t, "kosdaq_code_sample.mst.zip")
	httpmock.RegisterResponder(http.MethodGet, krxmaster.KosdaqURL,
		httpmock.NewBytesResponder(200, zipBytes))

	c := newTestClient(t)
	syms, err := c.FetchKosdaqSymbols(context.Background())
	require.NoError(t, err)
	require.Len(t, syms, 3)
	for _, s := range syms {
		assert.NotEmpty(t, s.ShortCode)
		assert.NotEmpty(t, s.KoreanName)
	}
}

func TestClient_FetchKospiSymbols_DownloadError(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(http.MethodGet, krxmaster.KospiURL,
		httpmock.NewStringResponder(500, "internal error"))

	c := newTestClient(t)
	_, err := c.FetchKospiSymbols(context.Background())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "HTTP 500")
}
