package store

import (
	"context"
	"database/sql"
	"fmt"
	"iter"
	"time"

	"github.com/kellegous/glue/fn"
	"github.com/kellegous/poop"
	"google.golang.org/protobuf/types/known/timestamppb"
	_ "modernc.org/sqlite"

	"github.com/kellegous/gz"
)

var ErrNotFound = sql.ErrNoRows

type Store struct {
	db    *sql.DB
	clock func() time.Time
}

func (s *Store) Close() error {
	return poop.Chain(s.db.Close())
}

func (s *Store) WithTx(
	ctx context.Context,
	fn func(ctx context.Context, tx *Tx) error,
) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return poop.Chain(err)
	}
	defer tx.Rollback()
	if err := fn(ctx, &Tx{tx: tx, store: s}); err != nil {
		return poop.Chain(err)
	}
	return tx.Commit()
}

func WithTx[T any](
	ctx context.Context,
	store *Store,
	fn func(ctx context.Context, tx *Tx) (T, error),
) (T, error) {
	var t T
	if err := store.WithTx(ctx, func(ctx context.Context, tx *Tx) error {
		var err error
		t, err = fn(ctx, tx)
		return err
	}); err != nil {
		return t, poop.Chain(err)
	}
	return t, nil
}

func (s *Store) BeginTx(
	ctx context.Context,
	opts *sql.TxOptions,
) (*Tx, error) {
	tx, err := s.db.BeginTx(ctx, opts)
	if err != nil {
		return nil, poop.Chain(err)
	}
	return &Tx{tx: tx, store: s}, nil
}

func (s *Store) GetBranch(ctx context.Context, name string) (*gz.Branch, error) {
	return getBranch(ctx, s.db, name)
}

func (s *Store) ListBranches(ctx context.Context) iter.Seq2[*gz.Branch, error] {
	return getBranches(ctx, s.db)
}

type Option func(*Store)

func WithClock(clock func() time.Time) Option {
	return func(s *Store) {
		s.clock = clock
	}
}

func Open(
	ctx context.Context,
	path string,
	opts ...Option,
) (*Store, error) {
	db, err := sql.Open("sqlite", fmt.Sprintf("file:%s?_journal_mode=WAL&_busy_timeout=5000", path))
	if err != nil {
		return nil, poop.Chain(err)
	}

	s, err := newStore(ctx, db, opts...)
	if err != nil {
		fn.WithCare(db.Close, &err)
		return nil, poop.Chain(err)
	}

	return s, nil

}

func InMemory(ctx context.Context, opts ...Option) (*Store, error) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		return nil, poop.Chain(err)
	}

	s, err := newStore(ctx, db, opts...)
	if err != nil {
		fn.WithCare(db.Close, &err)
		return nil, poop.Chain(err)
	}

	return s, nil
}

func newStore(ctx context.Context, db *sql.DB, opts ...Option) (*Store, error) {
	db.SetMaxOpenConns(1)

	// turn on foreign key support?
	if err := ensureSchema(ctx, db); err != nil {
		return nil, poop.Chain(err)
	}

	s := &Store{db: db, clock: time.Now}
	for _, opt := range opts {
		opt(s)
	}

	return s, nil
}

func scanBranch(row scanner) (*gz.Branch, error) {
	var (
		name           string
		description    string
		parentRef      string
		parentSha      []byte
		createdAt      int64
		updatedAt      int64
		lastAccessedAt int64
	)
	if err := row.Scan(
		&name,
		&description,
		&parentRef,
		&parentSha,
		&createdAt,
		&updatedAt,
		&lastAccessedAt,
	); err != nil {
		return nil, poop.Chain(err)
	}

	return &gz.Branch{
		Name:        name,
		Description: description,
		Parent: &gz.Parent{
			Ref: parentRef,
			Sha: parentSha,
		},
		CreatedAt:      timestamppb.New(int64ToTime(createdAt)),
		UpdatedAt:      timestamppb.New(int64ToTime(updatedAt)),
		LastAccessedAt: timestamppb.New(int64ToTime(lastAccessedAt)),
	}, nil
}

func int64ToTime(v int64) time.Time {
	if v == 0 {
		return time.Time{}
	}
	return time.Unix(0, v)
}

func timeToInt64(t time.Time) int64 {
	if t.IsZero() {
		return 0
	}
	return t.UnixNano()
}
