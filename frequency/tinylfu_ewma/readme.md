# TinyLFU

TinyLFU (Tiny Least Frequently Used) is a probabilistic data structure based on [Count-Min Sketch](../cmsketch)
that estimates element frequency in data streams while also being able to "forget" stored element data through
periodic decay or smooth fading.

## Usage

Minimal working example:

```go
package main

import (
	"time"

	"github.com/koykov/hash/xxhash"
  "github.com/koykov/pbtk/frequency/tinylfu_ewma"
)

const (
	confidence = 0.99999
	epsilon    = 0.00001
)

func main() {
	conf := tinylfu.NewConfig(confidence, epsilon, xxhash.Hasher64[[]byte]{}).
		WithEWMATau(60) // smoothing constant time 1 minute
	est, err := tinylfu.NewEstimator(conf)
	_ = err
	_ = est.AddN("foobar", 5)
	println("time=0", est.Estimate("foobar")) // time=0 +5.000000e+000
	time.Sleep(15 * time.Second)
	println("time=15", est.Estimate("foobar")) // time=15 +3.894004e+000
	time.Sleep(30 * time.Second)
	println("time=45", est.Estimate("foobar")) // time=45 +2.361833e+000
	time.Sleep(15 * time.Second)
	println("time=60", est.Estimate("foobar")) // time=60 +1.839397e+000

	_ = est.AddN("foobar", 30)
	println("time=0 (after update)", est.Estimate("foobar")) // time=0 (after update) +2.000000e+001 (20)
	time.Sleep(15 * time.Second)
	println("time=15 (after update)", est.Estimate("foobar")) // time=15 (after update) +1.557602e+001 (~15.576)
}
```

As seen in the example, the structure counts the frequency of element "foobar" with smooth fading. The theoretical basis
of the algorithm will be described in the next chapter.

The [Config](config.go) structure allows for more detailed configuration of TinyLFU (see description and comments).

## How it works

First, let's consider how Count-Min Sketch (next CMS) works. CMS is also a probabilistic structure for solving
the frequency estimation problem with given accuracy $ε$ (epsilon) and confidence $δ$. It represents a matrix of
counters
with width $w$ and depth $d$, calculated by the following:

$$
w = \lceil{e \over ϵ}\rceil
$$

$$
d = \lceil ln({1 \over {1-δ}})\rceil
$$

where $e$ is the base of natural logarithm (~2.718).

For each added/updated element $x$ with weight $Δ$:

* position $j$ is calculated based on $d$ hashes

$$
j = {hash_i(x) \bmod w}
$$

* counter at position $i$ and $j$ is incremented

$$
CMS[i][j] += Δ
$$

To calculate frequency $E$ of element $x$:

* looks for minimum value of all $d$ counters for element $x$

$$
E(x) = \min\limits_{i∈[0,d−1]} CMS[i][hash_i(x) \bmod w]
$$

The problem with CMS is overestimation (due to hash collisions). Classic TinyLFU solves this through periodic counters
decay - the struct contains a total counter of element additions/updates and, upon reaching a certain limit, multiplies
all counters by $decay factor$ (typically $0.5$). I didn't like this approach due to high time complexity (traverse
the entire CMS) and excessive use of atomic operations. Therefore, I adapted the EWMA (Exponentially Weighted Moving
Average)
formula as a more elegant solution - it avoids traversing the entire CMS and solves the locking problem. The formula is:

$$
freq = counter * e^{-{Δt \over τ}} + Δ * (1 - e^{-{Δt \over τ}})
$$

where:

* $e$ - base of natural logarithm (~2.718)
* $Δt$ - time since last element update
* $τ$ (tau) - smoothing constant (typically 1 second, but configurable)
* $counter$ - current counter value in CMS
* $Δ$ - weight of updated element

> [!NOTE]
> This same formula is used for Linux load average calculation.

Since counter isn't enough (time delta is required as well), the counter became composite value - a 64-bit unsigned
number
where:

* first 32 bits store time since TinyLFU start (for compactness - using seconds allows storing 2^32 seconds, ~136 years)
* last 32 bits store the counter value

The 32-bit counter limit restricts values to 4294967295, so consider this when designing your system.

Why such complexity? First, it enables atomic counter operations in concurrent environments. Second, it saves memory.
It also adds slight overhead for EWMA application - before calculation, the counter value must be split into two 32-bit
values, and after calculation, two 32-bit values (new update time and calculated counter value) must be combined into
one 64-bit value.

Thus, to calculate frequency, it's sufficient to apply:

$$
Δt, counter = decode(CMS[i][hash_i(x) \bmod w])
$$

$$
E(x) = \min\limits_{i∈[0,d−1]} counter * e^{-{Δt \over τ}}
$$

This way, the problem is solved through math transformations without heavy decay operations over counters.

## References

Note these references describe classic TinyLFU (with periodic counter decay). This version uses EWMA smoothing
(see chapter below), while the classic implementation is provided for reference.

* https://florian.github.io/count-min-sketch/ Count-Min Sketch description and interactive example
* https://arxiv.org/abs/1512.00727 TinyLFU description
* https://highscalability.com/design-of-a-modern-cache/ another TinyLFU description
* https://towardsdatascience.com/time-series-from-scratch-exponentially-weighted-moving-averages-ewma-theory-and-implementation-607661d574fe/
  EWMA theory
* https://www.brendangregg.com/blog/2017-08-08/linux-load-averages.html EWMA: Linux load average application
* https://observablehq.com/@stwind/exponentially-weighted-moving-average EWMA: interactive calculator
