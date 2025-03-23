package countminsketch

import "io"

type cnvector[T ~uint32 | ~uint64] struct {
	d, w, lim uint64
	buf       []T
}

func (vec *cnvector[T]) add(uint64) error {
	// todo implement me
	return nil
}

func (vec *cnvector[T]) estimate(uint64) uint64 {
	// todo implement me
	return 0
}

func (vec *cnvector[T]) reset() {
	// todo implement me
}

func (vec *cnvector[T]) readFrom(io.Reader) (n int64, err error) {
	// todo implement me
	return
}

func (vec *cnvector[T]) writeTo(io.Writer) (n int64, err error) {
	// todo implement me
	return
}

func newConcurrentVector[T ~uint32 | ~uint64](d, w, lim uint64) *cnvector[T] {
	return &cnvector[T]{
		d:   d,
		w:   w,
		lim: lim,
		buf: make([]T, d*w),
	}
}
