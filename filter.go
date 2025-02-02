package bloom

import "sync"

type Filter struct {
	once sync.Once
	conf *Config

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
	// todo calculate number of hashing steps
}
