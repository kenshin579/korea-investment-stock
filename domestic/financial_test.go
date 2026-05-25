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
	assert.True(t, decimal.NewFromInt(279600000).Equal(res.Output[0].SaleAccount))
	assert.True(t, decimal.NewFromInt(176000000).Equal(res.Output[0].SaleCost))
	assert.True(t, decimal.NewFromInt(32830000).Equal(res.Output[0].BsopPrti))
	assert.True(t, decimal.NewFromInt(23456000).Equal(res.Output[0].ThtrNtin))
}

func TestClient_InquireBalanceSheet(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/finance/balance-sheet`,
		httpmock.NewStringResponder(200, loadFixtureString(t, "balance_sheet_success.json")),
	)

	c := newTestClient(t)
	res, err := c.InquireBalanceSheet(context.Background(), domestic.InquireBalanceSheetParams{
		Symbol: "005930",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	require.Len(t, res.Output, 1)
	assert.Equal(t, "202412", res.Output[0].StacYymm)
	assert.True(t, decimal.NewFromInt(189000000).Equal(res.Output[0].Cras))
	assert.True(t, decimal.NewFromInt(434000000).Equal(res.Output[0].TotalAset))
	assert.True(t, decimal.NewFromInt(94000000).Equal(res.Output[0].TotalLblt))
	assert.True(t, decimal.NewFromInt(340000000).Equal(res.Output[0].TotalCptl))
}

func TestClient_InquireProfitRatio(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/finance/profit-ratio`,
		httpmock.NewStringResponder(200, loadFixtureString(t, "profit_ratio_success.json")),
	)

	c := newTestClient(t)
	res, err := c.InquireProfitRatio(context.Background(), domestic.InquireProfitRatioParams{
		Symbol: "005930",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	require.Len(t, res.Output, 1)
	assert.Equal(t, "202412", res.Output[0].StacYymm)
	assert.InDelta(t, 8.45, res.Output[0].CptlNtinRate, 0.001)
	assert.InDelta(t, 11.50, res.Output[0].SelfCptlNtinInrt, 0.001)
	assert.InDelta(t, 12.30, res.Output[0].SaleNtinRate, 0.001)
	assert.InDelta(t, 37.05, res.Output[0].SaleTotlRate, 0.001)
}

func TestClient_InquireGrowthRatio(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/finance/growth-ratio`,
		httpmock.NewStringResponder(200, loadFixtureString(t, "growth_ratio_success.json")),
	)

	c := newTestClient(t)
	res, err := c.InquireGrowthRatio(context.Background(), domestic.InquireGrowthRatioParams{
		Symbol: "005930",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	require.Len(t, res.Output, 1)
	assert.Equal(t, "202412", res.Output[0].StacYymm)
	assert.InDelta(t, 5.42, res.Output[0].Grs, 0.001)
	assert.InDelta(t, 12.30, res.Output[0].BsopPrfiInrt, 0.001)
	assert.InDelta(t, 8.50, res.Output[0].EqutInrt, 0.001)
	assert.InDelta(t, 10.20, res.Output[0].TotlAsetInrt, 0.001)
}

func TestClient_InquireOtherMajorRatios(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/finance/other-major-ratios`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "other_major_ratios_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireOtherMajorRatios(context.Background(), domestic.InquireOtherMajorRatiosParams{
		Symbol: "005930",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	// fid_div_cls_code 가 소문자임을 검증 (FID_DIV_CLS_CODE 가 아님)
	assert.Equal(t, "0", capturedQuery.Get("fid_div_cls_code"))
	assert.Equal(t, "", capturedQuery.Get("FID_DIV_CLS_CODE"))
	assert.Equal(t, "J", capturedQuery.Get("fid_cond_mrkt_div_code"))
	assert.Equal(t, "005930", capturedQuery.Get("fid_input_iscd"))

	require.Len(t, res.Output, 2)
	assert.Equal(t, "202412", res.Output[0].StacYymm)
	assert.Equal(t, "99.99", res.Output[0].PayoutRate) // string 보존 (KIS 비정상 출력)
	expectedEva, _ := decimal.NewFromString("1234567890")
	assert.True(t, expectedEva.Equal(res.Output[0].Eva))
	expectedEbitda, _ := decimal.NewFromString("9876543210")
	assert.True(t, expectedEbitda.Equal(res.Output[0].Ebitda))
	assert.InDelta(t, 8.45, res.Output[0].EvEbitda, 0.001)
}

func TestClient_InquireOtherMajorRatios_Quarter(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/finance/other-major-ratios`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "other_major_ratios_success.json")), nil
		},
	)

	c := newTestClient(t)
	_, err := c.InquireOtherMajorRatios(context.Background(), domestic.InquireOtherMajorRatiosParams{
		Symbol:  "005930",
		Quarter: true,
	})
	require.NoError(t, err)
	assert.Equal(t, "1", capturedQuery.Get("fid_div_cls_code")) // 1=분기
}

func TestClient_InquireOtherMajorRatios_InvalidJSON(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// envelope 은 valid 하지만 output 이 array 가 아닌 string → unmarshal 실패
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/finance/other-major-ratios`,
		httpmock.NewStringResponder(200, `{"rt_cd":"0","msg_cd":"X","msg1":"x","output":"not-array"}`),
	)

	c := newTestClient(t)
	_, err := c.InquireOtherMajorRatios(context.Background(), domestic.InquireOtherMajorRatiosParams{
		Symbol: "005930",
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "OtherMajorRatios")
}
