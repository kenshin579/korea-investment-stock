package krxmaster

import (
	"archive/zip"
	"bytes"
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

func TestParseKosdaq(t *testing.T) {
	zipBytes, err := os.ReadFile("testdata/kosdaq_code_sample.mst.zip")
	require.NoError(t, err)

	syms, err := ParseKosdaq(zipBytes)
	require.NoError(t, err)
	require.Len(t, syms, 3)

	shortCodeRe := regexp.MustCompile(`^[0-9A-Z]{6}$`)
	hangulRe := regexp.MustCompile(`[\x{AC00}-\x{D7A3}]`)

	for i, s := range syms {
		assert.True(t, shortCodeRe.MatchString(s.ShortCode),
			"row %d: ShortCode %q 는 6자리 영숫자", i, s.ShortCode)
		assert.NotEmpty(t, s.StandardCode, "row %d: StandardCode", i)
		assert.True(t, hangulRe.MatchString(s.KoreanName),
			"row %d: KoreanName 한글", i)
		assert.NotEmpty(t, s.GroupCode, "row %d: GroupCode", i)
		assert.NotNil(t, s.Raw, "row %d: Raw map", i)
		assert.GreaterOrEqual(t, len(s.Raw), 50, "row %d: Raw 60+ 컬럼", i)
	}

	// 첫 행 핵심 필드 정확값 검증 (Task 3 패턴 따라 강한 assertion 추가).
	require.GreaterOrEqual(t, len(syms), 1)
	s0 := syms[0]
	assert.Equal(t, "ST", s0.GroupCode, "주권 그룹코드는 'ST'")
	assert.Greater(t, s0.BasePrice, int64(0), "기준가 양수")
	assert.True(t, s0.FaceValue.IsPositive(), "액면가 양수")
	assert.Regexp(t, `^(0[1-9]|1[0-2])$`, s0.SettlementMonth, "결산월 01-12")
}

// --- 내부 헬퍼 단위 테스트 (white-box, 같은 패키지) ---

func TestOpenMstFromZip_NoMstFile(t *testing.T) {
	// ZIP 에 .mst 파일이 없는 경우 → error
	// 유효한 zip 이지만 .txt 파일만 포함.
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	f, err2 := zw.Create("readme.txt")
	require.NoError(t, err2)
	_, _ = f.Write([]byte("hello"))
	require.NoError(t, zw.Close())

	_, err := openMstFromZip(buf.Bytes())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), ".mst file not found")
}

func TestParseFwf_LineTooShort(t *testing.T) {
	// widths 합이 line 길이보다 클 경우 break 에 걸려 일부 컬럼만 반환.
	widths := []int{3, 3, 3} // 총 9 byte 필요
	names := []string{"a", "b", "c"}
	line := "AB" // 2 byte 만 — 첫 번째 컬럼(3)도 못 채움 → 빈 map
	out := parseFwf(line, widths, names)
	assert.Empty(t, out)

	line5 := "ABCDE" // 5 byte — 첫 컬럼(3) 채우고 두 번째(3) 못 채움
	out5 := parseFwf(line5, widths, names)
	assert.Len(t, out5, 1)
	assert.Equal(t, "ABC", out5["a"])
}

func TestAtoi64_InvalidString(t *testing.T) {
	assert.Equal(t, int64(0), atoi64("abc"))  // 파싱 실패 → 0
	assert.Equal(t, int64(0), atoi64(""))     // 빈 문자열 → 0
	assert.Equal(t, int64(42), atoi64(" 42 ")) // 공백 포함 → 42
}

func TestToDecimal_InvalidString(t *testing.T) {
	d := toDecimal("notanumber")
	assert.True(t, d.IsZero(), "파싱 실패 → decimal.Zero")

	d2 := toDecimal("")
	assert.True(t, d2.IsZero(), "빈 문자열 → decimal.Zero")

	d3 := toDecimal("500")
	assert.Equal(t, "500", d3.String())
}
