package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

var (
	ExecutionsCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "codeexec_executions_total",
			Help: "Amount of executions",
		},
		[]string{"language"},
	)

	ExecutionsDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "codeexec_execution_duration_seconds",
			Help:    "Execution duration in seconds",
			Buckets: []float64{1, 3, 10, 30},
		},
		[]string{"language"},
	)

	IndexPageCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "codeexec_index_page_served",
			Help: "Total number of index page served",
		},
	)
)

func InitMetrics() {
	log.Info("Registering metrics")
	prometheus.MustRegister(ExecutionsCounter)
	prometheus.MustRegister(ExecutionsDuration)
	prometheus.MustRegister(IndexPageCounter)
}
