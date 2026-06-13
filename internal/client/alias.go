package client

import (
	"context"

	"github.com/kellegous/gz/internal/store"
	"github.com/kellegous/poop"
)

func (c *Client) Alias(
	ctx context.Context,
	name string,
	aliases []string,
) error {
	return c.store.WithTx(ctx, func(ctx context.Context, tx *store.Tx) error {
		branch, err := tx.GetBranch(ctx, name)
		if err != nil {
			return poop.Chain(err)
		}

		for _, alias := range aliases {
			if err := tx.AliasBranch(ctx, branch.Name, alias); err != nil {
				return poop.Chain(err)
			}
		}

		return nil
	})
}
