package overseas_test

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kenshin579/korea-investment-stock/overseas"
)

func TestClient_InquireUpdownRate(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/ranking/updown-rate`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "updown_rate_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireUpdownRate(context.Background(), overseas.InquireUpdownRateParams{
		Excd: "NAS",
		Gubn: "1", // 상승율
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, "NAS", capturedQuery.Get("EXCD"))
	assert.Equal(t, "1", capturedQuery.Get("GUBN"))
	assert.Equal(t, "0", capturedQuery.Get("NDAY"))
	assert.Equal(t, "0", capturedQuery.Get("VOL_RANG"))

	require.Len(t, res.Output2, 2)
	assert.Equal(t, "NVDA", res.Output2[0].Symb)
	assert.Equal(t, "엔비디아", res.Output2[0].Name)
	d, _ := decimal.NewFromString("920.45")
	assert.True(t, d.Equal(res.Output2[0].Last))
	assert.InDelta(t, 5.16, res.Output2[0].Rate, 0.001)
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
	res, err := c.InquireMarketCap(context.Background(), overseas.InquireMarketCapParams{
		ExcdCode: "NAS",
		VolRang:  "0",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, "NAS", capturedQuery.Get("EXCD"))
	assert.Equal(t, "0", capturedQuery.Get("VOL_RANG"))

	// output1 검증
	assert.Equal(t, "2", res.Output1.Zdiv)
	assert.Equal(t, int64(2), res.Output1.Crec)
	assert.Equal(t, int64(500), res.Output1.Trec)
	assert.Equal(t, int64(30), res.Output1.Nrec)

	// output2[0] 검증
	require.Len(t, res.Output2, 2)
	assert.Equal(t, "AAPL", res.Output2[0].Symb)
	assert.Equal(t, "애플", res.Output2[0].Name)
	d, _ := decimal.NewFromString("189.30")
	assert.True(t, d.Equal(res.Output2[0].Last))
	assert.InDelta(t, 1.34, res.Output2[0].Rate, 0.001)
	assert.Equal(t, int64(55000000), res.Output2[0].Tvol)
	assert.Equal(t, int64(15634232000), res.Output2[0].Shar)
	tomv, _ := decimal.NewFromString("2958652560000")
	assert.True(t, tomv.Equal(res.Output2[0].Tomv))
	assert.InDelta(t, 6.85, res.Output2[0].Grav, 0.001)
	assert.Equal(t, int64(1), res.Output2[0].Rank)
	assert.Equal(t, "APPLE INC", res.Output2[0].Ename)
}

func TestClient_InquireTradeVol(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/ranking/trade-vol`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "trade_vol_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireTradeVol(context.Background(), overseas.InquireTradeVolParams{
		ExcdCode: "NAS",
		NDay:     "0",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, "NAS", capturedQuery.Get("EXCD"))
	assert.Equal(t, "0", capturedQuery.Get("NDAY"))
	assert.Equal(t, "0", capturedQuery.Get("VOL_RANG"))

	// output1 검증
	assert.Equal(t, int64(2), res.Output1.Crec)
	assert.Equal(t, int64(30), res.Output1.Nrec)

	// output2[0] 검증
	require.Len(t, res.Output2, 2)
	assert.Equal(t, "AAPL", res.Output2[0].Symb)
	d, _ := decimal.NewFromString("189.30")
	assert.True(t, d.Equal(res.Output2[0].Last))
	pask, _ := decimal.NewFromString("189.35")
	assert.True(t, pask.Equal(res.Output2[0].Pask))
	pbid, _ := decimal.NewFromString("189.25")
	assert.True(t, pbid.Equal(res.Output2[0].Pbid))
	assert.Equal(t, int64(55000000), res.Output2[0].Tvol)
	assert.Equal(t, int64(10411500000), res.Output2[0].Tamt)
	assert.Equal(t, int64(48000000), res.Output2[0].ATvol)
	assert.Equal(t, int64(1), res.Output2[0].Rank)
}

