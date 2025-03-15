# Xor Filter

XOR filters are static AMQ filters that are based on a Bloomier filter and use the idea of perfect hash tables.
Similar to cuckoo filters, they save fingerprints of the elements in a hash table.

This package implements XorBinaryFuse8 variant as the fastest and memory-efficient solution.

## Usage

The minimal working example:
```go
import (
    "github.com/koykov/approxity/amq/xor_filter"
    "github.com/koykov/hash/xxhash"
)

func main() {
    f, err := xor.NewFilter[string](cuckoo.NewConfig(xxhash.Hasher64[[]byte]{}), []string{
	    "foobar", "wilson",	
    })
    _ = err
    print(f.Contains("foobar")) // true
    print(f.Contains("qwerty")) // false
}
```

Similar to [bloom filter](../bloom_filter/readme.md#usage) xor allows to initiate config more detailed:
```go
import "github.com/koykov/approxity/amq/metrics/prometheus"

func main() {
    // set items number and hasher
    f, _ := cuckoo.NewFilter[string](cuckoo.NewConfig(xxhash.Hasher64[[]byte]{}, []string{...}).
        // cover with metrics
        WithMetricsWriter(prometheus.NewPrometheusMetrics("example_filter")))
	...
```
