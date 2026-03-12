package chiadapter

import (
	"net/http"

	"github.com/parevo/core/auth"
	"github.com/parevo/core/auth/adapters"
	nethttpadapter "github.com/parevo/core/auth/adapters/nethttp"
	"github.com/parevo/core/auth/ratelimit"
)

func AuthMiddleware(service *auth.Service, opts adapters.Options) func(http.Handler) http.Handler {
	return nethttpadapter.AuthMiddleware(service, opts)
}

func RequirePermission(permission string) func(http.Handler) http.Handler {
	return nethttpadapter.RequirePermission(permission)
}

func RequirePermissionWithService(service *auth.Service, permission string) func(http.Handler) http.Handler {
	return nethttpadapter.RequirePermissionWithService(service, permission)
}

func RateLimitByIP(limiter *ratelimit.MemoryLimiter) func(http.Handler) http.Handler {
	return nethttpadapter.RateLimitByIP(limiter)
}

func TenantRateLimitMiddleware(limiter *ratelimit.TenantLimiter) func(http.Handler) http.Handler {
	return nethttpadapter.TenantRateLimitMiddleware(limiter)
}
