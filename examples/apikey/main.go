// API key authentication example.
// Tokens with pk_ or sk_ prefix are validated as API keys instead of JWT.
package main

import (
	"fmt"
	"net/http"

	"github.com/parevo/core/auth"
	"github.com/parevo/core/auth/adapters"
	"github.com/parevo/core/auth/apikey"
	apikeymemory "github.com/parevo/core/auth/apikey/memory"
	nethttpadapter "github.com/parevo/core/auth/adapters/nethttp"
)

func main() {
	apiKeyStore := apikeymemory.NewStore()
	key, err := apikey.GenerateKey("pk_")
	if err != nil {
		panic(err)
	}
	apiKeyStore.Add(key, "user-api", "tenant-1")
	fmt.Printf("Sample API key (demo only): %s\n", key)

	svc, err := auth.NewServiceWithModules(auth.Config{
		Issuer:    "parevo",
		Audience:  "parevo-api",
		SecretKey: []byte("change-this-in-production"),
	}, auth.Modules{
		APIKey: apikey.NewService(apiKeyStore),
	})
	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()
	protected := nethttpadapter.AuthMiddleware(svc, adapters.Options{})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, _ := auth.ClaimsFromContext(r.Context())
		_, _ = fmt.Fprintf(w, "hello user=%s tenant=%s", claims.UserID, claims.TenantID)
	}))
	mux.Handle("/secure", protected)

	fmt.Println("Run: curl -H 'Authorization: Bearer <API_KEY>' http://localhost:8082/secure")
	_ = http.ListenAndServe(":8082", mux)
}
