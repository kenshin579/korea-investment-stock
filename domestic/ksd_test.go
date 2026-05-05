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
