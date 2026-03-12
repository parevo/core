package adapters

import "github.com/parevo/core/auth"

type ErrorHandler func(err error) (status int, body string)

type Options struct {
	TenantHeader   string
	OverridePolicy auth.TenantOverridePolicy
	ErrorHandler   ErrorHandler
}

func (o Options) WithDefaults() Options {
	normalized := o
	if normalized.TenantHeader == "" {
		normalized.TenantHeader = "X-Tenant-Id"
	}
	if normalized.ErrorHandler == nil {
		normalized.ErrorHandler = DefaultErrorHandler
	}
	return normalized
}

func DefaultErrorHandler(err error) (int, string) {
	status := auth.HTTPStatusForError(err)
	switch status {
	case 400:
		return status, "bad request"
	case 401:
		return status, "unauthorized"
	case 403:
		return status, "forbidden"
	default:
		return status, "internal server error"
	}
}
