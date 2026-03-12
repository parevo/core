package auth

import (
	"context"

	"github.com/parevo/core/permission"
	"github.com/parevo/core/storage"
	"github.com/parevo/core/tenant"
)

type APIKeyValidator interface {
	Validate(ctx context.Context, rawKey string) (userID, tenantID string, err error)
}

type MetricsRecorder interface {
	IncAuthSuccess(flow string)
	IncAuthFail(reason string)
	ObserveAuthDuration(operation string, seconds float64)
}

type Tracer interface {
	StartSpan(ctx context.Context, name string, attrs map[string]string) (context.Context, func())
}

type Modules struct {
	Tenant         *tenant.Service
	Permission     *permission.Service
	TenantOverride tenant.OverridePolicy
	SessionStore   storage.SessionStore
	RefreshStore   storage.RefreshTokenStore
	AuditLogger    AuditLogger
	Metrics        MetricsRecorder
	Tracer         Tracer
	APIKey         APIKeyValidator
}
