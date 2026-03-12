// Blacklist logout example.
// Logout adds the token jti to blacklist; subsequent requests with that token fail.
package main

import (
	"fmt"
	"net/http"

	"strings"

	"github.com/parevo/core/auth"
	"github.com/parevo/core/auth/adapters"
	"github.com/parevo/core/auth/blacklist"
	nethttpadapter "github.com/parevo/core/auth/adapters/nethttp"
	"github.com/parevo/core/storage/memory"
)

func extractToken(h string) string {
	parts := strings.Fields(strings.TrimSpace(h))
	if len(parts) >= 2 && strings.EqualFold(parts[0], "Bearer") {
		return parts[1]
	}
	if len(parts) == 1 {
		return parts[0]
	}
	return ""
}

func main() {
	blacklistStore := memory.NewBlacklistStore()
	blacklistSvc := blacklist.NewService(blacklistStore)

	sessionStore := &memory.SessionStore{}
	svc, err := auth.NewServiceWithModules(auth.Config{
		Issuer:    "parevo",
		Audience:  "parevo-api",
		SecretKey: []byte("change-this-in-production"),
	}, auth.Modules{
		SessionStore: sessionStore,
		Blacklist:    blacklistSvc,
	})
	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()

	// Login: issue token
	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		pair, err := svc.IssueTokenPair(r.Context(), auth.Claims{UserID: "u1", TenantID: "t1"})
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		_, _ = fmt.Fprintf(w, `{"access_token":"%s","refresh_token":"%s"}`, pair.AccessToken, pair.RefreshToken)
	})

	// Logout: blacklist the token
	mux.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		token := extractToken(r.Header.Get("Authorization"))
		if token == "" {
			http.Error(w, "missing token", 400)
			return
		}
		claims, err := svc.ParseAndValidate(token)
		if err != nil {
			http.Error(w, err.Error(), 401)
			return
		}
		// Blacklist by jti until token expires
		exp, _ := claims.GetExpirationTime()
		if exp != nil {
			_ = blacklistSvc.Revoke(r.Context(), claims.ID, exp.Time)
		}
		_, _ = w.Write([]byte(`{"ok":true}`))
	})

	// Protected
	protected := nethttpadapter.AuthMiddleware(svc, adapters.Options{})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, _ := auth.ClaimsFromContext(r.Context())
		_, _ = fmt.Fprintf(w, "hello %s", claims.UserID)
	}))
	mux.Handle("/me", protected)

	fmt.Println("Blacklist logout example: http://localhost:8085")
	fmt.Println("1. POST /login -> get token")
	fmt.Println("2. GET /me Authorization: Bearer <token> -> ok")
	fmt.Println("3. POST /logout Authorization: Bearer <token>")
	fmt.Println("4. GET /me Authorization: Bearer <token> -> 401 (blacklisted)")
	_ = http.ListenAndServe(":8085", mux)
}
