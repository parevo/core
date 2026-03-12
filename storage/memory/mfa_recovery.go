package memory

import (
	"context"
	"sync"
)

// MFARecoveryStore implements mfa.RecoveryStore in memory.
type MFARecoveryStore struct {
	mu     sync.RWMutex
	hashes map[string][]string // userID -> hashes
	used   map[string]map[string]bool // userID -> hash -> true
}

// NewMFARecoveryStore creates an in-memory recovery store.
func NewMFARecoveryStore() *MFARecoveryStore {
	return &MFARecoveryStore{
		hashes: make(map[string][]string),
		used:   make(map[string]map[string]bool),
	}
}

func (s *MFARecoveryStore) GetHashes(_ context.Context, userID string) ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return append([]string(nil), s.hashes[userID]...), nil
}

func (s *MFARecoveryStore) SetHashes(_ context.Context, userID string, hashes []string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.hashes[userID] = append([]string(nil), hashes...)
	s.used[userID] = make(map[string]bool)
	return nil
}

func (s *MFARecoveryStore) Consume(_ context.Context, userID, hash string) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.used[userID][hash] {
		return false, nil
	}
	for _, h := range s.hashes[userID] {
		if h == hash {
			if s.used[userID] == nil {
				s.used[userID] = make(map[string]bool)
			}
			s.used[userID][hash] = true
			return true, nil
		}
	}
	return false, nil
}
