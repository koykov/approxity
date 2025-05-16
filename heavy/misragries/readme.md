# Misra-Gries

Misra-Gries is a probabilistic algorithm for finding **frequent items** (heavy hitters) in a data stream.  
It solves the problem of **approximate counting** of elements that occur no less than $n/k$ times, where:
- $n$ is the total number of elements in the stream,
- $k$ is a parameter that determines accuracy and memory consumption.

The algorithm uses $O(k)$ memory and guarantees that all elements with frequency $\geq n/k$ will be found
(though it may also return some elements with lower frequency).

## Implementation Features

* Support for custom hash functions
* Sharding for multi-threaded environments
* Built-in metrics coverage

## Usage

```go
package main

import (
	"fmt"

	"github.com/koykov/hash/xxhash"
	"github.com/koykov/pbtk/heavy/misragries"
)

const K = 5

func main() {
	hasher := xxhash.Hasher64[[]byte]{}
	hitter, err := misragries.NewHitter[string](misragries.NewConfig(K, hasher))
	_ = err
	oftenKey := "key0"
	for i := 0; i < 1000; i++ {
		key := fmt.Sprintf("key%d", i)
		_ = hitter.Add(key)
		if i%2 == 0 {
			_ = hitter.Add(oftenKey)
		}
	}
	hits := hitter.Hits()
	fmt.Printf("%+v", hits) // [{Key:key0 Rate:242} {Key:key990 Rate:1} {Key:key992 Rate:1} {Key:key999 Rate:1} {Key:key996 Rate:1}]
}
```

## Use Cases

The Misra-Gries algorithm can be used in:
1. **Network traffic analysis** - detecting frequent IP addresses or URLs.
2. **Log processing** - identifying popular queries or errors.
3. **Recommendation systems** - discovering trending items or content.
4. **Stream analytics** - real-time anomaly monitoring.

## Conclusion

This implementation provides:

* Efficient detection of frequent items in data streams
* Multi-threading support through sharding
* Flexibility with custom hash functions
* Ready-to-use metrics for performance evaluation

The algorithm is particularly useful in scenarios where speed and memory efficiency are crucial, and some margin of error is acceptable.
