package smtp

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/smtp"
	"strings"

	"github.com/parevo/core/notification"
)

// Config holds SMTP connection settings.
type Config struct {
	Host     string // e.g. "smtp.example.com"
	Port     int    // e.g. 587
	Username string
	Password string
	From     string // default From address
	TLS      bool   // use STARTTLS
}

// EmailProvider sends email via SMTP.
type EmailProvider struct {
	cfg Config
}

// NewEmailProvider creates an SMTP email provider.
func NewEmailProvider(cfg Config) *EmailProvider {
	if cfg.Port == 0 {
		cfg.Port = 587
	}
	return &EmailProvider{cfg: cfg}
}

// SendEmail sends the email via SMTP.
func (p *EmailProvider) SendEmail(ctx context.Context, payload notification.EmailPayload) error {
	addr := fmt.Sprintf("%s:%d", p.cfg.Host, p.cfg.Port)
	auth := smtp.PlainAuth("", p.cfg.Username, p.cfg.Password, p.cfg.Host)

	from := p.cfg.From
	if from == "" {
		from = p.cfg.Username
	}

	msg := buildMessage(from, payload.To, payload.Subject, payload.Body, payload.HTML, payload.Cc, payload.Bcc)

	if p.cfg.TLS {
		tlsConfig := &tls.Config{ServerName: p.cfg.Host}
		conn, err := tls.Dial("tcp", addr, tlsConfig)
		if err != nil {
			return fmt.Errorf("smtp tls dial: %w", err)
		}
		defer func() { _ = conn.Close() }()
		client, err := smtp.NewClient(conn, p.cfg.Host)
		if err != nil {
			return fmt.Errorf("smtp client: %w", err)
		}
		defer func() { _ = client.Close() }()
		if err := send(client, auth, from, payload, msg); err != nil {
			return err
		}
		return nil
	}

	toAddrs := []string{payload.To}
	toAddrs = append(toAddrs, payload.Cc...)
	toAddrs = append(toAddrs, payload.Bcc...)
	return smtp.SendMail(addr, auth, from, toAddrs, msg)
}

func buildMessage(from, to, subject, body, html string, cc, bcc []string) []byte {
	var b strings.Builder
	b.WriteString("From: " + from + "\r\n")
	b.WriteString("To: " + to + "\r\n")
	if len(cc) > 0 {
		b.WriteString("Cc: " + strings.Join(cc, ", ") + "\r\n")
	}
	b.WriteString("Subject: " + subject + "\r\n")
	b.WriteString("MIME-Version: 1.0\r\n")
	if html != "" {
		b.WriteString("Content-Type: text/html; charset=UTF-8\r\n")
		b.WriteString("\r\n")
		b.WriteString(html)
	} else {
		b.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
		b.WriteString("\r\n")
		b.WriteString(body)
	}
	return []byte(b.String())
}

func send(client *smtp.Client, auth smtp.Auth, from string, payload notification.EmailPayload, msg []byte) error {
	if err := client.Auth(auth); err != nil {
		return fmt.Errorf("smtp auth: %w", err)
	}
	if err := client.Mail(from); err != nil {
		return fmt.Errorf("smtp mail: %w", err)
	}
	toAddrs := append([]string{payload.To}, payload.Cc...)
	toAddrs = append(toAddrs, payload.Bcc...)
	for _, addr := range toAddrs {
		if addr == "" {
			continue
		}
		if err := client.Rcpt(addr); err != nil {
			return fmt.Errorf("smtp rcpt: %w", err)
		}
	}
	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("smtp data: %w", err)
	}
	defer func() { _ = w.Close() }()
	if _, err := w.Write(msg); err != nil {
		return fmt.Errorf("smtp write: %w", err)
	}
	return nil
}

var _ notification.EmailProvider = (*EmailProvider)(nil)
