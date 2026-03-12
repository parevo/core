package query

// Op is a filter operator.
type Op string

const (
	OpEq    Op = "eq"
	OpNe    Op = "ne"
	OpGt    Op = "gt"
	OpGte   Op = "gte"
	OpLt    Op = "lt"
	OpLte   Op = "lte"
	OpLike  Op = "like"
	OpIn    Op = "in"
	OpNin   Op = "nin"
	OpNull  Op = "null"  // field is null
	OpNotNull Op = "nnull" // field is not null
)

// Filter is a single filter condition.
type Filter struct {
	Field string
	Op    Op
	Value any
}

// Filters is a collection of Filter (AND logic).
type Filters []Filter

// Add appends a filter.
func (f *Filters) Add(field string, op Op, value any) {
	*f = append(*f, Filter{Field: field, Op: op, Value: value})
}

// Eq adds equality filter.
func (f *Filters) Eq(field string, value any) {
	f.Add(field, OpEq, value)
}

// Ne adds not-equal filter.
func (f *Filters) Ne(field string, value any) {
	f.Add(field, OpNe, value)
}

// Gt adds greater-than filter.
func (f *Filters) Gt(field string, value any) {
	f.Add(field, OpGt, value)
}

// Gte adds greater-or-equal filter.
func (f *Filters) Gte(field string, value any) {
	f.Add(field, OpGte, value)
}

// Lt adds less-than filter.
func (f *Filters) Lt(field string, value any) {
	f.Add(field, OpLt, value)
}

// Lte adds less-or-equal filter.
func (f *Filters) Lte(field string, value any) {
	f.Add(field, OpLte, value)
}

// Like adds LIKE filter (value: string with % wildcards).
func (f *Filters) Like(field string, value string) {
	f.Add(field, OpLike, value)
}

// In adds IN filter (value: []any or []string).
func (f *Filters) In(field string, value any) {
	f.Add(field, OpIn, value)
}

// Nin adds NOT IN filter.
func (f *Filters) Nin(field string, value any) {
	f.Add(field, OpNin, value)
}

// Null adds IS NULL filter.
func (f *Filters) Null(field string) {
	f.Add(field, OpNull, nil)
}

// NotNull adds IS NOT NULL filter.
func (f *Filters) NotNull(field string) {
	f.Add(field, OpNotNull, nil)
}
