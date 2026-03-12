package mockprovider

import (
	"context"

	"github.com/parevo/core/storage"
)

type Provider struct {
	ProviderName string
}

func (p Provider) Name() string {
	if p.ProviderName == "" {
		return "mock"
	}
	return p.ProviderName
}

func (p Provider) ExchangeCode(_ context.Context, code string, _ string) (storage.SocialIdentity, error) {
	return storage.SocialIdentity{
		Provider:       p.Name(),
		ProviderUserID: "pid-" + code,
		Email:          code + "@example.com",
		EmailVerified:  true,
		Name:           "Mock User",
	}, nil
}
