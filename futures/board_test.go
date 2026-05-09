package futures_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jarcoal/httpmock"

	"github.com/kenshin579/korea-investment-stock/futures"
)

// ─── EP8: DisplayBoardTop ─────────────────────────────────────────────────────

func TestClient_DisplayBoardTop(t *testing.T) {
	client, transport := newTestClient(t)
	transport.RegisterResponder(
		http.MethodGet,
		"=~/uapi/domestic-futureoption/v1/quotations/display-board-top",
		httpmock.NewStringResponder(http.StatusOK, loadFixtureString(t, "display_board_top_success.json")),
	)

	got, err := client.DisplayBoardTop(context.Background(), futures.DisplayBoardTopParams{
		MarketCode: "F",
		Code:       "101V06",
	})
	require.NoError(t, err)
	require.NotNil(t, got)

	// output1 assertions
	assert.Equal(t, "KOSPI200선물", got.Output1.HtsKorIsnm)
	unasPrpr, _ := decimal.NewFromString("360.50")
	assert.True(t, unasPrpr.Equal(got.Output1.UnasPrpr))
	assert.Equal(t, "2", got.Output1.UnasPrdyVrssSign)
	assert.InDelta(t, 0.35, got.Output1.UnasPrdyCtrt, 0.001)
	assert.Equal(t, int64(523410), got.Output1.UnasAcmlVol)
	futsPrpr, _ := decimal.NewFromString("360.75")
	assert.True(t, futsPrpr.Equal(got.Output1.FutsPrpr))
	assert.Equal(t, "2", got.Output1.PrdyVrssSign)
	assert.InDelta(t, 0.42, got.Output1.FutsPrdyCtrt, 0.001)

	// output2 assertions — 월물별 잔존일수 (hts_rmnn_dynu 1 필드)
	require.Len(t, got.Output2, 3)
	assert.Equal(t, "30", got.Output2[0].HtsRmnnDynu)
	assert.Equal(t, "93", got.Output2[1].HtsRmnnDynu)
	assert.Equal(t, "184", got.Output2[2].HtsRmnnDynu)
}

// ─── EP9: DisplayBoardFutures ─────────────────────────────────────────────────

func TestClient_DisplayBoardFutures(t *testing.T) {
	client, transport := newTestClient(t)
	transport.RegisterResponder(
		http.MethodGet,
		"=~/uapi/domestic-futureoption/v1/quotations/display-board-futures",
		httpmock.NewStringResponder(http.StatusOK, loadFixtureString(t, "display_board_futures_success.json")),
	)

	got, err := client.DisplayBoardFutures(context.Background(), futures.DisplayBoardFuturesParams{
		MarketCode: "F",
	})
	require.NoError(t, err)
	require.NotNil(t, got)

	// output1 assertions
	require.Len(t, got.Output1, 2)
	item := got.Output1[0]
	assert.Equal(t, "101V06", item.FutsShrnIscd)
	assert.Equal(t, "KOSPI200선물 2506", item.HtsKorIsnm)
	futsPrpr, _ := decimal.NewFromString("360.75")
	assert.True(t, futsPrpr.Equal(item.FutsPrpr))
	assert.Equal(t, int64(523410), item.AcmlVol)
	assert.Equal(t, int64(310250), item.HtsOtstStplQty)
	assert.Equal(t, "30", item.HtsRmnnDynu)
	assert.Equal(t, int64(1250), item.TotalAskpRsqn)
	assert.InDelta(t, 0.43, item.AntcCntgPrdyCtrt, 0.001)

	// second item
	item2 := got.Output1[1]
	assert.Equal(t, "101V09", item2.FutsShrnIscd)
	assert.Equal(t, "93", item2.HtsRmnnDynu)
}

// ─── EP10: DisplayBoardOptionList ────────────────────────────────────────────

