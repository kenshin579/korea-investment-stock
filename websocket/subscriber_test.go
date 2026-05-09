package websocket

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSubscriber_Add(t *testing.T) {
	s := newSubscriber()
	added, err := s.Add("H0STCNT0", "005930")
	assert.NoError(t, err)
	assert.True(t, added)

	// 중복
	added, err = s.Add("H0STCNT0", "005930")
	assert.NoError(t, err)
	assert.False(t, added) // 이미 존재 → false
}

func TestSubscriber_Remove(t *testing.T) {
	s := newSubscriber()
	s.Add("H0STCNT0", "005930")
	removed := s.Remove("H0STCNT0", "005930")
	assert.True(t, removed)

	removed = s.Remove("H0STCNT0", "005930")
	assert.False(t, removed) // 이미 없음
}

func TestSubscriber_All(t *testing.T) {
	s := newSubscriber()
	s.Add("H0STCNT0", "005930")
	s.Add("H0STCNT0", "000660")
	s.Add("H0STASP0", "005930")

	all := s.All()
	assert.Len(t, all, 3)

	// 정렬되지 않은 list — set 비교
	keys := map[string]bool{}
	for _, sub := range all {
		keys[sub.TrID+":"+sub.TrKey] = true
	}
	assert.True(t, keys["H0STCNT0:005930"])
	assert.True(t, keys["H0STCNT0:000660"])
	assert.True(t, keys["H0STASP0:005930"])
}

func TestSubscriber_Concurrent(t *testing.T) {
	s := newSubscriber()
	done := make(chan struct{}, 100)
	for i := 0; i < 100; i++ {
		go func(i int) {
			s.Add("H0STCNT0", "00593"+string(rune('0'+i%10)))
			done <- struct{}{}
		}(i)
	}
	for i := 0; i < 100; i++ {
		<-done
	}
	// race detector 가 -race 모드에서 panic 없이 통과하면 OK
	assert.NotEmpty(t, s.All())
}
