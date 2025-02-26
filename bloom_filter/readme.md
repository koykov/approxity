# Bloom filter

Bloom filter is a good old implementation of AMQ, presented at 1970s. The filter is a bit array with length `m` and `k`
hash functions. It is extremely space efficient and is typically used to add elements to a set and test if an element is
in a set. Though, the elements themselves are not added to a set. A hash of the elements is added to the set instead:
```
pos = hash(element) % m
bitArray[pos] = 1
```

See [full description](https://en.wikipedia.org/wiki/Bloom_filter).

## Usage

The minimal working example:
```go
import (
    "github.com/koykov/amq/bloom_filter"
    "github.com/koykov/hash/xxhash"
)

func main() {
    f, err := bloom.NewFilter(bloom.NewConfig(1000, 0.01, xxhash.Hasher64[[]byte]{}))
    _ = err
    _ = f.Set("foobar")
    print(f.Contains("foobar")) // true
    print(f.Contains("qwerty")) // false
}
```
, but [initial config](config.go) allows to tune filter for better efficiency:
```go
import "github.com/koykov/amq/metrics/prometheus"

func func main() {
    // set filter size and hasher
    config := bloom.NewConfig(1000, 0.01, xxhash.Hasher64[[]byte]{Seed: 1234}).
        // switch to race protected bit array (atomic based)
        WithConcurrency().
        // cover with metrics
        WithMetricsWriter(prometheus.NewPrometheusMetrics("example_filter"))
    
    // filter is ready to use
    f, _ := NewFilter(config)
	...
}
```

### Optimal params calculation

There is no need to calculate optimal size `m` and number of hash functions `k` due to filter makes it itself using
desired number of items (`Config.ItemsNumber`) and false positive probability (`Config.FPP`) params.
