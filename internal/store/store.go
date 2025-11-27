package store

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/kellegous/poop"
	_ "modernc.org/sqlite"
)

type Store struct {
	db *sql.DB
}

func (s *Store) Close() error {
	return poop.Chain(s.db.Close())
}

func Open(ctx context.Context, path string) (*Store, error) {
	db, err := sql.Open("sqlite", fmt.Sprintf("file:%s?_journal_mode=WAL&_busy_timeout=5000", path))
	if err != nil {
		return nil, poop.Chain(err)
	}

	db.SetMaxOpenConns(1)

	if err := ensureSchema(ctx, db); err != nil {
		return nil, poop.Chain(err)
	}

	return &Store{
		db: db,
	}, nil
}
