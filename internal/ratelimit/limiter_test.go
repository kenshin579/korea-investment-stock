package ratelimit

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew_InvalidRate(t *testing.T) {
	assert.Panics(t, func() { New(0) })
	assert.Panics(t, func() { New(-1) })
}

func TestLimiter_Wait_Throttles(t *testing.T) {
	l := New(10) // 10 req/sec → min interval 100ms
	ctx := context.Background()

	require.NoError(t, l.Wait(ctx))
	start := time.Now()
	require.NoError(t, l.Wait(ctx))
	elapsed := time.Since(start)
	assert.GreaterOrEqual(t, elapsed, 90*time.Millisecond,
		"second Wait should sleep ~100ms")
}

func TestLimiter_Wait_ContextCancelled(t *testing.T) {
	l := New(1) // 1 req/sec → min interval 1s
	ctx := context.Background()
	require.NoError(t, l.Wait(ctx))

	ctxCancel, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	err := l.Wait(ctxCancel) // 1초 sleep 시작, 50ms 후 ctx done
	assert.ErrorIs(t, err, context.DeadlineExceeded)
}

func TestLimiter_Stats(t *testing.T) {
	l := New(1000) // 빠른 rate, throttle 거의 없음
	ctx := context.Background()
	for i := 0; i < 5; i++ {
		require.NoError(t, l.Wait(ctx))
	}
	s := l.Stats()
	assert.Equal(t, int64(5), s.TotalCalls)
}

func TestLimiter_ConcurrentSafe(t *testing.T) {
	l := New(1000)
	ctx := context.Background()
	done := make(chan struct{}, 10)
	for i := 0; i < 10; i++ {
		go func() {
			for j := 0; j < 5; j++ {
				_ = l.Wait(ctx)
			}
			done <- struct{}{}
		}()
	}
	for i := 0; i < 10; i++ {
		<-done
	}
	assert.Equal(t, int64(50), l.Stats().TotalCalls)
}

func TestLimiter_Wait_PreCancelledContext(t *testing.T) {
	l := New(1000) // 빠른 rate, throttle 거의 없음

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // 미리 취소

	err := l.Wait(ctx)
	assert.ErrorIs(t, err, context.Canceled, "pre-cancelled ctx should propagate error even on fast path")
}
