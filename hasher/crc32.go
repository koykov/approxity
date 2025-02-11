package hasher

import (
	"hash/crc32"
)

type CRC32 struct{}

func (CRC32) Sum64(data []byte) uint64 {
	return uint64(crc32.ChecksumIEEE(data))
}
