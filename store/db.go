package store

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"

	"github.com/Bz-Lxt/mailrelay/clock"

	_ "modernc.org/sqlite"
)

type DB struct {
	sql   *sql.DB
	clock clock.Clock
	path  string
}

func Open(path string, clk clock.Clock) (*DB, error) {
	if clk == nil {
		clk = clock.Beijing{}
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, err
	}
	sqlDB, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(1)
	if _, err := sqlDB.Exec(schemaSQL); err != nil {
		_ = sqlDB.Close()
		return nil, err
	}
	return &DB{sql: sqlDB, clock: clk, path: path}, nil
}

func (db *DB) Close() error {
	if db == nil || db.sql == nil {
		return nil
	}
	return db.sql.Close()
}

func (db *DB) SQL() *sql.DB { return db.sql }

func (db *DB) Clock() clock.Clock { return db.clock }

func (db *DB) SetMeta(ctx context.Context, key, value string) error {
	_, err := db.sql.ExecContext(ctx, `INSERT INTO meta(key,value) VALUES(?,?)
		ON CONFLICT(key) DO UPDATE SET value=excluded.value`, key, value)
	return err
}

func (db *DB) Meta(ctx context.Context, key string) (string, error) {
	var v string
	err := db.sql.QueryRowContext(ctx, `SELECT value FROM meta WHERE key=?`, key).Scan(&v)
	if err == sql.ErrNoRows {
		return "", nil
	}
	return v, err
}
