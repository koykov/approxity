package frequency

type MetricsWriter interface {
	Add(error) error
	Estimate(uint64) uint64
}

type SignedMetricsWriter interface {
	Add(error) error
	Estimate(int64) int64
}

type PreciseMetricsWriter interface {
	Add(error) error
	Estimate(float64) float64
}

type DummyMetricsWriter struct{}

func (w DummyMetricsWriter) Add(err error) error      { return err }
func (w DummyMetricsWriter) Estimate(n uint64) uint64 { return n }

type DummySignedMetricsWriter struct{}

func (w DummySignedMetricsWriter) Add(err error) error    { return err }
func (w DummySignedMetricsWriter) Estimate(n int64) int64 { return n }

type DummyPreciseMetricsWriter struct{}

func (w DummyPreciseMetricsWriter) Add(err error) error        { return err }
func (w DummyPreciseMetricsWriter) Estimate(n float64) float64 { return n }
