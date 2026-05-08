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

func TestClient_InquireVolumeRank(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/quotations/volume-rank`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "volume_rank_success.json")), nil
		},
	)

	c := newTestClient(t)
	rank, err := c.InquireVolumeRank(context.Background(), domestic.InquireVolumeRankParams{
		InputISCD: "0000",
	})
	require.NoError(t, err)
	require.NotNil(t, rank)

	// 필수 query 기본값 검증
	assert.Equal(t, "J", capturedQuery.Get("FID_COND_MRKT_DIV_CODE"))
	assert.Equal(t, "20171", capturedQuery.Get("FID_COND_SCR_DIV_CODE"))
	assert.Equal(t, "0000", capturedQuery.Get("FID_INPUT_ISCD"))
	assert.Equal(t, "0", capturedQuery.Get("FID_DIV_CLS_CODE"))
	assert.Equal(t, "0", capturedQuery.Get("FID_BLNG_CLS_CODE"))
	assert.Equal(t, "111111111", capturedQuery.Get("FID_TRGT_CLS_CODE"))
	assert.Equal(t, "0000000000", capturedQuery.Get("FID_TRGT_EXLS_CLS_CODE"))

	// 응답 검증
	require.Len(t, rank.Output, 2)
	assert.Equal(t, "삼성전자", rank.Output[0].HtsKorIsnm)
	assert.Equal(t, "005930", rank.Output[0].MkscShrnIscd)
	assert.Equal(t, "1", rank.Output[0].DataRank)
	assert.Equal(t, decimal.NewFromInt(75800), rank.Output[0].StckPrpr)
	assert.Equal(t, int64(12345678), rank.Output[0].AcmlVol)
	assert.InDelta(t, 0.21, rank.Output[0].VolTnrt, 0.001)
	assert.Equal(t, int64(938223456000), rank.Output[0].AcmlTrPbmn)
}

func TestClient_InquireVolumeRank_Variant(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/quotations/volume-rank`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "volume_rank_success.json")), nil
		},
	)

	c := newTestClient(t)
	rank, err := c.InquireVolumeRank(context.Background(), domestic.InquireVolumeRankParams{
		MarketCode: "NX",
		InputISCD:  "0000",
		DivCode:    "1", // 보통주
		BelongCode: "3", // 거래금액순
	})
	require.NoError(t, err)

	// Override 검증
	assert.Equal(t, "NX", capturedQuery.Get("FID_COND_MRKT_DIV_CODE"))
	assert.Equal(t, "1", capturedQuery.Get("FID_DIV_CLS_CODE"))
	assert.Equal(t, "3", capturedQuery.Get("FID_BLNG_CLS_CODE"))

	// negative decimal 검증
	require.GreaterOrEqual(t, len(rank.Output), 1)
	assert.True(t, rank.Output[0].PrdyVrss.IsNegative(), "PrdyVrss=-200 must be negative")
}

func TestClient_InquireMarketCap(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/ranking/market-cap`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "market_cap_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireMarketCap(context.Background(), domestic.InquireMarketCapParams{
		InputISCD: "0000",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, "J", capturedQuery.Get("fid_cond_mrkt_div_code"))
	assert.Equal(t, "20174", capturedQuery.Get("fid_cond_scr_div_code"))
	assert.Equal(t, "0000", capturedQuery.Get("fid_input_iscd"))

	require.Len(t, res.Output, 2)
	assert.Equal(t, "삼성전자", res.Output[0].HtsKorIsnm)
	assert.Equal(t, decimal.NewFromInt(75800), res.Output[0].StckPrpr)
	assert.Equal(t, int64(452329543), res.Output[0].StckAvls)
	assert.InDelta(t, 20.45, res.Output[0].MrktWholAvlsRlim, 0.001)
}

func TestClient_InquireDividendRate(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/ranking/dividend-rate`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "dividend_rate_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireDividendRate(context.Background(), domestic.InquireDividendRateParams{
		Sector:   "0001",
		FromDate: "20250101",
		ToDate:   "20251231",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	// 기본 query 검증
	assert.Equal(t, "0", capturedQuery.Get("GB1"))
	assert.Equal(t, "0001", capturedQuery.Get("UPJONG"))
	assert.Equal(t, "0", capturedQuery.Get("GB2"))
	assert.Equal(t, "1", capturedQuery.Get("GB3"))
	assert.Equal(t, "20250101", capturedQuery.Get("F_DT"))
	assert.Equal(t, "20251231", capturedQuery.Get("T_DT"))
	assert.Equal(t, "0", capturedQuery.Get("GB4"))

	require.Len(t, res.Output1, 2)
	assert.Equal(t, "1", res.Output1[0].Rank)
	assert.Equal(t, "005930", res.Output1[0].ShtCd)
	assert.Equal(t, "삼성전자", res.Output1[0].IsinName)
	assert.Equal(t, "20251231", res.Output1[0].RecordDate)
	assert.Equal(t, decimal.NewFromInt(1444), res.Output1[0].PerStoDiviAmt)
	assert.InDelta(t, 1.91, res.Output1[0].DiviRate, 0.001)
	assert.Equal(t, "현금배당", res.Output1[0].DiviKind)
}

