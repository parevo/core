package webauthn

import (
	"context"
	"encoding/json"
	"sync"
	"testing"
)

type mockCredStore struct {
	mu sync.Mutex
	m  map[string][]StoredCredential
}

func (m *mockCredStore) GetCredentials(ctx context.Context, userID string) ([]StoredCredential, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.m[userID], nil
}

func (m *mockCredStore) SaveCredential(ctx context.Context, userID string, cred StoredCredential) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.m == nil {
		m.m = make(map[string][]StoredCredential)
	}
	m.m[userID] = append(m.m[userID], cred)
	return nil
}

func (m *mockCredStore) DeleteCredential(ctx context.Context, userID, credentialID string) error { return nil }

func TestStubService_ReturnsNotSupported(t *testing.T) {
	store := &mockCredStore{}
	svc, err := NewService("example.com", "https://example.com", store)
	if err != nil {
		t.Fatalf("NewService failed: %v", err)
	}
	ctx := context.Background()

	_, err = svc.BeginRegistration(ctx, RegistrationOptions{UserID: "u1", UserName: "user"})
	if err != ErrWebAuthnNotSupported {
		t.Errorf("expected ErrWebAuthnNotSupported, got %v", err)
	}

	_, err = svc.BeginAssertion(ctx, AssertionOptions{UserID: "u1"})
	if err != ErrWebAuthnNotSupported {
		t.Errorf("expected ErrWebAuthnNotSupported, got %v", err)
	}

	_, err = svc.FinishAssertion(ctx, &AssertionSession{}, json.RawMessage("{}"))
	if err != ErrWebAuthnNotSupported {
		t.Errorf("expected ErrWebAuthnNotSupported, got %v", err)
	}
}
