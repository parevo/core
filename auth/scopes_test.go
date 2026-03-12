package auth

import (
	"testing"
)

func TestHasScope(t *testing.T) {
	claims := &Claims{Scopes: []string{"read:orders", "write:users"}}

	if !HasScope(claims, "read:orders") {
		t.Error("read:orders should match")
	}
	if HasScope(claims, "read:products") {
		t.Error("read:orders should not match read:products")
	}
	claims2 := &Claims{Scopes: []string{"read:*"}}
	if !HasScope(claims2, "read:products") {
		t.Error("read:* should match read:products")
	}
	if HasScope(claims, "write:orders") {
		t.Error("write:orders should not match")
	}
	if HasScope(nil, "read:orders") {
		t.Error("nil claims should return false")
	}
}

func TestHasScope_Wildcard(t *testing.T) {
	claims := &Claims{Scopes: []string{"*", "read:orders"}}

	if !HasScope(claims, "write:users") {
		t.Error("* should match any scope")
	}
}

func TestRequireScope(t *testing.T) {
	claims := &Claims{Scopes: []string{"read:orders", "write:users"}}

	if !RequireScope(claims, "read:orders", "write:users") {
		t.Error("both scopes present should pass")
	}
	if RequireScope(claims, "read:orders", "admin:all") {
		t.Error("missing admin:all should fail")
	}
}
