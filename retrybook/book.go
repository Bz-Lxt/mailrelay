package retrybook

type Attempt struct {
	ID    string
	Tries int
	Last  string
}

type Book struct{ items map[string]Attempt }

func New() *Book { return &Book{items: map[string]Attempt{}} }

func (b *Book) Note(id, at string) Attempt {
	cur := b.items[id]
	cur.ID = id
	cur.Tries++
	cur.Last = at
	b.items[id] = cur
	return cur
}

func (b *Book) Get(id string) (Attempt, bool) {
	v, ok := b.items[id]
	return v, ok
}

func (b *Book) Forget(id string) { delete(b.items, id) }

func (b *Book) Over(id string, limit int) bool {
	v, ok := b.items[id]
	if !ok {
		return false
	}
	return v.Tries >= limit
}

func (b *Book) Size() int { return len(b.items) }
