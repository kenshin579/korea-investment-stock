package token

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"golang.org/x/sync/singleflight"
)

// oauthPath 는 한국투자 OAuth 토큰 발급 엔드포인트.
const oauthPath = "/oauth2/tokenP"

// seoulTZ 는 한투 응답의 만료 시각이 사용하는 시간대.
var seoulTZ *time.Location

func init() {
	loc, err := time.LoadLocation("Asia/Seoul")
	if err != nil {
		// fallback (KST = UTC+9)
		loc = time.FixedZone("KST", 9*3600)
	}
	seoulTZ = loc
}

// Config 는 Manager 생성 옵션.
type Config struct {
	Storage    Storage
	BaseURL    string
	APIKey     string
	APISecret  string
	HTTPClient *http.Client
}

// Manager 는 OAuth 토큰의 발급/캐시/갱신을 담당.
// 동시 호출 시 한 번만 발급 (singleflight).
type Manager struct {
	cfg    Config
	mu     sync.RWMutex
	cached *AccessToken
	flight singleflight.Group
}

// NewManager 는 Manager 생성.
func NewManager(cfg Config) *Manager {
	if cfg.HTTPClient == nil {
		cfg.HTTPClient = http.DefaultClient
	}
	return &Manager{cfg: cfg}
}

// Get 은 유효한 토큰의 Bearer 문자열 반환.
// 캐시된 토큰이 없거나 만료 임박이면 자동 발급.
//
// 주의 (singleflight 동작):
// 동시 호출 시 첫 caller 의 ctx 가 inflight 발급 동안 사용됨. 첫 caller 가
// ctx cancel 하면 다른 waiter 들도 모두 cancellation 에러 받음. token 발급은
// 일반적으로 짧은 작업이라 caller 가 적절한 timeout 을 주면 문제 없음.
func (m *Manager) Get(ctx context.Context) (string, error) {
	if t := m.cachedValid(); t != nil {
		return t.Bearer(), nil
	}

	// storage 에서 로드 시도
	if t, err := m.cfg.Storage.Load(ctx); err == nil && t != nil && !t.IsExpired() {
		m.mu.Lock()
		m.cached = t
		m.mu.Unlock()
		return t.Bearer(), nil
	}

	// 새로 발급 (singleflight 로 동시 호출 1번만)
	v, err, _ := m.flight.Do("issue", func() (interface{}, error) {
		return m.issue(ctx)
	})
	if err != nil {
		return "", err
	}
	return v.(*AccessToken).Bearer(), nil
}

// Refresh 는 캐시 무시하고 강제로 새 토큰 발급.
func (m *Manager) Refresh(ctx context.Context) (string, error) {
	v, err, _ := m.flight.Do("refresh", func() (interface{}, error) {
		return m.issue(ctx)
	})
	if err != nil {
		return "", err
	}
	return v.(*AccessToken).Bearer(), nil
}

func (m *Manager) cachedValid() *AccessToken {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.cached == nil || m.cached.IsExpired() {
		return nil
	}
	return m.cached
}

type oauthResp struct {
	AccessToken          string `json:"access_token"`
	TokenType            string `json:"token_type"`
	ExpiresIn            int    `json:"expires_in"`
	AccessTokenExpiredAt string `json:"access_token_token_expired"`
}

func (m *Manager) issue(ctx context.Context) (*AccessToken, error) {
	body, _ := json.Marshal(map[string]string{
		"grant_type": "client_credentials",
		"appkey":     m.cfg.APIKey,
		"appsecret":  m.cfg.APISecret,
	})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		strings.TrimRight(m.cfg.BaseURL, "/")+oauthPath,
		bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := m.cfg.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("kis: token issue: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("kis: token issue: HTTP %d", resp.StatusCode)
	}
	var oauth oauthResp
	if err := json.NewDecoder(resp.Body).Decode(&oauth); err != nil {
		return nil, fmt.Errorf("kis: token decode: %w", err)
	}

	expiresAt, err := parseExpiry(oauth)
	if err != nil {
		return nil, err
	}
	tok := &AccessToken{
		Value:     oauth.AccessToken,
		TokenType: oauth.TokenType,
		ExpiresAt: expiresAt,
	}

	if err := m.cfg.Storage.Save(ctx, tok); err != nil {
		// 저장 실패는 warning, 발급은 성공
	}
	m.mu.Lock()
	m.cached = tok
	m.mu.Unlock()
	return tok, nil
}

// parseExpiry 는 한투 응답의 만료 시각 파싱.
// 우선순위: access_token_token_expired (Asia/Seoul) → expires_in (now + delta).
func parseExpiry(r oauthResp) (time.Time, error) {
	if r.AccessTokenExpiredAt != "" {
		t, err := time.ParseInLocation("2006-01-02 15:04:05", r.AccessTokenExpiredAt, seoulTZ)
		if err == nil {
			return t, nil
		}
	}
	if r.ExpiresIn > 0 {
		return time.Now().Add(time.Duration(r.ExpiresIn) * time.Second), nil
	}
	return time.Time{}, errors.New("kis: token response missing expiry")
}
