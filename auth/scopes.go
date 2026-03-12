package auth

import "strings"

// HasScope returns true if claims include the required scope.
// Supports wildcards: "read:*" matches "read:orders".
func HasScope(claims *Claims, required string) bool {
	if claims == nil {
		return false
	}
	for _, s := range claims.Scopes {
		if matchScope(s, required) {
			return true
		}
	}
	return false
}

// RequireScope checks all required scopes are present.
func RequireScope(claims *Claims, required ...string) bool {
	for _, r := range required {
		if !HasScope(claims, r) {
			return false
		}
	}
	return true
}

func matchScope(granted, required string) bool {
	if granted == required || granted == "*" || granted == "*:*" {
		return true
	}
	gParts := strings.SplitN(granted, ":", 2)
	rParts := strings.SplitN(required, ":", 2)
	if len(gParts) != 2 || len(rParts) != 2 {
		return false
	}
	if gParts[0] == "*" && gParts[1] == "*" {
		return true
	}
	if gParts[0] == rParts[0] && (gParts[1] == "*" || gParts[1] == rParts[1]) {
		return true
	}
	return false
}
