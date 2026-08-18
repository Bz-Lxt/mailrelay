package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/Bz-Lxt/mailrelay/config"
	"github.com/Bz-Lxt/mailrelay/engine"
	"github.com/Bz-Lxt/mailrelay/types"
)

func runProbe(name string) int {
	switch name {
	case "subscribe":
		return runSubscribeProbe()
	default:
		fmt.Println("[EXPECT] known probe name")
		fmt.Printf("[ACTUAL] unknown probe: %s\n", name)
		return 2
	}
}

func runSubscribeProbe() int {
	dir := filepath.Join(os.TempDir(), "mailrelay-subscribe-probe")
	_ = os.RemoveAll(dir)
	r, err := engine.Open(config.Config{Dir: dir, Addr: ":0"})
	if err != nil {
		fmt.Println("[EXPECT] engine opens for the probe")
		fmt.Printf("[ACTUAL] open failed: %v\n", err)
		return 2
	}
	defer r.Close()

	ctx, cancel := context.WithCancel(context.Background())
	ch := r.Subscribe(ctx)
	cancel()

	if _, err := r.Enqueue(context.Background(), types.Envelope{From: "a@local", To: "b@local", Subject: "probe", Body: "x"}); err != nil {
		fmt.Println("[EXPECT] enqueue succeeds so an event is published")
		fmt.Printf("[ACTUAL] enqueue failed: %v\n", err)
		return 2
	}

	select {
	case ev, ok := <-ch:
		if ok {
			fmt.Println("[EXPECT] canceled subscription receives no event")
			fmt.Printf("[ACTUAL] event delivered after cancel (kind=%s)\n", ev.Kind)
			return 1
		}
		fmt.Println("[EXPECT] canceled subscription receives no event")
		fmt.Println("[ACTUAL] subscription closed after cancel, no event")
		return 0
	case <-time.After(time.Second):
		fmt.Println("[EXPECT] canceled subscription receives no event")
		fmt.Println("[ACTUAL] no event delivered after cancel")
		return 0
	}
}
