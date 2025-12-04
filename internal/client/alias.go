package client

import (
	"context"

	"github.com/kellegous/poop"
)

func (c *Client) Alias(
	ctx context.Context,
	name string,
	aliases []string,
) error {
	branch, err := c.store.GetBranch(ctx, name)
	if err != nil {
		return poop.Chain(err)
	}

	return poop.Chain(c.store.AliasBranch(ctx, branch.Name, aliases))
}
