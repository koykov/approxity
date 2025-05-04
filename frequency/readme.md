# Frequency Estimation

This repository contains Go implementations of probabilistic data structures for solving frequency estimation problems.

## What is Frequency Estimation?

Frequency estimation is the task of counting or estimating the number of element occurrences in a data stream.
It addresses problems such as:

* Identifying "hot" (most frequent) items
* Filtering rare data
* Cache optimization (e.g., in LFU caches)

## Implemented Structures

* [**Count-Min Sketch**](cusketch) - A probabilistic structure providing an upper-bound frequency estimate with specified
  accuracy.
* [**Conservative Update Sketch**](cusketch) - A Count-Min Sketch modification that reduces error through
  conservative updates (only minimal counters are incremented).
* [**Count Sketch**](countsketch) - Unlike Count-Min Sketch, this structure can provide both upper and lower frequency bounds.
* [**TinyLFU**](tinylfu) - An adaptive frequency estimation structure optimized for cache usage.
* [**TinyLFU (EWMA version)**](tinylfu_ewma) - A TinyLFU variation using Exponential Weighted Moving Average for better
  adaptation to frequency distribution changes. This implementation is particularly recommended as it's significantly
  more optimal than classic `TinyLFU`.

## Implementation Features

All structures are designed for high-load concurrent environments:

* Using atomic operations instead of heavy locks
* Minimizing memory consumption/allocations
* Utilizing SIMD instructions for computation acceleration

### Initialization

Each package contains a `Config` structure for flexible configuration. See example [`Config`](cmsketch/config.go).
Common configuration options include:

* Required hash function (mandatory parameter)
* Confidence/Epsilon parameters for accuracy control and error tolerance
* Compact flag for using 32-bit counters
* Concurrent operations support (disabled by default)
* [`MetricsWriter`](metrics.go) parameter for metrics collection

### State Serialization

Structures support state serialization through `io.WriterTo` and restoration via `io.ReaderFrom`.
This solves the cold start problem - after system restart, you can continue collecting statistics from saved data
rather than starting from scratch.

### Unified Interface

All implementations conform to the [`Estimator`](interface.go) interface, which provides:

* Adding elements
* Estimating an element's frequency
* Structure clearing

This enables easy swapping between different structures without code changes and provides flexibility in choosing
the optimal algorithm for specific tasks.

Some structures implement `SignedEstimator` and `PreciseEstimator` (see [interface.go](interface.go)) due to implementation
specifics - the ability to provide negative and fractional frequency estimates.

### Monitoring and Metrics

The `Config` structure accepts a [`MetricsWriter`](metrics.go) implementation for writing metrics:

* Number of added elements
* Frequency histogram of queried elements

Similar to the `Estimator` interface, `MetricsWriter` also has `SignedMetricsWriter` and `PreciseMetricsWriter` variants
(see [metrics.go](metrics.go)).

Using metrics helps solve the "black box" problem - you can always evaluate how effectively the structure performs
its task and optimize configuration when needed.

An out-of-the-box [Prometheus](../metrics/prometheus/frequency.go) TSDB implementation is available.
Custom implementations can be created for other TSDBs (e.g., VictoriaMetrics) if required.

## Use Cases

* Cache optimization (LFU caches)
* Network traffic analysis (identifying frequent requests)
* Spam filtering (detecting common patterns)
* Recommendation systems (identifying popular content)
* Load balancing (detecting "hot" keys in distributed systems)

## Conclusion

These probabilistic data structures provide efficient solutions for estimating element frequencies in data streams.
Their features make them particularly valuable in high-load systems where both performance and minimal resource usage 
are crucial. The ability to easily switch between different implementations allows selecting the optimal solution for
specific use cases.
