package nethttpadapter

import (
	"net/http"

	"github.com/parevo/core/auth"
	"github.com/parevo/core/auth/ratelimit"
)

func TenantRateLimitMiddleware(limiter *ratelimit.TenantLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if limiter == nil {
				next.ServeHTTP(w, r)
				return
			}
			tenantID, _ := auth.TenantIDFromContext(r.Context())
			if tenantID == "" {
				tenantID = r.Header.Get("X-Tenant-Id")
			}
			subKey := clientIP(r)
			if !limiter.Allow(tenantID, subKey) {
				http.Error(w, "too many requests", http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
