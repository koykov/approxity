package amq

type MetricsWriter interface {
	Set()
	Unset()
	Contains(positive bool)
}

type DummyMetricsWriter struct{}

func (DummyMetricsWriter) Set()            {}
func (DummyMetricsWriter) Unset()          {}
func (DummyMetricsWriter) Contains(_ bool) {}
