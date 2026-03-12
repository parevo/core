package memory

import (
	"context"
	"sync"
	"time"

	"github.com/parevo/core/storage"
)

// SessionMetadataStore extends SessionStore with metadata.
type SessionMetadataStore struct {
	*SessionStore
	mu       sync.RWMutex
	metadata map[string]storage.SessionMetadata
}

// NewSessionMetadataStore creates a session store with metadata.
func NewSessionMetadataStore() *SessionMetadataStore {
	return &SessionMetadataStore{
		SessionStore: &SessionStore{
			Revoked:       make(map[string]bool),
			SessionToUser: make(map[string]string),
			UserSessions:  make(map[string]map[string]struct{}),
		},
		metadata: make(map[string]storage.SessionMetadata),
	}
}

func (s *SessionMetadataStore) SetMetadata(_ context.Context, sessionID, userID, ip, userAgent string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := time.Now()
	m := s.metadata[sessionID]
	m.SessionID = sessionID
	m.UserID = userID
	m.IP = ip
	m.UserAgent = userAgent
	m.LastActivity = now
	if m.CreatedAt.IsZero() {
		m.CreatedAt = now
	}
	s.metadata[sessionID] = m
	return nil
}

func (s *SessionMetadataStore) UpdateActivity(_ context.Context, sessionID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if m, ok := s.metadata[sessionID]; ok {
		m.LastActivity = time.Now()
		s.metadata[sessionID] = m
	}
	return nil
}

func (s *SessionMetadataStore) ListWithMetadata(_ context.Context, userID string) ([]storage.SessionMetadata, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var out []storage.SessionMetadata
	for _, m := range s.metadata {
		if m.UserID == userID {
			out = append(out, m)
		}
	}
	return out, nil
}
