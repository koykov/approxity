# HyperLogLog

HyperLogLog is an algorithm for the count-distinct problem, approximating the number of distinct elements in a multiset.
See [full description](https://en.wikipedia.org/wiki/HyperLogLog) for more details.

## Usage

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
