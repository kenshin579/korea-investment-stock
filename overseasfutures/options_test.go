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

// ─── EP1: InquireTimeOptchartprice ───────────────────────────────────────────

func TestClient_InquireTimeOptchartprice(t *testing.T) {
	client, transport := newTestClient(t)
	transport.RegisterResponder(
		http.MethodGet,
		"=~/uapi/overseas-futureoption/v1/quotations/inquire-time-optchartprice",
		httpmock.NewStringResponder(http.StatusOK, loadFixtureString(t, "inquire_time_optchartprice_success.json")),
	)

	params := overseasfutures.InquireTimeOptchartpriceParams{
		SrsCd:  "OESU24 C5500",
		ExchCd: "CME",
		QryTp:  "Q",
		QryCnt: "120",
		QryGap: "5",
	}
	got, err := client.InquireTimeOptchartprice(context.Background(), params)
	require.NoError(t, err)
	require.NotNil(t, got)

	// output2 (메타, 단일) assertions — 역전 패턴
	assert.Equal(t, "5", got.Output2.RetCnt)
	assert.Equal(t, "5", got.Output2.LastNCnt)
	assert.Equal(t, "20240906160000OESU24CME", got.Output2.IndexKey)

	// output1 (분봉 배열) assertions
	require.Len(t, got.Output1, 5)

	candle0 := got.Output1[0]
	assert.Equal(t, "20240906", candle0.DataDate)
	assert.Equal(t, "160000", candle0.DataTime)

	openPrice0, _ := decimal.NewFromString("7525")
	assert.True(t, openPrice0.Equal(candle0.OpenPrice))

	highPrice0, _ := decimal.NewFromString("7550")
	assert.True(t, highPrice0.Equal(candle0.HighPrice))

	lowPrice0, _ := decimal.NewFromString("7510")
	assert.True(t, lowPrice0.Equal(candle0.LowPrice))

	lastPrice0, _ := decimal.NewFromString("7540")
	assert.True(t, lastPrice0.Equal(candle0.LastPrice))

	assert.Equal(t, int64(120), candle0.LastQntt)
	assert.Equal(t, int64(45230), candle0.Vol)
	assert.Equal(t, "2", candle0.PrevDiffFlag)
	assert.InDelta(t, 0.20, float64(candle0.PrevDiffRate), 0.001)

	// 마지막 항목 확인
	candle4 := got.Output1[4]
	assert.Equal(t, "154000", candle4.DataTime)
}

func TestClient_InquireTimeOptchartprice_InvalidJSON(t *testing.T) {
	client, transport := newTestClient(t)
	transport.RegisterResponder(
		http.MethodGet,
		"=~/uapi/overseas-futureoption/v1/quotations/inquire-time-optchartprice",
		httpmock.NewStringResponder(http.StatusOK, `{"rt_cd":"0","output2":"bad"}`),
	)
	params := overseasfutures.InquireTimeOptchartpriceParams{SrsCd: "X", ExchCd: "CME"}
	_, err := client.InquireTimeOptchartprice(context.Background(), params)
	require.Error(t, err)
}

// ─── EP2: SearchOptDetail ────────────────────────────────────────────────────

