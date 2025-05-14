package heavy

type MetricsWriter interface {
	Add(err error) error
	Hits(min, max float64)
	Reset()
}

type DummyMetricsWriter struct{}

func (w *DummyMetricsWriter) Add(err error) error { return err }
func (w *DummyMetricsWriter) Hits(_, _ float64)   {}
func (w *DummyMetricsWriter) Reset()              {}
