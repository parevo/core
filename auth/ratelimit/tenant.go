package ratelimit

import (
	"sync"
	"time"
)

type TenantLimitConfig struct {
	Limit  int
	Window time.Duration
}

type TenantLimiter struct {
	mu        sync.Mutex
	defaultC  TenantLimitConfig
	overrides map[string]TenantLimitConfig
	counters  map[string]entry
}

func NewTenantLimiter(defaultLimit int, defaultWindow time.Duration) *TenantLimiter {
	if defaultLimit <= 0 {
		defaultLimit = 100
	}
	if defaultWindow <= 0 {
		defaultWindow = time.Minute
	}
	return &TenantLimiter{
		defaultC:  TenantLimitConfig{Limit: defaultLimit, Window: defaultWindow},
		overrides: map[string]TenantLimitConfig{},
		counters:  map[string]entry{},
	}
}

func (t *TenantLimiter) SetTenantLimit(tenantID string, limit int, window time.Duration) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.overrides == nil {
		t.overrides = map[string]TenantLimitConfig{}
	}
	t.overrides[tenantID] = TenantLimitConfig{Limit: limit, Window: window}
}

func (t *TenantLimiter) configFor(tenantID string) TenantLimitConfig {
	if c, ok := t.overrides[tenantID]; ok {
		return c
	}
	return t.defaultC
}

func (t *TenantLimiter) Allow(tenantID, subKey string) bool {
	if tenantID == "" {
		tenantID = "_default"
	}
	key := "tenant:" + tenantID + ":" + subKey
	cfg := t.configFor(tenantID)

	t.mu.Lock()
	defer t.mu.Unlock()

	now := time.Now().UTC()
	current := t.counters[key]
	if current.windowStart.IsZero() || now.Sub(current.windowStart) >= cfg.Window {
		t.counters[key] = entry{count: 1, windowStart: now}
		return true
	}
	if current.count >= cfg.Limit {
		return false
	}
	current.count++
	t.counters[key] = current
	return true
}
