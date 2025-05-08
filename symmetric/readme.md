# Symmetric Difference

This repository contains implementations of various symmetric difference algorithms optimized for high-load environments.

The symmetric difference between two sets $A$ and $B$ is the set of elements that belong to exactly one of the sets:

$$
A \triangle B = (A \setminus B) \cup (B \setminus A) = (A \cup B) \setminus (A \cap B)
$$

In practice, this is useful for:
* Identifying discrepancies between datasets
* Detecting changes between system states
* Implementing delta encoding in distributed systems

## Implementation Features

* All algorithms are designed for high-performance multi-threaded environments
* Support for custom LSH algorithms
* Implement a unified Differ interface, allowing easy algorithm swapping without application code changes
* SIMD processing of bit arrays and internal structure cleanup

## Use Cases

* Network change detection systems
    * Identifying configuration discrepancies between nodes
    * Detecting new/disappeared devices on the network
* Distributed databases
    * Comparing data segments between replicas
    * Identifying synchronization discrepancies
* Monitoring systems
    * Detecting metric changes between time points
    * Identifying anomalies in streaming data
* Blockchain and distributed ledgers
    * Comparing node states
    * Detecting transaction discrepancies

## Conclusion

The provided implementations offer:

* High performance in multi-threaded environments
* Flexibility in hash algorithm selection
* Unified interface for different approaches
* Optimized computations through SIMD

When selecting a specific algorithm, consider your data characteristics and accuracy requirements.
