package ratelimit

import (
	"testing"
	"time"
)

func TestLockoutManager(t *testing.T) {
	m := NewLockoutManager(1, time.Minute, time.Minute)
	key := "u@example.com"
	if m.IsLocked(key) {
		t.Fatalf("should not be locked initially")
	}
	m.RegisterFailure(key)
	if m.IsLocked(key) {
		t.Fatalf("first failure should not lock")
	}
	m.RegisterFailure(key)
	if !m.IsLocked(key) {
		t.Fatalf("second failure should lock")
	}
	m.RegisterSuccess(key)
	if m.IsLocked(key) {
		t.Fatalf("success should unlock")
	}
}
