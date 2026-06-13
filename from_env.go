package kis

import (
	"fmt"

	"github.com/kenshin579/korea-investment-stock/internal/token"
)

// storageFromConfig 는 Config 의 토큰 저장소 설정을 token.Storage 로 변환한다.
//   - redis: TOKEN_STORAGE=redis + REDIS_URL 지정 시.
//   - file:  TOKEN_FILE 지정 시 해당 경로의 FileStorage.
//   - 그 외: nil (NewClient 가 기본 FileStorage(~/.cache/kis/token.json) 사용).
func storageFromConfig(cfg *Config) (token.Storage, error) {
	if cfg.TokenStorage == "redis" && cfg.RedisURL != "" {
		s, err := newRedisStorage(cfg.RedisURL, cfg.RedisPassword)
		if err != nil {
			return nil, fmt.Errorf("kis: redis storage: %w", err)
		}
		return s, nil
	}
	if cfg.TokenFile != "" {
		return token.NewFileStorage(cfg.TokenFile), nil
	}
	return nil, nil
}

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

	s, err := storageFromConfig(cfg)
	if err != nil {
		return nil, err
	}
	if s != nil {
		baseOpts = append(baseOpts, WithTokenStorage(s))
	}

	allOpts := append(baseOpts, opts...)
	return NewClient(cfg.APIKey, cfg.APISecret, cfg.AccountNo, allOpts...)
}
