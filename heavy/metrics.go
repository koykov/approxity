package heavy

type MetricsWriter interface {
	Add(error) error
	Hits()
	Reset()
}

type DummyMetricsWriter struct{}

func (w *DummyMetricsWriter) Add(error) error { return nil }
func (w *DummyMetricsWriter) Hits()           {}
func (w *DummyMetricsWriter) Reset()          {}
