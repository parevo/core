package memory

import (
	"context"
	"sync"

	"github.com/parevo/core/auth/webauthn"
)

// WebAuthnCredentialStore implements webauthn.CredentialStore in memory.
type WebAuthnCredentialStore struct {
	mu   sync.RWMutex
	creds map[string][]webauthn.StoredCredential // userID -> credentials
}

// NewWebAuthnCredentialStore creates an in-memory credential store.
func NewWebAuthnCredentialStore() *WebAuthnCredentialStore {
	return &WebAuthnCredentialStore{creds: make(map[string][]webauthn.StoredCredential)}
}

func (s *WebAuthnCredentialStore) GetCredentials(_ context.Context, userID string) ([]webauthn.StoredCredential, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return append([]webauthn.StoredCredential(nil), s.creds[userID]...), nil
}

func (s *WebAuthnCredentialStore) SaveCredential(_ context.Context, userID string, cred webauthn.StoredCredential) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	list := s.creds[userID]
	for i, c := range list {
		if c.ID == cred.ID {
			list[i] = cred
			return nil
		}
	}
	s.creds[userID] = append(list, cred)
	return nil
}

func (s *WebAuthnCredentialStore) DeleteCredential(_ context.Context, userID, credentialID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	list := s.creds[userID]
	for i, c := range list {
		if c.ID == credentialID {
			s.creds[userID] = append(list[:i], list[i+1:]...)
			return nil
		}
	}
	return nil
}
