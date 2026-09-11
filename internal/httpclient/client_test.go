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

// 실제 한투는 "유효하지 않은 token"(EGW00121) 을 HTTP 500 + JSON body 로 반환한다.
// send() 가 이를 단순 5xx 로 보고 blind retry 하면 refresh 기회를 잃는다.
// → 토큰 오류 body 면 1회 refresh 후 재시도해야 한다 (blind 3회 재시도 X).
func TestClient_Do_InvalidTokenHTTP500AutoRefresh(t *testing.T) {
	tm := &stubTokenMgr{bearer: "Bearer T"}
	c := newTestClient(t, tm)
	calls := atomic.Int64{}
	httpmock.RegisterResponder(http.MethodGet, "=~/inquire-price",
		func(req *http.Request) (*http.Response, error) {
			calls.Add(1)
			if req.Header.Get("Authorization") == "Bearer T-refreshed" {
				return httpmock.NewStringResponse(200, `{"rt_cd":"0","msg_cd":"OK","msg1":"ok"}`), nil
			}
			// 오래된 토큰 → 한투가 HTTP 500 + EGW00121
			return httpmock.NewStringResponse(500, `{"rt_cd":"1","msg_cd":"EGW00121","msg1":"유효하지 않은 token 입니다."}`), nil
		})

	resp, err := c.Do(context.Background(), &Request{
		Method: http.MethodGet, Path: "/inquire-price", TrID: "FHKST01010100",
	})
	require.NoError(t, err)
	assert.Equal(t, "0", resp.RtCode)
	assert.Equal(t, int64(2), calls.Load(), "500 invalid-token → refresh 후 1회 재시도 (blind 5xx 재시도 아님)")
	assert.GreaterOrEqual(t, tm.calls.Load(), int64(2), "token Get + Refresh 모두 호출")
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

func TestClient_Do_TrCont(t *testing.T) {
	tm := &stubTokenMgr{bearer: "Bearer t"}
	c := newTestClient(t, tm)
	httpmock.RegisterResponder(http.MethodGet, "=~/inquire-balance",
		func(req *http.Request) (*http.Response, error) {
			assert.Equal(t, "N", req.Header.Get("tr_cont"), "연속조회 헤더 전달")
			resp := httpmock.NewStringResponse(200, `{"rt_cd":"0","msg_cd":"OK","msg1":"ok","output1":[]}`)
			resp.Header.Set("tr_cont", "M")
			return resp, nil
		})

	resp, err := c.Do(context.Background(), &Request{
		Method: http.MethodGet, Path: "/inquire-balance", TrID: "TTTC8434R", TrCont: "N",
	})
	require.NoError(t, err)
	assert.Equal(t, "M", resp.TrCont, "응답 헤더 tr_cont 캡처")
}

func TestClient_Do_TrContEmptyNotSent(t *testing.T) {
	tm := &stubTokenMgr{bearer: "Bearer t"}
	c := newTestClient(t, tm)
	httpmock.RegisterResponder(http.MethodGet, "=~/inquire-balance",
		func(req *http.Request) (*http.Response, error) {
			// net/http 는 헤더 키를 canonical 형태(Tr_cont)로 저장한다
			_, has := req.Header["Tr_cont"]
			assert.False(t, has, "초기 조회는 tr_cont 헤더를 보내지 않는다")
			return httpmock.NewStringResponse(200, `{"rt_cd":"0","msg_cd":"OK","msg1":"ok"}`), nil
		})
	resp, err := c.Do(context.Background(), &Request{Method: http.MethodGet, Path: "/inquire-balance", TrID: "T"})
	require.NoError(t, err)
	assert.Equal(t, "", resp.TrCont)
}

func TestSplitAccountNo(t *testing.T) {
	cases := []struct {
		in, cano, prdt string
		wantErr        bool
	}{
		{"12345678-01", "12345678", "01", false},
		{"1234567801", "12345678", "01", false},
		{" 12345678-01 ", "12345678", "01", false},
		{"1234-01", "", "", true},
		{"12345678-1", "", "", true},
		{"", "", "", true},
		{"abcdefgh-01", "", "", true},
		{"1234567a-01", "", "", true},
		{"12345678-0a", "", "", true},
	}
	for _, tc := range cases {
		cano, prdt, err := SplitAccountNo(tc.in)
		if tc.wantErr {
			assert.Error(t, err, tc.in)
			continue
		}
		require.NoError(t, err, tc.in)
		assert.Equal(t, tc.cano, cano)
		assert.Equal(t, tc.prdt, prdt)
	}
}

func TestClient_Account(t *testing.T) {
	c := newTestClient(t, &stubTokenMgr{bearer: "b"})
	cano, prdt, err := c.Account()
	require.NoError(t, err)
	assert.Equal(t, "12345678", cano)
	assert.Equal(t, "01", prdt)
}

func TestClient_IsPaper(t *testing.T) {
	paper := New(Config{
		BaseURL:   "https://openapivts.koreainvestment.com:29443",
		AppKey:    "ak",
		AppSecret: "as",
		AccountNo: "12345678-01",
		Limiter:   ratelimit.New(1000),
		TokenMgr:  &stubTokenMgr{bearer: "b"},
	})
	assert.True(t, paper.IsPaper(), "openapivts 도메인은 모의투자")

	live := New(Config{
		BaseURL:   "https://openapi.koreainvestment.com:9443",
		AppKey:    "ak",
		AppSecret: "as",
		AccountNo: "12345678-01",
		Limiter:   ratelimit.New(1000),
		TokenMgr:  &stubTokenMgr{bearer: "b"},
	})
	assert.False(t, live.IsPaper(), "openapi 도메인은 실전")
}

func TestHasNext(t *testing.T) {
	assert.True(t, HasNext("F"))
	assert.True(t, HasNext("M"))
	assert.False(t, HasNext("D"))
	assert.False(t, HasNext("E"))
	assert.False(t, HasNext(""))
}

// 토큰 재발급(1회 재시도) 경로에서도 tr_cont 요청 헤더가 재전송되고, 최종 응답의 tr_cont
// 헤더가 정상적으로 캡처되는지 확인한다.
func TestClient_Do_TrContSurvivesTokenRefresh(t *testing.T) {
	tm := &stubTokenMgr{bearer: "Bearer T"}
	c := newTestClient(t, tm)
	calls := atomic.Int64{}
	httpmock.RegisterResponder(http.MethodGet, "=~/inquire-balance",
		func(req *http.Request) (*http.Response, error) {
			assert.Equal(t, "N", req.Header.Get("tr_cont"), "재시도 요청에도 tr_cont 전달")
			n := calls.Add(1)
			if n == 1 {
				return httpmock.NewStringResponse(200, `{"rt_cd":"1","msg_cd":"EGW00123","msg1":"기간이 만료된 token 입니다"}`), nil
			}
			resp := httpmock.NewStringResponse(200, `{"rt_cd":"0","msg_cd":"OK","msg1":"ok"}`)
			resp.Header.Set("tr_cont", "M")
			return resp, nil
		})

	resp, err := c.Do(context.Background(), &Request{
		Method: http.MethodGet, Path: "/inquire-balance", TrID: "TTTC8434R", TrCont: "N",
	})
	require.NoError(t, err)
	assert.Equal(t, "0", resp.RtCode)
	assert.Equal(t, "M", resp.TrCont, "재발급 후 최종 응답 tr_cont 캡처")
}
