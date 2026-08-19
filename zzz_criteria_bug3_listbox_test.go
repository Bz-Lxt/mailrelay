package store_test

// 隐藏判据（外部测试包）：只断言 store 的公开 API 行为，编译器阻止访问未导出符号。
//
// 完整运行命令（工作目录为仓库根）：
//   go test ./... -count=20 -race -run TestListBoxReturnsOnlyStoredEnvelopes

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/Bz-Lxt/mailrelay/clock"
	"github.com/Bz-Lxt/mailrelay/store"
	"github.com/Bz-Lxt/mailrelay/types"
)

// TestListBoxReturnsOnlyStoredEnvelopes 入队两个信封后读取收件箱，要求返回的切片
// 恰好只含已落盘的两个信封，不多不少，也不夹带任何空信封。
func TestListBoxReturnsOnlyStoredEnvelopes(t *testing.T) {
	db, err := store.Open(filepath.Join(t.TempDir(), "relay.sqlite"), clock.Beijing{})
	if err != nil {
		t.Fatalf("store.Open: %v", err)
	}
	defer db.Close()

	ctx := context.Background()
	want := []string{"mail-aaa", "mail-bbb"}
	for i, id := range want {
		env := types.Envelope{
			ID:      id,
			From:    "a@local",
			To:      "b@local",
			Subject: "hello",
			Body:    "body",
			Status:  "queued",
			Box:     types.BoxInbound,
			Created: "2026-01-01 00:00:00",
			Updated: "2026-01-01 00:00:00",
		}
		if err := db.PutEnvelope(ctx, env); err != nil {
			t.Fatalf("PutEnvelope[%d]: %v", i, err)
		}
	}

	got, err := db.ListBox(ctx, types.BoxInbound)
	if err != nil {
		t.Fatalf("ListBox: %v", err)
	}

	if len(got) != len(want) {
		t.Fatalf("ListBox returned %d envelopes, want %d", len(got), len(want))
	}

	seen := make(map[string]bool, len(got))
	for i, e := range got {
		if e.ID == "" {
			t.Fatalf("ListBox[%d] is an empty envelope, want only stored rows", i)
		}
		if e.Box != types.BoxInbound {
			t.Fatalf("ListBox[%d].Box = %q, want %q", i, e.Box, types.BoxInbound)
		}
		if seen[e.ID] {
			t.Fatalf("ListBox[%d] duplicates id %q", i, e.ID)
		}
		seen[e.ID] = true
	}
	for _, id := range want {
		if !seen[id] {
			t.Fatalf("ListBox missing stored envelope %q, got %v", id, seen)
		}
	}
}
