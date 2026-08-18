package quota

import "sync"

type Counter struct {
	mu    sync.Mutex
	limit int
	used  int
}

func New(limit int) *Counter {
	if limit <= 0 {
		limit = 1024
	}
	return &Counter{limit: limit}
}

func (c *Counter) Acquire() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.used >= c.limit {
		return false
	}
	c.used++
	return true
}

func (c *Counter) Release() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.used > 0 {
		c.used--
	}
}

func (c *Counter) Used() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.used
}
