package store

import (
	"database/sql"
	"fmt"

	"github.com/kellegous/poop"
	_ "modernc.org/sqlite"
)

type Store struct {
	db *sql.DB
}

func Open(path string) (*Store, error) {
	db, err := sql.Open("sqlite", fmt.Sprintf("file:%s?_journal_mode=WAL&_busy_timeout=5000", path))
	if err != nil {
		return nil, poop.Chain(err)
	}

	db.SetMaxOpenConns(1)

	// TODO(kellegous): ensure schema

	return &Store{
		db: db,
	}, nil
}
