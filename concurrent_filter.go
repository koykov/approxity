package bloom

// ConcurrentFilter is a concurrent implementation of Bloom filter.
type ConcurrentFilter = Filter

// NewConcurrentFilter creates new concurrent Bloom filter.
func NewConcurrentFilter(config *Config) (*ConcurrentFilter, error) {
	if config == nil {
		return nil, ErrBadConfig
	}
	f := &ConcurrentFilter{
		conf: config.copy(),
		cncr: true, // The main difference from non-concurrent implementation.
	}
	f.once.Do(f.init)
	return f, nil
}
