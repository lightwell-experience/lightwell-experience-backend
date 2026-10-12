package instrumentation

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	NameSpace           = "lightwell_experience"
	HTTPStatusHistogram = "http_status_histogram"
)

type Metrics struct {
	HTTPStatusHistogram prometheus.HistogramVec
	reg                 *prometheus.Registry
}

// https://prometheus.io/docs/tutorials/understanding_metric_types/#types-of-metrics
func NewMetrics(reg *prometheus.Registry) *Metrics {
	if reg == nil {
		panic("reg cannot be nil")
	}

	metrics := &Metrics{
		reg: reg,
		HTTPStatusHistogram: *promauto.With(reg).NewHistogramVec(prometheus.HistogramOpts{
			Namespace: NameSpace,
			Name:      HTTPStatusHistogram,
			Help:      "Duration of HTTP requests",
			Buckets:   prometheus.DefBuckets,
		}, []string{"status", "method", "path"}),
	}

	reg.MustRegister(collectors.NewBuildInfoCollector())

	return metrics
}

func (m *Metrics) Registry() *prometheus.Registry {
	return m.reg
}
