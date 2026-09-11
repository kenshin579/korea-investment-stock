package domestic_test

import (
	"context"
	"net/http"
	"strings"
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

// KIS 는 예고 없이 숫자 포맷을 바꾼 전례가 있다. pchs_amt 같은 kistypes.Int 필드 하나가
// "700000.00" 처럼 소수점이 붙어 와도 전체 응답 파싱이 깨지지 않아야 한다.
func TestClient_InquireBalance_IntegralDecimalAmount(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	fixture := loadFixtureString(t, "inquire_balance_success.json")
	modified := strings.Replace(fixture, `"pchs_amt": "700000"`, `"pchs_amt": "700000.00"`, 1)
	require.Contains(t, modified, `"pchs_amt": "700000.00"`, "치환이 실제로 적용됐는지 확인")
	require.NotEqual(t, fixture, modified)

	httpmock.RegisterResponder(http.MethodGet, `=~/trading/inquire-balance`,
		httpmock.NewStringResponder(200, modified))

	c := newTestClient(t)
	res, err := c.InquireBalance(context.Background(), domestic.InquireBalanceParams{})
	require.NoError(t, err)
	require.Len(t, res.Output1, 2)
	assert.Equal(t, int64(700000), int64(res.Output1[0].PchsAmt))
}

// 모의투자 도메인(openapivts)이면 TR ID 가 VTTC8434R 로 바뀐다.
func TestClient_InquireBalance_PaperTrID(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(http.MethodGet, `=~/trading/inquire-balance`,
		func(req *http.Request) (*http.Response, error) {
			assert.Equal(t, "VTTC8434R", req.Header.Get("tr_id"), "모의투자 TR ID")
			return httpmock.NewStringResponse(200, loadFixtureString(t, "inquire_balance_empty.json")), nil
		})

	c := newPaperTestClient(t)
	_, err := c.InquireBalance(context.Background(), domestic.InquireBalanceParams{})
	require.NoError(t, err)
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

func TestClient_InquireBalanceAll_FollowsTrCont(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	calls := 0
	httpmock.RegisterResponder(http.MethodGet, `=~/trading/inquire-balance`,
		func(req *http.Request) (*http.Response, error) {
			calls++
			q := req.URL.Query()
			switch calls {
			case 1:
				assert.Equal(t, "", req.Header.Get("tr_cont"))
				assert.Equal(t, "", q.Get("CTX_AREA_NK100"))
				resp := httpmock.NewStringResponse(200, loadFixtureString(t, "inquire_balance_page1.json"))
				resp.Header.Set("tr_cont", "M")
				return resp, nil
			case 2:
				assert.Equal(t, "N", req.Header.Get("tr_cont"), "2페이지는 tr_cont=N")
				assert.Equal(t, "FK-PAGE1", q.Get("CTX_AREA_FK100"), "이전 응답 커서 전달")
				assert.Equal(t, "NK-PAGE1", q.Get("CTX_AREA_NK100"))
				resp := httpmock.NewStringResponse(200, loadFixtureString(t, "inquire_balance_success.json"))
				resp.Header.Set("tr_cont", "D")
				return resp, nil
			}
			t.Fatalf("unexpected call #%d", calls)
			return nil, nil
		})

	c := newTestClient(t)
	// 호출자가 커서를 넣어도 첫 페이지부터 다시 읽는다.
	res, err := c.InquireBalanceAll(context.Background(), domestic.InquireBalanceParams{TrCont: "N", CtxAreaNk100: "junk"})
	require.NoError(t, err)
	assert.Equal(t, 2, calls)
	require.Len(t, res.Output1, 3, "1 + 2 건 이어 붙임")
	assert.Equal(t, "000660", res.Output1[0].Pdno)
	assert.Equal(t, "005930", res.Output1[1].Pdno)
	assert.Equal(t, "379780", res.Output1[2].Pdno)
	require.Len(t, res.Output2, 1, "요약은 마지막 페이지 것")
	assert.Equal(t, int64(916000), int64(res.Output2[0].TotEvluAmt))
	assert.Equal(t, "D", res.TrCont)
}

func TestClient_InquireBalanceAll_SinglePage(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder(http.MethodGet, `=~/trading/inquire-balance`,
		httpmock.NewStringResponder(200, loadFixtureString(t, "inquire_balance_success.json")))

	c := newTestClient(t)
	res, err := c.InquireBalanceAll(context.Background(), domestic.InquireBalanceParams{})
	require.NoError(t, err)
	assert.Len(t, res.Output1, 2)
	assert.Equal(t, 1, httpmock.GetTotalCallCount(), "tr_cont 헤더 없음 = 1페이지로 끝")
}

func TestClient_InquireBalanceAll_PageCap(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder(http.MethodGet, `=~/trading/inquire-balance`,
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(200, loadFixtureString(t, "inquire_balance_page1.json"))
			resp.Header.Set("tr_cont", "M") // 영원히 다음 있음
			return resp, nil
		})

	c := newTestClient(t)
	_, err := c.InquireBalanceAll(context.Background(), domestic.InquireBalanceParams{})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "exceeded")
	assert.Equal(t, 100, httpmock.GetTotalCallCount())
}
