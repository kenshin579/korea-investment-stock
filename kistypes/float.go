// Package kistypes 는 KIS API 응답 파싱을 위한 공용 타입을 제공한다.
package kistypes

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// Float 는 KIS 응답의 부호(+/-) 붙은 숫자 문자열과 빈 문자열을 안전하게 파싱하는 float64.
//
// 표준 encoding/json 의 `,string` 태그는 따옴표 안 값이 JSON number 문법이길 요구하는데,
// JSON number 는 leading '+' 를 금지한다. KIS 해외 API 는 등락 필드를 "+1.26" 형태로 주므로
// 그 태그로는 상승 값에서 파싱이 실패한다. 이 타입이 그 대체다.
type Float float64

// UnmarshalJSON 은 다음을 허용한다:
//   - "null" / 빈 입력 → 0
//   - 따옴표로 감싼 숫자 문자열(KIS 기본): "+1.26", "-1.26", "1.26", ""(→0)
//   - 따옴표 없는 JSON number: 1.26
func (f *Float) UnmarshalJSON(b []byte) error {
	s := strings.TrimSpace(string(b))
	if s == "" || s == "null" {
		*f = 0
		return nil
	}
	s = strings.TrimSpace(strings.Trim(s, `"`))
	if s == "" {
		*f = 0
		return nil
	}
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return err
	}
	if math.IsNaN(v) || math.IsInf(v, 0) {
		return fmt.Errorf("kistypes.Float: non-finite value %q", s)
	}
	*f = Float(v)
	return nil
}
