package auth

import "context"

type contextKey string

const (
	claimsContextKey contextKey = "parevo.auth.claims"
	tenantContextKey contextKey = "parevo.auth.tenant_id"
	userContextKey   contextKey = "parevo.auth.user_id"
)

func WithClaims(ctx context.Context, claims *Claims) context.Context {
	ctx = context.WithValue(ctx, claimsContextKey, claims)
	ctx = context.WithValue(ctx, tenantContextKey, claims.TenantID)
	ctx = context.WithValue(ctx, userContextKey, claims.UserID)
	ctx = context.WithValue(ctx, "parevo.tenant_id", claims.TenantID)
	ctx = context.WithValue(ctx, "parevo.user_id", claims.UserID)
	return ctx
}

func ClaimsFromContext(ctx context.Context) (*Claims, bool) {
	claims, ok := ctx.Value(claimsContextKey).(*Claims)
	return claims, ok
}

func TenantIDFromContext(ctx context.Context) (string, bool) {
	tenantID, ok := ctx.Value(tenantContextKey).(string)
	return tenantID, ok
}

func UserIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(userContextKey).(string)
	return userID, ok
}
