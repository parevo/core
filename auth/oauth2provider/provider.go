package oauth2provider

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"
)

var (
	ErrInvalidClient   = errors.New("invalid client")
	ErrInvalidCode     = errors.New("invalid or expired authorization code")
	ErrInvalidScope    = errors.New("invalid scope")
	ErrAccessDenied    = errors.New("access denied")
	ErrConsentRequired = errors.New("consent required")
)

// Client represents an OAuth2 client application.
type Client struct {
	ID           string
	Secret       string
	RedirectURIs []string
	Scopes       []string
}

// ClientStore looks up OAuth2 clients.
type ClientStore interface {
	GetClient(ctx context.Context, clientID string) (*Client, bool, error)
	ValidateRedirect(ctx context.Context, clientID, redirectURI string) (bool, error)
}

// AuthCodeStore manages authorization codes.
type AuthCodeStore interface {
	Create(ctx context.Context, code, clientID, userID, redirectURI string, scopes []string, expiresAt time.Time) error
	Consume(ctx context.Context, code string) (clientID, userID, redirectURI string, scopes []string, err error)
}

// Provider implements OAuth2 authorization server flows.
type Provider struct {
	clients   ClientStore
	authCodes AuthCodeStore
	issuer    Issuer
}

// Issuer issues access and refresh tokens.
type Issuer interface {
	IssueAccessToken(userID, clientID string, scopes []string) (accessToken string, expiresIn int, err error)
	IssueRefreshToken(userID, clientID string, scopes []string) (refreshToken string, err error)
}

// NewProvider creates an OAuth2 provider.
func NewProvider(clients ClientStore, authCodes AuthCodeStore, issuer Issuer) *Provider {
	return &Provider{clients: clients, authCodes: authCodes, issuer: issuer}
}

// Authorize creates an authorization code for the user.
func (p *Provider) Authorize(ctx context.Context, clientID, redirectURI, userID string, scopes []string) (code string, err error) {
	_, ok, err := p.clients.GetClient(ctx, clientID)
	if err != nil || !ok {
		return "", ErrInvalidClient
	}
	valid, err := p.clients.ValidateRedirect(ctx, clientID, redirectURI)
	if err != nil || !valid {
		return "", ErrInvalidClient
	}
	code = generateCode()
	expiresAt := time.Now().Add(10 * time.Minute)
	if err := p.authCodes.Create(ctx, code, clientID, userID, redirectURI, scopes, expiresAt); err != nil {
		return "", err
	}
	return code, nil
}

// Exchange exchanges an auth code for tokens.
func (p *Provider) Exchange(ctx context.Context, code, clientID, clientSecret, redirectURI string) (accessToken, refreshToken string, expiresIn int, err error) {
	client, ok, err := p.clients.GetClient(ctx, clientID)
	if err != nil || !ok || client.Secret != clientSecret {
		return "", "", 0, ErrInvalidClient
	}
	cid, userID, ruri, scopes, err := p.authCodes.Consume(ctx, code)
	if err != nil || cid == "" || cid != clientID || ruri != redirectURI {
		return "", "", 0, ErrInvalidClient
	}
	_ = userID
	_ = scopes
	accessToken, expiresIn, err = p.issuer.IssueAccessToken(userID, clientID, scopes)
	if err != nil {
		return "", "", 0, err
	}
	refreshToken, _ = p.issuer.IssueRefreshToken(userID, clientID, scopes)
	return accessToken, refreshToken, expiresIn, nil
}

func generateCode() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}
