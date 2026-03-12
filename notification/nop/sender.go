package nop

import (
	"context"

	"github.com/parevo/core/notification"
)

// Sender is a no-op sender that discards all notifications.
// Use when notifications are disabled or not yet configured.
type Sender struct{}

// SendEmail does nothing.
func (s *Sender) SendEmail(_ context.Context, _ notification.EmailPayload) error {
	return nil
}

// SendSMS does nothing.
func (s *Sender) SendSMS(_ context.Context, _ notification.SMSPayload) error {
	return nil
}

// SendWebSocket does nothing.
func (s *Sender) SendWebSocket(_ context.Context, _ notification.WebSocketPayload) error {
	return nil
}

var _ notification.Sender = (*Sender)(nil)
