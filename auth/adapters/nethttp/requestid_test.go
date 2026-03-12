package nethttpadapter

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/parevo/core/auth"
)

func TestRequestIDMiddleware(t *testing.T) {
	handler := RequestIDMiddleware()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if requestID, ok := auth.RequestIDFromContext(r.Context()); !ok || requestID == "" {
			t.Fatalf("request id not found in context")
		}
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if rec.Header().Get("X-Request-Id") == "" {
		t.Fatalf("expected X-Request-Id header")
	}
}
