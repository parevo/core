package googleprovider

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"github.com/parevo/core/storage"
)

var ErrMissingCode = errors.New("missing authorization code")

type Provider struct {
	Config       oauth2.Config
	UserInfoURL  string
	HTTPClient   *http.Client
	ProviderName string
}

func New(clientID, clientSecret, redirectURL string, scopes ...string) Provider {
	if len(scopes) == 0 {
		scopes = []string{"openid", "email", "profile"}
	}
	return Provider{
		Config: oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Scopes:       scopes,
			Endpoint:     google.Endpoint,
		},
		UserInfoURL:  "https://openidconnect.googleapis.com/v1/userinfo",
		ProviderName: "google",
	}
}

func (p Provider) Name() string {
	if p.ProviderName == "" {
		return "google"
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

	userInfoURL := p.UserInfoURL
	if userInfoURL == "" {
		userInfoURL = "https://openidconnect.googleapis.com/v1/userinfo"
	}
	client := p.HTTPClient
	if client == nil {
		client = p.Config.Client(ctx, token)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, userInfoURL, nil)
	if err != nil {
		return storage.SocialIdentity{}, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return storage.SocialIdentity{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return storage.SocialIdentity{}, fmt.Errorf("userinfo request failed: %d", resp.StatusCode)
	}

	var payload struct {
		Sub           string `json:"sub"`
		Email         string `json:"email"`
		EmailVerified bool   `json:"email_verified"`
		Name          string `json:"name"`
		Picture       string `json:"picture"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return storage.SocialIdentity{}, err
	}
	return storage.SocialIdentity{
		Provider:       p.Name(),
		ProviderUserID: payload.Sub,
		Email:          payload.Email,
		EmailVerified:  payload.EmailVerified,
		Name:           payload.Name,
		AvatarURL:      payload.Picture,
	}, nil
}
