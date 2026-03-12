package memory

import (
	"context"
	"fmt"
	"sync"

	"github.com/parevo/core/notification"
)

// SentItem holds a sent notification for inspection.
type SentItem struct {
	Channel string
	Payload interface{}
}

// EmailProvider implements notification.EmailProvider with in-memory storage.
type EmailProvider struct {
	mu         sync.RWMutex
	Sent       []notification.EmailPayload
	SharedSent *[]SentItem // optional: append to shared slice (e.g. for MemorySender)
	LogFn      func(format string, args ...interface{})
}

// NewEmailProvider creates an in-memory email provider.
func NewEmailProvider() *EmailProvider {
	return &EmailProvider{
		Sent: make([]notification.EmailPayload, 0),
		LogFn: func(format string, args ...interface{}) {
			fmt.Printf("[notification/email] "+format+"\n", args...)
		},
	}
}

// SendEmail records the email and optionally logs.
func (p *EmailProvider) SendEmail(ctx context.Context, payload notification.EmailPayload) error {
	p.mu.Lock()
	p.Sent = append(p.Sent, payload)
	if p.SharedSent != nil {
		*p.SharedSent = append(*p.SharedSent, SentItem{Channel: "email", Payload: payload})
	}
	p.mu.Unlock()
	if p.LogFn != nil {
		p.LogFn("to=%s subject=%s", payload.To, payload.Subject)
	}
	return nil
}

// Reset clears sent items (for tests).
func (p *EmailProvider) Reset() {
	p.mu.Lock()
	p.Sent = p.Sent[:0]
	p.mu.Unlock()
}

var _ notification.EmailProvider = (*EmailProvider)(nil)
