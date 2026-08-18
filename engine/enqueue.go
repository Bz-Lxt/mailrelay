package engine

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/Bz-Lxt/mailrelay/clock"
	"github.com/Bz-Lxt/mailrelay/event"
	"github.com/Bz-Lxt/mailrelay/idgen"
	"github.com/Bz-Lxt/mailrelay/mimeutil"
	"github.com/Bz-Lxt/mailrelay/textutil"
	"github.com/Bz-Lxt/mailrelay/types"
	"github.com/Bz-Lxt/mailrelay/wal"
)

func (r *Relay) Enqueue(ctx context.Context, env types.Envelope) (types.Envelope, error) {
	if err := ctx.Err(); err != nil {
		return env, types.ErrCanceled
	}
	env.From = textutil.Clip(env.From, 200)
	env.To = textutil.Clip(env.To, 200)
	env.Subject = textutil.Clip(env.Subject, 200)
	if !textutil.NonEmpty(env.From) || !textutil.NonEmpty(env.To) {
		return env, types.ErrInvalid
	}
	if err := mimeutil.CheckAddress(env.To); err != nil {
		return env, err
	}
	if !r.quota.Acquire() {
		return env, types.ErrQuota
	}
	now := clock.Format(r.clock.Now())
	if env.ID == "" {
		env.ID = idgen.Next("mail", r.clock)
	}
	env.Status = "queued"
	env.Box = types.BoxInbound
	env.Created = now
	env.Updated = now
	body, _ := json.Marshal(env)
	if err := r.wal.Append(wal.Record{Op: "enqueue", Body: body, Stamp: now}); err != nil {
		r.quota.Release()
		return env, err
	}
	if err := r.db.PutEnvelope(ctx, env); err != nil {
		r.quota.Release()
		return env, err
	}
	r.metric.Add("enqueued", 1)
	r.events.Publish(event.Event{Kind: "enqueue", ID: env.ID})
	return env, nil
}

func (r *Relay) Defer(ctx context.Context, id, until string) error {
	if err := ctx.Err(); err != nil {
		return types.ErrCanceled
	}
	e, err := r.db.GetEnvelope(ctx, id)
	if err != nil {
		return err
	}
	e.Status = "deferred"
	e.Box = types.BoxDeferred
	e.DeferTo = until
	e.Updated = clock.Format(r.clock.Now())
	return r.db.PutEnvelope(ctx, e)
}

func (r *Relay) Bounce(ctx context.Context, id, reason string) error {
	if err := ctx.Err(); err != nil {
		return types.ErrCanceled
	}
	e, err := r.db.GetEnvelope(ctx, id)
	if err != nil {
		return err
	}
	e.Status = "bounced"
	e.Box = types.BoxBounced
	e.Reason = strings.TrimSpace(reason)
	e.Updated = clock.Format(r.clock.Now())
	r.metric.Add("bounced", 1)
	return r.db.PutEnvelope(ctx, e)
}
