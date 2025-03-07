package cardinality

type MetricsWriter interface {
	Add(error) error
	Estimate(uint64) uint64
}
