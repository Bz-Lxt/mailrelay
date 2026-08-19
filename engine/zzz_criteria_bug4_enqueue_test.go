package engine_test

// 隐藏判据（外部测试包）：只断言 engine 的公开 API 行为，编译器阻止访问未导出符号。
//
// 完整运行命令（工作目录为仓库根）：
//   go test ./... -count=20 -race -run TestEnqueueRejectsUndeliverableRecipient

import (
	"context"
	"testing"

	"github.com/Bz-Lxt/mailrelay/config"
	"github.com/Bz-Lxt/mailrelay/engine"
	"github.com/Bz-Lxt/mailrelay/types"
)

// TestEnqueueRejectsUndeliverableRecipient 要求收件人不像一个地址（没有 @）时
// Enqueue 必须返回 types.ErrInvalid 这个判别量。HTTP 层与调用方都靠这个判别量
// 把请求归为 4xx；一旦它被包成别的错误或换成别的值，对端就会看到 5xx。
func TestEnqueueRejectsUndeliverableRecipient(t *testing.T) {
	r, err := engine.Open(config.Config{Dir: t.TempDir(), Addr: ":0"})
	if err != nil {
		t.Fatalf("engine.Open: %v", err)
	}
	defer r.Close()

	// 发件人合法、收件人不是地址（缺 @）；服务应把这封信判为非法。
	if _, err := r.Enqueue(context.Background(), types.Envelope{From: "a@local", To: "garbage"}); err != types.ErrInvalid {
		t.Fatalf("Enqueue with undeliverable recipient: got %v, want types.ErrInvalid", err)
	}
}
