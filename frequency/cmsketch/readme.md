# Count-Min Sketch

Count-Min Sketch is a probabilistic data structure that estimates the frequency of elements in a data stream with minimal
memory usage. It provides approximate counts with controlled error bounds.

## How It Works

* **Initialization**:
  * CMS initializes two parameters:
    * $confidence$ - possibility that potential error will be in range of acceptable error rate
    * $ϵ (epsilon)$ - estimation precision (0..1)
  * Making a 2D array of counters $C[w][d]$, where
    * $w$ - number of counters (width) $$ w = {e \over ϵ} = {Euler \over epsilon} $$
    * $d$ - number of hash functions (height) $$ d = ln({1 \over δ}) = ln({1 \over {1-confidence}}) $$
* **Insertion**:
  * For item $x$ and its weight $Δ$:
    * for each hash function $h_i$ calculates index $j = h_i(x)$
    * counter $C[i][j]$ increases to $Δ$
* **Estimation**:
  * For item $x$:
    * for each hash function $h_i$ calculates index $j = h_i(x)$
    * estimation $E$ is a minimal value of all counters $C[i][j]$ $$E(x) = \min_{i=1..d}(C[i][j])$$

## Usage

The minimal working example:
```go
import (
    "github.com/koykov/approxity/frequency/cmsketch"
    "github.com/koykov/hash/xxhash"
)

const (
	confidence = 0.99
	epsilon = 0.01
)

func main() {
    est, err := cmsketch.NewEstimator[string](cmsketch.NewConfig(confidence, epsilon, xxhash.Hasher64[[]byte]{}))
    _ = err
    for i:=0; i<5; i++ {
        key := fmt.Sprintf("item-%d", i)
        _ = est.Add(key)
        if i == 3 {
            for j:=0; j<1e6; j++ {
                _ = est.Add(key)
            }
        }
    }
    println(est.Estimate("item-3")) // ~1000000
}
```
, but [initial config](config.go) allows to tune estimation for better efficiency:
```go
import "github.com/koykov/approxity/cardinality/metrics/prometheus"

func func main() {
    // set estimation precision and hasher
    config := cmsketch.NewConfig[string](18, xxhash.Hasher64[[]byte]{}).
        // use 32-bit counters instead of 64-bit to reduce memory usage
        WithCompact().
        // switch to race protected bit array (atomic based)
        WithConcurrency().
        // cover with metrics
        WithMetricsWriter(prometheus.NewPrometheusMetrics("example_estimation"))
    
    // estimation is ready to use
    est, _ := cmsketch.NewEstimator(config)
    ...
}
```

## Key Features

* Memory efficient.
* Easy to implement.
* Frequency may be overestimated due to CMS guarantees upper bound of error.
