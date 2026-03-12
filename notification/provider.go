package notification

import "context"

// EmailProvider sends email. Implement with SMTP, SendGrid, etc.
type EmailProvider interface {
	SendEmail(ctx context.Context, p EmailPayload) error
}

// SMSProvider sends SMS. Implement with Twilio, SNS, etc.
type SMSProvider interface {
	SendSMS(ctx context.Context, p SMSPayload) error
}

// WebSocketProvider delivers real-time in-app notifications.
// Implement with Redis Pub/Sub, in-memory hub, etc.
type WebSocketProvider interface {
	SendWebSocket(ctx context.Context, p WebSocketPayload) error
}
