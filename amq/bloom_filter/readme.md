# Bloom Filter

A Bloom filter is a probabilistic data structure designed to test whether an element is a member of a set.
The filter may yield false positives (claiming an element is present when it's not), but it never produces false negatives.

The key advantage over hash tables is its compact size (see comparative example below).

## Implementation Features

* Customizable hash function
* Concurrent operation support - parallel reads and writes
* Lock-free implementation (using only atomic operations)
* SIMD operations where applicable
* Counting Bloom Filter support (storage abstraction)

## Math basics

### Optimal Parameters Calculation

The bit array size $m$ and number of hash functions $k$ are calculated based on:
- $N$ — expected maximum number of elements
- $FPP$ — desired false positive probability

Optimal bit array size $m$:

$$
m = -\frac{N \cdot \ln(FPP)}{(\ln 2)^2}
$$

Optimal number of hash functions $k$:

$$
k = \frac{m}{N} \ln 2
$$

### Calculation Example

For $N = 1,000,000$ elements and $FPP = 0.01$ (1%):

1. Calculate $m$:
   $$
   m = -\frac{1,000,000 \cdot \ln(0.01)}{(\ln 2)^2} \approx -\frac{1,000,000 \cdot (-4.605)}{0.4805} \approx 9,583,000 \text{ bits} \approx 1.14 \text{ MB}
   $$

2. Calculate $k$:
   $$
   k = \frac{9,583,000}{1,000,000} \cdot 0.693 \approx 6.64 \approx 7 \text{ hash functions}
   $$

For comparison, a hash table storing 8-byte keys would require:
$$
1,000,000 \times 8 \text{ bytes} = 7.63 \text{ MB}
$$
for keys alone. The actual size would be larger due to additional data structures (buckets) and load factor (+30-50%).

## Usage Example

```go
package main

import (
	"github.com/koykov/pbtk/amq/bloom_filter"
	"github.com/koykov/pbtk/metrics/prometheus"
	"github.com/koykov/hash/xxhash"
)

const (
	N = 1000
	FPP = 0.01
)

func main() {
	hasher := xxhash.Hasher64[[]byte]{} // hash function
	config := bloom.NewConfig[string](N, FPP, hasher).
		WithConcurrency(). // switch to race protected bit array (atomic based)
		WithMetricsWriter(prometheus.NewAMQ("example_filter")) // cover with metrics
	// config.WithCBF() // switch to counting bloom filter
	f, err := bloom.NewFilter[string](config)
	_ = err
	_ = f.Set("foobar")
	print(f.Contains("foobar")) // true
	print(f.Contains("qwerty")) // false
}
```

## Applications

1. **Caching**: Quick existence checks before expensive operations
2. **Databases**: Pre-filtering disk queries
3. **Network Technologies**: URL verification in web crawlers
4. **Security Systems**: Banned password/token checks
5. **Blockchain**: Transaction verification optimization

## References

1. [Bloom filter](https://en.wikipedia.org/wiki/Bloom_filter)
2. [Counting Bloom filter](https://en.wikipedia.org/wiki/Counting_Bloom_filter)

## Conclusion

The Bloom filter provides an efficient trade-off between memory usage and error probability.
For scenarios where false positives are acceptable, it can reduce memory consumption by 10-30x compared to traditional data structures.
The atomic operation implementation enables efficient usage in concurrent environments without complex synchronization mechanisms.
