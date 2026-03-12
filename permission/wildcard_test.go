package permission

import "testing"

func TestMatchPermission(t *testing.T) {
	tests := []struct {
		granted, requested string
		want               bool
	}{
		{"orders:read", "orders:read", true},
		{"orders:*", "orders:read", true},
		{"*:read", "orders:read", true},
		{"*:*", "orders:read", true},
		{"*", "orders:read", true},
		{"orders:write", "orders:read", false},
	}
	for _, tt := range tests {
		if got := MatchPermission(tt.granted, tt.requested); got != tt.want {
			t.Errorf("MatchPermission(%q, %q) = %v, want %v", tt.granted, tt.requested, got, tt.want)
		}
	}
}
