package memory

import (
	"context"
	"testing"
	"time"
)

func TestMagicLinkStore(t *testing.T) {
	store := NewMagicLinkStore()
	ctx := context.Background()

	store.Create(ctx, "user@x.com", "token123", time.Now().Add(time.Hour))

	email, err := store.Consume(ctx, "token123")
	if err != nil {
		t.Fatalf("Consume failed: %v", err)
	}
	if email != "user@x.com" {
		t.Errorf("expected user@x.com, got %s", email)
	}

	email, _ = store.Consume(ctx, "token123")
	if email != "" {
		t.Error("reuse should return empty")
	}
}
