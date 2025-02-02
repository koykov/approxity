package bloom

import "sync"

type Filter struct {
	Size, Expected uint64

	once           sync.Once
	size, expected uint64
}

func New(size, expected uint64) *Filter {
	f := &Filter{
		Size:     size,
		Expected: expected,
	}
	f.once.Do(f.init)
	return f
}

func (f *Filter) init() {
	f.size, f.expected = f.Size, f.Expected
}
