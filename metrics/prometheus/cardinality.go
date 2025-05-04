package prometheus

import (
	"github.com/koykov/pbtk/cardinality"
	"github.com/prometheus/client_golang/prometheus"
)

type mwCardinality struct {
	name string
	uniq *prometheus.HistogramVec
}

func NewCardinality(name string) cardinality.MetricsWriter {
	return &mwCardinality{name: name, uniq: cardUniq}
}

// NewCardinalityWithBuckets creates a new cardinality metrics writer with custom buckets.
// Param metricName must be unique!
func NewCardinalityWithBuckets(name, metricName, metricDesc string, buckets ...float64) (cardinality.MetricsWriter, error) {
	metric := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    metricName,
		Help:    metricDesc,
		Buckets: buckets,
	}, []string{"name"})
	prometheus.MustRegister(metric)
	return &mwCardinality{name: name, uniq: metric}, nil
}

func (mw *mwCardinality) Add(err error) error {
	result := "success"
	if err != nil {
		result = "fail"
	}
	cardAdd.WithLabelValues(mw.name, result).Inc()
	return err
}

func (mw *mwCardinality) Estimate(value uint64) uint64 {
	mw.uniq.WithLabelValues(mw.name).Observe(float64(value))
	return value
}

func init() {
	cardAdd = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "cardinality_add",
		Help: "Indicates how many times new items was set.",
	}, []string{"name", "result"})

	buckets := []float64{1, 10, 100, 1000, 1e4, 1e5, 1e6, 1e7, 1e8, 1e9}
	cardUniq = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "cardinality_unique",
		Help:    "Indicates estimations for buckets.",
		Buckets: buckets,
	}, []string{"name"})

	prometheus.MustRegister(cardAdd, cardUniq)
}

var (
	cardAdd  *prometheus.CounterVec
	cardUniq *prometheus.HistogramVec

	_, _ = NewCardinality, NewCardinalityWithBuckets
)
