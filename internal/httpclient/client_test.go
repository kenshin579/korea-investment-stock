package httpclient

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kenshin579/korea-investment-stock/internal/ratelimit"
	"github.com/kenshin579/korea-investment-stock/internal/token"
)

type stubTokenMgr struct {
	bearer string
	calls  atomic.Int64
}

func (s *stubTokenMgr) Get(ctx context.Context) (string, error) {
	s.calls.Add(1)
	return s.bearer, nil
}

func (s *stubTokenMgr) Refresh(ctx context.Context) (string, error) {
	s.calls.Add(1)
	return s.bearer + "-refreshed", nil
}

func newTestClient(t *testing.T, tm TokenManager) *Client {
	t.Helper()
	c := New(Config{
		BaseURL:   "https://openapi.test",
		AppKey:    "ak",
		AppSecret: "as",
		AccountNo: "12345678-01",
		Limiter:   ratelimit.New(1000),
		TokenMgr:  tm,
		Retries:   2,
	})
	httpmock.ActivateNonDefault(c.resty.GetClient())
	t.Cleanup(httpmock.DeactivateAndReset)
	return c
}

func TestClient_Do_Success(t *testing.T) {
	tm := &stubTokenMgr{bearer: "Bearer T"}
	c := newTestClient(t, tm)
	httpmock.RegisterResponder(http.MethodGet, "=~/inquire-price",
		httpmock.NewStringResponder(200, `{"rt_cd":"0","msg_cd":"OK","msg1":"ok","output":{"x":"1"}}`))

	resp, err := c.Do(context.Background(), &Request{
		Method: http.MethodGet,
		Path:   "/uapi/domestic-stock/v1/quotations/inquire-price",
		TrID:   "FHKST01010100",
	})
	require.NoError(t, err)
	assert.Equal(t, "0", resp.RtCode)
	assert.Equal(t, int64(1), tm.calls.Load(), "token Get called once")
}

func TestClient_Do_APIError(t *testing.T) {
	tm := &stubTokenMgr{bearer: "Bearer T"}
	c := newTestClient(t, tm)
	httpmock.RegisterResponder(http.MethodGet, "=~/inquire-price",
		httpmock.NewStringResponder(200, `{"rt_cd":"1","msg_cd":"MCA00001","msg1":"잘못된 종목"}`))

	_, err := c.Do(context.Background(), &Request{
		Method: http.MethodGet,
		Path:   "/inquire-price",
		TrID:   "FHKST01010100",
	})
	require.Error(t, err)
	var apiErr *APIError
	require.ErrorAs(t, err, &apiErr)
	assert.Equal(t, "MCA00001", apiErr.MsgCode)
}

func TestClient_Do_TokenExpiredAutoRetry(t *testing.T) {
	tm := &stubTokenMgr{bearer: "Bearer T"}
	c := newTestClient(t, tm)
	calls := atomic.Int64{}
	httpmock.RegisterResponder(http.MethodGet, "=~/inquire-price",
		func(req *http.Request) (*http.Response, error) {
			n := calls.Add(1)
			if n == 1 {
				return httpmock.NewStringResponse(200, `{"rt_cd":"1","msg_cd":"EGW00123","msg1":"기간이 만료된 token 입니다"}`), nil
			}
			return httpmock.NewStringResponse(200, `{"rt_cd":"0","msg_cd":"OK","msg1":"ok"}`), nil
		})

	resp, err := c.Do(context.Background(), &Request{
		Method: http.MethodGet,
		Path:   "/inquire-price",
		TrID:   "FHKST01010100",
	})
	require.NoError(t, err)
	assert.Equal(t, "0", resp.RtCode)
	assert.GreaterOrEqual(t, tm.calls.Load(), int64(2), "token Get + Refresh both called")
}

func TestClient_Do_Retry5xx(t *testing.T) {
	tm := &stubTokenMgr{bearer: "Bearer T"}
	c := newTestClient(t, tm)
	calls := atomic.Int64{}
	httpmock.RegisterResponder(http.MethodGet, "=~/inquire-price",
		func(req *http.Request) (*http.Response, error) {
			n := calls.Add(1)
			if n < 2 {
				return httpmock.NewStringResponse(503, `service unavailable`), nil
			}
			return httpmock.NewStringResponse(200, `{"rt_cd":"0","msg_cd":"OK","msg1":"ok"}`), nil
		})

	resp, err := c.Do(context.Background(), &Request{
		Method: http.MethodGet, Path: "/inquire-price", TrID: "FHKST01010100",
	})
	require.NoError(t, err)
	assert.Equal(t, "0", resp.RtCode)
	assert.Equal(t, int64(2), calls.Load(), "5xx → retry")
}

func TestClient_Do_TokenError(t *testing.T) {
	tm := &errorTokenMgr{err: errors.New("oauth down")}
	c := newTestClient(t, tm)
	_, err := c.Do(context.Background(), &Request{
		Method: http.MethodGet, Path: "/inquire-price", TrID: "FHKST01010100",
	})
	require.Error(t, err)
	assert.True(t, strings.Contains(err.Error(), "oauth down"))
}

type errorTokenMgr struct{ err error }

func (e *errorTokenMgr) Get(ctx context.Context) (string, error)     { return "", e.err }
func (e *errorTokenMgr) Refresh(ctx context.Context) (string, error) { return "", e.err }

var _ TokenManager = (*stubTokenMgr)(nil)
var _ TokenManager = (*errorTokenMgr)(nil)
var _ = token.AccessToken{} // import 활용
