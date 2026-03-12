package ratelimit

import (
	"testing"
	"time"
)

func TestTenantLimiter(t *testing.T) {
	limiter := NewTenantLimiter(2, 100*time.Millisecond)

	if !limiter.Allow("t1", "ip1") {
		t.Fatal("first request should allow")
	}
	if !limiter.Allow("t1", "ip1") {
		t.Fatal("second request should allow")
	}
	if limiter.Allow("t1", "ip1") {
		t.Fatal("third request should deny (limit 2)")
	}

	if !limiter.Allow("t2", "ip1") {
		t.Fatal("different tenant should allow")
	}
}

func TestTenantLimiterOverride(t *testing.T) {
	limiter := NewTenantLimiter(2, 100*time.Millisecond)
	limiter.SetTenantLimit("vip", 5, 100*time.Millisecond)

	for i := 0; i < 5; i++ {
		if !limiter.Allow("vip", "ip1") {
			t.Fatalf("vip request %d should allow", i+1)
		}
	}
	if limiter.Allow("vip", "ip1") {
		t.Fatal("6th vip request should deny")
	}
}
