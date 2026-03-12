package query

// Cursor holds cursor-based pagination params.
type Cursor struct {
	After  string // cursor after which to fetch
	Before string // cursor before which to fetch (for prev page)
	Limit  int
}

// DefaultCursor returns Cursor with limit 20.
func DefaultCursor() Cursor {
	return Cursor{Limit: 20}
}

// WithLimit sets limit (max 1000, min 1).
func (c Cursor) WithLimit(n int) Cursor {
	if n <= 0 {
		n = 20
	}
	if n > 1000 {
		n = 1000
	}
	c.Limit = n
	return c
}

// WithAfter sets after cursor.
func (c Cursor) WithAfter(after string) Cursor {
	c.After = after
	return c
}

// WithBefore sets before cursor.
func (c Cursor) WithBefore(before string) Cursor {
	c.Before = before
	return c
}

// CursorResult holds cursor-paginated result.
type CursorResult struct {
	Items    any
	Next     string // cursor for next page
	Previous string // cursor for prev page
	HasNext  bool
	HasPrev  bool
}
