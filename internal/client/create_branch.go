package client

import (
	"context"

	"github.com/go-git/go-git/v6/plumbing"
	"github.com/kellegous/gz"
	"github.com/kellegous/poop"
)

func (c *Client) CreateBranch(
	ctx context.Context,
	name string,
	from string,
	aliases []string,
) (*gz.Branch, error) {
	var err error
	var ref *plumbing.Reference
	if from == "" {
		ref, err = c.repo.Head()
		if err != nil {
			return nil, poop.Chain(err)
		}
	} else {
		ref, err = c.repo.Reference(plumbing.NewBranchReferenceName(from), true)
		if err != nil {
			return nil, poop.Chain(err)
		}
	}

	args := []string{"checkout", "-b", name}
	if from != "" {
		args = append(args, from)
	}

	if err := c.gitCommand(ctx, args...).Run(); err != nil {
		return nil, poop.Chain(err)
	}

	branch, err := c.store.UpsertBranch(
		ctx,
		&gz.Branch{
			Name:   name,
			Parent: ref.Name().Short(),
		},
		aliases,
	)
	if err != nil {
		return nil, poop.Chain(err)
	}

	return branch, nil
}
