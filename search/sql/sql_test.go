package sql

import (
	"testing"
)

func TestFullTextWhere(t *testing.T) {
	ft := NewFullText("MATCH(name,description) AGAINST(? IN NATURAL LANGUAGE MODE)")

	where, args := ft.Where("foo")
	if where != "MATCH(name,description) AGAINST(? IN NATURAL LANGUAGE MODE)" {
		t.Errorf("unexpected where: %s", where)
	}
	if len(args) != 1 || args[0] != "foo" {
		t.Errorf("unexpected args: %v", args)
	}

	where2, args2 := ft.Where("")
	if where2 != "1=1" {
		t.Errorf("empty query should return 1=1, got %s", where2)
	}
	if len(args2) != 0 {
		t.Errorf("empty query should have no args, got %v", args2)
	}
}
