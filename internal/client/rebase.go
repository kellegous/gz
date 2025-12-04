package client

import (
	"context"
	"errors"
	"fmt"

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
		return poop.Chain(c.rebaseRoot(ctx, head.Name().Short()))
	}

	return poop.Chain(c.rebaseChild(ctx, branch))
}

func (c *Client) rebaseChild(
	ctx context.Context,
	branch *gz.Branch,
) error {
	return nil
}

func (c *Client) rebaseRoot(
	ctx context.Context,
	name string,
) error {
	return nil
}
