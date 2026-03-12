package consent

import (
	"context"
	"sync"
	"testing"
)

type mockConsentStore struct {
	mu  sync.RWMutex
	recs map[string]map[string]*ConsentRecord
}

func (m *mockConsentStore) Get(ctx context.Context, userID, clientID string) (*ConsentRecord, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if r, ok := m.recs[userID][clientID]; ok {
		cp := *r
		return &cp, nil
	}
	return nil, nil
}

func (m *mockConsentStore) Save(ctx context.Context, record ConsentRecord) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.recs == nil {
		m.recs = make(map[string]map[string]*ConsentRecord)
	}
	if m.recs[record.UserID] == nil {
		m.recs[record.UserID] = make(map[string]*ConsentRecord)
	}
	r := record
	m.recs[record.UserID][record.ClientID] = &r
	return nil
}

func (m *mockConsentStore) Revoke(ctx context.Context, userID, clientID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if recs, ok := m.recs[userID]; ok {
		delete(recs, clientID)
	}
	return nil
}

func TestConsent_CheckAndGrant(t *testing.T) {
	store := &mockConsentStore{}
	svc := NewService(store, []Scope{
		{ID: "read", Name: "Read", Description: "Read access"},
		{ID: "write", Name: "Write", Description: "Write access"},
	})
	ctx := context.Background()

	err := svc.Check(ctx, "u1", "client1", []string{"read"})
	if err != ErrConsentRequired {
		t.Errorf("no consent should require: %v", err)
	}

	if err := svc.Grant(ctx, "u1", "client1", []string{"read", "write"}); err != nil {
		t.Fatalf("Grant failed: %v", err)
	}

	if err := svc.Check(ctx, "u1", "client1", []string{"read"}); err != nil {
		t.Errorf("consent granted should pass: %v", err)
	}
	if err := svc.Check(ctx, "u1", "client1", []string{"read", "write"}); err != nil {
		t.Errorf("both scopes consented should pass: %v", err)
	}
}

func TestConsent_Revoke(t *testing.T) {
	store := &mockConsentStore{}
	svc := NewService(store, nil)
	ctx := context.Background()

	svc.Grant(ctx, "u1", "client1", []string{"read"})
	if err := svc.Revoke(ctx, "u1", "client1"); err != nil {
		t.Fatalf("Revoke failed: %v", err)
	}

	if err := svc.Check(ctx, "u1", "client1", []string{"read"}); err != ErrConsentRequired {
		t.Errorf("revoked consent should require again: %v", err)
	}
}
