package echoadapter

import (
	"github.com/labstack/echo/v4"

	"github.com/parevo/core/auth"
	"github.com/parevo/core/auth/adapters"
)

func AuthMiddleware(service *auth.Service, opts adapters.Options) echo.MiddlewareFunc {
	opts = opts.WithDefaults()

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx, claims, err := service.AuthenticateContext(
				c.Request().Context(),
				c.Request().Header.Get("Authorization"),
				c.Request().Header.Get(opts.TenantHeader),
				opts.OverridePolicy,
			)
			if err != nil {
				status, body := opts.ErrorHandler(err)
				return c.JSON(status, map[string]string{"error": body})
			}

			c.Set("user_id", claims.UserID)
			c.Set("tenant_id", claims.TenantID)
			c.Set("roles", claims.Roles)
			c.Set("permissions", claims.Permissions)
			c.SetRequest(c.Request().WithContext(ctx))

			return next(c)
		}
	}
}

func RequirePermission(permission string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if err := auth.RequirePermission(c.Request().Context(), permission); err != nil {
				status, body := adapters.DefaultErrorHandler(err)
				return c.JSON(status, map[string]string{"error": body})
			}
			return next(c)
		}
	}
}

func RequirePermissionWithService(service *auth.Service, permission string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if err := service.AuthorizePermission(c.Request().Context(), permission); err != nil {
				status, body := adapters.DefaultErrorHandler(err)
				return c.JSON(status, map[string]string{"error": body})
			}
			return next(c)
		}
	}
}
