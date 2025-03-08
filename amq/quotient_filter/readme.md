# Quotient filter

> [!CAUTION]
> Not production ready!
> This implementation eats too much memory and doesn't support simultaneous read/write operations.

Quotient filters use a hash table with quotienting, where each item’s hash is split into two parts:
* a quotient, which determines the bucket.
* a remainder, which is stored in the bucket.

The filter uses three metadata bits (is_occupied, is_continuation, is_shifted) per bucket to handle collisions and
maintain the logical order of elements.

See [full description](https://en.wikipedia.org/wiki/Quotient_filter).

## Usage

```go
import (
	"github.com/koykov/approxity/amq/quotient_filter"
	"github.com/koykov/hash/xxhash"
)

func main() {
	f, err := quotient.NewFilter[string](quotient.NewConfig(1e3, 0.01, xxhash.Hasher64[[]byte]{}))
    _ = err
	_ = f.Set("foobar")
	print(f.Contains("foobar")) // true
    print(f.Contains("qwerty")) // false
}
```

### Optimal params calculation

There is no need to calculate optimal size `m` due to filter makes it itself using desired number of items (`Config.ItemsNumber`),
false positive probability (`Config.FPP`) and load factor (`Config.LoadFactor`) params.
