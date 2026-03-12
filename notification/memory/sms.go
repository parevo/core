package memory

import (
	"context"
	"fmt"
	"sync"

	"github.com/parevo/core/notification"
)

// SMSProvider implements notification.SMSProvider with in-memory storage.
type SMSProvider struct {
	mu         sync.RWMutex
	Sent       []notification.SMSPayload
	SharedSent *[]SentItem // optional: append to shared slice
	LogFn      func(format string, args ...interface{})
}

// NewSMSProvider creates an in-memory SMS provider.
func NewSMSProvider() *SMSProvider {
	return &SMSProvider{
		Sent: make([]notification.SMSPayload, 0),
		LogFn: func(format string, args ...interface{}) {
			fmt.Printf("[notification/sms] "+format+"\n", args...)
		},
	}
}

// SendSMS records the SMS and optionally logs.
func (p *SMSProvider) SendSMS(ctx context.Context, payload notification.SMSPayload) error {
	p.mu.Lock()
	p.Sent = append(p.Sent, payload)
	if p.SharedSent != nil {
		*p.SharedSent = append(*p.SharedSent, SentItem{Channel: "sms", Payload: payload})
	}
	p.mu.Unlock()
	if p.LogFn != nil {
		p.LogFn("to=%s body=%s", payload.To, payload.Body)
	}
	return nil
}

// Reset clears sent items (for tests).
func (p *SMSProvider) Reset() {
	p.mu.Lock()
	p.Sent = p.Sent[:0]
	p.mu.Unlock()
}

var _ notification.SMSProvider = (*SMSProvider)(nil)
