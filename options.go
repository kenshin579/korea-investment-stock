package kis

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/kenshin579/korea-investment-stock/internal/token"
)

// WithBaseURL 은 한투 API base URL 을 지정 (default: RealEnv).
func WithBaseURL(url string) Option {
	return func(o *clientOptions) { o.baseURL = url }
}

// WithPaperEnv 는 모의투자 API 엔드포인트로 변경.
func WithPaperEnv() Option {
	return func(o *clientOptions) { o.baseURL = PaperEnv }
}

// WithRetries 는 5xx/429 재시도 횟수 (default: 3).
func WithRetries(n int) Option {
	return func(o *clientOptions) { o.retries = n }
}

// WithRateLimit 은 호출/초 한도 (default: 15).
func WithRateLimit(rps float64) Option {
	return func(o *clientOptions) { o.rateLimit = rps }
}

// WithHTTPClient 는 사용자 정의 *http.Client 주입 (custom transport, proxy 등).
func WithHTTPClient(c *http.Client) Option {
	return func(o *clientOptions) { o.httpClient = c }
}

// WithTokenStorage 는 사용자 정의 토큰 저장소 (default: FileStorage at ~/.cache/kis/token.json).
func WithTokenStorage(s token.Storage) Option {
	return func(o *clientOptions) { o.tokenStorage = s }
}

// WithMasterCacheDir 는 KOSPI/KOSDAQ 마스터 파일 캐시 디렉터리.
func WithMasterCacheDir(dir string) Option {
	return func(o *clientOptions) { o.masterCacheDir = dir }
}

// WithLogger 는 사용자 정의 slog logger.
func WithLogger(l *slog.Logger) Option {
	return func(o *clientOptions) { o.logger = l }
}

// WithTimeout 은 단일 HTTP 호출의 timeout (default: 30s).
func WithTimeout(d time.Duration) Option {
	return func(o *clientOptions) { o.timeout = d }
}

// WithUserAgent 는 User-Agent 헤더 (default: "korea-investment-stock-go").
func WithUserAgent(ua string) Option {
	return func(o *clientOptions) { o.userAgent = ua }
}
