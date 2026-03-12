package memory

import (
	"context"
	"sync"
	"time"
)

// MagicLinkStore implements magiclink.TokenStore in memory.
type MagicLinkStore struct {
	mu     sync.RWMutex
	tokens map[string]magicLinkEntry
}

type magicLinkEntry struct {
	Email     string
	ExpiresAt time.Time
}

// NewMagicLinkStore creates an in-memory magic link store.
func NewMagicLinkStore() *MagicLinkStore {
	return &MagicLinkStore{tokens: make(map[string]magicLinkEntry)}
}

func (s *MagicLinkStore) Create(_ context.Context, email, token string, expiresAt time.Time) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tokens[token] = magicLinkEntry{Email: email, ExpiresAt: expiresAt}
	return nil
}

func (s *MagicLinkStore) Consume(_ context.Context, token string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	e, ok := s.tokens[token]
	if !ok {
		return "", nil
	}
	if time.Now().After(e.ExpiresAt) {
		delete(s.tokens, token)
		return "", nil
	}
	delete(s.tokens, token)
	return e.Email, nil
}
