package mailrelay_test

import (
	"context"
	"testing"
	"time"

	"github.com/Bz-Lxt/mailrelay/config"
	"github.com/Bz-Lxt/mailrelay/engine"
	"github.com/Bz-Lxt/mailrelay/types"
)

func TestCanceledSubscriptionReceivesNoEvent(t *testing.T) {
	r, err := engine.Open(config.Config{Dir: t.TempDir(), Addr: ":0"})
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()

	ctx, cancel := context.WithCancel(context.Background())
	ch := r.Subscribe(ctx)
	cancel()

	if _, err := r.Enqueue(context.Background(), types.Envelope{From: "a@local", To: "b@local", Subject: "x", Body: "y"}); err != nil {
		t.Fatalf("enqueue: %v", err)
	}

	select {
	case ev, ok := <-ch:
		if ok {
			t.Fatalf("取消后的订阅仍投递事件: %+v", ev)
		}
	case <-time.After(time.Second):
	}
}
