package client

import (
	"context"

	"github.com/kellegous/poop"
)

func (c *Client) Reset(ctx context.Context) error {
	head, err := c.repo.Head()
	if err != nil {
		return poop.Chain(err)
	}

	branch, err := c.store.GetBranch(ctx, head.Name().Short())
	if err != nil {
		return poop.Chain(err)
	}

	return poop.Chain(c.gitCommand(
		ctx,
		"reset",
		"--hard",
		branch.Parent,
	).Run())
}
