package twilio

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/parevo/core/notification"
)

// Config holds Twilio credentials.
type Config struct {
	AccountSID string // Twilio Account SID
	AuthToken  string // Twilio Auth Token
	From       string // Twilio phone number (E.164, e.g. +15551234567)
}

// SMSProvider sends SMS via Twilio API.
type SMSProvider struct {
	cfg       Config
	client    *http.Client
	apiBase   string
}

// NewSMSProvider creates a Twilio SMS provider.
func NewSMSProvider(cfg Config) *SMSProvider {
	if cfg.AccountSID == "" || cfg.AuthToken == "" || cfg.From == "" {
		panic("twilio: AccountSID, AuthToken, and From are required")
	}
	return &SMSProvider{
		cfg:     cfg,
		client:  &http.Client{},
		apiBase: "https://api.twilio.com/2010-04-01/Accounts/" + cfg.AccountSID + "/Messages.json",
	}
}

// SendSMS sends the SMS via Twilio.
func (p *SMSProvider) SendSMS(ctx context.Context, payload notification.SMSPayload) error {
	form := url.Values{}
	form.Set("To", payload.To)
	form.Set("From", p.cfg.From)
	form.Set("Body", payload.Body)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.apiBase, strings.NewReader(form.Encode()))
	if err != nil {
		return fmt.Errorf("twilio request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(p.cfg.AccountSID, p.cfg.AuthToken)

	resp, err := p.client.Do(req)
	if err != nil {
		return fmt.Errorf("twilio do: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("twilio api: status %d", resp.StatusCode)
	}
	return nil
}

var _ notification.SMSProvider = (*SMSProvider)(nil)
