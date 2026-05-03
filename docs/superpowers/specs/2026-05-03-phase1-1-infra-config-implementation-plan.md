# Phase 1.1 — Infrastructure + Config Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** `korea-investment-stock` Go 라이브러리의 인프라 (rate limiter, token manager, HTTP client, master cache) 와 Config 진입점을 구현해 v0.1.0 release 준비.

**Architecture:** TDD 기반. 모든 인프라는 `internal/` 아래에 두어 외부 노출 안 함. root 의 `client.go` 가 인프라를 wiring 하고 외부에는 functional options + 3개 진입점 (`NewClient`, `NewClientFromEnv`, `NewClientFromYAML`) 노출. Phase 1.2 부터 메서드 추가가 깔끔하게 되도록 sub-package (`domestic/`, `overseas/`) 의 `Client` struct 를 placeholder 로 정의.

**Tech Stack:** Go 1.23+, `github.com/go-resty/resty/v2`, `github.com/shopspring/decimal`, `github.com/redis/go-redis/v9`, `gopkg.in/yaml.v3`, `github.com/jarcoal/httpmock`, `github.com/stretchr/testify`, `golang.org/x/sync/singleflight`

**참고 spec:**
- Phase 1: `docs/superpowers/specs/2026-05-03-phase1-api-coverage-design.md` (commit `66d9733`)
- Phase 0: `docs/superpowers/specs/2026-05-03-korea-investment-go-migration-design.md`

**Python 참조 코드:** `python-final` 태그의 `korea_investment_stock/{rate_limit,token,parsers}/` — 동작 명세 참고용

---

## 사전 정보

| 항목 | 값 |
|------|---|
| 작업 브랜치 | `docs/phase1-api-coverage-spec` (이미 spec 커밋 완료) |
| 시작 HEAD | `66d9733` (Phase 1 spec) |
| Release 목표 | `v0.1.0` (PR merge 후 태그) |
| PR 베이스 | `main` |
| 현재 main 상태 | Phase 0 후 = `client.go` 스켈레톤 + `domestic/`, `overseas/`, `internal/{httpclient,ratelimit,token}/doc.go` 만 존재 |

## 파일 구조

### 신규 (Foundation, root)
- `errors.go` — `APIError` 타입 + sentinel 에러
- `options.go` — `Option` 타입 + `clientOptions` struct + `WithXXX` functions
- `config.go` — `Config` struct + `LoadConfigFromEnv` + `LoadConfigFromYAML`

### 신규 (internal)
- `internal/ratelimit/limiter.go` — `Limiter` struct + `Wait(ctx)` + `Stats`
- `internal/ratelimit/limiter_test.go`
- `internal/token/storage.go` — `Storage` interface + `AccessToken` 타입
- `internal/token/file_storage.go` — `FileStorage` 구현
- `internal/token/file_storage_test.go`
- `internal/token/redis_storage.go` — `RedisStorage` 구현
- `internal/token/redis_storage_test.go` — miniredis 사용
- `internal/token/manager.go` — OAuth 발급 + singleflight + 자동 갱신
- `internal/token/manager_test.go`
- `internal/httpclient/client.go` — resty wrap + `Do(ctx, req)` + 재시도 + 토큰 만료 자동 감지
- `internal/httpclient/client_test.go`
- `internal/httpclient/hashkey.go` — Hashkey 발급
- `internal/httpclient/hashkey_test.go`
- `internal/mastercache/cache.go` — KOSPI/KOSDAQ ZIP 다운로드/디스크 캐시
- `internal/mastercache/cache_test.go`

### 수정 (root)
- `client.go` — NewClient 보강 (인프라 wiring, sub-client 주입), `RealEnv`/`PaperEnv` 보강
- `from_env.go` (신규) — `NewClientFromEnv`
- `from_yaml.go` (신규) — `NewClientFromYAML`
- `auth.go` (신규) — 외부 노출 `IssueAccessToken`

### 수정 (sub-packages)
- `domestic/client.go` (신규) — `Client` struct (placeholder, http 필드만), `New(http *httpclient.Client) *Client`
- `domestic/doc.go` (수정) — placeholder 코멘트 갱신
- `overseas/client.go` (신규) — 동일
- `overseas/doc.go` (수정) — 동일

### 신규 (examples)
- `examples/basic/main.go` — `NewClient` + 토큰 발급만 호출
- `examples/env_config/main.go` — `NewClientFromEnv`
- `examples/yaml_config/main.go` + `examples/yaml_config/config.yaml`

### 신규 (CI 보강 — 선택)
- `.github/workflows/test.yml` — go test + vet + coverage threshold

---

## Task 1: 의존성 설치 + go.sum 생성

**Files:**
- Modify: `go.mod`, `go.sum`

- [ ] **Step 1: 의존성 추가**

Run:
```bash
go get github.com/go-resty/resty/v2
go get github.com/shopspring/decimal
go get github.com/redis/go-redis/v9
go get github.com/alicebob/miniredis/v2
go get gopkg.in/yaml.v3
go get github.com/jarcoal/httpmock
go get github.com/stretchr/testify
go get golang.org/x/sync/singleflight
go mod tidy
```

- [ ] **Step 2: 검증**

Run:
```bash
go mod verify && cat go.mod
```
Expected: `all modules verified`. `go.mod` 의 require 블록에 위 8개 패키지 포함.

- [ ] **Step 3: Commit**

```bash
git add go.mod go.sum
git commit -m "[chore] Phase 1.1 의존성 추가

resty, decimal, go-redis, miniredis, yaml.v3, httpmock, testify, singleflight.

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>"
```

---

## Task 2: errors.go — APIError + sentinel

**Files:**
- Create: `errors.go`
- Create: `errors_test.go`

- [ ] **Step 1: 테스트 작성** — `errors_test.go`

```go
package kis

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAPIError_Error(t *testing.T) {
	err := &APIError{
		RtCode:  "1",
		MsgCode: "MCA00001",
		Message: "잘못된 요청",
		TrID:    "FHKST01010100",
	}
	assert.Contains(t, err.Error(), "MCA00001")
	assert.Contains(t, err.Error(), "잘못된 요청")
}

func TestAPIError_ErrorsAs(t *testing.T) {
	var err error = &APIError{RtCode: "1", MsgCode: "MCA00001", Message: "msg"}
	var apiErr *APIError
	assert.True(t, errors.As(err, &apiErr))
	assert.Equal(t, "MCA00001", apiErr.MsgCode)
}

func TestSentinelErrors(t *testing.T) {
	assert.Equal(t, "kis: token expired", ErrTokenExpired.Error())
	assert.Equal(t, "kis: rate limited", ErrRateLimited.Error())
	assert.Equal(t, "kis: resource not found", ErrNotFound.Error())
	assert.Equal(t, "kis: unauthorized", ErrUnauthorized.Error())
}
```

- [ ] **Step 2: 테스트 실행 → FAIL**

Run: `go test ./... -run 'TestAPIError|TestSentinel' -v`
Expected: 컴파일 실패 (APIError 미정의)

- [ ] **Step 3: 구현** — `errors.go`

```go
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
```

- [ ] **Step 4: 테스트 실행 → PASS**

Run: `go test ./... -run 'TestAPIError|TestSentinel' -v`
Expected: PASS (3 test cases)

- [ ] **Step 5: Commit**

```bash
git add errors.go errors_test.go
git commit -m "[feat] APIError 타입 + sentinel 에러 추가

errors.As 로 한투 API 에러 분기, errors.Is 로 sentinel 분기.

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>"
```

---

## Task 3: internal/ratelimit/limiter.go — token bucket rate limiter

**Files:**
- Create: `internal/ratelimit/limiter.go`
- Create: `internal/ratelimit/limiter_test.go`

- [ ] **Step 1: 테스트 작성** — `internal/ratelimit/limiter_test.go`

```go
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
```

- [ ] **Step 2: 테스트 실행 → FAIL**

Run: `go test ./internal/ratelimit/... -v`
Expected: 컴파일 실패 (Limiter 미정의)

- [ ] **Step 3: 구현** — `internal/ratelimit/limiter.go`

```go
// Package ratelimit 은 한투 API 호출 빈도를 제어하는 토큰 버킷 rate limiter 를 제공.
//
// 사용자에게 노출되지 않는 internal 패키지. kis.Client 가 내부적으로 사용.
package ratelimit

import (
	"context"
	"sync"
	"time"
)

// Limiter 는 thread-safe 토큰 버킷 rate limiter.
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
		return nil
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
```

- [ ] **Step 4: 테스트 실행 → PASS**

Run: `go test ./internal/ratelimit/... -v -race`
Expected: PASS (5 test cases, race 검증)

- [ ] **Step 5: Commit**

```bash
git add internal/ratelimit/
git commit -m "[feat] internal/ratelimit Limiter 추가

토큰 버킷 rate limiter, ctx 지원, thread-safe, 통계.

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>"
```

---

## Task 4: internal/token Storage interface + AccessToken type

**Files:**
- Create: `internal/token/storage.go`
- Create: `internal/token/storage_test.go`

- [ ] **Step 1: 테스트 작성** — `internal/token/storage_test.go`

```go
package token

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAccessToken_IsExpired(t *testing.T) {
	cases := []struct {
		name     string
		expires  time.Time
		expected bool
	}{
		{"expired", time.Now().Add(-1 * time.Hour), true},
		{"about to expire (within margin)", time.Now().Add(2 * time.Minute), true},
		{"valid", time.Now().Add(1 * time.Hour), false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tok := &AccessToken{ExpiresAt: tc.expires}
			assert.Equal(t, tc.expected, tok.IsExpired())
		})
	}
}

func TestAccessToken_Bearer(t *testing.T) {
	tok := &AccessToken{Value: "abc123"}
	assert.Equal(t, "Bearer abc123", tok.Bearer())
}
```

- [ ] **Step 2: 테스트 실행 → FAIL**

Run: `go test ./internal/token/... -v`
Expected: 컴파일 실패

- [ ] **Step 3: 구현** — `internal/token/storage.go`

```go
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
```

- [ ] **Step 4: 테스트 실행 → PASS**

Run: `go test ./internal/token/... -v`
Expected: PASS (4 sub-tests)

- [ ] **Step 5: Commit**

```bash
git add internal/token/storage.go internal/token/storage_test.go
git commit -m "[feat] token.Storage 인터페이스 + AccessToken 타입

만료 5분 전 선제 발급 마진, Bearer 헬퍼.

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>"
```

---

## Task 5: internal/token/file_storage.go — FileStorage

