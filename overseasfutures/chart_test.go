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

// ─── EP2: InquireTimeFuturechartprice ────────────────────────────────────────

func TestClient_InquireTimeFuturechartprice(t *testing.T) {
	client, transport := newTestClient(t)
	transport.RegisterResponder(
		http.MethodGet,
		"=~/uapi/overseas-futureoption/v1/quotations/inquire-time-futurechartprice",
		httpmock.NewStringResponder(http.StatusOK, loadFixtureString(t, "inquire_time_futurechartprice_success.json")),
	)

	params := overseasfutures.InquireTimeFuturechartpriceParams{
		SrsCd:         "CNHU24",
		ExchCd:        "CME",
		StartDateTime: "",
		CloseDateTime: "20231214",
		QryTp:         "Q",
		QryCnt:        "120",
		QryGap:        "5",
		IndexKey:      "",
	}
	got, err := client.InquireTimeFuturechartprice(context.Background(), params)
	require.NoError(t, err)
	require.NotNil(t, got)

	// output2 (메타, 단일) assertions — 역전 패턴
	assert.Equal(t, "5", got.Output2.RetCnt)
	assert.Equal(t, "5", got.Output2.LastNCnt)
	assert.Equal(t, "20231214160000ESM24CME", got.Output2.IndexKey)

	// output1 (분봉 배열) assertions
	require.Len(t, got.Output1, 5)

	candle0 := got.Output1[0]
	assert.Equal(t, "20231214", candle0.DataDate)
	assert.Equal(t, "160000", candle0.DataTime)

	openPrice0, _ := decimal.NewFromString("4720.50")
	assert.True(t, openPrice0.Equal(candle0.OpenPrice))

	highPrice0, _ := decimal.NewFromString("4725.00")
	assert.True(t, highPrice0.Equal(candle0.HighPrice))

	lowPrice0, _ := decimal.NewFromString("4718.25")
	assert.True(t, lowPrice0.Equal(candle0.LowPrice))

	lastPrice0, _ := decimal.NewFromString("4723.75")
	assert.True(t, lastPrice0.Equal(candle0.LastPrice))

	assert.Equal(t, int64(145), candle0.LastQntt)
	assert.Equal(t, int64(125430), candle0.Vol)
	assert.Equal(t, "2", candle0.PrevDiffFlag)
	assert.InDelta(t, 0.27, float64(candle0.PrevDiffRate), 0.001)

	// 마지막 항목 확인
	candle4 := got.Output1[4]
	assert.Equal(t, "154000", candle4.DataTime)
}

// ─── EP4: MonthlyCcnl ────────────────────────────────────────────────────────

func TestClient_MonthlyCcnl(t *testing.T) {
	client, transport := newTestClient(t)
	transport.RegisterResponder(
		http.MethodGet,
		"=~/uapi/overseas-futureoption/v1/quotations/monthly-ccnl",
		httpmock.NewStringResponder(http.StatusOK, loadFixtureString(t, "monthly_ccnl_success.json")),
	)

	params := overseasfutures.CcnlParams{
		SrsCd:         "6AM24",
		ExchCd:        "CME",
		CloseDateTime: "20240402",
		QryTp:         "Q",
		QryCnt:        "40",
	}
	got, err := client.MonthlyCcnl(context.Background(), params)
	require.NoError(t, err)
	require.NotNil(t, got)

	// output1 assertions
	assert.Equal(t, "5", got.Output1.TretCnt)
	assert.Equal(t, "0", got.Output1.LastNCnt)
	assert.Equal(t, "20240101000000ESM24CME", got.Output1.IndexKey)

	// output2 배열 assertions
	require.Len(t, got.Output2, 5)

	item0 := got.Output2[0]
	assert.Equal(t, "20240401", item0.DataDate)
	assert.Equal(t, "000000", item0.DataTime)

	openPrice0, _ := decimal.NewFromString("5265.00")
	assert.True(t, openPrice0.Equal(item0.OpenPrice))

	assert.Equal(t, int64(0), item0.LastQntt)
	assert.Equal(t, int64(1523456), item0.Vol)
	assert.Equal(t, "2", item0.PrevDiffFlag)
	assert.InDelta(t, 0.81, float64(item0.PrevDiffRate), 0.001)

	// 마지막 항목
	item4 := got.Output2[4]
	assert.Equal(t, "20231201", item4.DataDate)
}

