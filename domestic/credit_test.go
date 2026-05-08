package domestic_test

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kenshin579/korea-investment-stock/domestic"
)

func TestClient_InquireCreditByCompany(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/quotations/credit-by-company`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "credit_by_company_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireCreditByCompany(context.Background(), domestic.InquireCreditByCompanyParams{})
	require.NoError(t, err)
	require.NotNil(t, res)

	// 2 hardcoded
	assert.Equal(t, "20477", capturedQuery.Get("fid_cond_scr_div_code"))
	assert.Equal(t, "J", capturedQuery.Get("fid_cond_mrkt_div_code"))
	// 3 default 값
	assert.Equal(t, "0", capturedQuery.Get("fid_rank_sort_cls_code"))
	assert.Equal(t, "0", capturedQuery.Get("fid_slct_yn"))
	assert.Equal(t, "0000", capturedQuery.Get("fid_input_iscd"))

	require.Len(t, res.Output, 3)
	assert.Equal(t, "005930", res.Output[0].StckShrnIscd)
	assert.Equal(t, "삼성전자", res.Output[0].HtsKorIsnm)
	assert.InDelta(t, 60.00, res.Output[0].CrdtRate, 0.001)
	assert.Equal(t, "NAVER", res.Output[2].HtsKorIsnm)
	assert.InDelta(t, 55.00, res.Output[2].CrdtRate, 0.001)
}

func TestClient_InquireCreditByCompany_Overrides(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/quotations/credit-by-company`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "credit_by_company_success.json")), nil
		},
	)

	c := newTestClient(t)
	_, err := c.InquireCreditByCompany(context.Background(), domestic.InquireCreditByCompanyParams{
		SortCode:  "1",    // 이름순
		SelectYN:  "1",    // 신용주문불가
		InputISCD: "1001", // 코스닥
	})
	require.NoError(t, err)
	assert.Equal(t, "1", capturedQuery.Get("fid_rank_sort_cls_code"))
	assert.Equal(t, "1", capturedQuery.Get("fid_slct_yn"))
	assert.Equal(t, "1001", capturedQuery.Get("fid_input_iscd"))
}

func TestClient_InquireCreditByCompany_InvalidJSON(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/quotations/credit-by-company`,
		httpmock.NewStringResponder(200, `{"rt_cd":"0","msg_cd":"X","msg1":"x","output":"not-array"}`),
	)

	c := newTestClient(t)
	_, err := c.InquireCreditByCompany(context.Background(), domestic.InquireCreditByCompanyParams{})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "CreditByCompany")
}
