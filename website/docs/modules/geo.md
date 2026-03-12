# Geo Module

IP geolocation (ülke, bölge, şehir).

## Providers

- `geo/memory` — stub (dev/test)

## Usage

```go
import (
    "github.com/parevo/core/geo"
    "github.com/parevo/core/geo/memory"
)

provider := memory.NewProvider()
loc, _ := provider.Lookup(ctx, "8.8.8.8")
// loc.Country, loc.Region, loc.City

// Varsayılan değer ile
provider.Default = &geo.Location{Country: "US", Region: "CA"}
```

## Production

MaxMind GeoIP2 vb. için `geo.Provider` interface'ini implement eden bir adapter eklenebilir.
