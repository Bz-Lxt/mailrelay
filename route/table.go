package route

import "strings"

type Rule struct {
	Prefix string
	Box    string
}

type Table struct{ rules []Rule }

func New() *Table { return &Table{} }

func (t *Table) Add(prefix, box string) {
	t.rules = append(t.rules, Rule{Prefix: strings.ToLower(prefix), Box: box})
}

func (t *Table) Lookup(addr string) string {
	addr = strings.ToLower(strings.TrimSpace(addr))
	best := ""
	bestN := -1
	for _, r := range t.rules {
		if strings.HasSuffix(addr, r.Prefix) && len(r.Prefix) > bestN {
			best = r.Box
			bestN = len(r.Prefix)
		}
	}
	if best == "" {
		return "inbound"
	}
	return best
}

func (t *Table) Withdraw(prefix string) {
	prefix = strings.ToLower(prefix)
	out := t.rules[:0]
	for _, r := range t.rules {
		if r.Prefix != prefix {
			out = append(out, r)
		}
	}
	t.rules = out
}

func (t *Table) Dump() []Rule {
	out := make([]Rule, len(t.rules))
	copy(out, t.rules)
	return out
}
