package tenant

import (
	"context"
	"errors"

	"github.com/parevo/core/storage"
)

var ErrFeatureDisabled = errors.New("feature disabled for tenant")

// PolicyBinding binds a tenant to a plan with optional overrides.
type PolicyBinding struct {
	Plan      storage.Plan
	Overrides map[string]bool  // feature -> enabled
	Limits    map[string]int64  // limitKey -> value
}

// FeatureService checks tenant feature flags and limits.
type FeatureService struct {
	store storage.FeatureFlagStore
}

// NewFeatureService creates a FeatureService.
func NewFeatureService(store storage.FeatureFlagStore) *FeatureService {
	return &FeatureService{store: store}
}

// IsEnabled returns true if the feature is enabled for the tenant.
func (s *FeatureService) IsEnabled(ctx context.Context, tenantID, feature string) (bool, error) {
	if s.store == nil {
		return true, nil // default allow when no store
	}
	return s.store.IsEnabled(ctx, tenantID, feature)
}

// RequireEnabled returns ErrFeatureDisabled if the feature is not enabled.
func (s *FeatureService) RequireEnabled(ctx context.Context, tenantID, feature string) error {
	ok, err := s.IsEnabled(ctx, tenantID, feature)
	if err != nil {
		return err
	}
	if !ok {
		return ErrFeatureDisabled
	}
	return nil
}

// GetPlan returns the tenant's plan.
func (s *FeatureService) GetPlan(ctx context.Context, tenantID string) (storage.Plan, error) {
	if s.store == nil {
		return storage.PlanFree, nil
	}
	return s.store.GetPlan(ctx, tenantID)
}

// GetLimit returns the tenant's limit for the given key (e.g. "max_users", "api_calls_per_month").
func (s *FeatureService) GetLimit(ctx context.Context, tenantID, limitKey string) (int64, error) {
	if s.store == nil {
		return 0, nil
	}
	return s.store.GetLimit(ctx, tenantID, limitKey)
}

// CheckLimit returns nil if current < limit, otherwise ErrLimitExceeded.
var ErrLimitExceeded = errors.New("tenant limit exceeded")

func (s *FeatureService) CheckLimit(ctx context.Context, tenantID, limitKey string, current int64) error {
	limit, err := s.GetLimit(ctx, tenantID, limitKey)
	if err != nil {
		return err
	}
	if limit > 0 && current >= limit {
		return ErrLimitExceeded
	}
	return nil
}
