package token

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestRedis(t *testing.T) (*redis.Client, *miniredis.Miniredis) {
	t.Helper()
	mr := miniredis.RunT(t)
	rdb := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	t.Cleanup(func() { _ = rdb.Close() })
	return rdb, mr
}

func TestRedisStorage_SaveAndLoad(t *testing.T) {
	rdb, _ := newTestRedis(t)
	s := NewRedisStorage(rdb, "kis:token:test")

	tok := &AccessToken{
		Value:     "abc",
		TokenType: "Bearer",
		ExpiresAt: time.Now().Add(time.Hour).Round(time.Second),
	}
	require.NoError(t, s.Save(context.Background(), tok))

	loaded, err := s.Load(context.Background())
	require.NoError(t, err)
	require.NotNil(t, loaded)
	assert.Equal(t, "abc", loaded.Value)
}

func TestRedisStorage_LoadEmpty(t *testing.T) {
	rdb, _ := newTestRedis(t)
	s := NewRedisStorage(rdb, "kis:token:empty")
	loaded, err := s.Load(context.Background())
	require.NoError(t, err)
	assert.Nil(t, loaded)
}

func TestRedisStorage_Clear(t *testing.T) {
	rdb, _ := newTestRedis(t)
	s := NewRedisStorage(rdb, "kis:token:clear")
	ctx := context.Background()
	require.NoError(t, s.Save(ctx, &AccessToken{Value: "x", ExpiresAt: time.Now().Add(time.Hour)}))
	require.NoError(t, s.Clear(ctx))
	loaded, _ := s.Load(ctx)
	assert.Nil(t, loaded)
}

func TestRedisStorage_TTL(t *testing.T) {
	rdb, mr := newTestRedis(t)
	s := NewRedisStorage(rdb, "kis:token:ttl")
	ctx := context.Background()
	require.NoError(t, s.Save(ctx, &AccessToken{
		Value: "x", ExpiresAt: time.Now().Add(2 * time.Hour),
	}))
	ttl := mr.TTL("kis:token:ttl")
	assert.Greater(t, ttl, 90*time.Minute)
	assert.LessOrEqual(t, ttl, 2*time.Hour+1*time.Second)
}