func TestClient_InquireFluctuation(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/ranking/fluctuation`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "fluctuation_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireFluctuation(context.Background(), domestic.InquireFluctuationParams{
		InputISCD: "0000",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, "J", capturedQuery.Get("fid_cond_mrkt_div_code"))
	assert.Equal(t, "20170", capturedQuery.Get("fid_cond_scr_div_code"))
	assert.Equal(t, "0000", capturedQuery.Get("fid_input_iscd"))
	assert.Equal(t, "0", capturedQuery.Get("fid_rank_sort_cls_code"))

	require.Len(t, res.Output, 1)
	assert.Equal(t, "005930", res.Output[0].StckShrnIscd)
	assert.Equal(t, "삼성전자", res.Output[0].HtsKorIsnm)
	assert.Equal(t, decimal.NewFromInt(75800), res.Output[0].StckPrpr)
	assert.Equal(t, decimal.NewFromInt(76200), res.Output[0].StckHgpr)
	assert.Equal(t, decimal.NewFromInt(75500), res.Output[0].StckLwpr)
	assert.Equal(t, "131542", res.Output[0].HgprHour)
}

func TestClient_InquireFinanceRatioRanking(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/ranking/finance-ratio`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "finance_ratio_ranking_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireFinanceRatioRanking(context.Background(), domestic.InquireFinanceRatioRankingParams{
		Year:     "2024",
		Period:   "3",  // 결산
		RankSort: "11", // 안정성
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	// 5 개 hardcoded params 검증
	assert.Equal(t, "0", capturedQuery.Get("fid_trgt_cls_code"))
	assert.Equal(t, "20175", capturedQuery.Get("fid_cond_scr_div_code"))
	assert.Equal(t, "0", capturedQuery.Get("fid_div_cls_code"))
	assert.Equal(t, "0", capturedQuery.Get("fid_blng_cls_code"))
	assert.Equal(t, "0", capturedQuery.Get("fid_trgt_exls_cls_code"))
	// default 값 검증
	assert.Equal(t, "J", capturedQuery.Get("fid_cond_mrkt_div_code"))
	assert.Equal(t, "0000", capturedQuery.Get("fid_input_iscd"))
	// 사용자 입력 검증
	assert.Equal(t, "2024", capturedQuery.Get("fid_input_option_1"))
	assert.Equal(t, "3", capturedQuery.Get("fid_input_option_2"))
	assert.Equal(t, "11", capturedQuery.Get("fid_rank_sort_cls_code"))

	// 응답 검증
	require.Len(t, res.Output, 2)
	assert.Equal(t, int64(1), res.Output[0].DataRank)
	assert.Equal(t, "삼성전자", res.Output[0].HtsKorIsnm)
	assert.Equal(t, "005930", res.Output[0].MkscShrnIscd)
	assert.Equal(t, decimal.NewFromInt(75800), res.Output[0].StckPrpr)
	assert.Equal(t, decimal.NewFromInt(-200), res.Output[0].PrdyVrss)
	assert.Equal(t, "5", res.Output[0].PrdyVrssSign)
	assert.InDelta(t, -0.26, res.Output[0].PrdyCtrt, 0.001)
	assert.Equal(t, int64(12345678), res.Output[0].AcmlVol)
	assert.InDelta(t, 65.30, res.Output[0].Bis, 0.001)
	assert.InDelta(t, 53.75, res.Output[0].LbltRate, 0.001)
	assert.InDelta(t, 5.42, res.Output[0].Grs, 0.001)
	assert.Equal(t, "12", res.Output[0].StacMonth)
	assert.Equal(t, "0", res.Output[0].StacMonthClsCode)
	assert.Equal(t, int64(30), res.Output[0].IqryCsnu)
}

func TestClient_InquireFinanceRatioRanking_Overrides(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/ranking/finance-ratio`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "finance_ratio_ranking_success.json")), nil
		},
	)

	c := newTestClient(t)
	_, err := c.InquireFinanceRatioRanking(context.Background(), domestic.InquireFinanceRatioRankingParams{
		MarketCode: "NX",
		InputISCD:  "1001",
		PriceFrom:  "10000",
		PriceTo:    "100000",
		VolFrom:    "100000",
		Year:       "2024",
		Period:     "0",  // 1Q
		RankSort:   "20", // 활동성
	})
	require.NoError(t, err)
	assert.Equal(t, "NX", capturedQuery.Get("fid_cond_mrkt_div_code"))
	assert.Equal(t, "1001", capturedQuery.Get("fid_input_iscd"))
	assert.Equal(t, "10000", capturedQuery.Get("fid_input_price_1"))
	assert.Equal(t, "100000", capturedQuery.Get("fid_input_price_2"))
	assert.Equal(t, "100000", capturedQuery.Get("fid_vol_cnt"))
}

func TestClient_InquireFinanceRatioRanking_InvalidJSON(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// envelope 은 valid 하지만 output 이 array 가 아닌 string → unmarshal 실패
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/ranking/finance-ratio`,
		httpmock.NewStringResponder(200, `{"rt_cd":"0","msg_cd":"X","msg1":"x","output":"not-array"}`),
	)

	c := newTestClient(t)
	_, err := c.InquireFinanceRatioRanking(context.Background(), domestic.InquireFinanceRatioRankingParams{
		Year:     "2024",
		Period:   "3",
		RankSort: "11",
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "FinanceRatioRanking")
}

