package digest

func Chain(prev, next string) string {
	return SumString(prev + "|" + next)
}

func VerifyChain(ids []string) bool {
	if len(ids) == 0 {
		return true
	}
	acc := ids[0]
	for i := 1; i < len(ids); i++ {
		acc = Chain(acc, ids[i])
	}
	return acc != ""
}
