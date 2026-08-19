package idgen

import (
	"fmt"
	"sync/atomic"

	"github.com/Bz-Lxt/mailrelay/clock"
	"github.com/Bz-Lxt/mailrelay/digest"
)

var seq uint64

// Next returns a digest-based id for an envelope. The sequence counter must be
// advanced atomically: concurrent callers must each observe a distinct, strictly
// increasing n, otherwise two submissions in the same clock second can derive
// the same raw string (same prefix, same second-resolution timestamp, same n)
// and therefore the same id, which the store would collapse via
// ON CONFLICT(id) DO UPDATE — silently dropping envelopes.
func Next(prefix string, clk clock.Clock) string {
	if clk == nil {
		clk = clock.Beijing{}
	}
	n := atomic.AddUint64(&seq, 1) - 1
	raw := fmt.Sprintf("%s|%s|%d", prefix, clock.Format(clk.Now()), n)
	id := digest.SumString(raw)
	return id
}

func Short(id string) string { return digest.Short(id, 16) }
