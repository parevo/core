package main

import (
	"github.com/gin-gonic/gin"

	"github.com/parevo/core/auth"
	"github.com/parevo/core/auth/adapters"
	ginadapter "github.com/parevo/core/auth/adapters/gin"
	"github.com/parevo/core/permission"
	"github.com/parevo/core/storage/memory"
	"github.com/parevo/core/tenant"
)

func main() {
	tenantStore := &memory.TenantStore{
		SubjectTenants: map[string][]string{
			"user-1": {"tenant-a", "tenant-b"},
		},
	}
	permissionStore := &memory.PermissionStore{
		Grants: map[string]bool{
			"user-1|tenant-a|orders:read": true,
		},
	}

	svc, err := auth.NewServiceWithModules(auth.Config{
		Issuer:    "parevo",
		Audience:  "parevo-api",
		SecretKey: []byte("change-this-in-production"),
	}, auth.Modules{
		Tenant:         tenant.NewService(tenantStore),
		Permission:     permission.NewService(permissionStore),
		TenantOverride: tenant.StaticOverridePolicy{Allow: true},
	})
	if err != nil {
		panic(err)
	}

	r := gin.New()
	r.Use(gin.Recovery())
	r.GET("/secure", ginadapter.AuthMiddleware(svc, adapters.Options{}), ginadapter.RequirePermissionWithService(svc, "orders:read"), func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})
	_ = r.Run(":8081")
}
