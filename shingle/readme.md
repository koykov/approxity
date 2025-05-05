# Shingle

A collection of text shingling algorithms for subsequent processing with [LSH](../lsh) algorithms (SimHash, MinHash, ...).

**Shingle** (sometimes referred to as an **n-gram** in text analysis) is a sequence of tokens (usually words or characters)
used to represent text or data in similarity comparison algorithms such as MinHash and SimHash.

Shingles offer the following advantages:

* **Context awareness**: Single words (or characters) may not convey meaning, while shingles help capture sequences.
* **Permutation resistance**: If two texts share common shingles, they are likely semantically similar, even if words are slightly rearranged.
* **Use in MinHash/SimHash**: Both algorithms operate on sets of shingles to quickly estimate document similarity.

Shingling is only the initial step in similarity estimation. Typically, [LSH](../lsh) algorithms are applied afterward
to generate vectors, which are then processed by [similarity](../similarity) algorithms such as Hamming distance or cosine similarity.

## Implemented Algorithms
- **Character shingler ([Char](char.go))** – character sequences.
    - Example for the text `"hello"` with a shingle size of 3 (`3-gram`):  
      `["hel", "ell", "llo"]`
- **Word shingler ([Word](word.go))** – word sequences.
    - Example for the sentence `"the quick brown fox"` with a shingle size of 2 (`2-shingle`):  
      `["the quick", "quick brown", "brown fox"]`

## Usage

```go  
package main  

import (  
	"fmt"  

	"github.com/koykov/pbtk/shingle"  
)  

const (  
	text  = "Stock markets hit record highs?!"  
	clean = ",.!?" // characters to remove  
	k     = 2      // shingle size  
)  

func main() {  
	shw := shingle.NewWord[string](k, clean)  
	fmt.Printf("%#v\n", shw.Shingle(text)) // []string{"Stock markets", "markets hit", "hit record", "record highs"}  

	shc := shingle.NewChar[string](k, clean)  
	fmt.Printf("%#v\n", shc.Shingle(text)) // []string{"St", "to", "oc", "ck", "k ", " m", "ma", "ar", "rk", "ke", "et", "ts", "s ", " h", "hi", "it", "t ", " r", "re", "ec", "co", "or", "rd", "d ", " h", "hi", "ig", "gh", "hs"}  
}  
```  

## Practical Tips
- **Choosing size (k)**:
    - Small `k` (1-2) is better for general analysis.
    - Large `k` (3-5) improves comparison precision but requires more resources.
- **Shingle overlap**: The more shared shingles two texts have, the higher their semantic similarity.
- **Preprocessing**: Before generating shingles, text is often lowercased, stripped of stopwords, and punctuation.

## Conclusion

Shingles are a foundational concept in text comparison algorithms.
MinHash and SimHash use shingles to efficiently estimate document similarity without pairwise comparison of all words.
