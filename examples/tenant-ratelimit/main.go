// Tenant-based rate limit example.
// Per-tenant QPS limits. Tenant is determined by X-Tenant-Id header.
package main

import (
	"fmt"
	"net/http"
	"time"

	nethttpadapter "github.com/parevo/core/auth/adapters/nethttp"
	"github.com/parevo/core/auth/ratelimit"
)

func main() {
	// 3 requests/minute per tenant (default)
	limiter := ratelimit.NewTenantLimiter(3, time.Minute)
	// VIP tenant: 10 requests/minute
	limiter.SetTenantLimit("vip", 10, time.Minute)

	handler := nethttpadapter.TenantRateLimitMiddleware(limiter)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tenantID := r.Header.Get("X-Tenant-Id")
		if tenantID == "" {
			tenantID = "_default"
		}
		_, _ = w.Write([]byte(fmt.Sprintf("ok tenant=%s\n", tenantID)))
	}))

	mux := http.NewServeMux()
	mux.Handle("/api/", handler)
	fmt.Println("Tenant rate limit demo: curl -H 'X-Tenant-Id: vip' http://localhost:8083/api/")
	fmt.Println("Listen :8083")
	_ = http.ListenAndServe(":8083", mux)
}
