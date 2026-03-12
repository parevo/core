package memory

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/parevo/core/notification"
)

// WebSocketProvider implements notification.WebSocketProvider with in-memory storage.
type WebSocketProvider struct {
	mu         sync.RWMutex
	Sent       []notification.WebSocketPayload
	SharedSent *[]SentItem // optional: append to shared slice
	LogFn      func(format string, args ...interface{})
}

// NewWebSocketProvider creates an in-memory WebSocket provider.
func NewWebSocketProvider() *WebSocketProvider {
	return &WebSocketProvider{
		Sent: make([]notification.WebSocketPayload, 0),
		LogFn: func(format string, args ...interface{}) {
			fmt.Printf("[notification/websocket] "+format+"\n", args...)
		},
	}
}

// SendWebSocket records the payload and optionally logs.
func (p *WebSocketProvider) SendWebSocket(ctx context.Context, payload notification.WebSocketPayload) error {
	p.mu.Lock()
	p.Sent = append(p.Sent, payload)
	if p.SharedSent != nil {
		*p.SharedSent = append(*p.SharedSent, SentItem{Channel: "websocket", Payload: payload})
	}
	p.mu.Unlock()
	if p.LogFn != nil {
		b, _ := json.Marshal(payload.Data)
		p.LogFn("userIDs=%v event=%s data=%s", payload.UserIDs, payload.Event, string(b))
	}
	return nil
}

// Reset clears sent items (for tests).
func (p *WebSocketProvider) Reset() {
	p.mu.Lock()
	p.Sent = p.Sent[:0]
	p.mu.Unlock()
}

var _ notification.WebSocketProvider = (*WebSocketProvider)(nil)
