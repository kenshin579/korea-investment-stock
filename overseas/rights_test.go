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

func TestClient_InquireRightsByIce(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/quotations/rights-by-ice`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "rights_by_ice_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireRightsByIce(context.Background(), overseas.InquireRightsByIceParams{
		NCod:  "US",
		Symb:  "AAPL",
		StYmd: "20260401",
		EdYmd: "20260430",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, "US", capturedQuery.Get("NCOD"))
	assert.Equal(t, "AAPL", capturedQuery.Get("SYMB"))
	assert.Equal(t, "20260401", capturedQuery.Get("ST_YMD"))
	assert.Equal(t, "20260430", capturedQuery.Get("ED_YMD"))

	// output1 only — output2 없음 anomaly 검증
	require.Len(t, res.Output1, 2)
	assert.Equal(t, "20260401", res.Output1[0].AnnoDt)
	assert.Equal(t, "주식배당", res.Output1[0].CaTitle)
	assert.Equal(t, "20260430", res.Output1[0].PayDt)
}

func TestClient_InquirePeriodRights(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/quotations/period-rights`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "period_rights_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquirePeriodRights(context.Background(), overseas.InquirePeriodRightsParams{
		InqrStrtDt: "20260401",
		InqrEndDt:  "20260430",
		Pdno:       "AAPL",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, "20260401", capturedQuery.Get("INQR_STRT_DT"))
	assert.Equal(t, "20260430", capturedQuery.Get("INQR_END_DT"))
	assert.Equal(t, "AAPL", capturedQuery.Get("PDNO"))
	// cursor pagination 첫 조회 시 빈 값 전송 확인
	assert.Equal(t, "", capturedQuery.Get("CTX_AREA_NK50"))
	assert.Equal(t, "", capturedQuery.Get("CTX_AREA_FK50"))

	require.Len(t, res.Output, 2)
	assert.Equal(t, "20260401", res.Output[0].BassDt)
	assert.Equal(t, "AAPL", res.Output[0].Pdno)
	assert.Equal(t, "Y", res.Output[0].DfntYn)
	assert.Equal(t, "USD", res.Output[0].CrcyCd)
}

// Coverage 보강 — JSON unmarshal error path 검증.

func TestClient_InquireRightsByIce_InvalidJSON(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(http.MethodGet, `=~/quotations/rights-by-ice`,
		httpmock.NewStringResponder(200, `{"rt_cd":"0","msg_cd":"X","msg1":"x","output1": "not-array"}`))

	c := newTestClient(t)
	_, err := c.InquireRightsByIce(context.Background(), overseas.InquireRightsByIceParams{})
	require.Error(t, err)
}

func TestClient_InquirePeriodRights_InvalidJSON(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(http.MethodGet, `=~/quotations/period-rights`,
		httpmock.NewStringResponder(200, `{"rt_cd":"0","msg_cd":"X","msg1":"x","output": "not-array"}`))

	c := newTestClient(t)
	_, err := c.InquirePeriodRights(context.Background(), overseas.InquirePeriodRightsParams{})
	require.Error(t, err)
}
