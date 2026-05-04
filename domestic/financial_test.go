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

func TestClient_InquireFinancialRatio(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/finance/financial-ratio`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "financial_ratio_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireFinancialRatio(context.Background(), domestic.InquireFinancialRatioParams{
		Symbol: "005930",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, "0", capturedQuery.Get("FID_DIV_CLS_CODE"))
	assert.Equal(t, "J", capturedQuery.Get("fid_cond_mrkt_div_code"))
	assert.Equal(t, "005930", capturedQuery.Get("fid_input_iscd"))

	require.Len(t, res.Output, 2)
	assert.Equal(t, "202412", res.Output[0].StacYymm)
	assert.InDelta(t, 5.42, res.Output[0].Grs, 0.001)
	assert.InDelta(t, 11.50, res.Output[0].RoeVal, 0.001)
	assert.Equal(t, decimal.NewFromInt(6638), res.Output[0].Eps)
	assert.Equal(t, decimal.NewFromInt(57420), res.Output[0].Bps)
}

func TestClient_InquireFinancialRatio_Quarter(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/finance/financial-ratio`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "financial_ratio_success.json")), nil
		},
	)

	c := newTestClient(t)
	_, err := c.InquireFinancialRatio(context.Background(), domestic.InquireFinancialRatioParams{
		Symbol:  "005930",
		Quarter: true,
	})
	require.NoError(t, err)
	assert.Equal(t, "1", capturedQuery.Get("FID_DIV_CLS_CODE")) // 1=분기
}

func TestClient_InquireIncomeStatement(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/finance/income-statement`,
		httpmock.NewStringResponder(200, loadFixtureString(t, "income_statement_success.json")),
	)

	c := newTestClient(t)
	res, err := c.InquireIncomeStatement(context.Background(), domestic.InquireIncomeStatementParams{
		Symbol: "005930",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	require.Len(t, res.Output, 1)
	assert.Equal(t, "202412", res.Output[0].StacYymm)
	assert.Equal(t, int64(279600000), res.Output[0].SaleAccount)
	assert.Equal(t, int64(176000000), res.Output[0].SaleCost)
	assert.Equal(t, int64(32830000), res.Output[0].BsopPrti)
	assert.Equal(t, int64(23456000), res.Output[0].ThtrNtin)
}
