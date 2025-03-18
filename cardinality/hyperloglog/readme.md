# HyperLogLog

HyperLogLog is a probabilistic data structure used to estimate the number of unique elements (cardinality) in a large
dataset with minimal memory usage. It is particularly useful for scenarios where exact counting is impractical due to
memory constraints.

HyperLogLog is an algorithm for the count-distinct problem, approximating the number of distinct elements in a multiset.
See [full description](https://en.wikipedia.org/wiki/HyperLogLog) for more details.

## How It Works

* Hashing: Each element is hashed into a binary string.
* Bucketing: The hash is divided into buckets, and the number of leading zeros in the binary string is counted.
* Averaging: The harmonic mean of the counts across buckets is used to estimate the cardinality.

## Usage

* Initialization: Create a HyperLogLog structure with a specified number of buckets (e.g., 2^14 buckets for ~1.5% error rate).
* Adding Elements: Insert elements into the HyperLogLog structure.
* Estimating Cardinality: Retrieve the estimated number of unique elements.

The minimal working example:
```go
import (
    "github.com/koykov/approxity/cardinality/hyperloglog"
    "github.com/koykov/hash/xxhash"
)

func main() {
    est, err := hyperloglog.NewEstimator[string](hyperloglog.NewConfig(18, xxhash.Hasher64[[]byte]{}))
    _ = err
    for i:=0; i<5; i++ {
	    for j:=0; j<1e6; j++ {
		    _ = est.Add(fmt.Sprintf("item-%d", j))
	    }	
    }
	println(est.Cardinality()) // ~1000000
}
```
, but [initial config](config.go) allows to tune estimation for better efficiency:
```go
import "github.com/koykov/approxity/cardinality/metrics/prometheus"

func func main() {
    // set estimation precision and hasher
    config := hyperloglog.NewConfig[string](18, xxhash.Hasher64[[]byte]{}).
        // switch to race protected bit array (atomic based)
        WithConcurrency().
        // cover with metrics
        WithMetricsWriter(prometheus.NewPrometheusMetrics("example_estimation"))
    
    // estimation is ready to use
    est, _ := NewEstimator(config)
	...
}
```

## Key Features

* Memory Efficiency: Uses only a few kilobytes of memory, even for billions of elements.
* High Accuracy: Provides an estimate with a typical error rate of about 1-2%.
* Scalability: Efficiently handles large datasets.

## References

* https://en.wikipedia.org/wiki/HyperLogLog
* [Original paper](http://algo.inria.fr/flajolet/Publications/FlFuGaMe07.pdf)
