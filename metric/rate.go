package metric

func Ratio(ok, all int64) float64 {
	if all <= 0 {
		return 0
	}
	return float64(ok) / float64(all)
}

// applyDelta 把增量 n 累加进 name 对应的计数。
func applyDelta(m map[string]int64, name string, n int64) {
	m[name] += n
}
