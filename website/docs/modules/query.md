# Query Module

Pagination, filter, and sort helpers for list operations.

## Pagination

### Offset/Limit

```go
import "github.com/parevo/core/query"

page := query.DefaultPage().WithLimit(50).WithOffset(100)
limit, offset := page.SQL() // 50, 100
```

### Cursor-based

```go
cursor := query.DefaultCursor().WithAfter("last_id").WithLimit(20)
```

## Filtering

```go
var f query.Filters
f.Eq("status", "active")
f.Like("name", "%foo%")
f.In("tenant_id", []string{"t1", "t2"})
f.Gte("created_at", time.Now().Add(-24*time.Hour))
```

Operators: `eq`, `ne`, `gt`, `gte`, `lt`, `lte`, `like`, `in`, `nin`, `null`, `nnull`.

## Sort

```go
var s query.SortBy
s.Desc("created_at")
s.Asc("name")
```

## SQL building

```go
b := query.NewSQLBuilder(query.PlaceholderQ)  // MySQL
// b := query.NewSQLBuilder(query.PlaceholderDollar)  // Postgres
where, args := b.Where(filters)
orderBy := b.OrderBy(sortBy)
// SELECT * FROM users WHERE status = ? AND name LIKE ? ORDER BY created_at DESC
```

## Combined params

```go
params := query.NewParams().
    WithPage(query.DefaultPage().WithLimit(20)).
    WithFilter(filters).
    WithSort(sortBy)
```
