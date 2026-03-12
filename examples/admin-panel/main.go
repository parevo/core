// Admin panel example.
// Login required. Only users with admin:* permission can access the panel.
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/parevo/core/admin"
	"github.com/parevo/core/auth"
	"github.com/parevo/core/auth/adapters"
	nethttpadapter "github.com/parevo/core/auth/adapters/nethttp"
	"github.com/parevo/core/permission"
	"github.com/parevo/core/storage"
	"github.com/parevo/core/storage/memory"
	"github.com/parevo/core/tenant"
)

const (
	adminUser     = "admin"
	adminPassword = "admin123"
	adminTenant   = "t1"
	tokenCookie   = "parevo_admin_token"
)

// redirectToLoginOnAuthFail redirects to /login when inner handler would return 401/403.
// We can't intercept the response, so we use a wrapper that checks auth first and redirects.
func redirectToLoginOnAuthFail(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// If no token at all, redirect immediately
		if r.Header.Get("Authorization") == "" {
			if c, _ := r.Cookie(tokenCookie); c == nil || c.Value == "" {
				http.Redirect(w, r, "/login", http.StatusFound)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	tenantLifecycle := &memory.TenantLifecycleStore{
		Tenants: map[string]storage.TenantStatus{},
	}
	permStore := &memory.PermissionStore{
		Grants: map[string]bool{
			"admin|t1|admin:*": true,
		},
	}
	sessionStore := &memory.SessionStore{
		Revoked:       map[string]bool{},
		SessionToUser: map[string]string{},
		UserSessions:  map[string]map[string]struct{}{},
	}
	tenantStore := &memory.TenantStore{
		SubjectTenants: map[string][]string{"admin": {"t1"}},
	}
	userStore := &memory.UserStore{
		Users: map[string]storage.UserInfo{
			"admin":  {ID: "admin", Email: "admin@example.com", Name: "Admin User"},
			"user-1": {ID: "user-1", Email: "user1@example.com", Name: "Test User 1"},
		},
	}

	svc, _ := auth.NewServiceWithModules(auth.Config{
		Issuer:    "parevo",
		Audience:  "parevo-api",
		SecretKey: []byte("change-this-in-production"),
	}, auth.Modules{
		Tenant:         tenant.NewService(tenantStore),
		Permission:     permission.NewService(permStore),
		TenantOverride: tenant.StaticOverridePolicy{Allow: true},
	})

	// Middleware: use cookie as Bearer token when Authorization header is empty
	cookieAuth := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("Authorization") == "" {
				if c, err := r.Cookie(tokenCookie); err == nil && c.Value != "" {
					r.Header.Set("Authorization", "Bearer "+c.Value)
				}
			}
			next.ServeHTTP(w, r)
		})
	}

	adminHandler := cookieAuth(
		redirectToLoginOnAuthFail(
			nethttpadapter.AuthMiddleware(svc, adapters.Options{})(
			nethttpadapter.RequirePermissionWithService(svc, "admin:*")(
				admin.Handler(admin.Options{
					AuthService:        svc,
					TenantLifecycle:    tenantLifecycle,
					TenantList:         tenantLifecycle,
					PermissionGrant:    permStore,
					SessionStore:       sessionStore,
					SessionList:        sessionStore,
					UserList:           userStore,
					RequiredPermission: "admin:*",
					BasePath:           "/admin",
					LogoutPath:         "/logout",
				}),
		),
		),
	),
	)

	mux := http.NewServeMux()
	mux.Handle("/admin/", adminHandler)
	mux.Handle("/admin", adminHandler)

	// Login: GET form, POST credentials
	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			_, _ = w.Write([]byte(loginHTML))
			return
		}
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var req struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
			return
		}
		if strings.TrimSpace(req.Username) != adminUser || req.Password != adminPassword {
			http.Error(w, `{"error":"invalid credentials"}`, http.StatusUnauthorized)
			return
		}
		token, err := svc.IssueAccessToken(auth.Claims{
			UserID:   adminUser,
			TenantID: adminTenant,
			Roles:    []string{"admin"},
		})
		if err != nil {
			http.Error(w, `{"error":"token issue failed"}`, http.StatusInternalServerError)
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:     tokenCookie,
			Value:    token,
			Path:     "/",
			MaxAge:   3600,
			HttpOnly: true,
			Secure:   false,
			SameSite: http.SameSiteLaxMode,
		})
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"redirect":"/admin"}`))
	})
	mux.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:     tokenCookie,
			Value:    "",
			Path:     "/",
			MaxAge:   -1,
			HttpOnly: true,
			Expires:  time.Unix(0, 0),
		})
		http.Redirect(w, r, "/login", http.StatusFound)
	})

	// Redirect / to /login
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		http.Redirect(w, r, "/login", http.StatusFound)
	})

	fmt.Println("Admin panel: http://localhost:8086/admin")
	fmt.Println("Login: admin / admin123 (demo only)")
	_ = http.ListenAndServe(":8086", mux)
}

const loginHTML = `<!DOCTYPE html>
<html lang="en" class="dark">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>Admin Login</title>
  <script src="https://cdn.tailwindcss.com"></script>
  <script>tailwind.config = { darkMode: 'class' }</script>
</head>
<body class="min-h-screen bg-slate-900 text-slate-100 flex items-center justify-center">
  <div class="bg-slate-800 p-8 rounded-xl w-full max-w-sm shadow-xl">
    <h1 class="text-xl font-semibold mb-6">Admin Login</h1>
    <form id="form" class="space-y-4">
      <input type="text" name="username" placeholder="Username" required autocomplete="username" class="w-full px-4 py-3 rounded-lg bg-slate-700 border border-slate-600 text-slate-100 placeholder-slate-500 focus:ring-2 focus:ring-sky-500 focus:border-transparent">
      <input type="password" name="password" placeholder="Password" required autocomplete="current-password" class="w-full px-4 py-3 rounded-lg bg-slate-700 border border-slate-600 text-slate-100 placeholder-slate-500 focus:ring-2 focus:ring-sky-500 focus:border-transparent">
      <button type="submit" class="w-full py-3 rounded-lg bg-sky-500 hover:bg-sky-400 text-slate-900 font-semibold transition">Sign in</button>
    </form>
    <div id="err" class="mt-3 text-sm text-red-400"></div>
  </div>
  <script>
    document.getElementById('form').onsubmit = async (e) => {
      e.preventDefault();
      const fd = new FormData(e.target);
      const r = await fetch('/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ username: fd.get('username'), password: fd.get('password') })
      });
      const j = await r.json();
      if (r.ok && j.redirect) location.href = j.redirect;
      else document.getElementById('err').textContent = j.error || 'Login failed';
    };
  </script>
</body>
</html>`
