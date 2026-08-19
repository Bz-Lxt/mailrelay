package engine_test

// 运行命令：go test ./engine/ -count=20 -race -run TestEnqueueKeepsSubjectHeadAndFullDigest -v

import (
	"context"
	"strings"
	"testing"

	"github.com/Bz-Lxt/mailrelay/config"
	"github.com/Bz-Lxt/mailrelay/engine"
	"github.com/Bz-Lxt/mailrelay/types"
)

// 投递一封 subject 超过 200 字节的信封后，服务必须保留 subject 的开头 200 字节，
// 并且为信封分配完整长度的摘要标识，而不是只保留尾部窗口。
func TestEnqueueKeepsSubjectHeadAndFullDigest(t *testing.T) {
	r, err := engine.Open(config.Config{Dir: t.TempDir(), Addr: ":0"})
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer r.Close()

	long := "HEAD" + strings.Repeat("a", 246) // 250 bytes
	env, err := r.Enqueue(context.Background(), types.Envelope{
		From: "a@local", To: "b@local", Subject: long, Body: "x",
	})
	if err != nil {
		t.Fatalf("enqueue: %v", err)
	}

	// 开头标记必须存活：字段被截到 200 字节时应保留开头，而不是尾巴。
	wantSubject := "HEAD" + strings.Repeat("a", 196)
	if env.Subject != wantSubject {
		t.Fatalf("subject not clipped to head: got len=%d %q, want len=%d %q",
			len(env.Subject), env.Subject, len(wantSubject), wantSubject)
	}

	// 信封标识必须是完整长度的摘要，而不是只剩尾段的短串。
	if len(env.ID) != 64 {
		t.Fatalf("id not full digest: got len=%d %q, want 64", len(env.ID), env.ID)
	}

	// 持久化后再读，结论必须一致。
	got, err := r.Get(context.Background(), env.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.Subject != wantSubject || len(got.ID) != 64 {
		t.Fatalf("persisted envelope drifted: subject len=%d id len=%d",
			len(got.Subject), len(got.ID))
	}
}
