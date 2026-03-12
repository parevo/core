package graphql

import (
	"context"

	"github.com/parevo/core/auth"
	"github.com/parevo/core/auth/adapters"
	"github.com/parevo/core/tenant"
)

// ResolverContext injects auth into GraphQL resolver context.
type ResolverContext struct {
	Auth   *auth.Service
	Tenant *tenant.Service
	Opts   adapters.Options
}

// NewResolverContext creates a resolver context helper.
func NewResolverContext(authSvc *auth.Service, tenantSvc *tenant.Service, opts adapters.Options) *ResolverContext {
	opts = opts.WithDefaults()
	return &ResolverContext{Auth: authSvc, Tenant: tenantSvc, Opts: opts}
}

// AuthenticateRequest authenticates the request and returns a context with claims.
// Call this at the start of your GraphQL request handler (e.g. in the middleware before graphql.Execute).
func (r *ResolverContext) AuthenticateRequest(ctx context.Context, authHeader, tenantID string) (context.Context, error) {
	policy := r.Opts.OverridePolicy
	if policy == nil {
		policy = auth.StaticTenantOverridePolicy{Allow: false}
	}
	newCtx, _, err := r.Auth.AuthenticateContext(ctx, authHeader, tenantID, policy)
	return newCtx, err
}

// RequireAuth returns an error if the context has no valid claims.
func RequireAuth(ctx context.Context) error {
	_, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		return auth.ErrUnauthenticated
	}
	return nil
}

// RequireScope returns an error if claims don't have the required scope.
func RequireScope(ctx context.Context, scope string) error {
	claims, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		return auth.ErrUnauthenticated
	}
	if !auth.HasScope(claims, scope) {
		return auth.ErrForbidden
	}
	return nil
}
