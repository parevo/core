package memory

import (
	"context"
	"sync"
	"time"
)

// BlacklistStore implements blacklist.Store in memory.
type BlacklistStore struct {
	mu   sync.RWMutex
	set  map[string]time.Time
}

// NewBlacklistStore creates an in-memory blacklist store.
func NewBlacklistStore() *BlacklistStore {
	return &BlacklistStore{set: make(map[string]time.Time)}
}

func (s *BlacklistStore) Add(_ context.Context, jtiOrHash string, expiresAt time.Time) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.set[jtiOrHash] = expiresAt
	return nil
}

func (s *BlacklistStore) Contains(_ context.Context, jtiOrHash string) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	exp, ok := s.set[jtiOrHash]
	if !ok {
		return false, nil
	}
	if time.Now().After(exp) {
		s.mu.RUnlock()
		s.mu.Lock()
		delete(s.set, jtiOrHash)
		s.mu.Unlock()
		s.mu.RLock()
		return false, nil
	}
	return true, nil
}
