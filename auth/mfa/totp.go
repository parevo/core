package mfa

import (
	"context"
	"crypto/rand"
	"encoding/base32"
	"errors"
)

var (
	ErrTOTPInvalid = errors.New("invalid totp code")
	ErrTOTPDisabled = errors.New("totp not enabled")
)

type TOTPVerifier interface {
	Verify(secret, code string) bool
	GenerateSecret() (string, error)
}

type TOTPService struct {
	store   TOTPStore
	verifier TOTPVerifier
}

func NewTOTPService(store TOTPStore, verifier TOTPVerifier) *TOTPService {
	return &TOTPService{store: store, verifier: verifier}
}

func (s *TOTPService) Setup(ctx context.Context, userID string) (secret, qrURL string, err error) {
	if s.verifier == nil {
		return "", "", ErrTOTPInvalid
	}
	secret, err = s.verifier.GenerateSecret()
	if err != nil {
		return "", "", err
	}
	if err := s.store.SetSecret(ctx, userID, secret); err != nil {
		return "", "", err
	}
	qrURL = "otpauth://totp/Parevo:" + userID + "?secret=" + secret
	return secret, qrURL, nil
}

func (s *TOTPService) Verify(ctx context.Context, userID, code string) error {
	secret, enabled, err := s.store.GetSecret(ctx, userID)
	if err != nil {
		return err
	}
	if !enabled || secret == "" {
		return ErrTOTPDisabled
	}
	if !s.verifier.Verify(secret, code) {
		return ErrTOTPInvalid
	}
	return nil
}

func (s *TOTPService) Enable(ctx context.Context, userID string) error {
	return s.store.Enable(ctx, userID)
}

func (s *TOTPService) Disable(ctx context.Context, userID string) error {
	return s.store.Disable(ctx, userID)
}

func GenerateTOTPSecret() (string, error) {
	b := make([]byte, 20)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(b), nil
}
