package hasher

import (
	"hash/crc64"
	"sync"
)

type CRC64 struct {
	once  sync.Once
	poly  uint64
	table *crc64.Table
}

func (h *CRC64) Sum64(data []byte) uint64 {
	h.once.Do(func() {
		if h.poly == 0 {
			h.poly = crc64.ISO
		}
		if h.table == nil {
			h.table = crc64.MakeTable(h.poly)
		}
	})
	return crc64.Checksum(data, h.table)
}
