package permission

import (
	"context"
	"testing"

	"github.com/parevo/core/storage"
	"github.com/parevo/core/storage/memory"
)

func TestResourceOwnerCondition(t *testing.T) {
	cond := ResourceOwnerCondition{}
	ctx := context.Background()

	// Owner matches
	ok, err := cond.Allow(ctx, SubjectAttributes{
		Subject: storage.Subject{ID: "u1"},
	}, Resource{OwnerID: "u1"}, "edit")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Error("owner should be allowed")
	}

	// Owner does not match
	ok, err = cond.Allow(ctx, SubjectAttributes{
		Subject: storage.Subject{ID: "u2"},
	}, Resource{OwnerID: "u1"}, "edit")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Error("non-owner should be denied")
	}

	// Empty owner
	ok, err = cond.Allow(ctx, SubjectAttributes{Subject: storage.Subject{ID: "u1"}}, Resource{}, "edit")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Error("empty owner should deny")
	}
}

func TestDepartmentCondition(t *testing.T) {
	cond := DepartmentCondition{RequireMatch: true}
	ctx := context.Background()

	ok, err := cond.Allow(ctx, SubjectAttributes{Department: "eng"}, Resource{Department: "eng"}, "read")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Error("same department should allow")
	}

	ok, err = cond.Allow(ctx, SubjectAttributes{Department: "sales"}, Resource{Department: "eng"}, "read")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Error("different department should deny")
	}
}

func TestEnvironmentCondition(t *testing.T) {
	cond := EnvironmentCondition{AllowedEnvs: []string{"prod", "staging"}}
	ctx := context.Background()

	ok, err := cond.Allow(ctx, SubjectAttributes{}, Resource{Environment: "prod"}, "read")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Error("prod should be allowed")
	}

	ok, err = cond.Allow(ctx, SubjectAttributes{}, Resource{Environment: "dev"}, "read")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Error("dev should be denied")
	}
}

func TestABACPolicy_ModeAll(t *testing.T) {
	policy := &ABACPolicy{
		Conditions: []ABACCondition{
			ResourceOwnerCondition{},
			DepartmentCondition{RequireMatch: true},
		},
		Mode: ABACModeAll,
	}
	ctx := context.Background()

	subject := SubjectAttributes{
		Subject:    storage.Subject{ID: "u1"},
		Department: "eng",
	}
	resource := Resource{OwnerID: "u1", Department: "eng"}

	ok, err := policy.Allow(ctx, subject, resource, "edit")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Error("both conditions pass, should allow")
	}

	resource.Department = "sales"
	ok, err = policy.Allow(ctx, subject, resource, "edit")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Error("department mismatch should deny")
	}
}

func TestABACPolicy_EmptyConditions(t *testing.T) {
	policy := &ABACPolicy{Conditions: nil}
	ctx := context.Background()
	ok, err := policy.Allow(ctx, SubjectAttributes{}, Resource{}, "read")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Error("empty conditions should allow")
	}
}

func TestABACService_Check(t *testing.T) {
	permStore := &memory.PermissionStore{Grants: map[string]bool{"u1|t1|orders:read": true}}
	permSvc := NewService(permStore)
	abacSvc := NewABACService(permSvc)
	ctx := context.Background()

	subject := SubjectAttributes{Subject: storage.Subject{ID: "u1"}}
	resource := Resource{OwnerID: "u1"}

	err := abacSvc.Check(ctx, subject, "t1", "orders:read", resource, &ABACPolicy{
		Conditions: []ABACCondition{ResourceOwnerCondition{}},
		Mode:      ABACModeAll,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	resource.OwnerID = "u2"
	err = abacSvc.Check(ctx, subject, "t1", "orders:read", resource, &ABACPolicy{
		Conditions: []ABACCondition{ResourceOwnerCondition{}},
		Mode:      ABACModeAll,
	})
	if err != ErrABACDenied {
		t.Errorf("expected ErrABACDenied, got %v", err)
	}
}
