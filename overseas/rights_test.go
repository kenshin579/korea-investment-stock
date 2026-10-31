package overseas_test

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
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

// 실 응답(2026-09 시장 전체 배당, 2026-11-01 실측)으로 만든 픽스처 — 첫 페이지는 tr_cont "M" + 공백 포함 50자 커서.
func TestClient_InquirePeriodRights_ParsesCursorAndTrCont(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var gotTrCont string
	httpmock.RegisterResponder(http.MethodGet, `=~/quotations/period-rights`,
		func(req *http.Request) (*http.Response, error) {
			gotTrCont = req.Header.Get("tr_cont")
			resp := httpmock.NewStringResponse(200, loadFixtureString(t, "period_rights_page1.json"))
			resp.Header.Set("tr_cont", "M")
			return resp, nil
		})

	c := newTestClient(t)
	res, err := c.InquirePeriodRights(context.Background(), overseas.InquirePeriodRightsParams{
		RghtTypeCd: "03", InqrDvsnCd: "02", InqrStrtDt: "20260901", InqrEndDt: "20260930",
	})
	require.NoError(t, err)
	assert.Equal(t, "", gotTrCont, "첫 페이지는 tr_cont 헤더 없음")
	assert.Equal(t, "M", res.TrCont)
	assert.Equal(t, "20260902!^GOVZ!^529!^03                           ", res.CtxAreaNk50, "커서는 뒤 공백까지 그대로")
	assert.Equal(t, "03!^02!^20260901!^20260930!^!^                    ", res.CtxAreaFk50)
	require.Len(t, res.Output, 3)
	assert.Equal(t, "03", res.Output[0].RghtTypeCd)
	assert.Equal(t, "20260901", res.Output[0].AcplBassDt)
}

func TestClient_InquirePeriodRights_SendsTrContHeader(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var gotTrCont string
	var q url.Values
	httpmock.RegisterResponder(http.MethodGet, `=~/quotations/period-rights`,
		func(req *http.Request) (*http.Response, error) {
			gotTrCont, q = req.Header.Get("tr_cont"), req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "period_rights_last.json")), nil
		})

	c := newTestClient(t)
	_, err := c.InquirePeriodRights(context.Background(), overseas.InquirePeriodRightsParams{
		TrCont: "N", CtxAreaNk50: "NK ", CtxAreaFk50: "FK ",
	})
	require.NoError(t, err)
	assert.Equal(t, "N", gotTrCont)
	assert.Equal(t, "NK ", q.Get("CTX_AREA_NK50"))
	assert.Equal(t, "FK ", q.Get("CTX_AREA_FK50"))
}

func TestClient_InquirePeriodRightsAll_FollowsTrCont(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	calls := 0
	httpmock.RegisterResponder(http.MethodGet, `=~/quotations/period-rights`,
		func(req *http.Request) (*http.Response, error) {
			calls++
			q := req.URL.Query()
			switch calls {
			case 1:
				assert.Equal(t, "", req.Header.Get("tr_cont"))
				assert.Equal(t, "", q.Get("CTX_AREA_NK50"))
				assert.Equal(t, "", q.Get("CTX_AREA_FK50"))
				resp := httpmock.NewStringResponse(200, loadFixtureString(t, "period_rights_page1.json"))
				resp.Header.Set("tr_cont", "M")
				return resp, nil
			case 2:
				assert.Equal(t, "N", req.Header.Get("tr_cont"))
				assert.Equal(t, "20260902!^GOVZ!^529!^03                           ", q.Get("CTX_AREA_NK50"))
				assert.Equal(t, "03!^02!^20260901!^20260930!^!^                    ", q.Get("CTX_AREA_FK50"))
				assert.Equal(t, "03", q.Get("RGHT_TYPE_CD"), "조회 조건은 페이지마다 그대로")
				resp := httpmock.NewStringResponse(200, loadFixtureString(t, "period_rights_last.json"))
				resp.Header.Set("tr_cont", "D")
				return resp, nil
			}
			t.Fatalf("unexpected call #%d", calls)
			return nil, nil
		})

	c := newTestClient(t)
	res, err := c.InquirePeriodRightsAll(context.Background(), overseas.InquirePeriodRightsParams{
		RghtTypeCd: "03", InqrDvsnCd: "02", InqrStrtDt: "20260901", InqrEndDt: "20260930",
		TrCont: "N", CtxAreaNk50: "junk", CtxAreaFk50: "junk",
	})
	require.NoError(t, err)
	assert.Equal(t, 2, calls)
	require.Len(t, res.Output, 5)
	assert.Equal(t, "00003", res.Output[0].Pdno)
	assert.Equal(t, "TDAX", res.Output[3].Pdno)
	assert.Equal(t, "TFII", res.Output[4].Pdno)
	assert.Equal(t, "D", res.TrCont)
}

