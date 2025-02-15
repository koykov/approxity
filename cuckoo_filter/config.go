package cuckoo

type Config struct {
	Buckets         uint64
	BucketSize      uint64
	FingerprintSize uint64
	KicksLimit      uint64
}