**Files:**
- Create: `internal/token/file_storage.go`
- Create: `internal/token/file_storage_test.go`

- [ ] **Step 1: 테스트 작성**

```go
package token

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileStorage_SaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "token.json")
	s := NewFileStorage(path)
	ctx := context.Background()

	tok := &AccessToken{
		Value:     "abc123",
		TokenType: "Bearer",
		ExpiresAt: time.Now().Add(1 * time.Hour).Round(time.Second),
	}
	require.NoError(t, s.Save(ctx, tok))

	loaded, err := s.Load(ctx)
	require.NoError(t, err)
	require.NotNil(t, loaded)
	assert.Equal(t, "abc123", loaded.Value)
	assert.Equal(t, "Bearer", loaded.TokenType)
	assert.True(t, loaded.ExpiresAt.Equal(tok.ExpiresAt))
}

func TestFileStorage_LoadEmpty(t *testing.T) {
	dir := t.TempDir()
	s := NewFileStorage(filepath.Join(dir, "token.json"))
	loaded, err := s.Load(context.Background())
	require.NoError(t, err)
	assert.Nil(t, loaded)
}

func TestFileStorage_Clear(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "token.json")
	s := NewFileStorage(path)
	ctx := context.Background()

	require.NoError(t, s.Save(ctx, &AccessToken{Value: "x", ExpiresAt: time.Now().Add(time.Hour)}))
	require.NoError(t, s.Clear(ctx))

	_, err := os.Stat(path)
	assert.True(t, os.IsNotExist(err))
}

func TestFileStorage_Permissions(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "token.json")
	s := NewFileStorage(path)
	require.NoError(t, s.Save(context.Background(), &AccessToken{
		Value: "x", ExpiresAt: time.Now().Add(time.Hour),
	}))

	info, err := os.Stat(path)
	require.NoError(t, err)
	// 토큰 파일은 0600 (소유자만 읽기/쓰기) 권한이어야 함
	assert.Equal(t, os.FileMode(0600), info.Mode().Perm())
}
```

- [ ] **Step 2: 테스트 실행 → FAIL**

Run: `go test ./internal/token/... -v -run TestFileStorage`
Expected: 컴파일 실패 (NewFileStorage 미정의)

- [ ] **Step 3: 구현** — `internal/token/file_storage.go`

```go
package token

import (
	"context"
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"time"
)

// FileStorage 는 토큰을 로컬 파일에 JSON 으로 저장.
type FileStorage struct {
	path string
}

// NewFileStorage 는 지정된 경로에 토큰을 저장하는 FileStorage 생성.
// 파일은 0600 권한으로 작성.
func NewFileStorage(path string) *FileStorage {
	return &FileStorage{path: path}
}

type fileToken struct {
	Value     string    `json:"value"`
	TokenType string    `json:"token_type"`
	ExpiresAt time.Time `json:"expires_at"`
}

// Save 는 토큰을 파일에 저장.
func (s *FileStorage) Save(_ context.Context, token *AccessToken) error {
	if err := os.MkdirAll(filepath.Dir(s.path), 0700); err != nil {
		return err
	}
	data, err := json.Marshal(fileToken{
		Value: token.Value, TokenType: token.TokenType, ExpiresAt: token.ExpiresAt,
	})
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0600)
}

// Load 는 파일에서 토큰을 읽음. 파일이 없으면 nil, nil.
func (s *FileStorage) Load(_ context.Context) (*AccessToken, error) {
	data, err := os.ReadFile(s.path)
	if errors.Is(err, fs.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var ft fileToken
	if err := json.Unmarshal(data, &ft); err != nil {
		return nil, err
	}
	return &AccessToken{
		Value: ft.Value, TokenType: ft.TokenType, ExpiresAt: ft.ExpiresAt,
	}, nil
}

// Clear 는 파일을 삭제. 파일이 없으면 에러 없음.
func (s *FileStorage) Clear(_ context.Context) error {
	err := os.Remove(s.path)
	if errors.Is(err, fs.ErrNotExist) {
		return nil
	}
	return err
}
```

- [ ] **Step 4: 테스트 실행 → PASS**

Run: `go test ./internal/token/... -v -run TestFileStorage`
Expected: 4 PASS

- [ ] **Step 5: Commit**

```bash
git add internal/token/file_storage.go internal/token/file_storage_test.go
git commit -m "[feat] token.FileStorage 추가

JSON 직렬화, 0600 권한, 디렉터리 자동 생성.

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>"
```

---

## Task 6: internal/token/redis_storage.go — RedisStorage (miniredis 테스트)

**Files:**
- Create: `internal/token/redis_storage.go`
- Create: `internal/token/redis_storage_test.go`

- [ ] **Step 1: 테스트 작성**

```go
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
	// TTL 이 만료 시각과 비슷해야 함 (2시간 ± margin)
	assert.Greater(t, ttl, 90*time.Minute)
	assert.LessOrEqual(t, ttl, 2*time.Hour+1*time.Second)
}
```

- [ ] **Step 2: 테스트 실행 → FAIL**

Run: `go test ./internal/token/... -v -run TestRedisStorage`
Expected: 컴파일 실패

- [ ] **Step 3: 구현** — `internal/token/redis_storage.go`

```go
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
		ttl = time.Second // ttl 0 = 영구. 음수면 즉시 삭제. 안전 위해 1초.
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
```

- [ ] **Step 4: 테스트 실행 → PASS**

Run: `go test ./internal/token/... -v -run TestRedisStorage`
Expected: 4 PASS

- [ ] **Step 5: Commit**

```bash
git add internal/token/redis_storage.go internal/token/redis_storage_test.go
git commit -m "[feat] token.RedisStorage 추가

go-redis/v9 사용, TTL 자동 설정, miniredis 단위 테스트.

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>"
```

---

## Task 7: internal/token/manager.go — OAuth Manager + singleflight

**Files:**
- Create: `internal/token/manager.go`
- Create: `internal/token/manager_test.go`

- [ ] **Step 1: 테스트 작성** — `internal/token/manager_test.go`

```go
package token

import (
	"context"
	"net/http"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestManager_Get_FreshFetch(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder(http.MethodPost, "=~/oauth2/tokenP",
		httpmock.NewStringResponder(200, `{
			"access_token": "fresh-token",
			"token_type": "Bearer",
			"expires_in": 86400,
			"access_token_token_expired": "2099-12-31 23:59:59"
		}`))

	storage := NewFileStorage(t.TempDir() + "/token.json")
	m := NewManager(Config{
		Storage:   storage,
		BaseURL:   "https://openapi.koreainvestment.com:9443",
		APIKey:    "k",
		APISecret: "s",
		HTTPClient: http.DefaultClient,
	})
	bearer, err := m.Get(context.Background())
	require.NoError(t, err)
	assert.Equal(t, "Bearer fresh-token", bearer)
}

func TestManager_Get_UsesCache(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	calls := atomic.Int64{}
	httpmock.RegisterResponder(http.MethodPost, "=~/oauth2/tokenP",
		func(req *http.Request) (*http.Response, error) {
			calls.Add(1)
			return httpmock.NewStringResponse(200, `{
				"access_token": "cached",
				"token_type": "Bearer",
				"access_token_token_expired": "2099-12-31 23:59:59"
			}`), nil
		})

	storage := NewFileStorage(t.TempDir() + "/token.json")
	m := NewManager(Config{
		Storage: storage, BaseURL: "https://x", APIKey: "k", APISecret: "s",
		HTTPClient: http.DefaultClient,
	})
	for i := 0; i < 5; i++ {
		_, err := m.Get(context.Background())
		require.NoError(t, err)
	}
	assert.Equal(t, int64(1), calls.Load(), "cached token should be reused")
}

func TestManager_Refresh_Forces(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	calls := atomic.Int64{}
	httpmock.RegisterResponder(http.MethodPost, "=~/oauth2/tokenP",
		func(req *http.Request) (*http.Response, error) {
			calls.Add(1)
			return httpmock.NewStringResponse(200, `{
				"access_token": "refreshed",
				"token_type": "Bearer",
				"access_token_token_expired": "2099-12-31 23:59:59"
			}`), nil
		})

	storage := NewFileStorage(t.TempDir() + "/token.json")
	m := NewManager(Config{
		Storage: storage, BaseURL: "https://x", APIKey: "k", APISecret: "s",
		HTTPClient: http.DefaultClient,
	})
	require.NoError(t, m.warmup(context.Background()))
	_, err := m.Refresh(context.Background())
	require.NoError(t, err)
	assert.Equal(t, int64(2), calls.Load(), "warmup + Refresh = 2 calls")
}

func TestManager_Get_Singleflight(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	calls := atomic.Int64{}
	httpmock.RegisterResponder(http.MethodPost, "=~/oauth2/tokenP",
		func(req *http.Request) (*http.Response, error) {
			calls.Add(1)
			time.Sleep(50 * time.Millisecond) // 동시 발급 race window
			return httpmock.NewStringResponse(200, `{
				"access_token": "single",
				"token_type": "Bearer",
				"access_token_token_expired": "2099-12-31 23:59:59"
			}`), nil
		})

	storage := NewFileStorage(t.TempDir() + "/token.json")
	m := NewManager(Config{
		Storage: storage, BaseURL: "https://x", APIKey: "k", APISecret: "s",
		HTTPClient: http.DefaultClient,
	})

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _ = m.Get(context.Background())
		}()
	}
	wg.Wait()
	assert.Equal(t, int64(1), calls.Load(), "10 concurrent Get → 1 OAuth call (singleflight)")
}
```

- [ ] **Step 2: 테스트 실행 → FAIL**

Run: `go test ./internal/token/... -v -run TestManager`
Expected: 컴파일 실패 (Manager 미정의)

- [ ] **Step 3: 구현** — `internal/token/manager.go`

