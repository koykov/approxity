# MinHash

**MinHash** is a probabilistic algorithm designed to efficiently estimate the similarity between two sets (e.g., text
documents, user preferences, or web pages) by computing compact signatures. This implementation focuses on generating a
**MinHash signature**—a fixed-length vector of hash minima—which can later be used to approximate the **Jaccard
similarity coefficient** between sets.

The Jaccard similarity between two sets $A$ and $B$ is defined as:

$$
J(A, B) = \frac{|A \cap B|}{|A \cup B|}
$$

The MinHash signature allows approximating $J(A, B)$ by comparing the fraction of matching hash minima between the two
signatures.

## How It Works

### 1. Shingling (Text → Set of Shingles)

The input string is split into overlapping **shingles** (n-grams) of fixed length (e.g., 3 characters).

**Example**:  
For the string `"hello"` and shingle length `3`, the shingles are:  
`["hel", "ell", "llo"]`.

### 2. Hashing Shingles

For each shingle, `k` independent hash values are computed using a hash function (e.g., MurmurHash3) with different seed
values (`0` to `k-1`).

For a shingle $s$, the hash values are:  
$$
h_0(s), h_1(s), \dots, h_{k-1}(s)
$$

### 3. Building the MinHash Signature

For each of the `k` hash functions, the **minimum hash value** across all shingles is selected. The resulting signature
is a vector of `k` minima:

$$
\text{MinHash}(A) = \left[ \min_{s \in A} h_0(s), \min_{s \in A} h_1(s), \dots, \min_{s \in A} h_{k-1}(s) \right]
$$

### 4. Using the Signature for Jaccard Estimation

The MinHash signature can be passed to a **Jaccard similarity estimator**, where the similarity between two sets $A$
and $B$ is approximated by:

$$
J(A, B) \approx \frac{\text{Number of matching minima in } \text{MinHash}(A) \text{ and } \text{MinHash}(B)}{k}
$$

## Key Features

- **Efficiency**: Computes signatures in **O(n × k)**, where `n` is the number of shingles.
- **Fixed-size signature**: The output is always a vector of `k` integers, regardless of input size.
- **Probabilistic guarantees**: The estimate converges to the true Jaccard similarity as `k` increases.

## Parameters

| Description                                    | Recommended Value                    | Parameter          |  
|------------------------------------------------|--------------------------------------|--------------------|
| Size of n-grams (characters or words).         | `3`–`5` for chars, `1`–`2` for words | **Shingle length** |  
| Number of hash functions (controls precision). | `50`–`200` (higher = more accurate)  | **k**              |

## Usage

The minimal working example:

```go
import (
"github.com/koykov/hash/xxhash"
"github.com/koykov/pbtk/lsh/minhash"
"github.com/koykov/pbtk/shingle"
)

func main() {
hasher := xxhash.Hasher64[[]byte]{} // prepare hash function

const k = 50 // k == 50 hash functions

// init hash with 3-gram shingler
lsh0, _ := minhash.NewHasher[string](minhash.NewConfig[string](hasher, k, shingle.NewChar[string](3, "")))
_ = lsh0.Add("hello world")
vector0 := lsh0.Hash()
println(vector0) // [310722046221163110 644366121911734340 318615316911410328 145560956935501306 59683552839790942 98502252458282848 141515789415006756 293867436433651891 171773667406779622]

// init hash with 3-gram shingler
lsh1, _ := minhash.NewHasher[string](minhash.NewConfig[string](hasher, k, shingle.NewChar[string](3, "")))
_ = lsh1.Add("hello there")
vector1 := lsh1.Hash()
println(vector1) // [310722046221163110 644366121911734340 318615316911410328 145560956935501306 19320324343221180 235974345174438171 37615032563520863 53237528186642853 79373271539028071]

// apply Jaccard algorithm to vectors to get the similarity...
// See https://github.com/koykov/pbtk/similarity/jaccard
}
```

## Use Cases

- **Near-duplicate detection** (e.g., clustering similar documents).
- **Recommendation systems** (finding users with similar preferences).
- **Plagiarism detection** (comparing text similarity).

## Notes

- This implementation **only generates MinHash signatures** and does not compute Jaccard similarity directly.
- For comparing two signatures, use:  
  $$
  \text{Jaccard Estimate} = \frac{\text{Number of matching minima}}{k}
  $$
- Increasing `k` improves accuracy but requires more computation.

## References

* https://en.wikipedia.org/wiki/MinHash
* [Locality Sensitive Hashing (LSH): The Illustrated Guide](https://www.pinecone.io/learn/series/faiss/locality-sensitive-hashing/)
