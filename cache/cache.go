package cache

import (
	"context"
	"time"
)

// Store provides key-value cache operations with TTL.
type Store interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
}
