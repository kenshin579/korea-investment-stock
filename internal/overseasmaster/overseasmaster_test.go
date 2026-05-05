package overseasmaster

import (
	"os"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse_NAS(t *testing.T) {
	zipBytes, err := os.ReadFile("testdata/nas_code_sample.cod.zip")
	require.NoError(t, err)

	syms, err := Parse("nas", zipBytes)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(syms), 1, "nas sample 은 최소 1 행 포함")

	// 종목 코드는 영문 대문자 + 숫자 (NASDAQ 형식)
	symRe := regexp.MustCompile(`^[A-Z0-9.]+$`)
	for i, s := range syms {
		assert.True(t, symRe.MatchString(s.Symbol), "row %d: Symbol %q 는 영문/숫자/점", i, s.Symbol)
		assert.NotEmpty(t, s.EnglishName, "row %d: EnglishName 비어있음", i)
		assert.NotNil(t, s.Raw, "row %d: Raw map 비어있지 않음", i)
	}

	// First row spot-check (catches column-index regression — sample data is committed bytes, won't rotate)
	first := syms[0]
	assert.Equal(t, "AACB", first.Symbol)
	assert.Equal(t, "NAS", first.MarketCode)
	assert.Equal(t, "USD", first.Currency)
	assert.Equal(t, "N", first.SuspendedYn)
	assert.Equal(t, "10.3800", first.BasePrice)
	assert.Equal(t, "2", first.StockType)
}

func TestParse_InvalidMarket(t *testing.T) {
	_, err := Parse("invalid", []byte{})
	assert.Error(t, err)
}

func TestParse_EmptyZip(t *testing.T) {
	_, err := Parse("nas", nil)
	assert.Error(t, err)
}
