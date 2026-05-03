// Package ratelimit 은 한투 API 호출 빈도를 제어하는 rate limiter 를 제공.
//
// 요청 간격을 minInterval 로 강제하는 leaky-bucket(request-spacing) 방식.
//
// 사용자에게 노출되지 않는 internal 패키지. kis.Client 가 내부적으로 사용.
package ratelimit

import (
	"context"
	"sync"
	"time"
)

// Limiter 는 thread-safe leaky-bucket rate limiter.
// 연속 호출 간격을 minInterval(= 1s/callsPerSec) 이상으로 유지.
// callsPerSec 이 클수록 호출 빈도가 높아짐. Default 15.
type Limiter struct {
	callsPerSec    float64
	minInterval    time.Duration
	mu             sync.Mutex
	lastCall       time.Time
	totalCalls     int64
	throttledCalls int64
	totalWait      time.Duration
}

// Stats 는 Limiter 의 호출 통계.
type Stats struct {
	CallsPerSec    float64       // 설정된 호출 한도
	TotalCalls     int64         // 누적 호출 수
	ThrottledCalls int64         // 대기한 호출 수
	TotalWait      time.Duration // 누적 대기 시간
	AvgWait        time.Duration // 평균 대기 시간 (throttle 된 것 기준)
}

// New 는 callsPerSec 호출/초 의 Limiter 를 생성. 0 이하면 panic.
func New(callsPerSec float64) *Limiter {
	if callsPerSec <= 0 {
		panic("ratelimit: callsPerSec must be positive")
	}
	return &Limiter{
		callsPerSec: callsPerSec,
		minInterval: time.Duration(float64(time.Second) / callsPerSec),
	}
}

// Wait 는 다음 호출이 허용될 때까지 대기.
// ctx 가 done 되면 그 이유의 에러 반환 (sleep 인터럽트).
//
// 주의: ctx 가 sleep 도중 취소되어도 해당 슬롯(lastCall 예약)은 반환되지 않음.
// 취소 비율이 높은 환경에서는 슬롯 낭비가 발생할 수 있음.
// 또한 throttledCalls / totalWait 통계는 예약된 sleep 기준이라 ctx cancel 로
// 조기 반환해도 감소하지 않음.
func (l *Limiter) Wait(ctx context.Context) error {
	l.mu.Lock()
	now := time.Now()
	elapsed := now.Sub(l.lastCall)
	var sleep time.Duration
	if elapsed < l.minInterval {
		sleep = l.minInterval - elapsed
	}
	l.lastCall = now.Add(sleep)
	l.totalCalls++
	if sleep > 0 {
		l.throttledCalls++
		l.totalWait += sleep
	}
	l.mu.Unlock()

	if sleep <= 0 {
		return ctx.Err() // nil if active; propagate error if already cancelled
	}

	timer := time.NewTimer(sleep)
	defer timer.Stop()
	select {
	case <-timer.C:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Stats 는 현재 통계 스냅샷 반환.
func (l *Limiter) Stats() Stats {
	l.mu.Lock()
	defer l.mu.Unlock()
	var avg time.Duration
	if l.throttledCalls > 0 {
		avg = l.totalWait / time.Duration(l.throttledCalls)
	}
	return Stats{
		CallsPerSec:    l.callsPerSec,
		TotalCalls:     l.totalCalls,
		ThrottledCalls: l.throttledCalls,
		TotalWait:      l.totalWait,
		AvgWait:        avg,
	}
}
