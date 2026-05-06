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
