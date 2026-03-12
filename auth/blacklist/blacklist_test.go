package blacklist

import (
	"context"
	"testing"
	"time"

	"github.com/parevo/core/storage/memory"
)

func TestBlacklist_RevokeAndCheck(t *testing.T) {
	store := memory.NewBlacklistStore()
	svc := NewService(store)
	ctx := context.Background()

	expiresAt := time.Now().Add(time.Hour)
	if err := svc.Revoke(ctx, "jti-123", expiresAt); err != nil {
		t.Fatalf("Revoke failed: %v", err)
	}

	ok, err := svc.IsBlacklisted(ctx, "jti-123")
	if err != nil {
		t.Fatalf("IsBlacklisted failed: %v", err)
	}
	if !ok {
		t.Error("token should be blacklisted")
	}

	if err := svc.Check(ctx, "jti-123"); err != ErrTokenBlacklisted {
		t.Errorf("Check should return ErrTokenBlacklisted: %v", err)
	}
}

func TestBlacklist_NotBlacklisted(t *testing.T) {
	svc := NewService(memory.NewBlacklistStore())
	ctx := context.Background()

	ok, err := svc.IsBlacklisted(ctx, "unknown")
	if err != nil {
		t.Fatalf("IsBlacklisted failed: %v", err)
	}
	if ok {
		t.Error("unknown token should not be blacklisted")
	}

	if err := svc.Check(ctx, "unknown"); err != nil {
		t.Errorf("Check should pass for non-blacklisted: %v", err)
	}
}
