package bloom

type ConcurrentFilter = Filter

func NewConcurrentFilter(config *Config) (*ConcurrentFilter, error) {
	if config == nil {
		return nil, ErrBadConfig
	}
	f := &ConcurrentFilter{
		conf: config.copy(),
		cncr: true,
	}
	f.once.Do(f.init)
	return f, nil
}
