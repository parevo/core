package webhooks

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

func TestDispatcher_Dispatch(t *testing.T) {
	var received []Payload
	var mu sync.Mutex
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var p Payload
		if err := json.Unmarshal(body, &p); err == nil {
			mu.Lock()
			received = append(received, p)
			mu.Unlock()
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	d := NewDispatcher()
	d.Subscribe("ep1", Endpoint{URL: server.URL, Events: []string{EventUserCreated}})

	ctx := context.Background()
	err := d.Dispatch(ctx, EventUserCreated, map[string]any{"user_id": "u1"})
	if err != nil {
		t.Errorf("Dispatch failed: %v", err)
	}

	mu.Lock()
	n := len(received)
	mu.Unlock()
	if n != 1 {
		t.Errorf("expected 1 webhook call, got %d", n)
	}
}
