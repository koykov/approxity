package cardinality

type MetricsWriter interface {
	Add(error) error
	Estimate(uint64) uint64
}

type DummyMetricsWriter struct{}

func (w DummyMetricsWriter) Add(err error) error      { return err }
func (w DummyMetricsWriter) Estimate(n uint64) uint64 { return n }
