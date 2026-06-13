package kis

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kenshin579/korea-investment-stock/internal/token"
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

// TOKEN_FILE 이 지정되면 그 경로의 FileStorage 가 쓰여야 한다.
// (회귀: 과거 from_env 가 cfg.TokenFile 을 무시해 항상 기본 경로만 사용했다.)
func TestStorageFromConfig_TokenFileHonored(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "tok.json")

	s, err := storageFromConfig(&Config{TokenFile: path})
	require.NoError(t, err)
	require.NotNil(t, s)

	err = s.Save(context.Background(), &token.AccessToken{
		Value: "abc", TokenType: "Bearer", ExpiresAt: time.Now().Add(time.Hour),
	})
	require.NoError(t, err)

	_, statErr := os.Stat(path)
	require.NoError(t, statErr, "토큰이 설정한 TokenFile 경로에 저장되어야 함")
}

// token_file/redis 미지정이면 nil → NewClient 기본 파일 저장소에 위임.
func TestStorageFromConfig_DefaultNil(t *testing.T) {
	s, err := storageFromConfig(&Config{})
	require.NoError(t, err)
	require.Nil(t, s)
}
