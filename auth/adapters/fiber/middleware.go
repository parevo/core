package fiberadapter

import (
	"context"

	"github.com/gofiber/fiber/v2"

	"github.com/parevo/core/auth"
	"github.com/parevo/core/auth/adapters"
)

func AuthMiddleware(service *auth.Service, opts adapters.Options) fiber.Handler {
	opts = opts.WithDefaults()

	return func(c *fiber.Ctx) error {
		ctx, claims, err := service.AuthenticateContext(
			context.Background(),
			c.Get("Authorization"),
			c.Get(opts.TenantHeader),
			opts.OverridePolicy,
		)
		if err != nil {
			status, body := opts.ErrorHandler(err)
			return c.Status(status).JSON(fiber.Map{"error": body})
		}

		c.Locals("auth_context", ctx)
		c.Locals("user_id", claims.UserID)
		c.Locals("tenant_id", claims.TenantID)
		c.Locals("roles", claims.Roles)
		c.Locals("permissions", claims.Permissions)
		return c.Next()
	}
}

func Context(c *fiber.Ctx) context.Context {
	ctx, ok := c.Locals("auth_context").(context.Context)
	if !ok || ctx == nil {
		return context.Background()
	}
	return ctx
}

func RequirePermission(permission string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if err := auth.RequirePermission(Context(c), permission); err != nil {
			status, body := adapters.DefaultErrorHandler(err)
			return c.Status(status).JSON(fiber.Map{"error": body})
		}
		return c.Next()
	}
}

func RequirePermissionWithService(service *auth.Service, permission string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if err := service.AuthorizePermission(Context(c), permission); err != nil {
			status, body := adapters.DefaultErrorHandler(err)
			return c.Status(status).JSON(fiber.Map{"error": body})
		}
		return c.Next()
	}
}
