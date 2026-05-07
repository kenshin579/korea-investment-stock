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

func TestClient_InquireNearNewHighlow(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/ranking/near-new-highlow`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "near_new_highlow_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireNearNewHighlow(context.Background(), domestic.InquireNearNewHighlowParams{
		InputISCD:  "0000",
		PrcClsCode: "0",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	// query param 검증
	assert.Equal(t, "J", capturedQuery.Get("fid_cond_mrkt_div_code"))
	assert.Equal(t, "20187", capturedQuery.Get("fid_cond_scr_div_code"))
	assert.Equal(t, "0000", capturedQuery.Get("fid_input_iscd"))
	assert.Equal(t, "0", capturedQuery.Get("fid_prc_cls_code"))

	require.Len(t, res.Output, 2)

	// output[0] 필드 검증
	assert.Equal(t, "삼성전자", res.Output[0].HtsKorIsnm)
	assert.Equal(t, "005930", res.Output[0].MkscShrnIscd)
	d, _ := decimal.NewFromString("75800")
	assert.True(t, d.Equal(res.Output[0].StckPrpr))
	assert.Equal(t, int64(12345678), res.Output[0].AcmlVol)
	d2, _ := decimal.NewFromString("76500")
	assert.True(t, d2.Equal(res.Output[0].NewHgpr))
	assert.InDelta(t, 1.24, res.Output[0].HprcNearRate, 0.01)
	d3, _ := decimal.NewFromString("74000")
	assert.True(t, d3.Equal(res.Output[0].NewLwpr))
	assert.InDelta(t, 2.43, res.Output[0].LwprNearRate, 0.01)
}

func TestClient_InquireOvertimePrice(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/quotations/inquire-overtime-price`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "overtime_price_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireOvertimePrice(context.Background(), domestic.InquireOvertimePriceParams{
		Symbol: "005930",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, "J", capturedQuery.Get("FID_COND_MRKT_DIV_CODE"))
	assert.Equal(t, "005930", capturedQuery.Get("FID_INPUT_ISCD"))

	// output 필드 검증
	assert.Equal(t, "전기전자", res.Output.BstpKorIsnm)
	d, _ := decimal.NewFromString("75700")
	assert.True(t, d.Equal(res.Output.OvtmUntpPrpr))
	assert.Equal(t, int64(234567), res.Output.OvtmUntpVol)
	assert.InDelta(t, 20.00, res.Output.MargRate, 0.01)
	assert.Equal(t, "N", res.Output.TrhtYn)
	assert.Equal(t, "KOSPI", res.Output.RprsMrktKorName)
	d2, _ := decimal.NewFromString("75700")
	assert.True(t, d2.Equal(res.Output.Bidp))
}

func TestClient_InquireOvertimeAskingPrice(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/quotations/inquire-overtime-asking-price`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "overtime_asking_price_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireOvertimeAskingPrice(context.Background(), domestic.InquireOvertimeAskingPriceParams{
		Symbol: "005930",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, "J", capturedQuery.Get("FID_COND_MRKT_DIV_CODE"))
	assert.Equal(t, "005930", capturedQuery.Get("FID_INPUT_ISCD"))

	// output1 핵심 필드 검증
	assert.Equal(t, "180542", res.Output1.OvtmUntpLastHour)
	d, _ := decimal.NewFromString("75750")
	assert.True(t, d.Equal(res.Output1.OvtmUntpAskp1))
	d2, _ := decimal.NewFromString("75700")
	assert.True(t, d2.Equal(res.Output1.OvtmUntpBidp1))
	assert.Equal(t, int64(100), res.Output1.OvtmUntpAskpIcdc1)
	assert.Equal(t, int64(200), res.Output1.OvtmUntpBidpIcdc1)
	assert.Equal(t, int64(1200), res.Output1.OvtmUntpAskpRsqn1)
	assert.Equal(t, int64(2000), res.Output1.OvtmUntpBidpRsqn1)
	assert.Equal(t, int64(6150), res.Output1.OvtmUntpTotalAskpRsqn)
	assert.Equal(t, int64(8100), res.Output1.OvtmUntpTotalBidpRsqn)
	assert.Equal(t, int64(11100), res.Output1.TotalAskpRsqn)
	assert.Equal(t, int64(13400), res.Output1.TotalBidpRsqn)
}

func TestClient_InquireOvertimeVolume(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/ranking/overtime-volume`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "overtime_volume_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireOvertimeVolume(context.Background(), domestic.InquireOvertimeVolumeParams{
		InputISCD: "0000",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, "J", capturedQuery.Get("FID_COND_MRKT_DIV_CODE"))
	assert.Equal(t, "20235", capturedQuery.Get("FID_COND_SCR_DIV_CODE"))
	assert.Equal(t, "0000", capturedQuery.Get("FID_INPUT_ISCD"))

	// output1 검증
	assert.Equal(t, int64(12345678), res.Output1.OvtmUntpExchVol)
	assert.Equal(t, int64(9876543), res.Output1.OvtmUntpKosdaqVol)

	// output2 검증
	require.Len(t, res.Output2, 2)
	assert.Equal(t, "005930", res.Output2[0].StckShrnIscd)
	d, _ := decimal.NewFromString("75700")
	assert.True(t, d.Equal(res.Output2[0].OvtmUntpPrpr))
	assert.Equal(t, int64(234567), res.Output2[0].OvtmUntpVol)
	assert.InDelta(t, 1.90, res.Output2[0].OvtmVrssAcmlVolRlim, 0.01)
}

