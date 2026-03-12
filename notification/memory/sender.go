package memory

import (
	"context"
	"fmt"

	"github.com/parevo/core/notification"
)

// Sender is a convenience Sender that uses in-memory providers for all channels.
// Use for development and testing.
type Sender struct {
	sent   *[]SentItem
	sender notification.Sender
}

// NewSender creates an in-memory sender with shared Sent for all channels.
func NewSender() *Sender {
	sent := make([]SentItem, 0)
	logFn := func(ch, format string, args ...interface{}) {
		fmt.Printf("[notification/"+ch+"] "+format+"\n", args...)
	}
	ep := NewEmailProvider()
	ep.SharedSent = &sent
	ep.LogFn = func(format string, args ...interface{}) { logFn("email", format, args...) }
	sp := NewSMSProvider()
	sp.SharedSent = &sent
	sp.LogFn = func(format string, args ...interface{}) { logFn("sms", format, args...) }
	wp := NewWebSocketProvider()
	wp.SharedSent = &sent
	wp.LogFn = func(format string, args ...interface{}) { logFn("websocket", format, args...) }
	return &Sender{
		sent:   &sent,
		sender: notification.NewSender(ep, sp, wp),
	}
}

// SendEmail delegates to the email provider.
func (s *Sender) SendEmail(ctx context.Context, p notification.EmailPayload) error {
	return s.sender.SendEmail(ctx, p)
}

// SendSMS delegates to the SMS provider.
func (s *Sender) SendSMS(ctx context.Context, p notification.SMSPayload) error {
	return s.sender.SendSMS(ctx, p)
}

// SendWebSocket delegates to the WebSocket provider.
func (s *Sender) SendWebSocket(ctx context.Context, p notification.WebSocketPayload) error {
	return s.sender.SendWebSocket(ctx, p)
}

// Sent returns all sent notifications (for tests, debugging).
func (s *Sender) Sent() []SentItem {
	return *s.sent
}

// Reset clears the sent items (for tests).
func (s *Sender) Reset() {
	*s.sent = (*s.sent)[:0]
}

var _ notification.Sender = (*Sender)(nil)
