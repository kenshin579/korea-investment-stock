package kis

import (
	"context"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
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

func TestNewClient_TokenIssue(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder(http.MethodPost, "=~/oauth2/tokenP",
		httpmock.NewStringResponder(200, `{
			"access_token": "T",
			"token_type": "Bearer",
			"access_token_token_expired": "2099-12-31 23:59:59"
		}`))

	c, err := NewClient("k", "s", "acc",
		WithHTTPClient(&http.Client{Transport: httpmock.DefaultTransport}),
	)
	require.NoError(t, err)
	tok, err := c.IssueAccessToken(context.Background())
	require.NoError(t, err)
	assert.Equal(t, "Bearer T", tok)
}
