# LogLog

LogLog is a probabilistic data structure used to estimate the number of unique elements (cardinality) in a large dataset
with minimal memory usage. It is the predecessor to [HyperLogLog](../hyperloglog) and provides a simpler but less accurate
approach to cardinality estimation.

## How It Works

* **Hashing**: Each element $x$ is hashed into a binary string using a hash function $h(x)$. The hash function should
produce uniformly distributed outputs.

* **Bucketing**: The hash is divided into two parts:
  * The first $p$ bits determine the bucket index $j$ (where $m = 2^p$ is the number of buckets).
  * The remaining bits are used to count the number of leading zeros $ρ(w)$ in the binary representation.

* **Tracking Maximum Leading Zeros**: For each bucket $j$, the maximum number of leading zeros $M_j$ is tracked. The
cardinality $E$ is estimated using the geometric mean of $M_j$:

$$E=α_m⋅m⋅2^{{1 \over m} \sum_{j=1}^m {M_j} }$$

Where:
  * $α_m$ is a correction factor for small and large ranges (e.g., $α_{16}≈0.773$).
  * $m$ is the number of buckets.

* **Bias Correction**: For small cardinalities, a bias correction is applied to improve accuracy:

$$
E'=\left\{
\begin{array}{ll}
m⋅log({m \over V}) &\text{if }E \leq {5 \over 2}m \\
E &\text{otherwise},
\end{array}
\right.
$$

where $V$ is the number of buckets with $M_j=0$.

## Usage

The minimal working example:
```go
import (
    "github.com/koykov/pbtk/cardinality/loglog"
    "github.com/koykov/hash/xxhash"
)

func main() {
    est, err := loglog.NewEstimator[string](loglog.NewConfig(18, xxhash.Hasher64[[]byte]{}))
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
    config := loglog.NewConfig[string](18, xxhash.Hasher64[[]byte]{}).
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

* **Memory Efficiency**: Uses only a few kilobytes of memory, even for billions of elements.
* **Moderate Accuracy**: Provides an estimate with a typical error rate of about 2-3%.
* **Scalability**: Efficiently handles large datasets.
