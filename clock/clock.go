package clock

import (
	"sync"
	"time"
)

const Offset = 8 * time.Hour

type Clock interface{ Now() time.Time }

type Beijing struct{}

func (Beijing) Now() time.Time {
	return time.Now().In(time.FixedZone("CST", int(Offset.Seconds()))).Truncate(time.Second)
}

type Fixed struct{ T time.Time }

func (f Fixed) Now() time.Time { return f.T }

type Step struct {
	mu   sync.Mutex
	next time.Time
	step time.Duration
}

func NewStep(start time.Time, step time.Duration) *Step {
	if step <= 0 {
		step = time.Second
	}
	return &Step{next: start, step: step}
}

func (s *Step) Now() time.Time {
	s.mu.Lock()
	defer s.mu.Unlock()
	cur := s.next
	s.next = s.next.Add(s.step)
	return cur
}
