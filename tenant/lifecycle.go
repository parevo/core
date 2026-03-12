package tenant

import (
	"context"
	"errors"

	"github.com/parevo/core/storage"
)

var (
	ErrTenantExists   = errors.New("tenant already exists")
	ErrTenantNotFound = errors.New("tenant not found")
)

type LifecycleService struct {
	store storage.TenantLifecycleStore
}

func NewLifecycleService(store storage.TenantLifecycleStore) *LifecycleService {
	return &LifecycleService{store: store}
}

func (s *LifecycleService) Create(ctx context.Context, tenantID, name, ownerID string) error {
	if s.store == nil {
		return ErrTenantNotFound
	}
	return s.store.Create(ctx, tenantID, name, ownerID)
}

func (s *LifecycleService) Suspend(ctx context.Context, tenantID string) error {
	if s.store == nil {
		return ErrTenantNotFound
	}
	return s.store.Suspend(ctx, tenantID)
}

func (s *LifecycleService) Resume(ctx context.Context, tenantID string) error {
	if s.store == nil {
		return ErrTenantNotFound
	}
	return s.store.Resume(ctx, tenantID)
}

func (s *LifecycleService) Delete(ctx context.Context, tenantID string) error {
	if s.store == nil {
		return ErrTenantNotFound
	}
	return s.store.Delete(ctx, tenantID)
}

func (s *LifecycleService) Status(ctx context.Context, tenantID string) (storage.TenantStatus, error) {
	if s.store == nil {
		return "", ErrTenantNotFound
	}
	return s.store.Status(ctx, tenantID)
}
