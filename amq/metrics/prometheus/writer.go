package prometheus

import (
	"github.com/koykov/pbtk/amq"
	"github.com/prometheus/client_golang/prometheus"
)

type metricsWriter struct {
	name string
}

func NewMetricsWriter(name string) amq.MetricsWriter {
	return &metricsWriter{name: name}
}

func (mw *metricsWriter) Capacity(cap uint64) {
	mcap.WithLabelValues(mw.name).Set(float64(cap))
}

func (mw *metricsWriter) Set(err error) error {
	result := "success"
	if err != nil {
		result = "fail"
	}
	mset.WithLabelValues(mw.name, result).Inc()
	return err
}

func (mw *metricsWriter) Unset(err error) error {
	result := "success"
	if err != nil {
		result = "fail"
	}
	munset.WithLabelValues(mw.name, result).Inc()
	return err
}

func (mw *metricsWriter) Contains(positive bool) bool {
	result := "positive"
	if !positive {
		result = "negative"
	}
	mcontains.WithLabelValues(mw.name, result).Inc()
	return positive
}

func (mw *metricsWriter) Reset() {}

func init() {
	mcap = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "amq_capacity",
		Help: "Indicates how many items filter may contain.",
	}, []string{"name"})

	mset = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "amq_set",
		Help: "Indicates how many times new items was set to the filter.",
	}, []string{"name", "result"})

	munset = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "amq_unset",
		Help: "Indicates how many times an items was unset from the filter.",
	}, []string{"name", "result"})

	mcontains = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "amq_contains",
		Help: "Indicates how many times filter was checked and check result (positive/negative).",
	}, []string{"name", "result"})
}

var (
	mcap                    *prometheus.GaugeVec
	mset, munset, mcontains *prometheus.CounterVec

	_ = NewMetricsWriter
)
