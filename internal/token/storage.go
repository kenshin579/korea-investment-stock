// Package token 은 한국투자증권 OAuth 토큰의 발급/저장/갱신을 담당.
//
// 사용자에게 노출되지 않는 internal 패키지. kis.Client 가 내부적으로 사용.
package token

import (
	"context"
	"time"
)

// expiryMargin 은 만료 전 선제 발급 마진. 만료 5분 전부터 IsExpired = true.
const expiryMargin = 5 * time.Minute

// AccessToken 은 발급된 OAuth 토큰.
type AccessToken struct {
	Value     string    // raw token (Bearer prefix 없음)
	TokenType string    // "Bearer"
	ExpiresAt time.Time // 만료 시각 (Asia/Seoul)
}

// IsExpired 는 토큰이 만료되었거나 만료 임박(5분 이내) 인지 반환.
func (t *AccessToken) IsExpired() bool {
	return time.Until(t.ExpiresAt) <= expiryMargin
}

// Bearer 는 "Bearer <value>" 형태의 Authorization 헤더 값 반환.
func (t *AccessToken) Bearer() string {
	return "Bearer " + t.Value
}

// Storage 는 토큰 영구 저장소 인터페이스.
// FileStorage / RedisStorage 가 구현.
type Storage interface {
	// Save 는 토큰을 저장.
	Save(ctx context.Context, token *AccessToken) error
	// Load 는 저장된 토큰을 반환. 없으면 nil, nil.
	Load(ctx context.Context) (*AccessToken, error)
	// Clear 는 저장된 토큰 삭제.
	Clear(ctx context.Context) error
}
