package websocket

import (
	"strings"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --------------------------------------------------------------------------
// H0MFCNT0 — KRX야간선물 실시간종목체결 (49 fields)
// --------------------------------------------------------------------------

func TestDecodeKrxNightFuturesTrade(t *testing.T) {
	raw := loadFixture(t, "h0mfcnt0_success.txt")
	f, err := parseFrame(raw)
	require.NoError(t, err)
	require.Equal(t, "H0MFCNT0", f.TrID)

	events, err := decodeKrxNightFuturesTrade(f)
	require.NoError(t, err)
	require.Len(t, events, 1)

	ev := events[0]
	assert.Equal(t, "101W09000000", ev.Symbol)
	assert.Equal(t, "183500", ev.Time)
	// FUTS_PRDY_VRSS = -25.00
	assert.True(t, decimal.NewFromFloat(-25.00).Equal(ev.PrevDiff))
	// PRDY_VRSS_SIGN = 5
	assert.Equal(t, "5", ev.PrevDiffSign)
	// FUTS_PRPR = 2525.00
	assert.True(t, decimal.NewFromFloat(2525.00).Equal(ev.Price))
	// LAST_CNQN = 100
	assert.Equal(t, int64(100), ev.LastTradeVolume)
	// ACML_VOL = 45678
	assert.Equal(t, int64(45678), ev.AccumVolume)
	// DYNM_MXPR = 2600.00
	assert.True(t, decimal.NewFromFloat(2600.00).Equal(ev.DynamicUpperLimit))
	// DYNM_LLAM = 2450.00
	assert.True(t, decimal.NewFromFloat(2450.00).Equal(ev.DynamicLowerLimit))
	// DYNM_PRC_LIMT_YN = Y
	assert.Equal(t, "Y", ev.DynamicPriceLimitYN)
	assert.NotEmpty(t, ev.Raw)
	assert.Len(t, ev.Raw, 49)
}

func TestDecodeKrxNightFuturesTrade_FieldCountMismatch(t *testing.T) {
	f := frame{Kind: frameKindRealtime, TrID: "H0MFCNT0", Count: 1, Fields: []string{"a", "b", "c"}}
	_, err := decodeKrxNightFuturesTrade(f)
	require.Error(t, err)
}

func TestDecodeKrxNightFuturesTrade_BadNumeric(t *testing.T) {
	raw := strings.Replace(loadFixture(t, "h0mfcnt0_success.txt"), "2525.00", "xyz", 1)
	f, err := parseFrame(raw)
	require.NoError(t, err)

	events, err := decodeKrxNightFuturesTrade(f)
	require.NoError(t, err)
	require.Len(t, events, 1)
	// Price 파싱 실패 시 zero 반환
	assert.True(t, events[0].Price.IsZero())
}

// --------------------------------------------------------------------------
// H0MFASP0 — KRX야간선물 실시간호가 (38 fields)
// --------------------------------------------------------------------------

func TestDecodeKrxNightFuturesAsk(t *testing.T) {
	raw := loadFixture(t, "h0mfasp0_success.txt")
	f, err := parseFrame(raw)
	require.NoError(t, err)
	require.Equal(t, "H0MFASP0", f.TrID)

	events, err := decodeKrxNightFuturesAsk(f)
	require.NoError(t, err)
	require.Len(t, events, 1)

	ev := events[0]
	assert.Equal(t, "101W09000000", ev.Symbol)
	assert.Equal(t, "183500", ev.Time)
	// Ask[0] = 2526.00
	assert.True(t, decimal.NewFromFloat(2526.00).Equal(ev.Ask[0]))
	// Ask[4] = 2530.00
	assert.True(t, decimal.NewFromFloat(2530.00).Equal(ev.Ask[4]))
	// Bid[0] = 2524.00
	assert.True(t, decimal.NewFromFloat(2524.00).Equal(ev.Bid[0]))
	// 5단계 모두 non-zero
	for i := 0; i < 5; i++ {
		assert.False(t, ev.Ask[i].IsZero(), "Ask[%d] zero", i)
		assert.False(t, ev.Bid[i].IsZero(), "Bid[%d] zero", i)
	}
	// TOTAL_ASKP_RSQN = 1250
	assert.Equal(t, int64(1250), ev.TotalAskSize)
	// TOTAL_BIDP_RSQN = 2000
	assert.Equal(t, int64(2000), ev.TotalBidSize)
	assert.NotEmpty(t, ev.Raw)
	assert.Len(t, ev.Raw, 38)
}

// --------------------------------------------------------------------------
// H0EUCNT0 — KRX야간옵션 실시간체결가 (56 fields)
// --------------------------------------------------------------------------

func TestDecodeKrxNightOptionTrade(t *testing.T) {
	raw := loadFixture(t, "h0eucnt0_success.txt")
	f, err := parseFrame(raw)
	require.NoError(t, err)
	require.Equal(t, "H0EUCNT0", f.TrID)

	events, err := decodeKrxNightOptionTrade(f)
	require.NoError(t, err)
	require.Len(t, events, 1)

	ev := events[0]
	assert.Equal(t, "201W09250C", ev.Symbol)
	assert.Equal(t, "183500", ev.Time)
	// OPTN_PRPR = 3.50
	assert.True(t, decimal.NewFromFloat(3.50).Equal(ev.Price))
	// DELTA = 0.45
	assert.InDelta(t, 0.45, ev.Delta, 1e-9)
	// GAMA = 0.02
	assert.InDelta(t, 0.02, ev.Gamma, 1e-9)
	// DYNM_MXPR = 4.50
	assert.True(t, decimal.NewFromFloat(4.50).Equal(ev.DynamicUpperLimit))
	// DYNM_PRC_LIMT_YN = Y (docs anomaly: PRC_LIMT_YN comes before LLAM)
	assert.Equal(t, "Y", ev.DynamicPriceLimitYN)
	// DYNM_LLAM = 3.20
	assert.True(t, decimal.NewFromFloat(3.20).Equal(ev.DynamicLowerLimit))
	assert.NotEmpty(t, ev.Raw)
	assert.Len(t, ev.Raw, 56)
}

func TestDecodeKrxNightOptionTrade_FieldCountMismatch(t *testing.T) {
	f := frame{Kind: frameKindRealtime, TrID: "H0EUCNT0", Count: 1, Fields: []string{"a", "b"}}
	_, err := decodeKrxNightOptionTrade(f)
	require.Error(t, err)
}

// --------------------------------------------------------------------------
// H0EUASP0 — KRX야간옵션 실시간호가 (38 fields)
// --------------------------------------------------------------------------

func TestDecodeKrxNightOptionAsk(t *testing.T) {
	raw := loadFixture(t, "h0euasp0_success.txt")
	f, err := parseFrame(raw)
	require.NoError(t, err)
	require.Equal(t, "H0EUASP0", f.TrID)

	events, err := decodeKrxNightOptionAsk(f)
	require.NoError(t, err)
	require.Len(t, events, 1)

	ev := events[0]
	assert.Equal(t, "201W09250C", ev.Symbol)
	// OPTN_ASKP1 = 3.55
	assert.True(t, decimal.NewFromFloat(3.55).Equal(ev.Ask[0]))
	// OPTN_BIDP1 = 3.45
	assert.True(t, decimal.NewFromFloat(3.45).Equal(ev.Bid[0]))
	for i := 0; i < 5; i++ {
		assert.False(t, ev.Ask[i].IsZero(), "Ask[%d] zero", i)
		assert.False(t, ev.Bid[i].IsZero(), "Bid[%d] zero", i)
	}
	assert.Equal(t, int64(1900), ev.TotalAskSize)
	assert.Equal(t, int64(2250), ev.TotalBidSize)
	assert.Len(t, ev.Raw, 38)
}

// --------------------------------------------------------------------------
// H0EUANC0 — KRX야간옵션 실시간예상체결 (8 fields)
// --------------------------------------------------------------------------

func TestDecodeKrxNightOptionExpectTrade(t *testing.T) {
	raw := loadFixture(t, "h0euanc0_success.txt")
	f, err := parseFrame(raw)
	require.NoError(t, err)
	require.Equal(t, "H0EUANC0", f.TrID)

	events, err := decodeKrxNightOptionExpectTrade(f)
	require.NoError(t, err)
	require.Len(t, events, 1)

	ev := events[0]
	assert.Equal(t, "201W09250C", ev.Symbol)
	assert.Equal(t, "183500", ev.Time)
	// ANTC_CNPR = 3.50
	assert.True(t, decimal.NewFromFloat(3.50).Equal(ev.ExpectPrice))
	// ANTC_CNTG_VRSS = 0.25
	assert.True(t, decimal.NewFromFloat(0.25).Equal(ev.ExpectDiff))
	// ANTC_CNTG_VRSS_SIGN = 2
	assert.Equal(t, "2", ev.ExpectDiffSign)
	// ANTC_CNQN = 500
	assert.Equal(t, int64(500), ev.ExpectQuantity)
	assert.Len(t, ev.Raw, 8)
}

// --------------------------------------------------------------------------
// H0ZFCNT0 — 주식선물 실시간체결가 (49 fields)
// --------------------------------------------------------------------------

func TestDecodeStockFuturesTrade(t *testing.T) {
	raw := loadFixture(t, "h0zfcnt0_success.txt")
	f, err := parseFrame(raw)
	require.NoError(t, err)
	require.Equal(t, "H0ZFCNT0", f.TrID)

	events, err := decodeStockFuturesTrade(f)
	require.NoError(t, err)
	require.Len(t, events, 1)

	ev := events[0]
	assert.Equal(t, "KAK0F", ev.Symbol)
	assert.Equal(t, "133500", ev.Time)
	// STCK_PRPR = 73400
	assert.True(t, decimal.NewFromFloat(73400).Equal(ev.Price))
	// PRDY_VRSS = 1500
	assert.True(t, decimal.NewFromFloat(1500).Equal(ev.PrevDiff))
	// LAST_CNQN = 50
	assert.Equal(t, int64(50), ev.LastTradeVolume)
	// MRKT_BASIS = -50
	assert.True(t, decimal.NewFromFloat(-50).Equal(ev.MarketBasis))
	// DYNM_MXPR = 75000
	assert.True(t, decimal.NewFromFloat(75000).Equal(ev.DynamicUpperLimit))
	// DYNM_LLAM = 71000
	assert.True(t, decimal.NewFromFloat(71000).Equal(ev.DynamicLowerLimit))
	// DYNM_PRC_LIMT_YN = N
	assert.Equal(t, "N", ev.DynamicPriceLimitYN)
	assert.Len(t, ev.Raw, 49)
}

// --------------------------------------------------------------------------
// H0ZFASP0 — 주식선물 실시간호가 (68 fields)
// --------------------------------------------------------------------------

func TestDecodeStockFuturesAsk(t *testing.T) {
	raw := loadFixture(t, "h0zfasp0_success.txt")
	f, err := parseFrame(raw)
	require.NoError(t, err)
	require.Equal(t, "H0ZFASP0", f.TrID)

	events, err := decodeStockFuturesAsk(f)
	require.NoError(t, err)
	require.Len(t, events, 1)

	ev := events[0]
	assert.Equal(t, "KAK0F", ev.Symbol)
	// ASKP1 = 73500
	assert.True(t, decimal.NewFromFloat(73500).Equal(ev.Ask[0]))
	// ASKP10 = 73950
	assert.True(t, decimal.NewFromFloat(73950).Equal(ev.Ask[9]))
	// BIDP1 = 73300
	assert.True(t, decimal.NewFromFloat(73300).Equal(ev.Bid[0]))
	// BIDP10 = 72850
	assert.True(t, decimal.NewFromFloat(72850).Equal(ev.Bid[9]))
	// 10단계 모두 non-zero
	for i := 0; i < 10; i++ {
		assert.False(t, ev.Ask[i].IsZero(), "Ask[%d] zero", i)
		assert.False(t, ev.Bid[i].IsZero(), "Bid[%d] zero", i)
	}
	assert.Equal(t, int64(1250), ev.TotalAskSize)
	assert.Equal(t, int64(2350), ev.TotalBidSize)
	assert.Len(t, ev.Raw, 68)
}

// --------------------------------------------------------------------------
// H0ZFANC0 — 주식선물 실시간예상체결 (8 fields)
// --------------------------------------------------------------------------

func TestDecodeStockFuturesExpectTrade(t *testing.T) {
	raw := loadFixture(t, "h0zfanc0_success.txt")
	f, err := parseFrame(raw)
	require.NoError(t, err)
	require.Equal(t, "H0ZFANC0", f.TrID)

	events, err := decodeStockFuturesExpectTrade(f)
	require.NoError(t, err)
	require.Len(t, events, 1)

	ev := events[0]
	assert.Equal(t, "KAK0F000000", ev.Symbol)
	assert.Equal(t, "133500", ev.Time)
	// ANTC_CNPR = 73400
	assert.True(t, decimal.NewFromFloat(73400).Equal(ev.ExpectPrice))
	// ANTC_CNQN = 200
	assert.Equal(t, int64(200), ev.ExpectQuantity)
	assert.Len(t, ev.Raw, 8)
}

// --------------------------------------------------------------------------
// H0ZOCNT0 — 주식옵션 실시간체결가 (53 fields)
// --------------------------------------------------------------------------

func TestDecodeStockOptionTrade(t *testing.T) {
	raw := loadFixture(t, "h0zocnt0_success.txt")
	f, err := parseFrame(raw)
	require.NoError(t, err)
	require.Equal(t, "H0ZOCNT0", f.TrID)

	events, err := decodeStockOptionTrade(f)
	require.NoError(t, err)
	require.Len(t, events, 1)

	ev := events[0]
	assert.Equal(t, "KAK0C", ev.Symbol)
	assert.Equal(t, "133500", ev.Time)
	// OPTN_PRPR = 2.80
	assert.True(t, decimal.NewFromFloat(2.80).Equal(ev.Price))
	// DELTA = 0.38
	assert.InDelta(t, 0.38, ev.Delta, 1e-9)
	// GAMA = 0.015
	assert.InDelta(t, 0.015, ev.Gamma, 1e-9)
	// THETA = -0.065
	assert.InDelta(t, -0.065, ev.Theta, 1e-9)
	// PRDY_VOL_VRSS_ACML_VOL_RATE = 11.20
	assert.InDelta(t, 11.20, ev.PrevVolRate, 1e-9)
	// DYNM 필드 없음
	assert.Len(t, ev.Raw, 53)
}

func TestDecodeStockOptionTrade_FieldCountMismatch(t *testing.T) {
	f := frame{Kind: frameKindRealtime, TrID: "H0ZOCNT0", Count: 1, Fields: []string{"a", "b"}}
	_, err := decodeStockOptionTrade(f)
	require.Error(t, err)
}

// --------------------------------------------------------------------------
// H0ZOASP0 — 주식옵션 실시간호가 (68 fields)
// --------------------------------------------------------------------------

func TestDecodeStockOptionAsk(t *testing.T) {
	raw := loadFixture(t, "h0zoasp0_success.txt")
	f, err := parseFrame(raw)
	require.NoError(t, err)
	require.Equal(t, "H0ZOASP0", f.TrID)

	events, err := decodeStockOptionAsk(f)
	require.NoError(t, err)
	require.Len(t, events, 1)

	ev := events[0]
	assert.Equal(t, "KAK0C", ev.Symbol)
	// OPTN_ASKP1 = 2.85
	assert.True(t, decimal.NewFromFloat(2.85).Equal(ev.Ask1to5[0]))
	// OPTN_ASKP5 = 3.05
	assert.True(t, decimal.NewFromFloat(3.05).Equal(ev.Ask1to5[4]))
	// OPTN_BIDP1 = 2.75
	assert.True(t, decimal.NewFromFloat(2.75).Equal(ev.Bid1to5[0]))
	// OPTN_ASKP6 = 3.10
	assert.True(t, decimal.NewFromFloat(3.10).Equal(ev.Ask6to10[0]))
	// OPTN_ASKP10 = 3.30
	assert.True(t, decimal.NewFromFloat(3.30).Equal(ev.Ask6to10[4]))
	// OPTN_BIDP6 = 2.50
	assert.True(t, decimal.NewFromFloat(2.50).Equal(ev.Bid6to10[0]))
	// TOTAL_ASKP_RSQN = 1750
	assert.Equal(t, int64(1750), ev.TotalAskSize)
	// TOTAL_BIDP_RSQN = 2020
	assert.Equal(t, int64(2020), ev.TotalBidSize)
	assert.Len(t, ev.Raw, 68)
}

// --------------------------------------------------------------------------
// H0ZOANC0 — 주식옵션 실시간예상체결 (7 fields)
// --------------------------------------------------------------------------

func TestDecodeStockOptionExpectTrade(t *testing.T) {
	raw := loadFixture(t, "h0zoanc0_success.txt")
	f, err := parseFrame(raw)
	require.NoError(t, err)
	require.Equal(t, "H0ZOANC0", f.TrID)

	events, err := decodeStockOptionExpectTrade(f)
	require.NoError(t, err)
	require.Len(t, events, 1)

	ev := events[0]
	assert.Equal(t, "KAK0C000000", ev.Symbol)
	assert.Equal(t, "133500", ev.Time)
	// ANTC_CNPR = 2.80
	assert.True(t, decimal.NewFromFloat(2.80).Equal(ev.ExpectPrice))
	// ANTC_CNTG_VRSS = 0.20
	assert.True(t, decimal.NewFromFloat(0.20).Equal(ev.ExpectDiff))
	// ANTC_CNTG_VRSS_SIGN = 2
	assert.Equal(t, "2", ev.ExpectDiffSign)
	// ANTC_MKOP_CLS_CODE = 310
	assert.Equal(t, "310", ev.ExpectMarketCode)
	// ANTC_CNQN 없음 — 7 fields
	assert.Len(t, ev.Raw, 7)
}

// --------------------------------------------------------------------------
// FieldCountMismatch 보강 (coverage ≥70%)
// --------------------------------------------------------------------------

func TestDecodeKrxNightFuturesAsk_FieldCountMismatch(t *testing.T) {
	f := frame{Kind: frameKindRealtime, TrID: "H0MFASP0", Count: 1, Fields: []string{"a", "b"}}
	_, err := decodeKrxNightFuturesAsk(f)
	require.Error(t, err)
}

func TestDecodeKrxNightOptionAsk_FieldCountMismatch(t *testing.T) {
	f := frame{Kind: frameKindRealtime, TrID: "H0EUASP0", Count: 1, Fields: []string{"a"}}
	_, err := decodeKrxNightOptionAsk(f)
	require.Error(t, err)
}

func TestDecodeKrxNightOptionExpectTrade_FieldCountMismatch(t *testing.T) {
	f := frame{Kind: frameKindRealtime, TrID: "H0EUANC0", Count: 1, Fields: []string{"a"}}
	_, err := decodeKrxNightOptionExpectTrade(f)
	require.Error(t, err)
}

func TestDecodeStockFuturesTrade_FieldCountMismatch(t *testing.T) {
	f := frame{Kind: frameKindRealtime, TrID: "H0ZFCNT0", Count: 1, Fields: []string{"a"}}
	_, err := decodeStockFuturesTrade(f)
	require.Error(t, err)
}

func TestDecodeStockFuturesAsk_FieldCountMismatch(t *testing.T) {
	f := frame{Kind: frameKindRealtime, TrID: "H0ZFASP0", Count: 1, Fields: []string{"a"}}
	_, err := decodeStockFuturesAsk(f)
	require.Error(t, err)
}

func TestDecodeStockFuturesExpectTrade_FieldCountMismatch(t *testing.T) {
	f := frame{Kind: frameKindRealtime, TrID: "H0ZFANC0", Count: 1, Fields: []string{"a"}}
	_, err := decodeStockFuturesExpectTrade(f)
	require.Error(t, err)
}

func TestDecodeStockOptionAsk_FieldCountMismatch(t *testing.T) {
	f := frame{Kind: frameKindRealtime, TrID: "H0ZOASP0", Count: 1, Fields: []string{"a"}}
	_, err := decodeStockOptionAsk(f)
	require.Error(t, err)
}

func TestDecodeStockOptionExpectTrade_FieldCountMismatch(t *testing.T) {
	f := frame{Kind: frameKindRealtime, TrID: "H0ZOANC0", Count: 1, Fields: []string{"a"}}
	_, err := decodeStockOptionExpectTrade(f)
	require.Error(t, err)
}
