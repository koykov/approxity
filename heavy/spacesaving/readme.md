# Space-Saving Stream Summary

The Space-Saving algorithm solves the problem of finding the most frequent elements (heavy hitters) in a data stream while using limited memory.

## Key Features

* Suitable for **stream processing** (continuous data flow)
* Guarantees detection of high-frequency elements
* Uses a **fixed number of counters** (doesn't store the entire stream)

Counter update formula (classic version):
- If the element is already being tracked:

$$
\text{count}[x] \leftarrow \text{count}[x] + j
$$

- For new elements, replaces the element with the smallest counter:

$$
\text{count}[x_{\text{new}}] \leftarrow \text{count}[x_{\text{min}}] + j
$$

## Implementation Features

* Support for custom hash functions
* Sharding for multi-threaded environments
* EWMA smoothing to eliminate abnormal frequency spikes
* Built-in metric collection

## EWMA Smoothing for Peak Reduction

To prevent unrealistic counter growth, we use **Exponentially Weighted Moving Average (EWMA)**:

$$
\text{count}[x] \leftarrow \alpha \cdot j + (1 - \alpha) \cdot \text{count}[x]
$$

Parameters:
* **α (smoothing factor)** - determines how strongly new data affects the counter (typically `0.01 ≤ α ≤ 0.1`)
* Lower **α** values result in smoother changes (but slower response to new data)

## Usage Example

```go
package main

import (
	"fmt"

	"github.com/koykov/hash/xxhash"
	"github.com/koykov/pbtk/heavy/spacesaving"
)

const (
	K         = 5
	alphaEWMA = 0.999
)

func main() {
	hasher := xxhash.Hasher64[[]byte]{}
	hitter, err := spacesaving.NewHitter[string](spacesaving.NewConfig(K, hasher).
		WithEWMA(alphaEWMA))
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
	fmt.Printf("%+v", hits) // [{Key:key0 Rate:649.3094164130389} {Key:key992 Rate:1} {Key:key4 Rate:1} {Key:key24 Rate:1} {Key:key999 Rate:1}]
}
```

## Application Areas

1. **Network traffic analysis** (DDoS detection, top IP addresses)
2. **Recommendation systems** (most viewed products)
3. **Web application logs** (frequent error detection)
4. **Search query analysis** (trends, popular keywords)
5. **Sensor data processing** (anomalies, frequent events)

## Conclusion

This **Space-Saving** implementation with **EWMA, sharding, and metrics** support enables:
* Efficient detection of frequent elements in data streams
* Operation in **multi-threaded environments** with minimal overhead
* Prevention of abnormal counter growth through **EWMA**
* Easy monitoring integration thanks to **built-in metrics**

The algorithm is particularly useful for **real-time analytics**, **monitoring**, and **anomaly detection**.
