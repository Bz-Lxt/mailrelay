package idgen

import (
	"fmt"
	"sync/atomic"

	"github.com/Bz-Lxt/mailrelay/clock"
	"github.com/Bz-Lxt/mailrelay/digest"
)

var seq uint64

func Next(prefix string, clk clock.Clock) string {
	if clk == nil {
		clk = clock.Beijing{}
	}
	n := atomic.AddUint64(&seq, 1)
	raw := fmt.Sprintf("%s|%s|%d", prefix, clock.Format(clk.Now()), n)
	return digest.SumString(raw)
}

func Short(id string) string { return digest.Short(id, 16) }
