// API key authentication example.
// Tokens with pk_ or sk_ prefix are validated as API keys instead of JWT.
package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"sync"

	"github.com/parevo/core/auth"
	"github.com/parevo/core/auth/adapters"
	"github.com/parevo/core/auth/apikey"
	nethttpadapter "github.com/parevo/core/auth/adapters/nethttp"
)

type memAPIKeyStore struct {
	mu   sync.RWMutex
	keys map[string]struct{ userID, tenantID string }
}

func (s *memAPIKeyStore) Validate(_ context.Context, keyHash, _ string) (userID, tenantID string, err error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if e, ok := s.keys[keyHash]; ok {
		return e.userID, e.tenantID, nil
	}
	return "", "", apikey.ErrKeyNotFound
}

func (s *memAPIKeyStore) Add(rawKey, userID, tenantID string) {
	hash := sha256.Sum256([]byte(rawKey))
	keyHash := hex.EncodeToString(hash[:])
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.keys == nil {
		s.keys = make(map[string]struct{ userID, tenantID string })
	}
	s.keys[keyHash] = struct{ userID, tenantID string }{userID, tenantID}
}

func main() {
	apiKeyStore := &memAPIKeyStore{}
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