func TestClient_SearchOptDetail(t *testing.T) {
	client, transport := newTestClient(t)
	transport.RegisterResponder(
		http.MethodGet,
		"=~/uapi/overseas-futureoption/v1/quotations/search-opt-detail",
		httpmock.NewStringResponder(http.StatusOK, loadFixtureString(t, "search_opt_detail_success.json")),
	)

	params := overseasfutures.SearchOptDetailParams{
		Codes: []string{"OESU24 C5500", "OESU24 P5500"},
	}
	got, err := client.SearchOptDetail(context.Background(), params)
	require.NoError(t, err)
	require.NotNil(t, got)

	// output2 배열 assertions
	require.Len(t, got.Output2, 2)

	item0 := got.Output2[0]
	assert.Equal(t, "CME", item0.ExchCd)
	assert.Equal(t, "C", item0.ClasCd) // 옵션: length=1 (선물: length=3)
	assert.Equal(t, "USD", item0.CrcCd)

	sttlPrice0, _ := decimal.NewFromString("7525")
	assert.True(t, sttlPrice0.Equal(item0.SttlPrice))

	assert.Equal(t, "20240920", item0.SttlDate)

	trstMgn0, _ := decimal.NewFromString("3500.00")
	assert.True(t, trstMgn0.Equal(item0.TrstMgn))

	tickSz0, _ := decimal.NewFromString("0.05")
	assert.True(t, tickSz0.Equal(item0.TickSz))

	tickVal0, _ := decimal.NewFromString("12.50")
	assert.True(t, tickVal0.Equal(item0.TickVal))

	assert.Equal(t, "20240906", item0.MrktOpenDate)
	assert.Equal(t, "170000", item0.MrktOpenTime)
	assert.Equal(t, "20240920", item0.ExprDate)
	assert.Equal(t, "12", item0.RemnCnt)
	assert.Equal(t, "1", item0.StatTp)
	assert.Equal(t, "CASH", item0.StlTp)

	// 두 번째 항목 (Put)
	item1 := got.Output2[1]
	assert.Equal(t, "P", item1.ClasCd)

	sttlPrice1, _ := decimal.NewFromString("7475")
	assert.True(t, sttlPrice1.Equal(item1.SttlPrice))
}

func TestClient_SearchOptDetail_InvalidJSON(t *testing.T) {
	client, transport := newTestClient(t)
	transport.RegisterResponder(
		http.MethodGet,
		"=~/uapi/overseas-futureoption/v1/quotations/search-opt-detail",
		httpmock.NewStringResponder(http.StatusOK, `{"rt_cd":"0","output2":"not-an-array"}`),
	)
	params := overseasfutures.SearchOptDetailParams{Codes: []string{"OESU24 C5500"}}
	_, err := client.SearchOptDetail(context.Background(), params)
	require.Error(t, err)
}

// ─── EP3: OptMonthlyCcnl ─────────────────────────────────────────────────────

func TestClient_OptMonthlyCcnl(t *testing.T) {
	client, transport := newTestClient(t)
	transport.RegisterResponder(
		http.MethodGet,
		"=~/uapi/overseas-futureoption/v1/quotations/opt-monthly-ccnl",
		httpmock.NewStringResponder(http.StatusOK, loadFixtureString(t, "opt_monthly_ccnl_success.json")),
	)

	params := overseasfutures.CcnlParams{
		SrsCd:  "OESU24 C5500",
		ExchCd: "CME",
		QryTp:  "Q",
		QryCnt: "120",
	}
	got, err := client.OptMonthlyCcnl(context.Background(), params)
	require.NoError(t, err)
	require.NotNil(t, got)

	// output1 assertions — ret_cnt (옵션 공통)
	assert.Equal(t, "5", got.Output1.RetCnt)
	assert.Equal(t, "0", got.Output1.LastNCnt)
	assert.Equal(t, "20240101000000OESU24CME", got.Output1.IndexKey)

	// output2 배열 assertions
	require.Len(t, got.Output2, 5)

	item0 := got.Output2[0]
	assert.Equal(t, "20240901", item0.DataDate)
	assert.Equal(t, "000000", item0.DataTime)

	openPrice0, _ := decimal.NewFromString("7200")
	assert.True(t, openPrice0.Equal(item0.OpenPrice))

	lastPrice0, _ := decimal.NewFromString("7525")
	assert.True(t, lastPrice0.Equal(item0.LastPrice))

	assert.Equal(t, int64(0), item0.LastQntt)
	assert.Equal(t, int64(523456), item0.Vol)
	assert.Equal(t, "2", item0.PrevDiffFlag)
	assert.InDelta(t, 4.52, float64(item0.PrevDiffRate), 0.001)

	// 마지막 항목
	item4 := got.Output2[4]
	assert.Equal(t, "20240501", item4.DataDate)
}

func TestClient_OptMonthlyCcnl_InvalidJSON(t *testing.T) {
	client, transport := newTestClient(t)
	transport.RegisterResponder(
		http.MethodGet,
		"=~/uapi/overseas-futureoption/v1/quotations/opt-monthly-ccnl",
		httpmock.NewStringResponder(http.StatusOK, `{"rt_cd":"0","output1":"bad"}`),
	)
	_, err := client.OptMonthlyCcnl(context.Background(), overseasfutures.CcnlParams{})
	require.Error(t, err)
}

