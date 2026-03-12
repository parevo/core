package memory

import (
	"context"
	"sync"
	"time"

	"github.com/parevo/core/lock"
)

type lockEntry struct {
	until time.Time
}

// Locker implements lock.Locker with in-memory storage (single-instance only).
type Locker struct {
	mu    sync.Mutex
	locks map[string]lockEntry
}

// NewLocker creates an in-memory locker.
func NewLocker() *Locker {
	return &Locker{locks: make(map[string]lockEntry)}
}

// Lock acquires the lock. Returns true if acquired, false if already held.
func (l *Locker) Lock(ctx context.Context, key string, ttl time.Duration) (bool, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	now := time.Now()
	if e, ok := l.locks[key]; ok && now.Before(e.until) {
		return false, nil
	}
	l.locks[key] = lockEntry{until: now.Add(ttl)}
	return true, nil
}

// Unlock releases the lock.
func (l *Locker) Unlock(ctx context.Context, key string) error {
	l.mu.Lock()
	delete(l.locks, key)
	l.mu.Unlock()
	return nil
}

var _ lock.Locker = (*Locker)(nil)
