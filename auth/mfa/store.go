package mfa

import "context"

type TOTPStore interface {
	GetSecret(ctx context.Context, userID string) (secret string, enabled bool, err error)
	SetSecret(ctx context.Context, userID, secret string) error
	Enable(ctx context.Context, userID string) error
	Disable(ctx context.Context, userID string) error
}
