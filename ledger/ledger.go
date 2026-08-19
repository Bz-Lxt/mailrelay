package ledger

type Entry struct {
	ID   string
	Kind string
	Note string
}

type Pair struct {
	Key string
	Val int
}

type Book struct{ items []Entry }

func New() *Book { return &Book{} }

func (b *Book) Add(e Entry) { b.items = append(b.items, clone(e)) }

func (b *Book) All() []Entry {
	out := make([]Entry, len(b.items))
	copy(out, b.items)
	return out
}

// CountBy folds a slice of key/count pairs into a map keyed by Key. Pairs whose
// count is zero are skipped so the result only carries mailboxes that hold
// envelopes.
func CountBy(pairs []Pair) map[string]int {
	out := make(map[string]int)
	for _, p := range pairs {
		if p.Val == 0 {
			continue
		}
		out[p.Key] = p.Val
	}
	return out
}

func clone(e Entry) Entry { return e }
