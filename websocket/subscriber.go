package websocket

import "sync"

// subKey 는 한 구독을 (tr_id, tr_key) 로 식별.
type subKey struct {
	TrID  string
	TrKey string
}

// subscriber 는 활성 구독을 thread-safe 하게 추적.
type subscriber struct {
	mu sync.RWMutex
	m  map[subKey]struct{}
}

func newSubscriber() *subscriber {
	return &subscriber{m: make(map[subKey]struct{})}
}

// Add 는 (tr_id, tr_key) 를 등록. 새로 추가했으면 true, 이미 존재했으면 false.
func (s *subscriber) Add(trID, trKey string) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	k := subKey{TrID: trID, TrKey: trKey}
	if _, exists := s.m[k]; exists {
		return false, nil
	}
	s.m[k] = struct{}{}
	return true, nil
}

// Remove 는 (tr_id, tr_key) 를 해제. 존재했으면 true.
func (s *subscriber) Remove(trID, trKey string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	k := subKey{TrID: trID, TrKey: trKey}
	if _, exists := s.m[k]; !exists {
		return false
	}
	delete(s.m, k)
	return true
}

// All 은 현재 활성 구독 list 를 snapshot 으로 반환 (reconnect 시 RestoreAll 용).
func (s *subscriber) All() []subKey {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]subKey, 0, len(s.m))
	for k := range s.m {
		out = append(out, k)
	}
	return out
}
