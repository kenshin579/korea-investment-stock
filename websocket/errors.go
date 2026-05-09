package websocket

import (
	"errors"
	"fmt"
)

var (
	ErrWSNotConnected          = errors.New("kis ws: not connected")
	ErrWSGiveUp                = errors.New("kis ws: give up after max reconnect attempts")
	ErrWSApprovalFailed        = errors.New("kis ws: approval key issuance failed")
	ErrWSInvalidFrame          = errors.New("kis ws: invalid frame format")
	ErrWSDuplicateSub          = errors.New("kis ws: duplicate subscription")
	ErrWSEncryptedNotSupported = errors.New("kis ws: encrypted frames not supported in Phase 8")
)

// WSServerError 는 KIS 서버가 등록/해제 응답에서 반환한 에러.
type WSServerError struct {
	TrID  string // H0STCNT0
	MsgCd string // OPSP0001 등
	Msg   string // "ALREADY IN SUBSCRIBE" 등
}

func (e *WSServerError) Error() string {
	return fmt.Sprintf("kis ws server error: tr_id=%s msg_cd=%s msg=%s", e.TrID, e.MsgCd, e.Msg)
}
