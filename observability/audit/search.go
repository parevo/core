package audit

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"io"
	"time"
)

// Event represents an audit log entry for search/export.
type Event struct {
	ID        string
	Timestamp time.Time
	Event     string
	UserID    string
	TenantID  string
	IP        string
	Attrs     map[string]string
}

// Query filters audit events.
type Query struct {
	From      time.Time
	To        time.Time
	Event     string
	UserID    string
	TenantID  string
	Limit     int
	Offset    int
}

// EventStore provides audit event search and export.
type EventStore interface {
	Search(ctx context.Context, q Query) ([]Event, error)
	ExportJSON(ctx context.Context, q Query, w io.Writer) error
	ExportCSV(ctx context.Context, q Query, w io.Writer) error
}

// InMemoryEventStore is a simple in-memory implementation for development.
type InMemoryEventStore struct {
	events []Event
}

// NewInMemoryEventStore creates an in-memory audit store.
func NewInMemoryEventStore() *InMemoryEventStore {
	return &InMemoryEventStore{events: make([]Event, 0)}
}

// Append adds an event.
func (s *InMemoryEventStore) Append(e Event) {
	s.events = append(s.events, e)
}

func (s *InMemoryEventStore) Search(ctx context.Context, q Query) ([]Event, error) {
	_ = ctx
	var out []Event
	for _, e := range s.events {
		if !q.From.IsZero() && e.Timestamp.Before(q.From) {
			continue
		}
		if !q.To.IsZero() && e.Timestamp.After(q.To) {
			continue
		}
		if q.Event != "" && e.Event != q.Event {
			continue
		}
		if q.UserID != "" && e.UserID != q.UserID {
			continue
		}
		if q.TenantID != "" && e.TenantID != q.TenantID {
			continue
		}
		out = append(out, e)
	}
	offset := q.Offset
	if offset < 0 {
		offset = 0
	}
	if offset < len(out) {
		out = out[offset:]
	}
	if q.Limit > 0 && len(out) > q.Limit {
		out = out[:q.Limit]
	}
	return out, nil
}

func (s *InMemoryEventStore) ExportJSON(ctx context.Context, q Query, w io.Writer) error {
	events, err := s.Search(ctx, q)
	if err != nil {
		return err
	}
	return json.NewEncoder(w).Encode(events)
}

func (s *InMemoryEventStore) ExportCSV(ctx context.Context, q Query, w io.Writer) error {
	events, err := s.Search(ctx, q)
	if err != nil {
		return err
	}
	cw := csv.NewWriter(w)
	cw.Write([]string{"id", "timestamp", "event", "user_id", "tenant_id", "ip"})
	for _, e := range events {
		cw.Write([]string{e.ID, e.Timestamp.Format(time.RFC3339), e.Event, e.UserID, e.TenantID, e.IP})
	}
	cw.Flush()
	return cw.Error()
}
