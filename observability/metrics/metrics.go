package metrics

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	authSuccessTotal *prometheus.CounterVec
	authFailTotal    *prometheus.CounterVec
	authDuration     *prometheus.HistogramVec
	initOnce         sync.Once
)

func initMetrics() {
	authSuccessTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "parevo_auth_success_total",
			Help: "Total number of successful authentications",
		},
		[]string{"flow"},
	)
	authFailTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "parevo_auth_fail_total",
			Help: "Total number of failed authentications",
		},
		[]string{"reason"},
	)
	authDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "parevo_auth_duration_seconds",
			Help:    "Authentication duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"operation"},
	)
}

func Register(reg prometheus.Registerer) {
	initOnce.Do(func() {
		initMetrics()
		if reg == nil {
			reg = prometheus.DefaultRegisterer
		}
		reg.MustRegister(authSuccessTotal, authFailTotal, authDuration)
	})
}

func IncAuthSuccess(flow string) {
	initOnce.Do(func() { initMetrics() })
	authSuccessTotal.WithLabelValues(flow).Inc()
}

func IncAuthFail(reason string) {
	initOnce.Do(func() { initMetrics() })
	authFailTotal.WithLabelValues(reason).Inc()
}

func ObserveAuthDuration(operation string, seconds float64) {
	initOnce.Do(func() { initMetrics() })
	authDuration.WithLabelValues(operation).Observe(seconds)
}
