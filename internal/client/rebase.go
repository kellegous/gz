package client

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-git/go-git/v6/plumbing"
	"github.com/kellegous/poop"

	"github.com/kellegous/gz"
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
	head, err := c.repo.Head()
	if err != nil {
		return poop.Chain(err)
	}

	branch, err := c.store.GetBranch(ctx, head.Name().Short())
	if err != nil && !errors.Is(err, store.ErrNotFound) {
		return poop.Chain(err)
	}

	if branch == nil {
		if _, err := c.rebaseRoot(ctx, head.Name().Short(), opts); err != nil {
			return poop.Chain(err)
		}
	} else {
		if _, err := c.rebaseChild(ctx, branch, opts); err != nil {
			return poop.Chain(err)
		}
	}

	return nil
}

func (c *Client) rebaseChild(
	ctx context.Context,
	branch *gz.Branch,
	opts *RebaseOptions,
) (*plumbing.Reference, error) {
	child, err := c.store.GetBranch(ctx, branch.Parent)
	if err != nil && !errors.Is(err, store.ErrNotFound) {
		return nil, poop.Chain(err)
	}

	if child == nil {
		if _, err := c.rebaseRoot(ctx, branch.Parent, opts); err != nil {
			return nil, poop.Chain(err)
		}
	} else {
		if _, err := c.rebaseChild(ctx, child, opts); err != nil {
			return nil, poop.Chain(err)
		}
	}

	// TODO(kellegous): parent has been rebase, we may or may not need
	// to rebase the child. To determine this, we need to look at the
	// commit after the local commits to see if that is the HEAD of the
	// parent branch. If it is, we're all good. If it is not, we need to
	// rebase onto the child branch.

	return nil, nil
}

func (c *Client) rebaseRoot(
	ctx context.Context,
	name string,
	opts *RebaseOptions,
) (*plumbing.Reference, error) {
	if opts.Root == RootUpdateFetchAndRebase {
		if err := c.gitCommand(ctx, "fetch", "origin", name).Run(); err != nil {
			return nil, poop.Chain(err)
		}

		if err := c.gitCommand(ctx, "rebase", "origin/"+name).Run(); err != nil {
			return nil, poop.Chain(err)
		}
	}

	ref, err := c.repo.Reference(plumbing.NewBranchReferenceName(name), true)
	if err != nil {
		return nil, poop.Chain(err)
	}

	return ref, nil
}
