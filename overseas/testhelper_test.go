package overseas_test

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"

	"github.com/kenshin579/korea-investment-stock/internal/httpclient"
	"github.com/kenshin579/korea-investment-stock/internal/mastercache"
	"github.com/kenshin579/korea-investment-stock/internal/ratelimit"
	"github.com/kenshin579/korea-investment-stock/overseas"
)

const testBaseURL = "https://openapi.koreainvestment.com:9443"

func loadFixture(t *testing.T, name string) []byte {
	t.Helper()
	b, err := os.ReadFile(filepath.Join("testdata", name))
	require.NoError(t, err)
	return b
}

func loadFixtureString(t *testing.T, name string) string {
	return string(loadFixture(t, name))
}

type stubTokenManager struct{}

func (stubTokenManager) Get(ctx context.Context) (string, error)     { return "Bearer test", nil }
func (stubTokenManager) Refresh(ctx context.Context) (string, error) { return "Bearer test", nil }

func newTestClient(t *testing.T) *overseas.Client {
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
	return overseas.New(httpcli, master)
}
