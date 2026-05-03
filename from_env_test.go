package kis

import (
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewClientFromEnv_Success(t *testing.T) {
	t.Setenv("KOREA_INVESTMENT_API_KEY", "k")
	t.Setenv("KOREA_INVESTMENT_API_SECRET", "s")
	t.Setenv("KOREA_INVESTMENT_ACCOUNT_NO", "12345678-01")

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder(http.MethodPost, "=~/oauth2/tokenP",
		httpmock.NewStringResponder(200, `{"access_token":"x","token_type":"Bearer","access_token_token_expired":"2099-12-31 23:59:59"}`))

	c, err := NewClientFromEnv(WithHTTPClient(&http.Client{Transport: httpmock.DefaultTransport}))
	require.NoError(t, err)
	assert.NotNil(t, c.Domestic)
}

func TestNewClientFromEnv_MissingEnv(t *testing.T) {
	t.Setenv("KOREA_INVESTMENT_API_KEY", "")
	t.Setenv("KOREA_INVESTMENT_API_SECRET", "s")
	t.Setenv("KOREA_INVESTMENT_ACCOUNT_NO", "x")
	_, err := NewClientFromEnv()
	require.Error(t, err)
}
