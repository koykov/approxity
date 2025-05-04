package prometheus

import (
	"github.com/koykov/pbtk/frequency"
	"github.com/prometheus/client_golang/prometheus"
)

type mwFrequency struct {
	name string
	freq *prometheus.HistogramVec
}

func NewFrequency(name string) frequency.MetricsWriter {
	return &mwFrequency{name: name, freq: freqUniq}
}

// NewFrequencyWithBuckets creates a new frequency metrics writer with custom buckets.
// Param metricName must be unique!
func NewFrequencyWithBuckets(name, metricName, metricDesc string, buckets ...float64) (frequency.MetricsWriter, error) {
	metric := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    metricName,
		Help:    metricDesc,
		Buckets: buckets,
	}, []string{"name"})
	prometheus.MustRegister(metric)
	return &mwFrequency{name: name, freq: metric}, nil
}

func (mw *mwFrequency) Add(err error) error {
	result := "success"
	if err != nil {
		result = "fail"
	}
	freqAdd.WithLabelValues(mw.name, result).Inc()
	return err
}

func (mw *mwFrequency) Estimate(value uint64) uint64 {
	mw.freq.WithLabelValues(mw.name).Observe(float64(value))
	return value
}

func init() {
	freqAdd = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "frequency_add",
		Help: "Indicates how many times new items was set.",
	}, []string{"name", "result"})

	buckets := []float64{1, 10, 100, 1000, 1e4, 1e5, 1e6, 1e7, 1e8, 1e9}
	freqUniq = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "frequency_estimation",
		Help:    "Indicates estimations for buckets.",
		Buckets: buckets,
	}, []string{"name"})

	prometheus.MustRegister(freqAdd, freqUniq)
}

var (
	freqAdd  *prometheus.CounterVec
	freqUniq *prometheus.HistogramVec

	_, _ = NewFrequency, NewFrequencyWithBuckets
)
