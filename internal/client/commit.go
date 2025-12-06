package client

import (
	"context"

	"github.com/kellegous/gz/internal"
	"github.com/kellegous/poop"
)

type CommitOptions struct {
	All     bool
	Append  bool
	Message MessageOption
}

func (c *Client) Commit(ctx context.Context, opts *CommitOptions) (*internal.Branch, error) {
	head, err := c.repo.Head()
	if err != nil {
		return nil, poop.Chain(err)
	}

	branch, err := c.store.GetBranch(ctx, head.Name().Short())
	if err != nil {
		return nil, poop.Chain(err)
	}

	args := []string{"commit"}

	amend := len(branch.Commits) > 0 && !opts.Append

	if opts.All {
		args = append(args, "-a")
	}

	if amend {
		args = append(args, "--amend")
	}

	if m := opts.Message; m.valid {
		if t := m.text; t != "" {
			args = append(args, "-m", t)
		} else if amend {
			args = append(args, "--no-edit")
		}
	}

	if err := c.gitCommand(ctx, args...).Run(); err != nil {
		return nil, poop.Chain(err)
	}

	head, err = c.repo.Head()
	if err != nil {
		return nil, poop.Chain(err)
	}

	commits := branch.Commits
	if amend {
		commits = append(commits[:len(commits)-1], head.Hash().Bytes())
	} else {
		commits = append(commits, head.Hash().Bytes())
	}

	branch, err = c.store.UpdateBranch(ctx, &internal.Branch{
		Name:        branch.Name,
		Commits:     commits,
		Parent:      branch.Parent,
		Description: branch.Description,
	})
	if err != nil {
		return nil, poop.Chain(err)
	}

	return branch, nil
}
