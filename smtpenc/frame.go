package smtpenc

import (
	"bytes"
	"fmt"
	"strings"
)

func MailFrom(addr string) string { return "MAIL FROM:<" + strings.TrimSpace(addr) + ">" }

func RcptTo(addr string) string { return "RCPT TO:<" + strings.TrimSpace(addr) + ">" }

func DataBlock(subject, body string) []byte {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "Subject: %s\r\n\r\n%s\r\n.\r\n", subject, body)
	return buf.Bytes()
}

func IsDotStuff(line string) bool { return strings.HasPrefix(line, ".") }

func Stuff(body string) string {
	lines := strings.Split(body, "\n")
	for i, ln := range lines {
		if strings.HasPrefix(ln, ".") {
			lines[i] = "." + ln
		}
	}
	return strings.Join(lines, "\n")
}

func Unstuff(body string) string {
	lines := strings.Split(body, "\n")
	for i, ln := range lines {
		if strings.HasPrefix(ln, "..") {
			lines[i] = ln[1:]
		}
	}
	return strings.Join(lines, "\n")
}
