package client

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/go-git/go-git/v6"
	"github.com/kellegous/gz/internal/store"
	"github.com/kellegous/poop"
)

const defaultStorePath = ".git/gz.db"

type Client struct {
	path     string
	repo     *git.Repository
	workTree *git.Worktree
	store    *store.Store
	envFn    func() []string
}

func Open(
	ctx context.Context,
	root string,
	opts ...Option,
) (*Client, error) {
	o := Options{
		envFn: func() []string {
			return nil
		},
		storePath: defaultStorePath,
	}
	for _, opt := range opts {
		opt(&o)
	}

	repo, err := git.PlainOpen(root)
	if err != nil {
		return nil, poop.Chain(err)
	}

	workTree, err := repo.Worktree()
	if err != nil {
		return nil, poop.Chain(err)
	}

	s, err := store.Open(ctx, filepath.Join(root, o.storePath))
	if err != nil {
		return nil, poop.Chain(err)
	}

	return &Client{
		path:     root,
		repo:     repo,
		workTree: workTree,
		store:    s,
		envFn:    o.envFn,
	}, nil
}

func (c *Client) Close() error {
	return c.store.Close()
}

func (c *Client) gitCommand(ctx context.Context, args ...string) *exec.Cmd {
	cmd := exec.CommandContext(ctx, "git", args...)
	cmd.Dir = c.path
	cmd.Env = append(os.Environ(), c.envFn()...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd
}
