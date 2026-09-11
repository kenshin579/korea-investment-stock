package overseas_test

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kenshin579/korea-investment-stock/overseas"
)

func TestClient_InquireBalance(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(http.MethodGet, `=~/overseas-stock/v1/trading/inquire-balance`,
		func(req *http.Request) (*http.Response, error) {
			q := req.URL.Query()
			assert.Equal(t, "00000000", q.Get("CANO"))
			assert.Equal(t, "00", q.Get("ACNT_PRDT_CD"))
			assert.Equal(t, "NASD", q.Get("OVRS_EXCG_CD"))
			assert.Equal(t, "USD", q.Get("TR_CRCY_CD"))
			assert.Equal(t, "", q.Get("CTX_AREA_FK200"))
			assert.Equal(t, "", q.Get("CTX_AREA_NK200"))
			assert.Equal(t, "TTTS3012R", req.Header.Get("tr_id"))
			assert.Equal(t, "", req.Header.Get("tr_cont"))
			resp := httpmock.NewStringResponse(200, loadFixtureString(t, "inquire_balance_success.json"))
			resp.Header.Set("tr_cont", "D")
			return resp, nil
		})

	c := newTestClient(t)
	res, err := c.InquireBalance(context.Background(), overseas.InquireBalanceParams{OvrsExcgCd: "NASD", TrCrcyCd: "USD"})
	require.NoError(t, err)

	require.Len(t, res.Output1, 2)
	a := res.Output1[0]
	assert.Equal(t, "AAPL", a.OvrsPdno)
	assert.Equal(t, "APPLE INC", a.OvrsItemName)
	assert.Equal(t, "NASD", a.OvrsExcgCd)
	assert.Equal(t, "USD", a.TrCrcyCd)
	assert.InDelta(t, 10, float64(a.OvrsCblcQty), 0.0001)
	assert.InDelta(t, 200.0, float64(a.PchsAvgPric), 0.0001)
	assert.InDelta(t, 2120.5, float64(a.OvrsStckEvluAmt), 0.0001)
	assert.InDelta(t, 120.5, float64(a.FrcrEvluPflsAmt), 0.0001, "부호 + 문자열")
	assert.InDelta(t, 6.02, float64(a.EvluPflsRt), 0.001)
	assert.InDelta(t, -15.0, float64(res.Output1[1].FrcrEvluPflsAmt), 0.0001, "부호 - 문자열")

	assert.InDelta(t, 2600.0, float64(res.Output2.FrcrPchsAmt1), 0.0001)
	assert.InDelta(t, 105.5, float64(res.Output2.TotEvluPflsAmt), 0.0001)
	assert.InDelta(t, 4.05, float64(res.Output2.TotPftrt), 0.001)
	assert.Equal(t, "D", res.TrCont)
}

