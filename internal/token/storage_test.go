package token

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAccessToken_IsExpired(t *testing.T) {
	cases := []struct {
		name     string
		expires  time.Time
		expected bool
	}{
		{"expired", time.Now().Add(-1 * time.Hour), true},
		{"about to expire (within margin)", time.Now().Add(2 * time.Minute), true},
		{"valid", time.Now().Add(1 * time.Hour), false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tok := &AccessToken{ExpiresAt: tc.expires}
			assert.Equal(t, tc.expected, tok.IsExpired())
		})
	}
}

func TestAccessToken_Bearer(t *testing.T) {
	tok := &AccessToken{Value: "abc123"}
	assert.Equal(t, "Bearer abc123", tok.Bearer())
}
