package memory

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"sync"

	"github.com/parevo/core/auth/apikey"
)

// Store implements apikey.APIKeyStore with in-memory storage for dev/test.
type Store struct {
	mu   sync.RWMutex
	keys map[string]struct{ userID, tenantID string }
}

// NewStore creates an in-memory API key store.
func NewStore() *Store {
	return &Store{
		keys: make(map[string]struct{ userID, tenantID string }),
	}
}

// Validate checks the key hash and returns userID, tenantID if valid.
func (s *Store) Validate(_ context.Context, keyHash, _ string) (userID, tenantID string, err error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if e, ok := s.keys[keyHash]; ok {
		return e.userID, e.tenantID, nil
	}
	return "", "", apikey.ErrKeyNotFound
}

// Add registers a raw API key for the given user and tenant.
// The key is hashed (SHA-256) before storage.
func (s *Store) Add(rawKey, userID, tenantID string) {
	hash := sha256.Sum256([]byte(rawKey))
	keyHash := hex.EncodeToString(hash[:])
	s.mu.Lock()
	defer s.mu.Unlock()
	s.keys[keyHash] = struct{ userID, tenantID string }{userID, tenantID}
}

// Remove removes a key by its hash (e.g. after revocation).
func (s *Store) Remove(keyHash string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.keys, keyHash)
}

var _ apikey.APIKeyStore = (*Store)(nil)
