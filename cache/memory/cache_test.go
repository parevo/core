package memory

import (
	"context"
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	ctx := context.Background()
	c := NewCache(5 * time.Minute)

	// Set and Get
	err := c.Set(ctx, "k1", []byte("v1"), 10*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	v, err := c.Get(ctx, "k1")
	if err != nil {
		t.Fatal(err)
	}
	if string(v) != "v1" {
		t.Errorf("got %q", v)
	}

	// Delete
	err = c.Delete(ctx, "k1")
	if err != nil {
		t.Fatal(err)
	}
	v, _ = c.Get(ctx, "k1")
	if v != nil {
		t.Errorf("expected nil after delete, got %q", v)
	}

	// Expiry
	_ = c.Set(ctx, "k2", []byte("v2"), 10*time.Millisecond)
	time.Sleep(20 * time.Millisecond)
	v, _ = c.Get(ctx, "k2")
	if v != nil {
		t.Errorf("expected nil after expiry, got %q", v)
	}
}
