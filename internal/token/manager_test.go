package token

import (
	"context"
	"net/http"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// warmup 은 테스트 편의용 — 첫 토큰을 미리 받아둠.
func (m *Manager) warmup(ctx context.Context) error {
	_, err := m.Get(ctx)
	return err
}

func TestManager_Get_FreshFetch(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder(http.MethodPost, "=~/oauth2/tokenP",
		httpmock.NewStringResponder(200, `{
			"access_token": "fresh-token",
			"token_type": "Bearer",
			"expires_in": 86400,
			"access_token_token_expired": "2099-12-31 23:59:59"
		}`))

	storage := NewFileStorage(t.TempDir() + "/token.json")
	m := NewManager(Config{
		Storage:    storage,
		BaseURL:    "https://openapi.koreainvestment.com:9443",
		APIKey:     "k",
		APISecret:  "s",
		HTTPClient: http.DefaultClient,
	})
	bearer, err := m.Get(context.Background())
	require.NoError(t, err)
	assert.Equal(t, "Bearer fresh-token", bearer)
}

func TestManager_Get_UsesCache(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	calls := atomic.Int64{}
	httpmock.RegisterResponder(http.MethodPost, "=~/oauth2/tokenP",
		func(req *http.Request) (*http.Response, error) {
			calls.Add(1)
			return httpmock.NewStringResponse(200, `{
				"access_token": "cached",
				"token_type": "Bearer",
				"access_token_token_expired": "2099-12-31 23:59:59"
			}`), nil
		})

	storage := NewFileStorage(t.TempDir() + "/token.json")
	m := NewManager(Config{
		Storage: storage, BaseURL: "https://x", APIKey: "k", APISecret: "s",
		HTTPClient: http.DefaultClient,
	})
	for i := 0; i < 5; i++ {
		_, err := m.Get(context.Background())
		require.NoError(t, err)
	}
	assert.Equal(t, int64(1), calls.Load(), "cached token should be reused")
}

func TestManager_Refresh_Forces(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	calls := atomic.Int64{}
	httpmock.RegisterResponder(http.MethodPost, "=~/oauth2/tokenP",
		func(req *http.Request) (*http.Response, error) {
			calls.Add(1)
			return httpmock.NewStringResponse(200, `{
				"access_token": "refreshed",
				"token_type": "Bearer",
				"access_token_token_expired": "2099-12-31 23:59:59"
			}`), nil
		})

	storage := NewFileStorage(t.TempDir() + "/token.json")
	m := NewManager(Config{
		Storage: storage, BaseURL: "https://x", APIKey: "k", APISecret: "s",
		HTTPClient: http.DefaultClient,
	})
	require.NoError(t, m.warmup(context.Background()))
	_, err := m.Refresh(context.Background())
	require.NoError(t, err)
	assert.Equal(t, int64(2), calls.Load(), "warmup + Refresh = 2 calls")
}

func TestManager_Get_Singleflight(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	calls := atomic.Int64{}
	httpmock.RegisterResponder(http.MethodPost, "=~/oauth2/tokenP",
		func(req *http.Request) (*http.Response, error) {
			calls.Add(1)
			time.Sleep(50 * time.Millisecond) // 동시 발급 race window
			return httpmock.NewStringResponse(200, `{
				"access_token": "single",
				"token_type": "Bearer",
				"access_token_token_expired": "2099-12-31 23:59:59"
			}`), nil
		})

	storage := NewFileStorage(t.TempDir() + "/token.json")
	m := NewManager(Config{
		Storage: storage, BaseURL: "https://x", APIKey: "k", APISecret: "s",
		HTTPClient: http.DefaultClient,
	})

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _ = m.Get(context.Background())
		}()
	}
	wg.Wait()
	assert.Equal(t, int64(1), calls.Load(), "10 concurrent Get → 1 OAuth call (singleflight)")
}

func TestManager_Get_UsesStorageCache(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	calls := atomic.Int64{}
	httpmock.RegisterResponder(http.MethodPost, "=~/oauth2/tokenP",
		func(req *http.Request) (*http.Response, error) {
			calls.Add(1)
			return httpmock.NewStringResponse(200, `{
				"access_token": "should-not-be-used",
				"token_type": "Bearer",
				"access_token_token_expired": "2099-12-31 23:59:59"
			}`), nil
		})

	storage := NewFileStorage(t.TempDir() + "/token.json")
	// Storage 에 미리 유효한 토큰 저장
	require.NoError(t, storage.Save(context.Background(), &AccessToken{
		Value:     "from-storage",
		TokenType: "Bearer",
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}))

	m := NewManager(Config{
		Storage: storage, BaseURL: "https://x", APIKey: "k", APISecret: "s",
		HTTPClient: http.DefaultClient,
	})

	bearer, err := m.Get(context.Background())
	require.NoError(t, err)
	assert.Equal(t, "Bearer from-storage", bearer)
	assert.Equal(t, int64(0), calls.Load(), "should NOT call OAuth — token loaded from storage")
}
