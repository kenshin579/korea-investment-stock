package overseasfutures_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jarcoal/httpmock"

	"github.com/kenshin579/korea-investment-stock/overseasfutures"
)

// ─── EP1: InvestorUnpdTrend ──────────────────────────────────────────────────

func TestClient_InvestorUnpdTrend(t *testing.T) {
	client, transport := newTestClient(t)
	transport.RegisterResponder(
		http.MethodGet,
		"=~/uapi/overseas-futureoption/v1/quotations/investor-unpd-trend",
		httpmock.NewStringResponder(http.StatusOK, loadFixtureString(t, "investor_unpd_trend_success.json")),
	)

	params := overseasfutures.InvestorUnpdTrendParams{
		ProdIscd:  "ES",
		BsopDate:  "20240513",
		UpmuGubun: "0",
		CtsKey:    "",
	}
	got, err := client.InvestorUnpdTrend(context.Background(), params)
	require.NoError(t, err)
	require.NotNil(t, got)

	// output1 assertions
	assert.Equal(t, "3", got.Output1.RowCnt)

	// output2 배열 assertions
	require.Len(t, got.Output2, 3)

	item0 := got.Output2[0]
	assert.Equal(t, "ES", item0.ProdIscd)
	assert.Equal(t, "13874A", item0.CftcIscd)
	assert.Equal(t, "20240507", item0.BsopDate)
	assert.Equal(t, "318456", item0.BidpSpec)
	assert.Equal(t, "98765", item0.AskpSpec)
	assert.Equal(t, "12345", item0.SpreadSpec)
	assert.Equal(t, "234567", item0.BidpHedge)
	assert.Equal(t, "123456", item0.AskpHedge)
	assert.Equal(t, "2765432", item0.HtsOtstSmtn)
	assert.Equal(t, "5678", item0.BidpMissing)
	assert.Equal(t, "4321", item0.AskpMissing)
	assert.Equal(t, "98765", item0.BidpSpecCust)
	assert.Equal(t, "76543", item0.AskpSpecCust)
	assert.Equal(t, "3456", item0.SpreadSpecCust)
	assert.Equal(t, "54321", item0.BidpHedgeCust)
	assert.Equal(t, "43210", item0.AskpHedgeCust)
	assert.Equal(t, "654321", item0.CustSmtn)

	// 두 번째 항목
	item1 := got.Output2[1]
	assert.Equal(t, "20240430", item1.BsopDate)
	assert.Equal(t, "315000", item1.BidpSpec)

	// 세 번째 항목
	item2 := got.Output2[2]
	assert.Equal(t, "20240423", item2.BsopDate)
}

func TestClient_InvestorUnpdTrend_InvalidJSON(t *testing.T) {
	client, transport := newTestClient(t)
	transport.RegisterResponder(
		http.MethodGet,
		"=~/uapi/overseas-futureoption/v1/quotations/investor-unpd-trend",
		httpmock.NewStringResponder(http.StatusOK, `{"rt_cd":"0","output1":"bad"}`),
	)
	params := overseasfutures.InvestorUnpdTrendParams{ProdIscd: "ES", BsopDate: "20240513", UpmuGubun: "0"}
	_, err := client.InvestorUnpdTrend(context.Background(), params)
	require.Error(t, err)
}
