package tracing

import (
	"context"
)

type Tracer interface {
	StartSpan(ctx context.Context, name string, attrs map[string]string) (context.Context, func())
}

type noopTracer struct{}

func (noopTracer) StartSpan(ctx context.Context, _ string, _ map[string]string) (context.Context, func()) {
	return ctx, func() {}
}

var defaultTracer Tracer = noopTracer{}

func SetTracer(t Tracer) {
	if t != nil {
		defaultTracer = t
	}
}

func StartSpan(ctx context.Context, name string, attrs map[string]string) (context.Context, func()) {
	return defaultTracer.StartSpan(ctx, name, attrs)
}
