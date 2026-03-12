package memory

import (
	"context"
	"sync"
)

// IPFilterStore implements ipfilter.Store in memory.
type IPFilterStore struct {
	mu        sync.RWMutex
	allowed   map[string]map[string]bool // tenantID -> ip -> true
	blocked   map[string]map[string]bool // tenantID -> ip -> true
	globalAllowed map[string]bool
	globalBlocked map[string]bool
}

// NewIPFilterStore creates an in-memory IP filter store.
func NewIPFilterStore() *IPFilterStore {
	return &IPFilterStore{
		allowed:       make(map[string]map[string]bool),
		blocked:       make(map[string]map[string]bool),
		globalAllowed: make(map[string]bool),
		globalBlocked: make(map[string]bool),
	}
}

func (s *IPFilterStore) IsAllowed(_ context.Context, tenantID, ip string) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.globalAllowed[ip] {
		return true, nil
	}
	if m, ok := s.allowed[tenantID]; ok && m[ip] {
		return true, nil
	}
	return false, nil
}

func (s *IPFilterStore) IsBlocked(_ context.Context, tenantID, ip string) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.globalBlocked[ip] {
		return true, nil
	}
	if m, ok := s.blocked[tenantID]; ok && m[ip] {
		return true, nil
	}
	return false, nil
}

// AllowIP adds an IP to the allowlist.
func (s *IPFilterStore) AllowIP(tenantID, ip string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if tenantID == "" {
		s.globalAllowed[ip] = true
		return
	}
	if s.allowed[tenantID] == nil {
		s.allowed[tenantID] = make(map[string]bool)
	}
	s.allowed[tenantID][ip] = true
}

// BlockIP adds an IP to the blocklist.
func (s *IPFilterStore) BlockIP(tenantID, ip string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if tenantID == "" {
		s.globalBlocked[ip] = true
		return
	}
	if s.blocked[tenantID] == nil {
		s.blocked[tenantID] = make(map[string]bool)
	}
	s.blocked[tenantID][ip] = true
}