// ─── EP5: DailyCcnl ──────────────────────────────────────────────────────────

func TestClient_DailyCcnl(t *testing.T) {
	client, transport := newTestClient(t)
	transport.RegisterResponder(
		http.MethodGet,
		"=~/uapi/overseas-futureoption/v1/quotations/daily-ccnl",
		httpmock.NewStringResponder(http.StatusOK, loadFixtureString(t, "daily_ccnl_success.json")),
	)

	params := overseasfutures.CcnlParams{
		SrsCd:         "6AM24",
		ExchCd:        "CME",
		CloseDateTime: "20240402",
		QryTp:         "Q",
		QryCnt:        "40",
	}
	got, err := client.DailyCcnl(context.Background(), params)
	require.NoError(t, err)
	require.NotNil(t, got)

	// output1 assertions (tret_cnt)
	assert.Equal(t, "5", got.Output1.TretCnt)
	assert.Equal(t, "0", got.Output1.LastNCnt)
	assert.Equal(t, "20240402000000ESM24CME", got.Output1.IndexKey)

	// output2 배열 assertions
	require.Len(t, got.Output2, 5)

	item0 := got.Output2[0]
	assert.Equal(t, "20240402", item0.DataDate)

	lastPrice0, _ := decimal.NewFromString("5204.34")
	assert.True(t, lastPrice0.Equal(item0.LastPrice))

	assert.Equal(t, "5", item0.PrevDiffFlag)
	assert.InDelta(t, -1.30, float64(item0.PrevDiffRate), 0.001)
}

// ─── EP6: WeeklyCcnl ─────────────────────────────────────────────────────────

func TestClient_WeeklyCcnl(t *testing.T) {
	client, transport := newTestClient(t)
	transport.RegisterResponder(
		http.MethodGet,
		"=~/uapi/overseas-futureoption/v1/quotations/weekly-ccnl",
		httpmock.NewStringResponder(http.StatusOK, loadFixtureString(t, "weekly_ccnl_success.json")),
	)

	params := overseasfutures.CcnlParams{
		SrsCd:         "6AM24",
		ExchCd:        "CME",
		CloseDateTime: "20240402",
		QryTp:         "Q",
		QryCnt:        "40",
	}
	got, err := client.WeeklyCcnl(context.Background(), params)
	require.NoError(t, err)
	require.NotNil(t, got)

	// output1 assertions — ret_cnt anomaly 검증
	assert.Equal(t, "4", got.Output1.RetCnt)
	assert.Equal(t, "0", got.Output1.LastNCnt)
	assert.Equal(t, "20240402000000ESM24CME", got.Output1.IndexKey)

	// output2 배열 assertions
	require.Len(t, got.Output2, 4)

	item0 := got.Output2[0]
	assert.Equal(t, "20240401", item0.DataDate)

	lastPrice0, _ := decimal.NewFromString("5204.34")
	assert.True(t, lastPrice0.Equal(item0.LastPrice))

	assert.Equal(t, "5", item0.PrevDiffFlag)
	assert.InDelta(t, -0.59, float64(item0.PrevDiffRate), 0.001)
}

func TestClient_InquireTimeFuturechartprice_InvalidJSON(t *testing.T) {
	client, transport := newTestClient(t)
	transport.RegisterResponder(
		http.MethodGet,
		"=~/uapi/overseas-futureoption/v1/quotations/inquire-time-futurechartprice",
		httpmock.NewStringResponder(http.StatusOK, `{"rt_cd":"0","output2":"bad"}`),
	)
	params := overseasfutures.InquireTimeFuturechartpriceParams{SrsCd: "X", ExchCd: "CME"}
	_, err := client.InquireTimeFuturechartprice(context.Background(), params)
	require.Error(t, err)
}

