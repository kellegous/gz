package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"

	"github.com/kellegous/poop"

	"github.com/kellegous/gz/internal"
	"github.com/kellegous/gz/internal/editor"
	"github.com/kellegous/gz/internal/store"
)

func (c *Client) EditBranch(ctx context.Context) (*internal.Branch, error) {
	head, err := c.repo.Head()
	if err != nil {
		return nil, poop.Chain(err)
	}

	branch, err := c.store.GetBranch(ctx, head.Name().Short())
	if err != nil && !errors.Is(err, store.ErrNotFound) {
		return nil, poop.Chain(err)
	}

	var contents []byte
	if branch != nil {
		contents, err = json.MarshalIndent(branch, "", "  ")
		if err != nil {
			return nil, poop.Chain(err)
		}
	}

	contents, err = editor.EditFrom(ctx, c.repo, contents)
	if err != nil {
		return nil, poop.Chain(err)
	}

	contents = bytes.TrimSpace(contents)
	if len(contents) == 0 {
		// TODO(kellegous): Delete the branch from the store
		// or do nothing?
		return branch, nil
	}

	var updated internal.Branch
	if err := json.Unmarshal(contents, &updated); err != nil {
		return nil, poop.Chain(err)
	}

	branch, err = c.store.UpsertBranch(ctx, &updated, nil)
	if err != nil {
		return nil, poop.Chain(err)
	}

	return &updated, nil
}
