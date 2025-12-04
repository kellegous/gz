package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"iter"

	"github.com/kellegous/poop"
	"google.golang.org/protobuf/proto"
	_ "modernc.org/sqlite"

	"github.com/kellegous/gz"
)

var ErrNotFound = errors.New("not found")

type Store struct {
	db *sql.DB
}

func (s *Store) Close() error {
	return poop.Chain(s.db.Close())
}

func (s *Store) UpsertBranch(
	ctx context.Context,
	branch *gz.Branch,
	aliases []string,
) (*gz.Branch, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, poop.Chain(err)
	}
	defer tx.Rollback()

	branch, err = upsertBranch(ctx, tx, branch)
	if err != nil {
		return nil, poop.Chain(err)
	}

	for _, alias := range aliases {
		if err := aliasBranch(ctx, tx, branch.Name, alias); err != nil {
			return nil, poop.Chain(err)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, poop.Chain(err)
	}

	return branch, nil
}

func upsertBranch(ctx context.Context, tx dbOrTx, branch *gz.Branch) (*gz.Branch, error) {
	data, err := proto.Marshal(branch)
	if err != nil {
		return nil, poop.Chain(err)
	}

	branch, err = scanBranch(tx.QueryRowContext(
		ctx,
		`INSERT INTO branches (name, data)
		VALUES (:name, :data)
		ON CONFLICT(name) DO UPDATE SET data = :data
		RETURNING name, data`,
		sql.Named("name", branch.Name),
		sql.Named("data", data),
	))
	if err != nil {
		return nil, poop.Chain(err)
	}

	return branch, nil
}

func (s *Store) AliasBranch(ctx context.Context, name, alias string) error {
	return poop.Chain(aliasBranch(ctx, s.db, name, alias))
}

func aliasBranch(ctx context.Context, tx dbOrTx, name, alias string) error {
	_, err := tx.ExecContext(ctx, `
		INSERT INTO aliases (name, alias)
		VALUES (:name, :alias)
		ON CONFLICT (alias) DO UPDATE SET name = :name
	`, sql.Named("name", name), sql.Named("alias", alias))
	return poop.Chain(err)
}

func updateBranch(
	ctx context.Context,
	tx dbOrTx,
	branch *gz.Branch,
) (*gz.Branch, error) {
	data, err := proto.Marshal(branch)
	if err != nil {
		return nil, poop.Chain(err)
	}

	branch, err = scanBranch(tx.QueryRowContext(
		ctx,
		`UPDATE branches SET data = :data WHERE name = :name RETURNING name, data`,
		sql.Named("name", branch.Name),
		sql.Named("data", data),
	))
	if err != nil {
		return nil, poop.Chain(err)
	}

	return branch, nil
}

func (s *Store) UpdateBranch(
	ctx context.Context,
	branch *gz.Branch,
) (*gz.Branch, error) {
	return updateBranch(ctx, s.db, branch)
}

func getBranch(ctx context.Context, tx dbOrTx, name string) (*gz.Branch, error) {
	return scanBranch(tx.QueryRowContext(
		ctx,
		`SELECT name, data
		 FROM branches
		 WHERE
		 	name = :name
		   	OR
			name IN (SELECT name FROM aliases WHERE alias = :name)
		`,
		sql.Named("name", name),
	))
}

func (s *Store) GetBranch(ctx context.Context, name string) (*gz.Branch, error) {
	return getBranch(ctx, s.db, name)
}

func deleteBranch(ctx context.Context, tx dbOrTx, name string) (*gz.Branch, error) {
	return scanBranch(tx.QueryRowContext(
		ctx,
		`DELETE FROM branches WHERE name = :name RETURNING name, data`,
		sql.Named("name", name),
	))
}

func (s *Store) DeleteBranch(ctx context.Context, name string) (*gz.Branch, error) {
	return deleteBranch(ctx, s.db, name)
}

func (s *Store) ListBranches(ctx context.Context) iter.Seq2[*gz.Branch, error] {
	return func(yield func(*gz.Branch, error) bool) {
		rows, err := s.db.QueryContext(
			ctx,
			`SELECT name, data FROM branches ORDER BY name ASC`,
		)
		if err != nil {
			yield(nil, poop.Chain(err))
			return
		}
		defer rows.Close()

		for rows.Next() {
			branch, err := scanBranch(rows)
			if err != nil {
				yield(nil, poop.Chain(err))
				return
			}
			if !yield(branch, nil) {
				return
			}
		}

		if err := rows.Err(); err != nil {
			yield(nil, poop.Chain(err))
			return
		}
	}
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

func scanBranch(row scanner) (*gz.Branch, error) {
	var name string
	var data []byte
	if err := row.Scan(&name, &data); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, poop.Chain(err)
	}

	var branch gz.Branch
	if err := proto.Unmarshal(data, &branch); err != nil {
		return nil, poop.Chain(err)
	}

	return &branch, nil
}
