// Package query provides pagination, filtering, and sorting for list operations.
//
// # Pagination
//
// Offset/limit:
//
//	page := query.DefaultPage().WithLimit(50).WithOffset(100)
//	limit, offset := page.SQL()
//
// Cursor-based:
//
//	cursor := query.DefaultCursor().WithAfter("last_id").WithLimit(20)
//
// # Filtering
//
//	 filters := query.Filters{}
//	 filters.Eq("status", "active")
//	 filters.Like("name", "%foo%")
//	 filters.In("tenant_id", []string{"t1", "t2"})
//
// # SQL building
//
//	b := query.NewSQLBuilder(query.PlaceholderQ)
//	where, args := b.Where(filters)
//	orderBy := b.OrderBy(sortBy)
package query