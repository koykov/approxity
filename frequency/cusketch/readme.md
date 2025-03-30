# Conservative Update Sketch

A memory-efficient variant of [Count-Min Sketch](../cmsketch) that reduces overestimation errors by updating only
the minimum counter per row during increments. Unlike the standard Count-Min Sketch (which increments all corresponding
counters), Conservative Update provides more accurate frequency estimates for low-frequency items while maintaining
the same space complexity.

Ideal for applications where precise counts matter, such as network traffic monitoring or anomaly detection.
Trade-off: slightly slower updates due to the extra minimum-finding step.

Key Difference vs Count-Min Sketch:

* **Conservative Update**: Tighter error bounds (2-3× less overestimation).
* **Count-Min Sketch**: Faster updates, simpler implementation.

## Usage

The same as [Count-Min Sketch example](../cmsketch/readme.md#usage).

