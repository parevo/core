package tenant

import (
	"context"
	"testing"

	"github.com/parevo/core/storage"
	"github.com/parevo/core/storage/memory"
)

func TestFeatureService_NilStore(t *testing.T) {
	svc := NewFeatureService(nil)
	ctx := context.Background()

	ok, err := svc.IsEnabled(ctx, "t1", "saml")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Error("nil store should default allow")
	}

	plan, err := svc.GetPlan(ctx, "t1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if plan != storage.PlanFree {
		t.Errorf("expected PlanFree, got %v", plan)
	}

	err = svc.RequireEnabled(ctx, "t1", "saml")
	if err != nil {
		t.Errorf("nil store should allow: %v", err)
	}
}

func TestFeatureService_WithStore(t *testing.T) {
	store := memory.NewTenantFeatureStore()
	_ = store.SetPlan(context.Background(), "t1", storage.PlanPro)
	_ = store.SetFeature(context.Background(), "t1", "saml", true)
	_ = store.SetLimit(context.Background(), "t1", "max_users", 100)

	svc := NewFeatureService(store)
	ctx := context.Background()

	ok, err := svc.IsEnabled(ctx, "t1", "saml")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Error("saml should be enabled for t1")
	}

	plan, err := svc.GetPlan(ctx, "t1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if plan != storage.PlanPro {
		t.Errorf("expected PlanPro, got %v", plan)
	}

	limit, err := svc.GetLimit(ctx, "t1", "max_users")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if limit != 100 {
		t.Errorf("expected limit 100, got %d", limit)
	}

	err = svc.CheckLimit(ctx, "t1", "max_users", 50)
	if err != nil {
		t.Errorf("50 < 100 should pass: %v", err)
	}

	err = svc.CheckLimit(ctx, "t1", "max_users", 100)
	if err != ErrLimitExceeded {
		t.Errorf("expected ErrLimitExceeded, got %v", err)
	}
}

func TestFeatureService_RequireEnabled_Denied(t *testing.T) {
	store := memory.NewTenantFeatureStore()
	_ = store.SetPlan(context.Background(), "t1", storage.PlanFree)

	svc := NewFeatureService(store)
	ctx := context.Background()

	err := svc.RequireEnabled(ctx, "t1", "saml")
	if err != ErrFeatureDisabled {
		t.Errorf("expected ErrFeatureDisabled, got %v", err)
	}
}
