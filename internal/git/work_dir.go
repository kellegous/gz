package git

import (
	"context"
	"os"
	"os/exec"

	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing"
	"github.com/kellegous/poop"
)

type WorkDir struct {
	path string
	repo *git.Repository
	wt   *git.Worktree
}

func WorkDirAt(path string) (*WorkDir, error) {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return nil, poop.Chain(err)
	}

	wt, err := repo.Worktree()
	if err != nil {
		return nil, poop.Chain(err)
	}

	return &WorkDir{
		path: path,
		repo: repo,
		wt:   wt,
	}, nil
}

func (w *WorkDir) Worktree() *git.Worktree {
	return w.wt
}

func (w *WorkDir) Repository() *git.Repository {
	return w.repo
}

func (w *WorkDir) CreateBranch(
	ctx context.Context,
	name string,
	from string,
	opts ...GitOption,
) error {
	args := []string{"checkout", "-b", name}
	if from != "" {
		args = append(args, from)
	}
	return poop.Chain(w.gitCommand(ctx, args, opts...).Run())
}

type Msg struct {
	m string
}

func NoEdit() *Msg {
	return noEdit
}

func Message(m string) *Msg {
	return &Msg{m: m}
}

var noEdit = &Msg{}

type CommitOptions struct {
	All     bool
	Amend   bool
	Message *Msg
}

func (w *WorkDir) Commit(
	ctx context.Context,
	opts CommitOptions,
	gitOpts ...GitOption,
) (*plumbing.Reference, error) {
	args := []string{"commit"}

	if opts.All {
		args = append(args, "-a")
	}

	if opts.Amend {
		args = append(args, "--amend")
	}

	if m := opts.Message; m == noEdit {
		args = append(args, "--no-edit")
	} else if m != nil {
		args = append(args, "-m", m.m)
	}

	if err := w.gitCommand(ctx, args, gitOpts...).Run(); err != nil {
		return nil, poop.Chain(err)
	}

	head, err := w.repo.Head()
	if err != nil {
		return nil, poop.Chain(err)
	}

	return head, nil
}

func (w *WorkDir) gitCommand(
	ctx context.Context,
	args []string,
	opts ...GitOption,
) *exec.Cmd {
	var o GitOptions
	for _, opt := range opts {
		opt(&o)
	}

	cmd := exec.CommandContext(ctx, "git", args...)
	cmd.Dir = w.path
	cmd.Env = append(os.Environ(), o.env...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd
}
