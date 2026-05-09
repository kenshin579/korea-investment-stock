package websocket

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --------------------------------------------------------------------------
// HDFFF020 — 해외선물옵션 실시간체결가 (25 fields)
// --------------------------------------------------------------------------

func TestDecodeOverseasFuturesTrade(t *testing.T) {
	raw := loadFixture(t, "hdfff020_success.txt")
	f, err := parseFrame(raw)
	require.NoError(t, err)
	require.Equal(t, "HDFFF020", f.TrID)

	events, err := decodeOverseasFuturesTrade(f)
	require.NoError(t, err)
	require.Len(t, events, 1)

	ev := events[0]
	// [0] SERIES_CD = ESM24
	assert.Equal(t, "ESM24", ev.Symbol)
	// [1] BSNS_DATE = 20240315
	assert.Equal(t, "20240315", ev.BsnsDate)
	// [2] MRKT_OPEN_DATE = 20240314
	assert.Equal(t, "20240314", ev.MrktOpenDate)
	// [3] MRKT_OPEN_TIME = 173000
	assert.Equal(t, "173000", ev.MrktOpenTime)
	// [4] MRKT_CLOSE_DATE = 20240315
	assert.Equal(t, "20240315", ev.MrktCloseDate)
	// [5] MRKT_CLOSE_TIME = 160000
	assert.Equal(t, "160000", ev.MrktCloseTime)
	// [6] PREV_PRICE = 52150
	assert.True(t, decimal.NewFromInt(52150).Equal(ev.PrevPrice))
	// [7] RECV_DATE = 20240315
	assert.Equal(t, "20240315", ev.RecvDate)
	// [8] RECV_TIME = 091530
	assert.Equal(t, "091530", ev.RecvTime)
	// [9] ACTIVE_FLAG = 1
	assert.Equal(t, "1", ev.ActiveFlag)
	// [10] LAST_PRICE = 52200
	assert.True(t, decimal.NewFromInt(52200).Equal(ev.LastPrice))
	// [11] LAST_QNTT = 2
	assert.Equal(t, int64(2), ev.LastQntt)
	// [12] PREV_DIFF_PRICE = 50
	assert.True(t, decimal.NewFromInt(50).Equal(ev.PrevDiffPrice))
	// [13] PREV_DIFF_RATE = 0.10
	assert.InDelta(t, 0.10, ev.PrevDiffRate, 1e-9)
	// [14] OPEN_PRICE = 52100
	assert.True(t, decimal.NewFromInt(52100).Equal(ev.OpenPrice))
	// [15] HIGH_PRICE = 52300
	assert.True(t, decimal.NewFromInt(52300).Equal(ev.HighPrice))
	// [16] LOW_PRICE = 52050
	assert.True(t, decimal.NewFromInt(52050).Equal(ev.LowPrice))
	// [17] VOL = 12345
	assert.Equal(t, int64(12345), ev.Vol)
	// [18] PREV_SIGN = 2
	assert.Equal(t, "2", ev.PrevSign)
	// [19] QUOTSIGN = 2
	assert.Equal(t, "2", ev.QuotSign)
	// [20] RECV_TIME2 = 0015
	assert.Equal(t, "0015", ev.RecvTime2)
	// [21] PSTTL_PRICE = 52000
	assert.True(t, decimal.NewFromInt(52000).Equal(ev.PsttlPrice))
	// [22] PSTTL_SIGN = 2
	assert.Equal(t, "2", ev.PsttlSign)
	// [23] PSTTL_DIFF_PRICE = 50
	assert.True(t, decimal.NewFromInt(50).Equal(ev.PsttlDiffPrice))
	// [24] PSTTL_DIFF_RATE = 0.10
	assert.InDelta(t, 0.10, ev.PsttlDiffRate, 1e-9)
	assert.NotEmpty(t, ev.Raw)
	assert.Len(t, ev.Raw, 25)
}

func TestDecodeOverseasFuturesTrade_FieldCountMismatch(t *testing.T) {
	f := frame{Kind: frameKindRealtime, TrID: "HDFFF020", Count: 1, Fields: []string{"a", "b", "c"}}
	_, err := decodeOverseasFuturesTrade(f)
	require.Error(t, err)
}

