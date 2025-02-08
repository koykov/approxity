package bloom

import (
	"hash/crc32"
	"hash/crc64"
	"sync"

	"github.com/koykov/byteconv"
)

type hasherStringCRC32 struct{}

func (hasherStringCRC32) Hash(data any) uint64 {
	switch x := data.(type) {
	case string:
		return uint64(crc32.ChecksumIEEE(byteconv.S2B(x)))
	case *string:
		return uint64(crc32.ChecksumIEEE(byteconv.S2B(*x)))
	case []byte:
		return uint64(crc32.ChecksumIEEE(x))
	case *[]byte:
		return uint64(crc32.ChecksumIEEE(*x))
	}
	return 0
}

type hasherStringCRC64 struct {
	once  sync.Once
	poly  uint64
	table *crc64.Table
}

func (h *hasherStringCRC64) Hash(data any) uint64 {
	h.once.Do(func() {
		if h.poly == 0 {
			h.poly = crc64.ISO
		}
		if h.table == nil {
			h.table = crc64.MakeTable(h.poly)
		}
	})
	switch x := data.(type) {
	case string:
		return crc64.Checksum(byteconv.S2B(x), h.table)
	case *string:
		return crc64.Checksum(byteconv.S2B(*x), h.table)
	case []byte:
		return crc64.Checksum(x, h.table)
	case *[]byte:
		return crc64.Checksum(*x, h.table)
	}
	return 0
}
