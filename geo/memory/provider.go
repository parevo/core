package memory

import (
	"context"

	"github.com/parevo/core/geo"
)

// Provider implements geo.Provider with stub responses (dev/test).
type Provider struct {
	Default *geo.Location // optional default for all IPs
}

// NewProvider creates a stub geo provider.
func NewProvider() *Provider {
	return &Provider{}
}

// Lookup returns the default location or empty. For dev/test only.
func (p *Provider) Lookup(ctx context.Context, ip string) (*geo.Location, error) {
	if p.Default != nil {
		return &geo.Location{
			Country: p.Default.Country,
			Region:  p.Default.Region,
			City:    p.Default.City,
		}, nil
	}
	return &geo.Location{}, nil
}

var _ geo.Provider = (*Provider)(nil)
