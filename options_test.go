package kis

import (
	"log/slog"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestOptions_Apply(t *testing.T) {
	cfg := &clientOptions{}
	httpC := &http.Client{}
	logger := slog.Default()

	WithBaseURL("https://example.com")(cfg)
	WithRetries(5)(cfg)
	WithRateLimit(20)(cfg)
	WithHTTPClient(httpC)(cfg)
	WithLogger(logger)(cfg)
	WithTimeout(10 * time.Second)(cfg)
	WithUserAgent("test-ua")(cfg)
	WithMasterCacheDir("/tmp/kis")(cfg)

	assert.Equal(t, "https://example.com", cfg.baseURL)
	assert.Equal(t, 5, cfg.retries)
	assert.Equal(t, 20.0, cfg.rateLimit)
	assert.Same(t, httpC, cfg.httpClient)
	assert.Same(t, logger, cfg.logger)
	assert.Equal(t, 10*time.Second, cfg.timeout)
	assert.Equal(t, "test-ua", cfg.userAgent)
	assert.Equal(t, "/tmp/kis", cfg.masterCacheDir)
}

func TestWithPaperEnv(t *testing.T) {
	cfg := &clientOptions{}
	WithPaperEnv()(cfg)
	assert.Equal(t, PaperEnv, cfg.baseURL)
}
