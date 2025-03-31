package tinylfu

import "io"

func (e *estimator[T]) Estimate(key T) uint64               { return e.est.Estimate(key) }
func (e *estimator[T]) HEstimate(hkey uint64) uint64        { return e.est.HEstimate(hkey) }
func (e *estimator[T]) Reset()                              { e.est.Reset() }
func (e *estimator[T]) ReadFrom(r io.Reader) (int64, error) { return e.est.ReadFrom(r) }
func (e *estimator[T]) WriteTo(w io.Writer) (int64, error)  { return e.est.WriteTo(w) }
