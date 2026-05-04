package domestic_test

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"

	"github.com/kenshin579/korea-investment-stock/domestic"
	"github.com/kenshin579/korea-investment-stock/internal/httpclient"
	"github.com/kenshin579/korea-investment-stock/internal/mastercache"
	"github.com/kenshin579/korea-investment-stock/internal/ratelimit"
)

const testBaseURL = "https://openapi.koreainvestment.com:9443"

// loadFixture 는 testdata/<name> 파일 byte 를 로드.
func loadFixture(t *testing.T, name string) []byte {
	t.Helper()
	b, err := os.ReadFile(filepath.Join("testdata", name))
	require.NoError(t, err)
	return b
}

// loadFixtureString 은 string 으로 로드 (httpmock.NewStringResponder 용).
func loadFixtureString(t *testing.T, name string) string {
	return string(loadFixture(t, name))
}

// stubTokenManager 는 httpclient.TokenManager 의 stub. 항상 "Bearer test" 반환.
type stubTokenManager struct{}

func (stubTokenManager) Get(ctx context.Context) (string, error)     { return "Bearer test", nil }
func (stubTokenManager) Refresh(ctx context.Context) (string, error) { return "Bearer test", nil }

// newTestClient 는 httpmock 활성 상태에서 사용할 domestic.Client 생성.
// 호출자는 httpmock.Activate() / httpmock.DeactivateAndReset() 직접 관리.
func newTestClient(t *testing.T) *domestic.Client {
	t.Helper()
	httpClient := &http.Client{Transport: httpmock.DefaultTransport}
	httpcli := httpclient.New(httpclient.Config{
		BaseURL:    testBaseURL,
		AppKey:     "test-key",
		AppSecret:  "test-secret",
		AccountNo:  "00000000-00",
		Limiter:    ratelimit.New(1000),
		TokenMgr:   stubTokenManager{},
		Retries:    0,
		Timeout:    5 * time.Second,
		HTTPClient: httpClient,
	})
	master := mastercache.New(t.TempDir(), time.Hour)
	return domestic.New(httpcli, master)
}
