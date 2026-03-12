package nethttpadapter

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/parevo/core/auth/ratelimit"
)

func TestBruteForceMiddlewareBlocksLockedKey(t *testing.T) {
	lockout := ratelimit.NewLockoutManager(1, time.Minute, time.Minute)
	lockout.RegisterFailure("127.0.0.1")
	lockout.RegisterFailure("127.0.0.1")

	handler := BruteForceMiddleware(lockout, BruteForceConfig{})(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodPost, "/login", nil)
	req.RemoteAddr = "127.0.0.1:1234"
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusTooManyRequests {
		t.Fatalf("expected 429, got %d", rec.Code)
	}
}
