package nop

import (
	"context"
	"testing"

	"github.com/parevo/core/notification"
)

func TestSender(t *testing.T) {
	s := &Sender{}
	ctx := context.Background()

	if err := s.SendEmail(ctx, notification.EmailPayload{}); err != nil {
		t.Fatal(err)
	}
	if err := s.SendSMS(ctx, notification.SMSPayload{}); err != nil {
		t.Fatal(err)
	}
	if err := s.SendWebSocket(ctx, notification.WebSocketPayload{}); err != nil {
		t.Fatal(err)
	}
}
