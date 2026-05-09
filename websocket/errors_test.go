package websocket_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kenshin579/korea-investment-stock/websocket"
)

func TestWSServerError_Error(t *testing.T) {
	e := &websocket.WSServerError{TrID: "H0STCNT0", MsgCd: "OPSP0001", Msg: "ALREADY IN SUBSCRIBE"}
	assert.Contains(t, e.Error(), "H0STCNT0")
	assert.Contains(t, e.Error(), "OPSP0001")
	assert.Contains(t, e.Error(), "ALREADY IN SUBSCRIBE")
}

func TestSentinelErrors(t *testing.T) {
	assert.NotNil(t, websocket.ErrWSGiveUp)
	assert.NotNil(t, websocket.ErrWSApprovalFailed)
	assert.NotNil(t, websocket.ErrWSInvalidFrame)
	assert.NotNil(t, websocket.ErrWSNotConnected)
	assert.NotNil(t, websocket.ErrWSDuplicateSub)
	assert.NotNil(t, websocket.ErrWSEncryptedNotSupported)
	// errors.Is 동작 검증
	wrapped := errors.Join(websocket.ErrWSGiveUp, errors.New("net err"))
	assert.True(t, errors.Is(wrapped, websocket.ErrWSGiveUp))
}
