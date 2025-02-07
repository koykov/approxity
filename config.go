package bloom

type Policy uint

type Config struct {
	Size               uint64
	WriteAttemptsLimit uint64
	Hasher             Hasher
	HashChecksLimit    uint
}

func (c *Config) copy() *Config {
	cpy := *c
	return &cpy
}
