package webauthn

import (
	"context"
	"encoding/json"
	"errors"
)

var (
	ErrWebAuthnNotSupported = errors.New("webauthn not supported")
	ErrCredentialNotFound   = errors.New("credential not found")
)

// CredentialStore persists WebAuthn credentials for users.
type CredentialStore interface {
	GetCredentials(ctx context.Context, userID string) ([]StoredCredential, error)
	SaveCredential(ctx context.Context, userID string, cred StoredCredential) error
	DeleteCredential(ctx context.Context, userID, credentialID string) error
}

// StoredCredential represents a persisted WebAuthn credential.
type StoredCredential struct {
	ID              string
	PublicKey       []byte
	AttestationType string
	Transport       []string
	Flags           CredentialFlags
	AAGUID          []byte
	SignCount       uint32
	CloneWarning    bool
}

// CredentialFlags holds credential flags.
type CredentialFlags struct {
	UserPresent    bool
	UserVerified   bool
	BackupEligible bool
	BackupState   bool
}

// RegistrationOptions are options for beginning registration.
type RegistrationOptions struct {
	UserID         string
	UserName       string
	UserDisplayName string
	ResidentKey    bool
}

// RegistrationSession holds state for registration flow.
type RegistrationSession struct {
	UserID    string
	Challenge []byte
	Options   json.RawMessage
}

// AssertionOptions are options for beginning authentication.
type AssertionOptions struct {
	UserID        string // empty for usernameless
	AllowedCreds  []string
}

// AssertionSession holds state for assertion flow.
type AssertionSession struct {
	Challenge []byte
	Options   json.RawMessage
}

// Service provides WebAuthn registration and authentication.
// Implementations use github.com/go-webauthn/webauthn under the hood.
type Service interface {
	// BeginRegistration starts passkey registration.
	BeginRegistration(ctx context.Context, opts RegistrationOptions) (*RegistrationSession, error)
	// FinishRegistration completes registration with the attestation response.
	FinishRegistration(ctx context.Context, session *RegistrationSession, response json.RawMessage) error
	// BeginAssertion starts passkey authentication.
	BeginAssertion(ctx context.Context, opts AssertionOptions) (*AssertionSession, error)
	// FinishAssertion completes authentication; returns userID on success.
	FinishAssertion(ctx context.Context, session *AssertionSession, response json.RawMessage) (userID string, err error)
}
