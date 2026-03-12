package main

import (
	"fmt"
	"net/http"

	"github.com/parevo/core/auth"
	"github.com/parevo/core/auth/adapters"
	nethttpadapter "github.com/parevo/core/auth/adapters/nethttp"
)

func main() {
	svc, err := auth.NewService(auth.Config{
		Issuer:    "parevo",
		Audience:  "parevo-api",
		SecretKey: []byte("change-this-in-production"),
	})
	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()
	mux.Handle("/health", http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	}))

	protected := nethttpadapter.AuthMiddleware(svc, adapters.Options{})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, _ := auth.ClaimsFromContext(r.Context())
		_, _ = fmt.Fprintf(w, "hello user=%s tenant=%s", claims.UserID, claims.TenantID)
	}))

	mux.Handle("/secure", protected)
	fmt.Println("listening on :8080")
	_ = http.ListenAndServe(":8080", mux)
}
