// Package query provides pagination, filtering, and sorting helpers for list operations.
package query

// Page holds offset/limit pagination params.
type Page struct {
	Offset int
	Limit  int
}

// DefaultPage returns Page with limit 20, offset 0.
func DefaultPage() Page {
	return Page{Limit: 20}
}

// WithLimit sets limit (max 1000, min 1).
func (p Page) WithLimit(n int) Page {
	if n <= 0 {
		n = 20
	}
	if n > 1000 {
		n = 1000
	}
	p.Limit = n
	return p
}

// WithOffset sets offset.
func (p Page) WithOffset(n int) Page {
	if n < 0 {
		n = 0
	}
	p.Offset = n
	return p
}

// PageNum returns 1-based page number from offset/limit.
func (p Page) PageNum() int {
	if p.Limit <= 0 {
		return 1
	}
	return (p.Offset / p.Limit) + 1
}

// SQL returns "LIMIT ? OFFSET ?" args.
func (p Page) SQL() (limit, offset int) {
	if p.Limit <= 0 {
		p.Limit = 20
	}
	if p.Offset < 0 {
		p.Offset = 0
	}
	return p.Limit, p.Offset
}

// PageResult holds paginated result metadata.
type PageResult struct {
	Items   any   // slice of items
	Total   int64 // total count (optional, -1 = unknown)
	Offset  int
	Limit   int
	HasNext bool
}

// NewPageResult builds PageResult from items, total, and page.
func NewPageResult(items any, total int64, page Page) PageResult {
	limit, offset := page.SQL()
	hasNext := total >= 0 && int64(offset+limit) < total
	return PageResult{
		Items:   items,
		Total:   total,
		Offset:  offset,
		Limit:   limit,
		HasNext: hasNext,
	}
}