func TestClient_InquireBalance_RequiredParams(t *testing.T) {
	c := newTestClient(t)
	_, err := c.InquireBalance(context.Background(), overseas.InquireBalanceParams{TrCrcyCd: "USD"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "OvrsExcgCd")
	_, err = c.InquireBalance(context.Background(), overseas.InquireBalanceParams{OvrsExcgCd: "NASD"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "TrCrcyCd")
}

func TestClient_InquireBalance_PaperTrID(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder(http.MethodGet, `=~/overseas-stock/v1/trading/inquire-balance`,
		func(req *http.Request) (*http.Response, error) {
			assert.Equal(t, "VTTS3012R", req.Header.Get("tr_id"))
			return httpmock.NewStringResponse(200, loadFixtureString(t, "inquire_balance_success.json")), nil
		})
	c := newPaperTestClient(t)
	_, err := c.InquireBalance(context.Background(), overseas.InquireBalanceParams{OvrsExcgCd: "NASD", TrCrcyCd: "USD"})
	require.NoError(t, err)
}

func TestClient_InquireBalance_EmptyNumericStrings(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	body := strings.Replace(loadFixtureString(t, "inquire_balance_success.json"), `"tot_pftrt": "+4.05"`, `"tot_pftrt": ""`, 1)
	require.Contains(t, body, `"tot_pftrt": ""`)
	httpmock.RegisterResponder(http.MethodGet, `=~/overseas-stock/v1/trading/inquire-balance`, httpmock.NewStringResponder(200, body))
	c := newTestClient(t)
	res, err := c.InquireBalance(context.Background(), overseas.InquireBalanceParams{OvrsExcgCd: "NASD", TrCrcyCd: "USD"})
	require.NoError(t, err)
	assert.Equal(t, 0.0, float64(res.Output2.TotPftrt), `"" → 0`)
}

func TestClient_InquireBalanceAll_FollowsTrCont(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	calls := 0
	httpmock.RegisterResponder(http.MethodGet, `=~/overseas-stock/v1/trading/inquire-balance`,
		func(req *http.Request) (*http.Response, error) {
			calls++
			q := req.URL.Query()
			switch calls {
			case 1:
				assert.Equal(t, "", req.Header.Get("tr_cont"))
				assert.Equal(t, "", q.Get("CTX_AREA_FK200"))
				assert.Equal(t, "", q.Get("CTX_AREA_NK200"))
				resp := httpmock.NewStringResponse(200, loadFixtureString(t, "inquire_balance_page1.json"))
				resp.Header.Set("tr_cont", "F")
				return resp, nil
			case 2:
				assert.Equal(t, "N", req.Header.Get("tr_cont"))
				assert.Equal(t, "FK-PAGE1", q.Get("CTX_AREA_FK200"))
				assert.Equal(t, "NK-PAGE1", q.Get("CTX_AREA_NK200"))
				resp := httpmock.NewStringResponse(200, loadFixtureString(t, "inquire_balance_success.json"))
				resp.Header.Set("tr_cont", "E")
				return resp, nil
			}
			t.Fatalf("unexpected call #%d", calls)
			return nil, nil
		})

	c := newTestClient(t)
	res, err := c.InquireBalanceAll(context.Background(), overseas.InquireBalanceParams{
		OvrsExcgCd: "NASD", TrCrcyCd: "USD", TrCont: "N", CtxAreaFk200: "junk", CtxAreaNk200: "junk",
	})
	require.NoError(t, err)
	assert.Equal(t, 2, calls)
	require.Len(t, res.Output1, 3)
	assert.Equal(t, "QQQM", res.Output1[0].OvrsPdno)
	assert.Equal(t, "AAPL", res.Output1[1].OvrsPdno)
	assert.Equal(t, "O", res.Output1[2].OvrsPdno)
	assert.InDelta(t, 2600.0, float64(res.Output2.FrcrPchsAmt1), 0.0001, "요약은 마지막 페이지")
	assert.Equal(t, "E", res.TrCont)
}

func TestClient_InquireBalanceAll_SinglePage(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder(http.MethodGet, `=~/overseas-stock/v1/trading/inquire-balance`,
		httpmock.NewStringResponder(200, loadFixtureString(t, "inquire_balance_success.json")))
	c := newTestClient(t)
	res, err := c.InquireBalanceAll(context.Background(), overseas.InquireBalanceParams{OvrsExcgCd: "NASD", TrCrcyCd: "USD"})
	require.NoError(t, err)
	assert.Len(t, res.Output1, 2)
	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestClient_InquireBalanceAll_CursorNotAdvancing(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder(http.MethodGet, `=~/overseas-stock/v1/trading/inquire-balance`,
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(200, loadFixtureString(t, "inquire_balance_success.json")) // 커서 ""
			resp.Header.Set("tr_cont", "M")
			return resp, nil
		})
	c := newTestClient(t)
	_, err := c.InquireBalanceAll(context.Background(), overseas.InquireBalanceParams{OvrsExcgCd: "NASD", TrCrcyCd: "USD"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "cursor did not advance")
	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestClient_InquireBalanceAll_PageCap(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	calls := 0
	httpmock.RegisterResponder(http.MethodGet, `=~/overseas-stock/v1/trading/inquire-balance`,
		func(req *http.Request) (*http.Response, error) {
			calls++
			body := strings.Replace(loadFixtureString(t, "inquire_balance_page1.json"), "NK-PAGE1", fmt.Sprintf("NK-%d", calls), 1)
			resp := httpmock.NewStringResponse(200, body)
			resp.Header.Set("tr_cont", "M")
			return resp, nil
		})
	c := newTestClient(t)
	_, err := c.InquireBalanceAll(context.Background(), overseas.InquireBalanceParams{OvrsExcgCd: "NASD", TrCrcyCd: "USD"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "exceeded")
	assert.Equal(t, 100, calls)
}

func TestClient_InquireBalanceAll_PageErrorHasContext(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	calls := 0
	httpmock.RegisterResponder(http.MethodGet, `=~/overseas-stock/v1/trading/inquire-balance`,
		func(req *http.Request) (*http.Response, error) {
			calls++
			if calls == 1 {
				resp := httpmock.NewStringResponse(200, loadFixtureString(t, "inquire_balance_page1.json"))
				resp.Header.Set("tr_cont", "M")
				return resp, nil
			}
			return httpmock.NewStringResponse(200, `{"rt_cd":"1","msg_cd":"EGW00002","msg1":"서버 에러"}`), nil
		})
	c := newTestClient(t)
	_, err := c.InquireBalanceAll(context.Background(), overseas.InquireBalanceParams{OvrsExcgCd: "NASD", TrCrcyCd: "USD"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "page 2")
	assert.Contains(t, err.Error(), "EGW00002")
}
