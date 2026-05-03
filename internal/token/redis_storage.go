package token

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisStorage 는 토큰을 Redis 에 저장. 분산 환경에서 여러 인스턴스가 토큰 공유.
type RedisStorage struct {
	client *redis.Client
	key    string
}

// NewRedisStorage 는 지정된 redis client 와 key 를 사용하는 storage 생성.
func NewRedisStorage(client *redis.Client, key string) *RedisStorage {
	return &RedisStorage{client: client, key: key}
}

// Save 는 토큰을 Redis 에 저장. TTL 은 토큰의 ExpiresAt - Now.
func (s *RedisStorage) Save(ctx context.Context, token *AccessToken) error {
	data, err := json.Marshal(fileToken{
		Value: token.Value, TokenType: token.TokenType, ExpiresAt: token.ExpiresAt,
	})
	if err != nil {
		return err
	}
	ttl := time.Until(token.ExpiresAt)
	if ttl <= 0 {
		ttl = time.Second
	}
	return s.client.Set(ctx, s.key, data, ttl).Err()
}

// Load 는 토큰을 Redis 에서 읽음. 없으면 nil, nil.
func (s *RedisStorage) Load(ctx context.Context) (*AccessToken, error) {
	data, err := s.client.Get(ctx, s.key).Bytes()
	if errors.Is(err, redis.Nil) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var ft fileToken
	if err := json.Unmarshal(data, &ft); err != nil {
		return nil, err
	}
	return &AccessToken{Value: ft.Value, TokenType: ft.TokenType, ExpiresAt: ft.ExpiresAt}, nil
}

// Clear 는 토큰 키 삭제. 없어도 에러 없음.
func (s *RedisStorage) Clear(ctx context.Context) error {
	return s.client.Del(ctx, s.key).Err()
}
