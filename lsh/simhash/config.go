package simhash

import "github.com/koykov/pbtk"

type Config struct {
	// Hasher to calculate hash sum of the items.
	// Mandatory param.
	Hasher pbtk.Hasher
}

func NewConfig(hasher pbtk.Hasher) *Config {
	return &Config{
		Hasher: hasher,
	}
}