```go
package token

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"golang.org/x/sync/singleflight"
)

// oauthPath 는 한국투자 OAuth 토큰 발급 엔드포인트.
const oauthPath = "/oauth2/tokenP"

// seoulTZ 는 한투 응답의 만료 시각이 사용하는 시간대.
var seoulTZ *time.Location

func init() {
	loc, err := time.LoadLocation("Asia/Seoul")
	if err != nil {
		// fallback (KST = UTC+9)
		loc = time.FixedZone("KST", 9*3600)
	}
	seoulTZ = loc
}

// Config 는 Manager 생성 옵션.
type Config struct {
	Storage    Storage
	BaseURL    string
	APIKey     string
	APISecret  string
	HTTPClient *http.Client
}

// Manager 는 OAuth 토큰의 발급/캐시/갱신을 담당.
// 동시 호출 시 한 번만 발급 (singleflight).
type Manager struct {
	cfg     Config
	mu      sync.RWMutex
	cached  *AccessToken
	flight  singleflight.Group
}

// NewManager 는 Manager 생성.
func NewManager(cfg Config) *Manager {
	if cfg.HTTPClient == nil {
		cfg.HTTPClient = http.DefaultClient
	}
	return &Manager{cfg: cfg}
}

// Get 은 유효한 토큰의 Bearer 문자열 반환.
// 캐시된 토큰이 없거나 만료 임박이면 자동 발급.
func (m *Manager) Get(ctx context.Context) (string, error) {
	if t := m.cachedValid(); t != nil {
		return t.Bearer(), nil
	}

	// storage 에서 로드 시도
	if t, err := m.cfg.Storage.Load(ctx); err == nil && t != nil && !t.IsExpired() {
		m.mu.Lock()
		m.cached = t
		m.mu.Unlock()
		return t.Bearer(), nil
	}

	// 새로 발급 (singleflight 로 동시 호출 1번만)
	v, err, _ := m.flight.Do("issue", func() (interface{}, error) {
		return m.issue(ctx)
	})
	if err != nil {
		return "", err
	}
	return v.(*AccessToken).Bearer(), nil
}

// Refresh 는 캐시 무시하고 강제로 새 토큰 발급.
func (m *Manager) Refresh(ctx context.Context) (string, error) {
	v, err, _ := m.flight.Do("refresh", func() (interface{}, error) {
		return m.issue(ctx)
	})
	if err != nil {
		return "", err
	}
	return v.(*AccessToken).Bearer(), nil
}

// warmup 은 테스트 편의용 — 첫 토큰을 미리 받아둠.
func (m *Manager) warmup(ctx context.Context) error {
	_, err := m.Get(ctx)
	return err
}

func (m *Manager) cachedValid() *AccessToken {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.cached == nil || m.cached.IsExpired() {
		return nil
	}
	return m.cached
}

type oauthResp struct {
	AccessToken          string `json:"access_token"`
	TokenType            string `json:"token_type"`
	ExpiresIn            int    `json:"expires_in"`
	AccessTokenExpiredAt string `json:"access_token_token_expired"`
}

func (m *Manager) issue(ctx context.Context) (*AccessToken, error) {
	body, _ := json.Marshal(map[string]string{
		"grant_type": "client_credentials",
		"appkey":     m.cfg.APIKey,
		"appsecret":  m.cfg.APISecret,
	})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		strings.TrimRight(m.cfg.BaseURL, "/")+oauthPath,
		strings.NewReader(string(body)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := m.cfg.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("kis: token issue: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("kis: token issue: HTTP %d", resp.StatusCode)
	}
	var oauth oauthResp
	if err := json.NewDecoder(resp.Body).Decode(&oauth); err != nil {
		return nil, fmt.Errorf("kis: token decode: %w", err)
	}

	expiresAt, err := parseExpiry(oauth)
	if err != nil {
		return nil, err
	}
	tok := &AccessToken{
		Value:     oauth.AccessToken,
		TokenType: oauth.TokenType,
		ExpiresAt: expiresAt,
	}

	if err := m.cfg.Storage.Save(ctx, tok); err != nil {
		// 저장 실패는 warning, 발급은 성공
	}
	m.mu.Lock()
	m.cached = tok
	m.mu.Unlock()
	return tok, nil
}

// parseExpiry 는 한투 응답의 만료 시각 파싱.
// 우선순위: access_token_token_expired (Asia/Seoul) → expires_in (now + delta).
func parseExpiry(r oauthResp) (time.Time, error) {
	if r.AccessTokenExpiredAt != "" {
		t, err := time.ParseInLocation("2006-01-02 15:04:05", r.AccessTokenExpiredAt, seoulTZ)
		if err == nil {
			return t, nil
		}
	}
	if r.ExpiresIn > 0 {
		return time.Now().Add(time.Duration(r.ExpiresIn) * time.Second), nil
	}
	return time.Time{}, errors.New("kis: token response missing expiry")
}
```

- [ ] **Step 4: 테스트 실행 → PASS**

Run: `go test ./internal/token/... -v -run TestManager`
Expected: 4 PASS

- [ ] **Step 5: Commit**

```bash
git add internal/token/manager.go internal/token/manager_test.go
git commit -m "[feat] token.Manager OAuth 발급 + singleflight + 자동 캐시

만료 5분 전 선제 발급, 동시 호출 시 1번만 발급, Asia/Seoul 시간대 파싱.

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>"
```

---

## Task 8: internal/mastercache — KOSPI/KOSDAQ ZIP 디스크 캐시

**Files:**
- Create: `internal/mastercache/cache.go`
- Create: `internal/mastercache/cache_test.go`

- [ ] **Step 1: 테스트 작성**

```go
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
```

- [ ] **Step 2: 테스트 실행 → FAIL**

Run: `go test ./internal/mastercache/... -v`
Expected: 컴파일 실패

- [ ] **Step 3: 구현** — `internal/mastercache/cache.go`

```go
// Package mastercache 는 KRX 종목 마스터 파일 (KOSPI/KOSDAQ ZIP) 의 디스크 캐시.
//
// 다운로드 비용이 큰 마스터 파일을 1주일 단위로 재사용. 다운로드 실패 시 옛 캐시 fallback.
package mastercache

import (
	"context"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// FetchFunc 는 캐시 미스 시 호출되는 다운로드 함수.
type FetchFunc func(ctx context.Context) ([]byte, error)

// Cache 는 디스크 기반 file cache.
type Cache struct {
	dir string
	ttl time.Duration
	mu  sync.Mutex
}

// New 는 지정된 디렉터리와 TTL 로 Cache 생성.
// dir 가 빈 문자열이면 DefaultDir() 사용.
func New(dir string, ttl time.Duration) *Cache {
	if dir == "" {
		d, err := DefaultDir()
		if err == nil {
			dir = d
		}
	}
	return &Cache{dir: dir, ttl: ttl}
}

// DefaultDir 은 OS 별 기본 캐시 디렉터리. macOS: ~/Library/Caches/kis, Linux: ~/.cache/kis.
func DefaultDir() (string, error) {
	base, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(base, "kis"), nil
}

// Get 은 name 으로 캐시 조회. miss 또는 TTL 만료 시 fetch 호출 후 저장.
// fetch 실패하고 옛 캐시가 있으면 옛 캐시 반환.
func (c *Cache) Get(ctx context.Context, name string, fetch FetchFunc) ([]byte, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	path := filepath.Join(c.dir, name)
	info, statErr := os.Stat(path)
	hasCache := statErr == nil
	hot := hasCache && time.Since(info.ModTime()) < c.ttl

	if hot {
		return os.ReadFile(path)
	}

	// fetch (cold or expired)
	data, fetchErr := fetch(ctx)
	if fetchErr != nil {
		if hasCache {
			// fallback to stale
			return os.ReadFile(path)
		}
		return nil, fetchErr
	}

	if err := os.MkdirAll(c.dir, 0700); err != nil {
		return nil, err
	}
	if err := os.WriteFile(path, data, 0600); err != nil {
		return nil, err
	}
	return data, nil
}

// Clear 는 name 캐시 제거. 없으면 에러 없음.
func (c *Cache) Clear(name string) error {
	err := os.Remove(filepath.Join(c.dir, name))
	if errors.Is(err, fs.ErrNotExist) {
		return nil
	}
	return err
}
```

- [ ] **Step 4: 테스트 실행 → PASS**

Run: `go test ./internal/mastercache/... -v -race`
Expected: 5 PASS

- [ ] **Step 5: Commit**

```bash
git add internal/mastercache/
git commit -m "[feat] internal/mastercache 디스크 캐시 추가

KOSPI/KOSDAQ ZIP 다운로드 캐시. TTL 만료 시 재다운로드, 실패 시 stale fallback.

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>"
```

---

## Task 9: internal/httpclient/client.go — resty wrapper + Do

**Files:**
- Create: `internal/httpclient/client.go`
- Create: `internal/httpclient/client_test.go`

- [ ] **Step 1: 테스트 작성**

