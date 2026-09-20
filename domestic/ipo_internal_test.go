package domestic

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

// TestParsePaddedDecimal 은 parsePaddedDecimal 이 KIS 가 실제로 보내는(그리고
// 보낼 법한) 숫자 문자열 모양을 직접 검증한다. "-" → 0 은 이 helper 의 문서화된
// 동작인데 지금까지 어디서도 검증되지 않았다.
func TestParsePaddedDecimal(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want decimal.Decimal
	}{
		{"공백 좌측 패딩", "       19500", decimal.NewFromInt(19500)},
		{"0 패딩", "000000500", decimal.NewFromInt(500)},
		{"천단위 콤마", "19,500", decimal.NewFromInt(19500)},
		{"빈 문자열", "", decimal.Zero},
		{"공백만", "   ", decimal.Zero},
		{"sentinel 하이픈", "-", decimal.Zero},
		{"숫자 아님", "abc", decimal.Zero},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := parsePaddedDecimal(tc.in)
			assert.True(t, tc.want.Equal(got), "parsePaddedDecimal(%q) = %s, want %s", tc.in, got, tc.want)
		})
	}
}

// TestParsePaddedInt64 는 parsePaddedInt64 를 같은 케이스 집합으로 검증한다.
func TestParsePaddedInt64(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want int64
	}{
		{"공백 좌측 패딩", "           0", 0},
		{"공백 좌측 패딩(0 아님)", "     5063824", 5063824},
		{"0 패딩", "000000500", 500},
		{"천단위 콤마", "1,000,000", 1000000},
		{"빈 문자열", "", 0},
		{"공백만", "   ", 0},
		{"sentinel 하이픈", "-", 0},
		{"숫자 아님", "abc", 0},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, parsePaddedInt64(tc.in))
		})
	}
}
