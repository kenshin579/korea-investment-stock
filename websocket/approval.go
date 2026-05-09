package websocket

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

type approvalKeyManager struct {
	httpClient *http.Client
	baseURL    string
	appKey     string
	appSecret  string
	ttl        time.Duration

	mu     sync.Mutex
	cached string
	expiry time.Time
}

func newApprovalKeyManager(c *http.Client, baseURL, appKey, appSecret string, ttl time.Duration) *approvalKeyManager {
	if c == nil {
		c = http.DefaultClient
	}
	return &approvalKeyManager{
		httpClient: c,
		baseURL:    baseURL,
		appKey:     appKey,
		appSecret:  appSecret,
		ttl:        ttl,
	}
}

func (m *approvalKeyManager) Get(ctx context.Context) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.cached != "" && time.Now().Before(m.expiry) {
		return m.cached, nil
	}

	body, _ := json.Marshal(map[string]string{
		"grant_type": "client_credentials",
		"appkey":     m.appKey,
		"secretkey":  m.appSecret,
	})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, m.baseURL+"/oauth2/Approval", bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("%w: build req: %v", ErrWSApprovalFailed, err)
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	resp, err := m.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("%w: http: %v", ErrWSApprovalFailed, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		raw, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("%w: status %d body %s", ErrWSApprovalFailed, resp.StatusCode, string(raw))
	}

	var out struct {
		ApprovalKey string `json:"approval_key"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", fmt.Errorf("%w: decode: %v", ErrWSApprovalFailed, err)
	}
	if out.ApprovalKey == "" {
		return "", fmt.Errorf("%w: empty approval_key", ErrWSApprovalFailed)
	}

	m.cached = out.ApprovalKey
	m.expiry = time.Now().Add(m.ttl)
	return m.cached, nil
}
