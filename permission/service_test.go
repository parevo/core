package permission

import (
	"context"
	"testing"

	"github.com/parevo/core/storage"
	"github.com/parevo/core/storage/memory"
)

func TestCheckWithStore(t *testing.T) {
	svc := NewService(&memory.PermissionStore{
		Grants: map[string]bool{
			"u1|t1|orders:read": true,
		},
	})

	err := svc.Check(context.Background(), storage.Subject{ID: "u1"}, "t1", "orders:read")
	if err != nil {
		t.Fatalf("permission check failed: %v", err)
	}
}
