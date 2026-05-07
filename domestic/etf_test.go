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

func TestClient_InquireEtfPrice(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/etfetn/v1/quotations/inquire-price`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "etf_price_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireEtfPrice(context.Background(), domestic.InquireEtfPriceParams{
		Symbol: "069500",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	// query param 검증
	assert.Equal(t, "J", capturedQuery.Get("fid_cond_mrkt_div_code"))
	assert.Equal(t, "069500", capturedQuery.Get("fid_input_iscd"))

	// output 필드 검증
	d, _ := decimal.NewFromString("28350")
	assert.True(t, d.Equal(res.Output.StckPrpr))
	assert.Equal(t, "5", res.Output.PrdyVrssSign)
	dVrss, _ := decimal.NewFromString("-150")
	assert.True(t, dVrss.Equal(res.Output.PrdyVrss))
	assert.InDelta(t, -0.53, res.Output.PrdyCtrt, 0.001)
	assert.Equal(t, int64(12345678), res.Output.AcmlVol)
	assert.Equal(t, int64(10234567), res.Output.PrdyVol)

	dNav, _ := decimal.NewFromString("28342")
	assert.True(t, dNav.Equal(res.Output.Nav))
	dNavVrss, _ := decimal.NewFromString("-138")
	assert.True(t, dNavVrss.Equal(res.Output.NavPrdyVrss))
	assert.InDelta(t, -0.48, res.Output.NavPrdyCtrt, 0.001)
	assert.InDelta(t, 0.02, res.Output.TrcErrt, 0.001)

	assert.Equal(t, int64(50000000), res.Output.EtfCrclStcn)
	assert.Equal(t, int64(1417100000000), res.Output.EtfNtasTtam)
	assert.Equal(t, int64(50000), res.Output.EtfCuUnitScrtCnt)
	assert.Equal(t, int64(93), res.Output.EtfCnfgIssuCnt)

	assert.InDelta(t, -0.03, res.Output.Dprt, 0.001)
	assert.Equal(t, "삼성자산운용", res.Output.MbcrName)
	assert.InDelta(t, 0.58, res.Output.LpHldnRate, 0.001)
	assert.Equal(t, int64(290000), res.Output.LpHldnVol)
}
