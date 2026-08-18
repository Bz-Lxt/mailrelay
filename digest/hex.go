package digest

import (
	"encoding/hex"
	"fmt"
	"strings"
)

func Encode(b []byte) string { return hex.EncodeToString(b) }

func Decode(s string) ([]byte, error) {
	s = strings.TrimSpace(strings.ToLower(s))
	if len(s)%2 != 0 {
		return nil, fmt.Errorf("odd hex length")
	}
	return hex.DecodeString(s)
}

func Short(s string, n int) string {
	if n <= 0 || len(s) <= n {
		return s
	}
	return s[:n]
}

func EqualFold(a, b string) bool {
	return strings.EqualFold(strings.TrimSpace(a), strings.TrimSpace(b))
}
