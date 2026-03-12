package mfa

import (
	"testing"
	"time"

	"github.com/pquerna/otp/totp"
)

func TestPquernaVerifier(t *testing.T) {
	v := NewPquernaVerifier()
	secret, err := v.GenerateSecret()
	if err != nil {
		t.Fatalf("GenerateSecret: %v", err)
	}
	if secret == "" {
		t.Fatal("expected non-empty secret")
	}

	code, err := v.GenerateCode(secret, time.Now())
	if err != nil {
		t.Fatalf("GenerateCode: %v", err)
	}
	if !v.Verify(secret, code) {
		t.Fatal("Verify failed for valid code")
	}

	if v.Verify(secret, "000000") {
		t.Fatal("Verify should fail for invalid code")
	}

	if v.Verify(secret, "12345") {
		t.Fatal("Verify should fail for wrong-length code")
	}
}

func TestPquernaCompatibleWithTotp(t *testing.T) {
	v := NewPquernaVerifier()
	secret, _ := v.GenerateSecret()
	code, _ := totp.GenerateCode(secret, time.Now())
	if !v.Verify(secret, code) {
		t.Fatal("PquernaVerifier should accept totp.GenerateCode output")
	}
}
