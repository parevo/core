//go:build !webauthn
// +build !webauthn

package webauthn

import (
	"context"
	"encoding/json"
)

// NewService returns a stub that returns ErrWebAuthnNotSupported.
// Build with -tags webauthn and add github.com/go-webauthn/webauthn for full support.
func NewService(_, _ string, _ CredentialStore) (Service, error) {
	return &stubService{}, nil
}

type stubService struct{}

func (stubService) BeginRegistration(context.Context, RegistrationOptions) (*RegistrationSession, error) {
	return nil, ErrWebAuthnNotSupported
}

func (stubService) FinishRegistration(context.Context, *RegistrationSession, json.RawMessage) error {
	return ErrWebAuthnNotSupported
}

func (stubService) BeginAssertion(context.Context, AssertionOptions) (*AssertionSession, error) {
	return nil, ErrWebAuthnNotSupported
}

func (stubService) FinishAssertion(context.Context, *AssertionSession, json.RawMessage) (string, error) {
	return "", ErrWebAuthnNotSupported
}
