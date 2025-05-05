# Similarity Estimation

Similarity estimation is a set of techniques for measuring the degree of similarity between two datasets.  
These algorithms solve problems such as:
* Finding duplicates or near-duplicates
* Clustering similar documents
* Recommendation systems
* Plagiarism detection
* Removing similar records in datasets

## Implemented Algorithms
* **Hamming Distance** - Measures the number of differing bits between two vectors.  
  Effective for comparing binary data or fixed-length hashes.
* **Cosine Similarity** - Estimates similarity by the angle between vectors in multidimensional space.  
  Widely used for text data represented as feature vectors.
* **Jaccard Distance** - Computes the difference between sets as the proportion of non-matching elements.  
  Well-suited for comparing word sets or shingles.

## Implementation Features

* **High Performance**: Minimized memory allocations, efficient data structures used
* **Unified Interface**: All algorithms implement the `Estimator` interface
* **Interchangeability**: Algorithms can be swapped without modifying core code
* **LSH Integration**: Works with vectors produced by LSH algorithm

## Usage

```go
package main

import (
	"github.com/koykov/hash/xxhash"
	"github.com/koykov/pbtk/lsh/minhash"
	"github.com/koykov/pbtk/shingle"
	"github.com/koykov/pbtk/similarity/jaccard"
)

func main() {
	hasher := xxhash.Hasher64[[]byte]{}
	shingler := shingle.NewChar[[]byte](3, "") // 3-gram
	lsh, _ := minhash.NewHasher[[]byte](minhash.NewConfig[[]byte](hasher, 50, shingler))
	est, err := jaccard.NewEstimator[[]byte](jaccard.NewConfig[[]byte](lsh))
	_ = err

	e, _ := est.Estimate([]byte("Four children are doing backbends in the gym"), []byte("Four children are doing backbends in the park"))
	println(e) // 0.8478260869565217 (high similarity)

	est.Reset()
	e, _ = est.Estimate([]byte("A man is sitting near a bike and is writing a note"), []byte("A man is standing near a bike and is writing on a piece paper"))
	println(e) // 0.532258064516129 (medium similarity)

	est.Reset()
	e, _ = est.Estimate([]byte("One white dog and one black one are sitting side by side on the grass"), []byte("A black and a white dog are joyfully running on the grass"))
	println(e) // 0.44155844155844154 (low similarity)
}
```

## Use Cases

1. **Finding similar documents** in large text collections
2. **Duplicate detection** in e-commerce product catalogs
3. **Content recommendations** based on similarity to previously viewed items
4. **Plagiarism detection** in academic works
5. **News clustering** on the same topic from different sources

## Conclusion

The implemented algorithms provide a flexible toolkit for solving a wide range of text comparison tasks.  
Thanks to the unified interface and performance optimizations, they can be easily integrated into existing data processing systems.
