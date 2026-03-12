package tenant

import (
	"context"
	"testing"

	"github.com/parevo/core/storage"
	"github.com/parevo/core/storage/memory"
)

func TestLifecycleService(t *testing.T) {
	store := &memory.TenantLifecycleStore{}
	svc := NewLifecycleService(store)

	if err := svc.Create(context.Background(), "t1", "Tenant 1", "owner-1"); err != nil {
		t.Fatalf("create failed: %v", err)
	}
	status, err := svc.Status(context.Background(), "t1")
	if err != nil {
		t.Fatalf("status failed: %v", err)
	}
	if status != storage.TenantStatusActive {
		t.Fatalf("expected active, got %s", status)
	}
	if err := svc.Suspend(context.Background(), "t1"); err != nil {
		t.Fatalf("suspend failed: %v", err)
	}
	status, _ = svc.Status(context.Background(), "t1")
	if status != storage.TenantStatusSuspended {
		t.Fatalf("expected suspended, got %s", status)
	}
}
