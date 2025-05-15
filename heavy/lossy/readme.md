# Lossy Counting

Lossy Counting is a probabilistic algorithm for identifying frequent items in a data stream with controlled error margin.  
The algorithm solves the problem of finding items whose frequency exceeds a given minimum threshold $s$ (support), while guaranteeing that:
- All truly frequent items will be detected
- Each item's frequency may be underestimated by at most $εN$, where $ε$ is the allowable error margin

The key feature of this implementation is automatic control over the size of hash table $D$ that stores items and their counters.  
The algorithm maintains $D$'s size within $O(1/ε)$ bounds by periodically removing items satisfying the condition $f + Δ ≤ b$, where:
- $f$ — current frequency estimate of an item
- $Δ$ — maximum frequency error
- $b$ — current data block number

## Implementation Features

* Support for custom hash functions
* Sharding for multi-threaded environments
* Built-in metrics coverage

## Mathematical Foundations

The algorithm is based on the following theoretical principles:

1. **Accuracy guarantees**:
    - All items with frequency $f ≥ sN$ will be found
    - For each found item $f ≥ (s - ε)N$

2. **Memory estimation**:  
   The size of table $D$ is bounded by:

$$
|D| \leq \frac{1}{\varepsilon} \log(\varepsilon N)
$$

3. **Algorithm parameters**:
    - $ε$ — allowable error margin (must be less than $s$)
    - $w = \lceil 1/ε \rceil$ — processing block size

## Usage Example

```go
package main

import (
	"fmt"

	"github.com/koykov/hash/xxhash"
	"github.com/koykov/pbtk/heavy/lossy"
)

const (
	epsilon = 0.01
	support = 0.02
)

func main() {
	hasher := xxhash.Hasher64[[]byte]{}
	hitter, err := lossy.NewHitter[string](lossy.NewConfig(epsilon, support, hasher))
	_ = err
	oftenKey := "key0"
	for i := 0; i < 1000; i++ {
		key := fmt.Sprintf("key%d", i)
		_ = hitter.Add(key)
		if i%5 == 0 {
			_ = hitter.Add(oftenKey)
		}
	}
	hits := hitter.Hits()
	fmt.Printf("%+v", hits) // [{Key:key0 Rate:0.45318965517241383} {Key:key325 Rate:0.024132231404958678} ... {Key:key880 Rate:0.024132231404958678}]
}
```

## Application Areas

* Log analysis: detecting frequent errors or request patterns
* Recommendation systems: identifying popular items/content
* Network monitoring: DDoS attack detection (abnormally frequent requests)

## Conclusion

This Golang implementation of Lossy Counting provides:

* An efficient streaming algorithm for frequent item detection
* Scalability in multi-threaded environments
* Memory usage control
* Ready-to-use monitoring system integration

The algorithm is particularly useful in stream processing systems where real-time identification of popular items is required with guaranteed accuracy under resource constraints.
