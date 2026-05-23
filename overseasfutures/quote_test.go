package overseasfutures_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jarcoal/httpmock"

	"github.com/kenshin579/korea-investment-stock/overseasfutures"
)

// ─── EP10: InquirePrice ──────────────────────────────────────────────────────

func TestClient_InquirePrice(t *testing.T) {
	client, transport := newTestClient(t)
	transport.RegisterResponder(
		http.MethodGet,
		"=~/uapi/overseas-futureoption/v1/quotations/inquire-price",
		httpmock.NewStringResponder(http.StatusOK, loadFixtureString(t, "inquire_price_success.json")),
	)

	got, err := client.InquirePrice(context.Background(), "CNHU24")
	require.NoError(t, err)
	require.NotNil(t, got)

	// output1 assertions
	assert.Equal(t, "20240402", got.Output1.ProcDate)
	assert.Equal(t, "160023", got.Output1.ProcTime)

	highPrice, _ := decimal.NewFromString("5245.75")
	assert.True(t, highPrice.Equal(got.Output1.HighPrice))

	openPrice, _ := decimal.NewFromString("5230.50")
	assert.True(t, openPrice.Equal(got.Output1.OpenPrice))

	lastPrice, _ := decimal.NewFromString("5204.34")
	assert.True(t, lastPrice.Equal(got.Output1.LastPrice))

	assert.Equal(t, int64(1345678), got.Output1.Vol)
	assert.Equal(t, "5", got.Output1.PrevDiffFlag)
	assert.InDelta(t, -1.30, float64(got.Output1.PrevDiffRate), 0.001)

	assert.Equal(t, int64(45), got.Output1.BidQntt)
	assert.Equal(t, int64(38), got.Output1.AskQntt)
	assert.Equal(t, "CME", got.Output1.ExchCd)
	assert.Equal(t, "USD", got.Output1.CrcCd)
	assert.Equal(t, "20240621", got.Output1.ExprDate)
	assert.Equal(t, "43", got.Output1.RemnCnt)
	assert.Equal(t, int64(5), got.Output1.LastQntt)
	assert.Equal(t, int64(653), got.Output1.TotAskQntt)
	assert.Equal(t, int64(664), got.Output1.TotBidQntt)

	tickSize, _ := decimal.NewFromString("0.25")
	assert.True(t, tickSize.Equal(got.Output1.TickSize))

	assert.Equal(t, "20240509", got.Output1.OpenDate)
	assert.Equal(t, "180000", got.Output1.OpenTime)
	assert.Equal(t, "20240510", got.Output1.CloseDate)
	assert.Equal(t, "170000", got.Output1.CloseTime)
	assert.Equal(t, "20240402", got.Output1.Sbsnsdate)

	sttlPrice, _ := decimal.NewFromString("5215.25")
	assert.True(t, sttlPrice.Equal(got.Output1.SttlPrice))
}

func TestClient_InquirePrice_InvalidJSON(t *testing.T) {
	client, transport := newTestClient(t)
	transport.RegisterResponder(
		http.MethodGet,
		"=~/uapi/overseas-futureoption/v1/quotations/inquire-price",
		httpmock.NewStringResponder(http.StatusOK, `{"rt_cd":"0","output1":"not-an-object"}`),
	)
	_, err := client.InquirePrice(context.Background(), "CNHU24")
	require.Error(t, err)
}

// ─── EP9: StockDetail ────────────────────────────────────────────────────────

