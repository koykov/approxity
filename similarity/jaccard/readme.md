# Jaccard Similarity

The **Jaccard Similarity Coefficient** (also known as the **Jaccard Index**) is a measure of similarity between two sets.
It is widely used in:
- **Text mining** (document similarity, plagiarism detection)
- **Recommendation systems** (collaborative filtering)
- **Bioinformatics** (comparing genetic sequences)
- **Duplicate detection** (finding near-identical records)

The coefficient ranges from **0** (no similarity) to **1** (identical sets).

## Math basics

Given two sets **A** and **B**, the Jaccard Similarity is defined as:

$$
J(A, B) = \frac{|A \cap B|}{|A \cup B|}
$$

Where:
- $|A \cap B|$ = Number of common elements in **A** and **B**
- $|A \cup B|$ = Total number of unique elements in **A** and **B**

### Jaccard Distance

If a **dissimilarity measure** is needed, the **Jaccard Distance** can be computed as:

$$
J_{\text{distance}}(A, B) = 1 - J(A, B)
$$

A distance of **0** means the sets are identical, while **1** means they are completely different.

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
	est, _ := jaccard.NewEstimator[[]byte](jaccard.NewConfig[[]byte](lsh))

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

## Key Properties
1. **Range**: Always between **0** and **1**.
2. **Commutative**: $J(A, B) = J(B, A)$.
3. **Binary Features**: Works best with **set-based data** (presence/absence of elements).
4. **Efficiency**: Can be optimized using **hashing** (e.g., MinHash for large datasets).

## Applications
* Text Similarity
  - Split documents into **shingles (n-grams)**.
  - Compare sets of shingles using Jaccard Similarity.
* MinHash Approximation
  - For large datasets, **MinHash** can estimate Jaccard Similarity efficiently.
  - Uses **sketches (hash signatures)** instead of full sets.
* Duplicate Detection
  - Identify near-duplicate records in databases.
  - Useful in **data cleaning** and **record linkage**.
