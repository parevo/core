// Permission cache example.
// HasPermission results cached with TTL; invalidation on role changes.
package main

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/parevo/core/auth"
	"github.com/parevo/core/auth/adapters"
	ginadapter "github.com/parevo/core/auth/adapters/gin"
	"github.com/parevo/core/permission"
	"github.com/parevo/core/storage/memory"
	"github.com/parevo/core/tenant"
)

func main() {
	basePerm := &memory.PermissionStore{
		Grants: map[string]bool{
			"user-1|tenant-a|orders:read": true,
		},
	}
	cachedPerm := permission.NewCachedPermissionStore(basePerm, 5*time.Minute)

	tenantStore := &memory.TenantStore{
		SubjectTenants: map[string][]string{"user-1": {"tenant-a"}},
	}

	svc, _ := auth.NewServiceWithModules(auth.Config{
		Issuer:    "parevo",
		Audience:  "parevo-api",
		SecretKey: []byte("change-this-in-production"),
	}, auth.Modules{
		Tenant:         tenant.NewService(tenantStore),
		Permission:     permission.NewService(cachedPerm),
		TenantOverride: tenant.StaticOverridePolicy{Allow: true},
	})

	r := gin.New()
	r.Use(gin.Recovery())
	r.GET("/orders", ginadapter.AuthMiddleware(svc, adapters.Options{}),
		ginadapter.RequirePermissionWithService(svc, "orders:read"),
		func(c *gin.Context) {
			c.JSON(200, gin.H{"orders": []string{}})
		})

	// On role change: cachedPerm.InvalidateSubject("user-1")
	// Or tenant-scoped: cachedPerm.InvalidateSubjectTenant("user-1", "tenant-a")

	_ = r.Run(":8084")
}
