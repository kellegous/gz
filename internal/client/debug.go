package client

import (
	"context"

	"github.com/kellegous/gz"
	"github.com/kellegous/poop"
)

func (c *Client) EditBranch(ctx context.Context) (*gz.Branch, error) {
	return nil, poop.New("not implemented")
	// head, err := c.repo.Head()
	// if err != nil {
	// 	return nil, poop.Chain(err)
	// }

	// branch, err := c.store.GetBranch(ctx, head.Name().Short())
	// if err != nil && !errors.Is(err, store.ErrNotFound) {
	// 	return nil, poop.Chain(err)
	// }

	// var contents []byte
	// if branch != nil {
	// 	contents, err = json.MarshalIndent(branch, "", "  ")
	// 	if err != nil {
	// 		return nil, poop.Chain(err)
	// 	}
	// }

	// contents, err = editor.EditFrom(ctx, c.repo, contents)
	// if err != nil {
	// 	return nil, poop.Chain(err)
	// }

	// contents = bytes.TrimSpace(contents)
	// if len(contents) == 0 {
	// 	// TODO(kellegous): Delete the branch from the store
	// 	// or do nothing?
	// 	return branch, nil
	// }

	// var updated internal.Branch
	// if err := json.Unmarshal(contents, &updated); err != nil {
	// 	return nil, poop.Chain(err)
	// }

	// branch, err = c.store.UpsertBranch(ctx, &updated, nil)
	// if err != nil {
	// 	return nil, poop.Chain(err)
	// }

	// return &updated, nil
}

func (c *Client) GetBranch(ctx context.Context) (*gz.Branch, error) {
	head, err := c.repo.Head()
	if err != nil {
		return nil, poop.Chain(err)
	}

	return c.store.GetBranch(ctx, head.Name().Short())
}
