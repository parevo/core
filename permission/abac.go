package permission

import (
	"context"
	"errors"

	"github.com/parevo/core/storage"
)

var ErrABACDenied = errors.New("ABAC condition denied")

// Attribute represents a key-value attribute for ABAC evaluation.
type Attribute struct {
	Key   string
	Value string
}

// Resource represents the resource being accessed with attributes.
type Resource struct {
	ID          string
	OwnerID     string
	Department  string
	Environment string
	Attributes  map[string]string
}

// SubjectAttributes extends storage.Subject with ABAC attributes.
type SubjectAttributes struct {
	Subject     storage.Subject
	Department  string
	Environment string
	Attributes  map[string]string
}

// ABACCondition evaluates whether access is allowed based on attributes.
type ABACCondition interface {
	Allow(ctx context.Context, subject SubjectAttributes, resource Resource, permission string) (bool, error)
}

// ResourceOwnerCondition allows access if subject is the resource owner.
type ResourceOwnerCondition struct{}

func (ResourceOwnerCondition) Allow(_ context.Context, subject SubjectAttributes, resource Resource, _ string) (bool, error) {
	return resource.OwnerID != "" && resource.OwnerID == subject.Subject.ID, nil
}

// DepartmentCondition allows access if subject's department matches resource's department.
type DepartmentCondition struct {
	RequireMatch bool
}

func (c DepartmentCondition) Allow(_ context.Context, subject SubjectAttributes, resource Resource, _ string) (bool, error) {
	if resource.Department == "" || subject.Department == "" {
		return !c.RequireMatch, nil
	}
	return subject.Department == resource.Department, nil
}

// EnvironmentCondition allows access if subject's environment matches resource's environment.
type EnvironmentCondition struct {
	AllowedEnvs []string // e.g. ["prod", "staging"] - empty means all
}

func (c EnvironmentCondition) Allow(_ context.Context, subject SubjectAttributes, resource Resource, _ string) (bool, error) {
	env := resource.Environment
	if env == "" {
		env = subject.Environment
	}
	if len(c.AllowedEnvs) == 0 {
		return true, nil
	}
	for _, e := range c.AllowedEnvs {
		if e == env {
			return true, nil
		}
	}
	return false, nil
}

// AttributeMatchCondition allows access if subject has required attribute.
type AttributeMatchCondition struct {
	RequiredKey   string
	RequiredValue string
}

func (c AttributeMatchCondition) Allow(_ context.Context, subject SubjectAttributes, resource Resource, _ string) (bool, error) {
	if v, ok := subject.Attributes[c.RequiredKey]; ok && (c.RequiredValue == "" || v == c.RequiredValue) {
		return true, nil
	}
	if v, ok := resource.Attributes[c.RequiredKey]; ok && (c.RequiredValue == "" || v == c.RequiredValue) {
		return true, nil
	}
	return false, nil
}

// ABACPolicy combines multiple conditions with AND logic.
type ABACPolicy struct {
	Conditions []ABACCondition
	Mode      ABACMode
}

type ABACMode int

const (
	ABACModeAll  ABACMode = iota // all conditions must pass
	ABACModeAny                  // at least one condition must pass
)

func (p *ABACPolicy) Allow(ctx context.Context, subject SubjectAttributes, resource Resource, permission string) (bool, error) {
	if len(p.Conditions) == 0 {
		return true, nil
	}
	var pass int
	for _, cond := range p.Conditions {
		ok, err := cond.Allow(ctx, subject, resource, permission)
		if err != nil {
			return false, err
		}
		if ok {
			pass++
			if p.Mode == ABACModeAny {
				return true, nil
			}
		} else if p.Mode == ABACModeAll {
			return false, nil
		}
	}
	return p.Mode == ABACModeAll && pass == len(p.Conditions), nil
}

// ABACService combines RBAC permission check with ABAC conditions.
type ABACService struct {
	perm *Service
}

func NewABACService(perm *Service) *ABACService {
	return &ABACService{perm: perm}
}

// Check performs permission check and ABAC evaluation.
func (s *ABACService) Check(ctx context.Context, subject SubjectAttributes, tenantID, permission string, resource Resource, policy *ABACPolicy) error {
	if err := s.perm.Check(ctx, subject.Subject, tenantID, permission); err != nil {
		return err
	}
	if policy != nil {
		ok, err := policy.Allow(ctx, subject, resource, permission)
		if err != nil {
			return err
		}
		if !ok {
			return ErrABACDenied
		}
	}
	return nil
}
