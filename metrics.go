package amq

type MetricsWriter interface {
	Set(err error) error
	Unset(err error) error
	Contains(positive bool) bool
}

type DummyMetricsWriter struct{}

func (DummyMetricsWriter) Set(err error) error         { return err }
func (DummyMetricsWriter) Unset(err error) error       { return err }
func (DummyMetricsWriter) Contains(positive bool) bool { return positive }
