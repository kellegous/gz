package client_test

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/kellegous/poop"
)

type Repo struct {
	path string
}

type GitOptions struct {
	env []string
}

type GitOption func(*GitOptions)

func WithEnv(env []string) GitOption {
	return func(o *GitOptions) {
		o.env = env
	}
}

func WithCommitTime(t time.Time) GitOption {
	return func(o *GitOptions) {
		o.env = append(o.env, fmt.Sprintf("GIT_AUTHOR_DATE=%s", t.Format(time.RFC3339)))
		o.env = append(o.env, fmt.Sprintf("GIT_COMMITTER_DATE=%s", t.Format(time.RFC3339)))
	}
}

func (r *Repo) git(opts ...GitOption) func(ctx context.Context, args ...string) *exec.Cmd {
	return func(ctx context.Context, args ...string) *exec.Cmd {
		o := GitOptions{
			env: os.Environ(),
		}
		for _, opt := range opts {
			opt(&o)
		}

		c := exec.CommandContext(ctx, "git", args...)
		c.Dir = r.path
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		c.Stdin = os.Stdin
		if o.env != nil {
			c.Env = o.env
		}
		return c
	}
}

func (r *Repo) Main() *Branch {
	return r.GetBranch("main")
}

func (r *Repo) GetBranch(name string) *Branch {
	return &Branch{
		repo: r,
		Name: name,
	}
}

type Branch struct {
	repo *Repo
	Name string
}

func (b *Branch) Create(ctx context.Context, from string) error {
	git := b.repo.git()
	if from == "" {
		return poop.Chain(git(
			ctx,
			"checkout",
			"-b",
			b.Name,
		).Run())
	}
	return poop.Chain(git(
		ctx,
		"checkout",
		"-b",
		b.Name,
		from,
	).Run())
}

func (b *Branch) Checkout(ctx context.Context) error {
	git := b.repo.git()
	return poop.Chain(git(
		ctx,
		"checkout",
		b.Name,
	).Run())
}

func (b *Branch) Commit(
	ctx context.Context,
	message string,
	t time.Time,
	files ...*File,
) error {
	git := b.repo.git()
	for _, file := range files {
		if err := os.WriteFile(
			filepath.Join(b.repo.path, file.Path),
			[]byte(file.Content),
			0644,
		); err != nil {
			return poop.Chain(err)
		}

		if err := git(
			ctx,
			"add",
			file.Path,
		).Run(); err != nil {
			return poop.Chain(err)
		}
	}

	if err := b.repo.git(WithCommitTime(t))(
		ctx,
		"commit",
		"-m",
		message,
	).Run(); err != nil {
		return poop.Chain(err)
	}

	return nil
}

type File struct {
	Path    string
	Content string
}
