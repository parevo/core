package mfa

import (
	"context"
	"testing"

	"github.com/parevo/core/storage/memory"
)

func TestRecoveryService_GenerateAndVerify(t *testing.T) {
	store := memory.NewMFARecoveryStore()
	svc := NewRecoveryService(store, 5)
	ctx := context.Background()

	codes, err := svc.Generate(ctx, "u1")
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}
	if len(codes) != 5 {
		t.Errorf("expected 5 codes, got %d", len(codes))
	}

	// Verify first code
	err = svc.Verify(ctx, "u1", codes[0])
	if err != nil {
		t.Fatalf("Verify failed: %v", err)
	}

	// Reuse same code should fail
	err = svc.Verify(ctx, "u1", codes[0])
	if err != ErrRecoveryCodeInvalid {
		t.Errorf("expected ErrRecoveryCodeInvalid on reuse, got %v", err)
	}

	// Verify another code
	err = svc.Verify(ctx, "u1", codes[1])
	if err != nil {
		t.Fatalf("Verify second code failed: %v", err)
	}
}

func TestRecoveryService_InvalidCode(t *testing.T) {
	store := memory.NewMFARecoveryStore()
	svc := NewRecoveryService(store, 3)
	ctx := context.Background()

	_, _ = svc.Generate(ctx, "u1")

	err := svc.Verify(ctx, "u1", "invalid")
	if err != ErrRecoveryCodeInvalid {
		t.Errorf("expected ErrRecoveryCodeInvalid, got %v", err)
	}
}

func TestRecoveryService_NormalizeCode(t *testing.T) {
	store := memory.NewMFARecoveryStore()
	svc := NewRecoveryService(store, 1)
	ctx := context.Background()

	codes, _ := svc.Generate(ctx, "u1")
	code := codes[0]

	// With dashes and spaces
	err := svc.Verify(ctx, "u1", "  "+code+"  ")
	if err != nil {
		t.Fatalf("normalized code should verify: %v", err)
	}
}
