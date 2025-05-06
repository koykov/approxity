# Cuckoo Filter

Cuckoo Filter is a probabilistic data structure for membership testing that combines the advantages of Cuckoo Hashing
and Bloom filters. Unlike Bloom filters, Cuckoo Filter supports item deletion without losing accuracy.

## Implementation Features

* Support for custom hash functions
* Optional concurrency mode - safe asynchronous read/write operations
* Lock-free design through atomic operations
* Use of SIMD instructions to accelerate operations
* Fixed bucket size (4 elements)
* Automatic table size calculation based on expected element count and kicks limit

## Cuckoo Hashing Principle

* Each element can reside in one of two positions determined by hash functions
* Each position (bucket) contains several fingerprints of the element
* On collision, existing fingerprints get kicked to their alternate position

## Mathematical Foundations

### Optimal Table Size Formula

The `optimalM` function calculates the optimal table size (number of buckets) for a given maximum element count `n`:

$$
m = \frac{2^{\lceil \log_2(n) \rceil}}{b}
$$

Where:
- `b` - bucket size (fixed at 4 in this implementation)
- `⌈log₂(n)⌉` - rounded up to nearest power of two

This formula ensures:
* Table size is always a power of two, enabling fast bitwise operations instead of expensive divisions
* Approximately 95% load factor with bucket size of 4
* Minimized probability of infinite kicking loops

### False Positive Probability

The false positive probability for Cuckoo Filter is calculated as:

$$
\epsilon \approx \frac{2b}{2^f}
$$

Where:
- `f` - fingerprint length in bits
- `b` - bucket size

## Size Calculation Example

For `n = 1000`:
* Find nearest power of two: `2^10 = 1024`
* Divide by bucket size (4): `1024 / 4 = 256`
* Final table size: 256 buckets

## Usage Example

```go
package main

import (
	"github.com/koykov/hash/xxhash"
	"github.com/koykov/pbtk/amq/cuckoo_filter"
	"github.com/koykov/pbtk/metrics/prometheus"
)

const N = 1e7

func main() {
	hasher := xxhash.Hasher64[[]byte]{} // hash function
	config := cuckoo.NewConfig(N, hasher).
		WithKicksLimit(10).                                    // limit for cuckoo kicks to avoid infinite loop
		WithConcurrency().                                     // switch to race protected buckets array (atomic based)
		WithMetricsWriter(prometheus.NewAMQ("example_filter")) // cover with metrics
	f, err := cuckoo.NewFilter[string](config)
	_ = err
	_ = f.Set("foobar")
	print(f.Contains("foobar")) // true
	print(f.Contains("qwerty")) // false
}
```

## Applications

* Caching - quick membership checks
* Network applications - packet deduplication
* Databases - query acceleration by filtering unpromising queries
* Monitoring systems - tracking unique events
* Distributed systems - data deduplication

## Conclusion

This implementation provides an efficient tool for probabilistic data storage with concurrent access support.
The use of atomic operations and SIMD instructions enables high performance, while mathematically grounded parameter
selection ensures optimal memory usage.

The filter is particularly useful in scenarios requiring frequent data updates under memory constraints and high load conditions.
