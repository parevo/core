package ratelimit

import (
	"sync"
	"time"
)

type entry struct {
	count       int
	windowStart time.Time
}

type MemoryLimiter struct {
	mu       sync.Mutex
	limit    int
	window   time.Duration
	counters map[string]entry
}

func NewMemoryLimiter(limit int, window time.Duration) *MemoryLimiter {
	if limit <= 0 {
		limit = 1
	}
	if window <= 0 {
		window = time.Minute
	}
	return &MemoryLimiter{
		limit:    limit,
		window:   window,
		counters: map[string]entry{},
	}
}

func (l *MemoryLimiter) Allow(key string) bool {
	now := time.Now().UTC()
	l.mu.Lock()
	defer l.mu.Unlock()

	current := l.counters[key]
	if current.windowStart.IsZero() || now.Sub(current.windowStart) >= l.window {
		l.counters[key] = entry{count: 1, windowStart: now}
		return true
	}
	if current.count >= l.limit {
		return false
	}
	current.count++
	l.counters[key] = current
	return true
}
