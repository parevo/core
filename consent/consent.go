package consent

import (
	"context"
	"errors"
	"time"
)

var (
	ErrConsentRequired = errors.New("user consent required")
	ErrConsentDenied   = errors.New("user denied consent")
)

// Scope represents an OAuth2 scope the user can consent to.
type Scope struct {
	ID          string
	Name        string
	Description string
}

// ConsentRecord stores user consent for a client's scopes.
type ConsentRecord struct {
	UserID    string
	ClientID  string
	Scopes    []string
	GrantedAt time.Time
}

// Store persists consent records.
type Store interface {
	Get(ctx context.Context, userID, clientID string) (*ConsentRecord, error)
	Save(ctx context.Context, record ConsentRecord) error
	Revoke(ctx context.Context, userID, clientID string) error
}

// Service manages OAuth2 consent.
type Service struct {
	store  Store
	scopes map[string]Scope
}

// NewService creates a consent service.
func NewService(store Store, scopes []Scope) *Service {
	m := make(map[string]Scope)
	for _, s := range scopes {
		m[s.ID] = s
	}
	return &Service{store: store, scopes: m}
}

// Check returns nil if user has consented to all requested scopes.
func (s *Service) Check(ctx context.Context, userID, clientID string, requestedScopes []string) error {
	rec, err := s.store.Get(ctx, userID, clientID)
	if err != nil {
		return err
	}
	if rec == nil {
		return ErrConsentRequired
	}
	allowed := make(map[string]bool)
	for _, sc := range rec.Scopes {
		allowed[sc] = true
	}
	for _, sc := range requestedScopes {
		if !allowed[sc] {
			return ErrConsentRequired
		}
	}
	return nil
}

// Grant records user consent.
func (s *Service) Grant(ctx context.Context, userID, clientID string, scopes []string) error {
	return s.store.Save(ctx, ConsentRecord{
		UserID:    userID,
		ClientID:  clientID,
		Scopes:    scopes,
		GrantedAt: time.Now(),
	})
}

// Revoke removes consent.
func (s *Service) Revoke(ctx context.Context, userID, clientID string) error {
	return s.store.Revoke(ctx, userID, clientID)
}
