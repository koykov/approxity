# HyperBitBit

> [!CAUTION]
> Not production ready!
> Estimation is too inaccurate, especially on small datasets. Use with caution. 

A memory-efficient alternative to HyperLogLog.

See [full scription](https://www.birs.ca/workshops/2022/22w5004/files/Bob%20Sedgewick/HyperBit.pdf).

## Usage

The minimal working example:
```go
import (
    "github.com/koykov/approxity/cardinality/hyperbitbit"
    "github.com/koykov/hash/xxhash"
)

func main() {
    est, err := hyperbitbit.NewEstimator[string](hyperbitbit.NewConfig(5, xxhash.Hasher64[[]byte]{}))
    _ = err
    for i:=0; i<5; i++ {
	    for j:=0; j<1e6; j++ {
		    _ = est.Add(fmt.Sprintf("item-%d", j))
	    }	
    }
	println(est.Cardinality())
}
```
