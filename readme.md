# Probabilistic Toolkit

A collection of probabilistic data structures and algorithms. Solves the following tasks:

* AMQ (Approximate Membership Query)
* Cardinality estimation
* Frequency estimation
* Similarity estimation
* Symmetric difference
* LSH (Locality-Sensitive Hashing)

All solutions are designed for high-load environments and provide the following features:

* Compact structures with minimal memory footprint
* Minimal or zero memory allocations
* Lock-free operations via atomic operations
* Concurrency mode support for multithreaded environments
* SIMD optimizations
* Flexible initialization (all auxiliary structures are abstracted; e.g., any hash algorithm can be specified)
* Each structure implements a unified interface (within its problem domain), allowing easy switching between implementations
* Built-in metrics coverage

Full implementations tree:

* [AMQ](amq)
    * [Bloom filter](amq/bloom_filter)
    * [Counting Bloom filter](amq/bloom_filter)
    * [Cuckoo filter](amq/cuckoo_filter)
    * [Quotient filter](amq/quotient_filter)
    * [Xor filter](amq/xor_filter)
* [Cardinality estimation](cardinality)
    * [LogLog](cardinality/loglog)
    * [HyperLogLog](cardinality/hyperloglog)
    * [HyperBitBit](cardinality/hyperbitbit)
    * [Linear counting](cardinality/linear_counting)
* [Frequency estimation](frequency)
    * [Count-Min Sketch](frequency/cmsketch)
    * [Conservative Update Sketch](frequency/cusketch)
    * [Count Sketch](frequency/countsketch)
    * [TinyLFU](frequency/tinylfu)
    * [TinyLFU (EWMA)](frequency/tinylfu_ewma)
* [Similarity estimation](similarity)
    * [Cosine similarity](similarity/cosine)
    * [Jaccard similarity](similarity/jaccard)
    * [Hamming similarity](similarity/hamming)
* [Symmetric difference](symmetric)
    * [Odd Sketch](symmetric/oddsketch)
* [LSH](lsh)
    * [SimHash](lsh/simhash)
    * [MinHash](lsh/minhash)
    * [b-Bit MinHash](lsh/bbitminhash)
* [Shingle](shingle)
    * [Char](shingle/char.go)
    * [Word](shingle/word.go)

Below is a brief description of each task. For algorithm details, refer to the corresponding sections.

## AMQ (Approximate Membership Query)

AMQ structures solve the *membership testing* problem—determining whether a key belongs to a set. While hash tables
can solve this task, they are only feasible for small sets. AMQ structures, in contrast, can handle very large sets
at the cost of precision: false positives are possible, but false negatives are not.

[Detailed description](amq).

## Cardinality Estimation

Cardinality estimation structures solve the problem of counting unique keys in a set. Hash tables can also solve this task,
but their memory usage becomes prohibitive for large sets. Cardinality estimation structures minimize memory usage
while providing approximate results.

[Detailed description](cardinality).

## Frequency Estimation

Frequency estimation structures determine the frequency of keys in a set. Like other probabilistic structures,
they trade precision for minimal memory consumption.

[Detailed description](frequency).

## Similarity Estimation

Similarity estimation structures measure the similarity between two sets. In this package, sets are represented as string data,
enabling fuzzy text comparison. These algorithms require an auxiliary LSH structure and handle all stages
(shingling, hashing, and vectorization) transparently.

[Detailed description](similarity).

## LSH (Locality-Sensitive Hashing)

LSH is a nearest-neighbor search method that hashes input data so that similar items likely collide into the same bucket.
In this package, LSH operates only on string data and performs vectorization of pre-shingled texts.
The resulting vectors are then passed to similarity estimation structures.  

In practice, LSH is initialized with a Shingle algorithm and handles shingling transparently.

[Detailed description](lsh).

## Shingle

Shingling is a text tokenization method. This package implements character-based and word-based shingling.
Shingled texts are then passed to LSH for vectorization and further similarity estimation.

[Detailed description](shingle).

## Symmetric Difference

Symmetric difference structures measure the symmetric difference between two sets. This package implements algorithms for text data.
Symmetric difference is inversely related to similarity estimation: the more similar two texts are, the smaller their symmetric difference.
Like similarity estimation, these structures require an auxiliary LSH.

[Detailed description](symmetric).

## Conclusion

The implemented structures enable real-time analysis of large datasets or data streams with minimal resource usage and optimal performance.
Abstraction layers allow seamless switching between algorithms to choose the best fit for a task.
Concurrency mode support ensures lock-free operation in multithreaded environments, while built-in metrics help evaluate
and fine-tune configurations.
