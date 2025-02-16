package cuckoo

import "github.com/koykov/amq"

const (
	defaultBucketSize        = 4
	defaultKicksLimit        = 500
	defaultFalsePositiveRate = 0.01
	defaultSeed              = 2077
)

type Config struct {
	Size            uint64
	BucketSize      uint64
	FingerprintSize uint64
	KicksLimit      uint64
	Hasher          amq.Hasher
	Seed            uint64
}

func (c *Config) copy() *Config {
	cpy := *c
	return &cpy
}
