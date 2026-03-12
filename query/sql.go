package query

import (
	"fmt"
	"strings"
)

// ErrFieldNotAllowed is returned when a filter/sort field is not in the whitelist.
var ErrFieldNotAllowed = fmt.Errorf("query: field not allowed")

// SQLBuilder helps build SQL WHERE and ORDER BY from Filters and SortBy.
// Use with parameterized queries; placeholder ? for MySQL, $1,$2 for Postgres.
// Use WithWhitelist to prevent field injection when building from user input.
type SQLBuilder struct {
	placeholder func(int) string
	whitelist   map[string]struct{}
	args        []any
	argIndex    int
}

// NewSQLBuilder creates builder. Use PlaceholderQ for MySQL, PlaceholderDollar for Postgres.
func NewSQLBuilder(placeholder func(int) string) *SQLBuilder {
	if placeholder == nil {
		placeholder = PlaceholderQ
	}
	return &SQLBuilder{placeholder: placeholder}
}

// WithWhitelist restricts Filter.Field and Sort.Field to allowed columns. Prevents SQL injection from user-supplied field names.
func (b *SQLBuilder) WithWhitelist(fields []string) *SQLBuilder {
	b.whitelist = make(map[string]struct{})
	for _, f := range fields {
		b.whitelist[f] = struct{}{}
	}
	return b
}

func (b *SQLBuilder) isFieldAllowed(field string) bool {
	if b.whitelist == nil {
		return true
	}
	_, ok := b.whitelist[field]
	return ok
}

// PlaceholderDollar returns $1, $2, ... for Postgres.
func PlaceholderDollar(i int) string {
	return fmt.Sprintf("$%d", i)
}

// PlaceholderQ returns ? for MySQL.
func PlaceholderQ(int) string {
	return "?"
}

// Where builds WHERE clause from Filters. Returns "WHERE ..." or "" if no filters.
// If WithWhitelist was set, filters with disallowed fields are skipped.
func (b *SQLBuilder) Where(f Filters) (string, []any) {
	where, _, _ := b.WhereSafe(f)
	return where, b.args
}

// WhereSafe builds WHERE clause and returns error if any field is not in whitelist (when set).
func (b *SQLBuilder) WhereSafe(f Filters) (string, []any, error) {
	if len(f) == 0 {
		return "", nil, nil
	}
	b.args = nil
	b.argIndex = 0
	var parts []string
	for _, ff := range f {
		if b.whitelist != nil && !b.isFieldAllowed(ff.Field) {
			return "", nil, fmt.Errorf("%w: %q", ErrFieldNotAllowed, ff.Field)
		}
		part := b.buildCondition(ff)
		if part != "" {
			parts = append(parts, part)
		}
	}
	if len(parts) == 0 {
		return "", nil, nil
	}
	return " WHERE " + strings.Join(parts, " AND "), b.args, nil
}

func (b *SQLBuilder) buildCondition(f Filter) string {
	b.argIndex++
	ph := b.placeholder(b.argIndex)
	switch f.Op {
	case OpEq:
		if f.Value == nil {
			return f.Field + " IS NULL"
		}
		b.args = append(b.args, f.Value)
		return f.Field + " = " + ph
	case OpNe:
		if f.Value == nil {
			return f.Field + " IS NOT NULL"
		}
		b.args = append(b.args, f.Value)
		return f.Field + " != " + ph
	case OpGt:
		b.args = append(b.args, f.Value)
		return f.Field + " > " + ph
	case OpGte:
		b.args = append(b.args, f.Value)
		return f.Field + " >= " + ph
	case OpLt:
		b.args = append(b.args, f.Value)
		return f.Field + " < " + ph
	case OpLte:
		b.args = append(b.args, f.Value)
		return f.Field + " <= " + ph
	case OpLike:
		b.args = append(b.args, f.Value)
		return f.Field + " LIKE " + ph
	case OpIn:
		// Expand IN (?, ?, ?)
		switch v := f.Value.(type) {
		case []string:
			if len(v) == 0 {
				return "1=0"
			}
			placeholders := make([]string, len(v))
			for i := range v {
				b.argIndex++
				placeholders[i] = b.placeholder(b.argIndex)
				b.args = append(b.args, v[i])
			}
			return f.Field + " IN (" + strings.Join(placeholders, ", ") + ")"
		case []any:
			if len(v) == 0 {
				return "1=0"
			}
			placeholders := make([]string, len(v))
			for i := range v {
				b.argIndex++
				placeholders[i] = b.placeholder(b.argIndex)
				b.args = append(b.args, v[i])
			}
			return f.Field + " IN (" + strings.Join(placeholders, ", ") + ")"
		default:
			return ""
		}
	case OpNin:
		switch v := f.Value.(type) {
		case []string:
			if len(v) == 0 {
				return "1=1"
			}
			placeholders := make([]string, len(v))
			for i := range v {
				b.argIndex++
				placeholders[i] = b.placeholder(b.argIndex)
				b.args = append(b.args, v[i])
			}
			return f.Field + " NOT IN (" + strings.Join(placeholders, ", ") + ")"
		case []any:
			if len(v) == 0 {
				return "1=1"
			}
			placeholders := make([]string, len(v))
			for i := range v {
				b.argIndex++
				placeholders[i] = b.placeholder(b.argIndex)
				b.args = append(b.args, v[i])
			}
			return f.Field + " NOT IN (" + strings.Join(placeholders, ", ") + ")"
		default:
			return ""
		}
	case OpNull:
		return f.Field + " IS NULL"
	case OpNotNull:
		return f.Field + " IS NOT NULL"
	default:
		return ""
	}
}

// OrderBy builds ORDER BY clause from SortBy. Returns "ORDER BY ..." or "".
// If WithWhitelist was set, sorts with disallowed fields are skipped.
func (b *SQLBuilder) OrderBy(s SortBy) string {
	order, _ := b.OrderBySafe(s)
	return order
}

// OrderBySafe builds ORDER BY and returns error if any field is not in whitelist (when set).
func (b *SQLBuilder) OrderBySafe(s SortBy) (string, error) {
	if len(s) == 0 {
		return "", nil
	}
	var parts []string
	for _, ss := range s {
		if b.whitelist != nil && !b.isFieldAllowed(ss.Field) {
			return "", fmt.Errorf("%w: %q", ErrFieldNotAllowed, ss.Field)
		}
		dir := " ASC"
		if ss.Desc {
			dir = " DESC"
		}
		parts = append(parts, ss.Field+dir)
	}
	if len(parts) == 0 {
		return "", nil
	}
	return " ORDER BY " + strings.Join(parts, ", "), nil
}
