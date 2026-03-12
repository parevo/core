# Notification Module

Unified sender for email, SMS, WebSocket.

## Providers

- `notification/memory` — dev/test
- `notification/smtp` — SMTP email
- `notification/gmail` — Gmail
- `notification/ses` — Amazon SES
- `notification/twilio` — Twilio SMS

## Usage

```go
sender := notification.NewSender(
    gmail.NewEmailProvider(gmail.Config{Email: "...", AppPass: "..."}),
    twilio.NewSMSProvider(twilio.Config{...}),
    memory.NewWebSocketProvider(),
)
sender.SendEmail(ctx, notification.EmailPayload{To: "...", Subject: "...", Body: "..."})
```
