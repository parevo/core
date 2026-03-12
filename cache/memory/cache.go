package memory

import (
	"context"
	"sync"
	"time"

	"github.com/parevo/core/cache"
)

type entry struct {
	value []byte
	until time.Time
}

// Cache implements cache.Store with in-memory storage.
type Cache struct {
	mu       sync.RWMutex
	data     map[string]entry
	defaultTTL time.Duration
}

// NewCache creates an in-memory cache with optional default TTL.
func NewCache(defaultTTL time.Duration) *Cache {
	return &Cache{
		data:       make(map[string]entry),
		defaultTTL: defaultTTL,
	}
}

// Get returns the value for key, or nil if not found/expired.
func (c *Cache) Get(ctx context.Context, key string) ([]byte, error) {
	c.mu.RLock()
	e, ok := c.data[key]
	c.mu.RUnlock()
	if !ok || time.Now().After(e.until) {
		return nil, nil
	}
	return append([]byte(nil), e.value...), nil
}

// Set stores value with TTL. If ttl <= 0, defaultTTL is used.
func (c *Cache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	if ttl <= 0 {
		ttl = c.defaultTTL
	}
	until := time.Now().Add(ttl)
	c.mu.Lock()
	c.data[key] = entry{value: append([]byte(nil), value...), until: until}
	c.mu.Unlock()
	return nil
}

// Delete removes the key.
func (c *Cache) Delete(ctx context.Context, key string) error {
	c.mu.Lock()
	delete(c.data, key)
	c.mu.Unlock()
	return nil
}

var _ cache.Store = (*Cache)(nil)
