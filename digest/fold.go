package digest

import "strings"

func NormalizeID(raw string) string {
	raw = strings.TrimSpace(strings.ToLower(raw))
	out := make([]byte, 0, len(raw))
	for i := 0; i < len(raw); i++ {
		c := raw[i]
		if (c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') {
			out = append(out, c)
		}
	}
	return string(out)
}

func PadHex(s string, n int) string {
	if len(s) >= n {
		return s[:n]
	}
	return s + strings.Repeat("0", n-len(s))
}