func TestClient_InquireOvertimeFluctuation(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/ranking/overtime-fluctuation`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "overtime_fluctuation_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireOvertimeFluctuation(context.Background(), domestic.InquireOvertimeFluctuationParams{
		InputISCD:  "0000",
		DivClsCode: "2",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, "J", capturedQuery.Get("FID_COND_MRKT_DIV_CODE"))
	assert.Equal(t, "20234", capturedQuery.Get("FID_COND_SCR_DIV_CODE"))
	assert.Equal(t, "0000", capturedQuery.Get("FID_INPUT_ISCD"))
	assert.Equal(t, "2", capturedQuery.Get("FID_DIV_CLS_CODE"))

	// output1 검증
	assert.Equal(t, int64(5), res.Output1.OvtmUntpUplmIssuCnt)
	assert.Equal(t, int64(312), res.Output1.OvtmUntpAscnIssuCnt)
	assert.Equal(t, int64(22345678), res.Output1.OvtmUntpAcmlVol)

	// output2 검증
	require.Len(t, res.Output2, 2)
	assert.Equal(t, "005930", res.Output2[0].MkscShrnIscd)
	d, _ := decimal.NewFromString("75700")
	assert.True(t, d.Equal(res.Output2[0].OvtmUntpPrpr))
	assert.InDelta(t, -0.39, res.Output2[0].OvtmUntpPrdyCtrt, 0.01)
	assert.Equal(t, int64(234567), res.Output2[0].OvtmUntpVol)
	assert.InDelta(t, 1.90, res.Output2[0].OvtmVrssAcmlVolRlim, 0.01)
}

func TestClient_InquireVolumePower(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/ranking/volume-power`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "volume_power_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireVolumePower(context.Background(), domestic.InquireVolumePowerParams{
		Symbol:     "0001",
		DivClsCode: "0",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	// lowercase wire keys 확인 (UPPERCASE FID_ 아님)
	assert.Equal(t, "J", capturedQuery.Get("fid_cond_mrkt_div_code"))
	assert.Equal(t, "20168", capturedQuery.Get("fid_cond_scr_div_code"))
	assert.Equal(t, "0001", capturedQuery.Get("fid_input_iscd"))
	assert.Empty(t, capturedQuery.Get("FID_COND_MRKT_DIV_CODE"), "lowercase 사용 확인")

	require.Len(t, res.Output, 2)
	assert.Equal(t, "005930", res.Output[0].StckShrnIscd)
	assert.Equal(t, "1", res.Output[0].DataRank)

	wantPrpr, _ := decimal.NewFromString("82500")
	assert.True(t, wantPrpr.Equal(res.Output[0].StckPrpr))
	assert.InDelta(t, 125.30, res.Output[0].TdayRltv, 0.001)
	assert.Equal(t, int64(6800000), res.Output[0].SelnCnqnSmtn)
	assert.Equal(t, int64(7200000), res.Output[0].ShnuCnqnSmtn)
}

func TestClient_InquireBulkTransNum(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/ranking/bulk-trans-num`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "bulk_trans_num_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireBulkTransNum(context.Background(), domestic.InquireBulkTransNumParams{
		Symbol:       "0000",
		DivClsCode:   "0",
		RankSortCode: "0",
		BlngClsCode:  "0",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	// lowercase wire keys 확인
	assert.Equal(t, "J", capturedQuery.Get("fid_cond_mrkt_div_code"))
	assert.Equal(t, "11909", capturedQuery.Get("fid_cond_scr_div_code"))
	assert.Equal(t, "0000", capturedQuery.Get("fid_input_iscd"))

	require.Len(t, res.Output, 2)
	// mksc_shrn_iscd (stck_shrn_iscd 아님) 확인
	assert.Equal(t, "005930", res.Output[0].MkscShrnIscd)
	assert.Equal(t, "1", res.Output[0].DataRank)
	assert.Equal(t, int64(3200), res.Output[0].ShnuCntgCsnu)
	assert.Equal(t, int64(2800), res.Output[0].SelnCntgCsnu)
	assert.Equal(t, int64(400000), res.Output[0].NtbyCnqn)
}

func TestClient_InquireTradprtByamt(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/quotations/tradprt-byamt`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "tradprt_byamt_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireTradprtByamt(context.Background(), domestic.InquireTradprtByamtParams{
		Symbol: "005930",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, "J", capturedQuery.Get("FID_COND_MRKT_DIV_CODE"))
	assert.Equal(t, "11119", capturedQuery.Get("FID_COND_SCR_DIV_CODE"))
	assert.Equal(t, "005930", capturedQuery.Get("FID_INPUT_ISCD"))

	require.Len(t, res.Output, 2)
	assert.Equal(t, "1억원 이상", res.Output[0].PrprName)
	assert.Equal(t, int64(8500), res.Output[0].AcmlVol)
	assert.InDelta(t, 12.50, res.Output[0].WholNtbyQtyRate, 0.001)
	// whol_shun_vol_rate typo 필드 확인
	assert.InDelta(t, 45.20, res.Output[0].WholShunVolRate, 0.001)
	assert.InDelta(t, 42.30, res.Output[0].WholSelnVolRate, 0.001)

	wantAvrgPrpr, _ := decimal.NewFromString("150000000")
	assert.True(t, wantAvrgPrpr.Equal(res.Output[0].SmtnAvrgPrpr))
}

func TestClient_InquireHtsTopView(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(http.MethodGet, `=~/ranking/hts-top-view`,
		httpmock.NewStringResponder(http.StatusOK, loadFixtureString(t, "hts_top_view_success.json")))

	c := newTestClient(t)
	res, err := c.InquireHtsTopView(context.Background(), domestic.InquireHtsTopViewParams{})
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, "J", res.Output1.MrktDivClsCode)
	assert.Equal(t, "005930", res.Output1.MkscShrnIscd)
}

