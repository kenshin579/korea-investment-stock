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

func TestClient_InquireVolumeRank_Variant(t *testing.T) {
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
		MarketCode: "NX",
		InputISCD:  "0000",
		DivCode:    "1",      // 보통주
		BelongCode: "3",      // 거래금액순
	})
	require.NoError(t, err)

	// Override 검증
	assert.Equal(t, "NX", capturedQuery.Get("FID_COND_MRKT_DIV_CODE"))
	assert.Equal(t, "1", capturedQuery.Get("FID_DIV_CLS_CODE"))
	assert.Equal(t, "3", capturedQuery.Get("FID_BLNG_CLS_CODE"))

	// negative decimal 검증
	require.GreaterOrEqual(t, len(rank.Output), 1)
	assert.True(t, rank.Output[0].PrdyVrss.IsNegative(), "PrdyVrss=-200 must be negative")
}

func TestClient_InquireFluctuation(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/ranking/fluctuation`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "fluctuation_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireFluctuation(context.Background(), domestic.InquireFluctuationParams{
		InputISCD: "0000",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, "J", capturedQuery.Get("fid_cond_mrkt_div_code"))
	assert.Equal(t, "20170", capturedQuery.Get("fid_cond_scr_div_code"))
	assert.Equal(t, "0000", capturedQuery.Get("fid_input_iscd"))
	assert.Equal(t, "0", capturedQuery.Get("fid_rank_sort_cls_code"))

	require.Len(t, res.Output, 1)
	assert.Equal(t, "005930", res.Output[0].StckShrnIscd)
	assert.Equal(t, "삼성전자", res.Output[0].HtsKorIsnm)
	assert.Equal(t, decimal.NewFromInt(75800), res.Output[0].StckPrpr)
	assert.Equal(t, decimal.NewFromInt(76200), res.Output[0].StckHgpr)
	assert.Equal(t, decimal.NewFromInt(75500), res.Output[0].StckLwpr)
	assert.Equal(t, "131542", res.Output[0].HgprHour)
}
