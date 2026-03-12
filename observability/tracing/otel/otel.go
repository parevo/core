package otel

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/parevo/core/observability/tracing"
)

type otelTracer struct {
	tracer trace.Tracer
}

func New(serviceName string) *otelTracer {
	if serviceName == "" {
		serviceName = "parevo"
	}
	return &otelTracer{
		tracer: otel.Tracer(serviceName),
	}
}

func (t *otelTracer) StartSpan(ctx context.Context, name string, attrs map[string]string) (context.Context, func()) {
	var opts []trace.SpanStartOption
	if len(attrs) > 0 {
		kvs := make([]attribute.KeyValue, 0, len(attrs))
		for k, v := range attrs {
			kvs = append(kvs, attribute.String(k, v))
		}
		opts = append(opts, trace.WithAttributes(kvs...))
	}
	ctx, span := t.tracer.Start(ctx, name, opts...)
	return ctx, func() { span.End() }
}

func InitOtel(serviceName string) {
	tracing.SetTracer(New(serviceName))
}