func TestClient_InquirePbarTraRatio(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(http.MethodGet, `=~/quotations/pbar-tratio`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(http.StatusOK, loadFixtureString(t, "pbar_tratio_success.json")), nil
		})

	c := newTestClient(t)
	res, err := c.InquirePbarTraRatio(context.Background(), domestic.InquirePbarTraRatioParams{
		Symbol:     "005930",
		InputHour1: "153000",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, "J", capturedQuery.Get("FID_COND_MRKT_DIV_CODE"))
	assert.Equal(t, "11130", capturedQuery.Get("FID_COND_SCR_DIV_CODE"))
	assert.Equal(t, "005930", capturedQuery.Get("FID_INPUT_ISCD"))
	assert.Equal(t, "153000", capturedQuery.Get("FID_INPUT_HOUR_1"))

	assert.Equal(t, "KOSPI", res.Output1.RprsMrktKorName)
	assert.Equal(t, "005930", res.Output1.StckShrnIscd)
	assert.Equal(t, "삼성전자", res.Output1.HtsKorIsnm)
	wantPrpr, _ := decimal.NewFromString("82500")
	assert.True(t, wantPrpr.Equal(res.Output1.StckPrpr))
	assert.Equal(t, int64(12500000), res.Output1.AcmlVol)
	assert.Equal(t, int64(11000000), res.Output1.PrdyVol)
	assert.Equal(t, int64(5969782550), res.Output1.LstnStcn)
	wantWavg, _ := decimal.NewFromString("82350")
	assert.True(t, wantWavg.Equal(res.Output1.WghnAvrgStckPrc))

	require.Len(t, res.Output2, 2)
	assert.Equal(t, "1", res.Output2[0].DataRank)
	wantItemPrpr, _ := decimal.NewFromString("82500")
	assert.True(t, wantItemPrpr.Equal(res.Output2[0].StckPrpr))
	assert.Equal(t, int64(1500000), res.Output2[0].CntgVol)
	assert.InDelta(t, 12.00, res.Output2[0].AcmlVolRlim, 0.001)
}

