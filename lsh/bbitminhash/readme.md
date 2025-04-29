# b-Bit MinHash

**b-Bit MinHash** is a memory-efficient variant of the **MinHash** algorithm, designed to estimate the **Jaccard
similarity**
between two sets while significantly reducing storage requirements.

Let's remember how MinHash works:

MinHash is a probabilistic data structure used to estimate the Jaccard similarity between two sets $A$ and $B$:

$$
J(A, B) = \frac{|A \cap B|}{|A \cup B|}
$$

**Steps:**

1. Apply $k$ independent hash functions to each element in the set.
2. For each hash function, keep the **minimum hash value** observed.
3. The probability that the MinHash values of two sets match equals their Jaccard similarity.

**Problem:**  
Storing $k$ full hash values (typically 32 or 64 bits each) consumes significant memory.

## **How b-Bit MinHash Improves Efficiency**

Instead of storing the entire hash, **b-Bit MinHash** keeps only the **lowest $b$ bits** of each MinHash (
usually $b = 1, 2, 4$).

### **Key Differences from MinHash**

| Feature      | MinHash                | b-Bit MinHash                         |
|--------------|------------------------|---------------------------------------|
| **Storage**  | Full hash (32/64 bits) | Only $b$ bits per hash                |
| **Accuracy** | High                   | Slightly reduced (adjustable via $b$) |
| **Memory**   | High                   | **Dramatically lower**                |

## **Memory Savings Example**

Suppose we use:

- $k = 50$ hash functions
- Original MinHash: **64 bits per hash** → **3200 bits (400 bytes)**
- b-Bit MinHash ($b=4$): **4 bit per hash** → **200 bits (25 bytes)**

**Reduction:** **16x less memory!**

## Usage

The minimal working example:

```go
import (
"github.com/koykov/hash/xxhash"
"github.com/koykov/pbtk/lsh/bbitminhash"
"github.com/koykov/pbtk/shingle"
)

func main() {
hasher := xxhash.Hasher64[[]byte]{} // prepare hash function

const k = 50 // k == 50 hash functions
const b = 7  // b == 7 bits per hash

// init hash with 3-gram shingler
lsh0, _ := bbitminhash.NewHasher[string](bbitminhash.NewConfig[string](hasher, k, shingle.NewChar[string](3, ""), b))
_ = lsh0.Add("hello world")
vector0 := lsh0.Hash()
println(vector0) // [0 0 127 127 127 127 127 127 0]

// init hash with 3-gram shingler
lsh1, _ := bbitminhash.NewHasher[string](bbitminhash.NewConfig[string](hasher, k, shingle.NewChar[string](3, "")))
_ = lsh1.Add("hello there")
vector1 := lsh1.Hash()
println(vector1) // [0 0 127 127 127 127 127 0 0]

// apply Jaccard algorithm to vectors to get the similarity...
// See https://github.com/koykov/pbtk/similarity/jaccard
}
```

## **When to Use b-Bit MinHash?**

**Best for:**

- Large-scale similarity searches (e.g., near-duplicate detection)
- Applications where memory efficiency is critical

**Not ideal when:**

- Extremely high accuracy is required (use full MinHash)
- $b$ is too small (e.g., $b=1$ for very dissimilar sets)  
