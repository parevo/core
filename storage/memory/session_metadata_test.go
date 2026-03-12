package memory

import (
	"context"
	"testing"
)

func TestSessionMetadataStore(t *testing.T) {
	store := NewSessionMetadataStore()
	ctx := context.Background()

	_ = store.SetMetadata(ctx, "s1", "u1", "192.168.1.1", "Mozilla/5.0")
	_ = store.BindSessionToUser(ctx, "u1", "s1")

	meta, err := store.ListWithMetadata(ctx, "u1")
	if err != nil {
		t.Fatalf("ListWithMetadata failed: %v", err)
	}
	if len(meta) != 1 || meta[0].IP != "192.168.1.1" {
		t.Errorf("unexpected metadata: %v", meta)
	}
}
