package ginadapter

import (
	"github.com/gin-gonic/gin"

	"github.com/parevo/core/auth"
	"github.com/parevo/core/auth/adapters"
)

func AuthMiddleware(service *auth.Service, opts adapters.Options) gin.HandlerFunc {
	opts = opts.WithDefaults()

	return func(c *gin.Context) {
		ctx, claims, err := service.AuthenticateContext(
			c.Request.Context(),
			c.GetHeader("Authorization"),
			c.GetHeader(opts.TenantHeader),
			opts.OverridePolicy,
		)
		if err != nil {
			status, body := opts.ErrorHandler(err)
			c.AbortWithStatusJSON(status, gin.H{"error": body})
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("tenant_id", claims.TenantID)
		c.Set("roles", claims.Roles)
		c.Set("permissions", claims.Permissions)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func RequirePermission(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := auth.RequirePermission(c.Request.Context(), permission); err != nil {
			status, body := adapters.DefaultErrorHandler(err)
			c.AbortWithStatusJSON(status, gin.H{"error": body})
			return
		}
		c.Next()
	}
}

func RequirePermissionWithService(service *auth.Service, permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := service.AuthorizePermission(c.Request.Context(), permission); err != nil {
			status, body := adapters.DefaultErrorHandler(err)
			c.AbortWithStatusJSON(status, gin.H{"error": body})
			return
		}
		c.Next()
	}
}
