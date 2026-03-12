package auth

import (
	"context"
	"errors"
	"testing"
)

func TestRequirePermission(t *testing.T) {
	base := context.Background()
	ctx := WithClaims(base, &Claims{
		UserID:      "u1",
		TenantID:    "t1",
		Permissions: []string{"users:write"},
	})

	if err := RequirePermission(ctx, "users:write"); err != nil {
		t.Fatalf("expected permission granted: %v", err)
	}
	if err := RequirePermission(ctx, "users:delete"); !errors.Is(err, ErrForbidden) {
		t.Fatalf("expected forbidden, got: %v", err)
	}
}

func TestRequireTenant(t *testing.T) {
	ctx := WithClaims(context.Background(), &Claims{UserID: "u1", TenantID: "t1"})
	if err := RequireTenant(ctx, "t1"); err != nil {
		t.Fatalf("expected tenant check pass: %v", err)
	}
	if err := RequireTenant(ctx, "t2"); !errors.Is(err, ErrTenantMismatch) {
		t.Fatalf("expected tenant mismatch, got: %v", err)
	}
}
