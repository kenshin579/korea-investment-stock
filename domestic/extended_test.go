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

func TestClient_InquireNearNewHighlow(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/ranking/near-new-highlow`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "near_new_highlow_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireNearNewHighlow(context.Background(), domestic.InquireNearNewHighlowParams{
		InputISCD:  "0000",
		PrcClsCode: "0",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	// query param 검증
	assert.Equal(t, "J", capturedQuery.Get("fid_cond_mrkt_div_code"))
	assert.Equal(t, "20187", capturedQuery.Get("fid_cond_scr_div_code"))
	assert.Equal(t, "0000", capturedQuery.Get("fid_input_iscd"))
	assert.Equal(t, "0", capturedQuery.Get("fid_prc_cls_code"))

	require.Len(t, res.Output, 2)

	// output[0] 필드 검증
	assert.Equal(t, "삼성전자", res.Output[0].HtsKorIsnm)
	assert.Equal(t, "005930", res.Output[0].MkscShrnIscd)
	d, _ := decimal.NewFromString("75800")
	assert.True(t, d.Equal(res.Output[0].StckPrpr))
	assert.Equal(t, int64(12345678), res.Output[0].AcmlVol)
	d2, _ := decimal.NewFromString("76500")
	assert.True(t, d2.Equal(res.Output[0].NewHgpr))
	assert.InDelta(t, 1.24, res.Output[0].HprcNearRate, 0.01)
	d3, _ := decimal.NewFromString("74000")
	assert.True(t, d3.Equal(res.Output[0].NewLwpr))
	assert.InDelta(t, 2.43, res.Output[0].LwprNearRate, 0.01)
}