func TestClient_InquireTradePbmn(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/ranking/trade-pbmn`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "trade_pbmn_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireTradePbmn(context.Background(), overseas.InquireTradePbmnParams{
		ExcdCode: "NAS",
		NDay:     "0",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, "NAS", capturedQuery.Get("EXCD"))
	assert.Equal(t, "0", capturedQuery.Get("NDAY"))
	assert.Equal(t, "0", capturedQuery.Get("VOL_RANG"))

	// output1 검증
	assert.Equal(t, int64(30), res.Output1.Nrec)

	// output2[0] 검증 — MSFT 가 순위 1 (거래대금 기준)
	require.Len(t, res.Output2, 2)
	assert.Equal(t, "MSFT", res.Output2[0].Symb)
	d, _ := decimal.NewFromString("415.20")
	assert.True(t, d.Equal(res.Output2[0].Last))
	assert.Equal(t, int64(9134400000), res.Output2[0].Tamt)
	assert.Equal(t, int64(8500000000), res.Output2[0].ATamt) // a_tamt, not a_tvol
	assert.Equal(t, int64(1), res.Output2[0].Rank)
}

func TestClient_InquireVolumeSurge(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/ranking/volume-surge`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "volume_surge_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireVolumeSurge(context.Background(), overseas.InquireVolumeSurgeParams{
		ExcdCode: "NAS",
		MixN:     "0",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, "NAS", capturedQuery.Get("EXCD"))
	assert.Equal(t, "0", capturedQuery.Get("MIXN"))
	assert.Equal(t, "0", capturedQuery.Get("VOL_RANG"))

	// output1 검증 — 3-field MinSummary (crec/trec 없음)
	assert.Equal(t, "2", res.Output1.Zdiv)
	assert.Equal(t, int64(30), res.Output1.Nrec)

	// output2[0] 검증 — knam/enam (name/ename 아님)
	require.Len(t, res.Output2, 2)
	assert.Equal(t, "AAPL", res.Output2[0].Symb)
	assert.Equal(t, "애플", res.Output2[0].Knam)
	assert.Equal(t, "APPLE INC", res.Output2[0].Enam)
	d, _ := decimal.NewFromString("189.30")
	assert.True(t, d.Equal(res.Output2[0].Last))
	assert.Equal(t, int64(55000000), res.Output2[0].Tvol)
	assert.Equal(t, int64(8000000), res.Output2[0].NTvol)
	ndiff, _ := decimal.NewFromString("47000000")
	assert.True(t, ndiff.Equal(res.Output2[0].NDiff))
	assert.InDelta(t, 587.50, res.Output2[0].NRate, 0.01)
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
	res, err := c.InquireVolumePower(context.Background(), overseas.InquireVolumePowerParams{
		ExcdCode: "NAS",
		NDay:     "0",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, "NAS", capturedQuery.Get("EXCD"))
	assert.Equal(t, "0", capturedQuery.Get("NDAY")) // wire name NDAY (분 단위)
	assert.Equal(t, "0", capturedQuery.Get("VOL_RANG"))

	// output1 검증 — 3-field MinSummary
	assert.Equal(t, int64(30), res.Output1.Nrec)

	// output2[0] 검증 — knam/enam + tpow/powx
	require.Len(t, res.Output2, 2)
	assert.Equal(t, "AAPL", res.Output2[0].Symb)
	assert.Equal(t, "애플", res.Output2[0].Knam)
	assert.Equal(t, "APPLE INC", res.Output2[0].Enam)
	d, _ := decimal.NewFromString("189.30")
	assert.True(t, d.Equal(res.Output2[0].Last))
	assert.InDelta(t, 143.25, res.Output2[0].Tpow, 0.01)
	assert.InDelta(t, 138.90, res.Output2[0].Powx, 0.01)
}
