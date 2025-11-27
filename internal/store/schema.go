package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/kellegous/poop"
)

type dbOrTx interface {
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
}

func getSchemaVersion(ctx context.Context, db dbOrTx) (int, error) {
	var version int
	if err := db.QueryRowContext(ctx, "PRAGMA user_version").Scan(&version); errors.Is(err, sql.ErrNoRows) {
		return 0, nil
	} else if err != nil {
		return 0, poop.Chain(err)
	}
	return version, nil
}

func setSchemaVersion(ctx context.Context, db dbOrTx, version int) error {
	_, err := db.ExecContext(ctx, fmt.Sprintf("PRAGMA user_version = %d", version))
	return poop.Chain(err)
}

func ensureSchema(ctx context.Context, db *sql.DB) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return poop.Chain(err)
	}
	defer tx.Rollback()

	version, err := getSchemaVersion(ctx, tx)
	if err != nil {
		return poop.Chain(err)
	}

	for v := version; v < len(migrations); v++ {
		if err := migrations[v](ctx, tx); err != nil {
			return poop.Chain(err)
		}
		if err := setSchemaVersion(ctx, tx, v); err != nil {
			return poop.Chain(err)
		}
	}

	return poop.Chain(tx.Commit())
}

var migrations = []func(ctx context.Context, tx *sql.Tx) error{
	// Version 0: initial schema
	func(ctx context.Context, tx *sql.Tx) error {
		if _, err := tx.ExecContext(ctx, `
			CREATE TABLE IF NOT EXISTS branches (
				name TEXT PRIMARY KEY,
				data BLOB NOT NULL
			)
		`); err != nil {
			return poop.Chain(err)
		}
		return nil
	},
}
