package mimeutil

import "strings"

func QuotePrintable(s string) string {
	var b strings.Builder
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == ' ' || (c >= 33 && c <= 126 && c != '=') {
			b.WriteByte(c)
			continue
		}
		b.WriteByte('=')
		b.WriteByte("0123456789ABCDEF"[c>>4])
		b.WriteByte("0123456789ABCDEF"[c&15])
	}
	return b.String()
}

func UnwrapSubject(s string) string {
	s = strings.TrimSpace(s)
	for strings.HasPrefix(strings.ToLower(s), "re:") || strings.HasPrefix(strings.ToLower(s), "fwd:") {
		if i := strings.IndexByte(s, ':'); i >= 0 {
			s = strings.TrimSpace(s[i+1:])
		} else {
			break
		}
	}
	return s
}

func WordCount(s string) int {
	n := 0
	in := false
	for _, r := range s {
		if r == ' ' || r == '\n' || r == '\t' {
			in = false
			continue
		}
		if !in {
			n++
			in = true
		}
	}
	return n
}