// ─── EP4: OptDailyCcnl ───────────────────────────────────────────────────────

func TestClient_OptDailyCcnl(t *testing.T) {
	client, transport := newTestClient(t)
	transport.RegisterResponder(
		http.MethodGet,
		"=~/uapi/overseas-futureoption/v1/quotations/opt-daily-ccnl",
		httpmock.NewStringResponder(http.StatusOK, loadFixtureString(t, "opt_daily_ccnl_success.json")),
	)

	params := overseasfutures.CcnlParams{
		SrsCd:  "OESU24 C5500",
		ExchCd: "CME",
		QryTp:  "Q",
		QryCnt: "119",
	}
	got, err := client.OptDailyCcnl(context.Background(), params)
	require.NoError(t, err)
	require.NotNil(t, got)

	// output1 assertions — ret_cnt
	assert.Equal(t, "5", got.Output1.RetCnt)
	assert.Equal(t, "0", got.Output1.LastNCnt)
	assert.Equal(t, "20240906000000OESU24CME", got.Output1.IndexKey)

	// output2 배열 assertions
	require.Len(t, got.Output2, 5)

	item0 := got.Output2[0]
	assert.Equal(t, "20240906", item0.DataDate)

	lastPrice0, _ := decimal.NewFromString("7525")
	assert.True(t, lastPrice0.Equal(item0.LastPrice))

	assert.Equal(t, "2", item0.PrevDiffFlag)
	assert.InDelta(t, 0.60, float64(item0.PrevDiffRate), 0.001)

	// 마지막 항목
	item4 := got.Output2[4]
	assert.Equal(t, "20240902", item4.DataDate)
}

func TestClient_OptDailyCcnl_InvalidJSON(t *testing.T) {
	client, transport := newTestClient(t)
	transport.RegisterResponder(
		http.MethodGet,
		"=~/uapi/overseas-futureoption/v1/quotations/opt-daily-ccnl",
		httpmock.NewStringResponder(http.StatusOK, `{"rt_cd":"0","output1":"bad"}`),
	)
	_, err := client.OptDailyCcnl(context.Background(), overseasfutures.CcnlParams{})
	require.Error(t, err)
}

// ─── EP5: OptWeeklyCcnl ──────────────────────────────────────────────────────

func TestClient_OptWeeklyCcnl(t *testing.T) {
	client, transport := newTestClient(t)
	transport.RegisterResponder(
		http.MethodGet,
		"=~/uapi/overseas-futureoption/v1/quotations/opt-weekly-ccnl",
		httpmock.NewStringResponder(http.StatusOK, loadFixtureString(t, "opt_weekly_ccnl_success.json")),
	)

	params := overseasfutures.CcnlParams{
		SrsCd:  "OESU24 C5500",
		ExchCd: "CME",
		QryTp:  "Q",
		QryCnt: "40",
	}
	got, err := client.OptWeeklyCcnl(context.Background(), params)
	require.NoError(t, err)
	require.NotNil(t, got)

	// output1 assertions — ret_cnt
	assert.Equal(t, "5", got.Output1.RetCnt)
	assert.Equal(t, "0", got.Output1.LastNCnt)
	assert.Equal(t, "20240902000000OESU24CME", got.Output1.IndexKey)

	// output2 배열 assertions
	require.Len(t, got.Output2, 5)

	item0 := got.Output2[0]
	assert.Equal(t, "20240902", item0.DataDate)

	lastPrice0, _ := decimal.NewFromString("7525")
	assert.True(t, lastPrice0.Equal(item0.LastPrice))

	assert.Equal(t, "2", item0.PrevDiffFlag)
	assert.InDelta(t, 2.38, float64(item0.PrevDiffRate), 0.001)

	// 마지막 항목
	item4 := got.Output2[4]
	assert.Equal(t, "20240805", item4.DataDate)
}

