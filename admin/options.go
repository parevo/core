package admin

import (
	"net/http"

	"github.com/parevo/core/auth"
	"github.com/parevo/core/storage"
)

// Options configures the admin panel.
type Options struct {
	AuthService *auth.Service
	// Required for tenant management
	TenantLifecycle storage.TenantLifecycleStore
	TenantList      storage.TenantListStore
	// Required for permission management
	PermissionGrant storage.PermissionGrantStore
	// Required for session management
	SessionStore   storage.SessionStore
	SessionList    storage.SessionListStore
	// Optional: list registered users
	UserList storage.UserListStore
	// Permission required to access admin (e.g. "admin:*" or "realm:admin")
	RequiredPermission string
	// Base path for admin routes (default: /admin)
	BasePath string
	// Logout path for UI (e.g. /logout). If empty, no logout link shown.
	LogoutPath string
}

func (o *Options) basePath() string {
	if o.BasePath != "" {
		return o.BasePath
	}
	return "/admin"
}

// Handler returns an http.Handler for the admin panel.
// Mount it behind auth middleware. Example:
//
//	mux.Handle("/admin/", authMiddleware(admin.Handler(opts)))
func Handler(opts Options) http.Handler {
	return newAdminHandler(opts)
}
