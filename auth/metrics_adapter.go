package auth

import "github.com/parevo/core/observability/metrics"

type prometheusMetrics struct{}

func (prometheusMetrics) IncAuthSuccess(flow string) { metrics.IncAuthSuccess(flow) }
func (prometheusMetrics) IncAuthFail(reason string)  { metrics.IncAuthFail(reason) }
func (prometheusMetrics) ObserveAuthDuration(operation string, seconds float64) {
	metrics.ObserveAuthDuration(operation, seconds)
}

func PrometheusMetrics() MetricsRecorder {
	metrics.Register(nil)
	return prometheusMetrics{}
}
