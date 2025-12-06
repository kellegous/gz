package client

import (
	"context"

	"github.com/kellegous/gz/internal"
	"github.com/kellegous/poop"
)

func (c *Client) Reset(ctx context.Context) (*internal.Branch, error) {
	head, err := c.repo.Head()
	if err != nil {
		return nil, poop.Chain(err)
	}

	branch, err := c.store.GetBranch(ctx, head.Name().Short())
	if err != nil {
		return nil, poop.Chain(err)
	}

	if err := c.gitCommand(
		ctx,
		"reset",
		"--hard",
		branch.Parent,
	).Run(); err != nil {
		return nil, poop.Chain(err)
	}

	branch, err = c.store.UpdateBranch(ctx, &internal.Branch{
		Name:        branch.Name,
		Parent:      branch.Parent,
		Description: branch.Description,
	})
	if err != nil {
		return nil, poop.Chain(err)
	}

	return branch, nil
}
