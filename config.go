package bloom

type Policy uint

type Config struct {
	Concurrent struct {
		WriteAttemptsLimit uint64
	}
	Size            uint64
	Hasher          Hasher
	HashChecksLimit uint64
}

func (c *Config) copy() *Config {
	cpy := *c
	return &cpy
}
