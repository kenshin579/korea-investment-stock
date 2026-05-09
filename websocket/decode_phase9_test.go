package websocket

import (
	"strings"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --------------------------------------------------------------------------
// 체결가 (H0NXCNT0 / H0UNCNT0) — 46 fields
// --------------------------------------------------------------------------

func TestDecodeAltMarketTrade_NXT(t *testing.T) {
	raw := loadFixture(t, "h0nxcnt0_success.txt")
	f, err := parseFrame(raw)
	require.NoError(t, err)
	require.Equal(t, "H0NXCNT0", f.TrID)

	events, err := decodeAltMarketTrade(f)
	require.NoError(t, err)
	require.Len(t, events, 1)

	ev := events[0]
	assert.Equal(t, "005930", ev.Symbol)
	assert.Equal(t, "123929", ev.Time)
	assert.True(t, decimal.NewFromInt(73100).Equal(ev.Price))
	assert.Equal(t, "2", ev.PrevDiffSign)
	assert.Equal(t, int64(150), ev.TradeVolume)
	assert.Equal(t, "1", ev.TradeKind) // CNTG_CLS_CODE
	assert.True(t, decimal.NewFromInt(72500).Equal(ev.ViStandardPrice))
	assert.NotEmpty(t, ev.Raw)
}

func TestDecodeAltMarketTrade_Unified(t *testing.T) {
	raw := loadFixture(t, "h0uncnt0_success.txt")
	f, err := parseFrame(raw)
	require.NoError(t, err)
	require.Equal(t, "H0UNCNT0", f.TrID)

	events, err := decodeAltMarketTrade(f)
	require.NoError(t, err)
	require.Len(t, events, 1)
	assert.Equal(t, "005930", events[0].Symbol)
}

func TestDecodeAltMarketTrade_FieldCountMismatch(t *testing.T) {
	f := frame{Kind: frameKindRealtime, TrID: "H0NXCNT0", Count: 1, Fields: []string{"a", "b", "c"}}
	_, err := decodeAltMarketTrade(f)
	require.Error(t, err)
}

// --------------------------------------------------------------------------
// 호가 (H0NXASP0 / H0UNASP0) — 65 fields
// --------------------------------------------------------------------------

func TestDecodeAltMarketAsk_NXT(t *testing.T) {
	raw := loadFixture(t, "h0nxasp0_success.txt")
	f, err := parseFrame(raw)
	require.NoError(t, err)
	require.Equal(t, "H0NXASP0", f.TrID)

	events, err := decodeAltMarketAsk(f)
	require.NoError(t, err)
	require.Len(t, events, 1)

	ev := events[0]
	assert.Equal(t, "005930", ev.Symbol)
	for i := 0; i < 10; i++ {
		assert.False(t, ev.Ask[i].IsZero(), "Ask[%d] zero", i)
		assert.False(t, ev.Bid[i].IsZero(), "Bid[%d] zero", i)
	}
	// KRX/NXT 중간가 6 fields (NXT 변형 only)
	assert.True(t, decimal.NewFromInt(73050).Equal(ev.KrxMidPrice))
	assert.Equal(t, int64(65000), ev.KrxMidTotalSize)
	assert.Equal(t, "1", ev.KrxMidCode)
	assert.True(t, decimal.NewFromInt(73050).Equal(ev.NxtMidPrice))
	assert.Equal(t, int64(32000), ev.NxtMidTotalSize)
	assert.Equal(t, "2", ev.NxtMidCode)
}

func TestDecodeAltMarketAsk_Unified(t *testing.T) {
	raw := loadFixture(t, "h0unasp0_success.txt")
	f, err := parseFrame(raw)
	require.NoError(t, err)

	events, err := decodeAltMarketAsk(f)
	require.NoError(t, err)
	require.Len(t, events, 1)
	assert.Equal(t, "005930", events[0].Symbol)
}

// --------------------------------------------------------------------------
// 예상체결 (H0NXANC0 / H0UNANC0) — 46 fields (KRX 45 + VI_STND_PRC)
// --------------------------------------------------------------------------

func TestDecodeAltMarketExpectTrade_NXT(t *testing.T) {
	raw := loadFixture(t, "h0nxanc0_success.txt")
	f, err := parseFrame(raw)
	require.NoError(t, err)
	require.Equal(t, "H0NXANC0", f.TrID)

	events, err := decodeAltMarketExpectTrade(f)
	require.NoError(t, err)
	require.Len(t, events, 1)
	ev := events[0]
	assert.Equal(t, "005930", ev.Symbol)
	// VI_STND_PRC 가 KRX 와 차이 — 0 이 아닌 값이어야 함
	assert.False(t, ev.ViStandardPrice.IsZero())
}

func TestDecodeAltMarketExpectTrade_Unified(t *testing.T) {
	raw := loadFixture(t, "h0unanc0_success.txt")
	f, err := parseFrame(raw)
	require.NoError(t, err)

	events, err := decodeAltMarketExpectTrade(f)
	require.NoError(t, err)
	require.Len(t, events, 1)
	assert.Equal(t, "005930", events[0].Symbol)
}

// --------------------------------------------------------------------------
// 프로그램매매 (H0NXPGM0 / H0UNPGM0) — 11 fields, 신규 EP
// --------------------------------------------------------------------------

func TestDecodeProgramTrade_NXT(t *testing.T) {
	raw := loadFixture(t, "h0nxpgm0_success.txt")
	f, err := parseFrame(raw)
	require.NoError(t, err)
	require.Equal(t, "H0NXPGM0", f.TrID)

	events, err := decodeProgramTrade(f)
	require.NoError(t, err)
	require.Len(t, events, 1)
	ev := events[0]
	assert.Equal(t, "005930", ev.Symbol)
	assert.Equal(t, "123929", ev.Time)
	assert.Equal(t, int64(15000), ev.AskQuantity)
	assert.Equal(t, int64(1100000000), ev.AskValue)
	assert.Equal(t, int64(25000), ev.BidQuantity)
	assert.Equal(t, int64(10000), ev.NetQuantity)
	assert.Equal(t, int64(3000), ev.TotalNetQuantity)
}

func TestDecodeProgramTrade_Unified(t *testing.T) {
	raw := loadFixture(t, "h0unpgm0_success.txt")
	f, err := parseFrame(raw)
	require.NoError(t, err)

	events, err := decodeProgramTrade(f)
	require.NoError(t, err)
	require.Len(t, events, 1)
	assert.Equal(t, "005930", events[0].Symbol)
}

// --------------------------------------------------------------------------
// 회원사 (H0NXMBC0 / H0UNMBC0) — 78 fields, 신규 EP
// --------------------------------------------------------------------------

func TestDecodeMember_NXT(t *testing.T) {
	raw := loadFixture(t, "h0nxmbc0_success.txt")
	f, err := parseFrame(raw)
	require.NoError(t, err)
	require.Equal(t, "H0NXMBC0", f.TrID)

	events, err := decodeMember(f)
	require.NoError(t, err)
	require.Len(t, events, 1)
	ev := events[0]
	assert.Equal(t, "005930", ev.Symbol)
	assert.Equal(t, "한국증권", ev.SellBrokerNames[0])
	assert.Equal(t, "키움증권", ev.SellBrokerNames[4])
	assert.Equal(t, "모건스탠리", ev.BuyBrokerNames[0])
	assert.Equal(t, int64(150000), ev.TotalSellQty[0])
	assert.Equal(t, int64(180000), ev.TotalBuyQty[0])
	assert.Equal(t, "N", ev.SellGlobalYN[0])
	assert.Equal(t, "Y", ev.BuyGlobalYN[0])
	assert.Equal(t, "001", ev.SellBrokerCodes[0])
	assert.Equal(t, "010", ev.BuyBrokerCodes[0])
	assert.InDelta(t, 15.5, ev.SellRatio[0], 1e-6)
	assert.InDelta(t, 18.0, ev.BuyRatio[0], 1e-6)
	assert.Equal(t, int64(5000), ev.SellQtyChange[0])
	assert.Equal(t, int64(8000), ev.BuyQtyChange[0])
	assert.Equal(t, int64(350000), ev.GlobalTotalSellQty)
	assert.Equal(t, int64(380000), ev.GlobalTotalBuyQty)
	assert.Equal(t, int64(15000), ev.GlobalSellQtyChange)
	assert.Equal(t, int64(18000), ev.GlobalBuyQtyChange)
	assert.Equal(t, int64(30000), ev.GlobalNetBuyQty)
	assert.InDelta(t, 35.5, ev.GlobalSellRatio, 1e-6)
	assert.InDelta(t, 38.0, ev.GlobalBuyRatio, 1e-6)
	assert.Equal(t, "Korea Inv", ev.SellBrokerEngNames[0])
	assert.Equal(t, "Goldman Sachs", ev.BuyBrokerEngNames[4])
}

func TestDecodeMember_Unified(t *testing.T) {
	raw := loadFixture(t, "h0unmbc0_success.txt")
	f, err := parseFrame(raw)
	require.NoError(t, err)

	events, err := decodeMember(f)
	require.NoError(t, err)
	require.Len(t, events, 1)
	assert.Equal(t, "005930", events[0].Symbol)
}

// --------------------------------------------------------------------------
// BadNumeric (parse fail → zero) — 1 base 만 검증
// --------------------------------------------------------------------------

func TestDecodeAltMarketTrade_BadNumeric(t *testing.T) {
	raw := strings.Replace(loadFixture(t, "h0nxcnt0_success.txt"), "73100", "abc", 1)
	f, err := parseFrame(raw)
	require.NoError(t, err)

	events, err := decodeAltMarketTrade(f)
	require.NoError(t, err)
	require.Len(t, events, 1)
	assert.True(t, events[0].Price.IsZero())
}
