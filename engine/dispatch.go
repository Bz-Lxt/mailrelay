package engine

import (
	"context"

	"github.com/Bz-Lxt/mailrelay/clock"
	"github.com/Bz-Lxt/mailrelay/event"
	"github.com/Bz-Lxt/mailrelay/types"
)

func (r *Relay) Dispatch(ctx context.Context, n int) ([]types.Envelope, error) {
	if err := ctx.Err(); err != nil {
		return nil, types.ErrCanceled
	}
	if n <= 0 {
		n = 1
	}
	inbox, err := r.db.ListBox(ctx, types.BoxInbound)
	if err != nil {
		return nil, err
	}
	out := make([]types.Envelope, 0, n)
	for _, e := range inbox {
		if len(out) >= n {
			break
		}
		e.Status = "sent"
		e.Box = types.BoxSent
		e.Tries++
		e.Updated = clock.Format(r.clock.Now())
		if err := r.db.PutEnvelope(ctx, e); err != nil {
			return out, err
		}
		r.metric.Add("dispatched", 1)
		r.events.Publish(event.Event{Kind: "dispatch", ID: e.ID})
		out = append(out, e)
	}
	return out, nil
}

func (r *Relay) List(ctx context.Context, box string) ([]types.Envelope, error) {
	if box == "" {
		box = types.BoxInbound
	}
	return r.db.ListBox(ctx, box)
}

func (r *Relay) Get(ctx context.Context, id string) (types.Envelope, error) {
	return r.db.GetEnvelope(ctx, id)
}

// Mailbox returns a one-line summary (name + envelope count) for a box.
func (r *Relay) Mailbox(ctx context.Context, box string) (*types.Mailbox, error) {
	if box == "" {
		box = types.BoxInbound
	}
	n, err := r.db.CountBox(ctx, box)
	if err != nil {
		return nil, err
	}
	mb := describeBox(box)
	mb.Count = n
	return mb, nil
}

// describeBox returns the display descriptor for a mailbox.
func describeBox(box string) *types.Mailbox {
	return &types.Mailbox{Name: box}
}

func (r *Relay) Stats() map[string]int64 { return r.metric.Snapshot() }

func (r *Relay) Events() []event.Event { return r.events.List() }
