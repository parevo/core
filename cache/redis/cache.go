package redis

import (
	"context"
	"time"

	"github.com/parevo/core/cache"
	"github.com/redis/go-redis/v9"
)

const defaultPrefix = "cache:"

// Cache implements cache.Store using Redis.
type Cache struct {
	client *redis.Client
	prefix string
}

// NewCache creates a Redis-backed cache.
func NewCache(client *redis.Client, prefix string) *Cache {
	if prefix == "" {
		prefix = defaultPrefix
	}
	return &Cache{client: client, prefix: prefix}
}

func (c *Cache) key(k string) string {
	return c.prefix + k
}

// Get returns the value for key, or nil if not found.
func (c *Cache) Get(ctx context.Context, key string) ([]byte, error) {
	val, err := c.client.Get(ctx, c.key(key)).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	return val, err
}

// Set stores value with TTL.
func (c *Cache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	if ttl <= 0 {
		ttl = 0 // no expiry
	}
	return c.client.Set(ctx, c.key(key), value, ttl).Err()
}

// Delete removes the key.
func (c *Cache) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, c.key(key)).Err()
}

var _ cache.Store = (*Cache)(nil)
