package amq

type MetricsWriter interface {
	Capacity(cap uint64)
	Set(err error) error
	Unset(err error) error
	Contains(positive bool) bool
	Reset()
}

type DummyMetricsWriter struct{}

func (DummyMetricsWriter) Capacity(cap uint64)         {}
func (DummyMetricsWriter) Set(err error) error         { return err }
func (DummyMetricsWriter) Unset(err error) error       { return err }
func (DummyMetricsWriter) Contains(positive bool) bool { return positive }
func (DummyMetricsWriter) Reset()                      {}
