package engine

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/Bz-Lxt/mailrelay/config"
	"github.com/Bz-Lxt/mailrelay/types"
)

func TestEnqueueDispatchRestart(t *testing.T) {
	dir := t.TempDir()
	r, err := Open(config.Config{Dir: dir, Addr: ":0"})
	if err != nil {
		t.Fatal(err)
	}
	env, err := r.Enqueue(context.Background(), types.Envelope{From: "a@local", To: "b@local", Subject: "hi", Body: "ok"})
	if err != nil {
		t.Fatal(err)
	}
	sent, err := r.Dispatch(context.Background(), 1)
	if err != nil || len(sent) != 1 || sent[0].ID != env.ID {
		t.Fatalf("dispatch %+v %v", sent, err)
	}
	_ = r.Close()
	r2, err := Open(config.Config{Dir: dir, Addr: ":0"})
	if err != nil {
		t.Fatal(err)
	}
	defer r2.Close()
	got, err := r2.Get(context.Background(), env.ID)
	if err != nil || got.Box != types.BoxSent {
		t.Fatalf("restart %+v %v", got, err)
	}
	_ = filepath.Join(dir, "relay.sqlite")
}

func TestCanceledEnqueue(t *testing.T) {
	r, err := Open(config.Config{Dir: t.TempDir(), Addr: ":0"})
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := r.Enqueue(ctx, types.Envelope{From: "a@local", To: "b@local"}); err != types.ErrCanceled {
		t.Fatalf("want canceled got %v", err)
	}
}
