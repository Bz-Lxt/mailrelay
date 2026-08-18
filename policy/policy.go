package policy

func MaxTries(n int) int {
	if n <= 0 {
		return 5
	}
	if n > 32 {
		return 32
	}
	return n
}

func ShouldBounce(tries, limit int) bool { return tries >= MaxTries(limit) }
