package bloom

type Access uint

const (
	AccessReadWrite Access = iota
	AccessExclusiveReadOrWrite
)

type Config struct {
	Size       uint64
	Access     Access
	Hasher     Hasher
	HashChecks uint
}

func (c *Config) copy() *Config {
	cpy := *c
	return &cpy
}
