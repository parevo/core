package nethttpadapter

import (
	"net/http"

	"github.com/parevo/core/auth/ratelimit"
)

type BruteForceConfig struct {
	ExtractKey func(r *http.Request) string
}

func BruteForceMiddleware(lockout *ratelimit.LockoutManager, cfg BruteForceConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if lockout == nil {
				next.ServeHTTP(w, r)
				return
			}
			key := clientIP(r)
			if cfg.ExtractKey != nil {
				if customKey := cfg.ExtractKey(r); customKey != "" {
					key = customKey
				}
			}
			if lockout.IsLocked(key) {
				http.Error(w, "account temporarily locked", http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func BruteForceKey(r *http.Request, cfg BruteForceConfig) string {
	if cfg.ExtractKey != nil {
		if customKey := cfg.ExtractKey(r); customKey != "" {
			return customKey
		}
	}
	return clientIP(r)
}
