package websocket

import (
	"strings"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --------------------------------------------------------------------------
// HDFSCNT0 — 해외주식 실시간지연체결가 (26 fields)
// --------------------------------------------------------------------------

func TestDecodeOverseasTrade(t *testing.T) {
	raw := loadFixture(t, "hdfscnt0_success.txt")
	f, err := parseFrame(raw)
	require.NoError(t, err)
	require.Equal(t, "HDFSCNT0", f.TrID)

	events, err := decodeOverseasTrade(f)
	require.NoError(t, err)
	require.Len(t, events, 1)

	ev := events[0]
	assert.Equal(t, "DNASAAPL", ev.Symbol)
	assert.Equal(t, "AAPL", ev.SymbolCode)
	assert.Equal(t, "4", ev.Decimals)
	assert.Equal(t, "20260509", ev.LocalDate)
	assert.Equal(t, "093015", ev.LocalTime)
	assert.Equal(t, "223015", ev.KrTime)

	assert.True(t, decimal.NewFromFloat(195.32).Equal(ev.Open))
	assert.True(t, decimal.NewFromFloat(196.50).Equal(ev.High))
	assert.True(t, decimal.NewFromFloat(195.85).Equal(ev.Last))
	assert.Equal(t, "2", ev.PrevDiffSign)
	assert.True(t, decimal.NewFromFloat(0.53).Equal(ev.PrevDiff))
	assert.InDelta(t, 0.27, ev.ChangeRate, 1e-6)
	assert.True(t, decimal.NewFromFloat(195.83).Equal(ev.Bid))
	assert.True(t, decimal.NewFromFloat(195.86).Equal(ev.Ask))
	assert.Equal(t, int64(1500), ev.BidSize)
	assert.Equal(t, int64(1200), ev.AskSize)
	assert.Equal(t, int64(150), ev.TradeVolume)
	assert.Equal(t, int64(25000000), ev.AccumVolume)
	assert.Equal(t, int64(4900000000), ev.AccumValue)
	assert.Equal(t, int64(120000), ev.AskTradeVol) // BIVL
	assert.Equal(t, int64(150000), ev.BidTradeVol) // ASVL
	assert.InDelta(t, 102.5, ev.TradeStrength, 1e-6)
	assert.Equal(t, "1", ev.MarketKind)
	assert.NotEmpty(t, ev.Raw)
}

func TestDecodeOverseasTrade_FieldCountMismatch(t *testing.T) {
	f := frame{Kind: frameKindRealtime, TrID: "HDFSCNT0", Count: 1, Fields: []string{"a", "b"}}
	_, err := decodeOverseasTrade(f)
	require.Error(t, err)
}

func TestDecodeOverseasTrade_BadNumeric(t *testing.T) {
	raw := strings.Replace(loadFixture(t, "hdfscnt0_success.txt"), "195.85", "abc", 1)
	f, err := parseFrame(raw)
	require.NoError(t, err)
	events, err := decodeOverseasTrade(f)
	require.NoError(t, err)
	require.Len(t, events, 1)
	assert.True(t, events[0].Last.IsZero())
}

// --------------------------------------------------------------------------
// HDFSASP0 — 해외주식 실시간호가 (17 fields, 1호가만)
// --------------------------------------------------------------------------

func TestDecodeOverseasAsk(t *testing.T) {
	raw := loadFixture(t, "hdfsasp0_success.txt")
	f, err := parseFrame(raw)
	require.NoError(t, err)
	require.Equal(t, "HDFSASP0", f.TrID)

	events, err := decodeOverseasAsk(f)
	require.NoError(t, err)
	require.Len(t, events, 1)

	ev := events[0]
	assert.Equal(t, "DNASAAPL", ev.Symbol)
	assert.Equal(t, "AAPL", ev.SymbolCode)
	assert.Equal(t, "4", ev.Decimals)
	assert.Equal(t, "20260509", ev.LocalDayDate)
	assert.Equal(t, "093015", ev.LocalTime)
	assert.Equal(t, "20260509", ev.KrDate)
	assert.Equal(t, "223015", ev.KrTime)
	assert.Equal(t, int64(15000), ev.TotalBidSize)
	assert.Equal(t, int64(18000), ev.TotalAskSize)
	assert.Equal(t, int64(500), ev.TotalBidSizeChange)
	assert.Equal(t, int64(-300), ev.TotalAskSizeChange)
	assert.True(t, decimal.NewFromFloat(195.83).Equal(ev.Bid1))
	assert.True(t, decimal.NewFromFloat(195.86).Equal(ev.Ask1))
	assert.Equal(t, int64(1500), ev.Bid1Size)
	assert.Equal(t, int64(1200), ev.Ask1Size)
	assert.Equal(t, int64(200), ev.Bid1SizeChange)
	assert.Equal(t, int64(-100), ev.Ask1SizeChange)
}

func TestDecodeOverseasAsk_FieldCountMismatch(t *testing.T) {
	f := frame{Kind: frameKindRealtime, TrID: "HDFSASP0", Count: 1, Fields: []string{"a"}}
	_, err := decodeOverseasAsk(f)
	require.Error(t, err)
}
