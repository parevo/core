// Package lock provides distributed locking for rate limiting, job deduplication, and critical sections.
//
// Providers:
//   - lock/memory — in-memory (single-instance)
//   - lock/redis — Redis (multi-instance)
package lock
