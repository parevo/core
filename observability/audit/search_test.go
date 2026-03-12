package audit

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"
	"time"
)

func TestInMemoryEventStore_Search(t *testing.T) {
	store := NewInMemoryEventStore()
	store.Append(Event{ID: "1", Timestamp: time.Now(), Event: "auth_success", UserID: "u1", TenantID: "t1"})
	store.Append(Event{ID: "2", Timestamp: time.Now(), Event: "auth_failed", UserID: "u2", TenantID: "t1"})
	store.Append(Event{ID: "3", Timestamp: time.Now(), Event: "auth_success", UserID: "u1", TenantID: "t2"})

	ctx := context.Background()

	events, err := store.Search(ctx, Query{UserID: "u1"})
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}
	if len(events) != 2 {
		t.Errorf("expected 2 events for u1, got %d", len(events))
	}

	events, err = store.Search(ctx, Query{Event: "auth_failed"})
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}
	if len(events) != 1 || events[0].UserID != "u2" {
		t.Errorf("expected 1 auth_failed event, got %d", len(events))
	}
}

func TestInMemoryEventStore_ExportJSON(t *testing.T) {
	store := NewInMemoryEventStore()
	store.Append(Event{ID: "1", Event: "auth_success", UserID: "u1"})

	var buf bytes.Buffer
	err := store.ExportJSON(context.Background(), Query{}, &buf)
	if err != nil {
		t.Fatalf("ExportJSON failed: %v", err)
	}

	var events []Event
	if err := json.Unmarshal(buf.Bytes(), &events); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(events) != 1 || events[0].UserID != "u1" {
		t.Errorf("expected 1 event, got %v", events)
	}
}

func TestInMemoryEventStore_ExportCSV(t *testing.T) {
	store := NewInMemoryEventStore()
	store.Append(Event{ID: "1", Event: "auth_success", UserID: "u1", TenantID: "t1"})

	var buf bytes.Buffer
	err := store.ExportCSV(context.Background(), Query{}, &buf)
	if err != nil {
		t.Fatalf("ExportCSV failed: %v", err)
	}
	if buf.Len() == 0 {
		t.Error("expected non-empty CSV")
	}
}
