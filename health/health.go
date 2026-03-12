// Package health provides health check handlers for DB, Redis, blob, and custom checks.
//
// # Usage
//
//	h := health.NewChecker()
//	h.Add("db", health.PingDB(db))
//	h.Add("redis", health.PingRedis(redisClient))
//	h.Add("s3", health.PingBlob(blobStore, "bucket", "health"))
//
//	// In HTTP handler:
//	if !h.Check(ctx) {
//	    w.WriteHeader(http.StatusServiceUnavailable)
//	}
package health

import (
	"context"
	"database/sql"
	"io"
	"strings"
	"time"

	"github.com/parevo/core/blob"
	"github.com/redis/go-redis/v9"
)

// Checker runs named health checks.
type Checker struct {
	checks map[string]Check
}

// Check is a health check function. Returns nil if healthy.
type Check func(ctx context.Context) error

// NewChecker creates a health checker.
func NewChecker() *Checker {
	return &Checker{checks: make(map[string]Check)}
}

// Add registers a named check.
func (c *Checker) Add(name string, check Check) {
	c.checks[name] = check
}

// Check runs all checks. Returns false if any fails.
func (c *Checker) Check(ctx context.Context) bool {
	for _, check := range c.checks {
		if check(ctx) != nil {
			return false
		}
	}
	return true
}

// CheckWithResults returns per-check results.
func (c *Checker) CheckWithResults(ctx context.Context) map[string]error {
	results := make(map[string]error)
	for name, check := range c.checks {
		results[name] = check(ctx)
	}
	return results
}

// PingDB returns a check that pings the database.
func PingDB(db *sql.DB) Check {
	return func(ctx context.Context) error {
		return db.PingContext(ctx)
	}
}

// PingRedis returns a check that pings Redis.
func PingRedis(client *redis.Client) Check {
	return func(ctx context.Context) error {
		return client.Ping(ctx).Err()
	}
}

// PingRedisWithTimeout wraps PingRedis with a timeout.
func PingRedisWithTimeout(client *redis.Client, timeout time.Duration) Check {
	return func(ctx context.Context) error {
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()
		return client.Ping(ctx).Err()
	}
}

// PingBlob returns a check that verifies blob store accessibility (List call).
func PingBlob(store blob.Store, bucket, prefix string) Check {
	return func(ctx context.Context) error {
		_, err := store.List(ctx, bucket, prefix)
		return err
	}
}

// PingBlobPutGet returns a check that writes and reads a small object (more thorough).
func PingBlobPutGet(store blob.Store, bucket, key string) Check {
	return func(ctx context.Context) error {
		body := strings.NewReader("health")
		if err := store.Put(ctx, bucket, key, body, "text/plain"); err != nil {
			return err
		}
		rc, err := store.Get(ctx, bucket, key)
		if err != nil {
			return err
		}
		_, _ = io.Copy(io.Discard, rc)
		_ = rc.Close()
		_ = store.Delete(ctx, bucket, key)
		return nil
	}
}
