package apikey

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"
)

var (
	ErrInvalidAPIKey = errors.New("invalid api key")
	ErrKeyNotFound   = errors.New("api key not found")
)

type APIKeyStore interface {
	Validate(ctx context.Context, keyHash, prefix string) (userID, tenantID string, err error)
}

type Service struct {
	store APIKeyStore
}

func NewService(store APIKeyStore) *Service {
	return &Service{store: store}
}

func (s *Service) Validate(ctx context.Context, rawKey string) (userID, tenantID string, err error) {
	rawKey = strings.TrimSpace(strings.TrimPrefix(rawKey, "Bearer"))
	if rawKey == "" {
		return "", "", ErrInvalidAPIKey
	}
	if !strings.HasPrefix(rawKey, "pk_") && !strings.HasPrefix(rawKey, "sk_") {
		return "", "", ErrInvalidAPIKey
	}
	prefix := rawKey[:8]
	if s.store == nil {
		return "", "", ErrKeyNotFound
	}
	hash := hashKey(rawKey)
	return s.store.Validate(ctx, hash, prefix)
}

func hashKey(key string) string {
	h := sha256.Sum256([]byte(key))
	return hex.EncodeToString(h[:])
}

func GenerateKey(prefix string) (string, error) {
	b := make([]byte, 24)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return prefix + hex.EncodeToString(b), nil
}
