# Count Sketch

Count Sketch is a probabilistic streaming algorithm for frequency estimation that uses randomized hashing with
signs (±1) to provide unbiased estimates.

Unlike [Count-Min Sketch](../cmsketch), it cancels out noise by aggregating signed contributions across multiple
hash tables, making it robust for detecting low-frequency elements in adversarial or noisy data.

## How It Works

* **Initialization**:
    * Count Sketch initializes using two parameters:
        * $confidence$ - possibility that potential error will be in range of acceptable error rate
        * $ϵ (epsilon)$ - estimation precision (0..1)
    * Making a 2D array of counters $C[w][d]$, where
        * $w$ - number of counters (width) $$w = {1 \over ϵ^2} = {1 \over epsilon^2}$$
        * $d$ - number of hash functions (height) $$d = ln({1 \over δ}) = ln({1 \over {1-confidence}})$$
* **Insertion**:
    * For item $x$ and its weight $Δ$:
        * for each hash function $h_i$ calculates index $$j = {h_i(x) \mod w}$$
        * for each hash function $s_i$ calculates sign
          $$
          sign=\left\{
          \begin{array}{ll}
          1 &\text{if }s_i(x) \mod 2 == 0\\
          -1 &\text{otherwise},
          \end{array}
          \right.
          $$
        * counter $C[i][j]$ increases to $sign * Δ$
* **Estimation**:
    * For item $x$:
        * for each hash function $h_i$ calculates index $$j = h_i(x) \mod w$$
        * for each hash function $s_i$ calculates sign
          $$
          sign=\left\{
          \begin{array}{ll}
          1 &\text{if }s_i(x) \mod 2 == 0\\
          -1 &\text{otherwise},
          \end{array}
          \right.
          $$
        * estimation $E$ is a median value of all counters $C[i][j]$ $$E(x) = \mathrm{med}(C[i][0], C[i][1], \dots C[i][d-1])$$

## Usage

The minimal working example:
```go
import (
    "github.com/koykov/approxity/frequency/countsketch"
    "github.com/koykov/hash/xxhash"
)

const (
	confidence = 0.99
	epsilon = 0.01
)

func main() {
    est, err := countsketch.NewEstimator[string](countsketch.NewConfig(confidence, epsilon, xxhash.Hasher64[[]byte]{}))
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
    config := countsketch.NewConfig[string](18, xxhash.Hasher64[[]byte]{}).
        // use 32-bit counters instead of 64-bit to reduce memory usage
        WithCompact().
        // switch to race protected bit array (atomic based)
        WithConcurrency().
        // cover with metrics
        WithMetricsWriter(prometheus.NewPrometheusMetrics("example_estimation"))
    
    // estimation is ready to use
    est, _ := countsketch.NewEstimator(config)
    ...
}
```

## Key Features

* **Unbiased estimates**: Median of signed counters ensures $E[f̂(x)] = f(x)$.
* **Heavy hitters**: Ideal for identifying significant elements in skewed distributions.

## References

TODO...
