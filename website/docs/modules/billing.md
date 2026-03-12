# Billing Module

Tenant kullanım takibi ve plan limitleri.

## Usage

```go
import (
    "github.com/parevo/core/billing"
    "github.com/parevo/core/billing/memory"
)

store := memory.NewUsageStore()
store.Record(ctx, "tenant-1", "api_calls", 100)
store.Record(ctx, "tenant-1", "api_calls", 50)

used, _ := store.Usage(ctx, "tenant-1", "api_calls", startOfMonth, endOfMonth)
withinLimit := used < planLimit
```

## Plan Limits

`PlanLimits` struct'ı uygulama tarafında plan limitlerini tanımlamak için kullanılabilir.
