# Heavy Hitters

This repository contains implementations of probabilistic heavy hitters data structures in Go.

## What are Heavy Hitters?

Heavy Hitters is a class of algorithms designed to identify the most frequently occurring elements in a data stream.  
They solve the problem of efficiently finding elements whose frequency exceeds a given threshold while minimizing computational resources and memory usage.

## Implemented Algorithms

* [Space-Saving](spacesaving) - tracks a limited number of elements, replacing the least probable candidates when overflow occurs.
* [Misra-Gries](misragries) - provides an approximate solution using a fixed number of counters, ensuring guaranteed accuracy for high-frequency elements.
* [Lossy Counting](lossy) - gradually reduces counter weights, allowing efficient identification of frequent elements in infinite data streams with controlled error margins.

## Implementation Features

All structures are designed to work in high-load, multi-threaded environments:

* Sharding based on element hashes to reduce contention
* Minimized memory usage and allocations

### Initialization

Each package contains a `Config` structure for flexible configuration.  
Example [`Config`](spacesaving/config.go). Common config parameters include:

* Required hash function (mandatory)
* Number of counters or acceptable error bounds (mandatory)
* Number of shards (optional, default is 4)
* [`MetricsWriter`](metrics.go) parameter for metric collection

### Unified Interface

All implementations adhere to the [`Hitter`](interface.go) interface, which allows:

* Adding an element
* Retrieving frequency estimates for top elements

This enables easy swapping of one structure for another without modifying application code and provides flexibility in choosing the optimal algorithm for a specific task.

### Monitoring and Metrics

The `Config` structure allows passing a [`MetricsWriter`](metrics.go) implementation to collect metrics such as:

* Number of elements added
* Min/Max frequency
* Average frequency
* Sum of all frequencies
* Standard deviation (stddev)
* Median frequency
* Frequency variance
* Distribution asymmetry (skewness)
* Variation coefficient (normalized variance)
* Frequency percentiles (25, 50, 75, 90, 99)
* Frequency histogram

Using metrics helps solve the "black box" problem—you can always evaluate how effectively the structure performs and optimize its configuration if needed.

An out-of-the-box [Prometheus](../metrics/prometheus/heavy.go) TSDB implementation is included. If necessary, you can write a custom implementation for your preferred TSDB (e.g., VictoriaMetrics).

## Use Cases

* Network traffic analysis for anomaly detection
* Search query processing in large-scale search engines
* Click monitoring in online advertising
* Trend detection in social media
* Cache optimization in databases

## Conclusion

The implemented algorithms can be effectively used in high-load environments for processing large data streams.  
They balance result accuracy and resource consumption, making them suitable for real-time systems and distributed computing.
