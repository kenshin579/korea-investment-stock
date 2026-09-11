package kistypes_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kenshin579/korea-investment-stock/kistypes"
)

func TestInt_UnmarshalJSON(t *testing.T) {
	cases := []struct {
		name    string
		in      string
		want    int64
		wantErr bool
	}{
		{"양수 문자열", `"123"`, 123, false},
		{"플러스 부호", `"+123"`, 123, false},
		{"음수 문자열", `"-45"`, -45, false},
		{"빈 문자열", `""`, 0, false},
		{"null", `null`, 0, false},
		{"따옴표 없는 number", `678`, 678, false},
		{"숫자 아님", `"abc"`, 0, true},
		{"소수점 이하 0 (정수)", `"123.00"`, 123, false},
		{"음수 + 소수점 이하 0", `"-45.0"`, -45, false},
		{"플러스 부호 + 소수점 이하 0", `"+7.000"`, 7, false},
		{"소수점 이하 0 아님", `"123.45"`, 0, true},
		{"지수 표기 미허용", `"1e3"`, 0, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var v kistypes.Int
			err := json.Unmarshal([]byte(tc.in), &v)
			if tc.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.want, int64(v))
		})
	}
}
