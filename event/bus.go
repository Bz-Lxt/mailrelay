package event

import "sync"

type Event struct {
	Kind string
	ID   string
	Note string
}

type sub struct {
	ch   chan Event
	done <-chan struct{}
}

type Bus struct {
	mu   sync.Mutex
	cap  int
	buf  []Event
	head int
	subs []*sub
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
	} else {
		b.buf[b.head] = ev
		b.head = (b.head + 1) % b.cap
	}
	alive := make([]*sub, 0, len(b.subs))
	for _, s := range b.subs {
		select {
		case <-s.done:
			close(s.ch)
			continue
		default:
		}
		select {
		case s.ch <- ev:
		default:
		}
		alive = append(alive, s)
	}
	b.subs = alive
}

// Subscribe registers a subscriber whose delivery stops once done closes.
func (b *Bus) Subscribe(done <-chan struct{}) <-chan Event {
	b.mu.Lock()
	defer b.mu.Unlock()
	s := &sub{ch: make(chan Event, 16), done: done}
	b.subs = append(b.subs, s)
	return s.ch
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
