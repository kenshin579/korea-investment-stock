package kistypes_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kenshin579/korea-investment-stock/kistypes"
)

func TestFloat_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    float64
		wantErr bool
	}{
		{"plus sign", `"+1.26"`, 1.26, false},
		{"minus sign", `"-1.26"`, -1.26, false},
		{"quoted plain", `"1.26"`, 1.26, false},
		{"unquoted number", `1.26`, 1.26, false},
		{"empty string", `""`, 0, false},
		{"null", `null`, 0, false},
		{"zero", `"0"`, 0, false},
		{"invalid", `"abc"`, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var f kistypes.Float
			err := json.Unmarshal([]byte(tt.input), &f)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.InDelta(t, tt.want, float64(f), 1e-9)
		})
	}
}

func TestFloat_StructField(t *testing.T) {
	var s struct {
		R kistypes.Float `json:"r"`
	}
	err := json.Unmarshal([]byte(`{"r":"+1.26"}`), &s)
	require.NoError(t, err)
	assert.InDelta(t, 1.26, float64(s.R), 1e-9)
}
