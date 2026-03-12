package ses

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
	"github.com/parevo/core/notification"
)

// Config holds Amazon SES settings.
type Config struct {
	Region string // AWS region (e.g. us-east-1)
	From   string // Verified sender email
}

// EmailProvider sends email via Amazon SES.
type EmailProvider struct {
	client *sesv2.Client
	from   string
}

// NewEmailProvider creates an SES email provider.
// Uses default AWS credential chain (env, shared config, IAM role).
func NewEmailProvider(cfg Config) (*EmailProvider, error) {
	if cfg.Region == "" {
		cfg.Region = "us-east-1"
	}
	if cfg.From == "" {
		return nil, fmt.Errorf("ses: From is required")
	}
	awsCfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(cfg.Region))
	if err != nil {
		return nil, fmt.Errorf("ses config: %w", err)
	}
	return &EmailProvider{
		client: sesv2.NewFromConfig(awsCfg),
		from:   cfg.From,
	}, nil
}

// SendEmail sends the email via Amazon SES.
func (p *EmailProvider) SendEmail(ctx context.Context, payload notification.EmailPayload) error {
	body := &types.Body{}
	if payload.HTML != "" {
		body.Html = &types.Content{Data: aws.String(payload.HTML)}
	}
	if payload.Body != "" {
		body.Text = &types.Content{Data: aws.String(payload.Body)}
	}
	if body.Html == nil && body.Text == nil {
		body.Text = &types.Content{Data: aws.String("")}
	}

	dest := types.Destination{ToAddresses: []string{payload.To}}
	if len(payload.Cc) > 0 {
		dest.CcAddresses = payload.Cc
	}
	if len(payload.Bcc) > 0 {
		dest.BccAddresses = payload.Bcc
	}

	_, err := p.client.SendEmail(ctx, &sesv2.SendEmailInput{
		FromEmailAddress: aws.String(p.from),
		Destination:     &dest,
		Content: &types.EmailContent{
			Simple: &types.Message{
				Subject: &types.Content{Data: aws.String(payload.Subject)},
				Body:    body,
			},
		},
	})
	if err != nil {
		return fmt.Errorf("ses send: %w", err)
	}
	return nil
}

var _ notification.EmailProvider = (*EmailProvider)(nil)
