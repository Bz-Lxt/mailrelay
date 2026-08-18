package event

import "sync"

type Event struct {
	Kind string
	ID   string
	Note string
}

type Bus struct {
	mu   sync.Mutex
	cap  int
	buf  []Event
	head int
}

func New(cap int) *Bus {
	if cap < 4 {
		cap = 64
	}
	return &Bus{cap: cap, buf: make([]Event, 0, cap)}
}

func (b *Bus) Publish(ev Event) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if len(b.buf) < b.cap {
		b.buf = append(b.buf, ev)
		return
	}
	b.buf[b.head] = ev
	b.head = (b.head + 1) % b.cap
}

func (b *Bus) List() []Event {
	b.mu.Lock()
	defer b.mu.Unlock()
	if len(b.buf) < b.cap {
		out := make([]Event, len(b.buf))
		copy(out, b.buf)
		return out
	}
	out := make([]Event, 0, b.cap)
	out = append(out, b.buf[b.head:]...)
	out = append(out, b.buf[:b.head]...)
	return out
}
