package social

import (
	"context"
	"errors"
	"fmt"

	"github.com/parevo/core/storage"
)

var (
	ErrProviderNotFound = errors.New("provider not found")
	ErrInvalidIdentity  = errors.New("invalid social identity")
)

type Provider interface {
	Name() string
	ExchangeCode(ctx context.Context, code string, redirectURI string) (storage.SocialIdentity, error)
}

type TokenIssuer interface {
	IssueAccessToken(ctx context.Context, userID string, tenantID string) (string, error)
}

type Service struct {
	providers map[string]Provider
	store     storage.SocialAccountStore
	issuer    TokenIssuer
}

type CallbackResult struct {
	UserID      string
	AccessToken string
}

func NewService(store storage.SocialAccountStore, issuer TokenIssuer, providers ...Provider) *Service {
	providerMap := make(map[string]Provider, len(providers))
	for _, p := range providers {
		providerMap[p.Name()] = p
	}
	return &Service{
		providers: providerMap,
		store:     store,
		issuer:    issuer,
	}
}

func (s *Service) HandleCallback(ctx context.Context, providerName, code, redirectURI, tenantID string) (CallbackResult, error) {
	provider, ok := s.providers[providerName]
	if !ok {
		return CallbackResult{}, ErrProviderNotFound
	}

	identity, err := provider.ExchangeCode(ctx, code, redirectURI)
	if err != nil {
		return CallbackResult{}, err
	}
	if identity.Provider == "" || identity.ProviderUserID == "" || identity.Email == "" {
		return CallbackResult{}, ErrInvalidIdentity
	}

	userID, found, err := s.store.FindUserBySocial(ctx, identity.Provider, identity.ProviderUserID)
	if err != nil {
		return CallbackResult{}, err
	}
	if !found {
		userID, err = s.store.FindOrCreateUserByEmail(ctx, identity.Email, identity.Name)
		if err != nil {
			return CallbackResult{}, err
		}
		if err := s.store.LinkSocialAccount(ctx, userID, identity); err != nil {
			return CallbackResult{}, err
		}
	}

	token, err := s.issuer.IssueAccessToken(ctx, userID, tenantID)
	if err != nil {
		return CallbackResult{}, fmt.Errorf("issue token: %w", err)
	}

	return CallbackResult{
		UserID:      userID,
		AccessToken: token,
	}, nil
}
