# HyperBitBit

> [!CAUTION]
> Not production ready!
> Estimation is too inaccurate, especially on small datasets. Use with caution. 

HyperBitBit is a memory-efficient probabilistic algorithm for estimating the number of unique elements (cardinality) in
large datasets. It is an enhancement over the HyperLogLog algorithm, designed to reduce memory usage while maintaining
high accuracy.

See [full scription](https://www.birs.ca/workshops/2022/22w5004/files/Bob%20Sedgewick/HyperBit.pdf).

## How It Works

* Initialization: A small bit array (or register) is initialized to track information about the dataset.
* Hashing: Each element is hashed, and the algorithm uses the hash to update the bit array.
* Estimation: The cardinality is estimated based on the patterns of bits set in the array, using a probabilistic formula.

## Usage
* Initialize the HyperBitBit counter.
* Add elements to the counter using the `Add` method.
* Estimate the cardinality using the `Estimate` method.

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
	println(est.Estimate())
}
```

## Key Features

* Memory Efficiency: Uses significantly less memory than traditional methods.
* Scalability: Handles large datasets with ease.
* High Accuracy: Provides accurate estimates even with minimal memory usage.