// --------------------------------------------------------------------------
// HDFFF010 — 해외선물옵션 실시간호가 (35 fields, BID/ASK 교차 배열)
// --------------------------------------------------------------------------

func TestDecodeOverseasFuturesAsk(t *testing.T) {
	raw := loadFixture(t, "hdfff010_success.txt")
	f, err := parseFrame(raw)
	require.NoError(t, err)
	require.Equal(t, "HDFFF010", f.TrID)

	events, err := decodeOverseasFuturesAsk(f)
	require.NoError(t, err)
	require.Len(t, events, 1)

	ev := events[0]
	// [0] SERIES_CD = ESM24
	assert.Equal(t, "ESM24", ev.Symbol)
	// [1] RECV_DATE = 20240315
	assert.Equal(t, "20240315", ev.RecvDate)
	// [2] RECV_TIME = 091530000000 (12자리, 나노초 포함)
	assert.Equal(t, "091530000000", ev.RecvTime)
	// [3] PREV_PRICE = 52150
	assert.True(t, decimal.NewFromInt(52150).Equal(ev.PrevPrice))

	// BID 1단계: [4]=150, [5]=0000000001, [6]=52195
	assert.Equal(t, int64(150), ev.BidQntt[0])
	assert.Equal(t, "0000000001", ev.BidNum[0])
	assert.True(t, decimal.NewFromInt(52195).Equal(ev.BidPrice[0]))
	// ASK 1단계: [7]=120, [8]=0000000002, [9]=52200
	assert.Equal(t, int64(120), ev.AskQntt[0])
	assert.Equal(t, "0000000002", ev.AskNum[0])
	assert.True(t, decimal.NewFromInt(52200).Equal(ev.AskPrice[0]))

	// BID 2단계: [10]=100, [11]=0000000003, [12]=52185
	assert.Equal(t, int64(100), ev.BidQntt[1])
	assert.Equal(t, "0000000003", ev.BidNum[1])
	assert.True(t, decimal.NewFromInt(52185).Equal(ev.BidPrice[1]))
	// ASK 2단계: [13]=80, [14]=0000000004, [15]=52205
	assert.Equal(t, int64(80), ev.AskQntt[1])
	assert.Equal(t, "0000000004", ev.AskNum[1])
	assert.True(t, decimal.NewFromInt(52205).Equal(ev.AskPrice[1]))

	// BID 5단계: [28]=20, [29]=0000000009, [30]=52170
	// fixture: 4=150,5=0000000001,6=52195, 7=120,8=0000000002,9=52200
	//          10=100,11=0000000003,12=52185, 13=80,14=0000000004,15=52205
	//          16=60,17=0000000005,18=52180, 19=50,20=0000000006,21=52210
	//          22=40,23=0000000007,24=52175, 25=30,26=0000000008,27=52215
	//          28=20,29=0000000009,30=52170, 31=10,32=0000000010,33=52220
	assert.Equal(t, int64(20), ev.BidQntt[4])
	assert.Equal(t, "0000000009", ev.BidNum[4])
	assert.True(t, decimal.NewFromInt(52170).Equal(ev.BidPrice[4]))
	// ASK 5단계: [31]=10, [32]=0000000010, [33]=52220
	assert.Equal(t, int64(10), ev.AskQntt[4])
	assert.Equal(t, "0000000010", ev.AskNum[4])
	assert.True(t, decimal.NewFromInt(52220).Equal(ev.AskPrice[4]))

	// [34] STTL_PRICE = 52000
	assert.True(t, decimal.NewFromInt(52000).Equal(ev.SttlPrice))
	assert.NotEmpty(t, ev.Raw)
	assert.Len(t, ev.Raw, 35)
}

func TestDecodeOverseasFuturesAsk_FieldCountMismatch(t *testing.T) {
	f := frame{Kind: frameKindRealtime, TrID: "HDFFF010", Count: 1, Fields: []string{"a", "b", "c"}}
	_, err := decodeOverseasFuturesAsk(f)
	require.Error(t, err)
}
