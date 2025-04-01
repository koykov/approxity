package tinylfu

import (
	"io"

	"github.com/koykov/pbtk/frequency"
	"github.com/koykov/pbtk/frequency/cmsketch"
)

func (c *Config) WithConcurrency() *Config {
	c.Concurrent = &cmsketch.ConcurrentConfig{}
	return c
}

func (c *Config) WithWriteAttemptsLimit(limit uint64) *Config {
	if c.Concurrent == nil {
		c.Concurrent = &cmsketch.ConcurrentConfig{}
	}
	c.Concurrent.WriteAttemptsLimit = limit
	return c
}

func (c *Config) WithCompact() *Config {
	c.Compact = true
	return c
}

func WithMetricsWriter(conf *Config, mw frequency.MetricsWriter) *Config {
	conf.MetricsWriter = mw
	return conf
}

func (e *estimator[T]) Estimate(key T) uint64               { return e.est.Estimate(key) }
func (e *estimator[T]) HEstimate(hkey uint64) uint64        { return e.est.HEstimate(hkey) }
func (e *estimator[T]) Reset()                              { e.est.Reset() }
func (e *estimator[T]) ReadFrom(r io.Reader) (int64, error) { return e.est.ReadFrom(r) }
func (e *estimator[T]) WriteTo(w io.Writer) (int64, error)  { return e.est.WriteTo(w) }
