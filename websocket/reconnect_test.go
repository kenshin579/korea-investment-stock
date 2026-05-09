package websocket

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestReconnect_BackoffSequence(t *testing.T) {
	r := newReconnectController(reconnectOpts{
		Min:         1 * time.Second,
		Max:         30 * time.Second,
		MaxAttempts: 10,
	})

	// 1, 2, 4, 8, 16, 30, 30, 30, 30, 30 (cap)
	expected := []time.Duration{1, 2, 4, 8, 16, 30, 30, 30, 30, 30}
	for i, want := range expected {
		got, err := r.NextBackoff()
		assert.NoError(t, err, "attempt %d", i+1)
		assert.Equalf(t, want*time.Second, got, "attempt %d", i+1)
	}
}

func TestReconnect_GiveUp(t *testing.T) {
	r := newReconnectController(reconnectOpts{
		Min:         1 * time.Second,
		Max:         30 * time.Second,
		MaxAttempts: 3,
	})

	r.NextBackoff()
	r.NextBackoff()
	r.NextBackoff()
	_, err := r.NextBackoff()
	assert.True(t, errors.Is(err, ErrWSGiveUp))
}

func TestReconnect_Reset(t *testing.T) {
	r := newReconnectController(reconnectOpts{
		Min:         1 * time.Second,
		Max:         30 * time.Second,
		MaxAttempts: 10,
	})
	r.NextBackoff()
	r.NextBackoff()
	r.Reset()
	got, _ := r.NextBackoff()
	assert.Equal(t, 1*time.Second, got)
}
