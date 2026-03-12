package googleprovider

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"golang.org/x/oauth2"
)

func TestExchangeCode(t *testing.T) {
	tokenSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"access_token":"at-1","token_type":"Bearer","expires_in":3600}`))
	}))
	defer tokenSrv.Close()

	userInfoSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"sub":"g-123","email":"u@example.com","email_verified":true,"name":"User","picture":"https://example.com/p.png"}`))
	}))
	defer userInfoSrv.Close()

	provider := Provider{
		ProviderName: "google",
		Config: oauth2.Config{
			ClientID:     "x",
			ClientSecret: "y",
			RedirectURL:  "http://localhost/callback",
			Endpoint: oauth2.Endpoint{
				TokenURL: tokenSrv.URL,
			},
		},
		UserInfoURL: userInfoSrv.URL,
		HTTPClient:  userInfoSrv.Client(),
	}

	identity, err := provider.ExchangeCode(context.Background(), "code-1", "")
	if err != nil {
		t.Fatalf("exchange code failed: %v", err)
	}
	if identity.ProviderUserID != "g-123" || identity.Email != "u@example.com" {
		t.Fatalf("unexpected identity: %+v", identity)
	}

	if _, err := provider.ExchangeCode(context.Background(), "", ""); err == nil {
		t.Fatalf("expected missing code error")
	}
}
