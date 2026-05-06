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

func TestClient_InquireCcnl(t *testing.T) {
	client, transport := newTestClient(t)
	transport.RegisterResponder(
		http.MethodGet,
		"=~/uapi/domestic-bond/v1/quotations/inquire-ccnl",
		httpmock.NewStringResponder(http.StatusOK, loadFixtureString(t, "bond_ccnl_success.json")),
	)

	got, err := client.InquireCcnl(context.Background(), bonds.InquireCcnlParams{
		MarketCode: "B",
		Symbol:     "KR103501GCC7",
	})
	require.NoError(t, err)
	require.NotNil(t, got)

	ccnlPrpr, _ := decimal.NewFromString("9823.50")

	assert.Equal(t, "141523", got.StckCntgHour)
	assert.True(t, ccnlPrpr.Equal(got.BondPrpr))
	assert.Equal(t, int64(500), got.CntgVol)
	assert.Equal(t, int64(152000), got.AcmlVol)
	assert.InDelta(t, -0.13, got.PrdyCtrt, 0.001)
}

func TestClient_InquireAskingPrice(t *testing.T) {
	client, transport := newTestClient(t)
	transport.RegisterResponder(
		http.MethodGet,
		"=~/uapi/domestic-bond/v1/quotations/inquire-asking-price",
		httpmock.NewStringResponder(http.StatusOK, loadFixtureString(t, "bond_asking_price_success.json")),
	)

	got, err := client.InquireAskingPrice(context.Background(), bonds.InquireAskingPriceParams{
		MarketCode: "B",
		Symbol:     "KR103501GCC7",
	})
	require.NoError(t, err)
	require.NotNil(t, got)

	askp1, _ := decimal.NewFromString("9824.00")
	bidp1, _ := decimal.NewFromString("9823.00")

	assert.Equal(t, "141600", got.AsprAcptHour)
	assert.True(t, askp1.Equal(got.BondAskp1))
	assert.True(t, bidp1.Equal(got.BondBidp1))
	assert.Equal(t, int64(10000), got.AskpRsqn1)
	assert.Equal(t, int64(28000), got.TotalAskpRsqn)
	assert.InDelta(t, 3.61, got.SelnErnnRate1, 0.001)
	assert.InDelta(t, 3.62, got.ShnuErnnRate1, 0.001)
}

func TestClient_InquireDailyPrice(t *testing.T) {
	client, transport := newTestClient(t)
	transport.RegisterResponder(
		http.MethodGet,
		"=~/uapi/domestic-bond/v1/quotations/inquire-daily-price",
		httpmock.NewStringResponder(http.StatusOK, loadFixtureString(t, "bond_daily_price_success.json")),
	)

	got, err := client.InquireDailyPrice(context.Background(), bonds.InquireDailyPriceParams{
		MarketCode: "B",
		Symbol:     "KR103501GCC7",
	})
	require.NoError(t, err)
	require.NotNil(t, got)

	dailyPrpr, _ := decimal.NewFromString("9823.50")
	dailyHgpr, _ := decimal.NewFromString("9830.00")

	assert.Equal(t, "20260505", got.StckBsopDate)
	assert.True(t, dailyPrpr.Equal(got.BondPrpr))
	assert.Equal(t, int64(152000), got.AcmlVol)
	assert.InDelta(t, -0.13, got.PrdyCtrt, 0.001)
	assert.True(t, dailyHgpr.Equal(got.BondHgpr))
}

func TestClient_InquireAvgUnit(t *testing.T) {
	client, transport := newTestClient(t)
	transport.RegisterResponder(
		http.MethodGet,
		"=~/uapi/domestic-bond/v1/quotations/avg-unit",
		httpmock.NewStringResponder(http.StatusOK, loadFixtureString(t, "avg_unit_success.json")),
	)

	got, err := client.InquireAvgUnit(context.Background(), bonds.InquireAvgUnitParams{
		InqrStrtDt:   "20260401",
		InqrEndDt:    "20260505",
		Pdno:         "KR1010012345",
		PrdtTypeCd:   "301",
		VrfcKindCd:   "01",
		CtxAreaNk30:  "",
		CtxAreaFk100: "",
	})
	require.NoError(t, err)
	require.NotNil(t, got)

	// output1 array assertions
	require.GreaterOrEqual(t, len(got.Output1), 1)
	u := got.Output1[0]
	assert.Equal(t, "20260505", u.EvluDt)
	assert.Equal(t, "KR1010012345", u.Pdno)
	avgEvluUnpr, _ := decimal.NewFromString("9851.25")
	assert.True(t, avgEvluUnpr.Equal(u.AvgEvluUnpr))
	assert.Equal(t, "AAA", u.KisCrdtGradText)
	assert.InDelta(t, 3.500, u.KisErngRt, 0.001)

	// output2 array assertions
	require.GreaterOrEqual(t, len(got.Output2), 1)
	a := got.Output2[0]
	assert.Equal(t, "20260505", a.EvluDt)
	assert.Equal(t, int64(9851250), a.AvgEvluAmt)

	// output3 array assertions
	require.GreaterOrEqual(t, len(got.Output3), 1)
	p := got.Output3[0]
	assert.Equal(t, "20260505", p.EvluDt)
	avgEvluUnitPric, _ := decimal.NewFromString("9851.25")
	assert.True(t, avgEvluUnitPric.Equal(p.AvgEvluUnitPric))
	assert.Equal(t, "KRW", p.KisCrcyCd)
}

func TestClient_InquireDailyItemchartprice(t *testing.T) {
	client, transport := newTestClient(t)
	transport.RegisterResponder(
		http.MethodGet,
		"=~/uapi/domestic-bond/v1/quotations/inquire-daily-itemchartprice",
		httpmock.NewStringResponder(http.StatusOK, loadFixtureString(t, "bond_daily_itemchartprice_success.json")),
	)

	got, err := client.InquireDailyItemchartprice(context.Background(), bonds.InquireDailyItemchartpriceParams{
		MarketCode: "B",
		Symbol:     "KR103501GCC7",
	})
	require.NoError(t, err)
	require.NotNil(t, got)

	require.GreaterOrEqual(t, len(got.Output), 2)

	item0 := got.Output[0]
	prpr0, _ := decimal.NewFromString("9823.50")
	oprc0, _ := decimal.NewFromString("9825.00")

	assert.Equal(t, "20260505", item0.StckBsopDate)
	assert.True(t, prpr0.Equal(item0.BondPrpr))
	assert.True(t, oprc0.Equal(item0.BondOprc))
	assert.Equal(t, int64(152000), item0.AcmlVol)

	item1 := got.Output[1]
	prpr1, _ := decimal.NewFromString("9836.00")
	assert.Equal(t, "20260502", item1.StckBsopDate)
	assert.True(t, prpr1.Equal(item1.BondPrpr))
	assert.Equal(t, int64(98000), item1.AcmlVol)
}
