package dlcsketch

import (
	"github.com/koykov/pbtk"
	"github.com/koykov/pbtk/frequency"
)

type Config struct {
	N uint64 // desired number of unique items
	M uint64 // size of tables
	D uint64 // number of hashes
	C uint64 // collision probability
	// Hasher to calculate hash sum of the items.
	// Mandatory param.
	Hasher pbtk.Hasher
	// Enable compact mode.
	// By default, uses 64 bit per counter. This param allows to use 32 bit per counter.
	Compact bool
	// Setting up this section enables concurrent read/write operations.
	Concurrent *ConcurrentConfig
	// Metrics writer handler.
	MetricsWriter frequency.MetricsWriter
}

// ConcurrentConfig configures concurrent section of config.
type ConcurrentConfig struct {
	// How many write attempts may perform.
	WriteAttemptsLimit uint64
}
