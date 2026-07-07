package sessions

import (
	"fmt"
	"time"

	"github.com/fiorix/go-smpp/v2/smpp/pdu"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	labelStatus = "status"

	metricSuccess = "success"
	metricFailure = "failure"
)

type Metrics struct {
	sessionsActive    *prometheus.GaugeVec
	pdusReceivedTotal *prometheus.CounterVec
	pdusSentTotal     *prometheus.CounterVec
	bindAttemptsTotal *prometheus.CounterVec
	submitSmTotal     *prometheus.CounterVec
	submitSmDuration  prometheus.Histogram
	querySmTotal      *prometheus.CounterVec
	querySmDuration   prometheus.Histogram
	errorsTotal       *prometheus.CounterVec
}

func NewMetrics() *Metrics {
	return &Metrics{
		sessionsActive: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Name: "smpp_sessions_active",
			Help: "Current number of active SMPP sessions by state",
		}, []string{"state"}),
		pdusReceivedTotal: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "smpp_pdus_received_total",
			Help: "Total number of PDU requests received",
		}, []string{"command"}),
		pdusSentTotal: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "smpp_pdus_sent_total",
			Help: "Total number of PDU responses sent",
		}, []string{"command"}),
		bindAttemptsTotal: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "smpp_bind_attempts_total",
			Help: "Total number of bind attempts",
		}, []string{"type", labelStatus}),
		submitSmTotal: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "smpp_submit_sm_total",
			Help: "Total number of submit_sm operations",
		}, []string{labelStatus}),
		submitSmDuration: promauto.NewHistogram(prometheus.HistogramOpts{
			Name:    "smpp_submit_sm_duration_seconds",
			Help:    "Duration of submit_sm operations",
			Buckets: prometheus.DefBuckets,
		}),
		querySmTotal: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "smpp_query_sm_total",
			Help: "Total number of query_sm operations",
		}, []string{labelStatus}),
		querySmDuration: promauto.NewHistogram(prometheus.HistogramOpts{
			Name:    "smpp_query_sm_duration_seconds",
			Help:    "Duration of query_sm operations",
			Buckets: prometheus.DefBuckets,
		}),
		errorsTotal: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "smpp_errors_total",
			Help: "Total number of SMPP error responses by error code",
		}, []string{"code"}),
	}
}

func (m *Metrics) IncSessionsActive(s state) {
	if m == nil {
		return
	}
	m.sessionsActive.WithLabelValues(stateLabel(s)).Inc()
}

func (m *Metrics) DecSessionsActive(s state) {
	if m == nil {
		return
	}
	m.sessionsActive.WithLabelValues(stateLabel(s)).Dec()
}

func (m *Metrics) IncPDUsReceived(id pdu.ID) {
	if m == nil {
		return
	}
	m.pdusReceivedTotal.WithLabelValues(commandName(id)).Inc()
}

func (m *Metrics) IncPDUsSent(id pdu.ID) {
	if m == nil {
		return
	}
	m.pdusSentTotal.WithLabelValues(commandName(id)).Inc()
}

func (m *Metrics) IncBindAttempts(receiver, transceiverMode bool, success bool) {
	if m == nil {
		return
	}
	status := metricFailure
	if success {
		status = metricSuccess
	}
	m.bindAttemptsTotal.WithLabelValues(bindTypeLabel(receiver, transceiverMode), status).Inc()
}

func (m *Metrics) IncSubmitSM(success bool) {
	if m == nil {
		return
	}
	status := metricFailure
	if success {
		status = metricSuccess
	}
	m.submitSmTotal.WithLabelValues(status).Inc()
}

func (m *Metrics) StartSubmitSM() func() {
	if m == nil {
		return func() {}
	}
	start := time.Now()
	return func() {
		m.submitSmDuration.Observe(time.Since(start).Seconds())
	}
}

func (m *Metrics) IncQuerySM(success bool) {
	if m == nil {
		return
	}
	status := metricFailure
	if success {
		status = metricSuccess
	}
	m.querySmTotal.WithLabelValues(status).Inc()
}

func (m *Metrics) StartQuerySM() func() {
	if m == nil {
		return func() {}
	}
	start := time.Now()
	return func() {
		m.querySmDuration.Observe(time.Since(start).Seconds())
	}
}

func (m *Metrics) IncError(code uint32) {
	if m == nil {
		return
	}
	m.errorsTotal.WithLabelValues(fmt.Sprintf("0x%08X", code)).Inc()
}

func commandName(id pdu.ID) string {
	if name := id.String(); name != "" {
		return name
	}
	return fmt.Sprintf("0x%08X", uint32(id))
}

func stateLabel(s state) string {
	switch s {
	case StateOpen:
		return "open"
	case StateBoundTX:
		return "bound_tx"
	case StateBoundRX:
		return "bound_rx"
	case StateBoundTRX:
		return "bound_trx"
	default:
		return "unknown"
	}
}

func bindTypeLabel(receiver, transceiverMode bool) string {
	switch {
	case transceiverMode:
		return "trx"
	case receiver:
		return "rx"
	default:
		return "tx"
	}
}
