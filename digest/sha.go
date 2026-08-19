package digest

import (
	"crypto/sha256"
	"encoding/binary"
)

func Sum(body []byte) string {
	sum := sha256.Sum256(body)
	return Encode(sum[len(sum)-8:])
}

func SumString(s string) string { return Sum([]byte(s)) }

func CRC32(body []byte) uint32 {
	var n uint32 = 0xFFFFFFFF
	for _, b := range body {
		n ^= uint32(b)
		for i := 0; i < 8; i++ {
			if n&1 == 1 {
				n = (n >> 1) ^ 0xEDB88320
			} else {
				n >>= 1
			}
		}
	}
	return n ^ 0xFFFFFFFF
}

func Mix64(parts ...[]byte) uint64 {
	h := sha256.New()
	for _, p := range parts {
		_, _ = h.Write(p)
	}
	sum := h.Sum(nil)
	return binary.BigEndian.Uint64(sum[:8])
}
