// Package billing provides tenant usage tracking and plan limit checks.
//
// # Usage
//
//	store := memory.NewUsageStore()
//	store.Record(ctx, "tenant-1", "api_calls", 100)
//	used, _ := store.Usage(ctx, "tenant-1", "api_calls", start, end)
//	withinLimit := used < planLimit
package billing

import (
	"context"
	"time"
)

// UsageStore records and retrieves tenant usage for billing/limits.
type UsageStore interface {
	Record(ctx context.Context, tenantID, metric string, value int64) error
	Usage(ctx context.Context, tenantID, metric string, from, to time.Time) (int64, error)
}

// PlanLimits defines limits per plan (optional, app-specific).
type PlanLimits struct {
	APIRequestsPerMonth int64
	StorageMB           int64
	Users               int64
}
