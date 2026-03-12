package auth

import "context"

type AuditEvent string

const (
	AuditAuthSuccess    AuditEvent = "auth_success"
	AuditAuthFailed     AuditEvent = "auth_failed"
	AuditTenantMismatch AuditEvent = "tenant_mismatch"
	AuditPermissionDeny AuditEvent = "permission_denied"
	AuditSessionRevoked AuditEvent = "session_revoked"
)

type AuditLogger interface {
	Log(ctx context.Context, event AuditEvent, attrs map[string]string)
}

func (s *Service) audit(ctx context.Context, event AuditEvent, attrs map[string]string) {
	if s.modules.AuditLogger != nil {
		s.modules.AuditLogger.Log(ctx, event, attrs)
	}
	if s.modules.Metrics != nil {
		switch event {
		case AuditAuthSuccess:
			flow := attrs["flow"]
			if flow == "" {
				flow = "authenticate"
			}
			s.modules.Metrics.IncAuthSuccess(flow)
		case AuditAuthFailed:
			reason := attrs["reason"]
			if reason == "" {
				reason = "unknown"
			}
			s.modules.Metrics.IncAuthFail(reason)
		case AuditTenantMismatch:
			s.modules.Metrics.IncAuthFail("tenant_mismatch")
		case AuditPermissionDeny:
			s.modules.Metrics.IncAuthFail("permission_denied")
		}
	}
}
