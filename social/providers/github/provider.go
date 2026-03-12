package githubprovider

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
	"github.com/parevo/core/storage"
)

var ErrMissingCode = errors.New("missing authorization code")

var endpoint = oauth2.Endpoint{
	AuthURL:  "https://github.com/login/oauth/authorize",
	TokenURL: "https://github.com/login/oauth/access_token",
}

type Provider struct {
	Config       oauth2.Config
	UserInfoURL  string
	HTTPClient   *http.Client
	ProviderName string
}

func New(clientID, clientSecret, redirectURL string, scopes ...string) Provider {
	if len(scopes) == 0 {
		scopes = []string{"user:email", "read:user"}
	}
	return Provider{
		Config: oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Scopes:       scopes,
			Endpoint:     endpoint,
		},
		UserInfoURL:  "https://api.github.com/user",
		ProviderName: "github",
	}
}

func (p Provider) Name() string {
	if p.ProviderName == "" {
		return "github"
	}
	return p.ProviderName
}

func (p Provider) ExchangeCode(ctx context.Context, code string, _ string) (storage.SocialIdentity, error) {
	if code == "" {
		return storage.SocialIdentity{}, ErrMissingCode
	}
	token, err := p.Config.Exchange(ctx, code)
	if err != nil {
		return storage.SocialIdentity{}, err
	}

	client := p.HTTPClient
	if client == nil {
		client = p.Config.Client(ctx, token)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, p.UserInfoURL, nil)
	if err != nil {
		return storage.SocialIdentity{}, err
	}
	req.Header.Set("Accept", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return storage.SocialIdentity{}, err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode >= 300 {
		return storage.SocialIdentity{}, fmt.Errorf("userinfo request failed: %d", resp.StatusCode)
	}

	var payload struct {
		ID    int    `json:"id"`
		Login string `json:"login"`
		Email string `json:"email"`
		Name  string `json:"name"`
		Avatar string `json:"avatar_url"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return storage.SocialIdentity{}, err
	}
	providerUserID := fmt.Sprintf("%d", payload.ID)
	email := payload.Email
	if email == "" {
		email = payload.Login + "@users.noreply.github.com"
	}
	return storage.SocialIdentity{
		Provider:       p.Name(),
		ProviderUserID: providerUserID,
		Email:          email,
		EmailVerified:  payload.Email != "",
		Name:           payload.Name,
		AvatarURL:      payload.Avatar,
	}, nil
}
