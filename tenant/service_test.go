package tenant

import (
	"context"
	"testing"

	"github.com/parevo/core/storage"
	"github.com/parevo/core/storage/memory"
)

func TestResolveWithStore(t *testing.T) {
	svc := NewService(&memory.TenantStore{
		SubjectTenants: map[string][]string{
			"u1": {"t1", "t2"},
		},
	})

	tenantID, err := svc.Resolve(
		context.Background(),
		storage.Subject{ID: "u1", Roles: []string{"admin"}},
		"t1",
		"t2",
		StaticOverridePolicy{Allow: true},
	)
	if err != nil {
		t.Fatalf("resolve failed: %v", err)
	}
	if tenantID != "t2" {
		t.Fatalf("expected t2, got %s", tenantID)
	}
}