func TestClient_OptWeeklyCcnl_InvalidJSON(t *testing.T) {
	client, transport := newTestClient(t)
	transport.RegisterResponder(
		http.MethodGet,
		"=~/uapi/overseas-futureoption/v1/quotations/opt-weekly-ccnl",
		httpmock.NewStringResponder(http.StatusOK, `{"rt_cd":"0","output1":"bad"}`),
	)
	_, err := client.OptWeeklyCcnl(context.Background(), overseasfutures.CcnlParams{})
	require.Error(t, err)
}

// ─── EP6: OptTickCcnl ────────────────────────────────────────────────────────

func TestClient_OptTickCcnl(t *testing.T) {
	client, transport := newTestClient(t)
	transport.RegisterResponder(
		http.MethodGet,
		"=~/uapi/overseas-futureoption/v1/quotations/opt-tick-ccnl",
		httpmock.NewStringResponder(http.StatusOK, loadFixtureString(t, "opt_tick_ccnl_success.json")),
	)

	params := overseasfutures.CcnlParams{
		SrsCd:  "OESU24 C5500",
		ExchCd: "CME",
		QryTp:  "Q",
		QryCnt: "40",
		QryGap: "1",
	}
	got, err := client.OptTickCcnl(context.Background(), params)
	require.NoError(t, err)
	require.NotNil(t, got)

	// output1 assertions — ret_cnt
	assert.Equal(t, "5", got.Output1.RetCnt)
	assert.Equal(t, "5", got.Output1.LastNCnt)
	assert.Equal(t, "20240906160023OESU24CME", got.Output1.IndexKey)

	// output2 배열 assertions
	require.Len(t, got.Output2, 5)

	tick0 := got.Output2[0]
	assert.Equal(t, "20240906", tick0.DataDate)
	assert.Equal(t, "160023", tick0.DataTime)

	lastPrice0, _ := decimal.NewFromString("7540")
	assert.True(t, lastPrice0.Equal(tick0.LastPrice))

	assert.Equal(t, int64(5), tick0.LastQntt)
	assert.Equal(t, int64(45230), tick0.Vol)

	// 마지막 틱
	tick4 := got.Output2[4]
	assert.Equal(t, "155958", tick4.DataTime)
}

func TestClient_OptTickCcnl_InvalidJSON(t *testing.T) {
	client, transport := newTestClient(t)
	transport.RegisterResponder(
		http.MethodGet,
		"=~/uapi/overseas-futureoption/v1/quotations/opt-tick-ccnl",
		httpmock.NewStringResponder(http.StatusOK, `{"rt_cd":"0","output1":"bad"}`),
	)
	_, err := client.OptTickCcnl(context.Background(), overseasfutures.CcnlParams{})
	require.Error(t, err)
}

// ─── EP7: OptAskingPrice ─────────────────────────────────────────────────────

func TestClient_OptAskingPrice(t *testing.T) {
	client, transport := newTestClient(t)
	transport.RegisterResponder(
		http.MethodGet,
		"=~/uapi/overseas-futureoption/v1/quotations/opt-asking-price",
		httpmock.NewStringResponder(http.StatusOK, loadFixtureString(t, "opt_asking_price_success.json")),
	)

	got, err := client.OptAskingPrice(context.Background(), "OESM24 C5340")
	require.NoError(t, err)
	require.NotNil(t, got)

	// output1 assertions
	openPrice, _ := decimal.NewFromString("7480")
	assert.True(t, openPrice.Equal(got.Output1.OpenPrice))

	highPrice, _ := decimal.NewFromString("7560")
	assert.True(t, highPrice.Equal(got.Output1.HighPrice))

	// lowp_rice 오타 필드명 검증
	lowpRice, _ := decimal.NewFromString("7465")
	assert.True(t, lowpRice.Equal(got.Output1.LowpRice))

	lastPrice, _ := decimal.NewFromString("7525")
	assert.True(t, lastPrice.Equal(got.Output1.LastPrice))

	// sttl_price — 옵션 고유 필드 (선물 prev_price 와 다름)
	sttlPrice, _ := decimal.NewFromString("7500")
	assert.True(t, sttlPrice.Equal(got.Output1.SttlPrice))

	assert.Equal(t, int64(45230), got.Output1.Vol)
	assert.InDelta(t, 0.60, float64(got.Output1.PrevDiffRate), 0.001)
	assert.Equal(t, "20240906", got.Output1.QuotDate)
	assert.Equal(t, "160023", got.Output1.QuotTime)

	// output2 배열 assertions (1~5호가)
	require.Len(t, got.Output2, 5)

	q0 := got.Output2[0]
	assert.Equal(t, int64(25), q0.BidQntt)
	assert.Equal(t, "8", q0.BidNum)

	bidPrice0, _ := decimal.NewFromString("7520")
	assert.True(t, bidPrice0.Equal(q0.BidPrice))

	assert.Equal(t, int64(30), q0.AskQntt)
	assert.Equal(t, "10", q0.AskNum)

	askPrice0, _ := decimal.NewFromString("7525")
	assert.True(t, askPrice0.Equal(q0.AskPrice))
}

