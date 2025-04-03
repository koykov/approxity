# TinyLFU (Tiny Least Frequently Used)

An implementation of the TinyLFU - probabilistic data structure to solve frequency estimation problem.
TinyLFU is a space-efficient algorithm that approximates how frequently items appear in a data stream, solving
the frequency estimation problem with minimal memory footprint.

The implementation builds upon the [Count-Min Sketch](../cmsketch) (CMS) data structure, adding configurable decay
mechanisms to prevent overestimation and maintain accuracy over time.

## Relation to Count-Min Sketch

TinyLFU is essentially an enhanced Count-Min Sketch with decay capabilities:

* The core frequency counting is performed by an underlying CMS structure
* CMS is initialized with `confidence` and `epsilon` parameters that control its accuracy
* TinyLFU adds decay mechanisms to periodically reduce counter values, preventing unbounded growth
* Without decay parameters, TinyLFU behaves exactly like a standard [Count-Min Sketch](../cmsketch)

## Configuration Parameters

Core parameters (for CMS):

* `Confidence`: Probability that the estimate is within epsilon bounds
* `Epsilon`:  Error tolerance for frequency estimates

Decay parameters:

* `DecayLimit`: Number of new items after which decay is triggered
* `DecayInterval`: Time period after which decay occurs automatically
* `ForceDecayNotifier`: External signal interface for manual decay triggering (*)

These parameters compete with each over and control when decay operations occur. The first condition to be met triggers
the decay.

(*) **ForceDecayNotifier** signature:

```go
type ForceDecayNotifier interface {
Notify() <-chan struct{}
}
```

At any decay call:

* Resets the counter of new items
* Resets decay interval timer
* Performs the decay over CMS counters

> [!NOTE]
> If none of these parameters are set, no decay will occur and the structure will behave like a standard CMS
> (potentially overestimate frequencies).

## Decay Factors

These control how aggressively counters are reduced during decay:

* `DecayFactor`: (default: 0.5) Standard multiplier for counters (range 0..1)
* `SoftDecayFactor`: (default: 0.7) Used when:
  * Too little time has passed since last decay (<50% of `DecayInterval`)
  * Too few new items were added (<50% of `DecayLimit`)

## Usage

The minimal working example:

```go
import (
"github.com/koykov/pbtk/frequency/tinylfu"
"github.com/koykov/hash/xxhash"
)

const (
confidence = 0.99
epsilon = 0.01
)

func main() {
est, err := tinylfu.NewEstimator[string](tinylfu.NewConfig(confidence, epsilon, xxhash.Hasher64[[]byte]{}).
WithDecayLimit(100000))
_ = err
for i := 0; i<5; i++ {
key := fmt.Sprintf("item-%d", i)
_ = est.Add(key)
if i == 3 {
for j := 0; j<1e6; j++ {
_ = est.Add(key)
}
}
}
println(est.Estimate("item-3")) // ~1000000
}
```
