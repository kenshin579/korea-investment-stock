package domestic

import (
	"encoding/json"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLenient_EmptyNumbersBecomeZero(t *testing.T) {
	var it InvestorTradeByStockDailyItem
	raw := []byte(`{
		"stck_bsop_date":"20260713",
		"stck_clpr":"",
		"stck_oprc":"1200",
		"frgn_ntby_qty":"",
		"orgn_ntby_qty":"-5",
		"prdy_ctrt":""
	}`)
	require.NoError(t, json.Unmarshal(raw, &it))
	assert.True(t, it.StckClpr.Equal(decimal.Zero), "empty decimal -> 0")
	assert.True(t, it.StckOprc.Equal(decimal.NewFromInt(1200)))
	assert.Equal(t, int64(0), it.FrgnNtbyQty, "empty ,string int -> 0")
	assert.Equal(t, int64(-5), it.OrgnNtbyQty)
	assert.Equal(t, float64(0), it.PrdyCtrt, "empty ,string float -> 0")
	assert.Equal(t, "20260713", it.StckBsopDate)
}

func TestLenient_PreservesEmptyStringFields(t *testing.T) {
	var s InvestorTradeByStockDailySummary
	raw := []byte(`{"stck_prpr":"","rprs_mrkt_kor_name":"","prdy_vrss_sign":""}`)
	require.NoError(t, json.Unmarshal(raw, &s))
	assert.True(t, s.StckPrpr.Equal(decimal.Zero))
	assert.Equal(t, "", s.RprsMrktKorName, "string field must stay empty, not \"0\"")
	assert.Equal(t, "", s.PrdyVrssSign)
}

func TestLenient_NormalValuesUnaffected(t *testing.T) {
	var it InvestorTradeByStockDailyItem
	raw := []byte(`{"stck_clpr":"75800","acml_vol":"1000000"}`)
	require.NoError(t, json.Unmarshal(raw, &it))
	assert.True(t, it.StckClpr.Equal(decimal.NewFromInt(75800)))
	assert.Equal(t, int64(1000000), it.AcmlVol)
}

func TestLenient_ProgramItemEmptyNumbers(t *testing.T) {
	var it ProgramTradeByStockDailyItem
	raw := []byte(`{"stck_bsop_date":"20260713","stck_clpr":"","whol_smtn_ntby_qty":""}`)
	require.NoError(t, json.Unmarshal(raw, &it))
	assert.True(t, it.StckClpr.Equal(decimal.Zero))
	assert.Equal(t, int64(0), it.WholSmtnNtbyQty)
}