func TestClient_DisplayBoardOptionList(t *testing.T) {
	client, transport := newTestClient(t)
	transport.RegisterResponder(
		http.MethodGet,
		"=~/uapi/domestic-futureoption/v1/quotations/display-board-option-list",
		httpmock.NewStringResponder(http.StatusOK, loadFixtureString(t, "display_board_option_list_success.json")),
	)

	got, err := client.DisplayBoardOptionList(context.Background(), futures.DisplayBoardOptionListParams{})
	require.NoError(t, err)
	require.NotNil(t, got)

	// output1 assertions
	require.Len(t, got.Output1, 3)
	assert.Equal(t, "202506", got.Output1[0].MtrtYymmCode)
	assert.Equal(t, "202506", got.Output1[0].MtrtYymm)
	assert.Equal(t, "202509", got.Output1[1].MtrtYymmCode)
	assert.Equal(t, "202512", got.Output1[2].MtrtYymmCode)
}

// ─── EP11: DisplayBoardCallput ────────────────────────────────────────────────

func TestClient_DisplayBoardCallput(t *testing.T) {
	client, transport := newTestClient(t)
	transport.RegisterResponder(
		http.MethodGet,
		"=~/uapi/domestic-futureoption/v1/quotations/display-board-callput",
		httpmock.NewStringResponder(http.StatusOK, loadFixtureString(t, "display_board_callput_success.json")),
	)

	got, err := client.DisplayBoardCallput(context.Background(), futures.DisplayBoardCallputParams{
		MtrtCnt: "202506",
	})
	require.NoError(t, err)
	require.NotNil(t, got)

	// output1 (콜옵션) assertions
	require.Len(t, got.Output1, 2)
	call := got.Output1[0]
	acpr, _ := decimal.NewFromString("360.00")
	assert.True(t, acpr.Equal(call.Acpr))
	assert.Equal(t, "201V06360", call.OptnShrnIscd)
	optnPrpr, _ := decimal.NewFromString("4.25")
	assert.True(t, optnPrpr.Equal(call.OptnPrpr))
	assert.Equal(t, "2", call.PrdyVrssSign)
	assert.InDelta(t, 13.33, call.OptnPrdyCtrt, 0.001)
	assert.Equal(t, int64(12500), call.AcmlVol)
	assert.Equal(t, int64(8250), call.HtsOtstStplQty)
	assert.InDelta(t, 0.5421, call.DeltaVal, 0.0001)
	assert.InDelta(t, 18.25, call.HtsIntsVltl, 0.01)
	assert.Equal(t, "ATM", call.AtmClsName)
	assert.Equal(t, int64(980), call.TotalAskpRsqn)
	antcCnpr, _ := decimal.NewFromString("360.80")
	assert.True(t, antcCnpr.Equal(call.FutsAntcCnpr))

	// second call item
	call2 := got.Output1[1]
	acpr2, _ := decimal.NewFromString("362.50")
	assert.True(t, acpr2.Equal(call2.Acpr))
	assert.Equal(t, "OTM", call2.AtmClsName)

	// output2 (풋옵션) assertions
	require.Len(t, got.Output2, 1)
	put := got.Output2[0]
	acprPut, _ := decimal.NewFromString("360.00")
	assert.True(t, acprPut.Equal(put.Acpr))
	assert.Equal(t, "301V06360", put.OptnShrnIscd)
	assert.Equal(t, "5", put.PrdyVrssSign)
	assert.InDelta(t, -0.4579, put.DeltaVal, 0.0001)
	assert.InDelta(t, -0.0298, put.Rho, 0.0001)
	assert.Equal(t, "ATM", put.AtmClsName)
}

func TestClient_DisplayBoardCallput_InvalidJSON(t *testing.T) {
	client, transport := newTestClient(t)
	transport.RegisterResponder(
		http.MethodGet,
		"=~/uapi/domestic-futureoption/v1/quotations/display-board-callput",
		httpmock.NewStringResponder(http.StatusOK, "invalid json"),
	)

	got, err := client.DisplayBoardCallput(context.Background(), futures.DisplayBoardCallputParams{
		MtrtCnt: "202506",
	})
	require.Error(t, err)
	require.Nil(t, got)
}
