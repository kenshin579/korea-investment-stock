package domestic_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kenshin579/korea-investment-stock/domestic"
)

func TestClient_InquireBalance(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(http.MethodGet, `=~/trading/inquire-balance`,
		func(req *http.Request) (*http.Response, error) {
			q := req.URL.Query()
			assert.Equal(t, "00000000", q.Get("CANO"), "계좌번호 앞 8자리")
			assert.Equal(t, "00", q.Get("ACNT_PRDT_CD"), "계좌상품코드")
			assert.Equal(t, "N", q.Get("AFHR_FLPR_YN"), "기본값")
			assert.Equal(t, "02", q.Get("INQR_DVSN"), "기본값 종목별")
			assert.Equal(t, "01", q.Get("UNPR_DVSN"))
			assert.Equal(t, "N", q.Get("FUND_STTL_ICLD_YN"))
			assert.Equal(t, "N", q.Get("FNCG_AMT_AUTO_RDPT_YN"))
			assert.Equal(t, "00", q.Get("PRCS_DVSN"))
			assert.Equal(t, "", q.Get("CTX_AREA_FK100"))
			assert.Equal(t, "TTTC8434R", req.Header.Get("tr_id"))
			resp := httpmock.NewStringResponse(200, loadFixtureString(t, "inquire_balance_success.json"))
			resp.Header.Set("tr_cont", "D")
			return resp, nil
		})

	c := newTestClient(t)
	res, err := c.InquireBalance(context.Background(), domestic.InquireBalanceParams{})
	require.NoError(t, err)
	require.NotNil(t, res)

	require.Len(t, res.Output1, 2)
	assert.Equal(t, "005930", res.Output1[0].Pdno)
	assert.Equal(t, "삼성전자", res.Output1[0].PrdtName)
	assert.Equal(t, int64(10), int64(res.Output1[0].HldgQty))
	assert.InDelta(t, 70000.0, float64(res.Output1[0].PchsAvgPric), 0.0001)
	assert.Equal(t, int64(700000), int64(res.Output1[0].PchsAmt))
	assert.Equal(t, int64(750000), int64(res.Output1[0].EvluAmt))
	assert.Equal(t, int64(50000), int64(res.Output1[0].EvluPflsAmt))
	assert.InDelta(t, 7.14, float64(res.Output1[0].EvluPflsRt), 0.001)
	assert.Equal(t, int64(-100), int64(res.Output1[1].BfdyCprsIcdc), "음수 정수")
	assert.InDelta(t, -0.45, float64(res.Output1[1].FlttRt), 0.001)

	require.Len(t, res.Output2, 1)
	assert.Equal(t, int64(100000), int64(res.Output2[0].DncaTotAmt))
	assert.Equal(t, int64(816000), int64(res.Output2[0].EvluAmtSmtlAmt))
	assert.Equal(t, int64(916000), int64(res.Output2[0].TotEvluAmt))
	assert.Equal(t, "D", res.TrCont)
}

func TestClient_InquireBalance_Empty(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder(http.MethodGet, `=~/trading/inquire-balance`,
		httpmock.NewStringResponder(200, loadFixtureString(t, "inquire_balance_empty.json")))

	c := newTestClient(t)
	res, err := c.InquireBalance(context.Background(), domestic.InquireBalanceParams{})
	require.NoError(t, err)
	assert.Empty(t, res.Output1)
	require.Len(t, res.Output2, 1)
	assert.Equal(t, int64(0), int64(res.Output2[0].CmaEvluAmt), `"" → 0`)
	assert.Equal(t, 0.0, float64(res.Output2[0].AsstIcdcErngRt), `"" → 0`)
	assert.Equal(t, "", res.TrCont, "헤더 없으면 빈 값 = 마지막")
}

func TestClient_InquireBalance_ParamsOverride(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder(http.MethodGet, `=~/trading/inquire-balance`,
		func(req *http.Request) (*http.Response, error) {
			q := req.URL.Query()
			assert.Equal(t, "01", q.Get("INQR_DVSN"))
			assert.Equal(t, "X", q.Get("AFHR_FLPR_YN"))
			assert.Equal(t, "FK", q.Get("CTX_AREA_FK100"))
			assert.Equal(t, "NK", q.Get("CTX_AREA_NK100"))
			assert.Equal(t, "N", req.Header.Get("tr_cont"))
			return httpmock.NewStringResponse(200, loadFixtureString(t, "inquire_balance_empty.json")), nil
		})
	c := newTestClient(t)
	_, err := c.InquireBalance(context.Background(), domestic.InquireBalanceParams{
		InqrDvsn: "01", AfhrFlprYn: "X", CtxAreaFk100: "FK", CtxAreaNk100: "NK", TrCont: "N",
	})
	require.NoError(t, err)
}
