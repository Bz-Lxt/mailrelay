package idgen

func Many(prefix string, n int) []string {
	if n <= 0 {
		return nil
	}
	out := make([]string, 0, n)
	for i := 0; i < n; i++ {
		out = append(out, Next(prefix, nil))
	}
	return out
}
