package bloom

type Config struct {
	Concurrent *struct {
		// How many write attempts may perform.
		WriteAttemptsLimit uint64
	}
	// The size of the filter in bits.
	// Mandatory param.
	Size uint64
	// Hasher to calculate hash sum of the items.
	Hasher Hasher
	// How many hash checks filter may do to reduce false positives cases.
	HashChecksLimit uint64
}

func (c *Config) copy() *Config {
	cpy := *c
	return &cpy
}
