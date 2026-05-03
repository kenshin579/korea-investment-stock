package kis

import "fmt"

// NewClientFromEnv 는 KOREA_INVESTMENT_* 환경변수에서 credentials/설정 자동 감지 후 Client 생성.
// 추가 옵션은 functional options 로 override 가능.
func NewClientFromEnv(opts ...Option) (*Client, error) {
	cfg, err := LoadConfigFromEnv()
	if err != nil {
		return nil, err
	}
	return newFromConfig(cfg, opts...)
}

func newFromConfig(cfg *Config, opts ...Option) (*Client, error) {
	// Config → options 로 변환 (옵션이 마지막에 override)
	baseOpts := []Option{}
	if cfg.BaseURL != "" {
		baseOpts = append(baseOpts, WithBaseURL(cfg.BaseURL))
	}
	if cfg.RateLimit > 0 {
		baseOpts = append(baseOpts, WithRateLimit(cfg.RateLimit))
	}
	if cfg.Retries > 0 {
		baseOpts = append(baseOpts, WithRetries(cfg.Retries))
	}
	if cfg.MasterCacheDir != "" {
		baseOpts = append(baseOpts, WithMasterCacheDir(cfg.MasterCacheDir))
	}

	if cfg.TokenStorage == "redis" && cfg.RedisURL != "" {
		s, err := newRedisStorage(cfg.RedisURL, cfg.RedisPassword)
		if err != nil {
			return nil, fmt.Errorf("kis: redis storage: %w", err)
		}
		baseOpts = append(baseOpts, WithTokenStorage(s))
	}

	allOpts := append(baseOpts, opts...)
	return NewClient(cfg.APIKey, cfg.APISecret, cfg.AccountNo, allOpts...)
}
