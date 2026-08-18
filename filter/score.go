package filter

func Score(body string) int {
	n := 0
	for i := 0; i < len(body); i++ {
		c := body[i]
		if c >= 'A' && c <= 'Z' {
			n++
		}
		if c == '!' {
			n += 2
		}
	}
	if n > 100 {
		return 100
	}
	return n
}

func OverLimit(score, limit int) bool {
	if limit <= 0 {
		limit = 80
	}
	return score >= limit
}
