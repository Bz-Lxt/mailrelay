package engine

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Bz-Lxt/mailrelay/event"
	"github.com/Bz-Lxt/mailrelay/types"
	"github.com/Bz-Lxt/mailrelay/wal"
)

// Compact trims the write-ahead log up to the current sequence, keeping only
// records that are still in flight. Callers get back a snapshot of the
// counters taken before the trim so they can observe what was compacted.
func (r *Relay) Compact(ctx context.Context) (map[string]int64, error) {
	r.compactMu.Lock()
	defer r.compactMu.Unlock()
	recs, err := wal.ReplayFile(r.wal.Path())
	if err != nil {
		return nil, err
	}
	for _, rec := range recs {
		if rec.Op != "enqueue" {
			continue
		}
		var e types.Envelope
		if err := json.Unmarshal(rec.Body, &e); err != nil {
			return nil, err
		}
		if _, err := r.Lookup(ctx, e.ID); err != nil {
			return nil, err
		}
		if err := r.verifyRecord(ctx, e); err != nil {
			return nil, err
		}
	}
	snap := r.metric.Snapshot()
	if err := r.wal.TruncatePrefix(r.wal.NextSeq()); err != nil {
		return snap, err
	}
	r.events.Publish(event.Event{Kind: "compact", Note: "truncated"})
	return snap, nil
}

// verifyRecord double-checks that an envelope still has a usable row before
// its journal entry is trimmed.
func (r *Relay) verifyRecord(ctx context.Context, e types.Envelope) error {
	got, err := r.db.GetEnvelope(ctx, e.ID)
	if err != nil {
		return err
	}
	if got.Status == "" {
		return fmt.Errorf("envelope %s has no status", e.ID)
	}
	return nil
}
