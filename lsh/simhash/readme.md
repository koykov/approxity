# SimHash

**SimHash** is a locality-sensitive hashing (LSH) algorithm designed to detect near-duplicate content by generating
compact fingerprint vectors. Unlike [MinHash](../minhash) (which estimates Jaccard similarity), SimHash approximates
**cosine similarity** between documents, making it particularly effective for detecting small modifications in text.

This implementation focuses on computing the **SimHash fingerprint**—a fixed-length binary vector—which can be compared
later using Hamming distance to estimate document similarity.

## How It Works

### 1. Feature Extraction

The input string is processed into features (typically word n-grams or shingles). Common approaches:

- **Word-based**: Tokens ("the quick fox" → [["the quick"], ["quick", "fox"]], for k=2 shingles)
- **N-gram**: Character-level shingles (e.g., 3-grams of "hello" → ["hel", "ell", "llo""])

### 2. Weighted Hashing

1. For each feature, a standard hash (e.g., MurmurHash3) generates an **L-bit** binary vector.

### 3. Vector Aggregation

- For each of the **L** bit positions:
    - Sum weights of features with `1` in that position.
    - Subtract weights of features with `0`.
- The **sign** of the sum determines the bit value in the fingerprint:  
  $$
  \text{SimHash}_i =
  \begin{cases}
  1 & \text{if } \sum_{f} w_f \cdot h_f(i) \geq 0 \\
  0 & \text{otherwise}
  \end{cases}
  $$
  where $h_f(i)$ is the $i$-th bit of feature $f$'s hash, and $w_f$ is its weight.

### 4. Similarity Estimation

The **Hamming distance** (number of differing bits) between two SimHash fingerprints approximates their cosine
dissimilarity:  
$$
\text{Cosine Similarity} \approx 1 - \frac{\text{Hamming Distance}}{L}
$$

## Key Features

- **Efficiency**: Computes fingerprints in **O(nL)** time for $n$ features.
- **Fixed-size output**: Always generates an $L$-bit fingerprint (typical $L=64$ or $128$).
- **Robustness**: Small content changes (e.g., paraphrasing) produce small Hamming distance changes.

## Parameters

| Parameter                | Description                      | Recommended Value                            |  
|--------------------------|----------------------------------|----------------------------------------------|  
| **Feature type**         | Word tokens or character n-grams | Words for long docs, n-grams for short texts |  
| **Fingerprint size (L)** | Bit length of output hash        | `64` (balance of precision/speed)            |  
| **Weighting**            | Feature importance metric        | TF-IDF or binary (1/0)                       |  

## Usage

The minimal working example:

```go
import (
"github.com/koykov/hash/xxhash"
"github.com/koykov/pbtk/lsh/simhash"
"github.com/koykov/pbtk/shingle"
)

func main() {
hasher := xxhash.Hasher64[[]byte]{}

lsh0, _ := simhash.NewHasher[string](simhash.NewConfig[string](hasher, shingle.NewChar[string](3, "")))
_ = lsh0.Add("A sad man is crying")
vector0 := lsh0.Hash()
println(vector0) // [12297735558912453205]

lsh1, _ := simhash.NewHasher[string](simhash.NewConfig[string](hasher, shingle.NewChar[string](3, "")))
_ = lsh1.Add("A man is screaming")
vector1 := lsh1.Hash()
println(vector1) // [12297735558912453205]

// apply Hamming distance algorithm to vectors to get the similarity...
// See https://github.com/koykov/pbtk/similarity/hamming
}
```

## Use Cases

- **Duplicate detection** (e.g., crawling scraped web pages)
- **Plagiarism detection** (resilient to minor rewrites)
- **Version control** (identifying similar documents)

## Notes

- This implementation **only generates SimHash fingerprints**.
- For comparison, compute Hamming distance between fingerprints (e.g., using XOR + bit count).
- Larger $L$ improves accuracy but increases storage/memory usage.

## References

* https://matpalm.com/resemblance/simhash/
* [Locality Sensitive Hashing (LSH): The Illustrated Guide](https://www.pinecone.io/learn/series/faiss/locality-sensitive-hashing/)
