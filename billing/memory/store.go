package memory

import (
	"context"
	"sync"
	"time"

	"github.com/parevo/core/billing"
)

type usageEntry struct {
	tenantID string
	metric   string
	value    int64
	at       time.Time
}

// UsageStore implements billing.UsageStore with in-memory storage.
type UsageStore struct {
	mu     sync.RWMutex
	events []usageEntry
}

// NewUsageStore creates an in-memory usage store.
func NewUsageStore() *UsageStore {
	return &UsageStore{events: make([]usageEntry, 0)}
}

// Record adds a usage event.
func (s *UsageStore) Record(ctx context.Context, tenantID, metric string, value int64) error {
	s.mu.Lock()
	s.events = append(s.events, usageEntry{tenantID: tenantID, metric: metric, value: value, at: time.Now()})
	s.mu.Unlock()
	return nil
}

// Usage returns the sum of values for the tenant/metric in the time range.
func (s *UsageStore) Usage(ctx context.Context, tenantID, metric string, from, to time.Time) (int64, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var sum int64
	for _, e := range s.events {
		if e.tenantID == tenantID && e.metric == metric && !e.at.Before(from) && !e.at.After(to) {
			sum += e.value
		}
	}
	return sum, nil
}

var _ billing.UsageStore = (*UsageStore)(nil)
