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

// BlacklistChecker checks if a token (by jti or hash) is blacklisted. Optional; used for logout.
type BlacklistChecker interface {
	Check(ctx context.Context, jtiOrHash string) error
}

type Modules struct {
	Tenant         *tenant.Service
	Permission     *permission.Service
	TenantOverride tenant.OverridePolicy
	SessionStore   storage.SessionStore
	RefreshStore   storage.RefreshTokenStore
	Blacklist      BlacklistChecker
	AuditLogger    AuditLogger
	Metrics        MetricsRecorder
	Tracer         Tracer
	APIKey         APIKeyValidator
}
