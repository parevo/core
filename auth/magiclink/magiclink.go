package magiclink

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"
)

var (
	ErrTokenExpired = errors.New("magic link token expired")
	ErrTokenInvalid = errors.New("magic link token invalid")
)

// TokenStore persists magic link tokens.
type TokenStore interface {
	Create(ctx context.Context, email, token string, expiresAt time.Time) error
	Consume(ctx context.Context, token string) (email string, err error)
}

// Mailer sends magic link emails.
type Mailer interface {
	SendMagicLink(ctx context.Context, email, link string) error
}

// Config holds magic link settings.
type Config struct {
	TTL       time.Duration
	TokenLen  int
	BaseURL   string
	Path      string
}

// DefaultConfig returns sensible defaults.
func DefaultConfig(baseURL string) Config {
	return Config{
		TTL:      15 * time.Minute,
		TokenLen: 32,
		BaseURL:  baseURL,
		Path:     "/auth/magic-link/verify",
	}
}

// Service handles magic link generation and verification.
type Service struct {
	store  TokenStore
	mailer Mailer
	cfg    Config
}

// NewService creates a magic link service.
func NewService(store TokenStore, mailer Mailer, cfg Config) *Service {
	if cfg.TTL <= 0 {
		cfg.TTL = 15 * time.Minute
	}
	if cfg.TokenLen <= 0 {
		cfg.TokenLen = 32
	}
	return &Service{store: store, mailer: mailer, cfg: cfg}
}

// SendLink generates a token, stores it, and sends the email.
func (s *Service) SendLink(ctx context.Context, email string) (token string, err error) {
	b := make([]byte, s.cfg.TokenLen)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	token = hex.EncodeToString(b)
	expiresAt := time.Now().Add(s.cfg.TTL)
	if err := s.store.Create(ctx, email, token, expiresAt); err != nil {
		return "", err
	}
	link := fmt.Sprintf("%s%s?token=%s", s.cfg.BaseURL, s.cfg.Path, token)
	if s.mailer != nil {
		if err := s.mailer.SendMagicLink(ctx, email, link); err != nil {
			return "", err
		}
	}
	return token, nil
}

// Verify consumes the token and returns the email.
func (s *Service) Verify(ctx context.Context, token string) (email string, err error) {
	email, err = s.store.Consume(ctx, token)
	if err != nil {
		return "", err
	}
	if email == "" {
		return "", ErrTokenInvalid
	}
	return email, nil
}
