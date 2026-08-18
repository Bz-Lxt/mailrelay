package retrybook

func Delay(tries int) int {
	if tries <= 0 {
		return 1
	}
	n := 1
	for i := 0; i < tries && n < 128; i++ {
		n *= 2
	}
	return n
}

func Jitter(base, salt int) int {
	if base <= 0 {
		return 1
	}
	if salt < 0 {
		salt = -salt
	}
	return base + salt%base
}
