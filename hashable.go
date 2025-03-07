package approxity

// Hashable is a constraint that permits any type that can be hashed.
type Hashable interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64 |
		~string | ~[]byte | ~[]rune | ~bool
}
