package idgen

import (
	"fmt"

	"github.com/Bz-Lxt/mailrelay/clock"
	"github.com/Bz-Lxt/mailrelay/digest"
)

var seq uint64

func Next(prefix string, clk clock.Clock) string {
	if clk == nil {
		clk = clock.Beijing{}
	}
	n := seq
	raw := fmt.Sprintf("%s|%s|%d", prefix, clock.Format(clk.Now()), n)
	id := digest.SumString(raw)
	seq = n + 1
	return id
}

func Short(id string) string { return digest.Short(id, 16) }
