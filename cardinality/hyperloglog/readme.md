# HyperLogLog

HyperLogLog is a probabilistic data structure used to estimate the number of unique elements (cardinality) in a large
dataset with minimal memory usage. It is particularly useful for scenarios where exact counting is impractical due to
memory constraints.

## How It Works

* **Hashing**: Each element `x` is hashed into a binary string using a hash function $h(x)$. The hash function should produce
uniformly distributed outputs.
* **Bucketing**: The hash is divided into two parts:
  * The first $p$ bits determine the bucket index $j$ (where $m = 2^P$ is the number of buckets).
  * The remaining bits are used to count the number of leading zeros $ρ(w)$ in the binary representation.
* **Estimating Cardinality**: For each bucket $j$, the maximum number of leading zeros $M_j$ is tracked. The cardinality $E$
is estimated using the harmonic mean of $M_j$:

$$
E=α_m⋅m^2⋅\left( \sum_{j=1}^m 2^{-M_j} \right)^{-1}
$$

Where:
  * $α_m$ is a correction factor for small and large ranges, e.g.: $α_{16} ≈ 0.673$
  * $m$ is the number of buckets.


* **Bias correction**

For small cardinalities, a bias correction is applied to improve accuracy:

$$
E'=\left\{
\begin{array}{ll}
m⋅log({m \over V}) &\text{if }E \leq {5 \over 2}m \\
E &\text{otherwise},
\end{array}
\right.
$$

where $V$ is the number of buckets with $M_j = 0$.

## Usage

* **Initialization**: Create a HyperLogLog structure with a specified number of buckets (e.g., $2^{14}$ buckets for ~1.5% error rate).
* **Adding Elements**: Insert elements into the HyperLogLog structure.
* **Estimating Cardinality**: Retrieve the estimated number of unique elements.

The minimal working example:
```go
import (
    "github.com/koykov/pbtk/cardinality/hyperloglog"
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
import "github.com/koykov/pbtk/cardinality/metrics/prometheus"

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