func TestClient_InquireTradedByCompany(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/ranking/traded-by-company`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "traded_by_company_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireTradedByCompany(context.Background(), domestic.InquireTradedByCompanyParams{
		InputDate1: "20260501",
		InputDate2: "20260508",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	// 4 hardcoded
	assert.Equal(t, "0", capturedQuery.Get("fid_trgt_exls_cls_code"))
	assert.Equal(t, "20186", capturedQuery.Get("fid_cond_scr_div_code"))
	assert.Equal(t, "0", capturedQuery.Get("fid_trgt_cls_code"))
	assert.Equal(t, "0", capturedQuery.Get("fid_aply_rang_vol"))
	// default 값
	assert.Equal(t, "J", capturedQuery.Get("fid_cond_mrkt_div_code"))
	assert.Equal(t, "0", capturedQuery.Get("fid_div_cls_code"))
	assert.Equal(t, "0", capturedQuery.Get("fid_rank_sort_cls_code"))
	assert.Equal(t, "0000", capturedQuery.Get("fid_input_iscd"))
	// 사용자 입력
	assert.Equal(t, "20260501", capturedQuery.Get("fid_input_date_1"))
	assert.Equal(t, "20260508", capturedQuery.Get("fid_input_date_2"))

	require.Len(t, res.Output, 2)
	assert.Equal(t, int64(1), res.Output[0].DataRank)
	assert.Equal(t, "005930", res.Output[0].MkscShrnIscd)
	assert.Equal(t, "삼성전자", res.Output[0].HtsKorIsnm)
	assert.Equal(t, decimal.NewFromInt(75500), res.Output[0].StckPrpr)
	assert.Equal(t, decimal.NewFromInt(1500), res.Output[0].PrdyVrss)
	assert.InDelta(t, 2.03, res.Output[0].PrdyCtrt, 0.001)
	assert.Equal(t, int64(12345678), res.Output[0].AcmlVol)
	assert.Equal(t, int64(987654321000), res.Output[0].AcmlTrPbmn)
	assert.Equal(t, int64(5000000), res.Output[0].SelnCnqnSmtn)
	assert.Equal(t, int64(7345678), res.Output[0].ShnuCnqnSmtn)
	assert.Equal(t, int64(2345678), res.Output[0].NtbyCnqn)
	// 두 번째 행 — 음수 순매수
	assert.Equal(t, int64(-1234567), res.Output[1].NtbyCnqn)
}

func TestClient_InquireTradedByCompany_Overrides(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/ranking/traded-by-company`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "traded_by_company_success.json")), nil
		},
	)

	c := newTestClient(t)
	_, err := c.InquireTradedByCompany(context.Background(), domestic.InquireTradedByCompanyParams{
		MarketCode: "NX",
		DivCode:    "6", // 보통주
		SortCode:   "1", // 매수상위
		InputDate1: "20260101",
		InputDate2: "20260131",
		InputISCD:  "1001", // 코스닥
		PriceFrom:  "10000",
		PriceTo:    "100000",
	})
	require.NoError(t, err)
	assert.Equal(t, "NX", capturedQuery.Get("fid_cond_mrkt_div_code"))
	assert.Equal(t, "6", capturedQuery.Get("fid_div_cls_code"))
	assert.Equal(t, "1", capturedQuery.Get("fid_rank_sort_cls_code"))
	assert.Equal(t, "1001", capturedQuery.Get("fid_input_iscd"))
	assert.Equal(t, "10000", capturedQuery.Get("fid_aply_rang_prc_1"))
	assert.Equal(t, "100000", capturedQuery.Get("fid_aply_rang_prc_2"))
}

func TestClient_InquireTradedByCompany_InvalidJSON(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/ranking/traded-by-company`,
		httpmock.NewStringResponder(200, `{"rt_cd":"0","msg_cd":"X","msg1":"x","output":"not-array"}`),
	)

	c := newTestClient(t)
	_, err := c.InquireTradedByCompany(context.Background(), domestic.InquireTradedByCompanyParams{})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "TradedByCompany")
}
