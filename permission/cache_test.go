package permission

import (
	"context"
	"testing"
	"time"

	"github.com/parevo/core/storage/memory"
)

func TestCachedPermissionStore(t *testing.T) {
	ctx := context.Background()
	base := &memory.PermissionStore{Grants: map[string]bool{
		"u1|t1|orders:read": true,
	}}
	cached := NewCachedPermissionStore(base, 100*time.Millisecond)

	ok, err := cached.HasPermission(ctx, "u1", "t1", "orders:read", nil)
	if err != nil || !ok {
		t.Fatalf("expected ok=true, err=nil; got ok=%v err=%v", ok, err)
	}

	ok, err = cached.HasPermission(ctx, "u1", "t1", "orders:write", nil)
	if err != nil || ok {
		t.Fatalf("expected ok=false, err=nil; got ok=%v err=%v", ok, err)
	}

	ok, err = cached.HasPermission(ctx, "u1", "t1", "orders:write", nil)
	if err != nil || ok {
		t.Fatalf("expected ok=false (no grant yet), err=nil; got ok=%v err=%v", ok, err)
	}

	base.Grants["u1|t1|orders:write"] = true
	ok, err = cached.HasPermission(ctx, "u1", "t1", "orders:write", nil)
	if err != nil || ok {
		t.Fatalf("expected ok=false (cache hit with old result), err=nil; got ok=%v err=%v", ok, err)
	}

	cached.InvalidateSubjectTenant("u1", "t1")
	ok, err = cached.HasPermission(ctx, "u1", "t1", "orders:write", nil)
	if err != nil || !ok {
		t.Fatalf("expected ok=true after invalidation (cache miss, fresh from store), err=nil; got ok=%v err=%v", ok, err)
	}
}
