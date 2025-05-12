package heavy

type MetricsWriter interface {
	Add(error) error
	Hits()
}

type DummyMetricsWriter struct{}

func (w *DummyMetricsWriter) Add(error) error { return nil }
func (w *DummyMetricsWriter) Hits()           {}
