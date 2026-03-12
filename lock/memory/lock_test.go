package memory

import (
	"context"
	"testing"
	"time"
)

func TestLocker(t *testing.T) {
	ctx := context.Background()
	l := NewLocker()

	ok, err := l.Lock(ctx, "key1", 30*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Error("expected lock acquired")
	}

	// Second lock should fail
	ok2, _ := l.Lock(ctx, "key1", 30*time.Second)
	if ok2 {
		t.Error("expected lock not acquired")
	}

	l.Unlock(ctx, "key1")

	// After unlock, should acquire again
	ok3, _ := l.Lock(ctx, "key1", 30*time.Second)
	if !ok3 {
		t.Error("expected lock acquired after unlock")
	}
	l.Unlock(ctx, "key1")
}
