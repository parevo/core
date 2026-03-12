package oauth2provider

import (
	"context"
	"testing"

	"github.com/parevo/core/storage/memory"
)

type mockIssuer struct {
	accessToken  string
	refreshToken string
	expiresIn    int
}

func (m *mockIssuer) IssueAccessToken(userID, clientID string, scopes []string) (string, int, error) {
	return m.accessToken, m.expiresIn, nil
}

func (m *mockIssuer) IssueRefreshToken(userID, clientID string, scopes []string) (string, error) {
	return m.refreshToken, nil
}

type mockClientStore struct {
	clients map[string]struct {
		secret       string
		redirectURIs []string
	}
}

func (m *mockClientStore) GetClient(ctx context.Context, clientID string) (*Client, bool, error) {
	c, ok := m.clients[clientID]
	if !ok {
		return nil, false, nil
	}
	return &Client{ID: clientID, Secret: c.secret, RedirectURIs: c.redirectURIs}, true, nil
}

func (m *mockClientStore) ValidateRedirect(ctx context.Context, clientID, redirectURI string) (bool, error) {
	c, ok := m.clients[clientID]
	if !ok {
		return false, nil
	}
	for _, u := range c.redirectURIs {
		if u == redirectURI {
			return true, nil
		}
	}
	return false, nil
}

func TestProvider_AuthorizeAndExchange(t *testing.T) {
	clients := &mockClientStore{
		clients: map[string]struct {
			secret       string
			redirectURIs []string
		}{
			"client1": {secret: "secret1", redirectURIs: []string{"https://app.com/cb"}},
		},
	}
	authCodes := memory.NewOAuth2AuthCodeStore()
	issuer := &mockIssuer{accessToken: "at", refreshToken: "rt", expiresIn: 3600}

	provider := NewProvider(clients, authCodes, issuer)
	ctx := context.Background()

	code, err := provider.Authorize(ctx, "client1", "https://app.com/cb", "user1", []string{"read"})
	if err != nil {
		t.Fatalf("Authorize failed: %v", err)
	}
	if code == "" {
		t.Error("expected non-empty code")
	}

	access, refresh, expiresIn, err := provider.Exchange(ctx, code, "client1", "secret1", "https://app.com/cb")
	if err != nil {
		t.Fatalf("Exchange failed: %v", err)
	}
	if access != "at" || refresh != "rt" || expiresIn != 3600 {
		t.Errorf("unexpected tokens: access=%s refresh=%s expiresIn=%d", access, refresh, expiresIn)
	}
}

func TestProvider_Exchange_InvalidClient(t *testing.T) {
	clients := &mockClientStore{clients: map[string]struct {
		secret       string
		redirectURIs []string
	}{"client1": {"s", []string{"https://a.com/cb"}}}}
	authCodes := memory.NewOAuth2AuthCodeStore()
	provider := NewProvider(clients, authCodes, &mockIssuer{})
	ctx := context.Background()

	code, _ := provider.Authorize(ctx, "client1", "https://a.com/cb", "u1", nil)

	_, _, _, err := provider.Exchange(ctx, code, "client1", "wrong-secret", "https://a.com/cb")
	if err != ErrInvalidClient {
		t.Errorf("wrong secret should fail: %v", err)
	}

	_, _, _, err = provider.Exchange(ctx, "invalid-code", "client1", "s", "https://a.com/cb")
	if err != ErrInvalidClient {
		t.Errorf("invalid code should fail: %v", err)
	}
}