func TestClient_StockDetail(t *testing.T) {
	client, transport := newTestClient(t)
	transport.RegisterResponder(
		http.MethodGet,
		"=~/uapi/overseas-futureoption/v1/quotations/stock-detail",
		httpmock.NewStringResponder(http.StatusOK, loadFixtureString(t, "stock_detail_success.json")),
	)

	got, err := client.StockDetail(context.Background(), "CNHU24")
	require.NoError(t, err)
	require.NotNil(t, got)

	// output1 assertions
	assert.Equal(t, "CME", got.Output1.ExchCd)

	tickSz, _ := decimal.NewFromString("0.25")
	assert.True(t, tickSz.Equal(got.Output1.TickSz))

	assert.Equal(t, "2", got.Output1.DispDigit)

	trstMgn, _ := decimal.NewFromString("15500.00")
	assert.True(t, trstMgn.Equal(got.Output1.TrstMgn))

	assert.Equal(t, "20240621", got.Output1.SttlDate)

	prevPrice, _ := decimal.NewFromString("5272.75")
	assert.True(t, prevPrice.Equal(got.Output1.PrevPrice))

	assert.Equal(t, "USD", got.Output1.CrcCd)
	assert.Equal(t, "ES", got.Output1.ClasCd)

	tickVal, _ := decimal.NewFromString("12.50")
	assert.True(t, tickVal.Equal(got.Output1.TickVal))

	assert.Equal(t, "20240509", got.Output1.MrktOpenDate)
	assert.Equal(t, "180000", got.Output1.MrktOpenTime)
	assert.Equal(t, "20240510", got.Output1.MrktCloseDate)
	assert.Equal(t, "170000", got.Output1.MrktCloseTime)
	assert.Equal(t, "20231222", got.Output1.TrdFrDate)
	assert.Equal(t, "20240621", got.Output1.ExprDate)
	assert.Equal(t, "20240621", got.Output1.TrdToDate)
	assert.Equal(t, "43", got.Output1.RemnCnt)
	assert.Equal(t, "1", got.Output1.StatTp)

	ctrtSize, _ := decimal.NewFromString("50.00")
	assert.True(t, ctrtSize.Equal(got.Output1.CtrtSize))

	assert.Equal(t, "CASH", got.Output1.StlTp)
	assert.Equal(t, "", got.Output1.SprdSrsCd1)
	assert.Equal(t, "", got.Output1.SprdSrsCd2)
}

// ─── EP3: SearchContractDetail ───────────────────────────────────────────────

func TestClient_SearchContractDetail(t *testing.T) {
	client, transport := newTestClient(t)
	transport.RegisterResponder(
		http.MethodGet,
		"=~/uapi/overseas-futureoption/v1/quotations/search-contract-detail",
		httpmock.NewStringResponder(http.StatusOK, loadFixtureString(t, "search_contract_detail_success.json")),
	)

	params := overseasfutures.SearchContractDetailParams{
		Codes: []string{"ESM24", "NQM24"},
	}
	got, err := client.SearchContractDetail(context.Background(), params)
	require.NoError(t, err)
	require.NotNil(t, got)

	// output2 배열 assertions
	require.Len(t, got.Output2, 2)

	item0 := got.Output2[0]
	assert.Equal(t, "CME", item0.ExchCd)
	assert.Equal(t, "ES", item0.ClasCd)
	assert.Equal(t, "USD", item0.CrcCd)

	sttlPrice0, _ := decimal.NewFromString("5230.50")
	assert.True(t, sttlPrice0.Equal(item0.SttlPrice))

	assert.Equal(t, "20240621", item0.SttlDate)

	trstMgn0, _ := decimal.NewFromString("15500.00")
	assert.True(t, trstMgn0.Equal(item0.TrstMgn))

	tickSz0, _ := decimal.NewFromString("0.25")
	assert.True(t, tickSz0.Equal(item0.TickSz))

	tickVal0, _ := decimal.NewFromString("12.50")
	assert.True(t, tickVal0.Equal(item0.TickVal))

	ctrtSize0, _ := decimal.NewFromString("50.00")
	assert.True(t, ctrtSize0.Equal(item0.CtrtSize))

	assert.Equal(t, "CASH", item0.StlTp)
	assert.Equal(t, "GLOBEX", item0.SubExchNm)
	assert.Equal(t, "1", item0.StatTp)

	// 두 번째 항목 부분 검증
	item1 := got.Output2[1]
	assert.Equal(t, "NQ", item1.ClasCd)

	sttlPrice1, _ := decimal.NewFromString("18450.25")
	assert.True(t, sttlPrice1.Equal(item1.SttlPrice))
}

