package ipfilter

import (
	"context"
	"testing"

	"github.com/parevo/core/storage/memory"
)

func TestIPFilter_BlocklistMode(t *testing.T) {
	store := memory.NewIPFilterStore()
	store.BlockIP("", "192.168.1.100")
	store.BlockIP("t1", "10.0.0.5")

	svc := NewService(store, false)
	ctx := context.Background()

	if err := svc.Allow(ctx, "", "192.168.1.100"); err != ErrIPBlocked {
		t.Errorf("blocked IP should fail: %v", err)
	}
	if err := svc.Allow(ctx, "t1", "10.0.0.5"); err != ErrIPBlocked {
		t.Errorf("tenant blocked IP should fail: %v", err)
	}
	if err := svc.Allow(ctx, "", "192.168.1.1"); err != nil {
		t.Errorf("non-blocked IP should pass: %v", err)
	}
}

func TestIPFilter_AllowlistMode(t *testing.T) {
	store := memory.NewIPFilterStore()
	store.AllowIP("t1", "192.168.1.100")

	svc := NewService(store, true)
	ctx := context.Background()

	if err := svc.Allow(ctx, "t1", "192.168.1.100"); err != nil {
		t.Errorf("allowed IP should pass: %v", err)
	}
	if err := svc.Allow(ctx, "t1", "192.168.1.1"); err != ErrIPNotAllowed {
		t.Errorf("non-allowed IP should fail: %v", err)
	}
}

func TestIPFilter_EmptyIP(t *testing.T) {
	svc := NewService(memory.NewIPFilterStore(), false)
	ctx := context.Background()
	if err := svc.Allow(ctx, "", ""); err != nil {
		t.Errorf("empty IP should pass: %v", err)
	}
}
