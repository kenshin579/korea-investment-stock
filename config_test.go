package kis

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadConfigFromEnv_Required(t *testing.T) {
	t.Setenv("KOREA_INVESTMENT_API_KEY", "k")
	t.Setenv("KOREA_INVESTMENT_API_SECRET", "s")
	t.Setenv("KOREA_INVESTMENT_ACCOUNT_NO", "12345678-01")

	cfg, err := LoadConfigFromEnv()
	require.NoError(t, err)
	assert.Equal(t, "k", cfg.APIKey)
	assert.Equal(t, "s", cfg.APISecret)
	assert.Equal(t, "12345678-01", cfg.AccountNo)
}

func TestLoadConfigFromEnv_MissingKey(t *testing.T) {
	t.Setenv("KOREA_INVESTMENT_API_KEY", "")
	t.Setenv("KOREA_INVESTMENT_API_SECRET", "s")
	t.Setenv("KOREA_INVESTMENT_ACCOUNT_NO", "x")
	_, err := LoadConfigFromEnv()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "KOREA_INVESTMENT_API_KEY")
}

func TestLoadConfigFromEnv_Optional(t *testing.T) {
	t.Setenv("KOREA_INVESTMENT_API_KEY", "k")
	t.Setenv("KOREA_INVESTMENT_API_SECRET", "s")
	t.Setenv("KOREA_INVESTMENT_ACCOUNT_NO", "x")
	t.Setenv("KOREA_INVESTMENT_BASE_URL", "https://custom")
	t.Setenv("KOREA_INVESTMENT_TOKEN_STORAGE", "redis")
	t.Setenv("KOREA_INVESTMENT_REDIS_URL", "redis://1.2.3.4:6379/0")

	cfg, err := LoadConfigFromEnv()
	require.NoError(t, err)
	assert.Equal(t, "https://custom", cfg.BaseURL)
	assert.Equal(t, "redis", cfg.TokenStorage)
	assert.Equal(t, "redis://1.2.3.4:6379/0", cfg.RedisURL)
}

func TestLoadConfigFromYAML(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	require.NoError(t, os.WriteFile(path, []byte(`
api_key: yk
api_secret: ys
acc_no: "98765432-01"
base_url: https://yaml-base
token_storage_type: file
token_file: /tmp/yaml-token.json
rate_limit: 12.5
retries: 5
`), 0600))

	cfg, err := LoadConfigFromYAML(path)
	require.NoError(t, err)
	assert.Equal(t, "yk", cfg.APIKey)
	assert.Equal(t, "ys", cfg.APISecret)
	assert.Equal(t, "98765432-01", cfg.AccountNo)
	assert.Equal(t, "https://yaml-base", cfg.BaseURL)
	assert.Equal(t, "file", cfg.TokenStorage)
	assert.Equal(t, "/tmp/yaml-token.json", cfg.TokenFile)
	assert.Equal(t, 12.5, cfg.RateLimit)
	assert.Equal(t, 5, cfg.Retries)
}

func TestLoadConfigFromYAML_FileNotFound(t *testing.T) {
	_, err := LoadConfigFromYAML("/nonexistent/file.yaml")
	require.Error(t, err)
}
