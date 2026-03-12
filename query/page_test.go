package query

import (
	"testing"
)

func TestPage(t *testing.T) {
	p := DefaultPage().WithLimit(10).WithOffset(20)
	limit, offset := p.SQL()
	if limit != 10 || offset != 20 {
		t.Fatalf("expected limit=10 offset=20, got %d %d", limit, offset)
	}
	if p.PageNum() != 3 {
		t.Fatalf("expected page 3, got %d", p.PageNum())
	}
}

func TestPageResult(t *testing.T) {
	items := []string{"a", "b"}
	page := Page{Offset: 0, Limit: 1}
	res := NewPageResult(items, 2, page)
	if !res.HasNext {
		t.Fatal("expected HasNext true when offset+limit < total")
	}
	res2 := NewPageResult(items, 2, Page{Offset: 0, Limit: 2})
	if res2.HasNext {
		t.Fatal("expected HasNext false when offset+limit >= total")
	}
}