func TestClient_OptAskingPrice_InvalidJSON(t *testing.T) {
	client, transport := newTestClient(t)
	transport.RegisterResponder(
		http.MethodGet,
		"=~/uapi/overseas-futureoption/v1/quotations/opt-asking-price",
		httpmock.NewStringResponder(http.StatusOK, `{"rt_cd":"0","output1":"bad"}`),
	)
	_, err := client.OptAskingPrice(context.Background(), "OESM24 C5340")
	require.Error(t, err)
}

// ─── EP8: OptDetail ──────────────────────────────────────────────────────────

func TestClient_OptDetail(t *testing.T) {
	client, transport := newTestClient(t)
	transport.RegisterResponder(
		http.MethodGet,
		"=~/uapi/overseas-futureoption/v1/quotations/opt-detail",
		httpmock.NewStringResponder(http.StatusOK, loadFixtureString(t, "opt_detail_success.json")),
	)

	got, err := client.OptDetail(context.Background(), "OESU24 C5500")
	require.NoError(t, err)
	require.NotNil(t, got)

	// output1 assertions
	assert.Equal(t, "CME", got.Output1.ExchCd)
	assert.Equal(t, "C", got.Output1.ClasCd)
	assert.Equal(t, "USD", got.Output1.CrcCd)

	// sttl_price — 실제 전일종가 수신 (docs 주의: 정산가 X 전일종가 O)
	sttlPrice, _ := decimal.NewFromString("7480")
	assert.True(t, sttlPrice.Equal(got.Output1.SttlPrice))

	assert.Equal(t, "20240920", got.Output1.SttlDate)

	trstMgn, _ := decimal.NewFromString("3500.00")
	assert.True(t, trstMgn.Equal(got.Output1.TrstMgn))

	assert.Equal(t, "2", got.Output1.DispDigit)

	tickSz, _ := decimal.NewFromString("0.05")
	assert.True(t, tickSz.Equal(got.Output1.TickSz))

	tickVal, _ := decimal.NewFromString("12.50")
	assert.True(t, tickVal.Equal(got.Output1.TickVal))

	assert.Equal(t, "20240906", got.Output1.MrktOpenDate)
	assert.Equal(t, "170000", got.Output1.MrktOpenTime)
	assert.Equal(t, "20240907", got.Output1.MrktCloseDate)
	assert.Equal(t, "160000", got.Output1.MrktCloseTime)
	assert.Equal(t, "20240101", got.Output1.TrdFrDate)
	assert.Equal(t, "20240920", got.Output1.ExprDate)
	assert.Equal(t, "20240918", got.Output1.TrdToDate)
	assert.Equal(t, "12", got.Output1.RemnCnt)
	assert.Equal(t, "1", got.Output1.StatTp)

	ctrtSize, _ := decimal.NewFromString("50.00")
	assert.True(t, ctrtSize.Equal(got.Output1.CtrtSize))

	assert.Equal(t, "CASH", got.Output1.StlTp)
}

