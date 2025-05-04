# Cardinality Estimation

This repository contains Go implementations of probabilistic data structures for solving cardinality estimation problems.

## What is Cardinality Estimation?

Cardinality estimation refers to methods for approximate counting of unique elements (cardinality) in large datasets.
It solves problems like analyzing unique visits, counting distinct values in databases, network traffic monitoring,
and other cases where exact counting would require prohibitively large memory.

## Implemented Data Structures

* [**LogLog**](loglog) - A probabilistic algorithm that estimates cardinality based on maximum leading zeros in element
  hashes.
* [**HyperLogLog**](hyperloglog) - An improved version of LogLog with higher accuracy using harmonic mean.
* [**HyperBitBit**](hyperbitbit) - A compact HyperLogLog variation combining bitmaps for small cardinalities and
  probabilistic estimation for large ones.
* [**Linear Counting**](linear_counting) - A simple bitmap-based algorithm effective for small and medium-sized sets.

## Implementation Features

All implementations are designed for high-load systems and multi-threaded environments:

* Lock-free design using only atomic operations
    * High performance for concurrent read/write operations
* Minimal memory consumption and allocations
* SIMD operations where possible

### Initialization

Each package contains a `Config` structure for flexible configuration. See example [`Config`](hyperloglog/config.go).
Common configuration options include:

* Required hash function (mandatory parameter)
* Estimated maximum number of elements (set with some margin)
* Accuracy or acceptable collision probability
* Concurrent operations support mode (disabled by default)
* [`MetricsWriter`](metrics.go) parameter for metrics collection

### State Serialization

All structures support internal state serialization via `io.WriterTo` and restoration via `io.ReaderFrom`.
This solves the "cold start" problem - allowing to save accumulated statistics between system restarts and quickly resume
operation without losing estimation accuracy.

Intended usage includes periodic state saving to storage (e.g., file or cloud) and reading on startup.
When concurrent operations mode is enabled, these operations will be protected from data races.

### Unified Interface

All implementations share a common [`Estimator`](interface.go) interface that allows:

- Adding elements
- Estimating cardinality of all added elements
- Clearing the structure

This ensures interchangeability of data structures without modifying client code, making it easy to compare and select
the most suitable algorithm for specific tasks.

### Monitoring and Metrics

Through the `Config` structure, you can provide a [`MetricsWriter`](metrics.go) implementation to each structure that
will record:

- Number of elements added
- Cardinality distribution histogram

This helps solve the "black box" problem - metrics help evaluate whether the structure is optimally configured and
enable fine-tuning when needed.

An out-of-the-box [Prometheus](../metrics/prometheus/cardinality.go) TSDB implementation is included.
You can also implement custom versions for other TSDBs (e.g., VictoriaMetrics).

## Use Cases

* Unique visitor analysis for websites and mobile apps
* Network traffic monitoring and anomaly detection (e.g., DDoS attacks)
* Database query optimization (DINCT-value estimation)
* Streaming platforms - counting unique viewers
* Big data - fast set size estimation in ETL processes

## Conclusion

These implementations can help estimate set cardinality in cases where more primitive methods (like hash tables) isn't enough.
Use metrics for configuration and always stay informed about what data enters into your system.
