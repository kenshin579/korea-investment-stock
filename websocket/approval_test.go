package websocket

import (
	"context"
	"errors"
	"net/http"
	"sync/atomic"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestApproval_FetchAndCache(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var calls atomic.Int32
	httpmock.RegisterResponder(http.MethodPost, `=~/oauth2/Approval`,
		func(req *http.Request) (*http.Response, error) {
			calls.Add(1)
			return httpmock.NewStringResponse(200, `{"approval_key":"key-12345"}`), nil
		},
	)

	m := newApprovalKeyManager(http.DefaultClient, "https://api.example", "appkey", "appsecret", 23*time.Hour)

	k1, err := m.Get(context.Background())
	require.NoError(t, err)
	assert.Equal(t, "key-12345", k1)

	// 캐시 동작 — 두 번째 호출은 HTTP 안 침
	k2, err := m.Get(context.Background())
	require.NoError(t, err)
	assert.Equal(t, "key-12345", k2)
	assert.Equal(t, int32(1), calls.Load())
}

func TestApproval_Expiry(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	keys := []string{"key-A", "key-B"}
	idx := 0
	httpmock.RegisterResponder(http.MethodPost, `=~/oauth2/Approval`,
		func(req *http.Request) (*http.Response, error) {
			defer func() { idx++ }()
			return httpmock.NewStringResponse(200, `{"approval_key":"`+keys[idx]+`"}`), nil
		},
	)

	// TTL 0 → 매번 갱신
	m := newApprovalKeyManager(http.DefaultClient, "https://api.example", "appkey", "appsecret", 0)
	k1, _ := m.Get(context.Background())
	k2, _ := m.Get(context.Background())
	assert.NotEqual(t, k1, k2)
}

func TestApproval_Failure(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(http.MethodPost, `=~/oauth2/Approval`,
		httpmock.NewStringResponder(500, `{"error":"server"}`),
	)

	m := newApprovalKeyManager(http.DefaultClient, "https://api.example", "appkey", "appsecret", 23*time.Hour)
	_, err := m.Get(context.Background())
	require.Error(t, err)
	assert.True(t, errors.Is(err, ErrWSApprovalFailed))
}
