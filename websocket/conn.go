package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	wslib "github.com/coder/websocket"
)

type connManager struct {
	mu       sync.Mutex
	endpoint string
	conn     *wslib.Conn
}

func newConnManager(endpoint string) *connManager {
	return &connManager{endpoint: endpoint}
}

// Dial 은 WebSocket 연결.
func (cm *connManager) Dial(ctx context.Context) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	if cm.conn != nil {
		_ = cm.conn.Close(wslib.StatusNormalClosure, "redial")
		cm.conn = nil
	}
	c, _, err := wslib.Dial(ctx, cm.endpoint, &wslib.DialOptions{
		Subprotocols: nil,
	})
	if err != nil {
		return fmt.Errorf("kis ws: dial %s: %w", cm.endpoint, err)
	}
	c.SetReadLimit(1 << 20) // 1 MiB
	cm.conn = c
	return nil
}

// SendSubscribe 는 subscribe/unsubscribe frame 송신.
//
// trType: "1"=등록, "2"=해제.
func (cm *connManager) SendSubscribe(ctx context.Context, approvalKey, custType, trType, trID, trKey string) error {
	cm.mu.Lock()
	c := cm.conn
	cm.mu.Unlock()
	if c == nil {
		return ErrWSNotConnected
	}
	msg := map[string]any{
		"header": map[string]string{
			"approval_key": approvalKey,
			"custtype":     custType,
			"tr_type":      trType,
			"content-type": "utf-8",
		},
		"body": map[string]any{
			"input": map[string]string{
				"tr_id":  trID,
				"tr_key": trKey,
			},
		},
	}
	raw, _ := json.Marshal(msg)
	return c.Write(ctx, wslib.MessageText, raw)
}

// Read 는 다음 text frame 을 받음 (blocking).
func (cm *connManager) Read(ctx context.Context) (string, error) {
	cm.mu.Lock()
	c := cm.conn
	cm.mu.Unlock()
	if c == nil {
		return "", ErrWSNotConnected
	}
	_, raw, err := c.Read(ctx)
	if err != nil {
		return "", err
	}
	return string(raw), nil
}

// Pong — KIS 가 보낸 PINGPONG JSON 메시지 그대로 echo (서버가 PONG 으로 인식).
func (cm *connManager) Pong(ctx context.Context, raw string) error {
	cm.mu.Lock()
	c := cm.conn
	cm.mu.Unlock()
	if c == nil {
		return ErrWSNotConnected
	}
	return c.Write(ctx, wslib.MessageText, []byte(raw))
}

// Close — graceful close.
func (cm *connManager) Close() error {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	if cm.conn == nil {
		return nil
	}
	err := cm.conn.Close(wslib.StatusNormalClosure, "client shutdown")
	cm.conn = nil
	return err
}
