package auth

import (
	"context"

	"github.com/parevo/core/observability/tracing"
)

type tracingAdapter struct{}

func (tracingAdapter) StartSpan(ctx context.Context, name string, attrs map[string]string) (context.Context, func()) {
	return tracing.StartSpan(ctx, name, attrs)
}

func TracingAdapter() Tracer {
	return tracingAdapter{}
}
