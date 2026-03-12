package nop

import (
	"context"

	"github.com/parevo/core/notification"
)

// EmailProvider discards all emails.
type EmailProvider struct{}

// SendEmail does nothing.
func (EmailProvider) SendEmail(_ context.Context, _ notification.EmailPayload) error {
	return nil
}

// SMSProvider discards all SMS.
type SMSProvider struct{}

// SendSMS does nothing.
func (SMSProvider) SendSMS(_ context.Context, _ notification.SMSPayload) error {
	return nil
}

// WebSocketProvider discards all WebSocket notifications.
type WebSocketProvider struct{}

// SendWebSocket does nothing.
func (WebSocketProvider) SendWebSocket(_ context.Context, _ notification.WebSocketPayload) error {
	return nil
}

var (
	_ notification.EmailProvider     = EmailProvider{}
	_ notification.SMSProvider       = SMSProvider{}
	_ notification.WebSocketProvider = WebSocketProvider{}
)
