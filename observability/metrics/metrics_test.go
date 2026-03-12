package metrics

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
)

func TestMetrics(t *testing.T) {
	reg := prometheus.NewRegistry()
	Register(reg)
	IncAuthSuccess("authenticate")
	IncAuthFail("invalid_token")
	ObserveAuthDuration("parse", 0.001)
}
