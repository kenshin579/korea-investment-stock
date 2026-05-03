package kis

import (
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewClientFromYAML(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	require.NoError(t, os.WriteFile(path, []byte(`
api_key: yk
api_secret: ys
acc_no: "98765432-01"
`), 0600))

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder(http.MethodPost, "=~/oauth2/tokenP",
		httpmock.NewStringResponder(200, `{"access_token":"x","token_type":"Bearer","access_token_token_expired":"2099-12-31 23:59:59"}`))

	c, err := NewClientFromYAML(path,
		WithHTTPClient(&http.Client{Transport: httpmock.DefaultTransport}))
	require.NoError(t, err)
	assert.NotNil(t, c.Domestic)
}

func TestNewClientFromYAML_NotFound(t *testing.T) {
	_, err := NewClientFromYAML("/nonexistent.yaml")
	require.Error(t, err)
}
