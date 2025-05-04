# AMQ (Approximate Membership Query) Filters

This repository contains implementations of probabilistic data structures (AMQ filters) in Go, designed for highload
environments.

## What are AMQ Filters?

Approximate Membership Query (AMQ) structures are probabilistic data structures that efficiently test whether an element
belongs to a set. They solve the following problems:

- Filtering already processed elements
- Preliminary cache checking
- Reducing expensive database queries
- Decreasing network traffic in distributed systems

AMQ filters may yield false positives (but never false negatives) in exchange for compact data representation.

## Implemented Structures

- [**Bloom filter**](bloom_filter) - A classic probabilistic structure that uses multiple hash functions to set bits
  in a bit array.
- [**Cuckoo filter**](cuckoo_filter) - An improved version of Bloom filter that supports element deletion and offers
  higher storage density.
- [**Quotient filter**](quotient_filter) - A compact structure organizing data as a hash table with a special collision
  resolution method.
- [**Xor filter**](xor_filter) - One of the newest and most efficient structures, providing the lowest false-positive
  rate with compact storage.

## Implementation Features

All implementations are designed for highload systems and concurrent environments:

- High performance for parallel read/write operations
- Minimal locking and contention (using only atomic operations)
- Reduced memory allocations
- Minimal possible memory consumption
- SIMD operations where applicable

### Initialization

Each package contains a `Config` structure for flexible filter configuration. See example: [`Config`](bloom_filter/config.go).
Common configurable parameters include:

- Required hash function (mandatory)
- Expected number of elements to store (mandatory)
- Concurrent operations support (disabled by default)
- [`MetricsWriter`](metrics.go) parameter for metrics collection

### Serialization and Restoration

All structures support:
- Saving current state to `io.WriterTo`
- Restoring state from `io.ReaderFrom` (solving cold start problems)

This combination addresses the cold start problem - by periodically saving the filter state to storage (e.g., file or cloud),
you'll have a ready-to-use filter upon startup. When concurrent operations are enabled, these operations will be protected
from data races.

### Unified Interface

All filters implement a common [`Filter`](interface.go) interface that provides:

- Adding elements to the filter
- Checking element membership
- Removing elements (when supported, e.g., Bloom filters prohibit this)
- Getting current filter size/capacity
- Clearing the filter

This allows easy swapping of implementations without application code changes.

### Monitoring and Metrics

Through the `Config` structure, you can provide a [`MetricsWriter`](metrics.go) implementation that records:

- Number of elements added
- Number of elements removed
- Read operations count and their results

This approach solves the "black box" problem - metrics help evaluate filter efficiency and optimize its configuration.

The package contains builtin [Prometheus](../metrics/prometheus/amq.go) implementation.
Feel free to implement your own implementation (e.g. VictoriaMetrics...).

## Conclusion

AMQ filters are powerful tools for data filtering tasks, but they should only be used when simpler approaches
(e.g., hash tables) become inefficient in terms of memory or performance.
