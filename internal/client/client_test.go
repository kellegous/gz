package client_test

import (
	"bufio"
	"context"
	"fmt"
	"iter"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/kellegous/poop"
)

type Repo struct {
	path string
}

func CreateRepo(ctx context.Context) (*Repo, error) {
	dir, err := os.MkdirTemp("", "gz-repo-test-*")
	if err != nil {
		return nil, poop.Chain(err)
	}

	r := &Repo{
		path: dir,
	}

	if err := r.git()(ctx, "init").Run(); err != nil {
		return nil, poop.Chain(err)
	}

	return r, nil
}

func CreateRepoWith(
	ctx context.Context,
	fn func(*Repo) error,
) (*Repo, error) {
	r, err := CreateRepo(ctx)
	if err != nil {
		return nil, poop.Chain(err)
	}
	if err := fn(r); err != nil {
		return nil, poop.Chain(err)
	}
	return r, nil
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

func (r *Repo) Remove() error {
	return os.RemoveAll(r.path)
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

func (r *Repo) CreateBranch(
	ctx context.Context,
	name string,
	from string,
) error {
	git := r.git()
	if from == "" {
		return poop.Chain(git(ctx, "checkout", "-b", name).Run())
	}
	return poop.Chain(git(ctx, "checkout", "-b", name, from).Run())
}

func (r *Repo) CheckoutBranch(ctx context.Context, name string) error {
	return poop.Chain(r.git()(ctx, "checkout", name).Run())
}

func (r *Repo) Commit(
	ctx context.Context,
	message string,
	t time.Time,
	amend bool,
	files ...*File,
) error {
	git := r.git()

	for _, file := range files {
		if err := os.WriteFile(
			filepath.Join(r.path, file.Path),
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

	args := []string{"commit"}
	if amend {
		args = append(args, "--amend")
	}
	args = append(args, "-m", message)

	return poop.Chain(r.git(WithCommitTime(t))(ctx, args...).Run())
}

func (r *Repo) Commits(ctx context.Context, name string) iter.Seq2[*Commit, error] {
	git := r.git()
	return func(yield func(*Commit, error) bool) {
		var c *exec.Cmd
		if name != "" {
			c = git(ctx, "log", "--pretty=format:%H %at", name, "--")
		} else {
			c = git(ctx, "log", "--pretty=format:%H %at")
		}
		// TODO(kellegous): fix this.
		c.Stdout = nil

		r, err := c.StdoutPipe()
		if err != nil {
			yield(nil, poop.Chain(err))
			return
		}
		defer r.Close()

		if err := c.Start(); err != nil {
			yield(nil, poop.Chain(err))
			return
		}

		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			line := scanner.Text()
			hash, ts, ok := strings.Cut(line, " ")
			if !ok {
				yield(nil, poop.Newf("invalid commit line: %s", line))
				return
			}

			t, err := strconv.ParseInt(ts, 10, 64)
			if err != nil {
				yield(nil, poop.Chain(err))
				return
			}

			if !yield(&Commit{
				Hash: hash,
				Time: time.Unix(t, 0),
			}, nil) {
				return
			}
		}

		if err := scanner.Err(); err != nil {
			yield(nil, poop.Chain(err))
			return
		}
	}
}

type File struct {
	Path    string
	Content string
}

type Commit struct {
	Hash string
	Time time.Time
}

func TestRepo(t *testing.T) {
	ctx := t.Context()

	r, err := CreateRepo(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer r.Remove()

	if err := r.Commit(
		ctx,
		"initial commit",
		time.Now(),
		false,
		&File{
			Path:    "README.md",
			Content: "This is a README file",
		},
	); err != nil {
		t.Fatal(err)
	}

	for commit, err := range r.Commits(ctx, "") {
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("commit: %s", commit.Hash)
	}
}