```go
package httpclient

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kenshin579/korea-investment-stock/internal/ratelimit"
	"github.com/kenshin579/korea-investment-stock/internal/token"
)

type stubTokenMgr struct {
	bearer string
	calls  atomic.Int64
}

func (s *stubTokenMgr) Get(ctx context.Context) (string, error) {
	s.calls.Add(1)
	return s.bearer, nil
}

func (s *stubTokenMgr) Refresh(ctx context.Context) (string, error) {
	s.calls.Add(1)
	return s.bearer + "-refreshed", nil
}

func newTestClient(t *testing.T, tm TokenManager) *Client {
	t.Helper()
	c := New(Config{
		BaseURL:     "https://openapi.test",
		AppKey:      "ak",
		AppSecret:   "as",
		AccountNo:   "12345678-01",
		Limiter:     ratelimit.New(1000),
		TokenMgr:    tm,
		Retries:     2,
	})
	httpmock.ActivateNonDefault(c.resty.GetClient())
	t.Cleanup(httpmock.DeactivateAndReset)
	return c
}

func TestClient_Do_Success(t *testing.T) {
	tm := &stubTokenMgr{bearer: "Bearer T"}
	c := newTestClient(t, tm)
	httpmock.RegisterResponder(http.MethodGet, "=~/inquire-price",
		httpmock.NewStringResponder(200, `{"rt_cd":"0","msg_cd":"OK","msg1":"ok","output":{"x":"1"}}`))

	resp, err := c.Do(context.Background(), &Request{
		Method: http.MethodGet,
		Path:   "/uapi/domestic-stock/v1/quotations/inquire-price",
		TrID:   "FHKST01010100",
	})
	require.NoError(t, err)
	assert.Equal(t, "0", resp.RtCode)
	assert.Equal(t, int64(1), tm.calls.Load(), "token Get called once")
}

func TestClient_Do_APIError(t *testing.T) {
	tm := &stubTokenMgr{bearer: "Bearer T"}
	c := newTestClient(t, tm)
	httpmock.RegisterResponder(http.MethodGet, "=~/inquire-price",
		httpmock.NewStringResponder(200, `{"rt_cd":"1","msg_cd":"MCA00001","msg1":"잘못된 종목"}`))

	_, err := c.Do(context.Background(), &Request{
		Method: http.MethodGet,
		Path:   "/inquire-price",
		TrID:   "FHKST01010100",
	})
	require.Error(t, err)
	var apiErr *APIError
	require.ErrorAs(t, err, &apiErr)
	assert.Equal(t, "MCA00001", apiErr.MsgCode)
}

func TestClient_Do_TokenExpiredAutoRetry(t *testing.T) {
	tm := &stubTokenMgr{bearer: "Bearer T"}
	c := newTestClient(t, tm)
	calls := atomic.Int64{}
	httpmock.RegisterResponder(http.MethodGet, "=~/inquire-price",
		func(req *http.Request) (*http.Response, error) {
			n := calls.Add(1)
			if n == 1 {
				return httpmock.NewStringResponse(200, `{"rt_cd":"1","msg_cd":"EGW00123","msg1":"기간이 만료된 token 입니다"}`), nil
			}
			return httpmock.NewStringResponse(200, `{"rt_cd":"0","msg_cd":"OK","msg1":"ok"}`), nil
		})

	resp, err := c.Do(context.Background(), &Request{
		Method: http.MethodGet,
		Path:   "/inquire-price",
		TrID:   "FHKST01010100",
	})
	require.NoError(t, err)
	assert.Equal(t, "0", resp.RtCode)
	assert.GreaterOrEqual(t, tm.calls.Load(), int64(2), "token Get + Refresh both called")
}

func TestClient_Do_Retry5xx(t *testing.T) {
	tm := &stubTokenMgr{bearer: "Bearer T"}
	c := newTestClient(t, tm)
	calls := atomic.Int64{}
	httpmock.RegisterResponder(http.MethodGet, "=~/inquire-price",
		func(req *http.Request) (*http.Response, error) {
			n := calls.Add(1)
			if n < 2 {
				return httpmock.NewStringResponse(503, `service unavailable`), nil
			}
			return httpmock.NewStringResponse(200, `{"rt_cd":"0","msg_cd":"OK","msg1":"ok"}`), nil
		})

	resp, err := c.Do(context.Background(), &Request{
		Method: http.MethodGet, Path: "/inquire-price", TrID: "FHKST01010100",
	})
	require.NoError(t, err)
	assert.Equal(t, "0", resp.RtCode)
	assert.Equal(t, int64(2), calls.Load(), "5xx → retry")
}

func TestClient_Do_TokenError(t *testing.T) {
	tm := &errorTokenMgr{err: errors.New("oauth down")}
	c := newTestClient(t, tm)
	_, err := c.Do(context.Background(), &Request{
		Method: http.MethodGet, Path: "/inquire-price", TrID: "FHKST01010100",
	})
	require.Error(t, err)
	assert.True(t, strings.Contains(err.Error(), "oauth down"))
}

type errorTokenMgr struct{ err error }

func (e *errorTokenMgr) Get(ctx context.Context) (string, error)     { return "", e.err }
func (e *errorTokenMgr) Refresh(ctx context.Context) (string, error) { return "", e.err }

var _ TokenManager = (*stubTokenMgr)(nil)
var _ TokenManager = (*errorTokenMgr)(nil)
var _ = token.AccessToken{} // import 활용
```

- [ ] **Step 2: 테스트 실행 → FAIL**

Run: `go test ./internal/httpclient/... -v`
Expected: 컴파일 실패

- [ ] **Step 3: 구현** — `internal/httpclient/client.go`

```go
// Package httpclient 은 한투 API 호출용 resty 래퍼.
//
// rate limit / token / 재시도 / 에러 정규화를 한 곳에서 처리.
// 사용자에게 노출되지 않는 internal 패키지.
package httpclient

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/kenshin579/korea-investment-stock/internal/ratelimit"
)

// TokenManager 는 token.Manager 의 인터페이스 추상화 (테스트 편의).
type TokenManager interface {
	Get(ctx context.Context) (string, error)
	Refresh(ctx context.Context) (string, error)
}

// Config 는 Client 생성 옵션.
type Config struct {
	BaseURL    string
	AppKey     string
	AppSecret  string
	AccountNo  string
	Limiter    *ratelimit.Limiter
	TokenMgr   TokenManager
	Retries    int
	Timeout    time.Duration
	UserAgent  string
	HTTPClient *http.Client
}

// Client 는 한투 API 호출 단일 진입점.
type Client struct {
	cfg   Config
	resty *resty.Client
}

// New 는 Config 로 Client 생성.
func New(cfg Config) *Client {
	if cfg.Retries < 0 {
		cfg.Retries = 0
	}
	if cfg.Timeout <= 0 {
		cfg.Timeout = 30 * time.Second
	}
	if cfg.UserAgent == "" {
		cfg.UserAgent = "korea-investment-stock-go"
	}
	r := resty.New().
		SetBaseURL(strings.TrimRight(cfg.BaseURL, "/")).
		SetTimeout(cfg.Timeout).
		SetHeader("User-Agent", cfg.UserAgent)
	if cfg.HTTPClient != nil {
		r.SetTransport(cfg.HTTPClient.Transport)
	}
	return &Client{cfg: cfg, resty: r}
}

// Request 는 단일 한투 API 호출 요청.
type Request struct {
	Method string            // http.MethodGet 등
	Path   string            // "/uapi/..." (BaseURL 제외)
	TrID   string            // tr_id 헤더 (한투 transaction ID)
	Query  map[string]string // GET 쿼리 파라미터
	Body   any               // POST body (JSON 직렬화)
	// CustType: P (개인), B (법인). 빈 문자열이면 미지정.
	CustType string
}

// Response 는 한투 API 응답을 정규화한 결과.
type Response struct {
	RtCode  string          `json:"rt_cd"`
	MsgCode string          `json:"msg_cd"`
	Msg1    string          `json:"msg1"`
	Output  json.RawMessage `json:"output"`
	Output1 json.RawMessage `json:"output1"`
	Output2 json.RawMessage `json:"output2"`
	Raw     []byte          `json:"-"`
}

// APIError 는 한투 응답의 rt_cd != "0" 케이스. kis.APIError 와 동일 구조 (root 와 internal 분리).
type APIError struct {
	RtCode  string
	MsgCode string
	Message string
	TrID    string
}

func (e *APIError) Error() string {
	return "kis: API error [" + e.MsgCode + "] " + e.Message
}

// Do 는 단일 호출 + 재시도 + 토큰 만료 자동 재발급.
func (c *Client) Do(ctx context.Context, req *Request) (*Response, error) {
	if err := c.cfg.Limiter.Wait(ctx); err != nil {
		return nil, err
	}

	bearer, err := c.cfg.TokenMgr.Get(ctx)
	if err != nil {
		return nil, err
	}

	resp, err := c.send(ctx, req, bearer)
	if err != nil {
		return nil, err
	}

	// 토큰 만료 → 1회 재발급 후 재시도
	if resp != nil && isTokenExpired(resp) {
		newBearer, refErr := c.cfg.TokenMgr.Refresh(ctx)
		if refErr != nil {
			return nil, refErr
		}
		resp, err = c.send(ctx, req, newBearer)
		if err != nil {
			return nil, err
		}
	}

	if resp.RtCode != "0" {
		return nil, &APIError{RtCode: resp.RtCode, MsgCode: resp.MsgCode, Message: resp.Msg1, TrID: req.TrID}
	}
	return resp, nil
}

func (c *Client) send(ctx context.Context, req *Request, bearer string) (*Response, error) {
	var lastHTTPErr error
	for attempt := 0; attempt <= c.cfg.Retries; attempt++ {
		r := c.resty.R().
			SetContext(ctx).
			SetHeader("Authorization", bearer).
			SetHeader("appkey", c.cfg.AppKey).
			SetHeader("appsecret", c.cfg.AppSecret).
			SetHeader("tr_id", req.TrID).
			SetHeader("Content-Type", "application/json; charset=utf-8")
		if req.CustType != "" {
			r.SetHeader("custtype", req.CustType)
		}
		if len(req.Query) > 0 {
			r.SetQueryParams(req.Query)
		}
		if req.Body != nil {
			r.SetBody(req.Body)
		}

		httpResp, err := r.Execute(req.Method, req.Path)
		if err != nil {
			lastHTTPErr = err
			if attempt == c.cfg.Retries {
				return nil, fmt.Errorf("kis: http: %w", err)
			}
			time.Sleep(backoff(attempt))
			continue
		}

		if httpResp.StatusCode() >= 500 || httpResp.StatusCode() == http.StatusTooManyRequests {
			lastHTTPErr = fmt.Errorf("HTTP %d", httpResp.StatusCode())
			if attempt == c.cfg.Retries {
				return nil, fmt.Errorf("kis: http: %s after %d retries", lastHTTPErr, c.cfg.Retries)
			}
			time.Sleep(backoff(attempt))
			continue
		}

		raw := httpResp.Body()
		var resp Response
		if err := json.Unmarshal(raw, &resp); err != nil {
			return nil, fmt.Errorf("kis: parse: %w (body=%s)", err, string(raw))
		}
		resp.Raw = raw
		return &resp, nil
	}
	return nil, errors.New("unreachable")
}

func isTokenExpired(r *Response) bool {
	// 한투는 만료 시 msg_cd 가 EGW00123 또는 메시지에 "기간이 만료된 token" 포함.
	return r.MsgCode == "EGW00123" || strings.Contains(r.Msg1, "기간이 만료된 token")
}

func backoff(attempt int) time.Duration {
	// 0.5s, 1s, 2s, ...
	d := 500 * time.Millisecond
	for i := 0; i < attempt; i++ {
		d *= 2
	}
	return d
}
```

- [ ] **Step 4: 테스트 실행 → PASS**

Run: `go test ./internal/httpclient/... -v -race`
Expected: 5 PASS

- [ ] **Step 5: Commit**

```bash
git add internal/httpclient/client.go internal/httpclient/client_test.go
git commit -m "[feat] httpclient.Client resty 래퍼 추가

rate limit + token 자동 주입 + 만료 자동 재발급 + 5xx/429 재시도.

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>"
```

---

## Task 10: internal/httpclient/hashkey.go — Hashkey 발급

**Files:**
- Create: `internal/httpclient/hashkey.go`
- Create: `internal/httpclient/hashkey_test.go`

- [ ] **Step 1: 테스트 작성**