func TestClient_MonthlyCcnl_InvalidJSON(t *testing.T) {
	client, transport := newTestClient(t)
	transport.RegisterResponder(
		http.MethodGet,
		"=~/uapi/overseas-futureoption/v1/quotations/monthly-ccnl",
		httpmock.NewStringResponder(http.StatusOK, `{"rt_cd":"0","output1":"bad"}`),
	)
	_, err := client.MonthlyCcnl(context.Background(), overseasfutures.CcnlParams{})
	require.Error(t, err)
}

func TestClient_DailyCcnl_InvalidJSON(t *testing.T) {
	client, transport := newTestClient(t)
	transport.RegisterResponder(
		http.MethodGet,
		"=~/uapi/overseas-futureoption/v1/quotations/daily-ccnl",
		httpmock.NewStringResponder(http.StatusOK, `{"rt_cd":"0","output1":"bad"}`),
	)
	_, err := client.DailyCcnl(context.Background(), overseasfutures.CcnlParams{})
	require.Error(t, err)
}

func TestClient_WeeklyCcnl_InvalidJSON(t *testing.T) {
	client, transport := newTestClient(t)
	transport.RegisterResponder(
		http.MethodGet,
		"=~/uapi/overseas-futureoption/v1/quotations/weekly-ccnl",
		httpmock.NewStringResponder(http.StatusOK, `{"rt_cd":"0","output1":"bad"}`),
	)
	_, err := client.WeeklyCcnl(context.Background(), overseasfutures.CcnlParams{})
	require.Error(t, err)
}

// ─── EP7: TickCcnl ───────────────────────────────────────────────────────────

func TestClient_TickCcnl(t *testing.T) {
	client, transport := newTestClient(t)
	transport.RegisterResponder(
		http.MethodGet,
		"=~/uapi/overseas-futureoption/v1/quotations/tick-ccnl",
		httpmock.NewStringResponder(http.StatusOK, loadFixtureString(t, "tick_ccnl_success.json")),
	)

	params := overseasfutures.CcnlParams{
		SrsCd:         "6AM24",
		ExchCd:        "CME",
		CloseDateTime: "20240402",
		QryTp:         "Q",
		QryCnt:        "40",
		QryGap:        "1",
	}
	got, err := client.TickCcnl(context.Background(), params)
	require.NoError(t, err)
	require.NotNil(t, got)

	// output1 assertions (tret_cnt)
	assert.Equal(t, "5", got.Output1.TretCnt)
	assert.Equal(t, "5", got.Output1.LastNCnt)
	assert.Equal(t, "20240402160023ESM24CME", got.Output1.IndexKey)

	// output2 배열 assertions
	require.Len(t, got.Output2, 5)

	tick0 := got.Output2[0]
	assert.Equal(t, "20240402", tick0.DataDate)
	assert.Equal(t, "160023", tick0.DataTime)

	lastPrice0, _ := decimal.NewFromString("5204.34")
	assert.True(t, lastPrice0.Equal(tick0.LastPrice))

	assert.Equal(t, int64(5), tick0.LastQntt)
	assert.Equal(t, int64(1345678), tick0.Vol)

	// 마지막 틱
	tick4 := got.Output2[4]
	assert.Equal(t, "160003", tick4.DataTime)
}

func TestClient_TickCcnl_InvalidJSON(t *testing.T) {
	client, transport := newTestClient(t)
	transport.RegisterResponder(
		http.MethodGet,
		"=~/uapi/overseas-futureoption/v1/quotations/tick-ccnl",
		httpmock.NewStringResponder(http.StatusOK, `{"rt_cd":"0","output1":"bad"}`),
	)
	_, err := client.TickCcnl(context.Background(), overseasfutures.CcnlParams{})
	require.Error(t, err)
}
