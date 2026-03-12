package memory

import (
	"context"
	"sync"

	"github.com/parevo/core/storage"
)

// TenantFeatureStore implements storage.TenantFeatureStore with in-memory storage.
type TenantFeatureStore struct {
	mu       sync.RWMutex
	Plans    map[string]storage.Plan
	Features map[string]map[string]bool  // tenantID -> feature -> enabled
	Limits   map[string]map[string]int64  // tenantID -> limitKey -> value
}

// NewTenantFeatureStore creates an in-memory TenantFeatureStore.
func NewTenantFeatureStore() *TenantFeatureStore {
	return &TenantFeatureStore{
		Plans:    make(map[string]storage.Plan),
		Features: make(map[string]map[string]bool),
		Limits:   make(map[string]map[string]int64),
	}
}

func (s *TenantFeatureStore) GetPlan(_ context.Context, tenantID string) (storage.Plan, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if p, ok := s.Plans[tenantID]; ok {
		return p, nil
	}
	return storage.PlanFree, nil
}

func (s *TenantFeatureStore) IsEnabled(_ context.Context, tenantID, feature string) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if m, ok := s.Features[tenantID]; ok && m[feature] {
		return true, nil
	}
	switch s.Plans[tenantID] {
	case storage.PlanEnterprise:
		return true, nil
	case storage.PlanPro:
		return feature != "saml" && feature != "ldap", nil
	default:
		return feature == "basic_auth", nil
	}
}

func (s *TenantFeatureStore) GetLimit(_ context.Context, tenantID, limitKey string) (int64, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if m, ok := s.Limits[tenantID]; ok {
		if v, ok := m[limitKey]; ok {
			return v, nil
		}
	}
	switch s.Plans[tenantID] {
	case storage.PlanEnterprise:
		return 0, nil
	case storage.PlanPro:
		if limitKey == "max_users" {
			return 100, nil
		}
		if limitKey == "api_calls_per_month" {
			return 100000, nil
		}
	default:
		if limitKey == "max_users" {
			return 5, nil
		}
		if limitKey == "api_calls_per_month" {
			return 1000, nil
		}
	}
	return 0, nil
}

func (s *TenantFeatureStore) SetPlan(_ context.Context, tenantID string, plan storage.Plan) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Plans[tenantID] = plan
	return nil
}

func (s *TenantFeatureStore) SetFeature(_ context.Context, tenantID, feature string, enabled bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.Features[tenantID] == nil {
		s.Features[tenantID] = make(map[string]bool)
	}
	s.Features[tenantID][feature] = enabled
	return nil
}

func (s *TenantFeatureStore) SetLimit(_ context.Context, tenantID, limitKey string, value int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.Limits[tenantID] == nil {
		s.Limits[tenantID] = make(map[string]int64)
	}
	s.Limits[tenantID][limitKey] = value
	return nil
}
