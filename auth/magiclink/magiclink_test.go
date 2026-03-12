package magiclink

import (
	"context"
	"testing"
	"time"

	"github.com/parevo/core/storage/memory"
)

func TestMagicLink_SendAndVerify(t *testing.T) {
	store := memory.NewMagicLinkStore()
	svc := NewService(store, nil, Config{
		TTL:     15 * time.Minute,
		TokenLen: 16,
		BaseURL:  "https://app.example.com",
		Path:     "/verify",
	})
	ctx := context.Background()

	token, err := svc.SendLink(ctx, "user@example.com")
	if err != nil {
		t.Fatalf("SendLink failed: %v", err)
	}
	if token == "" {
		t.Error("expected non-empty token")
	}

	email, err := svc.Verify(ctx, token)
	if err != nil {
		t.Fatalf("Verify failed: %v", err)
	}
	if email != "user@example.com" {
		t.Errorf("expected user@example.com, got %s", email)
	}

	// Reuse should fail
	_, err = svc.Verify(ctx, token)
	if err != ErrTokenInvalid {
		t.Errorf("expected ErrTokenInvalid on reuse, got %v", err)
	}
}

func TestMagicLink_InvalidToken(t *testing.T) {
	store := memory.NewMagicLinkStore()
	svc := NewService(store, nil, DefaultConfig("https://app.com"))
	ctx := context.Background()

	_, err := svc.Verify(ctx, "nonexistent")
	if err != ErrTokenInvalid {
		t.Errorf("expected ErrTokenInvalid, got %v", err)
	}
}
