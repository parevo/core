package auth

import "errors"

var (
	ErrInvalidConfig   = errors.New("invalid config")
	ErrUnauthenticated = errors.New("unauthenticated")
	ErrForbidden       = errors.New("forbidden")
	ErrTenantMismatch  = errors.New("tenant mismatch")
	ErrMissingTenant   = errors.New("missing tenant")
	ErrInvalidToken    = errors.New("invalid token")
	ErrSessionRevoked  = errors.New("session revoked")
	ErrInvalidRefresh  = errors.New("invalid refresh token")
	ErrRefreshReuse    = errors.New("refresh token reuse detected")
)
