package nethttpadapter

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"

	"github.com/parevo/core/auth"
)

const requestIDHeader = "X-Request-Id"

func RequestIDMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get(requestIDHeader)
			if requestID == "" {
				requestID = generateRequestID()
			}
			w.Header().Set(requestIDHeader, requestID)
			ctx := auth.WithRequestID(r.Context(), requestID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func generateRequestID() string {
	var b [12]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "req-fallback"
	}
	return hex.EncodeToString(b[:])
}
