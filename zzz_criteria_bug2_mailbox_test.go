package engine_test

import (
	"context"
	"testing"

	"github.com/Bz-Lxt/mailrelay/config"
	"github.com/Bz-Lxt/mailrelay/engine"
	"github.com/Bz-Lxt/mailrelay/types"
)

// TestMailboxSummaryOutbound asserts that summarizing the outbound mailbox
// returns a non-nil descriptor with the right name and count. The call is
// wrapped so a runtime panic surfaces as a normal test failure rather than
// crashing the test binary.
func TestMailboxSummaryOutbound(t *testing.T) {
	r, err := engine.Open(config.Config{Dir: t.TempDir(), Addr: ":0"})
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()

	var got *types.Mailbox
	var panicked bool
	func() {
		defer func() {
			if rec := recover(); rec != nil {
				panicked = true
			}
		}()
		got, err = r.Mailbox(context.Background(), types.BoxOutbound)
	}()

	if panicked {
		t.Fatalf("Mailbox panicked for outbound box")
	}
	if err != nil {
		t.Fatalf("Mailbox returned error: %v", err)
	}
	if got == nil {
		t.Fatalf("Mailbox returned nil descriptor")
	}
	if got.Name != types.BoxOutbound {
		t.Fatalf("mailbox name = %q, want %q", got.Name, types.BoxOutbound)
	}
	if got.Count != 0 {
		t.Fatalf("mailbox count = %d, want 0", got.Count)
	}
}
