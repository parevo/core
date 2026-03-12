package gmail

import (
	"github.com/parevo/core/notification"
	"github.com/parevo/core/notification/smtp"
)

// Config holds Gmail SMTP settings.
// Use an App Password (not your regular Gmail password) when 2FA is enabled.
// Create at: https://myaccount.google.com/apppasswords
type Config struct {
	Email    string // Gmail address (e.g. user@gmail.com)
	AppPass  string // App Password (16 chars, no spaces)
}

// NewEmailProvider creates an email provider configured for Gmail SMTP.
func NewEmailProvider(cfg Config) notification.EmailProvider {
	return smtp.NewEmailProvider(smtp.Config{
		Host:     "smtp.gmail.com",
		Port:     587,
		Username: cfg.Email,
		Password: cfg.AppPass,
		From:     cfg.Email,
		TLS:      false, // Gmail uses STARTTLS on 587, smtp.SendMail handles it
	})
}
