package futures_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kenshin579/korea-investment-stock/futures"
)

// 회귀: 대비율(+0.85)·괴리율(빈 값) 등이 부호/빈문자열로 와도 파싱돼야 한다.
func TestInquirePriceOutput1_SignedFields(t *testing.T) {
	const body = `{"futs_prdy_ctrt":"+0.85","dprt":""}`
	var got futures.InquirePriceOutput1
	err := json.Unmarshal([]byte(body), &got)
	require.NoError(t, err)
	assert.InDelta(t, 0.85, float64(got.FutsPrdyCtrt), 1e-9)
	assert.InDelta(t, 0, float64(got.Dprt), 1e-9)
}
