// Package httpclient 은 한투 API 호출용 resty 래퍼.
//
// rate limit / token / 재시도 / 에러 정규화를 한 곳에서 처리.
// 사용자에게 노출되지 않는 internal 패키지.
package httpclient

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/kenshin579/korea-investment-stock/internal/ratelimit"
)

// TokenManager 는 token.Manager 의 인터페이스 추상화 (테스트 편의).
type TokenManager interface {
	Get(ctx context.Context) (string, error)
	Refresh(ctx context.Context) (string, error)
}

// Config 는 Client 생성 옵션.
type Config struct {
	BaseURL    string
	AppKey     string
	AppSecret  string
	AccountNo  string
	Limiter    *ratelimit.Limiter
	TokenMgr   TokenManager
	Retries    int
	Timeout    time.Duration
	UserAgent  string
	HTTPClient *http.Client
}

// Client 는 한투 API 호출 단일 진입점.
type Client struct {
	cfg   Config
	resty *resty.Client
}

// New 는 Config 로 Client 생성.
func New(cfg Config) *Client {
	if cfg.Retries < 0 {
		cfg.Retries = 0
	}
	if cfg.Timeout <= 0 {
		cfg.Timeout = 30 * time.Second
	}
	if cfg.UserAgent == "" {
		cfg.UserAgent = "korea-investment-stock-go"
	}
	// nil guard — 사용자가 Config 직접 만들 때 panic 방지.
	// Production code path 는 root kis.NewClient 가 모두 채워서 옴.
	if cfg.Limiter == nil {
		cfg.Limiter = ratelimit.New(15)
	}
	// TokenMgr 는 nil 이면 단언적으로 panic — 사용자에게 명확한 에러:
	// 토큰 매니저 없이는 한투 API 인증이 불가능하므로 default 도 의미 없음.
	// (root NewClient 는 항상 주입하므로 이 경로는 직접 New() 호출 시만 도달.)
	if cfg.TokenMgr == nil {
		panic("httpclient: TokenMgr is required (typically set by kis.NewClient)")
	}
	r := resty.New().
		SetBaseURL(strings.TrimRight(cfg.BaseURL, "/")).
		SetTimeout(cfg.Timeout).
		SetHeader("User-Agent", cfg.UserAgent)
	// Note: SetTransport only — resty manages its own Timeout (cfg.Timeout).
	// http.Client 의 Jar / CheckRedirect 는 의도적으로 forward 하지 않음.
	if cfg.HTTPClient != nil {
		r.SetTransport(cfg.HTTPClient.Transport)
	}
	return &Client{cfg: cfg, resty: r}
}

// Request 는 단일 한투 API 호출 요청.
type Request struct {
	Method string            // http.MethodGet 등
	Path   string            // "/uapi/..." (BaseURL 제외)
	TrID   string            // tr_id 헤더 (한투 transaction ID)
	Query  map[string]string // GET 쿼리 파라미터
	Body   any               // POST body (JSON 직렬화)
	// CustType: P (개인), B (법인). 빈 문자열이면 미지정.
	CustType string
}

// Response 는 한투 API 응답을 정규화한 결과.
type Response struct {
	RtCode  string          `json:"rt_cd"`
	MsgCode string          `json:"msg_cd"`
	Msg1    string          `json:"msg1"`
	Output  json.RawMessage `json:"output"`
	Output1 json.RawMessage `json:"output1"`
	Output2 json.RawMessage `json:"output2"`
	Raw     []byte          `json:"-"`
}

// APIError 는 한투 응답의 rt_cd != "0" 케이스 — internal 패키지 전용.
// error.Error() 에 msg_cd / msg1 이 포함되어 호출자가 메시지로 구분 가능.
type APIError struct {
	RtCode  string
	MsgCode string
	Message string
	TrID    string
}

func (e *APIError) Error() string {
	return "kis: API error [" + e.MsgCode + "] " + e.Message
}

// Do 는 단일 호출 + 재시도 + 토큰 만료 자동 재발급.
func (c *Client) Do(ctx context.Context, req *Request) (*Response, error) {
	if err := c.cfg.Limiter.Wait(ctx); err != nil {
		return nil, err
	}

	bearer, err := c.cfg.TokenMgr.Get(ctx)
	if err != nil {
		return nil, err
	}

	resp, err := c.send(ctx, req, bearer)
	if err != nil {
		return nil, err
	}

	// 토큰 만료 → 1회 재발급 후 재시도
	if resp != nil && isTokenExpired(resp) {
		newBearer, refErr := c.cfg.TokenMgr.Refresh(ctx)
		if refErr != nil {
			return nil, refErr
		}
		resp, err = c.send(ctx, req, newBearer)
		if err != nil {
			return nil, err
		}
		if resp == nil {
			return nil, errors.New("kis: nil response after refresh")
		}
	}

	if resp.RtCode != "0" {
		return nil, &APIError{RtCode: resp.RtCode, MsgCode: resp.MsgCode, Message: resp.Msg1, TrID: req.TrID}
	}
	return resp, nil
}

func (c *Client) send(ctx context.Context, req *Request, bearer string) (*Response, error) {
	var lastHTTPErr error
	for attempt := 0; attempt <= c.cfg.Retries; attempt++ {
		r := c.resty.R().
			SetContext(ctx).
			SetHeader("Authorization", bearer).
			SetHeader("appkey", c.cfg.AppKey).
			SetHeader("appsecret", c.cfg.AppSecret).
			SetHeader("tr_id", req.TrID).
			SetHeader("Content-Type", "application/json; charset=utf-8")
		if req.CustType != "" {
			r.SetHeader("custtype", req.CustType)
		}
		if len(req.Query) > 0 {
			r.SetQueryParams(req.Query)
		}
		if req.Body != nil {
			r.SetBody(req.Body)
		}

		httpResp, err := r.Execute(req.Method, req.Path)
		if err != nil {
			lastHTTPErr = err
			if attempt == c.cfg.Retries {
				return nil, fmt.Errorf("kis: http: %w", err)
			}
			timer := time.NewTimer(backoff(attempt))
			select {
			case <-timer.C:
			case <-ctx.Done():
				timer.Stop()
				return nil, ctx.Err()
			}
			continue
		}

		if httpResp.StatusCode() >= 500 || httpResp.StatusCode() == http.StatusTooManyRequests {
			lastHTTPErr = fmt.Errorf("HTTP %d", httpResp.StatusCode())
			if attempt == c.cfg.Retries {
				return nil, fmt.Errorf("kis: http: %s after %d retries", lastHTTPErr, c.cfg.Retries)
			}
			timer := time.NewTimer(backoff(attempt))
			select {
			case <-timer.C:
			case <-ctx.Done():
				timer.Stop()
				return nil, ctx.Err()
			}
			continue
		}

		raw := httpResp.Body()
		var resp Response
		if err := json.Unmarshal(raw, &resp); err != nil {
			return nil, fmt.Errorf("kis: parse: %w (body=%s)", err, string(raw))
		}
		resp.Raw = raw
		return &resp, nil
	}
	return nil, errors.New("unreachable")
}

func isTokenExpired(r *Response) bool {
	// 한투는 만료 시 msg_cd 가 EGW00123 또는 메시지에 "기간이 만료된 token" 포함.
	return r.MsgCode == "EGW00123" || strings.Contains(r.Msg1, "기간이 만료된 token")
}

func backoff(attempt int) time.Duration {
	// 0.5s, 1s, 2s, ...
	d := 500 * time.Millisecond
	for i := 0; i < attempt; i++ {
		d *= 2
	}
	return d
}
