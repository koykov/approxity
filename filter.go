package bloom

import (
	"sync"

	"github.com/koykov/bitvector"
)

type Filter struct {
	once sync.Once
	conf *Config
	vec  bitvector.Interface

	err error
}

func New(config *Config) *Filter {
	f := &Filter{
		conf: config.copy(),
	}
	f.once.Do(f.init)
	return f
}

func (f *Filter) init() {
	// todo implement me
}
