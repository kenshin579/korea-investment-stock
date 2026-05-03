package mastercache

import (
	"context"
	"errors"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCache_GetCold(t *testing.T) {
	c := New(t.TempDir(), 24*time.Hour)
	calls := atomic.Int64{}
	fetch := func(ctx context.Context) ([]byte, error) {
		calls.Add(1)
		return []byte("hello"), nil
	}
	data, err := c.Get(context.Background(), "test.bin", fetch)
	require.NoError(t, err)
	assert.Equal(t, []byte("hello"), data)
	assert.Equal(t, int64(1), calls.Load())
}

func TestCache_GetHot(t *testing.T) {
	c := New(t.TempDir(), 24*time.Hour)
	calls := atomic.Int64{}
	fetch := func(ctx context.Context) ([]byte, error) {
		calls.Add(1)
		return []byte("hello"), nil
	}
	_, _ = c.Get(context.Background(), "test.bin", fetch)
	data, _ := c.Get(context.Background(), "test.bin", fetch)
	assert.Equal(t, []byte("hello"), data)
	assert.Equal(t, int64(1), calls.Load(), "second Get should hit cache")
}

func TestCache_TTLExpired(t *testing.T) {
	c := New(t.TempDir(), 1*time.Millisecond)
	calls := atomic.Int64{}
	fetch := func(ctx context.Context) ([]byte, error) {
		calls.Add(1)
		return []byte("v" + strings.Repeat("x", 1)), nil
	}
	_, _ = c.Get(context.Background(), "test.bin", fetch)
	time.Sleep(10 * time.Millisecond)
	_, _ = c.Get(context.Background(), "test.bin", fetch)
	assert.Equal(t, int64(2), calls.Load(), "TTL expired → refetch")
}

func TestCache_FetchError_FallbackToStaleIfExists(t *testing.T) {
	c := New(t.TempDir(), 1*time.Millisecond)
	calls := atomic.Int64{}
	fetch := func(ctx context.Context) ([]byte, error) {
		calls.Add(1)
		if calls.Load() == 1 {
			return []byte("ok"), nil
		}
		return nil, errors.New("network down")
	}
	_, _ = c.Get(context.Background(), "test.bin", fetch)
	time.Sleep(10 * time.Millisecond)

	data, err := c.Get(context.Background(), "test.bin", fetch)
	require.NoError(t, err)
	assert.Equal(t, []byte("ok"), data, "fallback to stale on fetch failure")
}

func TestDefaultDir(t *testing.T) {
	dir, err := DefaultDir()
	require.NoError(t, err)
	assert.True(t, strings.Contains(dir, "kis"), "default dir should contain 'kis'")
}
