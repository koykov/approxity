# Hamming Similarity

Hamming Similarity is a measure used to compare two binary strings (or binary vectors) of equal length.
It quantifies how similar two strings are by counting the number of positions where corresponding bits match.

## Math basics

Given two binary strings $s_1$ and $s_2$ of length $n$, the Hamming Similarity is calculated as:

$$
\text{Hamming Similarity} = 1 - \frac{\text{Hamming Distance}(s_1, s_2)}{n}
$$

Where:
- $Hamming Distance$ = Number of positions where bits differ
- $n$ = Length of the strings

### Properties
- Range: `[0, 1]` where:
    - `0` = No matching bits (completely dissimilar)
    - `1` = All bits match (identical)
- Symmetric: `sim(s₁,s₂) = sim(s₂,s₁)`

## Algorithm Steps

1. **Length Validation**:
    - Both input strings must have equal length
2. **Bitwise Comparison**:
    - Compare bits at each position
3. **Distance Calculation**:
    - Count mismatched positions
4. **Similarity Normalization**:
    - Divide by length and subtract from 1

## Complexity
- Time: **O(n)** where n = bit length
- Space: **O(1)** (constant additional space)

## Example Calculation

For strings:
```
s₁ = 101010
s₂ = 100110
```

1. Hamming Distance = 2 (positions 3 and 4 differ)
2. Similarity = 1 - 2/6 ≈ 0.6667 (66.67% similar)

## Usage

The minimal working example:

```go
package main

import (
	"github.com/koykov/hash/xxhash"
	"github.com/koykov/pbtk/lsh/simhash"
	"github.com/koykov/pbtk/shingle"
	"github.com/koykov/pbtk/similarity/hamming"
)

func main() {
	hasher := xxhash.Hasher64[[]byte]{}
	shingler := shingle.NewChar[[]byte](3, "") // 3-gram
	lsh, _ := simhash.NewHasher[[]byte](simhash.NewConfig[[]byte](hasher, shingler))
	est, _ := hamming.NewEstimator[[]byte](hamming.NewConfig[[]byte](lsh))

	e, _ := est.Estimate([]byte("Four children are doing backbends in the gym"), []byte("Four children are doing backbends in the park"))
	println(e) // 1.0 (high similarity)

	est.Reset()
	e, _ = est.Estimate([]byte("A man is sitting near a bike and is writing a note"), []byte("A man is standing near a bike and is writing on a piece paper"))
	println(e) // 0.9375 (medium similarity)

	est.Reset()
	e, _ = est.Estimate([]byte("One white dog and one black one are sitting side by side on the grass"), []byte("A black and a white dog are joyfully running on the grass"))
	println(e) // 0.625 (low similarity)
}

```

## Applications

1. **Text Processing**:
    - Document similarity when using binary representations (e.g., SimHash)
2. **Bioinformatics**:
    - DNA sequence comparison
3. **Computer Vision**:
    - Image fingerprint comparison
4. **Cryptography**:
    - Error detection in binary transmissions
