package policy

func BackoffSeconds(tries int) int {
	if tries <= 1 {
		return 1
	}
	n := 1
	for i := 1; i < tries && n < 64; i++ {
		n *= 2
	}
	return n
}
