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

// ─── EP10: MarketTime ────────────────────────────────────────────────────────

func TestClient_MarketTime(t *testing.T) {
	client, transport := newTestClient(t)
	transport.RegisterResponder(
		http.MethodGet,
		"=~/uapi/overseas-futureoption/v1/quotations/market-time",
		httpmock.NewStringResponder(http.StatusOK, loadFixtureString(t, "market_time_success.json")),
	)

	params := overseasfutures.MarketTimeParams{
		FmPdgrCd: "",
		FmClasCd: "003",
		FmExcgCd: "CME",
		OptYn:    "Y",
	}
	got, err := client.MarketTime(context.Background(), params)
	require.NoError(t, err)
	require.NotNil(t, got)

	// output 배열 assertions (output 키 직접 사용 — 특이 패턴)
	require.Len(t, got.Output, 3)

	item0 := got.Output[0]
	assert.Equal(t, "ES", item0.FmPdgrCd)
	assert.Equal(t, "E-mini S&P500", item0.FmPdgrName)
	assert.Equal(t, "CME", item0.FmExcgCd)
	assert.Equal(t, "CME", item0.FmExcgName)
	assert.Equal(t, "옵션", item0.FuopDvsnName)
	assert.Equal(t, "003", item0.FmClasCd)
	assert.Equal(t, "지수", item0.FmClasName)
	assert.Equal(t, "180000", item0.AmMkmnStrtTmd)
	assert.Equal(t, "235959", item0.AmMkmnEndTmd)
	assert.Equal(t, "000000", item0.PmMkmnStrtTmd)
	assert.Equal(t, "170000", item0.PmMkmnEndTmd)
	assert.Equal(t, "000000", item0.MkmnNxdyStrtTmd)
	assert.Equal(t, "170000", item0.MkmnNxdyEndTmd)
	assert.Equal(t, "083000", item0.BaseMketStrtTmd)
	assert.Equal(t, "150000", item0.BaseMketEndTmd)

	// 두 번째 항목 (NQ)
	item1 := got.Output[1]
	assert.Equal(t, "NQ", item1.FmPdgrCd)
	assert.Equal(t, "E-mini NASDAQ-100", item1.FmPdgrName)

	// 세 번째 항목 (6E — 통화)
	item2 := got.Output[2]
	assert.Equal(t, "6E", item2.FmPdgrCd)
	assert.Equal(t, "001", item2.FmClasCd)
	assert.Equal(t, "통화", item2.FmClasName)
	assert.Equal(t, "070000", item2.BaseMketStrtTmd)
	assert.Equal(t, "160000", item2.BaseMketEndTmd)
}

func TestClient_MarketTime_InvalidJSON(t *testing.T) {
	client, transport := newTestClient(t)
	transport.RegisterResponder(
		http.MethodGet,
		"=~/uapi/overseas-futureoption/v1/quotations/market-time",
		httpmock.NewStringResponder(http.StatusOK, `{"rt_cd":"0","output":"not-an-array"}`),
	)
	params := overseasfutures.MarketTimeParams{FmExcgCd: "CME", OptYn: "Y"}
	_, err := client.MarketTime(context.Background(), params)
	require.Error(t, err)
}
