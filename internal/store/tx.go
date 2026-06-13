package store

import (
	"context"
	"database/sql"
	"iter"

	"github.com/kellegous/gz"
)

type Tx struct {
	tx    *sql.Tx
	store *Store
}

func (t *Tx) Commit() error {
	return t.tx.Commit()
}

func (t *Tx) Rollback() error {
	return t.tx.Rollback()
}

func (t *Tx) UpsertBranch(
	ctx context.Context,
	branch *gz.Branch,
) (*gz.Branch, error) {
	return upsertBranch(ctx, t.tx, branch)
}

func (t *Tx) GetBranch(
	ctx context.Context,
	name string,
) (*gz.Branch, error) {
	return getBranch(ctx, t.tx, name)
}

func (t *Tx) DeleteBranch(
	ctx context.Context,
	name string,
) (*gz.Branch, error) {
	return deleteBranch(ctx, t.tx, name)
}

func (t *Tx) ListBranches(
	ctx context.Context,
) iter.Seq2[*gz.Branch, error] {
	return getBranches(ctx, t.tx)
}

func (t *Tx) AliasBranch(
	ctx context.Context,
	name string,
	alias string,
) error {
	return aliasBranch(ctx, t.tx, name, alias)
}
