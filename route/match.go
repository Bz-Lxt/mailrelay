package route

import "strings"

func MatchDomain(addr, domain string) bool {
	addr = strings.ToLower(addr)
	domain = strings.ToLower(domain)
	if !strings.Contains(addr, "@") {
		return false
	}
	return strings.HasSuffix(addr, "@"+domain) || strings.HasSuffix(addr, "."+domain)
}

func LocalPart(addr string) string {
	i := strings.IndexByte(addr, '@')
	if i < 0 {
		return addr
	}
	return addr[:i]
}
