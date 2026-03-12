package ldap

import (
	"context"
	"errors"
)

var (
	ErrLDAPNotConfigured = errors.New("LDAP not configured")
	ErrInvalidCredentials = errors.New("invalid LDAP credentials")
	ErrUserNotFound      = errors.New("LDAP user not found")
)

// Config holds LDAP connection settings.
type Config struct {
	URL      string
	BindDN   string
	BindPass string
	BaseDN   string
	UserFilter string // e.g. "(uid=%s)"
}

// User represents an LDAP user.
type User struct {
	DN       string
	Username string
	Email    string
	Name     string
	Groups   []string
}

// Store provides LDAP operations. Implementations use ldap.v3 or similar.
type Store interface {
	Authenticate(ctx context.Context, username, password string) (*User, error)
	FindByUsername(ctx context.Context, username string) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
}

// Service provides LDAP authentication.
type Service struct {
	store Store
}

// NewService creates an LDAP service.
func NewService(store Store) *Service {
	return &Service{store: store}
}

// Authenticate validates credentials and returns user info.
func (s *Service) Authenticate(ctx context.Context, username, password string) (*User, error) {
	if s.store == nil {
		return nil, ErrLDAPNotConfigured
	}
	return s.store.Authenticate(ctx, username, password)
}

// FindUser looks up a user by username.
func (s *Service) FindUser(ctx context.Context, username string) (*User, error) {
	if s.store == nil {
		return nil, ErrLDAPNotConfigured
	}
	return s.store.FindByUsername(ctx, username)
}
