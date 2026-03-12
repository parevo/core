// Package geo provides IP geolocation lookup for country, region, and city.
//
// # Usage
//
//	provider := memory.NewProvider() // stub
//	loc, _ := provider.Lookup(ctx, "8.8.8.8")
//	// loc.Country, loc.Region, loc.City
package geo

import (
	"context"
)

// Location holds geolocation data for an IP.
type Location struct {
	Country string
	Region  string
	City    string
}

// Provider looks up IP geolocation.
type Provider interface {
	Lookup(ctx context.Context, ip string) (*Location, error)
}
