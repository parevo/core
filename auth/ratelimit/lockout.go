package ratelimit

import (
	"sync"
	"time"
)

type LockoutManager struct {
	limiter      *MemoryLimiter
	lockDuration time.Duration
	mu           sync.Mutex
	lockedUntil  map[string]time.Time
}

func NewLockoutManager(maxFailures int, window time.Duration, lockDuration time.Duration) *LockoutManager {
	if lockDuration <= 0 {
		lockDuration = 15 * time.Minute
	}
	return &LockoutManager{
		limiter:      NewMemoryLimiter(maxFailures, window),
		lockDuration: lockDuration,
		lockedUntil:  map[string]time.Time{},
	}
}

func (m *LockoutManager) IsLocked(key string) bool {
	now := time.Now().UTC()
	m.mu.Lock()
	defer m.mu.Unlock()
	until, ok := m.lockedUntil[key]
	if !ok {
		return false
	}
	if now.After(until) {
		delete(m.lockedUntil, key)
		return false
	}
	return true
}

func (m *LockoutManager) RegisterFailure(key string) {
	if key == "" {
		return
	}
	if m.limiter.Allow(key) {
		return
	}
	m.mu.Lock()
	m.lockedUntil[key] = time.Now().UTC().Add(m.lockDuration)
	m.mu.Unlock()
}

func (m *LockoutManager) RegisterSuccess(key string) {
	if key == "" {
		return
	}
	m.mu.Lock()
	delete(m.lockedUntil, key)
	m.mu.Unlock()
}