func TestClient_InquireExpPriceTrend(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(http.MethodGet, `=~/quotations/exp-price-trend`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(http.StatusOK, loadFixtureString(t, "exp_price_trend_success.json")), nil
		})

	c := newTestClient(t)
	res, err := c.InquireExpPriceTrend(context.Background(), domestic.InquireExpPriceTrendParams{
		Symbol: "005930",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	// lowercase fid_* 검증
	assert.Equal(t, "J", capturedQuery.Get("fid_cond_mrkt_div_code"))
	assert.Equal(t, "11810", capturedQuery.Get("fid_cond_scr_div_code"))
	assert.Equal(t, "005930", capturedQuery.Get("fid_input_iscd"))
	// UPPERCASE 는 비어 있음을 확인
	assert.Empty(t, capturedQuery.Get("FID_INPUT_ISCD"))

	// output1 검증
	assert.Equal(t, "코스피", res.Output1.RprsMrktKorName)
	wantCnpr, _ := decimal.NewFromString("82700")
	assert.True(t, wantCnpr.Equal(res.Output1.AntcCnpr))
	assert.Equal(t, int64(850000), res.Output1.AntcVol)
	assert.Equal(t, int64(70297500000), res.Output1.AntcTrPbmn)
	assert.InDelta(t, 0.24, res.Output1.AntcCntgPrdyCtrt, 0.001)

	// output2 검증
	require.Len(t, res.Output2, 2)
	assert.Equal(t, "20260507", res.Output2[0].StckBsopDate)
	assert.Equal(t, "153000", res.Output2[0].StckCntgHour)
	wantPrpr, _ := decimal.NewFromString("82500")
	assert.True(t, wantPrpr.Equal(res.Output2[0].StckPrpr))
	assert.Equal(t, int64(12500000), res.Output2[0].AcmlVol)
}

func TestClient_InquireExpTransUpdown(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(http.MethodGet, `=~/ranking/exp-trans-updown`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(http.StatusOK, loadFixtureString(t, "exp_trans_updown_success.json")), nil
		})

	c := newTestClient(t)
	res, err := c.InquireExpTransUpdown(context.Background(), domestic.InquireExpTransUpdownParams{
		Symbol: "0000",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	// lowercase fid_* 검증
	assert.Equal(t, "J", capturedQuery.Get("fid_cond_mrkt_div_code"))
	assert.Equal(t, "11820", capturedQuery.Get("fid_cond_scr_div_code"))
	assert.Equal(t, "0000", capturedQuery.Get("fid_input_iscd"))
	assert.Empty(t, capturedQuery.Get("FID_INPUT_ISCD"))

	require.Len(t, res.Output, 2)
	assert.Equal(t, "005930", res.Output[0].StckShrnIscd)
	assert.Equal(t, "삼성전자", res.Output[0].HtsKorIsnm)
	wantPrpr, _ := decimal.NewFromString("82500")
	assert.True(t, wantPrpr.Equal(res.Output[0].StckPrpr))
	assert.InDelta(t, 0.61, res.Output[0].PrdyCtrt, 0.001)
	assert.Equal(t, int64(70297500000), res.Output[0].AntcTrPbmn)
}

func TestClient_InquireShortSale(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/ranking/short-sale`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(http.StatusOK, loadFixtureString(t, "short_sale_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireShortSale(context.Background(), domestic.InquireShortSaleParams{
		Symbol:        "005930",
		PeriodDivCode: "D",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	// UPPERCASE FID_ 키 검증
	assert.Equal(t, "J", capturedQuery.Get("FID_COND_MRKT_DIV_CODE"))
	assert.Equal(t, "20482", capturedQuery.Get("FID_COND_SCR_DIV_CODE"))
	assert.Equal(t, "005930", capturedQuery.Get("FID_INPUT_ISCD"))
	assert.Equal(t, "D", capturedQuery.Get("FID_PERIOD_DIV_CODE"))

	require.Len(t, res.Output, 2)

	// output[0] 필드 검증
	assert.Equal(t, "005930", res.Output[0].MkscShrnIscd)
	assert.Equal(t, "삼성전자", res.Output[0].HtsKorIsnm)
	d, _ := decimal.NewFromString("75800")
	assert.True(t, d.Equal(res.Output[0].StckPrpr))
	dVrss, _ := decimal.NewFromString("-200")
	assert.True(t, dVrss.Equal(res.Output[0].PrdyVrss))
	assert.Equal(t, "5", res.Output[0].PrdyVrssSign)
	assert.InDelta(t, -0.26, res.Output[0].PrdyCtrt, 0.001)
	assert.Equal(t, int64(12345678), res.Output[0].AcmlVol)
	assert.Equal(t, int64(935714440400), res.Output[0].AcmlTrPbmn)
	assert.Equal(t, int64(123456), res.Output[0].SstsCntgQty)
	assert.InDelta(t, 1.00, res.Output[0].SstsVolRlim, 0.001)
	assert.Equal(t, int64(9357144404), res.Output[0].SstsTrPbmn)
	assert.InDelta(t, 1.00, res.Output[0].SstsTrPbmnRlim, 0.001)
	assert.Equal(t, "20260501", res.Output[0].StndDate1)
	assert.Equal(t, "20260507", res.Output[0].StndDate2)
	dAvrg, _ := decimal.NewFromString("75800")
	assert.True(t, dAvrg.Equal(res.Output[0].AvrgPrc))
}

func TestClient_InquireDailyShortSale(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/quotations/daily-short-sale`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(http.StatusOK, loadFixtureString(t, "daily_short_sale_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireDailyShortSale(context.Background(), domestic.InquireDailyShortSaleParams{
		Symbol:     "005930",
		InputDate1: "20260501",
		InputDate2: "20260507",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	// UPPERCASE FID_ 키 검증
	assert.Equal(t, "J", capturedQuery.Get("FID_COND_MRKT_DIV_CODE"))
	assert.Equal(t, "005930", capturedQuery.Get("FID_INPUT_ISCD"))
	assert.Equal(t, "20260501", capturedQuery.Get("FID_INPUT_DATE_1"))
	assert.Equal(t, "20260507", capturedQuery.Get("FID_INPUT_DATE_2"))

	// output1 (single obj) 검증
	d1, _ := decimal.NewFromString("75800")
	assert.True(t, d1.Equal(res.Output1.StckPrpr))
	dVrss, _ := decimal.NewFromString("-200")
	assert.True(t, dVrss.Equal(res.Output1.PrdyVrss))
	assert.Equal(t, "5", res.Output1.PrdyVrssSign)
	assert.InDelta(t, -0.26, res.Output1.PrdyCtrt, 0.001)
	assert.Equal(t, int64(12345678), res.Output1.AcmlVol)
	assert.Equal(t, int64(13456789), res.Output1.PrdyVol)

	// output2 (array) 검증
	require.Len(t, res.Output2, 2)
	assert.Equal(t, "20260507", res.Output2[0].StckBsopDate)
	d2, _ := decimal.NewFromString("75800")
	assert.True(t, d2.Equal(res.Output2[0].StckClpr))
	assert.InDelta(t, -0.26, res.Output2[0].PrdyCtrt, 0.001)
	assert.Equal(t, int64(12345678), res.Output2[0].AcmlVol)
	assert.Equal(t, int64(123456), res.Output2[0].SstsCntgQty)
	assert.InDelta(t, 1.90, res.Output2[0].AcmlSstsCntgQtyRlim, 0.001)
	assert.Equal(t, int64(935714440400), res.Output2[0].AcmlTrPbmn)
	assert.InDelta(t, 2.00, res.Output2[0].AcmlSstsTrPbmnRlim, 0.001)
	dOprc, _ := decimal.NewFromString("76000")
	assert.True(t, dOprc.Equal(res.Output2[0].StckOprc))
	dHgpr, _ := decimal.NewFromString("76200")
	assert.True(t, dHgpr.Equal(res.Output2[0].StckHgpr))
	dLwpr, _ := decimal.NewFromString("75500")
	assert.True(t, dLwpr.Equal(res.Output2[0].StckLwpr))
	dAvrg, _ := decimal.NewFromString("75800")
	assert.True(t, dAvrg.Equal(res.Output2[0].AvrgPrc))
}

func TestClient_InquireCreditBalance(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/ranking/credit-balance`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(http.StatusOK, loadFixtureString(t, "credit_balance_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireCreditBalance(context.Background(), domestic.InquireCreditBalanceParams{
		Symbol:       "005930",
		Option:       "5",
		RankSortCode: "0",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	// UPPERCASE FID_ 키 검증
	assert.Equal(t, "J", capturedQuery.Get("FID_COND_MRKT_DIV_CODE"))
	assert.Equal(t, "11701", capturedQuery.Get("FID_COND_SCR_DIV_CODE"))
	assert.Equal(t, "005930", capturedQuery.Get("FID_INPUT_ISCD"))
	assert.Equal(t, "5", capturedQuery.Get("FID_OPTION"))
	assert.Equal(t, "0", capturedQuery.Get("FID_RANK_SORT_CLS_CODE"))

	// output1 (header array) 검증
	require.Len(t, res.Output1, 2)
	assert.Equal(t, "0001", res.Output1[0].BstpClsCode)
	assert.Equal(t, "코스피", res.Output1[0].HtsKorIsnm)
	assert.Equal(t, "20260501", res.Output1[0].StndDate1)
	assert.Equal(t, "20260507", res.Output1[0].StndDate2)

	// output2 (balance array) 검증
	require.Len(t, res.Output2, 2)
	assert.Equal(t, "005930", res.Output2[0].MkscShrnIscd)
	assert.Equal(t, "삼성전자", res.Output2[0].HtsKorIsnm)
	dPrpr, _ := decimal.NewFromString("75800")
	assert.True(t, dPrpr.Equal(res.Output2[0].StckPrpr))
	dVrss, _ := decimal.NewFromString("-200")
	assert.True(t, dVrss.Equal(res.Output2[0].PrdyVrss))
	assert.Equal(t, "5", res.Output2[0].PrdyVrssSign)
	assert.InDelta(t, -0.26, res.Output2[0].PrdyCtrt, 0.001)
	assert.Equal(t, int64(12345678), res.Output2[0].AcmlVol)
	assert.Equal(t, int64(9876543), res.Output2[0].WholLoanRmndStcn)
	assert.Equal(t, int64(748681660), res.Output2[0].WholLoanRmndAmt)
	assert.InDelta(t, 0.08, res.Output2[0].WholLoanRmndRate, 0.001)
	assert.Equal(t, int64(1234567), res.Output2[0].WholStlnRmndStcn)
	assert.Equal(t, int64(93623580), res.Output2[0].WholStlnRmndAmt)
	assert.InDelta(t, 0.01, res.Output2[0].WholStlnRmndRate, 0.001)
	assert.InDelta(t, 0.12, res.Output2[0].NdayVrssLoanRmndInrt, 0.001)
	assert.InDelta(t, -0.05, res.Output2[0].NdayVrssStlnRmndInrt, 0.001)
}

func TestClient_InquireDailyCreditBalance(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/quotations/daily-credit-balance`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(http.StatusOK, loadFixtureString(t, "daily_credit_balance_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireDailyCreditBalance(context.Background(), domestic.InquireDailyCreditBalanceParams{
		Symbol:     "005930",
		InputDate1: "20260507",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	// lowercase fid_* 키 검증 (UPPERCASE 아님)
	assert.Equal(t, "J", capturedQuery.Get("fid_cond_mrkt_div_code"))
	assert.Equal(t, "20476", capturedQuery.Get("fid_cond_scr_div_code"))
	assert.Equal(t, "005930", capturedQuery.Get("fid_input_iscd"))
	assert.Equal(t, "20260507", capturedQuery.Get("fid_input_date_1"))
	assert.Empty(t, capturedQuery.Get("FID_INPUT_ISCD"), "uppercase 키는 비어야 함")

	require.Len(t, res.Output, 2)

	// output[0] 핵심 필드 검증
	assert.Equal(t, "20260507", res.Output[0].DealDate)
	dPrpr, _ := decimal.NewFromString("75800")
	assert.True(t, dPrpr.Equal(res.Output[0].StckPrpr))
	assert.Equal(t, "5", res.Output[0].PrdyVrssSign)
	dVrss, _ := decimal.NewFromString("-200")
	assert.True(t, dVrss.Equal(res.Output[0].PrdyVrss))
	assert.InDelta(t, -0.26, res.Output[0].PrdyCtrt, 0.001)
	assert.Equal(t, int64(12345678), res.Output[0].AcmlVol)
	assert.Equal(t, "20260509", res.Output[0].StlmDate)
	assert.Equal(t, int64(123456), res.Output[0].WholLoanNewStcn)
	assert.Equal(t, int64(98765), res.Output[0].WholLoanRdmpStcn)
	assert.Equal(t, int64(9876543), res.Output[0].WholLoanRmndStcn)
	assert.Equal(t, int64(9357144404), res.Output[0].WholLoanNewAmt)
	assert.Equal(t, int64(7487234056), res.Output[0].WholLoanRdmpAmt)
	assert.Equal(t, int64(748681660), res.Output[0].WholLoanRmndAmt)
	assert.InDelta(t, 0.08, res.Output[0].WholLoanRmndRate, 0.001)
	assert.InDelta(t, 98.50, res.Output[0].WholLoanGvrt, 0.001)
	assert.Equal(t, int64(12345), res.Output[0].WholStlnNewStcn)
	assert.Equal(t, int64(9876), res.Output[0].WholStlnRdmpStcn)
	assert.Equal(t, int64(1234567), res.Output[0].WholStlnRmndStcn)
	assert.Equal(t, int64(935714440), res.Output[0].WholStlnNewAmt)
	assert.Equal(t, int64(748481660), res.Output[0].WholStlnRdmpAmt)
	assert.Equal(t, int64(93623580), res.Output[0].WholStlnRmndAmt)
	assert.InDelta(t, 0.01, res.Output[0].WholStlnRmndRate, 0.001)
	assert.InDelta(t, 99.20, res.Output[0].WholStlnGvrt, 0.001)
	dOprc, _ := decimal.NewFromString("76000")
	assert.True(t, dOprc.Equal(res.Output[0].StckOprc))
	dHgpr, _ := decimal.NewFromString("76200")
	assert.True(t, dHgpr.Equal(res.Output[0].StckHgpr))
	dLwpr, _ := decimal.NewFromString("75500")
	assert.True(t, dLwpr.Equal(res.Output[0].StckLwpr))
}

func TestClient_InquireLendableByCompany(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/quotations/lendable-by-company`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(http.StatusOK, loadFixtureString(t, "lendable_by_company_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireLendableByCompany(context.Background(), domestic.InquireLendableByCompanyParams{
		ExcgDvsnCd:     "02",
		Pdno:           "005930",
		ThcoStlnPsblYn: "Y",
		InqrDvsn1:      "1",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	// non-FID UPPERCASE 키 검증
	assert.Equal(t, "02", capturedQuery.Get("EXCG_DVSN_CD"))
	assert.Equal(t, "005930", capturedQuery.Get("PDNO"))
	assert.Equal(t, "Y", capturedQuery.Get("THCO_STLN_PSBL_YN"))
	assert.Equal(t, "1", capturedQuery.Get("INQR_DVSN_1"))
	assert.Equal(t, "", capturedQuery.Get("CTX_AREA_FK200"))
	assert.Equal(t, "", capturedQuery.Get("CTX_AREA_NK100"))

	// output1 (array) 검증
	require.Len(t, res.Output1, 2)
	assert.Equal(t, "005930", res.Output1[0].Pdno)
	assert.Equal(t, "삼성전자", res.Output1[0].PrdtName)
	dPapr, _ := decimal.NewFromString("75800")
	assert.True(t, dPapr.Equal(res.Output1[0].Papr))
	dClpr, _ := decimal.NewFromString("76000")
	assert.True(t, dClpr.Equal(res.Output1[0].BfdyClpr))
	dSbst, _ := decimal.NewFromString("75800")
	assert.True(t, dSbst.Equal(res.Output1[0].SbstPrvs))
	assert.Equal(t, "정상", res.Output1[0].TrStopDvsnName)
	assert.Equal(t, "가능", res.Output1[0].PsblYnName)
	assert.Equal(t, int64(500000), res.Output1[0].LmtQty1)
	assert.Equal(t, int64(123456), res.Output1[0].UseQty1)
	assert.Equal(t, int64(376544), res.Output1[0].TradPsblQty2)
	assert.Equal(t, "01", res.Output1[0].RghtTypeCd)
	assert.Equal(t, "20260507", res.Output1[0].BassDt)
	assert.Equal(t, "Y", res.Output1[0].PsblYn)

	// output2 (summary) 검증
	assert.Equal(t, int64(800000), res.Output2.TotStupLmtQty)
	assert.Equal(t, int64(600000), res.Output2.BrchLmtQty)
	assert.Equal(t, int64(588890), res.Output2.RqstPsblQty)
}

func TestClient_InquireAfterHourBalance(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/ranking/after-hour-balance`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(http.StatusOK, loadFixtureString(t, "after_hour_balance_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireAfterHourBalance(context.Background(), domestic.InquireAfterHourBalanceParams{
		Symbol:       "0000",
		RankSortCode: "0",
		DivClsCode:   "0",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	// lowercase fid_* 키 검증 (UPPERCASE 아님)
	assert.Equal(t, "J", capturedQuery.Get("fid_cond_mrkt_div_code"))
	assert.Equal(t, "20176", capturedQuery.Get("fid_cond_scr_div_code"))
	assert.Equal(t, "0000", capturedQuery.Get("fid_input_iscd"))
	assert.Equal(t, "0", capturedQuery.Get("fid_rank_sort_cls_code"))
	assert.Equal(t, "0", capturedQuery.Get("fid_div_cls_code"))
	assert.Empty(t, capturedQuery.Get("FID_INPUT_ISCD"), "uppercase 키는 비어야 함")

	require.Len(t, res.Output, 2)

	// output[0] 핵심 필드 검증
	assert.Equal(t, "005930", res.Output[0].StckShrnIscd)
	assert.Equal(t, "1", res.Output[0].DataRank)
	assert.Equal(t, "삼성전자", res.Output[0].HtsKorIsnm)
	dPrpr, _ := decimal.NewFromString("75800")
	assert.True(t, dPrpr.Equal(res.Output[0].StckPrpr))
	dVrss, _ := decimal.NewFromString("-200")
	assert.True(t, dVrss.Equal(res.Output[0].PrdyVrss))
	assert.Equal(t, "5", res.Output[0].PrdyVrssSign)
	assert.InDelta(t, -0.26, res.Output[0].PrdyCtrt, 0.001)
	assert.Equal(t, int64(345678), res.Output[0].OvtmTotalAskpRsqn)
	assert.Equal(t, int64(456789), res.Output[0].OvtmTotalBidpRsqn)
	assert.Equal(t, int64(234567), res.Output[0].MkobOtcpVol)
	assert.Equal(t, int64(222222), res.Output[0].MkfaOtcpVol)
}

func TestClient_InquireQuoteBalance(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/ranking/quote-balance`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(http.StatusOK, loadFixtureString(t, "quote_balance_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireQuoteBalance(context.Background(), domestic.InquireQuoteBalanceParams{
		VolCnt:       "30",
		Symbol:       "0000",
		RankSortCode: "0",
		DivClsCode:   "0",
		TrgtClsCode:  "111111111",
		TrgtExlsCode: "000000000",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	// lowercase fid_* 키 검증 (UPPERCASE 아님)
	assert.Equal(t, "J", capturedQuery.Get("fid_cond_mrkt_div_code"))
	assert.Equal(t, "20172", capturedQuery.Get("fid_cond_scr_div_code"))
	assert.Equal(t, "0000", capturedQuery.Get("fid_input_iscd"))
	assert.Equal(t, "30", capturedQuery.Get("fid_vol_cnt"))
	assert.Equal(t, "0", capturedQuery.Get("fid_rank_sort_cls_code"))
	assert.Equal(t, "0", capturedQuery.Get("fid_div_cls_code"))
	assert.Empty(t, capturedQuery.Get("FID_INPUT_ISCD"), "uppercase 키는 비어야 함")

	require.Len(t, res.Output, 2)

	// output[0] 핵심 필드 검증
	assert.Equal(t, "005930", res.Output[0].MkscShrnIscd)
	assert.Equal(t, "1", res.Output[0].DataRank)
	assert.Equal(t, "삼성전자", res.Output[0].HtsKorIsnm)
	dPrpr, _ := decimal.NewFromString("75800")
	assert.True(t, dPrpr.Equal(res.Output[0].StckPrpr))
	dVrss, _ := decimal.NewFromString("-200")
	assert.True(t, dVrss.Equal(res.Output[0].PrdyVrss))
	assert.Equal(t, "5", res.Output[0].PrdyVrssSign)
	assert.InDelta(t, -0.26, res.Output[0].PrdyCtrt, 0.001)
	assert.Equal(t, int64(12345678), res.Output[0].AcmlVol)
	assert.Equal(t, int64(2345678), res.Output[0].TotalAskpRsqn)
	assert.Equal(t, int64(3456789), res.Output[0].TotalBidpRsqn)
	assert.Equal(t, int64(1111111), res.Output[0].TotalNtslBidpRsqn)
	assert.InDelta(t, 59.60, res.Output[0].ShnuRsqnRate, 0.001)
	assert.InDelta(t, 40.40, res.Output[0].SelnRsqnRate, 0.001)
}

func TestClient_InquireOvertimeExpTransFluct(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/ranking/overtime-exp-trans-fluct`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(http.StatusOK, loadFixtureString(t, "overtime_exp_trans_fluct_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireOvertimeExpTransFluct(context.Background(), domestic.InquireOvertimeExpTransFluctParams{
		Symbol:       "0000",
		RankSortCode: "0",
		DivClsCode:   "0",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	// UPPERCASE FID_ 키 검증
	assert.Equal(t, "J", capturedQuery.Get("FID_COND_MRKT_DIV_CODE"))
	assert.Equal(t, "11186", capturedQuery.Get("FID_COND_SCR_DIV_CODE"))
	assert.Equal(t, "0000", capturedQuery.Get("FID_INPUT_ISCD"))
	assert.Equal(t, "0", capturedQuery.Get("FID_RANK_SORT_CLS_CODE"))
	assert.Equal(t, "0", capturedQuery.Get("FID_DIV_CLS_CODE"))
	assert.Empty(t, capturedQuery.Get("fid_input_iscd"), "lowercase 키는 비어야 함")

	// output 은 단일 객체 (배열 아님) 검증
	assert.Equal(t, "1", res.Output.DataRank)
	assert.Equal(t, "00", res.Output.IscdStatClsCode)
	assert.Equal(t, "005930", res.Output.StckShrnIscd)
	assert.Equal(t, "삼성전자", res.Output.HtsKorIsnm)
	dCnpr, _ := decimal.NewFromString("75900")
	assert.True(t, dCnpr.Equal(res.Output.OvtmUntpAntcCnpr))
	dVrss, _ := decimal.NewFromString("100")
	assert.True(t, dVrss.Equal(res.Output.OvtmUntpAntcCntgVrss))
	// KIS docs 오타 보존: vrss+sign 연결 (밑줄 없음)
	assert.Equal(t, "2", res.Output.OvtmUntpAntcCntgVrsssign)
	assert.InDelta(t, 0.13, res.Output.OvtmUntpAntcCntgCtrt, 0.001)
	assert.Equal(t, int64(234567), res.Output.OvtmUntpAskpRsqn1)
	assert.Equal(t, int64(345678), res.Output.OvtmUntpBidpRsqn1)
	assert.Equal(t, int64(111111), res.Output.OvtmUntpAntcCnqn)
	assert.Equal(t, int64(98765), res.Output.ItmtVol)
	dPrpr, _ := decimal.NewFromString("75800")
	assert.True(t, dPrpr.Equal(res.Output.StckPrpr))
}

func TestClient_InquireMarketValue(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/ranking/market-value`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(http.StatusOK, loadFixtureString(t, "market_value_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireMarketValue(context.Background(), domestic.InquireMarketValueParams{
		Symbol:       "0000",
		RankSortCode: "0",
		DivClsCode:   "0",
		TrgtClsCode:  "111111111",
		TrgtExlsCode: "000000000",
		BlngClsCode:  "0",
		InputOption1: "2025",
		InputOption2: "1",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	// lowercase fid_* 키 검증 (UPPERCASE 아님)
	assert.Equal(t, "J", capturedQuery.Get("fid_cond_mrkt_div_code"))
	assert.Equal(t, "20179", capturedQuery.Get("fid_cond_scr_div_code"))
	assert.Equal(t, "0000", capturedQuery.Get("fid_input_iscd"))
	assert.Equal(t, "0", capturedQuery.Get("fid_rank_sort_cls_code"))
	assert.Equal(t, "0", capturedQuery.Get("fid_div_cls_code"))
	assert.Equal(t, "2025", capturedQuery.Get("fid_input_option_1"))
	assert.Equal(t, "1", capturedQuery.Get("fid_input_option_2"))
	assert.Empty(t, capturedQuery.Get("FID_INPUT_ISCD"), "uppercase 키는 비어야 함")

	require.Len(t, res.Output, 2)

	// output[0] 핵심 필드 검증 (20 fields)
	assert.Equal(t, "1", res.Output[0].DataRank)
	assert.Equal(t, "삼성전자", res.Output[0].HtsKorIsnm)
	assert.Equal(t, "005930", res.Output[0].MkscShrnIscd)
	dPrpr, _ := decimal.NewFromString("75800")
	assert.True(t, dPrpr.Equal(res.Output[0].StckPrpr))
	dVrss, _ := decimal.NewFromString("-200")
	assert.True(t, dVrss.Equal(res.Output[0].PrdyVrss))
	assert.Equal(t, "5", res.Output[0].PrdyVrssSign)
	assert.InDelta(t, -0.26, res.Output[0].PrdyCtrt, 0.001)
	assert.Equal(t, int64(12345678), res.Output[0].AcmlVol)
	assert.InDelta(t, 12.34, res.Output[0].Per, 0.001)
	assert.InDelta(t, 1.23, res.Output[0].Pbr, 0.001)
	assert.InDelta(t, 8.90, res.Output[0].Pcr, 0.001)
	assert.InDelta(t, 1.45, res.Output[0].Psr, 0.001)
	assert.InDelta(t, 6143.0, res.Output[0].Eps, 0.001)
	assert.InDelta(t, 1234567.0, res.Output[0].Eva, 0.001)
	assert.InDelta(t, 98765432.0, res.Output[0].Ebitda, 0.001)
	assert.InDelta(t, 7.65, res.Output[0].PvDivEbitda, 0.001)
	assert.InDelta(t, 23.45, res.Output[0].EbitdaDivFnncExpn, 0.001)
	assert.Equal(t, "12", res.Output[0].StacMonth)
	assert.Equal(t, "01", res.Output[0].StacMonthClsCode)
	assert.Equal(t, "1", res.Output[0].IqryCsnu)
}

func TestClient_InquireDisparity(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/ranking/disparity`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(http.StatusOK, loadFixtureString(t, "disparity_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireDisparity(context.Background(), domestic.InquireDisparityParams{
		Symbol:       "0000",
		HourClsCode:  "20",
		DivClsCode:   "0",
		RankSortCode: "0",
		TrgtClsCode:  "111111111",
		TrgtExlsCode: "000000000",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	// lowercase fid_* 키 검증
	assert.Equal(t, "J", capturedQuery.Get("fid_cond_mrkt_div_code"))
	assert.Equal(t, "20178", capturedQuery.Get("fid_cond_scr_div_code"))
	assert.Equal(t, "0000", capturedQuery.Get("fid_input_iscd"))
	assert.Equal(t, "20", capturedQuery.Get("fid_hour_cls_code"))
	assert.Equal(t, "0", capturedQuery.Get("fid_div_cls_code"))
	assert.Empty(t, capturedQuery.Get("FID_COND_MRKT_DIV_CODE"), "uppercase 키는 비어야 함")

	require.Len(t, res.Output, 2)

	// output[0] 핵심 필드 검증 (13 fields)
	assert.Equal(t, "005930", res.Output[0].MkscShrnIscd)
	assert.Equal(t, "1", res.Output[0].DataRank)
	assert.Equal(t, "삼성전자", res.Output[0].HtsKorIsnm)
	dPrpr, _ := decimal.NewFromString("75800")
	assert.True(t, dPrpr.Equal(res.Output[0].StckPrpr))
	dVrss, _ := decimal.NewFromString("-200")
	assert.True(t, dVrss.Equal(res.Output[0].PrdyVrss))
	assert.InDelta(t, -0.26, res.Output[0].PrdyCtrt, 0.001)
	assert.Equal(t, "5", res.Output[0].PrdyVrssSign)
	assert.Equal(t, int64(12345678), res.Output[0].AcmlVol)
	assert.InDelta(t, 101.23, res.Output[0].D5Dsrt, 0.001)
	assert.InDelta(t, 99.87, res.Output[0].D10Dsrt, 0.001)
	assert.InDelta(t, 98.54, res.Output[0].D20Dsrt, 0.001)
	assert.InDelta(t, 97.32, res.Output[0].D60Dsrt, 0.001)
	assert.InDelta(t, 96.10, res.Output[0].D120Dsrt, 0.001)
}

func TestClient_InquirePreferDisparateRatio(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/ranking/prefer-disparate-ratio`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(http.StatusOK, loadFixtureString(t, "prefer_disparate_ratio_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquirePreferDisparateRatio(context.Background(), domestic.InquirePreferDisparateRatioParams{
		Symbol:       "0000",
		DivClsCode:   "0",
		TrgtClsCode:  "111111111",
		TrgtExlsCode: "000000000",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	// lowercase fid_* 키 검증
	assert.Equal(t, "J", capturedQuery.Get("fid_cond_mrkt_div_code"))
	assert.Equal(t, "20177", capturedQuery.Get("fid_cond_scr_div_code"))
	assert.Equal(t, "0000", capturedQuery.Get("fid_input_iscd"))
	assert.Equal(t, "0", capturedQuery.Get("fid_div_cls_code"))
	assert.Empty(t, capturedQuery.Get("FID_COND_MRKT_DIV_CODE"), "uppercase 키는 비어야 함")

	require.Len(t, res.Output, 2)

	// output[0] 핵심 필드 검증 (17 fields)
	assert.Equal(t, "005930", res.Output[0].MkscShrnIscd)
	assert.Equal(t, "1", res.Output[0].DataRank)
	assert.Equal(t, "삼성전자", res.Output[0].HtsKorIsnm)
	dPrpr, _ := decimal.NewFromString("75800")
	assert.True(t, dPrpr.Equal(res.Output[0].StckPrpr))
	dVrss, _ := decimal.NewFromString("-200")
	assert.True(t, dVrss.Equal(res.Output[0].PrdyVrss))
	assert.Equal(t, "5", res.Output[0].PrdyVrssSign)
	assert.Equal(t, int64(12345678), res.Output[0].AcmlVol)
	assert.Equal(t, "005935", res.Output[0].PrstIscd)
	assert.Equal(t, "삼성전자우", res.Output[0].PrstKorIsnm)
	dPrstPrpr, _ := decimal.NewFromString("68500")
	assert.True(t, dPrstPrpr.Equal(res.Output[0].PrstPrpr))
	dPrstVrss, _ := decimal.NewFromString("-300")
	assert.True(t, dPrstVrss.Equal(res.Output[0].PrstPrdyVrss))
	assert.Equal(t, "5", res.Output[0].PrstPrdyVrssSign)
	assert.Equal(t, int64(234567), res.Output[0].PrstAcmlVol)
	dDiffPrpr, _ := decimal.NewFromString("7300")
	assert.True(t, dDiffPrpr.Equal(res.Output[0].DiffPrpr))
	assert.InDelta(t, 10.66, res.Output[0].Dprt, 0.001)
	assert.InDelta(t, -0.26, res.Output[0].PrdyCtrt, 0.001)
	assert.InDelta(t, -0.44, res.Output[0].PrstPrdyCtrt, 0.001)
}

func TestClient_InquireProfitAssetIndex(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/ranking/profit-asset-index`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(http.StatusOK, loadFixtureString(t, "profit_asset_index_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireProfitAssetIndex(context.Background(), domestic.InquireProfitAssetIndexParams{
		Symbol:       "0000",
		DivClsCode:   "0",
		TrgtClsCode:  "111111111",
		TrgtExlsCode: "000000000",
		RankSortCode: "0",
		BlngClsCode:  "0",
		InputOption1: "2025",
		InputOption2: "1",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	// lowercase fid_* 키 검증
	assert.Equal(t, "J", capturedQuery.Get("fid_cond_mrkt_div_code"))
	assert.Equal(t, "20173", capturedQuery.Get("fid_cond_scr_div_code"))
	assert.Equal(t, "0000", capturedQuery.Get("fid_input_iscd"))
	assert.Equal(t, "0", capturedQuery.Get("fid_div_cls_code"))
	assert.Equal(t, "2025", capturedQuery.Get("fid_input_option_1"))
	assert.Equal(t, "1", capturedQuery.Get("fid_input_option_2"))
	assert.Empty(t, capturedQuery.Get("FID_COND_MRKT_DIV_CODE"), "uppercase 키는 비어야 함")

	require.Len(t, res.Output, 2)

	// output[0] 핵심 필드 검증 (18 fields)
	assert.Equal(t, "1", res.Output[0].DataRank)
	assert.Equal(t, "삼성전자", res.Output[0].HtsKorIsnm)
	assert.Equal(t, "5", res.Output[0].PrdyVrssSign)
	assert.Equal(t, "005930", res.Output[0].MkscShrnIscd)
	dPrpr, _ := decimal.NewFromString("75800")
	assert.True(t, dPrpr.Equal(res.Output[0].StckPrpr))
	dVrss, _ := decimal.NewFromString("-200")
	assert.True(t, dVrss.Equal(res.Output[0].PrdyVrss))
	assert.InDelta(t, -0.26, res.Output[0].PrdyCtrt, 0.001)
	assert.Equal(t, int64(12345678), res.Output[0].AcmlVol)
	assert.Equal(t, int64(12345678901), res.Output[0].SaleTotlPrfi)
	assert.Equal(t, int64(9876543210), res.Output[0].BsopPrti)
	assert.Equal(t, int64(9876543210), res.Output[0].OpPrfi)
	assert.Equal(t, int64(8765432109), res.Output[0].ThtrNtin)
	assert.Equal(t, int64(987654321098), res.Output[0].TotalAset)
	assert.Equal(t, int64(345678901234), res.Output[0].TotalLblt)
	assert.Equal(t, int64(641975419864), res.Output[0].TotalCptl)
	assert.Equal(t, "12", res.Output[0].StacMonth)
	assert.Equal(t, "01", res.Output[0].StacMonthClsCode)
	assert.Equal(t, "1", res.Output[0].IqryCsnu)
}
