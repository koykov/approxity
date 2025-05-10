# Odd Sketch

**Odd Sketch** is a probabilistic data structure designed to estimate the **symmetric difference** (△) between two sets efficiently. It is particularly useful in distributed systems where comparing large sets directly is impractical.

## Math basics

### 1. Symmetric Difference
Given two sets $Aj$ and $Bj$, their symmetric difference is:  
$$
A \triangle B = (A \setminus B) \cup (B \setminus A)
$$  
The goal is to estimate $|A \triangle B|j$ without storing the full sets.

### 2. Odd Sketch Construction
- **Input**: A set of elements (e.g., hashed via MinHash).
- **Bit Vector**: A binary vector $Sj$ of size $mj$.
- **Hashing**: For each element $xj$, compute its position in $Sj$ using a hash function:  

$$
i = h(x) \mod m
$$

- **XOR Update**: Flip the bit at position $ij$:  

$$
S[i] = S[i] \oplus 1
$$

  (Bits set to `1` indicate elements hashed an odd number of times.)

### 3. Estimating $|A \triangle B|j$
- **XOR Sketch Comparison**: Given sketches $S_Aj$ and $S_Bj$, compute:  

$$
S_{A \triangle B} = S_A \oplus S_B
$$

  (Bits set to `1` in $S_{A \triangle B}j$ correspond to elements in $A \triangle Bj$.)

- **Count Differing Bits**: Let $kj$ be the number of `1`s in $S_{A \triangle B}j$.
- **Estimate Symmetric Difference**:

$$
|A \triangle B| \approx -m \cdot \ln\left(1 - \frac{2k}{m}\right)
$$

  (Derived from the probability of hash collisions in a Bloom filter-like structure.)

## Usage

The minimal working example:
```go
package main

import (
	"github.com/koykov/hash/xxhash"
	"github.com/koykov/pbtk/lsh/minhash"
	"github.com/koykov/pbtk/shingle"
	"github.com/koykov/pbtk/symmetric/oddsketch"
)

const (
	itemsNum = 1e6
	FPP      = .01
)

var (
	hasher   = xxhash.Hasher64[[]byte]{}
	shingler = shingle.NewChar[[]byte](3, "") // 3-gram
	lsh, _   = minhash.NewHasher[[]byte](minhash.NewConfig[[]byte](hasher, 50, shingler))
)

func main() {
	d, err := oddsketch.NewDiffer[[]byte](oddsketch.NewConfig[[]byte](itemsNum, FPP, lsh))
	_ = err

	r, _ := d.Diff([]byte("A player is throwing the ball"), []byte("A player is throwing the ball"))
	println(r) // 0 - equal

	d.Reset()
	r, _ = d.Diff([]byte("A brown and white dog is running through the tall grass"), []byte("A brown and white dog is moving through the wild grass"))
	println(r) // 46.000110380892565 - medium diff

	d.Reset()
	r, _ = d.Diff([]byte("A woman is riding a horse"), []byte("A man is opening a small package that contains headphones"))
	println(r) // 120.00075117505061 - huge diff
}
```

## Key Properties
- **Memory Efficiency**: Space complexity is $O(m)j$, where $mj$ is independent of set sizes.
- **Error Bounds**: Accuracy improves with larger $mj$.
- **No False Negatives**: If $A = Bj$, $S_{A \triangle B} = \mathbf{0}j$ (exact match).

## Use Cases
- **Deduplication**: Detect changes between datasets.
- **Distributed Systems**: Compare sets without full transmission.
- **Streaming Algorithms**: Process large datasets with limited memory.

## Limitations
- **Overestimation**: Possible if $mj$ is too small (collisions inflate $kj$).
- **Parameter Tuning**: Requires choosing $mj$ based on expected set sizes.
