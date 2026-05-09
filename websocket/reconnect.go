package websocket

import "time"

type reconnectOpts struct {
	Min         time.Duration // 1s
	Max         time.Duration // 30s
	MaxAttempts int           // 10
}

type reconnectController struct {
	opts     reconnectOpts
	attempts int
}

func newReconnectController(opts reconnectOpts) *reconnectController {
	return &reconnectController{opts: opts}
}

// NextBackoff returns the next sleep duration. If attempts > MaxAttempts, returns ErrWSGiveUp.
func (r *reconnectController) NextBackoff() (time.Duration, error) {
	r.attempts++
	if r.attempts > r.opts.MaxAttempts {
		return 0, ErrWSGiveUp
	}
	// 1 * 2^(attempts-1), capped at Max
	d := r.opts.Min << (r.attempts - 1)
	if d > r.opts.Max {
		d = r.opts.Max
	}
	return d, nil
}

// Reset resets the controller. Called on successful reconnection.
func (r *reconnectController) Reset() {
	r.attempts = 0
}
