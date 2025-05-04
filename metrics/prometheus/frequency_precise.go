package prometheus

import (
	"github.com/koykov/pbtk/frequency"
	"github.com/prometheus/client_golang/prometheus"
)

type mwPreciseFrequency struct {
	mwFrequency
}

func (mw *mwPreciseFrequency) Estimate(value float64) float64 {
	mw.freq.WithLabelValues(mw.name).Observe(value)
	return value
}

func NewPreciseFrequency(name string) frequency.PreciseMetricsWriter {
	return &mwPreciseFrequency{mwFrequency{name: name, freq: freqUniq}}
}

// NewPreciseFrequencyWithBuckets creates a new frequency metrics writer with custom buckets.
// Param metricName must be unique!
func NewPreciseFrequencyWithBuckets(name, metricName, metricDesc string, buckets ...float64) (frequency.PreciseMetricsWriter, error) {
	metric := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    metricName,
		Help:    metricDesc,
		Buckets: buckets,
	}, []string{"name"})
	prometheus.MustRegister(metric)
	return &mwPreciseFrequency{mwFrequency{name: name, freq: metric}}, nil
}

var _, _ = NewPreciseFrequency, NewPreciseFrequencyWithBuckets
