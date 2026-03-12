// Package lock provides distributed locking for rate limiting, job deduplication, and critical sections.
//
// # Usage
//
// Memory (single-instance):
//
//	l := memory.NewLocker()
//	ok, err := l.Lock(ctx, "job:123", 30*time.Second)
//	defer l.Unlock(ctx, "job:123")
//
// Redis (multi-instance):
//
//	l := redis.NewLocker(redisClient, "lock:")
//	ok, err := l.Lock(ctx, "job:123", 30*time.Second)
package lock

import (
	"context"
	"time"
)

// Locker provides distributed lock acquire/release.
type Locker interface {
	Lock(ctx context.Context, key string, ttl time.Duration) (bool, error)
	Unlock(ctx context.Context, key string) error
}
