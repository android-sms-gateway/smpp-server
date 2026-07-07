package smsgate

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	metricSuccess = "success"
	metricFailure = "failure"
)

type Metrics struct {
	requestsTotal          *prometheus.CounterVec
	requestDurationSeconds *prometheus.HistogramVec
}

func NewMetrics() *Metrics {
	return &Metrics{
		requestsTotal: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "smsgate_requests_total",
			Help: "Total number of SMS gateway API requests",
		}, []string{"operation", "status"}),
		requestDurationSeconds: promauto.NewHistogramVec(prometheus.HistogramOpts{
			Name:    "smsgate_request_duration_seconds",
			Help:    "Duration of SMS gateway API requests",
			Buckets: prometheus.DefBuckets,
		}, []string{"operation"}),
	}
}

func (m *Metrics) IncRequest(operation string, success bool) {
	if m == nil {
		return
	}
	status := metricFailure
	if success {
		status = metricSuccess
	}
	m.requestsTotal.WithLabelValues(operation, status).Inc()
}

func (m *Metrics) StartRequest(operation string) func() {
	if m == nil {
		return func() {}
	}
	start := time.Now()
	return func() {
		m.requestDurationSeconds.WithLabelValues(operation).Observe(time.Since(start).Seconds())
	}
}
