package webhooks

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"sync"
	"time"
)

var (
	ErrWebhookFailed = errors.New("webhook delivery failed")
)

// Event types for webhook payloads.
const (
	EventUserCreated     = "user.created"
	EventUserDeleted     = "user.deleted"
	EventSessionRevoked  = "session.revoked"
	EventPermissionGranted = "permission.granted"
	EventPermissionRevoked = "permission.revoked"
	EventTenantCreated   = "tenant.created"
	EventTenantSuspended = "tenant.suspended"
)

// Payload is the webhook payload.
type Payload struct {
	ID        string            `json:"id"`
	Event     string            `json:"event"`
	Timestamp time.Time         `json:"timestamp"`
	Data      map[string]any    `json:"data"`
}

// Endpoint is a webhook subscription.
type Endpoint struct {
	URL     string
	Secret  string
	Events  []string
}

// Dispatcher sends webhook events to subscribed endpoints.
type Dispatcher struct {
	mu         sync.RWMutex
	endpoints  map[string][]Endpoint
	httpClient *http.Client
}

// NewDispatcher creates a webhook dispatcher.
func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		endpoints:  make(map[string][]Endpoint),
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

// Subscribe adds an endpoint for events.
func (d *Dispatcher) Subscribe(endpointID string, ep Endpoint) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.endpoints[endpointID] = append(d.endpoints[endpointID], ep)
}

// Dispatch sends an event to all subscribed endpoints.
func (d *Dispatcher) Dispatch(ctx context.Context, event string, data map[string]any) error {
	payload := Payload{
		ID:        generateID(),
		Event:     event,
		Timestamp: time.Now(),
		Data:      data,
	}
	d.mu.RLock()
	all := make([]Endpoint, 0)
	for _, eps := range d.endpoints {
		for _, ep := range eps {
			if contains(ep.Events, event) {
				all = append(all, ep)
			}
		}
	}
	d.mu.RUnlock()

	var lastErr error
	for _, ep := range all {
		if err := d.send(ctx, ep, payload); err != nil {
			lastErr = err
		}
	}
	return lastErr
}

func (d *Dispatcher) send(ctx context.Context, ep Endpoint, payload Payload) error {
	body, _ := json.Marshal(payload)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, ep.URL, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Webhook-Event", payload.Event)
	// In real impl: add HMAC signature using ep.Secret
	_ = ep.Secret
	resp, err := d.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode >= 400 {
		return ErrWebhookFailed
	}
	return nil
}

func contains(s []string, x string) bool {
	for _, v := range s {
		if v == x {
			return true
		}
	}
	return false
}

func generateID() string {
	return time.Now().Format("20060102150405")
}
