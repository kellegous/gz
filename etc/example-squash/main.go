package main

import (
	"context"
	"flag"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/kellegous/poop"
)

func main() {
	if err := run(context.Background()); err != nil {
		poop.HitFan(err)
	}
}

func run(ctx context.Context) error {
	var path string
	flag.StringVar(&path, "path", "repo", "the path to the repository")
	flag.Parse()

	if err := os.MkdirAll(path, 0755); err != nil {
		return poop.Chain(err)
	}

	if err := runCommand(
		exec.CommandContext(ctx, "git", "init"),
		path,
	); err != nil {
		return poop.Chain(err)
	}

	if err := commit(
		ctx,
		path,
		"initial commit",
		&File{
			Path:    "README.md",
			Content: "This is a README file",
		},
	); err != nil {
		return poop.Chain(err)
	}

	{ // build feature branch

		if err := runCommand(
			exec.CommandContext(ctx, "git", "checkout", "-b", "feature-branch"),
			path,
		); err != nil {
			return poop.Chain(err)
		}

		if err := commit(
			ctx,
			path,
			"first commit",
			&File{
				Path:    "feature.txt",
				Content: "feature v1",
			},
		); err != nil {
			return poop.Chain(err)
		}

		if err := commit(
			ctx,
			path,
			"second commit",
			&File{
				Path:    "feature.txt",
				Content: "feature v2",
			},
		); err != nil {
			return poop.Chain(err)
		}
	}

	{ // add more commits to the main branch
		if err := runCommand(
			exec.CommandContext(ctx, "git", "checkout", "main"),
			path,
		); err != nil {
			return poop.Chain(err)
		}

		if err := commit(
			ctx,
			path,
			"third commit",
			&File{
				Path:    "main.txt",
				Content: "main v1",
			},
		); err != nil {
			return poop.Chain(err)
		}

		if err := runCommand(
			exec.CommandContext(ctx, "git", "merge", "--squash", "feature-branch"),
			path,
		); err != nil {
			return poop.Chain(err)
		}

		if err := runCommand(
			exec.CommandContext(ctx, "git", "commit", "--no-edit"),
			path,
		); err != nil {
			return poop.Chain(err)
		}
	}

	// 1. create a feature branch
	// 2. commit two changes to the feature branch
	// 3. checkout the main branch
	// 4. commit another change to the main branch
	// 5. git merge --squash feature-branch
	// 6. we must now get feature-branch to merge cleanly onto main

	return nil
}

type File struct {
	Path    string
	Content string
}

func commit(
	ctx context.Context,
	path string,
	message string,
	files ...*File,
) error {

	for _, file := range files {
		if filepath.IsAbs(file.Path) {
			return poop.Newf("file path is absolute: %s", file.Path)
		}

		if err := os.WriteFile(filepath.Join(path, file.Path), []byte(file.Content), 0644); err != nil {
			return poop.Chain(err)
		}

		if err := runCommand(
			exec.CommandContext(ctx, "git", "add", file.Path),
			path,
		); err != nil {
			return poop.Chain(err)
		}
	}

	return poop.Chain(runCommand(
		exec.CommandContext(ctx, "git", "commit", "-m", message),
		path,
	))
}

func runCommand(
	cmd *exec.Cmd,
	path string,
) error {
	cmd.Dir = path
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}
