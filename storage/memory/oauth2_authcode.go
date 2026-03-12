package memory

import (
	"context"
	"sync"
	"time"
)

type authCodeEntry struct {
	ClientID    string
	UserID      string
	RedirectURI string
	Scopes      []string
	ExpiresAt   time.Time
}

// OAuth2AuthCodeStore implements oauth2provider auth code storage in memory.
type OAuth2AuthCodeStore struct {
	mu    sync.RWMutex
	codes map[string]authCodeEntry
}

// NewOAuth2AuthCodeStore creates an in-memory auth code store.
func NewOAuth2AuthCodeStore() *OAuth2AuthCodeStore {
	return &OAuth2AuthCodeStore{codes: make(map[string]authCodeEntry)}
}

func (s *OAuth2AuthCodeStore) Create(_ context.Context, code, clientID, userID, redirectURI string, scopes []string, expiresAt time.Time) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.codes[code] = authCodeEntry{
		ClientID:    clientID,
		UserID:      userID,
		RedirectURI: redirectURI,
		Scopes:      append([]string(nil), scopes...),
		ExpiresAt:   expiresAt,
	}
	return nil
}

func (s *OAuth2AuthCodeStore) Consume(_ context.Context, code string) (clientID, userID, redirectURI string, scopes []string, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	e, ok := s.codes[code]
	if !ok {
		return "", "", "", nil, nil
	}
	delete(s.codes, code)
	if time.Now().After(e.ExpiresAt) {
		return "", "", "", nil, nil
	}
	return e.ClientID, e.UserID, e.RedirectURI, append([]string(nil), e.Scopes...), nil
}
