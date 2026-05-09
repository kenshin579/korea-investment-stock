// Package kis is a Go client for the Korea Investment Securities OpenAPI.
//
// See docs/superpowers/specs/2026-05-03-korea-investment-go-migration-design.md
// for the design rationale, and docs/superpowers/specs/2026-05-03-phase1-api-coverage-design.md
// for Phase 1 scope.
package kis

import (
	"errors"
	"log/slog"
	"net/http"
	"path/filepath"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/kenshin579/korea-investment-stock/bonds"
	"github.com/kenshin579/korea-investment-stock/domestic"
	"github.com/kenshin579/korea-investment-stock/futures"
	"github.com/kenshin579/korea-investment-stock/internal/httpclient"
	"github.com/kenshin579/korea-investment-stock/internal/mastercache"
	"github.com/kenshin579/korea-investment-stock/internal/ratelimit"
	"github.com/kenshin579/korea-investment-stock/internal/token"
	"github.com/kenshin579/korea-investment-stock/overseas"
	"github.com/kenshin579/korea-investment-stock/websocket"
)

// KIS OpenAPI base URLs.
const (
	RealEnv  = "https://openapi.koreainvestment.com:9443"
	PaperEnv = "https://openapivts.koreainvestment.com:29443"
)

// Client 는 kis 라이브러리의 단일 진입점.
type Client struct {
	apiKey    string
	apiSecret string
	accountNo string
	opts      clientOptions

	httpClient *httpclient.Client
	tokenMgr   *token.Manager
	masterC    *mastercache.Cache // Phase 1.2 의 FetchKospi/Kosdaq Symbols 가 사용 예정

	Domestic *domestic.Client
	Overseas *overseas.Client
	Bonds    *bonds.Client
	Futures  *futures.Client
	WS       *websocket.Client
}

// Option 은 functional option.
type Option func(*clientOptions)

type clientOptions struct {
	baseURL        string
	retries        int
	rateLimit      float64
	httpClient     *http.Client
	tokenStorage   token.Storage
	masterCacheDir string
	logger         *slog.Logger
	timeout        time.Duration
	userAgent      string
}

// NewClient 는 kis Client 생성 (직접 credentials 전달).
func NewClient(apiKey, apiSecret, accountNo string, opts ...Option) (*Client, error) {
	if apiKey == "" || apiSecret == "" || accountNo == "" {
		return nil, errors.New("kis: apiKey, apiSecret, and accountNo are required and must not be empty")
	}

	cfg := defaultOptions()
	for _, opt := range opts {
		opt(&cfg)
	}

	c := &Client{
		apiKey:    apiKey,
		apiSecret: apiSecret,
		accountNo: accountNo,
		opts:      cfg,
	}
	if err := c.wireInfra(); err != nil {
		return nil, err
	}
	c.Domestic = domestic.New(c.httpClient, c.masterC)
	c.Overseas = overseas.New(c.httpClient, c.masterC)
	c.Bonds = bonds.New(c.httpClient)
	c.Futures = futures.New(c.httpClient)

	// WebSocket endpoint: 실전 vs 모의투자
	wsEndpoint := "ws://ops.koreainvestment.com:21000"
	if cfg.baseURL == PaperEnv {
		wsEndpoint = "ws://ops.koreainvestment.com:31000"
	}
	c.WS = websocket.NewClient(websocket.Options{
		Endpoint:  wsEndpoint,
		BaseURL:   cfg.baseURL,
		AppKey:    apiKey,
		AppSecret: apiSecret,
	})
	return c, nil
}

func defaultOptions() clientOptions {
	return clientOptions{
		baseURL:   RealEnv,
		retries:   3,
		rateLimit: 15,
		timeout:   30 * time.Second,
		userAgent: "korea-investment-stock-go",
	}
}

func (c *Client) wireInfra() error {
	storage, err := c.resolveTokenStorage()
	if err != nil {
		return err
	}

	c.tokenMgr = token.NewManager(token.Config{
		Storage:    storage,
		BaseURL:    c.opts.baseURL,
		APIKey:     c.apiKey,
		APISecret:  c.apiSecret,
		HTTPClient: c.opts.httpClient,
	})

	c.httpClient = httpclient.New(httpclient.Config{
		BaseURL:    c.opts.baseURL,
		AppKey:     c.apiKey,
		AppSecret:  c.apiSecret,
		AccountNo:  c.accountNo,
		Limiter:    ratelimit.New(c.opts.rateLimit),
		TokenMgr:   c.tokenMgr,
		Retries:    c.opts.retries,
		Timeout:    c.opts.timeout,
		UserAgent:  c.opts.userAgent,
		HTTPClient: c.opts.httpClient,
	})

	// masterDir resolution 은 lenient: DefaultDir 실패 시 빈 문자열 → mastercache.New 가
	// os.TempDir()/kis 로 자동 fallback. token storage 와의 strict-vs-lenient 차이는
	// 의도적 — 토큰은 secret 이라 fallback 위치 특정 필수, master cache 는 단순 다운로드 캐시.
	masterDir := c.opts.masterCacheDir
	if masterDir == "" {
		d, _ := mastercache.DefaultDir()
		masterDir = d
	}
	c.masterC = mastercache.New(masterDir, 7*24*time.Hour)
	return nil
}

func (c *Client) resolveTokenStorage() (token.Storage, error) {
	if c.opts.tokenStorage != nil {
		return c.opts.tokenStorage, nil
	}
	// default: FileStorage at user cache dir
	dir, err := mastercache.DefaultDir() // 같은 위치 재사용
	if err != nil {
		return nil, err
	}
	return token.NewFileStorage(filepath.Join(dir, "token.json")), nil
}

// newRedisStorage 는 NewClientFromEnv 에서 redis 옵션 시 호출.
// (현재 client.go 내부에서는 사용되지 않음 — from_env.go (Task 16) 가 호출)
func newRedisStorage(url, password string) (token.Storage, error) {
	opts, err := redis.ParseURL(url)
	if err != nil {
		return nil, err
	}
	if password != "" {
		opts.Password = password
	}
	rdb := redis.NewClient(opts)
	return token.NewRedisStorage(rdb, "kis:token:default"), nil
}
