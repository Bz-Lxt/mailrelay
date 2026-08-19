package engine_test

// 校准判据：投递信封后调用压缩，必须在有积压记录的情况下按时返回。
// 运行命令：go test ./engine/ -count=20 -race -run TestCompactReturnsUnderEnqueue -v

import (
	"context"
	"testing"
	"time"

	"github.com/Bz-Lxt/mailrelay/config"
	"github.com/Bz-Lxt/mailrelay/engine"
	"github.com/Bz-Lxt/mailrelay/types"
)

func TestCompactReturnsUnderEnqueue(t *testing.T) {
	dir := t.TempDir()
	r, err := engine.Open(config.Config{Dir: dir, Addr: ":0"})
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()

	// 投递一封正常信封，让 WAL 里有一条积压记录，压缩路径才会去校验它。
	if _, err := r.Enqueue(context.Background(), types.Envelope{
		From: "a@local", To: "b@local", Subject: "seed", Body: "y",
	}); err != nil {
		t.Fatal(err)
	}

	done := make(chan struct{})
	go func() {
		defer close(done)
		_, _ = r.Compact(context.Background())
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatalf("Compact did not return within %s", 2*time.Second)
	}
}
