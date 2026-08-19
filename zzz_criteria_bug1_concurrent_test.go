package idgen_test

import (
	"sync"
	"testing"
	"time"

	"github.com/Bz-Lxt/mailrelay/clock"
	"github.com/Bz-Lxt/mailrelay/idgen"
)

// gatingClock holds every caller of Now until all of them have arrived.
// Because idgen reads the shared sequence before calling Now and only
// advances it afterwards, this forces every reader to observe the same
// value before any writer runs.
type gatingClock struct {
	arrived sync.WaitGroup
	stamp   time.Time
}

func newGatingClock(n int) *gatingClock {
	g := &gatingClock{stamp: time.Date(2026, 1, 1, 0, 0, 0, 0, clock.Zone())}
	g.arrived.Add(n)
	return g
}

func (g *gatingClock) Now() time.Time {
	g.arrived.Done()
	g.arrived.Wait()
	return g.stamp
}

// TestNextConcurrentIDsDistinct enqueues n generators at once and requires
// every produced id to be distinct. With correct synchronization the shared
// sequence advances before each id is built, so all ids differ.
func TestNextConcurrentIDsDistinct(t *testing.T) {
	const n = 8
	clk := newGatingClock(n)
	ids := make([]string, n)
	var wg sync.WaitGroup
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func(i int) {
			defer wg.Done()
			ids[i] = idgen.Next("mail", clk)
		}(i)
	}
	wg.Wait()

	seen := make(map[string]struct{}, n)
	for _, id := range ids {
		seen[id] = struct{}{}
	}
	if len(seen) != n {
		t.Fatalf("expected %d distinct ids, got %d: %v", n, len(seen), ids)
	}
}
