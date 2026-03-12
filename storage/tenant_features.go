package storage

import "context"

// Plan represents a tenant's subscription plan.
type Plan string

const (
	PlanFree       Plan = "free"
	PlanPro        Plan = "pro"
	PlanEnterprise Plan = "enterprise"
)

// FeatureFlagStore resolves tenant-level feature flags and plan limits.
type FeatureFlagStore interface {
	GetPlan(ctx context.Context, tenantID string) (Plan, error)
	IsEnabled(ctx context.Context, tenantID string, feature string) (bool, error)
	GetLimit(ctx context.Context, tenantID string, limitKey string) (int64, error)
}

// TenantFeatureStore extends FeatureFlagStore with write operations for admin.
type TenantFeatureStore interface {
	FeatureFlagStore
	SetPlan(ctx context.Context, tenantID string, plan Plan) error
	SetFeature(ctx context.Context, tenantID, feature string, enabled bool) error
	SetLimit(ctx context.Context, tenantID, limitKey string, value int64) error
}
