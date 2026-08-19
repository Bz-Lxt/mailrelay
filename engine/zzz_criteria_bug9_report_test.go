package engine_test

import (
	"context"
	"testing"

	"github.com/Bz-Lxt/mailrelay/config"
	"github.com/Bz-Lxt/mailrelay/engine"
	"github.com/Bz-Lxt/mailrelay/types"
)

// Report must return the per-mailbox counts after a message is queued: the
// inbound mailbox holds the freshly queued envelope, so its count is non-zero.
func TestReportCountsQueuedEnvelope(t *testing.T) {
	r, err := engine.Open(config.Config{Dir: t.TempDir(), Addr: ":0"})
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer r.Close()

	if _, err := r.Enqueue(context.Background(), types.Envelope{From: "a@local", To: "b@local", Subject: "hi", Body: "x"}); err != nil {
		t.Fatalf("enqueue: %v", err)
	}

	defer func() {
		if rec := recover(); rec != nil {
			t.Fatalf("report panicked: %v", rec)
		}
	}()
	got, err := r.Report(context.Background())
	if err != nil {
		t.Fatalf("report: %v", err)
	}
	if got[types.BoxInbound] != 1 {
		t.Fatalf("report inbound=%d, want 1 (full=%v)", got[types.BoxInbound], got)
	}
}
