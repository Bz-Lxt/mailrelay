package engine

import (
	"context"

	"github.com/Bz-Lxt/mailrelay/catalog"
	"github.com/Bz-Lxt/mailrelay/ledger"
)

// Report tallies how many envelopes sit in each known mailbox and returns the
// counts keyed by mailbox name. It is a read-only snapshot of queue state.
func (r *Relay) Report(ctx context.Context) (map[string]int, error) {
	pairs := make([]ledger.Pair, 0, len(catalog.Boxes))
	for _, box := range catalog.Boxes {
		n, err := r.db.CountBox(ctx, box)
		if err != nil {
			return nil, err
		}
		pairs = append(pairs, ledger.Pair{Key: box, Val: n})
	}
	result := ledger.CountBy(pairs)
	if result == nil {
		result = map[string]int{}
	}
	return result, nil
}
