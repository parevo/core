package memory

import (
	"context"
	"testing"
	"time"
)

func TestUsageStore(t *testing.T) {
	ctx := context.Background()
	store := NewUsageStore()

	_ = store.Record(ctx, "t1", "api_calls", 100)
	_ = store.Record(ctx, "t1", "api_calls", 50)
	_ = store.Record(ctx, "t2", "api_calls", 10)

	from := time.Now().Add(-1 * time.Hour)
	to := time.Now().Add(1 * time.Hour)

	used, err := store.Usage(ctx, "t1", "api_calls", from, to)
	if err != nil {
		t.Fatal(err)
	}
	if used != 150 {
		t.Errorf("got %d, want 150", used)
	}

	used2, _ := store.Usage(ctx, "t2", "api_calls", from, to)
	if used2 != 10 {
		t.Errorf("got %d, want 10", used2)
	}
}
