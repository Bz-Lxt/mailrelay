package config

import "strings"

func ValidDir(dir string) bool {
	dir = strings.TrimSpace(dir)
	if dir == "" || dir == "/" {
		return false
	}
	return !strings.Contains(dir, "\x00")
}

func ValidAddr(addr string) bool {
	return strings.Contains(addr, ":") && !strings.Contains(addr, " ")
}
