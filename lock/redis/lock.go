package redis

import (
	"context"
	"time"

	"github.com/parevo/core/lock"
	"github.com/redis/go-redis/v9"
)

const defaultPrefix = "lock:"

// Locker implements lock.Locker using Redis SET NX.
type Locker struct {
	client *redis.Client
	prefix string
}

// NewLocker creates a Redis-backed locker.
func NewLocker(client *redis.Client, prefix string) *Locker {
	if prefix == "" {
		prefix = defaultPrefix
	}
	return &Locker{client: client, prefix: prefix}
}

func (l *Locker) key(k string) string {
	return l.prefix + k
}

// Lock acquires the lock. Returns true if acquired, false if already held.
func (l *Locker) Lock(ctx context.Context, key string, ttl time.Duration) (bool, error) {
	err := l.client.SetArgs(ctx, l.key(key), "1", redis.SetArgs{Mode: "NX", TTL: ttl}).Err()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// Unlock releases the lock.
func (l *Locker) Unlock(ctx context.Context, key string) error {
	return l.client.Del(ctx, l.key(key)).Err()
}

var _ lock.Locker = (*Locker)(nil)
