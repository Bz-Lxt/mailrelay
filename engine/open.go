package engine

import (
	"context"
	"encoding/json"
	"path/filepath"
	"sync"

	"github.com/Bz-Lxt/mailrelay/clock"
	"github.com/Bz-Lxt/mailrelay/config"
	"github.com/Bz-Lxt/mailrelay/event"
	"github.com/Bz-Lxt/mailrelay/metric"
	"github.com/Bz-Lxt/mailrelay/quota"
	"github.com/Bz-Lxt/mailrelay/store"
	"github.com/Bz-Lxt/mailrelay/types"
	"github.com/Bz-Lxt/mailrelay/wal"
)

type Relay struct {
	mu     sync.Mutex
	cfg    config.Config
	db     *store.DB
	wal    *wal.Journal
	clock  clock.Clock
	quota  *quota.Counter
	metric *metric.Counters
	events *event.Bus
}

func Open(cfg config.Config) (*Relay, error) {
	cfg, err := config.Normalize(cfg)
	if err != nil {
		return nil, err
	}
	clk := clock.Beijing{}
	db, err := store.Open(filepath.Join(cfg.Dir, "relay.sqlite"), clk)
	if err != nil {
		return nil, err
	}
	j, err := wal.Open(filepath.Join(cfg.Dir, "journal.wal"))
	if err != nil {
		_ = db.Close()
		return nil, err
	}
	r := &Relay{
		cfg: cfg, db: db, wal: j, clock: clk,
		quota: quota.New(4096), metric: metric.New(), events: event.New(128),
	}
	if err := r.replay(); err != nil {
		_ = r.Close()
		return nil, err
	}
	return r, nil
}

func (r *Relay) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	_ = r.wal.Close()
	return r.db.Close()
}

func (r *Relay) replay() error {
	recs, err := wal.ReplayFile(r.wal.Path())
	if err != nil {
		return err
	}
	ctx := context.Background()
	for _, rec := range recs {
		if rec.Op != "enqueue" {
			continue
		}
		var e types.Envelope
		if err := json.Unmarshal(rec.Body, &e); err != nil {
			return err
		}
		if _, err := r.db.GetEnvelope(ctx, e.ID); err == nil {
			continue
		}
		if err := r.db.PutEnvelope(ctx, e); err != nil {
			return err
		}
	}
	return nil
}
