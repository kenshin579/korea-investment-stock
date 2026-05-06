package bonds_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kenshin579/korea-investment-stock/bonds"
)

func TestClient_SearchBondInfo(t *testing.T) {
	client, transport := newTestClient(t)
	transport.RegisterResponder(
		http.MethodGet,
		"=~/uapi/domestic-bond/v1/quotations/search-bond-info",
		httpmock.NewStringResponder(http.StatusOK, loadFixtureString(t, "search_bond_info_success.json")),
	)

	got, err := client.SearchBondInfo(context.Background(), bonds.SearchBondInfoParams{
		Pdno:       "KR103501GCC7",
		PrdtTypeCd: "300",
	})
	require.NoError(t, err)
	require.NotNil(t, got)

	assert.Equal(t, "KR103501GCC7", got.Pdno)
	assert.Equal(t, "국고채권 03500-5103(21-3)", got.KsdBondItemName)
	assert.Equal(t, "3.50", got.KsdRcvgBondSrfcInrt)
	assert.Equal(t, "국채", got.BondClsfKorName)
	assert.Equal(t, "Y", got.ElecSctyYn)
}

func TestClient_InquireIssueInfo(t *testing.T) {
	client, transport := newTestClient(t)
	transport.RegisterResponder(
		http.MethodGet,
		"=~/uapi/domestic-bond/v1/quotations/issue-info",
		httpmock.NewStringResponder(http.StatusOK, loadFixtureString(t, "issue_info_success.json")),
	)

	got, err := client.InquireIssueInfo(context.Background(), bonds.InquireIssueInfoParams{
		Pdno:       "KR103501GCC7",
		PrdtTypeCd: "300",
	})
	require.NoError(t, err)
	require.NotNil(t, got)

	assert.Equal(t, "KR103501GCC7", got.Pdno)
	assert.Equal(t, "국고채권 03500-5103(21-3)", got.PrdtName)
	assert.Equal(t, "3.50", got.SrfcInrt)
	assert.Equal(t, "AAA", got.KisCrdtGradText)
	assert.Equal(t, "Y", got.ElecSctyYn)
}

func TestClient_InquirePrice(t *testing.T) {
	client, transport := newTestClient(t)
	transport.RegisterResponder(
		http.MethodGet,
		"=~/uapi/domestic-bond/v1/quotations/inquire-price",
		httpmock.NewStringResponder(http.StatusOK, loadFixtureString(t, "bond_price_success.json")),
	)

	got, err := client.InquirePrice(context.Background(), bonds.InquirePriceParams{
		MarketCode: "B",
		Symbol:     "KR103501GCC7",
	})
	require.NoError(t, err)
	require.NotNil(t, got)

	prpr, _ := decimal.NewFromString("9823.50")
	vrss, _ := decimal.NewFromString("-12.50")

	assert.Equal(t, "KR103501GCC7", got.StndIscd)
	assert.True(t, prpr.Equal(got.BondPrpr))
	assert.True(t, vrss.Equal(got.BondPrdyVrss))
	assert.InDelta(t, -0.13, got.PrdyCtrt, 0.001)
	assert.Equal(t, int64(152000), got.AcmlVol)
}
