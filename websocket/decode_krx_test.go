package websocket

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func loadFixture(t *testing.T, name string) string {
	t.Helper()
	b, err := os.ReadFile(filepath.Join("testdata", name))
	require.NoError(t, err)
	return strings.TrimRight(string(b), "\n")
}

func TestDecodeKrxTrade_Single(t *testing.T) {
	raw := loadFixture(t, "h0stcnt0_success.txt")
	f, err := parseFrame(raw)
	require.NoError(t, err)
	require.Equal(t, "H0STCNT0", f.TrID)

	events, err := decodeKrxTrade(f)
	require.NoError(t, err)
	require.Len(t, events, 1)

	ev := events[0]
	assert.Equal(t, "005930", ev.Symbol)
	assert.Equal(t, "123929", ev.Time)
	assert.True(t, decimal.NewFromInt(73100).Equal(ev.Price))
	assert.Equal(t, "2", ev.PrevDiffSign)
	assert.Equal(t, int64(150), ev.TradeVolume)
	assert.NotEmpty(t, ev.Raw)
}

func TestDecodeKrxTrade_Paging(t *testing.T) {
	raw := loadFixture(t, "h0stcnt0_paging.txt")
	f, err := parseFrame(raw)
	require.NoError(t, err)

	events, err := decodeKrxTrade(f)
	require.NoError(t, err)
	require.Len(t, events, 2)
}

func TestDecodeKrxTrade_BadNumeric(t *testing.T) {
	raw := strings.Replace(loadFixture(t, "h0stcnt0_success.txt"), "73100", "abc", 1)
	f, err := parseFrame(raw)
	require.NoError(t, err)

	events, err := decodeKrxTrade(f)
	require.NoError(t, err)
	require.Len(t, events, 1)
	assert.True(t, events[0].Price.IsZero())
}

func TestDecodeKrxAsk(t *testing.T) {
	raw := loadFixture(t, "h0stasp0_success.txt")
	f, err := parseFrame(raw)
	require.NoError(t, err)

	events, err := decodeKrxAsk(f)
	require.NoError(t, err)
	require.Len(t, events, 1)
	assert.Equal(t, "005930", events[0].Symbol)
	// Ask[0..9] / Bid[0..9] populated
	for i := 0; i < 10; i++ {
		assert.False(t, events[0].Ask[i].IsZero(), "Ask[%d] zero", i)
	}
}

func TestDecodeKrxExpectTrade(t *testing.T) {
	raw := loadFixture(t, "h0stanc0_success.txt")
	f, err := parseFrame(raw)
	require.NoError(t, err)

	events, err := decodeKrxExpectTrade(f)
	require.NoError(t, err)
	require.Len(t, events, 1)
	assert.Equal(t, "005930", events[0].Symbol)
}

func TestDecodeKrxOvernightTrade(t *testing.T) {
	raw := loadFixture(t, "h0stoup0_success.txt")
	f, err := parseFrame(raw)
	require.NoError(t, err)

	events, err := decodeKrxOvernightTrade(f)
	require.NoError(t, err)
	require.Len(t, events, 1)
	assert.Equal(t, "005930", events[0].Symbol)
}

func TestDecodeKrxOvernightExpect(t *testing.T) {
	raw := loadFixture(t, "h0stoac0_success.txt")
	f, err := parseFrame(raw)
	require.NoError(t, err)

	events, err := decodeKrxOvernightExpect(f)
	require.NoError(t, err)
	require.Len(t, events, 1)
	assert.Equal(t, "005930", events[0].Symbol)
}

func TestDecodeKrxTrade_FieldCountMismatch(t *testing.T) {
	// 일부 필드 누락 → ErrWSInvalidFrame
	f := frame{Kind: frameKindRealtime, TrID: "H0STCNT0", Count: 1, Fields: []string{"a", "b", "c"}}
	_, err := decodeKrxTrade(f)
	require.Error(t, err)
}
