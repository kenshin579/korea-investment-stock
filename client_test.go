package kis

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewClient_RequiresAllParams(t *testing.T) {
	_, err := NewClient("", "s", "acc")
	assert.Error(t, err)
	_, err = NewClient("k", "", "acc")
	assert.Error(t, err)
	_, err = NewClient("k", "s", "")
	assert.Error(t, err)
}

func TestNewClient_AppliesOptions(t *testing.T) {
	c, err := NewClient("k", "s", "acc",
		WithBaseURL("https://x"),
		WithRetries(7),
		WithRateLimit(20),
	)
	require.NoError(t, err)
	assert.Equal(t, "https://x", c.opts.baseURL)
	assert.Equal(t, 7, c.opts.retries)
	assert.Equal(t, 20.0, c.opts.rateLimit)
	require.NotNil(t, c.Domestic)
	require.NotNil(t, c.Overseas)
}

func TestNewClient_Defaults(t *testing.T) {
	c, err := NewClient("k", "s", "acc")
	require.NoError(t, err)
	assert.Equal(t, RealEnv, c.opts.baseURL)
	assert.Equal(t, 3, c.opts.retries)
	assert.Equal(t, 15.0, c.opts.rateLimit)
}
