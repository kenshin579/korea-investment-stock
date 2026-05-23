package overseas_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kenshin579/korea-investment-stock/overseas"
)

// KIS 해외 응답은 등락 필드를 "+0.69" 처럼 부호 붙은 문자열로, 빈 값은 "" 로 준다.
// 회귀: float64,string 태그로는 '+' 에서 파싱이 깨졌다 (kistypes.Float 로 수정).
func TestPriceDetailSnapshot_SignedRate(t *testing.T) {
	const body = `{"rsym":"DNASAAPL","t_xrat":"+0.69","p_xrat":"","perx":"29.45","tvol":"100","pvol":"200","tomv":"1","pamt":"1","shar":"1","mcap":"1","tamt":"1"}`
	var got overseas.PriceDetailSnapshot
	err := json.Unmarshal([]byte(body), &got)
	require.NoError(t, err)
	assert.InDelta(t, 0.69, float64(got.TXrat), 1e-9)
	assert.InDelta(t, 0, float64(got.PXrat), 1e-9)
	assert.InDelta(t, 29.45, float64(got.Perx), 1e-9)
}

func TestDailyPriceCandle_SignedRate(t *testing.T) {
	const body = `{"xymd":"20260505","rate":"+0.69","tvol":"1","tamt":"1","vbid":"1","vask":"1"}`
	var got overseas.DailyPriceCandle
	err := json.Unmarshal([]byte(body), &got)
	require.NoError(t, err)
	assert.InDelta(t, 0.69, float64(got.Rate), 1e-9)
}
