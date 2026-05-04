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

func TestClient_InquirePubOffer(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/ksdinfo/pub-offer`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "pub_offer_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquirePubOffer(context.Background(), domestic.InquirePubOfferParams{
		FromDate: "20260501",
		ToDate:   "20260531",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	// query 키가 대문자 + 한글식
	assert.Equal(t, "", capturedQuery.Get("SHT_CD"))
	assert.Equal(t, "", capturedQuery.Get("CTS"))
	assert.Equal(t, "20260501", capturedQuery.Get("F_DT"))
	assert.Equal(t, "20260531", capturedQuery.Get("T_DT"))

	require.Len(t, res.Output1, 2)
	assert.Equal(t, "999998", res.Output1[0].ShtCd)
	assert.Equal(t, "샘플바이오", res.Output1[0].IsinName)
	assert.Equal(t, decimal.NewFromInt(30000), res.Output1[0].FixSubscrPri)
	assert.Equal(t, decimal.NewFromInt(100), res.Output1[0].FaceValue)
	assert.Equal(t, "20260505 ~ 20260506", res.Output1[0].SubscrDt)
	assert.Equal(t, "한국투자증권", res.Output1[0].LeadMgr)
}
