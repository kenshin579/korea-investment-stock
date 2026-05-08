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

func TestClient_InquireCompInterest(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/quotations/comp-interest`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "comp_interest_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireCompInterest(context.Background(), domestic.InquireCompInterestParams{})
	require.NoError(t, err)
	require.NotNil(t, res)

	// 4 hardcoded params 모두 UPPERCASE 검증
	assert.Equal(t, "I", capturedQuery.Get("FID_COND_MRKT_DIV_CODE"))
	assert.Equal(t, "20702", capturedQuery.Get("FID_COND_SCR_DIV_CODE"))
	assert.Equal(t, "1", capturedQuery.Get("FID_DIV_CLS_CODE"))
	// FID_DIV_CLS_CODE1 은 빈 값 — Get 결과 ""
	assert.Equal(t, "", capturedQuery.Get("FID_DIV_CLS_CODE1"))

	// output1 (single object) 검증
	assert.Equal(t, "01001", res.Output1.BcdtCode)
	assert.Equal(t, "국고채3년", res.Output1.HtsKorIsnm)
	expected, _ := decimal.NewFromString("3.421")
	assert.True(t, expected.Equal(res.Output1.BondMnrtPrpr))
	assert.Equal(t, "5", res.Output1.PrdyVrssSign)
	assert.InDelta(t, -0.35, res.Output1.PrdyCtrt, 0.001)
	assert.Equal(t, "20260508", res.Output1.StckBsopDate)

	// output2 (array) 검증
	require.Len(t, res.Output2, 2)
	assert.Equal(t, "01002", res.Output2[1].BcdtCode)
	assert.Equal(t, "국고채5년", res.Output2[1].HtsKorIsnm)
	assert.InDelta(t, 0.14, res.Output2[1].BstpNmixPrdyCtrt, 0.001)
}

func TestClient_InquireCompInterest_InvalidJSON(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// envelope 은 valid 하지만 output1 이 단일 object 가 아닌 array → unmarshal 실패
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/quotations/comp-interest`,
		httpmock.NewStringResponder(200, `{"rt_cd":"0","msg_cd":"X","msg1":"x","output1":["not-object"],"output2":[]}`),
	)

	c := newTestClient(t)
	_, err := c.InquireCompInterest(context.Background(), domestic.InquireCompInterestParams{})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "CompInterest")
}
