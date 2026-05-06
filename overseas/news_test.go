package overseas_test

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kenshin579/korea-investment-stock/overseas"
)

func TestClient_InquireNewsTitle(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/quotations/news-title`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "news_title_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireNewsTitle(context.Background(), overseas.InquireNewsTitleParams{
		NationCd:   "US",
		ExchangeCd: "NAS",
		DataDt:     "20260505",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, "US", capturedQuery.Get("NATION_CD"))
	assert.Equal(t, "NAS", capturedQuery.Get("EXCHANGE_CD"))
	assert.Equal(t, "20260505", capturedQuery.Get("DATA_DT"))

	require.Len(t, res.Outblock1, 2)
	assert.Equal(t, "AAPL", res.Outblock1[0].Symb)
	assert.Equal(t, "20260505120001", res.Outblock1[0].NewsKey)
	assert.Equal(t, "애플 2분기 실적 서프라이즈 발표", res.Outblock1[0].Title)
}

func TestClient_InquireBrknewsTitle(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/quotations/brknews-title`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "brknews_title_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireBrknewsTitle(context.Background(), overseas.InquireBrknewsTitleParams{
		InputDate1: "20260505",
		Symbol:     "AAPL",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	// FID_ prefix wire name 검증
	assert.Equal(t, "20260505", capturedQuery.Get("FID_INPUT_DATE_1"))
	assert.Equal(t, "AAPL", capturedQuery.Get("FID_INPUT_ISCD"))
	// hardcoded param 검증
	assert.Equal(t, "11801", capturedQuery.Get("FID_COND_SCR_DIV_CODE"))

	require.Len(t, res.Output, 2)
	assert.Equal(t, "미 연준 금리 동결 결정", res.Output[0].HtsPbntTitlCntt)
	assert.Equal(t, "AAPL", res.Output[0].Iscd1)
	assert.Equal(t, "애플", res.Output[0].KorIsnm1)
	// 빈 iscd/kor_isnm 필드 확인
	assert.Equal(t, "", res.Output[0].Iscd3)
	assert.Equal(t, "", res.Output[0].KorIsnm3)
}
