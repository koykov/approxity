# Cuckoo filter

The cuckoo filter is a minimized hash table that uses cuckoo hashing to resolve collisions. It minimizes its space
complexity by only keeping a fingerprint of the value to be stored in the set. Much like the bloom filter uses single
bits to store data and the counting bloom filter uses a small integer, the cuckoo filter uses a small `f`-bit fingerprint
to represent the data. The value of `f` is decided on the ideal false positive probability the programmer wants. This
implementation use `f` containing 8 bits (`uint8` value).

## Cuckoo hashing

In cuckoo hashing, each key is hashed by two different hash functions, so that the value can be assigned to one of two
buckets. The first bucket is tried first. If there's nothing there, then the value is placed in bucket 1. If there is
something there, bucket 2 is tried. If bucket 2 if empty, then the value is placed there. If bucket 2 is occupied, then
the occupant of bucket 2 is evicted and the value is placed there.

Now, there is a lone (key, value) pair that is unassigned. But, there are two hash functions. The same process begins
for the new key. If there is an open spot between the two possible buckets, then it is taken. However, if both buckets
are taken, one of the values is kicked out and the process repeats again.

## Usage

The minimal working example:
```go
import (
    "github.com/koykov/cuckoo_filter"
    "github.com/koykov/hash/metro"
)

func main() {
    f, err := cuckoo.NewFilter(cuckoo.NewConfig(1e7, metro.Hasher64[[]byte]{Seed: 1234}))
    _ = err
    _ = f.Set("foobar")
    print(f.Contains("foobar")) // true
    print(f.Contains("qwerty")) // false
}
```

Similar to [bloom filter](../bloom_filter/readme.md#usage) cuckoo allows to initiate config more detailed:
```go
import (
    "github.com/koykov/cuckoo_filter"
    "github.com/koykov/hash/metro"
)

func main() {
    // set filter size and hasher
    f, _ := cuckoo.NewFilter(cuckoo.NewConfig(1e7, metro.Hasher64[[]byte]{Seed: 1234}).
		// limit for cuckoo kicks to avoid infinite loop
        WithKicksLimit(10).
        // switch to race protected bit array (atomic based)
        WithConcurrency().
        // cover with metrics
        WithMetricsWriter(prometheus.NewPrometheusMetrics("example_filter")))
	...
```
