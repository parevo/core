package tenant

import (
	"context"
	"errors"

	"github.com/parevo/core/storage"
)

var (
	ErrAccessDenied = errors.New("tenant access denied")
	ErrNoTenant     = errors.New("tenant not found")
)

type OverridePolicy interface {
	CanOverride(subject storage.Subject, requestedTenantID string) bool
}

type StaticOverridePolicy struct {
	Allow bool
}

func (p StaticOverridePolicy) CanOverride(_ storage.Subject, _ string) bool {
	return p.Allow
}

type Service struct {
	store storage.TenantStore
}

func NewService(store storage.TenantStore) *Service {
	return &Service{store: store}
}

// Resolve chooses active tenant for subject.
// Priority: requestedTenantID (if allowed) -> defaultTenantID.
func (s *Service) Resolve(ctx context.Context, subject storage.Subject, defaultTenantID, requestedTenantID string, policy OverridePolicy) (string, error) {
	if requestedTenantID != "" && requestedTenantID != defaultTenantID {
		if policy == nil || !policy.CanOverride(subject, requestedTenantID) {
			return "", ErrAccessDenied
		}
		if err := s.ensureAccessible(ctx, subject.ID, requestedTenantID); err != nil {
			return "", err
		}
		return requestedTenantID, nil
	}

	if defaultTenantID == "" && requestedTenantID == "" {
		return "", nil // tenant-less mode
	}
	if defaultTenantID == "" {
		return "", ErrNoTenant
	}
	if err := s.ensureAccessible(ctx, subject.ID, defaultTenantID); err != nil {
		return "", err
	}
	return defaultTenantID, nil
}

func (s *Service) ensureAccessible(ctx context.Context, subjectID, tenantID string) error {
	if s.store == nil {
		return nil
	}
	tenants, err := s.store.ResolveSubjectTenants(ctx, subjectID)
	if err != nil {
		return err
	}
	for _, t := range tenants {
		if t == tenantID {
			return nil
		}
	}
	return ErrAccessDenied
}
