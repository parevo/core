package permission

import (
	"context"
	"testing"

	"github.com/parevo/core/storage"
	"github.com/parevo/core/storage/memory"
)

func BenchmarkCheck(b *testing.B) {
	svc := NewService(&memory.PermissionStore{
		Grants: map[string]bool{
			"u1|t1|orders:read": true,
		},
	})
	ctx := context.Background()
	subject := storage.Subject{ID: "u1"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = svc.Check(ctx, subject, "t1", "orders:read")
	}
}

func BenchmarkCheckWithWildcard(b *testing.B) {
	svc := NewService(nil)
	granted := []string{"orders:*", "users:read"}
	subject := storage.Subject{ID: "u1"}
	ctx := context.Background()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = svc.CheckWithWildcard(ctx, subject, "t1", "orders:create", granted)
	}
}

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
