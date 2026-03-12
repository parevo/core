package auth

import (
	"context"

	"github.com/parevo/core/observability/logging"
)

type StructuredAuditLogger struct {
	Logger *logging.Logger
}

func NewStructuredAuditLogger(logger *logging.Logger) *StructuredAuditLogger {
	return &StructuredAuditLogger{Logger: logger}
}

func (l *StructuredAuditLogger) Log(ctx context.Context, event AuditEvent, attrs map[string]string) {
	if l == nil || l.Logger == nil {
		return
	}
	fields := map[string]any{
		"event": string(event),
	}
	for k, v := range attrs {
		fields[k] = v
	}

	switch event {
	case AuditAuthFailed, AuditPermissionDeny:
		l.Logger.Warn(ctx, "auth.audit", fields)
	case AuditTenantMismatch:
		l.Logger.Warn(ctx, "auth.tenant", fields)
	default:
		l.Logger.Info(ctx, "auth.audit", fields)
	}
}
