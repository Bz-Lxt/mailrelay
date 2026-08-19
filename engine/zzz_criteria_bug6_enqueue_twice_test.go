package engine_test

// 隐藏判据（外部测试包）：只断言 engine 的公开 API 行为，编译器阻止访问未导出符号。
//
// 完整运行命令（工作目录为仓库根）：
//   go test ./... -count=20 -race -run TestEnqueueTwiceSameSession

import (
	"context"
	"testing"

	"github.com/Bz-Lxt/mailrelay/config"
	"github.com/Bz-Lxt/mailrelay/engine"
	"github.com/Bz-Lxt/mailrelay/types"
)

// TestEnqueueTwiceSameSession 在同一个 Relay 会话里连续投递两封信封。
// 正常实现下两次 Enqueue 都应成功，且收件箱里恰好留下两封；一旦底层把
// 在用的句柄提前释放掉，第二次 Enqueue 就会带着错误回来，收件箱也只会
// 留下第一封。本判据只看公开行为，不关心内部用什么字段实现。
func TestEnqueueTwiceSameSession(t *testing.T) {
	r, err := engine.Open(config.Config{Dir: t.TempDir(), Addr: ":0"})
	if err != nil {
		t.Fatalf("engine.Open: %v", err)
	}
	defer r.Close()

	ctx := context.Background()
	first := types.Envelope{From: "a@local", To: "b@local", Subject: "first", Body: "x"}
	if _, err := r.Enqueue(ctx, first); err != nil {
		t.Fatalf("first enqueue: %v", err)
	}
	second := types.Envelope{From: "a@local", To: "b@local", Subject: "second", Body: "y"}
	if _, err := r.Enqueue(ctx, second); err != nil {
		t.Fatalf("second enqueue: %v", err)
	}

	inbox, err := r.List(ctx, types.BoxInbound)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(inbox) != 2 {
		t.Fatalf("inbound count = %d, want 2", len(inbox))
	}
}
