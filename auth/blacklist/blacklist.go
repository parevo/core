package blacklist

import (
	"context"
	"errors"
	"time"
)

var ErrTokenBlacklisted = errors.New("token is blacklisted")

// Store persists blacklisted JWT IDs (jti) or token hashes.
type Store interface {
	Add(ctx context.Context, jtiOrHash string, expiresAt time.Time) error
	Contains(ctx context.Context, jtiOrHash string) (bool, error)
}

// Service provides JWT blacklist for immediate revoke (logout).
type Service struct {
	store Store
}

// NewService creates a blacklist service.
func NewService(store Store) *Service {
	return &Service{store: store}
}

// Revoke adds a token to the blacklist until it expires.
func (s *Service) Revoke(ctx context.Context, jtiOrHash string, expiresAt time.Time) error {
	return s.store.Add(ctx, jtiOrHash, expiresAt)
}

// IsBlacklisted returns true if the token should be rejected.
func (s *Service) IsBlacklisted(ctx context.Context, jtiOrHash string) (bool, error) {
	return s.store.Contains(ctx, jtiOrHash)
}

// Check returns ErrTokenBlacklisted if the token is blacklisted.
func (s *Service) Check(ctx context.Context, jtiOrHash string) error {
	ok, err := s.IsBlacklisted(ctx, jtiOrHash)
	if err != nil {
		return err
	}
	if ok {
		return ErrTokenBlacklisted
	}
	return nil
}
