package memory

import (
	"context"
	"testing"

	"github.com/parevo/core/geo"
)

func TestProvider(t *testing.T) {
	ctx := context.Background()
	p := NewProvider()

	loc, err := p.Lookup(ctx, "8.8.8.8")
	if err != nil {
		t.Fatal(err)
	}
	if loc == nil {
		t.Fatal("location is nil")
	}
}

func TestProviderWithDefault(t *testing.T) {
	ctx := context.Background()
	p := &Provider{Default: &geo.Location{Country: "US", Region: "CA", City: "SF"}}

	loc, err := p.Lookup(ctx, "1.2.3.4")
	if err != nil {
		t.Fatal(err)
	}
	if loc.Country != "US" || loc.Region != "CA" || loc.City != "SF" {
		t.Errorf("got %+v", loc)
	}
}
