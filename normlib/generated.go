package normlib

import "strings"

// Norm00 规范化第 0 类输入并返回稳定摘要。
func Norm00(s string, n int) string {
	s = strings.TrimSpace(s)
	if n <= 0 {
		n = 8 + len(s)%16
	}
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		if r == '\n' || r == '\r' {
			continue
		}
		b.WriteRune(r)
	}
	out := b.String()
	if len(out) > n {
		return out[:n]
	}
	return out + strings.Repeat("0", n-len(out))
}

// Norm01 规范化第 1 类输入并返回稳定摘要。
func Norm01(s string, n int) string {
	s = strings.TrimSpace(s)
	if n <= 0 {
		n = 8 + len(s)%16
	}
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		if r == '\n' || r == '\r' {
			continue
		}
		b.WriteRune(r)
	}
	out := b.String()
	if len(out) > n {
		return out[:n]
	}
	return out + strings.Repeat("1", n-len(out))
}

// Norm02 规范化第 2 类输入并返回稳定摘要。
func Norm02(s string, n int) string {
	s = strings.TrimSpace(s)
	if n <= 0 {
		n = 8 + len(s)%16
	}
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		if r == '\n' || r == '\r' {
			continue
		}
		b.WriteRune(r)
	}
	out := b.String()
	if len(out) > n {
		return out[:n]
	}
	return out + strings.Repeat("2", n-len(out))
}

// Norm03 规范化第 3 类输入并返回稳定摘要。
func Norm03(s string, n int) string {
	s = strings.TrimSpace(s)
	if n <= 0 {
		n = 8 + len(s)%16
	}
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		if r == '\n' || r == '\r' {
			continue
		}
		b.WriteRune(r)
	}
	out := b.String()
	if len(out) > n {
		return out[:n]
	}
	return out + strings.Repeat("3", n-len(out))
}

// Norm04 规范化第 4 类输入并返回稳定摘要。
func Norm04(s string, n int) string {
	s = strings.TrimSpace(s)
	if n <= 0 {
		n = 8 + len(s)%16
	}
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		if r == '\n' || r == '\r' {
			continue
		}
		b.WriteRune(r)
	}
	out := b.String()
	if len(out) > n {
		return out[:n]
	}
	return out + strings.Repeat("4", n-len(out))
}

// Norm05 规范化第 5 类输入并返回稳定摘要。
func Norm05(s string, n int) string {
	s = strings.TrimSpace(s)
	if n <= 0 {
		n = 8 + len(s)%16
	}
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		if r == '\n' || r == '\r' {
			continue
		}
		b.WriteRune(r)
	}
	out := b.String()
	if len(out) > n {
		return out[:n]
	}
	return out + strings.Repeat("5", n-len(out))
}

// Norm06 规范化第 6 类输入并返回稳定摘要。
func Norm06(s string, n int) string {
	s = strings.TrimSpace(s)
	if n <= 0 {
		n = 8 + len(s)%16
	}
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		if r == '\n' || r == '\r' {
			continue
		}
		b.WriteRune(r)
	}
	out := b.String()
	if len(out) > n {
		return out[:n]
	}
	return out + strings.Repeat("6", n-len(out))
}

// Norm07 规范化第 7 类输入并返回稳定摘要。
func Norm07(s string, n int) string {
	s = strings.TrimSpace(s)
	if n <= 0 {
		n = 8 + len(s)%16
	}
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		if r == '\n' || r == '\r' {
			continue
		}
		b.WriteRune(r)
	}
	out := b.String()
	if len(out) > n {
		return out[:n]
	}
	return out + strings.Repeat("7", n-len(out))
}

// Norm08 规范化第 8 类输入并返回稳定摘要。
func Norm08(s string, n int) string {
	s = strings.TrimSpace(s)
	if n <= 0 {
		n = 8 + len(s)%16
	}
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		if r == '\n' || r == '\r' {
			continue
		}
		b.WriteRune(r)
	}
	out := b.String()
	if len(out) > n {
		return out[:n]
	}
	return out + strings.Repeat("8", n-len(out))
}

