package graphql

import (
	"context"
	"testing"

	"github.com/parevo/core/auth"
	"github.com/parevo/core/auth/adapters"
	"github.com/parevo/core/tenant"
)

func TestRequireAuth(t *testing.T) {
	ctx := context.Background()

	if err := RequireAuth(ctx); err != auth.ErrUnauthenticated {
		t.Errorf("empty context should fail: %v", err)
	}

	ctx = auth.WithClaims(ctx, &auth.Claims{UserID: "u1"})
	if err := RequireAuth(ctx); err != nil {
		t.Errorf("context with claims should pass: %v", err)
	}
}

func TestRequireScope(t *testing.T) {
	ctx := context.Background()

	if err := RequireScope(ctx, "read:orders"); err != auth.ErrUnauthenticated {
		t.Errorf("empty context should fail: %v", err)
	}

	ctx = auth.WithClaims(ctx, &auth.Claims{UserID: "u1", Scopes: []string{"read:*"}})
	if err := RequireScope(ctx, "read:orders"); err != nil {
		t.Errorf("matching scope should pass: %v", err)
	}

	ctx = auth.WithClaims(ctx, &auth.Claims{UserID: "u1", Scopes: []string{"read:orders"}})
	if err := RequireScope(ctx, "write:users"); err != auth.ErrForbidden {
		t.Errorf("missing scope should fail: %v", err)
	}
}

func TestResolverContext_AuthenticateRequest(t *testing.T) {
	authSvc, _ := auth.NewService(auth.Config{
		Issuer:    "test",
		Audience:  "test",
		SecretKey: []byte("super-secret-key-at-least-32-bytes"),
	})
	rc := NewResolverContext(authSvc, &tenant.Service{}, adapters.Options{})

	token, _ := authSvc.IssueAccessToken(auth.Claims{UserID: "u1", TenantID: "t1"})
	ctx, err := rc.AuthenticateRequest(context.Background(), "Bearer "+token, "t1")
	if err != nil {
		t.Fatalf("AuthenticateRequest failed: %v", err)
	}
	if _, ok := auth.ClaimsFromContext(ctx); !ok {
		t.Error("context should have claims")
	}
}