func TestClient_InquirePeriodRightsAll_NoData(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder(http.MethodGet, `=~/quotations/period-rights`,
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(200, loadFixtureString(t, "period_rights_nodata.json"))
			resp.Header.Set("tr_cont", "D")
			return resp, nil
		})
	c := newTestClient(t)
	res, err := c.InquirePeriodRightsAll(context.Background(), overseas.InquirePeriodRightsParams{
		RghtTypeCd: "03", InqrDvsnCd: "02", InqrStrtDt: "20261101", InqrEndDt: "20261231",
	})
	require.NoError(t, err, "KIER2520(자료 없음)은 rt_cd 0 — 에러가 아니다")
	assert.Empty(t, res.Output)
	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestClient_InquirePeriodRightsAll_CursorNotAdvancing(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder(http.MethodGet, `=~/quotations/period-rights`,
		func(req *http.Request) (*http.Response, error) {
			// 커서를 넘겨도 같은 첫 페이지·같은 커서가 돌아오는 경우 — tr_cont 헤더를 빠뜨렸을 때 실 API 가 이렇다
			resp := httpmock.NewStringResponse(200, loadFixtureString(t, "period_rights_page1.json"))
			resp.Header.Set("tr_cont", "M")
			return resp, nil
		})
	c := newTestClient(t)
	_, err := c.InquirePeriodRightsAll(context.Background(), overseas.InquirePeriodRightsParams{RghtTypeCd: "03"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "cursor did not advance")
	assert.Equal(t, 2, httpmock.GetTotalCallCount())
}

func TestClient_InquirePeriodRightsAll_PageCap(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	calls := 0
	httpmock.RegisterResponder(http.MethodGet, `=~/quotations/period-rights`,
		func(req *http.Request) (*http.Response, error) {
			calls++
			body := strings.Replace(loadFixtureString(t, "period_rights_page1.json"), "GOVZ", fmt.Sprintf("G%d", calls), 1)
			resp := httpmock.NewStringResponse(200, body)
			resp.Header.Set("tr_cont", "M")
			return resp, nil
		})
	c := newTestClient(t)
	_, err := c.InquirePeriodRightsAll(context.Background(), overseas.InquirePeriodRightsParams{RghtTypeCd: "03"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "exceeded")
	assert.Equal(t, 1000, calls)
}

func TestClient_InquirePeriodRightsAll_PageErrorHasContext(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	calls := 0
	httpmock.RegisterResponder(http.MethodGet, `=~/quotations/period-rights`,
		func(req *http.Request) (*http.Response, error) {
			calls++
			if calls == 1 {
				resp := httpmock.NewStringResponse(200, loadFixtureString(t, "period_rights_page1.json"))
				resp.Header.Set("tr_cont", "M")
				return resp, nil
			}
			return httpmock.NewStringResponse(200, `{"rt_cd":"1","msg_cd":"EGW00002","msg1":"서버 에러"}`), nil
		})
	c := newTestClient(t)
	_, err := c.InquirePeriodRightsAll(context.Background(), overseas.InquirePeriodRightsParams{RghtTypeCd: "03"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "page 2")
	assert.Contains(t, err.Error(), "EGW00002")
}
