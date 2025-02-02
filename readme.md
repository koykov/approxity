# Bloom filter

A Bloom filter is a probabilistic data structure that allows for fast and space-efficient membership testing.

## Usage

```go
import "github.com/koykov/bloom"

f := bloom.New(100, 10)
f.Add("hello")
f.Add("world")

if f.Has("hello") {
    fmt.Println("hello is in the filter")
}
```
