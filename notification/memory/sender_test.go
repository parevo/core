package memory

import (
	"context"
	"testing"

	"github.com/parevo/core/notification"
)

func TestSender_SendEmail(t *testing.T) {
	s := NewSender()
	ctx := context.Background()

	err := s.SendEmail(ctx, notification.EmailPayload{
		To:      "a@b.com",
		Subject: "Test",
		Body:    "Hello",
	})
	if err != nil {
		t.Fatal(err)
	}
	sent := s.Sent()
	if len(sent) != 1 {
		t.Fatalf("want 1 sent, got %d", len(sent))
	}
	if sent[0].Channel != "email" {
		t.Errorf("want channel email, got %s", sent[0].Channel)
	}
}

func TestSender_SendSMS(t *testing.T) {
	s := NewSender()
	ctx := context.Background()

	err := s.SendSMS(ctx, notification.SMSPayload{To: "+1", Body: "Hi"})
	if err != nil {
		t.Fatal(err)
	}
	sent := s.Sent()
	if len(sent) != 1 || sent[0].Channel != "sms" {
		t.Fatalf("want 1 sms sent, got %d", len(sent))
	}
}

func TestSender_SendWebSocket(t *testing.T) {
	s := NewSender()
	ctx := context.Background()

	err := s.SendWebSocket(ctx, notification.WebSocketPayload{
		UserIDs: []string{"u1"},
		Event:   "test",
		Data:    map[string]interface{}{"x": 1},
	})
	if err != nil {
		t.Fatal(err)
	}
	sent := s.Sent()
	if len(sent) != 1 || sent[0].Channel != "websocket" {
		t.Fatalf("want 1 websocket sent, got %d", len(sent))
	}
}

func TestSender_Reset(t *testing.T) {
	s := NewSender()
	ctx := context.Background()
	_ = s.SendEmail(ctx, notification.EmailPayload{To: "a@b.com", Subject: "x", Body: "y"})
	s.Reset()
	if len(s.Sent()) != 0 {
		t.Fatalf("want 0 after reset, got %d", len(s.Sent()))
	}
}
