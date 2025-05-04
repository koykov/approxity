package prometheus

import (
	"github.com/koykov/pbtk/frequency"
	"github.com/prometheus/client_golang/prometheus"
)

type mwSignedFrequency struct {
	mwFrequency
}

func (mw *mwSignedFrequency) Estimate(value int64) int64 {
	mw.freq.WithLabelValues(mw.name).Observe(float64(value))
	return value
}

func NewSignedFrequency(name string) frequency.SignedMetricsWriter {
	return &mwSignedFrequency{mwFrequency{name: name, freq: freqUniq}}
}

// NewSignedFrequencyWithBuckets creates a new frequency metrics writer with custom buckets.
// Param metricName must be unique!
func NewSignedFrequencyWithBuckets(name, metricName, metricDesc string, buckets ...float64) (frequency.SignedMetricsWriter, error) {
	metric := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    metricName,
		Help:    metricDesc,
		Buckets: buckets,
	}, []string{"name"})
	prometheus.MustRegister(metric)
	return &mwSignedFrequency{mwFrequency{name: name, freq: metric}}, nil
}

var _, _ = NewSignedFrequency, NewSignedFrequencyWithBuckets
