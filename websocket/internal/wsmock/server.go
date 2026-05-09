// Package wsmock 은 KIS WebSocket 서버를 모방하는 local mock.
//
// 사용:
//
//	srv := wsmock.New(t)
//	defer srv.Close()
//	// srv.URL() 을 websocket.Options.Endpoint 로 사용
//
//	srv.SendRealtime(ctx, "H0STCNT0", "005930^123929^73100^...")
//	srv.SendText(ctx, `{"header":{"tr_id":"H0STCNT0"},"body":{"rt_cd":"0","msg_cd":"OPSP0000","msg1":"SUBSCRIBE SUCCESS"}}`)
package wsmock

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	wslib "github.com/coder/websocket"
)

type Server struct {
	t       *testing.T
	hs      *httptest.Server
	mu      sync.Mutex
	conn    *wslib.Conn
	connCh  chan struct{} // 클라이언트 connect 알림
	receive chan string   // 클라이언트 → mock 송신 메시지
}

func New(t *testing.T) *Server {
	t.Helper()
	s := &Server{
		t:       t,
		connCh:  make(chan struct{}, 16),
		receive: make(chan string, 100),
	}
	s.hs = httptest.NewServer(http.HandlerFunc(s.handle))
	return s
}

func (s *Server) URL() string {
	// httptest 는 http:// 반환 — ws:// 로 변환
	u := s.hs.URL
	return "ws" + u[len("http"):]
}

func (s *Server) Close() {
	s.mu.Lock()
	if s.conn != nil {
		_ = s.conn.Close(wslib.StatusNormalClosure, "test end")
	}
	s.mu.Unlock()
	s.hs.Close()
}

func (s *Server) handle(w http.ResponseWriter, r *http.Request) {
	c, err := wslib.Accept(w, r, &wslib.AcceptOptions{
		InsecureSkipVerify: true, // httptest origin 허용
	})
	if err != nil {
		s.t.Logf("wsmock accept: %v", err)
		return
	}
	s.mu.Lock()
	s.conn = c
	s.mu.Unlock()
	select {
	case s.connCh <- struct{}{}:
	default:
	}
	defer func() {
		_ = c.Close(wslib.StatusNormalClosure, "")
	}()

	// 클라이언트 → mock 메시지 read loop
	for {
		_, raw, err := c.Read(r.Context())
		if err != nil {
			return
		}
		select {
		case s.receive <- string(raw):
		default:
		}
	}
}

// WaitConnected 는 클라이언트 connect 까지 blocking.
func (s *Server) WaitConnected(ctx context.Context) error {
	select {
	case <-s.connCh:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// SendText 는 mock → 클라이언트 raw text frame 송신.
func (s *Server) SendText(ctx context.Context, msg string) error {
	s.mu.Lock()
	c := s.conn
	s.mu.Unlock()
	if c == nil {
		return http.ErrServerClosed
	}
	return c.Write(ctx, wslib.MessageText, []byte(msg))
}

// SendRealtime — 0|TR_ID|001|payload 형식 송신 helper.
func (s *Server) SendRealtime(ctx context.Context, trID, payload string) error {
	return s.SendText(ctx, "0|"+trID+"|001|"+payload)
}

// CloseConn — 클라이언트 측에서는 abnormal close 로 보임 (reconnect 시나리오 테스트용).
func (s *Server) CloseConn() {
	s.mu.Lock()
	c := s.conn
	s.conn = nil
	s.mu.Unlock()
	if c != nil {
		_ = c.Close(wslib.StatusAbnormalClosure, "test forced close")
	}
}

// Received 는 클라이언트 → mock 으로 송신된 메시지 channel.
func (s *Server) Received() <-chan string { return s.receive }
