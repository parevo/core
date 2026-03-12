// TOTP 2FA/MFA example.
// Secret generation and code verification with pquerna/otp.
package main

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/parevo/core/auth/mfa"
)

type memTOTPStore struct {
	mu      sync.RWMutex
	secrets map[string]string
	enabled map[string]bool
}

func (s *memTOTPStore) GetSecret(_ context.Context, userID string) (string, bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	secret := s.secrets[userID]
	return secret, s.enabled[userID], nil
}

func (s *memTOTPStore) SetSecret(_ context.Context, userID, secret string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.secrets == nil {
		s.secrets = make(map[string]string)
	}
	s.secrets[userID] = secret
	return nil
}

func (s *memTOTPStore) Enable(_ context.Context, userID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.enabled == nil {
		s.enabled = make(map[string]bool)
	}
	s.enabled[userID] = true
	return nil
}

func (s *memTOTPStore) Disable(_ context.Context, userID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.enabled[userID] = false
	return nil
}

func main() {
	store := &memTOTPStore{}
	totpSvc := mfa.NewTOTPService(store, mfa.NewPquernaVerifier())

	mux := http.NewServeMux()

	// Setup: generate secret + QR URL
	mux.HandleFunc("POST /mfa/setup", func(w http.ResponseWriter, r *http.Request) {
		userID := r.URL.Query().Get("user")
		if userID == "" {
			userID = "demo-user"
		}
		secret, qrURL, err := totpSvc.Setup(r.Context(), userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		_ = totpSvc.Enable(r.Context(), userID)
		_, _ = fmt.Fprintf(w, "Secret: %s\nQR URL: %s\n", secret, qrURL)
	})

	// Verify: validate 6-digit code
	mux.HandleFunc("POST /mfa/verify", func(w http.ResponseWriter, r *http.Request) {
		userID := r.URL.Query().Get("user")
		code := r.URL.Query().Get("code")
		if userID == "" {
			userID = "demo-user"
		}
		if code == "" {
			http.Error(w, "code required", http.StatusBadRequest)
			return
		}
		if err := totpSvc.Verify(r.Context(), userID, code); err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		_, _ = w.Write([]byte("TOTP verified\n"))
	})

	fmt.Println("TOTP MFA demo:")
	fmt.Println("  1. POST /mfa/setup?user=demo-user  -> get secret")
	fmt.Println("  2. Add secret to Google Authenticator")
	fmt.Println("  3. POST /mfa/verify?user=demo-user&code=123456")
	fmt.Println("Listen :8085")
	_ = http.ListenAndServe(":8085", mux)
}
