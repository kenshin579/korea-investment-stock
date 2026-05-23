package overseasfutures_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kenshin579/korea-investment-stock/overseasfutures"
)

// 회귀: 전일대비율(prev_diff_rate)이 "+1.50" 부호 문자열로 와도 파싱돼야 한다.
func TestInquirePriceOutput1_SignedDiffRate(t *testing.T) {
	const body = `{"prev_diff_rate":"+1.50"}`
	var got overseasfutures.InquirePriceOutput1
	err := json.Unmarshal([]byte(body), &got)
	require.NoError(t, err)
	assert.InDelta(t, 1.50, float64(got.PrevDiffRate), 1e-9)
}
