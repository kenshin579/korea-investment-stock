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

func TestClient_InquireAskingPriceExpCcn(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/quotations/inquire-asking-price-exp-ccn`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "asking_price_exp_ccn_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireAskingPriceExpCcn(context.Background(), domestic.InquireAskingPriceExpCcnParams{
		Symbol: "005930",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, "J", capturedQuery.Get("FID_COND_MRKT_DIV_CODE"))
	assert.Equal(t, "005930", capturedQuery.Get("FID_INPUT_ISCD"))

	// output1 검증
	assert.Equal(t, "131542", res.Output1.AsprAcptHour)
	d, _ := decimal.NewFromString("75900")
	assert.True(t, d.Equal(res.Output1.Askp1))
	d, _ = decimal.NewFromString("75800")
	assert.True(t, d.Equal(res.Output1.Bidp1))
	assert.Equal(t, int64(1500), res.Output1.AskpRsqn1)
	assert.Equal(t, int64(2500), res.Output1.BidpRsqn1)
	assert.Equal(t, int64(11100), res.Output1.TotalAskpRsqn)
	assert.Equal(t, int64(13400), res.Output1.TotalBidpRsqn)

	// output2 검증
	assert.Equal(t, "030", res.Output2.AntcMkopClsCode)
	d, _ = decimal.NewFromString("75800")
	assert.True(t, d.Equal(res.Output2.StckPrpr))
	assert.Equal(t, int64(12345678), res.Output2.AntcVol)
	assert.Equal(t, "005930", res.Output2.StckShrnIscd)
}
