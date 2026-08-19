package engine

import (
	"context"

	"github.com/Bz-Lxt/mailrelay/types"
)

// Lookup fetches a single envelope by id. Callers use it instead of reaching
// into the store so the relay can attach its own bookkeeping later.
func (r *Relay) Lookup(ctx context.Context, id string) (types.Envelope, error) {
	r.compactMu.Lock()
	defer r.compactMu.Unlock()
	return r.db.GetEnvelope(ctx, id)
}
