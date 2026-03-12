package notification

import "context"

// Sender sends notifications over email, SMS, and WebSocket.
type Sender interface {
	SendEmail(ctx context.Context, p EmailPayload) error
	SendSMS(ctx context.Context, p SMSPayload) error
	SendWebSocket(ctx context.Context, p WebSocketPayload) error
}

// CompositeSender delegates to separate providers.
func NewSender(email EmailProvider, sms SMSProvider, ws WebSocketProvider) Sender {
	return &compositeSender{
		email: email,
		sms:   sms,
		ws:    ws,
	}
}

type compositeSender struct {
	email EmailProvider
	sms   SMSProvider
	ws    WebSocketProvider
}

func (s *compositeSender) SendEmail(ctx context.Context, p EmailPayload) error {
	if s.email == nil {
		return nil
	}
	return s.email.SendEmail(ctx, p)
}

func (s *compositeSender) SendSMS(ctx context.Context, p SMSPayload) error {
	if s.sms == nil {
		return nil
	}
	return s.sms.SendSMS(ctx, p)
}

func (s *compositeSender) SendWebSocket(ctx context.Context, p WebSocketPayload) error {
	if s.ws == nil {
		return nil
	}
	return s.ws.SendWebSocket(ctx, p)
}
