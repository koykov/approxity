# LSH (Locality-Sensitive Hashing)

Locality-Sensitive Hashing (LSH) is an approximate nearest neighbor search method that hashes input data in such a way
that similar objects are likely to fall into the same "bucket". LSH solves the following problems:
* Duplicate detection
* Clustering similar documents
* Recommendation systems
* Plagiarism detection

## Implemented Algorithms

* **SimHash** - Generates a document fingerprint where similar documents have small Hamming distance between their hashes.
  Particularly effective for near-duplicate content detection.
* **MinHash** - Estimates set similarity based on the probability of minimum hash collisions. Optimal for comparing datasets
  (e.g., word sets in documents).
* **B-Bit MinHash** - An optimized version of MinHash that uses only the least significant bits of each hash, significantly
  reducing storage requirements while maintaining acceptable accuracy.

The implemented algorithms are designed for processing text data only.

## Implementation Features

* **High performance**: Memory allocations minimized, efficient data structures used
* **Unified interface**: All algorithms implement the [`Hasher`](interface.go) interface, allowing easy swapping between algorithms
* **Flexibility**: Supports various shinglers for text-to-shingle conversion and different shingle hashing algorithms
* **Scalability**: Optimized for working with large volumes of data

## Usage

```go
package main

import (
	"fmt"

	"github.com/koykov/hash/xxhash"
	"github.com/koykov/pbtk/lsh/minhash"
	"github.com/koykov/pbtk/shingle"
)

const k = 50 // number of hash functions

func main() {
	hasher := xxhash.Hasher64[[]byte]{}
	shingler := shingle.NewChar[[]byte](3, "") // 3-gram char shingler
	h, err := minhash.NewHasher[[]byte](minhash.NewConfig(hasher, k, shingler))
	_ = err
	_ = h.Add([]byte("A person in a black jacket is doing tricks on a motorbike"))
	fmt.Printf("%#v\n", h.Hash()) // []uint64{0xdb9aae3abdd1a77, 0x46b42729237c368, ..., 0x3f710684c9bab1f}

	h.Reset()
	_ = h.Add([]byte("Two men are taking a break from a trip on a snowy road"))
	fmt.Printf("%#v\n", h.Hash()) // []uint64{0x3e2eaf1df906ab, 0x4a12c40f3825533, ..., 0x5505351f1a39fa3}
}
```

## Application Examples

1. **Document duplicate detection** in large text corpora
2. **Recommendation systems** for finding similar content
3. **Plagiarism detection** in academic papers or source code
4. **News clustering** for grouping similar stories from different sources
5. **Finding similar users** by their activity or preferences

## Conclusion

This library provides efficient implementations of popular LSH algorithms ready for highload environments.
Thanks to its unified interface and performance optimizations, it can be easily integrated into existing systems
for similar object search tasks. The choice of a specific algorithm depends on data characteristics and
similarity estimation accuracy requirements.
