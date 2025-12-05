package client

import (
	"context"
	"errors"
	"fmt"

	"github.com/kellegous/poop"

	"github.com/kellegous/gz/internal/store"
)

type RootUpdate string

const (
	RootUpdateNothing        RootUpdate = ""
	RootUpdateFetchAndRebase RootUpdate = "fetch-and-rebase"
	RootUpdateRebase         RootUpdate = "rebase"
)

func (r *RootUpdate) Set(v string) error {
	switch v {
	case "fetch-and-rebase":
		*r = RootUpdateFetchAndRebase
		return nil
	case "rebase":
		*r = RootUpdateRebase
		return nil
	case "nothing":
		*r = RootUpdateNothing
		return nil
	}
	return fmt.Errorf("invalid root update: %s", v)
}

func (r *RootUpdate) String() string {
	switch *r {
	case RootUpdateFetchAndRebase:
		return "fetch-and-rebase"
	case RootUpdateRebase:
		return "rebase"
	case RootUpdateNothing:
		return "nothing"
	}
	panic("unreachable")
}

func (r *RootUpdate) Type() string {
	return "string"
}

type RebaseOptions struct {
	Root RootUpdate
}

func (c *Client) Rebase(ctx context.Context, opts *RebaseOptions) error {
	if _, err := c.rebase(ctx, opts); err != nil {
		return poop.Chain(err)
	}
	return nil
}

func (c *Client) rebase(
	ctx context.Context,
	opts *RebaseOptions,
) (bool, error) {
	head, err := c.repo.Head()
	if err != nil {
		return false, poop.Chain(err)
	}

	branch, err := c.store.GetBranch(ctx, head.Name().Short())
	if err != nil && !errors.Is(err, store.ErrNotFound) {
		return false, poop.Chain(err)
	}

	if branch == nil {
		return c.rebaseRoot(ctx, opts)
	}

	// checkout the parent branch
	if err := c.gitCommand(ctx, "checkout", branch.Parent).Run(); err != nil {
		return false, poop.Chain(err)
	}

	// force a rebase of the parent branch
	needsRebase, err := c.rebase(ctx, opts)
	if err != nil {
		return false, poop.Chain(err)
	}

	// return to the current branch
	if err := c.gitCommand(ctx, "checkout", branch.Name).Run(); err != nil {
		return false, poop.Chain(err)
	}

	if !needsRebase {
		return true, nil
	}

	if len(branch.Commits) == 0 {
		if err := c.gitCommand(ctx, "rebase", branch.Parent).Run(); err != nil {
			return false, poop.Chain(err)
		}
		return false, nil
	}

	if err := c.gitCommand(
		ctx,
		"rebase",
		"--onto",
		branch.Parent,
		fmt.Sprintf("HEAD~%d", len(branch.Commits)),
	).Run(); err != nil {
		return false, poop.Chain(err)
	}

	return true, nil
}

func (c *Client) rebaseRoot(
	ctx context.Context,
	opts *RebaseOptions,
) (bool, error) {
	switch opts.Root {
	case RootUpdateNothing:
		return false, nil
	case RootUpdateRebase:
		return true, nil
	case RootUpdateFetchAndRebase:
		// TODO(kellegous): fetch and rebase the current branch
		return true, nil
	}
	panic("unreachable")
}
