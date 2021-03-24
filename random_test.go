package random

import (
	"sync"
	"testing"
	"time"
)

//SecureRandomAlphaString
func TestSecureRandomAlphaString(t *testing.T) {
	t.Logf("Retrieval Reference Generation single (%v)", time.Now())
	uid := SecureRandomAlphaString(12)
	t.Logf("Retrieval Reference sample %v", uid)
	if uid == "" {
		t.Error("failed to get uid")
	}
}

//SecureRandomAlphaString
func TestRandomness(t *testing.T) {
	count := 10000000
	t.Logf("Retrieval Reference Generation count:%v (%v)", count, time.Now())
	v := make(map[string]int)
	for n := 0; n < count; n++ {
		uid := SecureRandomAlphaString(12)
		if uid == "" {
			t.Error("failed to get uid")
		}
		if _, fail := v[uid]; fail {
			t.Errorf("failed to generate %v unique", count)
		}
		v[uid] = n
	}
	t.Logf("Completed count:%v unique generations", count)
}
func TestRandString(t *testing.T) {
	t.Logf("Retrieval Reference %v", RandString(12))
}

func TestInParallel(t *testing.T) {
	start := time.Now()
	var wait sync.WaitGroup
	sm := NewSyncMap()
	count := 1000000   // hash count x worker count
	worker_count := 10 // worker count
	for i := 0; i < worker_count; i++ {
		wait.Add(1)
		go func() {
			for n := 0; n < count; n++ {
				uid := SecureRandomAlphaString(12)
				if _, fail := sm.Get(uid); fail {
					t.Errorf("collission detected %v %v %v", count, worker_count, uid)
				}
				sm.Set(uid, n)
			}
			wait.Done()
		}()
	}
	wait.Wait()
	t.Logf("concurrent duration %v, count %v, no collissions", time.Now().Sub(start), count*worker_count)
}

type SyncMap struct {
	lock sync.RWMutex
	m    map[string]int
}

func NewSyncMap() *SyncMap {
	return &SyncMap{m: make(map[string]int)}
}

func (s *SyncMap) Get(key string) (int, bool) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	value, ok := s.m[key]
	return value, ok
}

func (s *SyncMap) Set(key string, value int) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.m[key] = value
}
