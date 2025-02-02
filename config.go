package bloom

type Config struct {
	Size       uint64
	Buckets    uint
	Hasher     Hasher
	HashChecks uint
}

func (c *Config) copy() *Config {
	cpy := *c
	return &cpy
}
