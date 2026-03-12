package permission

import (
	"context"
	"errors"

	"github.com/parevo/core/storage"
)

var ErrPermissionDenied = errors.New("permission denied")

type Service struct {
	store storage.PermissionStore
}

func NewService(store storage.PermissionStore) *Service {
	return &Service{store: store}
}

func (s *Service) Check(ctx context.Context, subject storage.Subject, tenantID, permission string) error {
	if permission == "" {
		return ErrPermissionDenied
	}
	if s.store == nil {
		_ = ctx
		_ = subject
		_ = tenantID
		return ErrPermissionDenied
	}

	ok, err := s.store.HasPermission(ctx, subject.ID, tenantID, permission, subject.Roles)
	if err != nil {
		return err
	}
	if !ok {
		return ErrPermissionDenied
	}
	return nil
}

func (s *Service) CheckWithWildcard(_ context.Context, _ storage.Subject, _ string, permission string, grantedPermissions []string) bool {
	for _, g := range grantedPermissions {
		if MatchPermission(g, permission) {
			return true
		}
	}
	return false
}
