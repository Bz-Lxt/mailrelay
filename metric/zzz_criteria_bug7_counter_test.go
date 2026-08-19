package metric_test

import (
	"testing"

	"github.com/Bz-Lxt/mailrelay/metric"
)

func TestCounterAccumulatesAcrossAdds(t *testing.T) {
	c := metric.New()
	for i := 0; i < 5; i++ {
		c.Add("enqueued", 1)
	}
	got := c.Snapshot()["enqueued"]
	if got != 5 {
		t.Fatalf("enqueued counter want 5 got %d", got)
	}
}

func TestCounterAccumulatesDistinctKeys(t *testing.T) {
	c := metric.New()
	c.Add("enqueued", 1)
	c.Add("dispatched", 1)
	c.Add("enqueued", 1)
	c.Add("dispatched", 1)
	c.Add("enqueued", 1)
	snap := c.Snapshot()
	if got := snap["enqueued"]; got != 3 {
		t.Fatalf("enqueued counter want 3 got %d", got)
	}
	if got := snap["dispatched"]; got != 2 {
		t.Fatalf("dispatched counter want 2 got %d", got)
	}
}

func TestCounterSnapshotDoesNotResetState(t *testing.T) {
	c := metric.New()
	c.Add("enqueued", 1)
	c.Add("enqueued", 1)
	first := c.Snapshot()["enqueued"]
	c.Add("enqueued", 1)
	second := c.Snapshot()["enqueued"]
	if first != 2 {
		t.Fatalf("first snapshot want 2 got %d", first)
	}
	if second != 3 {
		t.Fatalf("second snapshot want 3 got %d", second)
	}
}
