package mimeutil

import (
	"strings"
	"unicode"
)

func Fold(name, value string) string {
	name = strings.TrimSpace(name)
	value = strings.TrimSpace(value)
	if name == "" {
		return value
	}
	return name + ": " + value
}

func Unfold(line string) (string, string) {
	i := strings.IndexByte(line, ':')
	if i < 0 {
		return "", strings.TrimSpace(line)
	}
	return strings.TrimSpace(line[:i]), strings.TrimSpace(line[i+1:])
}

func Sanitize(s string) string {
	var b strings.Builder
	for _, r := range s {
		if r == '\n' || r == '\r' {
			continue
		}
		if unicode.IsControl(r) {
			continue
		}
		b.WriteRune(r)
	}
	return strings.TrimSpace(b.String())
}

func SplitAddress(addr string) (local, domain string) {
	addr = strings.TrimSpace(addr)
	i := strings.LastIndexByte(addr, '@')
	if i < 0 {
		return addr, ""
	}
	return addr[:i], addr[i+1:]
}

func ValidAddress(addr string) bool {
	local, domain := SplitAddress(addr)
	return local != "" && domain != "" && strings.Contains(domain, ".")
}
