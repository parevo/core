package memory

import (
	"context"
	"sync"

	"github.com/parevo/core/consent"
)

// ConsentStore implements consent.Store in memory.
type ConsentStore struct {
	mu   sync.RWMutex
	recs map[string]map[string]*consent.ConsentRecord // userID -> clientID -> record
}

// NewConsentStore creates an in-memory consent store.
func NewConsentStore() *ConsentStore {
	return &ConsentStore{recs: make(map[string]map[string]*consent.ConsentRecord)}
}

func (s *ConsentStore) Get(_ context.Context, userID, clientID string) (*consent.ConsentRecord, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if m, ok := s.recs[userID]; ok {
		if r, ok := m[clientID]; ok {
			cp := *r
			return &cp, nil
		}
	}
	return nil, nil
}

func (s *ConsentStore) Save(_ context.Context, record consent.ConsentRecord) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.recs[record.UserID] == nil {
		s.recs[record.UserID] = make(map[string]*consent.ConsentRecord)
	}
	r := record
	s.recs[record.UserID][record.ClientID] = &r
	return nil
}

func (s *ConsentStore) Revoke(_ context.Context, userID, clientID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if m, ok := s.recs[userID]; ok {
		delete(m, clientID)
	}
	return nil
}