```go
package httpclient

import (
	"context"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_Hashkey(t *testing.T) {
	tm := &stubTokenMgr{bearer: "Bearer T"}
	c := newTestClient(t, tm)
	httpmock.RegisterResponder(http.MethodPost, "=~/uapi/hashkey",
		httpmock.NewStringResponder(200, `{"HASH":"abcdef","BODY":{}}`))

	hk, err := c.Hashkey(context.Background(), map[string]string{"k": "v"})
	require.NoError(t, err)
	assert.Equal(t, "abcdef", hk)
}
```

- [ ] **Step 2: 테스트 실행 → FAIL**

Run: `go test ./internal/httpclient/... -v -run TestClient_Hashkey`
Expected: 컴파일 실패

- [ ] **Step 3: 구현** — `internal/httpclient/hashkey.go`

```go
package httpclient

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const hashkeyPath = "/uapi/hashkey"

type hashkeyResp struct {
	Hash string `json:"HASH"`
}

// Hashkey 는 한투의 hashkey 엔드포인트로 body 를 보내고 hash 문자열 반환.
// 주문 등 일부 API 가 요구.
func (c *Client) Hashkey(ctx context.Context, body any) (string, error) {
	httpResp, err := c.resty.R().
		SetContext(ctx).
		SetHeader("appkey", c.cfg.AppKey).
		SetHeader("appsecret", c.cfg.AppSecret).
		SetHeader("Content-Type", "application/json; charset=utf-8").
		SetBody(body).
		Execute(http.MethodPost, hashkeyPath)
	if err != nil {
		return "", fmt.Errorf("kis: hashkey: %w", err)
	}
	if httpResp.StatusCode() != http.StatusOK {
		return "", fmt.Errorf("kis: hashkey: HTTP %d", httpResp.StatusCode())
	}
	var r hashkeyResp
	if err := json.Unmarshal(httpResp.Body(), &r); err != nil {
		return "", fmt.Errorf("kis: hashkey parse: %w", err)
	}
	return r.Hash, nil
}
```

- [ ] **Step 4: 테스트 실행 → PASS**

Run: `go test ./internal/httpclient/... -v -run TestClient_Hashkey`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add internal/httpclient/hashkey.go internal/httpclient/hashkey_test.go
git commit -m "[feat] httpclient.Hashkey 추가

주문 등 일부 한투 API 가 요구하는 hashkey 발급.

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>"
```

---

## Task 11: options.go — functional options

**Files:**
- Modify: `client.go` (이미 Phase 0 에서 작성된 `clientOptions` 보강)
- Create: `options.go`
- Create: `options_test.go`

- [ ] **Step 1: 테스트 작성** — `options_test.go`

```go
package kis

import (
	"log/slog"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestOptions_Apply(t *testing.T) {
	cfg := &clientOptions{}
	httpC := &http.Client{}
	logger := slog.Default()

	WithBaseURL("https://example.com")(cfg)
	WithRetries(5)(cfg)
	WithRateLimit(20)(cfg)
	WithHTTPClient(httpC)(cfg)
	WithLogger(logger)(cfg)
	WithTimeout(10 * time.Second)(cfg)
	WithUserAgent("test-ua")(cfg)
	WithMasterCacheDir("/tmp/kis")(cfg)

	assert.Equal(t, "https://example.com", cfg.baseURL)
	assert.Equal(t, 5, cfg.retries)
	assert.Equal(t, 20.0, cfg.rateLimit)
	assert.Same(t, httpC, cfg.httpClient)
	assert.Same(t, logger, cfg.logger)
	assert.Equal(t, 10*time.Second, cfg.timeout)
	assert.Equal(t, "test-ua", cfg.userAgent)
	assert.Equal(t, "/tmp/kis", cfg.masterCacheDir)
}

func TestWithPaperEnv(t *testing.T) {
	cfg := &clientOptions{}
	WithPaperEnv()(cfg)
	assert.Equal(t, PaperEnv, cfg.baseURL)
}
```

- [ ] **Step 2: 테스트 실행 → FAIL**

Run: `go test ./... -run TestOptions -v`
Expected: 컴파일 실패 (대부분 With 함수 미정의)

- [ ] **Step 3: client.go 의 clientOptions struct 보강**

`client.go` 의 `clientOptions` 부분을 다음으로 교체:

```go
type clientOptions struct {
	baseURL        string
	retries        int
	rateLimit      float64
	httpClient     *http.Client
	tokenStorage   token.Storage
	masterCacheDir string
	logger         *slog.Logger
	timeout        time.Duration
	userAgent      string
}
```

(필요한 import 추가: `log/slog`, `time`, `internal/token`. 다만 root 가 internal/token 에 의존하면 cycle 위험은 없는지 점검.)

- [ ] **Step 4: 구현** — `options.go`

```go
package kis

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/kenshin579/korea-investment-stock/internal/token"
)

// WithBaseURL 은 한투 API base URL 을 지정 (default: RealEnv).
func WithBaseURL(url string) Option {
	return func(o *clientOptions) { o.baseURL = url }
}

// WithPaperEnv 는 모의투자 API 엔드포인트로 변경.
func WithPaperEnv() Option {
	return func(o *clientOptions) { o.baseURL = PaperEnv }
}

// WithRetries 는 5xx/429 재시도 횟수 (default: 3).
func WithRetries(n int) Option {
	return func(o *clientOptions) { o.retries = n }
}

// WithRateLimit 은 호출/초 한도 (default: 15).
func WithRateLimit(rps float64) Option {
	return func(o *clientOptions) { o.rateLimit = rps }
}

// WithHTTPClient 는 사용자 정의 *http.Client 주입 (custom transport, proxy 등).
func WithHTTPClient(c *http.Client) Option {
	return func(o *clientOptions) { o.httpClient = c }
}

// WithTokenStorage 는 사용자 정의 토큰 저장소 (default: FileStorage at ~/.cache/kis/token.json).
func WithTokenStorage(s token.Storage) Option {
	return func(o *clientOptions) { o.tokenStorage = s }
}

// WithMasterCacheDir 는 KOSPI/KOSDAQ 마스터 파일 캐시 디렉터리.
func WithMasterCacheDir(dir string) Option {
	return func(o *clientOptions) { o.masterCacheDir = dir }
}

// WithLogger 는 사용자 정의 slog logger.
func WithLogger(l *slog.Logger) Option {
	return func(o *clientOptions) { o.logger = l }
}

// WithTimeout 은 단일 HTTP 호출의 timeout (default: 30s).
func WithTimeout(d time.Duration) Option {
	return func(o *clientOptions) { o.timeout = d }
}

// WithUserAgent 는 User-Agent 헤더 (default: "korea-investment-stock-go").
func WithUserAgent(ua string) Option {
	return func(o *clientOptions) { o.userAgent = ua }
}
```

- [ ] **Step 5: 테스트 실행 → PASS**

Run: `go test ./... -run TestOptions -v && go vet ./... && gofmt -l .`
Expected: PASS, vet/fmt clean

- [ ] **Step 6: Commit**

```bash
git add options.go options_test.go client.go
git commit -m "[feat] options.go 9개 functional options + clientOptions 보강

WithBaseURL, WithPaperEnv, WithRetries, WithRateLimit, WithHTTPClient,
WithTokenStorage, WithMasterCacheDir, WithLogger, WithTimeout, WithUserAgent.

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>"
```

---

## Task 12: config.go — Config struct + LoadConfigFromEnv + LoadConfigFromYAML

**Files:**
- Create: `config.go`
- Create: `config_test.go`

- [ ] **Step 1: 테스트 작성** — `config_test.go`

```go
package kis

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadConfigFromEnv_Required(t *testing.T) {
	t.Setenv("KOREA_INVESTMENT_API_KEY", "k")
	t.Setenv("KOREA_INVESTMENT_API_SECRET", "s")
	t.Setenv("KOREA_INVESTMENT_ACCOUNT_NO", "12345678-01")

	cfg, err := LoadConfigFromEnv()
	require.NoError(t, err)
	assert.Equal(t, "k", cfg.APIKey)
	assert.Equal(t, "s", cfg.APISecret)
	assert.Equal(t, "12345678-01", cfg.AccountNo)
}

func TestLoadConfigFromEnv_MissingKey(t *testing.T) {
	t.Setenv("KOREA_INVESTMENT_API_KEY", "")
	t.Setenv("KOREA_INVESTMENT_API_SECRET", "s")
	t.Setenv("KOREA_INVESTMENT_ACCOUNT_NO", "x")
	_, err := LoadConfigFromEnv()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "KOREA_INVESTMENT_API_KEY")
}

func TestLoadConfigFromEnv_Optional(t *testing.T) {
	t.Setenv("KOREA_INVESTMENT_API_KEY", "k")
	t.Setenv("KOREA_INVESTMENT_API_SECRET", "s")
	t.Setenv("KOREA_INVESTMENT_ACCOUNT_NO", "x")
	t.Setenv("KOREA_INVESTMENT_BASE_URL", "https://custom")
	t.Setenv("KOREA_INVESTMENT_TOKEN_STORAGE", "redis")
	t.Setenv("KOREA_INVESTMENT_REDIS_URL", "redis://1.2.3.4:6379/0")

	cfg, err := LoadConfigFromEnv()
	require.NoError(t, err)
	assert.Equal(t, "https://custom", cfg.BaseURL)
	assert.Equal(t, "redis", cfg.TokenStorage)
	assert.Equal(t, "redis://1.2.3.4:6379/0", cfg.RedisURL)
}

func TestLoadConfigFromYAML(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	require.NoError(t, os.WriteFile(path, []byte(`
api_key: yk
api_secret: ys
acc_no: "98765432-01"
base_url: https://yaml-base
token_storage_type: file
token_file: /tmp/yaml-token.json
rate_limit: 12.5
retries: 5
`), 0600))

	cfg, err := LoadConfigFromYAML(path)
	require.NoError(t, err)
	assert.Equal(t, "yk", cfg.APIKey)
	assert.Equal(t, "ys", cfg.APISecret)
	assert.Equal(t, "98765432-01", cfg.AccountNo)
	assert.Equal(t, "https://yaml-base", cfg.BaseURL)
	assert.Equal(t, "file", cfg.TokenStorage)
	assert.Equal(t, "/tmp/yaml-token.json", cfg.TokenFile)
	assert.Equal(t, 12.5, cfg.RateLimit)
	assert.Equal(t, 5, cfg.Retries)
}

func TestLoadConfigFromYAML_FileNotFound(t *testing.T) {
	_, err := LoadConfigFromYAML("/nonexistent/file.yaml")
	require.Error(t, err)
}
```

- [ ] **Step 2: 테스트 실행 → FAIL**

Run: `go test ./... -run TestLoadConfig -v`
Expected: 컴파일 실패

- [ ] **Step 3: 구현** — `config.go`

```go
package kis

