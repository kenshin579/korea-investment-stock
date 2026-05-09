package websocket

import (
	"strings"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --------------------------------------------------------------------------
// H0IFCNT0 — 지수선물 실시간체결가 (50 fields)
// --------------------------------------------------------------------------

func TestDecodeIndexFuturesTrade(t *testing.T) {
	raw := loadFixture(t, "h0ifcnt0_success.txt")
	f, err := parseFrame(raw)
	require.NoError(t, err)
	require.Equal(t, "H0IFCNT0", f.TrID)

	events, err := decodeIndexFuturesTrade(f)
	require.NoError(t, err)
	require.Len(t, events, 1)

	ev := events[0]
	// [0] FUTS_SHRN_ISCD
	assert.Equal(t, "101S12", ev.Symbol)
	// [1] BSOP_HOUR
	assert.Equal(t, "090130", ev.Time)
	// [2] FUTS_PRDY_VRSS = 5.00
	assert.True(t, decimal.NewFromFloat(5.00).Equal(ev.PrevDiff))
	// [3] PRDY_VRSS_SIGN = 2
	assert.Equal(t, "2", ev.PrevDiffSign)
	// [4] FUTS_PRDY_CTRT = 0.15
	assert.InDelta(t, 0.15, ev.PrevChangeRate, 1e-9)
	// [5] FUTS_PRPR = 3395.00
	assert.True(t, decimal.NewFromFloat(3395.00).Equal(ev.Price))
	// [6] FUTS_OPRC = 3380.00
	assert.True(t, decimal.NewFromFloat(3380.00).Equal(ev.Open))
	// [7] FUTS_HGPR = 3400.00
	assert.True(t, decimal.NewFromFloat(3400.00).Equal(ev.High))
	// [8] FUTS_LWPR = 3370.00
	assert.True(t, decimal.NewFromFloat(3370.00).Equal(ev.Low))
	// [9] LAST_CNQN = 150
	assert.Equal(t, int64(150), ev.LastTradeVolume)
	// [10] ACML_VOL = 78234
	assert.Equal(t, int64(78234), ev.AccumVolume)
	// [11] ACML_TR_PBMN = 265234560000
	assert.Equal(t, int64(265234560000), ev.AccumValue)
	// [18] HTS_OTST_STPL_QTY = 87654
	assert.Equal(t, int64(87654), ev.OpenInterestQty)
	// [30] CTTR = 118.50
	assert.InDelta(t, 118.50, ev.TradeStrength, 1e-9)
	// [43] TOTAL_ASKP_RSQN = 45000
	assert.Equal(t, int64(45000), ev.TotalAskSize)
	// [44] TOTAL_BIDP_RSQN = 62000
	assert.Equal(t, int64(62000), ev.TotalBidSize)
	// [45] PRDY_VOL_VRSS_ACML_VOL_RATE = 2.35
	assert.InDelta(t, 2.35, ev.PrevVolRate, 1e-9)
	// [46] DSCS_BLTR_ACML_QTY = 0
	assert.Equal(t, int64(0), ev.BlockTradeVolume)
	// [47] DYNM_MXPR = 3500.00
	assert.True(t, decimal.NewFromFloat(3500.00).Equal(ev.DynamicUpperLimit))
	// [48] DYNM_LLAM = 3290.00
	assert.True(t, decimal.NewFromFloat(3290.00).Equal(ev.DynamicLowerLimit))
	// [49] DYNM_PRC_LIMT_YN = N
	assert.Equal(t, "N", ev.DynamicPriceLimitYN)
	assert.NotEmpty(t, ev.Raw)
	assert.Len(t, ev.Raw, 50)
}

func TestDecodeIndexFuturesTrade_FieldCountMismatch(t *testing.T) {
	f := frame{Kind: frameKindRealtime, TrID: "H0IFCNT0", Count: 1, Fields: []string{"a", "b", "c"}}
	_, err := decodeIndexFuturesTrade(f)
	require.Error(t, err)
}

// --------------------------------------------------------------------------
// H0CFCNT0 — 상품선물 실시간체결가 (50 fields, alias of IndexFuturesTradeEvent)
// --------------------------------------------------------------------------

func TestDecodeCommodityFuturesTrade(t *testing.T) {
	raw := loadFixture(t, "h0cfcnt0_success.txt")
	f, err := parseFrame(raw)
	require.NoError(t, err)
	require.Equal(t, "H0CFCNT0", f.TrID)

	// 상품선물은 IndexFuturesTrade decoder 재사용 (alias schema)
	events, err := decodeIndexFuturesTrade(f)
	require.NoError(t, err)
	require.Len(t, events, 1)

	ev := events[0]
	// [0] FUTS_SHRN_ISCD = GCF24
	assert.Equal(t, "GCF24", ev.Symbol)
	// [1] BSOP_HOUR = 090130
	assert.Equal(t, "090130", ev.Time)
	// [5] FUTS_PRPR = 46800.00
	assert.True(t, decimal.NewFromFloat(46800.00).Equal(ev.Price))
	// [9] LAST_CNQN = 50
	assert.Equal(t, int64(50), ev.LastTradeVolume)
	// [10] ACML_VOL = 12345
	assert.Equal(t, int64(12345), ev.AccumVolume)
	// [47] DYNM_MXPR = 49000.00
	assert.True(t, decimal.NewFromFloat(49000.00).Equal(ev.DynamicUpperLimit))
	// [48] DYNM_LLAM = 44500.00
	assert.True(t, decimal.NewFromFloat(44500.00).Equal(ev.DynamicLowerLimit))
	// [49] DYNM_PRC_LIMT_YN = N
	assert.Equal(t, "N", ev.DynamicPriceLimitYN)
	// alias type 호환성 검증
	var _ CommodityFuturesTradeEvent = ev
	assert.Len(t, ev.Raw, 50)
}

// --------------------------------------------------------------------------
// H0IFASP0 — 지수선물 실시간호가 (38 fields)
// --------------------------------------------------------------------------

func TestDecodeIndexFuturesAsk(t *testing.T) {
	raw := loadFixture(t, "h0ifasp0_success.txt")
	f, err := parseFrame(raw)
	require.NoError(t, err)
	require.Equal(t, "H0IFASP0", f.TrID)

	events, err := decodeIndexFuturesAsk(f)
	require.NoError(t, err)
	require.Len(t, events, 1)

	ev := events[0]
	// [0] FUTS_SHRN_ISCD = 101S12
	assert.Equal(t, "101S12", ev.Symbol)
	// [1] BSOP_HOUR = 090130
	assert.Equal(t, "090130", ev.Time)
	// Ask[0] = 3397.00 (FUTS_ASKP1)
	assert.True(t, decimal.NewFromFloat(3397.00).Equal(ev.Ask[0]))
	// Ask[4] = 3401.00 (FUTS_ASKP5)
	assert.True(t, decimal.NewFromFloat(3401.00).Equal(ev.Ask[4]))
	// Bid[0] = 3394.00 (FUTS_BIDP1)
	assert.True(t, decimal.NewFromFloat(3394.00).Equal(ev.Bid[0]))
	// Bid[4] = 3390.00 (FUTS_BIDP5)
	assert.True(t, decimal.NewFromFloat(3390.00).Equal(ev.Bid[4]))
	// AskCsnu[0] = 45 (ASKP_CSNU1)
	assert.Equal(t, int64(45), ev.AskCsnu[0])
	// BidCsnu[0] = 60 (BIDP_CSNU1)
	assert.Equal(t, int64(60), ev.BidCsnu[0])
	// AskSize[0] = 380 (ASKP_RSQN1)
	assert.Equal(t, int64(380), ev.AskSize[0])
	// BidSize[0] = 520 (BIDP_RSQN1)
	assert.Equal(t, int64(520), ev.BidSize[0])
	// [32] TOTAL_ASKP_CSNU = 120
	assert.Equal(t, int64(120), ev.TotalAskCsnu)
	// [33] TOTAL_BIDP_CSNU = 155
	assert.Equal(t, int64(155), ev.TotalBidCsnu)
	// [34] TOTAL_ASKP_RSQN = 1010
	assert.Equal(t, int64(1010), ev.TotalAskSize)
	// [35] TOTAL_BIDP_RSQN = 1400
	assert.Equal(t, int64(1400), ev.TotalBidSize)
	// [36] TOTAL_ASKP_RSQN_ICDC = -150
	assert.Equal(t, int64(-150), ev.TotalAskSizeChg)
	// [37] TOTAL_BIDP_RSQN_ICDC = 200
	assert.Equal(t, int64(200), ev.TotalBidSizeChg)
	// 5단계 모두 non-zero
	for i := 0; i < 5; i++ {
		assert.False(t, ev.Ask[i].IsZero(), "Ask[%d] zero", i)
		assert.False(t, ev.Bid[i].IsZero(), "Bid[%d] zero", i)
	}
	assert.Len(t, ev.Raw, 38)
}

func TestDecodeIndexFuturesAsk_FieldCountMismatch(t *testing.T) {
	f := frame{Kind: frameKindRealtime, TrID: "H0IFASP0", Count: 1, Fields: []string{"a", "b"}}
	_, err := decodeIndexFuturesAsk(f)
	require.Error(t, err)
}

// --------------------------------------------------------------------------
// H0CFASP0 — 상품선물 실시간호가 (38 fields, alias of IndexFuturesAskEvent)
// --------------------------------------------------------------------------

func TestDecodeCommodityFuturesAsk(t *testing.T) {
	raw := loadFixture(t, "h0cfasp0_success.txt")
	f, err := parseFrame(raw)
	require.NoError(t, err)
	require.Equal(t, "H0CFASP0", f.TrID)

	// 상품선물은 IndexFuturesAsk decoder 재사용 (alias schema)
	events, err := decodeIndexFuturesAsk(f)
	require.NoError(t, err)
	require.Len(t, events, 1)

	ev := events[0]
	// [0] FUTS_SHRN_ISCD = GCF24
	assert.Equal(t, "GCF24", ev.Symbol)
	// [1] BSOP_HOUR = 090130
	assert.Equal(t, "090130", ev.Time)
	// Ask[0] = 46850.00 (FUTS_ASKP1)
	assert.True(t, decimal.NewFromFloat(46850.00).Equal(ev.Ask[0]))
	// Ask[4] = 47050.00 (FUTS_ASKP5)
	assert.True(t, decimal.NewFromFloat(47050.00).Equal(ev.Ask[4]))
	// Bid[0] = 46750.00 (FUTS_BIDP1)
	assert.True(t, decimal.NewFromFloat(46750.00).Equal(ev.Bid[0]))
	// [34] TOTAL_ASKP_RSQN = 570
	assert.Equal(t, int64(570), ev.TotalAskSize)
	// [35] TOTAL_BIDP_RSQN = 820
	assert.Equal(t, int64(820), ev.TotalBidSize)
	// alias type 호환성 검증
	var _ CommodityFuturesAskEvent = ev
	assert.Len(t, ev.Raw, 38)
}

// --------------------------------------------------------------------------
// H0IOCNT0 — 지수옵션 실시간체결가 (58 fields)
// --------------------------------------------------------------------------

func TestDecodeIndexOptionTrade(t *testing.T) {
	raw := loadFixture(t, "h0iocnt0_success.txt")
	f, err := parseFrame(raw)
	require.NoError(t, err)
	require.Equal(t, "H0IOCNT0", f.TrID)

	events, err := decodeIndexOptionTrade(f)
	require.NoError(t, err)
	require.Len(t, events, 1)

	ev := events[0]
	// [0] OPTN_SHRN_ISCD = 201S11305
	assert.Equal(t, "201S11305", ev.Symbol)
	// [1] BSOP_HOUR = 090130
	assert.Equal(t, "090130", ev.Time)
	// [2] OPTN_PRPR = 3.50
	assert.True(t, decimal.NewFromFloat(3.50).Equal(ev.Price))
	// [3] PRDY_VRSS_SIGN = 2
	assert.Equal(t, "2", ev.PrevDiffSign)
	// [4] OPTN_PRDY_VRSS = 0.30
	assert.True(t, decimal.NewFromFloat(0.30).Equal(ev.PrevDiff))
	// [5] PRDY_CTRT = 9.38
	assert.InDelta(t, 9.38, ev.PrevChangeRate, 1e-9)
	// [9] LAST_CNQN = 200
	assert.Equal(t, int64(200), ev.LastTradeVolume)
	// [10] ACML_VOL = 85432
	assert.Equal(t, int64(85432), ev.AccumVolume)
	// 옵션 그릭스
	// [28] DELTA = 0.4523
	assert.InDelta(t, 0.4523, ev.Delta, 1e-9)
	// [29] GAMA = 0.0182
	assert.InDelta(t, 0.0182, ev.Gamma, 1e-9)
	// [30] VEGA = 0.8234
	assert.InDelta(t, 0.8234, ev.Vega, 1e-9)
	// [31] THETA = -0.0752
	assert.InDelta(t, -0.0752, ev.Theta, 1e-9)
	// [32] RHO = 0.0091
	assert.InDelta(t, 0.0091, ev.Rho, 1e-9)
	// [33] HTS_INTS_VLTL = 18.50
	assert.InDelta(t, 18.50, ev.ImpliedVolatility, 1e-9)
	// [37] UNAS_HIST_VLTL = 16.30
	assert.InDelta(t, 16.30, ev.HistoricalVolatility, 1e-9)
	// [38] CTTR = 112.80
	assert.InDelta(t, 112.80, ev.TradeStrength, 1e-9)
	// [52] PRDY_VOL_VRSS_ACML_VOL_RATE = 3.85
	assert.InDelta(t, 3.85, ev.PrevVolRate, 1e-9)
	// [53] AVRG_VLTL = 15.20 (지수옵션 전용)
	assert.InDelta(t, 15.20, ev.AvgVolatility, 1e-9)
	// [54] DSCS_LRQN_VOL = 0
	assert.Equal(t, int64(0), ev.BlockTradeVolume)
	// [55] DYNM_MXPR = 5.00
	assert.True(t, decimal.NewFromFloat(5.00).Equal(ev.DynamicUpperLimit))
	// [56] DYNM_LLAM = 1.50
	assert.True(t, decimal.NewFromFloat(1.50).Equal(ev.DynamicLowerLimit))
	// [57] DYNM_PRC_LIMT_YN = N
	assert.Equal(t, "N", ev.DynamicPriceLimitYN)
	assert.NotEmpty(t, ev.Raw)
	assert.Len(t, ev.Raw, 58)
}

func TestDecodeIndexOptionTrade_FieldCountMismatch(t *testing.T) {
	f := frame{Kind: frameKindRealtime, TrID: "H0IOCNT0", Count: 1, Fields: []string{"a", "b", "c"}}
	_, err := decodeIndexOptionTrade(f)
	require.Error(t, err)
}

// --------------------------------------------------------------------------
// H0IOASP0 — 지수옵션 실시간호가 (38 fields)
// --------------------------------------------------------------------------

func TestDecodeIndexOptionAsk(t *testing.T) {
	raw := loadFixture(t, "h0ioasp0_success.txt")
	f, err := parseFrame(raw)
	require.NoError(t, err)
	require.Equal(t, "H0IOASP0", f.TrID)

	events, err := decodeIndexOptionAsk(f)
	require.NoError(t, err)
	require.Len(t, events, 1)

	ev := events[0]
	// [0] OPTN_SHRN_ISCD = 201S11305
	assert.Equal(t, "201S11305", ev.Symbol)
	// [1] BSOP_HOUR = 090130
	assert.Equal(t, "090130", ev.Time)
	// Ask[0] = 3.55 (OPTN_ASKP1)
	assert.True(t, decimal.NewFromFloat(3.55).Equal(ev.Ask[0]))
	// Ask[4] = 3.75 (OPTN_ASKP5)
	assert.True(t, decimal.NewFromFloat(3.75).Equal(ev.Ask[4]))
	// Bid[0] = 3.45 (OPTN_BIDP1)
	assert.True(t, decimal.NewFromFloat(3.45).Equal(ev.Bid[0]))
	// Bid[4] = 3.25 (OPTN_BIDP5)
	assert.True(t, decimal.NewFromFloat(3.25).Equal(ev.Bid[4]))
	// AskCsnu[0] = 35 (ASKP_CSNU1)
	assert.Equal(t, int64(35), ev.AskCsnu[0])
	// BidCsnu[0] = 50 (BIDP_CSNU1)
	assert.Equal(t, int64(50), ev.BidCsnu[0])
	// AskSize[0] = 450 (ASKP_RSQN1)
	assert.Equal(t, int64(450), ev.AskSize[0])
	// BidSize[0] = 680 (BIDP_RSQN1)
	assert.Equal(t, int64(680), ev.BidSize[0])
	// [32] TOTAL_ASKP_CSNU = 98
	assert.Equal(t, int64(98), ev.TotalAskCsnu)
	// [33] TOTAL_BIDP_CSNU = 132
	assert.Equal(t, int64(132), ev.TotalBidCsnu)
	// [34] TOTAL_ASKP_RSQN = 1240
	assert.Equal(t, int64(1240), ev.TotalAskSize)
	// [35] TOTAL_BIDP_RSQN = 1810
	assert.Equal(t, int64(1810), ev.TotalBidSize)
	// [36] TOTAL_ASKP_RSQN_ICDC = -200
	assert.Equal(t, int64(-200), ev.TotalAskSizeChg)
	// [37] TOTAL_BIDP_RSQN_ICDC = 300
	assert.Equal(t, int64(300), ev.TotalBidSizeChg)
	// 5단계 모두 non-zero
	for i := 0; i < 5; i++ {
		assert.False(t, ev.Ask[i].IsZero(), "Ask[%d] zero", i)
		assert.False(t, ev.Bid[i].IsZero(), "Bid[%d] zero", i)
	}
	assert.Len(t, ev.Raw, 38)
}

func TestDecodeIndexOptionAsk_FieldCountMismatch(t *testing.T) {
	f := frame{Kind: frameKindRealtime, TrID: "H0IOASP0", Count: 1, Fields: []string{"a", "b"}}
	_, err := decodeIndexOptionAsk(f)
	require.Error(t, err)
}

// --------------------------------------------------------------------------
// BadNumeric 보강 (parse fail → zero, coverage ≥70%)
// --------------------------------------------------------------------------

func TestDecodeIndexFuturesTrade_BadNumeric(t *testing.T) {
	raw := strings.Replace(loadFixture(t, "h0ifcnt0_success.txt"), "270.50", "abc", 1)
	f, err := parseFrame(raw)
	require.NoError(t, err)
	events, err := decodeIndexFuturesTrade(f)
	require.NoError(t, err)
	require.Len(t, events, 1)
}

func TestDecodeIndexOptionTrade_BadNumeric(t *testing.T) {
	raw := strings.Replace(loadFixture(t, "h0iocnt0_success.txt"), "0.5", "xyz", 1)
	f, err := parseFrame(raw)
	require.NoError(t, err)
	events, err := decodeIndexOptionTrade(f)
	require.NoError(t, err)
	require.Len(t, events, 1)
}
