package nethttpadapter

import (
	"net/http"

	"github.com/parevo/core/auth"
	"github.com/parevo/core/auth/adapters"
)

func AuthMiddleware(service *auth.Service, opts adapters.Options) func(http.Handler) http.Handler {
	opts = opts.WithDefaults()

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, _, err := service.AuthenticateContext(
				r.Context(),
				r.Header.Get("Authorization"),
				r.Header.Get(opts.TenantHeader),
				opts.OverridePolicy,
			)
			if err != nil {
				status, body := opts.ErrorHandler(err)
				http.Error(w, body, status)
				return
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RequirePermission(permission string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if err := auth.RequirePermission(r.Context(), permission); err != nil {
				status, body := adapters.DefaultErrorHandler(err)
				http.Error(w, body, status)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func RequirePermissionWithService(service *auth.Service, permission string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if err := service.AuthorizePermission(r.Context(), permission); err != nil {
				status, body := adapters.DefaultErrorHandler(err)
				http.Error(w, body, status)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
