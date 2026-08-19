package textutil

import "strings"

func Clip(s string, n int) string {
	s = strings.TrimSpace(s)
	if n <= 0 || len(s) <= n {
		return s
	}
	return s[len(s)-n:]
}

func NonEmpty(s string) bool { return strings.TrimSpace(s) != "" }

func SplitCSV(s string) []string {
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

func JoinCSV(items []string) string {
	clean := make([]string, 0, len(items))
	for _, it := range items {
		it = strings.TrimSpace(it)
		if it != "" {
			clean = append(clean, it)
		}
	}
	return strings.Join(clean, ",")
}