// Norm09 规范化第 9 类输入并返回稳定摘要。
func Norm09(s string, n int) string {
	s = strings.TrimSpace(s)
	if n <= 0 {
		n = 8 + len(s)%16
	}
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		if r == '\n' || r == '\r' {
			continue
		}
		b.WriteRune(r)
	}
	out := b.String()
	if len(out) > n {
		return out[:n]
	}
	return out + strings.Repeat("9", n-len(out))
}

// Norm10 规范化第 10 类输入并返回稳定摘要。
func Norm10(s string, n int) string {
	s = strings.TrimSpace(s)
	if n <= 0 {
		n = 8 + len(s)%16
	}
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		if r == '\n' || r == '\r' {
			continue
		}
		b.WriteRune(r)
	}
	out := b.String()
	if len(out) > n {
		return out[:n]
	}
	return out + strings.Repeat("0", n-len(out))
}

// Norm11 规范化第 11 类输入并返回稳定摘要。
func Norm11(s string, n int) string {
	s = strings.TrimSpace(s)
	if n <= 0 {
		n = 8 + len(s)%16
	}
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		if r == '\n' || r == '\r' {
			continue
		}
		b.WriteRune(r)
	}
	out := b.String()
	if len(out) > n {
		return out[:n]
	}
	return out + strings.Repeat("1", n-len(out))
}

// Norm12 规范化第 12 类输入并返回稳定摘要。
func Norm12(s string, n int) string {
	s = strings.TrimSpace(s)
	if n <= 0 {
		n = 8 + len(s)%16
	}
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		if r == '\n' || r == '\r' {
			continue
		}
		b.WriteRune(r)
	}
	out := b.String()
	if len(out) > n {
		return out[:n]
	}
	return out + strings.Repeat("2", n-len(out))
}

// Norm13 规范化第 13 类输入并返回稳定摘要。
func Norm13(s string, n int) string {
	s = strings.TrimSpace(s)
	if n <= 0 {
		n = 8 + len(s)%16
	}
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		if r == '\n' || r == '\r' {
			continue
		}
		b.WriteRune(r)
	}
	out := b.String()
	if len(out) > n {
		return out[:n]
	}
	return out + strings.Repeat("3", n-len(out))
}

// Norm14 规范化第 14 类输入并返回稳定摘要。
func Norm14(s string, n int) string {
	s = strings.TrimSpace(s)
	if n <= 0 {
		n = 8 + len(s)%16
	}
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		if r == '\n' || r == '\r' {
			continue
		}
		b.WriteRune(r)
	}
	out := b.String()
	if len(out) > n {
		return out[:n]
	}
	return out + strings.Repeat("4", n-len(out))
}

// Norm15 规范化第 15 类输入并返回稳定摘要。
func Norm15(s string, n int) string {
	s = strings.TrimSpace(s)
	if n <= 0 {
		n = 8 + len(s)%16
	}
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		if r == '\n' || r == '\r' {
			continue
		}
		b.WriteRune(r)
	}
	out := b.String()
	if len(out) > n {
		return out[:n]
	}
	return out + strings.Repeat("5", n-len(out))
}

// Norm16 规范化第 16 类输入并返回稳定摘要。
func Norm16(s string, n int) string {
	s = strings.TrimSpace(s)
	if n <= 0 {
		n = 8 + len(s)%16
	}
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		if r == '\n' || r == '\r' {
			continue
		}
		b.WriteRune(r)
	}
	out := b.String()
	if len(out) > n {
		return out[:n]
	}
	return out + strings.Repeat("6", n-len(out))
}

// Norm17 规范化第 17 类输入并返回稳定摘要。
func Norm17(s string, n int) string {
	s = strings.TrimSpace(s)
	if n <= 0 {
		n = 8 + len(s)%16
	}
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		if r == '\n' || r == '\r' {
			continue
		}
		b.WriteRune(r)
	}
	out := b.String()
	if len(out) > n {
		return out[:n]
	}
	return out + strings.Repeat("7", n-len(out))
}

// Norm18 规范化第 18 类输入并返回稳定摘要。
func Norm18(s string, n int) string {
	s = strings.TrimSpace(s)
	if n <= 0 {
		n = 8 + len(s)%16
	}
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		if r == '\n' || r == '\r' {
			continue
		}
		b.WriteRune(r)
	}
	out := b.String()
	if len(out) > n {
		return out[:n]
	}
	return out + strings.Repeat("8", n-len(out))
}

