# XOR Filter

The XOR Filter is a probabilistic data structure designed for efficient membership testing.
It serves as a compact alternative to Bloom filters, offering higher performance and lower false-positive rates.

The XOR Filter requires pre-construction on a fixed set of keys. The construction algorithm involves:
* Creating a dependency graph between keys
* Solving a system of equations where each key is associated with multiple positions in a bit array
* Populating the array such that XOR-ing the selected positions yields the key's hash

Adding new keys requires complete reconstruction of the structure as it would alter all existing dependencies.

## Implementation Features

* Implements XorBinaryFuse8
* Supports custom hash functions
* Filter reuse via `sync.Pool` to minimize allocations
* Lock-free operation
* Compact data storage

## Advantages

1. **Smaller size** compared to Bloom filters
2. **Fast lookups** — only 3 memory accesses and 2 XOR operations
3. **No false negatives** — if a key was added, it will always be found
4. **Low false-positive probability** — approximately 0.4%
5. **Concurrency** — implementation uses no locks

## Math basics

The XOR Filter is constructed in three stages:

1. **Key Distribution**:
   For each key $x$, positions are computed:

   $$
   h_1(x), h_2(x), h_3(x) \in \{0, \ldots, m-1\}
   $$

   where $m$ is the filter size, typically $m \approx 1.23 \cdot n$ for $n$ keys.

2. **Equation System Construction**:
   Each key $x$ corresponds to an equation:

   $$
   \text{fingerprint}(x) = \text{filter}[h_1(x)] \oplus \text{filter}[h_2(x)] \oplus \text{filter}[h_3(x)]
   $$
   
   where $\text{fingerprint}(x)$ is an 8-bit key hash.

3. **System Solution**:
   The system is solved using Gaussian elimination, guaranteeing:
 
   $$
   \text{Pr}(\text{false positive}) \leq \frac{1}{2^8} = \frac{1}{256} \approx 0.39\%
   $$

Membership testing for key $y$ is performed as:

$$
\text{Contains}(y) = \left(\text{filter}[h_1(y)] \oplus \text{filter}[h_2(y)] \oplus \text{filter}[h_3(y)]\right) == \text{fingerprint}(y)
$$

## Usage

```go
package main

import (
	"github.com/koykov/hash/xxhash"
	xor "github.com/koykov/pbtk/amq/xor_filter"
	"github.com/koykov/pbtk/metrics/prometheus"
)

func main() {
	hasher := xxhash.Hasher64[[]byte]{} // hash function
	config := xor.NewConfig(hasher).
		WithMetricsWriter(prometheus.NewAMQ("example_filter")) // cover with metrics
	f, err := xor.NewFilterWithKeys[[]byte](config, [][]byte{
		[]byte("foo"),
		[]byte("bar"),
	})
	_ = err

	println(f.Contains([]byte("foo"))) // true
	println(f.Contains([]byte("qwe"))) // false
}
```

## Use Cases

1. **Caching** — quick existence checks before expensive queries
2. **Databases** — query pre-filtering to reduce disk accesses
3. **Network filters** — blocking unwanted IPs or URLs
4. **Duplicate detection** — checking for duplicates in data streams
5. **Search engines** — document pre-selection for full-text search

## Conclusion

Xor Filter provides an efficient and compact method for membership testing with minimal allocations and high concurrent performance.
The filter is particularly useful in scenarios where lookup speed and memory efficiency are crucial, and the inability to dynamically add elements is an acceptable constraint.
