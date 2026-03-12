package query

import (
	"fmt"
	"strings"
)

// SQLBuilder helps build SQL WHERE and ORDER BY from Filters and SortBy.
// Use with parameterized queries; placeholder ? for MySQL, $1,$2 for Postgres.
type SQLBuilder struct {
	placeholder func(int) string
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

// PlaceholderDollar returns $1, $2, ... for Postgres.
func PlaceholderDollar(i int) string {
	return fmt.Sprintf("$%d", i)
}

// PlaceholderQ returns ? for MySQL.
func PlaceholderQ(int) string {
	return "?"
}

// Where builds WHERE clause from Filters. Returns "WHERE ..." or "" if no filters.
func (b *SQLBuilder) Where(f Filters) (string, []any) {
	if len(f) == 0 {
		return "", nil
	}
	b.args = nil
	b.argIndex = 0
	var parts []string
	for _, ff := range f {
		part := b.buildCondition(ff)
		if part != "" {
			parts = append(parts, part)
		}
	}
	if len(parts) == 0 {
		return "", nil
	}
	return " WHERE " + strings.Join(parts, " AND "), b.args
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
func (b *SQLBuilder) OrderBy(s SortBy) string {
	if len(s) == 0 {
		return ""
	}
	var parts []string
	for _, ss := range s {
		dir := " ASC"
		if ss.Desc {
			dir = " DESC"
		}
		parts = append(parts, ss.Field+dir)
	}
	return " ORDER BY " + strings.Join(parts, ", ")
}
