package smtpenc

import "strings"

func CodeOK(line string) bool {
	return strings.HasPrefix(line, "250") || strings.HasPrefix(line, "354")
}

func CodeFail(line string) bool {
	if len(line) < 3 {
		return false
	}
	return line[0] == '4' || line[0] == '5'
}

func ExtractCode(line string) string {
	if len(line) < 3 {
		return ""
	}
	return line[:3]
}
