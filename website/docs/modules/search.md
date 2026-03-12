# Search Module

Full-text ve yapılandırılmış arama arayüzü.

## SQL Full-Text

SQL FULLTEXT veya Postgres tsvector ile WHERE clause oluşturma:

```go
import (
    "github.com/parevo/core/search/sql"
)

ft := sql.NewFullText("MATCH(name,description) AGAINST(? IN NATURAL LANGUAGE MODE)")
where, args := ft.Where("foo")
// SELECT * FROM products WHERE MATCH(name,description) AGAINST(? IN NATURAL LANGUAGE MODE)
// args = ["foo"]
```

## SearchEngine Interface

Elasticsearch vb. tam entegrasyonlar için `search.SearchEngine` interface'ini implement edin.
