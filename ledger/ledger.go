package ledger

type Entry struct {
	ID   string
	Kind string
	Note string
}

type Book struct{ items []Entry }

func New() *Book { return &Book{} }

func (b *Book) Add(e Entry) { b.items = append(b.items, clone(e)) }

func (b *Book) All() []Entry {
	out := make([]Entry, len(b.items))
	copy(out, b.items)
	return out
}

func clone(e Entry) Entry { return e }
