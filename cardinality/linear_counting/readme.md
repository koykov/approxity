# Linear counting

Linear Counting is a probabilistic algorithm used to estimate the number of unique elements in a large dataset. It is
particularly useful when memory efficiency is crucial, as it approximates cardinality using a bit array instead of
storing all unique elements.

## How it works

* Initialization: A bit array of size `m` is created, initialized to 0.
* Hashing: Each element in the dataset is hashed, and the corresponding bit in the array is set to `1`.
* Estimation: The number of unique elements is estimated using the formula:

$$ N = |-m * ln(1 - {n \over m})| $$

where `n` is the number of bits set to `1`.

## Usage

The minimal working example:

```go
import (
	"github.com/koykov/approxity/linear_counting"
	"github.com/koykov/hash/xxhash"
)

const n = 10 // desired number of unique elements

func main() {
    est, err := linear.NewEstimator[[]byte](linear.NewConfig(n, xxhash.Hasher64[[]byte]{}))
    _ = err
    est.Add("foo")
    est.Add("bar")
    fmt.Println(est.Estimate()) // 2
}
```
, but [initial config](config.go) allows to tune filter for better efficiency:
```go
import "github.com/koykov/approxity/cardinality/metrics/prometheus"

func func main() {
    // set filter size and hasher
    config := linear.NewConfig[string](1000, xxhash.Hasher64[[]byte]{}).
        // switch to race protected bit array (atomic based)
        WithConcurrency().
		// set allowed collision probability
        WithCollisionProbability(0.05).
        // cover with metrics
        WithMetricsWriter(prometheus.NewPrometheusMetrics("example_filter"))
    
    // filter is ready to use
    est, _ := linear.NewEstimator(config)
    ...
}
```

## Key Features
* Memory Efficiency: Uses a bit array to track uniqueness.
* Scalability: Suitable for large datasets.
* Tunable Accuracy: Accuracy depends on the size of the bit array `m`.
