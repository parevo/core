package query

// Sort holds field + direction.
type Sort struct {
	Field string
	Desc  bool
}

// SortBy is a collection of Sort (order matters).
type SortBy []Sort

// Add appends a sort.
func (s *SortBy) Add(field string, desc bool) {
	*s = append(*s, Sort{Field: field, Desc: desc})
}

// Asc adds ascending sort.
func (s *SortBy) Asc(field string) {
	s.Add(field, false)
}

// Desc adds descending sort.
func (s *SortBy) Desc(field string) {
	s.Add(field, true)
}
