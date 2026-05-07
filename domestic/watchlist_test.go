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

func TestClient_InquireIntstockMultprice(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/quotations/intstock-multprice`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "intstock_multprice_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireIntstockMultprice(context.Background(), domestic.InquireIntstockMultpriceParams{
		MarketCodes: []string{"J", "J"},
		Symbols:     []string{"005930", "000660"},
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	// query param 검증 — 30쌍 번호 키
	assert.Equal(t, "J", capturedQuery.Get("FID_COND_MRKT_DIV_CODE_1"))
	assert.Equal(t, "005930", capturedQuery.Get("FID_INPUT_ISCD_1"))
	assert.Equal(t, "J", capturedQuery.Get("FID_COND_MRKT_DIV_CODE_2"))
	assert.Equal(t, "000660", capturedQuery.Get("FID_INPUT_ISCD_2"))

	// output 필드 검증
	assert.Equal(t, "코스피", res.Output.KospiKosdaqClsName)
	assert.Equal(t, "장중", res.Output.MrktTrtmClsName)
	assert.Equal(t, "005930", res.Output.InterShrnIscd)
	assert.Equal(t, "삼성전자", res.Output.InterKorIsnm)

	prpr, _ := decimal.NewFromString("57800")
	assert.True(t, prpr.Equal(res.Output.Inter2Prpr))

	vrss, _ := decimal.NewFromString("-300")
	assert.True(t, vrss.Equal(res.Output.Inter2PrdyVrss))

	assert.Equal(t, "5", res.Output.PrdyVrssSign)
	assert.InDelta(t, -0.52, res.Output.PrdyCtrt, 0.001)
	assert.Equal(t, int64(18234567), res.Output.AcmlVol)

	oprc, _ := decimal.NewFromString("57600")
	assert.True(t, oprc.Equal(res.Output.Inter2Oprc))

	hgpr, _ := decimal.NewFromString("58200")
	assert.True(t, hgpr.Equal(res.Output.Inter2Hgpr))

	assert.Equal(t, int64(345678), res.Output.SelnRsqn)
	assert.Equal(t, int64(234567), res.Output.ShnuRsqn)
	assert.Equal(t, int64(3456789), res.Output.TotalAskpRsqn)
	assert.Equal(t, int64(2345678), res.Output.TotalBidpRsqn)
	assert.Equal(t, int64(1054321098765), res.Output.AcmlTrPbmn)

	sdpr, _ := decimal.NewFromString("58100")
	assert.True(t, sdpr.Equal(res.Output.Inter2Sdpr))

	antcVrss, _ := decimal.NewFromString("-100")
	assert.True(t, antcVrss.Equal(res.Output.IntrAntcCntgVrss))
	assert.Equal(t, "5", res.Output.IntrAntcCntgVrssSign)
	assert.InDelta(t, -0.17, res.Output.IntrAntcCntgPrdyCtrt, 0.001)
	assert.Equal(t, int64(123456), res.Output.IntrAntcVol)
}
