package kis

import (
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config 는 라이브러리 설정. NewClientFromEnv / NewClientFromYAML 의 base.
type Config struct {
	APIKey         string  `yaml:"api_key"`
	APISecret      string  `yaml:"api_secret"`
	AccountNo      string  `yaml:"acc_no"`
	BaseURL        string  `yaml:"base_url"`
	TokenStorage   string  `yaml:"token_storage_type"` // "file" | "redis"
	TokenFile      string  `yaml:"token_file"`
	RedisURL       string  `yaml:"redis_url"`
	RedisPassword  string  `yaml:"redis_password"`
	MasterCacheDir string  `yaml:"master_cache_dir"`
	RateLimit      float64 `yaml:"rate_limit"`
	Retries        int     `yaml:"retries"`
}

// LoadConfigFromEnv 는 KOREA_INVESTMENT_* 환경변수에서 Config 로드.
// 필수: API_KEY, API_SECRET, ACCOUNT_NO.
func LoadConfigFromEnv() (*Config, error) {
	required := func(key string) (string, error) {
		v := os.Getenv(key)
		if v == "" {
			return "", fmt.Errorf("kis: env var %s is required", key)
		}
		return v, nil
	}

	apiKey, err := required("KOREA_INVESTMENT_API_KEY")
	if err != nil {
		return nil, err
	}
	apiSecret, err := required("KOREA_INVESTMENT_API_SECRET")
	if err != nil {
		return nil, err
	}
	accNo, err := required("KOREA_INVESTMENT_ACCOUNT_NO")
	if err != nil {
		return nil, err
	}

	return &Config{
		APIKey:        apiKey,
		APISecret:     apiSecret,
		AccountNo:     accNo,
		BaseURL:       os.Getenv("KOREA_INVESTMENT_BASE_URL"),
		TokenStorage:  os.Getenv("KOREA_INVESTMENT_TOKEN_STORAGE"),
		TokenFile:     os.Getenv("KOREA_INVESTMENT_TOKEN_FILE"),
		RedisURL:      os.Getenv("KOREA_INVESTMENT_REDIS_URL"),
		RedisPassword: os.Getenv("KOREA_INVESTMENT_REDIS_PASSWORD"),
	}, nil
}

// LoadConfigFromYAML 는 YAML 파일에서 Config 로드.
func LoadConfigFromYAML(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("kis: read config: %w", err)
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("kis: parse config: %w", err)
	}
	if cfg.APIKey == "" || cfg.APISecret == "" || cfg.AccountNo == "" {
		return nil, errors.New("kis: api_key, api_secret, acc_no are required in config")
	}
	return &cfg, nil
}
