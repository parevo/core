package main

import (
	"context"
	"fmt"

	"github.com/parevo/core/notification"
	"github.com/parevo/core/notification/memory"
)

func main() {
	// Option 1: All in-memory (dev/test)
	sender := memory.NewSender()

	// Option 2: Gmail + Twilio + memory WebSocket
	//   sender := notification.NewSender(
	//     gmail.NewEmailProvider(gmail.Config{Email: "you@gmail.com", AppPass: "..."}),
	//     twilio.NewSMSProvider(twilio.Config{AccountSID: "...", AuthToken: "...", From: "+15551234567"}),
	//     memory.NewWebSocketProvider(),
	//   )
	// Option 3: Amazon SES + Twilio
	//   sesProv, _ := ses.NewEmailProvider(ses.Config{Region: "us-east-1", From: "noreply@example.com"})
	//   sender := notification.NewSender(sesProv, twilio.NewSMSProvider(twilio.Config{...}), nop.WebSocketProvider{})

	ctx := context.Background()

	// Email (tenant-less: UserID only)
	if err := sender.SendEmail(ctx, notification.EmailPayload{
		To:      "user@example.com",
		Subject: "Welcome",
		Body:    "Welcome to the app.",
		UserID:  "user-1",
	}); err != nil {
		panic(err)
	}

	// SMS (tenant flow: UserID + TenantID)
	if err := sender.SendSMS(ctx, notification.SMSPayload{
		To:       "+15551234567",
		Body:     "Your code is 123456",
		UserID:   "user-1",
		TenantID: "tenant-a",
	}); err != nil {
		panic(err)
	}

	// WebSocket (tenant flow: UserIDs + TenantID for context)
	if err := sender.SendWebSocket(ctx, notification.WebSocketPayload{
		UserIDs:  []string{"user-1"},
		TenantID: "tenant-a",
		Event:    "notification.new",
		Data: map[string]interface{}{
			"title": "New message",
			"body":  "You have a new message",
		},
	}); err != nil {
		panic(err)
	}

	fmt.Printf("Sent %d notifications (check logs above)\n", len(sender.Sent()))
}
