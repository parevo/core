package memory

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"testing"
)

func TestMFARecoveryStore(t *testing.T) {
	store := NewMFARecoveryStore()
	ctx := context.Background()

	hashes := []string{"h1", "h2", "h3"}
	_ = store.SetHashes(ctx, "u1", hashes)

	got, err := store.GetHashes(ctx, "u1")
	if err != nil {
		t.Fatalf("GetHashes failed: %v", err)
	}
	if len(got) != 3 {
		t.Errorf("expected 3 hashes, got %d", len(got))
	}

	ok, err := store.Consume(ctx, "u1", "h1")
	if err != nil {
		t.Fatalf("Consume failed: %v", err)
	}
	if !ok {
		t.Error("first consume should succeed")
	}

	ok, err = store.Consume(ctx, "u1", "h1")
	if err != nil {
		t.Fatalf("Consume failed: %v", err)
	}
	if ok {
		t.Error("reuse should fail")
	}
}

func TestMFARecoveryStore_RealHash(t *testing.T) {
	store := NewMFARecoveryStore()
	ctx := context.Background()

	code := "abc123"
	h := sha256.Sum256([]byte(code))
	hash := hex.EncodeToString(h[:])
	_ = store.SetHashes(ctx, "u1", []string{hash})

	ok, err := store.Consume(ctx, "u1", hash)
	if err != nil {
		t.Fatalf("Consume failed: %v", err)
	}
	if !ok {
		t.Error("consume with real hash should succeed")
	}
}
