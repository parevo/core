// Package search provides a search interface for full-text and structured queries.
//
// # Usage
//
//	// SQL full-text (simple)
//	engine := search.NewSQLFullText("MATCH(name,description) AGAINST(? IN NATURAL LANGUAGE MODE)")
//	where, args := engine.Query("foo")
//
//	// Custom backends (Elasticsearch, etc.) implement SearchEngine.
package search

import (
	"context"
)

// SearchEngine performs search queries. Implement with SQL full-text, Elasticsearch, etc.
type SearchEngine interface {
	Search(ctx context.Context, query string, opts *SearchOptions) (*SearchResult, error)
}

// SearchOptions configures search behavior.
type SearchOptions struct {
	Limit   int
	Offset  int
	Filters map[string]string
}

// SearchResult holds search hits.
type SearchResult struct {
	Hits  []SearchHit
	Total int
}

// SearchHit is a single search result.
type SearchHit struct {
	ID     string
	Score  float64
	Fields map[string]any
}
