# Cosine Similarity

**Cosine Similarity** is a metric used to measure the similarity between two non-zero vectors in a multi-dimensional space. It calculates the cosine of the angle between them, providing a value between **-1** and **1**, where:
- **1** → Identical vectors (maximum similarity).
- **0** → Orthogonal vectors (no correlation).
- **-1** → Opposite vectors (maximum dissimilarity).

It is widely used in:
- **Natural Language Processing (NLP)** (document similarity, text clustering).
- **Recommendation Systems** (user/item matching).
- **Information Retrieval** (search engine ranking).
- **Computer Vision** (image embeddings comparison).

## Math basics

Given two vectors **A** and **B**, the cosine similarity is computed as:

$$
\text{Cosine Similarity}(A, B) = \frac{A \cdot B}{\|A\| \cdot \|B\|}
$$

Where:
- $A \cdot B$ → Dot product of **A** and **B**.
- $\|A\|$ and $\|B\|$ → Euclidean norms (magnitudes) of the vectors.

To calculate vectors from text data use [LSH](../../lsh) package.

### Key Properties
1. **Scale-Invariant**: Only considers vector direction, not magnitude.
2. **Efficient for Sparse Data**: Works well with high-dimensional vectors (e.g., TF-IDF, word embeddings).
3. **Bounded Output**: Always returns a value in **[-1, 1]**.

## Requirements
- Vectors **must be of the same dimensionality** (if not, pad the shorter one with zeros).
- Non-zero vectors (to avoid division by zero).

## Advantages & Limitations

| **Advantages** | **Limitations** |
|---------------|----------------|
| Robust to magnitude differences | Sensitive to feature alignment (requires consistent vector semantics) |  
| Works well in high-dimensional spaces | Not suitable for non-vector data (e.g., graphs) |  
| Fast to compute | May perform poorly if vectors are not normalized |  

## Usage

The minimal working example:

```go
import (
    "github.com/koykov/hash/xxhash"
    "github.com/koykov/pbtk/lsh/minhash"
    "github.com/koykov/pbtk/shingle"
    "github.com/koykov/pbtk/similarity/cosine"
)

func main() {
    hasher := xxhash.Hasher64[[]byte]{}
    shingler := shingle.NewChar[[]byte](3, "") // 3-gram
    lsh, _ := minhash.NewHasher[[]byte](minhash.NewConfig[[]byte](hasher, 50, shingler))
    est, _ := cosine.NewEstimator[[]byte](cosine.NewConfig[[]byte](lsh))
    
    e, _ := est.Estimate([]byte("Four children are doing backbends in the gym"), []byte("Four children are doing backbends in the park"))
    println(e) // 0.9953444221469755 (high similarity)
    
    est.Reset()
    e, _ = est.Estimate([]byte("A man is sitting near a bike and is writing a note"), []byte("A man is standing near a bike and is writing on a piece paper"))
    println(e) // 0.49682303908308706 (medium similarity)
    
    est.Reset()
    e, _ = est.Estimate([]byte("One white dog and one black one are sitting side by side on the grass"), []byte("A black and a white dog are joyfully running on the grass"))
    println(e) // 0.09107270108146269 (low similarity)
}
```

## Use cases

* Comparing text documents (e.g., search engines).
* Recommending similar items (e.g., collaborative filtering).
* Clustering or classification tasks where direction matters more than magnitude.
