package kistypes

import (
	"strconv"
	"strings"
)

// Int 는 KIS 응답의 부호(+/-) 붙은 정수 문자열과 빈 문자열을 안전하게 파싱하는 int64.
//
// 표준 encoding/json 의 `,string` 태그는 빈 문자열("")과 leading '+' 에서 실패한다.
// KIS 는 국채 등 일부 응답의 정수 필드를 빈 문자열로 내려주므로 이 타입이 그 대체다.
// Float 의 정수판.
type Int int64

// UnmarshalJSON 은 다음을 허용한다:
//   - "null" / 빈 입력 → 0
//   - 따옴표로 감싼 정수 문자열: "+123", "-45", "123", ""(→0)
//   - 따옴표 없는 JSON number: 678
func (i *Int) UnmarshalJSON(b []byte) error {
	s := strings.TrimSpace(string(b))
	if s == "" || s == "null" {
		*i = 0
		return nil
	}
	s = strings.TrimSpace(strings.Trim(s, `"`))
	if s == "" {
		*i = 0
		return nil
	}
	s = strings.TrimPrefix(s, "+")
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return err
	}
	*i = Int(v)
	return nil
}