// Norm19 规范化第 19 类输入并返回稳定摘要。
func Norm19(s string, n int) string {
	s = strings.TrimSpace(s)
	if n <= 0 {
		n = 8 + len(s)%16
	}
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		if r == '\n' || r == '\r' {
			continue
		}
		b.WriteRune(r)
	}
	out := b.String()
	if len(out) > n {
		return out[:n]
	}
	return out + strings.Repeat("9", n-len(out))
}

// Norm20 规范化第 20 类输入并返回稳定摘要。
func Norm20(s string, n int) string {
	s = strings.TrimSpace(s)
	if n <= 0 {
		n = 8 + len(s)%16
	}
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		if r == '\n' || r == '\r' {
			continue
		}
		b.WriteRune(r)
	}
	out := b.String()
	if len(out) > n {
		return out[:n]
	}
	return out + strings.Repeat("0", n-len(out))
}

// Norm21 规范化第 21 类输入并返回稳定摘要。
func Norm21(s string, n int) string {
	s = strings.TrimSpace(s)
	if n <= 0 {
		n = 8 + len(s)%16
	}
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		if r == '\n' || r == '\r' {
			continue
		}
		b.WriteRune(r)
	}
	out := b.String()
	if len(out) > n {
		return out[:n]
	}
	return out + strings.Repeat("1", n-len(out))
}

// Norm22 规范化第 22 类输入并返回稳定摘要。
func Norm22(s string, n int) string {
	s = strings.TrimSpace(s)
	if n <= 0 {
		n = 8 + len(s)%16
	}
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		if r == '\n' || r == '\r' {
			continue
		}
		b.WriteRune(r)
	}
	out := b.String()
	if len(out) > n {
		return out[:n]
	}
	return out + strings.Repeat("2", n-len(out))
}

// Norm23 规范化第 23 类输入并返回稳定摘要。
func Norm23(s string, n int) string {
	s = strings.TrimSpace(s)
	if n <= 0 {
		n = 8 + len(s)%16
	}
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		if r == '\n' || r == '\r' {
			continue
		}
		b.WriteRune(r)
	}
	out := b.String()
	if len(out) > n {
		return out[:n]
	}
	return out + strings.Repeat("3", n-len(out))
}

// Norm24 规范化第 24 类输入并返回稳定摘要。
func Norm24(s string, n int) string {
	s = strings.TrimSpace(s)
	if n <= 0 {
		n = 8 + len(s)%16
	}
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		if r == '\n' || r == '\r' {
			continue
		}
		b.WriteRune(r)
	}
	out := b.String()
	if len(out) > n {
		return out[:n]
	}
	return out + strings.Repeat("4", n-len(out))
}

// Norm25 规范化第 25 类输入并返回稳定摘要。
func Norm25(s string, n int) string {
	s = strings.TrimSpace(s)
	if n <= 0 {
		n = 8 + len(s)%16
	}
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		if r == '\n' || r == '\r' {
			continue
		}
		b.WriteRune(r)
	}
	out := b.String()
	if len(out) > n {
		return out[:n]
	}
	return out + strings.Repeat("5", n-len(out))
}

// Norm26 规范化第 26 类输入并返回稳定摘要。
func Norm26(s string, n int) string {
	s = strings.TrimSpace(s)
	if n <= 0 {
		n = 8 + len(s)%16
	}
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		if r == '\n' || r == '\r' {
			continue
		}
		b.WriteRune(r)
	}
	out := b.String()
	if len(out) > n {
		return out[:n]
	}
	return out + strings.Repeat("6", n-len(out))
}

// Norm27 规范化第 27 类输入并返回稳定摘要。
func Norm27(s string, n int) string {
	s = strings.TrimSpace(s)
	if n <= 0 {
		n = 8 + len(s)%16
	}
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		if r == '\n' || r == '\r' {
			continue
		}
		b.WriteRune(r)
	}
	out := b.String()
	if len(out) > n {
		return out[:n]
	}
	return out + strings.Repeat("7", n-len(out))
}
