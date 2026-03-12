package query

import (
	"strings"
	"testing"
)

func TestSQLBuilderWithWhitelist(t *testing.T) {
	b := NewSQLBuilder(PlaceholderQ).WithWhitelist([]string{"status", "name"})
	f := Filters{}
	f.Eq("status", "active")
	f.Eq("evil; DROP TABLE users", "x") // should fail
	_, _, err := b.WhereSafe(f)
	if err == nil {
		t.Fatal("expected error for disallowed field")
	}
	if !strings.Contains(err.Error(), "evil") {
		t.Fatalf("expected field name in error, got %v", err)
	}
}

func TestSQLBuilderWhitelistAllowed(t *testing.T) {
	b := NewSQLBuilder(PlaceholderQ).WithWhitelist([]string{"status", "name"})
	f := Filters{}
	f.Eq("status", "active")
	f.Like("name", "%x%")
	where, args, err := b.WhereSafe(f)
	if err != nil {
		t.Fatal(err)
	}
	if where == "" || len(args) != 2 {
		t.Fatalf("expected where and 2 args, got %q %d", where, len(args))
	}
}
