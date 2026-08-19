package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Bz-Lxt/mailrelay/config"
	"github.com/Bz-Lxt/mailrelay/engine"
	"github.com/Bz-Lxt/mailrelay/types"
)

func runProbe(name string) int {
	switch name {
	case "report":
		return runReportProbe()
	default:
		fmt.Println("[EXPECT] known probe name")
		fmt.Printf("[ACTUAL] unknown probe: %s\n", name)
		return 2
	}
}

func runReportProbe() int {
	dir := filepath.Join(os.TempDir(), "mailrelay-report-probe")
	_ = os.RemoveAll(dir)
	r, err := engine.Open(config.Config{Dir: dir, Addr: ":0"})
	if err != nil {
		fmt.Println("[EXPECT] engine opens for the probe")
		fmt.Printf("[ACTUAL] open failed: %v\n", err)
		return 2
	}
	defer r.Close()

	if _, err := r.Enqueue(context.Background(), types.Envelope{From: "a@local", To: "b@local", Subject: "probe", Body: "x"}); err != nil {
		fmt.Println("[EXPECT] enqueue succeeds so a mailbox count moves off zero")
		fmt.Printf("[ACTUAL] enqueue failed: %v\n", err)
		return 2
	}

	var panicked bool
	var got map[string]int
	func() {
		defer func() {
			if rec := recover(); rec != nil {
				panicked = true
				fmt.Println("[EXPECT] report returns the per-mailbox counts")
				fmt.Printf("[ACTUAL] report panicked: %v\n", rec)
			}
		}()
		got, err = r.Report(context.Background())
	}()
	if panicked {
		return 1
	}
	if err != nil {
		fmt.Println("[EXPECT] report returns the per-mailbox counts")
		fmt.Printf("[ACTUAL] report failed: %v\n", err)
		return 2
	}
	fmt.Println("[EXPECT] report returns the per-mailbox counts")
	fmt.Printf("[ACTUAL] report inbound=%d\n", got[types.BoxInbound])
	return 0
}
