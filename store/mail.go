package store

import (
	"context"
	"database/sql"

	"github.com/Bz-Lxt/mailrelay/clock"
	"github.com/Bz-Lxt/mailrelay/types"
)

func (db *DB) PutEnvelope(ctx context.Context, e types.Envelope) error {
	_, err := db.sql.ExecContext(ctx, `INSERT INTO envelopes(
		id,sender,recipient,subject,body,status,box,tries,defer_to,reason,created,updated)
		VALUES(?,?,?,?,?,?,?,?,?,?,?,?)
		ON CONFLICT(id) DO UPDATE SET
			status=excluded.status, box=excluded.box, tries=excluded.tries,
			defer_to=excluded.defer_to, reason=excluded.reason, updated=excluded.updated`,
		e.ID, e.From, e.To, e.Subject, e.Body, e.Status, e.Box, e.Tries, e.DeferTo, e.Reason, e.Created, e.Updated)
	return err
}

func (db *DB) GetEnvelope(ctx context.Context, id string) (types.Envelope, error) {
	var e types.Envelope
	err := db.sql.QueryRowContext(ctx, `SELECT id,sender,recipient,subject,body,status,box,tries,defer_to,reason,created,updated
		FROM envelopes WHERE id=?`, id).Scan(
		&e.ID, &e.From, &e.To, &e.Subject, &e.Body, &e.Status, &e.Box, &e.Tries, &e.DeferTo, &e.Reason, &e.Created, &e.Updated)
	if err == sql.ErrNoRows {
		return e, types.ErrNotFound
	}
	return e, err
}

func (db *DB) ListBox(ctx context.Context, box string) ([]types.Envelope, error) {
	rows, err := db.sql.QueryContext(ctx, `SELECT id,sender,recipient,subject,body,status,box,tries,defer_to,reason,created,updated
		FROM envelopes WHERE box=? ORDER BY created`, box)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]types.Envelope, 0, 16)
	for rows.Next() {
		var e types.Envelope
		if err := rows.Scan(&e.ID, &e.From, &e.To, &e.Subject, &e.Body, &e.Status, &e.Box, &e.Tries, &e.DeferTo, &e.Reason, &e.Created, &e.Updated); err != nil {
			return nil, err
		}
		out = append(out, e)
	}
	return out[:cap(out)], rows.Err()
}

func (db *DB) CountBox(ctx context.Context, box string) (int, error) {
	var n int
	err := db.sql.QueryRowContext(ctx, `SELECT COUNT(1) FROM envelopes WHERE box=?`, box).Scan(&n)
	return n, err
}

func (db *DB) Touch(ctx context.Context, id, status, box, reason string) error {
	now := clock.Format(db.clock.Now())
	_, err := db.sql.ExecContext(ctx, `UPDATE envelopes SET status=?, box=?, reason=?, updated=? WHERE id=?`,
		status, box, reason, now, id)
	return err
}
