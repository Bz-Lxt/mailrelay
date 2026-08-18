package metric

import "sync"

type Counters struct {
	mu sync.Mutex
	m  map[string]int64
}

func New() *Counters { return &Counters{m: map[string]int64{}} }

func (c *Counters) Add(name string, n int64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.m[name] += n
}

func (c *Counters) Get(name string) int64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.m[name]
}

func (c *Counters) Snapshot() map[string]int64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	out := make(map[string]int64, len(c.m))
	for k, v := range c.m {
		out[k] = v
	}
	return out
}