import (
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config 는 라이브러리 설정. NewClientFromEnv / NewClientFromYAML 의 base.
type Config struct {
	APIKey         string  `yaml:"api_key"`
	APISecret      string  `yaml:"api_secret"`
	AccountNo      string  `yaml:"acc_no"`
	BaseURL        string  `yaml:"base_url"`
	TokenStorage   string  `yaml:"token_storage_type"` // "file" | "redis"
	TokenFile      string  `yaml:"token_file"`
	RedisURL       string  `yaml:"redis_url"`
	RedisPassword  string  `yaml:"redis_password"`
	MasterCacheDir string  `yaml:"master_cache_dir"`
	RateLimit      float64 `yaml:"rate_limit"`
	Retries        int     `yaml:"retries"`
}

// LoadConfigFromEnv 는 KOREA_INVESTMENT_* 환경변수에서 Config 로드.
// 필수: API_KEY, API_SECRET, ACCOUNT_NO.
func LoadConfigFromEnv() (*Config, error) {
	required := func(key string) (string, error) {
		v := os.Getenv(key)
		if v == "" {
			return "", fmt.Errorf("kis: env var %s is required", key)
		}
		return v, nil
	}

	apiKey, err := required("KOREA_INVESTMENT_API_KEY")
	if err != nil {
		return nil, err
	}
	apiSecret, err := required("KOREA_INVESTMENT_API_SECRET")
	if err != nil {
		return nil, err
	}
	accNo, err := required("KOREA_INVESTMENT_ACCOUNT_NO")
	if err != nil {
		return nil, err
	}

	return &Config{
		APIKey:        apiKey,
		APISecret:     apiSecret,
		AccountNo:     accNo,
		BaseURL:       os.Getenv("KOREA_INVESTMENT_BASE_URL"),
		TokenStorage:  os.Getenv("KOREA_INVESTMENT_TOKEN_STORAGE"),
		TokenFile:     os.Getenv("KOREA_INVESTMENT_TOKEN_FILE"),
		RedisURL:      os.Getenv("KOREA_INVESTMENT_REDIS_URL"),
		RedisPassword: os.Getenv("KOREA_INVESTMENT_REDIS_PASSWORD"),
	}, nil
}

// LoadConfigFromYAML 는 YAML 파일에서 Config 로드.
func LoadConfigFromYAML(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("kis: read config: %w", err)
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("kis: parse config: %w", err)
	}
	if cfg.APIKey == "" || cfg.APISecret == "" || cfg.AccountNo == "" {
		return nil, errors.New("kis: api_key, api_secret, acc_no are required in config")
	}
	return &cfg, nil
}
```

- [ ] **Step 4: 테스트 실행 → PASS**

Run: `go test ./... -run TestLoadConfig -v`
Expected: 5 PASS

- [ ] **Step 5: Commit**

```bash
git add config.go config_test.go
git commit -m "[feat] Config + LoadConfigFromEnv + LoadConfigFromYAML

KOREA_INVESTMENT_* env vars 자동 감지 + YAML 파일 파싱.

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>"
```

---

## Task 13: domestic/client.go + overseas/client.go — placeholder Client struct

**Files:**
- Create: `domestic/client.go`
- Create: `overseas/client.go`
- Modify: `domestic/doc.go`
- Modify: `overseas/doc.go`

- [ ] **Step 1: domestic/client.go 작성**

```go
package domestic

import "github.com/kenshin579/korea-investment-stock/internal/httpclient"

// Client 는 국내주식 API sub-client. Phase 1.2 부터 메서드 추가.
//
// 사용자는 직접 생성하지 않고 kis.Client.Domestic 으로 접근.
type Client struct {
	http *httpclient.Client
}

// New 는 internal 용도. root kis.NewClient 가 호출.
func New(http *httpclient.Client) *Client {
	return &Client{http: http}
}
```

- [ ] **Step 2: overseas/client.go 작성**

```go
package overseas

import "github.com/kenshin579/korea-investment-stock/internal/httpclient"

// Client 는 해외주식 API sub-client. Phase 1.5 부터 메서드 추가.
//
// 사용자는 직접 생성하지 않고 kis.Client.Overseas 로 접근.
type Client struct {
	http *httpclient.Client
}

// New 는 internal 용도. root kis.NewClient 가 호출.
func New(http *httpclient.Client) *Client {
	return &Client{http: http}
}
```

- [ ] **Step 3: doc.go 갱신**

`domestic/doc.go`:
```go
// Package domestic provides methods and response types for domestic (Korean) stock APIs.
//
// Phase 0 placeholder + Phase 1.1 client struct. 메서드는 Phase 1.2 부터 추가.
package domestic
```

`overseas/doc.go`:
```go
// Package overseas provides methods and response types for overseas (US/HK/JP/CN/VN) stock APIs.
//
// Phase 0 placeholder + Phase 1.1 client struct. 메서드는 Phase 1.5 부터 추가.
package overseas
```

- [ ] **Step 4: 컴파일 검증**

Run: `go build ./...`
Expected: 출력 없음 (성공)

- [ ] **Step 5: Commit**

```bash
git add domestic overseas
git commit -m "[feat] domestic.Client + overseas.Client 추가 (Phase 1.1 placeholder)

shared httpclient 주입 받는 sub-client struct. 메서드는 Phase 1.2 부터.

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>"
```

---

## Task 14: client.go — NewClient 보강 (인프라 wiring)

**Files:**
- Modify: `client.go`
- Create: `client_test.go`

- [ ] **Step 1: 테스트 작성** — `client_test.go`

```go
package kis

