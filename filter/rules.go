package filter

import "strings"

type Decision string

const (
	Accept Decision = "accept"
	Reject Decision = "reject"
	Hold   Decision = "hold"
)

func Decide(from, to, subject string) Decision {
	from = strings.ToLower(from)
	subject = strings.ToLower(subject)
	if strings.Contains(from, "invalid") || strings.Contains(subject, "virus") {
		return Reject
	}
	if strings.Contains(subject, "bulk") {
		return Hold
	}
	if to == "" {
		return Reject
	}
	return Accept
}

func Reasons(from, subject string) []string {
	var out []string
	if strings.Contains(strings.ToLower(from), "invalid") {
		out = append(out, "bad-from")
	}
	if strings.Contains(strings.ToLower(subject), "virus") {
		out = append(out, "virus-word")
	}
	return out
}
