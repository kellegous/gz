package client

import (
	"context"
	"errors"

	"github.com/kellegous/gz/internal/store"
	"github.com/kellegous/poop"
)

func (c *Client) Checkout(ctx context.Context, name string) error {
	// look up the branch by name or alias
	branch, err := c.store.GetBranch(ctx, name)
	if err != nil && !errors.Is(err, store.ErrNotFound) {
		return poop.Chain(err)
	}

	// if found, checkout the branch by name
	if branch != nil {
		return poop.Chain(c.gitCommand(ctx, "checkout", branch.Name).Run())
	}

	return poop.Chain(c.gitCommand(ctx, "checkout", name).Run())
}