import (
	"context"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewClient_RequiresAllParams(t *testing.T) {
	_, err := NewClient("", "s", "acc")
	assert.Error(t, err)
	_, err = NewClient("k", "", "acc")
	assert.Error(t, err)
	_, err = NewClient("k", "s", "")
	assert.Error(t, err)
}

func TestNewClient_AppliesOptions(t *testing.T) {
	c, err := NewClient("k", "s", "acc",
		WithBaseURL("https://x"),
		WithRetries(7),
		WithRateLimit(20),
	)
	require.NoError(t, err)
	assert.Equal(t, "https://x", c.opts.baseURL)
	assert.Equal(t, 7, c.opts.retries)
	assert.Equal(t, 20.0, c.opts.rateLimit)
	require.NotNil(t, c.Domestic)
	require.NotNil(t, c.Overseas)
}

func TestNewClient_Defaults(t *testing.T) {
	c, err := NewClient("k", "s", "acc")
	require.NoError(t, err)
	assert.Equal(t, RealEnv, c.opts.baseURL)
	assert.Equal(t, 3, c.opts.retries)
	assert.Equal(t, 15.0, c.opts.rateLimit)
}

func TestNewClient_TokenIssue(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder(http.MethodPost, "=~/oauth2/tokenP",
		httpmock.NewStringResponder(200, `{
			"access_token": "T",
			"token_type": "Bearer",
			"access_token_token_expired": "2099-12-31 23:59:59"
		}`))

	c, err := NewClient("k", "s", "acc",
		WithHTTPClient(&http.Client{Transport: httpmock.DefaultTransport}),
	)
	require.NoError(t, err)
	tok, err := c.IssueAccessToken(context.Background())
	require.NoError(t, err)
	assert.Equal(t, "Bearer T", tok)
}
```

- [ ] **Step 2: 테스트 실행 → FAIL**

Run: `go test ./... -run TestNewClient -v`
Expected: 컴파일 실패 또는 missing method

- [ ] **Step 3: client.go 보강** — 다음 내용으로 전체 교체

```go
// Package kis is a Go client for the Korea Investment Securities OpenAPI.
//
// See docs/superpowers/specs/2026-05-03-korea-investment-go-migration-design.md
// for the design rationale, and docs/superpowers/specs/2026-05-03-phase1-api-coverage-design.md
// for Phase 1 scope.
package kis

import (
	"errors"
	"log/slog"
	"net/http"
	"path/filepath"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/kenshin579/korea-investment-stock/domestic"
	"github.com/kenshin579/korea-investment-stock/internal/httpclient"
	"github.com/kenshin579/korea-investment-stock/internal/mastercache"
	"github.com/kenshin579/korea-investment-stock/internal/ratelimit"
	"github.com/kenshin579/korea-investment-stock/internal/token"
	"github.com/kenshin579/korea-investment-stock/overseas"
)

// KIS OpenAPI base URLs.
const (
	RealEnv  = "https://openapi.koreainvestment.com:9443"
	PaperEnv = "https://openapivts.koreainvestment.com:29443"
)

// Client 는 kis 라이브러리의 단일 진입점.
type Client struct {
	apiKey    string
	apiSecret string
	accountNo string
	opts      clientOptions

	httpClient *httpclient.Client
	tokenMgr   *token.Manager
	masterC    *mastercache.Cache

	Domestic *domestic.Client
	Overseas *overseas.Client
}

// Option 은 functional option.
type Option func(*clientOptions)

type clientOptions struct {
	baseURL        string
	retries        int
	rateLimit      float64
	httpClient     *http.Client
	tokenStorage   token.Storage
	masterCacheDir string
	logger         *slog.Logger
	timeout        time.Duration
	userAgent      string
}

// NewClient 는 kis Client 생성 (직접 credentials 전달).
func NewClient(apiKey, apiSecret, accountNo string, opts ...Option) (*Client, error) {
	if apiKey == "" || apiSecret == "" || accountNo == "" {
		return nil, errors.New("kis: apiKey, apiSecret, and accountNo are required and must not be empty")
	}

	cfg := defaultOptions()
	for _, opt := range opts {
		opt(&cfg)
	}

	c := &Client{
		apiKey:    apiKey,
		apiSecret: apiSecret,
		accountNo: accountNo,
		opts:      cfg,
	}
	if err := c.wireInfra(); err != nil {
		return nil, err
	}
	c.Domestic = domestic.New(c.httpClient)
	c.Overseas = overseas.New(c.httpClient)
	return c, nil
}

func defaultOptions() clientOptions {
	return clientOptions{
		baseURL:   RealEnv,
		retries:   3,
		rateLimit: 15,
		timeout:   30 * time.Second,
		userAgent: "korea-investment-stock-go",
	}
}

func (c *Client) wireInfra() error {
	storage, err := c.resolveTokenStorage()
	if err != nil {
		return err
	}

	c.tokenMgr = token.NewManager(token.Config{
		Storage:    storage,
		BaseURL:    c.opts.baseURL,
		APIKey:     c.apiKey,
		APISecret:  c.apiSecret,
		HTTPClient: c.opts.httpClient,
	})

	c.httpClient = httpclient.New(httpclient.Config{
		BaseURL:    c.opts.baseURL,
		AppKey:     c.apiKey,
		AppSecret:  c.apiSecret,
		AccountNo:  c.accountNo,
		Limiter:    ratelimit.New(c.opts.rateLimit),
		TokenMgr:   c.tokenMgr,
		Retries:    c.opts.retries,
		Timeout:    c.opts.timeout,
		UserAgent:  c.opts.userAgent,
		HTTPClient: c.opts.httpClient,
	})

	masterDir := c.opts.masterCacheDir
	if masterDir == "" {
		d, _ := mastercache.DefaultDir()
		masterDir = d
	}
	c.masterC = mastercache.New(masterDir, 7*24*time.Hour)
	return nil
}

func (c *Client) resolveTokenStorage() (token.Storage, error) {
	if c.opts.tokenStorage != nil {
		return c.opts.tokenStorage, nil
	}
	// default: FileStorage at user cache dir
	dir, err := mastercache.DefaultDir() // 같은 위치 재사용
	if err != nil {
		return nil, err
	}
	return token.NewFileStorage(filepath.Join(dir, "token.json")), nil
}

// resolveTokenStorageRedis 는 NewClientFromEnv 에서 redis 옵션 시 호출.
// (현재 client.go 에는 사용되지 않음 — from_env.go 에서 import)
func newRedisStorage(url, password string) (token.Storage, error) {
	opts, err := redis.ParseURL(url)
	if err != nil {
		return nil, err
	}
	if password != "" {
		opts.Password = password
	}
	rdb := redis.NewClient(opts)
	return token.NewRedisStorage(rdb, "kis:token:default"), nil
}
```

- [ ] **Step 4: 테스트 실행 → PASS**

Run: `go test ./... -run TestNewClient -v && go build ./... && go vet ./...`
Expected: 4 PASS, build/vet clean

- [ ] **Step 5: Commit**

```bash
git add client.go client_test.go
git commit -m "[feat] NewClient 인프라 wiring 보강

token manager, http client, rate limiter, master cache 자동 wiring.
Domestic / Overseas sub-client 자동 주입.

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>"
```

---

## Task 15: auth.go — IssueAccessToken 외부 노출

**Files:**
- Create: `auth.go`

- [ ] **Step 1: 구현** — `auth.go`

```go
package kis

import "context"

// IssueAccessToken 은 OAuth 토큰을 발급하고 Bearer 문자열 반환.
// 일반적으로 사용자가 명시 호출할 필요는 없음 — 라이브러리가 자동 발급.
// 디버깅 목적이나 사전 warmup 시에 유용.
func (c *Client) IssueAccessToken(ctx context.Context) (string, error) {
	return c.tokenMgr.Get(ctx)
}

// RefreshAccessToken 은 캐시 무시하고 강제로 새 토큰 발급.
func (c *Client) RefreshAccessToken(ctx context.Context) (string, error) {
	return c.tokenMgr.Refresh(ctx)
}
```

- [ ] **Step 2: 검증**

`TestNewClient_TokenIssue` (Task 14) 가 이미 IssueAccessToken 호출 → 그 테스트가 통과하는지 확인.

Run: `go test ./... -run TestNewClient_TokenIssue -v`
Expected: PASS

- [ ] **Step 3: Commit**

```bash
git add auth.go
git commit -m "[feat] auth.go IssueAccessToken / RefreshAccessToken 외부 노출

디버깅/warmup 용. 일반 사용은 자동 발급에 의존.

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>"
```

---

## Task 16: from_env.go — NewClientFromEnv

**Files:**
- Create: `from_env.go`
- Create: `from_env_test.go`

- [ ] **Step 1: 테스트 작성** — `from_env_test.go`

```go
package kis

import (
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewClientFromEnv_Success(t *testing.T) {
	t.Setenv("KOREA_INVESTMENT_API_KEY", "k")
	t.Setenv("KOREA_INVESTMENT_API_SECRET", "s")
	t.Setenv("KOREA_INVESTMENT_ACCOUNT_NO", "12345678-01")

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder(http.MethodPost, "=~/oauth2/tokenP",
		httpmock.NewStringResponder(200, `{"access_token":"x","token_type":"Bearer","access_token_token_expired":"2099-12-31 23:59:59"}`))

	c, err := NewClientFromEnv(WithHTTPClient(&http.Client{Transport: httpmock.DefaultTransport}))
	require.NoError(t, err)
	assert.NotNil(t, c.Domestic)
}

func TestNewClientFromEnv_MissingEnv(t *testing.T) {
	t.Setenv("KOREA_INVESTMENT_API_KEY", "")
	t.Setenv("KOREA_INVESTMENT_API_SECRET", "s")
	t.Setenv("KOREA_INVESTMENT_ACCOUNT_NO", "x")
	_, err := NewClientFromEnv()
	require.Error(t, err)
}
```

- [ ] **Step 2: 테스트 실행 → FAIL**

Run: `go test ./... -run TestNewClientFromEnv -v`
Expected: 컴파일 실패

- [ ] **Step 3: 구현** — `from_env.go`

```go
package kis

import "fmt"

// NewClientFromEnv 는 KOREA_INVESTMENT_* 환경변수에서 credentials/설정 자동 감지 후 Client 생성.
// 추가 옵션은 functional options 로 override 가능.
func NewClientFromEnv(opts ...Option) (*Client, error) {
	cfg, err := LoadConfigFromEnv()
	if err != nil {
		return nil, err
	}
	return newFromConfig(cfg, opts...)
}

func newFromConfig(cfg *Config, opts ...Option) (*Client, error) {
	// Config → options 로 변환 (옵션이 마지막에 override)
	baseOpts := []Option{}
	if cfg.BaseURL != "" {
		baseOpts = append(baseOpts, WithBaseURL(cfg.BaseURL))
	}
	if cfg.RateLimit > 0 {
		baseOpts = append(baseOpts, WithRateLimit(cfg.RateLimit))
	}
	if cfg.Retries > 0 {
		baseOpts = append(baseOpts, WithRetries(cfg.Retries))
	}
	if cfg.MasterCacheDir != "" {
		baseOpts = append(baseOpts, WithMasterCacheDir(cfg.MasterCacheDir))
	}

	if cfg.TokenStorage == "redis" && cfg.RedisURL != "" {
		s, err := newRedisStorage(cfg.RedisURL, cfg.RedisPassword)
		if err != nil {
			return nil, fmt.Errorf("kis: redis storage: %w", err)
		}
		baseOpts = append(baseOpts, WithTokenStorage(s))
	}

	allOpts := append(baseOpts, opts...)
	return NewClient(cfg.APIKey, cfg.APISecret, cfg.AccountNo, allOpts...)
}
```

- [ ] **Step 4: 테스트 실행 → PASS**

Run: `go test ./... -run TestNewClientFromEnv -v`
Expected: 2 PASS

- [ ] **Step 5: Commit**

```bash
git add from_env.go from_env_test.go
git commit -m "[feat] NewClientFromEnv 진입점

KOREA_INVESTMENT_* env vars + functional options 조합 (options override).

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>"
```

---

## Task 17: from_yaml.go — NewClientFromYAML

**Files:**
- Create: `from_yaml.go`
- Create: `from_yaml_test.go`

- [ ] **Step 1: 테스트 작성** — `from_yaml_test.go`

```go
package kis

import (
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewClientFromYAML(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	require.NoError(t, os.WriteFile(path, []byte(`
api_key: yk
api_secret: ys
acc_no: "98765432-01"
`), 0600))

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder(http.MethodPost, "=~/oauth2/tokenP",
		httpmock.NewStringResponder(200, `{"access_token":"x","token_type":"Bearer","access_token_token_expired":"2099-12-31 23:59:59"}`))

	c, err := NewClientFromYAML(path,
		WithHTTPClient(&http.Client{Transport: httpmock.DefaultTransport}))
	require.NoError(t, err)
	assert.NotNil(t, c.Domestic)
}

func TestNewClientFromYAML_NotFound(t *testing.T) {
	_, err := NewClientFromYAML("/nonexistent.yaml")
	require.Error(t, err)
}
```

- [ ] **Step 2: 테스트 실행 → FAIL**

Run: `go test ./... -run TestNewClientFromYAML -v`
Expected: 컴파일 실패

- [ ] **Step 3: 구현** — `from_yaml.go`

```go
package kis

// NewClientFromYAML 은 YAML 파일에서 credentials/설정 로드 후 Client 생성.
// 추가 옵션은 functional options 로 override 가능.
func NewClientFromYAML(path string, opts ...Option) (*Client, error) {
	cfg, err := LoadConfigFromYAML(path)
	if err != nil {
		return nil, err
	}
	return newFromConfig(cfg, opts...)
}
```

- [ ] **Step 4: 테스트 실행 → PASS**

Run: `go test ./... -run TestNewClientFromYAML -v`
Expected: 2 PASS

- [ ] **Step 5: Commit**

```bash
git add from_yaml.go from_yaml_test.go
git commit -m "[feat] NewClientFromYAML 진입점

YAML 파일 로드 후 functional options 와 조합.

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>"
```

---

## Task 18: examples/

**Files:**
- Create: `examples/basic/main.go`
- Create: `examples/env_config/main.go`
- Create: `examples/yaml_config/main.go`
- Create: `examples/yaml_config/config.yaml`

- [ ] **Step 1: examples/basic/main.go**

```go
// Basic example: NewClient + IssueAccessToken.
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	kis "github.com/kenshin579/korea-investment-stock"
)

func main() {
	client, err := kis.NewClient(
		os.Getenv("KOREA_INVESTMENT_API_KEY"),
		os.Getenv("KOREA_INVESTMENT_API_SECRET"),
		os.Getenv("KOREA_INVESTMENT_ACCOUNT_NO"),
	)
	if err != nil {
		log.Fatal(err)
	}
	bearer, err := client.IssueAccessToken(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Token issued:", bearer[:20]+"...")
}
```

- [ ] **Step 2: examples/env_config/main.go**

```go
// env_config example: NewClientFromEnv.
package main

import (
	"context"
	"fmt"
	"log"

	kis "github.com/kenshin579/korea-investment-stock"
)

func main() {
	client, err := kis.NewClientFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	bearer, err := client.IssueAccessToken(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Token issued from env:", bearer[:20]+"...")
}
```

- [ ] **Step 3: examples/yaml_config/main.go**

```go
// yaml_config example: NewClientFromYAML.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	kis "github.com/kenshin579/korea-investment-stock"
)

func main() {
	path := flag.String("config", "./config.yaml", "path to config.yaml")
	flag.Parse()

	client, err := kis.NewClientFromYAML(*path)
	if err != nil {
		log.Fatal(err)
	}
	bearer, err := client.IssueAccessToken(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Token issued from yaml:", bearer[:20]+"...")
}
```

- [ ] **Step 4: examples/yaml_config/config.yaml**

```yaml
# Sample config for korea-investment-stock Go library.
# Copy and edit before running examples/yaml_config/main.go.

api_key: YOUR_API_KEY
api_secret: YOUR_API_SECRET
acc_no: "12345678-01"

# Optional
# base_url: https://openapi.koreainvestment.com:9443
# token_storage_type: file
# token_file: ~/.cache/kis/token.json
# rate_limit: 15
# retries: 3
```

- [ ] **Step 5: 컴파일 검증**

Run: `go build ./examples/... && echo OK`
Expected: `OK` (모든 example 컴파일 성공)

- [ ] **Step 6: Commit**

```bash
git add examples/
git commit -m "[feat] examples/ 추가 — basic, env_config, yaml_config

Phase 1.1 인프라 동작 확인용 minimal examples.

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>"
```

---

## Task 19: 최종 점검

- [ ] **Step 1: 빌드/vet/fmt**

Run:
```bash
go build ./... && go vet ./... && gofmt -l . | tee /tmp/fmt.out
```
Expected: 빌드/vet 출력 없음. `gofmt -l` 출력 빈 파일.

- [ ] **Step 2: 모든 테스트 통과 + race**

Run: `go test ./... -race -count=1`
Expected: 모든 패키지 PASS

- [ ] **Step 3: Coverage 측정**

Run:
```bash
go test ./... -coverprofile=/tmp/cov.out -covermode=atomic
go tool cover -func=/tmp/cov.out | tail -5
```
Expected: 마지막 줄 `total: (statements) ...` 의 비율이 ≥ 80%. 만약 부족하면 어느 패키지가 낮은지 확인 후 추가 테스트 (다음 task 에서).

- [ ] **Step 4: 디렉터리 구조 확인**

Run:
```bash
find . -name '*.go' -not -path './.git/*' -not -path './examples/*' | sort
```
Expected: 다음 패턴이 모두 존재 (개수 일치 안 해도 OK, 패턴만 확인)

```
./auth.go
./client.go
./client_test.go
./config.go
./config_test.go
./domestic/client.go
./domestic/doc.go
./errors.go
./errors_test.go
./from_env.go
./from_env_test.go
./from_yaml.go
./from_yaml_test.go
./internal/httpclient/client.go
./internal/httpclient/client_test.go
./internal/httpclient/hashkey.go
./internal/httpclient/hashkey_test.go
./internal/mastercache/cache.go
./internal/mastercache/cache_test.go
./internal/ratelimit/limiter.go
./internal/ratelimit/limiter_test.go
./internal/token/file_storage.go
./internal/token/file_storage_test.go
./internal/token/manager.go
./internal/token/manager_test.go
./internal/token/redis_storage.go
./internal/token/redis_storage_test.go
./internal/token/storage.go
./internal/token/storage_test.go
./options.go
./options_test.go
./overseas/client.go
./overseas/doc.go
```

- [ ] **Step 5: 커밋 history 확인**

Run: `git log main..HEAD --oneline | wc -l`
Expected: ≥ 18 commits (각 task 1 commit + plan/spec 커밋들).

(이 task 는 fix 가 필요하면 수행 후 commit; 이상 없으면 commit 없이 다음 단계로.)

---

## Task 20: PR 생성 및 push (사용자 승인 후)

> Claude 는 push / PR 생성을 사용자 명시적 승인 후에만 실행 (글로벌 정책).

- [ ] **Step 1: 사용자 승인 요청**

Claude 가 사용자에게:
- "Phase 1.1 모든 task 완료 + 최종 점검 통과. 지금 push + PR 생성하라" 또는
- "수정 사항 있다"

응답 받기.

- [ ] **Step 2: 승인 시 push**

Run:
```bash
git push -u origin docs/phase1-api-coverage-spec
```

- [ ] **Step 3: PR 생성** (gh CLI + HEREDOC, 사용자 정책)

```bash
gh pr create --base main --title "Phase 1.1: 인프라 + Config (v0.1.0)" --body "$(cat <<'EOF'
## Summary

`korea-investment-stock` Go 라이브러리의 첫 사용 가능 release (v0.1.0) — 인프라와 Config 진입점만. 메서드는 Phase 1.2 부터.

### 포함 내용

- **Phase 1 design spec** + **Phase 1.1 implementation plan** (`docs/superpowers/specs/2026-05-03-phase1-*.md`)
- **errors.go** — `APIError` + sentinel 에러 (`ErrTokenExpired`, `ErrRateLimited`, ...)
- **internal/ratelimit** — 토큰 버킷 rate limiter (ctx 지원, thread-safe, 통계)
- **internal/token** — Storage 인터페이스 + FileStorage + RedisStorage + Manager (OAuth + singleflight + 자동 갱신)
- **internal/httpclient** — resty wrap + 토큰 자동 주입 + 5xx/429 재시도 + 토큰 만료 자동 재발급 + Hashkey
- **internal/mastercache** — KOSPI/KOSDAQ ZIP 디스크 캐시 (TTL + stale fallback)
- **options.go** — 9 functional options (WithBaseURL/Retries/RateLimit/HTTPClient/TokenStorage/MasterCacheDir/Logger/Timeout/UserAgent)
- **config.go** — `Config` struct + `LoadConfigFromEnv` + `LoadConfigFromYAML`
- **client.go** — NewClient 인프라 wiring 보강 (httpclient, tokenMgr, masterCache, sub-client 자동 주입)
- **from_env.go / from_yaml.go** — 3 진입점 중 두 개
- **auth.go** — `IssueAccessToken` / `RefreshAccessToken` 외부 노출
- **domestic.Client / overseas.Client** — Phase 1.2~1.5 가 메서드를 채울 placeholder
- **examples/** — basic / env_config / yaml_config

### 검증

- ✅ `go build ./...`
- ✅ `go vet ./...`
- ✅ `gofmt -l .`
- ✅ `go test ./... -race -count=1`
- ✅ Coverage ≥ 80%

### 다음 단계

- **Phase 1.2** (`v0.2.0`): 국내 시세 + 심볼 + 차트 (8 메서드)
- **Phase 1.3 ~ 1.5** 순차 진행
- 모두 끝나면 **`v1.0.0`** = Python parity 완성

## Test plan

- [x] 모든 단위 테스트 통과
- [x] race detector 통과
- [ ] examples 실행 검증 (사용자 manual, 실제 KIS credentials 필요)
- [ ] PR merge 후 v0.1.0 태그
EOF
)"
```

- [ ] **Step 4: PR URL 사용자에게 전달**

Run: `gh pr view --json url --jq '.url'`

PR URL 사용자에게 보여주고 review 안내.

- [ ] **Step 5: PR merge 후 v0.1.0 태그 (사용자가 승인 후)**

PR merge 후:
```bash
git checkout main && git pull
git tag -a v0.1.0 -m "v0.1.0: 인프라 + Config (Phase 1.1)"
git push origin v0.1.0
```

---

## Self-Review

### 1. Spec coverage

| Phase 1 spec 요구사항 | 구현 task |
|---|---|
| Rate limiter (token bucket, ctx 지원, thread-safe, 통계) | Task 3 |
| Token manager (OAuth, 5분 선제 발급, singleflight) | Task 7 |
| Token Storage (FileStorage, RedisStorage) | Task 5, 6 |
| HTTP client (resty wrap, 재시도, 토큰 자동 재발급) | Task 9 |
| Hashkey | Task 10 |
| Master file cache | Task 8 |
| auth.go (외부 노출 토큰 발급) | Task 15 |
| Functional options 9개 | Task 11 |
| errors.go (APIError + sentinel) | Task 2 |
| Config + LoadFromEnv + LoadFromYAML | Task 12 |
| NewClient + NewClientFromEnv + NewClientFromYAML | Task 14, 16, 17 |
| domestic/overseas placeholder Client | Task 13 |
| Examples | Task 18 |
| 검증 | Task 19, 20 |

### 2. Placeholder scan

- 모든 step 에 실제 코드 / 명령어 / 검증 포함
- "TBD", "TODO", "implement later" 사용 없음
- (Task 19 의 coverage check 가 부족 시 "추가 테스트 다음 task 에서" 라는 모호한 표현이 있음 — 실제는 그 자리에서 처리해야 하지만, 내용이 사전에 모두 정의 안 되므로 implementer 가 케이스별 판단)

### 3. Type consistency

- `Client` struct 필드 (`apiKey`, `apiSecret`, `accountNo`, `opts`, `httpClient`, `tokenMgr`, `masterC`, `Domestic`, `Overseas`) 가 Task 14 에 정의됨. 후속 task 에서 동일하게 사용.
- `clientOptions` 필드들 (`baseURL`, `retries`, `rateLimit`, `httpClient`, `tokenStorage`, `masterCacheDir`, `logger`, `timeout`, `userAgent`) 가 Task 11 에 정의됨. Task 14 와 일치.
- `httpclient.Config` 필드명 일관성 OK.
- `token.Config` 필드명 일관성 OK.

### 4. 위험 / 결함

- Task 14 의 `wireInfra` 가 `c.tokenMgr` 를 만들 때 `c.opts.httpClient` (사용자 *http.Client) 를 그대로 넘김. `httpclient.New` 도 같은 `c.opts.httpClient` 사용. 두 곳에서 동일 transport 공유 — 의도. 단 사용자 transport 가 nil 이면 `http.DefaultClient` 사용 (token.NewManager 가 기본 처리).
- Task 14 의 `newRedisStorage` 가 `client.go` 에 정의되어 있는데 root 가 `redis` 패키지를 import 하면 라이브러리 의존성에 redis 강제 추가됨. 이는 의도 (토큰 storage 옵션). 사용자가 redis 안 쓰면 단지 import 만 남고 binary 에 미사용 코드는 linker 가 제거 — 실제 사이즈 영향 미미.

(self-review 통과)
