package permission

import (
	"context"
	"sync"
	"time"

	"github.com/parevo/core/storage"
)

type cacheEntry struct {
	ok    bool
	until time.Time
}

type CachedPermissionStore struct {
	mu     sync.RWMutex
	store  storage.PermissionStore
	ttl    time.Duration
	cache  map[string]cacheEntry
}

func NewCachedPermissionStore(store storage.PermissionStore, ttl time.Duration) *CachedPermissionStore {
	if ttl <= 0 {
		ttl = 5 * time.Minute
	}
	return &CachedPermissionStore{
		store: store,
		ttl:   ttl,
		cache: map[string]cacheEntry{},
	}
}

func cacheKey(subjectID, tenantID, permission string) string {
	return subjectID + "|" + tenantID + "|" + permission
}

func (c *CachedPermissionStore) HasPermission(ctx context.Context, subjectID, tenantID, permission string, roles []string) (bool, error) {
	key := cacheKey(subjectID, tenantID, permission)
	c.mu.RLock()
	e, ok := c.cache[key]
	c.mu.RUnlock()
	if ok && time.Now().Before(e.until) {
		return e.ok, nil
	}

	ok2, err := c.store.HasPermission(ctx, subjectID, tenantID, permission, roles)
	if err != nil {
		return false, err
	}
	until := time.Now().Add(c.ttl)
	c.mu.Lock()
	if c.cache == nil {
		c.cache = map[string]cacheEntry{}
	}
	c.cache[key] = cacheEntry{ok: ok2, until: until}
	c.mu.Unlock()
	return ok2, nil
}

func (c *CachedPermissionStore) InvalidateSubject(subjectID string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	prefix := subjectID + "|"
	for k := range c.cache {
		if len(k) > len(prefix) && k[:len(prefix)] == prefix {
			delete(c.cache, k)
		}
	}
}

func (c *CachedPermissionStore) InvalidateSubjectTenant(subjectID, tenantID string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	prefix := subjectID + "|" + tenantID + "|"
	for k := range c.cache {
		if len(k) > len(prefix) && k[:len(prefix)] == prefix {
			delete(c.cache, k)
		}
	}
}

func (c *CachedPermissionStore) InvalidateAll() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache = map[string]cacheEntry{}
}
