package domestic_test

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/kenshin579/korea-investment-stock/domestic"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInquireExpClosingPrice(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/quotations/exp-closing-price`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "exp_closing_price_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireExpClosingPrice(context.Background(), domestic.InquireExpClosingPriceParams{
		RankSortClsCode: "0",
		Symbol:          "0000",
		BlngClsCode:     "0",
	})
	require.NoError(t, err)
	require.Len(t, res.Output1, 2)

	// wire param keys 검증
	assert.Equal(t, "J", capturedQuery.Get("FID_COND_MRKT_DIV_CODE"))
	assert.Equal(t, "11173", capturedQuery.Get("FID_COND_SCR_DIV_CODE"))
	assert.Equal(t, "0", capturedQuery.Get("FID_RANK_SORT_CLS_CODE"))
	assert.Equal(t, "0000", capturedQuery.Get("FID_INPUT_ISCD"))
	assert.Equal(t, "0", capturedQuery.Get("FID_BLNG_CLS_CODE"))

	item := res.Output1[0]
	assert.Equal(t, "005930", item.StckShrnIscd)
	assert.Equal(t, "삼성전자", item.HtsKorIsnm)
	assert.Equal(t, "82500", item.StckPrpr.String())
	assert.Equal(t, "500", item.PrdyVrss.String())
	assert.Equal(t, "2", item.PrdyVrssSign)
	assert.InDelta(t, 0.61, item.PrdyCtrt, 0.001)
	assert.Equal(t, "1000", item.SdprVrssPrpr.String())
	assert.InDelta(t, 1.23, item.SdprVrssPrprRate, 0.001)
	assert.Equal(t, int64(125000), item.CntgVol)
}
