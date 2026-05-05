// File: domestic/ksd_test.go
package domestic_test

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kenshin579/korea-investment-stock/domestic"
)

func TestClient_InquireKsdDividend(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/ksdinfo/dividend`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "ksd_dividend_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireKsdDividend(context.Background(), domestic.InquireKsdDividendParams{
		FromDate: "20260101",
		ToDate:   "20260505",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, "0", capturedQuery.Get("GB1"))
	assert.Equal(t, "20260101", capturedQuery.Get("F_DT"))
	assert.Equal(t, "20260505", capturedQuery.Get("T_DT"))

	require.Len(t, res.Output1, 2)
	assert.Equal(t, "005930", res.Output1[0].ShtCd)
	assert.Equal(t, "삼성전자", res.Output1[0].IsinName)
	assert.Equal(t, "20260331", res.Output1[0].RecordDate)
}

func TestClient_InquireKsdBonusIssue(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/ksdinfo/bonus-issue`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "ksd_bonus_issue_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireKsdBonusIssue(context.Background(), domestic.InquireKsdBonusIssueParams{
		FromDate: "20260101",
		ToDate:   "20260505",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, "20260101", capturedQuery.Get("F_DT"))
	assert.Equal(t, "20260505", capturedQuery.Get("T_DT"))

	require.Len(t, res.Output1, 2)
	assert.Equal(t, "005930", res.Output1[0].ShtCd)
	assert.Equal(t, "삼성전자", res.Output1[0].IsinName)
	assert.Equal(t, "20260315", res.Output1[0].RecordDate)
}

func TestClient_InquireKsdPaidinCapin(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/ksdinfo/paidin-capin`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "ksd_paidin_capin_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireKsdPaidinCapin(context.Background(), domestic.InquireKsdPaidinCapinParams{
		FromDate: "20260101",
		ToDate:   "20260505",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, "1", capturedQuery.Get("GB1")) // default
	assert.Equal(t, "20260101", capturedQuery.Get("F_DT"))

	require.Len(t, res.Output, 2) // output (not Output1)
	assert.Equal(t, "005930", res.Output[0].ShtCd)
	assert.Equal(t, "삼성전자", res.Output[0].IsinName)
	assert.Equal(t, "68000", res.Output[0].FixPrice)
}

func TestClient_InquireKsdSharehldMeet(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/ksdinfo/sharehld-meet`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "ksd_sharehld_meet_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireKsdSharehldMeet(context.Background(), domestic.InquireKsdSharehldMeetParams{
		FromDate: "20260101",
		ToDate:   "20260505",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, "20260101", capturedQuery.Get("F_DT"))
	assert.Equal(t, "20260505", capturedQuery.Get("T_DT"))

	require.Len(t, res.Output1, 2)
	assert.Equal(t, "005930", res.Output1[0].ShtCd)
	assert.Equal(t, "정기주총", res.Output1[0].GenMeetType)
	assert.Equal(t, "20260326", res.Output1[0].GenMeetDt)
}

func TestClient_InquireKsdMergerSplit(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/ksdinfo/merger-split`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "ksd_merger_split_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireKsdMergerSplit(context.Background(), domestic.InquireKsdMergerSplitParams{
		FromDate: "20260101",
		ToDate:   "20260505",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, "20260101", capturedQuery.Get("F_DT"))

	require.Len(t, res.Output1, 2)
	assert.Equal(t, "005930", res.Output1[0].ShtCd)
	assert.Equal(t, "삼성전자", res.Output1[0].CustNm)       // cust_nm (합병측)
	assert.Equal(t, "흡수대상회사A", res.Output1[0].OppCustNm) // opp_cust_nm (피합병측)
	assert.Equal(t, "흡수합병", res.Output1[0].MergeType)
}

func TestClient_InquireKsdRevSplit(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/ksdinfo/rev-split`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "ksd_rev_split_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireKsdRevSplit(context.Background(), domestic.InquireKsdRevSplitParams{
		FromDate: "20260101",
		ToDate:   "20260505",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, "0", capturedQuery.Get("MARKET_GB")) // default
	assert.Equal(t, "20260101", capturedQuery.Get("F_DT"))

	require.Len(t, res.Output1, 2)
	assert.Equal(t, "005930", res.Output1[0].ShtCd)
	assert.Equal(t, "100", res.Output1[0].InterBfFaceAmt)
	assert.Equal(t, "500", res.Output1[0].InterAfFaceAmt)
}

func TestClient_InquireKsdForfeit(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/ksdinfo/forfeit`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "ksd_forfeit_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireKsdForfeit(context.Background(), domestic.InquireKsdForfeitParams{
		FromDate: "20260101",
		ToDate:   "20260505",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, "20260505", capturedQuery.Get("T_DT"))

	require.Len(t, res.Output1, 2)
	assert.Equal(t, "005930", res.Output1[0].ShtCd)
	assert.Equal(t, "68000", res.Output1[0].SubscrPrice)
	assert.Equal(t, "한국투자증권", res.Output1[0].LeadMgr)
}

func TestClient_InquireKsdMandDeposit(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/ksdinfo/mand-deposit`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "ksd_mand_deposit_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireKsdMandDeposit(context.Background(), domestic.InquireKsdMandDepositParams{
		FromDate: "20260101",
		ToDate:   "20260505",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, "20260505", capturedQuery.Get("T_DT"))

	require.Len(t, res.Output1, 2)
	assert.Equal(t, "005930", res.Output1[0].ShtCd)
	assert.Equal(t, "20260101", res.Output1[0].DepoDate) // depo_date (not record_date)
	assert.Equal(t, "의무보호예수", res.Output1[0].DepoReason)
	assert.Equal(t, "0.84", res.Output1[0].TotIssueQtyPerRate)
}

func TestClient_InquireKsdCapDcrs(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/ksdinfo/cap-dcrs`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "ksd_cap_dcrs_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireKsdCapDcrs(context.Background(), domestic.InquireKsdCapDcrsParams{
		FromDate: "20260101",
		ToDate:   "20260505",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, "20260101", capturedQuery.Get("F_DT"))

	require.Len(t, res.Output1, 2)
	assert.Equal(t, "005930", res.Output1[0].ShtCd)
	assert.Equal(t, "유상감자", res.Output1[0].ReduceCapType)
	assert.Equal(t, "주식병합", res.Output1[0].CompWay)
}

func TestClient_InquireKsdPurreq(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/ksdinfo/purreq`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "ksd_purreq_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireKsdPurreq(context.Background(), domestic.InquireKsdPurreqParams{
		FromDate: "20260101",
		ToDate:   "20260505",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, "20260505", capturedQuery.Get("T_DT"))

	require.Len(t, res.Output1, 2)
	assert.Equal(t, "005930", res.Output1[0].ShtCd)
	assert.Equal(t, "69000", res.Output1[0].BuyReqPrice)
	assert.Equal(t, "20260326", res.Output1[0].GetMeetDt)
}

func TestClient_InquireKsdListInfo(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/ksdinfo/list-info`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "ksd_list_info_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireKsdListInfo(context.Background(), domestic.InquireKsdListInfoParams{
		FromDate: "20260101",
		ToDate:   "20260505",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, "20260101", capturedQuery.Get("F_DT"))

	require.Len(t, res.Output1, 2)
	assert.Equal(t, "005930", res.Output1[0].ShtCd)
	assert.Equal(t, "20260102", res.Output1[0].ListDt) // list_dt (not record_date)
	assert.Equal(t, "유상증자", res.Output1[0].IssueType)
	assert.Equal(t, "68000", res.Output1[0].IssuePrice)
}
