package domestic_test

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kenshin579/korea-investment-stock/domestic"
)

func TestClient_InquireVolumeRank(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/quotations/volume-rank`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "volume_rank_success.json")), nil
		},
	)

	c := newTestClient(t)
	rank, err := c.InquireVolumeRank(context.Background(), domestic.InquireVolumeRankParams{
		InputISCD: "0000",
	})
	require.NoError(t, err)
	require.NotNil(t, rank)

	// 필수 query 기본값 검증
	assert.Equal(t, "J", capturedQuery.Get("FID_COND_MRKT_DIV_CODE"))
	assert.Equal(t, "20171", capturedQuery.Get("FID_COND_SCR_DIV_CODE"))
	assert.Equal(t, "0000", capturedQuery.Get("FID_INPUT_ISCD"))
	assert.Equal(t, "0", capturedQuery.Get("FID_DIV_CLS_CODE"))
	assert.Equal(t, "0", capturedQuery.Get("FID_BLNG_CLS_CODE"))
	assert.Equal(t, "111111111", capturedQuery.Get("FID_TRGT_CLS_CODE"))
	assert.Equal(t, "0000000000", capturedQuery.Get("FID_TRGT_EXLS_CLS_CODE"))

	// 응답 검증
	require.Len(t, rank.Output, 2)
	assert.Equal(t, "삼성전자", rank.Output[0].HtsKorIsnm)
	assert.Equal(t, "005930", rank.Output[0].MkscShrnIscd)
	assert.Equal(t, "1", rank.Output[0].DataRank)
	assert.Equal(t, decimal.NewFromInt(75800), rank.Output[0].StckPrpr)
	assert.Equal(t, int64(12345678), rank.Output[0].AcmlVol)
	assert.InDelta(t, 0.21, rank.Output[0].VolTnrt, 0.001)
	assert.Equal(t, int64(938223456000), rank.Output[0].AcmlTrPbmn)
}