func TestClient_OptDetail_InvalidJSON(t *testing.T) {
	client, transport := newTestClient(t)
	transport.RegisterResponder(
		http.MethodGet,
		"=~/uapi/overseas-futureoption/v1/quotations/opt-detail",
		httpmock.NewStringResponder(http.StatusOK, `{"rt_cd":"0","output1":"bad"}`),
	)
	_, err := client.OptDetail(context.Background(), "OESU24 C5500")
	require.Error(t, err)
}

// ─── EP9: OptPrice ───────────────────────────────────────────────────────────

func TestClient_OptPrice(t *testing.T) {
	client, transport := newTestClient(t)
	transport.RegisterResponder(
		http.MethodGet,
		"=~/uapi/overseas-futureoption/v1/quotations/opt-price",
		httpmock.NewStringResponder(http.StatusOK, loadFixtureString(t, "opt_price_success.json")),
	)

	got, err := client.OptPrice(context.Background(), "OESU24 C5500")
	require.NoError(t, err)
	require.NotNil(t, got)

	// output1 assertions
	assert.Equal(t, "20240906", got.Output1.ProcDate)
	assert.Equal(t, "160023", got.Output1.ProcTime)

	openPrice, _ := decimal.NewFromString("7480")
	assert.True(t, openPrice.Equal(got.Output1.OpenPrice))

	highPrice, _ := decimal.NewFromString("7560")
	assert.True(t, highPrice.Equal(got.Output1.HighPrice))

	lowPrice, _ := decimal.NewFromString("7465")
	assert.True(t, lowPrice.Equal(got.Output1.LowPrice))

	lastPrice, _ := decimal.NewFromString("7525")
	assert.True(t, lastPrice.Equal(got.Output1.LastPrice))

	assert.Equal(t, int64(45230), got.Output1.Vol)
	assert.Equal(t, "2", got.Output1.PrevDiffFlag)
	assert.InDelta(t, 0.60, float64(got.Output1.PrevDiffRate), 0.001)

	assert.Equal(t, int64(25), got.Output1.BidQntt)

	bidPrice, _ := decimal.NewFromString("7520")
	assert.True(t, bidPrice.Equal(got.Output1.BidPrice))

	assert.Equal(t, int64(30), got.Output1.AskQntt)

	askPrice, _ := decimal.NewFromString("7525")
	assert.True(t, askPrice.Equal(got.Output1.AskPrice))

	trstMgn, _ := decimal.NewFromString("3500.00")
	assert.True(t, trstMgn.Equal(got.Output1.TrstMgn))

	assert.Equal(t, "CME", got.Output1.ExchCd)
	assert.Equal(t, "USD", got.Output1.CrcCd)
	assert.Equal(t, "20240101", got.Output1.TrdFrDate)
	assert.Equal(t, "20240920", got.Output1.ExprDate)
	assert.Equal(t, "20240918", got.Output1.TrdToDate)
	assert.Equal(t, "12", got.Output1.RemnCnt)
	assert.Equal(t, int64(5), got.Output1.LastQntt)
	assert.Equal(t, int64(355), got.Output1.TotAskQntt)
	assert.Equal(t, int64(305), got.Output1.TotBidQntt)

	tickSize, _ := decimal.NewFromString("0.05")
	assert.True(t, tickSize.Equal(got.Output1.TickSize))

	assert.Equal(t, "20240906", got.Output1.OpenDate)
	assert.Equal(t, "170000", got.Output1.OpenTime)
	assert.Equal(t, "20240907", got.Output1.CloseDate)
	assert.Equal(t, "160000", got.Output1.CloseTime)
	assert.Equal(t, "20240906", got.Output1.Sbsnsdate)

	// sttl_price Optional 검증
	sttlPrice, _ := decimal.NewFromString("7500")
	assert.True(t, sttlPrice.Equal(got.Output1.SttlPrice))
}

func TestClient_OptPrice_InvalidJSON(t *testing.T) {
	client, transport := newTestClient(t)
	transport.RegisterResponder(
		http.MethodGet,
		"=~/uapi/overseas-futureoption/v1/quotations/opt-price",
		httpmock.NewStringResponder(http.StatusOK, `{"rt_cd":"0","output1":"not-an-object"}`),
	)
	_, err := client.OptPrice(context.Background(), "OESU24 C5500")
	require.Error(t, err)
}
