package domestic_test

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/kenshin579/korea-investment-stock/domestic"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_InquireInvestOpinion(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/quotations/invest-opinion`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "invest_opinion_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireInvestOpinion(context.Background(), domestic.InquireInvestOpinionParams{
		Symbol:    "005930",
		StartDate: "20260401",
		EndDate:   "20260506",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, "J", capturedQuery.Get("FID_COND_MRKT_DIV_CODE"))
	assert.Equal(t, "16633", capturedQuery.Get("FID_COND_SCR_DIV_CODE"))
	assert.Equal(t, "005930", capturedQuery.Get("FID_INPUT_ISCD"))
	assert.Equal(t, "20260401", capturedQuery.Get("FID_INPUT_DATE_1"))
	assert.Equal(t, "20260506", capturedQuery.Get("FID_INPUT_DATE_2"))

	require.Len(t, res.Output, 2)
	assert.Equal(t, "20260506", res.Output[0].StckBsopDate)
	assert.Equal(t, "삼성증권", res.Output[0].MbcrName)
	assert.Equal(t, "매수", res.Output[0].InvtOpnn)

	wantGoal, _ := decimal.NewFromString("95000")
	assert.True(t, wantGoal.Equal(res.Output[0].HtsGoalPrc))
	assert.InDelta(t, 12.73, res.Output[0].NdayDprt, 0.001)
	assert.InDelta(t, 9.09, res.Output[0].Dprt, 0.001)
}

func TestClient_InquireInvestOpbysec(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/quotations/invest-opbysec`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "invest_opbysec_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireInvestOpbysec(context.Background(), domestic.InquireInvestOpbysecParams{
		SecBrokerCode: "240", // 삼성증권 코드 예시
		DivClsCode:    "0",
		StartDate:     "20260401",
		EndDate:       "20260506",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	// 증권사코드가 FID_INPUT_ISCD 로 전송되는지 확인
	assert.Equal(t, "J", capturedQuery.Get("FID_COND_MRKT_DIV_CODE"))
	assert.Equal(t, "16634", capturedQuery.Get("FID_COND_SCR_DIV_CODE"))
	assert.Equal(t, "240", capturedQuery.Get("FID_INPUT_ISCD"))
	assert.Equal(t, "0", capturedQuery.Get("FID_DIV_CLS_CODE"))

	require.Len(t, res.Output, 2)
	assert.Equal(t, "005930", res.Output[0].StckShrnIscd)
	assert.Equal(t, "삼성전자", res.Output[0].HtsKorIsnm)
	assert.Equal(t, "삼성증권", res.Output[0].MbcrName)

	wantPrpr, _ := decimal.NewFromString("82500")
	assert.True(t, wantPrpr.Equal(res.Output[0].StckPrpr))
	assert.InDelta(t, 0.61, res.Output[0].PrdyCtrt, 0.001)
	assert.InDelta(t, 9.76, res.Output[0].Dprt, 0.001)
}
