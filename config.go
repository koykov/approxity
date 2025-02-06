package bloom

type Policy uint

const (
	PolicySimultaneousReadWrite Policy = iota
	PolicyExclusiveReadOrWrite
)

type Access uint

const (
	AccessReadWrite Access = iota
	AccessRead
	AccessWrite
)

type Config struct {
	Size       uint64
	Policy     Policy
	WriteLimit uint64
	Hasher     Hasher
	HashChecks uint
}

func (c *Config) copy() *Config {
	cpy := *c
	return &cpy
}
