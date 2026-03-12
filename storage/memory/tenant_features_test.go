package memory

import (
	"context"
	"testing"

	"github.com/parevo/core/storage"
)

func TestTenantFeatureStore(t *testing.T) {
	store := NewTenantFeatureStore()
	ctx := context.Background()

	store.SetPlan(ctx, "t1", storage.PlanPro)
	store.SetFeature(ctx, "t1", "saml", true)
	store.SetLimit(ctx, "t1", "max_users", 50)

	plan, err := store.GetPlan(ctx, "t1")
	if err != nil {
		t.Fatalf("GetPlan failed: %v", err)
	}
	if plan != storage.PlanPro {
		t.Errorf("expected PlanPro, got %v", plan)
	}

	ok, err := store.IsEnabled(ctx, "t1", "saml")
	if err != nil {
		t.Fatalf("IsEnabled failed: %v", err)
	}
	if !ok {
		t.Error("saml should be enabled")
	}

	limit, err := store.GetLimit(ctx, "t1", "max_users")
	if err != nil {
		t.Fatalf("GetLimit failed: %v", err)
	}
	if limit != 50 {
		t.Errorf("expected limit 50, got %d", limit)
	}
}
