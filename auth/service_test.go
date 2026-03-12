package auth

import (
	"context"
	"errors"
	"testing"

	"github.com/parevo/core/storage/memory"
)

func newTestService(t *testing.T) *Service {
	t.Helper()
	svc, err := NewService(Config{
		Issuer:    "parevo",
		Audience:  "parevo-api",
		SecretKey: []byte("super-secret-key"),
	})
	if err != nil {
		t.Fatalf("service init failed: %v", err)
	}
	return svc
}

func TestIssueAndParseToken(t *testing.T) {
	svc := newTestService(t)
	token, err := svc.IssueAccessToken(Claims{
		UserID:      "u1",
		TenantID:    "t1",
		Roles:       []string{"admin"},
		Permissions: []string{"orders:read"},
	})
	if err != nil {
		t.Fatalf("issue token failed: %v", err)
	}

	claims, err := svc.ParseAndValidate(token)
	if err != nil {
		t.Fatalf("parse token failed: %v", err)
	}
	if claims.UserID != "u1" || claims.TenantID != "t1" {
		t.Fatalf("unexpected claims: %+v", claims)
	}
}

func TestNewService_RequiresIssuerAudience(t *testing.T) {
	_, err := NewService(Config{
		SecretKey: []byte("super-secret-key"),
	})
	if err == nil {
		t.Fatalf("expected config validation error")
	}
}

func TestAuthenticateContextTenantOverrideDenied(t *testing.T) {
	svc := newTestService(t)
	token, _ := svc.IssueAccessToken(Claims{UserID: "u1", TenantID: "t1"})

	_, _, err := svc.AuthenticateContext(context.Background(), "Bearer "+token, "t2", StaticTenantOverridePolicy{Allow: false})
	if !errors.Is(err, ErrTenantMismatch) {
		t.Fatalf("expected tenant mismatch, got: %v", err)
	}
}

func TestAuthenticateContextTenantOverrideAllowed(t *testing.T) {
	svc := newTestService(t)
	token, _ := svc.IssueAccessToken(Claims{UserID: "u1", TenantID: "t1", Roles: []string{"superadmin"}})

	ctx, claims, err := svc.AuthenticateContext(context.Background(), "Bearer "+token, "t2", RoleBasedTenantOverridePolicy{
		AllowedRoles: map[string]struct{}{"superadmin": {}},
	})
	if err != nil {
		t.Fatalf("authenticate failed: %v", err)
	}
	if claims.TenantID != "t2" {
		t.Fatalf("expected overridden tenant t2, got %s", claims.TenantID)
	}

	tenantID, ok := TenantIDFromContext(ctx)
	if !ok || tenantID != "t2" {
		t.Fatalf("context tenant mismatch: %q", tenantID)
	}
}

func TestKeyRotationWithKID(t *testing.T) {
	svc, err := NewService(Config{
		Issuer:   "parevo",
		Audience: "parevo-api",
		SigningKeys: map[string][]byte{
			"k1": []byte("first-key"),
			"k2": []byte("second-key"),
		},
		ActiveKID: "k2",
	})
	if err != nil {
		t.Fatalf("service init failed: %v", err)
	}

	token, err := svc.IssueAccessToken(Claims{UserID: "u1", TenantID: "t1"})
	if err != nil {
		t.Fatalf("issue token failed: %v", err)
	}
	if _, err := svc.ParseAndValidate(token); err != nil {
		t.Fatalf("parse token with kid failed: %v", err)
	}
}

func TestRevokedSessionDenied(t *testing.T) {
	sessionStore := &memory.SessionStore{}
	svc, err := NewServiceWithModules(Config{
		Issuer:    "parevo",
		Audience:  "parevo-api",
		SecretKey: []byte("super-secret-key"),
	}, Modules{
		SessionStore: sessionStore,
	})
	if err != nil {
		t.Fatalf("service init failed: %v", err)
	}

	sessionID := "s1"
	token, err := svc.IssueAccessToken(Claims{UserID: "u1", TenantID: "t1", SessionID: sessionID})
	if err != nil {
		t.Fatalf("issue token failed: %v", err)
	}

	if err := svc.RevokeSession(context.Background(), sessionID); err != nil {
		t.Fatalf("revoke session failed: %v", err)
	}
	if _, err := svc.ParseAndValidate(token); !errors.Is(err, ErrSessionRevoked) {
		t.Fatalf("expected revoked session error, got: %v", err)
	}
}

func TestRotateRefreshTokenAndDetectReuse(t *testing.T) {
	sessionStore := &memory.SessionStore{}
	refreshStore := &memory.RefreshStore{}
	svc, err := NewServiceWithModules(Config{
		Issuer:    "parevo",
		Audience:  "parevo-api",
		SecretKey: []byte("super-secret-key"),
	}, Modules{
		SessionStore: sessionStore,
		RefreshStore: refreshStore,
	})
	if err != nil {
		t.Fatalf("service init failed: %v", err)
	}

	pair, err := svc.IssueTokenPair(context.Background(), Claims{
		UserID:    "u1",
		TenantID:  "t1",
		SessionID: "s1",
	})
	if err != nil {
		t.Fatalf("issue token pair failed: %v", err)
	}
	rotated, err := svc.RotateRefreshToken(context.Background(), pair.RefreshToken)
	if err != nil {
		t.Fatalf("rotate refresh failed: %v", err)
	}
	if rotated.AccessToken == "" || rotated.RefreshToken == "" {
		t.Fatalf("expected rotated pair")
	}

	_, err = svc.RotateRefreshToken(context.Background(), pair.RefreshToken)
	if !errors.Is(err, ErrRefreshReuse) {
		t.Fatalf("expected refresh reuse error, got: %v", err)
	}
}

func TestRevokeAllSessionsByUser(t *testing.T) {
	sessionStore := &memory.SessionStore{}
	refreshStore := &memory.RefreshStore{}
	svc, err := NewServiceWithModules(Config{
		Issuer:    "parevo",
		Audience:  "parevo-api",
		SecretKey: []byte("super-secret-key"),
	}, Modules{
		SessionStore: sessionStore,
		RefreshStore: refreshStore,
	})
	if err != nil {
		t.Fatalf("service init failed: %v", err)
	}

	_, err = svc.IssueTokenPair(context.Background(), Claims{UserID: "u1", TenantID: "t1", SessionID: "s1"})
	if err != nil {
		t.Fatalf("issue pair s1 failed: %v", err)
	}
	_, err = svc.IssueTokenPair(context.Background(), Claims{UserID: "u1", TenantID: "t1", SessionID: "s2"})
	if err != nil {
		t.Fatalf("issue pair s2 failed: %v", err)
	}

	if err := svc.RevokeAllSessionsByUser(context.Background(), "u1"); err != nil {
		t.Fatalf("revoke all failed: %v", err)
	}

	revoked1, _ := sessionStore.IsSessionRevoked(context.Background(), "s1")
	revoked2, _ := sessionStore.IsSessionRevoked(context.Background(), "s2")
	if !revoked1 || !revoked2 {
		t.Fatalf("expected both sessions revoked")
	}
}