func TestClient_SearchContractDetail_InvalidJSON(t *testing.T) {
	client, transport := newTestClient(t)
	transport.RegisterResponder(
		http.MethodGet,
		"=~/uapi/overseas-futureoption/v1/quotations/search-contract-detail",
		httpmock.NewStringResponder(http.StatusOK, `{"rt_cd":"0","output2":"not-an-array"}`),
	)
	params := overseasfutures.SearchContractDetailParams{Codes: []string{"ESM24"}}
	_, err := client.SearchContractDetail(context.Background(), params)
	require.Error(t, err)
}

// ─── EP8: InquireAskingPrice ─────────────────────────────────────────────────

func TestClient_InquireAskingPrice(t *testing.T) {
	client, transport := newTestClient(t)
	transport.RegisterResponder(
		http.MethodGet,
		"=~/uapi/overseas-futureoption/v1/quotations/inquire-asking-price",
		httpmock.NewStringResponder(http.StatusOK, loadFixtureString(t, "inquire_asking_price_success.json")),
	)

	got, err := client.InquireAskingPrice(context.Background(), "ESM24")
	require.NoError(t, err)
	require.NotNil(t, got)

	// output1 assertions
	openPrice, _ := decimal.NewFromString("5230.50")
	assert.True(t, openPrice.Equal(got.Output1.OpenPrice))

	highPrice, _ := decimal.NewFromString("5245.75")
	assert.True(t, highPrice.Equal(got.Output1.HighPrice))

	// lowp_rice 오타 필드명 검증
	lowpRice, _ := decimal.NewFromString("5198.25")
	assert.True(t, lowpRice.Equal(got.Output1.LowpRice))

	lastPrice, _ := decimal.NewFromString("5204.34")
	assert.True(t, lastPrice.Equal(got.Output1.LastPrice))

	prevPrice, _ := decimal.NewFromString("5272.75")
	assert.True(t, prevPrice.Equal(got.Output1.PrevPrice))

	assert.Equal(t, int64(1345678), got.Output1.Vol)
	assert.InDelta(t, -1.30, float64(got.Output1.PrevDiffRate), 0.001)
	assert.Equal(t, "20240402", got.Output1.QuotDate)
	assert.Equal(t, "160023", got.Output1.QuotTime)

	// output2 배열 assertions
	require.Len(t, got.Output2, 5)

	q0 := got.Output2[0]
	assert.Equal(t, int64(45), q0.BidQntt)
	assert.Equal(t, "12", q0.BidNum)

	bidPrice0, _ := decimal.NewFromString("5204.00")
	assert.True(t, bidPrice0.Equal(q0.BidPrice))

	assert.Equal(t, int64(38), q0.AskQntt)
	assert.Equal(t, "9", q0.AskNum)

	askPrice0, _ := decimal.NewFromString("5204.25")
	assert.True(t, askPrice0.Equal(q0.AskPrice))
}

func TestClient_InquireAskingPrice_InvalidJSON(t *testing.T) {
	client, transport := newTestClient(t)
	transport.RegisterResponder(
		http.MethodGet,
		"=~/uapi/overseas-futureoption/v1/quotations/inquire-asking-price",
		httpmock.NewStringResponder(http.StatusOK, `{"rt_cd":"0","output1":"bad"}`),
	)
	_, err := client.InquireAskingPrice(context.Background(), "ESM24")
	require.Error(t, err)
}

func TestClient_StockDetail_InvalidJSON(t *testing.T) {
	client, transport := newTestClient(t)
	transport.RegisterResponder(
		http.MethodGet,
		"=~/uapi/overseas-futureoption/v1/quotations/stock-detail",
		httpmock.NewStringResponder(http.StatusOK, `{"rt_cd":"0","output1":"bad"}`),
	)
	_, err := client.StockDetail(context.Background(), "CNHU24")
	require.Error(t, err)
}
