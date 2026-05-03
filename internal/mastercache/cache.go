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
