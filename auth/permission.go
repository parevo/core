package auth

import (
	"context"
	"errors"

	"github.com/parevo/core/permission"
	"github.com/parevo/core/storage"
)

func (s *Service) AuthorizePermission(ctx context.Context, perm string) error {
	if err := RequireAuth(ctx); err != nil {
		s.audit(ctx, AuditPermissionDeny, map[string]string{"reason": "unauthenticated", "permission": perm})
		return err
	}
	claims, _ := ClaimsFromContext(ctx)
	if s.modules.Permission == nil {
		for _, p := range claims.Permissions {
			if p == perm {
				return nil
			}
		}
		s.audit(ctx, AuditPermissionDeny, map[string]string{"user_id": claims.UserID, "tenant_id": claims.TenantID, "permission": perm})
		return ErrForbidden
	}

	subject := storage.Subject{
		ID:    claims.UserID,
		Roles: claims.Roles,
	}
	if err := s.modules.Permission.Check(ctx, subject, claims.TenantID, perm); err != nil {
		if errors.Is(err, permission.ErrPermissionDenied) {
			s.audit(ctx, AuditPermissionDeny, map[string]string{"user_id": claims.UserID, "tenant_id": claims.TenantID, "permission": perm})
			return ErrForbidden
		}
		return err
	}
	return nil
}
