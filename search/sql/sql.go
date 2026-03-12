package sql

// FullText provides SQL full-text search building. Use with MySQL FULLTEXT or Postgres tsvector.
type FullText struct {
	MatchExpr string // e.g. "MATCH(col1,col2) AGAINST(? IN NATURAL LANGUAGE MODE)"
}

// NewFullText creates a SQL full-text helper.
func NewFullText(matchExpr string) *FullText {
	return &FullText{MatchExpr: matchExpr}
}

// Where returns a WHERE clause and args for the query.
func (f *FullText) Where(query string) (string, []any) {
	if query == "" {
		return "1=1", nil
	}
	return f.MatchExpr, []any{query}
}
