package krxmaster

import (
	"os"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseKospi(t *testing.T) {
	zipBytes, err := os.ReadFile("testdata/kospi_code_sample.mst.zip")
	require.NoError(t, err)

	syms, err := ParseKospi(zipBytes)
	require.NoError(t, err)
	require.Len(t, syms, 3, "sample 은 첫 3 행만 포함")

	// KRX 데이터는 시간에 따라 변동. 종목코드/한글명을 strict 비교 대신 패턴 검증.
	shortCodeRe := regexp.MustCompile(`^[0-9A-Z]{6}$`)
	hangulRe := regexp.MustCompile(`[\x{AC00}-\x{D7A3}]`)

	for i, s := range syms {
		assert.True(t, shortCodeRe.MatchString(s.ShortCode),
			"row %d: ShortCode %q 는 6자리 영숫자", i, s.ShortCode)
		assert.NotEmpty(t, s.StandardCode, "row %d: StandardCode 비어있음", i)
		assert.True(t, hangulRe.MatchString(s.KoreanName),
			"row %d: KoreanName %q 에 한글 포함", i, s.KoreanName)
		assert.NotEmpty(t, s.GroupCode, "row %d: GroupCode 비어있음", i)
		assert.NotNil(t, s.Raw, "row %d: Raw map 비어있지 않음", i)
		assert.GreaterOrEqual(t, len(s.Raw), 60,
			"row %d: Raw 에 ~70 컬럼 (최소 60+)", i)
	}

	// 첫 행 (000020 동화약품) 의 핵심 필드 — 정확한 값 검증.
	// testdata 는 commit 된 binary, 변동 없음.
	require.GreaterOrEqual(t, len(syms), 1)
	s0 := syms[0]
	assert.Equal(t, "000020", s0.ShortCode)
	assert.Equal(t, "ST", s0.GroupCode, "주권 그룹코드는 'ST'")
	assert.Greater(t, s0.BasePrice, int64(0), "기준가 양수")
	assert.True(t, s0.FaceValue.IsPositive(), "액면가 양수")
	assert.Regexp(t, `^(0[1-9]|1[0-2])$`, s0.SettlementMonth, "결산월 01-12")
}

func TestParseKospi_InvalidZip(t *testing.T) {
	_, err := ParseKospi([]byte("not a zip"))
	assert.Error(t, err)
}

func TestParseKospi_EmptyZip(t *testing.T) {
	_, err := ParseKospi(nil)
	assert.Error(t, err)
}
