package query

import (
	"testing"
)

func TestFilters(t *testing.T) {
	var f Filters
	f.Eq("status", "active")
	f.Like("name", "%foo%")
	f.In("id", []string{"a", "b"})
	if len(f) != 3 {
		t.Fatalf("expected 3 filters, got %d", len(f))
	}
}

func TestSQLBuilderWhere(t *testing.T) {
	b := NewSQLBuilder(PlaceholderQ)
	f := Filters{}
	f.Eq("status", "active")
	f.Like("name", "%x%")
	where, args := b.Where(f)
	if where == "" {
		t.Fatal("expected non-empty where")
	}
	if len(args) != 2 {
		t.Fatalf("expected 2 args, got %d", len(args))
	}
}
