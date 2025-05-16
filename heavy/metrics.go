package heavy

type MetricsWriter interface {
	Add(err error) error
	Hits(hits []Freq)
	Reset()
}

type Freq interface {
	Freq() float64
}

type DummyMetricsWriter struct{}

func (w *DummyMetricsWriter) Add(err error) error { return err }
func (w *DummyMetricsWriter) Hits(_ []Freq)       {}
func (w *DummyMetricsWriter) Reset()              {}
