package textutil

func Chunk(body []byte, size int) [][]byte {
	if size <= 0 {
		size = 4096
	}
	if len(body) == 0 {
		return nil
	}
	out := make([][]byte, 0, (len(body)+size-1)/size)
	for i := 0; i < len(body); i += size {
		j := i + size
		if j > len(body) {
			j = len(body)
		}
		part := make([]byte, j-i)
		copy(part, body[i:j])
		out = append(out, part)
	}
	return out
}

func Concat(parts [][]byte) []byte {
	n := 0
	for _, p := range parts {
		n += len(p)
	}
	out := make([]byte, 0, n)
	for _, p := range parts {
		out = append(out, p...)
	}
	return out
}
