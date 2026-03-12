package auth

import "context"

func RequireAuth(ctx context.Context) error {
	claims, ok := ClaimsFromContext(ctx)
	if !ok || claims == nil {
		return ErrUnauthenticated
	}
	return nil
}

func RequireTenant(ctx context.Context, tenantID string) error {
	if err := RequireAuth(ctx); err != nil {
		return err
	}

	claims, _ := ClaimsFromContext(ctx)
	if claims.TenantID != tenantID {
		return ErrTenantMismatch
	}
	return nil
}

func RequirePermission(ctx context.Context, permission string) error {
	if err := RequireAuth(ctx); err != nil {
		return err
	}
	claims, _ := ClaimsFromContext(ctx)
	for _, p := range claims.Permissions {
		if p == permission {
			return nil
		}
	}
	return ErrForbidden
}
