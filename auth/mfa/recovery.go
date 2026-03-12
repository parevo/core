package mfa

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"
)

var (
	ErrRecoveryCodeInvalid = errors.New("recovery code invalid")
	ErrRecoveryCodeUsed    = errors.New("recovery code already used")
)

// RecoveryStore persists and validates recovery codes.
type RecoveryStore interface {
	GetHashes(ctx context.Context, userID string) ([]string, error)
	SetHashes(ctx context.Context, userID string, hashes []string) error
	Consume(ctx context.Context, userID string, hash string) (bool, error)
}

// RecoveryService generates and verifies MFA recovery codes.
type RecoveryService struct {
	store RecoveryStore
	count int
}

// NewRecoveryService creates a recovery code service.
func NewRecoveryService(store RecoveryStore, codeCount int) *RecoveryService {
	if codeCount <= 0 {
		codeCount = 10
	}
	return &RecoveryService{store: store, count: codeCount}
}

// Generate creates new recovery codes. Returns plain codes (show once) and stores hashes.
func (s *RecoveryService) Generate(ctx context.Context, userID string) ([]string, error) {
	codes := make([]string, s.count)
	hashes := make([]string, s.count)
	for i := 0; i < s.count; i++ {
		code := generateCode()
		codes[i] = code
		hashes[i] = hashCode(code)
	}
	if err := s.store.SetHashes(ctx, userID, hashes); err != nil {
		return nil, err
	}
	return codes, nil
}

// Verify consumes a recovery code if valid.
func (s *RecoveryService) Verify(ctx context.Context, userID, code string) error {
	code = normalizeCode(code)
	if code == "" {
		return ErrRecoveryCodeInvalid
	}
	h := hashCode(code)
	ok, err := s.store.Consume(ctx, userID, h)
	if err != nil {
		return err
	}
	if !ok {
		return ErrRecoveryCodeInvalid
	}
	return nil
}

func generateCode() string {
	b := make([]byte, 4)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

func hashCode(code string) string {
	h := sha256.Sum256([]byte(normalizeCode(code)))
	return hex.EncodeToString(h[:])
}

func normalizeCode(s string) string {
	return strings.ToLower(strings.ReplaceAll(strings.TrimSpace(s), "-", ""))
}
