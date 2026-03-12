package chiadapter

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/parevo/core/auth"
	"github.com/parevo/core/auth/adapters"
)

func TestAuthMiddleware_Unauthorized(t *testing.T) {
	svc, err := auth.NewService(auth.Config{
		Issuer:    "parevo",
		Audience:  "parevo-api",
		SecretKey: []byte("secret"),
	})
	if err != nil {
		t.Fatalf("new service failed: %v", err)
	}

	handler := AuthMiddleware(svc, adapters.Options{})(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/secure", nil)
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rec.Code)
	}
}
