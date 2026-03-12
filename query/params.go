package query

// Params combines Page, Filters, and SortBy for list queries.
type Params struct {
	Page   Page
	Filter Filters
	Sort   SortBy
}

// NewParams returns Params with defaults.
func NewParams() Params {
	return Params{
		Page:   DefaultPage(),
		Filter: Filters{},
		Sort:   SortBy{},
	}
}

// WithPage sets pagination.
func (p Params) WithPage(page Page) Params {
	p.Page = page
	return p
}

// WithFilter sets filters.
func (p Params) WithFilter(f Filters) Params {
	p.Filter = f
	return p
}

// WithSort sets sort order.
func (p Params) WithSort(s SortBy) Params {
	p.Sort = s
	return p
}
