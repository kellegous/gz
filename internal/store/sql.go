package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"iter"
	"strings"

	"github.com/kellegous/gz"
	"github.com/kellegous/poop"
)

var branchColumns = []string{
	"name",
	"description",
	"parent_ref",
	"parent_sha",
	"created_at",
	"updated_at",
	"last_accessed_at",
}

var upsertQuery = func() string {
	columns := strings.Join(branchColumns, ", ")
	return fmt.Sprintf(`
	INSERT INTO branches (
		%s
		) VALUES (
			:name,
			:description,
			:parent_ref,
			:parent_sha,
			:created_at,
			:updated_at,
			:last_accessed_at
		) ON CONFLICT(name) DO UPDATE SET
			description = :description,
			parent_ref = :parent_ref,
			parent_sha = :parent_sha,
			updated_at = :updated_at,
			last_accessed_at = :last_accessed_at
		WHERE name = :name
		RETURNING
		%s
	`, columns, columns)
}()

var getBranchQuery = func() string {
	columns := strings.Join(branchColumns, ", ")
	return fmt.Sprintf(`
		SELECT %s FROM branches
		WHERE name = :name
		   	OR
			name IN (SELECT name FROM aliases WHERE alias = :name)
	`, columns)
}()

var deleteBranchQuery = func() string {
	columns := strings.Join(branchColumns, ", ")
	return fmt.Sprintf(`
		DELETE FROM branches WHERE name = :name RETURNING %s
	`, columns)
}()

var listBranchesQuery = func() string {
	columns := strings.Join(branchColumns, ", ")
	return fmt.Sprintf(`
		SELECT %s FROM branches ORDER BY name ASC
	`, columns)
}()

func upsertBranch(
	ctx context.Context,
	tx dbOrTx,
	branch *gz.Branch,
) (*gz.Branch, error) {
	if branch.GetParent() == nil {
		return nil, poop.Chain(errors.New("parent is required"))
	}

	return scanBranch(tx.QueryRowContext(
		ctx,
		upsertQuery,
		sql.Named("name", branch.Name),
		sql.Named("description", branch.Description),
		sql.Named("parent_ref", branch.Parent.Ref),
		sql.Named("parent_sha", branch.Parent.Sha),
		sql.Named("created_at", timeToInt64(branch.CreatedAt.AsTime())),
		sql.Named("updated_at", timeToInt64(branch.UpdatedAt.AsTime())),
		sql.Named("last_accessed_at", timeToInt64(branch.LastAccessedAt.AsTime())),
	))
}

func getBranch(
	ctx context.Context,
	tx dbOrTx,
	name string,
) (*gz.Branch, error) {
	return scanBranch(tx.QueryRowContext(
		ctx,
		getBranchQuery,
		sql.Named("name", name),
	))
}

func getBranches(
	ctx context.Context,
	tx dbOrTx,
) iter.Seq2[*gz.Branch, error] {
	return func(yield func(*gz.Branch, error) bool) {
		rows, err := tx.QueryContext(
			ctx,
			listBranchesQuery,
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

func aliasBranch(
	ctx context.Context,
	tx dbOrTx,
	name,
	alias string,
) error {
	_, err := tx.ExecContext(ctx, `
		INSERT INTO aliases (name, alias)
		VALUES (:name, :alias)
		ON CONFLICT (alias) DO UPDATE SET name = :name
	`,
		sql.Named("name", name),
		sql.Named("alias", alias))
	return poop.Chain(err)
}

func deleteBranch(
	ctx context.Context,
	tx dbOrTx,
	name string,
) (*gz.Branch, error) {
	return scanBranch(tx.QueryRowContext(
		ctx,
		deleteBranchQuery,
		sql.Named("name", name),
	))
}
