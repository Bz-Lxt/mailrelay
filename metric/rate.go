package metric

func Ratio(ok, all int64) float64 {
	if all <= 0 {
		return 0
	}
	return float64(ok) / float64(all)
}
