package kis

import "errors"

// APIError 는 한국투자증권 API 가 비정상 응답(rt_cd != "0") 을 돌려줄 때 발생.
type APIError struct {
	RtCode  string // 한국투자 응답의 rt_cd
	MsgCode string // 한국투자 응답의 msg_cd
	Message string // 한국투자 응답의 msg1
	TrID    string // 디버깅용 — 어느 transaction 이 실패했는지
}

func (e *APIError) Error() string {
	return "kis: API error [" + e.MsgCode + "] " + e.Message
}

// Sentinel 에러. errors.Is 로 분기.
var (
	ErrTokenExpired = errors.New("kis: token expired")
	ErrRateLimited  = errors.New("kis: rate limited")
	ErrNotFound     = errors.New("kis: resource not found")
	ErrUnauthorized = errors.New("kis: unauthorized")
)
